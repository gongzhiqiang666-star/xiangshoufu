import 'package:flutter/material.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';

/// 钱包页面
class WalletPage extends StatefulWidget {
  const WalletPage({super.key});

  @override
  State<WalletPage> createState() => _WalletPageState();
}

class _WalletPageState extends State<WalletPage> {
  String _selectedChannel = '全部';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(title: const Text('我的钱包')),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // 总资产卡片
            _buildTotalAssets(),

            // 通道筛选
            _buildChannelFilter(),

            // 钱包卡片列表
            _buildWalletList(),

            // 底部操作按钮
            _buildBottomActions(),
          ],
        ),
      ),
    );
  }

  Widget _buildTotalAssets() {
    return Container(
      width: double.infinity,
      margin: const EdgeInsets.all(AppSpacing.md),
      padding: const EdgeInsets.all(AppSpacing.lg),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [AppColors.primary, AppColors.primaryDark],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: AppColors.primary.withOpacity(0.3),
            blurRadius: 20,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          const Text(
            '总资产',
            style: TextStyle(
              fontSize: 14,
              color: Colors.white70,
            ),
          ),
          const SizedBox(height: 8),
          const Text(
            '¥ 12,345.67',
            style: TextStyle(
              fontSize: 36,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            '累计提现: ¥88,500.00',
            style: TextStyle(
              fontSize: 13,
              color: Colors.white.withOpacity(0.8),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildChannelFilter() {
    final channels = ['全部', '拉卡拉', '随行付', '乐刷'];

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      child: Row(
        children: [
          const Text(
            '通道筛选: ',
            style: TextStyle(
              fontSize: 14,
              color: AppColors.textSecondary,
            ),
          ),
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(6),
              border: Border.all(color: AppColors.border),
            ),
            child: DropdownButton<String>(
              value: _selectedChannel,
              underline: const SizedBox(),
              isDense: true,
              icon: const Icon(Icons.keyboard_arrow_down, size: 20),
              items: channels.map((channel) {
                return DropdownMenuItem(
                  value: channel,
                  child: Text(channel, style: const TextStyle(fontSize: 14)),
                );
              }).toList(),
              onChanged: (value) {
                setState(() {
                  _selectedChannel = value ?? '全部';
                });
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildWalletList() {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        children: [
          _buildWalletCard(
            name: '分润钱包',
            channel: '拉卡拉',
            balance: 5680.00,
            threshold: 100,
            canWithdraw: true,
            gradientColors: AppColors.walletProfitGradient,
          ),
          _buildWalletCard(
            name: '服务费钱包',
            channel: '拉卡拉',
            balance: 3200.00,
            threshold: 200,
            canWithdraw: true,
            gradientColors: AppColors.walletServiceGradient,
          ),
          _buildWalletCard(
            name: '奖励钱包',
            channel: '拉卡拉',
            balance: 80.00,
            threshold: 100,
            canWithdraw: false,
            gradientColors: AppColors.walletRewardGradient,
          ),
        ],
      ),
    );
  }

  Widget _buildWalletCard({
    required String name,
    required String channel,
    required double balance,
    required double threshold,
    required bool canWithdraw,
    required List<Color> gradientColors,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppSpacing.cardGap),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: gradientColors,
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: gradientColors.first.withOpacity(0.3),
            blurRadius: 12,
            offset: const Offset(0, 6),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 头部
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                name,
                style: TextStyle(
                  fontSize: 15,
                  color: Colors.white.withOpacity(0.9),
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  channel,
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.white,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 余额
          Text(
            FormatUtils.formatYuan(balance),
            style: const TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '提现门槛: ¥${threshold.toStringAsFixed(0)}',
            style: TextStyle(
              fontSize: 12,
              color: Colors.white.withOpacity(0.7),
            ),
          ),
          const SizedBox(height: 16),

          // 提现按钮
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: canWithdraw ? () {} : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: gradientColors.first,
                disabledBackgroundColor: Colors.white.withOpacity(0.5),
                disabledForegroundColor: gradientColors.first.withOpacity(0.5),
                elevation: 0,
              ),
              child: Text(canWithdraw ? '申请提现' : '未达提现门槛'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBottomActions() {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Row(
        children: [
          Expanded(
            child: OutlinedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.list_alt, size: 18),
              label: const Text('钱包流水'),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: OutlinedButton.icon(
              onPressed: () {},
              icon: const Icon(Icons.history, size: 18),
              label: const Text('提现记录'),
            ),
          ),
        ],
      ),
    );
  }
}
