import '../../../../core/network/api_client.dart';
import '../models/channel_config_model.dart';

/// 通道配置服务
class ChannelService {
  final ApiClient _apiClient;

  ChannelService(this._apiClient);

  /// 获取通道完整配置（费率+押金+流量费）
  Future<ChannelFullConfig> getChannelFullConfig(int channelId) async {
    final response = await _apiClient.get('/api/v1/admin/channels/$channelId/full-config');
    final data = response.data;

    if (data is Map<String, dynamic>) {
      if (data['code'] == 0 && data['data'] != null) {
        return ChannelFullConfig.fromJson(data['data'] as Map<String, dynamic>);
      }
      throw Exception(data['message'] ?? '获取通道配置失败');
    }
    throw Exception('响应格式错误');
  }

  /// 获取通道费率配置列表
  Future<List<ChannelRateConfig>> getRateConfigs(int channelId) async {
    final response = await _apiClient.get('/api/v1/admin/channels/$channelId/rate-configs');
    final data = response.data;

    if (data is Map<String, dynamic>) {
      if (data['code'] == 0 && data['data'] != null) {
        final list = data['data'] as List;
        return list.map((e) => ChannelRateConfig.fromJson(e as Map<String, dynamic>)).toList();
      }
      throw Exception(data['message'] ?? '获取费率配置失败');
    }
    throw Exception('响应格式错误');
  }

  /// 获取通道押金档位列表
  Future<List<ChannelDepositTier>> getDepositTiers(int channelId) async {
    final response = await _apiClient.get('/api/v1/admin/channels/$channelId/deposit-tiers');
    final data = response.data;

    if (data is Map<String, dynamic>) {
      if (data['code'] == 0 && data['data'] != null) {
        final list = data['data'] as List;
        return list.map((e) => ChannelDepositTier.fromJson(e as Map<String, dynamic>)).toList();
      }
      throw Exception(data['message'] ?? '获取押金档位失败');
    }
    throw Exception('响应格式错误');
  }

  /// 获取通道流量费返现档位列表
  Future<List<ChannelSimCashbackTier>> getSimCashbackTiers(int channelId) async {
    final response = await _apiClient.get('/api/v1/admin/channels/$channelId/sim-cashback-tiers');
    final data = response.data;

    if (data is Map<String, dynamic>) {
      if (data['code'] == 0 && data['data'] != null) {
        final list = data['data'] as List;
        return list.map((e) => ChannelSimCashbackTier.fromJson(e as Map<String, dynamic>)).toList();
      }
      throw Exception(data['message'] ?? '获取流量费返现档位失败');
    }
    throw Exception('响应格式错误');
  }
}
