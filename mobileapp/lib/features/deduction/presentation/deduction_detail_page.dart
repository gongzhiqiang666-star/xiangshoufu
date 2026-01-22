import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../data/models/deduction_model.dart';
import 'providers/deduction_provider.dart';

/// 代扣计划详情页
class DeductionDetailPage extends ConsumerWidget {
  final int id;

  const DeductionDetailPage({
    super.key,
    required this.id,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(deductionPlanDetailProvider(id));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('代扣详情'),
        centerTitle: true,
      ),
      body: detailAsync.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
              const SizedBox(height: AppSpacing.md),
              Text('加载失败: $error'),
              const SizedBox(height: AppSpacing.md),
              ElevatedButton(
                onPressed: () => ref.invalidate(deductionPlanDetailProvider(id)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
        data: (detail) => _buildContent(context, ref, detail),
      ),
    );
  }

  Widget _buildContent(BuildContext context, WidgetRef ref, DeductionPlanDetail detail) {
    return SingleChildScrollView(
      child: Column(
        children: [
          _buildStatusCard(detail),
          _buildAmountCard(detail),
          _buildInfoCard(detail),
          _buildRecordsCard(detail),
          if (detail.status == 1 || detail.status == 3)
            _buildActionButtons(context, ref, detail),
          const SizedBox(height: AppSpacing.xl),
        ],
      ),
    );
  }

