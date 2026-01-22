import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/agent_channel_model.dart';
import '../../data/services/agent_channel_service.dart';

/// 代理商通道服务Provider
final agentChannelServiceProvider = Provider<AgentChannelService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return AgentChannelService(apiClient);
});

/// 代理商通道列表状态
class AgentChannelsState {
  final List<AgentChannel> channels;
  final bool isLoading;
  final String? error;

  AgentChannelsState({
    this.channels = const [],
    this.isLoading = false,
    this.error,
  });

  AgentChannelsState copyWith({
    List<AgentChannel>? channels,
    bool? isLoading,
    String? error,
  }) {
    return AgentChannelsState(
      channels: channels ?? this.channels,
      isLoading: isLoading ?? this.isLoading,
      error: error,
    );
  }
}

/// 代理商通道列表Notifier
class AgentChannelsNotifier extends StateNotifier<AgentChannelsState> {
  final AgentChannelService _service;
  final int? _agentId;

  AgentChannelsNotifier(this._service, this._agentId) : super(AgentChannelsState()) {
    loadChannels();
  }

  Future<void> loadChannels() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final channels = await _service.getAgentChannels(agentId: _agentId);
      state = state.copyWith(channels: channels, isLoading: false);
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> toggleChannelEnabled(AgentChannel channel) async {
    try {
      if (channel.isEnabled) {
        await _service.disableChannel(_agentId ?? channel.agentId, channel.channelId);
      } else {
        await _service.enableChannel(_agentId ?? channel.agentId, channel.channelId);
      }
      await loadChannels();
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  Future<void> toggleChannelVisibility(AgentChannel channel) async {
    try {
      await _service.setChannelVisibility(
        _agentId ?? channel.agentId,
        channel.channelId,
        !channel.isVisible,
      );
      await loadChannels();
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  Future<void> initChannels(int agentId) async {
    try {
      await _service.initAgentChannels(agentId);
      await loadChannels();
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }
}

/// 当前代理商通道列表Provider
final myAgentChannelsProvider = StateNotifierProvider<AgentChannelsNotifier, AgentChannelsState>((ref) {
  final service = ref.watch(agentChannelServiceProvider);
  return AgentChannelsNotifier(service, null);
});

/// 指定代理商通道列表Provider (用于查看下级代理商)
final agentChannelsProvider = StateNotifierProvider.family<AgentChannelsNotifier, AgentChannelsState, int>((ref, agentId) {
  final service = ref.watch(agentChannelServiceProvider);
  return AgentChannelsNotifier(service, agentId);
});

/// 已启用通道列表Provider
final enabledChannelsProvider = FutureProvider<List<AgentChannel>>((ref) async {
  final service = ref.watch(agentChannelServiceProvider);
  return service.getEnabledChannels();
});

/// 通道统计Provider
final agentChannelStatsProvider = FutureProvider.family<AgentChannelStats, int?>((ref, agentId) async {
  final service = ref.watch(agentChannelServiceProvider);
  return service.getAgentChannelStats(agentId: agentId);
});
