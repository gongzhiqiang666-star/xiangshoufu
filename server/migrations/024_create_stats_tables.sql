-- 统计汇总表迁移 - 支持首页与数据分析模块
-- 包含：代理商每日/每月汇总表、物化路径字段

-- ============================================
-- 1. 代理商每日统计汇总表
-- ============================================
CREATE TABLE IF NOT EXISTS agent_daily_stats (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    stat_date       DATE NOT NULL,
    scope           VARCHAR(10) NOT NULL,  -- 'direct' 直营 或 'team' 团队

    -- 交易统计
    trans_amount    BIGINT DEFAULT 0,      -- 交易金额(分)
    trans_count     INT DEFAULT 0,          -- 交易笔数

    -- 分润统计(按类型分)
    profit_trade    BIGINT DEFAULT 0,       -- 交易分润(分)
    profit_deposit  BIGINT DEFAULT 0,       -- 押金返现(分)
    profit_sim      BIGINT DEFAULT 0,       -- 流量返现(分)
    profit_reward   BIGINT DEFAULT 0,       -- 激活奖励(分)
    profit_total    BIGINT DEFAULT 0,       -- 总分润(分)

    -- 商户与终端统计
    merchant_count      INT DEFAULT 0,      -- 商户数量
    merchant_new        INT DEFAULT 0,      -- 新增商户数
    terminal_total      INT DEFAULT 0,      -- 终端总数
    terminal_activated  INT DEFAULT 0,      -- 已激活终端数
    terminal_new_activated INT DEFAULT 0,   -- 当日新激活数

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(agent_id, stat_date, scope)
);

-- 索引优化
CREATE INDEX IF NOT EXISTS idx_agent_daily_stats_date ON agent_daily_stats (stat_date, agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_daily_stats_agent ON agent_daily_stats (agent_id, stat_date);

COMMENT ON TABLE agent_daily_stats IS '代理商每日统计汇总表';
COMMENT ON COLUMN agent_daily_stats.scope IS '统计范围: direct=直营, team=团队';
COMMENT ON COLUMN agent_daily_stats.profit_trade IS '交易分润(分)';
COMMENT ON COLUMN agent_daily_stats.profit_deposit IS '押金返现(分)';
COMMENT ON COLUMN agent_daily_stats.profit_sim IS '流量返现(分)';
COMMENT ON COLUMN agent_daily_stats.profit_reward IS '激活奖励(分)';

-- ============================================
-- 2. 代理商每月统计汇总表
-- ============================================
CREATE TABLE IF NOT EXISTS agent_monthly_stats (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,
    stat_month      VARCHAR(7) NOT NULL,    -- '2026-01' 格式
    scope           VARCHAR(10) NOT NULL,   -- 'direct' 或 'team'

    -- 交易统计
    trans_amount    BIGINT DEFAULT 0,
    trans_count     INT DEFAULT 0,

    -- 分润统计
    profit_trade    BIGINT DEFAULT 0,
    profit_deposit  BIGINT DEFAULT 0,
    profit_sim      BIGINT DEFAULT 0,
    profit_reward   BIGINT DEFAULT 0,
    profit_total    BIGINT DEFAULT 0,

    -- 商户与终端统计
    merchant_count      INT DEFAULT 0,
    merchant_new        INT DEFAULT 0,
    terminal_total      INT DEFAULT 0,
    terminal_activated  INT DEFAULT 0,

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE(agent_id, stat_month, scope)
);

CREATE INDEX IF NOT EXISTS idx_agent_monthly_stats_month ON agent_monthly_stats (stat_month, agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_monthly_stats_agent ON agent_monthly_stats (agent_id, stat_month);

COMMENT ON TABLE agent_monthly_stats IS '代理商每月统计汇总表';

-- ============================================
-- 3. 确保agents表有agent_path字段(物化路径)
-- ============================================
-- 注意: agents表已有path字段，这里添加别名字段以保持向后兼容
DO $$
BEGIN
    -- 检查是否已有agent_path字段
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'agents' AND column_name = 'agent_path'
    ) THEN
        -- 添加agent_path字段
        ALTER TABLE agents ADD COLUMN agent_path VARCHAR(500);

        -- 从现有path字段初始化(如果path有数据)
        UPDATE agents SET agent_path = path WHERE path IS NOT NULL AND path != '';

        -- 创建索引
        CREATE INDEX IF NOT EXISTS idx_agents_agent_path ON agents (agent_path);

        RAISE NOTICE 'agent_path字段已添加';
    END IF;
END $$;

COMMENT ON COLUMN agents.agent_path IS '物化路径，如/1/5/12/表示层级关系';

-- ============================================
-- 4. 确保merchants表有merchant_type字段
-- ============================================
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'merchants' AND column_name = 'merchant_type'
    ) THEN
        ALTER TABLE merchants ADD COLUMN merchant_type VARCHAR(20) DEFAULT 'normal';
        CREATE INDEX IF NOT EXISTS idx_merchants_type ON merchants (merchant_type);
        RAISE NOTICE 'merchant_type字段已添加';
    END IF;

    -- 确保有type_updated_at字段
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_name = 'merchants' AND column_name = 'type_updated_at'
    ) THEN
        ALTER TABLE merchants ADD COLUMN type_updated_at TIMESTAMPTZ;
        RAISE NOTICE 'type_updated_at字段已添加';
    END IF;
