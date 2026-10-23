package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// CampaignHandler 处理「模型广告投放落地页」的公开读取与转化动作。
//
// 与站点级 settings 的区别：settings 是单站点平铺 KV，没有 host 维度；落地页需要
// 一个域名一套内容，因此单独成表并按 Host 解析（见 service.CampaignService）。
type CampaignHandler struct {
	campaignService *service.CampaignService
}

// NewCampaignHandler 创建落地页 handler。
func NewCampaignHandler(campaignService *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{campaignService: campaignService}
}

// campaignRequest 管理端新建/更新请求体。
type campaignRequest struct {
	Slug         string                       `json:"slug"`
	Name         string                       `json:"name"`
	Status       string                       `json:"status"`
	Host         string                       `json:"host"`
	Locale       string                       `json:"locale"`
	GroupID      *int64                       `json:"group_id"`
	Brand        service.CampaignBrand        `json:"brand"`
	Hero         service.CampaignHero         `json:"hero"`
	ContentMD    string                       `json:"content_md"`
	ModelEntries []service.CampaignModelEntry `json:"model_entries"`
	CodeSamples  []service.CampaignCodeSample `json:"code_samples"`
	CTA          service.CampaignCTA          `json:"cta"`
	Trial        service.CampaignTrial        `json:"trial"`
	SEO          service.CampaignSEO          `json:"seo"`
	// Translations 按语言的正文覆盖，key 为语言代码（en/zh/ja/ko/th）。
	Translations map[string]service.CampaignLocaleContent `json:"translations"`
	// HiddenModels 不在落地页展示的模型（别名/退役模型）。
	HiddenModels []string `json:"hidden_models"`
}

func (r *campaignRequest) toCampaign(id int64) *service.Campaign {
	return &service.Campaign{
		ID:           id,
		Slug:         r.Slug,
		Name:         r.Name,
		Status:       service.CampaignStatus(strings.TrimSpace(r.Status)),
		Host:         r.Host,
		Locale:       r.Locale,
		GroupID:      r.GroupID,
		Brand:        r.Brand,
		Hero:         r.Hero,
		ContentMD:    r.ContentMD,
		ModelEntries: r.ModelEntries,
		CodeSamples:  r.CodeSamples,
		CTA:          r.CTA,
		Trial:        r.Trial,
		SEO:          r.SEO,
		Translations: r.Translations,
		HiddenModels: r.HiddenModels,
	}
}

// Get 公开读取落地页。
//
// GET /api/v1/public/campaign
//
// 解析顺序：?c=<slug> 优先（同域名多活动 / 预览），否则按 Host 匹配 active 活动。
// 只返回 active；draft 预览走管理端 /admin/campaigns/:id/preview。
func (h *CampaignHandler) Get(c *gin.Context) {
	slug := strings.TrimSpace(c.Query("c"))
	if slug == "" {
		slug = strings.TrimSpace(c.Query("slug"))
	}

	view, err := h.campaignService.Resolve(c.Request.Context(), service.CampaignResolveOptions{
		Host:            campaignRequestHost(c),
		Slug:            slug,
		AllowInactive:   false,
		FallbackBaseURL: campaignRequestOrigin(c),
	})
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}

	// 落地页是投放入口，允许浏览器与 CDN 做短缓存；分组模型另有 60s 进程内缓存。
	c.Header("Cache-Control", "public, max-age=60")
	response.Success(c, view)
}

// ClaimKey 为当前登录用户在活动绑定的分组下发放一把 Key（幂等）。
//
// POST /api/v1/campaigns/:slug/key
func (h *CampaignHandler) ClaimKey(c *gin.Context) {
	subject, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok || subject.UserID <= 0 {
		response.Unauthorized(c, "Authentication required")
		return
	}
	slug := strings.TrimSpace(c.Param("slug"))
	if slug == "" {
		response.BadRequest(c, "campaign slug is required")
		return
	}

	view, err := h.campaignService.ClaimKey(c.Request.Context(), subject.UserID, slug, campaignRequestOrigin(c))
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	response.Success(c, view)
}

// AdminList 列出全部活动。
//
// GET /api/v1/admin/campaigns
func (h *CampaignHandler) AdminList(c *gin.Context) {
	items, err := h.campaignService.List(c.Request.Context())
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	if items == nil {
		items = []service.Campaign{}
	}
	c.Header("Cache-Control", "no-store")
	response.Success(c, gin.H{"items": items})
}

// AdminGet 按 ID 读取活动。
//
// GET /api/v1/admin/campaigns/:id
func (h *CampaignHandler) AdminGet(c *gin.Context) {
	id, ok := parseCampaignID(c)
	if !ok {
		return
	}
	campaign, err := h.campaignService.Get(c.Request.Context(), id)
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	if campaign == nil {
		response.NotFound(c, "Campaign not found")
		return
	}
	c.Header("Cache-Control", "no-store")
	response.Success(c, campaign)
}

