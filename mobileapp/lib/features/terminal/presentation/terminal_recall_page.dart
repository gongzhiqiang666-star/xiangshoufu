import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import 'providers/terminal_provider.dart';

/// 终端回拨页面
class TerminalRecallPage extends ConsumerStatefulWidget {
  final List<String> selectedSNs;

  const TerminalRecallPage({
    super.key,
    required this.selectedSNs,
  });

  @override
  ConsumerState<TerminalRecallPage> createState() => _TerminalRecallPageState();
}

class _TerminalRecallPageState extends ConsumerState<TerminalRecallPage> {
  String? _selectedAgentId;
  // 模拟的上级代理商列表（实际应从API获取当前用户的直属上级）
  final List<Map<String, dynamic>> _parentAgents = [
    {'id': '1001', 'name': '总代理', 'phone': '13800000000'},
  ];

  final _remarkController = TextEditingController();

  @override
  void initState() {
    super.initState();
    // 默认选中第一个上级
    if (_parentAgents.isNotEmpty) {
      _selectedAgentId = _parentAgents.first['id'];
    }
  }

  @override
  void dispose() {
    _remarkController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final recallState = ref.watch(terminalRecallProvider);

    // 监听状态变化
    ref.listen(terminalRecallProvider, (previous, next) {
      if (next.error != null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('回拨失败: ${next.error}')),
        );
      }

      if (!next.isSubmitting && previous?.isSubmitting == true && next.error == null) {
        if (next.failedCount == 0) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('成功回拨 ${next.successCount} 台终端')),
          );
          context.pop();
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('回拨完成: 成功${next.successCount}台, 失败${next.failedCount}台')),
          );
          context.pop();
        }
      }
    });

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('终端回拨')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildSelectedTerminals(),
            const SizedBox(height: AppSpacing.md),
            _buildAgentSelector(),
            const SizedBox(height: AppSpacing.md),
            _buildRemarkInput(),
            const SizedBox(height: AppSpacing.md),
            _buildWarningNotice(),
          ],
        ),
      ),
      bottomNavigationBar: _buildBottomBar(recallState.isSubmitting),
    );
  }

  Widget _buildSelectedTerminals() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '已选终端: ${widget.selectedSNs.length}台',
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'SN: ${widget.selectedSNs.join(", ")}',
            style: const TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAgentSelector() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '回拨给:',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          ..._parentAgents.map((agent) => _buildAgentItem(agent)),
        ],
      ),
    );
  }

  Widget _buildAgentItem(Map<String, dynamic> agent) {
    final isSelected = _selectedAgentId == agent['id'];

    return GestureDetector(
      onTap: () {
        setState(() {
          _selectedAgentId = agent['id'];
        });
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 8),
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: isSelected
              ? AppColors.primary.withValues(alpha: 0.05)
              : AppColors.background,
          border: Border.all(
            color: isSelected ? AppColors.primary : Colors.transparent,
          ),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Icon(
              isSelected ? Icons.radio_button_checked : Icons.radio_button_off,
              color: isSelected ? AppColors.primary : AppColors.textTertiary,
              size: 20,
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    '${agent['name']}',
                    style: const TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w500,
                      color: AppColors.textPrimary,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    '手机: ${agent['phone']}',
                    style: const TextStyle(
                      fontSize: 13,
                      color: AppColors.textSecondary,
                    ),
                  ),
                ],
              ),
            ),
            if (isSelected)
              const Text(
                '直属上级',
                style: TextStyle(
                  fontSize: 12,
                  color: AppColors.primary,
                ),
              ),
          ],
        ),
      ),
    );
  }

  Widget _buildRemarkInput() {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '备注',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          TextField(
            controller: _remarkController,
            maxLines: 3,
            decoration: const InputDecoration(
              hintText: '请输入回拨备注（选填）',
              border: OutlineInputBorder(),
              contentPadding: EdgeInsets.all(12),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildWarningNotice() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.warning.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: const Row(
        children: [
          Icon(Icons.warning_amber, color: AppColors.warning, size: 20),
          SizedBox(width: 8),
          Expanded(
            child: Text(
              '注意：只能回拨未激活的终端，APP仅支持回拨给直属上级',
              style: TextStyle(
                fontSize: 13,
                color: AppColors.warning,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomBar(bool isSubmitting) {
    return Container(
      padding: EdgeInsets.only(
        left: AppSpacing.md,
        right: AppSpacing.md,
        top: 12,
        bottom: MediaQuery.of(context).padding.bottom + 12,
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
      child: ElevatedButton(
        onPressed: isSubmitting || _selectedAgentId == null ? null : _handleRecall,
        child: isSubmitting
            ? const SizedBox(
                width: 20,
                height: 20,
                child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white),
              )
            : const Text('确认回拨'),
      ),
    );
  }

  void _handleRecall() {
    if (_selectedAgentId == null) return;

    final toAgentId = int.tryParse(_selectedAgentId!) ?? 0;

    ref.read(terminalRecallProvider.notifier).batchRecall(
      toAgentId: toAgentId,
      terminalSns: widget.selectedSNs,
      remark: _remarkController.text.isNotEmpty ? _remarkController.text : null,
    );
  }
}
