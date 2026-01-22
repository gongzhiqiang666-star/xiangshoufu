/// 代扣计划数据模型

// 代扣计划状态
enum DeductionPlanStatus {
  active(1, '进行中'),
  completed(2, '已完成'),
  paused(3, '已暂停'),
  cancelled(4, '已取消');

  final int value;
  final String label;
  const DeductionPlanStatus(this.value, this.label);

  static DeductionPlanStatus fromValue(int value) {
    return DeductionPlanStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => DeductionPlanStatus.active,
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
  final int totalPeriods;
  final int currentPeriod;
  final int periodAmount; // 分
  final int status;
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
    required this.totalPeriods,
    required this.currentPeriod,
    required this.periodAmount,
    required this.status,
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
      totalPeriods: json['total_periods'] ?? 0,
      currentPeriod: json['current_period'] ?? 0,
      periodAmount: json['period_amount'] ?? 0,
      status: json['status'] ?? 1,
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
    required super.totalPeriods,
    required super.currentPeriod,
    required super.periodAmount,
    required super.status,
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
      totalPeriods: json['total_periods'] ?? 0,
      currentPeriod: json['current_period'] ?? 0,
      periodAmount: json['period_amount'] ?? 0,
      status: json['status'] ?? 1,
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
  final int activeCount;
  final int completedCount;
  final int pausedCount;
  final int totalAmount;
  final int deductedAmount;
  final int remainingAmount;

  DeductionPlanStats({
    required this.totalCount,
    required this.activeCount,
    required this.completedCount,
    required this.pausedCount,
    required this.totalAmount,
    required this.deductedAmount,
    required this.remainingAmount,
  });

  factory DeductionPlanStats.fromJson(Map<String, dynamic> json) {
    return DeductionPlanStats(
      totalCount: json['total_count'] ?? 0,
      activeCount: json['active_count'] ?? 0,
      completedCount: json['completed_count'] ?? 0,
      pausedCount: json['paused_count'] ?? 0,
      totalAmount: json['total_amount'] ?? 0,
      deductedAmount: json['deducted_amount'] ?? 0,
      remainingAmount: json['remaining_amount'] ?? 0,
    );
  }

  double get totalAmountYuan => totalAmount / 100;
  double get deductedAmountYuan => deductedAmount / 100;
  double get remainingAmountYuan => remainingAmount / 100;
}
