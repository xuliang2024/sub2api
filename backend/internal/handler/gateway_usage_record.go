package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const (
	usageRecordSchemaVersion = 1
	// usageRecordMaxWait 限制 ?wait_ms 的上限，避免长轮询把连接和 goroutine 占住。
	usageRecordMaxWait = 3000 * time.Millisecond
	// usageRecordPollInterval 是长轮询的重试间隔（计费异步落账，实测 ~200ms 可见）。
	usageRecordPollInterval = 150 * time.Millisecond
)

var usageRecordRequestIDPattern = regexp.MustCompile(`^[A-Za-z0-9:_.\-]{8,200}$`)

type usageRecordPromptDetails struct {
	CachedTokens        int `json:"cached_tokens"`
	CacheCreationTokens int `json:"cache_creation_tokens"`
}

type usageRecordTokenBlock struct {
	PromptTokens        int                      `json:"prompt_tokens"`
	CompletionTokens    int                      `json:"completion_tokens"`
	TotalTokens         int                      `json:"total_tokens"`
	PromptTokensDetails usageRecordPromptDetails `json:"prompt_tokens_details"`
}

type usageRecordCostBlock struct {
	Currency          string  `json:"currency"`
	InputCost         float64 `json:"input_cost"`
	OutputCost        float64 `json:"output_cost"`
	CacheReadCost     float64 `json:"cache_read_cost"`
	CacheCreationCost float64 `json:"cache_creation_cost"`
	TotalCost         float64 `json:"total_cost"`
	ActualCost        float64 `json:"actual_cost"`
	RateMultiplier    float64 `json:"rate_multiplier"`
}

type usageRecordTimingBlock struct {
	DurationMs   *int `json:"duration_ms"`
	FirstTokenMs *int `json:"first_token_ms"`
	// SettleMs 是本次查询在服务端等待账单落账的毫秒数（客户端可用它衡量落账快慢）。
	SettleMs int `json:"settle_ms"`
}

// usageRecordResponse 是 /v1/sub2api/usage/:request_id 的响应体。
// 字段白名单：只暴露对账必需的用量与金额，不含 account_id / channel_id /
// upstream_model / account_cost / api_key / user 等上游与账号内部信息。
type usageRecordResponse struct {
	Object           string                 `json:"object"`
	SchemaVersion    int                    `json:"schema_version"`
	RequestID        string                 `json:"request_id"`
	Model            string                 `json:"model"`
	InboundEndpoint  *string                `json:"inbound_endpoint"`
	UpstreamEndpoint *string                `json:"upstream_endpoint"`
	Stream           bool                   `json:"stream"`
	BillingMode      *string                `json:"billing_mode"`
	Usage            usageRecordTokenBlock  `json:"usage"`
	Cost             usageRecordCostBlock   `json:"cost"`
	Timing           usageRecordTimingBlock `json:"timing"`
	CreatedAt        time.Time              `json:"created_at"`
}

// UsageRecord 返回单次请求的用量与扣费。
// GET /v1/sub2api/usage/:request_id
//
// 鉴权：用户 API Key（Authorization: Bearer / x-api-key），服务端从认证上下文取
// api_key_id 作为查询条件，不接受调用方传入。未命中与不属于本 Key 一律返回 404，
// 避免通过响应差异探测 request_id 是否存在。
func (h *GatewayHandler) UsageRecord(c *gin.Context) {
	apiKey, ok := middleware2.GetAPIKeyFromContext(c)
	if !ok || apiKey == nil {
		h.errorResponse(c, http.StatusUnauthorized, "authentication_error", "Invalid API key")
		return
	}

	requestID, ok := normalizeUsageRecordRequestID(c.Param("request_id"))
	if !ok {
		h.errorResponse(c, http.StatusBadRequest, "invalid_request_error", "request_id is invalid")
		return
	}

	if h.usageService == nil {
		h.errorResponse(c, http.StatusServiceUnavailable, "api_error", "Usage lookup is unavailable")
		return
	}

	started := time.Now()
	deadline := started.Add(parseUsageRecordWaitMs(c.Query("wait_ms")))

	var record *service.UsageLog
	for {
		log, err := h.usageService.GetByRequestIDAndAPIKey(c.Request.Context(), requestID, apiKey.ID)
		if err == nil {
			record = log
			break
		}
		if !errors.Is(err, service.ErrUsageLogNotFound) {
			h.errorResponse(c, http.StatusInternalServerError, "api_error", "Failed to load usage record")
			return
		}
		if !time.Now().Before(deadline) {
			// 计费是异步落账：让客户端稍后重试，而不是把"还没写进来"当成错误暴露细节。
			c.Header("Retry-After", "1")
			h.errorResponse(c, http.StatusNotFound, "usage_not_found", "Usage record is not available yet")
			return
		}
		select {
		case <-c.Request.Context().Done():
			return
		case <-time.After(usageRecordPollInterval):
		}
	}

	c.Header("Cache-Control", "no-store")
	c.JSON(http.StatusOK, buildUsageRecordResponse(record, time.Since(started)))
}

