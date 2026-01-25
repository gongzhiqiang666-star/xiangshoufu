import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../../../core/utils/format_utils.dart';
import '../../../router/app_router.dart';
import '../../marketing/presentation/widgets/banner_carousel.dart';
import '../../message/presentation/providers/message_provider.dart';
import '../domain/home_model.dart';
import 'providers/home_provider.dart';

/// 首页
class HomePage extends ConsumerStatefulWidget {
  const HomePage({super.key});

  @override
  ConsumerState<HomePage> createState() => _HomePageState();
}

class _HomePageState extends ConsumerState<HomePage> {
  @override
  Widget build(BuildContext context) {
    final homeState = ref.watch(homeProvider);
    final unreadCount = ref.watch(unreadCountProvider);

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        automaticallyImplyLeading: false, // 去掉左上角返回按钮
        title: const Text('享收付'),
        actions: [
          Stack(
            children: [
              IconButton(
                icon: const Icon(Icons.notifications_outlined),
                onPressed: () => context.push(RoutePaths.message),
              ),
              unreadCount.when(
                data: (count) => count > 0
                    ? Positioned(
                        right: 8,
                        top: 8,
                        child: Container(
                          padding: const EdgeInsets.all(4),
                          decoration: const BoxDecoration(
                            color: AppColors.danger,
                            shape: BoxShape.circle,
                          ),
                          child: Text(
                            count > 99 ? '99+' : '$count',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 10,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                      )
                    : const SizedBox.shrink(),
                loading: () => const SizedBox.shrink(),
                error: (_, __) => const SizedBox.shrink(),
              ),
            ],
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () async {
          await ref.read(homeProvider.notifier).refresh();
        },
        child: homeState.isLoading && homeState.overview == null
            ? const Center(child: CircularProgressIndicator())
            : homeState.error != null && homeState.overview == null
                ? _buildErrorWidget(homeState.error!)
                : SingleChildScrollView(
                    physics: const AlwaysScrollableScrollPhysics(),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        _buildBanner(),
                        const SizedBox(height: AppSpacing.md),
                        _buildTodayProfit(homeState),
                        const SizedBox(height: AppSpacing.md),
                        _buildProfitDetails(homeState),
                        const SizedBox(height: AppSpacing.md),
                        _buildQuickActions(),
                        const SizedBox(height: AppSpacing.md),
                        _buildRecentTransactions(homeState),
                        const SizedBox(height: AppSpacing.lg),
                      ],
                    ),
                  ),
      ),
    );
  }

