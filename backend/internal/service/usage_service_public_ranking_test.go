package service

import (
	"context"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	"github.com/stretchr/testify/require"
)

type publicRankingUsageRepo struct {
	UsageLogRepository
	limit int
}

func (r *publicRankingUsageRepo) GetUserSpendingRanking(ctx context.Context, startTime, endTime time.Time, limit int) (*usagestats.UserSpendingRankingResponse, error) {
	r.limit = limit
	return &usagestats.UserSpendingRankingResponse{
		Ranking: []usagestats.UserSpendingRankingItem{
			{UserID: 11, Email: "alice@example.com", ActualCost: 12.3, Requests: 7, Tokens: 9000},
			{UserID: 42, Email: "bob@domain.test", ActualCost: 8.1, Requests: 5, Tokens: 7000},
		},
		TotalActualCost: 20.4,
		TotalRequests:   12,
		TotalTokens:     16000,
	}, nil
}

func TestGetPublicUserSpendingRankingMasksEmailAndMarksCurrentUser(t *testing.T) {
	repo := &publicRankingUsageRepo{}
	svc := NewUsageService(repo, nil, nil, nil)

	got, err := svc.GetPublicUserSpendingRanking(context.Background(), 42, time.Time{}, time.Now(), 99)
	require.NoError(t, err)

	require.Equal(t, 10, repo.limit)
	require.Equal(t, 20.4, got.TotalActualCost)
	require.Len(t, got.Ranking, 2)
	require.Equal(t, usagestats.PublicUserSpendingRankingItem{
		Rank:          1,
		DisplayName:   "a***e@e***e.com",
		ActualCost:    12.3,
		Requests:      7,
		Tokens:        9000,
		IsCurrentUser: false,
	}, got.Ranking[0])
	require.Equal(t, "b***b@d***n.test", got.Ranking[1].DisplayName)
	require.True(t, got.Ranking[1].IsCurrentUser)
}
