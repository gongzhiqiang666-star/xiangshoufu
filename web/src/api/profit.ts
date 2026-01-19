import { get, post, put } from './request'
import type {
  Profit,
  ProfitStats,
  ProfitDailySummary,
  ProfitQueryParams,
  AdjustmentRequest,
  Adjustment,
  PaginatedResponse,
} from '@/types'

/**
 * 获取分润列表
 */
export function getProfits(params: ProfitQueryParams): Promise<PaginatedResponse<Profit>> {
  return get<PaginatedResponse<Profit>>('/v1/profits', params)
}

/**
 * 获取分润详情
 */
export function getProfitDetail(id: number): Promise<Profit> {
  return get<Profit>(`/v1/profits/${id}`)
}

/**
 * 获取分润详情 (别名)
 */
export function getProfit(id: number): Promise<Profit> {
  return get<Profit>(`/v1/profits/${id}`)
}

/**
 * 获取分润统计
 */
export function getProfitStats(params?: {
  channel_id?: number
  start_date?: string
  end_date?: string
}): Promise<ProfitStats> {
  return get<ProfitStats>('/v1/profits/stats', params)
}

/**
 * 获取每日分润汇总
 */
export function getProfitDailySummary(params?: {
  channel_id?: number
  days?: number
}): Promise<ProfitDailySummary[]> {
  return get<ProfitDailySummary[]>('/v1/profits/daily', params)
}

/**
 * 导出分润
 */
export function exportProfits(params: ProfitQueryParams): Promise<{ task_id: string }> {
  return post('/v1/profits/export', params)
}

// ======== 手动调账 ========

/**
 * 获取调账列表
 */
export function getAdjustments(params: {
  status?: string
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Adjustment>> {
  return get<PaginatedResponse<Adjustment>>('/v1/profits/adjustments', params)
}

/**
 * 创建调账申请
 */
export function createAdjustment(data: AdjustmentRequest): Promise<void> {
  return post<void>('/v1/profits/adjustments', data)
}

/**
 * 审核调账
 */
export function auditAdjustment(
  id: number,
  data: { status: 'approved' | 'rejected'; remark?: string }
): Promise<void> {
  return put<void>(`/v1/profits/adjustments/${id}/audit`, data)
}
