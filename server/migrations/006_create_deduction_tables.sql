-- 代扣管理相关表
-- 创建时间: 2026-01-18
-- 业务规则:
--   Q6: 伙伴代扣 - 任意代理商之间都可以发起（不限层级关系）
--   Q7: 每日扣款 - 每天固定时间检查余额并扣款
--   Q8: 多钱包扣款优先级 - 按余额从多到少扣

-- 代扣计划表
CREATE TABLE IF NOT EXISTS deduction_plans (
    id BIGSERIAL PRIMARY KEY,
    plan_no VARCHAR(64) NOT NULL UNIQUE,           -- 计划编号
    deductor_id BIGINT NOT NULL,                   -- 扣款方代理商ID
    deductee_id BIGINT NOT NULL,                   -- 被扣款方代理商ID
    plan_type SMALLINT NOT NULL,                   -- 1:货款代扣 2:伙伴代扣 3:押金代扣
    total_amount BIGINT NOT NULL,                  -- 总金额（分）
    deducted_amount BIGINT DEFAULT 0,              -- 已扣金额（分）
    remaining_amount BIGINT NOT NULL,              -- 剩余金额（分）
    total_periods INT NOT NULL,                    -- 总期数
    current_period INT DEFAULT 0,                  -- 当前期数
    period_amount BIGINT NOT NULL,                 -- 每期金额（分）
    status SMALLINT DEFAULT 1,                     -- 1:进行中 2:已完成 3:已暂停 4:已取消
    related_type VARCHAR(32),                      -- 关联类型: terminal_distribute, partner_loan
    related_id BIGINT,                             -- 关联ID
    remark VARCHAR(255),                           -- 备注
    created_by BIGINT,                             -- 创建人
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_deduction_plans_deductor ON deduction_plans(deductor_id);
CREATE INDEX idx_deduction_plans_deductee ON deduction_plans(deductee_id);
CREATE INDEX idx_deduction_plans_status ON deduction_plans(status);
CREATE INDEX idx_deduction_plans_related ON deduction_plans(related_type, related_id);

-- 代扣记录表
CREATE TABLE IF NOT EXISTS deduction_records (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL,                       -- 代扣计划ID
    plan_no VARCHAR(64),                           -- 计划编号
    deductor_id BIGINT NOT NULL,                   -- 扣款方
    deductee_id BIGINT NOT NULL,                   -- 被扣款方
    period_num INT NOT NULL,                       -- 期数
    amount BIGINT NOT NULL,                        -- 应扣金额（分）
    actual_amount BIGINT DEFAULT 0,                -- 实扣金额（分）
    status SMALLINT DEFAULT 0,                     -- 0:待扣款 1:成功 2:部分成功 3:失败
    wallet_details JSONB,                          -- 钱包扣款明细JSON
    fail_reason VARCHAR(255),                      -- 失败原因
    scheduled_at TIMESTAMP NOT NULL,               -- 计划扣款时间
    deducted_at TIMESTAMP,                         -- 实际扣款时间
    created_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_deduction_records_plan ON deduction_records(plan_id);
CREATE INDEX idx_deduction_records_deductee ON deduction_records(deductee_id);
CREATE INDEX idx_deduction_records_scheduled ON deduction_records(scheduled_at);
CREATE INDEX idx_deduction_records_status ON deduction_records(status);

-- 代扣链表（用于跨级下发）
CREATE TABLE IF NOT EXISTS deduction_chains (
    id BIGSERIAL PRIMARY KEY,
    chain_no VARCHAR(64) NOT NULL UNIQUE,          -- 代扣链编号
    distribute_id BIGINT NOT NULL,                 -- 终端下发记录ID
    terminal_sn VARCHAR(50),                       -- 终端SN
    total_levels INT NOT NULL,                     -- 总层级数
    total_amount BIGINT NOT NULL,                  -- 总金额
    status SMALLINT DEFAULT 1,                     -- 1:进行中 2:已完成 3:已取消
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_deduction_chains_distribute ON deduction_chains(distribute_id);
CREATE INDEX idx_deduction_chains_terminal ON deduction_chains(terminal_sn);

-- 代扣链节点表
CREATE TABLE IF NOT EXISTS deduction_chain_items (
    id BIGSERIAL PRIMARY KEY,
    chain_id BIGINT NOT NULL,                      -- 所属代扣链
    chain_no VARCHAR(64),                          -- 代扣链编号
    level INT NOT NULL,                            -- 层级（1,2,3...）
    from_agent_id BIGINT NOT NULL,                 -- 扣款方
    to_agent_id BIGINT NOT NULL,                   -- 收款方
    plan_id BIGINT,                                -- 关联的代扣计划ID
    amount BIGINT NOT NULL,                        -- 代扣金额
    status SMALLINT DEFAULT 0                      -- 0:待处理 1:已生成计划 2:已完成
);

-- 索引
CREATE INDEX idx_deduction_chain_items_chain ON deduction_chain_items(chain_id);
CREATE INDEX idx_deduction_chain_items_from ON deduction_chain_items(from_agent_id);
CREATE INDEX idx_deduction_chain_items_to ON deduction_chain_items(to_agent_id);

-- 添加注释
COMMENT ON TABLE deduction_plans IS '代扣计划表 - 支持伙伴代扣（任意代理商间）';
COMMENT ON COLUMN deduction_plans.plan_type IS '计划类型：1货款代扣 2伙伴代扣 3押金代扣';
COMMENT ON COLUMN deduction_plans.status IS '状态：1进行中 2已完成 3已暂停 4已取消';

COMMENT ON TABLE deduction_records IS '代扣记录表 - 每日扣款记录';
COMMENT ON COLUMN deduction_records.status IS '状态：0待扣款 1成功 2部分成功 3失败';
COMMENT ON COLUMN deduction_records.wallet_details IS '多钱包扣款明细，按余额从多到少扣';

COMMENT ON TABLE deduction_chains IS '代扣链表 - 跨级下发时自动生成A→B→C的代扣链';
COMMENT ON TABLE deduction_chain_items IS '代扣链节点表 - 每个节点对应一个代扣计划';
