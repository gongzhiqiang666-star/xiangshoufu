import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../models/wallet_model.dart';

/// Wallet服务Provider
final walletServiceProvider = Provider<WalletService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return WalletService(apiClient);
});

/// 钱包服务
class WalletService {
  final ApiClient _apiClient;

  WalletService(this._apiClient);

  // ========== 基础钱包 ==========

  /// 获取钱包汇总
  Future<WalletSummaryModel> getWalletSummary() async {
    final response = await _apiClient.get('/api/v1/wallets/summary');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => WalletSummaryModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 获取钱包列表
  Future<List<WalletModel>> getWallets() async {
    final response = await _apiClient.get('/api/v1/wallets');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final List<dynamic> list = apiResponse.data['list'] ?? [];
    return list.map((e) => WalletModel.fromJson(e)).toList();
  }

  /// 获取钱包流水
  Future<PaginatedResponse<WalletLogModel>> getWalletLogs({
    int? walletId,
    String? type,
    String? startDate,
    String? endDate,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (walletId != null) queryParams['wallet_id'] = walletId;
    if (type != null) queryParams['type'] = type;
    if (startDate != null) queryParams['start_date'] = startDate;
    if (endDate != null) queryParams['end_date'] = endDate;

    final response = await _apiClient.get(
      '/api/v1/wallets/logs',
      queryParameters: queryParams,
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return PaginatedResponse.fromJson(
      apiResponse.data,
      (json) => WalletLogModel.fromJson(json),
    );
  }

  /// 申请提现
  Future<void> applyWithdraw({
    required int walletId,
    required int amount,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/wallets/withdraw',
      data: {
        'wallet_id': walletId,
        'amount': amount,
      },
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }

  // ========== 充值钱包 ==========

  /// 获取钱包配置
  Future<AgentWalletConfigModel> getMyWalletConfig() async {
    final response = await _apiClient.get('/api/v1/charging-wallet/config');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => AgentWalletConfigModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 获取充值钱包汇总
  Future<ChargingWalletSummaryModel> getChargingWalletSummary() async {
    final response = await _apiClient.get('/api/v1/charging-wallet/summary');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => ChargingWalletSummaryModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 申请充值
  Future<String> createChargingDeposit({
    required int amount,
    required int paymentMethod,
    String? paymentRef,
    String? remark,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/charging-wallet/deposits',
      data: {
        'amount': amount,
        'payment_method': paymentMethod,
        if (paymentRef != null) 'payment_ref': paymentRef,
        if (remark != null) 'remark': remark,
      },
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data['deposit_no'] ?? '';
  }

  /// 发放奖励
  Future<String> issueChargingReward({
    required int toAgentId,
    required int amount,
    String? remark,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/charging-wallet/rewards',
      data: {
        'to_agent_id': toAgentId,
        'amount': amount,
        if (remark != null) 'remark': remark,
      },
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data['reward_no'] ?? '';
  }

  // ========== 沉淀钱包 ==========

  /// 获取沉淀钱包汇总
  Future<SettlementWalletSummaryModel> getSettlementWalletSummary() async {
    final response = await _apiClient.get('/api/v1/settlement-wallet/summary');
    final apiResponse = ApiResponse.fromJson(
      response.data,
      (data) => SettlementWalletSummaryModel.fromJson(data),
    );
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data!;
  }

  /// 获取下级余额明细
  Future<List<SubordinateBalanceModel>> getSubordinateBalances() async {
    final response = await _apiClient.get('/api/v1/settlement-wallet/subordinates');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    final List<dynamic> list = apiResponse.data['list'] ?? [];
    return list.map((e) => SubordinateBalanceModel.fromJson(e)).toList();
  }

  /// 使用沉淀款
  Future<String> useSettlement({
    required int amount,
    String? remark,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/settlement-wallet/use',
      data: {
        'amount': amount,
        if (remark != null) 'remark': remark,
      },
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data['usage_no'] ?? '';
  }

  /// 归还沉淀款
  Future<String> returnSettlement({
    required int amount,
    String? remark,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/settlement-wallet/return',
      data: {
        'amount': amount,
        if (remark != null) 'remark': remark,
      },
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data['usage_no'] ?? '';
  }

  /// 获取使用记录列表
  Future<PaginatedResponse<SettlementUsageModel>> getSettlementUsageList({
    int? usageType,
    int page = 1,
    int pageSize = 20,
  }) async {
    final queryParams = <String, dynamic>{
      'page': page,
      'page_size': pageSize,
    };
    if (usageType != null) queryParams['usage_type'] = usageType;

    final response = await _apiClient.get(
      '/api/v1/settlement-wallet/usages',
      queryParameters: queryParams,
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return PaginatedResponse.fromJson(
      apiResponse.data,
      (json) => SettlementUsageModel.fromJson(json),
    );
  }

  // ========== 钱包拆分相关 ==========

  /// 获取钱包列表（支持拆分模式）
  Future<WalletListWithSplitResponse> getWalletsWithSplit() async {
    final response = await _apiClient.get('/api/v1/wallets/with-split');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return WalletListWithSplitResponse.fromJson(apiResponse.data);
  }

  /// 检查是否按通道拆分
  Future<bool> checkSplitStatus() async {
    final response = await _apiClient.get('/api/v1/wallets/split-status');
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
    return apiResponse.data['split_by_channel'] ?? false;
  }

  /// 申请提现（支持拆分模式）
  Future<void> applyWithdrawWithChannel({
    required int walletId,
    required int amount,
    int? channelId,
  }) async {
    final response = await _apiClient.post(
      '/api/v1/wallets/withdraw',
      data: {
        'wallet_id': walletId,
        'amount': amount,
        if (channelId != null) 'channel_id': channelId,
      },
    );
    final apiResponse = ApiResponse.fromJson(response.data, null);
    if (!apiResponse.isSuccess) {
      throw ApiException(apiResponse.code, apiResponse.message);
    }
  }
}
