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
  final int? channelId;
  final String? brandCode;
  final String? modelCode;
  final String? statusGroup;
  final String? keyword;
  final String? error;

  TerminalListState({
    this.terminals = const [],
    this.isLoading = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.statusFilter,
    this.channelId,
    this.brandCode,
    this.modelCode,
    this.statusGroup,
    this.keyword,
    this.error,
  });

  TerminalListState copyWith({
    List<Terminal>? terminals,
    bool? isLoading,
    bool? hasMore,
    int? currentPage,
    int? statusFilter,
    int? channelId,
    String? brandCode,
    String? modelCode,
    String? statusGroup,
    String? keyword,
    String? error,
    bool clearChannelId = false,
    bool clearBrandCode = false,
    bool clearModelCode = false,
    bool clearStatusGroup = false,
    bool clearKeyword = false,
  }) {
    return TerminalListState(
      terminals: terminals ?? this.terminals,
      isLoading: isLoading ?? this.isLoading,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      statusFilter: statusFilter,
      channelId: clearChannelId ? null : (channelId ?? this.channelId),
      brandCode: clearBrandCode ? null : (brandCode ?? this.brandCode),
      modelCode: clearModelCode ? null : (modelCode ?? this.modelCode),
      statusGroup: clearStatusGroup ? null : (statusGroup ?? this.statusGroup),
      keyword: clearKeyword ? null : (keyword ?? this.keyword),
      error: error,
    );
  }
}

/// 终端列表Notifier
class TerminalListNotifier extends StateNotifier<TerminalListState> {
  final TerminalService _service;

  TerminalListNotifier(this._service) : super(TerminalListState());

