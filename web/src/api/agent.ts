import { get, post, put } from './request'
import type {
  Agent,
  AgentDetail,
  AgentStats,
  AgentTreeNode,
  AgentQueryParams,
  AgentPolicy,
  Wallet,
  Merchant,
  Terminal,
  Transaction,
  PaginatedResponse,
} from '@/types'

/**
 * 获取代理商列表（当前代理商的下级列表）
 */
export function getAgents(params: AgentQueryParams): Promise<PaginatedResponse<Agent>> {
  return get<PaginatedResponse<Agent>>('/v1/agents/subordinates', params)
}

/**
 * 获取代理商详情
 */
export function getAgentDetail(id: number): Promise<AgentDetail> {
  return get<AgentDetail>('/v1/agents/detail', { id })
}

/**
 * 获取下级代理商列表
 */
export function getSubordinates(params: AgentQueryParams & { parent_id?: number }): Promise<PaginatedResponse<Agent>> {
  return get<PaginatedResponse<Agent>>('/v1/agents/subordinates', params)
}

/**
 * 获取代理商统计
 */
export function getAgentStats(): Promise<AgentStats> {
  return get<AgentStats>('/v1/agents/stats')
}

/**
 * 获取团队层级树
 */
export function getAgentTree(agentId?: number): Promise<AgentTreeNode[]> {
  return get<AgentTreeNode[]>('/v1/agents/team-tree', { agent_id: agentId })
}

/**
 * 获取邀请码
 */
export function getInviteCode(): Promise<{ invite_code: string; qr_code_url: string }> {
  return get('/v1/agents/invite-code')
}

/**
 * 更新代理商资料
 */
export function updateAgentProfile(data: Partial<Agent>): Promise<void> {
  return put<void>('/v1/agents/profile', data)
}

/**
 * 更新代理商状态
 */
export function updateAgentStatus(id: number, status: number): Promise<void> {
  return put<void>(`/v1/agents/${id}/status`, { status })
}

/**
 * 搜索代理商（用于选择器）
 */
export async function searchAgents(keyword: string): Promise<Agent[]> {
  const res = await get<{ list: Agent[]; total: number }>('/v1/agents/search', { keyword, page_size: 20 })
  return res.list || []
}

/**
 * 搜索代理商列表（分页，用于管理端选择）
 */
export function searchAgentList(params: {
  keyword?: string
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Agent>> {
  return get<PaginatedResponse<Agent>>('/v1/admin/agents', params)
}

/**
 * 创建代理商
 */
export function createAgent(data: {
  agent_name: string
  contact_name: string
  contact_phone: string
  id_card_no?: string
  bank_name?: string
  bank_account?: string
  bank_card_no?: string
  parent_id?: number
}): Promise<AgentDetail> {
  return post<AgentDetail>('/v1/agents', data)
}

// ======== 代理商详情页Tab相关API ========

/**
 * 获取代理商政策列表
 */
export function getAgentPolicies(agentId: number): Promise<AgentPolicy[]> {
  return get<AgentPolicy[]>('/v1/policies/agent', { agent_id: agentId })
}

/**
 * 获取代理商钱包列表
 */
export function getAgentWallets(agentId: number): Promise<Wallet[]> {
  return get<Wallet[]>('/v1/wallets', { agent_id: agentId })
}

/**
 * 获取代理商的下级代理列表
 */
export function getAgentSubordinates(agentId: number, params?: {
  keyword?: string
  status?: number
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Agent>> {
  return get<PaginatedResponse<Agent>>('/v1/agents/subordinates', { ...params, parent_id: agentId })
}

/**
 * 获取代理商的商户列表
 */
export function getAgentMerchants(agentId: number, params?: {
  keyword?: string
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Merchant>> {
  return get<PaginatedResponse<Merchant>>('/v1/merchants', { ...params, agent_id: agentId })
}

/**
 * 获取代理商的终端列表
 */
export function getAgentTerminals(agentId: number, params?: {
  keyword?: string
  status?: number
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Terminal>> {
  return get<PaginatedResponse<Terminal>>('/v1/terminals', { ...params, agent_id: agentId })
}

/**
 * 获取代理商的交易记录
 */
export function getAgentTransactions(agentId: number, params?: {
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<Transaction>> {
  return get<PaginatedResponse<Transaction>>('/v1/transactions', { ...params, agent_id: agentId })
}
