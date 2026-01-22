import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../data/services/terminal_service.dart';
import '../../domain/models/terminal.dart';

/// 终端服务Provider
final terminalServiceProvider = Provider<TerminalService>((ref) {
  return TerminalService();
});

/// 终端统计Provider
final terminalStatsProvider = FutureProvider<TerminalStats>((ref) async {
  final service = ref.watch(terminalServiceProvider);
  return service.getTerminalStats();
});

/// 终端列表状态
class TerminalListState {
  final List<Terminal> terminals;
  final bool isLoading;
  final bool hasMore;
  final int currentPage;
  final int? statusFilter;
  final String? error;

  TerminalListState({
    this.terminals = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.statusFilter,
    this.error,
  });

  TerminalListState copyWith({
    List<Terminal>? terminals,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    int? statusFilter,
    String? error,
  }) {
    return TerminalListState(
      terminals: terminals ?? this.terminals,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      statusFilter: statusFilter,
      error: error,
    );
  }
}

/// 终端列表Notifier
class TerminalListNotifier extends StateNotifier<TerminalListState> {
  final TerminalService _service;

  TerminalListNotifier(this._service) : super(TerminalListState());

  /// 加载终端列表（首次或刷新）
  Future<void> loadTerminals({int? status}) async {
    state = state.copyWith(
      isLoading: true,
      error: null,
      statusFilter: status,
    );

    try {
      final response = await _service.getTerminals(
        status: status,
        page: 1,
        pageSize: 20,
      );

      state = state.copyWith(
        terminals: response.list,
        isLoading: false,
        hasMore: response.hasMore,
        currentPage: 1,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 加载更多
  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;

    state = state.copyWith(isLoading: true);

    try {
      final nextPage = state.currentPage + 1;
      final response = await _service.getTerminals(
        status: state.statusFilter,
        page: nextPage,
        pageSize: 20,
      );

      state = state.copyWith(
        terminals: [...state.terminals, ...response.list],
        isLoading: false,
        hasMore: response.hasMore,
        currentPage: nextPage,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  /// 刷新
  Future<void> refresh() async {
    await loadTerminals(status: state.statusFilter);
  }

  /// 设置状态筛选
  Future<void> setStatusFilter(int? status) async {
    await loadTerminals(status: status);
  }
}

/// 终端列表Provider
final terminalListProvider =
    StateNotifierProvider<TerminalListNotifier, TerminalListState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return TerminalListNotifier(service);
});

/// 终端下发状态
class TerminalDistributeState {
  final bool isSubmitting;
  final String? error;
  final TerminalDistribute? result;

  TerminalDistributeState({
    this.isSubmitting = false,
    this.error,
    this.result,
  });

  TerminalDistributeState copyWith({
    bool? isSubmitting,
    String? error,
    TerminalDistribute? result,
  }) {
    return TerminalDistributeState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      error: error,
      result: result,
    );
  }
}

/// 终端下发Notifier
class TerminalDistributeNotifier extends StateNotifier<TerminalDistributeState> {
  final TerminalService _service;
  final Ref _ref;

  TerminalDistributeNotifier(this._service, this._ref)
      : super(TerminalDistributeState());

  /// 提交下发
  Future<bool> submitDistribute({
    required int toAgentId,
    required String terminalSn,
    required int channelId,
    required int goodsPrice,
    required int deductionType,
    int? deductionPeriods,
    String? remark,
  }) async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final result = await _service.distributeTerminal(
        toAgentId: toAgentId,
        terminalSn: terminalSn,
        channelId: channelId,
        goodsPrice: goodsPrice,
        deductionType: deductionType,
        deductionPeriods: deductionPeriods,
        remark: remark,
      );

      state = state.copyWith(isSubmitting: false, result: result);

      // 刷新终端列表和统计
      _ref.invalidate(terminalStatsProvider);
      _ref.read(terminalListProvider.notifier).refresh();

      return true;
    } catch (e) {
      state = state.copyWith(isSubmitting: false, error: e.toString());
      return false;
    }
  }

