-- 代扣管理模块重构：添加冻结机制
-- 创建时间: 2026-01-27
-- 业务规则:
--   1. 合并货款代扣和伙伴代扣为统一的代扣管理
--   2. 扣款时机改为定时扣款（每天8:00）
--   3. 新增冻结机制：接收确认后开始冻结，入账时继续冻结
--   4. 冻结上限：冻结金额 ≤ 剩余待扣金额

-- 在 deduction_plans 表新增字段
ALTER TABLE deduction_plans ADD COLUMN IF NOT EXISTS need_accept BOOLEAN DEFAULT FALSE;
ALTER TABLE deduction_plans ADD COLUMN IF NOT EXISTS accepted_at TIMESTAMPTZ;
ALTER TABLE deduction_plans ADD COLUMN IF NOT EXISTS frozen_amount BIGINT DEFAULT 0;
ALTER TABLE deduction_plans ADD COLUMN IF NOT EXISTS deduction_source SMALLINT DEFAULT 3;

-- 更新状态默认值（0=待接收）
COMMENT ON COLUMN deduction_plans.status IS '状态：0待接收 1进行中 2已完成 3已暂停 4已取消 5已拒绝';
COMMENT ON COLUMN deduction_plans.need_accept IS '是否需要接收确认';
COMMENT ON COLUMN deduction_plans.accepted_at IS '接收确认时间';
COMMENT ON COLUMN deduction_plans.frozen_amount IS '已冻结金额（分）';
COMMENT ON COLUMN deduction_plans.deduction_source IS '扣款来源：1分润钱包 2服务费钱包 3两者都扣';

-- 代扣冻结明细表
CREATE TABLE IF NOT EXISTS deduction_freeze_logs (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES deduction_plans(id),
    agent_id BIGINT NOT NULL,                             -- 被扣款方代理商ID
    wallet_id BIGINT NOT NULL,                            -- 钱包ID
    wallet_type SMALLINT NOT NULL,                        -- 钱包类型
    channel_id BIGINT,                                    -- 通道ID
    freeze_amount BIGINT NOT NULL,                        -- 本次冻结金额（分）
    total_frozen BIGINT NOT NULL,                         -- 累计冻结金额（分）
    trigger_type VARCHAR(32),                             -- 触发类型: accept/income
    trigger_ref_id BIGINT,                                -- 触发来源ID
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_deduction_freeze_logs_plan ON deduction_freeze_logs(plan_id);
CREATE INDEX IF NOT EXISTS idx_deduction_freeze_logs_agent ON deduction_freeze_logs(agent_id);
CREATE INDEX IF NOT EXISTS idx_deduction_freeze_logs_created ON deduction_freeze_logs(created_at);

-- 添加注释
COMMENT ON TABLE deduction_freeze_logs IS '代扣冻结明细表 - 记录每次冻结操作的详情';
COMMENT ON COLUMN deduction_freeze_logs.trigger_type IS '触发类型：accept接收确认时冻结 income入账时冻结';

-- 在钱包表添加冻结金额字段（如果不存在）
ALTER TABLE wallets ADD COLUMN IF NOT EXISTS frozen_amount BIGINT DEFAULT 0;
COMMENT ON COLUMN wallets.frozen_amount IS '冻结金额（分），用于代扣冻结';

-- 创建索引优化查询
CREATE INDEX IF NOT EXISTS idx_deduction_plans_need_accept ON deduction_plans(need_accept) WHERE need_accept = TRUE;
CREATE INDEX IF NOT EXISTS idx_deduction_plans_frozen ON deduction_plans(frozen_amount) WHERE frozen_amount > 0;
