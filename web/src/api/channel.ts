import { get } from './request'
import type { RateTypeDefinition, Channel } from '@/types'

/**
 * 获取通道列表
 */
export function getChannelList(): Promise<Channel[]> {
  return get<Channel[]>('/v1/admin/channels')
}

/**
 * 获取通道详情
 */
export function getChannelDetail(channelId: number): Promise<Channel> {
  return get<Channel>(`/v1/admin/channels/${channelId}`)
}

/**
 * 获取通道费率类型列表
 */
export function getChannelRateTypes(channelId: number): Promise<RateTypeDefinition[]> {
  return get<RateTypeDefinition[]>(`/v1/admin/channels/${channelId}/rate-types`)
}
