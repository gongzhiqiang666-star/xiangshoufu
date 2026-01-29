/**
 * 代扣管理 API
 */
import { get, post } from './request'
import type { PaginatedResponse } from '@/types'
import type {
  DeductionPlan,
  DeductionPlanDetail,
  DeductionRecord,
  DeductionPlanQueryParams,
  DeductionPlanStats,
  CreateDeductionPlanRequest,
} from '@/types/deduction'

// 获取代扣计划列表
export function getDeductionPlans(
  params: DeductionPlanQueryParams
): Promise<PaginatedResponse<DeductionPlan>> {
  return get<PaginatedResponse<DeductionPlan>>('/v1/deduction/plans', params)
}

// 获取代扣计划详情
export function getDeductionPlanDetail(id: number): Promise<DeductionPlanDetail> {
  return get<DeductionPlanDetail>(`/v1/deduction/plans/${id}`)
}

// 创建代扣计划
export function createDeductionPlan(data: CreateDeductionPlanRequest): Promise<DeductionPlan> {
  return post<DeductionPlan>('/v1/deduction/plans', data)
}

// 暂停代扣计划
export function pauseDeductionPlan(id: number): Promise<void> {
  return post<void>(`/v1/deduction/plans/${id}/pause`)
}

// 恢复代扣计划
export function resumeDeductionPlan(id: number): Promise<void> {
  return post<void>(`/v1/deduction/plans/${id}/resume`)
}

// 取消代扣计划
export function cancelDeductionPlan(id: number): Promise<void> {
  return post<void>(`/v1/deduction/plans/${id}/cancel`)
}

// 获取代扣记录列表
export function getDeductionRecords(
  planId: number,
  params?: { page?: number; page_size?: number }
): Promise<PaginatedResponse<DeductionRecord>> {
  return get<PaginatedResponse<DeductionRecord>>(`/v1/deduction/plans/${planId}/records`, params)
}

// 获取代扣计划统计
export function getDeductionPlanStats(params?: {
  plan_type?: number
  start_date?: string
  end_date?: string
}): Promise<DeductionPlanStats> {
  return get<DeductionPlanStats>('/v1/deduction/plans/stats', params)
}

// 导出代扣计划
export function exportDeductionPlans(
  params: DeductionPlanQueryParams
): Promise<{ task_id: string }> {
  return post<{ task_id: string }>('/v1/deduction/plans/export', params)
}
