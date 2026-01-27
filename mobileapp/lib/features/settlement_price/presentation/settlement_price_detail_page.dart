import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/loading_indicator.dart';
import '../data/models/settlement_price_model.dart';
import 'providers/settlement_price_provider.dart';

/// 结算价详情页面
class SettlementPriceDetailPage extends ConsumerWidget {
  final int id;

  const SettlementPriceDetailPage({super.key, required this.id});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(settlementPriceDetailProvider(id));

    return Scaffold(
      appBar: AppBar(
        title: const Text('结算价详情'),
      ),
      body: detailAsync.when(
        loading: () => const LoadingIndicator(),
        error: (error, stack) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.error_outline, size: 48, color: Colors.red),
              const SizedBox(height: 16),
              Text('加载失败: $error'),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => ref.refresh(settlementPriceDetailProvider(id)),
                child: const Text('重试'),
              ),
            ],
          ),
        ),
        data: (data) => _buildContent(context, data),
      ),
    );
  }

  Widget _buildContent(BuildContext context, SettlementPriceModel item) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 基础信息卡片
          _buildSection('基础信息', [
            _buildInfoRow('通道', item.channelName.isNotEmpty ? item.channelName : '通道 ${item.channelId}'),
            _buildInfoRow('代理商', item.agentName.isNotEmpty ? item.agentName : '代理商 ${item.agentId}'),
            _buildInfoRow('版本', 'v${item.version}'),
            _buildInfoRow('状态', item.statusName),
          ]),

          const SizedBox(height: 16),

          // 费率配置卡片
          _buildSection('费率配置', [
            if (item.creditRate != null)
              _buildInfoRow('贷记卡费率', '${item.creditRate}%'),
            if (item.debitRate != null)
              _buildInfoRow('借记卡费率', '${item.debitRate}%'),
            if (item.debitCap != null)
              _buildInfoRow('借记卡封顶', '¥${item.debitCap}'),
            if (item.unionpayRate != null)
              _buildInfoRow('云闪付费率', '${item.unionpayRate}%'),
            if (item.wechatRate != null)
              _buildInfoRow('微信费率', '${item.wechatRate}%'),
            if (item.alipayRate != null)
              _buildInfoRow('支付宝费率', '${item.alipayRate}%'),
            if (item.rateConfigs.isNotEmpty)
              ...item.rateConfigs.entries.map((e) =>
                _buildInfoRow(e.key, '${e.value.rate}%'),
              ),
          ]),

          const SizedBox(height: 16),

          // 押金返现配置卡片
          if (item.depositCashbacks.isNotEmpty)
            _buildSection('押金返现配置', [
              ...item.depositCashbacks.map((dc) =>
                _buildInfoRow(
                  '押金 ¥${dc.depositAmountYuan.toStringAsFixed(0)}',
                  '返现 ¥${dc.cashbackAmountYuan.toStringAsFixed(0)}',
                ),
              ),
            ]),

          const SizedBox(height: 16),

          // 流量费返现配置卡片
          _buildSection('流量费返现配置', [
            _buildInfoRow('首次返现', '¥${item.simFirstCashbackYuan.toStringAsFixed(2)}'),
            _buildInfoRow('第2次返现', '¥${item.simSecondCashbackYuan.toStringAsFixed(2)}'),
            _buildInfoRow('第3次+返现', '¥${item.simThirdPlusCashbackYuan.toStringAsFixed(2)}'),
          ]),

          const SizedBox(height: 16),

          // 时间信息
          _buildSection('时间信息', [
            _buildInfoRow('创建时间', _formatDateTime(item.createdAt)),
            _buildInfoRow('更新时间', _formatDateTime(item.updatedAt)),
            if (item.effectiveAt != null)
              _buildInfoRow('生效时间', _formatDateTime(item.effectiveAt!)),
          ]),
        ],
      ),
    );
  }

  Widget _buildSection(String title, List<Widget> children) {
    // 过滤掉空的children
    final validChildren = children.where((child) => child is! SizedBox).toList();
    if (validChildren.isEmpty) {
      return const SizedBox.shrink();
    }

    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            title,
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.bold,
            ),
          ),
          const Divider(),
          ...children,
        ],
      ),
    );
  }

  Widget _buildInfoRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey[600],
            ),
          ),
          Text(
            value,
            style: const TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }

  String _formatDateTime(String dateStr) {
    if (dateStr.isEmpty) return '-';
    try {
      final date = DateTime.parse(dateStr);
      return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')} '
          '${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
    } catch (e) {
      return dateStr;
    }
  }
}
