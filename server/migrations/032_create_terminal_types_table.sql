-- 终端类型表
-- 用于管理各通道下的终端品牌和型号

CREATE TABLE IF NOT EXISTS terminal_types (
    id BIGSERIAL PRIMARY KEY,
    channel_id BIGINT NOT NULL REFERENCES channels(id) ON DELETE RESTRICT,
    channel_code VARCHAR(50) NOT NULL,
    brand_code VARCHAR(50) NOT NULL,        -- 品牌编码（如：NEWLAND）
    brand_name VARCHAR(100) NOT NULL,       -- 品牌名称（如：新大陆）
    model_code VARCHAR(50) NOT NULL,        -- 型号编码（如：ME31）
    model_name VARCHAR(100),                -- 型号名称（可选）
    description TEXT,                       -- 描述
    status SMALLINT DEFAULT 1,              -- 状态：1启用 0禁用
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(channel_id, brand_code, model_code)
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_terminal_types_channel ON terminal_types(channel_id);
CREATE INDEX IF NOT EXISTS idx_terminal_types_channel_code ON terminal_types(channel_code);
CREATE INDEX IF NOT EXISTS idx_terminal_types_status ON terminal_types(status);

-- 注释
COMMENT ON TABLE terminal_types IS '终端类型表';
COMMENT ON COLUMN terminal_types.channel_id IS '所属通道ID';
COMMENT ON COLUMN terminal_types.channel_code IS '通道编码';
COMMENT ON COLUMN terminal_types.brand_code IS '品牌编码';
COMMENT ON COLUMN terminal_types.brand_name IS '品牌名称';
COMMENT ON COLUMN terminal_types.model_code IS '型号编码';
COMMENT ON COLUMN terminal_types.model_name IS '型号名称';
COMMENT ON COLUMN terminal_types.description IS '描述';
COMMENT ON COLUMN terminal_types.status IS '状态：1启用 0禁用';
