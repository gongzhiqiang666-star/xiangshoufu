-- 商户表新增字段迁移
-- 创建日期: 2026-01-20
-- 说明: 添加商户类型、直营标识、激活时间、登记手机号等字段

-- 商户类型字段
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS merchant_type VARCHAR(20) DEFAULT 'normal';

-- 是否直营商户标识
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS is_direct BOOLEAN DEFAULT TRUE;

-- 激活时间(首次交易时间)
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS activated_at TIMESTAMPTZ;

-- 登记手机号（加密存储）
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS registered_phone VARCHAR(100);

-- 登记备注
ALTER TABLE merchants ADD COLUMN IF NOT EXISTS register_remark VARCHAR(500);

-- 添加索引以优化查询
CREATE INDEX IF NOT EXISTS idx_merchants_type ON merchants(merchant_type);
CREATE INDEX IF NOT EXISTS idx_merchants_direct ON merchants(is_direct);
CREATE INDEX IF NOT EXISTS idx_merchants_activated ON merchants(activated_at);

-- 添加字段注释
COMMENT ON COLUMN merchants.merchant_type IS '商户类型: loyal(忠诚)/quality(优质)/potential(潜力)/normal(一般)/low_active(低活跃)/inactive(无交易)';
COMMENT ON COLUMN merchants.is_direct IS '是否直营商户: true=直营(代理商直接拓展), false=团队(下级代理商拓展)';
COMMENT ON COLUMN merchants.activated_at IS '激活时间(首次交易时间)';
COMMENT ON COLUMN merchants.registered_phone IS '登记手机号(加密存储,用于联系商户)';
COMMENT ON COLUMN merchants.register_remark IS '登记备注';
