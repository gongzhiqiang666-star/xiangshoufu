import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';

/// 交易记录页面
class TransactionPage extends ConsumerStatefulWidget {
  const TransactionPage({super.key});

  @override
  ConsumerState<TransactionPage> createState() => _TransactionPageState();
}

class _TransactionPageState extends ConsumerState<TransactionPage> {
  final List<_TransactionItem> _transactions = [
    _TransactionItem(
      merchantName: '张三便利店',
      payType: '微信',
      amount: 128.50,
      time: '2024-01-15 14:30',
      status: '成功',
    ),
    _TransactionItem(
      merchantName: '李四超市',
      payType: '支付宝',
      amount: 56.00,
      time: '2024-01-15 13:20',
      status: '成功',
    ),
    _TransactionItem(
      merchantName: '王五餐饮',
      payType: '银行卡',
      amount: 320.00,
      time: '2024-01-15 12:15',
      status: '成功',
    ),
    _TransactionItem(
      merchantName: '赵六水果店',
      payType: '微信',
      amount: 45.80,
      time: '2024-01-15 11:05',
      status: '成功',
    ),
    _TransactionItem(
      merchantName: '钱七药店',
      payType: '支付宝',
      amount: 89.00,
      time: '2024-01-15 10:30',
      status: '成功',
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('交易记录'),
      ),
      body: _transactions.isEmpty
          ? _buildEmptyState()
          : ListView.separated(
              padding: const EdgeInsets.all(AppSpacing.md),
              itemCount: _transactions.length,
              separatorBuilder: (context, index) => const SizedBox(height: AppSpacing.sm),
              itemBuilder: (context, index) {
                return _buildTransactionCard(_transactions[index]);
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.receipt_long_outlined,
            size: 64,
            color: AppColors.textTertiary,
          ),
          const SizedBox(height: 16),
          Text(
            '暂无交易记录',
            style: TextStyle(
              fontSize: 16,
              color: AppColors.textSecondary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTransactionCard(_TransactionItem tx) {
    IconData icon;
    Color iconColor;

    switch (tx.payType) {
      case '微信':
        icon = Icons.wechat;
        iconColor = AppColors.wechatPay;
        break;
      case '支付宝':
        icon = Icons.account_balance_wallet;
        iconColor = AppColors.alipay;
        break;
      default:
        icon = Icons.credit_card;
        iconColor = AppColors.primary;
    }

    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Container(
            width: 44,
            height: 44,
            decoration: BoxDecoration(
              color: iconColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: iconColor, size: 22),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  tx.merchantName,
                  style: const TextStyle(
                    fontSize: 15,
                    fontWeight: FontWeight.w500,
                  ),
                ),
                const SizedBox(height: 4),
                Text(
                  '${tx.payType} · ${tx.time}',
                  style: TextStyle(
                    fontSize: 12,
                    color: AppColors.textTertiary,
                  ),
                ),
              ],
            ),
          ),
          Column(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              Text(
                FormatUtils.formatYuan(tx.amount),
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 4),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                decoration: BoxDecoration(
                  color: AppColors.success.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  tx.status,
                  style: TextStyle(
                    fontSize: 10,
                    color: AppColors.success,
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

class _TransactionItem {
  final String merchantName;
  final String payType;
  final double amount;
  final String time;
  final String status;

  _TransactionItem({
    required this.merchantName,
    required this.payType,
    required this.amount,
    required this.time,
    required this.status,
  });
}
