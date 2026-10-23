package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/handler"
	servermiddleware "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type usageRecordRouteRepo struct {
	service.UsageLogRepository
	apiKeyID int64
	record   *service.UsageLog
}

func (r *usageRecordRouteRepo) GetByRequestIDAndAPIKey(_ context.Context, requestID string, apiKeyID int64) (*service.UsageLog, error) {
	if r.record == nil || apiKeyID != r.apiKeyID || requestID != r.record.RequestID {
		return nil, service.ErrUsageLogNotFound
	}
	return r.record, nil
}

func newUsageRecordRouteTestRouter(usageRepo *usageRecordRouteRepo, injectKeyID int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	usageService := service.NewUsageService(usageRepo, nil, nil, nil)
	gatewayHandler := handler.NewGatewayHandler(
		nil, nil, nil, nil, nil, nil, nil, usageService, nil, nil, nil, nil, nil,
		&config.Config{}, nil,
	)

	router := gin.New()
	RegisterGatewayRoutes(
		router,
		&handler.Handlers{
			Gateway:       gatewayHandler,
			OpenAIGateway: &handler.OpenAIGatewayHandler{},
			AsyncImage:    handler.NewAsyncImageHandler(nil, nil),
		},
		servermiddleware.APIKeyAuthMiddleware(func(c *gin.Context) {
			if injectKeyID > 0 {
				c.Set(string(servermiddleware.ContextKeyAPIKey), &service.APIKey{ID: injectKeyID})
			}
			c.Next()
		}),
		nil, nil, nil, nil, nil,
		&config.Config{
			Gateway: config.GatewayConfig{
				MaxBodySize:     1024 * 1024,
				TextMaxBodySize: 1024 * 1024,
			},
		},
	)
	return router
}

func usageRecordRouteFixture() *service.UsageLog {
	inbound, upstream := "/v1/chat/completions", "/v1/responses"
	return &service.UsageLog{
		RequestID:        "client:1702da71-8f00-4883-aceb-441dbd75a2bc",
		Model:            "deepseek-v4.1-flash",
		InboundEndpoint:  &inbound,
		UpstreamEndpoint: &upstream,
		InputTokens:      41,
		OutputTokens:     79,
		TotalCost:        0.0000526549,
		ActualCost:       0.0000526549,
		RateMultiplier:   1,
		Stream:           true,
		CreatedAt:        time.Date(2026, 10, 24, 8, 30, 9, 0, time.FixedZone("JST", 9*3600)),
	}
}

func TestGatewayRoutesUsageRecordPathIsRegistered(t *testing.T) {
	router := newUsageRecordRouteTestRouter(&usageRecordRouteRepo{}, 0)

	for _, route := range router.Routes() {
		if route.Method == http.MethodGet && route.Path == "/v1/sub2api/usage/:request_id" {
			return
		}
	}
	t.Fatal("GET /v1/sub2api/usage/:request_id should be registered")
}

func TestGatewayRoutesUsageRecordEndToEnd(t *testing.T) {
	const rid = "client:1702da71-8f00-4883-aceb-441dbd75a2bc"

	t.Run("命中自己的记录", func(t *testing.T) {
		router := newUsageRecordRouteTestRouter(&usageRecordRouteRepo{apiKeyID: 8973, record: usageRecordRouteFixture()}, 8973)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+rid, nil))

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"object":"sub2api.usage_record"`)
		require.Contains(t, w.Body.String(), `"actual_cost":0.0000526549`)
	})

	t.Run("未鉴权", func(t *testing.T) {
		router := newUsageRecordRouteTestRouter(&usageRecordRouteRepo{apiKeyID: 8973, record: usageRecordRouteFixture()}, 0)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+rid, nil))

		require.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("查到别人的记录返回 404", func(t *testing.T) {
		router := newUsageRecordRouteTestRouter(&usageRecordRouteRepo{apiKeyID: 8973, record: usageRecordRouteFixture()}, 9999)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/sub2api/usage/"+rid, nil))

		require.Equal(t, http.StatusNotFound, w.Code)
		require.Equal(t, "1", w.Header().Get("Retry-After"))
		require.Contains(t, w.Body.String(), "usage_not_found")
	})

	t.Run("与其他子路径不冲突", func(t *testing.T) {
		router := newUsageRecordRouteTestRouter(&usageRecordRouteRepo{}, 8973)
		w := httptest.NewRecorder()
		// /v1/sub2api/billing 是兄弟路由，必须仍然可达（不应被 :request_id 吃掉）
		router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/v1/sub2api/billing", nil))
		require.NotEqual(t, http.StatusNotFound, w.Code)
	})
}
