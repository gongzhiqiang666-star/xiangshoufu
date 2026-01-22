import { get, post, put } from './request'
import type {
  Merchant,
  MerchantDetail,
  MerchantStats,
  MerchantQueryParams,
  PaginatedResponse,
} from '@/types'

// 扩展统计类型
export interface ExtendedMerchantStats {
  total_count: number
  active_count: number
  pending_count: number
  disabled_count: number
  direct_count: number
  team_count: number
  today_new_count: number
  loyal_count: number
  quality_count: number
  potential_count: number
  normal_count: number
  low_active_count: number
  inactive_count: number
}

/**
 * 获取商户列表
 */
export function getMerchants(params: MerchantQueryParams): Promise<PaginatedResponse<Merchant>> {
  return get<PaginatedResponse<Merchant>>('/v1/merchants', params)
}

/**
 * 获取商户详情
 */
export function getMerchantDetail(id: number): Promise<MerchantDetail> {
  return get<MerchantDetail>(`/v1/merchants/${id}`)
}

/**
 * 获取商户详情 (别名)
 */
export function getMerchant(id: number): Promise<MerchantDetail> {
  return get<MerchantDetail>(`/v1/merchants/${id}`)
}

/**
 * 获取商户关联终端
 */
export function getMerchantTerminals(merchantId: number): Promise<any[]> {
  return get<any[]>(`/v1/merchants/${merchantId}/terminals`)
}

/**
 * 获取商户统计
 */
export function getMerchantStats(): Promise<MerchantStats> {
  return get<MerchantStats>('/v1/merchants/stats')
}

/**
 * 获取商户交易列表
 */
export function getMerchantTransactions(
  merchantId: number,
  params: { page?: number; page_size?: number }
): Promise<PaginatedResponse<any>> {
  return get<PaginatedResponse<any>>(`/v1/merchants/${merchantId}/transactions`, params)
}

/**
 * 商户登记（记录完整手机号）
 */
export function registerMerchant(
  merchantId: number,
  data: { phone: string; remark?: string }
): Promise<void> {
  return post<void>(`/v1/merchants/${merchantId}/register`, data)
}

/**
 * 更新商户登记信息
 */
export function updateMerchantRegister(
  merchantId: number,
  data: { phone: string; remark?: string }
): Promise<void> {
  return put<void>(`/v1/merchants/${merchantId}/register`, data)
}

/**
 * 修改商户费率
 */
export function updateMerchantRate(
  merchantId: number,
  data: { credit_rate: number; debit_rate?: number }
): Promise<void> {
  return put<void>(`/v1/merchants/${merchantId}/rate`, data)
}

/**
 * 搜索商户
 */
export function searchMerchants(keyword: string): Promise<Merchant[]> {
  return get<Merchant[]>('/v1/merchants/search', { keyword, limit: 20 })
}

/**
 * 获取扩展统计（包含直营/团队、商户类型分布）
 */
export function getExtendedMerchantStats(): Promise<ExtendedMerchantStats> {
  return get<ExtendedMerchantStats>('/v1/merchants/stats/extended')
}

/**
 * 导出商户列表
 */
export function exportMerchants(params: MerchantQueryParams): Promise<Blob> {
  // 使用原生fetch获取blob
  const queryParams = new URLSearchParams()
  if (params.keyword) queryParams.append('keyword', params.keyword)
  if (params.merchant_type) queryParams.append('merchant_type', params.merchant_type)
  if (params.is_direct !== undefined) queryParams.append('is_direct', String(params.is_direct))
  if (params.status !== undefined) queryParams.append('status', String(params.status))

  return fetch(`/api/v1/merchants/export?${queryParams.toString()}`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${localStorage.getItem('token') || ''}`,
    },
  }).then(res => res.blob())
}

/**
 * 更新商户状态
 */
export function updateMerchantStatus(
  merchantId: number,
  status: number
): Promise<void> {
  return put<void>(`/v1/merchants/${merchantId}/status`, { status })
}
