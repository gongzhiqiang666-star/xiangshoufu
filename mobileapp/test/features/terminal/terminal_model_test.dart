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

    // ==================== 新增筛选和流动记录模型测试 ====================

    group('TerminalStatusGroup', () {
      test('fromValue returns correct status group', () {
        expect(TerminalStatusGroup.fromValue('all'), TerminalStatusGroup.all);
        expect(TerminalStatusGroup.fromValue('unstock'), TerminalStatusGroup.unstock);
        expect(TerminalStatusGroup.fromValue('stocked'), TerminalStatusGroup.stocked);
        expect(TerminalStatusGroup.fromValue('unbound'), TerminalStatusGroup.unbound);
        expect(TerminalStatusGroup.fromValue('inactive'), TerminalStatusGroup.inactive);
        expect(TerminalStatusGroup.fromValue('active'), TerminalStatusGroup.active);
      });

      test('fromValue returns all for unknown value', () {
        expect(TerminalStatusGroup.fromValue('unknown'), TerminalStatusGroup.all);
        expect(TerminalStatusGroup.fromValue(''), TerminalStatusGroup.all);
      });

      test('has correct value and label', () {
        expect(TerminalStatusGroup.all.value, 'all');
        expect(TerminalStatusGroup.all.label, '全部');
        expect(TerminalStatusGroup.unstock.value, 'unstock');
        expect(TerminalStatusGroup.unstock.label, '未出库');
        expect(TerminalStatusGroup.active.value, 'active');
        expect(TerminalStatusGroup.active.label, '已激活');
      });
    });

    group('ChannelOption', () {
      test('fromJson parses correctly', () {
        final json = {
          'channel_id': 1,
          'channel_code': 'HENGXINTONG',
        };

        final option = ChannelOption.fromJson(json);

        expect(option.channelId, 1);
        expect(option.channelCode, 'HENGXINTONG');
      });

      test('handles missing values', () {
        final json = <String, dynamic>{};
        final option = ChannelOption.fromJson(json);

        expect(option.channelId, 0);
        expect(option.channelCode, '');
      });
    });

    group('TerminalTypeOption', () {
      test('fromJson parses correctly', () {
        final json = {
          'channel_id': 1,
          'channel_code': 'HENGXINTONG',
          'brand_code': 'NEWLAND',
          'model_code': 'N910',
          'count': 50,
        };

        final option = TerminalTypeOption.fromJson(json);

        expect(option.channelId, 1);
        expect(option.channelCode, 'HENGXINTONG');
        expect(option.brandCode, 'NEWLAND');
        expect(option.modelCode, 'N910');
        expect(option.count, 50);
      });

      test('displayName returns correct format', () {
        final option = TerminalTypeOption(
          channelId: 1,
          channelCode: 'HENGXINTONG',
          brandCode: 'NEWLAND',
          modelCode: 'N910',
          count: 50,
        );

        expect(option.displayName, 'NEWLAND N910');
      });

      test('handles missing values', () {
        final json = <String, dynamic>{};
        final option = TerminalTypeOption.fromJson(json);

        expect(option.channelId, 0);
        expect(option.brandCode, '');
        expect(option.modelCode, '');
        expect(option.count, 0);
      });
    });

    group('StatusGroupCount', () {
      test('fromJson parses correctly', () {
        final json = {
          'key': 'active',
          'label': '已激活',
          'count': 100,
        };

        final count = StatusGroupCount.fromJson(json);

        expect(count.key, 'active');
        expect(count.label, '已激活');
        expect(count.count, 100);
      });

      test('handles missing values', () {
        final json = <String, dynamic>{};
        final count = StatusGroupCount.fromJson(json);

        expect(count.key, '');
        expect(count.label, '');
        expect(count.count, 0);
      });
    });

    group('TerminalFilterOptions', () {
      test('fromJson parses correctly', () {
        final json = {
          'channels': [
            {'channel_id': 1, 'channel_code': 'HENGXINTONG'},
            {'channel_id': 2, 'channel_code': 'LAKALA'},
          ],
          'terminal_types': [
            {'channel_id': 1, 'channel_code': 'HENGXINTONG', 'brand_code': 'NEWLAND', 'model_code': 'N910', 'count': 50},
          ],
          'status_groups': [
            {'key': 'all', 'label': '全部', 'count': 100},
            {'key': 'active', 'label': '已激活', 'count': 50},
          ],
        };

        final options = TerminalFilterOptions.fromJson(json);

        expect(options.channels.length, 2);
        expect(options.channels[0].channelCode, 'HENGXINTONG');
        expect(options.terminalTypes.length, 1);
        expect(options.terminalTypes[0].brandCode, 'NEWLAND');
        expect(options.statusGroups.length, 2);
        expect(options.statusGroups[1].key, 'active');
      });

      test('handles empty lists', () {
        final json = <String, dynamic>{};
        final options = TerminalFilterOptions.fromJson(json);

        expect(options.channels, isEmpty);
        expect(options.terminalTypes, isEmpty);
        expect(options.statusGroups, isEmpty);
      });
    });

    group('TerminalFlowLog', () {
      test('fromJson parses correctly', () {
        final json = {
          'id': 1,
          'log_type': 'distribute',
          'log_type_name': '下发',
          'from_agent_id': 100,
          'from_agent_name': '总代理',
          'to_agent_id': 101,
          'to_agent_name': '一级代理',
          'merchant_no': '',
          'status': 2,
          'status_name': '已确认',
          'remark': '测试备注',
          'created_at': '2024-01-15T10:00:00Z',
          'confirmed_at': '2024-01-15T11:00:00Z',
        };

        final log = TerminalFlowLog.fromJson(json);

        expect(log.id, 1);
        expect(log.logType, 'distribute');
        expect(log.logTypeName, '下发');
        expect(log.fromAgentId, 100);
        expect(log.fromAgentName, '总代理');
        expect(log.toAgentId, 101);
        expect(log.toAgentName, '一级代理');
        expect(log.status, 2);
        expect(log.statusName, '已确认');
        expect(log.remark, '测试备注');
        expect(log.confirmedAt, isNotNull);
      });

      test('handles missing optional values', () {
        final json = {
          'id': 1,
          'log_type': 'bind',
          'log_type_name': '绑定',
          'status': 2,
          'status_name': '已确认',
          'created_at': '2024-01-15T10:00:00Z',
        };

        final log = TerminalFlowLog.fromJson(json);

        expect(log.id, 1);
        expect(log.logType, 'bind');
        expect(log.fromAgentId, isNull);
        expect(log.fromAgentName, '');
        expect(log.toAgentId, isNull);
        expect(log.toAgentName, '');
        expect(log.merchantNo, '');
        expect(log.remark, '');
        expect(log.confirmedAt, isNull);
      });

      test('parses different log types', () {
        final logTypes = ['distribute', 'recall', 'bind', 'unbind', 'activate'];

        for (final type in logTypes) {
          final json = {
            'id': 1,
            'log_type': type,
            'log_type_name': type,
            'status': 2,
            'status_name': '已确认',
            'created_at': '2024-01-15T10:00:00Z',
          };

          final log = TerminalFlowLog.fromJson(json);
          expect(log.logType, type);
        }
      });
    });

    group('TerminalInfo', () {
      test('fromJson parses correctly', () {
        final json = {
          'terminal_sn': 'SN123456',
          'channel_id': 1,
          'channel_code': 'HENGXINTONG',
          'brand_code': 'NEWLAND',
          'model_code': 'N910',
        };

        final info = TerminalInfo.fromJson(json);

        expect(info.terminalSn, 'SN123456');
        expect(info.channelId, 1);
        expect(info.channelCode, 'HENGXINTONG');
        expect(info.brandCode, 'NEWLAND');
        expect(info.modelCode, 'N910');
      });

      test('handles missing values', () {
        final json = <String, dynamic>{};
        final info = TerminalInfo.fromJson(json);

        expect(info.terminalSn, '');
        expect(info.channelId, 0);
        expect(info.channelCode, '');
        expect(info.brandCode, '');
        expect(info.modelCode, '');
      });
    });

    group('TerminalFlowLogsResponse', () {
      test('fromJson parses correctly', () {
        final json = {
          'terminal': {
            'terminal_sn': 'SN123456',
            'channel_id': 1,
            'channel_code': 'HENGXINTONG',
            'brand_code': 'NEWLAND',
            'model_code': 'N910',
          },
          'list': [
            {
              'id': 1,
              'log_type': 'distribute',
              'log_type_name': '下发',
              'status': 2,
              'status_name': '已确认',
              'created_at': '2024-01-15T10:00:00Z',
            },
            {
              'id': 2,
              'log_type': 'bind',
              'log_type_name': '绑定',
              'status': 2,
              'status_name': '已确认',
              'created_at': '2024-01-16T10:00:00Z',
            },
          ],
          'total': 2,
          'page': 1,
          'page_size': 20,
        };

        final response = TerminalFlowLogsResponse.fromJson(json);

        expect(response.terminal.terminalSn, 'SN123456');
        expect(response.list.length, 2);
        expect(response.list[0].logType, 'distribute');
        expect(response.list[1].logType, 'bind');
        expect(response.total, 2);
        expect(response.page, 1);
        expect(response.pageSize, 20);
      });

      test('hasMore returns true when list size equals pageSize', () {
        final json = {
          'terminal': {'terminal_sn': 'SN123'},
          'list': List.generate(20, (i) => {
            'id': i,
            'log_type': 'distribute',
            'log_type_name': '下发',
            'status': 2,
            'status_name': '已确认',
            'created_at': '2024-01-15T10:00:00Z',
          }),
          'total': 50,
          'page': 1,
          'page_size': 20,
        };

        final response = TerminalFlowLogsResponse.fromJson(json);
        expect(response.hasMore, true);
      });

      test('hasMore returns false when list size is less than pageSize', () {
        final json = {
          'terminal': {'terminal_sn': 'SN123'},
          'list': [
            {
              'id': 1,
              'log_type': 'distribute',
              'log_type_name': '下发',
              'status': 2,
              'status_name': '已确认',
              'created_at': '2024-01-15T10:00:00Z',
            },
          ],
          'total': 1,
          'page': 1,
          'page_size': 20,
        };

        final response = TerminalFlowLogsResponse.fromJson(json);
        expect(response.hasMore, false);
      });

      test('handles empty list', () {
        final json = {
          'terminal': {'terminal_sn': 'SN123'},
          'list': [],
          'total': 0,
          'page': 1,
          'page_size': 20,
        };

        final response = TerminalFlowLogsResponse.fromJson(json);
        expect(response.list, isEmpty);
        expect(response.hasMore, false);
      });
    });
  });
}
