import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../data/home_api.dart';
import '../../domain/home_model.dart';

/// 首页状态
class HomeState {
  final bool isLoading;
  final String? error;
  final HomeOverviewData? overview;
  final List<RecentTransaction> recentTransactions;
  final String scope; // 'direct' 或 'team'

  HomeState({
    this.isLoading = false,
    this.error,
    this.overview,
    this.recentTransactions = const [],
    this.scope = 'direct',
  });

  HomeState copyWith({
    bool? isLoading,
    String? error,
    HomeOverviewData? overview,
    List<RecentTransaction>? recentTransactions,
    String? scope,
  }) {
    return HomeState(
      isLoading: isLoading ?? this.isLoading,
      error: error,
      overview: overview ?? this.overview,
      recentTransactions: recentTransactions ?? this.recentTransactions,
      scope: scope ?? this.scope,
    );
  }
}

/// 首页状态管理
class HomeNotifier extends StateNotifier<HomeState> {
  final HomeApi _api;

  HomeNotifier(this._api) : super(HomeState()) {
    // 初始化时加载数据
    loadData();
  }

  /// 加载首页数据
  Future<void> loadData() async {
    state = state.copyWith(isLoading: true, error: null);

    try {
      // 并行加载概览数据和最近交易
      final results = await Future.wait([
        _api.getOverview(scope: state.scope),
        _api.getRecentTransactions(limit: 5),
      ]);

      state = state.copyWith(
        isLoading: false,
        overview: results[0] as HomeOverviewData,
        recentTransactions: results[1] as List<RecentTransaction>,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 刷新数据
  Future<void> refresh() async {
    await loadData();
  }

  /// 切换统计范围
  Future<void> switchScope(String scope) async {
    if (state.scope == scope) return;
    state = state.copyWith(scope: scope);
    await loadData();
  }
}

/// 首页状态Provider
final homeProvider = StateNotifierProvider<HomeNotifier, HomeState>((ref) {
  final api = ref.watch(homeApiProvider);
  return HomeNotifier(api);
});

/// 今日收益(元)
final todayProfitProvider = Provider<double>((ref) {
  final state = ref.watch(homeProvider);
  return state.overview?.today.profitTotalYuan ?? 0;
});

/// 较昨日变化率
final profitChangeRateProvider = Provider<double>((ref) {
  final state = ref.watch(homeProvider);
  return state.overview?.profitChangeRate ?? 0;
});

/// 是否增长
final isProfitGrowthProvider = Provider<bool>((ref) {
  final state = ref.watch(homeProvider);
  return state.overview?.isProfitGrowth ?? true;
});
