# APP设计稿 - 代理商分润管理系统

## 一、技术选型

| 层级 | 技术选择 | 说明 |
|------|----------|------|
| **开发框架** | Flutter 3.x | 跨平台：iOS/Android/鸿蒙 |
| **开发语言** | Dart | Flutter 官方语言 |
| **状态管理** | Riverpod 2.0 | 简洁、类型安全 |
| **路由** | go_router | 官方推荐 |
| **HTTP** | dio | 强大的网络请求库 |
| **本地存储** | shared_preferences + sqflite | 简单配置 + 结构化数据 |
| **图表** | fl_chart | Flutter 原生图表 |

---

## 二、设计系统（Design System）

### 2.1 颜色规范

```dart
// lib/core/theme/app_colors.dart

class AppColors {
  // 主色系
  static const Color primary = Color(0xFF2563EB);        // 品牌蓝
  static const Color primaryLight = Color(0xFF60A5FA);   // 浅蓝
  static const Color primaryDark = Color(0xFF1D4ED8);    // 深蓝

  // 功能色
  static const Color success = Color(0xFF10B981);        // 成功绿
  static const Color warning = Color(0xFFF59E0B);        // 警告橙
  static const Color danger = Color(0xFFEF4444);         // 危险红
  static const Color info = Color(0xFF3B82F6);           // 信息蓝

  // 中性色
  static const Color textPrimary = Color(0xFF1F2937);    // 主文本
  static const Color textSecondary = Color(0xFF6B7280);  // 次要文本
  static const Color textTertiary = Color(0xFF9CA3AF);   // 辅助文本
  static const Color border = Color(0xFFE5E7EB);         // 边框
  static const Color divider = Color(0xFFF3F4F6);        // 分割线
  static const Color background = Color(0xFFF9FAFB);     // 背景
  static const Color cardBg = Color(0xFFFFFFFF);         // 卡片背景

  // 分润类型颜色
  static const Color profitTrade = Color(0xFF2563EB);    // 交易分润
  static const Color profitDeposit = Color(0xFF10B981);  // 押金返现
  static const Color profitSim = Color(0xFFF59E0B);      // 流量返现
  static const Color profitReward = Color(0xFF8B5CF6);   // 激活奖励
}
```

### 2.2 字体规范

```dart
// lib/core/theme/app_typography.dart

class AppTypography {
  // 标题
  static const TextStyle h1 = TextStyle(
    fontSize: 24,
    fontWeight: FontWeight.w700,
    height: 1.4,
    color: AppColors.textPrimary,
  );

  static const TextStyle h2 = TextStyle(
    fontSize: 20,
    fontWeight: FontWeight.w600,
    height: 1.4,
    color: AppColors.textPrimary,
  );

  static const TextStyle h3 = TextStyle(
    fontSize: 18,
    fontWeight: FontWeight.w600,
    height: 1.4,
    color: AppColors.textPrimary,
  );

  // 正文
  static const TextStyle body1 = TextStyle(
    fontSize: 16,
    fontWeight: FontWeight.w400,
    height: 1.5,
    color: AppColors.textPrimary,
  );

  static const TextStyle body2 = TextStyle(
    fontSize: 14,
    fontWeight: FontWeight.w400,
    height: 1.5,
    color: AppColors.textSecondary,
  );

  // 辅助
  static const TextStyle caption = TextStyle(
    fontSize: 12,
    fontWeight: FontWeight.w400,
    height: 1.4,
    color: AppColors.textTertiary,
  );

  // 金额
  static const TextStyle amount = TextStyle(
    fontSize: 28,
    fontWeight: FontWeight.w700,
    height: 1.2,
    color: AppColors.textPrimary,
  );

  static const TextStyle amountSmall = TextStyle(
    fontSize: 20,
    fontWeight: FontWeight.w600,
    height: 1.2,
    color: AppColors.textPrimary,
  );
}
```

### 2.3 间距规范

```dart
// lib/core/theme/app_spacing.dart

class AppSpacing {
  static const double xs = 4.0;
  static const double sm = 8.0;
  static const double md = 16.0;
  static const double lg = 24.0;
  static const double xl = 32.0;
  static const double xxl = 48.0;

  // 页面边距
  static const double pagePadding = 16.0;

  // 卡片间距
  static const double cardGap = 12.0;

  // 列表项间距
  static const double listItemGap = 8.0;
}
```

### 2.4 圆角规范

```dart
// lib/core/theme/app_radius.dart

class AppRadius {
  static const double xs = 4.0;
  static const double sm = 8.0;
  static const double md = 12.0;
  static const double lg = 16.0;
  static const double xl = 24.0;
  static const double full = 999.0;

  static BorderRadius get cardRadius => BorderRadius.circular(md);
  static BorderRadius get buttonRadius => BorderRadius.circular(sm);
  static BorderRadius get inputRadius => BorderRadius.circular(sm);
  static BorderRadius get tagRadius => BorderRadius.circular(xs);
}
```

### 2.5 阴影规范

```dart
// lib/core/theme/app_shadows.dart

class AppShadows {
  static List<BoxShadow> get sm => [
    BoxShadow(
      color: Colors.black.withOpacity(0.05),
      blurRadius: 4,
      offset: const Offset(0, 1),
    ),
  ];

  static List<BoxShadow> get md => [
    BoxShadow(
      color: Colors.black.withOpacity(0.08),
      blurRadius: 8,
      offset: const Offset(0, 2),
    ),
  ];

  static List<BoxShadow> get lg => [
    BoxShadow(
      color: Colors.black.withOpacity(0.1),
      blurRadius: 16,
      offset: const Offset(0, 4),
    ),
  ];
}
```

