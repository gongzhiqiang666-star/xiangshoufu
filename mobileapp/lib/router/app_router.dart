import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../features/auth/presentation/login_page.dart';
import '../features/auth/presentation/providers/auth_provider.dart';
import '../features/home/presentation/home_page.dart';
import '../features/terminal/presentation/terminal_page.dart';
import '../features/terminal/presentation/terminal_transfer_page.dart';
import '../features/terminal/presentation/terminal_detail_page.dart';
import '../features/terminal/presentation/terminal_recall_page.dart';
import '../features/cargo_deduction/presentation/cargo_deduction_page.dart';
import '../features/merchant/presentation/merchant_page.dart';
import '../features/merchant/presentation/merchant_detail_page.dart';
import '../features/data_analysis/presentation/data_analysis_page.dart';
import '../features/profit/presentation/profit_page.dart';
import '../features/wallet/presentation/wallet_page.dart';
import '../features/wallet/presentation/withdraw_page.dart';
import '../features/wallet/presentation/settlement_wallet_page.dart';
import '../features/wallet/presentation/charging_wallet_page.dart';
import '../features/wallet/presentation/issue_reward_page.dart';
import '../features/agent/presentation/agent_page.dart';
import '../features/agent/presentation/agent_channels_page.dart';
import '../features/policy/presentation/my_policy_page.dart';
import '../features/policy/presentation/subordinate_policy_page.dart';
import '../features/deduction/presentation/deduction_page.dart';
import '../features/deduction/presentation/deduction_detail_page.dart';
import '../features/goods_deduction/presentation/goods_deduction_page.dart';
import '../features/goods_deduction/presentation/goods_deduction_detail_page.dart';
import '../features/marketing/presentation/marketing_page.dart';
import '../features/message/presentation/message_page.dart';
import '../features/profile/presentation/profile_page.dart';
import '../features/profile/presentation/settings_page.dart';
import '../features/profile/presentation/bank_card_edit_page.dart';
import '../features/profile/presentation/invite_code_page.dart';
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
  static const String chargingWallet = '/wallet/charging';
  static const String settlementWallet = '/wallet/settlement';
  static const String issueReward = '/wallet/issue-reward';

  // 代理拓展
  static const String agent = '/agent';
  static const String agentAdd = '/agent/add';
  static const String agentDetail = '/agent/:id';
  static const String agentChannels = '/agent/:id/channels';
  static const String agentPolicy = '/agent/:id/policy';

  // 政策管理
  static const String myPolicy = '/policy/my';
  static const String subordinatePolicy = '/policy/subordinate/:id';

  // 代扣管理
  static const String deduction = '/deduction';
  static const String deductionDetail = '/deduction/:id';

  // 货款代扣（新版）
  static const String goodsDeduction = '/goods-deduction';
  static const String goodsDeductionDetail = '/goods-deduction/:id';

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
  final authState = ref.watch(authStateProvider);

  return GoRouter(
    initialLocation: RoutePaths.home,
    debugLogDiagnostics: true,

    // 重定向逻辑（检查登录状态）
    redirect: (context, state) {
      final isLoggedIn = authState.isAuthenticated;
      final isLoggingIn = state.matchedLocation == RoutePaths.login;

      // 未登录且不在登录页，跳转到登录页
      if (!isLoggedIn && !isLoggingIn) {
        return RoutePaths.login;
      }

      // 已登录但在登录页，跳转到首页
      if (isLoggedIn && isLoggingIn) {
        return RoutePaths.home;
      }

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
      GoRoute(
        path: RoutePaths.terminalRecall,
        name: 'terminalRecall',
        builder: (context, state) {
          final snList = state.extra as List<String>? ?? [];
          return TerminalRecallPage(selectedSNs: snList);
        },
      ),
      GoRoute(
        path: RoutePaths.terminalDetail,
        name: 'terminalDetail',
        builder: (context, state) {
          final id = state.pathParameters['id'] ?? '';
          return TerminalDetailPage(terminalId: id);
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
          final idStr = state.pathParameters['id'] ?? '0';
          final id = int.tryParse(idStr) ?? 0;
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
      GoRoute(
        path: RoutePaths.chargingWallet,
        name: 'chargingWallet',
        builder: (context, state) => const ChargingWalletPage(),
      ),
      GoRoute(
        path: RoutePaths.settlementWallet,
        name: 'settlementWallet',
        builder: (context, state) => const SettlementWalletPage(),
      ),
      GoRoute(
        path: RoutePaths.issueReward,
        name: 'issueReward',
        builder: (context, state) => const IssueRewardPage(),
      ),

      // ==================== 代理拓展 ====================
      GoRoute(
        path: RoutePaths.agent,
        name: 'agent',
        builder: (context, state) => const AgentPage(),
      ),
      GoRoute(
        path: RoutePaths.agentChannels,
        name: 'agentChannels',
        builder: (context, state) {
          final idStr = state.pathParameters['id'];
          final agentId = idStr != null ? int.tryParse(idStr) : null;
          return AgentChannelsPage(agentId: agentId);
        },
      ),
      GoRoute(
        path: RoutePaths.agentPolicy,
        name: 'agentPolicy',
        builder: (context, state) {
          final idStr = state.pathParameters['id'] ?? '0';
          final agentId = int.tryParse(idStr) ?? 0;
          final extra = state.extra as Map<String, dynamic>?;
          final agentName = extra?['name'] as String? ?? '下级代理商';
          final channelId = extra?['channelId'] as int?;
          return SubordinatePolicyPage(
            subordinateId: agentId,
            subordinateName: agentName,
            initialChannelId: channelId,
          );
        },
      ),

      // ==================== 政策管理 ====================
      GoRoute(
        path: RoutePaths.myPolicy,
        name: 'myPolicy',
        builder: (context, state) => const MyPolicyPage(),
      ),
      GoRoute(
        path: RoutePaths.subordinatePolicy,
        name: 'subordinatePolicy',
        builder: (context, state) {
          final idStr = state.pathParameters['id'] ?? '0';
          final subordinateId = int.tryParse(idStr) ?? 0;
          final extra = state.extra as Map<String, dynamic>?;
          final name = extra?['name'] as String? ?? '下级代理商';
          final channelId = extra?['channelId'] as int?;
          return SubordinatePolicyPage(
            subordinateId: subordinateId,
            subordinateName: name,
            initialChannelId: channelId,
          );
        },
      ),

      // ==================== 代扣管理 ====================
      GoRoute(
        path: RoutePaths.deduction,
        name: 'deduction',
        builder: (context, state) => const DeductionPage(),
      ),
      GoRoute(
        path: RoutePaths.deductionDetail,
        name: 'deductionDetail',
        builder: (context, state) {
          final idStr = state.pathParameters['id'] ?? '0';
          final id = int.tryParse(idStr) ?? 0;
          return DeductionDetailPage(id: id);
        },
      ),

      // ==================== 货款代扣 ====================
      GoRoute(
        path: RoutePaths.goodsDeduction,
        name: 'goodsDeduction',
        builder: (context, state) => const GoodsDeductionPage(),
      ),
      GoRoute(
        path: RoutePaths.goodsDeductionDetail,
        name: 'goodsDeductionDetail',
        builder: (context, state) {
          final idStr = state.pathParameters['id'] ?? '0';
          final id = int.tryParse(idStr) ?? 0;
          final extra = state.extra as Map<String, dynamic>?;
          final isSent = extra?['isSent'] as bool? ?? true;
          return GoodsDeductionDetailPage(id: id, isSent: isSent);
        },
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

      // ==================== 设置相关 ====================
      GoRoute(
        path: RoutePaths.settings,
        name: 'settings',
        builder: (context, state) => const SettingsPage(),
      ),
      GoRoute(
        path: RoutePaths.bankCard,
        name: 'bankCard',
        builder: (context, state) => const BankCardEditPage(),
      ),
      GoRoute(
        path: RoutePaths.inviteCode,
        name: 'inviteCode',
        builder: (context, state) => const InviteCodePage(),
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
