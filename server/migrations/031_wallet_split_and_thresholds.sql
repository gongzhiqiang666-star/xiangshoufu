-- ============================================================
-- 钱包拆分配置与提现门槛配置
-- 需求：按通道拆分钱包显示，提现门槛移至政策模版
-- ============================================================

-- 1. 代理商钱包拆分配置表
-- 控制代理商是否按通道拆分显示分润/服务费钱包
CREATE TABLE IF NOT EXISTS agent_wallet_split_configs (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL UNIQUE,           -- 被配置的代理商
    split_by_channel BOOLEAN DEFAULT FALSE,           -- 是否按通道拆分显示
    configured_by   BIGINT,                           -- 配置人（上级代理商ID或管理员ID）
    configured_at   TIMESTAMPTZ DEFAULT NOW(),        -- 配置时间
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_agent_wallet_split_agent ON agent_wallet_split_configs(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_wallet_split_configured_by ON agent_wallet_split_configs(configured_by);

-- 注释
COMMENT ON TABLE agent_wallet_split_configs IS '代理商钱包拆分配置 - 控制是否按通道拆分显示钱包';
COMMENT ON COLUMN agent_wallet_split_configs.agent_id IS '被配置的代理商ID';
COMMENT ON COLUMN agent_wallet_split_configs.split_by_channel IS '是否按通道拆分：true=分通道显示/提现，false=汇总显示';
COMMENT ON COLUMN agent_wallet_split_configs.configured_by IS '配置人ID（上级代理商或管理员）';
COMMENT ON COLUMN agent_wallet_split_configs.configured_at IS '配置时间';

-- 2. 政策模版提现门槛配置表
-- 存储各政策模版下不同钱包类型、不同通道的提现门槛
CREATE TABLE IF NOT EXISTS policy_withdraw_thresholds (
    id              BIGSERIAL PRIMARY KEY,
    template_id     BIGINT NOT NULL,                  -- 政策模版ID
    wallet_type     SMALLINT NOT NULL,                -- 钱包类型: 1分润 2服务费 3奖励
    channel_id      BIGINT DEFAULT 0,                 -- 通道ID: 0表示所有通道（不拆分时使用）
    threshold_amount BIGINT NOT NULL DEFAULT 10000,   -- 提现门槛金额（分），默认100元
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(template_id, wallet_type, channel_id)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_policy_withdraw_thresholds_template ON policy_withdraw_thresholds(template_id);
CREATE INDEX IF NOT EXISTS idx_policy_withdraw_thresholds_wallet_type ON policy_withdraw_thresholds(wallet_type);

-- 注释
COMMENT ON TABLE policy_withdraw_thresholds IS '政策模版提现门槛配置';
COMMENT ON COLUMN policy_withdraw_thresholds.template_id IS '政策模版ID，关联policy_templates表';
COMMENT ON COLUMN policy_withdraw_thresholds.wallet_type IS '钱包类型: 1=分润钱包, 2=服务费钱包, 3=奖励钱包';
COMMENT ON COLUMN policy_withdraw_thresholds.channel_id IS '通道ID: 0表示通用门槛（不拆分时使用），其他值表示特定通道门槛';
COMMENT ON COLUMN policy_withdraw_thresholds.threshold_amount IS '提现门槛金额（分）';

-- 3. 为现有政策模版初始化默认提现门槛
-- 分润钱包默认门槛100元，服务费钱包默认50元，奖励钱包默认100元
INSERT INTO policy_withdraw_thresholds (template_id, wallet_type, channel_id, threshold_amount)
SELECT id, 1, 0, 10000 FROM policy_templates WHERE id NOT IN (SELECT template_id FROM policy_withdraw_thresholds WHERE wallet_type = 1 AND channel_id = 0)
ON CONFLICT (template_id, wallet_type, channel_id) DO NOTHING;

INSERT INTO policy_withdraw_thresholds (template_id, wallet_type, channel_id, threshold_amount)
SELECT id, 2, 0, 5000 FROM policy_templates WHERE id NOT IN (SELECT template_id FROM policy_withdraw_thresholds WHERE wallet_type = 2 AND channel_id = 0)
ON CONFLICT (template_id, wallet_type, channel_id) DO NOTHING;

INSERT INTO policy_withdraw_thresholds (template_id, wallet_type, channel_id, threshold_amount)
SELECT id, 3, 0, 10000 FROM policy_templates WHERE id NOT IN (SELECT template_id FROM policy_withdraw_thresholds WHERE wallet_type = 3 AND channel_id = 0)
ON CONFLICT (template_id, wallet_type, channel_id) DO NOTHING;
