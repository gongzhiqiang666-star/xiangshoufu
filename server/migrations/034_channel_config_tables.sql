-- ============================================================
-- 通道配置管理表迁移
-- 建立 通道配置 → 政策模板 → 结算价 的三层继承链
-- ============================================================

-- 1. 创建通道费率配置表
CREATE TABLE IF NOT EXISTS channel_rate_configs (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id),
    rate_code VARCHAR(32) NOT NULL,           -- 费率编码 (CREDIT/DEBIT/WECHAT等)
    rate_name VARCHAR(64) NOT NULL,           -- 费率名称
    min_rate DECIMAL(10,4),                   -- 最低成本（通道底价）
    max_rate DECIMAL(10,4),                   -- 最高限制
    default_rate DECIMAL(10,4),               -- 默认费率
    sort_order INT DEFAULT 0,                 -- 排序
    status SMALLINT DEFAULT 1,                -- 状态：1启用 0禁用
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(channel_id, rate_code)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_channel_rate_configs_channel_id ON channel_rate_configs(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_rate_configs_status ON channel_rate_configs(status);

COMMENT ON TABLE channel_rate_configs IS '通道费率配置表';
COMMENT ON COLUMN channel_rate_configs.rate_code IS '费率编码，与支付回调的payTypeCode一致';
COMMENT ON COLUMN channel_rate_configs.min_rate IS '通道成本底价，政策模板和结算价不能低于此值';
COMMENT ON COLUMN channel_rate_configs.max_rate IS '费率上限，政策模板和结算价不能超过此值';

-- 2. 创建通道流量费返现档位表
CREATE TABLE IF NOT EXISTS channel_sim_cashback_tiers (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id),
    brand_code VARCHAR(32) DEFAULT '',        -- 品牌编码（空=通用）
    tier_order INT NOT NULL,                  -- 档位序号 1=首次, 2=第2次...
    tier_name VARCHAR(64) NOT NULL,           -- 档位名称
    is_last_tier BOOLEAN DEFAULT FALSE,       -- 是否最后档(N次及以后)
    max_cashback_amount BIGINT NOT NULL,      -- 返现上限（分）
    default_cashback BIGINT DEFAULT 0,        -- 默认返现（分）
    sim_fee_amount BIGINT NOT NULL,           -- 流量费金额（分）
    status SMALLINT DEFAULT 1,                -- 状态：1启用 0禁用
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(channel_id, brand_code, tier_order)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_channel_sim_cashback_tiers_channel_id ON channel_sim_cashback_tiers(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_sim_cashback_tiers_brand_code ON channel_sim_cashback_tiers(brand_code);

COMMENT ON TABLE channel_sim_cashback_tiers IS '通道流量费返现档位表';
COMMENT ON COLUMN channel_sim_cashback_tiers.tier_order IS '档位序号，1=首次缴费, 2=第2次缴费...';
COMMENT ON COLUMN channel_sim_cashback_tiers.is_last_tier IS '是否为最后档位，标记为true时表示该档位适用于N次及以后';
COMMENT ON COLUMN channel_sim_cashback_tiers.max_cashback_amount IS '该档位的返现上限，政策模板和结算价不能超过此值';

-- 3. 修改通道押金档位表，增加返现上限字段
ALTER TABLE channel_deposit_tiers
ADD COLUMN IF NOT EXISTS max_cashback_amount BIGINT DEFAULT 0,
ADD COLUMN IF NOT EXISTS default_cashback BIGINT DEFAULT 0;

COMMENT ON COLUMN channel_deposit_tiers.max_cashback_amount IS '该押金档位的返现上限（分）';
COMMENT ON COLUMN channel_deposit_tiers.default_cashback IS '该押金档位的默认返现金额（分）';

-- 4. 修改结算价表，增加N档流量费返现字段
ALTER TABLE settlement_prices
ADD COLUMN IF NOT EXISTS sim_cashbacks JSONB DEFAULT '[]';

COMMENT ON COLUMN settlement_prices.sim_cashbacks IS '流量费返现N档配置，格式: [{tier_order: 1, cashback_amount: 3000}, ...]';

-- 5. 修改流量费返现政策表（政策模板级别）
ALTER TABLE sim_cashback_policies
ADD COLUMN IF NOT EXISTS cashback_tiers JSONB DEFAULT '[]';

COMMENT ON COLUMN sim_cashback_policies.cashback_tiers IS '流量费返现N档配置，格式同settlement_prices.sim_cashbacks';

-- ============================================================
-- 初始化恒信通通道的默认配置数据
-- ============================================================

-- 获取恒信通通道ID
DO $$
DECLARE
    hxt_channel_id BIGINT;
BEGIN
    SELECT id INTO hxt_channel_id FROM channels WHERE channel_code = 'HENGXINTONG' LIMIT 1;

    IF hxt_channel_id IS NOT NULL THEN
        -- 插入费率配置
        INSERT INTO channel_rate_configs (channel_id, rate_code, rate_name, min_rate, max_rate, default_rate, sort_order, status)
        VALUES
            (hxt_channel_id, 'CREDIT', '贷记卡', 0.38, 0.60, 0.55, 1, 1),
            (hxt_channel_id, 'DEBIT', '借记卡', 0.35, 0.55, 0.50, 2, 1),
            (hxt_channel_id, 'UNIONPAY', '云闪付', 0.35, 0.55, 0.50, 3, 1),
            (hxt_channel_id, 'WECHAT', '微信支付', 0.35, 0.55, 0.50, 4, 1),
            (hxt_channel_id, 'ALIPAY', '支付宝', 0.35, 0.55, 0.50, 5, 1)
        ON CONFLICT (channel_id, rate_code) DO UPDATE SET
            rate_name = EXCLUDED.rate_name,
            min_rate = EXCLUDED.min_rate,
            max_rate = EXCLUDED.max_rate,
            default_rate = EXCLUDED.default_rate,
            sort_order = EXCLUDED.sort_order,
            updated_at = NOW();

        -- 更新押金档位的返现上限
        UPDATE channel_deposit_tiers
        SET
            max_cashback_amount = CASE
                WHEN deposit_amount = 9900 THEN 5000    -- 99元押金，返现上限50元
                WHEN deposit_amount = 19900 THEN 10000  -- 199元押金，返现上限100元
                WHEN deposit_amount = 29900 THEN 15000  -- 299元押金，返现上限150元
                ELSE deposit_amount / 2                 -- 默认上限为押金的一半
            END,
            default_cashback = CASE
                WHEN deposit_amount = 9900 THEN 4000    -- 99元押金，默认返现40元
                WHEN deposit_amount = 19900 THEN 8000   -- 199元押金，默认返现80元
                WHEN deposit_amount = 29900 THEN 12000  -- 299元押金，默认返现120元
                ELSE deposit_amount / 3                 -- 默认返现为押金的1/3
            END,
            updated_at = NOW()
        WHERE channel_id = hxt_channel_id;

        -- 插入流量费返现档位（4档：首次、二次、三次、四次及以后）
        INSERT INTO channel_sim_cashback_tiers (channel_id, brand_code, tier_order, tier_name, is_last_tier, max_cashback_amount, default_cashback, sim_fee_amount, status)
        VALUES
            (hxt_channel_id, '', 1, '首次缴费', FALSE, 3000, 2500, 3600, 1),
            (hxt_channel_id, '', 2, '第2次缴费', FALSE, 2500, 2000, 3600, 1),
            (hxt_channel_id, '', 3, '第3次缴费', FALSE, 2000, 1500, 3600, 1),
            (hxt_channel_id, '', 4, '第4次及以后', TRUE, 1500, 1000, 3600, 1)
        ON CONFLICT (channel_id, brand_code, tier_order) DO UPDATE SET
            tier_name = EXCLUDED.tier_name,
            is_last_tier = EXCLUDED.is_last_tier,
            max_cashback_amount = EXCLUDED.max_cashback_amount,
            default_cashback = EXCLUDED.default_cashback,
            sim_fee_amount = EXCLUDED.sim_fee_amount,
            updated_at = NOW();

        RAISE NOTICE '恒信通通道配置初始化完成，channel_id: %', hxt_channel_id;
    ELSE
        RAISE NOTICE '未找到恒信通通道，跳过初始化';
    END IF;
END $$;

-- ============================================================
-- 数据迁移：将现有结算价的流量费返现转换为新格式
-- ============================================================

-- 将旧的三档流量费返现迁移到新的N档格式
UPDATE settlement_prices
SET sim_cashbacks = jsonb_build_array(
    jsonb_build_object('tier_order', 1, 'cashback_amount', COALESCE(sim_first_cashback, 0)),
    jsonb_build_object('tier_order', 2, 'cashback_amount', COALESCE(sim_second_cashback, 0)),
    jsonb_build_object('tier_order', 3, 'cashback_amount', COALESCE(sim_third_plus_cashback, 0)),
    jsonb_build_object('tier_order', 4, 'cashback_amount', COALESCE(sim_third_plus_cashback, 0))
)
WHERE sim_cashbacks = '[]'::jsonb
  AND (sim_first_cashback > 0 OR sim_second_cashback > 0 OR sim_third_plus_cashback > 0);

-- 同样迁移政策模板的流量费返现
UPDATE sim_cashback_policies
SET cashback_tiers = jsonb_build_array(
    jsonb_build_object('tier_order', 1, 'cashback_amount', COALESCE(first_time_cashback, 0)),
    jsonb_build_object('tier_order', 2, 'cashback_amount', COALESCE(second_time_cashback, 0)),
    jsonb_build_object('tier_order', 3, 'cashback_amount', COALESCE(third_plus_cashback, 0)),
    jsonb_build_object('tier_order', 4, 'cashback_amount', COALESCE(third_plus_cashback, 0))
)
WHERE cashback_tiers = '[]'::jsonb
  AND (first_time_cashback > 0 OR second_time_cashback > 0 OR third_plus_cashback > 0);

RAISE NOTICE '通道配置管理表迁移完成';