### 2.6 主题配置

```dart
// lib/core/theme/app_theme.dart

import 'package:flutter/material.dart';

class AppTheme {
  static ThemeData get light => ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
      seedColor: AppColors.primary,
      brightness: Brightness.light,
    ),
    scaffoldBackgroundColor: AppColors.background,
    appBarTheme: const AppBarTheme(
      backgroundColor: Colors.white,
      foregroundColor: AppColors.textPrimary,
      elevation: 0,
      centerTitle: true,
      titleTextStyle: TextStyle(
        fontSize: 18,
        fontWeight: FontWeight.w600,
        color: AppColors.textPrimary,
      ),
    ),
    cardTheme: CardTheme(
      color: AppColors.cardBg,
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: AppRadius.cardRadius,
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: AppColors.primary,
        foregroundColor: Colors.white,
        minimumSize: const Size(double.infinity, 48),
        shape: RoundedRectangleBorder(
          borderRadius: AppRadius.buttonRadius,
        ),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w600,
        ),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: AppColors.primary,
        minimumSize: const Size(double.infinity, 48),
        side: const BorderSide(color: AppColors.primary),
        shape: RoundedRectangleBorder(
          borderRadius: AppRadius.buttonRadius,
        ),
      ),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: Colors.white,
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      border: OutlineInputBorder(
        borderRadius: AppRadius.inputRadius,
        borderSide: const BorderSide(color: AppColors.border),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: AppRadius.inputRadius,
        borderSide: const BorderSide(color: AppColors.border),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: AppRadius.inputRadius,
        borderSide: const BorderSide(color: AppColors.primary, width: 2),
      ),
    ),
    dividerTheme: const DividerThemeData(
      color: AppColors.divider,
      thickness: 1,
      space: 1,
    ),
  );
}
```

---

## 三、项目结构

```
mobile_app/
├── lib/
│   ├── main.dart                      # 入口文件
│   │
│   ├── core/                          # 核心模块
│   │   ├── theme/                     # 主题设计系统
│   │   │   ├── app_colors.dart
│   │   │   ├── app_typography.dart
│   │   │   ├── app_spacing.dart
│   │   │   ├── app_radius.dart
│   │   │   ├── app_shadows.dart
│   │   │   └── app_theme.dart
│   │   ├── network/                   # 网络层
│   │   │   ├── api_client.dart
│   │   │   ├── api_endpoints.dart
│   │   │   └── interceptors/
│   │   ├── storage/                   # 本地存储
│   │   ├── utils/                     # 工具类
│   │   │   ├── format_utils.dart      # 格式化工具
│   │   │   ├── validator.dart         # 校验工具
│   │   │   └── platform_utils.dart    # 平台判断
│   │   └── constants/                 # 常量定义
│   │
│   ├── shared/                        # 共享组件
│   │   ├── widgets/                   # 通用组件
│   │   │   ├── buttons/
│   │   │   ├── cards/
│   │   │   ├── inputs/
│   │   │   ├── dialogs/
│   │   │   ├── charts/
│   │   │   └── loading/
│   │   └── extensions/                # 扩展方法
│   │
│   ├── features/                      # 功能模块
│   │   ├── auth/                      # 认证模块
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   ├── home/                      # 首页
│   │   ├── agent/                     # 代理拓展
│   │   ├── terminal/                  # 终端管理
│   │   ├── cargo_deduction/           # 货款代扣
│   │   ├── merchant/                  # 商户管理
│   │   ├── data_analysis/             # 数据分析
│   │   ├── profit/                    # 收益统计
│   │   ├── wallet/                    # 钱包
│   │   ├── deduction/                 # 代扣管理
│   │   ├── marketing/                 # 营销海报
│   │   ├── message/                   # 消息通知
│   │   └── profile/                   # 我的信息
│   │
│   └── router/                        # 路由配置
│       └── app_router.dart
│
├── assets/                            # 资源文件
│   ├── images/
│   ├── icons/
│   └── fonts/
│
├── pubspec.yaml
└── README.md
```

---

## 四、核心组件库

### 4.1 统计卡片组件

```dart
// lib/shared/widgets/cards/stat_card.dart

import 'package:flutter/material.dart';

class StatCard extends StatelessWidget {
  final String title;
  final String value;
  final String? subtitle;
  final IconData? icon;
  final Color? iconColor;
  final Color? valueColor;
  final VoidCallback? onTap;

  const StatCard({
    super.key,
    required this.title,
    required this.value,
    this.subtitle,
    this.icon,
    this.iconColor,
    this.valueColor,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.all(AppSpacing.md),
        decoration: BoxDecoration(
          color: AppColors.cardBg,
          borderRadius: AppRadius.cardRadius,
          boxShadow: AppShadows.sm,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                if (icon != null) ...[
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: (iconColor ?? AppColors.primary).withOpacity(0.1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Icon(
                      icon,
                      size: 20,
                      color: iconColor ?? AppColors.primary,
                    ),
                  ),
                  const SizedBox(width: 12),
                ],
                Expanded(
                  child: Text(
                    title,
                    style: AppTypography.body2,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Text(
              value,
              style: AppTypography.amount.copyWith(
                color: valueColor ?? AppColors.textPrimary,
              ),
            ),
            if (subtitle != null) ...[
              const SizedBox(height: 4),
              Text(
                subtitle!,
                style: AppTypography.caption,
              ),
            ],
          ],
        ),
      ),
    );
  }
}
```

