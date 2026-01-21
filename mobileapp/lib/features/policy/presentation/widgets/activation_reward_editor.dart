import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../data/models/policy_model.dart';

/// 激活奖励编辑器
class ActivationRewardEditor extends StatefulWidget {
  final List<ActivationRewardItem> initialItems;
  final List<ActivationRewardItem> maxItems;
  final ValueChanged<List<ActivationRewardItem>> onChanged;

  const ActivationRewardEditor({
    super.key,
    required this.initialItems,
    required this.maxItems,
    required this.onChanged,
  });

  @override
  State<ActivationRewardEditor> createState() => _ActivationRewardEditorState();
}

class _ActivationRewardEditorState extends State<ActivationRewardEditor> {
  late List<_RewardEditItem> _editItems;

  @override
  void initState() {
    super.initState();
    _editItems = widget.initialItems.map((item) => _RewardEditItem(item)).toList();
    if (_editItems.isEmpty && widget.maxItems.isNotEmpty) {
      // 使用上级的奖励配置作为模板
      _editItems = widget.maxItems.map((item) => _RewardEditItem(
        ActivationRewardItem(
          rewardName: item.rewardName,
          minRegisterDays: item.minRegisterDays,
          maxRegisterDays: item.maxRegisterDays,
          targetAmount: item.targetAmount,
          rewardAmount: 0,
        ),
      )).toList();
    }
  }

  @override
  void dispose() {
    for (var item in _editItems) {
      item.dispose();
    }
    super.dispose();
  }

  void _notifyChange() {
    widget.onChanged(_editItems.map((e) => e.toItem()).toList());
  }

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildInfoCard(),
          const SizedBox(height: 16),
          if (_editItems.isEmpty)
            const Center(
              child: Padding(
                padding: EdgeInsets.all(32),
                child: Text('暂无激活奖励配置'),
              ),
            )
          else
            ..._editItems.asMap().entries.map((entry) {
              final index = entry.key;
              final maxItem = index < widget.maxItems.length ? widget.maxItems[index] : null;
              return _buildRewardCard(entry.value, maxItem);
            }),
        ],
      ),
    );
  }

  Widget _buildInfoCard() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.purple.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          Icon(Icons.info_outline, color: Colors.purple.shade400, size: 20),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '商户入网后达到交易量目标，给下级代理商发放奖励。奖励金额不能超过您的配置。',
              style: TextStyle(
                fontSize: 13,
                color: Colors.purple.shade400,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRewardCard(_RewardEditItem editItem, ActivationRewardItem? maxItem) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              editItem.item.rewardName,
              style: const TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              '入网 ${editItem.item.minRegisterDays}-${editItem.item.maxRegisterDays} 天内，交易满 ${editItem.item.targetAmountWan.toStringAsFixed(2)} 万元',
              style: const TextStyle(
                fontSize: 13,
                color: AppColors.textSecondary,
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: editItem.rewardController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: InputDecoration(
                labelText: '奖励金额',
                hintText: maxItem != null
                    ? '最高 ¥${maxItem.rewardAmountYuan.toStringAsFixed(2)}'
                    : '请输入奖励金额',
                suffixText: '元',
                border: const OutlineInputBorder(),
              ),
              onChanged: (_) => _notifyChange(),
            ),
          ],
        ),
      ),
    );
  }
}

class _RewardEditItem {
  final ActivationRewardItem item;
  final TextEditingController rewardController;

  _RewardEditItem(this.item)
      : rewardController = TextEditingController(
          text: item.rewardAmountYuan.toStringAsFixed(2),
        );

  void dispose() {
    rewardController.dispose();
  }

  ActivationRewardItem toItem() {
    return ActivationRewardItem(
      rewardName: item.rewardName,
      minRegisterDays: item.minRegisterDays,
      maxRegisterDays: item.maxRegisterDays,
      targetAmount: item.targetAmount,
      rewardAmount: ((double.tryParse(rewardController.text) ?? 0) * 100).round(),
    );
  }
}
