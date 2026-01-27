-- 结算价与调价记录表
-- 创建时间: 2026-01-27
-- 业务规则:
--   1. settlement_prices: 通道维度结算价（费率、押金返现、流量费返现）
--   2. agent_reward_settings: 代理商维度奖励配置（与通道无关）
--   3. price_change_logs: 统一调价记录（支持PC/APP端查询）
--   4. 分润计算必须基于结算价表，而非政策模版

-- ============================================================
-- 1. 通道结算价表
-- ============================================================
CREATE TABLE IF NOT EXISTS settlement_prices (
    id                      BIGSERIAL PRIMARY KEY,
    agent_id                BIGINT NOT NULL,                    -- 代理商ID
    channel_id              BIGINT NOT NULL,                    -- 通道ID
    template_id             BIGINT,                             -- 来源政策模版ID
    brand_code              VARCHAR(32) DEFAULT '',             -- 品牌编码

    -- 费率配置（动态JSONB）
    rate_configs            JSONB DEFAULT '{}',                 -- 费率配置 {"credit": {"rate": "0.60"}, ...}

    -- 旧字段保留兼容
    credit_rate             DECIMAL(10,4),                      -- 贷记卡费率
    debit_rate              DECIMAL(10,4),                      -- 借记卡费率
    debit_cap               DECIMAL(10,2),                      -- 借记卡封顶
    unionpay_rate           DECIMAL(10,4),                      -- 云闪付费率
    wechat_rate             DECIMAL(10,4),                      -- 微信费率
    alipay_rate             DECIMAL(10,4),                      -- 支付宝费率

    -- 押金返现配置（JSONB数组）
    deposit_cashbacks       JSONB DEFAULT '[]',                 -- [{"deposit_amount": 9900, "cashback_amount": 5000}, ...]

    -- 流量费返现配置
    sim_first_cashback      BIGINT DEFAULT 0,                   -- 首次返现金额（分）
    sim_second_cashback     BIGINT DEFAULT 0,                   -- 第2次返现金额（分）
    sim_third_plus_cashback BIGINT DEFAULT 0,                   -- 第3次及以后返现金额（分）

    -- 元数据
    version                 INT DEFAULT 1,                      -- 版本号（每次调价+1）
    status                  SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    effective_at            TIMESTAMPTZ DEFAULT NOW(),          -- 生效时间
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW(),
    created_by              BIGINT,                             -- 创建人ID
    updated_by              BIGINT                              -- 修改人ID
);

-- 索引
CREATE INDEX idx_settlement_prices_agent ON settlement_prices(agent_id);
CREATE INDEX idx_settlement_prices_channel ON settlement_prices(channel_id);
CREATE INDEX idx_settlement_prices_template ON settlement_prices(template_id);
CREATE INDEX idx_settlement_prices_status ON settlement_prices(status);
CREATE UNIQUE INDEX idx_settlement_prices_unique ON settlement_prices(agent_id, channel_id, brand_code) WHERE status = 1;

-- 注释
COMMENT ON TABLE settlement_prices IS '通道结算价表 - 代理商实际使用的价格配置';
COMMENT ON COLUMN settlement_prices.rate_configs IS '费率配置JSONB，格式: {"credit": {"rate": "0.60"}, "debit": {"rate": "0.50"}}';
COMMENT ON COLUMN settlement_prices.deposit_cashbacks IS '押金返现配置JSONB数组，格式: [{"deposit_amount": 9900, "cashback_amount": 5000}]';
COMMENT ON COLUMN settlement_prices.version IS '版本号，每次调价自增';

-- ============================================================
-- 2. 代理商奖励配置表
-- ============================================================
CREATE TABLE IF NOT EXISTS agent_reward_settings (
    id                      BIGSERIAL PRIMARY KEY,
    agent_id                BIGINT NOT NULL UNIQUE,             -- 代理商ID（唯一）
    template_id             BIGINT,                             -- 来源奖励政策模版ID

    -- 奖励金额（差额分配模式）
    reward_amount           BIGINT DEFAULT 0,                   -- 奖励金额（分），上级给下级配置

    -- 激活奖励配置（JSONB数组）
    activation_rewards      JSONB DEFAULT '[]',                 -- 激活奖励配置数组

    -- 元数据
    version                 INT DEFAULT 1,                      -- 版本号
    status                  SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    effective_at            TIMESTAMPTZ DEFAULT NOW(),          -- 生效时间
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW(),
    created_by              BIGINT,                             -- 创建人ID
    updated_by              BIGINT                              -- 修改人ID
);

