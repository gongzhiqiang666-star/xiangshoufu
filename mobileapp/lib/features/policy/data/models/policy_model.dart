import 'package:json_annotation/json_annotation.dart';

part 'policy_model.g.dart';

/// 费率配置
@JsonSerializable()
class RateConfig {
  @JsonKey(name: 'credit_rate')
  final String creditRate;
  @JsonKey(name: 'debit_rate')
  final String debitRate;
  @JsonKey(name: 'debit_cap')
  final String debitCap;
  @JsonKey(name: 'unionpay_rate')
  final String unionpayRate;
  @JsonKey(name: 'wechat_rate')
  final String wechatRate;
  @JsonKey(name: 'alipay_rate')
  final String alipayRate;

  RateConfig({
    required this.creditRate,
    required this.debitRate,
    required this.debitCap,
    required this.unionpayRate,
    required this.wechatRate,
    required this.alipayRate,
  });

  factory RateConfig.fromJson(Map<String, dynamic> json) =>
      _$RateConfigFromJson(json);
  Map<String, dynamic> toJson() => _$RateConfigToJson(this);

  RateConfig copyWith({
    String? creditRate,
    String? debitRate,
    String? debitCap,
    String? unionpayRate,
    String? wechatRate,
    String? alipayRate,
  }) {
    return RateConfig(
      creditRate: creditRate ?? this.creditRate,
      debitRate: debitRate ?? this.debitRate,
      debitCap: debitCap ?? this.debitCap,
      unionpayRate: unionpayRate ?? this.unionpayRate,
      wechatRate: wechatRate ?? this.wechatRate,
      alipayRate: alipayRate ?? this.alipayRate,
    );
  }
}

/// 押金返现配置项
@JsonSerializable()
class DepositCashbackItem {
  @JsonKey(name: 'deposit_amount')
  final int depositAmount;
  @JsonKey(name: 'cashback_amount')
  final int cashbackAmount;

  DepositCashbackItem({
    required this.depositAmount,
    required this.cashbackAmount,
  });

  factory DepositCashbackItem.fromJson(Map<String, dynamic> json) =>
      _$DepositCashbackItemFromJson(json);
  Map<String, dynamic> toJson() => _$DepositCashbackItemToJson(this);

  DepositCashbackItem copyWith({
    int? depositAmount,
    int? cashbackAmount,
  }) {
    return DepositCashbackItem(
      depositAmount: depositAmount ?? this.depositAmount,
      cashbackAmount: cashbackAmount ?? this.cashbackAmount,
    );
  }

  /// 获取押金金额（元）
  double get depositAmountYuan => depositAmount / 100;

  /// 获取返现金额（元）
  double get cashbackAmountYuan => cashbackAmount / 100;
}

/// 流量卡返现档位项（动态N档支持）
@JsonSerializable()
class SimCashbackTier {
  @JsonKey(name: 'tier_order')
  final int tierOrder;
  @JsonKey(name: 'tier_name')
  final String tierName;
  @JsonKey(name: 'cashback_amount')
  final int cashbackAmount;
  @JsonKey(name: 'is_last_tier')
  final bool isLastTier;

  SimCashbackTier({
    required this.tierOrder,
    required this.tierName,
    required this.cashbackAmount,
    this.isLastTier = false,
  });

  factory SimCashbackTier.fromJson(Map<String, dynamic> json) =>
      _$SimCashbackTierFromJson(json);
  Map<String, dynamic> toJson() => _$SimCashbackTierToJson(this);

  /// 获取返现金额（元）
  double get cashbackAmountYuan => cashbackAmount / 100;
}

/// 流量卡返现配置
@JsonSerializable()
class SimCashbackConfig {
  // 旧版固定3档字段（兼容）
  @JsonKey(name: 'first_time_cashback')
  final int firstTimeCashback;
  @JsonKey(name: 'second_time_cashback')
  final int secondTimeCashback;
  @JsonKey(name: 'third_plus_cashback')
  final int thirdPlusCashback;
  @JsonKey(name: 'sim_fee_amount')
  final int? simFeeAmount;

  // 新版动态N档字段
  @JsonKey(name: 'tiers')
  final List<SimCashbackTier>? tiers;

  SimCashbackConfig({
    required this.firstTimeCashback,
    required this.secondTimeCashback,
    required this.thirdPlusCashback,
    this.simFeeAmount,
    this.tiers,
  });

  factory SimCashbackConfig.fromJson(Map<String, dynamic> json) =>
      _$SimCashbackConfigFromJson(json);
  Map<String, dynamic> toJson() => _$SimCashbackConfigToJson(this);