  /// 加载终端列表（首次或刷新）
  Future<void> loadTerminals({
    int? status,
    int? channelId,
    String? brandCode,
    String? modelCode,
    String? statusGroup,
    String? keyword,
  }) async {
    state = state.copyWith(
      isLoading: true,
      error: null,
      statusFilter: status,
      channelId: channelId,
      brandCode: brandCode,
      modelCode: modelCode,
      statusGroup: statusGroup,
      keyword: keyword,
    );

    try {
      final response = await _service.getTerminals(
        status: status,
        channelId: channelId,
        brandCode: brandCode,
        modelCode: modelCode,
        statusGroup: statusGroup,
        keyword: keyword,
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
        channelId: state.channelId,
        brandCode: state.brandCode,
        modelCode: state.modelCode,
        statusGroup: state.statusGroup,
        keyword: state.keyword,
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
    await loadTerminals(
      status: state.statusFilter,
      channelId: state.channelId,
      brandCode: state.brandCode,
      modelCode: state.modelCode,
      statusGroup: state.statusGroup,
      keyword: state.keyword,
    );
  }

  /// 设置状态筛选
  Future<void> setStatusFilter(int? status) async {
    await loadTerminals(
      status: status,
      channelId: state.channelId,
      brandCode: state.brandCode,
      modelCode: state.modelCode,
      keyword: state.keyword,
    );
  }

  /// 设置状态分组筛选
  Future<void> setStatusGroup(String? statusGroup) async {
    await loadTerminals(
      channelId: state.channelId,
      brandCode: state.brandCode,
      modelCode: state.modelCode,
      statusGroup: statusGroup,
      keyword: state.keyword,
    );
  }

  /// 设置通道筛选
  Future<void> setChannelFilter(int? channelId) async {
    await loadTerminals(
      channelId: channelId,
      statusGroup: state.statusGroup,
      keyword: state.keyword,
    );
  }

  /// 设置终端类型筛选
  Future<void> setTerminalTypeFilter(String? brandCode, String? modelCode) async {
    await loadTerminals(
      channelId: state.channelId,
      brandCode: brandCode,
      modelCode: modelCode,
      statusGroup: state.statusGroup,
      keyword: state.keyword,
    );
  }

  /// 设置关键词搜索
  Future<void> setKeyword(String? keyword) async {
    await loadTerminals(
      channelId: state.channelId,
      brandCode: state.brandCode,
      modelCode: state.modelCode,
      statusGroup: state.statusGroup,
      keyword: keyword,
    );
  }

  /// 重置所有筛选条件
  Future<void> resetFilters() async {
    await loadTerminals();
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

// ==================== 划拨记录列表 ====================

/// 划拨记录列表状态
class DistributeListState {
  final List<TerminalDistribute> list;
  final bool isLoading;
  final bool isLoadingMore;
  final bool hasMore;
  final int currentPage;
  final String? error;

  DistributeListState({
    this.list = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.error,
  });

  DistributeListState copyWith({
    List<TerminalDistribute>? list,
    bool? isLoading,
    bool? isLoadingMore,
    bool? hasMore,
    int? currentPage,
    String? error,
  }) {
    return DistributeListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      error: error,
    );
  }
}

/// 划拨记录列表Notifier
class DistributeListNotifier extends StateNotifier<DistributeListState> {
  final TerminalService _service;
  final String _direction;
  final Ref _ref;

  DistributeListNotifier(this._service, this._direction, this._ref)
      : super(DistributeListState());

  /// 加载列表
  Future<void> loadList({bool refresh = false}) async {
    if (state.isLoading || state.isLoadingMore) return;
    if (!refresh && !state.hasMore) return;

    final isRefresh = refresh || state.list.isEmpty;

    state = state.copyWith(
      isLoading: isRefresh,
      isLoadingMore: !isRefresh,
      error: null,
    );

    try {
      final page = isRefresh ? 1 : state.currentPage + 1;
      final response = await _service.getDistributeList(
        direction: _direction,
        page: page,
        pageSize: 20,
      );

      state = state.copyWith(
        list: isRefresh ? response.list : [...state.list, ...response.list],
        isLoading: false,
        isLoadingMore: false,
        hasMore: response.hasMore,
        currentPage: page,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// 刷新
  Future<void> refresh() => loadList(refresh: true);

  /// 确认划拨
  Future<bool> confirmDistribute(int id) async {
    try {
      await _service.confirmDistribute(id);
      // 刷新列表
      await refresh();
      // 同时刷新另一个Tab的数据
      if (_direction == 'from') {
        _ref.read(receivedDistributesProvider.notifier).refresh();
      } else {
        _ref.read(sentDistributesProvider.notifier).refresh();
      }
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 拒绝划拨
  Future<bool> rejectDistribute(int id) async {
    try {
      await _service.rejectDistribute(id);
      // 刷新列表
      await refresh();
      // 同时刷新另一个Tab的数据
      if (_direction == 'from') {
        _ref.read(receivedDistributesProvider.notifier).refresh();
      } else {
        _ref.read(sentDistributesProvider.notifier).refresh();
      }
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 取消划拨
  Future<bool> cancelDistribute(int id) async {
    try {
      await _service.cancelDistribute(id);
      // 刷新列表
      await refresh();
      // 同时刷新另一个Tab的数据
      if (_direction == 'from') {
        _ref.read(receivedDistributesProvider.notifier).refresh();
      } else {
        _ref.read(sentDistributesProvider.notifier).refresh();
      }
      return true;
    } catch (e) {
      return false;
    }
  }
}

/// 我下发的划拨记录
final sentDistributesProvider =
    StateNotifierProvider<DistributeListNotifier, DistributeListState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return DistributeListNotifier(service, 'from', ref);
});

/// 下发给我的划拨记录
final receivedDistributesProvider =
    StateNotifierProvider<DistributeListNotifier, DistributeListState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return DistributeListNotifier(service, 'to', ref);
});

// ==================== 回拨记录列表 ====================

/// 回拨记录列表状态
class RecallListState {
  final List<TerminalRecall> list;
  final bool isLoading;
  final bool isLoadingMore;
  final bool hasMore;
  final int currentPage;
  final String? error;

  RecallListState({
    this.list = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.error,
  });

  RecallListState copyWith({
    List<TerminalRecall>? list,
    bool? isLoading,
    bool? isLoadingMore,
    bool? hasMore,
    int? currentPage,
    String? error,
  }) {
    return RecallListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      error: error,
    );
  }
}

/// 回拨记录列表Notifier
class RecallListNotifier extends StateNotifier<RecallListState> {
  final TerminalService _service;
  final String _direction;
  final Ref _ref;

  RecallListNotifier(this._service, this._direction, this._ref)
      : super(RecallListState());

  /// 加载列表
  Future<void> loadList({bool refresh = false}) async {
    if (state.isLoading || state.isLoadingMore) return;
    if (!refresh && !state.hasMore) return;

    final isRefresh = refresh || state.list.isEmpty;

    state = state.copyWith(
      isLoading: isRefresh,
      isLoadingMore: !isRefresh,
      error: null,
    );

    try {
      final page = isRefresh ? 1 : state.currentPage + 1;
      final response = await _service.getRecallList(
        direction: _direction,
        page: page,
        pageSize: 20,
      );

      state = state.copyWith(
        list: isRefresh ? response.list : [...state.list, ...response.list],
        isLoading: false,
        isLoadingMore: false,
        hasMore: response.hasMore,
        currentPage: page,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// 刷新
  Future<void> refresh() => loadList(refresh: true);

  /// 确认回拨
  Future<bool> confirmRecall(int id) async {
    try {
      await _service.confirmRecall(id);
      // 刷新列表
      await refresh();
      // 同时刷新另一个Tab的数据
      if (_direction == 'from') {
        _ref.read(receivedRecallsProvider.notifier).refresh();
      } else {
        _ref.read(sentRecallsProvider.notifier).refresh();
      }
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 拒绝回拨
  Future<bool> rejectRecall(int id) async {
    try {
      await _service.rejectRecall(id);
      // 刷新列表
      await refresh();
      // 同时刷新另一个Tab的数据
      if (_direction == 'from') {
        _ref.read(receivedRecallsProvider.notifier).refresh();
      } else {
        _ref.read(sentRecallsProvider.notifier).refresh();
      }
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 取消回拨
  Future<bool> cancelRecall(int id) async {
    try {
      await _service.cancelRecall(id);
      // 刷新列表
      await refresh();
      // 同时刷新另一个Tab的数据
      if (_direction == 'from') {
        _ref.read(receivedRecallsProvider.notifier).refresh();
      } else {
        _ref.read(sentRecallsProvider.notifier).refresh();
      }
      return true;
    } catch (e) {
      return false;
    }
  }
}

/// 我回拨的记录
final sentRecallsProvider =
    StateNotifierProvider<RecallListNotifier, RecallListState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return RecallListNotifier(service, 'from', ref);
});

/// 回拨给我的记录
final receivedRecallsProvider =
    StateNotifierProvider<RecallListNotifier, RecallListState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return RecallListNotifier(service, 'to', ref);
});

// ==================== 批量设置 ====================

/// 批量设置状态
class BatchSetState {
  final bool isSubmitting;
  final String? error;
  final int successCount;
  final int failedCount;

  BatchSetState({
    this.isSubmitting = false,
    this.error,
    this.successCount = 0,
    this.failedCount = 0,
  });

  BatchSetState copyWith({
    bool? isSubmitting,
    String? error,
    int? successCount,
    int? failedCount,
  }) {
    return BatchSetState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      error: error,
      successCount: successCount ?? this.successCount,
      failedCount: failedCount ?? this.failedCount,
    );
  }
}

/// 批量设置Notifier
class BatchSetNotifier extends StateNotifier<BatchSetState> {
  final TerminalService _service;
  final Ref _ref;

  BatchSetNotifier(this._service, this._ref) : super(BatchSetState());

  /// 批量设置费率
  Future<bool> batchSetRate({
    required List<String> terminalSns,
    required int creditRate,
  }) async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final result = await _service.batchSetRate(
        terminalSns: terminalSns,
        creditRate: creditRate,
      );

      final successCount = result['success_count'] as int? ?? terminalSns.length;
      final failedCount = result['failed_count'] as int? ?? 0;

      state = state.copyWith(
        isSubmitting: false,
        successCount: successCount,
        failedCount: failedCount,
      );

      // 刷新终端列表
      _ref.read(terminalListProvider.notifier).refresh();

      return failedCount == 0;
    } catch (e) {
      state = state.copyWith(isSubmitting: false, error: e.toString());
      return false;
    }
  }

  /// 批量设置流量费
  Future<bool> batchSetSimFee({
    required List<String> terminalSns,
    required int firstSimFee,
    required int nonFirstSimFee,
    required int simFeeIntervalDays,
  }) async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final result = await _service.batchSetSimFee(
        terminalSns: terminalSns,
        firstSimFee: firstSimFee,
        nonFirstSimFee: nonFirstSimFee,
        simFeeIntervalDays: simFeeIntervalDays,
      );

      final successCount = result['success_count'] as int? ?? terminalSns.length;
      final failedCount = result['failed_count'] as int? ?? 0;

      state = state.copyWith(
        isSubmitting: false,
        successCount: successCount,
        failedCount: failedCount,
      );

      // 刷新终端列表
      _ref.read(terminalListProvider.notifier).refresh();

      return failedCount == 0;
    } catch (e) {
      state = state.copyWith(isSubmitting: false, error: e.toString());
      return false;
    }
  }

  /// 批量设置押金
  Future<bool> batchSetDeposit({
    required List<String> terminalSns,
    required int depositAmount,
  }) async {
    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final result = await _service.batchSetDeposit(
        terminalSns: terminalSns,
        depositAmount: depositAmount,
      );

      final successCount = result['success_count'] as int? ?? terminalSns.length;
      final failedCount = result['failed_count'] as int? ?? 0;

      state = state.copyWith(
        isSubmitting: false,
        successCount: successCount,
        failedCount: failedCount,
      );

      // 刷新终端列表
      _ref.read(terminalListProvider.notifier).refresh();

      return failedCount == 0;
    } catch (e) {
      state = state.copyWith(isSubmitting: false, error: e.toString());
      return false;
    }
  }

  /// 重置状态
  void reset() {
    state = BatchSetState();
  }
}

/// 批量设置Provider
final batchSetProvider =
    StateNotifierProvider<BatchSetNotifier, BatchSetState>((ref) {
  final service = ref.watch(terminalServiceProvider);
  return BatchSetNotifier(service, ref);
});

// ==================== 筛选选项 ====================

/// 筛选选项Provider
final terminalFilterOptionsProvider = FutureProvider<TerminalFilterOptions>((ref) async {
  final service = ref.watch(terminalServiceProvider);
  return service.getFilterOptions();
});

/// 带参数的筛选选项Provider
final terminalFilterOptionsWithParamsProvider = FutureProvider.family<TerminalFilterOptions, Map<String, dynamic>>((ref, params) async {
  final service = ref.watch(terminalServiceProvider);
  return service.getFilterOptions(
    channelId: params['channel_id'] as int?,
    brandCode: params['brand_code'] as String?,
    modelCode: params['model_code'] as String?,
  );
});

// ==================== 流动记录 ====================

/// 流动记录列表状态
class FlowLogListState {
  final TerminalInfo? terminal;
  final List<TerminalFlowLog> list;
  final bool isLoading;
  final bool isLoadingMore;
  final bool hasMore;
  final int currentPage;
  final String logType;
  final String? error;

  FlowLogListState({
    this.terminal,
    this.list = const [],
    this.isLoading = false,
    this.isLoadingMore = false,
    this.hasMore = true,
    this.currentPage = 1,
    this.logType = 'all',
    this.error,
  });

  FlowLogListState copyWith({
    TerminalInfo? terminal,
    List<TerminalFlowLog>? list,
    bool? isLoading,
    bool? isLoadingMore,
    bool? hasMore,
    int? currentPage,
    String? logType,
    String? error,
  }) {
    return FlowLogListState(
      terminal: terminal ?? this.terminal,
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      isLoadingMore: isLoadingMore ?? this.isLoadingMore,
      hasMore: hasMore ?? this.hasMore,
      currentPage: currentPage ?? this.currentPage,
      logType: logType ?? this.logType,
      error: error,
    );
  }
}

/// 流动记录列表Notifier
class FlowLogListNotifier extends StateNotifier<FlowLogListState> {
  final TerminalService _service;
  final String _terminalSn;

  FlowLogListNotifier(this._service, this._terminalSn) : super(FlowLogListState());

  /// 加载列表
  Future<void> loadList({bool refresh = false, String? logType}) async {
    if (state.isLoading || state.isLoadingMore) return;
    if (!refresh && !state.hasMore && logType == null) return;

    final isRefresh = refresh || state.list.isEmpty || logType != null;
    final newLogType = logType ?? state.logType;

    state = state.copyWith(
      isLoading: isRefresh,
      isLoadingMore: !isRefresh,
      error: null,
      logType: newLogType,
    );

    try {
      final page = isRefresh ? 1 : state.currentPage + 1;
      final response = await _service.getFlowLogs(
        terminalSn: _terminalSn,
        logType: newLogType,
        page: page,
        pageSize: 20,
      );

      state = state.copyWith(
        terminal: response.terminal,
        list: isRefresh ? response.list : [...state.list, ...response.list],
        isLoading: false,
        isLoadingMore: false,
        hasMore: response.hasMore,
        currentPage: page,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        isLoadingMore: false,
        error: e.toString(),
      );
    }
  }

  /// 刷新
  Future<void> refresh() => loadList(refresh: true);

  /// 设置日志类型筛选
  Future<void> setLogType(String logType) => loadList(refresh: true, logType: logType);
}

/// 流动记录列表Provider（按终端SN）
final flowLogListProvider = StateNotifierProvider.family<FlowLogListNotifier, FlowLogListState, String>((ref, terminalSn) {
  final service = ref.watch(terminalServiceProvider);
  return FlowLogListNotifier(service, terminalSn);
});
