-- 030_update_merchant_types.sql
-- 商户类型分类调整：从6档改为5档
-- 原分类: loyal, quality, potential, normal, low_active, inactive
-- 新分类: quality, medium, normal, warning, churned

-- 开始事务
BEGIN;

-- 1. 更新商户类型映射
-- loyal -> quality (忠诚商户 -> 优质商户)
UPDATE merchants SET merchant_type = 'quality' WHERE merchant_type = 'loyal';

-- potential -> medium (潜力商户 -> 中等商户)
UPDATE merchants SET merchant_type = 'medium' WHERE merchant_type = 'potential';

-- low_active -> normal (低活跃 -> 普通商户，因为有交易)
UPDATE merchants SET merchant_type = 'normal' WHERE merchant_type = 'low_active';

-- inactive -> warning (30天无交易 -> 预警商户)
UPDATE merchants SET merchant_type = 'warning' WHERE merchant_type = 'inactive';

-- 2. 添加注释说明新的分类标准
COMMENT ON COLUMN merchants.merchant_type IS '商户类型（5档分类）: quality=优质(月均≥5万), medium=中等(3-5万), normal=普通(<3万), warning=预警(30天无交易), churned=流失(60天无交易)';

-- 3. 创建索引优化查询（如果不存在）
CREATE INDEX IF NOT EXISTS idx_merchants_merchant_type ON merchants(merchant_type);

-- 提交事务
COMMIT;

-- 验证迁移结果
SELECT merchant_type, COUNT(*) as count
FROM merchants
GROUP BY merchant_type
ORDER BY count DESC;
