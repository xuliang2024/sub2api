package admin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// PartnerHandler 合伙人制度管理端：合伙人设置、绑定关系、结算登记。
type PartnerHandler struct {
	partnerService *service.PartnerService
}

func NewPartnerHandler(partnerService *service.PartnerService) *PartnerHandler {
	return &PartnerHandler{partnerService: partnerService}
}

func partnerActorID(c *gin.Context) int64 {
	if subject, ok := middleware.GetAuthSubjectFromContext(c); ok {
		return subject.UserID
	}
	return 0
}

func partnerError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPartnerNotFound):
		response.NotFound(c, "Partner not found")
	case errors.Is(err, service.ErrPartnerNotDownline):
		response.BadRequest(c, "User is not bound to this partner")
	case errors.Is(err, service.ErrPartnerDisabled):
		response.BadRequest(c, "Partner is disabled")
	case errors.Is(err, service.ErrPartnerInvalidMinRate):
		response.BadRequest(c, "Invalid min rate or status")
	case errors.Is(err, service.ErrPartnerRateOutOfRange):
		response.BadRequest(c, "Rate multiplier out of partner allowed range")
	case errors.Is(err, service.ErrPartnerGroupNotAllowed):
		response.BadRequest(c, "Group is not available to the downline user")
	case errors.Is(err, service.ErrPartnerIsDownline):
		response.BadRequest(c, "User is bound as a downline, unbind first")
	case errors.Is(err, service.ErrPartnerDownlineIsPartner):
		response.BadRequest(c, "Partner cannot be bound as a downline")
	case errors.Is(err, service.ErrPartnerInvalidSettle):
		response.BadRequest(c, "Settlement amount is invalid or exceeds pending commission")
	default:
		response.ErrorFrom(c, err)
	}
}

func parsePartnerUserID(c *gin.Context, param string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid "+param)
		return 0, false
	}
	return id, true
}

// List GET /admin/partners
func (h *PartnerHandler) List(c *gin.Context) {
	page, pageSize := response.ParsePagination(c)
	items, total, err := h.partnerService.AdminListPartners(c.Request.Context(), c.Query("query"), c.Query("status"), page, pageSize)
	if err != nil {
		partnerError(c, err)
		return
	}
	response.Paginated(c, items, int64(total), page, pageSize)
}

type upsertPartnerRequest struct {
	MinRateMultiplier float64 `json:"min_rate_multiplier" binding:"required"`
	Status            string  `json:"status"`
}

// Upsert PUT /admin/partners/:userId（创建或更新合伙人）
func (h *PartnerHandler) Upsert(c *gin.Context) {
	userID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	var req upsertPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.partnerService.AdminUpsertPartner(c.Request.Context(), partnerActorID(c), userID, req.MinRateMultiplier, req.Status); err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"user_id": userID})
}

// Delete DELETE /admin/partners/:userId（取消资格，批量解绑并清除下线倍率）
func (h *PartnerHandler) Delete(c *gin.Context) {
	userID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	if err := h.partnerService.AdminDeletePartner(c.Request.Context(), partnerActorID(c), userID); err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"user_id": userID})
}

// ListRelations GET /admin/partners/:userId/relations
func (h *PartnerHandler) ListRelations(c *gin.Context) {
	userID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	items, err := h.partnerService.AdminListDownlines(c.Request.Context(), userID)
	if err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"items": items})
}

type bindPartnerRelationRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

// Bind POST /admin/partners/:userId/relations
func (h *PartnerHandler) Bind(c *gin.Context) {
	partnerUserID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	var req bindPartnerRelationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.partnerService.AdminBind(c.Request.Context(), partnerActorID(c), partnerUserID, req.UserID); err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"user_id": req.UserID})
}

// Unbind DELETE /admin/partners/:userId/relations/:targetUserId
func (h *PartnerHandler) Unbind(c *gin.Context) {
	partnerUserID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	targetUserID, ok := parsePartnerUserID(c, "targetUserId")
	if !ok {
		return
	}
	if err := h.partnerService.AdminUnbind(c.Request.Context(), partnerActorID(c), partnerUserID, targetUserID); err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"user_id": targetUserID})
}

type partnerSettlementRequest struct {
	Amount float64 `json:"amount" binding:"required"`
	Method string  `json:"method"`
	Note   string  `json:"note"`
}

// Settle POST /admin/partners/:userId/settlements（登记线下结算）
func (h *PartnerHandler) Settle(c *gin.Context) {
	partnerUserID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	var req partnerSettlementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.partnerService.AdminSettle(c.Request.Context(), partnerActorID(c), partnerUserID, req.Amount, strings.TrimSpace(req.Note)); err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"partner_user_id": partnerUserID})
}

// ListSettlements GET /admin/partners/:userId/settlements
func (h *PartnerHandler) ListSettlements(c *gin.Context) {
	partnerUserID, ok := parsePartnerUserID(c, "userId")
	if !ok {
		return
	}
	items, err := h.partnerService.AdminListSettlements(c.Request.Context(), partnerUserID)
	if err != nil {
		partnerError(c, err)
		return
	}
	response.Success(c, gin.H{"items": items})
}
