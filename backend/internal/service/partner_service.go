package service

import (
	"context"
	"time"
)

// PartnerGroupLister 提供"用户当前可用分组"（复用 APIKeyService.GetAvailableGroups）。
type PartnerGroupLister interface {
	GetAvailableGroups(ctx context.Context, userID int64) ([]Group, error)
}

// PartnerRateInput 合伙人设价条目；RateMultiplier 为 nil 表示清除该组定价（回分组默认）。
type PartnerRateInput struct {
	GroupID        int64    `json:"group_id"`
	RateMultiplier *float64 `json:"rate_multiplier"`
}

// PartnerService 合伙人业务逻辑。
type PartnerService struct {
	repo              PartnerRepository
	userRepo          UserRepository
	userGroupRateRepo UserGroupRateRepository
	groupLister       PartnerGroupLister
}

func NewPartnerService(
	repo PartnerRepository,
	userRepo UserRepository,
	userGroupRateRepo UserGroupRateRepository,
	groupLister PartnerGroupLister,
) *PartnerService {
	return &PartnerService{
		repo:              repo,
		userRepo:          userRepo,
		userGroupRateRepo: userGroupRateRepo,
		groupLister:       groupLister,
	}
}

func (s *PartnerService) getActivePartner(ctx context.Context, userID int64) (*Partner, error) {
	p, err := s.repo.GetPartnerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrPartnerNotFound
	}
	if !p.IsActive() {
		return nil, ErrPartnerDisabled
	}
	return p, nil
}

// IsPartner 用于菜单/列表标记：有合伙人记录即为 true（含 disabled）。
func (s *PartnerService) IsPartner(ctx context.Context, userID int64) bool {
	if s == nil || s.repo == nil || userID <= 0 {
		return false
	}
	p, err := s.repo.GetPartnerByUserID(ctx, userID)
	return err == nil && p != nil
}

// GetPartnerUserIDs 批量标记（管理端用户列表）。
func (s *PartnerService) GetPartnerUserIDs(ctx context.Context, userIDs []int64) (map[int64]bool, error) {
	if s == nil || s.repo == nil {
		return map[int64]bool{}, nil
	}
	return s.repo.GetPartnerUserIDs(ctx, userIDs)
}

// GetOverview 合伙人面板概览（非合伙人返回 ErrPartnerNotFound）。
func (s *PartnerService) GetOverview(ctx context.Context, userID int64) (*PartnerOverview, error) {
	p, err := s.repo.GetPartnerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrPartnerNotFound
	}
	total, month, err := s.repo.GetCommissionTotals(ctx, userID)
	if err != nil {
		return nil, err
	}
	downlines, err := s.repo.ListDownlines(ctx, userID)
	if err != nil {
		return nil, err
	}
	pending := total - p.SettledTotal
	if pending < 0 {
		pending = 0
	}
	return &PartnerOverview{
		IsPartner:         true,
		MinRateMultiplier: p.MinRateMultiplier,
		Status:            p.Status,
		DownlineCount:     len(downlines),
		TotalCommission:   total,
		MonthCommission:   month,
		SettledTotal:      p.SettledTotal,
		PendingSettlement: pending,
	}, nil
}

// ListDownlines 合伙人查看自己的下线（须 active）。
func (s *PartnerService) ListDownlines(ctx context.Context, partnerUserID int64) ([]PartnerDownline, error) {
	if _, err := s.getActivePartner(ctx, partnerUserID); err != nil {
		return nil, err
	}
	return s.listDownlines(ctx, partnerUserID)
}

// AdminListDownlines 管理端查看下线（disabled 也可见）。
func (s *PartnerService) AdminListDownlines(ctx context.Context, partnerUserID int64) ([]PartnerDownline, error) {
	p, err := s.repo.GetPartnerByUserID(ctx, partnerUserID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrPartnerNotFound
	}
	return s.listDownlines(ctx, partnerUserID)
}

func (s *PartnerService) listDownlines(ctx context.Context, partnerUserID int64) ([]PartnerDownline, error) {
	items, err := s.repo.ListDownlines(ctx, partnerUserID)
	if err != nil {
		return nil, err
	}
	stats, err := s.repo.GetDownlineStats(ctx, partnerUserID, time.Now().AddDate(0, 0, -30))
	if err != nil {
		return nil, err
	}
	for i := range items {
		d := &items[i]
		if st, ok := stats[d.UserID]; ok {
			d.TotalCost30d = st.ActualCost
			d.Commission30d = st.Commission
		}
		d.Rates = make([]PartnerDownlineRate, 0)
		if s.groupLister == nil {
			continue
		}
		groups, err := s.groupLister.GetAvailableGroups(ctx, d.UserID)
		if err != nil {
			return nil, err
		}
		rateMap := map[int64]float64{}
		if s.userGroupRateRepo != nil {
			if m, err := s.userGroupRateRepo.GetByUserID(ctx, d.UserID); err == nil {
				rateMap = m
			}
		}
		for _, g := range groups {
			rate := PartnerDownlineRate{
				GroupID:     g.ID,
				GroupName:   g.Name,
				DefaultRate: g.RateMultiplier,
			}
			if v, ok := rateMap[g.ID]; ok {
				vv := v
				rate.RateMultiplier = &vv
			}
			d.Rates = append(d.Rates, rate)
		}
	}
	return items, nil
}

