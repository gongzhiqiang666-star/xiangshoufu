import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/features/terminal/domain/models/terminal.dart';

void main() {
  group('Terminal Model Tests', () {
    group('TerminalStatus', () {
      test('fromValue returns correct status', () {
        expect(TerminalStatus.fromValue(1), TerminalStatus.pending);
        expect(TerminalStatus.fromValue(2), TerminalStatus.allocated);
        expect(TerminalStatus.fromValue(3), TerminalStatus.bound);
        expect(TerminalStatus.fromValue(4), TerminalStatus.activated);
        expect(TerminalStatus.fromValue(5), TerminalStatus.unbound);
        expect(TerminalStatus.fromValue(6), TerminalStatus.recycled);
      });

      test('fromValue returns pending for unknown value', () {
        expect(TerminalStatus.fromValue(999), TerminalStatus.pending);
      });
    });

    group('Terminal', () {
      test('fromJson parses correctly', () {
        final json = {
          'id': 1,
          'terminal_sn': 'SN123456',
          'channel_id': 1,
          'channel_code': 'HENGXINTONG',
          'brand_code': 'NEWLAND',
          'model_code': 'N910',
          'owner_agent_id': 100,
          'merchant_id': 200,
          'merchant_no': 'M123',
          'status': 4,
          'activated_at': '2024-01-15T10:00:00Z',
          'bound_at': '2024-01-10T10:00:00Z',
          'sim_fee_count': 2,
          'created_at': '2024-01-01T00:00:00Z',
          'updated_at': '2024-01-15T10:00:00Z',
        };

        final terminal = Terminal.fromJson(json);

        expect(terminal.id, 1);
        expect(terminal.terminalSn, 'SN123456');
        expect(terminal.channelId, 1);
        expect(terminal.channelCode, 'HENGXINTONG');
        expect(terminal.brandCode, 'NEWLAND');
        expect(terminal.modelCode, 'N910');
        expect(terminal.ownerAgentId, 100);
        expect(terminal.merchantId, 200);
        expect(terminal.merchantNo, 'M123');
        expect(terminal.status, TerminalStatus.activated);
        expect(terminal.simFeeCount, 2);
      });

      test('isActivated returns true for activated terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.activated,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.isActivated, true);
      });

      test('isActivated returns false for non-activated terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.bound,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.isActivated, false);
      });

      test('canRecall returns false for activated terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.activated,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.canRecall, false);
      });

      test('canRecall returns true for non-activated terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.bound,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.canRecall, true);
      });

      test('canDistribute returns true for pending terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.canDistribute, true);
      });

      test('canDistribute returns true for allocated terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.allocated,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.canDistribute, true);
      });

      test('canDistribute returns false for bound terminal', () {
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.bound,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );

        expect(terminal.canDistribute, false);
      });

      test('toJson returns correct map', () {
        final now = DateTime.now();
        final terminal = Terminal(
          id: 1,
          terminalSn: 'SN123',
          channelId: 1,
          channelCode: 'TEST',
          ownerAgentId: 1,
          status: TerminalStatus.pending,
          createdAt: now,
          updatedAt: now,
        );

        final json = terminal.toJson();

        expect(json['id'], 1);
        expect(json['terminal_sn'], 'SN123');
        expect(json['channel_id'], 1);
        expect(json['channel_code'], 'TEST');
        expect(json['status'], 1);
      });
    });

    group('TerminalStats', () {
      test('fromJson parses correctly', () {
        final json = {
          'total': 100,
          'pending_count': 20,
          'allocated_count': 30,
          'bound_count': 10,
          'activated_count': 35,
          'unbound_count': 5,
          'yesterday_activated': 3,
          'today_activated': 5,
          'month_activated': 50,
        };

        final stats = TerminalStats.fromJson(json);

        expect(stats.total, 100);
        expect(stats.pendingCount, 20);
        expect(stats.allocatedCount, 30);
        expect(stats.boundCount, 10);
        expect(stats.activatedCount, 35);
        expect(stats.unboundCount, 5);
        expect(stats.yesterdayActivated, 3);
        expect(stats.todayActivated, 5);
        expect(stats.monthActivated, 50);
      });

      test('inactiveCount is calculated correctly', () {
        final stats = TerminalStats(
          total: 100,
          pendingCount: 20,
          allocatedCount: 30,
          boundCount: 10,
          activatedCount: 35,
          unboundCount: 5,
          yesterdayActivated: 3,
          todayActivated: 5,
          monthActivated: 50,
        );

        expect(stats.inactiveCount, 65); // 100 - 35
      });

      test('stockCount is calculated correctly', () {
        final stats = TerminalStats(
          total: 100,
          pendingCount: 20,
          allocatedCount: 30,
          boundCount: 10,
          activatedCount: 35,
          unboundCount: 5,
          yesterdayActivated: 3,
          todayActivated: 5,
          monthActivated: 50,
        );

        expect(stats.stockCount, 50); // 20 + 30
      });
    });

    group('TerminalDistribute', () {
      test('fromJson parses correctly', () {
        final json = {
          'id': 1,
          'distribute_no': 'D202401010001',
          'from_agent_id': 100,
          'to_agent_id': 101,
          'terminal_sn': 'SN123',
          'channel_id': 1,
          'is_cross_level': false,
          'goods_price': 5000,
          'deduction_type': 3,
          'status': 1,
          'source': 2,
          'created_at': '2024-01-01T00:00:00Z',
        };

        final distribute = TerminalDistribute.fromJson(json);

        expect(distribute.id, 1);
        expect(distribute.distributeNo, 'D202401010001');
        expect(distribute.fromAgentId, 100);
        expect(distribute.toAgentId, 101);
        expect(distribute.terminalSn, 'SN123');
        expect(distribute.goodsPrice, 5000);
        expect(distribute.deductionType, 3);
        expect(distribute.status, 1);
      });

      test('statusLabel returns correct labels', () {
        expect(
          TerminalDistribute(
            id: 1, distributeNo: '', fromAgentId: 1, toAgentId: 2,
            terminalSn: '', channelId: 1, isCrossLevel: false,
            goodsPrice: 0, deductionType: 1, status: 1, source: 1,
            createdAt: DateTime.now(),
          ).statusLabel,
          '待确认',
        );

        expect(
          TerminalDistribute(
            id: 1, distributeNo: '', fromAgentId: 1, toAgentId: 2,
            terminalSn: '', channelId: 1, isCrossLevel: false,
            goodsPrice: 0, deductionType: 1, status: 2, source: 1,
            createdAt: DateTime.now(),
          ).statusLabel,
          '已确认',
        );

        expect(
          TerminalDistribute(
            id: 1, distributeNo: '', fromAgentId: 1, toAgentId: 2,
            terminalSn: '', channelId: 1, isCrossLevel: false,
            goodsPrice: 0, deductionType: 1, status: 3, source: 1,
            createdAt: DateTime.now(),
          ).statusLabel,
          '已拒绝',
        );

        expect(
          TerminalDistribute(
            id: 1, distributeNo: '', fromAgentId: 1, toAgentId: 2,
            terminalSn: '', channelId: 1, isCrossLevel: false,
            goodsPrice: 0, deductionType: 1, status: 4, source: 1,
            createdAt: DateTime.now(),
          ).statusLabel,
          '已取消',
        );
      });
    });

    group('TerminalRecall', () {
      test('fromJson parses correctly', () {
        final json = {
          'id': 1,
          'recall_no': 'R202401010001',
          'from_agent_id': 101,
          'to_agent_id': 100,
          'terminal_sn': 'SN123',
          'channel_id': 1,
          'is_cross_level': false,
          'status': 1,
          'source': 2,
          'created_at': '2024-01-01T00:00:00Z',
        };

        final recall = TerminalRecall.fromJson(json);

        expect(recall.id, 1);
        expect(recall.recallNo, 'R202401010001');
        expect(recall.fromAgentId, 101);
        expect(recall.toAgentId, 100);
        expect(recall.terminalSn, 'SN123');
        expect(recall.status, 1);
      });

      test('statusLabel returns correct labels', () {
        expect(
          TerminalRecall(
            id: 1, recallNo: '', fromAgentId: 1, toAgentId: 2,
            terminalSn: '', channelId: 1, isCrossLevel: false,
            status: 1, source: 1, createdAt: DateTime.now(),
          ).statusLabel,
          '待确认',
        );

        expect(
          TerminalRecall(
            id: 1, recallNo: '', fromAgentId: 1, toAgentId: 2,
            terminalSn: '', channelId: 1, isCrossLevel: false,
            status: 2, source: 1, createdAt: DateTime.now(),
          ).statusLabel,
          '已确认',
        );
      });
    });
  });
}
