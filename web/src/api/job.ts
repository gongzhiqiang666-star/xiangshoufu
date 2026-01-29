// 定时任务管理API
import { get, post, put, del } from './request'
import type { PaginatedResponse } from '@/types'
import type {
  JobListItem,
  JobDetail,
  JobExecutionLog,
  JobExecutionStats,
  UpdateJobConfigRequest,
  AlertConfig,
  AlertLog,
  CreateAlertConfigRequest,
  UpdateAlertConfigRequest
} from '@/types/job'

// ==================== 任务管理 ====================

// 获取任务列表
export function getJobs() {
  return get<JobListItem[]>('/v1/admin/jobs')
}

// 获取任务详情
export function getJob(name: string) {
  return get<JobDetail>(`/v1/admin/jobs/${name}`)
}

// 更新任务配置
export function updateJobConfig(name: string, data: UpdateJobConfigRequest) {
  return put<void>(`/v1/admin/jobs/${name}/config`, data)
}

// 手动触发任务
export function triggerJob(name: string) {
  return post<void>(`/v1/admin/jobs/${name}/trigger`)
}

// 启用/禁用任务
export function enableJob(name: string, isEnabled: boolean) {
  return put<void>(`/v1/admin/jobs/${name}/enable`, { is_enabled: isEnabled })
}

// ==================== 执行日志 ====================

// 获取执行日志列表
export function getJobLogs(params: {
  job_name?: string
  start_date?: string
  end_date?: string
  status?: number
  page?: number
  page_size?: number
}) {
  return get<PaginatedResponse<JobExecutionLog>>('/v1/admin/job-logs', params)
}

// 获取日志详情
export function getJobLog(id: number) {
  return get<JobExecutionLog>(`/v1/admin/job-logs/${id}`)
}

// 获取任务执行统计
export function getJobStats(params: {
  start_date: string
  end_date: string
}) {
  return get<JobExecutionStats[]>('/v1/admin/job-logs/stats', params)
}

// ==================== 告警配置 ====================

// 获取告警配置列表
export function getAlertConfigs() {
  return get<AlertConfig[]>('/v1/admin/alert-configs')
}

// 创建告警配置
export function createAlertConfig(data: CreateAlertConfigRequest) {
  return post<{ id: number }>('/v1/admin/alert-configs', data)
}

// 获取告警配置详情
export function getAlertConfig(id: number) {
  return get<AlertConfig>(`/v1/admin/alert-configs/${id}`)
}

// 更新告警配置
export function updateAlertConfig(id: number, data: UpdateAlertConfigRequest) {
  return put<void>(`/v1/admin/alert-configs/${id}`, data)
}

// 删除告警配置
export function deleteAlertConfig(id: number) {
  return del<void>(`/v1/admin/alert-configs/${id}`)
}

// 启用/禁用告警配置
export function enableAlertConfig(id: number, isEnabled: boolean) {
  return put<void>(`/v1/admin/alert-configs/${id}/enable`, { is_enabled: isEnabled })
}

// 测试告警配置
export function testAlertConfig(id: number) {
  return post<void>(`/v1/admin/alert-configs/${id}/test`)
}

// ==================== 告警日志 ====================

// 获取告警日志列表
export function getAlertLogs(params: {
  job_name?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}) {
  return get<PaginatedResponse<AlertLog>>('/v1/admin/alert-logs', params)
}

// 获取告警日志详情
export function getAlertLog(id: number) {
  return get<AlertLog>(`/v1/admin/alert-logs/${id}`)
}
