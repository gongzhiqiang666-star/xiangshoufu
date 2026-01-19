import 'package:flutter/material.dart';

/// 圆角规范
class AppRadius {
  AppRadius._();

  // ==================== 基础圆角 ====================
  /// 最小圆角 4px - 标签
  static const double xs = 4.0;
  /// 小圆角 8px - 按钮、输入框
  static const double sm = 8.0;
  /// 中圆角 12px - 卡片
  static const double md = 12.0;
  /// 大圆角 16px - 底部弹窗
  static const double lg = 16.0;
  /// 超大圆角 24px - 特殊卡片
  static const double xl = 24.0;
  /// 全圆角
  static const double full = 999.0;

  // ==================== 预设BorderRadius ====================
  /// 卡片圆角
  static BorderRadius get card => BorderRadius.circular(md);
  /// 按钮圆角
  static BorderRadius get button => BorderRadius.circular(sm);
  /// 输入框圆角
  static BorderRadius get input => BorderRadius.circular(sm);
  /// 标签圆角
  static BorderRadius get tag => BorderRadius.circular(xs);
  /// 底部弹窗顶部圆角
  static BorderRadius get bottomSheet => const BorderRadius.only(
        topLeft: Radius.circular(16),
        topRight: Radius.circular(16),
      );
  /// 图片圆角
  static BorderRadius get image => BorderRadius.circular(sm);
  /// 头像圆角
  static BorderRadius get avatar => BorderRadius.circular(full);
  /// 图标背景圆角
  static BorderRadius get iconBg => BorderRadius.circular(sm);
}

/// 阴影规范
class AppShadows {
  AppShadows._();

  /// 小阴影 - 按钮、标签
  static List<BoxShadow> get sm => [
        BoxShadow(
          color: Colors.black.withOpacity(0.05),
          blurRadius: 4,
          offset: const Offset(0, 1),
        ),
      ];

  /// 中阴影 - 卡片
  static List<BoxShadow> get md => [
        BoxShadow(
          color: Colors.black.withOpacity(0.08),
          blurRadius: 8,
          offset: const Offset(0, 2),
        ),
      ];

  /// 大阴影 - 弹窗、悬浮按钮
  static List<BoxShadow> get lg => [
        BoxShadow(
          color: Colors.black.withOpacity(0.1),
          blurRadius: 16,
          offset: const Offset(0, 4),
        ),
      ];

  /// 特大阴影 - 模态框
  static List<BoxShadow> get xl => [
        BoxShadow(
          color: Colors.black.withOpacity(0.15),
          blurRadius: 24,
          offset: const Offset(0, 8),
        ),
      ];

  /// 内阴影 - 输入框聚焦
  static List<BoxShadow> get inner => [
        BoxShadow(
          color: Colors.black.withOpacity(0.05),
          blurRadius: 2,
          offset: const Offset(0, 1),
          spreadRadius: -1,
        ),
      ];

  /// 无阴影
  static List<BoxShadow> get none => [];
}
