package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/handler"

	"github.com/gin-gonic/gin"
)

// 落地页路由的注册不能 panic：gin 对「同一位置的通配符」很敏感，
// /campaigns/:id 与 /campaigns/:id/preview 这类嵌套如果前缀处理不当会在启动时直接崩。
// 这个测试只验证路由树能建起来，以及关键路径能匹配到对应 handler。
func TestRegisterCampaignRoutesBuildsTree(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// 拦截中间件必须在注册路由之前挂上：gin 的路由在注册时就固化了当时的 handler 链，
	// 之后再 Use 不会作用到已有路由。这里只关心路由树，不触发业务逻辑。
	matched := map[string]bool{}
	router.Use(func(c *gin.Context) {
		matched[c.FullPath()] = true
		c.AbortWithStatus(http.StatusNoContent)
	})

	// handler 用空壳即可。
	routes := router.Group("/api/v1")
	admin := routes.Group("/admin")
	registerAdminCampaignRoutes(admin, &handler.Handlers{})

	cases := []struct {
		method string
		path   string
		want   string
	}{
		{http.MethodGet, "/api/v1/admin/campaigns", "/api/v1/admin/campaigns"},
		{http.MethodPost, "/api/v1/admin/campaigns", "/api/v1/admin/campaigns"},
		{http.MethodGet, "/api/v1/admin/campaigns/7", "/api/v1/admin/campaigns/:id"},
		{http.MethodPut, "/api/v1/admin/campaigns/7", "/api/v1/admin/campaigns/:id"},
		{http.MethodDelete, "/api/v1/admin/campaigns/7", "/api/v1/admin/campaigns/:id"},
		{http.MethodGet, "/api/v1/admin/campaigns/7/preview", "/api/v1/admin/campaigns/:id/preview"},
	}

	for _, tc := range cases {
		matched = map[string]bool{}
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(tc.method, tc.path, nil))
		if recorder.Code == http.StatusNotFound {
			t.Errorf("%s %s 未匹配到任何路由", tc.method, tc.path)
			continue
		}
		if !matched[tc.want] {
			t.Errorf("%s %s 匹配到了 %v，期望 %s", tc.method, tc.path, matched, tc.want)
		}
	}
}
