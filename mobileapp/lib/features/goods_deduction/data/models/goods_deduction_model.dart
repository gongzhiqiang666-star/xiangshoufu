/// 货款代扣数据模型

// 货款代扣状态
enum GoodsDeductionStatus {
  pendingAccept(1, '待接收'),
  inProgress(2, '进行中'),
  completed(3, '已完成'),
  rejected(4, '已拒绝');

  final int value;
  final String label;
  const GoodsDeductionStatus(this.value, this.label);

  static GoodsDeductionStatus fromValue(int value) {
    return GoodsDeductionStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => GoodsDeductionStatus.pendingAccept,
    );
  }
}

// 扣款来源
enum DeductionSource {
  profit(1, '分润钱包'),
  serviceFee(2, '服务费钱包'),
  both(3, '分润+服务费');

  final int value;
  final String label;
  const DeductionSource(this.value, this.label);

  static DeductionSource fromValue(int value) {
    return DeductionSource.values.firstWhere(
      (e) => e.value == value,
      orElse: () => DeductionSource.both,
    );
  }
}

// 货款代扣模型
class GoodsDeduction {
  final int id;
  final String deductionNo;
  final int fromAgentId;
  final String fromAgentName;
  final int toAgentId;
  final String toAgentName;
  final int totalAmount; // 分
  final int deductedAmount; // 分
  final int remainingAmount; // 分
  final int deductionSource;
  final String sourceName;
  final int terminalCount;
  final int unitPrice; // 分
  final int status;
  final String statusName;
  final double progress;
  final bool agreementSigned;
  final String? agreementUrl;
  final int? distributeId;
  final String? remark;
  final String createdAt;
  final String? acceptedAt;
  final String? completedAt;

  GoodsDeduction({
    required this.id,
    required this.deductionNo,
    required this.fromAgentId,
    required this.fromAgentName,
    required this.toAgentId,
    required this.toAgentName,
    required this.totalAmount,
    required this.deductedAmount,
    required this.remainingAmount,
    required this.deductionSource,
    required this.sourceName,
    required this.terminalCount,
    required this.unitPrice,
    required this.status,
    required this.statusName,
    required this.progress,
    required this.agreementSigned,
    this.agreementUrl,
    this.distributeId,
    this.remark,
    required this.createdAt,
    this.acceptedAt,
    this.completedAt,
  });

  factory GoodsDeduction.fromJson(Map<String, dynamic> json) {
    return GoodsDeduction(
      id: json['id'] ?? 0,
      deductionNo: json['deduction_no'] ?? '',
      fromAgentId: json['from_agent_id'] ?? 0,
      fromAgentName: json['from_agent_name'] ?? '',
      toAgentId: json['to_agent_id'] ?? 0,
      toAgentName: json['to_agent_name'] ?? '',
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
      deductionSource: json['deduction_source'] ?? 3,
      sourceName: json['source_name'] ?? '',
      terminalCount: json['terminal_count'] ?? 0,
      unitPrice: json['unit_price'] ?? 0,
      status: json['status'] ?? 1,
      statusName: json['status_name'] ?? '',
      progress: (json['progress'] ?? 0).toDouble(),
      agreementSigned: json['agreement_signed'] ?? false,
      agreementUrl: json['agreement_url'],
      distributeId: json['distribute_id'],
      remark: json['remark'],
      createdAt: json['created_at'] ?? '',
      acceptedAt: json['accepted_at'],
      completedAt: json['completed_at'],
    );
  }

  // 获取总金额（元）
  double get totalAmountYuan => totalAmount / 100;

  // 获取已扣金额（元）
  double get deductedAmountYuan => deductedAmount / 100;

  // 获取剩余金额（元）
  double get remainingAmountYuan => remainingAmount / 100;

  // 获取单价（元）
  double get unitPriceYuan => unitPrice / 100;

  // 获取状态枚举
  GoodsDeductionStatus get statusEnum => GoodsDeductionStatus.fromValue(status);

  // 获取来源枚举
  DeductionSource get sourceEnum => DeductionSource.fromValue(deductionSource);
}

// 货款代扣详情模型
class GoodsDeductionDetail extends GoodsDeduction {
  final List<GoodsDeductionTerminal> terminals;
  final List<GoodsDeductionDetailRecord> details;

  GoodsDeductionDetail({
    required super.id,
    required super.deductionNo,
    required super.fromAgentId,
    required super.fromAgentName,
    required super.toAgentId,
    required super.toAgentName,
    required super.totalAmount,
    required super.deductedAmount,
    required super.remainingAmount,
    required super.deductionSource,
    required super.sourceName,
    required super.terminalCount,
    required super.unitPrice,
    required super.status,
    required super.statusName,
    required super.progress,
    required super.agreementSigned,
    super.agreementUrl,
    super.distributeId,
    super.remark,
    required super.createdAt,
    super.acceptedAt,
    super.completedAt,
    this.terminals = const [],
    this.details = const [],
  });

