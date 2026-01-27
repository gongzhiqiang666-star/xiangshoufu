import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:xiangshoufu_app/features/settlement_price/data/models/settlement_price_model.dart';
import 'package:xiangshoufu_app/features/settlement_price/presentation/providers/settlement_price_provider.dart';

void main() {
  group('SettlementPriceListState', () {
    // ✅ 正常流程
    test('should have correct initial state', () {
      final state = SettlementPriceListState();

      expect(state.list, isEmpty);
      expect(state.isLoading, false);
      expect(state.error, isNull);
      expect(state.page, 1);
      expect(state.hasMore, true);
    });

    test('should create state with custom values', () {
      final mockList = [
        _createMockSettlementPrice(id: 1),
        _createMockSettlementPrice(id: 2),
      ];

      final state = SettlementPriceListState(
        list: mockList,
        isLoading: true,
        error: 'Test error',
        page: 2,
        hasMore: false,
      );

      expect(state.list.length, 2);
      expect(state.isLoading, true);
      expect(state.error, 'Test error');
      expect(state.page, 2);
      expect(state.hasMore, false);
    });

    // ✅ 边界情况
    test('should handle empty list', () {
      final state = SettlementPriceListState(
        list: [],
        hasMore: false,
      );

      expect(state.list, isEmpty);
      expect(state.hasMore, false);
    });

    // ✅ copyWith测试
    test('should copyWith correctly', () {
      final state = SettlementPriceListState();
      final newState = state.copyWith(
        isLoading: true,
        page: 2,
      );

      expect(newState.isLoading, true);
      expect(newState.page, 2);
      expect(newState.list, isEmpty);
      expect(newState.error, isNull);
    });

    test('should preserve values not specified in copyWith', () {
      final mockList = [_createMockSettlementPrice(id: 1)];
      final state = SettlementPriceListState(
        list: mockList,
        error: 'Original error',
      );

      final newState = state.copyWith(isLoading: true);

      expect(newState.list, mockList);
      expect(newState.error, 'Original error');
      expect(newState.isLoading, true);
    });
  });

  group('PriceChangeLogListState', () {
    // ✅ 正常流程
    test('should have correct initial state', () {
      final state = PriceChangeLogListState();

      expect(state.list, isEmpty);
      expect(state.isLoading, false);
      expect(state.error, isNull);
      expect(state.page, 1);
      expect(state.hasMore, true);
    });

    test('should create state with custom values', () {
      final mockList = [
        _createMockPriceChangeLog(id: 1),
        _createMockPriceChangeLog(id: 2),
      ];

      final state = PriceChangeLogListState(
        list: mockList,
        isLoading: false,
        page: 3,
        hasMore: true,
      );

      expect(state.list.length, 2);
      expect(state.page, 3);
    });

    // ✅ copyWith测试
    test('should copyWith correctly', () {
      final state = PriceChangeLogListState();
      final newState = state.copyWith(
        error: 'Network error',
        hasMore: false,
      );

      expect(newState.error, 'Network error');
      expect(newState.hasMore, false);
      expect(newState.page, 1);
    });
  });

  group('Provider definitions', () {
    test('settlementPriceServiceProvider should be defined', () {
      // 验证Provider定义存在
      expect(settlementPriceServiceProvider, isNotNull);
    });

    test('settlementPriceListProvider should be defined', () {
      expect(settlementPriceListProvider, isNotNull);
    });

    test('settlementPriceDetailProvider should be a family provider', () {
      // Family provider可以接受参数
      final provider = settlementPriceDetailProvider(1);
      expect(provider, isNotNull);
    });

    test('priceChangeLogListProvider should be defined', () {
      expect(priceChangeLogListProvider, isNotNull);
    });

    test('priceChangeLogDetailProvider should be a family provider', () {
      final provider = priceChangeLogDetailProvider(1);
      expect(provider, isNotNull);
    });
  });

  group('SettlementPriceListNotifier state transitions', () {
    late ProviderContainer container;

    setUp(() {
      container = ProviderContainer();
    });

    tearDown(() {
      container.dispose();
    });

    test('initial state should be correct', () {
      final state = container.read(settlementPriceListProvider);

      expect(state.list, isEmpty);
      expect(state.isLoading, false);
      expect(state.page, 1);
    });
  });

  group('PriceChangeLogListNotifier state transitions', () {
    late ProviderContainer container;

    setUp(() {
      container = ProviderContainer();
    });

    tearDown(() {
      container.dispose();
    });

    test('initial state should be correct', () {
      final state = container.read(priceChangeLogListProvider);

      expect(state.list, isEmpty);
      expect(state.isLoading, false);
      expect(state.page, 1);
    });
  });

  group('SettlementPriceListState pagination', () {
    test('should track page number correctly', () {
      var state = SettlementPriceListState(page: 1);
      expect(state.page, 1);

      state = state.copyWith(page: 2);
      expect(state.page, 2);

      state = state.copyWith(page: 3);
      expect(state.page, 3);
    });

    test('should track hasMore flag correctly', () {
      var state = SettlementPriceListState(hasMore: true);
      expect(state.hasMore, true);

      state = state.copyWith(hasMore: false);
      expect(state.hasMore, false);
    });

    test('should accumulate list items on load more', () {
      final list1 = [_createMockSettlementPrice(id: 1)];
      var state = SettlementPriceListState(list: list1, page: 1);
      expect(state.list.length, 1);

      final list2 = [...list1, _createMockSettlementPrice(id: 2)];
      state = state.copyWith(list: list2, page: 2);
      expect(state.list.length, 2);
    });
  });

  group('Error handling states', () {
    test('should set error and stop loading on failure', () {
      final state = SettlementPriceListState(
        isLoading: false,
        error: '网络错误',
      );

      expect(state.isLoading, false);
      expect(state.error, '网络错误');
    });

    test('should clear error by creating new state', () {
      var state = SettlementPriceListState(error: '之前的错误');
      expect(state.error, isNotNull);

      // 创建新状态来清除错误（模拟刷新成功）
      state = SettlementPriceListState(list: [_createMockSettlementPrice(id: 1)]);
      expect(state.error, isNull);
      expect(state.list.length, 1);
    });

    test('PriceChangeLogListState should handle errors', () {
      final state = PriceChangeLogListState(
        error: '加载失败',
        isLoading: false,
      );

      expect(state.error, '加载失败');
      expect(state.isLoading, false);
    });
  });
}

