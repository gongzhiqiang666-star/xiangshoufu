import { get, post } from './request'
import type { Message, UnreadCount, PaginatedResponse, PaginationParams } from '@/types'

/**
 * 获取消息列表
 */
export function getMessages(
  params: PaginationParams & { type?: string; is_read?: boolean }
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
 * 标记消息为已读
 */
export function markAsRead(ids: number[]): Promise<void> {
  return post<void>('/v1/messages/mark-read', { ids })
}

/**
 * 标记所有消息为已读
 */
export function markAllAsRead(): Promise<void> {
  return post<void>('/v1/messages/mark-all-read')
}
