-- 费率类型动态化迁移
-- 创建时间: 2026-01-25
-- 业务规则:
--   通道给什么费率类型，我们就存什么费率类型，原封不动
--   费率范围约束：通道底线(min_rate) ≤ 政策模版成本 ≤ 代理商结算价 ≤ 商户费率 ≤ 通道上限(max_rate)

-- ============================================================
-- 1. 为恒信通配置费率类型（在channels.config中添加rate_types）
-- ============================================================
UPDATE channels SET config = jsonb_set(
    COALESCE(config::jsonb, '{}'::jsonb),
    '{rate_types}',
    '[
        {"code": "CREDIT", "name": "贷记卡费率", "sort_order": 1, "min_rate": "0.50", "max_rate": "0.68"},
        {"code": "DEBIT", "name": "借记卡费率", "sort_order": 2, "min_rate": "0.45", "max_rate": "0.60"},
        {"code": "WECHAT", "name": "微信费率", "sort_order": 3, "min_rate": "0.30", "max_rate": "0.60"},
        {"code": "ALIPAY", "name": "支付宝费率", "sort_order": 4, "min_rate": "0.30", "max_rate": "0.60"},
        {"code": "UNIONPAY", "name": "云闪付费率", "sort_order": 5, "min_rate": "0.30", "max_rate": "0.60"}
    ]'::jsonb
) WHERE channel_code = 'HENGXINTONG';

-- ============================================================
-- 2. 添加 rate_configs 字段到 policy_templates 表
-- ============================================================
ALTER TABLE policy_templates
    ADD COLUMN IF NOT EXISTS rate_configs JSONB DEFAULT '{}';

COMMENT ON COLUMN policy_templates.rate_configs IS '费率配置JSON，key为费率类型code，value为{rate: "0.55"}';

-- ============================================================
-- 3. 迁移现有政策模版数据到 rate_configs
-- ============================================================
UPDATE policy_templates SET rate_configs = jsonb_build_object(
    'CREDIT', jsonb_build_object('rate', COALESCE(credit_rate::text, '0')),
    'DEBIT', jsonb_build_object('rate', COALESCE(debit_rate::text, '0')),
    'WECHAT', jsonb_build_object('rate', COALESCE(wechat_rate::text, '0')),
    'ALIPAY', jsonb_build_object('rate', COALESCE(alipay_rate::text, '0')),
    'UNIONPAY', jsonb_build_object('rate', COALESCE(unionpay_rate::text, '0'))
) WHERE rate_configs = '{}' OR rate_configs IS NULL;

-- ============================================================
-- 4. 添加 rate_configs 字段到 agent_policies 表
-- ============================================================
ALTER TABLE agent_policies
    ADD COLUMN IF NOT EXISTS rate_configs JSONB DEFAULT '{}';

COMMENT ON COLUMN agent_policies.rate_configs IS '费率配置JSON，key为费率类型code';

-- ============================================================
-- 5. 迁移现有代理商政策数据到 rate_configs
-- ============================================================
UPDATE agent_policies SET rate_configs = jsonb_build_object(
    'CREDIT', jsonb_build_object('rate', COALESCE(credit_rate::text, '0')),
    'DEBIT', jsonb_build_object('rate', COALESCE(debit_rate::text, '0')),
    'WECHAT', jsonb_build_object('rate', COALESCE(wechat_rate::text, '0')),
    'ALIPAY', jsonb_build_object('rate', COALESCE(alipay_rate::text, '0')),
    'UNIONPAY', jsonb_build_object('rate', COALESCE(unionpay_rate::text, '0'))
) WHERE rate_configs = '{}' OR rate_configs IS NULL;

-- ============================================================
-- 6. 添加 rate_deltas 字段到 rate_stage_policies 表
-- ============================================================
ALTER TABLE rate_stage_policies
    ADD COLUMN IF NOT EXISTS rate_deltas JSONB DEFAULT '{}';

COMMENT ON COLUMN rate_stage_policies.rate_deltas IS '费率阶梯调整值JSON，key为费率类型code，value为调整值字符串';

-- ============================================================
-- 7. 迁移现有费率阶梯数据到 rate_deltas
-- ============================================================
UPDATE rate_stage_policies SET rate_deltas = jsonb_build_object(
    'CREDIT', COALESCE(credit_rate_delta::text, '0'),
    'DEBIT', COALESCE(debit_rate_delta::text, '0'),
    'WECHAT', COALESCE(wechat_rate_delta::text, '0'),
    'ALIPAY', COALESCE(alipay_rate_delta::text, '0'),
    'UNIONPAY', COALESCE(unionpay_rate_delta::text, '0')
) WHERE rate_deltas = '{}' OR rate_deltas IS NULL;

