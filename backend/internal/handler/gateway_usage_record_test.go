package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// fakeUsageLogRepo 只实现被测方法；其余方法借助嵌入的 nil 接口满足编译。
type fakeUsageLogRepo struct {
	service.UsageLogRepository
	apiKeyID int64
	record   *service.UsageLog
	// notFoundTimes > 0 时，前 N 次调用返回 NotFound，用于模拟异步落账。
	notFoundTimes int
	calls         int
	err           error
}

func (f *fakeUsageLogRepo) GetByRequestIDAndAPIKey(_ context.Context, requestID string, apiKeyID int64) (*service.UsageLog, error) {
	f.calls++
	if f.err != nil {
		return nil, f.err
	}
	if f.apiKeyID != 0 && apiKeyID != f.apiKeyID {
		return nil, service.ErrUsageLogNotFound
	}
	if f.notFoundTimes > 0 && f.calls <= f.notFoundTimes {
		return nil, service.ErrUsageLogNotFound
	}
	if f.record == nil || (f.record.RequestID != requestID && f.record.RequestID != "client:"+requestID) {
		return nil, service.ErrUsageLogNotFound
	}
	return f.record, nil
}

func newUsageRecordTestRouter(t *testing.T, h *GatewayHandler, apiKey *service.APIKey) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/v1/sub2api/usage/:request_id", func(c *gin.Context) {
		if apiKey != nil {
			c.Set(string(middleware2.ContextKeyAPIKey), apiKey)
		}
		h.UsageRecord(c)
	})
	return r
}

func newUsageRecordHandler(repo *fakeUsageLogRepo) *GatewayHandler {
	return &GatewayHandler{usageService: service.NewUsageService(repo, nil, nil, nil)}
}

func sampleUsageRecord() *service.UsageLog {
	reqID := "client:1702da71-8f00-4883-aceb-441dbd75a2bc"
	inbound, upstream, mode := "/v1/chat/completions", "/v1/responses", "token"
	dur, ttft := 1617, 848
	return &service.UsageLog{
		ID:               2499839,
		UserID:           1,
		APIKeyID:         8973,
		AccountID:        1409,
		RequestID:        reqID,
		Model:            "deepseek-v4.1-flash",
		InboundEndpoint:  &inbound,
		UpstreamEndpoint: &upstream,
		BillingMode:      &mode,
		InputTokens:      41,
		OutputTokens:     79,
		CacheReadTokens:  0,
		InputCost:        0.0000060472,
		OutputCost:       0.0000466077,
		TotalCost:        0.0000526549,
		ActualCost:       0.0000526549,
		RateMultiplier:   1,
		Stream:           true,
		DurationMs:       &dur,
		FirstTokenMs:     &ttft,
		CreatedAt:        time.Date(2026, 10, 24, 8, 30, 9, 0, time.FixedZone("JST", 9*3600)),
	}
}

func TestNormalizeUsageRecordRequestID(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		want   string
		wantOK bool
	}{
		{"裸 UUID 补 client 前缀", "1702da71-8f00-4883-aceb-441dbd75a2bc", "client:1702da71-8f00-4883-aceb-441dbd75a2bc", true},
		{"已带 client 前缀", "client:1702da71", "client:1702da71", true},
		{"web_search 原样", "web_search:abc12345", "web_search:abc12345", true},
		{"local 原样", "local:abc12345", "local:abc12345", true},
		{"空串", "", "", false},
		{"太短", "abc", "", false},
		{"含斜杠", "client:../../etc/passwd", "", false},
		{"含空格", "client:abc 123", "", false},
		{"超长", strings.Repeat("a", 201), "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := normalizeUsageRecordRequestID(tc.input)
			require.Equal(t, tc.wantOK, ok)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseUsageRecordWaitMs(t *testing.T) {
	require.Equal(t, time.Duration(0), parseUsageRecordWaitMs(""))
	require.Equal(t, time.Duration(0), parseUsageRecordWaitMs("abc"))
	require.Equal(t, time.Duration(0), parseUsageRecordWaitMs("-5"))
	require.Equal(t, 500*time.Millisecond, parseUsageRecordWaitMs("500"))
	// 上限截断
	require.Equal(t, usageRecordMaxWait, parseUsageRecordWaitMs("999999"))
}

func TestBuildUsageRecordResponse_ReconstructsPromptTokens(t *testing.T) {
	rec := sampleUsageRecord()
	rec.InputTokens = 180
	rec.CacheReadTokens = 256
	rec.CacheCreationTokens = 0
	rec.OutputTokens = 128

	out := buildUsageRecordResponse(rec, 220*time.Millisecond)

	// 库里的 input 不含缓存，接口要还原成客户端当时看到的 prompt_tokens
	require.Equal(t, 436, out.Usage.PromptTokens)
	require.Equal(t, 128, out.Usage.CompletionTokens)
	require.Equal(t, 564, out.Usage.TotalTokens)
	require.Equal(t, 256, out.Usage.PromptTokensDetails.CachedTokens)
	require.Equal(t, 220, out.Timing.SettleMs)
	require.Equal(t, "sub2api.usage_record", out.Object)
}

func TestUsageRecord_UnauthorizedWithoutAPIKey(t *testing.T) {
	r := newUsageRecordTestRouter(t, newUsageRecordHandler(&fakeUsageLogRepo{}), nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/client:abcdefgh", nil))
	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "authentication_error")
}

func TestUsageRecord_InvalidRequestID(t *testing.T) {
	r := newUsageRecordTestRouter(t, newUsageRecordHandler(&fakeUsageLogRepo{}), &service.APIKey{ID: 1})
	for _, bad := range []string{"abc", "client:has space", strings.Repeat("a", 201)} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+url.PathEscape(bad), nil))
		require.Equal(t, http.StatusBadRequest, rec.Code, "id=%q", bad)
		require.Contains(t, rec.Body.String(), "invalid_request_error")
	}
}