  factory GoodsDeductionDetail.fromJson(Map<String, dynamic> json) {
    return GoodsDeductionDetail(
      id: json['id'] ?? 0,
      deductionNo: json['deduction_no'] ?? '',
      fromAgentId: json['from_agent_id'] ?? 0,
      fromAgentName: json['from_agent_name'] ?? '',
      toAgentId: json['to_agent_id'] ?? 0,
      toAgentName: json['to_agent_name'] ?? '',
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
      deductionSource: json['deduction_source'] ?? 3,
      sourceName: json['source_name'] ?? '',
      terminalCount: json['terminal_count'] ?? 0,
      unitPrice: json['unit_price'] ?? 0,
      status: json['status'] ?? 1,
      statusName: json['status_name'] ?? '',
      progress: (json['progress'] ?? 0).toDouble(),
      agreementSigned: json['agreement_signed'] ?? false,
      agreementUrl: json['agreement_url'],
      distributeId: json['distribute_id'],
      remark: json['remark'],
      createdAt: json['created_at'] ?? '',
      acceptedAt: json['accepted_at'],
      completedAt: json['completed_at'],
      terminals: (json['terminals'] as List<dynamic>?)
              ?.map((e) => GoodsDeductionTerminal.fromJson(e))
              .toList() ??
          [],
      details: (json['details'] as List<dynamic>?)
              ?.map((e) => GoodsDeductionDetailRecord.fromJson(e))
              .toList() ??
          [],
    );
  }
}

// 货款代扣终端
class GoodsDeductionTerminal {
  final int id;
  final int deductionId;
  final int terminalId;
  final String terminalSn;
  final int unitPrice;
  final String createdAt;

  GoodsDeductionTerminal({
    required this.id,
    required this.deductionId,
    required this.terminalId,
    required this.terminalSn,
    required this.unitPrice,
    required this.createdAt,
  });

  factory GoodsDeductionTerminal.fromJson(Map<String, dynamic> json) {
    return GoodsDeductionTerminal(
      id: json['id'] ?? 0,
      deductionId: json['deduction_id'] ?? 0,
      terminalId: json['terminal_id'] ?? 0,
      terminalSn: json['terminal_sn'] ?? '',
      unitPrice: json['unit_price'] ?? 0,
      createdAt: json['created_at'] ?? '',
    );
  }

  double get unitPriceYuan => unitPrice / 100;
}

// 货款代扣扣款明细
class GoodsDeductionDetailRecord {
  final int id;
  final int deductionId;
  final String deductionNo;
  final int amount;
  final int walletType;
  final String walletTypeName;
  final int? channelId;
  final String? channelName;
  final int walletBalanceBefore;
  final int walletBalanceAfter;
  final int cumulativeDeducted;
  final int remainingAfter;
  final String triggerType;
  final String createdAt;

  GoodsDeductionDetailRecord({
    required this.id,
    required this.deductionId,
    required this.deductionNo,
    required this.amount,
    required this.walletType,
    required this.walletTypeName,
    this.channelId,
    this.channelName,
    required this.walletBalanceBefore,
    required this.walletBalanceAfter,
    required this.cumulativeDeducted,
    required this.remainingAfter,
    required this.triggerType,
    required this.createdAt,
  });

  factory GoodsDeductionDetailRecord.fromJson(Map<String, dynamic> json) {
    return GoodsDeductionDetailRecord(
      id: json['id'] ?? 0,
      deductionId: json['deduction_id'] ?? 0,
      deductionNo: json['deduction_no'] ?? '',
      amount: json['amount'] ?? 0,
      walletType: json['wallet_type'] ?? 0,
      walletTypeName: json['wallet_type_name'] ?? '',
      channelId: json['channel_id'],
      channelName: json['channel_name'],
      walletBalanceBefore: json['wallet_balance_before'] ?? 0,
      walletBalanceAfter: json['wallet_balance_after'] ?? 0,
      cumulativeDeducted: json['cumulative_deducted'] ?? 0,
      remainingAfter: json['remaining_after'] ?? 0,
      triggerType: json['trigger_type'] ?? '',
      createdAt: json['created_at'] ?? '',
    );
  }

  double get amountYuan => amount / 100;
  double get walletBalanceBeforeYuan => walletBalanceBefore / 100;
  double get walletBalanceAfterYuan => walletBalanceAfter / 100;
  double get cumulativeDeductedYuan => cumulativeDeducted / 100;
  double get remainingAfterYuan => remainingAfter / 100;
}

// 货款代扣统计汇总
class GoodsDeductionSummary {
  final int totalCount;
  final int pendingCount;
  final int inProgressCount;
  final int completedCount;
  final int totalAmount;
  final int deductedAmount;
  final int remainingAmount;

  GoodsDeductionSummary({
    required this.totalCount,
    required this.pendingCount,
    required this.inProgressCount,
    required this.completedCount,
    required this.totalAmount,
    required this.deductedAmount,
    required this.remainingAmount,
  });

  factory GoodsDeductionSummary.fromJson(Map<String, dynamic> json) {
    return GoodsDeductionSummary(
      totalCount: json['total_count'] ?? 0,
      pendingCount: json['pending_count'] ?? 0,
      inProgressCount: json['in_progress_count'] ?? 0,
      completedCount: json['completed_count'] ?? 0,
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
    );
  }

  double get totalAmountYuan => totalAmount / 100;
  double get deductedAmountYuan => deductedAmount / 100;
  double get remainingAmountYuan => remainingAmount / 100;
}

// 货款代扣列表响应
class GoodsDeductionListResponse {
  final List<GoodsDeduction> list;
  final int total;
  final int page;
  final int pageSize;

  GoodsDeductionListResponse({
    required this.list,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory GoodsDeductionListResponse.fromJson(Map<String, dynamic> json) {
    return GoodsDeductionListResponse(
      list: (json['list'] as List<dynamic>?)
              ?.map((e) => GoodsDeduction.fromJson(e))
              .toList() ??
          [],
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      pageSize: json['page_size'] ?? 10,
    );
  }
}
