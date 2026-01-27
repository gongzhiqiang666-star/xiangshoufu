import 'package:json_annotation/json_annotation.dart';

part 'merchant_model.g.dart';

/// 商户模型
@JsonSerializable()
class Merchant {
  final int id;
  @JsonKey(name: 'merchant_no')
  final String merchantNo;
  @JsonKey(name: 'merchant_name')
  final String merchantName;
  @JsonKey(name: 'agent_id')
  final int agentId;
  @JsonKey(name: 'terminal_sn')
  final String? terminalSn;
  final int status;
  @JsonKey(name: 'status_name')
  final String? statusName;
  @JsonKey(name: 'merchant_type')
  final String merchantType;
  @JsonKey(name: 'is_direct')
  final bool isDirect;
  @JsonKey(name: 'owner_type')
  final String? ownerType;
  @JsonKey(name: 'credit_rate')
  final String? creditRate;
  @JsonKey(name: 'debit_rate')
  final String? debitRate;
  @JsonKey(name: 'activated_at')
  final String? activatedAt;
  @JsonKey(name: 'registered_phone')
  final String? registeredPhone;
  @JsonKey(name: 'created_at')
  final String? createdAt;

  Merchant({
    required this.id,
    required this.merchantNo,
    required this.merchantName,
    required this.agentId,
    this.terminalSn,
    required this.status,
    this.statusName,
    required this.merchantType,
    required this.isDirect,
    this.ownerType,
    this.creditRate,
    this.debitRate,
    this.activatedAt,
    this.registeredPhone,
    this.createdAt,
  });

  factory Merchant.fromJson(Map<String, dynamic> json) =>
      _$MerchantFromJson(json);

  Map<String, dynamic> toJson() => _$MerchantToJson(this);

  /// 获取商户类型显示名称
  /// 5档分类：优质/中等/普通/预警/流失
  String get merchantTypeName {
    switch (merchantType) {
      case 'quality':
        return '优质商户';
      case 'medium':
        return '中等商户';
      case 'normal':
        return '普通商户';
      case 'warning':
        return '预警商户';
      case 'churned':
        return '流失商户';
      default:
        return merchantType;
    }
  }

  /// 获取归属类型显示名称
  String get ownerTypeName => isDirect ? '直营' : '团队';
}

/// 商户统计
@JsonSerializable()
class MerchantStats {
  @JsonKey(name: 'total_count')
  final int totalCount;
  @JsonKey(name: 'active_count')
  final int activeCount;
  @JsonKey(name: 'pending_count')
  final int pendingCount;
  @JsonKey(name: 'disabled_count')
  final int disabledCount;
  @JsonKey(name: 'direct_count')
  final int directCount;
  @JsonKey(name: 'team_count')
  final int teamCount;
  @JsonKey(name: 'today_new_count')
  final int todayNewCount;

  MerchantStats({
    required this.totalCount,
    required this.activeCount,
    required this.pendingCount,
    required this.disabledCount,
    required this.directCount,
    required this.teamCount,
    required this.todayNewCount,
  });

  factory MerchantStats.fromJson(Map<String, dynamic> json) =>
      _$MerchantStatsFromJson(json);

  Map<String, dynamic> toJson() => _$MerchantStatsToJson(this);
}

/// 商户详情
@JsonSerializable()
class MerchantDetail {
  final int id;
  @JsonKey(name: 'merchant_no')
  final String merchantNo;
  @JsonKey(name: 'merchant_name')
  final String merchantName;
  @JsonKey(name: 'agent_id')
  final int agentId;
  @JsonKey(name: 'agent_name')
  final String? agentName;
  @JsonKey(name: 'agent_level')
  final int? agentLevel;
  @JsonKey(name: 'channel_id')
  final int? channelId;
  @JsonKey(name: 'channel_name')
  final String? channelName;
  @JsonKey(name: 'terminal_sn')
  final String? terminalSn;
  final int status;
  @JsonKey(name: 'status_name')
  final String? statusName;
  @JsonKey(name: 'approve_status')
  final int? approveStatus;
  @JsonKey(name: 'legal_name')
  final String? legalName;
  @JsonKey(name: 'legal_id_card')
  final String? legalIdCard;
  final String? mcc;
  @JsonKey(name: 'credit_rate')
  final String? creditRate;
  @JsonKey(name: 'debit_rate')
  final String? debitRate;
  @JsonKey(name: 'merchant_type')
  final String merchantType;
  @JsonKey(name: 'is_direct')
  final bool isDirect;
  @JsonKey(name: 'activated_at')
  final String? activatedAt;
  @JsonKey(name: 'registered_phone')
  final String? registeredPhone;
  @JsonKey(name: 'register_remark')
  final String? registerRemark;
  @JsonKey(name: 'month_amount')
  final int? monthAmount;
  @JsonKey(name: 'month_count')
  final int? monthCount;
  @JsonKey(name: 'terminal_count')
  final int? terminalCount;
  @JsonKey(name: 'created_at')
  final String? createdAt;
  @JsonKey(name: 'updated_at')
  final String? updatedAt;

  MerchantDetail({
    required this.id,
    required this.merchantNo,
    required this.merchantName,
    required this.agentId,
    this.agentName,
    this.agentLevel,
    this.channelId,
    this.channelName,
    this.terminalSn,
    required this.status,
    this.statusName,
    this.approveStatus,
    this.legalName,
    this.legalIdCard,
    this.mcc,
    this.creditRate,
    this.debitRate,
    required this.merchantType,
    required this.isDirect,
    this.activatedAt,
    this.registeredPhone,
    this.registerRemark,
    this.monthAmount,
    this.monthCount,
    this.terminalCount,
    this.createdAt,
    this.updatedAt,
  });

  factory MerchantDetail.fromJson(Map<String, dynamic> json) =>
      _$MerchantDetailFromJson(json);

  Map<String, dynamic> toJson() => _$MerchantDetailToJson(this);

  /// 格式化月交易额
  String get monthAmountFormatted {
    if (monthAmount == null) return '0.00';
    return (monthAmount! / 100).toStringAsFixed(2);
  }
}
