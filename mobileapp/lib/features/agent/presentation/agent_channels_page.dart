import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'providers/agent_channel_provider.dart';
import '../data/models/agent_channel_model.dart';

/// 代理商通道管理页面
class AgentChannelsPage extends ConsumerWidget {
  final int? agentId; // 如果为空，则查看当前登录代理商的通道

  const AgentChannelsPage({
    super.key,
    this.agentId,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = agentId != null
        ? ref.watch(agentChannelsProvider(agentId!))
        : ref.watch(myAgentChannelsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('通道管理'),
        actions: [
          IconButton(
            icon: const Icon(Icons.refresh),
            onPressed: () {
              if (agentId != null) {
                ref.read(agentChannelsProvider(agentId!).notifier).loadChannels();
              } else {
                ref.read(myAgentChannelsProvider.notifier).loadChannels();
              }
            },
          ),
        ],
      ),
      body: _buildBody(context, ref, state),
    );
  }

  Widget _buildBody(BuildContext context, WidgetRef ref, AgentChannelsState state) {
    if (state.isLoading && state.channels.isEmpty) {
      return const Center(child: CircularProgressIndicator());
    }

    if (state.error != null && state.channels.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(
              '加载失败',
              style: Theme.of(context).textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              state.error!,
              style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey),
            ),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () {
                if (agentId != null) {
                  ref.read(agentChannelsProvider(agentId!).notifier).loadChannels();
                } else {
                  ref.read(myAgentChannelsProvider.notifier).loadChannels();
                }
              },
              child: const Text('重试'),
            ),
          ],
        ),
      );
    }

    if (state.channels.isEmpty) {
      return _buildEmptyState(context, ref);
    }

    return RefreshIndicator(
      onRefresh: () async {
        if (agentId != null) {
          await ref.read(agentChannelsProvider(agentId!).notifier).loadChannels();
        } else {
          await ref.read(myAgentChannelsProvider.notifier).loadChannels();
        }
      },
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: state.channels.length,
        itemBuilder: (context, index) {
          final channel = state.channels[index];
          return _buildChannelCard(context, ref, channel);
        },
      ),
    );
  }

  Widget _buildEmptyState(BuildContext context, WidgetRef ref) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.settings_ethernet_outlined,
            size: 64,
            color: Colors.grey[400],
          ),
          const SizedBox(height: 16),
          Text(
            '暂无通道配置',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(color: Colors.grey),
          ),
          const SizedBox(height: 8),
          Text(
            '请联系上级初始化通道配置',
            style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey),
          ),
        ],
      ),
    );
  }

  Widget _buildChannelCard(BuildContext context, WidgetRef ref, AgentChannel channel) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        channel.channelName,
                        style: Theme.of(context).textTheme.titleMedium?.copyWith(
                              fontWeight: FontWeight.bold,
                            ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '编码: ${channel.channelCode}',
                        style: Theme.of(context).textTheme.bodySmall?.copyWith(
                              color: Colors.grey,
                            ),
                      ),
                    ],
                  ),
                ),
                _buildStatusChip(context, channel),
              ],
            ),
            const SizedBox(height: 12),
            Row(
              children: [
                _buildInfoItem(
                  context,
                  '启用状态',
                  channel.isEnabled ? '已启用' : '已禁用',
                  channel.isEnabled ? Colors.green : Colors.red,
                ),
                const SizedBox(width: 24),
                _buildInfoItem(
                  context,
                  'APP可见',
                  channel.isVisible ? '可见' : '隐藏',
                  channel.isVisible ? Colors.blue : Colors.grey,
                ),
              ],
            ),
            if (channel.enabledAt != null) ...[
              const SizedBox(height: 8),
              Text(
                '启用时间: ${_formatDate(channel.enabledAt!)}',
                style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildStatusChip(BuildContext context, AgentChannel channel) {
    if (!channel.isEnabled) {
      return Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        decoration: BoxDecoration(
          color: Colors.red.withValues(alpha: 0.1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Text(
          '已禁用',
          style: TextStyle(
            fontSize: 12,
            color: Colors.red,
            fontWeight: FontWeight.bold,
          ),
        ),
      );
    }

    if (!channel.isVisible) {
      return Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
        decoration: BoxDecoration(
          color: Colors.orange.withValues(alpha: 0.1),
          borderRadius: BorderRadius.circular(12),
        ),
        child: const Text(
          '隐藏中',
          style: TextStyle(
            fontSize: 12,
            color: Colors.orange,
            fontWeight: FontWeight.bold,
          ),
        ),
      );
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: Colors.green.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: const Text(
        '正常',
        style: TextStyle(
          fontSize: 12,
          color: Colors.green,
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }

  Widget _buildInfoItem(BuildContext context, String label, String value, Color color) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.bodySmall?.copyWith(color: Colors.grey),
        ),
        const SizedBox(height: 2),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: color,
                fontWeight: FontWeight.w500,
              ),
        ),
      ],
    );
  }

  String _formatDate(String dateStr) {
    try {
      final date = DateTime.parse(dateStr);
      return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
    } catch (e) {
      return dateStr;
    }
  }
}
