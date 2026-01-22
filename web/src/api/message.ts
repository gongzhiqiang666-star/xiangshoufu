import { get, post, del, put } from './request'
import type { Message, UnreadCount, PaginatedResponse, PaginationParams, MessageStats, MessageType, MessageCategory, SendMessageRequest } from '@/types'

/**
 * 获取消息列表
 */
export function getMessages(
  params: PaginationParams & { type?: string; category?: string; is_read?: boolean }
): Promise<PaginatedResponse<Message>> {
  return get<PaginatedResponse<Message>>('/v1/messages', params)
}

/**
 * 获取未读消息数量
 */
export function getUnreadCount(): Promise<UnreadCount> {
  return get<UnreadCount>('/v1/messages/unread-count')
}

/**
 * 获取消息统计
 */
export function getMessageStats(): Promise<MessageStats> {
  return get<MessageStats>('/v1/messages/stats')
}

/**
 * 获取消息类型和分类
 */
export function getMessageTypes(): Promise<{ types: MessageType[]; categories: MessageCategory[] }> {
  return get<{ types: MessageType[]; categories: MessageCategory[] }>('/v1/messages/types')
}

/**
 * 获取消息详情
 */
export function getMessageDetail(id: number): Promise<Message> {
  return get<Message>(`/v1/messages/${id}`)
}

/**
 * 标记消息为已读
 */
export function markAsRead(id: number): Promise<void> {
  return put<void>(`/v1/messages/${id}/read`)
}

/**
 * 标记所有消息为已读
 */
export function markAllAsRead(): Promise<void> {
  return put<void>('/v1/messages/read-all')
}

// ============================================================
// 管理端API
// ============================================================

/**
 * 获取管理端消息列表
 */
export function getAdminMessages(
  params: PaginationParams
): Promise<PaginatedResponse<Message>> {
  return get<PaginatedResponse<Message>>('/v1/admin/messages', params)
}

/**
 * 获取管理端消息类型
 */
export function getAdminMessageTypes(): Promise<MessageType[]> {
  return get<MessageType[]>('/v1/admin/messages/types')
}

/**
 * 获取管理端消息详情
 */
export function getAdminMessageDetail(id: number): Promise<Message> {
  return get<Message>(`/v1/admin/messages/${id}`)
}

/**
 * 发送消息
 */
export function sendMessage(data: SendMessageRequest): Promise<{ sent_count: number }> {
  return post<{ sent_count: number }>('/v1/admin/messages', data)
}

/**
 * 删除消息
 */
export function deleteMessage(id: number): Promise<void> {
  return del<void>(`/v1/admin/messages/${id}`)
}
