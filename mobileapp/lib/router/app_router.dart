import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../features/auth/presentation/login_page.dart';
import '../features/home/presentation/home_page.dart';
import '../features/terminal/presentation/terminal_page.dart';
import '../features/terminal/presentation/terminal_transfer_page.dart';
import '../features/cargo_deduction/presentation/cargo_deduction_page.dart';
import '../features/merchant/presentation/merchant_page.dart';
import '../features/merchant/presentation/merchant_detail_page.dart';
import '../features/data_analysis/presentation/data_analysis_page.dart';
import '../features/profit/presentation/profit_page.dart';
import '../features/wallet/presentation/wallet_page.dart';
import '../features/wallet/presentation/withdraw_page.dart';
import '../features/agent/presentation/agent_page.dart';
import '../features/deduction/presentation/deduction_page.dart';
import '../features/marketing/presentation/marketing_page.dart';
import '../features/message/presentation/message_page.dart';
import '../features/profile/presentation/profile_page.dart';
import '../shared/widgets/main_scaffold.dart';

/// 路由路径常量
class RoutePaths {
  RoutePaths._();

  // 认证
  static const String login = '/login';
  static const String register = '/register';

  // 主页面（底部导航）
  static const String home = '/';
  static const String terminal = '/terminal';
  static const String dataAnalysis = '/data-analysis';
  static const String wallet = '/wallet';
  static const String profile = '/profile';

  // 终端管理
  static const String terminalTransfer = '/terminal/transfer';
  static const String terminalRecall = '/terminal/recall';
  static const String terminalDetail = '/terminal/:id';

  // 货款代扣
  static const String cargoDeduction = '/cargo-deduction';

  // 商户管理
  static const String merchant = '/merchant';
  static const String merchantDetail = '/merchant/:id';

  // 收益统计
  static const String profit = '/profit';

  // 钱包
  static const String withdraw = '/wallet/withdraw';
  static const String walletFlow = '/wallet/flow';
  static const String withdrawRecord = '/wallet/withdraw-record';

  // 代理拓展
  static const String agent = '/agent';
  static const String agentAdd = '/agent/add';
  static const String agentDetail = '/agent/:id';

  // 代扣管理
  static const String deduction = '/deduction';

  // 营销海报
  static const String marketing = '/marketing';

  // 消息通知
  static const String message = '/message';
  static const String messageDetail = '/message/:id';

  // 设置
  static const String settings = '/settings';
  static const String bankCard = '/settings/bank-card';
  static const String inviteCode = '/settings/invite-code';
}

/// 路由配置Provider
final appRouterProvider = Provider<GoRouter>((ref) {
  return GoRouter(
    initialLocation: RoutePaths.home,
    debugLogDiagnostics: true,

    // 重定向逻辑（检查登录状态）
    redirect: (context, state) {
      // TODO: 检查登录状态
      // final isLoggedIn = ref.read(authStateProvider);
      // if (!isLoggedIn && !state.matchedLocation.startsWith('/login')) {
      //   return RoutePaths.login;
      // }
      return null;
    },

    routes: [
      // ==================== 认证页面 ====================
      GoRoute(
        path: RoutePaths.login,
        name: 'login',
        builder: (context, state) => const LoginPage(),
      ),

      // ==================== 主框架（底部导航） ====================
      StatefulShellRoute.indexedStack(
        builder: (context, state, navigationShell) {
          return MainScaffold(navigationShell: navigationShell);
        },
        branches: [
          // 首页
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: RoutePaths.home,
                name: 'home',
                builder: (context, state) => const HomePage(),
              ),
            ],
          ),
          // 终端
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: RoutePaths.terminal,
                name: 'terminal',
                builder: (context, state) => const TerminalPage(),
              ),
            ],
          ),
          // 数据
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: RoutePaths.dataAnalysis,
                name: 'dataAnalysis',
                builder: (context, state) => const DataAnalysisPage(),
              ),
            ],
          ),
          // 钱包
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: RoutePaths.wallet,
                name: 'wallet',
                builder: (context, state) => const WalletPage(),
              ),
            ],
          ),
          // 我的
          StatefulShellBranch(
            routes: [
              GoRoute(
                path: RoutePaths.profile,
                name: 'profile',
                builder: (context, state) => const ProfilePage(),
              ),
            ],
          ),
        ],
      ),

      // ==================== 终端管理 ====================
      GoRoute(
        path: RoutePaths.terminalTransfer,
        name: 'terminalTransfer',
        builder: (context, state) {
          final snList = state.extra as List<String>? ?? [];
          return TerminalTransferPage(selectedSNs: snList);
        },
      ),

      // ==================== 货款代扣 ====================
      GoRoute(
        path: RoutePaths.cargoDeduction,
        name: 'cargoDeduction',
        builder: (context, state) => const CargoDeductionPage(),
      ),

      // ==================== 商户管理 ====================
      GoRoute(
        path: RoutePaths.merchant,
        name: 'merchant',
        builder: (context, state) => const MerchantPage(),
      ),
      GoRoute(
        path: RoutePaths.merchantDetail,
        name: 'merchantDetail',
        builder: (context, state) {
          final id = state.pathParameters['id'] ?? '';
          return MerchantDetailPage(merchantId: id);
        },
      ),

      // ==================== 收益统计 ====================
      GoRoute(
        path: RoutePaths.profit,
        name: 'profit',
        builder: (context, state) => const ProfitPage(),
      ),

      // ==================== 钱包相关 ====================
      GoRoute(
        path: RoutePaths.withdraw,
        name: 'withdraw',
        builder: (context, state) {
          final walletId = state.extra as String? ?? '';
          return WithdrawPage(walletId: walletId);
        },
      ),

      // ==================== 代理拓展 ====================
      GoRoute(
        path: RoutePaths.agent,
        name: 'agent',
        builder: (context, state) => const AgentPage(),
      ),

      // ==================== 代扣管理 ====================
      GoRoute(
        path: RoutePaths.deduction,
        name: 'deduction',
        builder: (context, state) => const DeductionPage(),
      ),

      // ==================== 营销海报 ====================
      GoRoute(
        path: RoutePaths.marketing,
        name: 'marketing',
        builder: (context, state) => const MarketingPage(),
      ),

      // ==================== 消息通知 ====================
      GoRoute(
        path: RoutePaths.message,
        name: 'message',
        builder: (context, state) => const MessagePage(),
      ),
    ],

    // 错误页面
    errorBuilder: (context, state) => Scaffold(
      appBar: AppBar(title: const Text('页面未找到')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.error_outline, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            Text('页面不存在: ${state.matchedLocation}'),
            const SizedBox(height: 16),
            ElevatedButton(
              onPressed: () => context.go(RoutePaths.home),
              child: const Text('返回首页'),
            ),
          ],
        ),
      ),
    ),
  );
});
