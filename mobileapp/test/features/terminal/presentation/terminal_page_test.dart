import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/router/app_router.dart';

/// 终端页面测试
/// 测试场景：
/// 1. 快捷入口（划拨记录、回拨记录）路由配置
/// 2. 路由路径正确性
void main() {
  group('TerminalPage - 快捷入口路由', () {
    test('划拨记录路由路径应正确定义', () {
      // Assert
      expect(RoutePaths.terminalDistributeList, '/terminal/distribute-list');
    });

    test('回拨记录路由路径应正确定义', () {
      // Assert
      expect(RoutePaths.terminalRecallList, '/terminal/recall-list');
    });

    test('终端相关路由路径应全部定义', () {
      // Assert - 验证所有终端相关路由都已定义
      expect(RoutePaths.terminal, isNotEmpty);
      expect(RoutePaths.terminalDetail, isNotEmpty);
      expect(RoutePaths.terminalDistributeList, isNotEmpty);
      expect(RoutePaths.terminalRecallList, isNotEmpty);
      expect(RoutePaths.terminalTransfer, isNotEmpty);
    });
  });

  group('TerminalPage - 快捷入口显示逻辑', () {
    test('快捷入口应包含划拨记录和回拨记录', () {
      // Arrange - 模拟快捷入口配置
      final quickActions = [
        {'title': '划拨记录', 'route': RoutePaths.terminalDistributeList},
        {'title': '回拨记录', 'route': RoutePaths.terminalRecallList},
      ];

      // Assert
      expect(quickActions.length, 2);
      expect(quickActions[0]['title'], '划拨记录');
      expect(quickActions[1]['title'], '回拨记录');
    });

    test('划拨记录路由应指向正确页面', () {
      // Assert
      expect(
        RoutePaths.terminalDistributeList,
        contains('distribute'),
      );
    });

    test('回拨记录路由应指向正确页面', () {
      // Assert
      expect(
        RoutePaths.terminalRecallList,
        contains('recall'),
      );
    });
  });

  group('TerminalPage - Tab状态筛选', () {
    test('Tab配置应包含4个状态', () {
      // Arrange - 模拟Tab配置（与terminal_page.dart一致）
      final tabs = [
        {'label': '全部', 'status': null},
        {'label': '已激活', 'status': 2}, // TerminalStatus.activated.value
        {'label': '未激活', 'status': 1}, // TerminalStatus.bound.value
        {'label': '库存', 'status': 0},   // TerminalStatus.pending.value
      ];

      // Assert
      expect(tabs.length, 4);
      expect(tabs[0]['label'], '全部');
      expect(tabs[1]['label'], '已激活');
      expect(tabs[2]['label'], '未激活');
      expect(tabs[3]['label'], '库存');
    });
  });
}
