import { get, post } from './request'
import type { AgentChannel, AgentChannelStats } from '@/types'

/**
 * 获取代理商通道列表
 */
export function getAgentChannels(agentId?: number): Promise<AgentChannel[]> {
  return get<AgentChannel[]>('/v1/agent-channels', agentId ? { agent_id: agentId } : {})
}

/**
 * 获取已启用的通道列表（用于APP端）
 */
export function getEnabledChannels(): Promise<AgentChannel[]> {
  return get<AgentChannel[]>('/v1/agent-channels/enabled')
}

/**
 * 获取代理商通道统计
 */
export function getAgentChannelStats(agentId?: number): Promise<AgentChannelStats> {
  return get<AgentChannelStats>('/v1/agent-channels/stats', agentId ? { agent_id: agentId } : {})
}

/**
 * 启用代理商通道
 */
export function enableChannel(agentId: number, channelId: number): Promise<void> {
  return post<void>('/v1/agent-channels/enable', {
    agent_id: agentId,
    channel_id: channelId,
  })
}

/**
 * 禁用代理商通道
 */
export function disableChannel(agentId: number, channelId: number): Promise<void> {
  return post<void>('/v1/agent-channels/disable', {
    agent_id: agentId,
    channel_id: channelId,
  })
}

/**
 * 设置通道可见性
 */
export function setChannelVisibility(
  agentId: number,
  channelId: number,
  isVisible: boolean
): Promise<void> {
  return post<void>('/v1/agent-channels/visibility', {
    agent_id: agentId,
    channel_id: channelId,
    is_visible: isVisible,
  })
}

/**
 * 批量启用通道
 */
export function batchEnableChannels(agentId: number, channelIds: number[]): Promise<void> {
  return post<void>('/v1/agent-channels/batch-enable', {
    agent_id: agentId,
    channel_ids: channelIds,
  })
}

/**
 * 批量禁用通道
 */
export function batchDisableChannels(agentId: number, channelIds: number[]): Promise<void> {
  return post<void>('/v1/agent-channels/batch-disable', {
    agent_id: agentId,
    channel_ids: channelIds,
  })
}

/**
 * 初始化代理商通道配置
 */
export function initAgentChannels(agentId: number): Promise<void> {
  return post<void>('/v1/agent-channels/init', {
    agent_id: agentId,
  })
}
