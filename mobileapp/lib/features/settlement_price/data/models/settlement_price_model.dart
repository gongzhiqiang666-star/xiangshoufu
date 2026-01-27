/// 变更类型枚举
enum ChangeType {
  init(1, '初始化'),
  rate(2, '费率调整'),
  deposit(3, '押金返现调整'),
  sim(4, '流量费返现调整'),
  activation(5, '激活奖励调整'),
  batch(6, '批量调整'),
  sync(7, '模板同步');

  final int value;
  final String label;

  const ChangeType(this.value, this.label);

  static ChangeType fromValue(int value) {
    return ChangeType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ChangeType.init,
    );
  }
}

/// 配置类型枚举
enum ConfigType {
  settlement(1, '结算价'),
  reward(2, '奖励配置');

  final int value;
  final String label;

  const ConfigType(this.value, this.label);

  static ConfigType fromValue(int value) {
    return ConfigType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => ConfigType.settlement,
    );
  }
}

/// 费率配置
class RateConfig {
  final String rate;

  RateConfig({required this.rate});

  factory RateConfig.fromJson(Map<String, dynamic> json) {
    return RateConfig(
      rate: json['rate']?.toString() ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {'rate': rate};
  }
}

/// 押金返现配置项
class DepositCashbackItem {
  final int depositAmount;
  final int cashbackAmount;

  DepositCashbackItem({
    required this.depositAmount,
    required this.cashbackAmount,
  });

  double get depositAmountYuan => depositAmount / 100.0;
  double get cashbackAmountYuan => cashbackAmount / 100.0;

  factory DepositCashbackItem.fromJson(Map<String, dynamic> json) {
    return DepositCashbackItem(
      depositAmount: json['deposit_amount'] ?? 0,
      cashbackAmount: json['cashback_amount'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'deposit_amount': depositAmount,
      'cashback_amount': cashbackAmount,
    };
  }
}

/// 激活奖励配置项
class ActivationRewardItem {
  final String rewardName;
  final int minRegisterDays;
  final int maxRegisterDays;
  final int targetAmount;
  final int rewardAmount;
  final int priority;

  ActivationRewardItem({
    required this.rewardName,
    required this.minRegisterDays,
    required this.maxRegisterDays,
    required this.targetAmount,
    required this.rewardAmount,
    required this.priority,
  });

  double get targetAmountYuan => targetAmount / 100.0;
  double get rewardAmountYuan => rewardAmount / 100.0;

  factory ActivationRewardItem.fromJson(Map<String, dynamic> json) {
    return ActivationRewardItem(
      rewardName: json['reward_name'] ?? '',
      minRegisterDays: json['min_register_days'] ?? 0,
      maxRegisterDays: json['max_register_days'] ?? 0,
      targetAmount: json['target_amount'] ?? 0,
      rewardAmount: json['reward_amount'] ?? 0,
      priority: json['priority'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'reward_name': rewardName,
      'min_register_days': minRegisterDays,
      'max_register_days': maxRegisterDays,
      'target_amount': targetAmount,
      'reward_amount': rewardAmount,
      'priority': priority,
    };
  }
}

/// 结算价模型
class SettlementPriceModel {
  final int id;
  final int agentId;
  final String agentName;
  final int channelId;
  final String channelName;
  final int? templateId;
  final String brandCode;
  final Map<String, RateConfig> rateConfigs;
  final String? creditRate;
  final String? debitRate;
  final String? debitCap;
  final String? unionpayRate;
  final String? wechatRate;
  final String? alipayRate;
  final List<DepositCashbackItem> depositCashbacks;
  final int simFirstCashback;
  final int simSecondCashback;
  final int simThirdPlusCashback;
  final int version;
  final int status;
  final String? effectiveAt;
  final String createdAt;
  final String updatedAt;

  SettlementPriceModel({
    required this.id,
    required this.agentId,
    required this.agentName,
    required this.channelId,
    required this.channelName,
    this.templateId,
    required this.brandCode,
    required this.rateConfigs,
    this.creditRate,
    this.debitRate,
    this.debitCap,
    this.unionpayRate,
    this.wechatRate,
    this.alipayRate,
    required this.depositCashbacks,
    required this.simFirstCashback,
    required this.simSecondCashback,
    required this.simThirdPlusCashback,
    required this.version,
    required this.status,
    this.effectiveAt,
    required this.createdAt,
    required this.updatedAt,
  });

  double get simFirstCashbackYuan => simFirstCashback / 100.0;
  double get simSecondCashbackYuan => simSecondCashback / 100.0;
  double get simThirdPlusCashbackYuan => simThirdPlusCashback / 100.0;
  String get statusName => status == 1 ? '启用' : '禁用';

