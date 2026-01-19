-- 终端管理相关表
-- 创建时间: 2026-01-18
-- 业务规则:
--   Q16: 跨级下发时系统自动按层级生成A→B→C的货款代扣链
--   Q29: APP不能跨级，PC可以跨级（保留整个层级关系）

-- 终端表
CREATE TABLE IF NOT EXISTS terminals (
    id BIGSERIAL PRIMARY KEY,
    terminal_sn VARCHAR(50) NOT NULL UNIQUE,       -- 终端序列号
    channel_id BIGINT NOT NULL,                    -- 所属通道
    channel_code VARCHAR(32),                      -- 通道编码
    brand_code VARCHAR(32),                        -- 品牌编码
    model_code VARCHAR(32),                        -- 型号编码
    owner_agent_id BIGINT,                         -- 当前所属代理商
    merchant_id BIGINT,                            -- 绑定的商户
    merchant_no VARCHAR(64),                       -- 商户号
    status SMALLINT DEFAULT 1,                     -- 1:待分配 2:已分配 3:已绑定 4:已激活 5:已解绑 6:已回收
    activated_at TIMESTAMP,                        -- 激活时间
    bound_at TIMESTAMP,                            -- 绑定时间
    sim_fee_count INT DEFAULT 0,                   -- 流量费缴费次数
    last_sim_fee_at TIMESTAMP,                     -- 最后缴费时间
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_terminals_channel ON terminals(channel_id);
CREATE INDEX idx_terminals_owner ON terminals(owner_agent_id);
CREATE INDEX idx_terminals_merchant ON terminals(merchant_id);
CREATE INDEX idx_terminals_merchant_no ON terminals(merchant_no);
CREATE INDEX idx_terminals_status ON terminals(status);

-- 终端下发记录表
CREATE TABLE IF NOT EXISTS terminal_distributes (
    id BIGSERIAL PRIMARY KEY,
    distribute_no VARCHAR(64) NOT NULL UNIQUE,     -- 下发单号
    from_agent_id BIGINT NOT NULL,                 -- 下发方代理商
    to_agent_id BIGINT NOT NULL,                   -- 接收方代理商
    terminal_sn VARCHAR(50) NOT NULL,              -- 终端SN
    channel_id BIGINT NOT NULL,                    -- 通道ID
    is_cross_level BOOLEAN DEFAULT FALSE,          -- 是否跨级下发
    cross_level_path VARCHAR(500),                 -- 跨级路径 /A/B/C/
    goods_price BIGINT NOT NULL,                   -- 货款金额（分）
    deduction_type SMALLINT NOT NULL,              -- 1:一次性付款 2:分期代扣
    deduction_plan_id BIGINT,                      -- 关联代扣计划ID
    chain_id BIGINT,                               -- 关联代扣链ID（跨级时）
    status SMALLINT DEFAULT 1,                     -- 1:待确认 2:已确认 3:已拒绝 4:已取消
    source SMALLINT NOT NULL,                      -- 1:APP 2:PC
    remark VARCHAR(255),
    created_by BIGINT,
    confirmed_by BIGINT,
    created_at TIMESTAMP DEFAULT NOW(),
    confirmed_at TIMESTAMP
);

-- 索引
CREATE INDEX idx_terminal_distributes_from ON terminal_distributes(from_agent_id);
CREATE INDEX idx_terminal_distributes_to ON terminal_distributes(to_agent_id);
CREATE INDEX idx_terminal_distributes_terminal ON terminal_distributes(terminal_sn);
CREATE INDEX idx_terminal_distributes_status ON terminal_distributes(status);

-- 添加注释
COMMENT ON TABLE terminals IS '终端/机具表';
COMMENT ON COLUMN terminals.status IS '状态：1待分配 2已分配 3已绑定 4已激活 5已解绑 6已回收';
COMMENT ON COLUMN terminals.sim_fee_count IS '流量费缴费次数，用于计算返现档次';

COMMENT ON TABLE terminal_distributes IS '终端下发记录表';
COMMENT ON COLUMN terminal_distributes.is_cross_level IS '是否跨级下发，APP不能跨级，PC可以跨级';
COMMENT ON COLUMN terminal_distributes.source IS '来源：1APP（不能跨级） 2PC（可以跨级）';
COMMENT ON COLUMN terminal_distributes.deduction_type IS '货款扣款方式：1一次性付款 2分期代扣';
