package service

import (
	"context"
	"errors"
	"strconv"
	"time"
)

// 合伙人（Partner）制度：平台给合伙人设定最低分组倍率（成本价），
// 合伙人可给绑定下线在 [min_rate, 1.0] 区间定价，差价即合伙人收益。
const (
	PartnerStatusActive   = "active"
	PartnerStatusDisabled = "disabled"

	PartnerBoundSourceAdmin  = "admin"
	PartnerBoundSourceInvite = "invite"

	PartnerSettleMethodOffline = "offline"
)

// SetByPartner 返回 user_group_rate_multipliers.set_by 的合伙人标记。
func SetByPartner(partnerUserID int64) string {
	return "partner:" + strconv.FormatInt(partnerUserID, 10)
}

var (
	ErrPartnerNotFound        = errors.New("partner not found")
	ErrPartnerDisabled        = errors.New("partner is disabled")
	ErrPartnerInvalidMinRate  = errors.New("min rate must be in (0, 1]")
	ErrPartnerRateOutOfRange  = errors.New("rate multiplier out of partner allowed range")
	ErrPartnerAlreadyDownline = errors.New("user is already a partner downline")
	ErrPartnerIsDownline      = errors.New("user is bound as a downline, unbind first")
	ErrPartnerDownlineIsPartner = errors.New("partner cannot be bound as a downline")
	ErrPartnerNotDownline     = errors.New("user is not bound to this partner")
	ErrPartnerGroupNotAllowed = errors.New("group is not available to the downline user")
	ErrPartnerInvalidSettle   = errors.New("settlement amount exceeds pending commission")
)

// Partner 合伙人记录。
type Partner struct {
	UserID            int64     `json:"user_id"`
	MinRateMultiplier float64   `json:"min_rate_multiplier"`
	Status            string    `json:"status"`
	SettledTotal      float64   `json:"settled_total"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (p *Partner) IsActive() bool {
	return p != nil && p.Status == PartnerStatusActive
}

// PartnerListItem 管理端合伙人列表项。
type PartnerListItem struct {
	Partner
	Email             string  `json:"email"`
	Username          string  `json:"username"`
	DownlineCount     int     `json:"downline_count"`
	TotalCommission   float64 `json:"total_commission"`
	PendingSettlement float64 `json:"pending_settlement"`
}

// PartnerDownlineRate 下线用户的"用户×分组"倍率。
type PartnerDownlineRate struct {
	GroupID        int64    `json:"group_id"`
	GroupName      string   `json:"group_name"`
	RateMultiplier *float64 `json:"rate_multiplier"`
	DefaultRate    float64  `json:"default_rate"`
}

// PartnerDownline 下线用户及其定价/消费信息。
type PartnerDownline struct {
	UserID        int64                 `json:"user_id"`
	Email         string                `json:"email"`
	Username      string                `json:"username"`
	BoundSource   string                `json:"bound_source"`
	CreatedAt     time.Time             `json:"created_at"`
	Rates         []PartnerDownlineRate `json:"rates"`
	TotalCost30d  float64               `json:"total_cost_30d"`
	Commission30d float64               `json:"commission_30d"`
}

// PartnerSettlement 结算台账记录。
type PartnerSettlement struct {
	ID            int64     `json:"id"`
	PartnerUserID int64     `json:"partner_user_id"`
	Amount        float64   `json:"amount"`
	Method        string    `json:"method"`
	Note          string    `json:"note"`
	OperatorID    int64     `json:"operator_id"`
	CreatedAt     time.Time `json:"created_at"`
}

// PartnerCommissionRecord 收益明细（usage_logs 快照 + 用户邮箱）。
type PartnerCommissionRecord struct {
	ID             int64      `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UserID         int64      `json:"user_id"`
	UserEmail      string     `json:"user_email"`
	Model          string     `json:"model"`
	TotalCost      float64    `json:"total_cost"`
	BaseRate       *float64   `json:"base_rate"`
	RateFloor      *float64   `json:"rate_floor"`
	RateMultiplier float64    `json:"rate_multiplier"`
	Commission     float64    `json:"commission"`
}

// PartnerDownlineStats 下线消费/收益聚合。
type PartnerDownlineStats struct {
	ActualCost float64
	Commission float64
}

