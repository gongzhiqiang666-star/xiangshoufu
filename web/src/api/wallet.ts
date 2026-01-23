import { get, post, put } from './request'
import type {
  Wallet,
  WalletLog,
  WalletSummary,
  WithdrawRequest,
  Withdrawal,
  WalletLogQueryParams,
  PaginatedResponse,
} from '@/types'

/**
 * 获取钱包列表
 */
export async function getWallets(params?: { channel_id?: number }): Promise<Wallet[]> {
  const res = await get<{ list: Wallet[] }>('/v1/wallets', params)
  return res.list || []
}

/**
 * 获取单个钱包详情
 */
export function getWallet(id: number): Promise<Wallet> {
  return get<Wallet>(`/v1/wallets/${id}`)
}

/**
 * 获取钱包汇总
 */
export function getWalletSummary(): Promise<WalletSummary> {
  return get<WalletSummary>('/v1/wallets/summary')
}

/**
 * 获取钱包流水
 */
export function getWalletLogs(
  walletId: number,
  params: WalletLogQueryParams
): Promise<PaginatedResponse<WalletLog>> {
  return get<PaginatedResponse<WalletLog>>(`/v1/wallets/${walletId}/logs`, params)
}

/**
 * 申请提现
 */
export function applyWithdraw(data: WithdrawRequest): Promise<void> {
  return post<void>(`/v1/wallets/${data.wallet_id}/withdraw`, data)
}

// ======== 提现管理 ========

/**
 * 获取提现列表
 */
export function getWithdrawals(params: {
  status?: string
  channel_id?: number
  agent_id?: number
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Withdrawal>> {
  return get<PaginatedResponse<Withdrawal>>('/v1/withdrawals', params)
}

/**
 * 获取提现详情
 */
export function getWithdrawalDetail(id: number): Promise<Withdrawal> {
  return get<Withdrawal>(`/v1/withdrawals/${id}`)
}

/**
 * 审核提现
 */
export function auditWithdrawal(
  id: number,
  data: { status: 'approved' | 'rejected'; remark?: string }
): Promise<void> {
  return put<void>(`/v1/withdrawals/${id}/audit`, data)
}

/**
 * 批量审核提现
 */
export function batchAuditWithdrawals(data: {
  ids: number[]
  status: 'approved' | 'rejected'
  remark?: string
}): Promise<void> {
  return put<void>('/v1/withdrawals/batch-audit', data)
}

/**
 * 获取提现门槛配置
 */
export function getWithdrawThresholds(): Promise<
  { channel_id: number; wallet_type: string; threshold: number }[]
> {
  return get('/v1/withdrawals/thresholds')
}

/**
 * 更新提现门槛配置
 */
export function updateWithdrawThreshold(data: {
  channel_id: number
  wallet_type: string
  threshold: number
}): Promise<void> {
  return put<void>('/v1/withdrawals/thresholds', data)
}
