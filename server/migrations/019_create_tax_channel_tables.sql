-- 税筹通道表
-- 描述: 管理税筹通道配置，用于计算提现时的税费扣除

CREATE TABLE IF NOT EXISTS tax_channels (
    id              BIGSERIAL PRIMARY KEY,
    channel_code    VARCHAR(32) NOT NULL UNIQUE,
    channel_name    VARCHAR(100) NOT NULL,

    -- 扣费规则
    fee_type        SMALLINT NOT NULL DEFAULT 2,       -- 1=付款扣 2=出款扣
    tax_rate        DECIMAL(5,4) NOT NULL DEFAULT 0,   -- 税率 如0.09表示9%
    fixed_fee       BIGINT DEFAULT 0,                  -- 固定费用(分) 如300表示3元/笔

    -- 接口配置
    api_url         VARCHAR(255),
    api_key         VARCHAR(255),
    api_secret      VARCHAR(255),

    -- 状态
    status          SMALLINT DEFAULT 1,                -- 1=启用 0=禁用
    remark          VARCHAR(500),

    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_tax_channels_status ON tax_channels(status);
CREATE INDEX IF NOT EXISTS idx_tax_channels_code ON tax_channels(channel_code);

-- 通道-税筹通道关联表
-- 描述: 不同支付通道的不同钱包类型可以走不同的税筹通道
CREATE TABLE IF NOT EXISTS channel_tax_mappings (
    id              BIGSERIAL PRIMARY KEY,
    channel_id      BIGINT NOT NULL,                   -- 支付通道ID
    wallet_type     SMALLINT NOT NULL,                 -- 钱包类型: 1分润 2服务费 3奖励 4充值 5沉淀
    tax_channel_id  BIGINT NOT NULL,                   -- 税筹通道ID

    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW(),

    CONSTRAINT uk_channel_wallet UNIQUE(channel_id, wallet_type),
    CONSTRAINT fk_tax_channel FOREIGN KEY (tax_channel_id) REFERENCES tax_channels(id)
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_channel_tax_mappings_channel ON channel_tax_mappings(channel_id);
CREATE INDEX IF NOT EXISTS idx_channel_tax_mappings_tax_channel ON channel_tax_mappings(tax_channel_id);

-- 添加注释
COMMENT ON TABLE tax_channels IS '税筹通道配置表';
COMMENT ON COLUMN tax_channels.fee_type IS '扣费类型: 1=付款扣(充值时扣) 2=出款扣(提现时扣)';
COMMENT ON COLUMN tax_channels.tax_rate IS '税率，如0.0900表示9%';
COMMENT ON COLUMN tax_channels.fixed_fee IS '固定费用(分)，如300表示3元/笔';

COMMENT ON TABLE channel_tax_mappings IS '支付通道与税筹通道关联表';
COMMENT ON COLUMN channel_tax_mappings.wallet_type IS '钱包类型: 1=分润 2=服务费 3=奖励 4=充值 5=沉淀';

-- 插入默认税筹通道
INSERT INTO tax_channels (channel_code, channel_name, fee_type, tax_rate, fixed_fee, status, remark)
VALUES
    ('DEFAULT', '默认税筹通道', 2, 0.0900, 0, 1, '默认9%税率，无固定费用'),
    ('LOWRATE', '低税率通道', 2, 0.0600, 0, 1, '6%税率，适用于小额提现'),
    ('HIGHRATE', '高税率通道', 2, 0.1200, 300, 1, '12%税率+3元固定费用')
ON CONFLICT (channel_code) DO NOTHING;
