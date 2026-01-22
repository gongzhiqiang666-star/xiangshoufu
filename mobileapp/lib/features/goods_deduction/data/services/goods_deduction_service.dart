import '../../../../core/network/api_client.dart';
import '../models/goods_deduction_model.dart';

/// 货款代扣服务
class GoodsDeductionService {
  final ApiClient _apiClient;

  GoodsDeductionService(this._apiClient);

  /// 获取我发起的货款代扣列表
  Future<GoodsDeductionListResponse> getSentDeductions({
    int page = 1,
    int pageSize = 10,
    int? status,
    int? deductionSource,
    String? startDate,
    String? endDate,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (status != null) queryParams['status'] = status;
    if (deductionSource != null) queryParams['deduction_source'] = deductionSource;
    if (startDate != null) queryParams['start_date'] = startDate;
    if (endDate != null) queryParams['end_date'] = endDate;

    final response = await _apiClient.get(
      '/api/v1/goods-deduction/sent',
      queryParameters: queryParams,
    );
    return GoodsDeductionListResponse.fromJson(response.data['data']);
  }

  /// 获取我接收的货款代扣列表
  Future<GoodsDeductionListResponse> getReceivedDeductions({
    int page = 1,
    int pageSize = 10,
    int? status,
    int? deductionSource,
    String? startDate,
    String? endDate,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (status != null) queryParams['status'] = status;
    if (deductionSource != null) queryParams['deduction_source'] = deductionSource;
    if (startDate != null) queryParams['start_date'] = startDate;
    if (endDate != null) queryParams['end_date'] = endDate;

    final response = await _apiClient.get(
      '/api/v1/goods-deduction/received',
      queryParameters: queryParams,
    );
    return GoodsDeductionListResponse.fromJson(response.data['data']);
  }

  /// 获取货款代扣详情
  Future<GoodsDeductionDetail> getDeductionDetail(int id) async {
    final response = await _apiClient.get('/api/v1/goods-deduction/$id');
    return GoodsDeductionDetail.fromJson(response.data['data']);
  }

  /// 接收货款代扣
  Future<void> acceptDeduction(int id) async {
    await _apiClient.post('/api/v1/goods-deduction/$id/accept');
  }

  /// 拒绝货款代扣
  Future<void> rejectDeduction(int id, String reason) async {
    await _apiClient.post(
      '/api/v1/goods-deduction/$id/reject',
      data: {'reason': reason},
    );
  }

  /// 获取货款代扣统计汇总
  Future<GoodsDeductionSummary> getSummary({String? type}) async {
    final queryParams = <String, dynamic>{};
    if (type != null) queryParams['type'] = type;

    final response = await _apiClient.get(
      '/api/v1/goods-deduction/summary',
      queryParameters: queryParams,
    );
    return GoodsDeductionSummary.fromJson(response.data['data']);
  }

  /// 获取扣款明细列表
  Future<List<GoodsDeductionDetailRecord>> getDeductionDetails(
    int deductionId, {
    int page = 1,
    int pageSize = 50,
  }) async {
    final response = await _apiClient.get(
      '/api/v1/goods-deduction/$deductionId/details',
      queryParameters: {
        'page': page,
        'page_size': pageSize,
      },
    );
    final list = response.data['data']['list'] as List<dynamic>? ?? [];
    return list.map((e) => GoodsDeductionDetailRecord.fromJson(e)).toList();
  }
}
