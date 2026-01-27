import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../core/theme/app_colors.dart';
import '../../../core/widgets/app_card.dart';
import '../../../core/widgets/empty_state.dart';
import '../../../core/widgets/loading_indicator.dart';
import '../data/models/settlement_price_model.dart';
import 'providers/settlement_price_provider.dart';

/// 调价记录列表页面
/// 支持两种模式：
/// 1. 不传agentId - 显示所有调价记录
/// 2. 传入agentId - 显示指定下级代理商的调价记录
class PriceChangeLogListPage extends ConsumerStatefulWidget {
  final int? agentId;
  final String? agentName;

  const PriceChangeLogListPage({
    super.key,
    this.agentId,
    this.agentName,
  });

  @override
  ConsumerState<PriceChangeLogListPage> createState() => _PriceChangeLogListPageState();
}

class _PriceChangeLogListPageState extends ConsumerState<PriceChangeLogListPage> {
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    _scrollController.addListener(_onScroll);
    // 初次加载数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _refresh();
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollController.position.pixels >= _scrollController.position.maxScrollExtent - 200) {
      _loadMore();
    }
  }

  Future<void> _refresh() async {
    if (widget.agentId != null) {
      await ref.read(agentPriceChangeLogListProvider(widget.agentId!).notifier).refresh();
    } else {
      await ref.read(priceChangeLogListProvider.notifier).refresh();
    }
  }

  Future<void> _loadMore() async {
    if (widget.agentId != null) {
      await ref.read(agentPriceChangeLogListProvider(widget.agentId!).notifier).loadMore();
    } else {
      await ref.read(priceChangeLogListProvider.notifier).loadMore();
    }
  }

  @override
  Widget build(BuildContext context) {
    // 根据是否有agentId选择不同的provider
    final state = widget.agentId != null
        ? ref.watch(agentPriceChangeLogListProvider(widget.agentId!))
        : ref.watch(priceChangeLogListProvider);

    final title = widget.agentName != null ? '${widget.agentName}的调价记录' : '调价记录';

    return Scaffold(
      appBar: AppBar(
        title: Text(title),
      ),
      body: RefreshIndicator(
        onRefresh: _refresh,
        child: _buildContent(state),
      ),
    );
  }

  Widget _buildContent(dynamic state) {
    // 处理两种不同的state类型
    final List<PriceChangeLogModel> list;
    final bool isLoading;
    final String? error;
    final bool hasMore;

    if (state is AgentPriceChangeLogListState) {
      list = state.list;
      isLoading = state.isLoading;
      error = state.error;
      hasMore = state.hasMore;
    } else if (state is PriceChangeLogListState) {
      list = state.list;
      isLoading = state.isLoading;
      error = state.error;
      hasMore = state.hasMore;
    } else {
      return const SizedBox.shrink();
    }

    if (isLoading && list.isEmpty) {
      return const LoadingIndicator();
    }

    if (error != null && list.isEmpty) {
      return EmptyState(
        icon: Icons.error_outline,
        title: '加载失败',
        subtitle: error,
        actionText: '重试',
        onAction: _refresh,
      );
    }

    if (list.isEmpty) {
      return const EmptyState(
        icon: Icons.history,
        title: '暂无调价记录',
        subtitle: '还没有任何调价操作',
      );
    }

    return ListView.builder(
      controller: _scrollController,
      padding: const EdgeInsets.all(16),
      itemCount: list.length + (hasMore ? 1 : 0),
      itemBuilder: (context, index) {
        if (index >= list.length) {
          return const Padding(
            padding: EdgeInsets.symmetric(vertical: 16),
            child: Center(child: CircularProgressIndicator()),
          );
        }
        return _buildChangeLogCard(list[index]);
      },
    );
  }

  Widget _buildChangeLogCard(PriceChangeLogModel item) {
    return AppCard(
      margin: const EdgeInsets.only(bottom: 12),
      onTap: () => _showDetailDialog(item),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 头部：变更类型和配置类型
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  _buildChangeTypeBadge(item.changeType, item.changeTypeName),
                  const SizedBox(width: 8),
                  _buildConfigTypeBadge(item.configType, item.configTypeName),
                ],
              ),
              Text(
                item.source,
                style: TextStyle(
                  fontSize: 12,
                  color: item.source == 'PC' ? Colors.blue : Colors.green,
                ),
              ),
            ],
          ),
          const SizedBox(height: 12),

          // 变更摘要
          Text(
            item.changeSummary,
            style: const TextStyle(fontSize: 14),
            maxLines: 2,
            overflow: TextOverflow.ellipsis,
          ),
          const SizedBox(height: 8),

          // 通道信息
          if (item.channelName.isNotEmpty) ...[
            Text(
              '通道: ${item.channelName}',
              style: TextStyle(
                fontSize: 12,
                color: Colors.grey[600],
              ),
            ),
            const SizedBox(height: 4),
          ],

          // 操作人和时间
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '操作人: ${item.operatorName}',
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey[600],
                ),
              ),
              Text(
                _formatDateTime(item.createdAt),
                style: TextStyle(
                  fontSize: 12,
                  color: Colors.grey[600],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildChangeTypeBadge(int type, String name) {
    Color color;
    switch (type) {
      case 1:
        color = Colors.grey;
        break;
      case 2:
        color = AppColors.primary;
        break;
      case 3:
        color = AppColors.success;
        break;
      case 4:
        color = AppColors.warning;
        break;
      case 5:
        color = AppColors.error;
        break;
      default:
        color = Colors.grey;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        name,
        style: TextStyle(
          fontSize: 12,
          color: color,
        ),
      ),
    );
  }

  Widget _buildConfigTypeBadge(int type, String name) {
    final color = type == 1 ? AppColors.primary : AppColors.warning;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withOpacity(0.1),
        borderRadius: BorderRadius.circular(4),
      ),
      child: Text(
        name,
        style: TextStyle(
          fontSize: 12,
          color: color,
        ),
      ),
    );
  }

  String _formatDateTime(String dateStr) {
    if (dateStr.isEmpty) return '';
    try {
      final date = DateTime.parse(dateStr);
      return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')} '
          '${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
    } catch (e) {
      return dateStr;
    }
  }

  void _showDetailDialog(PriceChangeLogModel item) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(item.changeTypeName),
        content: SingleChildScrollView(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisSize: MainAxisSize.min,
            children: [
              _buildDetailRow('配置类型', item.configTypeName),
              _buildDetailRow('通道', item.channelName.isNotEmpty ? item.channelName : '-'),
              _buildDetailRow('变更字段', item.fieldName),
              _buildDetailRow('变更摘要', item.changeSummary),
              _buildDetailRow('操作人', item.operatorName),
              _buildDetailRow('操作来源', item.source),
              _buildDetailRow('操作时间', _formatDateTime(item.createdAt)),
              if (item.oldValue != null && item.oldValue!.isNotEmpty) ...[
                const SizedBox(height: 12),
                const Text('变更前:', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 4),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Colors.grey[100],
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    item.oldValue!,
                    style: const TextStyle(fontSize: 12, fontFamily: 'monospace'),
                  ),
                ),
              ],
              if (item.newValue != null && item.newValue!.isNotEmpty) ...[
                const SizedBox(height: 12),
                const Text('变更后:', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 4),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: Colors.grey[100],
                    borderRadius: BorderRadius.circular(4),
                  ),
                  child: Text(
                    item.newValue!,
                    style: const TextStyle(fontSize: 12, fontFamily: 'monospace'),
                  ),
                ),
              ],
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(
              '$label:',
              style: TextStyle(
                fontSize: 14,
                color: Colors.grey[600],
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontSize: 14),
            ),
          ),
        ],
      ),
    );
  }
}
