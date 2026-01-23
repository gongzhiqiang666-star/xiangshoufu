/// 测试工具函数集合
///
/// 提供测试中常用的Mock对象、工具函数和辅助方法

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

/// 创建带 ProviderScope 的测试 Widget
///
/// [child] 要测试的 Widget
/// [overrides] Provider 覆盖列表
Widget createTestApp({
  required Widget child,
  List<Override> overrides = const [],
}) {
  return ProviderScope(
    overrides: overrides,
    child: MaterialApp(
      home: Scaffold(body: child),
    ),
  );
}

/// 创建模拟用户数据
Map<String, dynamic> createMockUser({
  int? id,
  String? username,
  String? realName,
  String? phone,
  int? roleType,
}) {
  return {
    'id': id ?? 1,
    'username': username ?? 'testuser',
    'real_name': realName ?? '测试用户',
    'phone': phone ?? '13800138000',
    'role_type': roleType ?? 1,
  };
}

/// 创建模拟登录响应
Map<String, dynamic> createMockLoginResponse({
  String? accessToken,
  int? expiresIn,
  Map<String, dynamic>? user,
  Map<String, dynamic>? agent,
}) {
  return {
    'access_token': accessToken ?? 'mock-token-12345',
    'expires_in': expiresIn ?? 86400,
    'user': user ?? createMockUser(),
    'agent': agent ?? {'agent_name': '测试代理商'},
  };
}

/// 创建模拟仪表盘数据
Map<String, dynamic> createMockDashboardData({
  int? todayAmount,
  int? todayCount,
  int? monthAmount,
  int? monthCount,
  double? trend,
}) {
  return {
    'today_amount': todayAmount ?? 100000,
    'today_count': todayCount ?? 50,
    'month_amount': monthAmount ?? 3000000,
    'month_count': monthCount ?? 1500,
    'trend': trend ?? 0.15,
  };
}

/// 创建模拟交易数据
Map<String, dynamic> createMockTransaction({
  String? id,
  int? amount,
  String? status,
  String? merchantName,
  DateTime? createdAt,
}) {
  return {
    'id': id ?? 'TX20240101001',
    'amount': amount ?? 10000,
    'status': status ?? 'success',
    'merchant_name': merchantName ?? '测试商户',
    'created_at': (createdAt ?? DateTime.now()).toIso8601String(),
  };
}

/// 创建模拟钱包数据
Map<String, dynamic> createMockWallet({
  int? balance,
  int? frozenAmount,
  int? totalIncome,
  int? totalWithdraw,
}) {
  return {
    'balance': balance ?? 500000,
    'frozen_amount': frozenAmount ?? 10000,
    'total_income': totalIncome ?? 1000000,
    'total_withdraw': totalWithdraw ?? 500000,
  };
}

/// 等待异步操作完成
Future<void> pumpAndSettle(WidgetTester tester, {Duration? duration}) async {
  await tester.pump(duration ?? const Duration(milliseconds: 100));
  await tester.pumpAndSettle();
}

/// 创建 ProviderContainer 用于 Provider 测试
ProviderContainer createTestContainer({
  List<Override> overrides = const [],
}) {
  return ProviderContainer(overrides: overrides);
}