// AdminPreview 以落地页视角渲染任意状态的活动（含 draft），供后台预览。
//
// GET /api/v1/admin/campaigns/:id/preview
func (h *CampaignHandler) AdminPreview(c *gin.Context) {
	id, ok := parseCampaignID(c)
	if !ok {
		return
	}
	campaign, err := h.campaignService.Get(c.Request.Context(), id)
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	if campaign == nil {
		response.NotFound(c, "Campaign not found")
		return
	}
	view, err := h.campaignService.Resolve(c.Request.Context(), service.CampaignResolveOptions{
		Slug:            campaign.Slug,
		AllowInactive:   true,
		FallbackBaseURL: campaignRequestOrigin(c),
	})
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	response.Success(c, view)
}

// AdminCreate 新建活动。
//
// POST /api/v1/admin/campaigns
func (h *CampaignHandler) AdminCreate(c *gin.Context) {
	var req campaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid campaign payload")
		return
	}
	created, err := h.campaignService.Create(c.Request.Context(), req.toCampaign(0), adminOperatorID(c))
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	response.Created(c, created)
}

// AdminUpdate 全量更新活动。
//
// PUT /api/v1/admin/campaigns/:id
func (h *CampaignHandler) AdminUpdate(c *gin.Context) {
	id, ok := parseCampaignID(c)
	if !ok {
		return
	}
	var req campaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "invalid campaign payload")
		return
	}
	updated, err := h.campaignService.Update(c.Request.Context(), req.toCampaign(id))
	if err != nil {
		h.respondCampaignError(c, err)
		return
	}
	response.Success(c, updated)
}

// AdminDelete 删除活动。
//
// DELETE /api/v1/admin/campaigns/:id
func (h *CampaignHandler) AdminDelete(c *gin.Context) {
	id, ok := parseCampaignID(c)
	if !ok {
		return
	}
	if err := h.campaignService.Delete(c.Request.Context(), id); err != nil {
		h.respondCampaignError(c, err)
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *CampaignHandler) respondCampaignError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrCampaignNotFound), errors.Is(err, service.ErrCampaignNotActive):
		response.NotFound(c, "Campaign not found")
	case errors.Is(err, service.ErrCampaignNoGroup):
		response.BadRequest(c, "该活动尚未绑定模型分组，无法发放 Key")
	case errors.Is(err, service.ErrCampaignHostTaken):
		response.Error(c, http.StatusConflict, "该域名已被其它已上线活动占用")
	case errors.Is(err, service.ErrCampaignSlugTaken):
		response.Error(c, http.StatusConflict, "该 slug 已被占用")
	case errors.Is(err, service.ErrCampaignSlugInvalid):
		response.BadRequest(c, "slug 只能包含小写字母、数字与中划线，且不能以中划线开头或结尾")
	case errors.Is(err, service.ErrCampaignNameEmpty):
		response.BadRequest(c, "活动名称不能为空")
	case errors.Is(err, service.ErrCampaignStatusBad):
		response.BadRequest(c, "活动状态只能是 draft / active / archived")
	default:
		response.ErrorFrom(c, err)
	}
}

// campaignRequestHost 取请求域名：优先 X-Forwarded-Host（多级代理取第一个）。
func campaignRequestHost(c *gin.Context) string {
	raw := strings.TrimSpace(c.GetHeader("X-Forwarded-Host"))
	if raw != "" {
		if idx := strings.Index(raw, ","); idx >= 0 {
			raw = raw[:idx]
		}
		return service.NormalizeCampaignHost(raw)
	}
	return service.NormalizeCampaignHost(c.Request.Host)
}

// campaignRequestOrigin 还原对外 origin，用于 api_base_url 未配置时的兜底。
func campaignRequestOrigin(c *gin.Context) string {
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	if proto := strings.TrimSpace(c.GetHeader("X-Forwarded-Proto")); proto != "" {
		if idx := strings.Index(proto, ","); idx >= 0 {
			proto = proto[:idx]
		}
		scheme = strings.ToLower(strings.TrimSpace(proto))
	}
	if scheme != "http" && scheme != "https" {
		scheme = "https"
	}
	host := campaignRequestHost(c)
	if host == "" {
		return ""
	}
	return scheme + "://" + host
}

func parseCampaignID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "invalid campaign id")
		return 0, false
	}
	return id, true
}

// adminOperatorID 取当前管理员 ID 用于 created_by；取不到返回 0。
func adminOperatorID(c *gin.Context) int64 {
	if subject, ok := middleware.GetAuthSubjectFromContext(c); ok {
		return subject.UserID
	}
	return 0
}
