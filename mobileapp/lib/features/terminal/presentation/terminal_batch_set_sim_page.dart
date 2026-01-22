import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import 'providers/terminal_provider.dart';

/// 批量设置流量费页面
class TerminalBatchSetSimPage extends ConsumerStatefulWidget {
  final List<String> selectedSNs;

  const TerminalBatchSetSimPage({
    super.key,
    required this.selectedSNs,
  });

  @override
  ConsumerState<TerminalBatchSetSimPage> createState() =>
      _TerminalBatchSetSimPageState();
}

class _TerminalBatchSetSimPageState
    extends ConsumerState<TerminalBatchSetSimPage> {
  // 首次流量费选项（分）
  static const List<Map<String, dynamic>> _firstSimFeeOptions = [
    {'label': '¥48', 'value': 4800},
    {'label': '¥60', 'value': 6000},
    {'label': '¥69', 'value': 6900},
    {'label': '¥79', 'value': 7900},
    {'label': '¥89', 'value': 8900},
    {'label': '¥99', 'value': 9900},
  ];

  // 非首次流量费选项（分）
  static const List<Map<String, dynamic>> _nonFirstSimFeeOptions = [
    {'label': '¥48', 'value': 4800},
    {'label': '¥60', 'value': 6000},
    {'label': '¥69', 'value': 6900},
    {'label': '¥79', 'value': 7900},
    {'label': '¥89', 'value': 8900},
    {'label': '¥99', 'value': 9900},
  ];

  // 间隔天数选项
  static const List<Map<String, dynamic>> _intervalDaysOptions = [
    {'label': '180天', 'value': 180},
    {'label': '210天', 'value': 210},
    {'label': '240天', 'value': 240},
    {'label': '270天', 'value': 270},
    {'label': '300天', 'value': 300},
    {'label': '330天', 'value': 330},
    {'label': '360天', 'value': 360},
  ];

  int? _selectedFirstSimFee;
  int? _selectedNonFirstSimFee;
  int? _selectedIntervalDays;

  @override
  Widget build(BuildContext context) {
    final batchSetState = ref.watch(batchSetProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('批量设置流量费'),
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

                  // 首次流量费选择
                  _buildSectionTitle('首次流量费'),
                  const SizedBox(height: AppSpacing.sm),
                  _buildOptionSelector(
                    _firstSimFeeOptions,
                    _selectedFirstSimFee,
                    (value) => setState(() => _selectedFirstSimFee = value),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // 非首次流量费选择
                  _buildSectionTitle('非首次流量费'),
                  const SizedBox(height: AppSpacing.sm),
                  _buildOptionSelector(
                    _nonFirstSimFeeOptions,
                    _selectedNonFirstSimFee,
                    (value) => setState(() => _selectedNonFirstSimFee = value),
                  ),
                  const SizedBox(height: AppSpacing.lg),

                  // 间隔天数选择
                  _buildSectionTitle('收费间隔'),
                  const SizedBox(height: AppSpacing.sm),
                  _buildOptionSelector(
                    _intervalDaysOptions,
                    _selectedIntervalDays,
                    (value) => setState(() => _selectedIntervalDays = value),
                  ),
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

  Widget _buildOptionSelector(
    List<Map<String, dynamic>> options,
    int? selectedValue,
    ValueChanged<int> onSelect,
  ) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Wrap(
        spacing: AppSpacing.sm,
        runSpacing: AppSpacing.sm,
        children: options.map((option) {
          final isSelected = selectedValue == option['value'];
          return GestureDetector(
            onTap: () => onSelect(option['value'] as int),
            child: Container(
              padding: const EdgeInsets.symmetric(
                horizontal: AppSpacing.md,
                vertical: AppSpacing.sm,
              ),
              decoration: BoxDecoration(
                color: isSelected ? AppColors.primary : AppColors.background,
                borderRadius: BorderRadius.circular(8),
                border: Border.all(
                  color: isSelected ? AppColors.primary : AppColors.border,
                ),
              ),
              child: Text(
                option['label'] as String,
                style: TextStyle(
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                  color: isSelected ? Colors.white : AppColors.textPrimary,
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
              '流量费将按设置的间隔天数自动收取。首次流量费在终端激活时收取，非首次流量费在之后的每个周期收取。',
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
    final isValid = _selectedFirstSimFee != null &&
        _selectedNonFirstSimFee != null &&
        _selectedIntervalDays != null;

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
          onPressed: !isValid || state.isSubmitting ? null : _handleSubmit,
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
    if (_selectedFirstSimFee == null ||
        _selectedNonFirstSimFee == null ||
        _selectedIntervalDays == null) {
      return;
    }

    final success = await ref.read(batchSetProvider.notifier).batchSetSimFee(
      terminalSns: widget.selectedSNs,
      firstSimFee: _selectedFirstSimFee!,
      nonFirstSimFee: _selectedNonFirstSimFee!,
      simFeeIntervalDays: _selectedIntervalDays!,
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
