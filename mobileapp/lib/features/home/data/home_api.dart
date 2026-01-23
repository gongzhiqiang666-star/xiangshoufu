import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../domain/home_model.dart';

/// 首页API服务
class HomeApi {
  final Dio _dio;

  HomeApi(this._dio);

  /// 获取首页概览数据
  Future<HomeOverviewData> getOverview({String scope = 'direct'}) async {
    final response = await _dio.get(
      '/api/v1/dashboard/overview',
      queryParameters: {'scope': scope},
    );

    if (response.data['code'] == 0) {
      return HomeOverviewData.fromJson(response.data['data']);
    }
    throw Exception(response.data['message'] ?? '获取数据失败');
  }

  /// 获取最近交易列表
  Future<List<RecentTransaction>> getRecentTransactions({int limit = 10}) async {
    final response = await _dio.get(
      '/api/v1/dashboard/recent-transactions',
      queryParameters: {'limit': limit},
    );

    if (response.data['code'] == 0) {
      final list = response.data['data']['transactions'] as List? ?? [];
      return list.map((e) => RecentTransaction.fromJson(e)).toList();
    }
    throw Exception(response.data['message'] ?? '获取交易列表失败');
  }

  /// 获取图表趋势数据
  Future<List<TrendPoint>> getTrendData({
    int days = 7,
    String scope = 'direct',
  }) async {
    final response = await _dio.get(
      '/api/v1/dashboard/charts',
      queryParameters: {'days': days, 'scope': scope},
    );

    if (response.data['code'] == 0) {
      final list = response.data['data']['trans_trend'] as List? ?? [];
      return list.map((e) => TrendPoint.fromJson(e)).toList();
    }
    throw Exception(response.data['message'] ?? '获取趋势数据失败');
  }

  /// 获取通道统计
  Future<List<ChannelStats>> getChannelStats({
    String scope = 'direct',
    String period = 'month',
  }) async {
    final response = await _dio.get(
      '/api/v1/dashboard/channel-stats',
      queryParameters: {'scope': scope, 'period': period},
    );

    if (response.data['code'] == 0) {
      final list = response.data['data']['channel_stats'] as List? ?? [];
      return list.map((e) => ChannelStats.fromJson(e)).toList();
    }
    throw Exception(response.data['message'] ?? '获取通道统计失败');
  }

  /// 获取商户分布
  Future<List<MerchantDistribution>> getMerchantDistribution({
    String scope = 'direct',
  }) async {
    final response = await _dio.get(
      '/api/v1/dashboard/merchant-distribution',
      queryParameters: {'scope': scope},
    );

    if (response.data['code'] == 0) {
      final list = response.data['data']['distribution'] as List? ?? [];
      return list.map((e) => MerchantDistribution.fromJson(e)).toList();
    }
    throw Exception(response.data['message'] ?? '获取商户分布失败');
  }
}

/// 趋势数据点
class TrendPoint {
  final String date;
  final int transAmount;
  final double transAmountYuan;
  final int transCount;
  final int profitTotal;
  final double profitTotalYuan;

  TrendPoint({
    required this.date,
    required this.transAmount,
    required this.transAmountYuan,
    required this.transCount,
    required this.profitTotal,
    required this.profitTotalYuan,
  });

  factory TrendPoint.fromJson(Map<String, dynamic> json) {
    return TrendPoint(
      date: json['date'] ?? '',
      transAmount: json['trans_amount'] ?? 0,
      transAmountYuan: (json['trans_amount_yuan'] ?? 0).toDouble(),
      transCount: json['trans_count'] ?? 0,
      profitTotal: json['profit_total'] ?? 0,
      profitTotalYuan: (json['profit_total_yuan'] ?? 0).toDouble(),
    );
  }
}

/// 通道统计
class ChannelStats {
  final int channelId;
  final String channelCode;
  final String channelName;
  final int transAmount;
  final int transCount;
  final double percentage;

  ChannelStats({
    required this.channelId,
    required this.channelCode,
    required this.channelName,
    required this.transAmount,
    required this.transCount,
    required this.percentage,
  });

  factory ChannelStats.fromJson(Map<String, dynamic> json) {
    return ChannelStats(
      channelId: json['channel_id'] ?? 0,
      channelCode: json['channel_code'] ?? '',
      channelName: json['channel_name'] ?? '',
      transAmount: json['trans_amount'] ?? 0,
      transCount: json['trans_count'] ?? 0,
      percentage: (json['percentage'] ?? 0).toDouble(),
    );
  }
}

/// 商户分布
class MerchantDistribution {
  final String merchantType;
  final String typeName;
  final int count;
  final double percentage;

  MerchantDistribution({
    required this.merchantType,
    required this.typeName,
    required this.count,
    required this.percentage,
  });

  factory MerchantDistribution.fromJson(Map<String, dynamic> json) {
    return MerchantDistribution(
      merchantType: json['merchant_type'] ?? '',
      typeName: json['type_name'] ?? '',
      count: json['count'] ?? 0,
      percentage: (json['percentage'] ?? 0).toDouble(),
    );
  }
}

/// HomeApi Provider
final homeApiProvider = Provider<HomeApi>((ref) {
  final dio = ref.watch(dioProvider);
  return HomeApi(dio);
});
