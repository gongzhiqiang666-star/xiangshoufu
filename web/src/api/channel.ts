import { get, post, put, del } from './request'
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

// ============================================================
// 通道配置管理API
// ============================================================

/** 费率配置 */
export interface ChannelRateConfig {
  id: number
  channel_id: number
  rate_code: string
  rate_name: string
  min_rate: string
  max_rate: string
  default_rate: string
  sort_order: number
  status: number
}

/** 押金档位 */
export interface ChannelDepositTier {
  id: number
  channel_id: number
  brand_code: string
  tier_code: string
  deposit_amount: number
  tier_name: string
  max_cashback_amount: number
  default_cashback: number
  sort_order: number
  status: number
}

/** 流量费返现档位 */
export interface ChannelSimCashbackTier {
  id: number
  channel_id: number
  brand_code: string
  tier_order: number
  tier_name: string
  is_last_tier: boolean
  max_cashback_amount: number
  default_cashback: number
  sim_fee_amount: number
  status: number
}

/** 通道完整配置 */
export interface ChannelFullConfig {
  channel_id: number
  channel_code: string
  channel_name: string
  rate_configs: ChannelRateConfig[]
  deposit_tiers: ChannelDepositTier[]
  sim_cashback_tiers: ChannelSimCashbackTier[]
}

/**
 * 获取通道完整配置
 */
export function getChannelFullConfig(channelId: number): Promise<ChannelFullConfig> {
  return get<ChannelFullConfig>(`/v1/admin/channels/${channelId}/full-config`)
}

/**
 * 获取通道费率配置列表
 */
export function getChannelRateConfigs(channelId: number): Promise<ChannelRateConfig[]> {
  return get<ChannelRateConfig[]>(`/v1/admin/channels/${channelId}/rate-configs`)
}

/**
 * 创建通道费率配置
 */
export function createChannelRateConfig(channelId: number, data: {
  rate_code: string
  rate_name: string
  min_rate: string
  max_rate: string
  default_rate?: string
  sort_order?: number
}): Promise<ChannelRateConfig> {
  return post<ChannelRateConfig>(`/v1/admin/channels/${channelId}/rate-configs`, data)
}

/**
 * 更新通道费率配置
 */
export function updateChannelRateConfig(channelId: number, configId: number, data: {
  rate_name?: string
  min_rate?: string
  max_rate?: string
  default_rate?: string
  sort_order?: number
  status?: number
}): Promise<void> {
  return put<void>(`/v1/admin/channels/${channelId}/rate-configs/${configId}`, data)
}

/**
 * 删除通道费率配置
 */
export function deleteChannelRateConfig(channelId: number, configId: number): Promise<void> {
  return del<void>(`/v1/admin/channels/${channelId}/rate-configs/${configId}`)
}

/**
 * 获取通道押金档位列表
 */
export function getChannelDepositTiers(channelId: number): Promise<ChannelDepositTier[]> {
  return get<ChannelDepositTier[]>(`/v1/admin/channels/${channelId}/deposit-tiers`)
}

/**
 * 更新通道押金档位
 */
export function updateChannelDepositTier(channelId: number, tierId: number, data: {
  max_cashback_amount: number
  default_cashback: number
  status?: number
}): Promise<void> {
  return put<void>(`/v1/admin/channels/${channelId}/deposit-tiers/${tierId}`, data)
}

/**
 * 获取通道流量费返现档位列表
 */
export function getChannelSimCashbackTiers(channelId: number): Promise<ChannelSimCashbackTier[]> {
  return get<ChannelSimCashbackTier[]>(`/v1/admin/channels/${channelId}/sim-cashback-tiers`)
}

/**
 * 批量设置通道流量费返现档位
 */
export function batchSetChannelSimCashbackTiers(channelId: number, tiers: {
  tier_order: number
  tier_name: string
  is_last_tier: boolean
  max_cashback_amount: number
  default_cashback: number
  sim_fee_amount: number
}[]): Promise<void> {
  return post<void>(`/v1/admin/channels/${channelId}/sim-cashback-tiers/batch`, { tiers })
}
