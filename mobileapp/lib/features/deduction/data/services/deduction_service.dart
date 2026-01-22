import '../../../../core/network/api_client.dart';
import '../models/deduction_model.dart';

/// 代扣管理服务
class DeductionService {
  final ApiClient _apiClient;

  DeductionService(this._apiClient);

  /// 获取代扣计划列表
  Future<DeductionPlanListResponse> getDeductionPlans({
    int page = 1,
    int pageSize = 10,
    int? planType,
    int? status,
    int? deducteeId,
    String? startDate,
    String? endDate,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (planType != null) queryParams['plan_type'] = planType;
    if (status != null) queryParams['status'] = status;
    if (deducteeId != null) queryParams['deductee_id'] = deducteeId;
    if (startDate != null) queryParams['start_date'] = startDate;
    if (endDate != null) queryParams['end_date'] = endDate;

    final response = await _apiClient.get(
      '/api/v1/deduction/plans',
      queryParameters: queryParams,
    );
    return DeductionPlanListResponse.fromJson(response.data['data']);
  }

  /// 获取代扣计划详情
  Future<DeductionPlanDetail> getPlanDetail(int id) async {
    final response = await _apiClient.get('/api/v1/deduction/plans/$id');
    return DeductionPlanDetail.fromJson(response.data['data']);
  }

  /// 创建代扣计划
  Future<DeductionPlan> createPlan({
    required int deducteeId,
    required int planType,
    required int totalAmount,
    required int totalPeriods,
    String? remark,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/deduction/plans',
      data: {
        'deductee_id': deducteeId,
        'plan_type': planType,
        'total_amount': totalAmount,
        'total_periods': totalPeriods,
        if (remark != null) 'remark': remark,
      },
    );
    return DeductionPlan.fromJson(response.data['data']);
  }

  /// 暂停代扣计划
  Future<void> pausePlan(int id) async {
    await _apiClient.post('/api/v1/deduction/plans/$id/pause');
  }

  /// 恢复代扣计划
  Future<void> resumePlan(int id) async {
    await _apiClient.post('/api/v1/deduction/plans/$id/resume');
  }

  /// 取消代扣计划
  Future<void> cancelPlan(int id) async {
    await _apiClient.post('/api/v1/deduction/plans/$id/cancel');
  }

  /// 获取代扣记录列表
  Future<List<DeductionRecord>> getRecords(int planId, {
    int page = 1,
    int pageSize = 50,
  }) async {
    final response = await _apiClient.get(
      '/api/v1/deduction/plans/$planId/records',
      queryParameters: {
        'page': page,
        'page_size': pageSize,
      },
    );
    final list = response.data['data']['list'] as List<dynamic>? ?? [];
    return list.map((e) => DeductionRecord.fromJson(e)).toList();
  }

  /// 获取代扣计划统计
  Future<DeductionPlanStats> getStats({
    int? planType,
    String? startDate,
    String? endDate,
  }) async {
    final queryParams = <String, dynamic>{};
    if (planType != null) queryParams['plan_type'] = planType;
    if (startDate != null) queryParams['start_date'] = startDate;
    if (endDate != null) queryParams['end_date'] = endDate;

    final response = await _apiClient.get(
      '/api/v1/deduction/plans/stats',
      queryParameters: queryParams,
    );
    return DeductionPlanStats.fromJson(response.data['data']);
  }
}
