-- 246_campaigns.sql
--
-- 模型广告投放落地页（campaign）配置。
--
-- 背景：站点原有 settings 是单站点平铺 KV（无 site_id / host 维度），因此无法承载
-- "一个域名一套落地页"的多活动需求。这里单独建表，按 host 解析，互不干扰。
--
-- 约定：
--   * slug   —— 后台/预览用的稳定标识，同时作为「同域名多活动」时的 ?c= 覆盖参数。
--   * host   —— 绑定的独立域名（不含端口、小写）。为空表示仅通过 slug 访问。
--   * group_id —— 决定该落地页可用的模型、价格与倍率；NULL 表示仅做展示不建 Key。
--   * 结构化配置统一用 JSONB 存储，避免为每个营销区块加列；服务层负责解析与校验。
--
-- 幂等：全部 IF NOT EXISTS，可重复执行。

CREATE TABLE IF NOT EXISTS campaigns (
    id            BIGSERIAL PRIMARY KEY,
    slug          VARCHAR(64)  NOT NULL,
    name          VARCHAR(128) NOT NULL,
    -- draft: 仅管理员可预览；active: 对外；archived: 下线但保留配置
    status        VARCHAR(16)  NOT NULL DEFAULT 'draft',
    host          VARCHAR(255),
    group_id      BIGINT REFERENCES groups(id) ON DELETE SET NULL,

    -- brand: 站点名/Logo/主题色/页脚等外观
    brand         JSONB NOT NULL DEFAULT '{}'::jsonb,
    -- hero: 首屏文案与主视觉
    hero          JSONB NOT NULL DEFAULT '{}'::jsonb,
    -- content_md: 使用说明（Markdown 原文，前端 marked + DOMPurify 渲染）
    content_md    TEXT  NOT NULL DEFAULT '',
    -- model_entries: 主推模型清单 [{model, display_name, tagline, badge, featured, order}]
    model_entries JSONB NOT NULL DEFAULT '[]'::jsonb,
    -- code_samples: 接入示例 [{id, label, lang, code}]，支持 {{BASE_URL}}/{{API_KEY}}/{{MODEL}} 占位
    code_samples  JSONB NOT NULL DEFAULT '[]'::jsonb,
    -- cta: 转化动作 {mode, label, sublabel, topup_url}
    cta           JSONB NOT NULL DEFAULT '{}'::jsonb,
    -- trial: 注册赠送 {balance, note}
    trial         JSONB NOT NULL DEFAULT '{}'::jsonb,
    -- seo: {title, description, og_image, keywords}
    seo           JSONB NOT NULL DEFAULT '{}'::jsonb,

    created_by    BIGINT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- slug 全局唯一（同一活动可能既绑域名又能用 slug 直达）
CREATE UNIQUE INDEX IF NOT EXISTS idx_campaigns_slug
    ON campaigns (slug);

-- 一个域名只允许绑定一个 active 活动；draft/archived 不占用域名。
CREATE UNIQUE INDEX IF NOT EXISTS idx_campaigns_active_host
    ON campaigns (lower(host))
    WHERE host IS NOT NULL AND host <> '' AND status = 'active';

CREATE INDEX IF NOT EXISTS idx_campaigns_status ON campaigns (status);
CREATE INDEX IF NOT EXISTS idx_campaigns_group  ON campaigns (group_id);

COMMENT ON TABLE  campaigns IS '模型广告投放落地页配置（按 host 或 slug 解析）';
COMMENT ON COLUMN campaigns.host IS '绑定的独立域名（小写、不含端口）；空表示仅 slug 访问';
COMMENT ON COLUMN campaigns.model_entries IS '主推模型清单 JSON 数组';
COMMENT ON COLUMN campaigns.trial IS '注册赠送额度 {balance, note}';

-- 注册归因：记录用户来自哪个 campaign，用于投放 ROI 归集。
ALTER TABLE users ADD COLUMN IF NOT EXISTS signup_campaign_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_users_signup_campaign
    ON users (signup_campaign_id)
    WHERE signup_campaign_id IS NOT NULL;

COMMENT ON COLUMN users.signup_campaign_id IS '注册来源的广告活动 ID（campaigns.id），仅首登写入';

-- 一键建 Key 的幂等凭据：同一 campaign + 同一用户只发放一次，重复点击返回同一把 Key。
CREATE TABLE IF NOT EXISTS campaign_key_grants (
    id          BIGSERIAL PRIMARY KEY,
    campaign_id BIGINT NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key_id  BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_campaign_key_grants UNIQUE (campaign_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_campaign_key_grants_user ON campaign_key_grants (user_id);

COMMENT ON TABLE campaign_key_grants IS '落地页一键建 Key 的幂等记录（每活动每用户一把）';