// ============================================================
// Helper functions for creating mock data
// ============================================================

SettlementPriceModel _createMockSettlementPrice({
  required int id,
  int agentId = 100,
  int channelId = 1,
}) {
  return SettlementPriceModel.fromJson({
    'id': id,
    'agent_id': agentId,
    'agent_name': '测试代理商$id',
    'channel_id': channelId,
    'channel_name': '恒信通',
    'brand_code': 'HXT',
    'rate_configs': {'credit': {'rate': '0.60'}},
    'deposit_cashbacks': [{'deposit_amount': 9900, 'cashback_amount': 5000}],
    'sim_first_cashback': 5000,
    'sim_second_cashback': 3000,
    'sim_third_plus_cashback': 2000,
    'version': 1,
    'status': 1,
    'created_at': '2024-01-01T00:00:00Z',
    'updated_at': '2024-01-01T00:00:00Z',
  });
}

PriceChangeLogModel _createMockPriceChangeLog({
  required int id,
  int agentId = 100,
  int changeType = 2,
}) {
  return PriceChangeLogModel.fromJson({
    'id': id,
    'agent_id': agentId,
    'agent_name': '测试代理商',
    'channel_id': 1,
    'channel_name': '恒信通',
    'change_type': changeType,
    'change_type_name': '费率调整',
    'config_type': 1,
    'config_type_name': '结算价',
    'field_name': 'credit_rate',
    'old_value': '0.60',
    'new_value': '0.55',
    'change_summary': '贷记卡费率: 0.60% → 0.55%',
    'operator_name': 'admin',
    'source': 'PC',
    'created_at': '2024-01-02T10:00:00Z',
  });
}
