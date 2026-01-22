-- 021_alter_terminal_distributes_add_goods_deduction.sql
-- 终端划拨表新增货款代扣相关字段

-- 新增扣款来源字段（货款代扣来源: 1=分润 2=服务费 3=两者）
ALTER TABLE terminal_distributes
ADD COLUMN IF NOT EXISTS deduction_source SMALLINT DEFAULT 3;

-- 新增关联货款代扣ID字段
ALTER TABLE terminal_distributes
ADD COLUMN IF NOT EXISTS goods_deduction_id BIGINT REFERENCES goods_deductions(id);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_terminal_distributes_goods_deduction ON terminal_distributes(goods_deduction_id);

-- 添加字段注释
COMMENT ON COLUMN terminal_distributes.deduction_source IS '货款代扣来源: 1=分润钱包 2=服务费钱包 3=两者都扣';
COMMENT ON COLUMN terminal_distributes.goods_deduction_id IS '关联货款代扣ID';
