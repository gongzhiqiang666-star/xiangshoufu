import 'package:flutter/material.dart';
import '../../../../core/theme/app_colors.dart';
import '../../../../core/theme/app_spacing.dart';
import '../../data/models/policy_model.dart';

/// 押金返现编辑器
class DepositCashbackEditor extends StatefulWidget {
  final List<DepositCashbackItem> initialItems;
  final List<DepositCashbackItem> maxItems;
  final ValueChanged<List<DepositCashbackItem>> onChanged;

  const DepositCashbackEditor({
    super.key,
    required this.initialItems,
    required this.maxItems,
    required this.onChanged,
  });

  @override
  State<DepositCashbackEditor> createState() => _DepositCashbackEditorState();
}

class _DepositCashbackEditorState extends State<DepositCashbackEditor> {
  late List<DepositCashbackItem> _items;
  final Map<int, TextEditingController> _controllers = {};

  @override
  void initState() {
    super.initState();
    _items = List.from(widget.initialItems);
    // 确保包含常见押金档位
    _ensureDefaultItems();
    for (var item in _items) {
      _controllers[item.depositAmount] = TextEditingController(
        text: item.cashbackAmountYuan.toStringAsFixed(2),
      );
    }
  }

  void _ensureDefaultItems() {
    final defaultDeposits = [0, 9900, 19900, 29900]; // 分
    for (var deposit in defaultDeposits) {
      if (!_items.any((item) => item.depositAmount == deposit)) {
        _items.add(DepositCashbackItem(
          depositAmount: deposit,
          cashbackAmount: 0,
        ));
      }
    }
    _items.sort((a, b) => a.depositAmount.compareTo(b.depositAmount));
  }

  @override
  void dispose() {
    for (var controller in _controllers.values) {
      controller.dispose();
    }
    super.dispose();
  }

  void _updateItem(int depositAmount, String cashbackStr) {
    final cashback = (double.tryParse(cashbackStr) ?? 0) * 100;
    final index = _items.indexWhere((item) => item.depositAmount == depositAmount);
    if (index >= 0) {
      _items[index] = DepositCashbackItem(
        depositAmount: depositAmount,
        cashbackAmount: cashback.round(),
      );
      widget.onChanged(_items);
    }
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
          ..._items.map((item) => _buildCashbackField(item)),
        ],
      ),
    );
  }

  Widget _buildInfoCard() {
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: AppColors.success.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Row(
        children: [
          const Icon(Icons.info_outline, color: AppColors.success, size: 20),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              '商户押金收取后，按此配置返现给下级代理商。返现金额不能超过您的配置。',
              style: TextStyle(
                fontSize: 13,
                color: AppColors.success.withOpacity(0.8),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCashbackField(DepositCashbackItem item) {
    final maxItem = widget.maxItems.firstWhere(
      (m) => m.depositAmount == item.depositAmount,
      orElse: () => DepositCashbackItem(depositAmount: item.depositAmount, cashbackAmount: 0),
    );

    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextField(
        controller: _controllers[item.depositAmount],
        keyboardType: const TextInputType.numberWithOptions(decimal: true),
        decoration: InputDecoration(
          labelText: '押金 ¥${item.depositAmountYuan.toInt()} 返现',
          hintText: '最高 ¥${maxItem.cashbackAmountYuan.toStringAsFixed(2)}',
          suffixText: '元',
          border: const OutlineInputBorder(),
          filled: true,
          fillColor: Colors.white,
        ),
        onChanged: (value) => _updateItem(item.depositAmount, value),
      ),
    );
  }
}
