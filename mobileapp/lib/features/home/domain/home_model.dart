/// 首页数据模型
class HomeOverviewData {
  final TodayStats today;
  final TodayStats yesterday;
  final MonthStats month;
  final TeamStats team;
  final TerminalStats terminal;
  final WalletStats wallet;

  HomeOverviewData({
    required this.today,
    required this.yesterday,
    required this.month,
    required this.team,
    required this.terminal,
    required this.wallet,
  });

  factory HomeOverviewData.fromJson(Map<String, dynamic> json) {
    return HomeOverviewData(
      today: TodayStats.fromJson(json['today'] ?? {}),
      yesterday: TodayStats.fromJson(json['yesterday'] ?? {}),
      month: MonthStats.fromJson(json['month'] ?? {}),
      team: TeamStats.fromJson(json['team'] ?? {}),
      terminal: TerminalStats.fromJson(json['terminal'] ?? {}),
      wallet: WalletStats.fromJson(json['wallet'] ?? {}),
    );
  }

  /// 计算较昨日变化率
  double get profitChangeRate {
    if (yesterday.profitTotal == 0) return 0;
    return ((today.profitTotal - yesterday.profitTotal) / yesterday.profitTotal) * 100;
  }

  /// 较昨日是否增长
  bool get isProfitGrowth => today.profitTotal >= yesterday.profitTotal;
}

/// 今日统计
class TodayStats {
  final int transAmount;
  final double transAmountYuan;
  final int transCount;
  final int profitTotal;
  final double profitTotalYuan;
  final int profitTrade;
  final int profitDeposit;
  final int profitSim;
  final int profitReward;

  TodayStats({
    required this.transAmount,
    required this.transAmountYuan,
    required this.transCount,
    required this.profitTotal,
    required this.profitTotalYuan,
    required this.profitTrade,
    required this.profitDeposit,
    required this.profitSim,
    required this.profitReward,
  });

  factory TodayStats.fromJson(Map<String, dynamic> json) {
    return TodayStats(
      transAmount: json['trans_amount'] ?? 0,
      transAmountYuan: (json['trans_amount_yuan'] ?? 0).toDouble(),
      transCount: json['trans_count'] ?? 0,
      profitTotal: json['profit_total'] ?? 0,
      profitTotalYuan: (json['profit_total_yuan'] ?? 0).toDouble(),
      profitTrade: json['profit_trade'] ?? 0,
      profitDeposit: json['profit_deposit'] ?? 0,
      profitSim: json['profit_sim'] ?? 0,
      profitReward: json['profit_reward'] ?? 0,
    );
  }

  /// 分润分类(元)
  double get profitTradeYuan => profitTrade / 100;
  double get profitDepositYuan => profitDeposit / 100;
  double get profitSimYuan => profitSim / 100;
  double get profitRewardYuan => profitReward / 100;
}

/// 月度统计
class MonthStats {
  final int transAmount;
  final double transAmountYuan;
  final int transCount;
  final int profitTotal;
  final double profitTotalYuan;
  final int merchantNew;

  MonthStats({
    required this.transAmount,
    required this.transAmountYuan,
    required this.transCount,
    required this.profitTotal,
    required this.profitTotalYuan,
    required this.merchantNew,
  });

  factory MonthStats.fromJson(Map<String, dynamic> json) {
    return MonthStats(
      transAmount: json['trans_amount'] ?? 0,
      transAmountYuan: (json['trans_amount_yuan'] ?? 0).toDouble(),
      transCount: json['trans_count'] ?? 0,
      profitTotal: json['profit_total'] ?? 0,
      profitTotalYuan: (json['profit_total_yuan'] ?? 0).toDouble(),
      merchantNew: json['merchant_new'] ?? 0,
    );
  }
}

/// 团队统计
class TeamStats {
  final int directAgentCount;
  final int directMerchantCount;
  final int teamAgentCount;
  final int teamMerchantCount;

  TeamStats({
    required this.directAgentCount,
    required this.directMerchantCount,
    required this.teamAgentCount,
    required this.teamMerchantCount,
  });

  factory TeamStats.fromJson(Map<String, dynamic> json) {
    return TeamStats(
      directAgentCount: json['direct_agent_count'] ?? 0,
      directMerchantCount: json['direct_merchant_count'] ?? 0,
      teamAgentCount: json['team_agent_count'] ?? 0,
      teamMerchantCount: json['team_merchant_count'] ?? 0,
    );
  }
}

/// 终端统计
class TerminalStats {
  final int total;
  final int activated;
  final int todayActivated;
  final int monthActivated;

  TerminalStats({
    required this.total,
    required this.activated,
    required this.todayActivated,
    required this.monthActivated,
  });

  factory TerminalStats.fromJson(Map<String, dynamic> json) {
    return TerminalStats(
      total: json['total'] ?? 0,
      activated: json['activated'] ?? 0,
      todayActivated: json['today_activated'] ?? 0,
      monthActivated: json['month_activated'] ?? 0,
    );
  }
}

/// 钱包统计
class WalletStats {
  final int totalBalance;
  final double totalBalanceYuan;

  WalletStats({
    required this.totalBalance,
    required this.totalBalanceYuan,
  });

  factory WalletStats.fromJson(Map<String, dynamic> json) {
    return WalletStats(
      totalBalance: json['total_balance'] ?? 0,
      totalBalanceYuan: (json['total_balance_yuan'] ?? 0).toDouble(),
    );
  }
}

/// 最近交易
class RecentTransaction {
  final int id;
  final String merchantName;
  final int payType;
  final String payTypeName;
  final int amount;
  final double amountYuan;
  final DateTime tradeTime;
  final String timeAgo;

  RecentTransaction({
    required this.id,
    required this.merchantName,
    required this.payType,
    required this.payTypeName,
    required this.amount,
    required this.amountYuan,
    required this.tradeTime,
    required this.timeAgo,
  });

  factory RecentTransaction.fromJson(Map<String, dynamic> json) {
    return RecentTransaction(
      id: json['id'] ?? 0,
      merchantName: json['merchant_name'] ?? '',
      payType: json['pay_type'] ?? 0,
      payTypeName: json['pay_type_name'] ?? '',
      amount: json['amount'] ?? 0,
      amountYuan: (json['amount_yuan'] ?? 0).toDouble(),
      tradeTime: DateTime.tryParse(json['trade_time'] ?? '') ?? DateTime.now(),
      timeAgo: json['time_ago'] ?? '',
    );
  }
}
