-- 037_add_rate_limits.sql
-- 添加高调费率和P+0加价上限字段

-- 在 channel_rate_configs 表增加上限字段
ALTER TABLE channel_rate_configs
ADD COLUMN IF NOT EXISTS max_high_rate VARCHAR(10) DEFAULT NULL,  -- 高调费率上限（百分比）
ADD COLUMN IF NOT EXISTS max_d0_extra BIGINT DEFAULT NULL;        -- P+0加价上限（分）

-- 添加字段注释
COMMENT ON COLUMN channel_rate_configs.max_high_rate IS '高调费率上限（百分比，如0.65表示0.65%）';
COMMENT ON COLUMN channel_rate_configs.max_d0_extra IS 'P+0加价上限（分）';

-- 创建通道配置变更影响检查记录表（用于记录变更影响）
CREATE TABLE IF NOT EXISTS channel_config_change_logs (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id),
    change_type VARCHAR(50) NOT NULL,          -- 变更类型: RATE_RANGE, DEPOSIT_LIMIT, SIM_LIMIT
    rate_code VARCHAR(32),                      -- 费率编码（费率变更时）
    old_value JSONB,                            -- 旧值
    new_value JSONB,                            -- 新值
    affected_templates INT DEFAULT 0,           -- 受影响的政策模版数量
    affected_settlements INT DEFAULT 0,         -- 受影响的结算价数量
    affected_agents INT DEFAULT 0,              -- 受影响的代理商数量
    impact_details JSONB,                       -- 影响详情
    operator_id BIGINT NOT NULL,                -- 操作人ID
    operator_name VARCHAR(100),                 -- 操作人名称
    created_at TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_channel_config_change_logs_channel ON channel_config_change_logs(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_config_change_logs_created ON channel_config_change_logs(created_at);

COMMENT ON TABLE channel_config_change_logs IS '通道配置变更记录表';
