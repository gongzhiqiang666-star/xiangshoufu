-- 定时任务管理表
-- 创建时间: 2025年
-- 功能: 任务配置、执行日志、告警配置、告警记录

-- ============================================
-- 1. 定时任务配置表
-- ============================================
CREATE TABLE IF NOT EXISTS job_configs (
    id              BIGSERIAL PRIMARY KEY,
    job_name        VARCHAR(100) NOT NULL UNIQUE,   -- 任务名称（唯一标识）
    job_desc        VARCHAR(255),                    -- 任务描述
    cron_expr       VARCHAR(50),                     -- Cron表达式（预留）
    interval_seconds INT DEFAULT 300,                -- 执行间隔(秒)，默认5分钟
    is_enabled      BOOLEAN DEFAULT true,            -- 是否启用
    max_retries     INT DEFAULT 3,                   -- 最大重试次数
    retry_interval  INT DEFAULT 60,                  -- 初始重试间隔(秒)
    alert_threshold INT DEFAULT 3,                   -- 连续失败N次告警
    timeout_seconds INT DEFAULT 3600,                -- 任务超时时间(秒)，默认1小时
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE job_configs IS '定时任务配置表';
COMMENT ON COLUMN job_configs.job_name IS '任务名称，唯一标识';
COMMENT ON COLUMN job_configs.job_desc IS '任务描述';
COMMENT ON COLUMN job_configs.cron_expr IS 'Cron表达式（预留）';
COMMENT ON COLUMN job_configs.interval_seconds IS '执行间隔(秒)';
COMMENT ON COLUMN job_configs.is_enabled IS '是否启用';
COMMENT ON COLUMN job_configs.max_retries IS '最大重试次数';
COMMENT ON COLUMN job_configs.retry_interval IS '初始重试间隔(秒)';
COMMENT ON COLUMN job_configs.alert_threshold IS '连续失败N次触发告警';
COMMENT ON COLUMN job_configs.timeout_seconds IS '任务超时时间(秒)';

-- ============================================
-- 2. 任务执行日志表
-- ============================================
CREATE TABLE IF NOT EXISTS job_execution_logs (
    id              BIGSERIAL PRIMARY KEY,
    job_name        VARCHAR(100) NOT NULL,           -- 任务名称
    started_at      TIMESTAMPTZ NOT NULL,            -- 开始时间
    ended_at        TIMESTAMPTZ,                     -- 结束时间
    duration_ms     BIGINT,                          -- 执行耗时(毫秒)
    status          SMALLINT NOT NULL DEFAULT 3,     -- 1成功 2失败 3运行中
    processed_count INT DEFAULT 0,                   -- 处理条数
    success_count   INT DEFAULT 0,                   -- 成功条数
    fail_count      INT DEFAULT 0,                   -- 失败条数
    error_message   TEXT,                            -- 错误信息
    error_stack     TEXT,                            -- 错误堆栈
    retry_count     INT DEFAULT 0,                   -- 当前重试次数
    trigger_type    SMALLINT DEFAULT 1,              -- 1自动触发 2手动触发
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE job_execution_logs IS '任务执行日志表';
COMMENT ON COLUMN job_execution_logs.job_name IS '任务名称';
COMMENT ON COLUMN job_execution_logs.started_at IS '开始时间';
COMMENT ON COLUMN job_execution_logs.ended_at IS '结束时间';
COMMENT ON COLUMN job_execution_logs.duration_ms IS '执行耗时(毫秒)';
COMMENT ON COLUMN job_execution_logs.status IS '状态：1成功 2失败 3运行中';
COMMENT ON COLUMN job_execution_logs.processed_count IS '处理条数';
COMMENT ON COLUMN job_execution_logs.success_count IS '成功条数';
COMMENT ON COLUMN job_execution_logs.fail_count IS '失败条数';
COMMENT ON COLUMN job_execution_logs.error_message IS '错误信息';
COMMENT ON COLUMN job_execution_logs.error_stack IS '错误堆栈';
COMMENT ON COLUMN job_execution_logs.retry_count IS '当前重试次数';
COMMENT ON COLUMN job_execution_logs.trigger_type IS '触发类型：1自动触发 2手动触发';

-- 索引
CREATE INDEX IF NOT EXISTS idx_job_logs_job_name ON job_execution_logs(job_name);
CREATE INDEX IF NOT EXISTS idx_job_logs_started_at ON job_execution_logs(started_at);
CREATE INDEX IF NOT EXISTS idx_job_logs_status ON job_execution_logs(status);
CREATE INDEX IF NOT EXISTS idx_job_logs_created_at ON job_execution_logs(created_at);

-- ============================================
-- 3. 告警配置表
-- ============================================
CREATE TABLE IF NOT EXISTS alert_configs (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL,           -- 配置名称
    channel_type    SMALLINT NOT NULL,               -- 1钉钉 2企微 3邮件
    webhook_url     VARCHAR(500),                    -- Webhook地址（钉钉/企微）
    webhook_secret  VARCHAR(200),                    -- Webhook密钥（钉钉签名）
    email_addresses TEXT,                            -- 邮箱地址(逗号分隔)
    email_smtp_host VARCHAR(100),                    -- SMTP服务器
    email_smtp_port INT DEFAULT 465,                 -- SMTP端口
    email_username  VARCHAR(100),                    -- SMTP用户名
    email_password  VARCHAR(200),                    -- SMTP密码（加密存储）
    is_enabled      BOOLEAN DEFAULT true,
    created_by      BIGINT,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE alert_configs IS '告警配置表';
COMMENT ON COLUMN alert_configs.name IS '配置名称';
COMMENT ON COLUMN alert_configs.channel_type IS '通道类型：1钉钉 2企微 3邮件';
COMMENT ON COLUMN alert_configs.webhook_url IS 'Webhook地址';
COMMENT ON COLUMN alert_configs.webhook_secret IS 'Webhook密钥（钉钉签名用）';
COMMENT ON COLUMN alert_configs.email_addresses IS '邮箱地址，多个用逗号分隔';
COMMENT ON COLUMN alert_configs.email_smtp_host IS 'SMTP服务器地址';
COMMENT ON COLUMN alert_configs.email_smtp_port IS 'SMTP端口';
COMMENT ON COLUMN alert_configs.email_username IS 'SMTP用户名';
COMMENT ON COLUMN alert_configs.email_password IS 'SMTP密码（加密存储）';
COMMENT ON COLUMN alert_configs.is_enabled IS '是否启用';
COMMENT ON COLUMN alert_configs.created_by IS '创建人ID';

-- ============================================
-- 4. 告警记录表
-- ============================================
CREATE TABLE IF NOT EXISTS alert_logs (
    id              BIGSERIAL PRIMARY KEY,
    job_name        VARCHAR(100) NOT NULL,           -- 任务名称
    alert_type      SMALLINT NOT NULL,               -- 1任务失败 2连续失败 3任务超时
    channel_type    SMALLINT NOT NULL,               -- 1钉钉 2企微 3邮件
    config_id       BIGINT,                          -- 关联的告警配置ID
    title           VARCHAR(200),                    -- 告警标题
    message         TEXT NOT NULL,                   -- 告警内容
    send_status     SMALLINT DEFAULT 0,              -- 0待发送 1已发送 2发送失败
    send_at         TIMESTAMPTZ,                     -- 发送时间
    error_message   TEXT,                            -- 发送失败原因
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE alert_logs IS '告警记录表';
COMMENT ON COLUMN alert_logs.job_name IS '任务名称';
COMMENT ON COLUMN alert_logs.alert_type IS '告警类型：1任务失败 2连续失败 3任务超时';
COMMENT ON COLUMN alert_logs.channel_type IS '通道类型：1钉钉 2企微 3邮件';
COMMENT ON COLUMN alert_logs.config_id IS '关联的告警配置ID';
COMMENT ON COLUMN alert_logs.title IS '告警标题';
COMMENT ON COLUMN alert_logs.message IS '告警内容';
COMMENT ON COLUMN alert_logs.send_status IS '发送状态：0待发送 1已发送 2发送失败';
COMMENT ON COLUMN alert_logs.send_at IS '发送时间';
COMMENT ON COLUMN alert_logs.error_message IS '发送失败原因';

-- 索引
CREATE INDEX IF NOT EXISTS idx_alert_logs_job_name ON alert_logs(job_name);
CREATE INDEX IF NOT EXISTS idx_alert_logs_created_at ON alert_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_alert_logs_send_status ON alert_logs(send_status);

-- ============================================
-- 5. 任务失败计数表（用于连续失败告警）
-- ============================================
CREATE TABLE IF NOT EXISTS job_fail_counters (
    id                  BIGSERIAL PRIMARY KEY,
    job_name            VARCHAR(100) NOT NULL UNIQUE,  -- 任务名称
    consecutive_fails   INT DEFAULT 0,                  -- 连续失败次数
    last_fail_at        TIMESTAMPTZ,                    -- 最后失败时间
    last_success_at     TIMESTAMPTZ,                    -- 最后成功时间
    last_alert_at       TIMESTAMPTZ,                    -- 最后告警时间
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE job_fail_counters IS '任务失败计数表';
COMMENT ON COLUMN job_fail_counters.job_name IS '任务名称';
COMMENT ON COLUMN job_fail_counters.consecutive_fails IS '连续失败次数';
COMMENT ON COLUMN job_fail_counters.last_fail_at IS '最后失败时间';
COMMENT ON COLUMN job_fail_counters.last_success_at IS '最后成功时间';
COMMENT ON COLUMN job_fail_counters.last_alert_at IS '最后告警时间';

-- ============================================
-- 6. 初始化默认任务配置
-- ============================================
INSERT INTO job_configs (job_name, job_desc, interval_seconds, max_retries, retry_interval, alert_threshold) VALUES
    ('ProfitCalculatorJob', '分润计算定时任务（兜底重试）', 300, 3, 60, 3),
    ('CallbackRetryJob', '回调重试定时任务', 300, 3, 60, 3),
    ('MessageCleanupJob', '消息清理定时任务', 86400, 3, 300, 3),
    ('DataArchiverJob', '数据归档定时任务', 86400, 3, 300, 3),
    ('PartitionManagerJob', '分区管理定时任务', 2592000, 3, 300, 3),
    ('MerchantTypeCalculatorJob', '商户类型计算定时任务', 86400, 3, 300, 3),
    ('DeductionJob', '每日代扣定时任务', 86400, 3, 300, 3),
    ('SimCashbackJob', '流量费返现定时任务（兜底）', 600, 3, 60, 3),
    ('RewardCheckJob', '激活奖励检查定时任务', 86400, 3, 300, 3),
    ('DepositCashbackJob', '押金返现处理定时任务', 600, 3, 60, 3),
    ('ActivationRewardSettleJob', '激活奖励入账定时任务', 600, 3, 60, 3),
    ('SimCashbackSettleJob', '流量费返现入账定时任务', 600, 3, 60, 3),
    ('WalletBalanceCheckJob', '钱包余额一致性检查任务', 86400, 3, 300, 3)
ON CONFLICT (job_name) DO NOTHING;

-- ============================================
-- 7. 90天日志自动清理函数
-- ============================================
CREATE OR REPLACE FUNCTION cleanup_old_job_logs()
RETURNS INTEGER AS $$
DECLARE
    deleted_count INTEGER;
BEGIN
    DELETE FROM job_execution_logs
    WHERE created_at < NOW() - INTERVAL '90 days';
    GET DIAGNOSTICS deleted_count = ROW_COUNT;

    DELETE FROM alert_logs
    WHERE created_at < NOW() - INTERVAL '90 days';

    RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION cleanup_old_job_logs() IS '清理90天前的任务执行日志和告警记录';
