import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import 'providers/goods_deduction_provider.dart';
import 'widgets/agreement_dialog.dart';

/// 货款代扣详情页
class GoodsDeductionDetailPage extends ConsumerWidget {
  final int id;
  final bool isSent;

  const GoodsDeductionDetailPage({
    super.key,
    required this.id,
    required this.isSent,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(goodsDeductionDetailProvider(id));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('货款代扣详情'),
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
                onPressed: () => ref.invalidate(goodsDeductionDetailProvider(id)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
        data: (detail) => _buildContent(context, ref, detail),
      ),
    );
  }

  Widget _buildContent(BuildContext context, WidgetRef ref, dynamic detail) {
    return SingleChildScrollView(
      child: Column(
        children: [
          _buildStatusCard(detail),
          _buildAmountCard(detail),
          _buildInfoCard(detail),
          if (detail.terminals.isNotEmpty) _buildTerminalsCard(detail),
          _buildDetailsCard(detail),
          if (!isSent && detail.status == 1) _buildActionButtons(context, ref, detail),
          const SizedBox(height: AppSpacing.xl),
        ],
      ),
    );
  }

  Widget _buildStatusCard(dynamic detail) {
    Color statusColor;
    IconData statusIcon;

    switch (detail.status) {
      case 1:
        statusColor = AppColors.warning;
        statusIcon = Icons.schedule;
        break;
      case 2:
        statusColor = AppColors.primary;
        statusIcon = Icons.sync;
        break;
      case 3:
        statusColor = AppColors.success;
        statusIcon = Icons.check_circle;
        break;
      case 4:
        statusColor = AppColors.danger;
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
                Text(
                  detail.statusName,
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w600,
                    color: statusColor,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  _getStatusDescription(detail.status),
                  style: const TextStyle(
                    fontSize: 14,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),
          _buildProgressCircle(detail.progress),
        ],
      ),
    );
  }

  String _getStatusDescription(int status) {
    switch (status) {
      case 1:
        return '等待对方接收确认';
      case 2:
        return '正在自动扣款中';
      case 3:
        return '代扣已完成';
      case 4:
        return '对方已拒绝此代扣';
      default:
        return '';
    }
  }

  Widget _buildProgressCircle(double progress) {
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
              value: progress / 100,
              strokeWidth: 6,
              backgroundColor: AppColors.border,
              valueColor: AlwaysStoppedAnimation<Color>(
                progress >= 100 ? AppColors.success : AppColors.primary,
              ),
            ),
          ),
          Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(
                '${progress.toStringAsFixed(1)}%',
                style: const TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              const Text(
                '进度',
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

  Widget _buildAmountCard(dynamic detail) {
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

  Widget _buildInfoCard(dynamic detail) {
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
          _buildInfoRow('代扣编号', detail.deductionNo),
          _buildInfoRow('发起方', detail.fromAgentName),
          _buildInfoRow('接收方', detail.toAgentName),
          _buildInfoRow('终端数量', '${detail.terminalCount} 台'),
          _buildInfoRow('终端单价', '¥${FormatUtils.formatYuan(detail.unitPriceYuan)}'),
          _buildInfoRow('扣款来源', detail.sourceName),
          _buildInfoRow('创建时间', detail.createdAt.substring(0, 16)),
          if (detail.acceptedAt != null)
            _buildInfoRow('接收时间', detail.acceptedAt.substring(0, 16)),
          if (detail.completedAt != null)
            _buildInfoRow('完成时间', detail.completedAt.substring(0, 16)),
          if (detail.remark != null && detail.remark.isNotEmpty)
            _buildInfoRow('备注', detail.remark),
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

  Widget _buildTerminalsCard(dynamic detail) {
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
                '关联终端',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              Text(
                '共 ${detail.terminals.length} 台',
                style: const TextStyle(
                  fontSize: 14,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.sm),
          const Divider(),
          ...detail.terminals.take(5).map<Widget>((terminal) => Padding(
                padding: const EdgeInsets.symmetric(vertical: AppSpacing.xs),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      terminal.terminalSn,
                      style: const TextStyle(
                        fontSize: 14,
                        color: AppColors.textPrimary,
                      ),
                    ),
                    Text(
                      '¥${FormatUtils.formatYuan(terminal.unitPriceYuan)}',
                      style: const TextStyle(
                        fontSize: 14,
                        color: AppColors.textSecondary,
                      ),
                    ),
                  ],
                ),
              )),
          if (detail.terminals.length > 5)
            Padding(
              padding: const EdgeInsets.only(top: AppSpacing.sm),
              child: Center(
                child: Text(
                  '... 还有 ${detail.terminals.length - 5} 台',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textTertiary,
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildDetailsCard(dynamic detail) {
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
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                '扣款明细',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                  color: AppColors.textPrimary,
                ),
              ),
              Text(
                '共 ${detail.details.length} 条',
                style: const TextStyle(
                  fontSize: 14,
                  color: AppColors.textSecondary,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.sm),
          if (detail.details.isEmpty)
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
            ...detail.details.map<Widget>((record) => _buildDetailItem(record)),
        ],
      ),
    );
  }

  Widget _buildDetailItem(dynamic record) {
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
            width: 8,
            height: 8,
            decoration: const BoxDecoration(
              color: AppColors.success,
              shape: BoxShape.circle,
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
                      '¥${FormatUtils.formatYuan(record.amountYuan)}',
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w600,
                        color: AppColors.success,
                      ),
                    ),
                    Text(
                      record.walletTypeName,
                      style: const TextStyle(
                        fontSize: 12,
                        color: AppColors.primary,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  record.createdAt.substring(0, 16),
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textTertiary,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildActionButtons(BuildContext context, WidgetRef ref, dynamic detail) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Row(
        children: [
          Expanded(
            child: OutlinedButton(
              onPressed: () => _handleReject(context, ref, detail),
              style: OutlinedButton.styleFrom(
                foregroundColor: AppColors.danger,
                side: const BorderSide(color: AppColors.danger),
                padding: const EdgeInsets.symmetric(vertical: 14),
              ),
              child: const Text('拒绝'),
            ),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: ElevatedButton(
              onPressed: () => _handleAccept(context, ref, detail),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.success,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 14),
              ),
              child: const Text('接收'),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _handleAccept(BuildContext context, WidgetRef ref, dynamic detail) async {
    final agreed = await showAgreementDialog(
      context: context,
      title: '代扣服务协议',
      content: getDefaultAgreementContent(
        fromAgentName: detail.fromAgentName,
        toAgentName: detail.toAgentName,
        totalAmount: detail.totalAmountYuan,
        terminalCount: detail.terminalCount,
      ),
    );

    if (agreed == true) {
      final service = ref.read(goodsDeductionServiceProvider);
      try {
        await service.acceptDeduction(detail.id);
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('接收成功，代扣已开始'),
              backgroundColor: AppColors.success,
            ),
          );
          ref.invalidate(goodsDeductionDetailProvider(id));
          ref.read(receivedDeductionsProvider.notifier).loadDeductions(refresh: true);
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

  Future<void> _handleReject(BuildContext context, WidgetRef ref, dynamic detail) async {
    final reasonController = TextEditingController();

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('拒绝货款代扣'),
        content: TextField(
          controller: reasonController,
          maxLines: 3,
          decoration: const InputDecoration(
            hintText: '请输入拒绝原因',
            border: OutlineInputBorder(),
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () {
              if (reasonController.text.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('请输入拒绝原因')),
                );
                return;
              }
              Navigator.of(context).pop(true);
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppColors.danger,
              foregroundColor: Colors.white,
            ),
            child: const Text('确认拒绝'),
          ),
        ],
      ),
    );

    if (confirmed == true && reasonController.text.isNotEmpty) {
      final service = ref.read(goodsDeductionServiceProvider);
      try {
        await service.rejectDeduction(detail.id, reasonController.text);
        if (context.mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('已拒绝'),
              backgroundColor: AppColors.success,
            ),
          );
          ref.invalidate(goodsDeductionDetailProvider(id));
          ref.read(receivedDeductionsProvider.notifier).loadDeductions(refresh: true);
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

    reasonController.dispose();
  }
}
