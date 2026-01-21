-- 货款代扣相关表
-- 创建时间: 2026-01-22
-- 业务规则:
--   1. 扣款规则: 分润优先 - 先扣分润钱包，扣完再扣服务费钱包
--   2. 扣款时机: 实时扣款 - 钱包入账时立即触发扣款
--   3. 部分扣款: 有多少扣多少 - 余额不足时部分扣除，剩余下次继续扣
--   4. 扣款上限: 无上限 - 每次入账时全部扣除，直到扣完为止
--   5. 扣款优先级: 货款代扣 > 上级代扣 > 伙伴代扣
--   6. 与提现关系: 待扣金额占用钱包余额，影响可提现金额
--   7. 接收确认: 下级拒绝接收则终端划拨失败
--   8. 修改/取消: 设置后不可修改，只能扣完或线下协商

-- 货款代扣表
CREATE TABLE IF NOT EXISTS goods_deductions (
    id BIGSERIAL PRIMARY KEY,
    deduction_no VARCHAR(64) NOT NULL UNIQUE,             -- 代扣编号
    from_agent_id BIGINT NOT NULL,                        -- 上级代理商ID（扣款方/发起方）
    to_agent_id BIGINT NOT NULL,                          -- 下级代理商ID（被扣款方/接收方）
    total_amount BIGINT NOT NULL,                         -- 代扣总金额（分）
    deducted_amount BIGINT DEFAULT 0,                     -- 已扣金额（分）
    remaining_amount BIGINT NOT NULL,                     -- 剩余金额（分）
    deduction_source SMALLINT NOT NULL DEFAULT 3,         -- 扣款来源: 1=分润钱包 2=服务费钱包 3=两者都扣
    terminal_count INT NOT NULL DEFAULT 0,                -- 终端数量
    unit_price BIGINT NOT NULL DEFAULT 0,                 -- 单价（分）
    status SMALLINT NOT NULL DEFAULT 1,                   -- 状态: 1=待接收 2=进行中 3=已完成 4=已拒绝
    agreement_signed BOOLEAN DEFAULT FALSE,               -- 是否签署协议
    agreement_url VARCHAR(500),                           -- 协议文件URL
    distribute_id BIGINT,                                 -- 关联的终端划拨ID
    remark VARCHAR(500),                                  -- 备注
    created_by BIGINT,                                    -- 创建人
    accepted_at TIMESTAMPTZ,                              -- 接收时间
    completed_at TIMESTAMPTZ,                             -- 完成时间
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_goods_deductions_from_agent ON goods_deductions(from_agent_id);
CREATE INDEX idx_goods_deductions_to_agent ON goods_deductions(to_agent_id);
CREATE INDEX idx_goods_deductions_status ON goods_deductions(status);
CREATE INDEX idx_goods_deductions_distribute ON goods_deductions(distribute_id);

-- 货款代扣明细表（每次扣款记录）
CREATE TABLE IF NOT EXISTS goods_deduction_details (
    id BIGSERIAL PRIMARY KEY,
    deduction_id BIGINT NOT NULL REFERENCES goods_deductions(id),  -- 关联货款代扣ID
    deduction_no VARCHAR(64),                             -- 代扣编号（冗余）
    amount BIGINT NOT NULL,                               -- 本次扣款金额（分）
    wallet_type SMALLINT NOT NULL,                        -- 扣款钱包类型: 1=分润 2=服务费
    channel_id BIGINT,                                    -- 通道ID（钱包对应的通道）
    wallet_balance_before BIGINT NOT NULL,                -- 扣款前余额（分）
    wallet_balance_after BIGINT NOT NULL,                 -- 扣款后余额（分）
    cumulative_deducted BIGINT NOT NULL,                  -- 累计已扣金额（分）
    remaining_after BIGINT NOT NULL,                      -- 扣款后剩余待扣（分）
    trigger_type VARCHAR(32),                             -- 触发类型: profit_income, service_fee_income
    trigger_transaction_id BIGINT,                        -- 触发扣款的交易ID
    trigger_profit_id BIGINT,                             -- 触发扣款的分润记录ID
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_goods_deduction_details_deduction ON goods_deduction_details(deduction_id);
CREATE INDEX idx_goods_deduction_details_trigger_tx ON goods_deduction_details(trigger_transaction_id);
CREATE INDEX idx_goods_deduction_details_created ON goods_deduction_details(created_at);

-- 货款代扣终端关联表
CREATE TABLE IF NOT EXISTS goods_deduction_terminals (
    id BIGSERIAL PRIMARY KEY,
    deduction_id BIGINT NOT NULL REFERENCES goods_deductions(id),  -- 关联货款代扣ID
    terminal_id BIGINT NOT NULL,                          -- 终端ID
    terminal_sn VARCHAR(50),                              -- 终端SN（冗余）
    unit_price BIGINT NOT NULL,                           -- 单价（分）
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_goods_deduction_terminals_deduction ON goods_deduction_terminals(deduction_id);
CREATE INDEX idx_goods_deduction_terminals_terminal ON goods_deduction_terminals(terminal_id);

-- 货款代扣通知表（用于APP推送记录）
CREATE TABLE IF NOT EXISTS goods_deduction_notifications (
    id BIGSERIAL PRIMARY KEY,
    deduction_id BIGINT NOT NULL REFERENCES goods_deductions(id),
    detail_id BIGINT REFERENCES goods_deduction_details(id),  -- 关联扣款明细
    agent_id BIGINT NOT NULL,                             -- 接收通知的代理商
    notify_type SMALLINT NOT NULL,                        -- 通知类型: 1=待接收 2=扣款通知 3=完成通知
    title VARCHAR(100),                                   -- 通知标题
    content VARCHAR(500),                                 -- 通知内容
    is_read BOOLEAN DEFAULT FALSE,                        -- 是否已读
    read_at TIMESTAMPTZ,                                  -- 阅读时间
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_goods_deduction_notifications_agent ON goods_deduction_notifications(agent_id);
CREATE INDEX idx_goods_deduction_notifications_unread ON goods_deduction_notifications(agent_id, is_read) WHERE is_read = FALSE;

-- 添加注释
COMMENT ON TABLE goods_deductions IS '货款代扣表 - 终端划拨时设置的货款代扣，与代扣管理模块独立';
COMMENT ON COLUMN goods_deductions.deduction_source IS '扣款来源：1分润钱包 2服务费钱包 3两者都扣（优先分润）';
COMMENT ON COLUMN goods_deductions.status IS '状态：1待接收 2进行中 3已完成 4已拒绝';

COMMENT ON TABLE goods_deduction_details IS '货款代扣明细表 - 每次钱包入账触发的实时扣款记录';
COMMENT ON COLUMN goods_deduction_details.wallet_type IS '钱包类型：1分润 2服务费';
COMMENT ON COLUMN goods_deduction_details.trigger_type IS '触发类型：profit_income分润入账 service_fee_income服务费入账';

COMMENT ON TABLE goods_deduction_terminals IS '货款代扣终端关联表 - 记录代扣关联的终端及单价';

COMMENT ON TABLE goods_deduction_notifications IS '货款代扣通知表 - 记录待接收、扣款、完成等通知';
COMMENT ON COLUMN goods_deduction_notifications.notify_type IS '通知类型：1待接收 2扣款通知 3完成通知';
