import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/settlement_price_model.dart';
import '../../data/services/settlement_price_service.dart';

/// 结算价服务Provider
final settlementPriceServiceProvider = Provider<SettlementPriceService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return SettlementPriceService(apiClient);
});

// ==================== 我的结算价（只读） ====================

/// 我的结算价列表状态
class MySettlementPriceListState {
  final List<SettlementPriceModel> list;
  final bool isLoading;
  final String? error;
  final int page;
  final bool hasMore;

  MySettlementPriceListState({
    this.list = const [],
    this.isLoading = false,
    this.error,
    this.page = 1,
    this.hasMore = true,
  });

  MySettlementPriceListState copyWith({
    List<SettlementPriceModel>? list,
    bool? isLoading,
    String? error,
    int? page,
    bool? hasMore,
  }) {
    return MySettlementPriceListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 我的结算价列表Notifier（获取当前用户自己的结算价，只读）
class MySettlementPriceListNotifier extends StateNotifier<MySettlementPriceListState> {
  final SettlementPriceService _service;

  MySettlementPriceListNotifier(this._service) : super(MySettlementPriceListState());

  Future<void> refresh() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      // 不传agentId表示获取当前登录用户自己的结算价
      final response = await _service.getSettlementPrices(page: 1);
      state = state.copyWith(
        list: response.list,
        isLoading: false,
        page: 1,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;

    state = state.copyWith(isLoading: true);
    try {
      final nextPage = state.page + 1;
      final response = await _service.getSettlementPrices(page: nextPage);
      state = state.copyWith(
        list: [...state.list, ...response.list],
        isLoading: false,
        page: nextPage,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

/// 我的结算价列表Provider
final mySettlementPriceListProvider =
    StateNotifierProvider<MySettlementPriceListNotifier, MySettlementPriceListState>((ref) {
  final service = ref.watch(settlementPriceServiceProvider);
  return MySettlementPriceListNotifier(service);
});

// ==================== 下级代理商结算价（可编辑） ====================

/// 下级代理商结算价列表状态
class AgentSettlementPriceListState {
  final List<SettlementPriceModel> list;
  final bool isLoading;
  final String? error;
  final int page;
  final bool hasMore;

  AgentSettlementPriceListState({
    this.list = const [],
    this.isLoading = false,
    this.error,
    this.page = 1,
    this.hasMore = true,
  });

  AgentSettlementPriceListState copyWith({
    List<SettlementPriceModel>? list,
    bool? isLoading,
    String? error,
    int? page,
    bool? hasMore,
  }) {
    return AgentSettlementPriceListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 下级代理商结算价列表Notifier（获取指定下级代理商的结算价，可编辑）
class AgentSettlementPriceListNotifier extends StateNotifier<AgentSettlementPriceListState> {
  final SettlementPriceService _service;
  final int agentId;

  AgentSettlementPriceListNotifier(this._service, this.agentId) : super(AgentSettlementPriceListState());

  Future<void> refresh() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _service.getSettlementPrices(agentId: agentId, page: 1);
      state = state.copyWith(
        list: response.list,
        isLoading: false,
        page: 1,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;

    state = state.copyWith(isLoading: true);
    try {
      final nextPage = state.page + 1;
      final response = await _service.getSettlementPrices(agentId: agentId, page: nextPage);
      state = state.copyWith(
        list: [...state.list, ...response.list],
        isLoading: false,
        page: nextPage,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

/// 下级代理商结算价列表Provider（使用family模式，支持传入agentId）
final agentSettlementPriceListProvider =
    StateNotifierProvider.family<AgentSettlementPriceListNotifier, AgentSettlementPriceListState, int>((ref, agentId) {
  final service = ref.watch(settlementPriceServiceProvider);
  return AgentSettlementPriceListNotifier(service, agentId);
});

// ==================== 结算价更新操作 ====================

/// 更新费率
final updateRateProvider = FutureProvider.family<SettlementPriceModel, UpdateRateParams>((ref, params) async {
  final service = ref.watch(settlementPriceServiceProvider);
  return service.updateRate(params.id, params.data);
});

/// 更新押金返现
final updateDepositCashbackProvider = FutureProvider.family<SettlementPriceModel, UpdateDepositParams>((ref, params) async {
  final service = ref.watch(settlementPriceServiceProvider);
  return service.updateDepositCashback(params.id, params.data);
});

/// 更新流量费返现
final updateSimCashbackProvider = FutureProvider.family<SettlementPriceModel, UpdateSimParams>((ref, params) async {
  final service = ref.watch(settlementPriceServiceProvider);
  return service.updateSimCashback(params.id, params.data);
});

/// 更新费率参数
class UpdateRateParams {
  final int id;
  final Map<String, dynamic> data;

  UpdateRateParams({required this.id, required this.data});

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is UpdateRateParams && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// 更新押金返现参数
class UpdateDepositParams {
  final int id;
  final Map<String, dynamic> data;

  UpdateDepositParams({required this.id, required this.data});

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is UpdateDepositParams && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;
}

/// 更新流量费返现参数
class UpdateSimParams {
  final int id;
  final Map<String, dynamic> data;

  UpdateSimParams({required this.id, required this.data});

  @override
  bool operator ==(Object other) =>
      identical(this, other) ||
      other is UpdateSimParams && runtimeType == other.runtimeType && id == other.id;

  @override
  int get hashCode => id.hashCode;
}

// ==================== 下级代理商调价记录 ====================

/// 下级代理商调价记录状态
class AgentPriceChangeLogListState {
  final List<PriceChangeLogModel> list;
  final bool isLoading;
  final String? error;
  final int page;
  final bool hasMore;

  AgentPriceChangeLogListState({
    this.list = const [],
    this.isLoading = false,
    this.error,
    this.page = 1,
    this.hasMore = true,
  });

  AgentPriceChangeLogListState copyWith({
    List<PriceChangeLogModel>? list,
    bool? isLoading,
    String? error,
    int? page,
    bool? hasMore,
  }) {
    return AgentPriceChangeLogListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 下级代理商调价记录Notifier
class AgentPriceChangeLogListNotifier extends StateNotifier<AgentPriceChangeLogListState> {
  final SettlementPriceService _service;
  final int agentId;

  AgentPriceChangeLogListNotifier(this._service, this.agentId) : super(AgentPriceChangeLogListState());

  Future<void> refresh() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _service.getPriceChangeLogs(agentId: agentId, page: 1);
      state = state.copyWith(
        list: response.list,
        isLoading: false,
        page: 1,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;

    state = state.copyWith(isLoading: true);
    try {
      final nextPage = state.page + 1;
      final response = await _service.getPriceChangeLogs(agentId: agentId, page: nextPage);
      state = state.copyWith(
        list: [...state.list, ...response.list],
        isLoading: false,
        page: nextPage,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

/// 下级代理商调价记录Provider
final agentPriceChangeLogListProvider =
    StateNotifierProvider.family<AgentPriceChangeLogListNotifier, AgentPriceChangeLogListState, int>((ref, agentId) {
  final service = ref.watch(settlementPriceServiceProvider);
  return AgentPriceChangeLogListNotifier(service, agentId);
});

// ==================== 原有Provider（保持兼容） ====================

/// 结算价列表状态
class SettlementPriceListState {
  final List<SettlementPriceModel> list;
  final bool isLoading;
  final String? error;
  final int page;
  final bool hasMore;

  SettlementPriceListState({
    this.list = const [],
    this.isLoading = false,
    this.error,
    this.page = 1,
    this.hasMore = true,
  });

  SettlementPriceListState copyWith({
    List<SettlementPriceModel>? list,
    bool? isLoading,
    String? error,
    int? page,
    bool? hasMore,
  }) {
    return SettlementPriceListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 结算价列表Notifier
class SettlementPriceListNotifier extends StateNotifier<SettlementPriceListState> {
  final SettlementPriceService _service;

  SettlementPriceListNotifier(this._service) : super(SettlementPriceListState());

  Future<void> refresh() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _service.getSettlementPrices(page: 1);
      state = state.copyWith(
        list: response.list,
        isLoading: false,
        page: 1,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;

    state = state.copyWith(isLoading: true);
    try {
      final nextPage = state.page + 1;
      final response = await _service.getSettlementPrices(page: nextPage);
      state = state.copyWith(
        list: [...state.list, ...response.list],
        isLoading: false,
        page: nextPage,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

/// 结算价列表Provider
final settlementPriceListProvider =
    StateNotifierProvider<SettlementPriceListNotifier, SettlementPriceListState>((ref) {
  final service = ref.watch(settlementPriceServiceProvider);
  return SettlementPriceListNotifier(service);
});

/// 结算价详情Provider
final settlementPriceDetailProvider = FutureProvider.family<SettlementPriceModel, int>((ref, id) async {
  final service = ref.watch(settlementPriceServiceProvider);
  return service.getSettlementPrice(id);
});

/// 调价记录列表状态
class PriceChangeLogListState {
  final List<PriceChangeLogModel> list;
  final bool isLoading;
  final String? error;
  final int page;
  final bool hasMore;

  PriceChangeLogListState({
    this.list = const [],
    this.isLoading = false,
    this.error,
    this.page = 1,
    this.hasMore = true,
  });

  PriceChangeLogListState copyWith({
    List<PriceChangeLogModel>? list,
    bool? isLoading,
    String? error,
    int? page,
    bool? hasMore,
  }) {
    return PriceChangeLogListState(
      list: list ?? this.list,
      isLoading: isLoading ?? this.isLoading,
      error: error ?? this.error,
      page: page ?? this.page,
      hasMore: hasMore ?? this.hasMore,
    );
  }
}

/// 调价记录列表Notifier
class PriceChangeLogListNotifier extends StateNotifier<PriceChangeLogListState> {
  final SettlementPriceService _service;

  PriceChangeLogListNotifier(this._service) : super(PriceChangeLogListState());

  Future<void> refresh() async {
    state = state.copyWith(isLoading: true, error: null);
    try {
      final response = await _service.getPriceChangeLogs(page: 1);
      state = state.copyWith(
        list: response.list,
        isLoading: false,
        page: 1,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || !state.hasMore) return;

    state = state.copyWith(isLoading: true);
    try {
      final nextPage = state.page + 1;
      final response = await _service.getPriceChangeLogs(page: nextPage);
      state = state.copyWith(
        list: [...state.list, ...response.list],
        isLoading: false,
        page: nextPage,
        hasMore: response.list.length >= 20,
      );
    } catch (e) {
      state = state.copyWith(isLoading: false, error: e.toString());
    }
  }
}

/// 调价记录列表Provider
final priceChangeLogListProvider =
    StateNotifierProvider<PriceChangeLogListNotifier, PriceChangeLogListState>((ref) {
  final service = ref.watch(settlementPriceServiceProvider);
  return PriceChangeLogListNotifier(service);
});

/// 调价记录详情Provider
final priceChangeLogDetailProvider = FutureProvider.family<PriceChangeLogModel, int>((ref, id) async {
  final service = ref.watch(settlementPriceServiceProvider);
  return service.getPriceChangeLog(id);
});
