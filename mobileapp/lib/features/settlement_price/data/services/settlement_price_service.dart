import 'package:dio/dio.dart';
import '../../../../core/network/api_client.dart';
import '../models/settlement_price_model.dart';

/// 结算价服务
class SettlementPriceService {
  final ApiClient _apiClient;

  SettlementPriceService(this._apiClient);

  /// 获取结算价列表
  Future<SettlementPriceListResponse> getSettlementPrices({
    int? agentId,
    int? channelId,
    int? status,
    int page = 1,
    int pageSize = 20,
  }) async {
    final params = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (agentId != null) params['agent_id'] = agentId;
    if (channelId != null) params['channel_id'] = channelId;
    if (status != null) params['status'] = status;

    final response = await _apiClient.get('/api/v1/settlement-prices', queryParameters: params);
    return SettlementPriceListResponse.fromJson(response.data);
  }

  /// 获取结算价详情
  Future<SettlementPriceModel> getSettlementPrice(int id) async {
    final response = await _apiClient.get('/api/v1/settlement-prices/$id');
    return SettlementPriceModel.fromJson(response.data);
  }

  /// 更新费率
  Future<SettlementPriceModel> updateRate(int id, Map<String, dynamic> data) async {
    final response = await _apiClient.put('/api/v1/settlement-prices/$id/rate', data: data);
    return SettlementPriceModel.fromJson(response.data);
  }

  /// 更新押金返现
  Future<SettlementPriceModel> updateDepositCashback(int id, Map<String, dynamic> data) async {
    final response = await _apiClient.put('/api/v1/settlement-prices/$id/deposit', data: data);
    return SettlementPriceModel.fromJson(response.data);
  }

  /// 更新流量费返现
  Future<SettlementPriceModel> updateSimCashback(int id, Map<String, dynamic> data) async {
    final response = await _apiClient.put('/api/v1/settlement-prices/$id/sim', data: data);
    return SettlementPriceModel.fromJson(response.data);
  }

  /// 获取调价记录列表
  Future<PriceChangeLogListResponse> getPriceChangeLogs({
    int? agentId,
    int? channelId,
    int? changeType,
    String? startDate,
    String? endDate,
    int page = 1,
    int pageSize = 20,
  }) async {
    final params = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (agentId != null) params['agent_id'] = agentId;
    if (channelId != null) params['channel_id'] = channelId;
    if (changeType != null) params['change_type'] = changeType;
    if (startDate != null) params['start_date'] = startDate;
    if (endDate != null) params['end_date'] = endDate;

    final response = await _apiClient.get('/api/v1/price-change-logs', queryParameters: params);
    return PriceChangeLogListResponse.fromJson(response.data);
  }

  /// 获取调价记录详情
  Future<PriceChangeLogModel> getPriceChangeLog(int id) async {
    final response = await _apiClient.get('/api/v1/price-change-logs/$id');
    return PriceChangeLogModel.fromJson(response.data);
  }
}
