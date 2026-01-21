import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/models/merchant_model.dart';
import '../../data/repositories/merchant_repository.dart';

/// 商户仓储 Provider
final merchantRepositoryProvider = Provider<MerchantRepository>((ref) {
  return MerchantRepository();
});

/// 商户统计 Provider
final merchantStatsProvider = FutureProvider<MerchantStats>((ref) async {
  final repository = ref.watch(merchantRepositoryProvider);
  return repository.getMerchantStats();
});

/// 商户列表状态
class MerchantListState {
  final List<Merchant> merchants;
  final bool isLoading;
  final bool hasMore;
  final int page;
  final String? error;

  MerchantListState({
    this.merchants = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.page = 1,
    this.error,
  });

  MerchantListState copyWith({
    List<Merchant>? merchants,
    bool? isLoading,
    bool? hasMore,
    int? page,
    String? error,
  }) {
    return MerchantListState(
      merchants: merchants ?? this.merchants,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      page: page ?? this.page,
      error: error,
    );
  }
}

/// 商户列表 Notifier
class MerchantListNotifier extends StateNotifier<MerchantListState> {
  final MerchantRepository _repository;
  final bool isDirect;
  String? _keyword;
  String? _merchantType;

  MerchantListNotifier(this._repository, this.isDirect)
      : super(MerchantListState());

  /// 加载商户列表
  Future<void> loadMerchants({bool refresh = false}) async {
    if (state.isLoading) return;

    if (refresh) {
      state = state.copyWith(page: 1, hasMore: true, merchants: []);
    }

    if (!state.hasMore && !refresh) return;

    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _repository.getMerchants(
        isDirect: isDirect,
        keyword: _keyword,
        merchantType: _merchantType,
        page: state.page,
        pageSize: 20,
      );

      final newMerchants = refresh
          ? response.list
          : [...state.merchants, ...response.list];

      state = state.copyWith(
        merchants: newMerchants,
        isLoading: false,
        hasMore: response.hasMore,
        page: state.page + 1,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 刷新列表
  Future<void> refresh() async {
    await loadMerchants(refresh: true);
  }

  /// 设置搜索关键词
  void setKeyword(String? keyword) {
    _keyword = keyword;
    loadMerchants(refresh: true);
  }

  /// 设置商户类型筛选
  void setMerchantType(String? type) {
    _merchantType = type;
    loadMerchants(refresh: true);
  }
}

/// 直营商户列表 Provider
final directMerchantListProvider =
    StateNotifierProvider<MerchantListNotifier, MerchantListState>((ref) {
  final repository = ref.watch(merchantRepositoryProvider);
  return MerchantListNotifier(repository, true)..loadMerchants();
});

/// 团队商户列表 Provider
final teamMerchantListProvider =
    StateNotifierProvider<MerchantListNotifier, MerchantListState>((ref) {
  final repository = ref.watch(merchantRepositoryProvider);
  return MerchantListNotifier(repository, false)..loadMerchants();
});

/// 商户详情 Provider
final merchantDetailProvider =
    FutureProvider.family<MerchantDetail, int>((ref, id) async {
  final repository = ref.watch(merchantRepositoryProvider);
  return repository.getMerchantDetail(id);
});
