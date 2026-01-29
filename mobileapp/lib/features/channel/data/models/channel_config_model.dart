/// 通道费率配置
class ChannelRateConfig {
  final int id;
  final int channelId;
  final String rateCode;
  final String rateName;
  final String minRate;
  final String maxRate;
  final String defaultRate;
  final int sortOrder;
  final int status;

  ChannelRateConfig({
    required this.id,
    required this.channelId,
    required this.rateCode,
    required this.rateName,
    required this.minRate,
    required this.maxRate,
    required this.defaultRate,
    required this.sortOrder,
    required this.status,
  });

  factory ChannelRateConfig.fromJson(Map<String, dynamic> json) {
    return ChannelRateConfig(
      id: json['id'] ?? 0,
      channelId: json['channel_id'] ?? 0,
      rateCode: json['rate_code'] ?? '',
      rateName: json['rate_name'] ?? '',
      minRate: json['min_rate'] ?? '0',
      maxRate: json['max_rate'] ?? '100',
      defaultRate: json['default_rate'] ?? '0',
      sortOrder: json['sort_order'] ?? 0,
      status: json['status'] ?? 1,
    );
  }

  /// 获取最小费率数值
  double get minRateValue => double.tryParse(minRate) ?? 0;

  /// 获取最大费率数值
  double get maxRateValue => double.tryParse(maxRate) ?? 100;

  /// 获取默认费率数值
  double get defaultRateValue => double.tryParse(defaultRate) ?? 0;
}

/// 通道押金档位
class ChannelDepositTier {
  final int id;
  final int channelId;
  final String brandCode;
  final String tierCode;
  final int depositAmount;       // 分
  final String tierName;
  final int maxCashbackAmount;   // 分
  final int defaultCashback;     // 分
  final int sortOrder;
  final int status;

  ChannelDepositTier({
    required this.id,
    required this.channelId,
    required this.brandCode,
    required this.tierCode,
    required this.depositAmount,
    required this.tierName,
    required this.maxCashbackAmount,
    required this.defaultCashback,
    required this.sortOrder,
    required this.status,
  });

  factory ChannelDepositTier.fromJson(Map<String, dynamic> json) {
    return ChannelDepositTier(
      id: json['id'] ?? 0,
      channelId: json['channel_id'] ?? 0,
      brandCode: json['brand_code'] ?? '',
      tierCode: json['tier_code'] ?? '',
      depositAmount: json['deposit_amount'] ?? 0,
      tierName: json['tier_name'] ?? '',
      maxCashbackAmount: json['max_cashback_amount'] ?? 0,
      defaultCashback: json['default_cashback'] ?? 0,
      sortOrder: json['sort_order'] ?? 0,
      status: json['status'] ?? 1,
    );
  }

  /// 押金金额（元）
  double get depositAmountYuan => depositAmount / 100.0;

  /// 最高返现（元）
  double get maxCashbackAmountYuan => maxCashbackAmount / 100.0;

  /// 默认返现（元）
  double get defaultCashbackYuan => defaultCashback / 100.0;
}

/// 通道流量费返现档位
class ChannelSimCashbackTier {
  final int id;
  final int channelId;
  final String brandCode;
  final int tierOrder;
  final String tierName;
  final bool isLastTier;
  final int maxCashbackAmount;   // 分
  final int defaultCashback;     // 分
  final int simFeeAmount;        // 分
  final int status;

  ChannelSimCashbackTier({
    required this.id,
    required this.channelId,
    required this.brandCode,
    required this.tierOrder,
    required this.tierName,
    required this.isLastTier,
    required this.maxCashbackAmount,
    required this.defaultCashback,
    required this.simFeeAmount,
    required this.status,
  });

