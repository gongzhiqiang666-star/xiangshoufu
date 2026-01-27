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

  /// 获取我接收的代扣列表
  Future<DeductionPlanListResponse> getReceivedDeductions({
    int page = 1,
    int pageSize = 10,
    String? status, // 支持逗号分隔多状态，如 "0,1"
    int? planType,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (status != null) queryParams['status'] = status;
    if (planType != null) queryParams['plan_type'] = planType;

    final response = await _apiClient.get(
      '/api/v1/deduction/received',
      queryParameters: queryParams,
    );
    return DeductionPlanListResponse.fromJson(response.data['data']);
  }

  /// 获取我发起的代扣列表
  Future<DeductionPlanListResponse> getSentDeductions({
    int page = 1,
    int pageSize = 10,
    String? status, // 支持逗号分隔多状态
    int? planType,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (status != null) queryParams['status'] = status;
    if (planType != null) queryParams['plan_type'] = planType;

    final response = await _apiClient.get(
      '/api/v1/deduction/sent',
      queryParameters: queryParams,
    );
    return DeductionPlanListResponse.fromJson(response.data['data']);
  }

  /// 获取代扣摘要（我接收的/我发起的统计）
  Future<DeductionSummary> getDeductionSummary() async {
    final response = await _apiClient.get('/api/v1/deduction/summary');
    return DeductionSummary.fromJson(response.data['data']);
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

  /// 创建需要接收确认的代扣计划
  Future<DeductionPlan> createPlanWithAccept({
    required int deducteeId,
    required int planType,
    required int totalAmount,
    required int totalPeriods,
    int deductionSource = 3, // 1=分润 2=服务费 3=两者
    String? remark,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/deduction/plans/with-accept',
      data: {
        'deductee_id': deducteeId,
        'plan_type': planType,
        'total_amount': totalAmount,
        'total_periods': totalPeriods,
        'deduction_source': deductionSource,
        if (remark != null) 'remark': remark,
      },
    );
    return DeductionPlan.fromJson(response.data['data']);
  }

  /// 接收确认代扣计划
  Future<void> acceptPlan(int id) async {
    await _apiClient.post('/api/v1/deduction/plans/$id/accept');
  }

  /// 拒绝代扣计划
  Future<void> rejectPlan(int id, {String? reason}) async {
    await _apiClient.post(
      '/api/v1/deduction/plans/$id/reject',
      data: reason != null ? {'reason': reason} : null,
    );
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
