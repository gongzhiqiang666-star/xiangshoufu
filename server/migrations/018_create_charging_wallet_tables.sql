-- 充值钱包相关表
-- 充值钱包用于代理商向下级发放奖励，需要PC端开通权限后APP端显示

-- 代理商特殊钱包配置表
CREATE TABLE IF NOT EXISTS agent_wallet_configs (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL UNIQUE,

    -- 充值钱包配置
    charging_wallet_enabled    BOOLEAN DEFAULT FALSE,   -- 是否开通充值钱包
    charging_wallet_limit      BIGINT DEFAULT 0,        -- 充值钱包限额(分)

    -- 沉淀钱包配置
    settlement_wallet_enabled  BOOLEAN DEFAULT FALSE,   -- 是否开通沉淀钱包
    settlement_ratio           INT DEFAULT 30,          -- 沉淀比例(百分比，默认30%)

    -- 审计字段
    enabled_by      BIGINT,                             -- 开通人
    enabled_at      TIMESTAMP,                          -- 开通时间
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_wallet_configs_agent ON agent_wallet_configs(agent_id);

-- 充值钱包充值记录表
CREATE TABLE IF NOT EXISTS charging_wallet_deposits (
    id              BIGSERIAL PRIMARY KEY,
    deposit_no      VARCHAR(50) NOT NULL UNIQUE,        -- 充值单号
    agent_id        BIGINT NOT NULL,                    -- 代理商ID
    amount          BIGINT NOT NULL,                    -- 充值金额(分)
    payment_method  SMALLINT DEFAULT 1,                 -- 支付方式: 1=银行转账 2=微信 3=支付宝
    payment_ref     VARCHAR(100),                       -- 支付流水号
    status          SMALLINT DEFAULT 0,                 -- 状态: 0=待确认 1=已确认 2=已拒绝

    -- 审核信息
    confirmed_by    BIGINT,                             -- 确认人
    confirmed_at    TIMESTAMP,                          -- 确认时间
    reject_reason   VARCHAR(500),                       -- 拒绝原因

    remark          VARCHAR(500),                       -- 备注
    created_by      BIGINT,                             -- 创建人
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_charging_deposits_agent ON charging_wallet_deposits(agent_id);
CREATE INDEX IF NOT EXISTS idx_charging_deposits_status ON charging_wallet_deposits(status);

-- 充值钱包奖励发放记录表
CREATE TABLE IF NOT EXISTS charging_wallet_rewards (
    id              BIGSERIAL PRIMARY KEY,
    reward_no       VARCHAR(50) NOT NULL UNIQUE,        -- 奖励单号
    from_agent_id   BIGINT NOT NULL,                    -- 发放方代理商ID
    to_agent_id     BIGINT NOT NULL,                    -- 接收方代理商ID
    amount          BIGINT NOT NULL,                    -- 奖励金额(分)
    reward_type     SMALLINT DEFAULT 1,                 -- 奖励类型: 1=手动发放 2=自动发放(政策触发)
    policy_id       BIGINT,                             -- 关联的奖励政策ID

    -- 状态
    status          SMALLINT DEFAULT 1,                 -- 状态: 1=已发放 2=已撤销
    revoked_at      TIMESTAMP,                          -- 撤销时间
    revoke_reason   VARCHAR(500),                       -- 撤销原因

    remark          VARCHAR(500),                       -- 备注
    created_by      BIGINT,                             -- 创建人
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_charging_rewards_from ON charging_wallet_rewards(from_agent_id);
CREATE INDEX IF NOT EXISTS idx_charging_rewards_to ON charging_wallet_rewards(to_agent_id);

-- 沉淀钱包使用记录表
CREATE TABLE IF NOT EXISTS settlement_wallet_usages (
    id              BIGSERIAL PRIMARY KEY,
    usage_no        VARCHAR(50) NOT NULL UNIQUE,        -- 使用单号
    agent_id        BIGINT NOT NULL,                    -- 代理商ID
    amount          BIGINT NOT NULL,                    -- 使用金额(分)
    usage_type      SMALLINT DEFAULT 1,                 -- 类型: 1=使用 2=归还

    -- 来源明细(哪些下级的钱)
    source_details  JSONB,                              -- 来源明细JSON

    -- 状态
    status          SMALLINT DEFAULT 1,                 -- 状态: 1=正常 2=待归还
    return_deadline TIMESTAMP,                          -- 归还截止时间
    returned_at     TIMESTAMP,                          -- 归还时间

    remark          VARCHAR(500),                       -- 备注
    created_by      BIGINT,                             -- 创建人
    created_at      TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_settlement_usages_agent ON settlement_wallet_usages(agent_id);
CREATE INDEX IF NOT EXISTS idx_settlement_usages_status ON settlement_wallet_usages(status);

-- 添加注释
COMMENT ON TABLE agent_wallet_configs IS '代理商特殊钱包配置表';
COMMENT ON COLUMN agent_wallet_configs.charging_wallet_enabled IS '是否开通充值钱包';
COMMENT ON COLUMN agent_wallet_configs.settlement_wallet_enabled IS '是否开通沉淀钱包';
COMMENT ON COLUMN agent_wallet_configs.settlement_ratio IS '沉淀钱包可用比例(百分比)';

COMMENT ON TABLE charging_wallet_deposits IS '充值钱包充值记录表';
COMMENT ON TABLE charging_wallet_rewards IS '充值钱包奖励发放记录表';
COMMENT ON TABLE settlement_wallet_usages IS '沉淀钱包使用记录表';

-- 钱包类型说明:
-- 1 = 分润钱包 (profit wallet)
-- 2 = 服务费钱包 (service fee wallet, includes SIM fee + deposit cashback)
-- 3 = 奖励钱包 (reward wallet)
-- 4 = 充值钱包 (charging wallet) - 新增
-- 5 = 沉淀钱包 (settlement wallet) - 新增