  factory ChannelSimCashbackTier.fromJson(Map<String, dynamic> json) {
    return ChannelSimCashbackTier(
      id: json['id'] ?? 0,
      channelId: json['channel_id'] ?? 0,
      brandCode: json['brand_code'] ?? '',
      tierOrder: json['tier_order'] ?? 0,
      tierName: json['tier_name'] ?? '',
      isLastTier: json['is_last_tier'] ?? false,
      maxCashbackAmount: json['max_cashback_amount'] ?? 0,
      defaultCashback: json['default_cashback'] ?? 0,
      simFeeAmount: json['sim_fee_amount'] ?? 0,
      status: json['status'] ?? 1,
    );
  }

  /// 最高返现（元）
  double get maxCashbackAmountYuan => maxCashbackAmount / 100.0;

  /// 默认返现（元）
  double get defaultCashbackYuan => defaultCashback / 100.0;

  /// 流量费金额（元）
  double get simFeeAmountYuan => simFeeAmount / 100.0;
}

/// 通道完整配置
class ChannelFullConfig {
  final int channelId;
  final String channelCode;
  final String channelName;
  final List<ChannelRateConfig> rateConfigs;
  final List<ChannelDepositTier> depositTiers;
  final List<ChannelSimCashbackTier> simCashbackTiers;

  ChannelFullConfig({
    required this.channelId,
    required this.channelCode,
    required this.channelName,
    required this.rateConfigs,
    required this.depositTiers,
    required this.simCashbackTiers,
  });

  factory ChannelFullConfig.fromJson(Map<String, dynamic> json) {
    List<ChannelRateConfig> rateConfigs = [];
    if (json['rate_configs'] != null && json['rate_configs'] is List) {
      rateConfigs = (json['rate_configs'] as List)
          .map((e) => ChannelRateConfig.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    List<ChannelDepositTier> depositTiers = [];
    if (json['deposit_tiers'] != null && json['deposit_tiers'] is List) {
      depositTiers = (json['deposit_tiers'] as List)
          .map((e) => ChannelDepositTier.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    List<ChannelSimCashbackTier> simCashbackTiers = [];
    if (json['sim_cashback_tiers'] != null && json['sim_cashback_tiers'] is List) {
      simCashbackTiers = (json['sim_cashback_tiers'] as List)
          .map((e) => ChannelSimCashbackTier.fromJson(e as Map<String, dynamic>))
          .toList();
    }

    return ChannelFullConfig(
      channelId: json['channel_id'] ?? 0,
      channelCode: json['channel_code'] ?? '',
      channelName: json['channel_name'] ?? '',
      rateConfigs: rateConfigs,
      depositTiers: depositTiers,
      simCashbackTiers: simCashbackTiers,
    );
  }

  /// 根据费率编码获取费率配置
  ChannelRateConfig? getRateConfigByCode(String rateCode) {
    try {
      return rateConfigs.firstWhere((c) => c.rateCode == rateCode);
    } catch (_) {
      return null;
    }
  }

  /// 根据押金金额获取押金档位
  ChannelDepositTier? getDepositTierByAmount(int depositAmount) {
    try {
      return depositTiers.firstWhere((t) => t.depositAmount == depositAmount);
    } catch (_) {
      return null;
    }
  }

  /// 根据档位序号获取流量费返现档位
  ChannelSimCashbackTier? getSimCashbackTierByOrder(int tierOrder) {
    try {
      // 首先精确匹配
      final exactMatch = simCashbackTiers.where((t) => t.tierOrder == tierOrder);
      if (exactMatch.isNotEmpty) {
        return exactMatch.first;
      }
      // 如果没有精确匹配，返回最后一档（isLastTier=true 或 tierOrder最大的）
      final lastTier = simCashbackTiers.where((t) => t.isLastTier);
      if (lastTier.isNotEmpty) {
        return lastTier.first;
      }
      // 返回tierOrder最大的
      if (simCashbackTiers.isNotEmpty) {
        return simCashbackTiers.reduce((a, b) => a.tierOrder > b.tierOrder ? a : b);
      }
      return null;
    } catch (_) {
      return null;
    }
  }
}
