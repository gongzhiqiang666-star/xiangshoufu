import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../data/models/wallet_model.dart';
import 'providers/wallet_provider.dart';
import '../data/services/wallet_service.dart';

/// 充值钱包页面
class ChargingWalletPage extends ConsumerStatefulWidget {
  const ChargingWalletPage({super.key});

  @override
  ConsumerState<ChargingWalletPage> createState() => _ChargingWalletPageState();
}

class _ChargingWalletPageState extends ConsumerState<ChargingWalletPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final _amountController = TextEditingController();
  final _remarkController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _amountController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final summaryAsync = ref.watch(chargingWalletSummaryProvider);
    final configAsync = ref.watch(walletConfigProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('充值钱包'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              ref.invalidate(chargingWalletSummaryProvider);
              ref.invalidate(walletConfigProvider);
            },
          ),
        ],
        bottom: TabBar(
          controller: _tabController,
          tabs: const [
            Tab(text: '钱包概览'),
            Tab(text: '奖励发放'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          // 钱包概览
          _buildOverviewTab(summaryAsync, configAsync),
          // 奖励发放
          _buildRewardTab(summaryAsync),
        ],
      ),
    );
  }

  Widget _buildOverviewTab(
    AsyncValue<ChargingWalletSummaryModel> summaryAsync,
    AsyncValue<AgentWalletConfigModel> configAsync,
  ) {
    return RefreshIndicator(
      onRefresh: () async {
        ref.invalidate(chargingWalletSummaryProvider);
        ref.invalidate(walletConfigProvider);
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

            // 配置信息
            configAsync.when(
              data: (config) => _buildConfigCard(config),
              loading: () => const SizedBox.shrink(),
              error: (_, __) => const SizedBox.shrink(),
            ),

            // 操作按钮
            summaryAsync.when(
              data: (summary) => _buildActionButtons(summary),
              loading: () => const SizedBox.shrink(),
              error: (_, __) => const SizedBox.shrink(),
            ),

            // 奖励统计
            summaryAsync.when(
              data: (summary) => _buildRewardStatsCard(summary),
              loading: () => const SizedBox.shrink(),
              error: (_, __) => const SizedBox.shrink(),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSummaryCard(ChargingWalletSummaryModel summary) {
    return Container(
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF11998e), Color(0xFF38ef7d)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF11998e).withOpacity(0.3),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        children: [
          // 标题
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: const [
              Icon(Icons.account_balance_wallet, color: Colors.white70, size: 20),
              SizedBox(width: 8),
              Text(
                '充值钱包余额',
                style: TextStyle(
                  fontSize: 14,
                  color: Colors.white70,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),

          // 余额
          Text(
            FormatUtils.formatYuan(summary.balanceYuan),
            style: const TextStyle(
              fontSize: 36,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 20),

          // 奖励总金额
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.2),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.card_giftcard, color: Colors.white, size: 18),
                const SizedBox(width: 8),
                Text(
                  '累计奖励: ${FormatUtils.formatYuan(summary.totalRewardYuan)}',
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

  Widget _buildConfigCard(AgentWalletConfigModel config) {
    if (!config.chargingWalletEnabled) {
      return Container(
        margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
        padding: const EdgeInsets.all(AppSpacing.md),
        decoration: BoxDecoration(
          color: Colors.orange.shade50,
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: Colors.orange.shade200),
        ),
        child: Row(
          children: [
            Icon(Icons.info_outline, color: Colors.orange.shade700),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                '充值钱包未开通，请联系上级或管理员开通',
                style: TextStyle(color: Colors.orange.shade700),
              ),
            ),
          ],
        ),
      );
    }

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
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
          Row(
            children: [
              Icon(Icons.check_circle, color: Colors.green.shade600, size: 20),
              const SizedBox(width: 8),
              const Text(
                '充值钱包已开通',
                style: TextStyle(
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
              ),
            ],
          ),
          if (config.chargingWalletLimit > 0) ...[
            const SizedBox(height: 8),
            Text(
              '充值限额: ${FormatUtils.formatYuan(config.chargingWalletLimitYuan)}',
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 13,
              ),
            ),
          ],
          if (config.enabledAt != null) ...[
            const SizedBox(height: 4),
            Text(
              '开通时间: ${config.enabledAt}',
              style: const TextStyle(
                color: AppColors.textSecondary,
                fontSize: 13,
              ),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildActionButtons(ChargingWalletSummaryModel summary) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Row(
        children: [
          Expanded(
            child: ElevatedButton.icon(
              onPressed: () => _showDepositDialog(),
              icon: const Icon(Icons.add_circle_outline),
              label: const Text('申请充值'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 14),
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: ElevatedButton.icon(
              onPressed: summary.balance > 0
                  ? () {
                      _tabController.animateTo(1);
                    }
                  : null,
              icon: const Icon(Icons.card_giftcard),
              label: const Text('发放奖励'),
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.success,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 14),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildRewardStatsCard(ChargingWalletSummaryModel summary) {
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
            '奖励统计',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: _buildStatItem(
                  '手动发放',
                  FormatUtils.formatYuan(summary.totalIssuedYuan),
                  Icons.send,
                  Colors.blue,
                ),
              ),
              Container(
                width: 1,
                height: 40,
                color: Colors.grey.shade200,
              ),
              Expanded(
                child: _buildStatItem(
                  '系统自动',
                  FormatUtils.formatYuan(summary.totalAutoRewardYuan),
                  Icons.autorenew,
                  Colors.green,
                ),
              ),
              Container(
                width: 1,
                height: 40,
                color: Colors.grey.shade200,
              ),
              Expanded(
                child: _buildStatItem(
                  '奖励总计',
                  FormatUtils.formatYuan(summary.totalRewardYuan),
                  Icons.summarize,
                  Colors.orange,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildStatItem(String label, String value, IconData icon, Color color) {
    return Column(
      children: [
        Icon(icon, color: color, size: 24),
        const SizedBox(height: 8),
        Text(
          value,
          style: TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: const TextStyle(
            fontSize: 12,
            color: AppColors.textSecondary,
          ),
        ),
      ],
    );
  }

  Widget _buildRewardTab(AsyncValue<ChargingWalletSummaryModel> summaryAsync) {
    return summaryAsync.when(
      data: (summary) => _buildIssueRewardForm(summary),
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (e, _) => Center(child: Text('加载失败: $e')),
    );
  }

  Widget _buildIssueRewardForm(ChargingWalletSummaryModel summary) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 可用余额提示
          Container(
            padding: const EdgeInsets.all(AppSpacing.md),
            decoration: BoxDecoration(
              color: Colors.green.shade50,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(color: Colors.green.shade200),
            ),
            child: Row(
              children: [
                Icon(Icons.account_balance_wallet, color: Colors.green.shade700),
                const SizedBox(width: 12),
                Text(
                  '可用余额: ${FormatUtils.formatYuan(summary.balanceYuan)}',
                  style: TextStyle(
                    color: Colors.green.shade700,
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // 说明
          Container(
            padding: const EdgeInsets.all(AppSpacing.md),
            decoration: BoxDecoration(
              color: Colors.blue.shade50,
              borderRadius: BorderRadius.circular(12),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.info_outline, color: Colors.blue.shade700, size: 20),
                    const SizedBox(width: 8),
                    Text(
                      '奖励发放说明',
                      style: TextStyle(
                        color: Colors.blue.shade700,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 8),
                Text(
                  '1. 只能给直属下级代理商发放奖励\n'
                  '2. 发放金额从充值钱包扣除\n'
                  '3. 奖励将进入下级的奖励钱包',
                  style: TextStyle(
                    color: Colors.blue.shade700,
                    fontSize: 13,
                    height: 1.5,
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // 表单
          const Text(
            '发放奖励',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 16),

          // 选择下级代理商（简化版，实际需要选择器）
          TextField(
            decoration: const InputDecoration(
              labelText: '下级代理商ID',
              hintText: '请输入下级代理商ID',
              border: OutlineInputBorder(),
              prefixIcon: Icon(Icons.person),
            ),
            keyboardType: TextInputType.number,
            onChanged: (value) {
              // 保存代理商ID
            },
          ),
          const SizedBox(height: 16),

          TextField(
            controller: _amountController,
            keyboardType: const TextInputType.numberWithOptions(decimal: true),
            decoration: const InputDecoration(
              labelText: '发放金额 (元)',
              hintText: '请输入金额',
              border: OutlineInputBorder(),
              prefixText: '¥',
              prefixIcon: Icon(Icons.monetization_on),
            ),
          ),
          const SizedBox(height: 16),

          TextField(
            controller: _remarkController,
            decoration: const InputDecoration(
              labelText: '备注',
              hintText: '可选，填写发放原因',
              border: OutlineInputBorder(),
              prefixIcon: Icon(Icons.note),
            ),
            maxLines: 2,
          ),
          const SizedBox(height: 24),

          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: summary.balance > 0 ? () => _submitReward(summary) : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: AppColors.primary,
                foregroundColor: Colors.white,
                padding: const EdgeInsets.symmetric(vertical: 16),
              ),
              child: const Text(
                '确认发放',
                style: TextStyle(fontSize: 16),
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

  void _showDepositDialog() {
    _amountController.text = '';
    _remarkController.text = '';
    int selectedMethod = 1; // 默认银行转账

    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, setDialogState) => AlertDialog(
          title: const Text('申请充值'),
          content: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  '选择付款方式',
                  style: TextStyle(fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 8),
                Wrap(
                  spacing: 8,
                  children: [
                    ChoiceChip(
                      label: const Text('银行转账'),
                      selected: selectedMethod == 1,
                      onSelected: (selected) {
                        setDialogState(() => selectedMethod = 1);
                      },
                    ),
                    ChoiceChip(
                      label: const Text('微信'),
                      selected: selectedMethod == 2,
                      onSelected: (selected) {
                        setDialogState(() => selectedMethod = 2);
                      },
                    ),
                    ChoiceChip(
                      label: const Text('支付宝'),
                      selected: selectedMethod == 3,
                      onSelected: (selected) {
                        setDialogState(() => selectedMethod = 3);
                      },
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: _amountController,
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: const InputDecoration(
                    labelText: '充值金额 (元)',
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
                const SizedBox(height: 12),
                Text(
                  '提交后请按指定方式付款，管理员确认后到账',
                  style: TextStyle(
                    color: Colors.grey.shade600,
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: const Text('取消'),
            ),
            ElevatedButton(
              onPressed: () => _submitDeposit(selectedMethod),
              child: const Text('提交申请'),
            ),
          ],
        ),
      ),
    );
  }

  Future<void> _submitDeposit(int paymentMethod) async {
    final amountText = _amountController.text.trim();
    if (amountText.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入金额')),
      );
      return;
    }

    final amount = double.tryParse(amountText);
    if (amount == null || amount < 1) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('金额最少1元')),
      );
      return;
    }

    try {
      final walletService = ref.read(walletServiceProvider);
      final depositNo = await walletService.createChargingDeposit(
        amount: (amount * 100).round(),
        paymentMethod: paymentMethod,
        remark: _remarkController.text.trim().isEmpty ? null : _remarkController.text.trim(),
      );
      if (mounted) {
        Navigator.pop(context);
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('充值申请已提交，单号: $depositNo')),
        );
        ref.invalidate(chargingWalletSummaryProvider);
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('提交失败: $e')),
        );
      }
    }
  }

  Future<void> _submitReward(ChargingWalletSummaryModel summary) async {
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

    if (amount * 100 > summary.balance) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('发放金额超过可用余额')),
      );
      return;
    }

    // TODO: 获取选择的代理商ID，这里暂时提示
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('请先选择下级代理商')),
    );
  }
}
