import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/goods_deduction_model.dart';
import '../../data/services/goods_deduction_service.dart';

/// 货款代扣服务 Provider
final goodsDeductionServiceProvider = Provider<GoodsDeductionService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return GoodsDeductionService(apiClient);
});

/// 货款代扣统计 Provider
final goodsDeductionSummaryProvider = FutureProvider.family<GoodsDeductionSummary, String?>((ref, type) async {
  final service = ref.watch(goodsDeductionServiceProvider);
  return service.getSummary(type: type);
});

/// 货款代扣详情 Provider
final goodsDeductionDetailProvider = FutureProvider.family<GoodsDeductionDetail, int>((ref, id) async {
  final service = ref.watch(goodsDeductionServiceProvider);
  return service.getDeductionDetail(id);
});

/// 发起的货款代扣列表状态
class SentDeductionsState {
  final List<GoodsDeduction> list;
  final int total;
  final bool isLoading;
  final bool isLoadingMore;
  final String? error;
  final int currentPage;
  final int? statusFilter;
  final bool hasMore;

  SentDeductionsState({
    this.list = const [],
    this.total = 0,
    this.isLoading = false,
    this.isLoadingMore = false,
    this.error,
    this.currentPage = 1,
    this.statusFilter,
    this.hasMore = true,
  });

  SentDeductionsState copyWith({
    List<GoodsDeduction>? list,
    int? total,
    bool? isLoading,
    bool? isLoadingMore,
    String? error,
    int? currentPage,
    int? statusFilter,
    bool? hasMore,
  }) {
    return SentDeductionsState(
      list: list ?? this.list,
      total: total ?? this.total,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      error: error,
      currentPage: currentPage ?? this.currentPage,
      statusFilter: statusFilter ?? this.statusFilter,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 发起的货款代扣列表 Notifier
class SentDeductionsNotifier extends StateNotifier<SentDeductionsState> {
  final GoodsDeductionService _service;

  SentDeductionsNotifier(this._service) : super(SentDeductionsState());

  /// 加载列表
  Future<void> loadDeductions({bool refresh = false}) async {
    if (state.isLoading || state.isLoadingMore) return;

    if (refresh) {
      state = state.copyWith(isLoading: true, error: null, currentPage: 1);
    } else {
      state = state.copyWith(isLoadingMore: true, error: null);
    }

    try {
      final page = refresh ? 1 : state.currentPage;
      final response = await _service.getSentDeductions(
        page: page,
        pageSize: 10,
        status: state.statusFilter,
      );

      final newList = refresh ? response.list : [...state.list, ...response.list];
      final hasMore = newList.length < response.total;

      state = state.copyWith(
        list: newList,
        total: response.total,
        isLoading: false,
        isLoadingMore: false,
        currentPage: page + 1,
        hasMore: hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// 设置状态筛选
  void setStatusFilter(int? status) {
    state = state.copyWith(statusFilter: status);
    loadDeductions(refresh: true);
  }
}

final sentDeductionsProvider =
    StateNotifierProvider<SentDeductionsNotifier, SentDeductionsState>((ref) {
  final service = ref.watch(goodsDeductionServiceProvider);
  return SentDeductionsNotifier(service);
});

/// 接收的货款代扣列表状态
class ReceivedDeductionsState {
  final List<GoodsDeduction> list;
  final int total;
  final bool isLoading;
  final bool isLoadingMore;
  final String? error;
  final int currentPage;
  final int? statusFilter;
  final bool hasMore;

  ReceivedDeductionsState({
    this.list = const [],
    this.total = 0,
    this.isLoading = false,
    this.isLoadingMore = false,
    this.error,
    this.currentPage = 1,
    this.statusFilter,
    this.hasMore = true,
  });

  ReceivedDeductionsState copyWith({
    List<GoodsDeduction>? list,
    int? total,
    bool? isLoading,
    bool? isLoadingMore,
    String? error,
    int? currentPage,
    int? statusFilter,
    bool? hasMore,
  }) {
    return ReceivedDeductionsState(
      list: list ?? this.list,
      total: total ?? this.total,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      error: error,
      currentPage: currentPage ?? this.currentPage,
      statusFilter: statusFilter ?? this.statusFilter,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 接收的货款代扣列表 Notifier
class ReceivedDeductionsNotifier extends StateNotifier<ReceivedDeductionsState> {
  final GoodsDeductionService _service;

  ReceivedDeductionsNotifier(this._service) : super(ReceivedDeductionsState());

  /// 加载列表
  Future<void> loadDeductions({bool refresh = false}) async {
    if (state.isLoading || state.isLoadingMore) return;

    if (refresh) {
      state = state.copyWith(isLoading: true, error: null, currentPage: 1);
    } else {
      state = state.copyWith(isLoadingMore: true, error: null);
    }

    try {
      final page = refresh ? 1 : state.currentPage;
      final response = await _service.getReceivedDeductions(
        page: page,
        pageSize: 10,
        status: state.statusFilter,
      );

      final newList = refresh ? response.list : [...state.list, ...response.list];
      final hasMore = newList.length < response.total;

      state = state.copyWith(
        list: newList,
        total: response.total,
        isLoading: false,
        isLoadingMore: false,
        currentPage: page + 1,
        hasMore: hasMore,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// 设置状态筛选
  void setStatusFilter(int? status) {
    state = state.copyWith(statusFilter: status);
    loadDeductions(refresh: true);
  }

  /// 接收货款代扣
  Future<bool> acceptDeduction(int id) async {
    try {
      await _service.acceptDeduction(id);
      loadDeductions(refresh: true);
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 拒绝货款代扣
  Future<bool> rejectDeduction(int id, String reason) async {
    try {
      await _service.rejectDeduction(id, reason);
      loadDeductions(refresh: true);
      return true;
    } catch (e) {
      return false;
    }
  }
}

final receivedDeductionsProvider =
    StateNotifierProvider<ReceivedDeductionsNotifier, ReceivedDeductionsState>((ref) {
  final service = ref.watch(goodsDeductionServiceProvider);
  return ReceivedDeductionsNotifier(service);
});
