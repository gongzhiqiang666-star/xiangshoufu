import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../data/models/wallet_model.dart';
import 'providers/wallet_provider.dart';
import '../data/services/wallet_service.dart';

/// 沉淀钱包页面
class SettlementWalletPage extends ConsumerStatefulWidget {
  const SettlementWalletPage({super.key});

  @override
  ConsumerState<SettlementWalletPage> createState() => _SettlementWalletPageState();
}

class _SettlementWalletPageState extends ConsumerState<SettlementWalletPage> {
  final _amountController = TextEditingController();
  final _remarkController = TextEditingController();

  @override
  void dispose() {
    _amountController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final summaryAsync = ref.watch(settlementWalletSummaryProvider);
    final subordinatesAsync = ref.watch(subordinateBalancesProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('沉淀钱包'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              ref.invalidate(settlementWalletSummaryProvider);
              ref.invalidate(subordinateBalancesProvider);
            },
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(settlementWalletSummaryProvider);
          ref.invalidate(subordinateBalancesProvider);
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          child: Column(
            children: [
              // 汇总卡片
              summaryAsync.when(
                data: (summary) => _buildSummaryCard(summary),
                loading: () => const Center(
                  child: Padding(
                    padding: EdgeInsets.all(32),
                    child: CircularProgressIndicator(),
                  ),
                ),
                error: (e, _) => _buildErrorCard('加载失败: $e'),
              ),

              // 操作按钮
              summaryAsync.when(
                data: (summary) => _buildActionButtons(summary),
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
              ),

              // 下级余额明细
              _buildSubordinateSection(subordinatesAsync),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSummaryCard(SettlementWalletSummaryModel summary) {
    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF667eea), Color(0xFF764ba2)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF667eea).withOpacity(0.3),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        children: [
          // 沉淀比例
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.percent, color: Colors.white70, size: 16),
              const SizedBox(width: 4),
              Text(
                '沉淀比例: ${summary.settlementRatio}%',
                style: const TextStyle(
                  fontSize: 14,
                  color: Colors.white70,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 统计项
          Row(
            children: [
              Expanded(
                child: _buildStatItem(
                  '下级未提现',
                  FormatUtils.formatYuan(summary.subordinateTotalBalanceYuan),
                  Colors.white,
                ),
              ),
              Expanded(
                child: _buildStatItem(
                  '可用额度',
                  FormatUtils.formatYuan(summary.availableAmountYuan),
                  Colors.greenAccent,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),
          Row(
            children: [
              Expanded(
                child: _buildStatItem(
                  '已使用',
                  FormatUtils.formatYuan(summary.usedAmountYuan),
                  Colors.orangeAccent,
                ),
              ),
              Expanded(
                child: _buildStatItem(
                  '待归还',
                  FormatUtils.formatYuan(summary.pendingReturnAmountYuan),
                  Colors.redAccent,
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 剩余可用
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.account_balance_wallet, color: Colors.white, size: 18),
                const SizedBox(width: 8),
                Text(
                  '剩余可用: ${FormatUtils.formatYuan(summary.remainingAmountYuan)}',
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value, Color valueColor) {
    return Column(
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: 12,
            color: Colors.white.withOpacity(0.7),
          ),
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: TextStyle(
            fontSize: 18,
            fontWeight: FontWeight.bold,
            color: valueColor,
          ),
        ),
      ],
    );
  }

  Widget _buildActionButtons(SettlementWalletSummaryModel summary) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      child: Row(
        children: [
          Expanded(
            child: ElevatedButton.icon(
              onPressed: summary.remainingAmount > 0
                  ? () => _showUseDialog(summary)
                  : null,
              icon: const Icon(Icons.arrow_downward),
              label: const Text('使用沉淀款'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: ElevatedButton.icon(
              onPressed: summary.usedAmount > 0
                  ? () => _showReturnDialog(summary)
                  : null,
              icon: const Icon(Icons.arrow_upward),
              label: const Text('归还沉淀款'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.success,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSubordinateSection(AsyncValue<List<SubordinateBalanceModel>> subordinatesAsync) {
    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '下级余额明细',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 12),
          subordinatesAsync.when(
            data: (subordinates) {
              if (subordinates.isEmpty) {
                return const Center(
                  child: Padding(
                    padding: EdgeInsets.all(32),
                    child: Text(
                      '暂无下级余额数据',
                      style: TextStyle(color: AppColors.textSecondary),
                    ),
                  ),
                );
              }
              return ListView.separated(
                shrinkWrap: true,
                physics: const NeverScrollableScrollPhysics(),
                itemCount: subordinates.length,
                separatorBuilder: (_, __) => const Divider(height: 1),
                itemBuilder: (context, index) {
                  final item = subordinates[index];
                  return ListTile(
                    contentPadding: EdgeInsets.zero,
                    leading: CircleAvatar(
                      backgroundColor: AppColors.primary.withOpacity(0.1),
                      child: Text(
                        item.agentName.isNotEmpty ? item.agentName[0] : '?',
                        style: const TextStyle(color: AppColors.primary),
                      ),
                    ),
                    title: Text(item.agentName),
                    subtitle: Text('ID: ${item.agentId}'),
                    trailing: Text(
                      FormatUtils.formatYuan(item.availableBalanceYuan),
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                        color: AppColors.primary,
                      ),
                    ),
                  );
                },
              );
            },
            loading: () => const Center(
              child: Padding(
                padding: EdgeInsets.all(32),
                child: CircularProgressIndicator(),
              ),
            ),
            error: (e, _) => Center(
              child: Padding(
                padding: const EdgeInsets.all(32),
                child: Text('加载失败: $e'),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorCard(String message) {
    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        color: Colors.red.shade50,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          const Icon(Icons.error_outline, color: Colors.red),
          const SizedBox(width: 8),
          Expanded(child: Text(message)),
        ],
      ),
    );
  }

  void _showUseDialog(SettlementWalletSummaryModel summary) {
    _amountController.text = '';
    _remarkController.text = '';

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('使用沉淀款'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '可用额度: ${FormatUtils.formatYuan(summary.remainingAmountYuan)}',
              style: const TextStyle(color: AppColors.success),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _amountController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: const InputDecoration(
                labelText: '使用金额 (元)',
                hintText: '请输入金额',
                border: OutlineInputBorder(),
                prefixText: '¥',
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _remarkController,
              decoration: const InputDecoration(
                labelText: '备注',
                hintText: '可选',
                border: OutlineInputBorder(),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => _submitUse(summary),
            child: const Text('确认使用'),
          ),
        ],
      ),
    );
  }

  void _showReturnDialog(SettlementWalletSummaryModel summary) {
    _amountController.text = '';
    _remarkController.text = '';

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('归还沉淀款'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              '待归还: ${FormatUtils.formatYuan(summary.pendingReturnAmountYuan)}',
              style: const TextStyle(color: AppColors.warning),
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _amountController,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: const InputDecoration(
                labelText: '归还金额 (元)',
                hintText: '请输入金额',
                border: OutlineInputBorder(),
                prefixText: '¥',
              ),
            ),
            const SizedBox(height: 12),
            TextField(
              controller: _remarkController,
              decoration: const InputDecoration(
                labelText: '备注',
                hintText: '可选',
                border: OutlineInputBorder(),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () => _submitReturn(),
            child: const Text('确认归还'),
          ),
        ],
      ),
    );
  }

  Future<void> _submitUse(SettlementWalletSummaryModel summary) async {
    final amountText = _amountController.text.trim();
    if (amountText.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入金额')),
      );
      return;
    }

    final amount = double.tryParse(amountText);
    if (amount == null || amount <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入有效金额')),
      );
      return;
    }

    if (amount * 100 > summary.remainingAmount) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('使用金额超过可用额度')),
      );
      return;
    }

    try {
      final walletService = ref.read(walletServiceProvider);
      await walletService.useSettlement(
        amount: (amount * 100).round(),
        remark: _remarkController.text.trim().isEmpty ? null : _remarkController.text.trim(),
      );
      if (mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('使用成功')),
        );
        ref.invalidate(settlementWalletSummaryProvider);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('使用失败: $e')),
        );
      }
    }
  }

  Future<void> _submitReturn() async {
    final amountText = _amountController.text.trim();
    if (amountText.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入金额')),
      );
      return;
    }

    final amount = double.tryParse(amountText);
    if (amount == null || amount <= 0) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入有效金额')),
      );
      return;
    }

    try {
      final walletService = ref.read(walletServiceProvider);
      await walletService.returnSettlement(
        amount: (amount * 100).round(),
        remark: _remarkController.text.trim().isEmpty ? null : _remarkController.text.trim(),
      );
      if (mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('归还成功')),
        );
        ref.invalidate(settlementWalletSummaryProvider);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('归还失败: $e')),
        );
      }
    }
  }
}
