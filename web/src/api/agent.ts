import { get, post, put } from './request'
import type {
  Agent,
  AgentDetail,
  AgentStats,
  AgentTreeNode,
  AgentQueryParams,
  PaginatedResponse,
} from '@/types'

/**
 * 获取代理商列表
 */
export function getAgents(params: AgentQueryParams): Promise<PaginatedResponse<Agent>> {
  return get<PaginatedResponse<Agent>>('/v1/agents', params)
}

/**
 * 获取代理商详情
 */
export function getAgentDetail(id: number): Promise<AgentDetail> {
  return get<AgentDetail>(`/v1/agents/${id}`)
}

/**
 * 获取下级代理商列表
 */
export function getSubordinates(params: AgentQueryParams): Promise<PaginatedResponse<Agent>> {
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
export function searchAgents(keyword: string): Promise<Agent[]> {
  return get<Agent[]>('/v1/agents/search', { keyword, limit: 20 })
}
