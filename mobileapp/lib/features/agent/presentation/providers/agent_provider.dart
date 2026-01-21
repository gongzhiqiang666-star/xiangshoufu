import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../../core/network/api_client.dart';
import '../../data/models/agent_model.dart';
import '../../data/services/agent_service.dart';

/// AgentService Provider
final agentServiceProvider = Provider<AgentService>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return AgentService(apiClient);
});

/// 当前代理商详情
final myProfileProvider = FutureProvider<AgentDetail>((ref) async {
  final service = ref.watch(agentServiceProvider);
  return service.getMyProfile();
});

/// 邀请码信息
final inviteCodeProvider = FutureProvider<InviteCodeInfo>((ref) async {
  final service = ref.watch(agentServiceProvider);
  return service.getInviteCode();
});

/// 团队统计
final teamStatsProvider = FutureProvider<TeamStats>((ref) async {
  final service = ref.watch(agentServiceProvider);
  return service.getTeamStats();
});

/// 下级代理商列表状态
class SubordinatesState {
  final List<AgentInfo> list;
  final int total;
  final bool isLoading;
  final String? error;
  final int currentPage;

  SubordinatesState({
    this.list = const [],
    this.total = 0,
    this.isLoading = false,
    this.error,
    this.currentPage = 1,
  });

  SubordinatesState copyWith({
    List<AgentInfo>? list,
    int? total,
    bool? isLoading,
    String? error,
    int? currentPage,
  }) {
    return SubordinatesState(
      list: list ?? this.list,
      total: total ?? this.total,
      isLoading: isLoading ?? this.isLoading,
      error: error,
      currentPage: currentPage ?? this.currentPage,
    );
  }
}

/// 下级代理商列表Notifier
class SubordinatesNotifier extends StateNotifier<SubordinatesState> {
  final AgentService _service;

  SubordinatesNotifier(this._service) : super(SubordinatesState());

  Future<void> loadSubordinates({bool refresh = false}) async {
    if (state.isLoading) return;

    final page = refresh ? 1 : state.currentPage;
    state = state.copyWith(isLoading: true, error: null);

    try {
      final response = await _service.getSubordinates(page: page);

      if (refresh) {
        state = state.copyWith(
          list: response.list,
          total: response.total,
          isLoading: false,
          currentPage: 1,
        );
      } else {
        state = state.copyWith(
          list: response.list,
          total: response.total,
          isLoading: false,
        );
      }
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }

  Future<void> loadMore() async {
    if (state.isLoading || state.list.length >= state.total) return;

    state = state.copyWith(isLoading: true);

    try {
      final nextPage = state.currentPage + 1;
      final response = await _service.getSubordinates(page: nextPage);

      state = state.copyWith(
        list: [...state.list, ...response.list],
        total: response.total,
        isLoading: false,
        currentPage: nextPage,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        error: e.toString(),
      );
    }
  }
}

final subordinatesProvider =
    StateNotifierProvider<SubordinatesNotifier, SubordinatesState>((ref) {
  final service = ref.watch(agentServiceProvider);
  return SubordinatesNotifier(service);
});

/// 创建代理商状态
class CreateAgentState {
  final bool isSubmitting;
  final String? error;
  final AgentDetail? result;

  CreateAgentState({
    this.isSubmitting = false,
    this.error,
    this.result,
  });

  CreateAgentState copyWith({
    bool? isSubmitting,
    String? error,
    AgentDetail? result,
  }) {
    return CreateAgentState(
      isSubmitting: isSubmitting ?? this.isSubmitting,
      error: error,
      result: result,
    );
  }
}

class CreateAgentNotifier extends StateNotifier<CreateAgentState> {
  final AgentService _service;

  CreateAgentNotifier(this._service) : super(CreateAgentState());

  Future<bool> createAgent(CreateAgentRequest request) async {
    if (state.isSubmitting) return false;

    state = state.copyWith(isSubmitting: true, error: null);

    try {
      final result = await _service.createAgent(request);
      state = state.copyWith(
        isSubmitting: false,
        result: result,
      );
      return true;
    } catch (e) {
      state = state.copyWith(
        isSubmitting: false,
        error: e.toString(),
      );
      return false;
    }
  }

  void reset() {
    state = CreateAgentState();
  }
}

final createAgentProvider =
    StateNotifierProvider<CreateAgentNotifier, CreateAgentState>((ref) {
  final service = ref.watch(agentServiceProvider);
  return CreateAgentNotifier(service);
});
