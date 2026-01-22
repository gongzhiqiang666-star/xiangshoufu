import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/models/wallet_model.dart';
import '../../data/services/wallet_service.dart';

/// 钱包汇总Provider
final walletSummaryProvider = FutureProvider<WalletSummaryModel>((ref) async {
  final walletService = ref.watch(walletServiceProvider);
  return walletService.getWalletSummary();
});

/// 钱包列表Provider
final walletsProvider = FutureProvider<List<WalletModel>>((ref) async {
  final walletService = ref.watch(walletServiceProvider);
  return walletService.getWallets();
});

/// 钱包配置Provider
final walletConfigProvider = FutureProvider<AgentWalletConfigModel>((ref) async {
  final walletService = ref.watch(walletServiceProvider);
  return walletService.getMyWalletConfig();
});

/// 充值钱包汇总Provider
final chargingWalletSummaryProvider = FutureProvider<ChargingWalletSummaryModel>((ref) async {
  final walletService = ref.watch(walletServiceProvider);
  return walletService.getChargingWalletSummary();
});

/// 沉淀钱包汇总Provider
final settlementWalletSummaryProvider = FutureProvider<SettlementWalletSummaryModel>((ref) async {
  final walletService = ref.watch(walletServiceProvider);
  return walletService.getSettlementWalletSummary();
});

/// 下级余额列表Provider
final subordinateBalancesProvider = FutureProvider<List<SubordinateBalanceModel>>((ref) async {
  final walletService = ref.watch(walletServiceProvider);
  return walletService.getSubordinateBalances();
});

/// 钱包流水查询参数
class WalletLogsParams {
  final int? walletId;
  final String? type;
  final String? startDate;
  final String? endDate;
  final int page;
  final int pageSize;

  WalletLogsParams({
    this.walletId,
    this.type,
    this.startDate,
    this.endDate,
    this.page = 1,
    this.pageSize = 20,
  });

  WalletLogsParams copyWith({
    int? walletId,
    String? type,
    String? startDate,
    String? endDate,
    int? page,
    int? pageSize,
  }) {
    return WalletLogsParams(
      walletId: walletId ?? this.walletId,
      type: type ?? this.type,
      startDate: startDate ?? this.startDate,
      endDate: endDate ?? this.endDate,
      page: page ?? this.page,
      pageSize: pageSize ?? this.pageSize,
    );
  }
}

/// 钱包流水查询参数Provider
final walletLogsParamsProvider = StateProvider<WalletLogsParams>((ref) {
  return WalletLogsParams();
});

/// 钱包流水Provider
final walletLogsProvider = FutureProvider.family<List<WalletLogModel>, WalletLogsParams>((ref, params) async {
  final walletService = ref.watch(walletServiceProvider);
  final response = await walletService.getWalletLogs(
    walletId: params.walletId,
    type: params.type,
    startDate: params.startDate,
    endDate: params.endDate,
    page: params.page,
    pageSize: params.pageSize,
  );
  return response.list;
});

/// 沉淀使用记录Provider
final settlementUsagesProvider = FutureProvider.family<List<SettlementUsageModel>, int?>((ref, usageType) async {
  final walletService = ref.watch(walletServiceProvider);
  final response = await walletService.getSettlementUsageList(
    usageType: usageType,
    page: 1,
    pageSize: 50,
  );
  return response.list;
});