// SetDownlineRates 合伙人给下线设价：绑定归属 + 分组可用性 + [min_rate, 1.0] 范围校验。
func (s *PartnerService) SetDownlineRates(ctx context.Context, partnerUserID, userID int64, rates []PartnerRateInput) error {
	p, err := s.getActivePartner(ctx, partnerUserID)
	if err != nil {
		return err
	}
	rel, err := s.repo.GetRelationByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if rel == nil || rel.PartnerUserID != partnerUserID {
		return ErrPartnerNotDownline
	}
	allowed := map[int64]bool{}
	if s.groupLister != nil {
		groups, err := s.groupLister.GetAvailableGroups(ctx, userID)
		if err != nil {
			return err
		}
		for _, g := range groups {
			allowed[g.ID] = true
		}
	}
	for _, input := range rates {
		if input.GroupID <= 0 {
			continue
		}
		if len(allowed) > 0 && !allowed[input.GroupID] {
			return ErrPartnerGroupNotAllowed
		}
		if input.RateMultiplier == nil {
			if err := s.repo.ClearDownlineRate(ctx, partnerUserID, userID, input.GroupID); err != nil {
				return err
			}
			continue
		}
		rate := *input.RateMultiplier
		if rate < p.MinRateMultiplier-1e-9 || rate > 1.0+1e-9 {
			return ErrPartnerRateOutOfRange
		}
		if err := s.repo.SetDownlineRate(ctx, partnerUserID, userID, input.GroupID, rate); err != nil {
			return err
		}
	}
	detail := map[string]any{"partner_user_id": partnerUserID, "rates": rates}
	return s.repo.WriteAudit(ctx, partnerUserID, "set_rate", &userID, detail)
}

// BindByInviter 注册邀请自动绑定：inviter 是 active 合伙人时建立绑定。
// 绑定失败不阻断注册（由调用方记录日志）。
func (s *PartnerService) BindByInviter(ctx context.Context, inviterUserID, newUserID int64) error {
	if s == nil || s.repo == nil || inviterUserID <= 0 || newUserID <= 0 {
		return nil
	}
	p, err := s.repo.GetPartnerByUserID(ctx, inviterUserID)
	if err != nil || p == nil || !p.IsActive() {
		return err
	}
	return s.repo.BindRelation(ctx, inviterUserID, newUserID, PartnerBoundSourceInvite, inviterUserID)
}

// AdminListPartners 管理端合伙人列表。
func (s *PartnerService) AdminListPartners(ctx context.Context, query, status string, page, pageSize int) ([]PartnerListItem, int, error) {
	return s.repo.ListPartners(ctx, query, status, page, pageSize)
}

// AdminUpsertPartner 设置/更新合伙人（含最低倍率上调后的存量 clamp）。
func (s *PartnerService) AdminUpsertPartner(ctx context.Context, actorID, userID int64, minRate float64, status string) error {
	if minRate <= 0 || minRate > 1 {
		return ErrPartnerInvalidMinRate
	}
	if status == "" {
		status = PartnerStatusActive
	}
	if status != PartnerStatusActive && status != PartnerStatusDisabled {
		return ErrPartnerInvalidMinRate
	}
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return err
	}
	rel, err := s.repo.GetRelationByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if rel != nil && rel.PartnerUserID != userID {
		// 已是别人的下线，须先解绑才能设为合伙人（互斥约束）。
		return ErrPartnerIsDownline
	}
	return s.repo.UpsertPartnerWithClamp(ctx, &Partner{
		UserID:            userID,
		MinRateMultiplier: minRate,
		Status:            status,
	}, actorID)
}

// AdminDeletePartner 取消合伙人资格（批量解绑 + 清除下线倍率）。
func (s *PartnerService) AdminDeletePartner(ctx context.Context, actorID, userID int64) error {
	return s.repo.DeletePartnerAndUnbind(ctx, userID, actorID)
}

// AdminBind 管理端手动绑定/换绑。
func (s *PartnerService) AdminBind(ctx context.Context, actorID, partnerUserID, userID int64) error {
	if _, err := s.getActivePartner(ctx, partnerUserID); err != nil {
		return err
	}
	if _, err := s.userRepo.GetByID(ctx, userID); err != nil {
		return err
	}
	return s.repo.BindRelation(ctx, partnerUserID, userID, PartnerBoundSourceAdmin, actorID)
}

// AdminUnbind 管理端解绑。
func (s *PartnerService) AdminUnbind(ctx context.Context, actorID, partnerUserID, userID int64) error {
	return s.repo.UnbindRelation(ctx, partnerUserID, userID, actorID)
}

// AdminSettle 登记线下结算。
func (s *PartnerService) AdminSettle(ctx context.Context, actorID, partnerUserID int64, amount float64, note string) error {
	if amount <= 0 {
		return ErrPartnerInvalidSettle
	}
	return s.repo.CreateSettlement(ctx, &PartnerSettlement{
		PartnerUserID: partnerUserID,
		Amount:        amount,
		Method:        PartnerSettleMethodOffline,
		Note:          note,
		OperatorID:    actorID,
	})
}

// AdminListSettlements 结算记录。
func (s *PartnerService) AdminListSettlements(ctx context.Context, partnerUserID int64) ([]PartnerSettlement, error) {
	return s.repo.ListSettlements(ctx, partnerUserID, 50)
}

// ListCommissionRecords 合伙人收益明细。
func (s *PartnerService) ListCommissionRecords(ctx context.Context, partnerUserID int64, filter PartnerCommissionFilter) ([]PartnerCommissionRecord, int, float64, error) {
	if _, err := s.getActivePartner(ctx, partnerUserID); err != nil {
		return nil, 0, 0, err
	}
	return s.repo.ListCommissionRecords(ctx, partnerUserID, filter)
}
