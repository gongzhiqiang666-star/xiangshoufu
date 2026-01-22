/**
 * 货款代扣 API
 */
import { get, post } from './request'
import type { PaginatedResponse } from '@/types'
import type {
  GoodsDeduction,
  GoodsDeductionDetail,
  GoodsDeductionDetailRecord,
  GoodsDeductionQueryParams,
  GoodsDeductionSummary,
  CreateGoodsDeductionRequest,
  RejectGoodsDeductionRequest,
} from '@/types/deduction'

// 获取货款代扣列表（通用）
export function getGoodsDeductions(
  params: GoodsDeductionQueryParams
): Promise<PaginatedResponse<GoodsDeduction>> {
  return get<PaginatedResponse<GoodsDeduction>>('/v1/goods-deduction', params)
}

// 获取我发起的货款代扣
export function getSentGoodsDeductions(
  params: GoodsDeductionQueryParams
): Promise<PaginatedResponse<GoodsDeduction>> {
  return get<PaginatedResponse<GoodsDeduction>>('/v1/goods-deduction/sent', params)
}

// 获取我接收的货款代扣
export function getReceivedGoodsDeductions(
  params: GoodsDeductionQueryParams
): Promise<PaginatedResponse<GoodsDeduction>> {
  return get<PaginatedResponse<GoodsDeduction>>('/v1/goods-deduction/received', params)
}

// 获取货款代扣详情
export function getGoodsDeductionDetail(id: number): Promise<GoodsDeductionDetail> {
  return get<GoodsDeductionDetail>(`/v1/goods-deduction/${id}`)
}

// 创建货款代扣
export function createGoodsDeduction(
  data: CreateGoodsDeductionRequest
): Promise<GoodsDeduction> {
  return post<GoodsDeduction>('/v1/goods-deduction', data)
}

// 接收货款代扣
export function acceptGoodsDeduction(id: number): Promise<void> {
  return post<void>(`/v1/goods-deduction/${id}/accept`)
}

// 拒绝货款代扣
export function rejectGoodsDeduction(
  id: number,
  data: RejectGoodsDeductionRequest
): Promise<void> {
  return post<void>(`/v1/goods-deduction/${id}/reject`, data)
}

// 获取货款代扣扣款明细
export function getGoodsDeductionDetails(
  id: number,
  params?: { page?: number; page_size?: number }
): Promise<PaginatedResponse<GoodsDeductionDetailRecord>> {
  return get<PaginatedResponse<GoodsDeductionDetailRecord>>(
    `/v1/goods-deduction/${id}/details`,
    params
  )
}

// 获取货款代扣统计汇总
export function getGoodsDeductionSummary(type?: 'sent' | 'received'): Promise<GoodsDeductionSummary> {
  const params = type ? { type } : {}
  return get<GoodsDeductionSummary>('/v1/goods-deduction/summary', params)
}

// 导出货款代扣
export function exportGoodsDeductions(
  params: GoodsDeductionQueryParams & { type?: 'sent' | 'received' }
): Promise<{ task_id: string }> {
  return post<{ task_id: string }>('/v1/goods-deduction/export', params)
}
