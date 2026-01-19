-- 收享付 - 原始回调数据表
-- 用于存储所有通道的原始回调数据，支持按日期分区

-- 创建原始回调日志表（分区表）
CREATE TABLE IF NOT EXISTS raw_callback_logs (
    id              BIGSERIAL,
    channel_code    VARCHAR(32) NOT NULL,           -- 通道编码: HENGXINTONG, LAKALA等
    action_type     VARCHAR(64) NOT NULL,           -- 回调类型: merc_income, sn_bind, pos_order等
    raw_request     JSONB NOT NULL,                 -- 原始JSON请求体
    sign_verified   BOOLEAN DEFAULT FALSE,          -- 签名验证结果
    process_status  SMALLINT DEFAULT 0,             -- 0:待处理 1:处理成功 2:处理失败
    error_message   TEXT,                           -- 处理失败时的错误信息
    retry_count     SMALLINT DEFAULT 0,             -- 重试次数
    idempotent_key  VARCHAR(128) NOT NULL,          -- 幂等键 (channel_code + action_type + 业务唯一键)
    client_ip       VARCHAR(45),                    -- 请求来源IP
    received_at     TIMESTAMPTZ DEFAULT NOW(),      -- 接收时间
    processed_at    TIMESTAMPTZ,                    -- 处理完成时间
    created_date    DATE DEFAULT CURRENT_DATE,      -- 分区键
    PRIMARY KEY (id, created_date)
) PARTITION BY RANGE (created_date);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_raw_callback_channel_action ON raw_callback_logs(channel_code, action_type);
CREATE INDEX IF NOT EXISTS idx_raw_callback_status ON raw_callback_logs(process_status, created_date);
CREATE INDEX IF NOT EXISTS idx_raw_callback_idempotent ON raw_callback_logs(idempotent_key);
CREATE INDEX IF NOT EXISTS idx_raw_callback_received ON raw_callback_logs(received_at);

-- 创建初始分区（当前月份及未来3个月）
DO $$
DECLARE
    start_date DATE := DATE_TRUNC('month', CURRENT_DATE);
    end_date DATE;
    partition_name TEXT;
BEGIN
    FOR i IN 0..3 LOOP
        end_date := start_date + INTERVAL '1 month';
        partition_name := 'raw_callback_logs_' || TO_CHAR(start_date, 'YYYY_MM');

        EXECUTE format(
            'CREATE TABLE IF NOT EXISTS %I PARTITION OF raw_callback_logs
             FOR VALUES FROM (%L) TO (%L)',
            partition_name, start_date, end_date
        );

        start_date := end_date;
    END LOOP;
END $$;

-- 添加注释
COMMENT ON TABLE raw_callback_logs IS '原始回调日志表 - 存储所有通道的原始回调数据';
COMMENT ON COLUMN raw_callback_logs.channel_code IS '通道编码: HENGXINTONG(恒信通), LAKALA(拉卡拉)等';
COMMENT ON COLUMN raw_callback_logs.action_type IS '回调类型: merc_income(商户入网), sn_bind(终端绑定), pos_order(交易)等';
COMMENT ON COLUMN raw_callback_logs.process_status IS '处理状态: 0-待处理, 1-处理成功, 2-处理失败';
COMMENT ON COLUMN raw_callback_logs.idempotent_key IS '幂等键，防止重复处理';