// normalizeUsageRecordRequestID 归一化请求 ID。
// 网关在响应头 X-Client-Request-ID 里给出裸 UUID，而库里存的是 "client:<uuid>"，
// 所以裸 UUID 需要补前缀；其他命名空间（local: / generated: / web_search: 等）原样匹配。
func normalizeUsageRecordRequestID(raw string) (string, bool) {
	raw = strings.TrimSpace(raw)
	if !usageRecordRequestIDPattern.MatchString(raw) {
		return "", false
	}
	for _, prefix := range []string{"client:", "local:", "generated:", "web_search:", "grok-video:", "grok_audio:"} {
		if strings.HasPrefix(raw, prefix) {
			return raw, true
		}
	}
	return "client:" + raw, true
}

func parseUsageRecordWaitMs(raw string) time.Duration {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}
	ms, err := strconv.Atoi(raw)
	if err != nil || ms <= 0 {
		return 0
	}
	wait := time.Duration(ms) * time.Millisecond
	if wait > usageRecordMaxWait {
		wait = usageRecordMaxWait
	}
	return wait
}

func buildUsageRecordResponse(log *service.UsageLog, settle time.Duration) usageRecordResponse {
	if log == nil {
		return usageRecordResponse{Object: "sub2api.usage_record", SchemaVersion: usageRecordSchemaVersion}
	}

	// 库里的 InputTokens 已剔除缓存命中/写入部分；客户端看到的是含缓存的总量，
	// 这里按 OpenAI usage 语义还原，保证与当时响应里的 prompt_tokens 一致。
	promptTokens := log.InputTokens + log.CacheReadTokens + log.CacheCreationTokens

	return usageRecordResponse{
		Object:           "sub2api.usage_record",
		SchemaVersion:    usageRecordSchemaVersion,
		RequestID:        log.RequestID,
		Model:            log.Model,
		InboundEndpoint:  log.InboundEndpoint,
		UpstreamEndpoint: log.UpstreamEndpoint,
		Stream:           log.Stream,
		BillingMode:      log.BillingMode,
		Usage: usageRecordTokenBlock{
			PromptTokens:     promptTokens,
			CompletionTokens: log.OutputTokens,
			TotalTokens:      promptTokens + log.OutputTokens,
			PromptTokensDetails: usageRecordPromptDetails{
				CachedTokens:        log.CacheReadTokens,
				CacheCreationTokens: log.CacheCreationTokens,
			},
		},
		Cost: usageRecordCostBlock{
			Currency:          "USD",
			InputCost:         log.InputCost,
			OutputCost:        log.OutputCost,
			CacheReadCost:     log.CacheReadCost,
			CacheCreationCost: log.CacheCreationCost,
			TotalCost:         log.TotalCost,
			ActualCost:        log.ActualCost,
			RateMultiplier:    log.RateMultiplier,
		},
		Timing: usageRecordTimingBlock{
			DurationMs:   log.DurationMs,
			FirstTokenMs: log.FirstTokenMs,
			SettleMs:     int(settle / time.Millisecond),
		},
		CreatedAt: log.CreatedAt,
	}
}
