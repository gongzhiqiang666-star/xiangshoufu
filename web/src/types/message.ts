// 消息类型编号
export type MessageTypeValue = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8

// 消息类型信息
export interface MessageType {
  value: MessageTypeValue
  label: string
  category: string
}

// 消息分类
export interface MessageCategory {
  value: string
  label: string
}

// 消息信息
export interface Message {
  id: number
  agent_id?: number
  message_type: MessageTypeValue
  type_name?: string
  title: string
  content: string
  is_read: boolean
  is_pushed?: boolean
  related_id?: number
  related_type?: string
  expire_at?: string
  created_at: string
}

// 未读消息统计
export interface UnreadCount {
  count: number
}

// 消息统计
export interface MessageStats {
  total: number
  unread_total: number
  profit_count: number
  register_count: number
  consumption_count: number
  system_count: number
}

// 发送消息请求
export interface SendMessageRequest {
  title: string
  content: string
  message_type: MessageTypeValue
  expire_days?: number
  send_scope: 'all' | 'agents' | 'level'
  agent_ids?: number[]
  level?: number
}

// 消息类型配置
export const MESSAGE_TYPE_CONFIG: Record<MessageTypeValue, { label: string; color: string; category: string }> = {
  1: { label: '交易分润', color: '#67c23a', category: 'profit' },
  2: { label: '激活奖励', color: '#e6a23c', category: 'profit' },
  3: { label: '押金返现', color: '#409eff', category: 'profit' },
  4: { label: '流量返现', color: '#909399', category: 'profit' },
  5: { label: '退款撤销', color: '#f56c6c', category: 'system' },
  6: { label: '系统公告', color: '#409eff', category: 'system' },
  7: { label: '新代理注册', color: '#67c23a', category: 'register' },
  8: { label: '交易通知', color: '#e6a23c', category: 'consumption' },
}

// 消息分类配置
export const MESSAGE_CATEGORY_CONFIG: Record<string, { label: string; types: MessageTypeValue[] }> = {
  all: { label: '全部', types: [1, 2, 3, 4, 5, 6, 7, 8] },
  profit: { label: '分润', types: [1, 2, 3, 4] },
  register: { label: '注册', types: [7] },
  consumption: { label: '消费', types: [8] },
  system: { label: '系统', types: [5, 6] },
}
