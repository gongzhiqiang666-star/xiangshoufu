import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/deduction_model.dart';
import '../../data/services/deduction_service.dart';

/// 代扣服务 Provider
final deductionServiceProvider = Provider<DeductionService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return DeductionService(apiClient);
});

/// 代扣计划详情 Provider
final deductionPlanDetailProvider = FutureProvider.family<DeductionPlanDetail, int>((ref, id) async {
  final service = ref.watch(deductionServiceProvider);
  return service.getPlanDetail(id);
});

/// 代扣计划统计 Provider
final deductionStatsProvider = FutureProvider<DeductionPlanStats>((ref) async {
  final service = ref.watch(deductionServiceProvider);
  return service.getStats();
});

/// 代扣计划列表状态
class DeductionPlansState {
  final List<DeductionPlan> list;
  final int total;
  final bool isLoading;
  final bool isLoadingMore;
  final String? error;
  final int currentPage;
  final int? statusFilter;
  final int? typeFilter;
  final bool hasMore;

  DeductionPlansState({
    this.list = const [],
    this.total = 0,
    this.isLoading = false,
    this.isLoadingMore = false,
    this.error,
    this.currentPage = 1,
    this.statusFilter,
    this.typeFilter,
    this.hasMore = true,
  });

  DeductionPlansState copyWith({
    List<DeductionPlan>? list,
    int? total,
    bool? isLoading,
    bool? isLoadingMore,
    String? error,
    int? currentPage,
    int? statusFilter,
    int? typeFilter,
    bool? hasMore,
  }) {
    return DeductionPlansState(
      list: list ?? this.list,
      total: total ?? this.total,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      error: error,
      currentPage: currentPage ?? this.currentPage,
      statusFilter: statusFilter ?? this.statusFilter,
      typeFilter: typeFilter ?? this.typeFilter,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 代扣计划列表 Notifier
class DeductionPlansNotifier extends StateNotifier<DeductionPlansState> {
  final DeductionService _service;

  DeductionPlansNotifier(this._service) : super(DeductionPlansState());

  /// 加载列表
  Future<void> loadPlans({bool refresh = false}) async {
    if (state.isLoading || state.isLoadingMore) return;

    if (refresh) {
      state = state.copyWith(isLoading: true, error: null, currentPage: 1);
    } else {
      state = state.copyWith(isLoadingMore: true, error: null);
    }

    try {
      final page = refresh ? 1 : state.currentPage;
      final response = await _service.getDeductionPlans(
        page: page,
        pageSize: 10,
        planType: state.typeFilter,
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
    loadPlans(refresh: true);
  }

  /// 设置类型筛选
  void setTypeFilter(int? type) {
    state = state.copyWith(typeFilter: type);
    loadPlans(refresh: true);
  }

  /// 暂停代扣计划
  Future<bool> pausePlan(int id) async {
    try {
      await _service.pausePlan(id);
      loadPlans(refresh: true);
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 恢复代扣计划
  Future<bool> resumePlan(int id) async {
    try {
      await _service.resumePlan(id);
      loadPlans(refresh: true);
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 取消代扣计划
  Future<bool> cancelPlan(int id) async {
    try {
      await _service.cancelPlan(id);
      loadPlans(refresh: true);
      return true;
    } catch (e) {
      return false;
    }
  }
}

final deductionPlansProvider =
    StateNotifierProvider<DeductionPlansNotifier, DeductionPlansState>((ref) {
  final service = ref.watch(deductionServiceProvider);
  return DeductionPlansNotifier(service);
});
