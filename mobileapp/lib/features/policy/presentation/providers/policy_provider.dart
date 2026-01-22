import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/models/policy_model.dart';
import '../../data/services/policy_service.dart';

/// 我的政策列表Provider
final myPoliciesProvider = FutureProvider<List<AgentPolicy>>((ref) async {
  final service = ref.watch(policyServiceProvider);
  return service.getMyPolicies();
});

/// 可用通道列表Provider
final availableChannelsProvider = FutureProvider<List<ChannelInfo>>((ref) async {
  final service = ref.watch(policyServiceProvider);
  return service.getAvailableChannels();
});

/// 下级政策参数
class SubordinatePolicyParams {
  final int subordinateId;
  final int channelId;

  SubordinatePolicyParams({
    required this.subordinateId,
    required this.channelId,
  });

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is SubordinatePolicyParams &&
          runtimeType == other.runtimeType &&
          subordinateId == other.subordinateId &&
          channelId == other.channelId;

  @override
  int get hashCode => subordinateId.hashCode ^ channelId.hashCode;
}

/// 下级代理商政策Provider
final subordinatePolicyProvider =
    FutureProvider.family<AgentPolicy, SubordinatePolicyParams>((ref, params) async {
  final service = ref.watch(policyServiceProvider);
  return service.getSubordinatePolicy(params.subordinateId, params.channelId);
});

/// 政策限制Provider
final policyLimitsProvider =
    FutureProvider.family<PolicyLimits, int>((ref, channelId) async {
  final service = ref.watch(policyServiceProvider);
  return service.getPolicyLimits(channelId);
});

/// 选中的通道ID
final selectedChannelIdProvider = StateProvider<int?>((ref) => null);
