import { get, post, put, del } from './request'
import type {
  TaxChannel,
  CreateTaxChannelParams,
  UpdateTaxChannelParams,
  ChannelTaxMapping,
  SetChannelTaxMappingParams,
  CalculateTaxParams,
  TaxCalculationResult,
  TaxChannelStatus,
} from '@/types'

// ========== 税筹通道管理 ==========

/**
 * 创建税筹通道
 */
export function createTaxChannel(data: CreateTaxChannelParams): Promise<{ id: number }> {
  return post<{ id: number }>('/v1/tax-channels', data)
}

/**
 * 更新税筹通道
 */
export function updateTaxChannel(id: number, data: UpdateTaxChannelParams): Promise<void> {
  return put<void>(`/v1/tax-channels/${id}`, data)
}

/**
 * 获取税筹通道详情
 */
export function getTaxChannel(id: number): Promise<TaxChannel> {
  return get<TaxChannel>(`/v1/tax-channels/${id}`)
}

/**
 * 获取税筹通道列表
 */
export function getTaxChannelList(status?: TaxChannelStatus): Promise<{ list: TaxChannel[]; total: number }> {
  const params: Record<string, unknown> = {}
  if (status !== undefined) {
    params.status = status
  }
  return get<{ list: TaxChannel[]; total: number }>('/v1/tax-channels', params)
}

/**
 * 删除税筹通道
 */
export function deleteTaxChannel(id: number): Promise<void> {
  return del<void>(`/v1/tax-channels/${id}`)
}

// ========== 通道-税筹通道映射 ==========

/**
 * 设置通道税筹映射
 */
export function setChannelTaxMapping(data: SetChannelTaxMappingParams): Promise<void> {
  return post<void>('/v1/tax-channels/mappings', data)
}

/**
 * 获取通道税筹映射
 */
export function getChannelTaxMappings(channelId: number): Promise<{ list: ChannelTaxMapping[] }> {
  return get<{ list: ChannelTaxMapping[] }>(`/v1/tax-channels/mappings/channel/${channelId}`)
}

/**
 * 删除通道税筹映射
 */
export function deleteChannelTaxMapping(id: number): Promise<void> {
  return del<void>(`/v1/tax-channels/mappings/${id}`)
}

// ========== 税费计算 ==========

/**
 * 计算提现税费
 */
export function calculateWithdrawalTax(data: CalculateTaxParams): Promise<TaxCalculationResult> {
  return post<TaxCalculationResult>('/v1/tax-channels/calculate', data)
}