-- ============================================================
-- 8. 为联动通道配置费率类型（如果存在则更新，不存在则插入）
-- ============================================================
INSERT INTO channels (channel_code, channel_name, status, config) VALUES (
    'LIANDONG', '联动', 1,
    '{
        "rate_types": [
            {"code": "WECHAT", "name": "微信", "sort_order": 1, "min_rate": "0.30", "max_rate": "0.60"},
            {"code": "ALIPAY", "name": "支付宝", "sort_order": 2, "min_rate": "0.30", "max_rate": "0.60"},
            {"code": "POS_DC", "name": "普通刷卡-借记", "sort_order": 3, "min_rate": "0.45", "max_rate": "0.60"},
            {"code": "POS_CC", "name": "普通刷卡-贷记", "sort_order": 4, "min_rate": "0.50", "max_rate": "0.68"},
            {"code": "POS_DISCOUNT_CC", "name": "特惠", "sort_order": 10, "min_rate": "0.48", "max_rate": "0.60"},
            {"code": "POS_DISCOUNT_GF_CC", "name": "特惠GF", "sort_order": 11, "min_rate": "0.45", "max_rate": "0.58"},
            {"code": "POS_DISCOUNT_PA_CC", "name": "特惠PA", "sort_order": 12, "min_rate": "0.45", "max_rate": "0.58"},
            {"code": "POS_DISCOUNT_MS_CC", "name": "特惠MS", "sort_order": 13, "min_rate": "0.45", "max_rate": "0.58"},
            {"code": "UNIONPAY_DOWN_CC", "name": "云闪付1000-贷记", "sort_order": 5, "min_rate": "0.30", "max_rate": "0.60"},
            {"code": "UNIONPAY_DOWN_DC", "name": "云闪付1000-借记", "sort_order": 6, "min_rate": "0.30", "max_rate": "0.60"}
        ]
    }'::jsonb
) ON CONFLICT (channel_code) DO UPDATE SET
    config = jsonb_set(
        COALESCE(channels.config::jsonb, '{}'::jsonb),
        '{rate_types}',
        '[
            {"code": "WECHAT", "name": "微信", "sort_order": 1, "min_rate": "0.30", "max_rate": "0.60"},
            {"code": "ALIPAY", "name": "支付宝", "sort_order": 2, "min_rate": "0.30", "max_rate": "0.60"},
            {"code": "POS_DC", "name": "普通刷卡-借记", "sort_order": 3, "min_rate": "0.45", "max_rate": "0.60"},
            {"code": "POS_CC", "name": "普通刷卡-贷记", "sort_order": 4, "min_rate": "0.50", "max_rate": "0.68"},
            {"code": "POS_DISCOUNT_CC", "name": "特惠", "sort_order": 10, "min_rate": "0.48", "max_rate": "0.60"},
            {"code": "POS_DISCOUNT_GF_CC", "name": "特惠GF", "sort_order": 11, "min_rate": "0.45", "max_rate": "0.58"},
            {"code": "POS_DISCOUNT_PA_CC", "name": "特惠PA", "sort_order": 12, "min_rate": "0.45", "max_rate": "0.58"},
            {"code": "POS_DISCOUNT_MS_CC", "name": "特惠MS", "sort_order": 13, "min_rate": "0.45", "max_rate": "0.58"},
            {"code": "UNIONPAY_DOWN_CC", "name": "云闪付1000-贷记", "sort_order": 5, "min_rate": "0.30", "max_rate": "0.60"},
            {"code": "UNIONPAY_DOWN_DC", "name": "云闪付1000-借记", "sort_order": 6, "min_rate": "0.30", "max_rate": "0.60"}
        ]'::jsonb
    ),
    updated_at = NOW();

-- ============================================================
-- 注意：旧字段暂时保留，待验证迁移成功后在后续版本中删除
-- 以下为将来删除旧字段的SQL（暂不执行）：
-- ============================================================
-- ALTER TABLE policy_templates
--     DROP COLUMN IF EXISTS credit_rate,
--     DROP COLUMN IF EXISTS debit_rate,
--     DROP COLUMN IF EXISTS debit_cap,
--     DROP COLUMN IF EXISTS unionpay_rate,
--     DROP COLUMN IF EXISTS wechat_rate,
--     DROP COLUMN IF EXISTS alipay_rate;

-- ALTER TABLE agent_policies
--     DROP COLUMN IF EXISTS credit_rate,
--     DROP COLUMN IF EXISTS debit_rate,
--     DROP COLUMN IF EXISTS debit_cap,
--     DROP COLUMN IF EXISTS unionpay_rate,
--     DROP COLUMN IF EXISTS wechat_rate,
--     DROP COLUMN IF EXISTS alipay_rate;

-- ALTER TABLE rate_stage_policies
--     DROP COLUMN IF EXISTS credit_rate_delta,
--     DROP COLUMN IF EXISTS debit_rate_delta,
--     DROP COLUMN IF EXISTS unionpay_rate_delta,
--     DROP COLUMN IF EXISTS wechat_rate_delta,
--     DROP COLUMN IF EXISTS alipay_rate_delta;
