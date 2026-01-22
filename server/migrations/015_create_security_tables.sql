-- 审计日志表
-- 满足三级等保审计要求，记录敏感操作
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    log_type SMALLINT NOT NULL,           -- 日志类型
    log_level SMALLINT NOT NULL DEFAULT 1, -- 日志级别：1信息 2警告 3严重
    user_id BIGINT,                        -- 操作用户ID
    username VARCHAR(64),                  -- 操作用户名
    agent_id BIGINT,                       -- 代理商ID
    agent_name VARCHAR(128),               -- 代理商名称
    target_type VARCHAR(32),               -- 操作目标类型
    target_id BIGINT,                      -- 操作目标ID
    target_name VARCHAR(128),              -- 操作目标名称
    action VARCHAR(64) NOT NULL,           -- 操作动作
    description VARCHAR(512),              -- 操作描述
    old_value TEXT,                        -- 变更前的值（JSON）
    new_value TEXT,                        -- 变更后的值（JSON）
    ip VARCHAR(64),                        -- 操作IP
    user_agent VARCHAR(256),               -- 用户代理
    request_path VARCHAR(256),             -- 请求路径
    request_method VARCHAR(16),            -- 请求方法
    result SMALLINT DEFAULT 1,             -- 结果：1成功 2失败
    error_msg VARCHAR(512),                -- 错误信息
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_audit_logs_log_type ON audit_logs(log_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_agent_id ON audit_logs(agent_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_target_id ON audit_logs(target_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_audit_logs_ip ON audit_logs(ip);

-- IP黑名单表
CREATE TABLE IF NOT EXISTS ip_blacklist (
    id BIGSERIAL PRIMARY KEY,
    ip VARCHAR(64) NOT NULL UNIQUE,        -- IP地址
    reason VARCHAR(256),                   -- 封禁原因
    blocked_by BIGINT,                     -- 操作人
    expires_at TIMESTAMP WITH TIME ZONE,   -- 过期时间
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ip_blacklist_ip ON ip_blacklist(ip);
CREATE INDEX IF NOT EXISTS idx_ip_blacklist_expires_at ON ip_blacklist(expires_at);

-- 登录尝试记录表（用于登录失败锁定）
CREATE TABLE IF NOT EXISTS login_attempts (
    id BIGSERIAL PRIMARY KEY,
    identifier VARCHAR(128) NOT NULL,      -- 标识符（用户名或IP）
    identifier_type VARCHAR(16) NOT NULL,  -- 类型：username/ip
    fail_count INT NOT NULL DEFAULT 0,     -- 失败次数
    last_fail_at TIMESTAMP WITH TIME ZONE, -- 最后失败时间
    locked_until TIMESTAMP WITH TIME ZONE, -- 锁定到期时间
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_login_attempts_identifier ON login_attempts(identifier, identifier_type);

-- 安全配置表
CREATE TABLE IF NOT EXISTS security_configs (
    id BIGSERIAL PRIMARY KEY,
    config_key VARCHAR(64) NOT NULL UNIQUE, -- 配置键
    config_value TEXT NOT NULL,             -- 配置值（JSON）
    description VARCHAR(256),               -- 描述
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 插入默认安全配置
INSERT INTO security_configs (config_key, config_value, description) VALUES
    ('password_policy', '{"min_length":8,"max_length":32,"require_upper":true,"require_lower":true,"require_digit":true,"require_special":false}', '密码策略'),
    ('login_lock_policy', '{"max_fails":5,"lock_duration_minutes":30}', '登录锁定策略'),
    ('session_policy', '{"access_token_hours":2,"refresh_token_days":7}', '会话策略'),
    ('rate_limit_policy', '{"global_rate":1000,"ip_rate":100}', '限流策略')
ON CONFLICT (config_key) DO NOTHING;

COMMENT ON TABLE audit_logs IS '审计日志表 - 记录敏感操作，满足三级等保审计要求';
COMMENT ON TABLE ip_blacklist IS 'IP黑名单表';
COMMENT ON TABLE login_attempts IS '登录尝试记录表';
COMMENT ON TABLE security_configs IS '安全配置表';
