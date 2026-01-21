import 'package:dio/dio.dart';
import '../../../../core/network/api_client.dart';
import '../models/agent_model.dart';

/// 代理商API服务
class AgentService {
  final ApiClient _apiClient;

  AgentService(this._apiClient);

  /// 获取当前代理商详情
  Future<AgentDetail> getMyProfile() async {
    final response = await _apiClient.get('/api/v1/agents/detail');
    return AgentDetail.fromJson(response.data['data']);
  }

  /// 获取代理商详情
  Future<AgentDetail> getAgentDetail(int agentId) async {
    final response = await _apiClient.get('/api/v1/agents/$agentId');
    return AgentDetail.fromJson(response.data['data']);
  }

  /// 获取邀请码信息
  Future<InviteCodeInfo> getInviteCode() async {
    final response = await _apiClient.get('/api/v1/agents/invite-code');
    return InviteCodeInfo.fromJson(response.data['data']);
  }

  /// 获取团队统计
  Future<TeamStats> getTeamStats() async {
    final response = await _apiClient.get('/api/v1/agents/stats');
    final data = response.data['data'];
    return TeamStats(
      directAgentCount: data['direct_agent_count'] ?? 0,
      teamAgentCount: data['team_agent_count'] ?? 0,
      directMerchantCount: data['direct_merchant_count'] ?? 0,
      teamMerchantCount: data['team_merchant_count'] ?? 0,
      todayNewAgents: data['today_new_agents'] ?? 0,
      monthNewAgents: data['month_new_agents'] ?? 0,
    );
  }

  /// 获取直属下级代理商列表
  Future<SubordinateListResponse> getSubordinates({
    int page = 1,
    int pageSize = 10,
    String? keyword,
    int? status,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (keyword != null && keyword.isNotEmpty) {
      queryParams['keyword'] = keyword;
    }
    if (status != null) {
      queryParams['status'] = status;
    }

    final response = await _apiClient.get(
      '/api/v1/agents/subordinates',
      queryParameters: queryParams,
    );

    final data = response.data['data'];
    final list = (data['list'] as List<dynamic>?)
            ?.map((e) => AgentInfo.fromJson(e as Map<String, dynamic>))
            .toList() ??
        [];

    return SubordinateListResponse(
      list: list,
      total: data['total'] ?? 0,
    );
  }

  /// 创建代理商
  Future<AgentDetail> createAgent(CreateAgentRequest request) async {
    final response = await _apiClient.post(
      '/api/v1/agents',
      data: request.toJson(),
    );
    return AgentDetail.fromJson(response.data['data']);
  }

  /// 搜索代理商
  Future<List<AgentInfo>> searchAgents(String keyword) async {
    final response = await _apiClient.get(
      '/api/v1/agents/search',
      queryParameters: {'keyword': keyword},
    );

    final data = response.data['data'];
    final list = (data['list'] as List<dynamic>?)
            ?.map((e) => AgentInfo.fromJson(e as Map<String, dynamic>))
            .toList() ??
        [];

    return list;
  }

  /// 更新代理商资料
  Future<void> updateProfile({
    String? agentName,
    String? contactName,
    String? contactPhone,
    String? bankName,
    String? bankAccount,
    String? bankCardNo,
  }) async {
    final data = <String, dynamic>{};
    if (agentName != null) data['agent_name'] = agentName;
    if (contactName != null) data['contact_name'] = contactName;
    if (contactPhone != null) data['contact_phone'] = contactPhone;
    if (bankName != null) data['bank_name'] = bankName;
    if (bankAccount != null) data['bank_account'] = bankAccount;
    if (bankCardNo != null) data['bank_card_no'] = bankCardNo;

    await _apiClient.put('/api/v1/agents/profile', data: data);
  }
}
