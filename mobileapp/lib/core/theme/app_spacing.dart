import 'package:flutter/material.dart';

/// 间距规范
/// 设计系统 - 统一间距保证视觉一致性
class AppSpacing {
  AppSpacing._();

  // ==================== 基础间距 ====================
  /// 最小间距 4px
  static const double xs = 4.0;
  /// 小间距 8px
  static const double sm = 8.0;
  /// 中间距 16px - 常用
  static const double md = 16.0;
  /// 大间距 24px
  static const double lg = 24.0;
  /// 超大间距 32px
  static const double xl = 32.0;
  /// 最大间距 48px
  static const double xxl = 48.0;

  // ==================== 页面间距 ====================
  /// 页面水平内边距
  static const double pagePaddingH = 16.0;
  /// 页面垂直内边距
  static const double pagePaddingV = 16.0;
  /// 页面边距
  static const EdgeInsets pagePadding = EdgeInsets.all(16.0);
  /// 页面水平边距
  static const EdgeInsets pageHorizontal = EdgeInsets.symmetric(horizontal: 16.0);

  // ==================== 卡片间距 ====================
  /// 卡片内边距
  static const double cardPadding = 16.0;
  /// 卡片间距
  static const double cardGap = 12.0;
  /// 卡片边距
  static const EdgeInsets cardMargin = EdgeInsets.only(bottom: 12.0);

  // ==================== 列表间距 ====================
  /// 列表项间距
  static const double listItemGap = 8.0;
  /// 列表项内边距
  static const EdgeInsets listItemPadding = EdgeInsets.symmetric(
    horizontal: 16.0,
    vertical: 12.0,
  );

  // ==================== 表单间距 ====================
  /// 表单项间距
  static const double formItemGap = 16.0;
  /// 表单标签与输入框间距
  static const double formLabelGap = 8.0;

  // ==================== 组件内间距 ====================
  /// 图标与文字间距
  static const double iconTextGap = 8.0;
  /// 按钮内边距
  static const EdgeInsets buttonPadding = EdgeInsets.symmetric(
    horizontal: 24.0,
    vertical: 12.0,
  );
  /// 标签内边距
  static const EdgeInsets tagPadding = EdgeInsets.symmetric(
    horizontal: 8.0,
    vertical: 4.0,
  );

  // ==================== 网格间距 ====================
  /// 宫格间距
  static const double gridGap = 12.0;
  /// 2列网格
  static const int gridColumns2 = 2;
  /// 4列网格
  static const int gridColumns4 = 4;

  // ==================== 安全区域 ====================
  /// 底部安全区域额外间距
  static const double bottomSafeArea = 34.0;
  /// 底部按钮区域高度
  static const double bottomBarHeight = 80.0;
}
