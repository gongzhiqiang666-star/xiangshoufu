import { get, post } from './request'
import type {
  AgentWalletConfig,
  EnableChargingWalletParams,
  ChargingWalletSummary,
  ChargingWalletDeposit,
  CreateChargingDepositParams,
  ChargingWalletReward,
  IssueChargingRewardParams,
  ChargingDepositQueryParams,
  ChargingRewardQueryParams,
  PaginatedResponse,
} from '@/types'

// ========== 钱包配置 ==========

/**
 * 获取当前代理商的钱包配置
 */
export function getMyWalletConfig(): Promise<AgentWalletConfig> {
  return get<AgentWalletConfig>('/v1/charging-wallet/config')
}

/**
 * 获取指定代理商的钱包配置(管理员)
 */
export function getAgentWalletConfig(agentId: number): Promise<AgentWalletConfig> {
  return get<AgentWalletConfig>(`/v1/charging-wallet/config/${agentId}`)
}

/**
 * 开通充值钱包(管理员)
 */
export function enableChargingWallet(data: EnableChargingWalletParams): Promise<void> {
  return post<void>('/v1/charging-wallet/enable', data)
}

/**
 * 关闭充值钱包(管理员)
 */
export function disableChargingWallet(agentId: number): Promise<void> {
  return post<void>(`/v1/charging-wallet/disable/${agentId}`)
}

// ========== 充值钱包汇总 ==========

/**
 * 获取充值钱包汇总
 */
export function getChargingWalletSummary(): Promise<ChargingWalletSummary> {
  return get<ChargingWalletSummary>('/v1/charging-wallet/summary')
}

// ========== 充值操作 ==========

/**
 * 申请充值
 */
export function createChargingDeposit(data: CreateChargingDepositParams): Promise<{ deposit_no: string }> {
  return post<{ deposit_no: string }>('/v1/charging-wallet/deposits', data)
}

/**
 * 获取充值记录列表
 */
export function getChargingDepositList(
  params: ChargingDepositQueryParams
): Promise<PaginatedResponse<ChargingWalletDeposit>> {
  return get<PaginatedResponse<ChargingWalletDeposit>>('/v1/charging-wallet/deposits', params)
}

/**
 * 获取待审核充值记录(管理员)
 */
export function getPendingChargingDeposits(
  params: { page?: number; page_size?: number }
): Promise<PaginatedResponse<ChargingWalletDeposit>> {
  return get<PaginatedResponse<ChargingWalletDeposit>>('/v1/charging-wallet/deposits/pending', params)
}

/**
 * 确认充值(管理员)
 */
export function confirmChargingDeposit(id: number): Promise<void> {
  return post<void>(`/v1/charging-wallet/deposits/${id}/confirm`)
}

/**
 * 拒绝充值(管理员)
 */
export function rejectChargingDeposit(id: number, reason?: string): Promise<void> {
  return post<void>(`/v1/charging-wallet/deposits/${id}/reject`, { reason })
}

// ========== 奖励发放 ==========

/**
 * 发放奖励给下级
 */
export function issueChargingReward(data: IssueChargingRewardParams): Promise<{ reward_no: string }> {
  return post<{ reward_no: string }>('/v1/charging-wallet/rewards', data)
}

/**
 * 获取奖励记录列表
 */
export function getChargingRewardList(
  params: ChargingRewardQueryParams
): Promise<PaginatedResponse<ChargingWalletReward>> {
  return get<PaginatedResponse<ChargingWalletReward>>('/v1/charging-wallet/rewards', params)
}