### 4.2 交易列表项组件

```dart
// lib/shared/widgets/cards/transaction_item.dart

import 'package:flutter/material.dart';

class TransactionItem extends StatelessWidget {
  final String merchantName;
  final String amount;
  final String time;
  final String type; // 'credit' | 'debit' | 'wechat' | 'alipay'
  final VoidCallback? onTap;

  const TransactionItem({
    super.key,
    required this.merchantName,
    required this.amount,
    required this.time,
    required this.type,
    this.onTap,
  });

  IconData get _icon {
    switch (type) {
      case 'wechat':
        return Icons.wechat;
      case 'alipay':
        return Icons.account_balance_wallet;
      case 'debit':
        return Icons.credit_card;
      default:
        return Icons.credit_card;
    }
  }

  Color get _iconColor {
    switch (type) {
      case 'wechat':
        return const Color(0xFF07C160);
      case 'alipay':
        return const Color(0xFF1677FF);
      default:
        return AppColors.primary;
    }
  }

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: onTap,
      child: Container(
        padding: const EdgeInsets.symmetric(
          horizontal: AppSpacing.md,
          vertical: AppSpacing.sm,
        ),
        child: Row(
          children: [
            // 图标
            Container(
              width: 44,
              height: 44,
              decoration: BoxDecoration(
                color: _iconColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(
                _icon,
                color: _iconColor,
                size: 22,
              ),
            ),
            const SizedBox(width: 12),

            // 信息
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    merchantName,
                    style: AppTypography.body1.copyWith(
                      fontWeight: FontWeight.w500,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 2),
                  Text(
                    time,
                    style: AppTypography.caption,
                  ),
                ],
              ),
            ),

            // 金额
            Text(
              amount,
              style: AppTypography.body1.copyWith(
                fontWeight: FontWeight.w600,
                color: AppColors.textPrimary,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
```

### 4.3 钱包卡片组件

```dart
// lib/shared/widgets/cards/wallet_card.dart

import 'package:flutter/material.dart';

class WalletCard extends StatelessWidget {
  final String walletName;
  final String channelName;
  final String balance;
  final String threshold;
  final bool canWithdraw;
  final VoidCallback? onWithdraw;

  const WalletCard({
    super.key,
    required this.walletName,
    required this.channelName,
    required this.balance,
    required this.threshold,
    required this.canWithdraw,
    this.onWithdraw,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: AppSpacing.cardGap),
      padding: const EdgeInsets.all(AppSpacing.md),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [
            AppColors.primary,
            AppColors.primaryDark,
          ],
        ),
        borderRadius: AppRadius.cardRadius,
        boxShadow: AppShadows.md,
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 头部
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                walletName,
                style: AppTypography.body1.copyWith(
                  color: Colors.white.withOpacity(0.9),
                ),
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.2),
                  borderRadius: BorderRadius.circular(4),
                ),
                child: Text(
                  channelName,
                  style: AppTypography.caption.copyWith(
                    color: Colors.white,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),

          // 余额
          Text(
            '¥ $balance',
            style: const TextStyle(
              fontSize: 32,
              fontWeight: FontWeight.w700,
              color: Colors.white,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            '提现门槛: ¥$threshold',
            style: AppTypography.caption.copyWith(
              color: Colors.white.withOpacity(0.7),
            ),
          ),
          const SizedBox(height: 16),

          // 提现按钮
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: canWithdraw ? onWithdraw : null,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.white,
                foregroundColor: AppColors.primary,
                disabledBackgroundColor: Colors.white.withOpacity(0.5),
                disabledForegroundColor: AppColors.primary.withOpacity(0.5),
              ),
              child: Text(canWithdraw ? '申请提现' : '未达提现门槛'),
            ),
          ),
        ],
      ),
    );
  }
}
```

### 4.4 分润类型标签组件

```dart
// lib/shared/widgets/tags/profit_type_tag.dart

import 'package:flutter/material.dart';

enum ProfitType {
  trade,    // 交易分润
  deposit,  // 押金返现
  sim,      // 流量返现
  reward,   // 激活奖励
}

class ProfitTypeTag extends StatelessWidget {
  final ProfitType type;

  const ProfitTypeTag({super.key, required this.type});

  String get _label {
    switch (type) {
      case ProfitType.trade:
        return '交易分润';
      case ProfitType.deposit:
        return '押金返现';
      case ProfitType.sim:
        return '流量返现';
      case ProfitType.reward:
        return '激活奖励';
    }
  }

  Color get _color {
    switch (type) {
      case ProfitType.trade:
        return AppColors.profitTrade;
      case ProfitType.deposit:
        return AppColors.profitDeposit;
      case ProfitType.sim:
        return AppColors.profitSim;
      case ProfitType.reward:
        return AppColors.profitReward;
    }
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: _color.withOpacity(0.1),
        borderRadius: AppRadius.tagRadius,
      ),
      child: Text(
        _label,
        style: AppTypography.caption.copyWith(
          color: _color,
          fontWeight: FontWeight.w500,
        ),
      ),
    );
  }
}
```

### 4.5 空状态组件

