-- ============================================================
-- 全局提现门槛配置
-- 需求：在钱包管理模块统一配置提现门槛，与政策模版无关
-- ============================================================

-- 1. 全局提现门槛配置表
CREATE TABLE IF NOT EXISTS global_withdraw_thresholds (
    id              BIGSERIAL PRIMARY KEY,
    wallet_type     SMALLINT NOT NULL,                -- 钱包类型: 1分润 2服务费 3奖励
    channel_id      BIGINT DEFAULT 0,                 -- 通道ID: 0表示通用门槛
    threshold_amount BIGINT NOT NULL DEFAULT 10000,   -- 提现门槛金额（分），默认100元
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(wallet_type, channel_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_global_withdraw_thresholds_wallet_type ON global_withdraw_thresholds(wallet_type);
CREATE INDEX IF NOT EXISTS idx_global_withdraw_thresholds_channel ON global_withdraw_thresholds(channel_id);

-- 注释
COMMENT ON TABLE global_withdraw_thresholds IS '全局提现门槛配置 - 统一管理各钱包类型的提现门槛';
COMMENT ON COLUMN global_withdraw_thresholds.wallet_type IS '钱包类型: 1=分润钱包, 2=服务费钱包, 3=奖励钱包';
COMMENT ON COLUMN global_withdraw_thresholds.channel_id IS '通道ID: 0表示通用门槛，其他值表示特定通道门槛（优先级高于通用门槛）';
COMMENT ON COLUMN global_withdraw_thresholds.threshold_amount IS '提现门槛金额（分）';

-- 2. 初始化默认门槛配置
-- 分润钱包默认门槛100元，服务费钱包默认50元，奖励钱包默认100元
INSERT INTO global_withdraw_thresholds (wallet_type, channel_id, threshold_amount) VALUES
    (1, 0, 10000),  -- 分润钱包通用门槛 100元
    (2, 0, 5000),   -- 服务费钱包通用门槛 50元
    (3, 0, 10000)   -- 奖励钱包通用门槛 100元
ON CONFLICT (wallet_type, channel_id) DO NOTHING;
