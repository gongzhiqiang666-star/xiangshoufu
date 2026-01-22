import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../../../router/app_router.dart';
import '../data/models/wallet_model.dart';
import 'providers/wallet_provider.dart';

/// 钱包页面
class WalletPage extends ConsumerStatefulWidget {
  const WalletPage({super.key});

  @override
  ConsumerState<WalletPage> createState() => _WalletPageState();
}

class _WalletPageState extends ConsumerState<WalletPage> {
  String _selectedChannel = '全部';

  @override
  Widget build(BuildContext context) {
    final summaryAsync = ref.watch(walletSummaryProvider);
    final walletsAsync = ref.watch(walletsProvider);
    final configAsync = ref.watch(walletConfigProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: const Text('我的钱包'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              ref.invalidate(walletSummaryProvider);
              ref.invalidate(walletsProvider);
              ref.invalidate(walletConfigProvider);
            },
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          ref.invalidate(walletSummaryProvider);
          ref.invalidate(walletsProvider);
          ref.invalidate(walletConfigProvider);
        },
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          child: Column(
            children: [
              // 总资产卡片
              summaryAsync.when(
                data: (summary) => _buildTotalAssets(summary),
                loading: () => _buildTotalAssetsLoading(),
                error: (e, _) => _buildErrorCard('加载失败: $e'),
              ),

              // 特殊钱包入口（充值钱包、沉淀钱包）
              configAsync.when(
                data: (config) => _buildSpecialWalletEntries(config),
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
              ),

              // 通道筛选
              walletsAsync.when(
                data: (wallets) => _buildChannelFilter(wallets),
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
              ),

              // 钱包卡片列表
              walletsAsync.when(
                data: (wallets) => _buildWalletList(wallets),
                loading: () => _buildWalletListLoading(),
                error: (e, _) => _buildErrorCard('加载失败: $e'),
              ),

              // 底部操作按钮
              _buildBottomActions(),

              const SizedBox(height: 20),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTotalAssets(WalletSummaryModel summary) {
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
          Text(
            FormatUtils.formatYuan(summary.totalBalanceYuan),
            style: const TextStyle(
              fontSize: 36,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              _buildAssetItem('可用', summary.totalAvailableYuan),
              Container(
                width: 1,
                height: 30,
                color: Colors.white24,
              ),
              _buildAssetItem('冻结', summary.totalFrozenYuan),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildAssetItem(String label, double amount) {
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
          FormatUtils.formatYuan(amount),
          style: const TextStyle(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
      ],
    );
  }

  Widget _buildTotalAssetsLoading() {
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
      ),
      child: const Center(
        child: Padding(
          padding: EdgeInsets.all(20),
          child: CircularProgressIndicator(color: Colors.white),
        ),
      ),
    );
  }

  Widget _buildSpecialWalletEntries(AgentWalletConfigModel config) {
    // 只有开通了相关钱包才显示入口
    if (!config.chargingWalletEnabled && !config.settlementWalletEnabled) {
      return const SizedBox.shrink();
    }

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      child: Row(
        children: [
          if (config.chargingWalletEnabled)
            Expanded(
              child: _buildSpecialWalletCard(
                title: '充值钱包',
                icon: Icons.account_balance_wallet,
                color: const Color(0xFF11998e),
                onTap: () => context.push(RoutePaths.chargingWallet),
              ),
            ),
          if (config.chargingWalletEnabled && config.settlementWalletEnabled)
            const SizedBox(width: 12),
          if (config.settlementWalletEnabled)
            Expanded(
              child: _buildSpecialWalletCard(
                title: '沉淀钱包',
                icon: Icons.savings,
                color: const Color(0xFF667eea),
                onTap: () => context.push(RoutePaths.settlementWallet),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildSpecialWalletCard({
    required String title,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(AppSpacing.md),
        decoration: BoxDecoration(
          color: color.withOpacity(0.1),
          borderRadius: BorderRadius.circular(12),
          border: Border.all(color: color.withOpacity(0.3)),
        ),
        child: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: color,
                borderRadius: BorderRadius.circular(8),
              ),
              child: Icon(icon, color: Colors.white, size: 20),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Text(
                title,
                style: TextStyle(
                  color: color,
                  fontWeight: FontWeight.bold,
                  fontSize: 14,
                ),
              ),
            ),
            Icon(Icons.arrow_forward_ios, color: color, size: 16),
          ],
        ),
      ),
    );
  }

  Widget _buildChannelFilter(List<WalletModel> wallets) {
    // 提取所有通道名称
    final channelSet = <String>{'全部'};
    for (final wallet in wallets) {
      if (wallet.channelName.isNotEmpty) {
        channelSet.add(wallet.channelName);
      }
    }
    final channels = channelSet.toList();

    return Container(
      padding: const EdgeInsets.all(AppSpacing.md),
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
              value: channels.contains(_selectedChannel) ? _selectedChannel : '全部',
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

  Widget _buildWalletList(List<WalletModel> wallets) {
    // 按通道筛选
    final filteredWallets = _selectedChannel == '全部'
        ? wallets
        : wallets.where((w) => w.channelName == _selectedChannel).toList();

    if (filteredWallets.isEmpty) {
      return Container(
        padding: const EdgeInsets.all(32),
        child: const Center(
          child: Text(
            '暂无钱包数据',
            style: TextStyle(color: AppColors.textSecondary),
          ),
        ),
      );
    }

    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        children: filteredWallets.map((wallet) {
          return _buildWalletCard(wallet);
        }).toList(),
      ),
    );
  }

  Widget _buildWalletCard(WalletModel wallet) {
    final gradientColors = _getGradientColors(wallet.walletType);
    final canWithdraw = wallet.available > 0;

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
                wallet.walletTypeName,
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
                  wallet.channelName.isNotEmpty ? wallet.channelName : '通用',
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
            FormatUtils.formatYuan(wallet.balanceYuan),
            style: const TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 4),
          Row(
            children: [
              Text(
                '可用: ${FormatUtils.formatYuan(wallet.availableYuan)}',
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.white.withOpacity(0.7),
                ),
              ),
              if (wallet.frozen > 0) ...[
                const SizedBox(width: 16),
                Text(
                  '冻结: ${FormatUtils.formatYuan(wallet.frozenYuan)}',
                  style: TextStyle(
                    fontSize: 12,
                    color: Colors.white.withOpacity(0.7),
                  ),
                ),
              ],
            ],
          ),
          const SizedBox(height: 16),

          // 提现按钮
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: canWithdraw
                  ? () => context.push(
                        RoutePaths.withdraw,
                        extra: wallet.id.toString(),
                      )
                  : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: gradientColors.first,
                disabledBackgroundColor: Colors.white.withOpacity(0.5),
                disabledForegroundColor: gradientColors.first.withOpacity(0.5),
                elevation: 0,
              ),
              child: Text(canWithdraw ? '申请提现' : '暂无可提金额'),
            ),
          ),
        ],
      ),
    );
  }

  List<Color> _getGradientColors(int walletType) {
    switch (walletType) {
      case 1: // 分润钱包
        return AppColors.walletProfitGradient;
      case 2: // 服务费钱包
        return AppColors.walletServiceGradient;
      case 3: // 奖励钱包
        return AppColors.walletRewardGradient;
      case 4: // 充值钱包
        return [const Color(0xFF11998e), const Color(0xFF38ef7d)];
      case 5: // 沉淀钱包
        return [const Color(0xFF667eea), const Color(0xFF764ba2)];
      default:
        return [AppColors.primary, AppColors.primaryDark];
    }
  }

  Widget _buildWalletListLoading() {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        children: List.generate(
          3,
          (index) => Container(
            margin: const EdgeInsets.only(bottom: AppSpacing.cardGap),
            height: 150,
            decoration: BoxDecoration(
              color: Colors.grey.shade200,
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Center(
              child: CircularProgressIndicator(),
            ),
          ),
        ),
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
          TextButton(
            onPressed: () {
              ref.invalidate(walletSummaryProvider);
              ref.invalidate(walletsProvider);
            },
            child: const Text('重试'),
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
              onPressed: () {
                // TODO: 跳转到钱包流水页面
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('钱包流水功能开发中')),
                );
              },
              icon: const Icon(Icons.list_alt, size: 18),
              label: const Text('钱包流水'),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: OutlinedButton.icon(
              onPressed: () {
                // TODO: 跳转到提现记录页面
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('提现记录功能开发中')),
                );
              },
              icon: const Icon(Icons.history, size: 18),
              label: const Text('提现记录'),
            ),
          ),
        ],
      ),
    );
  }
}
