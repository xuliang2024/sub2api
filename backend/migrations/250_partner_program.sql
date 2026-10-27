-- 合伙人（Partner）制度：合伙人、绑定关系、结算台账、操作审计、usage 收益快照。

CREATE TABLE IF NOT EXISTS partners (
    id                  BIGSERIAL PRIMARY KEY,
    user_id             BIGINT NOT NULL UNIQUE REFERENCES users(id),
    min_rate_multiplier DECIMAL(10,4) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    settled_total       DECIMAL(20,8) NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT partners_min_rate_range CHECK (min_rate_multiplier > 0 AND min_rate_multiplier <= 1),
    CONSTRAINT partners_status_check CHECK (status IN ('active', 'disabled'))
);
CREATE INDEX IF NOT EXISTS idx_partners_user_id ON partners(user_id);

CREATE TABLE IF NOT EXISTS partner_relations (
    id              BIGSERIAL PRIMARY KEY,
    partner_user_id BIGINT NOT NULL REFERENCES users(id),
    user_id         BIGINT NOT NULL UNIQUE REFERENCES users(id),
    bound_source    VARCHAR(20) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT partner_relations_source_check CHECK (bound_source IN ('admin', 'invite')),
    CONSTRAINT partner_relations_no_self CHECK (partner_user_id <> user_id),
    UNIQUE (partner_user_id, user_id)
);
CREATE INDEX IF NOT EXISTS idx_partner_relations_partner ON partner_relations(partner_user_id);

CREATE TABLE IF NOT EXISTS partner_settlements (
    id              BIGSERIAL PRIMARY KEY,
    partner_user_id BIGINT NOT NULL REFERENCES users(id),
    amount          DECIMAL(20,8) NOT NULL,
    method          VARCHAR(20) NOT NULL DEFAULT 'offline',
    note            TEXT,
    operator_id     BIGINT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT partner_settlements_amount_positive CHECK (amount > 0)
);
CREATE INDEX IF NOT EXISTS idx_partner_settlements_partner ON partner_settlements(partner_user_id, created_at);

CREATE TABLE IF NOT EXISTS partner_audit_logs (
    id              BIGSERIAL PRIMARY KEY,
    actor_id        BIGINT NOT NULL,
    action          VARCHAR(40) NOT NULL,
    target_user_id  BIGINT,
    detail          JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_partner_audit_actor ON partner_audit_logs(actor_id, created_at);
CREATE INDEX IF NOT EXISTS idx_partner_audit_target ON partner_audit_logs(target_user_id, created_at);

-- usage 收益快照：命中合伙人定价的消费记录 partner_id；commission 可为 0（下线消费统计仍需 partner_id）。
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS partner_id BIGINT;
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS partner_commission DECIMAL(20,8) NOT NULL DEFAULT 0;
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS partner_base_rate DECIMAL(10,4);
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS partner_rate_floor DECIMAL(10,4);

-- 用户级倍率来源标记：admin（默认）/ partner:<user_id>。
ALTER TABLE user_group_rate_multipliers ADD COLUMN IF NOT EXISTS set_by VARCHAR(40) NOT NULL DEFAULT 'admin';
