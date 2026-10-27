package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

type partnerRepository struct {
	db  *sql.DB
	sql sqlExecutor
}

// NewPartnerRepository 创建合伙人仓储。
func NewPartnerRepository(db *sql.DB) service.PartnerRepository {
	return &partnerRepository{db: db, sql: db}
}

const partnerColumns = `user_id, min_rate_multiplier, status, settled_total, created_at, updated_at`

func (r *partnerRepository) GetPartnerByUserID(ctx context.Context, userID int64) (*service.Partner, error) {
	var p service.Partner
	err := scanSingleRow(ctx, r.sql, `SELECT `+partnerColumns+` FROM partners WHERE user_id = $1`, []any{userID},
		&p.UserID, &p.MinRateMultiplier, &p.Status, &p.SettledTotal, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *partnerRepository) UpsertPartnerWithClamp(ctx context.Context, p *service.Partner, actorID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO partners (user_id, min_rate_multiplier, status, created_at, updated_at)
		VALUES ($1, $2, $3, now(), now())
		ON CONFLICT (user_id) DO UPDATE SET
			min_rate_multiplier = EXCLUDED.min_rate_multiplier,
			status = EXCLUDED.status,
			updated_at = now()
	`, p.UserID, p.MinRateMultiplier, p.Status); err != nil {
		return err
	}

	// 上调 min_rate 时，把该合伙人设置的低于新下限的下线倍率抬升到新下限，
	// 否则用户消费价低于平台保底收款价，平台倒贴。
	// 注意：参数必须从 $1 起连续编号，PG 不接受跳号（$1 缺失会报
	// "could not determine data type of parameter $1"）。
	if _, err := tx.ExecContext(ctx, `
		UPDATE user_group_rate_multipliers
		SET rate_multiplier = $1, updated_at = now()
		WHERE set_by = $2 AND rate_multiplier IS NOT NULL AND rate_multiplier < $1
	`, p.MinRateMultiplier, service.SetByPartner(p.UserID)); err != nil {
		return err
	}

	if err := writePartnerAudit(ctx, tx, actorID, "upsert_partner", &p.UserID, map[string]any{
		"min_rate_multiplier": p.MinRateMultiplier,
		"status":              p.Status,
	}); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *partnerRepository) DeletePartnerAndUnbind(ctx context.Context, partnerUserID, actorID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var downlineIDs []int64
	rows, err := tx.QueryContext(ctx, `SELECT user_id FROM partner_relations WHERE partner_user_id = $1`, partnerUserID)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			_ = rows.Close()
			return err
		}
		downlineIDs = append(downlineIDs, id)
	}
	_ = rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	if err := clearPartnerRates(ctx, tx, service.SetByPartner(partnerUserID), downlineIDs); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM partner_relations WHERE partner_user_id = $1`, partnerUserID); err != nil {
		return err
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM partners WHERE user_id = $1`, partnerUserID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return service.ErrPartnerNotFound
	}
	if err := writePartnerAudit(ctx, tx, actorID, "delete_partner", &partnerUserID, map[string]any{
		"unbound_user_ids": downlineIDs,
	}); err != nil {
		return err
	}
	return tx.Commit()
}

// clearPartnerRates 清除某合伙人给下线设置的倍率（置 NULL 并归还 admin 来源），
// 随后清理这些用户名下 rate/rpm 全空的行。downlineIDs 为空时仅按 set_by 清理。
func clearPartnerRates(ctx context.Context, tx *sql.Tx, setBy string, downlineIDs []int64) error {
	if _, err := tx.ExecContext(ctx, `
		UPDATE user_group_rate_multipliers
		SET rate_multiplier = NULL, set_by = 'admin', updated_at = now()
		WHERE set_by = $1
	`, setBy); err != nil {
		return err
	}
	if len(downlineIDs) == 0 {
		return nil
	}
	_, err := tx.ExecContext(ctx, `
		DELETE FROM user_group_rate_multipliers
		WHERE user_id = ANY($1) AND rate_multiplier IS NULL AND rpm_override IS NULL
	`, pq.Array(downlineIDs))
	return err
}

func (r *partnerRepository) ListPartners(ctx context.Context, query, status string, page, pageSize int) ([]service.PartnerListItem, int, error) {
	where := "WHERE u.deleted_at IS NULL"
	args := make([]any, 0, 4)
	if q := strings.TrimSpace(query); q != "" {
		args = append(args, "%"+q+"%")
		where += fmt.Sprintf(" AND (u.email ILIKE $%d OR COALESCE(u.username,'') ILIKE $%d)", len(args), len(args))
	}
	if s := strings.TrimSpace(status); s != "" {
		args = append(args, s)
		where += fmt.Sprintf(" AND p.status = $%d", len(args))
	}

	var total int
	if err := scanSingleRow(ctx, r.sql, `SELECT COUNT(*) FROM partners p JOIN users u ON u.id = p.user_id `+where, args, &total); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	rows, err := r.sql.QueryContext(ctx, `
		SELECT p.user_id, u.email, COALESCE(u.username,''), p.min_rate_multiplier, p.status, p.settled_total,
			p.created_at, p.updated_at,
			(SELECT COUNT(*) FROM partner_relations pr WHERE pr.partner_user_id = p.user_id) AS downline_count,
			COALESCE((SELECT SUM(ul.partner_commission) FROM usage_logs ul WHERE ul.partner_id = p.user_id), 0) AS total_commission
		FROM partners p
		JOIN users u ON u.id = p.user_id
		`+where+`
		ORDER BY p.created_at DESC, p.user_id DESC
		LIMIT $`+fmt.Sprint(len(listArgs)-1)+` OFFSET $`+fmt.Sprint(len(listArgs)), listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.PartnerListItem, 0)
	for rows.Next() {
		var it service.PartnerListItem
		if err := rows.Scan(&it.UserID, &it.Email, &it.Username, &it.MinRateMultiplier, &it.Status, &it.SettledTotal,
			&it.CreatedAt, &it.UpdatedAt, &it.DownlineCount, &it.TotalCommission); err != nil {
			return nil, 0, err
		}
		it.PendingSettlement = it.TotalCommission - it.SettledTotal
		if it.PendingSettlement < 0 {
			it.PendingSettlement = 0
		}
		items = append(items, it)
	}
	return items, total, rows.Err()
}

func (r *partnerRepository) GetPartnerUserIDs(ctx context.Context, userIDs []int64) (map[int64]bool, error) {
	result := make(map[int64]bool, len(userIDs))
	if len(userIDs) == 0 {
		return result, nil
	}
	rows, err := r.sql.QueryContext(ctx, `SELECT user_id FROM partners WHERE user_id = ANY($1)`, pq.Array(userIDs))
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		result[id] = true
	}
	return result, rows.Err()
}

func (r *partnerRepository) GetRelationByUserID(ctx context.Context, userID int64) (*service.PartnerRelation, error) {
	var rel service.PartnerRelation
	err := scanSingleRow(ctx, r.sql, `
		SELECT id, partner_user_id, user_id, bound_source, created_at
		FROM partner_relations WHERE user_id = $1
	`, []any{userID}, &rel.ID, &rel.PartnerUserID, &rel.UserID, &rel.BoundSource, &rel.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &rel, nil
}

func (r *partnerRepository) BindRelation(ctx context.Context, partnerUserID, userID int64, source string, actorID int64) error {
	if partnerUserID <= 0 || userID <= 0 || partnerUserID == userID {
		return service.ErrPartnerNotDownline
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var isPartner int
	err = tx.QueryRowContext(ctx, `SELECT 1 FROM partners WHERE user_id = $1`, userID).Scan(&isPartner)
	if err == nil {
		return service.ErrPartnerDownlineIsPartner
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	var oldPartnerID int64
	var oldSetBy string
	err = tx.QueryRowContext(ctx, `SELECT partner_user_id FROM partner_relations WHERE user_id = $1 FOR UPDATE`, userID).Scan(&oldPartnerID)
	switch {
	case err == nil && oldPartnerID == partnerUserID:
		// 幂等：已绑定同一合伙人，仅更新来源。
		if _, err := tx.ExecContext(ctx, `UPDATE partner_relations SET bound_source = $3 WHERE partner_user_id = $1 AND user_id = $2`, partnerUserID, userID, source); err != nil {
			return err
		}
		return tx.Commit()
	case err == nil:
		oldSetBy = service.SetByPartner(oldPartnerID)
	case !errors.Is(err, sql.ErrNoRows):
		return err
	}

	if oldSetBy != "" {
		if err := clearPartnerRates(ctx, tx, oldSetBy, []int64{userID}); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM partner_relations WHERE user_id = $1`, userID); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO partner_relations (partner_user_id, user_id, bound_source, created_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (user_id) DO UPDATE SET partner_user_id = EXCLUDED.partner_user_id, bound_source = EXCLUDED.bound_source
	`, partnerUserID, userID, source); err != nil {
		return err
	}
	if err := writePartnerAudit(ctx, tx, actorID, "bind", &userID, map[string]any{
		"from_partner_id": oldPartnerID,
		"to_partner_id":   partnerUserID,
		"source":          source,
	}); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *partnerRepository) UnbindRelation(ctx context.Context, partnerUserID, userID, actorID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	res, err := tx.ExecContext(ctx, `DELETE FROM partner_relations WHERE partner_user_id = $1 AND user_id = $2`, partnerUserID, userID)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return service.ErrPartnerNotDownline
	}
	if err := clearPartnerRates(ctx, tx, service.SetByPartner(partnerUserID), []int64{userID}); err != nil {
		return err
	}
	if err := writePartnerAudit(ctx, tx, actorID, "unbind", &userID, map[string]any{"partner_user_id": partnerUserID}); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *partnerRepository) ListDownlines(ctx context.Context, partnerUserID int64) ([]service.PartnerDownline, error) {
	rows, err := r.sql.QueryContext(ctx, `
		SELECT pr.user_id, u.email, COALESCE(u.username,''), pr.bound_source, pr.created_at
		FROM partner_relations pr
		JOIN users u ON u.id = pr.user_id AND u.deleted_at IS NULL
		WHERE pr.partner_user_id = $1
		ORDER BY pr.created_at DESC, pr.user_id DESC
	`, partnerUserID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.PartnerDownline, 0)
	for rows.Next() {
		var d service.PartnerDownline
		if err := rows.Scan(&d.UserID, &d.Email, &d.Username, &d.BoundSource, &d.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, d)
	}
	return items, rows.Err()
}

func (r *partnerRepository) SetDownlineRate(ctx context.Context, partnerUserID, userID, groupID int64, rate float64) error {
	_, err := r.sql.ExecContext(ctx, `
		INSERT INTO user_group_rate_multipliers (user_id, group_id, rate_multiplier, set_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, now(), now())
		ON CONFLICT (user_id, group_id) DO UPDATE SET
			rate_multiplier = EXCLUDED.rate_multiplier,
			set_by = EXCLUDED.set_by,
			updated_at = now()
	`, userID, groupID, rate, service.SetByPartner(partnerUserID))
	return err
}

func (r *partnerRepository) ClearDownlineRate(ctx context.Context, partnerUserID, userID, groupID int64) error {
	if _, err := r.sql.ExecContext(ctx, `
		UPDATE user_group_rate_multipliers
		SET rate_multiplier = NULL, set_by = 'admin', updated_at = now()
		WHERE user_id = $1 AND group_id = $2 AND set_by = $3
	`, userID, groupID, service.SetByPartner(partnerUserID)); err != nil {
		return err
	}
	_, err := r.sql.ExecContext(ctx, `
		DELETE FROM user_group_rate_multipliers
		WHERE user_id = $1 AND group_id = $2 AND rate_multiplier IS NULL AND rpm_override IS NULL
	`, userID, groupID)
	return err
}

func (r *partnerRepository) ResolveUsageShare(ctx context.Context, userID, groupID int64) (*service.PartnerUsageShare, error) {
	var share service.PartnerUsageShare
	err := scanSingleRow(ctx, r.sql, `
		SELECT p.user_id, p.min_rate_multiplier, ugr.rate_multiplier
		FROM partner_relations pr
		JOIN partners p ON p.user_id = pr.partner_user_id AND p.status = 'active'
		JOIN user_group_rate_multipliers ugr
			ON ugr.user_id = pr.user_id AND ugr.group_id = $2
			AND ugr.set_by = 'partner:' || p.user_id::text
		WHERE pr.user_id = $1 AND ugr.rate_multiplier IS NOT NULL
	`, []any{userID, groupID}, &share.PartnerUserID, &share.MinRate, &share.BaseRate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &share, nil
}

func (r *partnerRepository) GetCommissionTotals(ctx context.Context, partnerUserID int64) (float64, float64, error) {
	var total, month float64
	err := scanSingleRow(ctx, r.sql, `
		SELECT COALESCE(SUM(partner_commission), 0),
			COALESCE(SUM(partner_commission) FILTER (WHERE created_at >= date_trunc('month', now())), 0)
		FROM usage_logs WHERE partner_id = $1
	`, []any{partnerUserID}, &total, &month)
	return total, month, err
}

func (r *partnerRepository) GetDownlineStats(ctx context.Context, partnerUserID int64, since time.Time) (map[int64]service.PartnerDownlineStats, error) {
	rows, err := r.sql.QueryContext(ctx, `
		SELECT user_id, COALESCE(SUM(actual_cost), 0), COALESCE(SUM(partner_commission), 0)
		FROM usage_logs
		WHERE partner_id = $1 AND created_at >= $2
		GROUP BY user_id
	`, partnerUserID, since)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := make(map[int64]service.PartnerDownlineStats)
	for rows.Next() {
		var userID int64
		var st service.PartnerDownlineStats
		if err := rows.Scan(&userID, &st.ActualCost, &st.Commission); err != nil {
			return nil, err
		}
		result[userID] = st
	}
	return result, rows.Err()
}

func (r *partnerRepository) ListCommissionRecords(ctx context.Context, partnerUserID int64, filter service.PartnerCommissionFilter) ([]service.PartnerCommissionRecord, int, float64, error) {
	where := "WHERE ul.partner_id = $1"
	args := []any{partnerUserID}
	if filter.UserID > 0 {
		args = append(args, filter.UserID)
		where += fmt.Sprintf(" AND ul.user_id = $%d", len(args))
	}
	if filter.StartTime != nil {
		args = append(args, *filter.StartTime)
		where += fmt.Sprintf(" AND ul.created_at >= $%d", len(args))
	}
	if filter.EndTime != nil {
		args = append(args, *filter.EndTime)
		where += fmt.Sprintf(" AND ul.created_at < $%d", len(args))
	}

	var total int
	var totalCommission float64
	if err := scanSingleRow(ctx, r.sql,
		`SELECT COUNT(*), COALESCE(SUM(ul.partner_commission), 0) FROM usage_logs ul `+where, args, &total, &totalCommission,
	); err != nil {
		return nil, 0, 0, err
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 || pageSize > 200 {
		pageSize = 20
	}
	listArgs := append(append([]any{}, args...), pageSize, (page-1)*pageSize)
	rows, err := r.sql.QueryContext(ctx, `
		SELECT ul.id, ul.created_at, ul.user_id, COALESCE(u.email, ''), ul.model, ul.total_cost,
			ul.partner_base_rate, ul.partner_rate_floor, ul.rate_multiplier, ul.partner_commission
		FROM usage_logs ul
		LEFT JOIN users u ON u.id = ul.user_id
		`+where+`
		ORDER BY ul.created_at DESC, ul.id DESC
		LIMIT $`+fmt.Sprint(len(listArgs)-1)+` OFFSET $`+fmt.Sprint(len(listArgs)), listArgs...)
	if err != nil {
		return nil, 0, 0, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.PartnerCommissionRecord, 0)
	for rows.Next() {
		var rec service.PartnerCommissionRecord
		if err := rows.Scan(&rec.ID, &rec.CreatedAt, &rec.UserID, &rec.UserEmail, &rec.Model, &rec.TotalCost,
			&rec.BaseRate, &rec.RateFloor, &rec.RateMultiplier, &rec.Commission); err != nil {
			return nil, 0, 0, err
		}
		items = append(items, rec)
	}
	return items, total, totalCommission, rows.Err()
}

func (r *partnerRepository) CreateSettlement(ctx context.Context, s *service.PartnerSettlement) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var settledTotal float64
	err = tx.QueryRowContext(ctx, `SELECT settled_total FROM partners WHERE user_id = $1 FOR UPDATE`, s.PartnerUserID).Scan(&settledTotal)
	if errors.Is(err, sql.ErrNoRows) {
		return service.ErrPartnerNotFound
	}
	if err != nil {
		return err
	}
	var totalCommission float64
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(SUM(partner_commission), 0) FROM usage_logs WHERE partner_id = $1`, s.PartnerUserID).Scan(&totalCommission); err != nil {
		return err
	}
	if s.Amount <= 0 || s.Amount > totalCommission-settledTotal+1e-9 {
		return service.ErrPartnerInvalidSettle
	}

	if err := tx.QueryRowContext(ctx, `
		INSERT INTO partner_settlements (partner_user_id, amount, method, note, operator_id, created_at)
		VALUES ($1, $2, $3, $4, $5, now())
		RETURNING id, created_at
	`, s.PartnerUserID, s.Amount, s.Method, s.Note, s.OperatorID).Scan(&s.ID, &s.CreatedAt); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE partners SET settled_total = settled_total + $2, updated_at = now() WHERE user_id = $1`, s.PartnerUserID, s.Amount); err != nil {
		return err
	}
	if err := writePartnerAudit(ctx, tx, s.OperatorID, "settle", &s.PartnerUserID, map[string]any{
		"amount": s.Amount,
		"method": s.Method,
		"note":   s.Note,
	}); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *partnerRepository) ListSettlements(ctx context.Context, partnerUserID int64, limit int) ([]service.PartnerSettlement, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := r.sql.QueryContext(ctx, `
		SELECT id, partner_user_id, amount, method, COALESCE(note, ''), COALESCE(operator_id, 0), created_at
		FROM partner_settlements WHERE partner_user_id = $1
		ORDER BY created_at DESC, id DESC LIMIT $2
	`, partnerUserID, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]service.PartnerSettlement, 0)
	for rows.Next() {
		var s service.PartnerSettlement
		if err := rows.Scan(&s.ID, &s.PartnerUserID, &s.Amount, &s.Method, &s.Note, &s.OperatorID, &s.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, s)
	}
	return items, rows.Err()
}

func (r *partnerRepository) WriteAudit(ctx context.Context, actorID int64, action string, targetUserID *int64, detail map[string]any) error {
	return writePartnerAudit(ctx, r.sql, actorID, action, targetUserID, detail)
}

func writePartnerAudit(ctx context.Context, exec sqlExecutor, actorID int64, action string, targetUserID *int64, detail map[string]any) error {
	var payload any
	if detail != nil {
		raw, err := json.Marshal(detail)
		if err != nil {
			return err
		}
		payload = raw
	}
	var target any
	if targetUserID != nil {
		target = *targetUserID
	}
	_, err := exec.ExecContext(ctx, `
		INSERT INTO partner_audit_logs (actor_id, action, target_user_id, detail, created_at)
		VALUES ($1, $2, $3, $4, now())
	`, actorID, action, target, payload)
	return err
}
