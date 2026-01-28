import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import '../../../../core/network/api_client.dart';
import '../../domain/models/terminal.dart';

/// 终端服务
class TerminalService {
  final Dio _dio = ApiClient().dio;

  /// 获取终端列表（支持多条件筛选）
  Future<PaginatedResponse<Terminal>> getTerminals({
    int? status,
    int? channelId,
    String? brandCode,
    String? modelCode,
    String? statusGroup,
    String? keyword,
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final response = await _dio.get(
        '/api/v1/terminals',
        queryParameters: {
          if (status != null) 'status': status,
          if (channelId != null) 'channel_id': channelId,
          if (brandCode != null && brandCode.isNotEmpty) 'brand_code': brandCode,
          if (modelCode != null && modelCode.isNotEmpty) 'model_code': modelCode,
          if (statusGroup != null && statusGroup.isNotEmpty) 'status_group': statusGroup,
          if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
          'page': page,
          'page_size': pageSize,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return PaginatedResponse.fromJson(
        apiResponse.data as Map<String, dynamic>,
        (json) => Terminal.fromJson(json),
      );
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 获取筛选选项
  Future<TerminalFilterOptions> getFilterOptions({
    int? channelId,
    String? brandCode,
    String? modelCode,
  }) async {
    try {
      final response = await _dio.get(
        '/api/v1/terminals/filter-options',
        queryParameters: {
          if (channelId != null) 'channel_id': channelId,
          if (brandCode != null && brandCode.isNotEmpty) 'brand_code': brandCode,
          if (modelCode != null && modelCode.isNotEmpty) 'model_code': modelCode,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return TerminalFilterOptions.fromJson(apiResponse.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 获取终端流动记录
  Future<TerminalFlowLogsResponse> getFlowLogs({
    required String terminalSn,
    String logType = 'all',
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final response = await _dio.get(
        '/api/v1/terminals/$terminalSn/flow-logs',
        queryParameters: {
          'log_type': logType,
          'page': page,
          'page_size': pageSize,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return TerminalFlowLogsResponse.fromJson(apiResponse.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 获取终端详情
  Future<Terminal> getTerminalDetail(String sn) async {
    try {
      final response = await _dio.get('/api/v1/terminals/$sn');

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return Terminal.fromJson(apiResponse.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 获取终端统计
  Future<TerminalStats> getTerminalStats() async {
    try {
      final response = await _dio.get('/api/v1/terminals/stats');

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return TerminalStats.fromJson(apiResponse.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 终端下发
  Future<TerminalDistribute> distributeTerminal({
    required int toAgentId,
    required String terminalSn,
    required int channelId,
    required int goodsPrice,
    required int deductionType,
    int? deductionPeriods,
    String? remark,
  }) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminal/distribute',
        data: {
          'to_agent_id': toAgentId,
          'terminal_sn': terminalSn,
          'channel_id': channelId,
          'goods_price': goodsPrice,
          'deduction_type': deductionType,
          if (deductionPeriods != null) 'deduction_periods': deductionPeriods,
          if (remark != null) 'remark': remark,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return TerminalDistribute.fromJson(
          apiResponse.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 批量终端下发
  Future<List<TerminalDistribute>> batchDistributeTerminals({
    required int toAgentId,
    required List<String> terminalSns,
    required int channelId,
    required int goodsPrice,
    required int deductionType,
    int? deductionPeriods,
    String? remark,
  }) async {
    final results = <TerminalDistribute>[];
    for (final sn in terminalSns) {
      try {
        final result = await distributeTerminal(
          toAgentId: toAgentId,
          terminalSn: sn,
          channelId: channelId,
          goodsPrice: goodsPrice,
          deductionType: deductionType,
          deductionPeriods: deductionPeriods,
          remark: remark,
        );
        results.add(result);
      } catch (e) {
        // 记录失败但继续处理其他终端
        // ignore: avoid_print
        debugPrint('Failed to distribute terminal $sn: $e');
      }
    }
    return results;
  }

  /// 确认下发
  Future<void> confirmDistribute(int distributeId) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminal/distribute/$distributeId/confirm',
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 拒绝下发
  Future<void> rejectDistribute(int distributeId) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminal/distribute/$distributeId/reject',
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 取消下发
  Future<void> cancelDistribute(int distributeId) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminal/distribute/$distributeId/cancel',
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 获取下发列表
  Future<PaginatedResponse<TerminalDistribute>> getDistributeList({
    required String direction, // 'from' | 'to'
    List<int>? status,
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final response = await _dio.get(
        '/api/v1/terminal/distribute',
        queryParameters: {
          'direction': direction,
          if (status != null && status.isNotEmpty) 'status': status.join(','),
          'page': page,
          'page_size': pageSize,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return PaginatedResponse.fromJson(
        apiResponse.data as Map<String, dynamic>,
        (json) => TerminalDistribute.fromJson(json),
      );
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 终端回拨
  Future<TerminalRecall> recallTerminal({
    required int toAgentId,
    required String terminalSn,
    int? channelId,
    String? remark,
  }) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/recall',
        data: {
          'to_agent_id': toAgentId,
          'terminal_sn': terminalSn,
          if (channelId != null) 'channel_id': channelId,
          if (remark != null) 'remark': remark,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return TerminalRecall.fromJson(apiResponse.data as Map<String, dynamic>);
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 批量终端回拨
  Future<Map<String, dynamic>> batchRecallTerminals({
    required int toAgentId,
    required List<String> terminalSns,
    String? remark,
  }) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/batch-recall',
        data: {
          'to_agent_id': toAgentId,
          'terminal_sns': terminalSns,
          if (remark != null) 'remark': remark,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return apiResponse.data as Map<String, dynamic>;
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 确认回拨
  Future<void> confirmRecall(int recallId) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/recall/$recallId/confirm',
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 拒绝回拨
  Future<void> rejectRecall(int recallId) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/recall/$recallId/reject',
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 取消回拨
  Future<void> cancelRecall(int recallId) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/recall/$recallId/cancel',
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 获取回拨列表
  Future<PaginatedResponse<TerminalRecall>> getRecallList({
    required String direction, // 'from' | 'to'
    int page = 1,
    int pageSize = 20,
  }) async {
    try {
      final response = await _dio.get(
        '/api/v1/terminals/recall',
        queryParameters: {
          'direction': direction,
          'page': page,
          'page_size': pageSize,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return PaginatedResponse.fromJson(
        apiResponse.data as Map<String, dynamic>,
        (json) => TerminalRecall.fromJson(json),
      );
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 批量设置费率
  /// [terminalSns] 终端SN列表
  /// [creditRate] 信用卡费率（万分比，如55表示0.55%）
  Future<Map<String, dynamic>> batchSetRate({
    required List<String> terminalSns,
    required int creditRate,
  }) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/batch-set-rate',
        data: {
          'terminal_sns': terminalSns,
          'credit_rate': creditRate,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return apiResponse.data as Map<String, dynamic>? ?? {};
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 批量设置流量费
  /// [terminalSns] 终端SN列表
  /// [firstSimFee] 首次流量费（分）
  /// [nonFirstSimFee] 非首次流量费（分）
  /// [simFeeIntervalDays] 流量费间隔天数
  Future<Map<String, dynamic>> batchSetSimFee({
    required List<String> terminalSns,
    required int firstSimFee,
    required int nonFirstSimFee,
    required int simFeeIntervalDays,
  }) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/batch-set-sim',
        data: {
          'terminal_sns': terminalSns,
          'first_sim_fee': firstSimFee,
          'non_first_sim_fee': nonFirstSimFee,
          'sim_fee_interval_days': simFeeIntervalDays,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return apiResponse.data as Map<String, dynamic>? ?? {};
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }

  /// 批量设置押金
  /// [terminalSns] 终端SN列表
  /// [depositAmount] 押金金额（分，0表示无押金）
  Future<Map<String, dynamic>> batchSetDeposit({
    required List<String> terminalSns,
    required int depositAmount,
  }) async {
    try {
      final response = await _dio.post(
        '/api/v1/terminals/batch-set-deposit',
        data: {
          'terminal_sns': terminalSns,
          'deposit_amount': depositAmount,
        },
      );

      final apiResponse = ApiResponse.fromJson(response.data, null);
      if (!apiResponse.isSuccess) {
        throw ApiException(apiResponse.code, apiResponse.message);
      }

      return apiResponse.data as Map<String, dynamic>? ?? {};
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? '网络错误');
    }
  }
}
