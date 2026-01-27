import { get, put } from './request'
import type {
  AgentWalletSplitConfig,
  PolicyWithdrawThreshold,
  WalletListWithSplitResponse,
  SetSplitConfigRequest,
  SetWithdrawThresholdRequest,
  WithdrawWithChannelRequest,
} from '@/types/wallet'

// ========== 钱包拆分配置 ==========

/**
 * 获取代理商钱包拆分配置
 */
export function getSplitConfig(agentId: number): Promise<AgentWalletSplitConfig> {
  return get<{ data: AgentWalletSplitConfig }>(`/v1/agents/${agentId}/wallet-split`).then(res => res.data)
}

/**
 * 设置代理商钱包拆分配置
 */
export function setSplitConfig(agentId: number, data: SetSplitConfigRequest): Promise<void> {
  return put<void>(`/v1/agents/${agentId}/wallet-split`, data)
}

/**
 * 检查代理商是否按通道拆分
 */
export function checkSplitStatus(agentId: number): Promise<{ split_by_channel: boolean }> {
  return get<{ data: { split_by_channel: boolean } }>(`/v1/agents/${agentId}/wallet-split/status`).then(res => res.data)
}

// ========== 提现门槛配置 ==========

/**
 * 获取政策模版提现门槛配置
 */
export function getWithdrawThresholds(templateId: number): Promise<PolicyWithdrawThreshold[]> {
  return get<{ data: PolicyWithdrawThreshold[] }>(`/v1/policy-templates/${templateId}/withdraw-thresholds`).then(res => res.data || [])
}

/**
 * 设置政策模版提现门槛
 */
export function setWithdrawThreshold(templateId: number, data: SetWithdrawThresholdRequest): Promise<void> {
  return put<void>(`/v1/policy-templates/${templateId}/withdraw-thresholds`, data)
}

/**
 * 批量设置政策模版提现门槛
 */
export function batchSetWithdrawThresholds(templateId: number, thresholds: SetWithdrawThresholdRequest[]): Promise<void> {
  return put<void>(`/v1/policy-templates/${templateId}/withdraw-thresholds/batch`, { thresholds })
}

// ========== 钱包展示（支持拆分模式） ==========

/**
 * 获取钱包列表（支持拆分模式）
 */
export function getWalletsWithSplit(): Promise<WalletListWithSplitResponse> {
  return get<{ data: WalletListWithSplitResponse }>('/v1/wallets/with-split').then(res => res.data)
}

/**
 * 提现（支持拆分模式）
 */
export function withdrawWithChannel(data: WithdrawWithChannelRequest): Promise<void> {
  return put<void>(`/v1/wallets/${data.wallet_id}/withdraw-channel`, data)
}
