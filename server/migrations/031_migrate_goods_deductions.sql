-- 031_migrate_goods_deductions.sql
-- 迁移现有货款代扣数据到统一的代扣计划表
-- 执行时间：业务低峰期
-- 注意：此脚本需要在 030_add_deduction_freeze.sql 执行后运行

BEGIN;

-- 1. 将进行中和待接收的货款代扣迁移到 deduction_plans
-- 注意：只迁移未完成的记录，已完成的保留在原表供历史查询
INSERT INTO deduction_plans (
    plan_no,
    deductor_id,
    deductee_id,
    plan_type,
    total_amount,
    deducted_amount,
    remaining_amount,
    total_periods,
    current_period,
    period_amount,
    status,
    need_accept,
    accepted_at,
    deduction_source,
    remark,
    created_by,
    created_at,
    updated_at
)
SELECT
    deduction_no,
    from_agent_id,
    to_agent_id,
    1, -- plan_type: 1=货款代扣
    total_amount,
    deducted_amount,
    remaining_amount,
    total_periods,
    current_period,
    CASE
        WHEN total_periods > 0 THEN total_amount / total_periods
        ELSE total_amount
    END, -- period_amount
    CASE status
        WHEN 1 THEN 0  -- 待接收 -> 待接收
        WHEN 2 THEN 1  -- 进行中 -> 进行中
        WHEN 4 THEN 5  -- 已拒绝 -> 已拒绝
        ELSE status
    END, -- status
    TRUE, -- need_accept: 货款代扣需要确认
    accepted_at,
    COALESCE(deduction_source, 3), -- deduction_source: 默认两者
    remark,
    created_by,
    created_at,
    updated_at
FROM goods_deductions
WHERE status IN (1, 2) -- 只迁移待接收和进行中的
ON CONFLICT (plan_no) DO NOTHING; -- 避免重复迁移

-- 2. 为迁移的计划生成代扣记录（如果原表有明细记录）
-- 这里假设原货款代扣表有对应的明细记录表
-- 如果没有明细记录表，则需要根据计划生成记录

-- 3. 更新原货款代扣表状态标记为已迁移（可选）
-- 添加迁移标记字段（如果需要）
-- ALTER TABLE goods_deductions ADD COLUMN IF NOT EXISTS migrated_to_plan_id BIGINT;

-- 4. 记录迁移日志
DO $$
DECLARE
    migrated_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO migrated_count
    FROM deduction_plans dp
    INNER JOIN goods_deductions gd ON dp.plan_no = gd.deduction_no
    WHERE gd.status IN (1, 2);

    RAISE NOTICE '已迁移 % 条货款代扣记录到统一代扣计划表', migrated_count;
END $$;

COMMIT;

-- 验证迁移结果
-- SELECT
--     'goods_deductions' as source,
--     status,
--     COUNT(*) as count
-- FROM goods_deductions
-- WHERE status IN (1, 2)
-- GROUP BY status
-- UNION ALL
-- SELECT
--     'deduction_plans' as source,
--     status,
--     COUNT(*) as count
-- FROM deduction_plans
-- WHERE plan_type = 1
-- GROUP BY status;