```dart
// lib/shared/widgets/empty/empty_state.dart

import 'package:flutter/material.dart';

class EmptyState extends StatelessWidget {
  final String title;
  final String? description;
  final IconData icon;
  final String? buttonText;
  final VoidCallback? onButtonTap;

  const EmptyState({
    super.key,
    required this.title,
    this.description,
    this.icon = Icons.inbox_outlined,
    this.buttonText,
    this.onButtonTap,
  });

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              icon,
              size: 64,
              color: AppColors.textTertiary,
            ),
            const SizedBox(height: 16),
            Text(
              title,
              style: AppTypography.h3.copyWith(
                color: AppColors.textSecondary,
              ),
              textAlign: TextAlign.center,
            ),
            if (description != null) ...[
              const SizedBox(height: 8),
              Text(
                description!,
                style: AppTypography.body2,
                textAlign: TextAlign.center,
              ),
            ],
            if (buttonText != null && onButtonTap != null) ...[
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: onButtonTap,
                child: Text(buttonText!),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
```

---

## 五、页面设计

### 5.1 首页

```
┌─────────────────────────────────────┐
│ ≡  代理商分润系统           🔔 (3)  │
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐ │
│ │    [轮播图/滚动图]              │ │
│ │    ◉ ○ ○                       │ │
│ └─────────────────────────────────┘ │
│                                     │
│  今日收益                           │
│ ┌───────────────────────────────┐   │
│ │       ¥ 1,234.56              │   │
│ │  较昨日 ↑12.5%                 │   │
│ └───────────────────────────────┘   │
│                                     │
│ ┌───────────┐ ┌───────────┐         │
│ │ 交易分润   │ │ 押金返现   │         │
│ │ ¥856.00  │ │ ¥150.00  │         │
│ └───────────┘ └───────────┘         │
│ ┌───────────┐ ┌───────────┐         │
│ │ 流量返现   │ │ 激活奖励   │         │
│ │ ¥138.56  │ │ ¥90.00   │         │
│ └───────────┘ └───────────┘         │
│                                     │
│  快捷入口                           │
│ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐   │
│ │ 📱  │ │ 👥  │ │ 📊  │ │ 💰  │   │
│ │终端  │ │商户  │ │数据  │ │钱包  │   │
│ └─────┘ └─────┘ └─────┘ └─────┘   │
│ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐   │
│ │ 📤  │ │ 🎫  │ │ 📢  │ │ 👤  │   │
│ │代扣  │ │海报  │ │消息  │ │我的  │   │
│ └─────┘ └─────┘ └─────┘ └─────┘   │
│                                     │
│  最近交易                     查看更多>│
│ ┌───────────────────────────────┐   │
│ │ 商户A    ¥1,000.00   10:30   │   │
│ │ 商户B    ¥2,500.00   10:25   │   │
│ │ 商户C    ¥800.00     10:20   │   │
│ └───────────────────────────────┘   │
│                                     │
├─────────────────────────────────────┤
│  🏠    📱    📊    💰    👤       │
│  首页   终端   数据   钱包   我的    │
└─────────────────────────────────────┘
```

### 5.2 终端管理

```
┌─────────────────────────────────────┐
│ ←  终端管理                          │
├─────────────────────────────────────┤
│                                     │
│  终端统计                           │
│ ┌───────────┐ ┌───────────┐         │
│ │   200     │ │   180     │         │
│ │  终端总数  │ │  已激活    │         │
│ └───────────┘ └───────────┘         │
│ ┌───────────┐ ┌───────────┐         │
│ │    20     │ │     5     │         │
│ │  未激活    │ │  今日激活  │         │
│ └───────────┘ └───────────┘         │
│                                     │
│  [ 全部 | 已激活 | 未激活 | 库存 ]    │
│                                     │
│ ┌───────────────────────────────┐   │
│ │ SN: 12345678                  │   │
│ │ 商户: 张三商店                  │   │
│ │ 状态: ✓ 已激活   激活时间: 1月20日│   │
│ │                    [详情] [设置]│   │
│ └───────────────────────────────┘   │
│                                     │
│ ┌───────────────────────────────┐   │
│ │ SN: 12345679                  │   │
│ │ 商户: -                        │   │
│ │ 状态: ○ 未激活                  │   │
│ │                    [划拨] [回拨]│   │
│ └───────────────────────────────┘   │
│                                     │
│  (更多终端...)                      │
│                                     │
├─────────────────────────────────────┤
│      [划拨]           [回拨]         │
└─────────────────────────────────────┘
```

### 5.3 终端划拨（不可跨级）

```
┌─────────────────────────────────────┐
│ ←  终端划拨                          │
├─────────────────────────────────────┤
│                                     │
│  已选终端: 3台                       │
│  SN: 12345679, 12345680, 12345681   │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  划拨给:                            │
│ ┌───────────────────────────────┐   │
│ │ 🔍 搜索直属下级代理商           │   │
│ └───────────────────────────────┘   │
│                                     │
│  直属下级代理商                      │
│ ┌───────────────────────────────┐   │
│ │ ○ 李四 (A002)                  │   │
│ │   手机: 139****9999            │   │
│ └───────────────────────────────┘   │
│ ┌───────────────────────────────┐   │
│ │ ● 王五 (A003)    ← 已选择      │   │
│ │   手机: 137****7777            │   │
│ └───────────────────────────────┘   │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  ☐ 设置货款代扣                     │
│                                     │
│    单价: ¥ [50] 元/台               │
│    总金额: ¥150                     │
│                                     │
│    扣款来源:                        │
│    ☑ 分润钱包                       │
│    ☐ 服务费钱包                     │
│    ☐ 奖励钱包                       │
│                                     │
│  ⚠️ APP仅支持划拨给直属下级          │
│                                     │
│      [确认划拨]                      │
│                                     │
└─────────────────────────────────────┘
```

