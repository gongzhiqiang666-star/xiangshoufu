import { get, post, put } from './request'
import type {
  Terminal,
  TerminalDetail,
  TerminalStats,
  TerminalHistory,
  TerminalDispatchParams,
  TerminalQueryParams,
  PaginatedResponse,
} from '@/types'

/**
 * 获取终端列表
 */
export function getTerminals(params: TerminalQueryParams): Promise<PaginatedResponse<Terminal>> {
  return get<PaginatedResponse<Terminal>>('/v1/terminals', params)
}

/**
 * 获取终端详情
 */
export function getTerminalDetail(sn: string): Promise<TerminalDetail> {
  return get<TerminalDetail>(`/v1/terminals/${sn}`)
}

/**
 * 获取终端详情 (按ID)
 */
export function getTerminal(id: number): Promise<TerminalDetail> {
  return get<TerminalDetail>(`/v1/terminals/id/${id}`)
}

/**
 * 获取终端交易列表
 */
export function getTerminalTransactions(
  terminalId: number,
  params: { page?: number; page_size?: number }
): Promise<PaginatedResponse<any>> {
  return get<PaginatedResponse<any>>(`/v1/terminals/${terminalId}/transactions`, params)
}

/**
 * 获取终端统计
 */
export function getTerminalStats(): Promise<TerminalStats> {
  return get<TerminalStats>('/v1/terminals/stats')
}

/**
 * 获取终端流转历史
 */
export function getTerminalHistory(sn: string): Promise<TerminalHistory[]> {
  return get<TerminalHistory[]>(`/v1/terminals/${sn}/history`)
}

/**
 * 终端入库
 */
export function importTerminals(data: { channel_id: number; sn_list: string[] }): Promise<{
  success: number
  failed: number
  errors: string[]
}> {
  return post('/v1/terminals/import', data)
}

/**
 * 终端下发（支持跨级）
 */
export function dispatchTerminals(data: TerminalDispatchParams): Promise<void> {
  return post<void>('/v1/terminals/dispatch', data)
}

/**
 * 终端回拨（支持跨级）
 */
export function recallTerminals(data: {
  terminal_ids: number[]
  to_agent_id: number
}): Promise<void> {
  return post<void>('/v1/terminals/recall', data)
}

/**
 * 批量设置终端策略
 */
export function batchConfigTerminals(data: {
  terminal_ids: number[]
  credit_rate?: number
  deposit_amount?: number
  sim_first_fee?: number
  sim_renewal_fee?: number
  sim_renewal_days?: number
}): Promise<void> {
  return put<void>('/v1/terminals/batch-config', data)
}

/**
 * 搜索终端
 */
export function searchTerminals(keyword: string): Promise<Terminal[]> {
  return get<Terminal[]>('/v1/terminals/search', { keyword, limit: 20 })
}