  Widget _buildStatusCard(DeductionPlanDetail detail) {
    Color statusColor;
    IconData statusIcon;

    switch (detail.status) {
      case 1:
        statusColor = AppColors.primary;
        statusIcon = Icons.sync;
        break;
      case 2:
        statusColor = AppColors.success;
        statusIcon = Icons.check_circle;
        break;
      case 3:
        statusColor = AppColors.warning;
        statusIcon = Icons.pause_circle;
        break;
      case 4:
        statusColor = AppColors.textTertiary;
        statusIcon = Icons.cancel;
        break;
      default:
        statusColor = AppColors.textTertiary;
        statusIcon = Icons.help_outline;
    }

    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        color: AppColors.cardBg,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Container(
            width: 56,
            height: 56,
            decoration: BoxDecoration(
              color: statusColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(28),
            ),
            child: Icon(statusIcon, color: statusColor, size: 28),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Text(
                      detail.statusEnum.label,
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: statusColor,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: AppColors.primary.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        detail.typeEnum.label,
                        style: const TextStyle(
                          fontSize: 10,
                          color: AppColors.primary,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 4),
                Text(
                  '代扣总额 ¥${FormatUtils.formatYuan(detail.totalAmountYuan)}',
                  style: const TextStyle(
                    fontSize: 14,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          _buildProgressCircle(detail),
        ],
      ),
    );
  }

  Widget _buildProgressCircle(DeductionPlanDetail detail) {
    return SizedBox(
      width: 80,
      height: 80,
      child: Stack(
        alignment: Alignment.center,
        children: [
          SizedBox(
            width: 80,
            height: 80,
            child: CircularProgressIndicator(
              value: detail.progress / 100,
              strokeWidth: 6,
              backgroundColor: AppColors.border,
              valueColor: AlwaysStoppedAnimation<Color>(
                detail.progress >= 100 ? AppColors.success : AppColors.primary,
              ),
            ),
          ),
          Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                '${detail.currentPeriod}/${detail.totalPeriods}',
                style: const TextStyle(
                  fontSize: 14,
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
        ],
      ),
    );
  }

  Widget _buildAmountCard(DeductionPlanDetail detail) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.cardBg,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: _buildAmountItem('代扣总额', detail.totalAmountYuan),
          ),
          Container(
            width: 1,
            height: 40,
            color: AppColors.divider,
          ),
          Expanded(
            child: _buildAmountItem(
              '已扣金额',
              detail.deductedAmountYuan,
              color: AppColors.success,
            ),
          ),
          Container(
            width: 1,
            height: 40,
            color: AppColors.divider,
          ),
          Expanded(
            child: _buildAmountItem(
              '剩余待扣',
              detail.remainingAmountYuan,
              color: AppColors.danger,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAmountItem(String label, double amount, {Color? color}) {
    return Column(
      children: [
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textTertiary,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          '¥${FormatUtils.formatYuan(amount)}',
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.w600,
            color: color ?? AppColors.textPrimary,
          ),
        ),
      ],
    );
  }

  Widget _buildInfoCard(DeductionPlanDetail detail) {
    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.cardBg,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '基本信息',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: AppSpacing.md),
          _buildInfoRow('计划编号', detail.planNo),
          _buildInfoRow('扣款方', detail.deductorName),
          _buildInfoRow('被扣款方', detail.deducteeName),
          _buildInfoRow('总期数', '${detail.totalPeriods} 期'),
          _buildInfoRow('当前期数', '第 ${detail.currentPeriod} 期'),
          _buildInfoRow('每期金额', '¥${FormatUtils.formatYuan(detail.periodAmountYuan)}'),
          _buildInfoRow('创建时间', detail.createdAt.substring(0, 16)),
          if (detail.completedAt != null)
            _buildInfoRow('完成时间', detail.completedAt!.substring(0, 16)),
          if (detail.remark != null && detail.remark!.isNotEmpty)
            _buildInfoRow('备注', detail.remark!),
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(
              label,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textSecondary,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontSize: 14,
                color: AppColors.textPrimary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRecordsCard(DeductionPlanDetail detail) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.cardBg,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                '扣款记录',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              Text(
                '共 ${detail.records.length} 条',
                style: const TextStyle(
                  fontSize: 14,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.sm),
          if (detail.records.isEmpty)
            const Padding(
              padding: EdgeInsets.symmetric(vertical: AppSpacing.lg),
              child: Center(
                child: Text(
                  '暂无扣款记录',
                  style: TextStyle(
                    fontSize: 14,
                    color: AppColors.textTertiary,
                  ),
                ),
              ),
            )
          else
            ...detail.records.map((record) => _buildRecordItem(record)),
        ],
      ),
    );
  }

  Widget _buildRecordItem(DeductionRecord record) {
    Color statusColor;
    switch (record.status) {
      case 0:
        statusColor = AppColors.textTertiary;
        break;
      case 1:
        statusColor = AppColors.success;
        break;
      case 2:
        statusColor = AppColors.warning;
        break;
      case 3:
        statusColor = AppColors.danger;
        break;
      default:
        statusColor = AppColors.textTertiary;
    }

    return Container(
      padding: const EdgeInsets.symmetric(vertical: AppSpacing.sm),
      decoration: const BoxDecoration(
        border: Border(
          bottom: BorderSide(color: AppColors.divider, width: 0.5),
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: statusColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(18),
            ),
            child: Center(
              child: Text(
                '${record.periodNum}',
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: statusColor,
                ),
              ),
            ),
          ),
          const SizedBox(width: AppSpacing.sm),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '第 ${record.periodNum} 期',
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w500,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: statusColor.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        record.statusEnum.label,
                        style: TextStyle(
                          fontSize: 10,
                          color: statusColor,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 2),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '应扣: ¥${FormatUtils.formatYuan(record.amountYuan)}',
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.textSecondary,
                      ),
                    ),
                    Text(
                      '实扣: ¥${FormatUtils.formatYuan(record.actualAmountYuan)}',
                      style: TextStyle(
                        fontSize: 12,
                        fontWeight: FontWeight.w500,
                        color: record.actualAmount > 0 ? AppColors.success : AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
                if (record.deductedAt != null)
                  Padding(
                    padding: const EdgeInsets.only(top: 2),
                    child: Text(
                      '扣款时间: ${record.deductedAt!.substring(0, 16)}',
                      style: const TextStyle(
                        fontSize: 11,
                        color: AppColors.textTertiary,
                      ),
                    ),
                  ),
                if (record.failReason != null && record.failReason!.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 2),
                    child: Text(
                      '失败原因: ${record.failReason}',
                      style: const TextStyle(
                        fontSize: 11,
                        color: AppColors.danger,
                      ),
                    ),
                  ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActionButtons(BuildContext context, WidgetRef ref, DeductionPlanDetail detail) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Row(
        children: [
          if (detail.status == 1)
            Expanded(
              child: OutlinedButton(
                onPressed: () => _handlePause(context, ref, detail),
                style: OutlinedButton.styleFrom(
                  foregroundColor: AppColors.warning,
                  side: const BorderSide(color: AppColors.warning),
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
                child: const Text('暂停'),
              ),
            ),
          if (detail.status == 3)
            Expanded(
              child: ElevatedButton(
                onPressed: () => _handleResume(context, ref, detail),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppColors.success,
                  foregroundColor: Colors.white,
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
                child: const Text('恢复'),
              ),
            ),
          if (detail.status == 1 || detail.status == 3) ...[
            const SizedBox(width: AppSpacing.md),
            Expanded(
              child: OutlinedButton(
                onPressed: () => _handleCancel(context, ref, detail),
                style: OutlinedButton.styleFrom(
                  foregroundColor: AppColors.danger,
                  side: const BorderSide(color: AppColors.danger),
                  padding: const EdgeInsets.symmetric(vertical: 14),
                ),
                child: const Text('取消'),
              ),
            ),
          ],
        ],
      ),
    );
  }

  Future<void> _handlePause(BuildContext context, WidgetRef ref, DeductionPlanDetail detail) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('暂停代扣'),
        content: const Text('确定要暂停此代扣计划吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.warning,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定暂停'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final service = ref.read(deductionServiceProvider);
      try {
        await service.pausePlan(detail.id);
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('暂停成功'),
              backgroundColor: AppColors.success,
            ),
          );
          ref.invalidate(deductionPlanDetailProvider(id));
          ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true);
        }
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('操作失败: $e'),
              backgroundColor: AppColors.danger,
            ),
          );
        }
      }
    }
  }

  Future<void> _handleResume(BuildContext context, WidgetRef ref, DeductionPlanDetail detail) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('恢复代扣'),
        content: const Text('确定要恢复此代扣计划吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.success,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定恢复'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final service = ref.read(deductionServiceProvider);
      try {
        await service.resumePlan(detail.id);
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('恢复成功'),
              backgroundColor: AppColors.success,
            ),
          );
          ref.invalidate(deductionPlanDetailProvider(id));
          ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true);
        }
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('操作失败: $e'),
              backgroundColor: AppColors.danger,
            ),
          );
        }
      }
    }
  }

  Future<void> _handleCancel(BuildContext context, WidgetRef ref, DeductionPlanDetail detail) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('取消代扣'),
        content: const Text('确定要取消此代扣计划吗？取消后不可恢复。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('返回'),
          ),
          ElevatedButton(
            onPressed: () => Navigator.of(context).pop(true),
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
              foregroundColor: Colors.white,
            ),
            child: const Text('确定取消'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final service = ref.read(deductionServiceProvider);
      try {
        await service.cancelPlan(detail.id);
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('已取消'),
              backgroundColor: AppColors.success,
            ),
          );
          ref.invalidate(deductionPlanDetailProvider(id));
          ref.read(deductionPlansProvider.notifier).loadPlans(refresh: true);
          Navigator.of(context).pop();
        }
      } catch (e) {
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('操作失败: $e'),
              backgroundColor: AppColors.danger,
            ),
          );
        }
      }
    }
  }
}
