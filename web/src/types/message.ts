// 消息类型
export type MessageType = 'system' | 'alert' | 'notification'

// 消息信息
export interface Message {
  id: number
  type: MessageType
  title: string
  content: string
  is_read: boolean
  created_at: string
}

// 未读消息统计
export interface UnreadCount {
  total: number
  system: number
  alert: number
  notification: number
}