### 5.4 货款代扣（独立模块）

```
┌─────────────────────────────────────┐
│ ←  货款代扣                          │
├─────────────────────────────────────┤
│                                     │
│  [ 待接收 | 进行中 | 已完成 ]         │
│                                     │
│  待接收 (2)                          │
│ ┌───────────────────────────────┐   │
│ │ 来自: 总部 (上级)               │   │
│ │ 终端: 10台 × ¥50 = ¥500       │   │
│ │ 扣款来源: 分润钱包              │   │
│ │ 时间: 2024-01-20 10:30        │   │
│ │              [拒绝]  [接收]    │   │
│ └───────────────────────────────┘   │
│                                     │
│  进行中 (1)                          │
│ ┌───────────────────────────────┐   │
│ │ 来自: 总部 (上级)               │   │
│ │ 总金额: ¥1,000                 │   │
│ │ 已扣: ¥350 / 待扣: ¥650        │   │
│ │ ████████░░░░░░░░░ 35%         │   │
│ │ 扣款来源: 分润钱包+服务费钱包    │   │
│ └───────────────────────────────┘   │
│                                     │
│  已完成 (5)                          │
│ ┌───────────────────────────────┐   │
│ │ 来自: 总部 (上级)               │   │
│ │ 总金额: ¥500    已扣完成        │   │
│ │ 完成时间: 2024-01-18           │   │
│ └───────────────────────────────┘   │
│                                     │
│  (更多...)                          │
│                                     │
└─────────────────────────────────────┘
```

### 5.5 商户管理

```
┌─────────────────────────────────────┐
│ ←  商户管理                          │
├─────────────────────────────────────┤
│                                     │
│  [ 直营 | 团队 ]                     │
│                                     │
│ ┌───────────────────────────────┐   │
│ │ 🔍 搜索商户名称/编号/机具号      │   │
│ └───────────────────────────────┘   │
│                                     │
│  直营商户 (45)                       │
│                                     │
│ ┌───────────────────────────────┐   │
│ │ 张三商店                        │   │
│ │ 编号: M001  机具: SN12345678   │   │
│ │ 本月交易: ¥125,000             │   │
│ │ 费率: 0.55%   状态: 活跃        │   │
│ │                        [详情] >│   │
│ └───────────────────────────────┘   │
│                                     │
│ ┌───────────────────────────────┐   │
│ │ 李四超市                        │   │
│ │ 编号: M002  机具: SN12345679   │   │
│ │ 本月交易: ¥86,500              │   │
│ │ 费率: 0.58%   状态: 活跃        │   │
│ │                        [详情] >│   │
│ └───────────────────────────────┘   │
│                                     │
│ ┌───────────────────────────────┐   │
│ │ 王五便利店                      │   │
│ │ 编号: M003  机具: SN12345680   │   │
│ │ 本月交易: ¥0                   │   │
│ │ 费率: 0.60%   状态: ⚠️ 30天无交易│   │
│ │                        [详情] >│   │
│ └───────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
```

### 5.6 商户详情

```
┌─────────────────────────────────────┐
│ ←  商户详情                          │
├─────────────────────────────────────┤
│                                     │
│  张三商店                            │
│  编号: M001                          │
│                                     │
│  基本信息                            │
│ ┌───────────────────────────────┐   │
│ │ 手机号: 138****8888            │   │
│ │ 机具号: SN12345678             │   │
│ │ 激活时间: 2024-01-15           │   │
│ │ 首次流量费: ¥79 (2024-01-15)   │   │
│ └───────────────────────────────┘   │
│                                     │
│  费率设置                            │
│ ┌───────────────────────────────┐   │
│ │ 刷卡费率: 0.55%        [修改]  │   │
│ │ 扫码费率: 0.38%        [修改]  │   │
│ └───────────────────────────────┘   │
│                                     │
│  交易统计                            │
│ ┌───────────────────────────────┐   │
│ │ 累计交易:    ¥1,250,000        │   │
│ │ 本月交易:    ¥125,000          │   │
│ │ ├ 贷记卡:    ¥80,000           │   │
│ │ ├ 借记卡:    ¥30,000           │   │
│ │ ├ 微信:      ¥10,000           │   │
│ │ └ 支付宝:    ¥5,000            │   │
│ └───────────────────────────────┘   │
│                                     │
│  近7天交易趋势                       │
│ ┌───────────────────────────────┐   │
│ │       [折线图]                 │   │
│ └───────────────────────────────┘   │
│                                     │
│  交易记录                     查看更多>│
│ ┌───────────────────────────────┐   │
│ │ 01-20 10:30  刷卡  ¥1,500.00  │   │
│ │ 01-20 09:15  微信  ¥320.00    │   │
│ │ 01-19 18:20  刷卡  ¥2,800.00  │   │
│ └───────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
```

### 5.7 钱包

