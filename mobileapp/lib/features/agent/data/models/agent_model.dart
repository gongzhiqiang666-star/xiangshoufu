import 'package:json_annotation/json_annotation.dart';

part 'agent_model.g.dart';

/// 代理商简要信息
@JsonSerializable()
class AgentInfo {
  final int id;
  @JsonKey(name: 'agent_no')
  final String agentNo;
  @JsonKey(name: 'agent_name')
  final String agentName;
  @JsonKey(name: 'contact_phone')
  final String contactPhone;
  final int level;
  final int status;
  @JsonKey(name: 'status_name')
  final String? statusName;
  @JsonKey(name: 'direct_agent_count')
  final int directAgentCount;
  @JsonKey(name: 'direct_merchant_count')
  final int directMerchantCount;
  @JsonKey(name: 'register_time')
  final String? registerTime;

  AgentInfo({
    required this.id,
    required this.agentNo,
    required this.agentName,
    required this.contactPhone,
    required this.level,
    required this.status,
    this.statusName,
    this.directAgentCount = 0,
    this.directMerchantCount = 0,
    this.registerTime,
  });

  factory AgentInfo.fromJson(Map<String, dynamic> json) => _$AgentInfoFromJson(json);
  Map<String, dynamic> toJson() => _$AgentInfoToJson(this);
}

/// 代理商详情
@JsonSerializable()
class AgentDetail {
  final int id;
  @JsonKey(name: 'agent_no')
  final String agentNo;
  @JsonKey(name: 'agent_name')
  final String agentName;
  @JsonKey(name: 'contact_name')
  final String? contactName;
  @JsonKey(name: 'contact_phone')
  final String contactPhone;
  @JsonKey(name: 'id_card_no')
  final String? idCardNo;
  @JsonKey(name: 'parent_id')
  final int? parentId;
  @JsonKey(name: 'parent_name')
  final String? parentName;
  final int level;
  final int status;
  @JsonKey(name: 'status_name')
  final String? statusName;
  @JsonKey(name: 'invite_code')
  final String? inviteCode;
  @JsonKey(name: 'qr_code_url')
  final String? qrCodeUrl;
  @JsonKey(name: 'bank_name')
  final String? bankName;
  @JsonKey(name: 'bank_account')
  final String? bankAccount;
  @JsonKey(name: 'bank_card_no')
  final String? bankCardNo;
  @JsonKey(name: 'direct_agent_count')
  final int directAgentCount;
  @JsonKey(name: 'team_agent_count')
  final int teamAgentCount;
  @JsonKey(name: 'direct_merchant_count')
  final int directMerchantCount;
  @JsonKey(name: 'team_merchant_count')
  final int teamMerchantCount;
  @JsonKey(name: 'register_time')
  final String? registerTime;

  AgentDetail({
    required this.id,
    required this.agentNo,
    required this.agentName,
    this.contactName,
    required this.contactPhone,
    this.idCardNo,
    this.parentId,
    this.parentName,
    required this.level,
    required this.status,
    this.statusName,
    this.inviteCode,
    this.qrCodeUrl,
    this.bankName,
    this.bankAccount,
    this.bankCardNo,
    this.directAgentCount = 0,
    this.teamAgentCount = 0,
    this.directMerchantCount = 0,
    this.teamMerchantCount = 0,
    this.registerTime,
  });

  factory AgentDetail.fromJson(Map<String, dynamic> json) => _$AgentDetailFromJson(json);
  Map<String, dynamic> toJson() => _$AgentDetailToJson(this);
}

/// 邀请码信息
@JsonSerializable()
class InviteCodeInfo {
  @JsonKey(name: 'invite_code')
  final String inviteCode;
  @JsonKey(name: 'invite_link')
  final String inviteLink;
  @JsonKey(name: 'qr_code_url')
  final String? qrCodeUrl;

  InviteCodeInfo({
    required this.inviteCode,
    required this.inviteLink,
    this.qrCodeUrl,
  });

  factory InviteCodeInfo.fromJson(Map<String, dynamic> json) => _$InviteCodeInfoFromJson(json);
  Map<String, dynamic> toJson() => _$InviteCodeInfoToJson(this);
}

/// 团队统计
@JsonSerializable()
class TeamStats {
  @JsonKey(name: 'direct_agent_count')
  final int directAgentCount;
  @JsonKey(name: 'team_agent_count')
  final int teamAgentCount;
  @JsonKey(name: 'direct_merchant_count')
  final int directMerchantCount;
  @JsonKey(name: 'team_merchant_count')
  final int teamMerchantCount;
  @JsonKey(name: 'today_new_agents')
  final int todayNewAgents;
  @JsonKey(name: 'month_new_agents')
  final int monthNewAgents;

  TeamStats({
    this.directAgentCount = 0,
    this.teamAgentCount = 0,
    this.directMerchantCount = 0,
    this.teamMerchantCount = 0,
    this.todayNewAgents = 0,
    this.monthNewAgents = 0,
  });

  factory TeamStats.fromJson(Map<String, dynamic> json) => _$TeamStatsFromJson(json);
  Map<String, dynamic> toJson() => _$TeamStatsToJson(this);
}

/// 创建代理商请求
@JsonSerializable()
class CreateAgentRequest {
  @JsonKey(name: 'agent_name')
  final String agentName;
  @JsonKey(name: 'contact_name')
  final String contactName;
  @JsonKey(name: 'contact_phone')
  final String contactPhone;
  @JsonKey(name: 'id_card_no')
  final String? idCardNo;
  @JsonKey(name: 'bank_name')
  final String? bankName;
  @JsonKey(name: 'bank_account')
  final String? bankAccount;
  @JsonKey(name: 'bank_card_no')
  final String? bankCardNo;
  @JsonKey(name: 'parent_id')
  final int? parentId;

  CreateAgentRequest({
    required this.agentName,
    required this.contactName,
    required this.contactPhone,
    this.idCardNo,
    this.bankName,
    this.bankAccount,
    this.bankCardNo,
    this.parentId,
  });

  factory CreateAgentRequest.fromJson(Map<String, dynamic> json) => _$CreateAgentRequestFromJson(json);
  Map<String, dynamic> toJson() => _$CreateAgentRequestToJson(this);
}

/// 下级代理商列表响应
@JsonSerializable()
class SubordinateListResponse {
  final List<AgentInfo> list;
  final int total;

  SubordinateListResponse({
    required this.list,
    required this.total,
  });

  factory SubordinateListResponse.fromJson(Map<String, dynamic> json) => _$SubordinateListResponseFromJson(json);
  Map<String, dynamic> toJson() => _$SubordinateListResponseToJson(this);
}
