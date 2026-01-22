import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import 'providers/terminal_provider.dart';

/// 批量设置费率页面
class TerminalBatchSetRatePage extends ConsumerStatefulWidget {
  final List<String> selectedSNs;

  const TerminalBatchSetRatePage({
    super.key,
    required this.selectedSNs,
  });

  @override
  ConsumerState<TerminalBatchSetRatePage> createState() =>
      _TerminalBatchSetRatePageState();
}

class _TerminalBatchSetRatePageState
    extends ConsumerState<TerminalBatchSetRatePage> {
  // 费率选项（万分比）
  static const List<Map<String, dynamic>> _rateOptions = [
    {'label': '0.53%', 'value': 53},
    {'label': '0.54%', 'value': 54},
    {'label': '0.55%', 'value': 55},
    {'label': '0.56%', 'value': 56},
    {'label': '0.57%', 'value': 57},
    {'label': '0.58%', 'value': 58},
    {'label': '0.59%', 'value': 59},
    {'label': '0.60%', 'value': 60},
  ];

  int? _selectedRate;

  @override
  Widget build(BuildContext context) {
    final batchSetState = ref.watch(batchSetProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('批量设置费率'),
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

                  // 费率选择
                  _buildSectionTitle('选择费率'),
                  const SizedBox(height: AppSpacing.sm),
                  _buildRateSelector(),
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
        color: AppColors.primary.withOpacity(0.1),
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

  Widget _buildRateSelector() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Wrap(
        spacing: AppSpacing.sm,
        runSpacing: AppSpacing.sm,
        children: _rateOptions.map((option) {
          final isSelected = _selectedRate == option['value'];
          return GestureDetector(
            onTap: () {
              setState(() {
                _selectedRate = option['value'] as int;
              });
            },
            child: Container(
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.md,
                vertical: AppSpacing.sm,
              ),
              decoration: BoxDecoration(
                color: isSelected
                    ? AppColors.primary
                    : AppColors.background,
                borderRadius: BorderRadius.circular(8),
                border: Border.all(
                  color: isSelected
                      ? AppColors.primary
                      : AppColors.border,
                ),
              ),
              child: Text(
                option['label'] as String,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: isSelected
                      ? Colors.white
                      : AppColors.textPrimary,
                ),
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
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(Icons.info_outline, color: AppColors.warning, size: 20),
          SizedBox(width: AppSpacing.sm),
          Expanded(
            child: Text(
              '设置的费率将应用于所有选中的终端。费率设置后，商户刷卡时将按此费率收取手续费。',
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
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, -2),
          ),
        ],
      ),
      child: SizedBox(
        width: double.infinity,
        height: 48,
        child: ElevatedButton(
          onPressed: _selectedRate == null || state.isSubmitting
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
    if (_selectedRate == null) return;

    final success = await ref.read(batchSetProvider.notifier).batchSetRate(
      terminalSns: widget.selectedSNs,
      creditRate: _selectedRate!,
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
