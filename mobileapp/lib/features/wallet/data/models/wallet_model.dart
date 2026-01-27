/// 钱包类型
enum WalletType {
  profit(1, '分润钱包'),
  service(2, '服务费钱包'),
  reward(3, '奖励钱包'),
  charging(4, '充值钱包'),
  settlement(5, '沉淀钱包');

  final int value;
  final String label;
  const WalletType(this.value, this.label);

  static WalletType fromValue(int value) {
    return WalletType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => WalletType.profit,
    );
  }
}

/// 钱包信息
class WalletModel {
  final int id;
  final int agentId;
  final String agentName;
  final int channelId;
  final String channelName;
  final int walletType;
  final int balance; // 分
  final int frozen; // 分
  final int available; // 分
  final int totalIncome; // 分
  final int totalWithdraw; // 分
  final String updatedAt;

  WalletModel({
    required this.id,
    required this.agentId,
    required this.agentName,
    required this.channelId,
    required this.channelName,
    required this.walletType,
    required this.balance,
    required this.frozen,
    required this.available,
    required this.totalIncome,
    required this.totalWithdraw,
    required this.updatedAt,
  });

  factory WalletModel.fromJson(Map<String, dynamic> json) {
    return WalletModel(
      id: json['id'] ?? 0,
      agentId: json['agent_id'] ?? 0,
      agentName: json['agent_name'] ?? '',
      channelId: json['channel_id'] ?? 0,
      channelName: json['channel_name'] ?? '',
      walletType: json['wallet_type'] ?? 1,
      balance: json['balance'] ?? 0,
      frozen: json['frozen'] ?? 0,
      available: json['available'] ?? 0,
      totalIncome: json['total_income'] ?? 0,
      totalWithdraw: json['total_withdraw'] ?? 0,
      updatedAt: json['updated_at'] ?? '',
    );
  }

  double get balanceYuan => balance / 100;
  double get frozenYuan => frozen / 100;
  double get availableYuan => available / 100;
  double get totalIncomeYuan => totalIncome / 100;
  double get totalWithdrawYuan => totalWithdraw / 100;

  String get walletTypeName => WalletType.fromValue(walletType).label;
}

/// 钱包汇总
class WalletSummaryModel {
  final int profitBalance;
  final int serviceBalance;
  final int rewardBalance;
  final int rechargeBalance;
  final int depositBalance;
  final int totalBalance;
  final int totalAvailable;
  final int totalFrozen;

  WalletSummaryModel({
    required this.profitBalance,
    required this.serviceBalance,
    required this.rewardBalance,
    required this.rechargeBalance,
    required this.depositBalance,
    required this.totalBalance,
    required this.totalAvailable,
    required this.totalFrozen,
  });

  factory WalletSummaryModel.fromJson(Map<String, dynamic> json) {
    return WalletSummaryModel(
      profitBalance: json['profit_balance'] ?? 0,
      serviceBalance: json['service_balance'] ?? 0,
      rewardBalance: json['reward_balance'] ?? 0,
      rechargeBalance: json['recharge_balance'] ?? 0,
      depositBalance: json['deposit_balance'] ?? 0,
      totalBalance: json['total_balance'] ?? 0,
      totalAvailable: json['total_available'] ?? json['available_balance'] ?? 0,
      totalFrozen: json['total_frozen'] ?? 0,
    );
  }

  double get profitBalanceYuan => profitBalance / 100;
  double get serviceBalanceYuan => serviceBalance / 100;
  double get rewardBalanceYuan => rewardBalance / 100;
  double get rechargeBalanceYuan => rechargeBalance / 100;
  double get depositBalanceYuan => depositBalance / 100;
  double get totalBalanceYuan => totalBalance / 100;
  double get totalAvailableYuan => totalAvailable / 100;
  double get totalFrozenYuan => totalFrozen / 100;
}

/// 钱包流水
class WalletLogModel {
  final int id;
  final int walletId;
  final String logNo;
  final String type;
  final int amount; // 分
  final int balanceBefore;
  final int balanceAfter;
  final String relatedNo;
  final String remark;
  final String createdAt;

  WalletLogModel({
    required this.id,
    required this.walletId,
    required this.logNo,
    required this.type,
    required this.amount,
    required this.balanceBefore,
    required this.balanceAfter,
    required this.relatedNo,
    required this.remark,
    required this.createdAt,
  });

  factory WalletLogModel.fromJson(Map<String, dynamic> json) {
    return WalletLogModel(
      id: json['id'] ?? 0,
      walletId: json['wallet_id'] ?? 0,
      logNo: json['log_no'] ?? '',
      type: json['type'] ?? '',
      amount: json['amount'] ?? 0,
      balanceBefore: json['balance_before'] ?? 0,
      balanceAfter: json['balance_after'] ?? 0,
      relatedNo: json['related_no'] ?? '',
      remark: json['remark'] ?? '',
      createdAt: json['created_at'] ?? '',
    );
  }

  double get amountYuan => amount / 100;
  double get balanceBeforeYuan => balanceBefore / 100;
  double get balanceAfterYuan => balanceAfter / 100;