```
┌─────────────────────────────────────┐
│ ←  我的钱包                          │
├─────────────────────────────────────┤
│                                     │
│  总资产                              │
│ ┌───────────────────────────────┐   │
│ │        ¥ 12,345.67            │   │
│ │      累计提现: ¥88,500.00      │   │
│ └───────────────────────────────┘   │
│                                     │
│  通道筛选: [全部▼]                   │
│                                     │
│  ┌─────────────────────────────┐    │
│  │        分润钱包               │    │
│  │        拉卡拉                │    │
│  │                              │    │
│  │     ¥ 5,680.00              │    │
│  │    提现门槛: ¥100            │    │
│  │                              │    │
│  │    [申请提现]                │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │        服务费钱包             │    │
│  │        拉卡拉                │    │
│  │                              │    │
│  │     ¥ 3,200.00              │    │
│  │    提现门槛: ¥200            │    │
│  │                              │    │
│  │    [申请提现]                │    │
│  └─────────────────────────────┘    │
│                                     │
│  ┌─────────────────────────────┐    │
│  │        奖励钱包               │    │
│  │        拉卡拉                │    │
│  │                              │    │
│  │     ¥ 1,500.00              │    │
│  │    提现门槛: ¥50             │    │
│  │                              │    │
│  │    [申请提现]                │    │
│  └─────────────────────────────┘    │
│                                     │
│  [钱包流水]     [提现记录]           │
│                                     │
└─────────────────────────────────────┘
```

### 5.8 收益统计

```
┌─────────────────────────────────────┐
│ ←  收益统计                          │
├─────────────────────────────────────┤
│                                     │
│  今日收益                            │
│ ┌───────────────────────────────┐   │
│ │        ¥ 1,234.56             │   │
│ │      较昨日 ↑12.5%             │   │
│ └───────────────────────────────┘   │
│                                     │
│  收益明细                            │
│ ┌─────────┐ ┌─────────┐             │
│ │ 交易分润 │ │ 押金返现 │             │
│ │ ¥856.00│ │ ¥150.00│             │
│ └─────────┘ └─────────┘             │
│ ┌─────────┐ ┌─────────┐             │
│ │ 流量返现 │ │ 激活奖励 │             │
│ │ ¥138.56│ │ ¥90.00 │             │
│ └─────────┘ └─────────┘             │
│                                     │
│  收益趋势  [7天] [30天]              │
│ ┌───────────────────────────────┐   │
│ │                               │   │
│ │       [折线图]                │   │
│ │                               │   │
│ └───────────────────────────────┘   │
│                                     │
│  月收益  [近6月] [近1年] [近2年]      │
│ ┌───────────────────────────────┐   │
│ │ 2024-01   ¥32,500.00         │   │
│ │ 2023-12   ¥28,800.00         │   │
│ │ 2023-11   ¥30,200.00         │   │
│ │ 2023-10   ¥26,500.00         │   │
│ │ 2023-09   ¥24,100.00         │   │
│ │ 2023-08   ¥22,800.00         │   │
│ └───────────────────────────────┘   │
│                                     │
└─────────────────────────────────────┘
```

### 5.9 代理拓展

```
┌─────────────────────────────────────┐
│ ←  代理拓展                          │
├─────────────────────────────────────┤
│                                     │
│  我的推广码                          │
│ ┌───────────────────────────────┐   │
│ │                               │   │
│ │        [二维码图片]            │   │
│ │                               │   │
│ │    邀请码: ZHANG001           │   │
│ │                               │   │
│ │  [保存到相册]  [复制链接]       │   │
│ └───────────────────────────────┘   │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  我的团队                            │
│  直属代理: 12    团队代理: 156       │
│                                     │
│  [+ 手动添加代理]                    │
│                                     │
│  直属代理商列表                       │
│ ┌───────────────────────────────┐   │
│ │ 李四 (A002)                    │   │
│ │ 手机: 139****9999              │   │
│ │ 入网: 2024-01-10              │   │
│ │ 下级: 15人   商户: 45个         │   │
│ └───────────────────────────────┘   │
│ ┌───────────────────────────────┐   │
│ │ 王五 (A003)                    │   │
│ │ 手机: 137****7777              │   │
│ │ 入网: 2024-01-05              │   │
│ │ 下级: 8人    商户: 32个         │   │
│ └───────────────────────────────┘   │
│                                     │
│  (更多...)                          │
│                                     │
└─────────────────────────────────────┘
```

### 5.10 我的信息

```
┌─────────────────────────────────────┐
│ ←  我的信息                          │
├─────────────────────────────────────┤
│                                     │
│         ┌─────────┐                 │
│         │  头像    │                 │
│         └─────────┘                 │
│            张三                      │
│         服务商编号: A001             │
│                                     │
│  ─────────────────────────────────  │
│                                     │
│  基本信息                            │
│ ┌───────────────────────────────┐   │
│ │ 姓名        张三               │   │
│ │ 手机号      138****8888        │   │
│ │ 身份证      110***********34   │   │
│ │ 入网时间    2024-01-15         │   │
│ └───────────────────────────────┘   │
│                                     │
│  结算信息                            │
│ ┌───────────────────────────────┐   │
│ │ 开户行      中国银行            │   │
│ │ 银行卡号    ****5678           │   │
│ │                        [更改] >│   │
│ └───────────────────────────────┘   │
│                                     │
│  费率成本                            │
│ ┌───────────────────────────────┐   │
│ │ 拉卡拉                         │   │
│ │ 贷记卡: 0.51%  借记卡: 0.51%   │   │
│ ├───────────────────────────────┤   │
│ │ 随行付                         │   │
│ │ 贷记卡: 0.52%  借记卡: 0.52%   │   │
│ └───────────────────────────────┘   │
│                                     │
│  我的邀请码                          │
│ ┌───────────────────────────────┐   │
│ │ ZHANG001          [自定义靓号]>│   │
│ └───────────────────────────────┘   │
│                                     │
│  [退出登录]                          │
│                                     │
└─────────────────────────────────────┘
```

