import { get, post, put } from './request'
import type {
  Terminal,
  TerminalDetail,
  TerminalStats,
  TerminalHistory,
  TerminalDispatchParams,
  TerminalQueryParams,
  TerminalImportParams,
  TerminalImportResult,
  TerminalRecallParams,
  TerminalRecallResult,
  TerminalRecall,
  TerminalPolicy,
  PolicyOptions,
  BatchSetRateParams,
  BatchSetSimFeeParams,
  BatchSetDepositParams,
  BatchPolicyResult,
  TerminalDistribute,
  TerminalDistributeParams,
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
export function importTerminals(data: TerminalImportParams): Promise<TerminalImportResult> {
  return post<TerminalImportResult>('/v1/terminals/import', data)
}

/**
 * 终端回拨（批量）
 */
export function recallTerminals(data: TerminalRecallParams): Promise<TerminalRecallResult> {
  return post<TerminalRecallResult>('/v1/terminals/batch-recall', data)
}

/**
 * 获取回拨列表
 */
export function getRecallList(
  params: { direction: 'from' | 'to'; page?: number; page_size?: number }
): Promise<PaginatedResponse<TerminalRecall>> {
  return get<PaginatedResponse<TerminalRecall>>('/v1/terminals/recall', params)
}

/**
 * 确认回拨
 */
export function confirmRecall(id: number): Promise<void> {
  return post<void>(`/v1/terminals/recall/${id}/confirm`)
}

/**
 * 拒绝回拨
 */
export function rejectRecall(id: number): Promise<void> {
  return post<void>(`/v1/terminals/recall/${id}/reject`)
}

/**
 * 取消回拨
 */
export function cancelRecall(id: number): Promise<void> {
  return post<void>(`/v1/terminals/recall/${id}/cancel`)
}

/**
 * 终端下发（支持跨级）
 */
export function dispatchTerminals(data: TerminalDispatchParams): Promise<void> {
  return post<void>('/v1/terminals/dispatch', data)
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

// ========== 终端政策设置 ==========

/**
 * 获取政策选项
 */
export function getPolicyOptions(): Promise<PolicyOptions> {
  return get<PolicyOptions>('/v1/terminals/policy-options')
}

/**
 * 获取终端政策
 */
export function getTerminalPolicy(sn: string): Promise<TerminalPolicy> {
  return get<TerminalPolicy>(`/v1/terminals/${sn}/policy`)
}

/**
 * 批量设置费率
 */
export function batchSetRate(data: BatchSetRateParams): Promise<BatchPolicyResult> {
  return post<BatchPolicyResult>('/v1/terminals/batch-set-rate', data)
}

/**
 * 批量设置SIM卡费用
 */
export function batchSetSimFee(data: BatchSetSimFeeParams): Promise<BatchPolicyResult> {
  return post<BatchPolicyResult>('/v1/terminals/batch-set-sim', data)
}

/**
 * 批量设置押金
 */
export function batchSetDeposit(data: BatchSetDepositParams): Promise<BatchPolicyResult> {
  return post<BatchPolicyResult>('/v1/terminals/batch-set-deposit', data)
}

// ========== 终端下发/回拨记录 ==========

/**
 * 终端下发
 */
export function distributeTerminal(data: TerminalDistributeParams): Promise<TerminalDistribute> {
  return post<TerminalDistribute>('/v1/terminal/distribute', data)
}

/**
 * 获取下发列表
 */
export function getDistributeList(
  params: { direction: 'from' | 'to'; status?: number[]; page?: number; page_size?: number }
): Promise<PaginatedResponse<TerminalDistribute>> {
  return get<PaginatedResponse<TerminalDistribute>>('/v1/terminal/distribute', params)
}

/**
 * 确认下发
 */
export function confirmDistribute(id: number): Promise<void> {
  return post<void>(`/v1/terminal/distribute/${id}/confirm`)
}

/**
 * 拒绝下发
 */
export function rejectDistribute(id: number): Promise<void> {
  return post<void>(`/v1/terminal/distribute/${id}/reject`)
}

/**
 * 取消下发
 */
export function cancelDistribute(id: number): Promise<void> {
  return post<void>(`/v1/terminal/distribute/${id}/cancel`)
}
