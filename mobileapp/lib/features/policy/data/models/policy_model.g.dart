// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'policy_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

RateConfig _$RateConfigFromJson(Map<String, dynamic> json) => RateConfig(
      creditRate: json['credit_rate'] as String,
      debitRate: json['debit_rate'] as String,
      debitCap: json['debit_cap'] as String,
      unionpayRate: json['unionpay_rate'] as String,
      wechatRate: json['wechat_rate'] as String,
      alipayRate: json['alipay_rate'] as String,
    );

Map<String, dynamic> _$RateConfigToJson(RateConfig instance) =>
    <String, dynamic>{
      'credit_rate': instance.creditRate,
      'debit_rate': instance.debitRate,
      'debit_cap': instance.debitCap,
      'unionpay_rate': instance.unionpayRate,
      'wechat_rate': instance.wechatRate,
      'alipay_rate': instance.alipayRate,
    };

DepositCashbackItem _$DepositCashbackItemFromJson(Map<String, dynamic> json) =>
    DepositCashbackItem(
      depositAmount: (json['deposit_amount'] as num).toInt(),
      cashbackAmount: (json['cashback_amount'] as num).toInt(),
    );

Map<String, dynamic> _$DepositCashbackItemToJson(
        DepositCashbackItem instance) =>
    <String, dynamic>{
      'deposit_amount': instance.depositAmount,
      'cashback_amount': instance.cashbackAmount,
    };

SimCashbackTier _$SimCashbackTierFromJson(Map<String, dynamic> json) =>
    SimCashbackTier(
      tierOrder: (json['tier_order'] as num).toInt(),
      tierName: json['tier_name'] as String,
      cashbackAmount: (json['cashback_amount'] as num).toInt(),
      isLastTier: json['is_last_tier'] as bool? ?? false,
    );

Map<String, dynamic> _$SimCashbackTierToJson(SimCashbackTier instance) =>
    <String, dynamic>{
      'tier_order': instance.tierOrder,
      'tier_name': instance.tierName,
      'cashback_amount': instance.cashbackAmount,
      'is_last_tier': instance.isLastTier,
    };

