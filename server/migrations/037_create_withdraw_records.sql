-- 037_create_withdraw_records.sql
-- 提现记录表

CREATE TABLE IF NOT EXISTS withdraw_records (
    id BIGSERIAL PRIMARY KEY,
    withdraw_no VARCHAR(50) NOT NULL,
    agent_id BIGINT NOT NULL,
    wallet_id BIGINT NOT NULL,
    wallet_type SMALLINT NOT NULL DEFAULT 1,
    channel_id BIGINT NOT NULL DEFAULT 0,
    tax_channel_id BIGINT,

    -- 金额信息
    amount BIGINT NOT NULL DEFAULT 0,
    tax_fee BIGINT NOT NULL DEFAULT 0,
    fixed_fee BIGINT NOT NULL DEFAULT 0,
    actual_amount BIGINT NOT NULL DEFAULT 0,

    -- 结算卡信息
    bank_name VARCHAR(100),
    bank_account VARCHAR(50),
    account_name VARCHAR(50),

    -- 状态信息
    status SMALLINT NOT NULL DEFAULT 0,
    reject_reason VARCHAR(500),
    fail_reason VARCHAR(500),

    -- 审核信息
    audited_by BIGINT,
    audited_at TIMESTAMP,
    audit_remark VARCHAR(500),

    -- 打款信息
    paid_at TIMESTAMP,
    paid_ref VARCHAR(100),
    paid_remark VARCHAR(500),

    -- 时间戳
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE UNIQUE INDEX IF NOT EXISTS idx_withdraw_records_withdraw_no ON withdraw_records(withdraw_no);
CREATE INDEX IF NOT EXISTS idx_withdraw_records_agent_id ON withdraw_records(agent_id);
CREATE INDEX IF NOT EXISTS idx_withdraw_records_wallet_id ON withdraw_records(wallet_id);
CREATE INDEX IF NOT EXISTS idx_withdraw_records_status ON withdraw_records(status);
CREATE INDEX IF NOT EXISTS idx_withdraw_records_created_at ON withdraw_records(created_at);

-- 注释
COMMENT ON TABLE withdraw_records IS '提现记录表';
COMMENT ON COLUMN withdraw_records.withdraw_no IS '提现单号';
COMMENT ON COLUMN withdraw_records.agent_id IS '代理商ID';
COMMENT ON COLUMN withdraw_records.wallet_id IS '钱包ID';
COMMENT ON COLUMN withdraw_records.wallet_type IS '钱包类型: 1-分润钱包 2-服务费钱包 3-奖励钱包';
COMMENT ON COLUMN withdraw_records.channel_id IS '支付通道ID';
COMMENT ON COLUMN withdraw_records.tax_channel_id IS '税筹通道ID';
COMMENT ON COLUMN withdraw_records.amount IS '提现金额(分)';
COMMENT ON COLUMN withdraw_records.tax_fee IS '税费(分)';
COMMENT ON COLUMN withdraw_records.fixed_fee IS '固定手续费(分)';
COMMENT ON COLUMN withdraw_records.actual_amount IS '实际到账金额(分)';
COMMENT ON COLUMN withdraw_records.status IS '状态: 0-待审核 1-已审核 2-已打款 3-已拒绝 4-打款失败 5-已取消';
