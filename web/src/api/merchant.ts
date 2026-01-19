import { get, post, put } from './request'
import type {
  Merchant,
  MerchantDetail,
  MerchantStats,
  MerchantQueryParams,
  PaginatedResponse,
} from '@/types'

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