### 5.11 消息通知

```
┌─────────────────────────────────────┐
│ ←  消息通知                          │
├─────────────────────────────────────┤
│                                     │
│  [ 全部 | 分润 | 注册 | 消费 | 系统 ] │
│                                     │
│  今天                                │
│ ┌───────────────────────────────┐   │
│ │ 💰 分润到账                    │   │
│ │ 您有一笔¥8.00的交易分润已入账   │   │
│ │ 今天 10:30              ● 未读│   │
│ └───────────────────────────────┘   │
│ ┌───────────────────────────────┐   │
│ │ 👤 新代理商注册                │   │
│ │ 李四(139****9999)已注册成为... │   │
│ │ 今天 09:15              ● 未读│   │
│ └───────────────────────────────┘   │
│                                     │
│  昨天                                │
│ ┌───────────────────────────────┐   │
│ │ 💳 交易通知                    │   │
│ │ 商户"张三商店"完成一笔¥1,500...│   │
│ │ 昨天 18:20                    │   │
│ └───────────────────────────────┘   │
│                                     │
│  更早                                │
│ ┌───────────────────────────────┐   │
│ │ 📢 系统公告                    │   │
│ │ 系统将于1月25日进行升级维护...  │   │
│ │ 3天前                         │   │
│ └───────────────────────────────┘   │
│                                     │
│  ⚠️ 消息3天后自动过期               │
│                                     │
└─────────────────────────────────────┘
```

### 5.12 营销海报

```
┌─────────────────────────────────────┐
│ ←  营销海报                          │
├─────────────────────────────────────┤
│                                     │
│  分类: [全部▼]                       │
│                                     │
│ ┌───────────────┐ ┌───────────────┐ │
│ │               │ │               │ │
│ │  [海报图片1]   │ │  [海报图片2]   │ │
│ │               │ │               │ │
│ │    新年活动    │ │    招商合作    │ │
│ │   [保存]      │ │   [保存]      │ │
│ └───────────────┘ └───────────────┘ │
│                                     │
│ ┌───────────────┐ ┌───────────────┐ │
│ │               │ │               │ │
│ │  [海报图片3]   │ │  [海报图片4]   │ │
│ │               │ │               │ │
│ │    优惠活动    │ │    产品介绍    │ │
│ │   [保存]      │ │   [保存]      │ │
│ └───────────────┘ └───────────────┘ │
│                                     │
│ ┌───────────────┐ ┌───────────────┐ │
│ │               │ │               │ │
│ │  [海报图片5]   │ │  [海报图片6]   │ │
│ │               │ │               │ │
│ │    加盟政策    │ │    品牌宣传    │ │
│ │   [保存]      │ │   [保存]      │ │
│ └───────────────┘ └───────────────┘ │
│                                     │
└─────────────────────────────────────┘
```

---

## 六、多平台适配

### 6.1 平台判断工具

```dart
// lib/core/utils/platform_utils.dart

import 'dart:io';

class PlatformUtils {
  static bool get isIOS => Platform.isIOS;
  static bool get isAndroid => Platform.isAndroid;
  static bool get isHarmonyOS {
    // HarmonyOS 设备标识判断
    return Platform.operatingSystem == 'harmonyos' ||
           Platform.environment.containsKey('HARMONYOS_VERSION');
  }

  static String get platformName {
    if (isIOS) return 'iOS';
    if (isHarmonyOS) return 'HarmonyOS';
    if (isAndroid) return 'Android';
    return 'Unknown';
  }
}
```

### 6.2 推送服务适配

```dart
// lib/core/services/push_service.dart

abstract class PushService {
  Future<void> init();
  Future<String?> getToken();
  void onMessageReceived(Function(Map<String, dynamic>) callback);
}

// iOS - APNs
class IOSPushService implements PushService {
  @override
  Future<void> init() async {
    // 初始化 APNs
  }

  @override
  Future<String?> getToken() async {
    // 获取 APNs Token
    return null;
  }

  @override
  void onMessageReceived(Function(Map<String, dynamic>) callback) {
    // 监听消息
  }
}

// Android - FCM 或 极光推送
class AndroidPushService implements PushService {
  @override
  Future<void> init() async {
    // 初始化 FCM / 极光
  }

  @override
  Future<String?> getToken() async {
    return null;
  }

  @override
  void onMessageReceived(Function(Map<String, dynamic>) callback) {}
}

// HarmonyOS - 华为 Push Kit
class HarmonyPushService implements PushService {
  @override
  Future<void> init() async {
    // 初始化华为 Push Kit
  }

  @override
  Future<String?> getToken() async {
    return null;
  }

  @override
  void onMessageReceived(Function(Map<String, dynamic>) callback) {}
}

// 工厂方法
PushService createPushService() {
  if (PlatformUtils.isIOS) {
    return IOSPushService();
  } else if (PlatformUtils.isHarmonyOS) {
    return HarmonyPushService();
  } else {
    return AndroidPushService();
  }
}
```

---

