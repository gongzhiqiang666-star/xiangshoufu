-- 政策模版完整表 - 押金返现、激活奖励、费率阶梯
-- 创建时间: 2026-01-22
-- 业务规则:
--   Q31: 政策模版按通道区分，分4块：成本(费率)、押金返现、流量卡返现、奖励返现
--   Q32: 押金返现 - 商户押金收取后返现给代理商，入服务费钱包
--   Q33: 激活奖励 - 按入网时间+交易量条件触发奖励，入奖励钱包
--   Q34: 费率阶梯 - 按商户/代理商入网时间自动调整费率

-- ============================================================
-- 1. 押金返现政策表
-- ============================================================
CREATE TABLE IF NOT EXISTS deposit_cashback_policies (
    id                  BIGSERIAL PRIMARY KEY,
    template_id         BIGINT NOT NULL,                    -- 政策模板ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    brand_code          VARCHAR(32),                        -- 品牌编码
    deposit_amount      BIGINT NOT NULL,                    -- 押金金额（分）：0/9900/19900/29900
    cashback_amount     BIGINT NOT NULL,                    -- 返现金额（分）
    status              SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_deposit_cashback_policies_template ON deposit_cashback_policies(template_id);
CREATE INDEX idx_deposit_cashback_policies_channel ON deposit_cashback_policies(channel_id);
CREATE UNIQUE INDEX idx_deposit_cashback_policies_unique ON deposit_cashback_policies(template_id, channel_id, deposit_amount) WHERE status = 1;

-- 注释
COMMENT ON TABLE deposit_cashback_policies IS '押金返现政策表 - 按押金金额配置返现';
COMMENT ON COLUMN deposit_cashback_policies.deposit_amount IS '押金金额（分）：0/9900/19900/29900';
COMMENT ON COLUMN deposit_cashback_policies.cashback_amount IS '返现金额（分）';

-- ============================================================
-- 2. 押金返现记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS deposit_cashback_records (
    id                  BIGSERIAL PRIMARY KEY,
    terminal_id         BIGINT NOT NULL,                    -- 终端ID
    terminal_sn         VARCHAR(50) NOT NULL,               -- 终端SN
    merchant_id         BIGINT NOT NULL,                    -- 商户ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    agent_id            BIGINT NOT NULL,                    -- 获得返现的代理商
    deposit_amount      BIGINT NOT NULL,                    -- 押金金额（分）
    self_cashback       BIGINT NOT NULL,                    -- 自身返现配置金额（分）
    upper_cashback      BIGINT NOT NULL,                    -- 上级应返金额（用于级差计算）（分）
    actual_cashback     BIGINT NOT NULL,                    -- 实际返现金额（级差）（分）
    source_agent_id     BIGINT,                             -- 下级代理商ID（级差来源）
    wallet_type         SMALLINT DEFAULT 2,                 -- 钱包类型：2-服务费钱包
    wallet_status       SMALLINT DEFAULT 0,                 -- 0:待入账 1:已入账
    trigger_type        SMALLINT DEFAULT 1,                 -- 触发类型：1-押金扣款 2-手动触发
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    processed_at        TIMESTAMPTZ
);

-- 索引
CREATE INDEX idx_deposit_cashback_records_terminal ON deposit_cashback_records(terminal_id);
CREATE INDEX idx_deposit_cashback_records_merchant ON deposit_cashback_records(merchant_id);
CREATE INDEX idx_deposit_cashback_records_agent ON deposit_cashback_records(agent_id);
CREATE INDEX idx_deposit_cashback_records_status ON deposit_cashback_records(wallet_status);
CREATE INDEX idx_deposit_cashback_records_created ON deposit_cashback_records(created_at);

-- 注释
COMMENT ON TABLE deposit_cashback_records IS '押金返现记录表 - 按级差计算';
COMMENT ON COLUMN deposit_cashback_records.actual_cashback IS '实际返现金额 = 自身配置 - 下级配置（级差）';
COMMENT ON COLUMN deposit_cashback_records.wallet_type IS '钱包类型：2-服务费钱包';

-- ============================================================
-- 3. 激活奖励政策表
-- ============================================================
CREATE TABLE IF NOT EXISTS activation_reward_policies (
    id                  BIGSERIAL PRIMARY KEY,
    template_id         BIGINT NOT NULL,                    -- 政策模板ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    brand_code          VARCHAR(32),                        -- 品牌编码
    reward_name         VARCHAR(100) NOT NULL,              -- 奖励名称
    min_register_days   INT NOT NULL DEFAULT 0,             -- 最小入网天数
    max_register_days   INT NOT NULL DEFAULT 30,            -- 最大入网天数
    target_amount       BIGINT NOT NULL,                    -- 目标交易量（分）
    reward_amount       BIGINT NOT NULL,                    -- 奖励金额（分）
    reward_type         SMALLINT DEFAULT 1,                 -- 奖励类型：1-固定金额 2-交易量比例
    is_cumulative       BOOLEAN DEFAULT FALSE,              -- 是否累计（多档可叠加）
    priority            INT DEFAULT 0,                      -- 优先级（数字越大优先级越高）
    status              SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_activation_reward_policies_template ON activation_reward_policies(template_id);
CREATE INDEX idx_activation_reward_policies_channel ON activation_reward_policies(channel_id);
CREATE INDEX idx_activation_reward_policies_priority ON activation_reward_policies(priority DESC);

-- 注释
COMMENT ON TABLE activation_reward_policies IS '激活奖励政策表 - 按入网天数+交易量触发';
COMMENT ON COLUMN activation_reward_policies.min_register_days IS '最小入网天数（含）';
COMMENT ON COLUMN activation_reward_policies.max_register_days IS '最大入网天数（含）';
COMMENT ON COLUMN activation_reward_policies.target_amount IS '目标交易量（分）';
COMMENT ON COLUMN activation_reward_policies.reward_amount IS '奖励金额（分）';

-- ============================================================
-- 4. 激活奖励记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS activation_reward_records (
    id                  BIGSERIAL PRIMARY KEY,
    policy_id           BIGINT NOT NULL,                    -- 奖励政策ID
    terminal_id         BIGINT NOT NULL,                    -- 终端ID
    terminal_sn         VARCHAR(50) NOT NULL,               -- 终端SN
    merchant_id         BIGINT NOT NULL,                    -- 商户ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    agent_id            BIGINT NOT NULL,                    -- 获得奖励的代理商
    register_days       INT NOT NULL,                       -- 入网天数
    trade_amount        BIGINT NOT NULL,                    -- 达成交易量（分）
    target_amount       BIGINT NOT NULL,                    -- 目标交易量（分）
    self_reward         BIGINT NOT NULL,                    -- 自身奖励配置金额（分）
    upper_reward        BIGINT NOT NULL,                    -- 上级应返金额（用于级差计算）（分）
    actual_reward       BIGINT NOT NULL,                    -- 实际奖励金额（级差）（分）
    source_agent_id     BIGINT,                             -- 下级代理商ID（级差来源）
    wallet_type         SMALLINT DEFAULT 3,                 -- 钱包类型：3-奖励钱包
    wallet_status       SMALLINT DEFAULT 0,                 -- 0:待入账 1:已入账
    check_date          DATE NOT NULL,                      -- 检查日期
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    processed_at        TIMESTAMPTZ
);

-- 索引
CREATE INDEX idx_activation_reward_records_policy ON activation_reward_records(policy_id);
CREATE INDEX idx_activation_reward_records_terminal ON activation_reward_records(terminal_id);
CREATE INDEX idx_activation_reward_records_merchant ON activation_reward_records(merchant_id);
CREATE INDEX idx_activation_reward_records_agent ON activation_reward_records(agent_id);
CREATE INDEX idx_activation_reward_records_status ON activation_reward_records(wallet_status);
CREATE INDEX idx_activation_reward_records_date ON activation_reward_records(check_date);
CREATE UNIQUE INDEX idx_activation_reward_records_unique ON activation_reward_records(policy_id, terminal_id, check_date);

-- 注释
COMMENT ON TABLE activation_reward_records IS '激活奖励记录表 - 按级差计算';
COMMENT ON COLUMN activation_reward_records.actual_reward IS '实际奖励金额 = 自身配置 - 下级配置（级差）';
COMMENT ON COLUMN activation_reward_records.wallet_type IS '钱包类型：3-奖励钱包';

-- ============================================================
-- 5. 费率阶梯政策表（代理商调价）
-- ============================================================
CREATE TABLE IF NOT EXISTS rate_stage_policies (
    id                  BIGSERIAL PRIMARY KEY,
    template_id         BIGINT NOT NULL,                    -- 政策模板ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    brand_code          VARCHAR(32),                        -- 品牌编码
    stage_name          VARCHAR(100) NOT NULL,              -- 阶梯名称
    apply_to            SMALLINT NOT NULL DEFAULT 1,        -- 应用对象：1-商户 2-代理商
    min_days            INT NOT NULL DEFAULT 0,             -- 最小入网天数
    max_days            INT NOT NULL,                       -- 最大入网天数（-1表示无限）
    credit_rate_delta   DECIMAL(10,4) DEFAULT 0,            -- 贷记卡费率调整值
    debit_rate_delta    DECIMAL(10,4) DEFAULT 0,            -- 借记卡费率调整值
    unionpay_rate_delta DECIMAL(10,4) DEFAULT 0,            -- 云闪付费率调整值
    wechat_rate_delta   DECIMAL(10,4) DEFAULT 0,            -- 微信费率调整值
    alipay_rate_delta   DECIMAL(10,4) DEFAULT 0,            -- 支付宝费率调整值
    priority            INT DEFAULT 0,                      -- 优先级
    status              SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_rate_stage_policies_template ON rate_stage_policies(template_id);
CREATE INDEX idx_rate_stage_policies_channel ON rate_stage_policies(channel_id);
CREATE INDEX idx_rate_stage_policies_apply ON rate_stage_policies(apply_to);

-- 注释
COMMENT ON TABLE rate_stage_policies IS '费率阶梯政策表 - 按入网天数自动调整费率';
COMMENT ON COLUMN rate_stage_policies.apply_to IS '应用对象：1-商户 2-代理商';
COMMENT ON COLUMN rate_stage_policies.credit_rate_delta IS '贷记卡费率调整值（正数加费率，负数减费率）';

-- ============================================================
-- 6. 扩展 agent_policies 表
-- ============================================================
-- 添加新字段
ALTER TABLE agent_policies
    ADD COLUMN IF NOT EXISTS unionpay_rate DECIMAL(10,4),
    ADD COLUMN IF NOT EXISTS wechat_rate DECIMAL(10,4),
    ADD COLUMN IF NOT EXISTS alipay_rate DECIMAL(10,4),
    ADD COLUMN IF NOT EXISTS debit_cap DECIMAL(10,2),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ DEFAULT NOW();

-- 注释
COMMENT ON COLUMN agent_policies.unionpay_rate IS '云闪付费率';
COMMENT ON COLUMN agent_policies.wechat_rate IS '微信扫码费率';
COMMENT ON COLUMN agent_policies.alipay_rate IS '支付宝扫码费率';
COMMENT ON COLUMN agent_policies.debit_cap IS '借记卡封顶金额';

-- ============================================================
-- 7. 代理商押金返现政策表（代理商个性化配置）
-- ============================================================
CREATE TABLE IF NOT EXISTS agent_deposit_cashback_policies (
    id                  BIGSERIAL PRIMARY KEY,
    agent_id            BIGINT NOT NULL,                    -- 代理商ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    deposit_amount      BIGINT NOT NULL,                    -- 押金金额（分）
    cashback_amount     BIGINT NOT NULL,                    -- 返现金额（分）
    status              SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(agent_id, channel_id, deposit_amount)
);

CREATE INDEX idx_agent_deposit_cashback_agent ON agent_deposit_cashback_policies(agent_id);
CREATE INDEX idx_agent_deposit_cashback_channel ON agent_deposit_cashback_policies(channel_id);

COMMENT ON TABLE agent_deposit_cashback_policies IS '代理商押金返现政策表 - 代理商个性化配置';

-- ============================================================
-- 8. 代理商流量卡返现政策表（代理商个性化配置）
-- ============================================================
CREATE TABLE IF NOT EXISTS agent_sim_cashback_policies (
    id                      BIGSERIAL PRIMARY KEY,
    agent_id                BIGINT NOT NULL,                    -- 代理商ID
    channel_id              BIGINT NOT NULL,                    -- 通道ID
    brand_code              VARCHAR(32),                        -- 品牌编码
    first_time_cashback     BIGINT NOT NULL,                    -- 首次返现金额（分）
    second_time_cashback    BIGINT NOT NULL,                    -- 第2次返现金额（分）
    third_plus_cashback     BIGINT NOT NULL,                    -- 第3次及以后返现金额（分）
    status                  SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    created_at              TIMESTAMPTZ DEFAULT NOW(),
    updated_at              TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(agent_id, channel_id, brand_code)
);

CREATE INDEX idx_agent_sim_cashback_agent ON agent_sim_cashback_policies(agent_id);
CREATE INDEX idx_agent_sim_cashback_channel ON agent_sim_cashback_policies(channel_id);

COMMENT ON TABLE agent_sim_cashback_policies IS '代理商流量卡返现政策表 - 代理商个性化配置';

-- ============================================================
-- 9. 代理商激活奖励政策表（代理商个性化配置）
-- ============================================================
CREATE TABLE IF NOT EXISTS agent_activation_reward_policies (
    id                  BIGSERIAL PRIMARY KEY,
    agent_id            BIGINT NOT NULL,                    -- 代理商ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    brand_code          VARCHAR(32),                        -- 品牌编码
    reward_name         VARCHAR(100) NOT NULL,              -- 奖励名称
    min_register_days   INT NOT NULL DEFAULT 0,             -- 最小入网天数
    max_register_days   INT NOT NULL DEFAULT 30,            -- 最大入网天数
    target_amount       BIGINT NOT NULL,                    -- 目标交易量（分）
    reward_amount       BIGINT NOT NULL,                    -- 奖励金额（分）
    priority            INT DEFAULT 0,                      -- 优先级
    status              SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_agent_activation_reward_agent ON agent_activation_reward_policies(agent_id);
CREATE INDEX idx_agent_activation_reward_channel ON agent_activation_reward_policies(channel_id);

COMMENT ON TABLE agent_activation_reward_policies IS '代理商激活奖励政策表 - 代理商个性化配置';

-- ============================================================
-- 示例数据
-- ============================================================
-- 押金返现政策示例
-- INSERT INTO deposit_cashback_policies (template_id, channel_id, deposit_amount, cashback_amount)
-- VALUES
--     (1, 1, 0, 0),           -- 无押金无返现
--     (1, 1, 9900, 5000),     -- 押金99元返现50元
--     (1, 1, 19900, 12000),   -- 押金199元返现120元
--     (1, 1, 29900, 20000);   -- 押金299元返现200元

-- 激活奖励政策示例
-- INSERT INTO activation_reward_policies (template_id, channel_id, reward_name, min_register_days, max_register_days, target_amount, reward_amount)
-- VALUES
--     (1, 1, '新机激活奖励', 0, 30, 100000000, 5000),    -- 30天内交易100万奖励50元
--     (1, 1, '快速激活奖励', 0, 7, 50000000, 3000);      -- 7天内交易50万奖励30元
