/// 终端状态枚举
enum TerminalStatus {
  pending(1, '待分配'),
  allocated(2, '已分配'),
  bound(3, '已绑定'),
  activated(4, '已激活'),
  unbound(5, '已解绑'),
  recycled(6, '已回收');

  final int value;
  final String label;
  const TerminalStatus(this.value, this.label);

  static TerminalStatus fromValue(int value) {
    return TerminalStatus.values.firstWhere(
      (e) => e.value == value,
      orElse: () => TerminalStatus.pending,
    );
  }
}

/// 终端信息
class Terminal {
  final int id;
  final String terminalSn;
  final int channelId;
  final String channelCode;
  final String? brandCode;
  final String? modelCode;
  final int ownerAgentId;
  final int? merchantId;
  final String? merchantNo;
  final TerminalStatus status;
  final DateTime? activatedAt;
  final DateTime? boundAt;
  final int simFeeCount;
  final DateTime? lastSimFeeAt;
  final DateTime createdAt;
  final DateTime updatedAt;

  Terminal({
    required this.id,
    required this.terminalSn,
    required this.channelId,
    required this.channelCode,
    this.brandCode,
    this.modelCode,
    required this.ownerAgentId,
    this.merchantId,
    this.merchantNo,
    required this.status,
    this.activatedAt,
    this.boundAt,
    this.simFeeCount = 0,
    this.lastSimFeeAt,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Terminal.fromJson(Map<String, dynamic> json) {
    return Terminal(
      id: json['id'] ?? 0,
      terminalSn: json['terminal_sn'] ?? '',
      channelId: json['channel_id'] ?? 0,
      channelCode: json['channel_code'] ?? '',
      brandCode: json['brand_code'],
      modelCode: json['model_code'],
      ownerAgentId: json['owner_agent_id'] ?? 0,
      merchantId: json['merchant_id'],
      merchantNo: json['merchant_no'],
      status: TerminalStatus.fromValue(json['status'] ?? 1),
      activatedAt: json['activated_at'] != null
          ? DateTime.tryParse(json['activated_at'])
          : null,
      boundAt: json['bound_at'] != null
          ? DateTime.tryParse(json['bound_at'])
          : null,
      simFeeCount: json['sim_fee_count'] ?? 0,
      lastSimFeeAt: json['last_sim_fee_at'] != null
          ? DateTime.tryParse(json['last_sim_fee_at'])
          : null,
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      updatedAt: DateTime.tryParse(json['updated_at'] ?? '') ?? DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'terminal_sn': terminalSn,
      'channel_id': channelId,
      'channel_code': channelCode,
      'brand_code': brandCode,
      'model_code': modelCode,
      'owner_agent_id': ownerAgentId,
      'merchant_id': merchantId,
      'merchant_no': merchantNo,
      'status': status.value,
      'activated_at': activatedAt?.toIso8601String(),
      'bound_at': boundAt?.toIso8601String(),
      'sim_fee_count': simFeeCount,
      'last_sim_fee_at': lastSimFeeAt?.toIso8601String(),
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }

  /// 是否已激活
  bool get isActivated => status == TerminalStatus.activated;

  /// 是否可以回拨（只有未激活的终端可以回拨）
  bool get canRecall => status != TerminalStatus.activated;

  /// 是否可以下发（只有待分配/已分配的终端可以下发）
  bool get canDistribute =>
      status == TerminalStatus.pending || status == TerminalStatus.allocated;
}

/// 终端统计
class TerminalStats {
  final int total;
  final int pendingCount;
  final int allocatedCount;
  final int boundCount;
  final int activatedCount;
  final int unboundCount;
  final int yesterdayActivated;
  final int todayActivated;
  final int monthActivated;

  TerminalStats({
    required this.total,
    required this.pendingCount,
    required this.allocatedCount,
    required this.boundCount,
    required this.activatedCount,
    required this.unboundCount,
    required this.yesterdayActivated,
    required this.todayActivated,
    required this.monthActivated,
  });

  factory TerminalStats.fromJson(Map<String, dynamic> json) {
    return TerminalStats(
      total: json['total'] ?? 0,
      pendingCount: json['pending_count'] ?? 0,
      allocatedCount: json['allocated_count'] ?? 0,
      boundCount: json['bound_count'] ?? 0,
      activatedCount: json['activated_count'] ?? 0,
      unboundCount: json['unbound_count'] ?? 0,
      yesterdayActivated: json['yesterday_activated'] ?? 0,
      todayActivated: json['today_activated'] ?? 0,
      monthActivated: json['month_activated'] ?? 0,
    );
  }

  /// 未激活数量
  int get inactiveCount => total - activatedCount;

  /// 库存数量（未绑定）
  int get stockCount => pendingCount + allocatedCount;
}

/// 终端下发记录
class TerminalDistribute {
  final int id;
  final String distributeNo;
  final int fromAgentId;
  final int toAgentId;
  final String terminalSn;
  final int channelId;
  final bool isCrossLevel;
  final String? crossLevelPath;
  final int goodsPrice;
  final int deductionType;
  final int status;
  final int source;
  final String? remark;
  final DateTime createdAt;
  final DateTime? confirmedAt;

  TerminalDistribute({
    required this.id,
    required this.distributeNo,
    required this.fromAgentId,
    required this.toAgentId,
    required this.terminalSn,
    required this.channelId,
    required this.isCrossLevel,
    this.crossLevelPath,
    required this.goodsPrice,
    required this.deductionType,
    required this.status,
    required this.source,
    this.remark,
    required this.createdAt,
    this.confirmedAt,
  });

  factory TerminalDistribute.fromJson(Map<String, dynamic> json) {
    return TerminalDistribute(
      id: json['id'] ?? 0,
      distributeNo: json['distribute_no'] ?? '',
      fromAgentId: json['from_agent_id'] ?? 0,
      toAgentId: json['to_agent_id'] ?? 0,
      terminalSn: json['terminal_sn'] ?? '',
      channelId: json['channel_id'] ?? 0,
      isCrossLevel: json['is_cross_level'] ?? false,
      crossLevelPath: json['cross_level_path'],
      goodsPrice: json['goods_price'] ?? 0,
      deductionType: json['deduction_type'] ?? 1,
      status: json['status'] ?? 1,
      source: json['source'] ?? 1,
      remark: json['remark'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      confirmedAt: json['confirmed_at'] != null
          ? DateTime.tryParse(json['confirmed_at'])
          : null,
    );
  }

  String get statusLabel {
    switch (status) {
      case 1:
        return '待确认';
      case 2:
        return '已确认';
      case 3:
        return '已拒绝';
      case 4:
        return '已取消';
      default:
        return '未知';
    }
  }
}

/// 终端回拨记录
class TerminalRecall {
  final int id;
  final String recallNo;
  final int fromAgentId;
  final int toAgentId;
  final String terminalSn;
  final int channelId;
  final bool isCrossLevel;
  final String? crossLevelPath;
  final int status;
  final int source;
  final String? remark;
  final DateTime createdAt;
  final DateTime? confirmedAt;

  TerminalRecall({
    required this.id,
    required this.recallNo,
    required this.fromAgentId,
    required this.toAgentId,
    required this.terminalSn,
    required this.channelId,
    required this.isCrossLevel,
    this.crossLevelPath,
    required this.status,
    required this.source,
    this.remark,
    required this.createdAt,
    this.confirmedAt,
  });

  factory TerminalRecall.fromJson(Map<String, dynamic> json) {
    return TerminalRecall(
      id: json['id'] ?? 0,
      recallNo: json['recall_no'] ?? '',
      fromAgentId: json['from_agent_id'] ?? 0,
      toAgentId: json['to_agent_id'] ?? 0,
      terminalSn: json['terminal_sn'] ?? '',
      channelId: json['channel_id'] ?? 0,
      isCrossLevel: json['is_cross_level'] ?? false,
      crossLevelPath: json['cross_level_path'],
      status: json['status'] ?? 1,
      source: json['source'] ?? 1,
      remark: json['remark'],
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      confirmedAt: json['confirmed_at'] != null
          ? DateTime.tryParse(json['confirmed_at'])
          : null,
    );
  }

  String get statusLabel {
    switch (status) {
      case 1:
        return '待确认';
      case 2:
        return '已确认';
      case 3:
        return '已拒绝';
      case 4:
        return '已取消';
      default:
        return '未知';
    }
  }
}

/// 终端状态分组
enum TerminalStatusGroup {
  all('all', '全部'),
  unstock('unstock', '未出库'),
  stocked('stocked', '已出库'),
  unbound('unbound', '未绑定'),
  inactive('inactive', '未激活'),
  active('active', '已激活');

  final String value;
  final String label;
  const TerminalStatusGroup(this.value, this.label);

  static TerminalStatusGroup fromValue(String value) {
    return TerminalStatusGroup.values.firstWhere(
      (e) => e.value == value,
      orElse: () => TerminalStatusGroup.all,
    );
  }
}

/// 终端筛选选项
class TerminalFilterOptions {
  final List<ChannelOption> channels;
  final List<TerminalTypeOption> terminalTypes;
  final List<StatusGroupCount> statusGroups;

  TerminalFilterOptions({
    required this.channels,
    required this.terminalTypes,
    required this.statusGroups,
  });

  factory TerminalFilterOptions.fromJson(Map<String, dynamic> json) {
    return TerminalFilterOptions(
      channels: (json['channels'] as List?)
              ?.map((e) => ChannelOption.fromJson(e))
              .toList() ??
          [],
      terminalTypes: (json['terminal_types'] as List?)
              ?.map((e) => TerminalTypeOption.fromJson(e))
              .toList() ??
          [],
      statusGroups: (json['status_groups'] as List?)
              ?.map((e) => StatusGroupCount.fromJson(e))
              .toList() ??
          [],
    );
  }
}

/// 通道选项
class ChannelOption {
  final int channelId;
  final String channelCode;

  ChannelOption({
    required this.channelId,
    required this.channelCode,
  });

  factory ChannelOption.fromJson(Map<String, dynamic> json) {
    return ChannelOption(
      channelId: json['channel_id'] ?? 0,
      channelCode: json['channel_code'] ?? '',
    );
  }
}

/// 终端类型选项
class TerminalTypeOption {
  final int channelId;
  final String channelCode;
  final String brandCode;
  final String modelCode;
  final int count;

  TerminalTypeOption({
    required this.channelId,
    required this.channelCode,
    required this.brandCode,
    required this.modelCode,
    required this.count,
  });

  factory TerminalTypeOption.fromJson(Map<String, dynamic> json) {
    return TerminalTypeOption(
      channelId: json['channel_id'] ?? 0,
      channelCode: json['channel_code'] ?? '',
      brandCode: json['brand_code'] ?? '',
      modelCode: json['model_code'] ?? '',
      count: json['count'] ?? 0,
    );
  }

  /// 获取显示名称
  String get displayName => '$brandCode $modelCode';
}

/// 状态分组统计
class StatusGroupCount {
  final String key;
  final String label;
  final int count;

  StatusGroupCount({
    required this.key,
    required this.label,
    required this.count,
  });

  factory StatusGroupCount.fromJson(Map<String, dynamic> json) {
    return StatusGroupCount(
      key: json['key'] ?? '',
      label: json['label'] ?? '',
      count: json['count'] ?? 0,
    );
  }
}

/// 终端流动记录
class TerminalFlowLog {
  final int id;
  final String logType;        // distribute/recall/bind/unbind/activate
  final String logTypeName;    // 下发/回拨/绑定/解绑/激活
  final int? fromAgentId;
  final String fromAgentName;
  final int? toAgentId;
  final String toAgentName;
  final String merchantNo;
  final int status;
  final String statusName;
  final String remark;
  final DateTime createdAt;
  final DateTime? confirmedAt;

  TerminalFlowLog({
    required this.id,
    required this.logType,
    required this.logTypeName,
    this.fromAgentId,
    this.fromAgentName = '',
    this.toAgentId,
    this.toAgentName = '',
    this.merchantNo = '',
    required this.status,
    required this.statusName,
    this.remark = '',
    required this.createdAt,
    this.confirmedAt,
  });

  factory TerminalFlowLog.fromJson(Map<String, dynamic> json) {
    return TerminalFlowLog(
      id: json['id'] ?? 0,
      logType: json['log_type'] ?? '',
      logTypeName: json['log_type_name'] ?? '',
      fromAgentId: json['from_agent_id'],
      fromAgentName: json['from_agent_name'] ?? '',
      toAgentId: json['to_agent_id'],
      toAgentName: json['to_agent_name'] ?? '',
      merchantNo: json['merchant_no'] ?? '',
      status: json['status'] ?? 0,
      statusName: json['status_name'] ?? '',
      remark: json['remark'] ?? '',
      createdAt: DateTime.tryParse(json['created_at'] ?? '') ?? DateTime.now(),
      confirmedAt: json['confirmed_at'] != null
          ? DateTime.tryParse(json['confirmed_at'])
          : null,
    );
  }
}

/// 终端流动记录响应
class TerminalFlowLogsResponse {
  final TerminalInfo terminal;
  final List<TerminalFlowLog> list;
  final int total;
  final int page;
  final int pageSize;

  TerminalFlowLogsResponse({
    required this.terminal,
    required this.list,
    required this.total,
    required this.page,
    required this.pageSize,
  });

  factory TerminalFlowLogsResponse.fromJson(Map<String, dynamic> json) {
    return TerminalFlowLogsResponse(
      terminal: TerminalInfo.fromJson(json['terminal'] ?? {}),
      list: (json['list'] as List?)
              ?.map((e) => TerminalFlowLog.fromJson(e))
              .toList() ??
          [],
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      pageSize: json['page_size'] ?? 20,
    );
  }

  bool get hasMore => list.length >= pageSize;
}

/// 终端基本信息（用于流动记录响应）
class TerminalInfo {
  final String terminalSn;
  final int channelId;
  final String channelCode;
  final String brandCode;
  final String modelCode;

  TerminalInfo({
    required this.terminalSn,
    required this.channelId,
    required this.channelCode,
    required this.brandCode,
    required this.modelCode,
  });

  factory TerminalInfo.fromJson(Map<String, dynamic> json) {
    return TerminalInfo(
      terminalSn: json['terminal_sn'] ?? '',
      channelId: json['channel_id'] ?? 0,
      channelCode: json['channel_code'] ?? '',
      brandCode: json['brand_code'] ?? '',
      modelCode: json['model_code'] ?? '',
    );
  }
}
