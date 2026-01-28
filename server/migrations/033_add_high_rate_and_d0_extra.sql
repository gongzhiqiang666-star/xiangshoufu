-- 033_add_high_rate_and_d0_extra.sql
-- 押金自定义档位、高调政策、P+0加价 功能迁移

-- ============================================================
-- 1. 通道押金档位表
-- ============================================================
CREATE TABLE IF NOT EXISTS channel_deposit_tiers (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL,                    -- 通道ID
    brand_code VARCHAR(32) DEFAULT '',             -- 品牌编码（空表示通用）
    tier_code VARCHAR(32) NOT NULL,                -- 档位编码（如 TIER_99）
    deposit_amount BIGINT NOT NULL,                -- 押金金额（分）
    tier_name VARCHAR(100) NOT NULL,               -- 档位名称（如 99元档）
    sort_order INT DEFAULT 0,                      -- 排序
    status SMALLINT DEFAULT 1,                     -- 状态：1启用 0禁用
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uk_channel_deposit_tier UNIQUE(channel_id, brand_code, tier_code)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_channel_deposit_tiers_channel ON channel_deposit_tiers(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_deposit_tiers_status ON channel_deposit_tiers(status);

COMMENT ON TABLE channel_deposit_tiers IS '通道押金档位配置表';
COMMENT ON COLUMN channel_deposit_tiers.channel_id IS '通道ID';
COMMENT ON COLUMN channel_deposit_tiers.brand_code IS '品牌编码，空表示通用';
COMMENT ON COLUMN channel_deposit_tiers.tier_code IS '档位编码，如TIER_99';
COMMENT ON COLUMN channel_deposit_tiers.deposit_amount IS '押金金额（分）';
COMMENT ON COLUMN channel_deposit_tiers.tier_name IS '档位名称，如99元档';
COMMENT ON COLUMN channel_deposit_tiers.sort_order IS '排序，数字越小越靠前';
COMMENT ON COLUMN channel_deposit_tiers.status IS '状态：1启用 0禁用';

-- ============================================================
-- 2. 政策模版表扩展 - 增加高调费率和P+0加价配置
-- ============================================================
ALTER TABLE policy_templates
    ADD COLUMN IF NOT EXISTS high_rate_configs JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS d0_extra_configs JSONB DEFAULT '{}';

COMMENT ON COLUMN policy_templates.high_rate_configs IS '高调费率配置（按费率类型），格式：{"CREDIT":{"rate":"0.03"},"DEBIT":{"rate":"0.02"}}';
COMMENT ON COLUMN policy_templates.d0_extra_configs IS 'P+0加价配置（按费率类型，单位分），格式：{"CREDIT":{"extra_fee":100},"DEBIT":{"extra_fee":50}}';

-- ============================================================
-- 3. 结算价表扩展 - 增加高调费率和P+0加价配置
-- ============================================================
ALTER TABLE settlement_prices
    ADD COLUMN IF NOT EXISTS high_rate_configs JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS d0_extra_configs JSONB DEFAULT '{}';

COMMENT ON COLUMN settlement_prices.high_rate_configs IS '高调费率配置（按费率类型），格式：{"CREDIT":{"rate":"0.03"},"DEBIT":{"rate":"0.02"}}';
COMMENT ON COLUMN settlement_prices.d0_extra_configs IS 'P+0加价配置（按费率类型，单位分），格式：{"CREDIT":{"extra_fee":100},"DEBIT":{"extra_fee":50}}';

-- ============================================================
-- 4. 分润记录表扩展 - 增加高调分润和P+0分润字段
-- ============================================================
ALTER TABLE profit_records
    ADD COLUMN IF NOT EXISTS high_rate_profit BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS d0_extra_profit BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS high_rate_self DECIMAL(10,4),
    ADD COLUMN IF NOT EXISTS high_rate_lower DECIMAL(10,4),
    ADD COLUMN IF NOT EXISTS d0_extra_self BIGINT DEFAULT 0,
    ADD COLUMN IF NOT EXISTS d0_extra_lower BIGINT DEFAULT 0;

COMMENT ON COLUMN profit_records.high_rate_profit IS '高调分润金额（分）';
COMMENT ON COLUMN profit_records.d0_extra_profit IS 'P+0分润金额（分）';
COMMENT ON COLUMN profit_records.high_rate_self IS '自身高调费率';
COMMENT ON COLUMN profit_records.high_rate_lower IS '下级高调费率';
COMMENT ON COLUMN profit_records.d0_extra_self IS '自身P+0加价配置（分）';
COMMENT ON COLUMN profit_records.d0_extra_lower IS '下级P+0加价配置（分）';

-- ============================================================
-- 5. 索引优化
-- ============================================================
-- 为分润记录表增加高调和P+0分润查询索引
CREATE INDEX IF NOT EXISTS idx_profit_records_high_rate ON profit_records(high_rate_profit) WHERE high_rate_profit > 0;
CREATE INDEX IF NOT EXISTS idx_profit_records_d0_extra ON profit_records(d0_extra_profit) WHERE d0_extra_profit > 0;

-- ============================================================
-- 6. 初始化常用押金档位数据（示例）
-- ============================================================
-- 注意：实际使用时需要根据各通道的实际押金档位配置
-- INSERT INTO channel_deposit_tiers (channel_id, tier_code, deposit_amount, tier_name, sort_order) VALUES
-- (1, 'TIER_99', 9900, '99元档', 1),
-- (1, 'TIER_199', 19900, '199元档', 2),
-- (1, 'TIER_299', 29900, '299元档', 3);