  String get typeName {
    switch (type) {
      case 'income':
        return '收入';
      case 'expense':
        return '支出';
      case 'freeze':
        return '冻结';
      case 'unfreeze':
        return '解冻';
      case 'withdraw':
        return '提现';
      default:
        return type;
    }
  }
}

/// 代理商钱包配置
class AgentWalletConfigModel {
  final int agentId;
  final bool chargingWalletEnabled;
  final int chargingWalletLimit;
  final bool settlementWalletEnabled;
  final int settlementRatio;
  final String? enabledAt;

  AgentWalletConfigModel({
    required this.agentId,
    required this.chargingWalletEnabled,
    required this.chargingWalletLimit,
    required this.settlementWalletEnabled,
    required this.settlementRatio,
    this.enabledAt,
  });

  factory AgentWalletConfigModel.fromJson(Map<String, dynamic> json) {
    return AgentWalletConfigModel(
      agentId: json['agent_id'] ?? 0,
      chargingWalletEnabled: json['charging_wallet_enabled'] ?? false,
      chargingWalletLimit: json['charging_wallet_limit'] ?? 0,
      settlementWalletEnabled: json['settlement_wallet_enabled'] ?? false,
      settlementRatio: json['settlement_ratio'] ?? 30,
      enabledAt: json['enabled_at'],
    );
  }

  double get chargingWalletLimitYuan => chargingWalletLimit / 100;
}

/// 充值钱包汇总
class ChargingWalletSummaryModel {
  final int balance;
  final int totalIssued; // 手动发放奖励总额
  final int totalAutoReward; // 系统自动奖励总额
  final int totalReward; // 奖励总金额（手动+自动）

  ChargingWalletSummaryModel({
    required this.balance,
    required this.totalIssued,
    required this.totalAutoReward,
    required this.totalReward,
  });

  factory ChargingWalletSummaryModel.fromJson(Map<String, dynamic> json) {
    return ChargingWalletSummaryModel(
      balance: json['balance'] ?? 0,
      totalIssued: json['total_issued'] ?? 0,
      totalAutoReward: json['total_auto_reward'] ?? 0,
      totalReward: json['total_reward'] ?? 0,
    );
  }

  double get balanceYuan => balance / 100;
  double get totalIssuedYuan => totalIssued / 100;
  double get totalAutoRewardYuan => totalAutoReward / 100;
  double get totalRewardYuan => totalReward / 100;
}

/// 沉淀钱包汇总
class SettlementWalletSummaryModel {
  final int subordinateTotalBalance;
  final int settlementRatio;
  final int availableAmount;
  final int usedAmount;
  final int pendingReturnAmount;

  SettlementWalletSummaryModel({
    required this.subordinateTotalBalance,
    required this.settlementRatio,
    required this.availableAmount,
    required this.usedAmount,
    required this.pendingReturnAmount,
  });

  factory SettlementWalletSummaryModel.fromJson(Map<String, dynamic> json) {
    return SettlementWalletSummaryModel(
      subordinateTotalBalance: json['subordinate_total_balance'] ?? 0,
      settlementRatio: json['settlement_ratio'] ?? 30,
      availableAmount: json['available_amount'] ?? 0,
      usedAmount: json['used_amount'] ?? 0,
      pendingReturnAmount: json['pending_return_amount'] ?? 0,
    );
  }

  double get subordinateTotalBalanceYuan => subordinateTotalBalance / 100;
  double get availableAmountYuan => availableAmount / 100;
  double get usedAmountYuan => usedAmount / 100;
  double get pendingReturnAmountYuan => pendingReturnAmount / 100;

  int get remainingAmount => availableAmount - usedAmount;
  double get remainingAmountYuan => remainingAmount / 100;
}

/// 下级余额信息
class SubordinateBalanceModel {
  final int agentId;
  final String agentName;
  final int availableBalance;

  SubordinateBalanceModel({
    required this.agentId,
    required this.agentName,
    required this.availableBalance,
  });

  factory SubordinateBalanceModel.fromJson(Map<String, dynamic> json) {
    return SubordinateBalanceModel(
      agentId: json['agent_id'] ?? 0,
      agentName: json['agent_name'] ?? '',
      availableBalance: json['available_balance'] ?? 0,
    );
  }

  double get availableBalanceYuan => availableBalance / 100;
}

/// 沉淀使用记录
class SettlementUsageModel {
  final int id;
  final String usageNo;
  final int agentId;
  final int amount;
  final int usageType; // 1=使用 2=归还
  final String usageTypeName;
  final int status; // 1=正常 2=待归还
  final String statusName;
  final String? returnDeadline;
  final String? returnedAt;
  final String? remark;
  final String createdAt;

  SettlementUsageModel({
    required this.id,
    required this.usageNo,
    required this.agentId,
    required this.amount,
    required this.usageType,
    required this.usageTypeName,
    required this.status,
    required this.statusName,
    this.returnDeadline,
    this.returnedAt,
    this.remark,
    required this.createdAt,
  });

