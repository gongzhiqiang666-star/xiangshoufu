import 'package:flutter/material.dart';

/// 应用颜色规范
/// 设计系统 - 符合中国金融类APP设计规范
class AppColors {
  AppColors._();

  // ==================== 主色系 ====================
  /// 品牌蓝 - 主色调
  static const Color primary = Color(0xFF2563EB);
  /// 浅蓝 - 次要强调
  static const Color primaryLight = Color(0xFF60A5FA);
  /// 深蓝 - 按压状态
  static const Color primaryDark = Color(0xFF1D4ED8);

  // ==================== 功能色 ====================
  /// 成功绿 - 收益增长、成功状态
  static const Color success = Color(0xFF10B981);
  /// 警告橙 - 提醒、警告状态
  static const Color warning = Color(0xFFF59E0B);
  /// 危险红 - 错误、亏损状态
  static const Color danger = Color(0xFFEF4444);
  /// 信息蓝 - 提示信息
  static const Color info = Color(0xFF3B82F6);

  // ==================== 中性色 ====================
  /// 主文本 - 标题、重要内容
  static const Color textPrimary = Color(0xFF1F2937);
  /// 次要文本 - 正文、说明
  static const Color textSecondary = Color(0xFF6B7280);
  /// 辅助文本 - 提示、时间戳
  static const Color textTertiary = Color(0xFF9CA3AF);
  /// 禁用文本
  static const Color textDisabled = Color(0xFFD1D5DB);

  /// 边框色
  static const Color border = Color(0xFFE5E7EB);
  /// 分割线
  static const Color divider = Color(0xFFF3F4F6);
  /// 页面背景
  static const Color background = Color(0xFFF9FAFB);
  /// 卡片背景
  static const Color cardBg = Color(0xFFFFFFFF);

  // ==================== 分润类型专属色 ====================
  /// 交易分润 - 蓝色
  static const Color profitTrade = Color(0xFF2563EB);
  /// 押金返现 - 绿色
  static const Color profitDeposit = Color(0xFF10B981);
  /// 流量返现 - 橙色
  static const Color profitSim = Color(0xFFF59E0B);
  /// 激活奖励 - 紫色
  static const Color profitReward = Color(0xFF8B5CF6);

  // ==================== 支付类型色 ====================
  /// 微信支付
  static const Color wechatPay = Color(0xFF07C160);
  /// 支付宝
  static const Color alipay = Color(0xFF1677FF);
  /// 银联
  static const Color unionPay = Color(0xFFE60012);

  // ==================== 钱包卡片渐变色 ====================
  /// 分润钱包渐变
  static const List<Color> walletProfitGradient = [
    Color(0xFF2563EB),
    Color(0xFF1D4ED8),
  ];
  /// 服务费钱包渐变
  static const List<Color> walletServiceGradient = [
    Color(0xFF10B981),
    Color(0xFF059669),
  ];
  /// 奖励钱包渐变
  static const List<Color> walletRewardGradient = [
    Color(0xFF8B5CF6),
    Color(0xFF7C3AED),
  ];

  // ==================== 状态色 ====================
  /// 已激活
  static const Color statusActivated = Color(0xFF10B981);
  /// 未激活
  static const Color statusInactive = Color(0xFF9CA3AF);
  /// 进行中
  static const Color statusPending = Color(0xFFF59E0B);
  /// 已完成
  static const Color statusCompleted = Color(0xFF2563EB);
}
