// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'merchant_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

Merchant _$MerchantFromJson(Map<String, dynamic> json) => Merchant(
      id: (json['id'] as num).toInt(),
      merchantNo: json['merchant_no'] as String,
      merchantName: json['merchant_name'] as String,
      agentId: (json['agent_id'] as num).toInt(),
      terminalSn: json['terminal_sn'] as String?,
      status: (json['status'] as num).toInt(),
      statusName: json['status_name'] as String?,
      merchantType: json['merchant_type'] as String,
      isDirect: json['is_direct'] as bool,
      ownerType: json['owner_type'] as String?,
      creditRate: json['credit_rate'] as String?,
      debitRate: json['debit_rate'] as String?,
      activatedAt: json['activated_at'] as String?,
      registeredPhone: json['registered_phone'] as String?,
      createdAt: json['created_at'] as String?,
    );

Map<String, dynamic> _$MerchantToJson(Merchant instance) => <String, dynamic>{
      'id': instance.id,
      'merchant_no': instance.merchantNo,
      'merchant_name': instance.merchantName,
      'agent_id': instance.agentId,
      'terminal_sn': instance.terminalSn,
      'status': instance.status,
      'status_name': instance.statusName,
      'merchant_type': instance.merchantType,
      'is_direct': instance.isDirect,
      'owner_type': instance.ownerType,
      'credit_rate': instance.creditRate,
      'debit_rate': instance.debitRate,
      'activated_at': instance.activatedAt,
      'registered_phone': instance.registeredPhone,
      'created_at': instance.createdAt,
    };

MerchantStats _$MerchantStatsFromJson(Map<String, dynamic> json) =>
    MerchantStats(
      totalCount: (json['total_count'] as num).toInt(),
      activeCount: (json['active_count'] as num).toInt(),
      pendingCount: (json['pending_count'] as num).toInt(),
      disabledCount: (json['disabled_count'] as num).toInt(),
      directCount: (json['direct_count'] as num).toInt(),
      teamCount: (json['team_count'] as num).toInt(),
      todayNewCount: (json['today_new_count'] as num).toInt(),
    );

Map<String, dynamic> _$MerchantStatsToJson(MerchantStats instance) =>
    <String, dynamic>{
      'total_count': instance.totalCount,
      'active_count': instance.activeCount,
      'pending_count': instance.pendingCount,
      'disabled_count': instance.disabledCount,
      'direct_count': instance.directCount,
      'team_count': instance.teamCount,
      'today_new_count': instance.todayNewCount,
    };

MerchantDetail _$MerchantDetailFromJson(Map<String, dynamic> json) =>
    MerchantDetail(
      id: (json['id'] as num).toInt(),
      merchantNo: json['merchant_no'] as String,
      merchantName: json['merchant_name'] as String,
      agentId: (json['agent_id'] as num).toInt(),
      agentName: json['agent_name'] as String?,
      agentLevel: (json['agent_level'] as num?)?.toInt(),
      channelId: (json['channel_id'] as num?)?.toInt(),
      channelName: json['channel_name'] as String?,
      terminalSn: json['terminal_sn'] as String?,
      status: (json['status'] as num).toInt(),
      statusName: json['status_name'] as String?,
      approveStatus: (json['approve_status'] as num?)?.toInt(),
      legalName: json['legal_name'] as String?,
      legalIdCard: json['legal_id_card'] as String?,
      mcc: json['mcc'] as String?,
      creditRate: json['credit_rate'] as String?,
      debitRate: json['debit_rate'] as String?,
      merchantType: json['merchant_type'] as String,
      isDirect: json['is_direct'] as bool,
      activatedAt: json['activated_at'] as String?,
      registeredPhone: json['registered_phone'] as String?,
      registerRemark: json['register_remark'] as String?,
      monthAmount: (json['month_amount'] as num?)?.toInt(),
      monthCount: (json['month_count'] as num?)?.toInt(),
      terminalCount: (json['terminal_count'] as num?)?.toInt(),
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$MerchantDetailToJson(MerchantDetail instance) =>
    <String, dynamic>{
      'id': instance.id,
      'merchant_no': instance.merchantNo,
      'merchant_name': instance.merchantName,
      'agent_id': instance.agentId,
      'agent_name': instance.agentName,
      'agent_level': instance.agentLevel,
      'channel_id': instance.channelId,
      'channel_name': instance.channelName,
      'terminal_sn': instance.terminalSn,
      'status': instance.status,
      'status_name': instance.statusName,
      'approve_status': instance.approveStatus,
      'legal_name': instance.legalName,
      'legal_id_card': instance.legalIdCard,
      'mcc': instance.mcc,
      'credit_rate': instance.creditRate,
      'debit_rate': instance.debitRate,
      'merchant_type': instance.merchantType,
      'is_direct': instance.isDirect,
      'activated_at': instance.activatedAt,
      'registered_phone': instance.registeredPhone,
      'register_remark': instance.registerRemark,
      'month_amount': instance.monthAmount,
      'month_count': instance.monthCount,
      'terminal_count': instance.terminalCount,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };
