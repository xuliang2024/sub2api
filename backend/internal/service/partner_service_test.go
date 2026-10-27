package service

import (
	"context"
	"errors"
	"testing"
)

type setRateCall struct {
	partnerID, userID, groupID int64
	rate                       float64
}

type partnerServiceRepoStub struct {
	PartnerRepository
	partner       *Partner
	relation      *PartnerRelation
	setRateCalls  []setRateCall
	clearCalls    int
	bindCalls     []setRateCall
	bindSource    string
	upsertCalls   int
	upsertPartner *Partner
}

func (s *partnerServiceRepoStub) GetPartnerByUserID(ctx context.Context, userID int64) (*Partner, error) {
	if s.partner != nil && s.partner.UserID == userID {
		return s.partner, nil
	}
	return nil, nil
}

func (s *partnerServiceRepoStub) GetRelationByUserID(ctx context.Context, userID int64) (*PartnerRelation, error) {
	if s.relation != nil && s.relation.UserID == userID {
		return s.relation, nil
	}
	return nil, nil
}

func (s *partnerServiceRepoStub) SetDownlineRate(ctx context.Context, partnerUserID, userID, groupID int64, rate float64) error {
	s.setRateCalls = append(s.setRateCalls, setRateCall{partnerUserID, userID, groupID, rate})
	return nil
}

func (s *partnerServiceRepoStub) ClearDownlineRate(ctx context.Context, partnerUserID, userID, groupID int64) error {
	s.clearCalls++
	return nil
}

func (s *partnerServiceRepoStub) WriteAudit(ctx context.Context, actorID int64, action string, targetUserID *int64, detail map[string]any) error {
	return nil
}

func (s *partnerServiceRepoStub) BindRelation(ctx context.Context, partnerUserID, userID int64, source string, actorID int64) error {
	s.bindCalls = append(s.bindCalls, setRateCall{partnerID: partnerUserID, userID: userID})
	s.bindSource = source
	return nil
}

func (s *partnerServiceRepoStub) UpsertPartnerWithClamp(ctx context.Context, p *Partner, actorID int64) error {
	s.upsertCalls++
	s.upsertPartner = p
	return nil
}

type partnerGroupListerStub struct {
	groups []Group
}

func (s *partnerGroupListerStub) GetAvailableGroups(ctx context.Context, userID int64) ([]Group, error) {
	return s.groups, nil
}

type partnerUserRepoStub struct {
	UserRepository
}

func (s *partnerUserRepoStub) GetByID(ctx context.Context, id int64) (*User, error) {
	return &User{ID: id}, nil
}

func newPartnerServiceForTest(repo *partnerServiceRepoStub) *PartnerService {
	return NewPartnerService(repo, &partnerUserRepoStub{}, nil, &partnerGroupListerStub{
		groups: []Group{{ID: 31, Name: "国产全模型", RateMultiplier: 1}, {ID: 25, Name: "聚合", RateMultiplier: 1}},
	})
}

func activePartner() *Partner {
	return &Partner{UserID: 42, MinRateMultiplier: 0.5, Status: PartnerStatusActive}
}

func TestSetDownlineRatesRangeAndBinding(t *testing.T) {
	repo := &partnerServiceRepoStub{
		partner:  activePartner(),
		relation: &PartnerRelation{PartnerUserID: 42, UserID: 100},
	}
	svc := newPartnerServiceForTest(repo)
	rate := func(v float64) *float64 { return &v }

	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: rate(0.8)}}); err != nil {
		t.Fatalf("valid rate rejected: %v", err)
	}
	if len(repo.setRateCalls) != 1 || repo.setRateCalls[0].rate != 0.8 {
		t.Fatalf("set rate not recorded: %+v", repo.setRateCalls)
	}

	// 边界：下限=最低倍率（不赚钱）、上限=1.0 都允许
	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: rate(0.5)}}); err != nil {
		t.Fatalf("floor rate should be allowed: %v", err)
	}
	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: rate(1.0)}}); err != nil {
		t.Fatalf("1.0 should be allowed: %v", err)
	}

	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: rate(0.4)}}); !errors.Is(err, ErrPartnerRateOutOfRange) {
		t.Fatalf("below floor should be rejected: %v", err)
	}
	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: rate(1.2)}}); !errors.Is(err, ErrPartnerRateOutOfRange) {
		t.Fatalf("above 1.0 should be rejected: %v", err)
	}
	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 99, RateMultiplier: rate(0.8)}}); !errors.Is(err, ErrPartnerGroupNotAllowed) {
		t.Fatalf("unavailable group should be rejected: %v", err)
	}
	if err := svc.SetDownlineRates(context.Background(), 42, 999, []PartnerRateInput{{GroupID: 31, RateMultiplier: rate(0.8)}}); !errors.Is(err, ErrPartnerNotDownline) {
		t.Fatalf("unbound user should be rejected: %v", err)
	}
	// nil = 清除该组定价
	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: nil}}); err != nil {
		t.Fatalf("clear rate failed: %v", err)
	}
	if repo.clearCalls != 1 {
		t.Fatalf("clear rate should be called once, got %d", repo.clearCalls)
	}
}

