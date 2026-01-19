-- 用户认证相关表
-- 执行: psql -d xiangshoufu -f migrations/009_create_users_table.sql

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    password        VARCHAR(100) NOT NULL,
    salt            VARCHAR(32) NOT NULL,
    agent_id        BIGINT REFERENCES agents(id),
    role_type       SMALLINT DEFAULT 1,  -- 1普通用户 2管理员
    status          SMALLINT DEFAULT 1,  -- 1正常 2禁用
    last_login_at   TIMESTAMPTZ,
    last_login_ip   VARCHAR(50),
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_agent ON users(agent_id);

-- 刷新令牌表
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token           VARCHAR(64) NOT NULL UNIQUE,
    expires_at      TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires ON refresh_tokens(expires_at);

-- 登录日志表
CREATE TABLE IF NOT EXISTS login_logs (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users(id),
    username        VARCHAR(50) NOT NULL,
    login_ip        VARCHAR(50),
    user_agent      VARCHAR(255),
    status          SMALLINT NOT NULL,  -- 1成功 2失败
    fail_msg        VARCHAR(255),
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_login_logs_user ON login_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_login_logs_time ON login_logs(created_at);

-- 添加表注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON TABLE refresh_tokens IS '刷新令牌表';
COMMENT ON TABLE login_logs IS '登录日志表';

COMMENT ON COLUMN users.role_type IS '角色类型: 1普通用户 2管理员';
COMMENT ON COLUMN users.status IS '状态: 1正常 2禁用';
COMMENT ON COLUMN login_logs.status IS '登录状态: 1成功 2失败';

-- 创建默认管理员用户（密码: admin123）
-- 注意：生产环境请立即修改密码
-- 密码哈希: sha256("admin123" + salt)
INSERT INTO users (username, password, salt, agent_id, role_type, status)
VALUES (
    'admin',
    'e99a18c428cb38d5f260853678922e03abd82979c5b3f4c0bdb4d94e41827c4d', -- 示例哈希，需要实际生成
    'default_salt_12345678',
    NULL,
    2,  -- 管理员
    1   -- 正常
) ON CONFLICT (username) DO NOTHING;
