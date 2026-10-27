CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_usage_logs_partner_created
    ON usage_logs (partner_id, created_at DESC, id DESC);
