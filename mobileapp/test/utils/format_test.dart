/// 格式化工具函数测试
///
/// 覆盖: 正常流程、边界情况、错误处理、特殊输入

import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/utils/format.dart';

void main() {
  group('formatAmount', () {
    // ✅ 正常流程 (Happy Path)
    test('should format 100 cents to "1.00"', () {
      expect(formatAmount(100), equals('1.00'));
    });

    test('should format 12345 cents to "123.45"', () {
      expect(formatAmount(12345), equals('123.45'));
    });

    // ✅ 边界情况 (Edge Cases)
    test('should handle zero', () {
      expect(formatAmount(0), equals('0.00'));
    });

    test('should handle small values (1 cent)', () {
      expect(formatAmount(1), equals('0.01'));
    });

    // ✅ 错误处理 (Error Handling)
    test('should handle negative values', () {
      expect(formatAmount(-100), equals('-1.00'));
    });

    // ✅ 特殊输入 (Special Inputs)
    test('should format large numbers with thousand separators', () {
      expect(formatAmount(10000000), equals('100,000.00'));
    });

    test('should respect custom decimals parameter', () {
      expect(formatAmount(12345, decimals: 0), equals('123'));
      expect(formatAmount(12345, decimals: 1), equals('123.5'));
      expect(formatAmount(12345, decimals: 3), equals('123.450'));
    });

    test('should show currency symbol when requested', () {
      expect(formatAmount(100, showSymbol: true), equals('¥1.00'));
    });
  });

  group('formatNumber', () {
    // ✅ 正常流程
    test('should format regular numbers with locale', () {
      expect(formatNumber(1234), equals('1,234'));
    });

    // ✅ 边界情况
    test('should handle zero', () {
      expect(formatNumber(0), equals('0'));
    });

    test('should format numbers >= 10000 with 万 unit', () {
      expect(formatNumber(10000), equals('1.00万'));
      expect(formatNumber(12345), equals('1.23万'));
    });

    test('should format numbers >= 100000000 with 亿 unit', () {
      expect(formatNumber(100000000), equals('1.00亿'));
      expect(formatNumber(123456789), equals('1.23亿'));
    });

    // ✅ 特殊输入
    test('should handle boundary between units', () {
      expect(formatNumber(9999), equals('9,999'));
      expect(formatNumber(10000), equals('1.00万'));
    });
  });

  group('formatPercent', () {
    // ✅ 正常流程
    test('should format 0.5 to "50.00%"', () {
      expect(formatPercent(0.5), equals('50.00%'));
    });

    test('should format 0.1234 to "12.34%"', () {
      expect(formatPercent(0.1234), equals('12.34%'));
    });

    // ✅ 边界情况
    test('should handle zero', () {
      expect(formatPercent(0), equals('0.00%'));
    });

    test('should handle 100%', () {
      expect(formatPercent(1), equals('100.00%'));
    });

    // ✅ 特殊输入
    test('should handle values > 100%', () {
      expect(formatPercent(1.5), equals('150.00%'));
    });

    test('should respect custom decimals', () {
      expect(formatPercent(0.12345, decimals: 0), equals('12%'));
      expect(formatPercent(0.12345, decimals: 1), equals('12.3%'));
    });
  });

  group('formatDateTime', () {
    final testDate = DateTime(2024, 3, 15, 10, 30, 45);

    // ✅ 正常流程
    test('should format date to datetime by default', () {
      expect(formatDateTime(testDate), equals('2024-03-15 10:30:45'));
    });

    // ✅ 不同格式
    test('should format to date only', () {
      expect(formatDateTime(testDate, format: 'date'), equals('2024-03-15'));
    });

    test('should format to time only', () {
      expect(formatDateTime(testDate, format: 'time'), equals('10:30:45'));
    });

    // ✅ 边界情况
    test('should pad single digit months and days', () {
      final date = DateTime(2024, 1, 5, 9, 5, 5);
      expect(formatDateTime(date, format: 'date'), equals('2024-01-05'));
      expect(formatDateTime(date, format: 'time'), equals('09:05:05'));
    });
  });

  group('formatDate', () {
    // ✅ 正常流程
    test('should format valid date string', () {
      expect(formatDate('2024-03-15T10:30:45'), equals('2024-03-15 10:30:45'));
    });

    // ✅ 空值处理
    test('should return "-" for null', () {
      expect(formatDate(null), equals('-'));
    });

    test('should return "-" for empty string', () {
      expect(formatDate(''), equals('-'));
    });

    // ✅ 错误处理
    test('should return "-" for invalid date string', () {
      expect(formatDate('invalid-date'), equals('-'));
    });
  });

  group('calculateTrend', () {
    // ✅ 正常流程
    test('should calculate positive trend', () {
      expect(calculateTrend(150, 100), equals(50.0));
    });

    test('should calculate negative trend', () {
      expect(calculateTrend(50, 100), equals(-50.0));
    });

    // ✅ 边界情况
    test('should return 0 when both values are 0', () {
      expect(calculateTrend(0, 0), equals(0.0));
    });

    test('should return 100 when previous is 0 and current is positive', () {
      expect(calculateTrend(100, 0), equals(100.0));
    });

    // ✅ 特殊输入
    test('should handle same values (0% change)', () {
      expect(calculateTrend(100, 100), equals(0.0));
    });
  });

  group('formatPhone', () {
    // ✅ 正常流程
    test('should mask middle 4 digits', () {
      expect(formatPhone('13800138000'), equals('138****8000'));
    });

    // ✅ 边界情况
    test('should return original if not 11 digits', () {
      expect(formatPhone('1234567'), equals('1234567'));
      expect(formatPhone(''), equals(''));
    });
  });

  group('formatBankCard', () {
    // ✅ 正常流程
    test('should show only last 4 digits', () {
      expect(formatBankCard('6222021234567890'), equals('**** **** **** 7890'));
    });

    // ✅ 边界情况
    test('should return original if less than 4 digits', () {
      expect(formatBankCard('123'), equals('123'));
    });
  });
}