func TestUsageRecord_ServiceUnavailable(t *testing.T) {
	r := newUsageRecordTestRouter(t, &GatewayHandler{}, &service.APIKey{ID: 1})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/client:abcdefgh", nil))
	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestUsageRecord_NotFoundSetsRetryAfter(t *testing.T) {
	r := newUsageRecordTestRouter(t, newUsageRecordHandler(&fakeUsageLogRepo{}), &service.APIKey{ID: 1})
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/client:abcdefgh", nil))
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "1", rec.Header().Get("Retry-After"))
	require.Contains(t, rec.Body.String(), "usage_not_found")
}

// 关键安全用例：别人的 request_id 与"根本不存在"必须完全不可区分。
func TestUsageRecord_OtherKeyIsIndistinguishable(t *testing.T) {
	repo := &fakeUsageLogRepo{apiKeyID: 8973, record: sampleUsageRecord()}

	otherKey := newUsageRecordTestRouter(t, newUsageRecordHandler(repo), &service.APIKey{ID: 9999})
	recOther := httptest.NewRecorder()
	otherKey.ServeHTTP(recOther, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+sampleUsageRecord().RequestID, nil))

	missing := newUsageRecordTestRouter(t, newUsageRecordHandler(&fakeUsageLogRepo{apiKeyID: 8973, record: sampleUsageRecord()}), &service.APIKey{ID: 8973})
	recMissing := httptest.NewRecorder()
	missing.ServeHTTP(recMissing, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/client:ffffffff-0000-0000-0000-000000000000", nil))

	require.Equal(t, http.StatusNotFound, recOther.Code)
	require.Equal(t, recMissing.Code, recOther.Code)
	require.Equal(t, recMissing.Body.String(), recOther.Body.String(), "越权与不存在的响应体必须一致")
}

func TestUsageRecord_Success(t *testing.T) {
	repo := &fakeUsageLogRepo{apiKeyID: 8973, record: sampleUsageRecord()}
	r := newUsageRecordTestRouter(t, newUsageRecordHandler(repo), &service.APIKey{ID: 8973})

	rec := httptest.NewRecorder()
	// 裸 UUID 也要能命中
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/1702da71-8f00-4883-aceb-441dbd75a2bc", nil))

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "no-store", rec.Header().Get("Cache-Control"))

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "sub2api.usage_record", body["object"])
	require.Equal(t, "client:1702da71-8f00-4883-aceb-441dbd75a2bc", body["request_id"])
	require.Equal(t, "/v1/chat/completions", body["inbound_endpoint"])
	require.Equal(t, "/v1/responses", body["upstream_endpoint"])

	usage := body["usage"].(map[string]any)
	require.EqualValues(t, 41, usage["prompt_tokens"])
	require.EqualValues(t, 79, usage["completion_tokens"])
	require.EqualValues(t, 120, usage["total_tokens"])

	cost := body["cost"].(map[string]any)
	require.Equal(t, "USD", cost["currency"])
	require.InDelta(t, 0.0000526549, cost["actual_cost"], 1e-12)
	require.InDelta(t, 1, cost["rate_multiplier"], 1e-12)

	// 字段白名单：不得出现上游/账号/Key 内部信息
	raw := rec.Body.String()
	for _, forbidden := range []string{"account_id", "channel_id", "upstream_model", "account_cost", "api_key", "\"user\"", "api_key_id", "user_id"} {
		require.NotContains(t, raw, forbidden)
	}
}

func TestUsageRecord_WaitMsPollsUntilSettled(t *testing.T) {
	repo := &fakeUsageLogRepo{apiKeyID: 8973, record: sampleUsageRecord(), notFoundTimes: 2}
	r := newUsageRecordTestRouter(t, newUsageRecordHandler(repo), &service.APIKey{ID: 8973})

	started := time.Now()
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+sampleUsageRecord().RequestID+"?wait_ms=2000", nil))
	elapsed := time.Since(started)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, 3, repo.calls, "应当重试到第 3 次才拿到记录")
	require.GreaterOrEqual(t, elapsed, 2*usageRecordPollInterval)
	require.Less(t, elapsed, 2*time.Second)
}

func TestUsageRecord_RepoErrorIsInternalError(t *testing.T) {
	repo := &fakeUsageLogRepo{apiKeyID: 8973, err: errors.New("db down")}
	r := newUsageRecordTestRouter(t, newUsageRecordHandler(repo), &service.APIKey{ID: 8973})

	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+sampleUsageRecord().RequestID, nil))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.NotContains(t, rec.Body.String(), "db down", "内部错误细节不能回传客户端")
}
