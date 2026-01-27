import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/empty_state.dart';
import '../../../core/widgets/loading_indicator.dart';
import '../data/models/settlement_price_model.dart';
import 'providers/settlement_price_provider.dart';

/// 结算价列表页面
class SettlementPriceListPage extends ConsumerStatefulWidget {
  const SettlementPriceListPage({super.key});

  @override
  ConsumerState<SettlementPriceListPage> createState() => _SettlementPriceListPageState();
}

class _SettlementPriceListPageState extends ConsumerState<SettlementPriceListPage> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    // 初次加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      ref.read(settlementPriceListProvider.notifier).refresh();
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent - 200) {
      ref.read(settlementPriceListProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(settlementPriceListProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('结算价管理'),
        actions: [
          IconButton(
            icon: const Icon(Icons.history),
            onPressed: () {
              Navigator.pushNamed(context, '/price-change-logs');
            },
            tooltip: '调价记录',
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: () => ref.read(settlementPriceListProvider.notifier).refresh(),
        child: _buildContent(state),
      ),
    );
  }

  Widget _buildContent(SettlementPriceListState state) {
    if (state.isLoading && state.list.isEmpty) {
      return const LoadingIndicator();
    }

    if (state.error != null && state.list.isEmpty) {
      return EmptyState(
        icon: Icons.error_outline,
        title: '加载失败',
        subtitle: state.error,
        actionText: '重试',
        onAction: () => ref.read(settlementPriceListProvider.notifier).refresh(),
      );
    }

    if (state.list.isEmpty) {
      return const EmptyState(
        icon: Icons.price_change_outlined,
        title: '暂无结算价数据',
        subtitle: '还没有配置结算价',
      );
    }

    return ListView.builder(
      controller: _scrollController,
      padding: const EdgeInsets.all(16),
      itemCount: state.list.length + (state.hasMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index >= state.list.length) {
          return const Padding(
            padding: EdgeInsets.symmetric(vertical: 16),
            child: Center(child: CircularProgressIndicator()),
          );
        }
        return _buildSettlementPriceCard(state.list[index]);
      },
    );
  }

  Widget _buildSettlementPriceCard(SettlementPriceModel item) {
    return AppCard(
      margin: const EdgeInsets.only(bottom: 12),
      onTap: () {
        Navigator.pushNamed(context, '/settlement-prices/${item.id}');
      },
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 头部：通道名称和状态
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                item.channelName.isNotEmpty ? item.channelName : '通道 ${item.channelId}',
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                ),
              ),
              _buildStatusBadge(item.status),
            ],
          ),
          const SizedBox(height: 12),

          // 费率信息
          if (item.creditRate != null || item.debitRate != null) ...[
            Row(
              children: [
                if (item.creditRate != null) ...[
                  _buildRateChip('贷记卡', '${item.creditRate}%'),
                  const SizedBox(width: 8),
                ],
                if (item.debitRate != null)
                  _buildRateChip('借记卡', '${item.debitRate}%'),
              ],
            ),
            const SizedBox(height: 8),
          ],

          // 押金返现
          if (item.depositCashbacks.isNotEmpty) ...[
            Wrap(
              spacing: 8,
              runSpacing: 4,
              children: item.depositCashbacks.map((dc) {
                return _buildCashbackChip(
                  '¥${dc.depositAmountYuan.toStringAsFixed(0)}',
                  '返¥${dc.cashbackAmountYuan.toStringAsFixed(0)}',
                );
              }).toList(),
            ),
            const SizedBox(height: 8),
          ],

          // 流量费返现
          if (item.simFirstCashback > 0 || item.simSecondCashback > 0) ...[
            Row(
              children: [
                _buildSimChip('首次', '¥${item.simFirstCashbackYuan.toStringAsFixed(0)}'),
                const SizedBox(width: 8),
                _buildSimChip('第2次', '¥${item.simSecondCashbackYuan.toStringAsFixed(0)}'),
                const SizedBox(width: 8),
                _buildSimChip('第3次+', '¥${item.simThirdPlusCashbackYuan.toStringAsFixed(0)}'),
              ],
            ),
          ],

          const SizedBox(height: 8),
          // 版本信息
          Text(
            '版本 v${item.version}',
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey[600],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildStatusBadge(int status) {
    final isEnabled = status == 1;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: isEnabled ? AppColors.success.withOpacity(0.1) : AppColors.error.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        isEnabled ? '启用' : '禁用',
        style: TextStyle(
          fontSize: 12,
          color: isEnabled ? AppColors.success : AppColors.error,
        ),
      ),
    );
  }

  Widget _buildRateChip(String label, String value) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AppColors.primary.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        '$label: $value',
        style: TextStyle(
          fontSize: 12,
          color: AppColors.primary,
        ),
      ),
    );
  }

  Widget _buildCashbackChip(String deposit, String cashback) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: AppColors.success.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        '$deposit → $cashback',
        style: TextStyle(
          fontSize: 12,
          color: AppColors.success,
        ),
      ),
    );
  }

  Widget _buildSimChip(String label, String value) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(
        color: AppColors.warning.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        '$label$value',
        style: TextStyle(
          fontSize: 11,
          color: AppColors.warning,
        ),
      ),
    );
  }
}
