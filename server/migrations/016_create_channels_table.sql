-- 通道表
-- 创建时间: 2026-01-22
-- 说明: 存储支付通道基础信息

-- ============================================================
-- 通道表
-- ============================================================
CREATE TABLE IF NOT EXISTS channels (
    id                  BIGSERIAL PRIMARY KEY,
    channel_code        VARCHAR(32) NOT NULL UNIQUE,        -- 通道编码（如：HENGXINTONG）
    channel_name        VARCHAR(64) NOT NULL,               -- 通道名称（如：恒信通）
    description         TEXT,                               -- 通道描述
    status              SMALLINT DEFAULT 1,                 -- 1:启用 0:禁用
    priority            INTEGER DEFAULT 0,                  -- 优先级（越大越优先）
    config              JSONB,                              -- 通道配置（JSON格式）
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_channels_status ON channels(status);

-- 注释
COMMENT ON TABLE channels IS '支付通道表';
COMMENT ON COLUMN channels.channel_code IS '通道编码';
COMMENT ON COLUMN channels.channel_name IS '通道名称';
COMMENT ON COLUMN channels.status IS '状态：1启用 0禁用';

-- 插入默认通道数据
INSERT INTO channels (channel_code, channel_name, description, status, priority) VALUES
    ('HENGXINTONG', '恒信通', '恒信通支付通道', 1, 100),
    ('LAKALA', '拉卡拉', '拉卡拉支付通道', 0, 90),
    ('YEAHKA', '乐刷', '乐刷支付通道', 0, 80),
    ('SUIXINGFU', '随行付', '随行付支付通道', 0, 70),
    ('LIANLIAN', '连连支付', '连连支付通道', 0, 60),
    ('SANDPAY', '杉德支付', '杉德支付通道', 0, 50),
    ('FUIOU', '富友支付', '富友支付通道', 0, 40),
    ('HEEPAY', '汇付天下', '汇付天下支付通道', 0, 30)
ON CONFLICT (channel_code) DO NOTHING;
