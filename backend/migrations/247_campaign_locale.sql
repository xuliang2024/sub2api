-- 247_campaign_locale.sql
--
-- 落地页界面语言。
--
-- 默认 en：这套落地页的主要场景是国产模型卖给海外用户，英文是默认面向。
-- 需要中文投放时在后台把 locale 改成 zh 即可（界面文案内置 en/zh 两套）。
--
-- 幂等：IF NOT EXISTS。

ALTER TABLE campaigns
    ADD COLUMN IF NOT EXISTS locale VARCHAR(8) NOT NULL DEFAULT 'en';

COMMENT ON COLUMN campaigns.locale IS '界面语言：en（默认，海外投放）/ zh';
