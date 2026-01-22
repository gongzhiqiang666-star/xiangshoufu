import 'package:flutter_test/flutter_test.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:xiangshoufu_app/features/terminal/domain/models/terminal.dart';
import 'package:xiangshoufu_app/features/terminal/presentation/providers/terminal_provider.dart';

void main() {
  group('Terminal Provider Tests', () {
    group('TerminalListState', () {
      test('initial state has correct defaults', () {
        final state = TerminalListState();

        expect(state.terminals, isEmpty);
        expect(state.isLoading, false);
        expect(state.hasMore, true);
        expect(state.currentPage, 1);
        expect(state.statusFilter, isNull);
        expect(state.error, isNull);
      });

      test('copyWith creates new instance with updated values', () {
        final state = TerminalListState();
        final terminals = [
          Terminal(
            id: 1,
            terminalSn: 'SN001',
            channelId: 1,
            channelCode: 'TEST',
            ownerAgentId: 1,
            status: TerminalStatus.pending,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ),
        ];

        final newState = state.copyWith(
          terminals: terminals,
          isLoading: true,
          hasMore: false,
          currentPage: 2,
          statusFilter: 4,
          error: 'Test error',
        );

        expect(newState.terminals.length, 1);
        expect(newState.isLoading, true);
        expect(newState.hasMore, false);
        expect(newState.currentPage, 2);
        expect(newState.statusFilter, 4);
        expect(newState.error, 'Test error');
      });

      test('copyWith preserves original values when not specified', () {
        final terminals = [
          Terminal(
            id: 1,
            terminalSn: 'SN001',
            channelId: 1,
            channelCode: 'TEST',
            ownerAgentId: 1,
            status: TerminalStatus.pending,
            createdAt: DateTime.now(),
            updatedAt: DateTime.now(),
          ),
        ];

        final state = TerminalListState(
          terminals: terminals,
          isLoading: true,
          hasMore: false,
          currentPage: 5,
        );

        final newState = state.copyWith(isLoading: false);

        expect(newState.terminals.length, 1);
        expect(newState.isLoading, false);
        expect(newState.hasMore, false);
        expect(newState.currentPage, 5);
      });
    });

    group('TerminalDistributeState', () {
      test('initial state has correct defaults', () {
        final state = TerminalDistributeState();

        expect(state.isSubmitting, false);
        expect(state.error, isNull);
        expect(state.result, isNull);
      });

      test('copyWith creates new instance with updated values', () {
        final state = TerminalDistributeState();
        final result = TerminalDistribute(
          id: 1,
          distributeNo: 'D001',
          fromAgentId: 1,
          toAgentId: 2,
          terminalSn: 'SN001',
          channelId: 1,
          isCrossLevel: false,
          goodsPrice: 5000,
          deductionType: 1,
          status: 1,
          source: 2,
          createdAt: DateTime.now(),
        );

        final newState = state.copyWith(
          isSubmitting: true,
          error: 'Test error',
          result: result,
        );

        expect(newState.isSubmitting, true);
        expect(newState.error, 'Test error');
        expect(newState.result, isNotNull);
        expect(newState.result!.distributeNo, 'D001');
      });
    });

    group('TerminalRecallState', () {
      test('initial state has correct defaults', () {
        final state = TerminalRecallState();

        expect(state.isSubmitting, false);
        expect(state.error, isNull);
        expect(state.successCount, 0);
        expect(state.failedCount, 0);
        expect(state.errors, isEmpty);
      });

      test('copyWith creates new instance with updated values', () {
        final state = TerminalRecallState();

        final newState = state.copyWith(
          isSubmitting: true,
          successCount: 5,
          failedCount: 2,
          errors: ['Error 1', 'Error 2'],
        );

        expect(newState.isSubmitting, true);
        expect(newState.successCount, 5);
        expect(newState.failedCount, 2);
        expect(newState.errors.length, 2);
      });
    });

    group('selectedTerminalsProvider', () {
      test('initial state is empty list', () {
        final container = ProviderContainer();
        addTearDown(container.dispose);

        final selected = container.read(selectedTerminalsProvider);

        expect(selected, isEmpty);
      });

      test('can update selected terminals', () {
        final container = ProviderContainer();
        addTearDown(container.dispose);

        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        container.read(selectedTerminalsProvider.notifier).state = [terminal];

        final selected = container.read(selectedTerminalsProvider);

        expect(selected.length, 1);
        expect(selected.first.terminalSn, 'SN001');
      });

      test('can add multiple terminals', () {
        final container = ProviderContainer();
        addTearDown(container.dispose);

        final terminal1 = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        final terminal2 = Terminal(
          id: 2,
          terminalSn: 'SN002',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.allocated,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        container.read(selectedTerminalsProvider.notifier).state = [terminal1, terminal2];

        final selected = container.read(selectedTerminalsProvider);

        expect(selected.length, 2);
      });

      test('can clear selected terminals', () {
        final container = ProviderContainer();
        addTearDown(container.dispose);

        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        container.read(selectedTerminalsProvider.notifier).state = [terminal];
        expect(container.read(selectedTerminalsProvider).length, 1);

        container.read(selectedTerminalsProvider.notifier).state = [];
        expect(container.read(selectedTerminalsProvider), isEmpty);
      });
    });

    group('Terminal business rules', () {
      test('activated terminals cannot be recalled', () {
        final container = ProviderContainer();
        addTearDown(container.dispose);

        final activatedTerminal = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.activated,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(activatedTerminal.canRecall, false);
      });

      test('pending terminals can be distributed', () {
        final pendingTerminal = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(pendingTerminal.canDistribute, true);
      });

      test('bound terminals cannot be distributed', () {
        final boundTerminal = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.bound,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(boundTerminal.canDistribute, false);
      });

      test('non-activated terminals can be recalled', () {
        final pendingTerminal = Terminal(
          id: 1,
          terminalSn: 'SN001',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        final boundTerminal = Terminal(
          id: 2,
          terminalSn: 'SN002',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.bound,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(pendingTerminal.canRecall, true);
        expect(boundTerminal.canRecall, true);
      });
    });
  });
}