-- 索引
CREATE INDEX idx_agent_reward_settings_template ON agent_reward_settings(template_id);
CREATE INDEX idx_agent_reward_settings_status ON agent_reward_settings(status);

-- 注释
COMMENT ON TABLE agent_reward_settings IS '代理商奖励配置表 - 与通道无关的奖励配置';
COMMENT ON COLUMN agent_reward_settings.reward_amount IS '奖励金额（分），用于差额分配模式';
COMMENT ON COLUMN agent_reward_settings.activation_rewards IS '激活奖励配置JSONB数组';

-- ============================================================
-- 3. 调价记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS price_change_logs (
    id                      BIGSERIAL PRIMARY KEY,
    agent_id                BIGINT NOT NULL,                    -- 代理商ID
    channel_id              BIGINT,                             -- 通道ID（结算价调整时有值）
    settlement_price_id     BIGINT,                             -- 结算价ID
    reward_setting_id       BIGINT,                             -- 奖励配置ID

    -- 变更类型
    change_type             SMALLINT NOT NULL,                  -- 变更类型：1-初始化 2-费率调整 3-押金返现调整 4-流量费返现调整 5-激活奖励调整 6-批量调整 7-模板同步
    config_type             SMALLINT NOT NULL,                  -- 配置类型：1-结算价 2-奖励配置

    -- 变更内容
    field_name              VARCHAR(100),                       -- 变更字段名
    old_value               TEXT,                               -- 旧值
    new_value               TEXT,                               -- 新值
    change_summary          VARCHAR(500),                       -- 变更摘要（前端展示用）

    -- 完整快照
    snapshot_before         JSONB,                              -- 变更前完整快照
    snapshot_after          JSONB,                              -- 变更后完整快照

    -- 操作信息
    operator_type           SMALLINT DEFAULT 1,                 -- 操作者类型：1-管理员 2-代理商
    operator_id             BIGINT NOT NULL,                    -- 操作者ID
    operator_name           VARCHAR(100),                       -- 操作者名称
    source                  VARCHAR(20) DEFAULT 'PC',           -- 来源：PC/APP
    ip_address              VARCHAR(50),                        -- IP地址
    user_agent              VARCHAR(500),                       -- 用户代理

    created_at              TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_price_change_logs_agent ON price_change_logs(agent_id);
CREATE INDEX idx_price_change_logs_channel ON price_change_logs(channel_id);
CREATE INDEX idx_price_change_logs_settlement ON price_change_logs(settlement_price_id);
CREATE INDEX idx_price_change_logs_reward ON price_change_logs(reward_setting_id);
CREATE INDEX idx_price_change_logs_type ON price_change_logs(change_type);
CREATE INDEX idx_price_change_logs_config_type ON price_change_logs(config_type);
CREATE INDEX idx_price_change_logs_operator ON price_change_logs(operator_id);
CREATE INDEX idx_price_change_logs_created ON price_change_logs(created_at DESC);

-- 注释
COMMENT ON TABLE price_change_logs IS '调价记录表 - 记录结算价和奖励配置的所有变更';
COMMENT ON COLUMN price_change_logs.change_type IS '变更类型：1-初始化 2-费率调整 3-押金返现调整 4-流量费返现调整 5-激活奖励调整 6-批量调整 7-模板同步';
COMMENT ON COLUMN price_change_logs.config_type IS '配置类型：1-结算价 2-奖励配置';
COMMENT ON COLUMN price_change_logs.source IS '操作来源：PC-PC端管理系统 APP-移动端APP';

-- ============================================================
-- 4. 变更类型枚举说明
-- ============================================================
-- change_type:
--   1 = INIT          初始化（从模板创建）
--   2 = RATE          费率调整
--   3 = DEPOSIT       押金返现调整
--   4 = SIM           流量费返现调整
--   5 = ACTIVATION    激活奖励调整
--   6 = BATCH         批量调整
--   7 = SYNC          模板同步

-- config_type:
--   1 = SETTLEMENT    结算价
--   2 = REWARD        奖励配置

-- operator_type:
--   1 = ADMIN         管理员
--   2 = AGENT         代理商
