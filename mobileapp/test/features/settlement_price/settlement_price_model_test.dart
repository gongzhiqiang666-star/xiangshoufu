import 'package:flutter_test/flutter_test.dart';
import 'package:xiangshoufu_app/features/settlement_price/data/models/settlement_price_model.dart';

void main() {
  group('SettlementPriceModel', () {
    // ✅ 正常流程
    test('should parse from JSON correctly', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'agent_name': '测试代理商',
        'channel_id': 1,
        'channel_name': '恒信通',
        'template_id': 10,
        'brand_code': 'HXT',
        'rate_configs': {
          'credit': {'rate': '0.60'},
          'debit': {'rate': '0.50'},
        },
        'credit_rate': '0.60',
        'debit_rate': '0.50',
        'deposit_cashbacks': [
          {'deposit_amount': 9900, 'cashback_amount': 5000},
          {'deposit_amount': 19900, 'cashback_amount': 12000},
        ],
        'sim_first_cashback': 5000,
        'sim_second_cashback': 3000,
        'sim_third_plus_cashback': 2000,
        'version': 1,
        'status': 1,
        'effective_at': '2024-01-01T00:00:00Z',
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final model = SettlementPriceModel.fromJson(json);

      expect(model.id, 1);
      expect(model.agentId, 100);
      expect(model.agentName, '测试代理商');
      expect(model.channelId, 1);
      expect(model.channelName, '恒信通');
      expect(model.creditRate, '0.60');
      expect(model.debitRate, '0.50');
      expect(model.rateConfigs.length, 2);
      expect(model.rateConfigs['credit']?.rate, '0.60');
      expect(model.depositCashbacks.length, 2);
      expect(model.depositCashbacks[0].depositAmount, 9900);
      expect(model.depositCashbacks[0].cashbackAmount, 5000);
      expect(model.simFirstCashback, 5000);
      expect(model.version, 1);
      expect(model.status, 1);
    });

    test('should calculate yuan values correctly', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'agent_name': '',
        'channel_id': 1,
        'channel_name': '',
        'brand_code': '',
        'rate_configs': {},
        'deposit_cashbacks': [],
        'sim_first_cashback': 5000,
        'sim_second_cashback': 3000,
        'sim_third_plus_cashback': 2000,
        'version': 1,
        'status': 1,
        'created_at': '',
        'updated_at': '',
      };

      final model = SettlementPriceModel.fromJson(json);

      expect(model.simFirstCashbackYuan, 50.0);
      expect(model.simSecondCashbackYuan, 30.0);
      expect(model.simThirdPlusCashbackYuan, 20.0);
    });

    // ✅ 边界情况
    test('should handle null optional fields', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'channel_id': 1,
        'version': 1,
        'status': 1,
      };

      final model = SettlementPriceModel.fromJson(json);

      expect(model.agentName, '');
      expect(model.channelName, '');
      expect(model.creditRate, isNull);
      expect(model.debitRate, isNull);
      expect(model.templateId, isNull);
      expect(model.effectiveAt, isNull);
    });

    test('should handle empty deposit cashbacks', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'channel_id': 1,
        'deposit_cashbacks': [],
        'version': 1,
        'status': 1,
      };

      final model = SettlementPriceModel.fromJson(json);

      expect(model.depositCashbacks, isEmpty);
    });

    test('should handle empty rate configs', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'channel_id': 1,
        'rate_configs': {},
        'version': 1,
        'status': 1,
      };

      final model = SettlementPriceModel.fromJson(json);

      expect(model.rateConfigs, isEmpty);
    });

    // ✅ 状态名称
    test('should return correct status name', () {
      final enabledJson = {'id': 1, 'agent_id': 1, 'channel_id': 1, 'status': 1, 'version': 1};
      final disabledJson = {'id': 1, 'agent_id': 1, 'channel_id': 1, 'status': 0, 'version': 1};

      final enabled = SettlementPriceModel.fromJson(enabledJson);
      final disabled = SettlementPriceModel.fromJson(disabledJson);

      expect(enabled.statusName, '启用');
      expect(disabled.statusName, '禁用');
    });
  });

  group('DepositCashbackItem', () {
    test('should parse from JSON correctly', () {
      final json = {
        'deposit_amount': 9900,
        'cashback_amount': 5000,
      };

      final item = DepositCashbackItem.fromJson(json);

      expect(item.depositAmount, 9900);
      expect(item.cashbackAmount, 5000);
      expect(item.depositAmountYuan, 99.0);
      expect(item.cashbackAmountYuan, 50.0);
    });

    test('should serialize to JSON correctly', () {
      final item = DepositCashbackItem(
        depositAmount: 9900,
        cashbackAmount: 5000,
      );

      final json = item.toJson();

      expect(json['deposit_amount'], 9900);
      expect(json['cashback_amount'], 5000);
    });
  });

  group('RateConfig', () {
    test('should parse from JSON correctly', () {
      final json = {'rate': '0.60'};

      final config = RateConfig.fromJson(json);

      expect(config.rate, '0.60');
    });

    test('should serialize to JSON correctly', () {
      final config = RateConfig(rate: '0.55');

      final json = config.toJson();

      expect(json['rate'], '0.55');
    });
  });

  group('PriceChangeLogModel', () {
    test('should parse from JSON correctly', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'agent_name': '测试代理商',
        'channel_id': 1,
        'channel_name': '恒信通',
        'change_type': 2,
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
      };

      final log = PriceChangeLogModel.fromJson(json);

      expect(log.id, 1);
      expect(log.agentId, 100);
      expect(log.changeType, 2);
      expect(log.changeTypeName, '费率调整');
      expect(log.configType, 1);
      expect(log.configTypeName, '结算价');
      expect(log.oldValue, '0.60');
      expect(log.newValue, '0.55');
      expect(log.operatorName, 'admin');
      expect(log.source, 'PC');
    });

    test('should return correct change type enum', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'agent_name': '',
        'channel_name': '',
        'change_type': 2,
        'change_type_name': '',
        'config_type': 1,
        'config_type_name': '',
        'field_name': '',
        'change_summary': '',
        'operator_name': '',
        'source': '',
        'created_at': '',
      };

      final log = PriceChangeLogModel.fromJson(json);

      expect(log.changeTypeEnum, ChangeType.rate);
    });

    test('should handle null optional fields', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'change_type': 1,
        'config_type': 1,
      };

      final log = PriceChangeLogModel.fromJson(json);

      expect(log.channelId, isNull);
      expect(log.oldValue, isNull);
      expect(log.newValue, isNull);
    });
  });

  group('ChangeType enum', () {
    test('should return correct enum from value', () {
      expect(ChangeType.fromValue(1), ChangeType.init);
      expect(ChangeType.fromValue(2), ChangeType.rate);
      expect(ChangeType.fromValue(3), ChangeType.deposit);
      expect(ChangeType.fromValue(4), ChangeType.sim);
      expect(ChangeType.fromValue(5), ChangeType.activation);
      expect(ChangeType.fromValue(6), ChangeType.batch);
      expect(ChangeType.fromValue(7), ChangeType.sync);
    });

    test('should return default for unknown value', () {
      expect(ChangeType.fromValue(999), ChangeType.init);
    });

    test('should have correct labels', () {
      expect(ChangeType.init.label, '初始化');
      expect(ChangeType.rate.label, '费率调整');
      expect(ChangeType.deposit.label, '押金返现调整');
      expect(ChangeType.sim.label, '流量费返现调整');
      expect(ChangeType.activation.label, '激活奖励调整');
    });
  });

  group('ConfigType enum', () {
    test('should return correct enum from value', () {
      expect(ConfigType.fromValue(1), ConfigType.settlement);
      expect(ConfigType.fromValue(2), ConfigType.reward);
    });

    test('should return default for unknown value', () {
      expect(ConfigType.fromValue(999), ConfigType.settlement);
    });

    test('should have correct labels', () {
      expect(ConfigType.settlement.label, '结算价');
      expect(ConfigType.reward.label, '奖励配置');
    });
  });

  group('SettlementPriceListResponse', () {
    test('should parse from JSON correctly', () {
      final json = {
        'list': [
          {'id': 1, 'agent_id': 100, 'channel_id': 1, 'version': 1, 'status': 1},
          {'id': 2, 'agent_id': 101, 'channel_id': 1, 'version': 1, 'status': 1},
        ],
        'total': 2,
        'page': 1,
        'size': 20,
      };

      final response = SettlementPriceListResponse.fromJson(json);

      expect(response.list.length, 2);
      expect(response.total, 2);
      expect(response.page, 1);
      expect(response.size, 20);
    });

    test('should handle empty list', () {
      final json = {
        'list': [],
        'total': 0,
        'page': 1,
        'size': 20,
      };

      final response = SettlementPriceListResponse.fromJson(json);

      expect(response.list, isEmpty);
      expect(response.total, 0);
    });
  });

  group('PriceChangeLogListResponse', () {
    test('should parse from JSON correctly', () {
      final json = {
        'list': [
          {
            'id': 1,
            'agent_id': 100,
            'change_type': 2,
            'config_type': 1,
          },
        ],
        'total': 1,
        'page': 1,
        'size': 20,
      };

      final response = PriceChangeLogListResponse.fromJson(json);

      expect(response.list.length, 1);
      expect(response.total, 1);
    });
  });

  group('AgentRewardSettingModel', () {
    test('should parse from JSON correctly', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'agent_name': '测试代理商',
        'template_id': 10,
        'template_name': '标准奖励模版',
        'reward_amount': 10000,
        'activation_rewards': [
          {
            'reward_name': '首次激活奖励',
            'min_register_days': 0,
            'max_register_days': 30,
            'target_amount': 100000,
            'reward_amount': 5000,
            'priority': 1,
          },
        ],
        'version': 1,
        'status': 1,
        'created_at': '2024-01-01T00:00:00Z',
        'updated_at': '2024-01-01T00:00:00Z',
      };

      final model = AgentRewardSettingModel.fromJson(json);

      expect(model.id, 1);
      expect(model.agentId, 100);
      expect(model.rewardAmount, 10000);
      expect(model.rewardAmountYuan, 100.0);
      expect(model.activationRewards.length, 1);
      expect(model.activationRewards[0].rewardName, '首次激活奖励');
    });

    test('should handle empty activation rewards', () {
      final json = {
        'id': 1,
        'agent_id': 100,
        'reward_amount': 0,
        'activation_rewards': [],
        'version': 1,
        'status': 1,
      };

      final model = AgentRewardSettingModel.fromJson(json);

      expect(model.activationRewards, isEmpty);
    });
  });

  group('ActivationRewardItem', () {
    test('should parse from JSON correctly', () {
      final json = {
        'reward_name': '首次激活奖励',
        'min_register_days': 0,
        'max_register_days': 30,
        'target_amount': 100000,
        'reward_amount': 5000,
        'priority': 1,
      };

      final item = ActivationRewardItem.fromJson(json);

      expect(item.rewardName, '首次激活奖励');
      expect(item.minRegisterDays, 0);
      expect(item.maxRegisterDays, 30);
      expect(item.targetAmount, 100000);
      expect(item.targetAmountYuan, 1000.0);
      expect(item.rewardAmount, 5000);
      expect(item.rewardAmountYuan, 50.0);
      expect(item.priority, 1);
    });

    test('should serialize to JSON correctly', () {
      final item = ActivationRewardItem(
        rewardName: '测试奖励',
        minRegisterDays: 0,
        maxRegisterDays: 30,
        targetAmount: 100000,
        rewardAmount: 5000,
        priority: 1,
      );

      final json = item.toJson();

      expect(json['reward_name'], '测试奖励');
      expect(json['target_amount'], 100000);
      expect(json['reward_amount'], 5000);
    });
  });
}
