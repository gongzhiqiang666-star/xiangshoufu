import 'package:json_annotation/json_annotation.dart';

part 'agent_channel_model.g.dart';

/// 代理商通道配置
@JsonSerializable()
class AgentChannel {
  final int id;
  @JsonKey(name: 'agent_id')
  final int agentId;
  @JsonKey(name: 'channel_id')
  final int channelId;
  @JsonKey(name: 'is_enabled')
  final bool isEnabled;
  @JsonKey(name: 'is_visible')
  final bool isVisible;
  @JsonKey(name: 'enabled_at')
  final String? enabledAt;
  @JsonKey(name: 'disabled_at')
  final String? disabledAt;
  @JsonKey(name: 'enabled_by')
  final int? enabledBy;
  @JsonKey(name: 'disabled_by')
  final int? disabledBy;
  final String? remark;
  @JsonKey(name: 'created_at')
  final String createdAt;
  @JsonKey(name: 'updated_at')
  final String updatedAt;

  // 关联字段
  @JsonKey(name: 'channel_code')
  final String channelCode;
  @JsonKey(name: 'channel_name')
  final String channelName;

  AgentChannel({
    required this.id,
    required this.agentId,
    required this.channelId,
    required this.isEnabled,
    required this.isVisible,
    this.enabledAt,
    this.disabledAt,
    this.enabledBy,
    this.disabledBy,
    this.remark,
    required this.createdAt,
    required this.updatedAt,
    required this.channelCode,
    required this.channelName,
  });

  factory AgentChannel.fromJson(Map<String, dynamic> json) => _$AgentChannelFromJson(json);
  Map<String, dynamic> toJson() => _$AgentChannelToJson(this);
}

/// 代理商通道统计
@JsonSerializable()
class AgentChannelStats {
  @JsonKey(name: 'total_channels')
  final int totalChannels;
  @JsonKey(name: 'enabled_channels')
  final int enabledChannels;
  @JsonKey(name: 'visible_channels')
  final int visibleChannels;

  AgentChannelStats({
    this.totalChannels = 0,
    this.enabledChannels = 0,
    this.visibleChannels = 0,
  });

  factory AgentChannelStats.fromJson(Map<String, dynamic> json) => _$AgentChannelStatsFromJson(json);
  Map<String, dynamic> toJson() => _$AgentChannelStatsToJson(this);
}