  SimCashbackConfig copyWith({
    int? firstTimeCashback,
    int? secondTimeCashback,
    int? thirdPlusCashback,
    int? simFeeAmount,
    List<SimCashbackTier>? tiers,
  }) {
    return SimCashbackConfig(
      firstTimeCashback: firstTimeCashback ?? this.firstTimeCashback,
      secondTimeCashback: secondTimeCashback ?? this.secondTimeCashback,
      thirdPlusCashback: thirdPlusCashback ?? this.thirdPlusCashback,
      simFeeAmount: simFeeAmount ?? this.simFeeAmount,
      tiers: tiers ?? this.tiers,
    );
  }

  /// 获取首次返现金额（元）
  double get firstTimeCashbackYuan => firstTimeCashback / 100;

  /// 获取第二次返现金额（元）
  double get secondTimeCashbackYuan => secondTimeCashback / 100;

  /// 获取第三次及以后返现金额（元）
  double get thirdPlusCashbackYuan => thirdPlusCashback / 100;

  /// 获取流量费金额（元）
  double get simFeeAmountYuan => (simFeeAmount ?? 9900) / 100;

  /// 获取动态档位列表（优先使用新版tiers，否则从旧字段构造）
  List<SimCashbackTier> get dynamicTiers {
    if (tiers != null && tiers!.isNotEmpty) {
      return tiers!;
    }
    // 从旧版固定字段构造默认3档
    return [
      SimCashbackTier(tierOrder: 1, tierName: '首次', cashbackAmount: firstTimeCashback),
      SimCashbackTier(tierOrder: 2, tierName: '二次', cashbackAmount: secondTimeCashback),
      SimCashbackTier(tierOrder: 3, tierName: '后续', cashbackAmount: thirdPlusCashback, isLastTier: true),
    ];
  }
}

/// 激活奖励配置项
@JsonSerializable()
class ActivationRewardItem {
  @JsonKey(name: 'reward_name')
  final String rewardName;
  @JsonKey(name: 'min_register_days')
  final int minRegisterDays;
  @JsonKey(name: 'max_register_days')
  final int maxRegisterDays;
  @JsonKey(name: 'target_amount')
  final int targetAmount;
  @JsonKey(name: 'reward_amount')
  final int rewardAmount;
  final int priority;

  ActivationRewardItem({
    required this.rewardName,
    required this.minRegisterDays,
    required this.maxRegisterDays,
    required this.targetAmount,
    required this.rewardAmount,
    this.priority = 0,
  });

  factory ActivationRewardItem.fromJson(Map<String, dynamic> json) =>
      _$ActivationRewardItemFromJson(json);
  Map<String, dynamic> toJson() => _$ActivationRewardItemToJson(this);

  ActivationRewardItem copyWith({
    String? rewardName,
    int? minRegisterDays,
    int? maxRegisterDays,
    int? targetAmount,
    int? rewardAmount,
    int? priority,
  }) {
    return ActivationRewardItem(
      rewardName: rewardName ?? this.rewardName,
      minRegisterDays: minRegisterDays ?? this.minRegisterDays,
      maxRegisterDays: maxRegisterDays ?? this.maxRegisterDays,
      targetAmount: targetAmount ?? this.targetAmount,
      rewardAmount: rewardAmount ?? this.rewardAmount,
      priority: priority ?? this.priority,
    );
  }

  /// 获取目标交易量（万元）
  double get targetAmountWan => targetAmount / 1000000;

  /// 获取奖励金额（元）
  double get rewardAmountYuan => rewardAmount / 100;

  /// 获取奖励摘要
  String get summary =>
      '入网$minRegisterDays-$maxRegisterDays天内，交易满${targetAmountWan.toStringAsFixed(2)}万元，奖励${rewardAmountYuan.toStringAsFixed(2)}元';
}

/// 代理商政策完整信息
@JsonSerializable()
class AgentPolicy {
  @JsonKey(name: 'agent_id')
  final int agentId;
  @JsonKey(name: 'channel_id')
  final int channelId;
  @JsonKey(name: 'channel_name')
  final String? channelName;
  @JsonKey(name: 'template_id')
  final int? templateId;
  @JsonKey(name: 'template_name')
  final String? templateName;

  // 费率配置
  @JsonKey(name: 'credit_rate')
  final String creditRate;
  @JsonKey(name: 'debit_rate')
  final String debitRate;
  @JsonKey(name: 'debit_cap')
  final String debitCap;
  @JsonKey(name: 'unionpay_rate')
  final String unionpayRate;
  @JsonKey(name: 'wechat_rate')
  final String wechatRate;
  @JsonKey(name: 'alipay_rate')
  final String alipayRate;

  // 押金返现
  @JsonKey(name: 'deposit_cashbacks')
  final List<DepositCashbackItem>? depositCashbacks;