## 七、工具函数

### 7.1 金额格式化

```dart
// lib/core/utils/format_utils.dart

class FormatUtils {
  /// 格式化金额 (分 -> 元)
  static String formatAmount(int? cents, {bool showSign = false}) {
    if (cents == null) return '¥0.00';
    final yuan = cents / 100;
    final sign = showSign && yuan > 0 ? '+' : '';
    return '$sign¥${yuan.toStringAsFixed(2)}';
  }

  /// 格式化金额 (元)
  static String formatYuan(double? yuan, {bool showSign = false}) {
    if (yuan == null) return '¥0.00';
    final sign = showSign && yuan > 0 ? '+' : '';
    return '$sign¥${yuan.toStringAsFixed(2)}';
  }

  /// 格式化大金额 (万)
  static String formatLargeAmount(double yuan) {
    if (yuan >= 10000) {
      return '¥${(yuan / 10000).toStringAsFixed(2)}万';
    }
    return '¥${yuan.toStringAsFixed(2)}';
  }

  /// 格式化费率
  static String formatRate(double rate) {
    return '${(rate * 100).toStringAsFixed(2)}%';
  }

  /// 格式化手机号脱敏
  static String maskPhone(String phone) {
    if (phone.length != 11) return phone;
    return '${phone.substring(0, 3)}****${phone.substring(7)}';
  }

  /// 格式化身份证脱敏
  static String maskIdCard(String idCard) {
    if (idCard.length != 18) return idCard;
    return '${idCard.substring(0, 3)}***********${idCard.substring(14)}';
  }

  /// 格式化银行卡脱敏
  static String maskBankCard(String cardNo) {
    if (cardNo.length < 4) return cardNo;
    return '****${cardNo.substring(cardNo.length - 4)}';
  }

  /// 格式化日期
  static String formatDate(DateTime date, {String pattern = 'yyyy-MM-dd'}) {
    // 使用 intl 包的 DateFormat
    return '${date.year}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}';
  }

  /// 格式化相对时间
  static String formatRelativeTime(DateTime date) {
    final now = DateTime.now();
    final diff = now.difference(date);

    if (diff.inDays == 0) {
      if (diff.inHours == 0) {
        return '${diff.inMinutes}分钟前';
      }
      return '${diff.inHours}小时前';
    } else if (diff.inDays == 1) {
      return '昨天';
    } else if (diff.inDays < 7) {
      return '${diff.inDays}天前';
    } else {
      return formatDate(date);
    }
  }
}
```

---

## 八、开发流程（单人+AI）

### 8.1 工作分工

| 工作内容 | 负责方 | 说明 |
|----------|--------|------|
| Figma 线框图 | 👤 人工 | 快速勾勒页面布局 |
| 设计系统定义 | 🤖 AI | 颜色、字体、间距规范 |
| 组件代码生成 | 🤖 AI | Flutter 组件库代码 |
| 页面布局代码 | 🤖 AI | 基于线框图生成 |
| 业务逻辑 | 👤 人工 + 🤖 AI | 协作完成 |
| API 对接 | 🤖 AI | 生成模型和请求代码 |
| 样式微调 | 👤 人工 | 细节打磨 |
| 测试 | 👤 人工 | 功能验收 |

### 8.2 开发阶段

| 阶段 | 内容 | 工期 |
|------|------|------|
| **Phase 1** | 项目搭建、设计系统、组件库 | 1周 |
| **Phase 2** | 认证模块、首页 | 1周 |
| **Phase 3** | 终端管理、货款代扣 | 1.5周 |
| **Phase 4** | 商户管理、数据分析 | 1.5周 |
| **Phase 5** | 钱包、收益统计 | 1周 |
| **Phase 6** | 代理拓展、消息、营销 | 1周 |
| **Phase 7** | 我的信息、设置 | 0.5周 |
| **Phase 8** | 联调测试、优化 | 2周 |

**总计**: 约10周

### 8.3 Figma 工作流

1. **创建组件库**
   - 按钮（主要、次要、文本）
   - 卡片（统计、列表项、钱包）
   - 输入框、选择器
   - 标签、徽章
   - 图表占位符

2. **设计页面流程**
   - 先画主流程页面
   - 使用组件拼装
   - 标注交互说明

3. **交付给AI**
   - 截图页面布局
   - 描述交互逻辑
   - AI 生成代码

---

## 九、关键依赖包

```yaml
# pubspec.yaml

dependencies:
  flutter:
    sdk: flutter

  # 状态管理
  flutter_riverpod: ^2.4.0

  # 路由
  go_router: ^12.0.0

  # 网络请求
  dio: ^5.3.0

  # 本地存储
  shared_preferences: ^2.2.0
  sqflite: ^2.3.0

  # 图表
  fl_chart: ^0.64.0

  # 图片
  cached_network_image: ^3.3.0
  image_gallery_saver: ^2.0.3

  # 二维码
  qr_flutter: ^4.1.0

  # 刷新
  pull_to_refresh: ^2.0.0

  # 工具
  intl: ^0.18.0
  url_launcher: ^6.2.0
  package_info_plus: ^5.0.1

  # 推送 (按平台选择)
  firebase_messaging: ^14.7.0  # Android
  flutter_local_notifications: ^16.1.0

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^3.0.0
  build_runner: ^2.4.0
  json_serializable: ^6.7.0
```

---

*文档版本: v1.0*
*最后更新: 2025-01-18*
