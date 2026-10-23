-- 248_campaign_translations.sql
--
-- 落地页正文的多语言覆盖。
--
-- 为什么需要：界面文案（按钮/表头/错误提示）由前端内置词表提供，但标题、副标题、
-- 使用说明属于运营内容，只有活动配置里才有。若只有 UI 翻译，切到泰语会出现
-- 「按钮是泰语、标题还是英语」的割裂感，投放页不可用。
--
-- 结构：{"ja": {"hero_title": "...", "hero_subtitle": "...", "content_md": "..."}, ...}
-- 未配置的语言回落活动的主文案（campaigns.hero / campaigns.content_md）。
--
-- 幂等：IF NOT EXISTS。

ALTER TABLE campaigns
    ADD COLUMN IF NOT EXISTS translations JSONB NOT NULL DEFAULT '{}'::jsonb;

COMMENT ON COLUMN campaigns.translations IS
    '按语言的正文覆盖：{lang: {hero_title, hero_subtitle, content_md}}';
