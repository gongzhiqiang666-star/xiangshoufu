import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../core/theme/app_colors.dart';
import '../../../core/theme/app_spacing.dart';
import '../data/models/settlement_price_model.dart';
import 'providers/settlement_price_provider.dart';

/// 下级代理商结算价列表页面
/// 入口：代理商详情页 → 结算价管理
/// 展示指定下级代理商各通道的结算价，可编辑
class AgentSettlementPriceListPage extends ConsumerStatefulWidget {
  final int agentId;
  final String? agentName;

  const AgentSettlementPriceListPage({
    super.key,
    required this.agentId,
    this.agentName,
  });

  @override
  ConsumerState<AgentSettlementPriceListPage> createState() => _AgentSettlementPriceListPageState();
}

class _AgentSettlementPriceListPageState extends ConsumerState<AgentSettlementPriceListPage> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    // 初始加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(agentSettlementPriceListProvider(widget.agentId).notifier).refresh();
    });

    // 滚动监听，加载更多
    _scrollController.addListener(_onScroll);
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent - 200) {
      ref.read(agentSettlementPriceListProvider(widget.agentId).notifier).loadMore();
    }
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(agentSettlementPriceListProvider(widget.agentId));

    return Scaffold(
      backgroundColor: AppColors.background,
      appBar: AppBar(
        title: Text(widget.agentName != null ? '${widget.agentName}的结算价' : '结算价管理'),
        actions: [
          IconButton(
            icon: const Icon(Icons.history),
            tooltip: '调价记录',
            onPressed: () => context.push('/agents/${widget.agentId}/price-change-logs'),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () => ref.read(agentSettlementPriceListProvider(widget.agentId).notifier).refresh(),
        child: _buildContent(state),
      ),
    );
  }

  Widget _buildContent(AgentSettlementPriceListState state) {
    if (state.isLoading && state.list.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.list.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 48, color: AppColors.danger),
            const SizedBox(height: 16),
            Text('加载失败: ${state.error}', textAlign: TextAlign.center),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => ref.read(agentSettlementPriceListProvider(widget.agentId).notifier).refresh(),
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (state.list.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.inbox_outlined, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              '暂无结算价配置',
              style: TextStyle(fontSize: 16, color: Colors.grey[600]),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      controller: _scrollController,
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: state.list.length + (state.hasMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index == state.list.length) {
          return const Padding(
            padding: EdgeInsets.all(16.0),
            child: Center(child: CircularProgressIndicator()),
          );
        }
        return _buildPriceCard(state.list[index]);
      },
    );
  }

  /// 构建结算价卡片（可编辑）
  Widget _buildPriceCard(SettlementPriceModel item) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppSpacing.md),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 通道头部（带编辑按钮）
          _buildCardHeader(item),
          const Divider(height: 1, color: AppColors.divider),

          // 费率配置
          _buildRateSection(item),

          // 押金返现配置
          if (item.depositCashbacks.isNotEmpty) ...[
            const Divider(height: 1, indent: 16, endIndent: 16, color: AppColors.divider),
            _buildDepositSection(item),
          ],

          // 流量费返现配置
          if (item.simCashbackTiers.any((tier) => tier.cashbackAmount > 0)) ...[
            const Divider(height: 1, indent: 16, endIndent: 16, color: AppColors.divider),
            _buildSimSection(item),
          ],
        ],
      ),
    );
  }

  /// 卡片头部：通道名称、状态和编辑按钮
  Widget _buildCardHeader(SettlementPriceModel item) {
    final isActive = item.status == 1;

    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Row(
        children: [
          // 通道图标
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: AppColors.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(10),
            ),
            child: const Icon(Icons.payment, color: AppColors.primary, size: 22),
          ),
          const SizedBox(width: 12),

          // 通道名称
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  item.channelName,
                  style: const TextStyle(
                    fontSize: 16,
                    fontWeight: FontWeight.w600,
                    color: AppColors.textPrimary,
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  '版本 v${item.version}',
                  style: const TextStyle(
                    fontSize: 12,
                    color: AppColors.textSecondary,
                  ),
                ),
              ],
            ),
          ),

          // 状态标签
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
            decoration: BoxDecoration(
              color: isActive ? AppColors.success.withOpacity(0.1) : AppColors.danger.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Text(
              item.statusName,
              style: TextStyle(
                fontSize: 12,
                color: isActive ? AppColors.success : AppColors.danger,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
          const SizedBox(width: 8),

          // 编辑按钮
          IconButton(
            icon: const Icon(Icons.edit_outlined, color: AppColors.primary),
            tooltip: '编辑结算价',
            onPressed: () => _navigateToEdit(item),
          ),
        ],
      ),
    );
  }

  /// 跳转到编辑页
  void _navigateToEdit(SettlementPriceModel item) {
    context.push(
      '/agents/${widget.agentId}/settlement-prices/${item.id}/edit',
      extra: {'agentName': widget.agentName, 'channelName': item.channelName},
    );
  }

  /// 费率配置区域
  Widget _buildRateSection(SettlementPriceModel item) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '费率配置',
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: item.rateConfigs.entries.map((e) =>
              _buildRateChip(e.key, '${e.value.rate}%')
            ).toList(),
          ),
        ],
      ),
    );
  }

  /// 费率标签
  Widget _buildRateChip(String label, String value) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.background,
        borderRadius: BorderRadius.circular(8),
        border: Border.all(color: AppColors.divider),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(
            label,
            style: const TextStyle(
              fontSize: 12,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(width: 6),
          Text(
            value,
            style: const TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w600,
              color: AppColors.primary,
            ),
          ),
        ],
      ),
    );
  }

  /// 押金返现配置区域
  Widget _buildDepositSection(SettlementPriceModel item) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '押金返现',
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: item.depositCashbacks.map((deposit) {
              return _buildDepositChip(deposit);
            }).toList(),
          ),
        ],
      ),
    );
  }

  /// 押金返现标签
  Widget _buildDepositChip(DepositCashbackItem item) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: AppColors.profitReward.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(
        '¥${item.depositAmountYuan.toStringAsFixed(0)} → 返¥${item.cashbackAmountYuan.toStringAsFixed(0)}',
        style: const TextStyle(
          fontSize: 13,
          fontWeight: FontWeight.w500,
          color: AppColors.profitReward,
        ),
      ),
    );
  }

  /// 流量费返现配置区域
  Widget _buildSimSection(SettlementPriceModel item) {
    return Padding(
      padding: const EdgeInsets.all(AppSpacing.md),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '流量费返现',
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: AppColors.textPrimary,
            ),
          ),
          const SizedBox(height: 12),
          Row(
            children: item.simCashbackTiers
                .where((tier) => tier.cashbackAmount > 0)
                .map((tier) => Expanded(child: _buildSimItem(tier.tierName, tier.cashbackAmountYuan)))
                .toList(),
          ),
        ],
      ),
    );
  }

  /// 流量费返现项
  Widget _buildSimItem(String label, double amount) {
    return Container(
      margin: const EdgeInsets.only(right: 8),
      padding: const EdgeInsets.symmetric(vertical: 8),
      decoration: BoxDecoration(
        color: AppColors.info.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Column(
        children: [
          Text(
            label,
            style: const TextStyle(
              fontSize: 12,
              color: AppColors.textSecondary,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '¥${amount.toStringAsFixed(0)}',
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
              color: AppColors.info,
            ),
          ),
        ],
      ),
    );
  }
}
