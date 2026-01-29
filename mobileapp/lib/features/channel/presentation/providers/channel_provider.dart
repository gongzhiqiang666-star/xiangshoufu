import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/channel_config_model.dart';
import '../../data/services/channel_service.dart';

/// 通道服务Provider
final channelServiceProvider = Provider<ChannelService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return ChannelService(apiClient);
});

/// 通道完整配置Provider（使用family模式，支持传入channelId）
final channelFullConfigProvider = FutureProvider.family<ChannelFullConfig, int>((ref, channelId) async {
  final service = ref.watch(channelServiceProvider);
  return service.getChannelFullConfig(channelId);
});

/// 通道费率配置Provider
final channelRateConfigsProvider = FutureProvider.family<List<ChannelRateConfig>, int>((ref, channelId) async {
  final service = ref.watch(channelServiceProvider);
  return service.getRateConfigs(channelId);
});

/// 通道押金档位Provider
final channelDepositTiersProvider = FutureProvider.family<List<ChannelDepositTier>, int>((ref, channelId) async {
  final service = ref.watch(channelServiceProvider);
  return service.getDepositTiers(channelId);
});

/// 通道流量费返现档位Provider
final channelSimCashbackTiersProvider = FutureProvider.family<List<ChannelSimCashbackTier>, int>((ref, channelId) async {
  final service = ref.watch(channelServiceProvider);
  return service.getSimCashbackTiers(channelId);
});
