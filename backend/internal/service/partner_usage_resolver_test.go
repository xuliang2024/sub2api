package service

import (
	"context"
	"errors"
	"testing"
	"time"
)

type stubPartnerUsageRepo struct {
	PartnerRepository
	share *PartnerUsageShare
	err   error
	calls int
}

func (s *stubPartnerUsageRepo) ResolveUsageShare(ctx context.Context, userID, groupID int64) (*PartnerUsageShare, error) {
	s.calls++
	return s.share, s.err
}

func TestApplyPartnerCommission(t *testing.T) {
	share := &PartnerUsageShare{PartnerUserID: 42, MinRate: 0.5, BaseRate: 0.8}
	log := &UsageLog{}
	ApplyPartnerCommission(log, share, 1.0)
	if log.PartnerID == nil || *log.PartnerID != 42 {
		t.Fatalf("partner id not set: %+v", log.PartnerID)
	}
	if log.PartnerBaseRate == nil || *log.PartnerBaseRate != 0.8 {
		t.Fatalf("base rate snapshot wrong: %+v", log.PartnerBaseRate)
	}
	if log.PartnerRateFloor == nil || *log.PartnerRateFloor != 0.5 {
		t.Fatalf("rate floor snapshot wrong: %+v", log.PartnerRateFloor)
	}
	if log.PartnerCommission != 0.3 {
		// 1.0 * (0.8-0.5) = 0.3（浮点差经 roundTo 8 位归一）
		t.Fatalf("commission = %v, want 0.3", log.PartnerCommission)
	}
}

func TestApplyPartnerCommissionZeroAndClamp(t *testing.T) {
	// 设到最低倍率 = 不赚钱，但仍记录 partner_id 供下线消费统计
	log := &UsageLog{}
	ApplyPartnerCommission(log, &PartnerUsageShare{PartnerUserID: 7, MinRate: 0.5, BaseRate: 0.5}, 2.0)
	if log.PartnerID == nil || log.PartnerCommission != 0 {
		t.Fatalf("floor rate should record partner with zero commission: %+v", log)
	}
	// 低于最低倍率：clamp 为 0
	log2 := &UsageLog{}
	ApplyPartnerCommission(log2, &PartnerUsageShare{PartnerUserID: 7, MinRate: 0.8, BaseRate: 0.6}, 2.0)
	if log2.PartnerCommission != 0 {
		t.Fatalf("below floor should clamp to zero: %v", log2.PartnerCommission)
	}
	// 零费用：无收益
	log3 := &UsageLog{}
	ApplyPartnerCommission(log3, &PartnerUsageShare{PartnerUserID: 7, MinRate: 0.5, BaseRate: 0.8}, 0)
	if log3.PartnerCommission != 0 {
		t.Fatalf("zero cost should produce zero commission: %v", log3.PartnerCommission)
	}
	// nil 快照：不写入
	log4 := &UsageLog{}
	ApplyPartnerCommission(log4, nil, 1.0)
	if log4.PartnerID != nil {
		t.Fatalf("nil share should not set partner: %+v", log4.PartnerID)
	}
}

func TestPartnerUsageResolverCacheAndFallback(t *testing.T) {
	repo := &stubPartnerUsageRepo{share: &PartnerUsageShare{PartnerUserID: 1, MinRate: 0.5, BaseRate: 0.8}}
	resolver := NewPartnerUsageResolver(repo, time.Minute)
	if got := resolver.Resolve(context.Background(), 100, 31); got == nil || got.PartnerUserID != 1 {
		t.Fatalf("resolve failed: %+v", got)
	}
	if got := resolver.Resolve(context.Background(), 100, 31); got == nil {
		t.Fatalf("cached resolve failed")
	}
	if repo.calls != 1 {
		t.Fatalf("expected 1 repo call (cache hit on second), got %d", repo.calls)
	}

	// 未命中（无绑定）应缓存 nil，不重复查询
	repo2 := &stubPartnerUsageRepo{}
	resolver2 := NewPartnerUsageResolver(repo2, time.Minute)
	if got := resolver2.Resolve(context.Background(), 100, 31); got != nil {
		t.Fatalf("expected nil share")
	}
	_ = resolver2.Resolve(context.Background(), 100, 31)
	if repo2.calls != 1 {
		t.Fatalf("nil share should be cached, calls = %d", repo2.calls)
	}

	// 出错降级为 nil，不阻塞计费
	repo3 := &stubPartnerUsageRepo{err: errors.New("db down")}
	resolver3 := NewPartnerUsageResolver(repo3, time.Minute)
	if got := resolver3.Resolve(context.Background(), 100, 31); got != nil {
		t.Fatalf("error should degrade to nil")
	}

	// nil resolver 安全
	var nilResolver *PartnerUsageResolver
	if got := nilResolver.Resolve(context.Background(), 100, 31); got != nil {
		t.Fatalf("nil resolver should return nil")
	}
}

func TestPartnerIndependentMediaBilling(t *testing.T) {
	group := &Group{ImageRateIndependent: true, VideoRateIndependent: true}
	if !partnerIndependentMediaBilling(1, 0, group) {
		t.Fatal("independent image billing should be excluded")
	}
	if !partnerIndependentMediaBilling(0, 2, group) {
		t.Fatal("independent video billing should be excluded")
	}
	if partnerIndependentMediaBilling(0, 0, group) {
		t.Fatal("token billing should not be excluded")
	}
	if partnerIndependentMediaBilling(1, 1, &Group{}) {
		t.Fatal("non-independent media should not be excluded")
	}
}
