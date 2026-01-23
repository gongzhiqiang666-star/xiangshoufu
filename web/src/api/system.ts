// 系统管理API
import request from './request'
import type { PageResult } from '@/types'
import type { SystemUser, CreateUserRequest, UpdateUserRequest, OperationLog } from '@/types/system'

// 获取用户列表
export function getUsers(params: {
  keyword?: string
  role?: string
  status?: number
  page?: number
  page_size?: number
}) {
  return request.get<PageResult<SystemUser>>('/v1/system/users', { params })
}

// 获取单个用户
export function getUser(id: number) {
  return request.get<SystemUser>(`/v1/system/users/${id}`)
}

// 创建用户
export function createUser(data: CreateUserRequest) {
  return request.post<SystemUser>('/v1/system/users', data)
}

// 更新用户
export function updateUser(id: number, data: UpdateUserRequest) {
  return request.put<SystemUser>(`/v1/system/users/${id}`, data)
}

// 删除用户
export function deleteUser(id: number) {
  return request.delete(`/v1/system/users/${id}`)
}

// 重置密码
export function resetPassword(id: number, password: string) {
  return request.post(`/v1/system/users/${id}/reset-password`, { password })
}

// 启用/禁用用户
export function toggleUserStatus(id: number, status: number) {
  return request.patch(`/v1/system/users/${id}/status`, { status })
}

// 获取操作日志列表
export function getLogs(params: {
  user_id?: number
  username?: string
  module?: string
  action?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}) {
  return request.get<PageResult<OperationLog>>('/v1/system/logs', { params })
}

// 获取单条日志详情
export function getLog(id: number) {
  return request.get<OperationLog>(`/v1/system/logs/${id}`)
}

// 导出日志
export function exportLogs(params: {
  user_id?: number
  module?: string
  start_date?: string
  end_date?: string
}) {
  return request.get('/v1/system/logs/export', {
    params,
    responseType: 'blob',
  })
}
