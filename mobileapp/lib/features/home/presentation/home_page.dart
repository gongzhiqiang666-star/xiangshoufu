import 'package:flutter/material.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';

/// 首页
class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('享收付'),
        leading: IconButton(
          icon: const Icon(Icons.menu),
          onPressed: () {},
        ),
        actions: [
          Stack(
            children: [
              IconButton(
                icon: const Icon(Icons.notifications_outlined),
                onPressed: () {},
              ),
              Positioned(
                right: 8,
                top: 8,
                child: Container(
                  padding: const EdgeInsets.all(4),
                  decoration: const BoxDecoration(
                    color: AppColors.danger,
                    shape: BoxShape.circle,
                  ),
                  child: const Text(
                    '3',
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await Future.delayed(const Duration(seconds: 1));
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              _buildBanner(),
              const SizedBox(height: AppSpacing.md),
              _buildTodayProfit(),
              const SizedBox(height: AppSpacing.md),
              _buildProfitDetails(),
              const SizedBox(height: AppSpacing.md),
              _buildQuickActions(),
              const SizedBox(height: AppSpacing.md),
              _buildRecentTransactions(),
              const SizedBox(height: AppSpacing.lg),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildBanner() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      height: 140,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(12),
        gradient: const LinearGradient(
          colors: [AppColors.primary, AppColors.primaryDark],
        ),
      ),
      child: const Center(
        child: Text(
          '欢迎使用享收付',
          style: TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold),
        ),
      ),
    );
  }

  Widget _buildTodayProfit() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          const Text('今日收益', style: TextStyle(fontSize: 14, color: AppColors.textSecondary)),
          const SizedBox(height: 8),
          const Text(
            '¥ 1,234.56',
            style: TextStyle(fontSize: 32, fontWeight: FontWeight.bold, color: AppColors.textPrimary),
          ),
          const SizedBox(height: 4),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('较昨日 ', style: TextStyle(fontSize: 12, color: AppColors.textTertiary)),
              Text('↑12.5%', style: TextStyle(fontSize: 12, color: AppColors.success, fontWeight: FontWeight.w500)),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildProfitDetails() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      child: GridView.count(
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        crossAxisCount: 2,
        crossAxisSpacing: AppSpacing.cardGap,
        mainAxisSpacing: AppSpacing.cardGap,
        childAspectRatio: 1.6,
        children: [
          _buildProfitCard('交易分润', 856.00, AppColors.profitTrade, Icons.swap_horiz),
          _buildProfitCard('押金返现', 150.00, AppColors.profitDeposit, Icons.monetization_on_outlined),
          _buildProfitCard('流量返现', 138.56, AppColors.profitSim, Icons.signal_cellular_alt),
          _buildProfitCard('激活奖励', 90.00, AppColors.profitReward, Icons.card_giftcard),
        ],
      ),
    );
  }

  Widget _buildProfitCard(String title, double amount, Color color, IconData icon) {
    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Row(
            children: [
              Container(
                width: 32,
                height: 32,
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Icon(icon, color: color, size: 18),
              ),
              const SizedBox(width: 8),
              Text(title, style: const TextStyle(fontSize: 13, color: AppColors.textSecondary)),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            FormatUtils.formatYuan(amount),
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: color),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActions() {
    final actions = [
      {'icon': Icons.devices, 'label': '终端'},
      {'icon': Icons.people_outline, 'label': '商户'},
      {'icon': Icons.analytics_outlined, 'label': '数据'},
      {'icon': Icons.account_balance_wallet_outlined, 'label': '钱包'},
      {'icon': Icons.upload_outlined, 'label': '代扣'},
      {'icon': Icons.image_outlined, 'label': '海报'},
      {'icon': Icons.notifications_outlined, 'label': '消息'},
      {'icon': Icons.person_outline, 'label': '我的'},
    ];

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: GridView.builder(
        shrinkWrap: true,
        physics: const NeverScrollableScrollPhysics(),
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 4,
          childAspectRatio: 1,
        ),
        itemCount: actions.length,
        itemBuilder: (context, index) {
          final action = actions[index];
          return Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 44,
                height: 44,
                decoration: BoxDecoration(
                  color: AppColors.background,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(action['icon'] as IconData, color: AppColors.primary, size: 24),
              ),
              const SizedBox(height: 6),
              Text(action['label'] as String, style: const TextStyle(fontSize: 12, color: AppColors.textSecondary)),
            ],
          );
        },
      ),
    );
  }

  Widget _buildRecentTransactions() {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(AppSpacing.md),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('最近交易', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600)),
                Row(
                  children: const [
                    Text('查看更多', style: TextStyle(fontSize: 13, color: AppColors.textTertiary)),
                    Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
                  ],
                ),
              ],
            ),
          ),
          const Divider(height: 1, color: AppColors.divider),
          _buildTransactionItem('张三商店', '刷卡', 1500.00, '10:30'),
          const Divider(height: 1, indent: 16, endIndent: 16, color: AppColors.divider),
          _buildTransactionItem('李四超市', '微信', 320.00, '10:25'),
          const Divider(height: 1, indent: 16, endIndent: 16, color: AppColors.divider),
          _buildTransactionItem('王五便利店', '支付宝', 850.00, '10:20'),
        ],
      ),
    );
  }

  Widget _buildTransactionItem(String name, String type, double amount, String time) {
    IconData icon;
    Color iconColor;

    switch (type) {
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

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: iconColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: Icon(icon, color: iconColor, size: 20),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(name, style: const TextStyle(fontSize: 15, fontWeight: FontWeight.w500)),
                const SizedBox(height: 2),
                Text(time, style: const TextStyle(fontSize: 12, color: AppColors.textTertiary)),
              ],
            ),
          ),
          Text(FormatUtils.formatYuan(amount), style: const TextStyle(fontSize: 15, fontWeight: FontWeight.w600)),
        ],
      ),
    );
  }
}
