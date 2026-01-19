-- 收享付 - 消息通知表

CREATE TABLE IF NOT EXISTS messages (
    id              BIGSERIAL PRIMARY KEY,
    agent_id        BIGINT NOT NULL,                -- 接收人（代理商ID）

    -- 消息内容
    message_type    SMALLINT NOT NULL,              -- 1:分润 2:激活奖励 3:押金返现 4:流量返现 5:退款撤销 6:系统公告
    title           VARCHAR(64) NOT NULL,           -- 消息标题
    content         TEXT,                           -- 消息内容

    -- 状态
    is_read         BOOLEAN DEFAULT FALSE,          -- 是否已读
    is_pushed       BOOLEAN DEFAULT FALSE,          -- 是否已推送到APP

    -- 关联信息
    related_id      BIGINT,                         -- 关联ID（交易ID/分润ID等）
    related_type    VARCHAR(32),                    -- 关联类型: transaction, profit_record, device_fee等

    -- 时间
    expire_at       TIMESTAMPTZ,                    -- 过期时间（3天后自动清理）
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

-- 创建索引（复合索引优化APP消息列表查询）
CREATE INDEX IF NOT EXISTS idx_messages_agent_read ON messages(agent_id, is_read, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_messages_expire ON messages(expire_at) WHERE expire_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_messages_type ON messages(message_type);

-- 添加注释
COMMENT ON TABLE messages IS '消息通知表 - 存储代理商的站内消息';
COMMENT ON COLUMN messages.message_type IS '消息类型: 1-分润, 2-激活奖励, 3-押金返现, 4-流量返现, 5-退款撤销, 6-系统公告';
COMMENT ON COLUMN messages.expire_at IS '过期时间，超过此时间的消息可被定时任务清理';
