// 代理商信息
export interface Agent {
  id: number
  agent_code: string
  name: string
  phone: string
  id_card_no: string
  parent_id: number
  parent_name: string
  level: number
  invite_code: string
  status: number
  channel_id: number
  channel_name: string
  created_at: string
  updated_at: string
}

// 代理商详情
export interface AgentDetail extends Agent {
  // 结算信息
  bank_name: string
  bank_card_no: string
  bank_branch: string

  // 统计数据
  direct_agent_count: number
  team_agent_count: number
  direct_merchant_count: number
  team_merchant_count: number
  terminal_total: number
  terminal_activated: number
  month_transaction_amount: number
  total_profit: number
}

// 代理商统计
export interface AgentStats {
  total: number
  active: number
  disabled: number
  today_new: number
  month_new: number
}

// 代理商树节点
export interface AgentTreeNode {
  id: number
  agent_code: string
  name: string
  level: number
  children?: AgentTreeNode[]
}

// 代理商查询参数
export interface AgentQueryParams {
  channel_id?: number
  status?: number
  keyword?: string
  page?: number
  page_size?: number
}

// 代理商通道配置
export interface AgentChannel {
  id: number
  agent_id: number
  channel_id: number
  is_enabled: boolean
  is_visible: boolean
  enabled_at: string | null
  disabled_at: string | null
  enabled_by: number | null
  disabled_by: number | null
  remark: string
  created_at: string
  updated_at: string
  // 关联字段
  channel_code: string
  channel_name: string
}

// 代理商通道统计
export interface AgentChannelStats {
  total_channels: number
  enabled_channels: number
  visible_channels: number
}