  Widget _buildErrorWidget(String error) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 48, color: AppColors.textTertiary),
          const SizedBox(height: 16),
          Text('加载失败', style: TextStyle(color: AppColors.textSecondary)),
          const SizedBox(height: 8),
          TextButton(
            onPressed: () => ref.read(homeProvider.notifier).refresh(),
            child: const Text('点击重试'),
          ),
        ],
      ),
    );
  }

  /// 构建Banner轮播
  Widget _buildBanner() {
    return BannerCarousel(
      height: 140,
      autoPlayInterval: 5000,
      onInternalLinkTap: (route) {
        context.push(route);
      },
    );
  }

  Widget _buildTodayProfit(HomeState state) {
    final overview = state.overview;
    final todayProfit = overview?.today.profitTotalYuan ?? 0;
    final changeRate = overview?.profitChangeRate ?? 0;
    final isGrowth = overview?.isProfitGrowth ?? true;

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
          Text(
            FormatUtils.formatYuan(todayProfit),
            style: const TextStyle(fontSize: 32, fontWeight: FontWeight.bold, color: AppColors.textPrimary),
          ),
          const SizedBox(height: 4),
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('较昨日 ', style: TextStyle(fontSize: 12, color: AppColors.textTertiary)),
              Text(
                '${isGrowth ? "↑" : "↓"}${changeRate.abs().toStringAsFixed(1)}%',
                style: TextStyle(
                  fontSize: 12,
                  color: isGrowth ? AppColors.success : AppColors.danger,
                  fontWeight: FontWeight.w500,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildProfitDetails(HomeState state) {
    final today = state.overview?.today;

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
          _buildProfitCard('交易分润', today?.profitTradeYuan ?? 0, AppColors.profitTrade, Icons.swap_horiz),
          _buildProfitCard('押金返现', today?.profitDepositYuan ?? 0, AppColors.profitDeposit, Icons.monetization_on_outlined),
          _buildProfitCard('流量返现', today?.profitSimYuan ?? 0, AppColors.profitSim, Icons.signal_cellular_alt),
          _buildProfitCard('激活奖励', today?.profitRewardYuan ?? 0, AppColors.profitReward, Icons.card_giftcard),
        ],
      ),
    );
  }

  Widget _buildProfitCard(String title, double amount, Color color, IconData icon) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.center,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Container(
                width: 24,
                height: 24,
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(6),
                ),
                child: Icon(icon, color: color, size: 14),
              ),
              const SizedBox(width: 6),
              Text(
                title,
                style: const TextStyle(fontSize: 12, color: AppColors.textSecondary),
              ),
            ],
          ),
          const SizedBox(height: 8),
          Text(
            FormatUtils.formatYuan(amount),
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: color),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActions() {
    final actions = [
      {'icon': Icons.devices, 'label': '终端', 'route': RoutePaths.terminal},
      {'icon': Icons.people_outline, 'label': '商户', 'route': RoutePaths.merchant},
      {'icon': Icons.group_add_outlined, 'label': '代理拓展', 'route': RoutePaths.agent},
      {'icon': Icons.account_balance_wallet_outlined, 'label': '钱包', 'route': RoutePaths.wallet},
      {'icon': Icons.upload_outlined, 'label': '代扣', 'route': RoutePaths.deduction},
      {'icon': Icons.image_outlined, 'label': '海报', 'route': RoutePaths.marketing},
      {'icon': Icons.notifications_outlined, 'label': '消息', 'route': RoutePaths.message},
      {'icon': Icons.person_outline, 'label': '我的', 'route': RoutePaths.profile},
    ];

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.transparent,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
           BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Material(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        clipBehavior: Clip.antiAlias,
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.md),
          child: GridView.builder(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
              crossAxisCount: 4,
              childAspectRatio: 0.85,
            ),
            itemCount: actions.length,
            itemBuilder: (context, index) {
              final action = actions[index];
              return InkWell(
                onTap: () {
                  final route = action['route'] as String;
                  if (route == RoutePaths.terminal ||
                      route == RoutePaths.wallet ||
                      route == RoutePaths.profile ||
                      route == RoutePaths.dataAnalysis) {
                    context.go(route);
                  } else {
                    context.push(route);
                  }
                },
                borderRadius: BorderRadius.circular(12),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Container(
                      width: 36,
                      height: 36,
                      decoration: BoxDecoration(
                        color: AppColors.background,
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: Icon(action['icon'] as IconData, color: AppColors.primary, size: 20),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      action['label'] as String,
                      style: const TextStyle(fontSize: 11, color: AppColors.textSecondary),
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              );
            },
          ),
        ),
      ),
    );
  }

  Widget _buildRecentTransactions(HomeState state) {
    final transactions = state.recentTransactions;

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
                GestureDetector(
                  onTap: () => context.push(RoutePaths.transaction),
                  child: Row(
                    children: const [
                      Text('查看更多', style: TextStyle(fontSize: 13, color: AppColors.textTertiary)),
                      Icon(Icons.chevron_right, size: 18, color: AppColors.textTertiary),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const Divider(height: 1, color: AppColors.divider),
          if (transactions.isEmpty)
            const Padding(
              padding: EdgeInsets.all(32),
              child: Text('暂无交易记录', style: TextStyle(color: AppColors.textTertiary)),
            )
          else
            ...transactions.asMap().entries.map((entry) {
              final index = entry.key;
              final tx = entry.value;
              return Column(
                children: [
                  _buildTransactionItem(tx.merchantName, tx.payTypeName, tx.amountYuan, tx.timeAgo),
                  if (index < transactions.length - 1)
                    const Divider(height: 1, indent: 16, endIndent: 16, color: AppColors.divider),
                ],
              );
            }).toList(),
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
