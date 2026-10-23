package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

// campaignRepository 落地页配置的原生 SQL 实现。
//
// 为什么不用 ent：campaigns 是配置型表，字段几乎全是 JSONB，用 ent 生成代码收益低
// 而成本高（需要重新 generate 整个 ent 包）。仓储只做 CRUD，与服务层解耦。
type campaignRepository struct {
	client *dbent.Client
	sql    sqlExecutor
}

// NewCampaignRepository 创建落地页仓储。
func NewCampaignRepository(client *dbent.Client, sqlDB *sql.DB) service.CampaignRepository {
	return &campaignRepository{client: client, sql: sqlDB}
}

const campaignSelectColumns = `
	id, slug, name, status, COALESCE(host, '') AS host, group_id, locale,
	brand, hero, content_md, model_entries, code_samples, cta, trial, seo, translations,
	created_by, created_at, updated_at`

// GetByID 按主键查询。
func (r *campaignRepository) GetByID(ctx context.Context, id int64) (*service.Campaign, error) {
	return r.queryOne(ctx, `SELECT`+campaignSelectColumns+` FROM campaigns WHERE id = $1`, id)
}

// GetBySlug 按 slug 查询（含非 active，用于管理端与预览）。
func (r *campaignRepository) GetBySlug(ctx context.Context, slug string) (*service.Campaign, error) {
	return r.queryOne(ctx, `SELECT`+campaignSelectColumns+` FROM campaigns WHERE slug = $1`, strings.TrimSpace(slug))
}

// GetActiveByHost 按域名查询对外可见的活动。
func (r *campaignRepository) GetActiveByHost(ctx context.Context, host string) (*service.Campaign, error) {
	if strings.TrimSpace(host) == "" {
		return nil, nil
	}
	return r.queryOne(ctx,
		`SELECT`+campaignSelectColumns+` FROM campaigns
		 WHERE lower(host) = $1 AND status = 'active'
		 ORDER BY id DESC LIMIT 1`, strings.ToLower(strings.TrimSpace(host)))
}

