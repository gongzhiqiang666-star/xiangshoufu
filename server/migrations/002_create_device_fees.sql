-- 享收付 - 流量费/服务费记录表

CREATE TABLE IF NOT EXISTS device_fees (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,                -- 通道ID
    channel_code    VARCHAR(32) NOT NULL,           -- 通道编码

    -- 关联信息
    terminal_sn     VARCHAR(50) NOT NULL,           -- 机具SN号
    merchant_no     VARCHAR(64),                    -- 商户号
    agent_id        BIGINT,                         -- 直属代理商ID

    -- 扣费信息
    order_no        VARCHAR(64) NOT NULL UNIQUE,    -- 订单号（幂等键）
    fee_type        SMALLINT NOT NULL,              -- 1:服务费 2:流量费/通讯费
    fee_amount      BIGINT NOT NULL,                -- 扣费金额（分）

    -- 返现状态
    cashback_status SMALLINT DEFAULT 0,             -- 0:待计算 1:已返现 2:不返现
    cashback_amount BIGINT DEFAULT 0,               -- 返现金额（分）

    -- 时间
    charging_time   TIMESTAMPTZ NOT NULL,           -- 扣款时间
    received_at     TIMESTAMPTZ DEFAULT NOW(),      -- 接收时间

    -- 扩展字段
    brand_code      VARCHAR(32),                    -- 品牌编号
    ext_data        JSONB,                          -- 通道特有扩展字段

    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_device_fees_terminal ON device_fees(terminal_sn);
CREATE INDEX IF NOT EXISTS idx_device_fees_merchant ON device_fees(merchant_no);
CREATE INDEX IF NOT EXISTS idx_device_fees_agent ON device_fees(agent_id);
CREATE INDEX IF NOT EXISTS idx_device_fees_status ON device_fees(cashback_status);
CREATE INDEX IF NOT EXISTS idx_device_fees_time ON device_fees(charging_time);

-- 添加注释
COMMENT ON TABLE device_fees IS '流量费/服务费记录表';
COMMENT ON COLUMN device_fees.fee_type IS '费用类型: 1-服务费, 2-流量费/通讯费';
COMMENT ON COLUMN device_fees.cashback_status IS '返现状态: 0-待计算, 1-已返现, 2-不返现';
