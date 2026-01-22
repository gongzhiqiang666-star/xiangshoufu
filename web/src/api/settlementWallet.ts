import { get, post, put } from './request'
import type {
  EnableSettlementWalletParams,
  SettlementWalletSummary,
  SubordinateBalance,
  SettlementWalletUsage,
  UseSettlementParams,
  ReturnSettlementParams,
  SettlementUsageQueryParams,
  PaginatedResponse,
} from '@/types'

// ========== 沉淀钱包配置 ==========

/**
 * 开通沉淀钱包(管理员)
 */
export function enableSettlementWallet(data: EnableSettlementWalletParams): Promise<void> {
  return post<void>('/v1/settlement-wallet/enable', data)
}

/**
 * 关闭沉淀钱包(管理员)
 */
export function disableSettlementWallet(agentId: number): Promise<void> {
  return post<void>(`/v1/settlement-wallet/disable/${agentId}`)
}

/**
 * 更新沉淀比例
 */
export function updateSettlementRatio(agentId: number, ratio: number): Promise<void> {
  return put<void>(`/v1/settlement-wallet/ratio/${agentId}`, { ratio })
}

// ========== 沉淀钱包汇总 ==========

/**
 * 获取沉淀钱包汇总
 */
export function getSettlementWalletSummary(): Promise<SettlementWalletSummary> {
  return get<SettlementWalletSummary>('/v1/settlement-wallet/summary')
}

/**
 * 获取下级余额明细
 */
export function getSubordinateBalances(): Promise<{ list: SubordinateBalance[] }> {
  return get<{ list: SubordinateBalance[] }>('/v1/settlement-wallet/subordinates')
}

// ========== 使用和归还 ==========

/**
 * 使用沉淀款
 */
export function useSettlement(data: UseSettlementParams): Promise<{ usage_no: string }> {
  return post<{ usage_no: string }>('/v1/settlement-wallet/use', data)
}

/**
 * 归还沉淀款
 */
export function returnSettlement(data: ReturnSettlementParams): Promise<{ usage_no: string }> {
  return post<{ usage_no: string }>('/v1/settlement-wallet/return', data)
}

/**
 * 获取使用记录列表
 */
export function getSettlementUsageList(
  params: SettlementUsageQueryParams
): Promise<PaginatedResponse<SettlementWalletUsage>> {
  return get<PaginatedResponse<SettlementWalletUsage>>('/v1/settlement-wallet/usages', params)
}