func TestSetDownlineRatesDisabledPartner(t *testing.T) {
	repo := &partnerServiceRepoStub{
		partner:  &Partner{UserID: 42, MinRateMultiplier: 0.5, Status: PartnerStatusDisabled},
		relation: &PartnerRelation{PartnerUserID: 42, UserID: 100},
	}
	svc := newPartnerServiceForTest(repo)
	rate := 0.8
	if err := svc.SetDownlineRates(context.Background(), 42, 100, []PartnerRateInput{{GroupID: 31, RateMultiplier: &rate}}); !errors.Is(err, ErrPartnerDisabled) {
		t.Fatalf("disabled partner should be rejected: %v", err)
	}
}

func TestBindByInviter(t *testing.T) {
	// inviter 不是合伙人：不建立绑定
	repo := &partnerServiceRepoStub{}
	svc := newPartnerServiceForTest(repo)
	if err := svc.BindByInviter(context.Background(), 42, 100); err != nil {
		t.Fatalf("non-partner inviter should be a no-op: %v", err)
	}
	if len(repo.bindCalls) != 0 {
		t.Fatalf("should not bind for non-partner inviter")
	}

	// inviter 是 active 合伙人：自动绑定，来源 invite
	repo2 := &partnerServiceRepoStub{partner: activePartner()}
	svc2 := newPartnerServiceForTest(repo2)
	if err := svc2.BindByInviter(context.Background(), 42, 100); err != nil {
		t.Fatalf("bind failed: %v", err)
	}
	if len(repo2.bindCalls) != 1 || repo2.bindSource != PartnerBoundSourceInvite {
		t.Fatalf("invite binding not recorded: %+v source=%s", repo2.bindCalls, repo2.bindSource)
	}

	// inviter 被停用：不绑定
	repo3 := &partnerServiceRepoStub{partner: &Partner{UserID: 42, MinRateMultiplier: 0.5, Status: PartnerStatusDisabled}}
	svc3 := newPartnerServiceForTest(repo3)
	if err := svc3.BindByInviter(context.Background(), 42, 100); err != nil {
		t.Fatalf("disabled inviter should be a no-op: %v", err)
	}
	if len(repo3.bindCalls) != 0 {
		t.Fatalf("should not bind for disabled inviter")
	}
}

func TestAdminUpsertPartnerValidation(t *testing.T) {
	repo := &partnerServiceRepoStub{}
	svc := newPartnerServiceForTest(repo)
	if err := svc.AdminUpsertPartner(context.Background(), 1, 42, 0, ""); !errors.Is(err, ErrPartnerInvalidMinRate) {
		t.Fatalf("zero min rate should be rejected: %v", err)
	}
	if err := svc.AdminUpsertPartner(context.Background(), 1, 42, 1.2, ""); !errors.Is(err, ErrPartnerInvalidMinRate) {
		t.Fatalf("min rate above 1 should be rejected: %v", err)
	}
	if err := svc.AdminUpsertPartner(context.Background(), 1, 42, 0.5, ""); err != nil {
		t.Fatalf("valid upsert failed: %v", err)
	}
	if repo.upsertCalls != 1 || repo.upsertPartner.Status != PartnerStatusActive {
		t.Fatalf("upsert not recorded correctly: %+v", repo.upsertPartner)
	}

	// 已是别人下线：须先解绑
	repo2 := &partnerServiceRepoStub{relation: &PartnerRelation{PartnerUserID: 7, UserID: 42}}
	svc2 := newPartnerServiceForTest(repo2)
	if err := svc2.AdminUpsertPartner(context.Background(), 1, 42, 0.5, ""); !errors.Is(err, ErrPartnerIsDownline) {
		t.Fatalf("downline user should not become partner directly: %v", err)
	}
}
