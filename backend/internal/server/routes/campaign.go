package routes

import (
	"github.com/Wei-Shaw/sub2api/internal/handler"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// RegisterCampaignRoutes 注册广告落地页的公开读取与转化路由。
//
// 公开读取挂 OptionalJWT：匿名也要能看落地页（投放的第一屏），带 token 时后续
// 可以按用户做差异化（当前未用）。转化动作 POST /campaigns/:slug/key 必须登录。
func RegisterCampaignRoutes(
	v1 *gin.RouterGroup,
	h *handler.Handlers,
	jwtAuth middleware.JWTAuthMiddleware,
	optionalJWT middleware.OptionalJWTAuthMiddleware,
	settingService *service.SettingService,
	panelRateLimiter *middleware.PanelRateLimiter,
) {
	public := v1.Group("/public/campaign")
	public.Use(panelRateLimiter.PublicIP())
	public.Use(gin.HandlerFunc(optionalJWT))
	public.Use(middleware.BackendModeUserGuard(settingService))
	{
		// 按 Host（或 ?c=slug）解析落地页配置
		public.GET("", h.Campaign.Get)
	}

	claim := v1.Group("/campaigns")
	claim.Use(gin.HandlerFunc(jwtAuth))
	{
		// 一键建 Key（幂等：同一活动同一用户返回同一把）
		claim.POST("/:slug/key", h.Campaign.ClaimKey)
	}
}

// registerAdminCampaignRoutes 注册落地页管理路由。
//
// 复用 RegisterAdminRoutes 里已挂好的 adminAuth / 限流 / 审计 / 合规校验，
// 因此这里只做路径分组，不重复挂中间件。
func registerAdminCampaignRoutes(admin *gin.RouterGroup, h *handler.Handlers) {
	campaigns := admin.Group("/campaigns")
	{
		campaigns.GET("", h.Campaign.AdminList)
		campaigns.POST("", h.Campaign.AdminCreate)
		campaigns.GET("/:id", h.Campaign.AdminGet)
		campaigns.PUT("/:id", h.Campaign.AdminUpdate)
		campaigns.DELETE("/:id", h.Campaign.AdminDelete)
		// 以落地页视角渲染（含 draft），用于后台预览
		campaigns.GET("/:id/preview", h.Campaign.AdminPreview)
	}
}