END $$;

-- ============================================
-- 5. 通道统计视图(便于查询通道占比)
-- ============================================
CREATE OR REPLACE VIEW v_channel_daily_stats AS
SELECT
    agent_id,
    DATE(trade_time) as stat_date,
    channel_id,
    channel_code,
    SUM(amount) as trans_amount,
    COUNT(*) as trans_count
FROM transactions
WHERE profit_status != 2  -- 排除已退款
GROUP BY agent_id, DATE(trade_time), channel_id, channel_code;

COMMENT ON VIEW v_channel_daily_stats IS '通道每日交易统计视图';

-- ============================================
-- 6. 初始化现有代理商的agent_path
-- ============================================
-- 递归更新所有代理商的agent_path
WITH RECURSIVE agent_tree AS (
    -- 根节点(无父级)
    SELECT id, parent_id, CONCAT('/', id, '/') as calculated_path
    FROM agents
    WHERE parent_id IS NULL

    UNION ALL

    -- 子节点
    SELECT a.id, a.parent_id, CONCAT(t.calculated_path, a.id, '/')
    FROM agents a
    INNER JOIN agent_tree t ON a.parent_id = t.id
)
UPDATE agents
SET agent_path = agent_tree.calculated_path
FROM agent_tree
WHERE agents.id = agent_tree.id
  AND (agents.agent_path IS NULL OR agents.agent_path = '');

-- ============================================
-- 7. 创建更新agent_path的触发器函数
-- ============================================
CREATE OR REPLACE FUNCTION update_agent_path()
RETURNS TRIGGER AS $$
DECLARE
    parent_path VARCHAR(500);
BEGIN
    IF NEW.parent_id IS NULL THEN
        NEW.agent_path := CONCAT('/', NEW.id, '/');
    ELSE
        SELECT agent_path INTO parent_path FROM agents WHERE id = NEW.parent_id;
        IF parent_path IS NULL THEN
            parent_path := CONCAT('/', NEW.parent_id, '/');
        END IF;
        NEW.agent_path := CONCAT(parent_path, NEW.id, '/');
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器(如果不存在)
DROP TRIGGER IF EXISTS trg_update_agent_path ON agents;
CREATE TRIGGER trg_update_agent_path
    BEFORE INSERT OR UPDATE OF parent_id ON agents
    FOR EACH ROW
    EXECUTE FUNCTION update_agent_path();

COMMENT ON FUNCTION update_agent_path() IS '自动维护代理商物化路径';