SimCashbackConfig _$SimCashbackConfigFromJson(Map<String, dynamic> json) =>
    SimCashbackConfig(
      firstTimeCashback: (json['first_time_cashback'] as num).toInt(),
      secondTimeCashback: (json['second_time_cashback'] as num).toInt(),
      thirdPlusCashback: (json['third_plus_cashback'] as num).toInt(),
      simFeeAmount: (json['sim_fee_amount'] as num?)?.toInt(),
      tiers: (json['tiers'] as List<dynamic>?)
          ?.map((e) => SimCashbackTier.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$SimCashbackConfigToJson(SimCashbackConfig instance) =>
    <String, dynamic>{
      'first_time_cashback': instance.firstTimeCashback,
      'second_time_cashback': instance.secondTimeCashback,
      'third_plus_cashback': instance.thirdPlusCashback,
      'sim_fee_amount': instance.simFeeAmount,
      'tiers': instance.tiers,
    };

ActivationRewardItem _$ActivationRewardItemFromJson(
        Map<String, dynamic> json) =>
    ActivationRewardItem(
      rewardName: json['reward_name'] as String,
      minRegisterDays: (json['min_register_days'] as num).toInt(),
      maxRegisterDays: (json['max_register_days'] as num).toInt(),
      targetAmount: (json['target_amount'] as num).toInt(),
      rewardAmount: (json['reward_amount'] as num).toInt(),
      priority: (json['priority'] as num?)?.toInt() ?? 0,
    );

Map<String, dynamic> _$ActivationRewardItemToJson(
        ActivationRewardItem instance) =>
    <String, dynamic>{
      'reward_name': instance.rewardName,
      'min_register_days': instance.minRegisterDays,
      'max_register_days': instance.maxRegisterDays,
      'target_amount': instance.targetAmount,
      'reward_amount': instance.rewardAmount,
      'priority': instance.priority,
    };

AgentPolicy _$AgentPolicyFromJson(Map<String, dynamic> json) => AgentPolicy(
      agentId: (json['agent_id'] as num).toInt(),
      channelId: (json['channel_id'] as num).toInt(),
      channelName: json['channel_name'] as String?,
      templateId: (json['template_id'] as num?)?.toInt(),
      templateName: json['template_name'] as String?,
      creditRate: json['credit_rate'] as String,
      debitRate: json['debit_rate'] as String,
      debitCap: json['debit_cap'] as String,
      unionpayRate: json['unionpay_rate'] as String,
      wechatRate: json['wechat_rate'] as String,
      alipayRate: json['alipay_rate'] as String,
      depositCashbacks: (json['deposit_cashbacks'] as List<dynamic>?)
          ?.map((e) => DepositCashbackItem.fromJson(e as Map<String, dynamic>))
          .toList(),
      simCashback: json['sim_cashback'] == null
          ? null
          : SimCashbackConfig.fromJson(
              json['sim_cashback'] as Map<String, dynamic>),
      activationRewards: (json['activation_rewards'] as List<dynamic>?)
          ?.map((e) => ActivationRewardItem.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$AgentPolicyToJson(AgentPolicy instance) =>
    <String, dynamic>{
      'agent_id': instance.agentId,
      'channel_id': instance.channelId,
      'channel_name': instance.channelName,
      'template_id': instance.templateId,
      'template_name': instance.templateName,
      'credit_rate': instance.creditRate,
      'debit_rate': instance.debitRate,
      'debit_cap': instance.debitCap,
      'unionpay_rate': instance.unionpayRate,
      'wechat_rate': instance.wechatRate,
      'alipay_rate': instance.alipayRate,
      'deposit_cashbacks': instance.depositCashbacks,
      'sim_cashback': instance.simCashback,
      'activation_rewards': instance.activationRewards,
    };

PolicyLimits _$PolicyLimitsFromJson(Map<String, dynamic> json) => PolicyLimits(
      minCreditRate: json['min_credit_rate'] as String,
      minDebitRate: json['min_debit_rate'] as String,
      minUnionpayRate: json['min_unionpay_rate'] as String,
      minWechatRate: json['min_wechat_rate'] as String,
      minAlipayRate: json['min_alipay_rate'] as String,
      maxDepositCashbacks: (json['max_deposit_cashbacks'] as List<dynamic>?)
          ?.map((e) => DepositCashbackItem.fromJson(e as Map<String, dynamic>))
          .toList(),
      maxSimCashback: json['max_sim_cashback'] == null
          ? null
          : SimCashbackConfig.fromJson(
              json['max_sim_cashback'] as Map<String, dynamic>),
      maxActivationRewards: (json['max_activation_rewards'] as List<dynamic>?)
          ?.map((e) => ActivationRewardItem.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$PolicyLimitsToJson(PolicyLimits instance) =>
    <String, dynamic>{
      'min_credit_rate': instance.minCreditRate,
      'min_debit_rate': instance.minDebitRate,
      'min_unionpay_rate': instance.minUnionpayRate,
      'min_wechat_rate': instance.minWechatRate,
      'min_alipay_rate': instance.minAlipayRate,
      'max_deposit_cashbacks': instance.maxDepositCashbacks,
      'max_sim_cashback': instance.maxSimCashback,
      'max_activation_rewards': instance.maxActivationRewards,
    };

UpdateSubordinatePolicyRequest _$UpdateSubordinatePolicyRequestFromJson(
        Map<String, dynamic> json) =>
    UpdateSubordinatePolicyRequest(
      channelId: (json['channel_id'] as num).toInt(),
      creditRate: json['credit_rate'] as String?,
      debitRate: json['debit_rate'] as String?,
      debitCap: json['debit_cap'] as String?,
      unionpayRate: json['unionpay_rate'] as String?,
      wechatRate: json['wechat_rate'] as String?,
      alipayRate: json['alipay_rate'] as String?,
      depositCashbacks: (json['deposit_cashbacks'] as List<dynamic>?)
          ?.map((e) => DepositCashbackItem.fromJson(e as Map<String, dynamic>))
          .toList(),
      simCashback: json['sim_cashback'] == null
          ? null
          : SimCashbackConfig.fromJson(
              json['sim_cashback'] as Map<String, dynamic>),
      activationRewards: (json['activation_rewards'] as List<dynamic>?)
          ?.map((e) => ActivationRewardItem.fromJson(e as Map<String, dynamic>))
          .toList(),
    );

Map<String, dynamic> _$UpdateSubordinatePolicyRequestToJson(
        UpdateSubordinatePolicyRequest instance) =>
    <String, dynamic>{
      'channel_id': instance.channelId,
      'credit_rate': instance.creditRate,
      'debit_rate': instance.debitRate,
      'debit_cap': instance.debitCap,
      'unionpay_rate': instance.unionpayRate,
      'wechat_rate': instance.wechatRate,
      'alipay_rate': instance.alipayRate,
      'deposit_cashbacks': instance.depositCashbacks,
      'sim_cashback': instance.simCashback,
      'activation_rewards': instance.activationRewards,
    };

ChannelInfo _$ChannelInfoFromJson(Map<String, dynamic> json) => ChannelInfo(
      id: (json['id'] as num).toInt(),
      channelCode: json['channel_code'] as String,
      channelName: json['channel_name'] as String,
      status: (json['status'] as num).toInt(),
    );

Map<String, dynamic> _$ChannelInfoToJson(ChannelInfo instance) =>
    <String, dynamic>{
      'id': instance.id,
      'channel_code': instance.channelCode,
      'channel_name': instance.channelName,
      'status': instance.status,
    };
