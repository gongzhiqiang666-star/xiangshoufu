import { get, post, put, del } from './request'
import type {
  PolicyTemplate,
  PolicyTemplateDetail,
  PolicyTemplateForm,
  AgentPolicy,
  PolicyQueryParams,
  PaginatedResponse,
} from '@/types'

/**
 * 获取政策模板列表
 */
export function getPolicyTemplates(
  params: PolicyQueryParams
): Promise<PaginatedResponse<PolicyTemplate>> {
  return get<PaginatedResponse<PolicyTemplate>>('/v1/policies/templates', params)
}

/**
 * 获取政策模板详情
 */
export function getPolicyTemplateDetail(id: number): Promise<PolicyTemplateDetail> {
  return get<PolicyTemplateDetail>(`/v1/policies/templates/${id}`)
}

/**
 * 获取政策模板详情 (别名)
 */
export function getPolicyTemplate(id: number): Promise<PolicyTemplateDetail> {
  return get<PolicyTemplateDetail>(`/v1/policies/templates/${id}`)
}

/**
 * 创建政策模板
 */
export function createPolicyTemplate(data: PolicyTemplateForm): Promise<{ id: number }> {
  return post<{ id: number }>('/v1/policies/templates', data)
}

/**
 * 更新政策模板
 */
export function updatePolicyTemplate(id: number, data: PolicyTemplateForm): Promise<void> {
  return put<void>(`/v1/policies/templates/${id}`, data)
}

/**
 * 删除政策模板
 */
export function deletePolicyTemplate(id: number): Promise<void> {
  return del<void>(`/v1/policies/templates/${id}`)
}

/**
 * 复制政策模板
 */
export function copyPolicyTemplate(id: number, newName: string): Promise<{ id: number }> {
  return post<{ id: number }>(`/v1/policies/templates/${id}/copy`, { new_name: newName })
}

/**
 * 设置默认模板
 */
export function setDefaultTemplate(id: number): Promise<void> {
  return put<void>(`/v1/policies/templates/${id}/default`)
}

// ======== 代理商政策 ========

/**
 * 获取代理商政策列表
 */
export function getAgentPolicies(params: {
  agent_id?: number
  channel_id?: number
  page?: number
  page_size?: number
}): Promise<PaginatedResponse<AgentPolicy>> {
  return get<PaginatedResponse<AgentPolicy>>('/v1/policies/agents', params)
}

/**
 * 获取我的政策
 */
export function getMyPolicies(): Promise<AgentPolicy[]> {
  return get<AgentPolicy[]>('/v1/policies/my')
}

/**
 * 分配政策给代理商
 */
export function applyPolicyToAgent(
  agentId: number,
  data: { channel_id: number; template_id: number }
): Promise<void> {
  return post<void>(`/v1/policies/agents/${agentId}/apply`, data)
}

/**
 * 获取通道默认模板
 */
export function getChannelDefaultTemplate(channelId: number): Promise<PolicyTemplate | null> {
  return get<PolicyTemplate | null>(`/v1/policies/channels/${channelId}/default`)
}

/**
 * 分配政策给代理商（带费率调整）
 */
export function assignAgentPolicy(
  agentId: number,
  data: {
    channel_id: number
    template_id: number
    credit_rate: number
    debit_rate: number
  }
): Promise<void> {
  return post<void>(`/v1/agents/${agentId}/policies`, data)
}

/**
 * 根据通道ID获取政策模板列表
 */
export function getTemplatesByChannel(channelId: number): Promise<PolicyTemplate[]> {
  return get<PolicyTemplate[]>('/v1/policies/templates', { channel_id: channelId })
}
