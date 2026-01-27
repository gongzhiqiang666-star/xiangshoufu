-- 修改代理商奖励配置表：从比例分配改为差额分配模式
-- 创建时间: 2026-01-26
-- 业务规则变更:
--   原规则：比例分配 - 每级代理商按比例瓜分奖励池
--   新规则：差额分配 - 上级给下级配置固定金额，下级拿配置金额，上级拿差额
--   例如：A给B配置100元，B给C配置30元，C达标时：C得30，B得70（100-30），总金额100

-- ============================================================
-- 1. 修改 agent_reward_rates 表
-- ============================================================

-- 删除旧的约束
ALTER TABLE agent_reward_rates DROP CONSTRAINT IF EXISTS chk_agent_reward_rates_rate;

-- 删除旧的唯一索引
DROP INDEX IF EXISTS idx_agent_reward_rates_agent;

-- 添加 template_id 列（差额分配需要按模版配置）
ALTER TABLE agent_reward_rates ADD COLUMN IF NOT EXISTS template_id BIGINT;

-- 将 reward_rate 列重命名为 reward_amount，并修改类型
-- 注意：如果有数据需要迁移，需要先处理
ALTER TABLE agent_reward_rates DROP COLUMN IF EXISTS reward_rate;
ALTER TABLE agent_reward_rates ADD COLUMN IF NOT EXISTS reward_amount BIGINT NOT NULL DEFAULT 0;

-- 添加新的约束
ALTER TABLE agent_reward_rates ADD CONSTRAINT chk_agent_reward_rates_amount CHECK (reward_amount >= 0);

-- 创建新的唯一索引（agent_id + template_id 唯一）
CREATE UNIQUE INDEX IF NOT EXISTS idx_agent_reward_rates_agent_template ON agent_reward_rates(agent_id, template_id);

-- 创建模版索引
CREATE INDEX IF NOT EXISTS idx_agent_reward_rates_template ON agent_reward_rates(template_id);

-- 更新表注释
COMMENT ON TABLE agent_reward_rates IS '代理商奖励金额配置表 - 差额分配模式：上级给下级配置的奖励金额';
COMMENT ON COLUMN agent_reward_rates.agent_id IS '代理商ID（被配置的下级代理商）';
COMMENT ON COLUMN agent_reward_rates.template_id IS '奖励模版ID（配置针对哪个模版）';
COMMENT ON COLUMN agent_reward_rates.reward_amount IS '奖励金额（分）- 上级给下级配置的金额，下级拿此金额，上级拿差额';

-- ============================================================
-- 2. 更新 reward_distributions 表注释
-- ============================================================
COMMENT ON TABLE reward_distributions IS '奖励发放记录表 - 差额分配模式的级差分配记录';
COMMENT ON COLUMN reward_distributions.reward_rate IS '该代理商获得金额占总奖励的比例（仅用于统计展示）';
COMMENT ON COLUMN reward_distributions.reward_amount IS '实际获得的奖励金额（分）= 上级配置金额 - 给下级配置金额';

-- ============================================================
-- 3. 更新 reward_overflow_logs 表注释
-- ============================================================
COMMENT ON TABLE reward_overflow_logs IS '奖励分配异常日志 - 记录差额分配过程中的异常情况';
COMMENT ON COLUMN reward_overflow_logs.total_rate IS '分配总金额占奖励金额的比例（用于异常诊断）';
COMMENT ON COLUMN reward_overflow_logs.error_message IS '异常信息：如配置金额不足、分配总额超出等';