  factory SettlementUsageModel.fromJson(Map<String, dynamic> json) {
    return SettlementUsageModel(
      id: json['id'] ?? 0,
      usageNo: json['usage_no'] ?? '',
      agentId: json['agent_id'] ?? 0,
      amount: json['amount'] ?? 0,
      usageType: json['usage_type'] ?? 1,
      usageTypeName: json['usage_type_name'] ?? '',
      status: json['status'] ?? 1,
      statusName: json['status_name'] ?? '',
      returnDeadline: json['return_deadline'],
      returnedAt: json['returned_at'],
      remark: json['remark'],
      createdAt: json['created_at'] ?? '',
    );
  }

  double get amountYuan => amount / 100;
}

// ========== 钱包拆分相关模型 ==========

/// 子钱包（按通道拆分时的二级钱包）
class SubWalletModel {
  final int channelId;
  final String channelName;
  final int balance;
  final int frozenAmount;
  final int withdrawThreshold;
  final bool canWithdraw;

  SubWalletModel({
    required this.channelId,
    required this.channelName,
    required this.balance,
    required this.frozenAmount,
    required this.withdrawThreshold,
    required this.canWithdraw,
  });

  factory SubWalletModel.fromJson(Map<String, dynamic> json) {
    return SubWalletModel(
      channelId: json['channel_id'] ?? 0,
      channelName: json['channel_name'] ?? '',
      balance: json['balance'] ?? 0,
      frozenAmount: json['frozen_amount'] ?? 0,
      withdrawThreshold: json['withdraw_threshold'] ?? 0,
      canWithdraw: json['can_withdraw'] ?? false,
    );
  }

  double get balanceYuan => balance / 100;
  double get frozenAmountYuan => frozenAmount / 100;
  double get withdrawThresholdYuan => withdrawThreshold / 100;
  int get availableAmount => balance - frozenAmount;
  double get availableAmountYuan => availableAmount / 100;
}

/// 钱包展示（支持拆分模式）
class WalletDisplayModel {
  final int walletType;
  final String walletTypeName;
  final int balance;
  final int frozenAmount;
  final int totalIncome;
  final int totalWithdraw;
  final int withdrawThreshold;
  final bool canWithdraw;
  final List<SubWalletModel>? subWallets;

  WalletDisplayModel({
    required this.walletType,
    required this.walletTypeName,
    required this.balance,
    required this.frozenAmount,
    required this.totalIncome,
    required this.totalWithdraw,
    required this.withdrawThreshold,
    required this.canWithdraw,
    this.subWallets,
  });

  factory WalletDisplayModel.fromJson(Map<String, dynamic> json) {
    List<SubWalletModel>? subWallets;
    if (json['sub_wallets'] != null) {
      subWallets = (json['sub_wallets'] as List)
          .map((e) => SubWalletModel.fromJson(e))
          .toList();
    }

    return WalletDisplayModel(
      walletType: json['wallet_type'] ?? 1,
      walletTypeName: json['wallet_type_name'] ?? '',
      balance: json['balance'] ?? 0,
      frozenAmount: json['frozen_amount'] ?? 0,
      totalIncome: json['total_income'] ?? 0,
      totalWithdraw: json['total_withdraw'] ?? 0,
      withdrawThreshold: json['withdraw_threshold'] ?? 0,
      canWithdraw: json['can_withdraw'] ?? false,
      subWallets: subWallets,
    );
  }

  double get balanceYuan => balance / 100;
  double get frozenAmountYuan => frozenAmount / 100;
  double get totalIncomeYuan => totalIncome / 100;
  double get totalWithdrawYuan => totalWithdraw / 100;
  double get withdrawThresholdYuan => withdrawThreshold / 100;
  int get availableAmount => balance - frozenAmount;
  double get availableAmountYuan => availableAmount / 100;
  bool get hasSplit => subWallets != null && subWallets!.isNotEmpty;
}

/// 钱包列表响应（支持拆分模式）
class WalletListWithSplitResponse {
  final bool splitByChannel;
  final List<WalletDisplayModel> wallets;

  WalletListWithSplitResponse({
    required this.splitByChannel,
    required this.wallets,
  });

  factory WalletListWithSplitResponse.fromJson(Map<String, dynamic> json) {
    return WalletListWithSplitResponse(
      splitByChannel: json['split_by_channel'] ?? false,
      wallets: (json['wallets'] as List? ?? [])
          .map((e) => WalletDisplayModel.fromJson(e))
          .toList(),
    );
  }
}

/// 代理商钱包拆分配置
class AgentWalletSplitConfigModel {
  final int agentId;
  final bool splitByChannel;
  final int? configuredBy;
  final String? configuredAt;

  AgentWalletSplitConfigModel({
    required this.agentId,
    required this.splitByChannel,
    this.configuredBy,
    this.configuredAt,
  });

  factory AgentWalletSplitConfigModel.fromJson(Map<String, dynamic> json) {
    return AgentWalletSplitConfigModel(
      agentId: json['agent_id'] ?? 0,
      splitByChannel: json['split_by_channel'] ?? false,
      configuredBy: json['configured_by'],
      configuredAt: json['configured_at'],
    );
  }
}
