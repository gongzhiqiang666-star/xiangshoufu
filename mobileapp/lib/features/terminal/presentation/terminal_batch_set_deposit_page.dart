import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import 'providers/terminal_provider.dart';

/// 批量设置押金页面
class TerminalBatchSetDepositPage extends ConsumerStatefulWidget {
  final List<String> selectedSNs;

  const TerminalBatchSetDepositPage({
    super.key,
    required this.selectedSNs,
  });

  @override
  ConsumerState<TerminalBatchSetDepositPage> createState() =>
      _TerminalBatchSetDepositPageState();
}

class _TerminalBatchSetDepositPageState
    extends ConsumerState<TerminalBatchSetDepositPage> {
  // 押金选项（分）
  static const List<Map<String, dynamic>> _depositOptions = [
    {'label': '无押金', 'value': 0},
    {'label': '¥99', 'value': 9900},
    {'label': '¥199', 'value': 19900},
    {'label': '¥299', 'value': 29900},
  ];

  int? _selectedDeposit;

  @override
  Widget build(BuildContext context) {
    final batchSetState = ref.watch(batchSetProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('批量设置押金'),
        centerTitle: true,
      ),
      body: Column(
        children: [
          Expanded(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(AppSpacing.md),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // 终端数量提示
                  _buildTerminalCount(),
                  const SizedBox(height: AppSpacing.lg),

                  // 押金选择
                  _buildSectionTitle('选择押金'),
                  const SizedBox(height: AppSpacing.sm),
                  _buildDepositSelector(),
                  const SizedBox(height: AppSpacing.lg),

                  // 说明
                  _buildNote(),
                ],
              ),
            ),
          ),

          // 底部按钮
          _buildBottomButton(batchSetState),
        ],
      ),
    );
  }

  Widget _buildTerminalCount() {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.primary.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          const Icon(Icons.devices, color: AppColors.primary),
          const SizedBox(width: AppSpacing.sm),
          Text(
            '已选择 ${widget.selectedSNs.length} 台终端',
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w500,
              color: AppColors.primary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Text(
      title,
      style: const TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
      ),
    );
  }

  Widget _buildDepositSelector() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: _depositOptions.map((option) {
          final isSelected = _selectedDeposit == option['value'];
          return GestureDetector(
            onTap: () {
              setState(() {
                _selectedDeposit = option['value'] as int;
              });
            },
            child: Container(
              width: double.infinity,
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.md,
                vertical: AppSpacing.md,
              ),
              margin: EdgeInsets.only(
                bottom: option != _depositOptions.last ? AppSpacing.sm : 0,
              ),
              decoration: BoxDecoration(
                color: isSelected
                    ? AppColors.primary.withValues(alpha: 0.1)
                    : AppColors.background,
                borderRadius: BorderRadius.circular(8),
                border: Border.all(
                  color: isSelected ? AppColors.primary : AppColors.border,
                  width: isSelected ? 2 : 1,
                ),
              ),
              child: Row(
                children: [
                  Icon(
                    isSelected
                        ? Icons.radio_button_checked
                        : Icons.radio_button_off,
                    color: isSelected
                        ? AppColors.primary
                        : AppColors.textSecondary,
                    size: 22,
                  ),
                  const SizedBox(width: AppSpacing.sm),
                  Text(
                    option['label'] as String,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                      color: isSelected
                          ? AppColors.primary
                          : AppColors.textPrimary,
                    ),
                  ),
                  const Spacer(),
                  if (option['value'] == 0)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 2,
                      ),
                      decoration: BoxDecoration(
                        color: AppColors.success.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: const Text(
                        '推荐',
                        style: TextStyle(
                          fontSize: 12,
                          color: AppColors.success,
                        ),
                      ),
                    ),
                ],
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildNote() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: AppColors.warning.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(Icons.info_outline, color: AppColors.warning, size: 20),
          SizedBox(width: AppSpacing.sm),
          Expanded(
            child: Text(
              '押金将在商户激活终端时从交易中扣除。无押金模式下，商户无需支付押金即可使用终端。',
              style: TextStyle(
                fontSize: 13,
                color: AppColors.textSecondary,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomButton(BatchSetState state) {
    return Container(
      padding: EdgeInsets.only(
        left: AppSpacing.md,
        right: AppSpacing.md,
        top: AppSpacing.md,
        bottom: MediaQuery.of(context).padding.bottom + AppSpacing.md,
      ),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SizedBox(
        width: double.infinity,
        height: 48,
        child: ElevatedButton(
          onPressed: _selectedDeposit == null || state.isSubmitting
              ? null
              : _handleSubmit,
          child: state.isSubmitting
              ? const SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(
                    strokeWidth: 2,
                    valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                  ),
                )
              : const Text('确认设置'),
        ),
      ),
    );
  }

  Future<void> _handleSubmit() async {
    if (_selectedDeposit == null) return;

    final success = await ref.read(batchSetProvider.notifier).batchSetDeposit(
      terminalSns: widget.selectedSNs,
      depositAmount: _selectedDeposit!,
    );

    if (mounted) {
      final state = ref.read(batchSetProvider);
      if (success) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('设置成功，共${state.successCount}台终端'),
            backgroundColor: AppColors.success,
          ),
        );
        // 清空选中状态并返回
        ref.read(selectedTerminalsProvider.notifier).state = [];
        context.pop();
      } else if (state.failedCount > 0) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
                '部分设置失败: 成功${state.successCount}台, 失败${state.failedCount}台'),
            backgroundColor: Colors.orange,
          ),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('设置失败: ${state.error ?? "未知错误"}'),
            backgroundColor: AppColors.danger,
          ),
        );
      }
    }
  }
}
