import { get, post, put, del } from './request'

// 奖励模版相关类型
export interface RewardStage {
  id?: number
  template_id?: number
  stage_order: number
  start_value: number
  end_value: number
  target_value: number
  reward_amount: number
}

export interface RewardTemplate {
  id: number
  name: string
  time_type: 'days' | 'months'
  dimension_type: 'amount' | 'count'
  trans_types: string
  amount_min: number | null
  amount_max: number | null
  allow_gap: boolean
  enabled: boolean
  created_at: string
  updated_at: string
  stages?: RewardStage[]
}

export interface CreateRewardTemplateRequest {
  name: string
  time_type: 'days' | 'months'
  dimension_type: 'amount' | 'count'
  trans_types: string
  amount_min?: number
  amount_max?: number
  allow_gap: boolean
  stages: Omit<RewardStage, 'id' | 'template_id'>[]
}

export interface UpdateRewardTemplateRequest extends CreateRewardTemplateRequest {}

export interface AgentRewardRate {
  id: number
  agent_id: number
  template_id: number
  reward_amount: number  // 奖励金额（分）- 差额分配模式
  created_at: string
  updated_at: string
}

export interface TerminalRewardProgress {
  id: number
  terminal_sn: string
  template_id: number
  bind_agent_id: number
  bind_time: string
  current_stage: number
  status: 'active' | 'completed' | 'terminated'
  stage_rewards?: TerminalStageReward[]
}

export interface TerminalStageReward {
  id: number
  terminal_sn: string
  stage_order: number
  stage_start: string
  stage_end: string
  target_value: number
  actual_value: number
  is_achieved: boolean
  reward_amount: number
  status: 'pending' | 'achieved' | 'failed' | 'gap_blocked'
}

export interface RewardOverflowLog {
  id: number
  terminal_sn: string
  stage_reward_id: number
  total_rate: number
  agent_chain: string
  resolved: boolean
  resolved_at: string | null
  resolved_by: string | null
  created_at: string
}

export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// ======== 奖励模版管理 ========

export function getRewardTemplates(params: {
  enabled?: boolean
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<RewardTemplate>> {
  return get<PaginatedResponse<RewardTemplate>>('/v1/rewards/templates', params)
}

export function getRewardTemplateDetail(id: number): Promise<RewardTemplate> {
  return get<RewardTemplate>(`/v1/rewards/templates/${id}`)
}

export function createRewardTemplate(data: CreateRewardTemplateRequest): Promise<{ id: number }> {
  return post<{ id: number }>('/v1/rewards/templates', data)
}

export function updateRewardTemplate(id: number, data: UpdateRewardTemplateRequest): Promise<void> {
  return put<void>(`/v1/rewards/templates/${id}`, data)
}

export function deleteRewardTemplate(id: number): Promise<void> {
  return del<void>(`/v1/rewards/templates/${id}`)
}

export function updateRewardTemplateStatus(id: number, enabled: boolean): Promise<void> {
  return put<void>(`/v1/rewards/templates/${id}/status`, { enabled })
}

// ======== 代理商奖励金额配置（差额分配模式） ========

export function getAgentRewardAmount(agentId: number, templateId: number): Promise<AgentRewardRate> {
  return get<AgentRewardRate>(`/v1/rewards/agents/${agentId}/amount`, { template_id: templateId })
}

export function setAgentRewardAmount(agentId: number, templateId: number, rewardAmount: number): Promise<void> {
  return put<void>(`/v1/rewards/agents/${agentId}/amount`, { template_id: templateId, reward_amount: rewardAmount })
}

// ======== 终端奖励进度 ========

export function getTerminalRewardProgress(terminalSn: string): Promise<TerminalRewardProgress> {
  return get<TerminalRewardProgress>(`/v1/rewards/terminals/${terminalSn}/progress`)
}

export function initTerminalRewardProgress(data: {
  terminal_sn: string
  terminal_id?: number
  agent_id: number
  template_id: number
}): Promise<TerminalRewardProgress> {
  return post<TerminalRewardProgress>('/v1/rewards/terminals/progress', data)
}

// ======== 溢出日志 ========

export function getOverflowLogs(params: {
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<RewardOverflowLog>> {
  return get<PaginatedResponse<RewardOverflowLog>>('/v1/rewards/overflow-logs', params)
}

export function resolveOverflowLog(id: number): Promise<void> {
  return post<void>(`/v1/rewards/overflow-logs/${id}/resolve`)
}
