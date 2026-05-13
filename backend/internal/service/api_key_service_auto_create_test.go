//go:build unit

package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type autoCreateAPIKeyRepo struct {
	apiKeyRepoStub
	keys      []APIKey
	nextID    int64
	createNum int
}

func (r *autoCreateAPIKeyRepo) Create(_ context.Context, key *APIKey) error {
	if key.ID == 0 {
		if r.nextID == 0 {
			r.nextID = 100
		}
		key.ID = r.nextID
		r.nextID++
	}
	now := time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
	if key.CreatedAt.IsZero() {
		key.CreatedAt = now
	}
	if key.UpdatedAt.IsZero() {
		key.UpdatedAt = now
	}
	clone := *key
	r.keys = append(r.keys, clone)
	r.createNum++
	return nil
}

func (r *autoCreateAPIKeyRepo) ListByUserID(_ context.Context, userID int64, params pagination.PaginationParams, filters APIKeyListFilters) ([]APIKey, *pagination.PaginationResult, error) {
	matched := make([]APIKey, 0, len(r.keys))
	for i := range r.keys {
		key := r.keys[i]
		if key.UserID != userID {
			continue
		}
		if filters.Search != "" && !strings.Contains(strings.ToLower(key.Name), strings.ToLower(filters.Search)) && !strings.Contains(strings.ToLower(key.Key), strings.ToLower(filters.Search)) {
			continue
		}
		if filters.Status != "" && key.Status != filters.Status {
			continue
		}
		if filters.GroupID != nil {
			if *filters.GroupID == 0 {
				if key.GroupID != nil {
					continue
				}
			} else if key.GroupID == nil || *key.GroupID != *filters.GroupID {
				continue
			}
		}
		matched = append(matched, key)
	}

	total := int64(len(matched))
	start := params.Offset()
	if start > len(matched) {
		start = len(matched)
	}
	end := start + params.Limit()
	if end > len(matched) {
		end = len(matched)
	}
	pageSize := params.Limit()
	pages := 1
	if pageSize > 0 && total > 0 {
		pages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}
	return append([]APIKey(nil), matched[start:end]...), &pagination.PaginationResult{
		Total:    total,
		Page:     params.Page,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

func (r *autoCreateAPIKeyRepo) CountByUserID(_ context.Context, userID int64) (int64, error) {
	var count int64
	for i := range r.keys {
		if r.keys[i].UserID == userID {
			count++
		}
	}
	return count, nil
}

func (r *autoCreateAPIKeyRepo) ExistsByKey(_ context.Context, key string) (bool, error) {
	for i := range r.keys {
		if r.keys[i].Key == key {
			return true, nil
		}
	}
	return false, nil
}

type autoCreateGroupRepo struct {
	groupRepoNoop
	groups []Group
}

func (r *autoCreateGroupRepo) ListActive(context.Context) ([]Group, error) {
	return append([]Group(nil), r.groups...), nil
}

func (r *autoCreateGroupRepo) GetByID(_ context.Context, id int64) (*Group, error) {
	for i := range r.groups {
		if r.groups[i].ID == id {
			clone := r.groups[i]
			return &clone, nil
		}
	}
	return nil, ErrGroupNotFound
}

func newAutoCreateAPIKeyService(repo *autoCreateAPIKeyRepo, groups []Group) *APIKeyService {
	userRepo := &userRepoStub{
		user: &User{
			ID:     42,
			Email:  "user@example.com",
			Role:   RoleUser,
			Status: StatusActive,
		},
	}
	return NewAPIKeyService(
		repo,
		userRepo,
		&autoCreateGroupRepo{groups: groups},
		userSubRepoNoop{},
		nil,
		&apiKeyCacheStub{},
		&config.Config{Default: config.DefaultConfig{APIKeyPrefix: "sk-"}},
	)
}

func TestAPIKeyServiceListAutoCreatesDefaultCodexKey(t *testing.T) {
	repo := &autoCreateAPIKeyRepo{}
	codexGroupID := int64(4)
	svc := newAutoCreateAPIKeyService(repo, []Group{
		{ID: 3, Name: "codex", Platform: PlatformAnthropic, Status: StatusActive, SubscriptionType: SubscriptionTypeStandard},
		{ID: codexGroupID, Name: "codex", Platform: PlatformOpenAI, Status: StatusActive, SubscriptionType: SubscriptionTypeStandard},
	})

	keys, page, err := svc.List(context.Background(), 42, pagination.PaginationParams{Page: 1, PageSize: 20, SortBy: "created_at", SortOrder: "desc"}, APIKeyListFilters{})
	require.NoError(t, err)
	require.Len(t, keys, 1)
	require.Equal(t, int64(1), page.Total)
	require.Equal(t, 1, repo.createNum)
	require.Equal(t, apiKeyDefaultAutoName, keys[0].Name)
	require.NotNil(t, keys[0].GroupID)
	require.Equal(t, codexGroupID, *keys[0].GroupID)
	require.True(t, strings.HasPrefix(keys[0].Key, "sk-"))
}

func TestAPIKeyServiceListDoesNotAutoCreateForFilteredEmptyList(t *testing.T) {
	repo := &autoCreateAPIKeyRepo{}
	svc := newAutoCreateAPIKeyService(repo, []Group{
		{ID: 4, Name: "codex", Platform: PlatformOpenAI, Status: StatusActive, SubscriptionType: SubscriptionTypeStandard},
	})

	keys, page, err := svc.List(context.Background(), 42, pagination.PaginationParams{Page: 1, PageSize: 20}, APIKeyListFilters{Search: "missing"})
	require.NoError(t, err)
	require.Empty(t, keys)
	require.Equal(t, int64(0), page.Total)
	require.Equal(t, 0, repo.createNum)
}
