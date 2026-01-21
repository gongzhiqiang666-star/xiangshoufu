-- 终端回拨记录表
CREATE TABLE IF NOT EXISTS terminal_recalls (
    id BIGSERIAL PRIMARY KEY,
    recall_no VARCHAR(64) NOT NULL UNIQUE,        -- 回拨单号
    from_agent_id BIGINT NOT NULL,                -- 回拨方代理商（当前持有者）
    to_agent_id BIGINT NOT NULL,                  -- 接收方代理商（上级）
    terminal_sn VARCHAR(50) NOT NULL,             -- 终端SN
    channel_id BIGINT NOT NULL,                   -- 通道ID
    is_cross_level BOOLEAN DEFAULT FALSE,         -- 是否跨级回拨
    cross_level_path VARCHAR(500),                -- 跨级路径
    status SMALLINT DEFAULT 1,                    -- 1:待确认 2:已确认 3:已拒绝 4:已取消
    source SMALLINT NOT NULL,                     -- 1:APP 2:PC
    remark VARCHAR(255),
    created_by BIGINT,
    confirmed_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_terminal_recalls_from_agent ON terminal_recalls(from_agent_id);
CREATE INDEX IF NOT EXISTS idx_terminal_recalls_to_agent ON terminal_recalls(to_agent_id);
CREATE INDEX IF NOT EXISTS idx_terminal_recalls_terminal_sn ON terminal_recalls(terminal_sn);
CREATE INDEX IF NOT EXISTS idx_terminal_recalls_status ON terminal_recalls(status);
CREATE INDEX IF NOT EXISTS idx_terminal_recalls_created_at ON terminal_recalls(created_at);

COMMENT ON TABLE terminal_recalls IS '终端回拨记录表';
COMMENT ON COLUMN terminal_recalls.recall_no IS '回拨单号';
COMMENT ON COLUMN terminal_recalls.from_agent_id IS '回拨方代理商（当前持有者）';
COMMENT ON COLUMN terminal_recalls.to_agent_id IS '接收方代理商（上级）';
COMMENT ON COLUMN terminal_recalls.terminal_sn IS '终端SN';
COMMENT ON COLUMN terminal_recalls.channel_id IS '通道ID';
COMMENT ON COLUMN terminal_recalls.is_cross_level IS '是否跨级回拨';
COMMENT ON COLUMN terminal_recalls.cross_level_path IS '跨级路径';
COMMENT ON COLUMN terminal_recalls.status IS '1:待确认 2:已确认 3:已拒绝 4:已取消';
COMMENT ON COLUMN terminal_recalls.source IS '1:APP 2:PC';

-- 终端入库记录表
CREATE TABLE IF NOT EXISTS terminal_import_records (
    id BIGSERIAL PRIMARY KEY,
    import_no VARCHAR(64) NOT NULL UNIQUE,        -- 入库批次号
    channel_id BIGINT NOT NULL,                   -- 通道ID
    channel_code VARCHAR(32),                     -- 通道编码
    brand_code VARCHAR(32),                       -- 品牌编码
    model_code VARCHAR(32),                       -- 型号编码
    total_count INT NOT NULL,                     -- 导入总数
    success_count INT DEFAULT 0,                  -- 成功数
    failed_count INT DEFAULT 0,                   -- 失败数
    failed_sns TEXT,                              -- 失败SN列表（JSON数组）
    owner_agent_id BIGINT NOT NULL,               -- 入库代理商
    created_by BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_terminal_import_records_owner ON terminal_import_records(owner_agent_id);
CREATE INDEX IF NOT EXISTS idx_terminal_import_records_channel ON terminal_import_records(channel_id);
CREATE INDEX IF NOT EXISTS idx_terminal_import_records_created_at ON terminal_import_records(created_at);

COMMENT ON TABLE terminal_import_records IS '终端入库记录表';
COMMENT ON COLUMN terminal_import_records.import_no IS '入库批次号';
COMMENT ON COLUMN terminal_import_records.channel_id IS '通道ID';
COMMENT ON COLUMN terminal_import_records.channel_code IS '通道编码';
COMMENT ON COLUMN terminal_import_records.brand_code IS '品牌编码';
COMMENT ON COLUMN terminal_import_records.model_code IS '型号编码';
COMMENT ON COLUMN terminal_import_records.total_count IS '导入总数';
COMMENT ON COLUMN terminal_import_records.success_count IS '成功数';
COMMENT ON COLUMN terminal_import_records.failed_count IS '失败数';
COMMENT ON COLUMN terminal_import_records.failed_sns IS '失败SN列表（JSON数组）';
COMMENT ON COLUMN terminal_import_records.owner_agent_id IS '入库代理商';