  factory SettlementPriceModel.fromJson(Map<String, dynamic> json) {
    Map<String, RateConfig> rateConfigs = {};
    if (json['rate_configs'] != null && json['rate_configs'] is Map) {
      final rawConfigs = json['rate_configs'] as Map;
      rawConfigs.forEach((key, value) {
        if (key is String && value is Map) {
          rateConfigs[key] = RateConfig.fromJson(Map<String, dynamic>.from(value));
        }
      });
    }

    List<DepositCashbackItem> depositCashbacks = [];
    if (json['deposit_cashbacks'] != null && json['deposit_cashbacks'] is List) {
      depositCashbacks = (json['deposit_cashbacks'] as List)
          .map((e) => DepositCashbackItem.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    return SettlementPriceModel(
      id: json['id'] ?? 0,
      agentId: json['agent_id'] ?? 0,
      agentName: json['agent_name'] ?? '',
      channelId: json['channel_id'] ?? 0,
      channelName: json['channel_name'] ?? '',
      templateId: json['template_id'],
      brandCode: json['brand_code'] ?? '',
      rateConfigs: rateConfigs,
      creditRate: json['credit_rate'],
      debitRate: json['debit_rate'],
      debitCap: json['debit_cap'],
      unionpayRate: json['unionpay_rate'],
      wechatRate: json['wechat_rate'],
      alipayRate: json['alipay_rate'],
      depositCashbacks: depositCashbacks,
      simFirstCashback: json['sim_first_cashback'] ?? 0,
      simSecondCashback: json['sim_second_cashback'] ?? 0,
      simThirdPlusCashback: json['sim_third_plus_cashback'] ?? 0,
      version: json['version'] ?? 1,
      status: json['status'] ?? 0,
      effectiveAt: json['effective_at'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

/// 代理商奖励设置模型
class AgentRewardSettingModel {
  final int id;
  final int agentId;
  final String? agentName;
  final int? templateId;
  final String? templateName;
  final int rewardAmount;
  final List<ActivationRewardItem> activationRewards;
  final int version;
  final int status;
  final String createdAt;
  final String updatedAt;

  AgentRewardSettingModel({
    required this.id,
    required this.agentId,
    this.agentName,
    this.templateId,
    this.templateName,
    required this.rewardAmount,
    required this.activationRewards,
    required this.version,
    required this.status,
    required this.createdAt,
    required this.updatedAt,
  });

  double get rewardAmountYuan => rewardAmount / 100.0;

  factory AgentRewardSettingModel.fromJson(Map<String, dynamic> json) {
    List<ActivationRewardItem> activationRewards = [];
    if (json['activation_rewards'] != null && json['activation_rewards'] is List) {
      activationRewards = (json['activation_rewards'] as List)
          .map((e) => ActivationRewardItem.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    return AgentRewardSettingModel(
      id: json['id'] ?? 0,
      agentId: json['agent_id'] ?? 0,
      agentName: json['agent_name'],
      templateId: json['template_id'],
      templateName: json['template_name'],
      rewardAmount: json['reward_amount'] ?? 0,
      activationRewards: activationRewards,
      version: json['version'] ?? 1,
      status: json['status'] ?? 0,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

/// 调价记录模型
class PriceChangeLogModel {
  final int id;
  final int agentId;
  final String agentName;
  final int? channelId;
  final String channelName;
  final int changeType;
  final String changeTypeName;
  final int configType;
  final String configTypeName;
  final String fieldName;
  final String? oldValue;
  final String? newValue;
  final String changeSummary;
  final String operatorName;
  final String source;
  final String createdAt;

  PriceChangeLogModel({
    required this.id,
    required this.agentId,
    required this.agentName,
    this.channelId,
    required this.channelName,
    required this.changeType,
    required this.changeTypeName,
    required this.configType,
    required this.configTypeName,
    required this.fieldName,
    this.oldValue,
    this.newValue,
    required this.changeSummary,
    required this.operatorName,
    required this.source,
    required this.createdAt,
  });

  ChangeType get changeTypeEnum => ChangeType.fromValue(changeType);
  ConfigType get configTypeEnum => ConfigType.fromValue(configType);

  factory PriceChangeLogModel.fromJson(Map<String, dynamic> json) {
    return PriceChangeLogModel(
      id: json['id'] ?? 0,
      agentId: json['agent_id'] ?? 0,
      agentName: json['agent_name'] ?? '',
      channelId: json['channel_id'],
      channelName: json['channel_name'] ?? '',
      changeType: json['change_type'] ?? 1,
      changeTypeName: json['change_type_name'] ?? '',
      configType: json['config_type'] ?? 1,
      configTypeName: json['config_type_name'] ?? '',
      fieldName: json['field_name'] ?? '',
      oldValue: json['old_value'],
      newValue: json['new_value'],
      changeSummary: json['change_summary'] ?? '',
      operatorName: json['operator_name'] ?? '',
      source: json['source'] ?? '',
      createdAt: json['created_at'] ?? '',
    );
  }
}

/// 结算价列表响应
class SettlementPriceListResponse {
  final List<SettlementPriceModel> list;
  final int total;
  final int page;
  final int size;

  SettlementPriceListResponse({
    required this.list,
    required this.total,
    required this.page,
    required this.size,
  });

  factory SettlementPriceListResponse.fromJson(Map<String, dynamic> json) {
    List<SettlementPriceModel> list = [];
    if (json['list'] != null && json['list'] is List) {
      list = (json['list'] as List)
          .map((e) => SettlementPriceModel.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    return SettlementPriceListResponse(
      list: list,
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      size: json['size'] ?? 20,
    );
  }
}

/// 调价记录列表响应
class PriceChangeLogListResponse {
  final List<PriceChangeLogModel> list;
  final int total;
  final int page;
  final int size;

  PriceChangeLogListResponse({
    required this.list,
    required this.total,
    required this.page,
    required this.size,
  });

  factory PriceChangeLogListResponse.fromJson(Map<String, dynamic> json) {
    List<PriceChangeLogModel> list = [];
    if (json['list'] != null && json['list'] is List) {
      list = (json['list'] as List)
          .map((e) => PriceChangeLogModel.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    return PriceChangeLogListResponse(
      list: list,
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      size: json['size'] ?? 20,
    );
  }
}
