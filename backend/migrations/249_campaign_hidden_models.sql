-- 249_campaign_hidden_models.sql
--
-- 落地页隐藏模型。
--
-- 背景：渠道定价里可能存在「历史别名」——例如 deepseek-v4.1-flash 的同价别名
-- deepseek-flash / deepseek-v4-flash / deepseek-v4-1-flash-260910。它们用于承接旧请求的
-- 计费，但对外宣传时不应该出现（后台渠道页也会合并显示）。
--
-- 这些别名属于部署事实，代码里硬编码去重会绑死某个厂商，因此做成活动级配置：
-- 运营填要隐藏的模型名，落地页的卡片与模型广场都会过滤掉。
--
-- 幂等：IF NOT EXISTS。

ALTER TABLE campaigns
    ADD COLUMN IF NOT EXISTS hidden_models JSONB NOT NULL DEFAULT '[]'::jsonb;

COMMENT ON COLUMN campaigns.hidden_models IS
    '不在落地页展示的模型名数组（别名/退役模型），卡片与广场都会过滤';
