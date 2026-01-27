/// 代扣计划数据模型

// 代扣计划状态
enum DeductionPlanStatus {
  pendingAccept(0, '待接收'),
  active(1, '进行中'),
  completed(2, '已完成'),
  paused(3, '已暂停'),
  cancelled(4, '已取消'),
  rejected(5, '已拒绝');

  final int value;
  final String label;
  const DeductionPlanStatus(this.value, this.label);

  static DeductionPlanStatus fromValue(int value) {
    return DeductionPlanStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => DeductionPlanStatus.active,
    );
  }

  /// 是否可以接收确认
  bool get canAccept => this == DeductionPlanStatus.pendingAccept;

  /// 是否可以拒绝
  bool get canReject => this == DeductionPlanStatus.pendingAccept;

  /// 是否可以暂停
  bool get canPause => this == DeductionPlanStatus.active;

  /// 是否可以恢复
  bool get canResume => this == DeductionPlanStatus.paused;

  /// 是否可以取消
  bool get canCancel => this == DeductionPlanStatus.active || this == DeductionPlanStatus.paused;
}

// 代扣来源
enum DeductionSource {
  profit(1, '分润'),
  serviceFee(2, '服务费'),
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

// 代扣计划类型
enum DeductionPlanType {
  goods(1, '货款代扣'),
  partner(2, '伙伴代扣'),
  deposit(3, '押金代扣');

  final int value;
  final String label;
  const DeductionPlanType(this.value, this.label);

  static DeductionPlanType fromValue(int value) {
    return DeductionPlanType.values.firstWhere(
      (e) => e.value == value,
      orElse: () => DeductionPlanType.partner,
    );
  }
}

// 代扣记录状态
enum DeductionRecordStatus {
  pending(0, '待扣款'),
  success(1, '成功'),
  partialSuccess(2, '部分成功'),
  failed(3, '失败');

  final int value;
  final String label;
  const DeductionRecordStatus(this.value, this.label);

  static DeductionRecordStatus fromValue(int value) {
    return DeductionRecordStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => DeductionRecordStatus.pending,
    );
  }
}

// 代扣计划模型
class DeductionPlan {
  final int id;
  final String planNo;
  final int deductorId;
  final String deductorName;
  final int deducteeId;
  final String deducteeName;
  final int planType;
  final int totalAmount; // 分
  final int deductedAmount; // 分
  final int remainingAmount; // 分
  final int frozenAmount; // 已冻结金额（分）
  final int totalPeriods;
  final int currentPeriod;
  final int periodAmount; // 分
  final int status;
  final bool needAccept; // 是否需要接收确认
  final String? acceptedAt; // 接收时间
  final int deductionSource; // 扣款来源：1=分润 2=服务费 3=两者
  final String? relatedType;
  final int? relatedId;
  final String? remark;
  final int createdBy;
  final String createdAt;
  final String updatedAt;
  final String? completedAt;

  DeductionPlan({
    required this.id,
    required this.planNo,
    required this.deductorId,
    required this.deductorName,
    required this.deducteeId,
    required this.deducteeName,
    required this.planType,
    required this.totalAmount,
    required this.deductedAmount,
    required this.remainingAmount,
    this.frozenAmount = 0,
    required this.totalPeriods,
    required this.currentPeriod,
    required this.periodAmount,
    required this.status,
    this.needAccept = false,
    this.acceptedAt,
    this.deductionSource = 3,
    this.relatedType,
    this.relatedId,
    this.remark,
    required this.createdBy,
    required this.createdAt,
    required this.updatedAt,
    this.completedAt,
  });

  factory DeductionPlan.fromJson(Map<String, dynamic> json) {
    return DeductionPlan(
      id: json['id'] ?? 0,
      planNo: json['plan_no'] ?? '',
      deductorId: json['deductor_id'] ?? 0,
      deductorName: json['deductor_name'] ?? '',
      deducteeId: json['deductee_id'] ?? 0,
      deducteeName: json['deductee_name'] ?? '',
      planType: json['plan_type'] ?? 2,
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
      frozenAmount: json['frozen_amount'] ?? 0,
      totalPeriods: json['total_periods'] ?? 0,
      currentPeriod: json['current_period'] ?? 0,
      periodAmount: json['period_amount'] ?? 0,
      status: json['status'] ?? 1,
      needAccept: json['need_accept'] ?? false,
      acceptedAt: json['accepted_at'],
      deductionSource: json['deduction_source'] ?? 3,
      relatedType: json['related_type'],
      relatedId: json['related_id'],
      remark: json['remark'],
      createdBy: json['created_by'] ?? 0,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      completedAt: json['completed_at'],
    );
  }