  /// 批量下发
  Future<int> batchDistribute({
    required int toAgentId,
    required List<String> terminalSns,
    required int channelId,
    required int goodsPrice,
    required int deductionType,
    int? deductionPeriods,
    String? remark,
  }) async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final results = await _service.batchDistributeTerminals(
        toAgentId: toAgentId,
        terminalSns: terminalSns,
        channelId: channelId,
        goodsPrice: goodsPrice,
        deductionType: deductionType,
        deductionPeriods: deductionPeriods,
        remark: remark,
      );

      state = state.copyWith(isSubmitting: false);

      // 刷新终端列表和统计
      _ref.invalidate(terminalStatsProvider);
      _ref.read(terminalListProvider.notifier).refresh();

      return results.length;
    } catch (e) {
      state = state.copyWith(isSubmitting: false, error: e.toString());
      return 0;
    }
  }

  /// 重置状态
  void reset() {
    state = TerminalDistributeState();
  }
}

/// 终端下发Provider
final terminalDistributeProvider =
    StateNotifierProvider<TerminalDistributeNotifier, TerminalDistributeState>(
        (ref) {
  final service = ref.watch(terminalServiceProvider);
  return TerminalDistributeNotifier(service, ref);
});

/// 终端回拨状态
class TerminalRecallState {
  final bool isSubmitting;
  final String? error;
  final int successCount;
  final int failedCount;
  final List<String> errors;

  TerminalRecallState({
    this.isSubmitting = false,
    this.error,
    this.successCount = 0,
    this.failedCount = 0,
    this.errors = const [],
  });

  TerminalRecallState copyWith({
    bool? isSubmitting,
    String? error,
    int? successCount,
    int? failedCount,
    List<String>? errors,
  }) {
    return TerminalRecallState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      error: error,
      successCount: successCount ?? this.successCount,
      failedCount: failedCount ?? this.failedCount,
      errors: errors ?? this.errors,
    );
  }
}

/// 终端回拨Notifier
class TerminalRecallNotifier extends StateNotifier<TerminalRecallState> {
  final TerminalService _service;
  final Ref _ref;

  TerminalRecallNotifier(this._service, this._ref)
      : super(TerminalRecallState());

  /// 批量回拨
  Future<bool> batchRecall({
    required int toAgentId,
    required List<String> terminalSns,
    String? remark,
  }) async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final result = await _service.batchRecallTerminals(
        toAgentId: toAgentId,
        terminalSns: terminalSns,
        remark: remark,
      );

      final successCount = result['success_count'] as int? ?? 0;
      final failedCount = result['failed_count'] as int? ?? 0;
      final errors = (result['errors'] as List?)?.cast<String>() ?? [];

      state = state.copyWith(
        isSubmitting: false,
        successCount: successCount,
        failedCount: failedCount,
        errors: errors,
      );

      // 刷新终端列表和统计
      _ref.invalidate(terminalStatsProvider);
      _ref.read(terminalListProvider.notifier).refresh();

      return failedCount == 0;
    } catch (e) {
      state = state.copyWith(isSubmitting: false, error: e.toString());
      return false;
    }
  }

  /// 重置状态
  void reset() {
    state = TerminalRecallState();
  }
}

/// 终端回拨Provider
final terminalRecallProvider =
    StateNotifierProvider<TerminalRecallNotifier, TerminalRecallState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return TerminalRecallNotifier(service, ref);
});

/// 终端详情Provider
final terminalDetailProvider = FutureProvider.family<Terminal, String>((ref, sn) async {
  final service = ref.watch(terminalServiceProvider);
  return service.getTerminalDetail(sn);
});

/// 选中的终端列表
final selectedTerminalsProvider = StateProvider<List<Terminal>>((ref) => []);
