-- 奖励模块重构表结构
-- 创建时间: 2026-01-25
-- 业务规则:
--   1. 奖励与通道解耦，在政策模版中单独维护
--   2. 奖励按"阶段+维度"两个维度配置
--   3. 阶段：按天数或按自然月
--   4. 维度：按金额或按笔数
--   5. 级差奖励：固定池分配模式
--   6. 政策快照：激活时绑定政策，全程使用

-- ============================================================
-- 1. 奖励政策模版表
-- ============================================================
CREATE TABLE IF NOT EXISTS reward_policy_templates (
    id                  BIGSERIAL PRIMARY KEY,
    name                VARCHAR(100) NOT NULL,              -- 模版名称
    time_type           VARCHAR(20) NOT NULL,               -- 时间类型：'days' 或 'months'
    dimension_type      VARCHAR(20) NOT NULL,               -- 维度类型：'amount' 或 'count'
    trans_types         VARCHAR(100),                       -- 交易类型，逗号分隔：'scan,debit,credit'
    amount_min          BIGINT DEFAULT 0,                   -- 交易金额下限（分）（筛选条件）
    amount_max          BIGINT,                             -- 交易金额上限（分）（筛选条件，NULL表示无上限）
    allow_gap           BOOLEAN DEFAULT FALSE,              -- 断档开关：是否允许断档
    enabled             BOOLEAN DEFAULT TRUE,               -- 是否启用
    description         TEXT,                               -- 描述说明
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_reward_policy_templates_enabled ON reward_policy_templates(enabled);
CREATE INDEX idx_reward_policy_templates_created ON reward_policy_templates(created_at);

-- 注释
COMMENT ON TABLE reward_policy_templates IS '奖励政策模版表 - 与通道解耦的独立奖励配置';
COMMENT ON COLUMN reward_policy_templates.time_type IS '时间类型：days-按天数，months-按自然月';
COMMENT ON COLUMN reward_policy_templates.dimension_type IS '维度类型：amount-按金额，count-按笔数';
COMMENT ON COLUMN reward_policy_templates.trans_types IS '参与计算的交易类型，逗号分隔';
COMMENT ON COLUMN reward_policy_templates.amount_min IS '交易金额下限（分），用于筛选参与奖励计算的交易';
COMMENT ON COLUMN reward_policy_templates.amount_max IS '交易金额上限（分），NULL表示无上限';
COMMENT ON COLUMN reward_policy_templates.allow_gap IS '断档开关：true-允许断档（跳过未达标阶段继续），false-不允许';

-- ============================================================
-- 2. 奖励阶段配置表
-- ============================================================
CREATE TABLE IF NOT EXISTS reward_stages (
    id                  BIGSERIAL PRIMARY KEY,
    template_id         BIGINT NOT NULL REFERENCES reward_policy_templates(id) ON DELETE CASCADE,
    stage_order         INT NOT NULL,                       -- 阶段顺序（从1开始）
    start_value         INT NOT NULL,                       -- 开始值（天数或月份）
    end_value           INT NOT NULL,                       -- 结束值（天数或月份）
    target_value        BIGINT NOT NULL,                    -- 达标值（金额分或笔数）
    reward_amount       BIGINT NOT NULL,                    -- 奖励金额（分）
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_reward_stages_order CHECK (stage_order > 0),
    CONSTRAINT chk_reward_stages_values CHECK (end_value >= start_value),
    CONSTRAINT chk_reward_stages_target CHECK (target_value > 0),
    CONSTRAINT chk_reward_stages_reward CHECK (reward_amount > 0)
);

-- 索引
CREATE INDEX idx_reward_stages_template ON reward_stages(template_id);
CREATE UNIQUE INDEX idx_reward_stages_template_order ON reward_stages(template_id, stage_order);

-- 注释
COMMENT ON TABLE reward_stages IS '奖励阶段配置表 - 定义每个阶段的时间范围和达标条件';
COMMENT ON COLUMN reward_stages.stage_order IS '阶段顺序，从1开始递增';
COMMENT ON COLUMN reward_stages.start_value IS '开始值：按天数时表示第N天，按月时表示第N月';
COMMENT ON COLUMN reward_stages.end_value IS '结束值：按天数时表示第N天，按月时表示第N月';
COMMENT ON COLUMN reward_stages.target_value IS '达标值：按金额时单位为分，按笔数时为交易笔数';
COMMENT ON COLUMN reward_stages.reward_amount IS '奖励金额（分）';

-- ============================================================
-- 3. 代理商奖励比例配置表
-- ============================================================
CREATE TABLE IF NOT EXISTS agent_reward_rates (
    id                  BIGSERIAL PRIMARY KEY,
    agent_id            BIGINT NOT NULL,                    -- 代理商ID
    reward_rate         DECIMAL(5,4) NOT NULL,              -- 奖励比例，如0.1000表示10%
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_agent_reward_rates_rate CHECK (reward_rate >= 0 AND reward_rate <= 1)
);

-- 索引
CREATE UNIQUE INDEX idx_agent_reward_rates_agent ON agent_reward_rates(agent_id);

-- 注释
COMMENT ON TABLE agent_reward_rates IS '代理商奖励比例配置表 - 固定池分配模式';
COMMENT ON COLUMN agent_reward_rates.reward_rate IS '奖励比例，0.1000表示10%，链上所有比例之和不能超过100%';

-- ============================================================
-- 4. 终端奖励进度表
-- ============================================================
CREATE TABLE IF NOT EXISTS terminal_reward_progress (
    id                  BIGSERIAL PRIMARY KEY,
    terminal_sn         VARCHAR(50) NOT NULL,               -- 终端SN
    terminal_id         BIGINT,                             -- 终端ID（可选）
    template_id         BIGINT NOT NULL REFERENCES reward_policy_templates(id),
    template_snapshot   JSONB NOT NULL,                     -- 激活时的政策快照（用旧政策）
    bind_agent_id       BIGINT NOT NULL,                    -- 绑定时的代理商ID（终端转移归原代理商）
    bind_time           TIMESTAMPTZ NOT NULL,               -- 绑定时间
    current_stage       INT DEFAULT 1,                      -- 当前阶段
    last_achieved_stage INT DEFAULT 0,                      -- 最后达标阶段
    status              VARCHAR(20) DEFAULT 'active',       -- 状态：active/completed/terminated
    completed_at        TIMESTAMPTZ,                        -- 完成时间
    terminated_at       TIMESTAMPTZ,                        -- 终止时间（中途解绑）
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_terminal_reward_progress_sn ON terminal_reward_progress(terminal_sn);
CREATE INDEX idx_terminal_reward_progress_template ON terminal_reward_progress(template_id);
CREATE INDEX idx_terminal_reward_progress_agent ON terminal_reward_progress(bind_agent_id);
CREATE INDEX idx_terminal_reward_progress_bind_time ON terminal_reward_progress(bind_time);
CREATE INDEX idx_terminal_reward_progress_status ON terminal_reward_progress(status);
CREATE INDEX idx_terminal_reward_progress_status_bind ON terminal_reward_progress(status, bind_time);
CREATE UNIQUE INDEX idx_terminal_reward_progress_active ON terminal_reward_progress(terminal_sn, template_id) WHERE status = 'active';

-- 注释
COMMENT ON TABLE terminal_reward_progress IS '终端奖励进度表 - 跟踪终端在奖励模版中的进度';
COMMENT ON COLUMN terminal_reward_progress.template_snapshot IS '激活时的政策快照，政策变更时仍用旧政策';
COMMENT ON COLUMN terminal_reward_progress.bind_agent_id IS '绑定时的代理商ID，终端转移时奖励归原代理商';
COMMENT ON COLUMN terminal_reward_progress.current_stage IS '当前正在进行的阶段';
COMMENT ON COLUMN terminal_reward_progress.last_achieved_stage IS '最后成功达标的阶段';
COMMENT ON COLUMN terminal_reward_progress.status IS 'active-进行中，completed-已完成所有阶段，terminated-中途终止';

-- ============================================================
-- 5. 终端阶段奖励记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS terminal_stage_rewards (
    id                  BIGSERIAL PRIMARY KEY,
    progress_id         BIGINT NOT NULL REFERENCES terminal_reward_progress(id) ON DELETE CASCADE,
    terminal_sn         VARCHAR(50) NOT NULL,               -- 终端SN
    stage_order         INT NOT NULL,                       -- 阶段顺序
    stage_start         TIMESTAMPTZ NOT NULL,               -- 阶段开始时间
    stage_end           TIMESTAMPTZ NOT NULL,               -- 阶段结束时间
    target_value        BIGINT NOT NULL,                    -- 目标值（金额分或笔数）
    actual_value        BIGINT NOT NULL DEFAULT 0,          -- 实际值（金额分或笔数）
    is_achieved         BOOLEAN NOT NULL DEFAULT FALSE,     -- 是否达标
    reward_amount       BIGINT,                             -- 应发奖励金额（分）
    status              VARCHAR(20) NOT NULL DEFAULT 'pending', -- 状态
    gap_blocked         BOOLEAN DEFAULT FALSE,              -- 是否被断档阻断
    settled_at          TIMESTAMPTZ,                        -- 结算时间
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_terminal_stage_rewards_progress ON terminal_stage_rewards(progress_id);
CREATE INDEX idx_terminal_stage_rewards_sn ON terminal_stage_rewards(terminal_sn);
CREATE INDEX idx_terminal_stage_rewards_status ON terminal_stage_rewards(status);
CREATE INDEX idx_terminal_stage_rewards_stage_end ON terminal_stage_rewards(stage_end);
CREATE UNIQUE INDEX idx_terminal_stage_rewards_unique ON terminal_stage_rewards(progress_id, stage_order);

-- 注释
COMMENT ON TABLE terminal_stage_rewards IS '终端阶段奖励记录表 - 记录每个阶段的达标情况';
COMMENT ON COLUMN terminal_stage_rewards.status IS 'pending-待检查，achieved-已达标，failed-未达标，gap_blocked-被断档阻断，settled-已结算';
COMMENT ON COLUMN terminal_stage_rewards.gap_blocked IS '断档阻断：当allow_gap=false时，前一阶段未达标会阻断后续阶段';

-- ============================================================
-- 6. 奖励发放记录表（级差分配）
-- ============================================================
CREATE TABLE IF NOT EXISTS reward_distributions (
    id                  BIGSERIAL PRIMARY KEY,
    stage_reward_id     BIGINT NOT NULL REFERENCES terminal_stage_rewards(id) ON DELETE CASCADE,
    terminal_sn         VARCHAR(50) NOT NULL,               -- 终端SN
    agent_id            BIGINT NOT NULL,                    -- 代理商ID
    agent_level         INT NOT NULL,                       -- 层级：1=终端归属，2=上级，3=上上级...
    reward_rate         DECIMAL(5,4) NOT NULL,              -- 奖励比例
    reward_amount       BIGINT NOT NULL,                    -- 奖励金额（分）
    wallet_record_id    BIGINT,                             -- 关联钱包记录ID
    wallet_status       SMALLINT DEFAULT 0,                 -- 钱包状态：0-待入账，1-已入账
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_reward_distributions_stage ON reward_distributions(stage_reward_id);
CREATE INDEX idx_reward_distributions_terminal ON reward_distributions(terminal_sn);
CREATE INDEX idx_reward_distributions_agent ON reward_distributions(agent_id);
CREATE INDEX idx_reward_distributions_wallet_status ON reward_distributions(wallet_status);

-- 注释
COMMENT ON TABLE reward_distributions IS '奖励发放记录表 - 固定池分配模式的级差分配记录';
COMMENT ON COLUMN reward_distributions.agent_level IS '代理商层级：1=终端归属代理商，2=上级，3=上上级...';
COMMENT ON COLUMN reward_distributions.reward_rate IS '该代理商配置的奖励比例';
COMMENT ON COLUMN reward_distributions.reward_amount IS '实际获得的奖励金额（分）';

-- ============================================================
-- 7. 奖励池溢出异常日志表
-- ============================================================
CREATE TABLE IF NOT EXISTS reward_overflow_logs (
    id                  BIGSERIAL PRIMARY KEY,
    terminal_sn         VARCHAR(50) NOT NULL,               -- 终端SN
    stage_reward_id     BIGINT,                             -- 关联阶段奖励记录
    agent_chain         JSONB NOT NULL,                     -- 代理商链及比例配置
    total_rate          DECIMAL(5,4) NOT NULL,              -- 总比例（超过1.0000）
    reward_amount       BIGINT NOT NULL,                    -- 原应发奖励金额（分）
    error_message       TEXT,                               -- 错误信息
    resolved            BOOLEAN DEFAULT FALSE,              -- 是否已解决
    resolved_at         TIMESTAMPTZ,                        -- 解决时间
    resolved_by         VARCHAR(50),                        -- 解决人
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_reward_overflow_logs_terminal ON reward_overflow_logs(terminal_sn);
CREATE INDEX idx_reward_overflow_logs_resolved ON reward_overflow_logs(resolved);
CREATE INDEX idx_reward_overflow_logs_created ON reward_overflow_logs(created_at);

-- 注释
COMMENT ON TABLE reward_overflow_logs IS '奖励池溢出异常日志 - 记录链上比例之和超过100%的异常情况';
COMMENT ON COLUMN reward_overflow_logs.agent_chain IS '代理商链及各级比例配置的JSON记录';
COMMENT ON COLUMN reward_overflow_logs.total_rate IS '链上所有代理商比例之和，超过1.0000为溢出';

-- ============================================================
-- 8. 视图：待处理的阶段奖励
-- ============================================================
CREATE OR REPLACE VIEW v_pending_stage_rewards AS
SELECT
    tsr.id,
    tsr.progress_id,
    tsr.terminal_sn,
    tsr.stage_order,
    tsr.stage_start,
    tsr.stage_end,
    tsr.target_value,
    tsr.actual_value,
    tsr.status,
    trp.template_id,
    trp.bind_agent_id,
    trp.template_snapshot,
    rpt.dimension_type,
    rpt.trans_types,
    rpt.amount_min,
    rpt.amount_max,
    rpt.allow_gap
FROM terminal_stage_rewards tsr
JOIN terminal_reward_progress trp ON tsr.progress_id = trp.id
JOIN reward_policy_templates rpt ON trp.template_id = rpt.id
WHERE tsr.status = 'pending'
  AND tsr.stage_end < NOW()
  AND trp.status = 'active';

COMMENT ON VIEW v_pending_stage_rewards IS '待处理的阶段奖励视图 - 用于定时任务查询';

-- ============================================================
-- 示例数据（注释，实际使用时可取消注释）
-- ============================================================
-- 奖励政策模版示例
-- INSERT INTO reward_policy_templates (name, time_type, dimension_type, trans_types, amount_min, amount_max, allow_gap)
-- VALUES
--     ('标准激活奖励(按天数/金额)', 'days', 'amount', 'scan,debit,credit', 0, NULL, FALSE),
--     ('快速激活奖励(按天数/笔数)', 'days', 'count', 'scan,debit,credit', 20000, NULL, TRUE),
--     ('月度考核奖励(按自然月/金额)', 'months', 'amount', 'credit', 0, 20000000, FALSE);

-- 奖励阶段配置示例
-- INSERT INTO reward_stages (template_id, stage_order, start_value, end_value, target_value, reward_amount)
-- VALUES
--     (1, 1, 1, 10, 1000000, 5000),      -- 第1-10天，交易满1万元，奖励50元
--     (1, 2, 11, 20, 2000000, 10000),    -- 第11-20天，交易满2万元，奖励100元
--     (1, 3, 21, 30, 3000000, 15000);    -- 第21-30天，交易满3万元，奖励150元

-- 代理商奖励比例配置示例
-- INSERT INTO agent_reward_rates (agent_id, reward_rate)
-- VALUES
--     (1, 0.0500),  -- 顶级代理商：5%
--     (2, 0.1000),  -- 中级代理商：10%
--     (3, 0.0000);  -- 终端归属代理商：0%（拿剩余部分）