  // 获取总金额（元）
  double get totalAmountYuan => totalAmount / 100;

  // 获取已扣金额（元）
  double get deductedAmountYuan => deductedAmount / 100;

  // 获取剩余金额（元）
  double get remainingAmountYuan => remainingAmount / 100;

  // 获取每期金额（元）
  double get periodAmountYuan => periodAmount / 100;

  // 获取状态枚举
  DeductionPlanStatus get statusEnum => DeductionPlanStatus.fromValue(status);

  // 获取类型枚举
  DeductionPlanType get typeEnum => DeductionPlanType.fromValue(planType);

  // 获取进度百分比
  double get progress {
    if (totalAmount <= 0) return 0;
    return (deductedAmount / totalAmount) * 100;
  }

  // 获取已冻结金额（元）
  double get frozenAmountYuan => frozenAmount / 100;

  // 获取扣款来源枚举
  DeductionSource get deductionSourceEnum => DeductionSource.fromValue(deductionSource);

  // 是否需要接收确认
  bool get isPendingAccept => statusEnum == DeductionPlanStatus.pendingAccept;

  // 是否已被拒绝
  bool get isRejected => statusEnum == DeductionPlanStatus.rejected;
}

// 代扣计划详情
class DeductionPlanDetail extends DeductionPlan {
  final List<DeductionRecord> records;

  DeductionPlanDetail({
    required super.id,
    required super.planNo,
    required super.deductorId,
    required super.deductorName,
    required super.deducteeId,
    required super.deducteeName,
    required super.planType,
    required super.totalAmount,
    required super.deductedAmount,
    required super.remainingAmount,
    super.frozenAmount,
    required super.totalPeriods,
    required super.currentPeriod,
    required super.periodAmount,
    required super.status,
    super.needAccept,
    super.acceptedAt,
    super.deductionSource,
    super.relatedType,
    super.relatedId,
    super.remark,
    required super.createdBy,
    required super.createdAt,
    required super.updatedAt,
    super.completedAt,
    this.records = const [],
  });

  factory DeductionPlanDetail.fromJson(Map<String, dynamic> json) {
    return DeductionPlanDetail(
      id: json['id'] ?? 0,
      planNo: json['plan_no'] ?? '',
      deductorId: json['deductor_id'] ?? 0,
      deductorName: json['deductor_name'] ?? '',
      deducteeId: json['deductee_id'] ?? 0,
      deducteeName: json['deductee_name'] ?? '',
      planType: json['plan_type'] ?? 2,
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
      frozenAmount: json['frozen_amount'] ?? 0,
      totalPeriods: json['total_periods'] ?? 0,
      currentPeriod: json['current_period'] ?? 0,
      periodAmount: json['period_amount'] ?? 0,
      status: json['status'] ?? 1,
      needAccept: json['need_accept'] ?? false,
      acceptedAt: json['accepted_at'],
      deductionSource: json['deduction_source'] ?? 3,
      relatedType: json['related_type'],
      relatedId: json['related_id'],
      remark: json['remark'],
      createdBy: json['created_by'] ?? 0,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      completedAt: json['completed_at'],
      records: (json['records'] as List<dynamic>?)
              ?.map((e) => DeductionRecord.fromJson(e))
              .toList() ??
          [],
    );
  }
}

// 代扣记录
class DeductionRecord {
  final int id;
  final int planId;
  final String planNo;
  final int deductorId;
  final int deducteeId;
  final int periodNum;
  final int amount; // 应扣金额（分）
  final int actualAmount; // 实扣金额（分）
  final int status;
  final String? failReason;
  final String scheduledAt;
  final String? deductedAt;
  final String createdAt;

  DeductionRecord({
    required this.id,
    required this.planId,
    required this.planNo,
    required this.deductorId,
    required this.deducteeId,
    required this.periodNum,
    required this.amount,
    required this.actualAmount,
    required this.status,
    this.failReason,
    required this.scheduledAt,
    this.deductedAt,
    required this.createdAt,
  });

