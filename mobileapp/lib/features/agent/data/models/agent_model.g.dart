// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'agent_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AgentInfo _$AgentInfoFromJson(Map<String, dynamic> json) => AgentInfo(
      id: (json['id'] as num).toInt(),
      agentNo: json['agent_no'] as String,
      agentName: json['agent_name'] as String,
      contactPhone: json['contact_phone'] as String,
      level: (json['level'] as num).toInt(),
      status: (json['status'] as num).toInt(),
      statusName: json['status_name'] as String?,
      directAgentCount: (json['direct_agent_count'] as num?)?.toInt() ?? 0,
      directMerchantCount:
          (json['direct_merchant_count'] as num?)?.toInt() ?? 0,
      registerTime: json['register_time'] as String?,
    );

Map<String, dynamic> _$AgentInfoToJson(AgentInfo instance) => <String, dynamic>{
      'id': instance.id,
      'agent_no': instance.agentNo,
      'agent_name': instance.agentName,
      'contact_phone': instance.contactPhone,
      'level': instance.level,
      'status': instance.status,
      'status_name': instance.statusName,
      'direct_agent_count': instance.directAgentCount,
      'direct_merchant_count': instance.directMerchantCount,
      'register_time': instance.registerTime,
    };

AgentDetail _$AgentDetailFromJson(Map<String, dynamic> json) => AgentDetail(
      id: (json['id'] as num).toInt(),
      agentNo: json['agent_no'] as String,
      agentName: json['agent_name'] as String,
      contactName: json['contact_name'] as String?,
      contactPhone: json['contact_phone'] as String,
      idCardNo: json['id_card_no'] as String?,
      parentId: (json['parent_id'] as num?)?.toInt(),
      parentName: json['parent_name'] as String?,
      level: (json['level'] as num).toInt(),
      status: (json['status'] as num).toInt(),
      statusName: json['status_name'] as String?,
      inviteCode: json['invite_code'] as String?,
      qrCodeUrl: json['qr_code_url'] as String?,
      bankName: json['bank_name'] as String?,
      bankAccount: json['bank_account'] as String?,
      bankCardNo: json['bank_card_no'] as String?,
      directAgentCount: (json['direct_agent_count'] as num?)?.toInt() ?? 0,
      teamAgentCount: (json['team_agent_count'] as num?)?.toInt() ?? 0,
      directMerchantCount:
          (json['direct_merchant_count'] as num?)?.toInt() ?? 0,
      teamMerchantCount: (json['team_merchant_count'] as num?)?.toInt() ?? 0,
      registerTime: json['register_time'] as String?,
    );

Map<String, dynamic> _$AgentDetailToJson(AgentDetail instance) =>
    <String, dynamic>{
      'id': instance.id,
      'agent_no': instance.agentNo,
      'agent_name': instance.agentName,
      'contact_name': instance.contactName,
      'contact_phone': instance.contactPhone,
      'id_card_no': instance.idCardNo,
      'parent_id': instance.parentId,
      'parent_name': instance.parentName,
      'level': instance.level,
      'status': instance.status,
      'status_name': instance.statusName,
      'invite_code': instance.inviteCode,
      'qr_code_url': instance.qrCodeUrl,
      'bank_name': instance.bankName,
      'bank_account': instance.bankAccount,
      'bank_card_no': instance.bankCardNo,
      'direct_agent_count': instance.directAgentCount,
      'team_agent_count': instance.teamAgentCount,
      'direct_merchant_count': instance.directMerchantCount,
      'team_merchant_count': instance.teamMerchantCount,
      'register_time': instance.registerTime,
    };

InviteCodeInfo _$InviteCodeInfoFromJson(Map<String, dynamic> json) =>
    InviteCodeInfo(
      inviteCode: json['invite_code'] as String,
      inviteLink: json['invite_link'] as String,
      qrCodeUrl: json['qr_code_url'] as String?,
    );

Map<String, dynamic> _$InviteCodeInfoToJson(InviteCodeInfo instance) =>
    <String, dynamic>{
      'invite_code': instance.inviteCode,
      'invite_link': instance.inviteLink,
      'qr_code_url': instance.qrCodeUrl,
    };

TeamStats _$TeamStatsFromJson(Map<String, dynamic> json) => TeamStats(
      directAgentCount: (json['direct_agent_count'] as num?)?.toInt() ?? 0,
      teamAgentCount: (json['team_agent_count'] as num?)?.toInt() ?? 0,
      directMerchantCount:
          (json['direct_merchant_count'] as num?)?.toInt() ?? 0,
      teamMerchantCount: (json['team_merchant_count'] as num?)?.toInt() ?? 0,
      todayNewAgents: (json['today_new_agents'] as num?)?.toInt() ?? 0,
      monthNewAgents: (json['month_new_agents'] as num?)?.toInt() ?? 0,
    );

Map<String, dynamic> _$TeamStatsToJson(TeamStats instance) => <String, dynamic>{
      'direct_agent_count': instance.directAgentCount,
      'team_agent_count': instance.teamAgentCount,
      'direct_merchant_count': instance.directMerchantCount,
      'team_merchant_count': instance.teamMerchantCount,
      'today_new_agents': instance.todayNewAgents,
      'month_new_agents': instance.monthNewAgents,
    };

CreateAgentRequest _$CreateAgentRequestFromJson(Map<String, dynamic> json) =>
    CreateAgentRequest(
      agentName: json['agent_name'] as String,
      contactName: json['contact_name'] as String,
      contactPhone: json['contact_phone'] as String,
      idCardNo: json['id_card_no'] as String?,
      bankName: json['bank_name'] as String?,
      bankAccount: json['bank_account'] as String?,
      bankCardNo: json['bank_card_no'] as String?,
      parentId: (json['parent_id'] as num?)?.toInt(),
    );

Map<String, dynamic> _$CreateAgentRequestToJson(CreateAgentRequest instance) =>
    <String, dynamic>{
      'agent_name': instance.agentName,
      'contact_name': instance.contactName,
      'contact_phone': instance.contactPhone,
      'id_card_no': instance.idCardNo,
      'bank_name': instance.bankName,
      'bank_account': instance.bankAccount,
      'bank_card_no': instance.bankCardNo,
      'parent_id': instance.parentId,
    };

SubordinateListResponse _$SubordinateListResponseFromJson(
        Map<String, dynamic> json) =>
    SubordinateListResponse(
      list: (json['list'] as List<dynamic>)
          .map((e) => AgentInfo.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: (json['total'] as num).toInt(),
    );

Map<String, dynamic> _$SubordinateListResponseToJson(
        SubordinateListResponse instance) =>
    <String, dynamic>{
      'list': instance.list,
      'total': instance.total,
    };
