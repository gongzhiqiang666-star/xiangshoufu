-- 流量费返现相关表
-- 创建时间: 2026-01-18
-- 业务规则:
--   Q30: 流量费返现三档：首次/2次/2+N次，按级差计算

-- 流量费返现政策表（三档配置）
CREATE TABLE IF NOT EXISTS sim_cashback_policies (
    id BIGSERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL,                   -- 政策模板ID
    channel_id BIGINT NOT NULL,                    -- 通道ID
    brand_code VARCHAR(32),                        -- 品牌编码
    first_time_cashback BIGINT NOT NULL,           -- 首次返现金额（分）
    second_time_cashback BIGINT NOT NULL,          -- 第2次返现金额（分）
    third_plus_cashback BIGINT NOT NULL,           -- 第3次及以后返现金额（分）
    sim_fee_amount BIGINT NOT NULL,                -- 流量费金额（分）
    status SMALLINT DEFAULT 1,                     -- 1:启用 0:禁用
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_sim_cashback_policies_template ON sim_cashback_policies(template_id);
CREATE INDEX idx_sim_cashback_policies_channel ON sim_cashback_policies(channel_id);
CREATE UNIQUE INDEX idx_sim_cashback_policies_unique ON sim_cashback_policies(template_id, channel_id, brand_code) WHERE status = 1;

-- 流量费返现记录表
CREATE TABLE IF NOT EXISTS sim_cashback_records (
    id BIGSERIAL PRIMARY KEY,
    device_fee_id BIGINT NOT NULL,                 -- 关联流量费记录ID
    terminal_sn VARCHAR(50) NOT NULL,              -- 终端SN
    channel_id BIGINT NOT NULL,                    -- 通道ID
    agent_id BIGINT NOT NULL,                      -- 获得返现的代理商
    sim_fee_count INT NOT NULL,                    -- 当前是第几次缴费
    sim_fee_amount BIGINT NOT NULL,                -- 流量费金额（分）
    cashback_tier SMALLINT NOT NULL,               -- 返现档次 1:首次 2:第2次 3:第3次及以后
    self_cashback BIGINT NOT NULL,                 -- 自身返现配置金额（分）
    upper_cashback BIGINT NOT NULL,                -- 上级应返金额（用于级差计算）（分）
    actual_cashback BIGINT NOT NULL,               -- 实际返现金额（级差）（分）
    source_agent_id BIGINT,                        -- 下级代理商ID（级差来源）
    wallet_type SMALLINT DEFAULT 1,                -- 钱包类型
    wallet_status SMALLINT DEFAULT 0,              -- 0:待入账 1:已入账
    created_at TIMESTAMP DEFAULT NOW(),
    processed_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_sim_cashback_records_device_fee ON sim_cashback_records(device_fee_id);
CREATE INDEX idx_sim_cashback_records_terminal ON sim_cashback_records(terminal_sn);
CREATE INDEX idx_sim_cashback_records_agent ON sim_cashback_records(agent_id);
CREATE INDEX idx_sim_cashback_records_created ON sim_cashback_records(created_at);

-- 添加注释
COMMENT ON TABLE sim_cashback_policies IS '流量费返现政策表 - 三档配置（首次/2次/2+N次）';
COMMENT ON COLUMN sim_cashback_policies.first_time_cashback IS '首次缴费返现金额（分）';
COMMENT ON COLUMN sim_cashback_policies.second_time_cashback IS '第2次缴费返现金额（分）';
COMMENT ON COLUMN sim_cashback_policies.third_plus_cashback IS '第3次及以后缴费返现金额（分）';

COMMENT ON TABLE sim_cashback_records IS '流量费返现记录表 - 按级差计算';
COMMENT ON COLUMN sim_cashback_records.cashback_tier IS '返现档次：1首次 2第2次 3第3次及以后';
COMMENT ON COLUMN sim_cashback_records.actual_cashback IS '实际返现金额 = 自身配置 - 下级配置（级差）';

-- 示例数据：流量费返现政策配置
-- INSERT INTO sim_cashback_policies (template_id, channel_id, brand_code, first_time_cashback, second_time_cashback, third_plus_cashback, sim_fee_amount)
-- VALUES
--     (1, 1, 'HENGXINTONG', 5000, 3000, 2000, 9900),  -- 模板1: 首次50元, 第2次30元, 第3次+20元
--     (2, 1, 'HENGXINTONG', 4000, 2500, 1500, 9900),  -- 模板2: 首次40元, 第2次25元, 第3次+15元
--     (3, 1, 'HENGXINTONG', 3000, 2000, 1000, 9900);  -- 模板3: 首次30元, 第2次20元, 第3次+10元
