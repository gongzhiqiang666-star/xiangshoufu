import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../models/policy_model.dart';

/// 政策服务Provider
final policyServiceProvider = Provider<PolicyService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return PolicyService(apiClient);
});

/// 政策服务
class PolicyService {
  final ApiClient _apiClient;

  PolicyService(this._apiClient);

  /// 获取我的政策列表
  Future<List<AgentPolicy>> getMyPolicies() async {
    final response = await _apiClient.get('/api/v1/policies/my');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final List<dynamic> list = apiResponse.data['list'] ?? [];
    return list.map((e) => AgentPolicy.fromJson(e)).toList();
  }

  /// 获取可用通道列表
  Future<List<ChannelInfo>> getAvailableChannels() async {
    final response = await _apiClient.get('/api/v1/channels');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final List<dynamic> list = apiResponse.data['list'] ?? [];
    return list.map((e) => ChannelInfo.fromJson(e)).toList();
  }

  /// 获取下级代理商政策
  Future<AgentPolicy> getSubordinatePolicy(int subordinateId, int channelId) async {
    final response = await _apiClient.get(
      '/api/v1/subordinates/$subordinateId/policy',
      queryParameters: {'channel_id': channelId},
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return AgentPolicy.fromJson(apiResponse.data);
  }

  /// 获取政策限制（我的政策作为限制）
  Future<PolicyLimits> getPolicyLimits(int channelId) async {
    final response = await _apiClient.get(
      '/api/v1/policies/limits',
      queryParameters: {'channel_id': channelId},
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return PolicyLimits.fromJson(apiResponse.data);
  }

  /// 更新下级代理商政策
  Future<void> updateSubordinatePolicy(
    int subordinateId,
    UpdateSubordinatePolicyRequest request,
  ) async {
    final response = await _apiClient.put(
      '/api/v1/subordinates/$subordinateId/policy',
      data: request.toJson(),
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }
}
