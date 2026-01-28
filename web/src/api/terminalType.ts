import { get, post, put, patch, del } from './request'
import type {
  TerminalType,
  TerminalTypeQueryParams,
  CreateTerminalTypeParams,
  UpdateTerminalTypeParams,
} from '@/types/terminalType'
import type { PaginatedResponse } from '@/types'

/**
 * 获取终端类型列表
 */
export function getTerminalTypes(params: TerminalTypeQueryParams): Promise<PaginatedResponse<TerminalType>> {
  return get<PaginatedResponse<TerminalType>>('/v1/admin/terminal-types', params)
}

/**
 * 获取终端类型详情
 */
export function getTerminalType(id: number): Promise<TerminalType> {
  return get<TerminalType>(`/v1/admin/terminal-types/${id}`)
}

/**
 * 创建终端类型
 */
export function createTerminalType(data: CreateTerminalTypeParams): Promise<TerminalType> {
  return post<TerminalType>('/v1/admin/terminal-types', data)
}

/**
 * 更新终端类型
 */
export function updateTerminalType(id: number, data: UpdateTerminalTypeParams): Promise<TerminalType> {
  return put<TerminalType>(`/v1/admin/terminal-types/${id}`, data)
}

/**
 * 更新终端类型状态
 */
export function updateTerminalTypeStatus(id: number, status: number): Promise<void> {
  return patch<void>(`/v1/admin/terminal-types/${id}/status`, { status })
}

/**
 * 删除终端类型
 */
export function deleteTerminalType(id: number): Promise<void> {
  return del<void>(`/v1/admin/terminal-types/${id}`)
}

/**
 * 根据通道ID获取终端类型列表（用于下拉选择）
 */
export function getTerminalTypesByChannel(channelId: number): Promise<TerminalType[]> {
  return get<TerminalType[]>(`/v1/admin/terminal-types/by-channel/${channelId}`)
}
