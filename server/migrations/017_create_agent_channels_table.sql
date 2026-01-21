-- 017_create_agent_channels_table.sql
-- 代理商通道启用配置表
-- 记录代理商可以使用哪些通道，APP端只显示已开通的通道

-- 代理商通道配置表
CREATE TABLE IF NOT EXISTS agent_channels (
    id                  BIGSERIAL PRIMARY KEY,
    agent_id            BIGINT NOT NULL,                    -- 代理商ID
    channel_id          BIGINT NOT NULL,                    -- 通道ID
    is_enabled          BOOLEAN DEFAULT TRUE,               -- 是否启用
    is_visible          BOOLEAN DEFAULT TRUE,               -- 对代理商是否可见
    enabled_at          TIMESTAMPTZ,                        -- 启用时间
    disabled_at         TIMESTAMPTZ,                        -- 禁用时间
    enabled_by          BIGINT,                             -- 启用人（用户ID）
    disabled_by         BIGINT,                             -- 禁用人（用户ID）
    remark              TEXT,                               -- 备注
    created_at          TIMESTAMPTZ DEFAULT NOW(),
    updated_at          TIMESTAMPTZ DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_agent_channels_agent ON agent_channels(agent_id);
CREATE INDEX idx_agent_channels_channel ON agent_channels(channel_id);
CREATE UNIQUE INDEX idx_agent_channels_unique ON agent_channels(agent_id, channel_id);
CREATE INDEX idx_agent_channels_enabled ON agent_channels(agent_id, is_enabled) WHERE is_enabled = TRUE;

-- 注释
COMMENT ON TABLE agent_channels IS '代理商通道配置表';
COMMENT ON COLUMN agent_channels.agent_id IS '代理商ID';
COMMENT ON COLUMN agent_channels.channel_id IS '通道ID';
COMMENT ON COLUMN agent_channels.is_enabled IS '是否启用';
COMMENT ON COLUMN agent_channels.is_visible IS '对代理商是否可见（APP端是否显示）';
COMMENT ON COLUMN agent_channels.enabled_at IS '启用时间';
COMMENT ON COLUMN agent_channels.disabled_at IS '禁用时间';
COMMENT ON COLUMN agent_channels.enabled_by IS '启用人（用户ID）';
COMMENT ON COLUMN agent_channels.disabled_by IS '禁用人（用户ID）';
COMMENT ON COLUMN agent_channels.remark IS '备注';
