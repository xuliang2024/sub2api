package handler

import (
	"errors"
	"strconv"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// PartnerHandler 合伙人面板（用户端）：概览、下线定价、收益明细。
type PartnerHandler struct {
	partnerService   *service.PartnerService
	affiliateService *service.AffiliateService
}

func NewPartnerHandler(partnerService *service.PartnerService, affiliateService *service.AffiliateService) *PartnerHandler {
	return &PartnerHandler{partnerService: partnerService, affiliateService: affiliateService}
}

func partnerSubjectID(c *gin.Context) (int64, bool) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok || subject.UserID <= 0 {
		response.Unauthorized(c, "User not authenticated")
		return 0, false
	}
	return subject.UserID, true
}

func partnerUserError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPartnerNotFound):
		response.Forbidden(c, "Partner feature is not enabled for this account")
	case errors.Is(err, service.ErrPartnerDisabled):
		response.Forbidden(c, "Partner account is disabled")
	case errors.Is(err, service.ErrPartnerRateOutOfRange):
		response.BadRequest(c, "Rate multiplier out of allowed range")
	case errors.Is(err, service.ErrPartnerNotDownline):
		response.BadRequest(c, "User is not bound to you")
	case errors.Is(err, service.ErrPartnerGroupNotAllowed):
		response.BadRequest(c, "Group is not available to this user")
	default:
		response.ErrorFrom(c, err)
	}
}

// Overview GET /user/partner
func (h *PartnerHandler) Overview(c *gin.Context) {
	userID, ok := partnerSubjectID(c)
	if !ok {
		return
	}
	overview, err := h.partnerService.GetOverview(c.Request.Context(), userID)
	if err != nil {
		partnerUserError(c, err)
		return
	}
	// 邀请链接复用 affiliate 邀请码（best-effort，未开启返利时为空）。
	if h.affiliateService != nil {
		if detail, err := h.affiliateService.GetAffiliateDetail(c.Request.Context(), userID); err == nil && detail != nil {
			overview.AffCode = detail.AffCode
		}
	}
	response.Success(c, overview)
}

// ListUsers GET /user/partner/users
func (h *PartnerHandler) ListUsers(c *gin.Context) {
	userID, ok := partnerSubjectID(c)
	if !ok {
		return
	}
	items, err := h.partnerService.ListDownlines(c.Request.Context(), userID)
	if err != nil {
		partnerUserError(c, err)
		return
	}
	response.Success(c, gin.H{"items": items})
}

type setPartnerRatesRequest struct {
	Rates []service.PartnerRateInput `json:"rates" binding:"required"`
}

// SetRates PUT /user/partner/users/:id/rates
func (h *PartnerHandler) SetRates(c *gin.Context) {
	partnerUserID, ok := partnerSubjectID(c)
	if !ok {
		return
	}
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || targetUserID <= 0 {
		response.BadRequest(c, "Invalid user id")
		return
	}
	var req setPartnerRatesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if len(req.Rates) == 0 {
		response.BadRequest(c, "rates is empty")
		return
	}
	if err := h.partnerService.SetDownlineRates(c.Request.Context(), partnerUserID, targetUserID, req.Rates); err != nil {
		partnerUserError(c, err)
		return
	}
	response.Success(c, gin.H{"user_id": targetUserID})
}

// ListCommission GET /user/partner/commission
func (h *PartnerHandler) ListCommission(c *gin.Context) {
	partnerUserID, ok := partnerSubjectID(c)
	if !ok {
		return
	}
	page, pageSize := response.ParsePagination(c)
	filter := service.PartnerCommissionFilter{Page: page, PageSize: pageSize}
	if raw := c.Query("user_id"); raw != "" {
		if id, err := strconv.ParseInt(raw, 10, 64); err == nil && id > 0 {
			filter.UserID = id
		}
	}
	if t, ok := parsePartnerDate(c.Query("start_date"), false); ok {
		filter.StartTime = &t
	}
	if t, ok := parsePartnerDate(c.Query("end_date"), true); ok {
		filter.EndTime = &t
	}
	items, total, totalCommission, err := h.partnerService.ListCommissionRecords(c.Request.Context(), partnerUserID, filter)
	if err != nil {
		partnerUserError(c, err)
		return
	}
	response.Success(c, gin.H{
		"items":            items,
		"total":            total,
		"total_commission": totalCommission,
		"page":             page,
		"page_size":        pageSize,
	})
}

// parsePartnerDate 解析 YYYY-MM-DD；end 时取次日 0 点（右开区间）。
func parsePartnerDate(raw string, end bool) (time.Time, bool) {
	if raw == "" {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return time.Time{}, false
	}
	if end {
		t = t.AddDate(0, 0, 1)
	}
	return t, true
}