// List 返回全部活动（管理端），按更新时间倒序。
func (r *campaignRepository) List(ctx context.Context) ([]service.Campaign, error) {
	rows, err := r.sql.QueryContext(ctx,
		`SELECT`+campaignSelectColumns+` FROM campaigns ORDER BY updated_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list campaigns: %w", err)
	}
	defer func() { _ = rows.Close() }()

	out := make([]service.Campaign, 0, 8)
	for rows.Next() {
		campaign, err := scanCampaign(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *campaign)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate campaigns: %w", err)
	}
	return out, nil
}

// Create 新建活动并回填自增字段。
func (r *campaignRepository) Create(ctx context.Context, campaign *service.Campaign) (*service.Campaign, error) {
	brand, hero, entries, samples, cta, trial, seo, translations, hidden, err := marshalCampaignJSON(campaign)
	if err != nil {
		return nil, err
	}
	rows, err := r.sql.QueryContext(ctx,
		`INSERT INTO campaigns
			(slug, name, status, host, group_id, locale, brand, hero, content_md, model_entries, code_samples, cta, trial, seo, translations, created_by)
		 VALUES ($1,$2,$3,NULLIF($4,''),$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		 RETURNING`+campaignSelectColumns,
		campaign.Slug, campaign.Name, string(campaign.Status), campaign.Host, campaign.GroupID, campaign.Locale,
		brand, hero, campaign.ContentMD, entries, samples, cta, trial, seo, translations, hidden, campaign.CreatedBy)
	if err != nil {
		return nil, translateCampaignWriteError(err)
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("insert campaign: %w", err)
		}
		return nil, fmt.Errorf("insert campaign: no row returned")
	}
	return scanCampaign(rows)
}

// Update 全量更新配置字段。
func (r *campaignRepository) Update(ctx context.Context, campaign *service.Campaign) (*service.Campaign, error) {
	brand, hero, entries, samples, cta, trial, seo, translations, hidden, err := marshalCampaignJSON(campaign)
	if err != nil {
		return nil, err
	}
	rows, err := r.sql.QueryContext(ctx,
		`UPDATE campaigns SET
			slug = $2, name = $3, status = $4, host = NULLIF($5,''), group_id = $6,
			locale = $7, brand = $8, hero = $9, content_md = $10, model_entries = $11,
			code_samples = $12, cta = $13, trial = $14, seo = $15, translations = $16,
			hidden_models = $17,
			updated_at = NOW()
		 WHERE id = $1
		 RETURNING`+campaignSelectColumns,
		campaign.ID, campaign.Slug, campaign.Name, string(campaign.Status), campaign.Host, campaign.GroupID,
		campaign.Locale, brand, hero, campaign.ContentMD, entries, samples, cta, trial, seo, translations, hidden)
	if err != nil {
		return nil, translateCampaignWriteError(err)
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("update campaign: %w", err)
		}
		return nil, service.ErrCampaignNotFound
	}
	return scanCampaign(rows)
}

// Delete 删除活动；campaign_key_grants 由外键级联清理。
func (r *campaignRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.sql.ExecContext(ctx, `DELETE FROM campaigns WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete campaign: %w", err)
	}
	if affected, err := res.RowsAffected(); err == nil && affected == 0 {
		return service.ErrCampaignNotFound
	}
	return nil
}

// GetKeyGrantAPIKeyID 查询该活动已发放给该用户的 Key。
func (r *campaignRepository) GetKeyGrantAPIKeyID(ctx context.Context, campaignID, userID int64) (int64, bool, error) {
	rows, err := r.sql.QueryContext(ctx,
		`SELECT api_key_id FROM campaign_key_grants WHERE campaign_id = $1 AND user_id = $2`,
		campaignID, userID)
	if err != nil {
		return 0, false, fmt.Errorf("get campaign key grant: %w", err)
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		return 0, false, nil
	}
	var apiKeyID int64
	if err := rows.Scan(&apiKeyID); err != nil {
		return 0, false, fmt.Errorf("scan campaign key grant: %w", err)
	}
	return apiKeyID, true, nil
}

// CreateKeyGrant 记录发放关系；并发冲突由唯一约束兜底。
func (r *campaignRepository) CreateKeyGrant(ctx context.Context, campaignID, userID, apiKeyID int64) error {
	_, err := r.sql.ExecContext(ctx,
		`INSERT INTO campaign_key_grants (campaign_id, user_id, api_key_id) VALUES ($1,$2,$3)`,
		campaignID, userID, apiKeyID)
	if err != nil {
		return fmt.Errorf("create campaign key grant: %w", err)
	}
	return nil
}

// AttributeSignup 写入注册归因；已有归因时保持不变（首站优先）。
func (r *campaignRepository) AttributeSignup(ctx context.Context, campaignID, userID int64) error {
	_, err := r.sql.ExecContext(ctx,
		`UPDATE users SET signup_campaign_id = $1
		 WHERE id = $2 AND signup_campaign_id IS NULL`,
		campaignID, userID)
	if err != nil {
		return fmt.Errorf("attribute campaign signup: %w", err)
	}
	return nil
}

func (r *campaignRepository) queryOne(ctx context.Context, query string, args ...any) (*service.Campaign, error) {
	rows, err := r.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query campaign: %w", err)
	}
	defer func() { _ = rows.Close() }()
	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("query campaign: %w", err)
		}
		return nil, nil
	}
	return scanCampaign(rows)
}

