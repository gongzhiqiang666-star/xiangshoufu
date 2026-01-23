/// 格式化工具函数
///
/// 提供金额、数字、百分比、日期等格式化功能

import 'package:intl/intl.dart';

/// 格式化金额（分转元）
///
/// [amount] 金额（分）
/// [decimals] 小数位数，默认2位
/// [showSymbol] 是否显示货币符号，默认false
///
/// 示例:
/// ```dart
/// formatAmount(100) // "1.00"
/// formatAmount(12345) // "123.45"
/// formatAmount(10000000) // "100,000.00"
/// formatAmount(100, showSymbol: true) // "¥1.00"
/// ```
String formatAmount(int amount, {int decimals = 2, bool showSymbol = false}) {
  final yuan = amount / 100;
  final formatter = NumberFormat.currency(
    locale: 'zh_CN',
    symbol: showSymbol ? '¥' : '',
    decimalDigits: decimals,
  );
  return formatter.format(yuan).trim();
}

/// 格式化大数字（带单位）
///
/// [num] 数字
///
/// 示例:
/// ```dart
/// formatNumber(1234) // "1,234"
/// formatNumber(10000) // "1.00万"
/// formatNumber(100000000) // "1.00亿"
/// ```
String formatNumber(int num) {
  if (num >= 100000000) {
    return '${(num / 100000000).toStringAsFixed(2)}亿';
  }
  if (num >= 10000) {
    return '${(num / 10000).toStringAsFixed(2)}万';
  }
  return NumberFormat('#,###', 'zh_CN').format(num);
}

/// 格式化百分比
///
/// [value] 值（0-1之间）
/// [decimals] 小数位数，默认2位
///
/// 示例:
/// ```dart
/// formatPercent(0.5) // "50.00%"
/// formatPercent(0.1234) // "12.34%"
/// formatPercent(0.1234, decimals: 1) // "12.3%"
/// ```
String formatPercent(double value, {int decimals = 2}) {
  return '${(value * 100).toStringAsFixed(decimals)}%';
}

/// 格式化日期时间
///
/// [date] 日期
/// [format] 格式类型：'date' | 'datetime' | 'time'
///
/// 示例:
/// ```dart
/// formatDateTime(DateTime.now()) // "2024-03-15 10:30:45"
/// formatDateTime(DateTime.now(), format: 'date') // "2024-03-15"
/// formatDateTime(DateTime.now(), format: 'time') // "10:30:45"
/// ```
String formatDateTime(DateTime date, {String format = 'datetime'}) {
  switch (format) {
    case 'date':
      return DateFormat('yyyy-MM-dd').format(date);
    case 'time':
      return DateFormat('HH:mm:ss').format(date);
    case 'datetime':
    default:
      return DateFormat('yyyy-MM-dd HH:mm:ss').format(date);
  }
}

/// 格式化日期（带空值处理）
///
/// [date] 日期字符串或null
///
/// 示例:
/// ```dart
/// formatDate('2024-03-15T10:30:45') // "2024-03-15 10:30:45"
/// formatDate(null) // "-"
/// formatDate('') // "-"
/// ```
String formatDate(String? date) {
  if (date == null || date.isEmpty) {
    return '-';
  }
  try {
    final dateTime = DateTime.parse(date);
    return formatDateTime(dateTime);
  } catch (e) {
    return '-';
  }
}

/// 计算趋势百分比
///
/// [current] 当前值
/// [previous] 前期值
///
/// 示例:
/// ```dart
/// calculateTrend(150, 100) // 50.0 (上涨50%)
/// calculateTrend(50, 100) // -50.0 (下跌50%)
/// calculateTrend(100, 0) // 100.0
/// ```
double calculateTrend(int current, int previous) {
  if (previous == 0) {
    return current > 0 ? 100.0 : 0.0;
  }
  return ((current - previous) / previous) * 100;
}

/// 格式化手机号（隐藏中间4位）
///
/// [phone] 手机号
///
/// 示例:
/// ```dart
/// formatPhone('13800138000') // "138****8000"
/// ```
String formatPhone(String phone) {
  if (phone.length != 11) {
    return phone;
  }
  return '${phone.substring(0, 3)}****${phone.substring(7)}';
}

/// 格式化银行卡号（只显示后4位）
///
/// [cardNo] 银行卡号
///
/// 示例:
/// ```dart
/// formatBankCard('6222021234567890') // "**** **** **** 7890"
/// ```
String formatBankCard(String cardNo) {
  if (cardNo.length < 4) {
    return cardNo;
  }
  final lastFour = cardNo.substring(cardNo.length - 4);
  return '**** **** **** $lastFour';
}
