-- 核心业务表 - 基于代理商分润管理系统完整设计方案

-- 支付通道表
CREATE TABLE IF NOT EXISTS channels (
    id              BIGSERIAL PRIMARY KEY,
    channel_code    VARCHAR(32) NOT NULL UNIQUE,
    channel_name    VARCHAR(100) NOT NULL,
    credit_rate     DECIMAL(10,4),
    debit_rate      DECIMAL(10,4),
    debit_cap       DECIMAL(10,2),
    unionpay_rate   DECIMAL(10,4),
    wechat_rate     DECIMAL(10,4),
    alipay_rate     DECIMAL(10,4),
    api_url         VARCHAR(255),
    api_key         VARCHAR(255),
    api_secret      VARCHAR(255),
    status          SMALLINT DEFAULT 1,
    is_visible      BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 代理商表
CREATE TABLE IF NOT EXISTS agents (
    id              BIGSERIAL PRIMARY KEY,
    agent_no        VARCHAR(32) NOT NULL UNIQUE,
    agent_name      VARCHAR(100) NOT NULL,
    parent_id       BIGINT REFERENCES agents(id),
    path            VARCHAR(500) DEFAULT '',
    level           INT DEFAULT 1,
    default_rate    DECIMAL(10,4),
    contact_name    VARCHAR(50),
    contact_phone   VARCHAR(20) NOT NULL,
    id_card_no      VARCHAR(18),
    bank_name       VARCHAR(100),
    bank_account    VARCHAR(30),
    bank_card_no    VARCHAR(25),
    invite_code     VARCHAR(20) UNIQUE,
    qr_code_url     VARCHAR(255),
    status          SMALLINT DEFAULT 1,
    register_time   TIMESTAMPTZ DEFAULT NOW(),
    direct_agent_count   INT DEFAULT 0,
    direct_merchant_count INT DEFAULT 0,
    team_agent_count     INT DEFAULT 0,
    team_merchant_count  INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agents_parent ON agents(parent_id);
CREATE INDEX IF NOT EXISTS idx_agents_path ON agents(path);

-- 代理商政策表
CREATE TABLE IF NOT EXISTS agent_policies (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    template_id     BIGINT NOT NULL,
    credit_rate     DECIMAL(10,4),
    debit_rate      DECIMAL(10,4),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(agent_id, channel_id)
);

-- 商户表
CREATE TABLE IF NOT EXISTS merchants (
    id              BIGSERIAL PRIMARY KEY,
    merchant_no     VARCHAR(64) NOT NULL UNIQUE,
    merchant_name   VARCHAR(100),
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    terminal_sn     VARCHAR(50),
    status          SMALLINT DEFAULT 1,
    approve_status  SMALLINT DEFAULT 1,
    legal_name      VARCHAR(50),
    legal_id_card   VARCHAR(18),
    mcc             VARCHAR(10),
    credit_rate     DECIMAL(10,4),
    debit_rate      DECIMAL(10,4),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_merchants_agent ON merchants(agent_id);
CREATE INDEX IF NOT EXISTS idx_merchants_terminal ON merchants(terminal_sn);

-- 终端表
CREATE TABLE IF NOT EXISTS terminals (
    id              BIGSERIAL PRIMARY KEY,
    sn              VARCHAR(50) NOT NULL UNIQUE,
    terminal_no     VARCHAR(20),
    channel_id      BIGINT NOT NULL,
    owner_agent_id  BIGINT,
    merchant_id     BIGINT,
    credit_rate     DECIMAL(10,4),
    debit_rate      DECIMAL(10,4),
    deposit_type    SMALLINT DEFAULT 0,
    deposit_amount  DECIMAL(10,2) DEFAULT 0,
    sim_first_fee   DECIMAL(10,2),
    sim_next_fee    DECIMAL(10,2),
    sim_interval    INT,
    status          SMALLINT DEFAULT 0,
    dispatch_time   TIMESTAMPTZ,
    bind_time       TIMESTAMPTZ,
    activate_time   TIMESTAMPTZ,
    first_trade_time TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_terminals_owner ON terminals(owner_agent_id);
CREATE INDEX IF NOT EXISTS idx_terminals_merchant ON terminals(merchant_id);

-- 交易流水表
CREATE TABLE IF NOT EXISTS transactions (
    id              BIGSERIAL PRIMARY KEY,
    trade_no        VARCHAR(64),
    order_no        VARCHAR(64) NOT NULL UNIQUE,
    channel_id      BIGINT NOT NULL,
    channel_code    VARCHAR(32),
    terminal_sn     VARCHAR(50) NOT NULL,
    merchant_id     BIGINT,
    agent_id        BIGINT NOT NULL,
    trade_type      SMALLINT NOT NULL,
    pay_type        SMALLINT NOT NULL,
    card_type       SMALLINT,
    amount          BIGINT NOT NULL,
    fee             BIGINT,
    rate            VARCHAR(20),
    d0_fee          BIGINT DEFAULT 0,
    high_rate       DECIMAL(10,4),
    card_no         VARCHAR(32),
    profit_status   SMALLINT DEFAULT 0,
    refund_status   SMALLINT DEFAULT 0,
    trade_time      TIMESTAMPTZ NOT NULL,
    received_at     TIMESTAMPTZ DEFAULT NOW(),
    brand_code      VARCHAR(32),
    ext_data        JSONB
);

CREATE INDEX IF NOT EXISTS idx_transactions_agent ON transactions(agent_id);
CREATE INDEX IF NOT EXISTS idx_transactions_merchant ON transactions(merchant_id);
CREATE INDEX IF NOT EXISTS idx_transactions_terminal ON transactions(terminal_sn);
CREATE INDEX IF NOT EXISTS idx_transactions_profit ON transactions(profit_status);
CREATE INDEX IF NOT EXISTS idx_transactions_time ON transactions(trade_time);

-- 分润明细表
CREATE TABLE IF NOT EXISTS profit_records (
    id              BIGSERIAL PRIMARY KEY,
    transaction_id  BIGINT NOT NULL,
    order_no        VARCHAR(64) NOT NULL,
    agent_id        BIGINT NOT NULL,
    profit_type     SMALLINT NOT NULL,
    trade_amount    BIGINT NOT NULL,
    self_rate       VARCHAR(20),
    lower_rate      VARCHAR(20),
    rate_diff       VARCHAR(20),
    profit_amount   BIGINT NOT NULL,
    source_merchant_id BIGINT,
    source_agent_id    BIGINT,
    channel_id         BIGINT,
    wallet_type     SMALLINT NOT NULL,
    wallet_status   SMALLINT DEFAULT 0,
    is_revoked      BOOLEAN DEFAULT FALSE,
    revoked_at      TIMESTAMPTZ,
    revoke_reason   VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_profit_records_tx ON profit_records(transaction_id);
CREATE INDEX IF NOT EXISTS idx_profit_records_agent ON profit_records(agent_id);
CREATE INDEX IF NOT EXISTS idx_profit_records_revoked ON profit_records(is_revoked) WHERE is_revoked = TRUE;

-- 钱包表
CREATE TABLE IF NOT EXISTS wallets (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    channel_id      BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,
    balance         BIGINT DEFAULT 0,
    frozen_amount   BIGINT DEFAULT 0,
    total_income    BIGINT DEFAULT 0,
    total_withdraw  BIGINT DEFAULT 0,
    withdraw_threshold BIGINT DEFAULT 10000,
    version         INT DEFAULT 0,
    updated_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(agent_id, channel_id, wallet_type)
);

-- 钱包流水表
CREATE TABLE IF NOT EXISTS wallet_logs (
    id              BIGSERIAL PRIMARY KEY,
    wallet_id       BIGINT NOT NULL,
    agent_id        BIGINT NOT NULL,
    wallet_type     SMALLINT NOT NULL,
    log_type        SMALLINT NOT NULL,
    amount          BIGINT NOT NULL,
    balance_before  BIGINT NOT NULL,
    balance_after   BIGINT NOT NULL,
    ref_type        VARCHAR(20),
    ref_id          BIGINT,
    remark          VARCHAR(500),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_wallet_logs_wallet ON wallet_logs(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallet_logs_agent ON wallet_logs(agent_id);

-- 政策模板表
CREATE TABLE IF NOT EXISTS policy_templates (
    id              BIGSERIAL PRIMARY KEY,
    template_name   VARCHAR(100) NOT NULL,
    channel_id      BIGINT NOT NULL,
    is_default      BOOLEAN DEFAULT FALSE,
    credit_rate     DECIMAL(10,4),
    debit_rate      DECIMAL(10,4),
    debit_cap       DECIMAL(10,2),
    unionpay_rate   DECIMAL(10,4),
    wechat_rate     DECIMAL(10,4),
    alipay_rate     DECIMAL(10,4),
    status          SMALLINT DEFAULT 1,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 插入默认通道数据
INSERT INTO channels (channel_code, channel_name, status) VALUES
    ('HENGXINTONG', '恒信通', 1),
    ('LAKALA', '拉卡拉', 1),
    ('YEAHKA', '乐刷', 1),
    ('SUIXINGFU', '随行付', 1),
    ('LIANLIAN', '连连支付', 1),
    ('SANDPAY', '杉德支付', 1),
    ('FUIOU', '富友支付', 1),
    ('HEEPAY', '汇付天下', 1)
ON CONFLICT (channel_code) DO NOTHING;

-- 添加表注释
COMMENT ON TABLE channels IS '支付通道表';
COMMENT ON TABLE agents IS '代理商表';
COMMENT ON TABLE merchants IS '商户表';
COMMENT ON TABLE terminals IS '终端/机具表';
COMMENT ON TABLE transactions IS '交易流水表';
COMMENT ON TABLE profit_records IS '分润明细表';
COMMENT ON TABLE wallets IS '钱包表';
COMMENT ON TABLE wallet_logs IS '钱包流水表';
