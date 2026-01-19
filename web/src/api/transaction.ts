import { get, post } from './request'
import type {
  Transaction,
  TransactionDetail,
  TransactionStats,
  TransactionTrend,
  TransactionQueryParams,
  PaginatedResponse,
} from '@/types'

/**
 * 获取交易列表
 */
export function getTransactions(
  params: TransactionQueryParams
): Promise<PaginatedResponse<Transaction>> {
  return get<PaginatedResponse<Transaction>>('/v1/transactions', params)
}

/**
 * 获取交易详情
 */
export function getTransactionDetail(id: number): Promise<TransactionDetail> {
  return get<TransactionDetail>(`/v1/transactions/${id}`)
}

/**
 * 获取交易详情 (别名)
 */
export function getTransaction(id: number): Promise<TransactionDetail> {
  return get<TransactionDetail>(`/v1/transactions/${id}`)
}

/**
 * 获取交易统计
 */
export function getTransactionStats(params?: {
  channel_id?: number
  start_date?: string
  end_date?: string
}): Promise<TransactionStats> {
  return get<TransactionStats>('/v1/transactions/stats', params)
}

/**
 * 获取交易趋势
 */
export function getTransactionTrend(params?: {
  channel_id?: number
  days?: number
  owner_type?: 'all' | 'direct' | 'team'
}): Promise<TransactionTrend> {
  return get<TransactionTrend>('/v1/transactions/trend', params)
}

/**
 * 导出交易
 */
export function exportTransactions(params: TransactionQueryParams): Promise<{ task_id: string }> {
  return post('/v1/transactions/export', params)
}

/**
 * 获取导出任务状态
 */
export function getExportStatus(taskId: string): Promise<{
  status: 'processing' | 'completed' | 'failed'
  download_url?: string
}> {
  return get(`/v1/exports/${taskId}`)
}