func scanCampaign(rows *sql.Rows) (*service.Campaign, error) {
	var (
		campaign  service.Campaign
		status    string
		groupID   sql.NullInt64
		createdBy sql.NullInt64
		brand     []byte
		hero      []byte
		entries   []byte
		samples   []byte
		cta       []byte
		trial     []byte
		seo       []byte
		trans     []byte
		hidden    []byte
	)
	if err := rows.Scan(
		&campaign.ID, &campaign.Slug, &campaign.Name, &status, &campaign.Host, &groupID, &campaign.Locale,
		&brand, &hero, &campaign.ContentMD, &entries, &samples, &cta, &trial, &seo, &trans, &hidden,
		&createdBy, &campaign.CreatedAt, &campaign.UpdatedAt,
	); err != nil {
		return nil, fmt.Errorf("scan campaign: %w", err)
	}
	campaign.Status = service.CampaignStatus(status)
	if groupID.Valid {
		id := groupID.Int64
		campaign.GroupID = &id
	}
	if createdBy.Valid {
		id := createdBy.Int64
		campaign.CreatedBy = &id
	}
	decodeCampaignJSON(brand, &campaign.Brand)
	decodeCampaignJSON(hero, &campaign.Hero)
	decodeCampaignJSON(entries, &campaign.ModelEntries)
	decodeCampaignJSON(samples, &campaign.CodeSamples)
	decodeCampaignJSON(cta, &campaign.CTA)
	decodeCampaignJSON(trial, &campaign.Trial)
	decodeCampaignJSON(seo, &campaign.SEO)
	decodeCampaignJSON(trans, &campaign.Translations)
	decodeCampaignJSON(hidden, &campaign.HiddenModels)
	return &campaign, nil
}

// decodeCampaignJSON 对配置型 JSON 采取「宽容解析」：脏数据不应让整个落地页 500。
func decodeCampaignJSON(raw []byte, target any) {
	if len(raw) == 0 {
		return
	}
	_ = json.Unmarshal(raw, target)
}

func marshalCampaignJSON(campaign *service.Campaign) (brand, hero, entries, samples, cta, trial, seo, translations, hidden []byte, err error) {
	if brand, err = json.Marshal(campaign.Brand); err != nil {
		return
	}
	if hero, err = json.Marshal(campaign.Hero); err != nil {
		return
	}
	if entries, err = json.Marshal(nonNilSlice(campaign.ModelEntries)); err != nil {
		return
	}
	if samples, err = json.Marshal(nonNilSlice(campaign.CodeSamples)); err != nil {
		return
	}
	if cta, err = json.Marshal(campaign.CTA); err != nil {
		return
	}
	if trial, err = json.Marshal(campaign.Trial); err != nil {
		return
	}
	if seo, err = json.Marshal(campaign.SEO); err != nil {
		return
	}
	translations, err = json.Marshal(nonNilTranslations(campaign.Translations))
	if err != nil {
		return
	}
	if hidden, err = json.Marshal(nonNilSlice(campaign.HiddenModels)); err != nil {
		return
	}
	return brand, hero, entries, samples, cta, trial, seo, translations, hidden, nil
}

// nonNilTranslations 保证写入 jsonb 的是 {} 而不是 null。
func nonNilTranslations(value map[string]service.CampaignLocaleContent) map[string]service.CampaignLocaleContent {
	if value == nil {
		return map[string]service.CampaignLocaleContent{}
	}
	return value
}

// nonNilSlice 保证写入 jsonb 的是 [] 而不是 null——前端与 SQL 默认值都按数组处理。
func nonNilSlice[T any](value []T) []T {
	if value == nil {
		return []T{}
	}
	return value
}

// translateCampaignWriteError 把唯一约束冲突翻译成业务错误，避免把 pq 原始错误透给前端。
func translateCampaignWriteError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && string(pqErr.Code) == "23505" {
		switch {
		case strings.Contains(pqErr.Constraint, "host"):
			return service.ErrCampaignHostTaken
		case strings.Contains(pqErr.Constraint, "slug"):
			return service.ErrCampaignSlugTaken
		}
	}
	return fmt.Errorf("write campaign: %w", err)
}
