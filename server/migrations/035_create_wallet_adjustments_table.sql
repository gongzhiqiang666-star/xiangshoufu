-- 035_create_wallet_adjustments_table.sql
-- 钱包调账记录表

CREATE TABLE IF NOT EXISTS wallet_adjustments (
    id BIGSERIAL PRIMARY KEY,
    adjustment_no VARCHAR(50) NOT NULL UNIQUE,           -- 调账单号
    agent_id BIGINT NOT NULL,                            -- 代理商ID
    wallet_id BIGINT NOT NULL,                           -- 钱包ID
    wallet_type SMALLINT NOT NULL,                       -- 钱包类型: 1分润 2服务费 3奖励 4充值 5沉淀
    channel_id BIGINT NOT NULL DEFAULT 0,                -- 通道ID: 0表示不区分通道
    amount BIGINT NOT NULL,                              -- 调账金额(分): 正数增加, 负数扣减
    balance_before BIGINT NOT NULL,                      -- 调账前余额
    balance_after BIGINT NOT NULL,                       -- 调账后余额
    reason VARCHAR(500) NOT NULL,                        -- 调账原因
    operator_id BIGINT NOT NULL,                         -- 操作人ID
    operator_name VARCHAR(50),                           -- 操作人名称
    status SMALLINT NOT NULL DEFAULT 1,                  -- 状态: 1已生效 2待审批 3已驳回(预留)
    approved_by BIGINT,                                  -- 审批人ID(预留)
    approved_at TIMESTAMP,                               -- 审批时间(预留)
    reject_reason VARCHAR(500),                          -- 驳回原因(预留)
    wallet_log_id BIGINT,                                -- 关联的钱包流水ID
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_wallet_adjustments_agent_id ON wallet_adjustments(agent_id);
CREATE INDEX idx_wallet_adjustments_wallet_id ON wallet_adjustments(wallet_id);
CREATE INDEX idx_wallet_adjustments_created_at ON wallet_adjustments(created_at DESC);
CREATE INDEX idx_wallet_adjustments_status ON wallet_adjustments(status);

COMMENT ON TABLE wallet_adjustments IS '钱包调账记录表';
COMMENT ON COLUMN wallet_adjustments.adjustment_no IS '调账单号';
COMMENT ON COLUMN wallet_adjustments.agent_id IS '代理商ID';
COMMENT ON COLUMN wallet_adjustments.wallet_id IS '钱包ID';
COMMENT ON COLUMN wallet_adjustments.wallet_type IS '钱包类型: 1分润 2服务费 3奖励 4充值 5沉淀';
COMMENT ON COLUMN wallet_adjustments.channel_id IS '通道ID: 0表示不区分通道';
COMMENT ON COLUMN wallet_adjustments.amount IS '调账金额(分): 正数增加, 负数扣减';
COMMENT ON COLUMN wallet_adjustments.balance_before IS '调账前余额';
COMMENT ON COLUMN wallet_adjustments.balance_after IS '调账后余额';
COMMENT ON COLUMN wallet_adjustments.reason IS '调账原因';
COMMENT ON COLUMN wallet_adjustments.operator_id IS '操作人ID';
COMMENT ON COLUMN wallet_adjustments.operator_name IS '操作人名称';
COMMENT ON COLUMN wallet_adjustments.status IS '状态: 1已生效 2待审批 3已驳回';
COMMENT ON COLUMN wallet_adjustments.wallet_log_id IS '关联的钱包流水ID';
