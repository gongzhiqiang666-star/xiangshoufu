// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'agent_channel_model.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

AgentChannel _$AgentChannelFromJson(Map<String, dynamic> json) => AgentChannel(
      id: (json['id'] as num).toInt(),
      agentId: (json['agent_id'] as num).toInt(),
      channelId: (json['channel_id'] as num).toInt(),
      isEnabled: json['is_enabled'] as bool,
      isVisible: json['is_visible'] as bool,
      enabledAt: json['enabled_at'] as String?,
      disabledAt: json['disabled_at'] as String?,
      enabledBy: (json['enabled_by'] as num?)?.toInt(),
      disabledBy: (json['disabled_by'] as num?)?.toInt(),
      remark: json['remark'] as String?,
      createdAt: json['created_at'] as String,
      updatedAt: json['updated_at'] as String,
      channelCode: json['channel_code'] as String,
      channelName: json['channel_name'] as String,
    );

Map<String, dynamic> _$AgentChannelToJson(AgentChannel instance) =>
    <String, dynamic>{
      'id': instance.id,
      'agent_id': instance.agentId,
      'channel_id': instance.channelId,
      'is_enabled': instance.isEnabled,
      'is_visible': instance.isVisible,
      'enabled_at': instance.enabledAt,
      'disabled_at': instance.disabledAt,
      'enabled_by': instance.enabledBy,
      'disabled_by': instance.disabledBy,
      'remark': instance.remark,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
      'channel_code': instance.channelCode,
      'channel_name': instance.channelName,
    };

AgentChannelStats _$AgentChannelStatsFromJson(Map<String, dynamic> json) =>
    AgentChannelStats(
      totalChannels: (json['total_channels'] as num?)?.toInt() ?? 0,
      enabledChannels: (json['enabled_channels'] as num?)?.toInt() ?? 0,
      visibleChannels: (json['visible_channels'] as num?)?.toInt() ?? 0,
    );

Map<String, dynamic> _$AgentChannelStatsToJson(AgentChannelStats instance) =>
    <String, dynamic>{
      'total_channels': instance.totalChannels,
      'enabled_channels': instance.enabledChannels,
      'visible_channels': instance.visibleChannels,
    };
