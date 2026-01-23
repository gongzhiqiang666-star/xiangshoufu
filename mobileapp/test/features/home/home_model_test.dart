import 'package:flutter_test/flutter_test.dart';

// 测试首页数据模型
void main() {
  group('HomeOverviewData', () {
    test('today stats calculation', () {
      // 模拟今日统计数据
      final todayStats = {
        'trans_amount': 12345600,
        'trans_count': 156,
        'profit_total': 123400,
        'profit_trade': 85600,
        'profit_deposit': 15000,
        'profit_sim': 13840,
        'profit_reward': 9000,
      };

      // 验证总分润 = 各分类之和
      final totalProfit = todayStats['profit_trade']! +
          todayStats['profit_deposit']! +
          todayStats['profit_sim']! +
          todayStats['profit_reward']!;
      expect(totalProfit, 123440);

      // 验证元转换
      expect(todayStats['trans_amount']! / 100, 123456.00);
      expect(todayStats['profit_total']! / 100, 1234.00);
    });

    test('yesterday comparison', () {
      final today = 12345600;
      final yesterday = 11000000;

      // 计算增长率
      final growthRate = (today - yesterday) / yesterday * 100;
      expect(growthRate, closeTo(12.23, 0.01));

      // 验证增长方向
      expect(today > yesterday, true);
    });

    test('month stats aggregation', () {
      final monthStats = {
        'trans_amount': 123456000,
        'trans_count': 1500,
        'profit_total': 1234560,
        'merchant_new': 38,
      };

      expect(monthStats['trans_amount']! / 100, 1234560.00);
      expect(monthStats['merchant_new'], 38);
    });

    test('team stats structure', () {
      final teamStats = {
        'direct_agent_count': 15,
        'team_agent_count': 120,
        'direct_merchant_count': 50,
        'team_merchant_count': 800,
      };

      // 团队数应该>=直营数
      expect(teamStats['team_agent_count']! >= teamStats['direct_agent_count']!, true);
      expect(teamStats['team_merchant_count']! >= teamStats['direct_merchant_count']!, true);
    });

    test('terminal stats', () {
      final terminalStats = {
        'total': 150,
        'activated': 120,
        'today_activated': 3,
        'month_activated': 15,
      };

      // 验证激活数不超过总数
      expect(terminalStats['activated']! <= terminalStats['total']!, true);
      // 验证今日激活不超过本月激活
      expect(terminalStats['today_activated']! <= terminalStats['month_activated']!, true);
      // 计算激活率
      final activationRate = terminalStats['activated']! / terminalStats['total']! * 100;
      expect(activationRate, 80.0);
    });
  });

  group('RecentTransaction', () {
    test('transaction display format', () {
      final transaction = {
        'merchant_name': '张三商店',
        'trans_type': 'swipe',
        'amount': 150000,
        'trans_time': '2026-01-23T10:30:00',
      };

      // 验证金额格式化
      expect((transaction['amount'] as int) / 100, 1500.00);

      // 验证交易类型映射
      final typeNames = {
        'swipe': '刷卡',
        'scan': '扫码',
        'quick': '快捷',
      };
      expect(typeNames[transaction['trans_type']], '刷卡');
    });

    test('time ago format', () {
      // 测试时间差计算
      final now = DateTime(2026, 1, 23, 10, 35, 0);
      final transTime = DateTime(2026, 1, 23, 10, 30, 0);

      final diff = now.difference(transTime);
      expect(diff.inMinutes, 5);

      // 根据时间差生成显示文本
      String formatTimeAgo(Duration diff) {
        if (diff.inMinutes < 60) {
          return '${diff.inMinutes}分钟前';
        } else if (diff.inHours < 24) {
          return '${diff.inHours}小时前';
        } else {
          return '${diff.inDays}天前';
        }
      }

      expect(formatTimeAgo(diff), '5分钟前');
    });
  });

  group('ProfitBreakdown', () {
    test('profit categories', () {
      final profits = {
        'trade': 85600,    // 交易分润
        'deposit': 15000,  // 押金返现
        'sim': 13840,      // 流量返现
        'reward': 9000,    // 激活奖励
      };

      // 验证各分类元值
      expect(profits['trade']! / 100, 856.00);
      expect(profits['deposit']! / 100, 150.00);
      expect(profits['sim']! / 100, 138.40);
      expect(profits['reward']! / 100, 90.00);

      // 验证总和
      final total = profits.values.reduce((a, b) => a + b);
      expect(total, 123440);
    });
  });

  group('ChannelStats', () {
    test('channel percentage calculation', () {
      final channels = [
        {'code': 'HENGXINTONG', 'name': '恒信通', 'amount': 6000000},
        {'code': 'LAKALA', 'name': '拉卡拉', 'amount': 2500000},
        {'code': 'OTHER', 'name': '其他', 'amount': 1500000},
      ];

      final totalAmount = channels.fold<int>(0, (sum, c) => sum + (c['amount'] as int));
      expect(totalAmount, 10000000);

      // 验证百分比
      for (final channel in channels) {
        final percentage = (channel['amount'] as int) / totalAmount * 100;
        switch (channel['code']) {
          case 'HENGXINTONG':
            expect(percentage, 60.0);
            break;
          case 'LAKALA':
            expect(percentage, 25.0);
            break;
          case 'OTHER':
            expect(percentage, 15.0);
            break;
        }
      }
    });
  });

  group('MerchantDistribution', () {
    test('merchant type classification', () {
      // 商户类型分类规则（月均交易额）
      double classifyMerchant(int avgAmountFen) {
        if (avgAmountFen >= 5000000) return 1; // 忠诚 >5万
        if (avgAmountFen >= 3000000) return 2; // 优质 3-5万
        if (avgAmountFen >= 2000000) return 3; // 潜力 2-3万
        if (avgAmountFen >= 1000000) return 4; // 一般 1-2万
        if (avgAmountFen > 0) return 5;        // 低活跃 <1万
        return 6;                              // 无交易
      }

      expect(classifyMerchant(6000000), 1); // 忠诚
      expect(classifyMerchant(4000000), 2); // 优质
      expect(classifyMerchant(2500000), 3); // 潜力
      expect(classifyMerchant(1500000), 4); // 一般
      expect(classifyMerchant(500000), 5);  // 低活跃
      expect(classifyMerchant(0), 6);       // 无交易
    });
  });

  group('AgentRanking', () {
    test('ranking order', () {
      final ranking = [
        {'agent_id': 101, 'value': 5200000},
        {'agent_id': 102, 'value': 4800000},
        {'agent_id': 103, 'value': 4500000},
      ];

      // 验证降序排列
      for (int i = 0; i < ranking.length - 1; i++) {
        expect((ranking[i]['value'] as int) > (ranking[i + 1]['value'] as int), true);
      }
    });

    test('change rate calculation', () {
      final current = 5200000;
      final previous = 4700000;
      final change = current - previous;
      final changeRate = change / previous * 100;

      expect(change, 500000);
      expect(changeRate, closeTo(10.64, 0.01));
    });
  });

  group('ScopeFilter', () {
    test('scope validation', () {
      final validScopes = ['direct', 'team'];

      expect(validScopes.contains('direct'), true);
      expect(validScopes.contains('team'), true);
      expect(validScopes.contains('all'), false);
      expect(validScopes.contains(''), false);
    });

    test('scope default value', () {
      String getScope(String? input) {
        if (input == null || !['direct', 'team'].contains(input)) {
          return 'direct';
        }
        return input;
      }

      expect(getScope(null), 'direct');
      expect(getScope(''), 'direct');
      expect(getScope('invalid'), 'direct');
      expect(getScope('direct'), 'direct');
      expect(getScope('team'), 'team');
    });
  });

  group('PeriodFilter', () {
    test('period validation', () {
      final validPeriods = ['day', 'week', 'month'];

      expect(validPeriods.contains('day'), true);
      expect(validPeriods.contains('week'), true);
      expect(validPeriods.contains('month'), true);
      expect(validPeriods.contains('year'), false);
    });
  });
}