  // 流量卡返现
  @JsonKey(name: 'sim_cashback')
  final SimCashbackConfig? simCashback;

  // 激活奖励
  @JsonKey(name: 'activation_rewards')
  final List<ActivationRewardItem>? activationRewards;

  AgentPolicy({
    required this.agentId,
    required this.channelId,
    this.channelName,
    this.templateId,
    this.templateName,
    required this.creditRate,
    required this.debitRate,
    required this.debitCap,
    required this.unionpayRate,
    required this.wechatRate,
    required this.alipayRate,
    this.depositCashbacks,
    this.simCashback,
    this.activationRewards,
  });

  factory AgentPolicy.fromJson(Map<String, dynamic> json) =>
      _$AgentPolicyFromJson(json);
  Map<String, dynamic> toJson() => _$AgentPolicyToJson(this);

  /// 获取费率配置
  RateConfig get rateConfig => RateConfig(
        creditRate: creditRate,
        debitRate: debitRate,
        debitCap: debitCap,
        unionpayRate: unionpayRate,
        wechatRate: wechatRate,
        alipayRate: alipayRate,
      );
}

/// 政策限制（上级的政策配置，作为调整范围限制）
@JsonSerializable()
class PolicyLimits {
  // 费率限制（最低费率）
  @JsonKey(name: 'min_credit_rate')
  final String minCreditRate;
  @JsonKey(name: 'min_debit_rate')
  final String minDebitRate;
  @JsonKey(name: 'min_unionpay_rate')
  final String minUnionpayRate;
  @JsonKey(name: 'min_wechat_rate')
  final String minWechatRate;
  @JsonKey(name: 'min_alipay_rate')
  final String minAlipayRate;

  // 押金返现限制（最大返现）
  @JsonKey(name: 'max_deposit_cashbacks')
  final List<DepositCashbackItem>? maxDepositCashbacks;

  // 流量卡返现限制（最大返现）
  @JsonKey(name: 'max_sim_cashback')
  final SimCashbackConfig? maxSimCashback;

  // 激活奖励限制（最大奖励）
  @JsonKey(name: 'max_activation_rewards')
  final List<ActivationRewardItem>? maxActivationRewards;

  PolicyLimits({
    required this.minCreditRate,
    required this.minDebitRate,
    required this.minUnionpayRate,
    required this.minWechatRate,
    required this.minAlipayRate,
    this.maxDepositCashbacks,
    this.maxSimCashback,
    this.maxActivationRewards,
  });

  factory PolicyLimits.fromJson(Map<String, dynamic> json) =>
      _$PolicyLimitsFromJson(json);
  Map<String, dynamic> toJson() => _$PolicyLimitsToJson(this);
}

/// 更新下级政策请求
@JsonSerializable()
class UpdateSubordinatePolicyRequest {
  @JsonKey(name: 'channel_id')
  final int channelId;

  // 费率配置
  @JsonKey(name: 'credit_rate')
  final String? creditRate;
  @JsonKey(name: 'debit_rate')
  final String? debitRate;
  @JsonKey(name: 'debit_cap')
  final String? debitCap;
  @JsonKey(name: 'unionpay_rate')
  final String? unionpayRate;
  @JsonKey(name: 'wechat_rate')
  final String? wechatRate;
  @JsonKey(name: 'alipay_rate')
  final String? alipayRate;

  // 押金返现
  @JsonKey(name: 'deposit_cashbacks')
  final List<DepositCashbackItem>? depositCashbacks;

  // 流量卡返现
  @JsonKey(name: 'sim_cashback')
  final SimCashbackConfig? simCashback;

  // 激活奖励
  @JsonKey(name: 'activation_rewards')
  final List<ActivationRewardItem>? activationRewards;

  UpdateSubordinatePolicyRequest({
    required this.channelId,
    this.creditRate,
    this.debitRate,
    this.debitCap,
    this.unionpayRate,
    this.wechatRate,
    this.alipayRate,
    this.depositCashbacks,
    this.simCashback,
    this.activationRewards,
  });

  factory UpdateSubordinatePolicyRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateSubordinatePolicyRequestFromJson(json);
  Map<String, dynamic> toJson() => _$UpdateSubordinatePolicyRequestToJson(this);
}

/// 通道信息
@JsonSerializable()
class ChannelInfo {
  final int id;
  @JsonKey(name: 'channel_code')
  final String channelCode;
  @JsonKey(name: 'channel_name')
  final String channelName;
  final int status;

  ChannelInfo({
    required this.id,
    required this.channelCode,
    required this.channelName,
    required this.status,
  });

  factory ChannelInfo.fromJson(Map<String, dynamic> json) =>
      _$ChannelInfoFromJson(json);
  Map<String, dynamic> toJson() => _$ChannelInfoToJson(this);
}
