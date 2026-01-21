-- 终端政策设置表
-- 用于存储每个终端的费率、SIM卡费用、押金等政策配置

CREATE TABLE IF NOT EXISTS terminal_policies (
    id BIGSERIAL PRIMARY KEY,
    terminal_sn VARCHAR(50) NOT NULL UNIQUE,
    channel_id BIGINT NOT NULL,
    agent_id BIGINT NOT NULL,

    -- 费率设置 (万分比)
    credit_rate INT DEFAULT 0,           -- 贷记卡费率，如 55 表示万分之55
    debit_rate INT DEFAULT 0,            -- 借记卡费率
    debit_cap INT DEFAULT 0,             -- 借记卡封顶(分)
    unionpay_rate INT DEFAULT 0,         -- 银联云闪付费率
    wechat_rate INT DEFAULT 0,           -- 微信扫码费率
    alipay_rate INT DEFAULT 0,           -- 支付宝扫码费率

    -- SIM卡费用设置 (分)
    first_sim_fee INT DEFAULT 0,         -- 首次流量费
    non_first_sim_fee INT DEFAULT 0,     -- 非首次流量费
    sim_fee_interval_days INT DEFAULT 0, -- 非首次间隔天数

    -- 押金设置 (分)
    deposit_amount INT DEFAULT 0,        -- 押金金额，0表示无押金

    -- 同步状态
    is_synced BOOLEAN DEFAULT FALSE,     -- 是否已同步到通道
    synced_at TIMESTAMP,                 -- 同步时间
    sync_error VARCHAR(500),             -- 同步错误信息

    -- 审计字段
    created_by BIGINT,
    updated_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_terminal_policies_channel_id ON terminal_policies(channel_id);
CREATE INDEX IF NOT EXISTS idx_terminal_policies_agent_id ON terminal_policies(agent_id);
CREATE INDEX IF NOT EXISTS idx_terminal_policies_is_synced ON terminal_policies(is_synced);

-- 注释
COMMENT ON TABLE terminal_policies IS '终端政策设置表';
COMMENT ON COLUMN terminal_policies.terminal_sn IS '终端序列号';
COMMENT ON COLUMN terminal_policies.channel_id IS '所属通道ID';
COMMENT ON COLUMN terminal_policies.agent_id IS '所属代理商ID';
COMMENT ON COLUMN terminal_policies.credit_rate IS '贷记卡费率(万分比)';
COMMENT ON COLUMN terminal_policies.debit_rate IS '借记卡费率(万分比)';
COMMENT ON COLUMN terminal_policies.debit_cap IS '借记卡封顶金额(分)';
COMMENT ON COLUMN terminal_policies.unionpay_rate IS '银联云闪付费率(万分比)';
COMMENT ON COLUMN terminal_policies.wechat_rate IS '微信扫码费率(万分比)';
COMMENT ON COLUMN terminal_policies.alipay_rate IS '支付宝扫码费率(万分比)';
COMMENT ON COLUMN terminal_policies.first_sim_fee IS '首次流量费(分)';
COMMENT ON COLUMN terminal_policies.non_first_sim_fee IS '非首次流量费(分)';
COMMENT ON COLUMN terminal_policies.sim_fee_interval_days IS '非首次流量费间隔天数';
COMMENT ON COLUMN terminal_policies.deposit_amount IS '押金金额(分)，0表示无押金';
COMMENT ON COLUMN terminal_policies.is_synced IS '是否已同步到支付通道';
COMMENT ON COLUMN terminal_policies.synced_at IS '最后同步时间';
COMMENT ON COLUMN terminal_policies.sync_error IS '同步错误信息';
