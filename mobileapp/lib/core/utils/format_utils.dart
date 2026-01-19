import 'package:intl/intl.dart';

/// 格式化工具类
/// 统一处理金额、日期、手机号等格式化
class FormatUtils {
  FormatUtils._();

  // ==================== 金额格式化 ====================

  /// 格式化金额（分 -> 元）
  /// [cents] 金额（分）
  /// [showSign] 是否显示正负号
  /// [showSymbol] 是否显示¥符号
  static String formatCents(
    int? cents, {
    bool showSign = false,
    bool showSymbol = true,
  }) {
    if (cents == null) return showSymbol ? '¥0.00' : '0.00';
    final yuan = cents / 100;
    return formatYuan(yuan, showSign: showSign, showSymbol: showSymbol);
  }

  /// 格式化金额（元）
  /// [yuan] 金额（元）
  /// [showSign] 是否显示正负号
  /// [showSymbol] 是否显示¥符号
  static String formatYuan(
    double? yuan, {
    bool showSign = false,
    bool showSymbol = true,
  }) {
    if (yuan == null) return showSymbol ? '¥0.00' : '0.00';
    final formatter = NumberFormat('#,##0.00', 'zh_CN');
    final sign = showSign && yuan > 0 ? '+' : '';
    final symbol = showSymbol ? '¥' : '';
    return '$sign$symbol${formatter.format(yuan)}';
  }

  /// 格式化大金额（自动转换为万）
  /// [yuan] 金额（元）
  static String formatLargeAmount(double? yuan) {
    if (yuan == null) return '¥0.00';
    if (yuan.abs() >= 10000) {
      final wan = yuan / 10000;
      return '¥${wan.toStringAsFixed(2)}万';
    }
    return formatYuan(yuan);
  }

  /// 格式化金额（简短模式：1.2k, 3.5w）
  static String formatShortAmount(double? yuan) {
    if (yuan == null) return '¥0';
    if (yuan.abs() >= 100000000) {
      return '¥${(yuan / 100000000).toStringAsFixed(1)}亿';
    }
    if (yuan.abs() >= 10000) {
      return '¥${(yuan / 10000).toStringAsFixed(1)}万';
    }
    if (yuan.abs() >= 1000) {
      return '¥${(yuan / 1000).toStringAsFixed(1)}k';
    }
    return '¥${yuan.toStringAsFixed(0)}';
  }

  // ==================== 费率格式化 ====================

  /// 格式化费率（小数 -> 百分比）
  /// [rate] 费率（如0.0055表示0.55%）
  static String formatRate(double? rate) {
    if (rate == null) return '0.00%';
    return '${(rate * 100).toStringAsFixed(2)}%';
  }

  /// 格式化万分比费率
  /// [bps] 万分比（如55表示万分之55）
  static String formatBps(int? bps) {
    if (bps == null) return '0.00%';
    return '${(bps / 100).toStringAsFixed(2)}%';
  }

  // ==================== 数字格式化 ====================

  /// 格式化数字（千分位）
  static String formatNumber(int? number) {
    if (number == null) return '0';
    final formatter = NumberFormat('#,###', 'zh_CN');
    return formatter.format(number);
  }

  /// 格式化百分比变化
  /// [change] 变化比例（如0.125表示12.5%）
  static String formatChange(double? change) {
    if (change == null) return '0.0%';
    final sign = change >= 0 ? '+' : '';
    return '$sign${(change * 100).toStringAsFixed(1)}%';
  }

  // ==================== 日期时间格式化 ====================

  /// 格式化日期
  /// [date] 日期
  /// [pattern] 格式模式，默认 yyyy-MM-dd
  static String formatDate(DateTime? date, {String pattern = 'yyyy-MM-dd'}) {
    if (date == null) return '';
    return DateFormat(pattern, 'zh_CN').format(date);
  }

  /// 格式化时间
  static String formatTime(DateTime? date) {
    if (date == null) return '';
    return DateFormat('HH:mm', 'zh_CN').format(date);
  }

  /// 格式化日期时间
  static String formatDateTime(DateTime? date) {
    if (date == null) return '';
    return DateFormat('yyyy-MM-dd HH:mm', 'zh_CN').format(date);
  }

  /// 格式化相对时间
  static String formatRelativeTime(DateTime? date) {
    if (date == null) return '';
    final now = DateTime.now();
    final diff = now.difference(date);

    if (diff.inSeconds < 60) {
      return '刚刚';
    } else if (diff.inMinutes < 60) {
      return '${diff.inMinutes}分钟前';
    } else if (diff.inHours < 24) {
      return '${diff.inHours}小时前';
    } else if (diff.inDays == 1) {
      return '昨天 ${formatTime(date)}';
    } else if (diff.inDays == 2) {
      return '前天 ${formatTime(date)}';
    } else if (diff.inDays < 7) {
      return '${diff.inDays}天前';
    } else if (date.year == now.year) {
      return DateFormat('MM-dd HH:mm', 'zh_CN').format(date);
    } else {
      return formatDateTime(date);
    }
  }

  /// 格式化月份
  static String formatMonth(DateTime? date) {
    if (date == null) return '';
    return DateFormat('yyyy年MM月', 'zh_CN').format(date);
  }

  // ==================== 脱敏格式化 ====================

  /// 手机号脱敏
  /// 138****8888
  static String maskPhone(String? phone) {
    if (phone == null || phone.length != 11) return phone ?? '';
    return '${phone.substring(0, 3)}****${phone.substring(7)}';
  }

  /// 身份证号脱敏
  /// 110***********1234
  static String maskIdCard(String? idCard) {
    if (idCard == null || idCard.length != 18) return idCard ?? '';
    return '${idCard.substring(0, 3)}***********${idCard.substring(14)}';
  }

  /// 银行卡号脱敏
  /// **** **** **** 5678
  static String maskBankCard(String? cardNo) {
    if (cardNo == null || cardNo.length < 4) return cardNo ?? '';
    return '**** **** **** ${cardNo.substring(cardNo.length - 4)}';
  }

  /// 姓名脱敏
  /// 张*三
  static String maskName(String? name) {
    if (name == null || name.isEmpty) return '';
    if (name.length == 2) {
      return '${name[0]}*';
    } else if (name.length > 2) {
      return '${name[0]}${'*' * (name.length - 2)}${name[name.length - 1]}';
    }
    return name;
  }

  // ==================== SN号格式化 ====================

  /// 格式化SN号（添加空格分隔）
  /// 1234 5678 9012
  static String formatSN(String? sn) {
    if (sn == null || sn.isEmpty) return '';
    final buffer = StringBuffer();
    for (var i = 0; i < sn.length; i++) {
      if (i > 0 && i % 4 == 0) {
        buffer.write(' ');
      }
      buffer.write(sn[i]);
    }
    return buffer.toString();
  }
}