// PartnerUsageShare 计费侧快照：命中"该下线当前绑定合伙人设置的倍率"时返回。
type PartnerUsageShare struct {
	PartnerUserID int64
	MinRate       float64
	BaseRate      float64
}

// PartnerCommissionFilter 收益明细筛选。
type PartnerCommissionFilter struct {
	UserID    int64
	StartTime *time.Time
	EndTime   *time.Time
	Page      int
	PageSize  int
}

// PartnerOverview 合伙人面板概览。
type PartnerOverview struct {
	IsPartner         bool    `json:"is_partner"`
	MinRateMultiplier float64 `json:"min_rate_multiplier"`
	Status            string  `json:"status"`
	DownlineCount     int     `json:"downline_count"`
	TotalCommission   float64 `json:"total_commission"`
	MonthCommission   float64 `json:"month_commission"`
	SettledTotal      float64 `json:"settled_total"`
	PendingSettlement float64 `json:"pending_settlement"`
	AffCode           string  `json:"aff_code"`
}

// PartnerRepository 合伙人数据访问接口。
type PartnerRepository interface {
	GetPartnerByUserID(ctx context.Context, userID int64) (*Partner, error)
	// UpsertPartnerWithClamp 创建/更新合伙人（min_rate、status），
	// 同事务将该合伙人设置的低于 min_rate 的下线倍率抬升到 min_rate，并写审计。
	UpsertPartnerWithClamp(ctx context.Context, p *Partner, actorID int64) error
	// DeletePartnerAndUnbind 取消合伙人资格：同事务批量解绑全部下线并清除其合伙人倍率。
	DeletePartnerAndUnbind(ctx context.Context, partnerUserID, actorID int64) error
	ListPartners(ctx context.Context, query, status string, page, pageSize int) ([]PartnerListItem, int, error)
	GetPartnerUserIDs(ctx context.Context, userIDs []int64) (map[int64]bool, error)

	GetRelationByUserID(ctx context.Context, userID int64) (*PartnerRelation, error)
	// BindRelation 绑定/换绑：同事务清除原合伙人的倍率并建立新绑定，写审计。
	BindRelation(ctx context.Context, partnerUserID, userID int64, source string, actorID int64) error
	// UnbindRelation 解绑并清除该合伙人给此用户设置的倍率。
	UnbindRelation(ctx context.Context, partnerUserID, userID, actorID int64) error
	ListDownlines(ctx context.Context, partnerUserID int64) ([]PartnerDownline, error)

	// SetDownlineRate 合伙人给下线设置"用户×分组"倍率（set_by=partner:<id>）。
	SetDownlineRate(ctx context.Context, partnerUserID, userID, groupID int64, rate float64) error
	// ClearDownlineRate 清除该组合伙人定价（回分组默认，来源归还 admin）。
	ClearDownlineRate(ctx context.Context, partnerUserID, userID, groupID int64) error

	// ResolveUsageShare 计费侧：仅当该用户×分组命中"当前绑定合伙人设置的倍率"时返回快照。
	ResolveUsageShare(ctx context.Context, userID, groupID int64) (*PartnerUsageShare, error)

	GetCommissionTotals(ctx context.Context, partnerUserID int64) (total, month float64, err error)
	GetDownlineStats(ctx context.Context, partnerUserID int64, since time.Time) (map[int64]PartnerDownlineStats, error)
	ListCommissionRecords(ctx context.Context, partnerUserID int64, filter PartnerCommissionFilter) ([]PartnerCommissionRecord, int, float64, error)

	// CreateSettlement 登记结算：插台账并累加 settled_total（同事务校验待结算余额）。
	CreateSettlement(ctx context.Context, s *PartnerSettlement) error
	ListSettlements(ctx context.Context, partnerUserID int64, limit int) ([]PartnerSettlement, error)

	WriteAudit(ctx context.Context, actorID int64, action string, targetUserID *int64, detail map[string]any) error
}

// PartnerRelation 绑定关系。
type PartnerRelation struct {
	ID            int64     `json:"id"`
	PartnerUserID int64     `json:"partner_user_id"`
	UserID        int64     `json:"user_id"`
	BoundSource   string    `json:"bound_source"`
	CreatedAt     time.Time `json:"created_at"`
}
