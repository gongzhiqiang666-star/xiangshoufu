import '../../../../core/network/api_client.dart';
import '../models/agent_channel_model.dart';

/// 代理商通道API服务
class AgentChannelService {
  final ApiClient _apiClient;

  AgentChannelService(this._apiClient);

  /// 获取代理商通道列表
  Future<List<AgentChannel>> getAgentChannels({int? agentId}) async {
    final queryParams = <String, dynamic>{};
    if (agentId != null) {
      queryParams['agent_id'] = agentId;
    }

    final response = await _apiClient.get(
      '/api/v1/agent-channels',
      queryParameters: queryParams,
    );

    final data = response.data['data'] as List<dynamic>?;
    return data?.map((e) => AgentChannel.fromJson(e as Map<String, dynamic>)).toList() ?? [];
  }

  /// 获取已启用的通道列表（用于APP端显示）
  Future<List<AgentChannel>> getEnabledChannels() async {
    final response = await _apiClient.get('/api/v1/agent-channels/enabled');

    final data = response.data['data'] as List<dynamic>?;
    return data?.map((e) => AgentChannel.fromJson(e as Map<String, dynamic>)).toList() ?? [];
  }

  /// 获取代理商通道统计
  Future<AgentChannelStats> getAgentChannelStats({int? agentId}) async {
    final queryParams = <String, dynamic>{};
    if (agentId != null) {
      queryParams['agent_id'] = agentId;
    }

    final response = await _apiClient.get(
      '/api/v1/agent-channels/stats',
      queryParameters: queryParams,
    );

    return AgentChannelStats.fromJson(response.data['data']);
  }

  /// 启用通道
  Future<void> enableChannel(int agentId, int channelId) async {
    await _apiClient.post(
      '/api/v1/agent-channels/enable',
      data: {
        'agent_id': agentId,
        'channel_id': channelId,
      },
    );
  }

  /// 禁用通道
  Future<void> disableChannel(int agentId, int channelId) async {
    await _apiClient.post(
      '/api/v1/agent-channels/disable',
      data: {
        'agent_id': agentId,
        'channel_id': channelId,
      },
    );
  }

  /// 设置通道可见性
  Future<void> setChannelVisibility(int agentId, int channelId, bool isVisible) async {
    await _apiClient.post(
      '/api/v1/agent-channels/visibility',
      data: {
        'agent_id': agentId,
        'channel_id': channelId,
        'is_visible': isVisible,
      },
    );
  }

  /// 批量启用通道
  Future<void> batchEnableChannels(int agentId, List<int> channelIds) async {
    await _apiClient.post(
      '/api/v1/agent-channels/batch-enable',
      data: {
        'agent_id': agentId,
        'channel_ids': channelIds,
      },
    );
  }

  /// 批量禁用通道
  Future<void> batchDisableChannels(int agentId, List<int> channelIds) async {
    await _apiClient.post(
      '/api/v1/agent-channels/batch-disable',
      data: {
        'agent_id': agentId,
        'channel_ids': channelIds,
      },
    );
  }

  /// 初始化代理商通道配置
  Future<void> initAgentChannels(int agentId) async {
    await _apiClient.post(
      '/api/v1/agent-channels/init',
      data: {
        'agent_id': agentId,
      },
    );
  }
}
