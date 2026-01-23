-- 费率同步日志表
-- 记录费率修改同步到支付通道的操作日志

CREATE TABLE IF NOT EXISTS rate_sync_logs (
    id BIGSERIAL PRIMARY KEY,

    -- 关联信息
    merchant_id BIGINT NOT NULL,                    -- 商户ID
    merchant_no VARCHAR(64) NOT NULL,               -- 商户号
    terminal_sn VARCHAR(64),                        -- 终端SN
    channel_code VARCHAR(32) NOT NULL,              -- 通道编码
    agent_id BIGINT NOT NULL,                       -- 操作代理商ID

    -- 费率信息（修改前）
    old_credit_rate DECIMAL(10,4),                  -- 原贷记卡费率
    old_debit_rate DECIMAL(10,4),                   -- 原借记卡费率
    old_debit_cap BIGINT,                           -- 原借记卡封顶（分）
    old_wechat_rate DECIMAL(10,4),                  -- 原微信费率
    old_alipay_rate DECIMAL(10,4),                  -- 原支付宝费率
    old_unionpay_rate DECIMAL(10,4),                -- 原云闪付费率

    -- 费率信息（修改后）
    new_credit_rate DECIMAL(10,4),                  -- 新贷记卡费率
    new_debit_rate DECIMAL(10,4),                   -- 新借记卡费率
    new_debit_cap BIGINT,                           -- 新借记卡封顶（分）
    new_wechat_rate DECIMAL(10,4),                  -- 新微信费率
    new_alipay_rate DECIMAL(10,4),                  -- 新支付宝费率
    new_unionpay_rate DECIMAL(10,4),                -- 新云闪付费率

    -- 同步状态
    sync_status SMALLINT NOT NULL DEFAULT 0,        -- 0-待同步 1-同步中 2-同步成功 3-同步失败
    channel_trade_no VARCHAR(128),                  -- 通道返回的流水号
    error_message TEXT,                             -- 错误信息
    retry_count INT NOT NULL DEFAULT 0,             -- 重试次数
    max_retries INT NOT NULL DEFAULT 3,             -- 最大重试次数
    next_retry_at TIMESTAMP WITH TIME ZONE,         -- 下次重试时间

    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    synced_at TIMESTAMP WITH TIME ZONE              -- 同步成功时间
);

-- 创建索引
CREATE INDEX idx_rate_sync_logs_merchant_id ON rate_sync_logs(merchant_id);
CREATE INDEX idx_rate_sync_logs_merchant_no ON rate_sync_logs(merchant_no);
CREATE INDEX idx_rate_sync_logs_channel_code ON rate_sync_logs(channel_code);
CREATE INDEX idx_rate_sync_logs_agent_id ON rate_sync_logs(agent_id);
CREATE INDEX idx_rate_sync_logs_sync_status ON rate_sync_logs(sync_status);
CREATE INDEX idx_rate_sync_logs_created_at ON rate_sync_logs(created_at);
CREATE INDEX idx_rate_sync_logs_next_retry ON rate_sync_logs(next_retry_at) WHERE sync_status = 3 AND retry_count < max_retries;

-- 添加注释
COMMENT ON TABLE rate_sync_logs IS '费率同步日志表';
COMMENT ON COLUMN rate_sync_logs.sync_status IS '同步状态: 0-待同步 1-同步中 2-同步成功 3-同步失败';
