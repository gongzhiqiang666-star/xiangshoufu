import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../../../core/utils/format_utils.dart';
import '../../data/models/deduction_model.dart';

/// 代扣计划卡片组件
class DeductionCard extends StatelessWidget {
  final DeductionPlan plan;
  final VoidCallback? onTap;
  final VoidCallback? onAccept;
  final VoidCallback? onReject;
  final VoidCallback? onPause;
  final VoidCallback? onResume;
  final VoidCallback? onCancel;
  final bool showDeductor; // 是否显示扣款方（我接收的代扣时显示）

  const DeductionCard({
    super.key,
    required this.plan,
    this.onTap,
    this.onAccept,
    this.onReject,
    this.onPause,
    this.onResume,
    this.onCancel,
    this.showDeductor = false,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        margin: const EdgeInsets.only(bottom: AppSpacing.cardGap),
        padding: const EdgeInsets.all(AppSpacing.cardPadding),
        decoration: BoxDecoration(
          color: AppColors.cardBg,
          borderRadius: BorderRadius.circular(12),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 8,
              offset: const Offset(0, 2),
            ),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildHeader(),
            const SizedBox(height: AppSpacing.sm),
            _buildAmountRow(),
            const SizedBox(height: AppSpacing.sm),
            _buildProgress(),
            const SizedBox(height: AppSpacing.sm),
            _buildFooter(),
            // 待接收、进行中、已暂停状态显示操作按钮
            if (plan.status == 0 || plan.status == 1 || plan.status == 3) ...[
              const SizedBox(height: AppSpacing.sm),
              _buildActions(),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Row(
      children: [
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                plan.deducteeName,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              const SizedBox(height: 2),
              Text(
                plan.planNo,
                style: const TextStyle(
                  fontSize: 12,
                  color: AppColors.textTertiary,
                ),
              ),
            ],
          ),
        ),
        _buildTypeTag(),
        const SizedBox(width: 8),
        _buildStatusTag(),
      ],
    );
  }

  Widget _buildTypeTag() {
    Color bgColor;
    Color textColor;

    switch (plan.planType) {
      case 1:
        bgColor = AppColors.primary.withOpacity(0.1);
        textColor = AppColors.primary;
        break;
      case 2:
        bgColor = AppColors.success.withOpacity(0.1);
        textColor = AppColors.success;
        break;
      case 3:
        bgColor = AppColors.warning.withOpacity(0.1);
        textColor = AppColors.warning;
        break;
      default:
        bgColor = AppColors.textTertiary.withOpacity(0.1);
        textColor = AppColors.textTertiary;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        plan.typeEnum.label,
        style: TextStyle(
          fontSize: 10,
          color: textColor,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Widget _buildStatusTag() {
    Color bgColor;
    Color textColor;

    switch (plan.status) {
      case 0: // 待接收
        bgColor = AppColors.info.withOpacity(0.1);
        textColor = AppColors.info;
        break;
      case 1: // 进行中
        bgColor = AppColors.primary.withOpacity(0.1);
        textColor = AppColors.primary;
        break;
      case 2: // 已完成
        bgColor = AppColors.success.withOpacity(0.1);
        textColor = AppColors.success;
        break;
      case 3: // 已暂停
        bgColor = AppColors.warning.withOpacity(0.1);
        textColor = AppColors.warning;
        break;
      case 4: // 已取消
        bgColor = AppColors.textTertiary.withOpacity(0.1);
        textColor = AppColors.textTertiary;
        break;
      case 5: // 已拒绝
        bgColor = AppColors.danger.withOpacity(0.1);
        textColor = AppColors.danger;
        break;
      default:
        bgColor = AppColors.textTertiary.withOpacity(0.1);
        textColor = AppColors.textTertiary;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: bgColor,
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        plan.statusEnum.label,
        style: TextStyle(
          fontSize: 12,
          color: textColor,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }

  Widget _buildAmountRow() {
    return Row(
      children: [
        Expanded(
          child: _buildAmountItem('代扣总额', plan.totalAmountYuan),
        ),
        Expanded(
          child: _buildAmountItem(
            '已扣金额',
            plan.deductedAmountYuan,
            color: AppColors.success,
          ),
        ),
        Expanded(
          child: _buildAmountItem(
            '剩余待扣',
            plan.remainingAmountYuan,
            color: AppColors.danger,
          ),
        ),
      ],
    );
  }

  Widget _buildAmountItem(String label, double amount, {Color? color}) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textTertiary,
          ),
        ),
        const SizedBox(height: 2),
        Text(
          '¥${FormatUtils.formatYuan(amount)}',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w600,
            color: color ?? AppColors.textPrimary,
          ),
        ),
      ],
    );
  }

  Widget _buildProgress() {
    return Row(
      children: [
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    '扣款进度',
                    style: const TextStyle(
                      fontSize: 12,
                      color: AppColors.textTertiary,
                    ),
                  ),
                  Text(
                    '${plan.progress.toStringAsFixed(1)}%',
                    style: const TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                      color: AppColors.textPrimary,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 4),
              ClipRRect(
                borderRadius: BorderRadius.circular(4),
                child: LinearProgressIndicator(
                  value: plan.progress / 100,
                  minHeight: 6,
                  backgroundColor: AppColors.border,
                  valueColor: AlwaysStoppedAnimation<Color>(
                    plan.progress >= 100 ? AppColors.success : AppColors.primary,
                  ),
                ),
              ),
            ],
          ),
        ),
        const SizedBox(width: AppSpacing.md),
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: AppColors.background,
            borderRadius: BorderRadius.circular(8),
          ),
          child: Column(
            children: [
              Text(
                '${plan.currentPeriod}/${plan.totalPeriods}',
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              const Text(
                '期',
                style: TextStyle(
                  fontSize: 10,
                  color: AppColors.textTertiary,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildFooter() {
    return Row(
      children: [
        Icon(
          Icons.calendar_today,
          size: 14,
          color: AppColors.textTertiary,
        ),
        const SizedBox(width: 4),
        Text(
          plan.createdAt.substring(0, 10),
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textTertiary,
          ),
        ),
        const Spacer(),
        Text(
          '每期 ¥${FormatUtils.formatYuan(plan.periodAmountYuan)}',
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildActions() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.end,
      children: [
        // 待接收状态：显示接收和拒绝按钮
        if (plan.status == 0) ...[
          TextButton(
            onPressed: onReject,
            style: TextButton.styleFrom(
              foregroundColor: AppColors.danger,
              padding: const EdgeInsets.symmetric(horizontal: 12),
            ),
            child: const Text('拒绝'),
          ),
          ElevatedButton(
            onPressed: onAccept,
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.success,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              minimumSize: const Size(0, 32),
            ),
            child: const Text('接收确认'),
          ),
        ],
        // 进行中状态：显示暂停和取消按钮
        if (plan.status == 1) ...[
          TextButton(
            onPressed: onPause,
            style: TextButton.styleFrom(
              foregroundColor: AppColors.warning,
              padding: const EdgeInsets.symmetric(horizontal: 12),
            ),
            child: const Text('暂停'),
          ),
          TextButton(
            onPressed: onCancel,
            style: TextButton.styleFrom(
              foregroundColor: AppColors.danger,
              padding: const EdgeInsets.symmetric(horizontal: 12),
            ),
            child: const Text('取消'),
          ),
        ],
        // 已暂停状态：显示恢复和取消按钮
        if (plan.status == 3) ...[
          TextButton(
            onPressed: onResume,
            style: TextButton.styleFrom(
              foregroundColor: AppColors.success,
              padding: const EdgeInsets.symmetric(horizontal: 12),
            ),
            child: const Text('恢复'),
          ),
          TextButton(
            onPressed: onCancel,
            style: TextButton.styleFrom(
              foregroundColor: AppColors.danger,
              padding: const EdgeInsets.symmetric(horizontal: 12),
            ),
            child: const Text('取消'),
          ),
        ],
      ],
    );
  }
}
