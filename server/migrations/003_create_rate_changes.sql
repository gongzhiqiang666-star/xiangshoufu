-- 享收付 - 费率变更记录表

CREATE TABLE IF NOT EXISTS rate_changes (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,                -- 通道ID
    channel_code    VARCHAR(32) NOT NULL,           -- 通道编码

    -- 关联信息
    terminal_sn     VARCHAR(50) NOT NULL,           -- 机具SN号
    merchant_no     VARCHAR(64) NOT NULL,           -- 商户号
    agent_id        BIGINT,                         -- 代理商ID

    -- 费率信息（百分比，如0.60表示0.60%）
    credit_rate         DECIMAL(10,4),              -- 贷记卡费率
    credit_extra_rate   DECIMAL(10,4),              -- 贷记卡额外费率
    debit_rate          DECIMAL(10,4),              -- 借记卡费率
    alipay_rate         DECIMAL(10,4),              -- 支付宝费率
    wechat_rate         DECIMAL(10,4),              -- 微信费率
    unionpay_rate       DECIMAL(10,4),              -- 银联云闪付费率

    -- 调价费率
    credit_addition_rate    DECIMAL(10,4),          -- 贷记卡调价费率
    unionpay_addition_rate  DECIMAL(10,4),          -- 银联调价费率
    alipay_addition_rate    DECIMAL(10,4),          -- 支付宝调价费率
    wechat_addition_rate    DECIMAL(10,4),          -- 微信调价费率

    -- 同步状态
    sync_status     SMALLINT DEFAULT 0,             -- 0:待同步 1:已同步到终端表

    -- 时间
    received_at     TIMESTAMPTZ DEFAULT NOW(),      -- 接收时间

    -- 扩展字段
    brand_code      VARCHAR(32),                    -- 品牌编号
    ext_data        JSONB,                          -- 通道特有扩展字段

    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_rate_changes_terminal ON rate_changes(terminal_sn);
CREATE INDEX IF NOT EXISTS idx_rate_changes_merchant ON rate_changes(merchant_no);
CREATE INDEX IF NOT EXISTS idx_rate_changes_time ON rate_changes(received_at);

-- 添加注释
COMMENT ON TABLE rate_changes IS '费率变更记录表 - 记录通道推送的费率变更信息';
COMMENT ON COLUMN rate_changes.sync_status IS '同步状态: 0-待同步, 1-已同步到终端表';