  factory DeductionRecord.fromJson(Map<String, dynamic> json) {
    return DeductionRecord(
      id: json['id'] ?? 0,
      planId: json['plan_id'] ?? 0,
      planNo: json['plan_no'] ?? '',
      deductorId: json['deductor_id'] ?? 0,
      deducteeId: json['deductee_id'] ?? 0,
      periodNum: json['period_num'] ?? 0,
      amount: json['amount'] ?? 0,
      actualAmount: json['actual_amount'] ?? 0,
      status: json['status'] ?? 0,
      failReason: json['fail_reason'],
      scheduledAt: json['scheduled_at'] ?? '',
      deductedAt: json['deducted_at'],
      createdAt: json['created_at'] ?? '',
    );
  }

  double get amountYuan => amount / 100;
  double get actualAmountYuan => actualAmount / 100;

  DeductionRecordStatus get statusEnum => DeductionRecordStatus.fromValue(status);
}

// 代扣计划列表响应
class DeductionPlanListResponse {
  final List<DeductionPlan> list;
  final int total;
  final int page;
  final int pageSize;

  DeductionPlanListResponse({
    required this.list,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory DeductionPlanListResponse.fromJson(Map<String, dynamic> json) {
    return DeductionPlanListResponse(
      list: (json['list'] as List<dynamic>?)
              ?.map((e) => DeductionPlan.fromJson(e))
              .toList() ??
          [],
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      pageSize: json['page_size'] ?? 10,
    );
  }
}

// 代扣计划统计
class DeductionPlanStats {
  final int totalCount;
  final int pendingAcceptCount; // 待接收数量
  final int activeCount;
  final int completedCount;
  final int pausedCount;
  final int rejectedCount; // 已拒绝数量
  final int totalAmount;
  final int deductedAmount;
  final int remainingAmount;
  final int frozenAmount; // 冻结金额

  DeductionPlanStats({
    required this.totalCount,
    this.pendingAcceptCount = 0,
    required this.activeCount,
    required this.completedCount,
    required this.pausedCount,
    this.rejectedCount = 0,
    required this.totalAmount,
    required this.deductedAmount,
    required this.remainingAmount,
    this.frozenAmount = 0,
  });

  factory DeductionPlanStats.fromJson(Map<String, dynamic> json) {
    return DeductionPlanStats(
      totalCount: json['total_count'] ?? 0,
      pendingAcceptCount: json['pending_accept_count'] ?? 0,
      activeCount: json['active_count'] ?? 0,
      completedCount: json['completed_count'] ?? 0,
      pausedCount: json['paused_count'] ?? 0,
      rejectedCount: json['rejected_count'] ?? 0,
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
      frozenAmount: json['frozen_amount'] ?? 0,
    );
  }

  double get totalAmountYuan => totalAmount / 100;
  double get deductedAmountYuan => deductedAmount / 100;
  double get remainingAmountYuan => remainingAmount / 100;
  double get frozenAmountYuan => frozenAmount / 100;
}

// 代扣摘要（我接收的/我发起的统计）
class DeductionSummary {
  final int receivedPendingCount; // 待接收数量
  final int receivedActiveCount; // 进行中数量
  final int receivedTotalAmount; // 待扣总额
  final int sentPendingCount; // 发起待确认数量
  final int sentActiveCount; // 发起进行中数量
  final int sentTotalAmount; // 发起总额

  DeductionSummary({
    required this.receivedPendingCount,
    required this.receivedActiveCount,
    required this.receivedTotalAmount,
    required this.sentPendingCount,
    required this.sentActiveCount,
    required this.sentTotalAmount,
  });

  factory DeductionSummary.fromJson(Map<String, dynamic> json) {
    return DeductionSummary(
      receivedPendingCount: json['received_pending_count'] ?? 0,
      receivedActiveCount: json['received_active_count'] ?? 0,
      receivedTotalAmount: json['received_total_amount'] ?? 0,
      sentPendingCount: json['sent_pending_count'] ?? 0,
      sentActiveCount: json['sent_active_count'] ?? 0,
      sentTotalAmount: json['sent_total_amount'] ?? 0,
    );
  }

  double get receivedTotalAmountYuan => receivedTotalAmount / 100;
  double get sentTotalAmountYuan => sentTotalAmount / 100;
}
