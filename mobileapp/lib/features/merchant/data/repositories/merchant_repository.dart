import 'package:dio/dio.dart';
import '../models/merchant_model.dart';
import '../../../../core/network/api_client.dart';

/// 商户仓储
class MerchantRepository {
  final Dio _dio;

  MerchantRepository([Dio? dio]) : _dio = dio ?? ApiClient().dio;

  /// 获取商户列表
  Future<PaginatedResponse<Merchant>> getMerchants({
    bool? isDirect,
    String? keyword,
    String? merchantType,
    int page = 1,
    int pageSize = 20,
  }) async {
    final response = await _dio.get('/api/v1/merchants', queryParameters: {
      if (isDirect != null) 'is_direct': isDirect,
      if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      if (merchantType != null && merchantType.isNotEmpty)
        'merchant_type': merchantType,
      'page': page,
      'page_size': pageSize,
    });

    final data = response.data['data'];
    return PaginatedResponse<Merchant>.fromJson(
      data,
      (json) => Merchant.fromJson(json),
    );
  }

  /// 获取商户详情
  Future<MerchantDetail> getMerchantDetail(int id) async {
    final response = await _dio.get('/api/v1/merchants/$id');
    return MerchantDetail.fromJson(response.data['data']);
  }

  /// 获取商户统计
  Future<MerchantStats> getMerchantStats() async {
    final response = await _dio.get('/api/v1/merchants/stats/extended');
    return MerchantStats.fromJson(response.data['data']);
  }

  /// 商户登记
  Future<void> registerMerchant(int id, String phone, String? remark) async {
    await _dio.post('/api/v1/merchants/$id/register', data: {
      'phone': phone,
      if (remark != null && remark.isNotEmpty) 'remark': remark,
    });
  }

  /// 更新商户费率
  Future<void> updateMerchantRate(
    int id,
    double creditRate,
    double debitRate,
  ) async {
    await _dio.put('/api/v1/merchants/$id/rate', data: {
      'credit_rate': creditRate,
      'debit_rate': debitRate,
    });
  }

  /// 更新商户状态
  Future<void> updateMerchantStatus(int id, int status) async {
    await _dio.put('/api/v1/merchants/$id/status', data: {
      'status': status,
    });
  }
}
