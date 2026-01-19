// 终端状态
export type TerminalStatus = 'stock' | 'distributed' | 'activated' | 'returned'

// 终端信息
export interface Terminal {
  id: number
  sn: string
  channel_id: number
  channel_name: string
  owner_agent_id: number
  owner_agent_name: string
  merchant_id: number
  merchant_name: string
  status: TerminalStatus
  activated_at: string
  created_at: string
}

// 终端详情
export interface TerminalDetail extends Terminal {
  // 费率设置
  credit_rate: number
  debit_rate: number
  debit_cap: number

  // SIM卡设置
  sim_first_fee: number
  sim_renewal_fee: number
  sim_renewal_days: number

  // 押金设置
  deposit_amount: number

  // 流转历史
  history: TerminalHistory[]
}

// 终端流转历史
export interface TerminalHistory {
  id: number
  action: 'import' | 'dispatch' | 'recall' | 'activate' | 'bind' | 'unbind'
  from_agent_id: number
  from_agent_name: string
  to_agent_id: number
  to_agent_name: string
  operator_id: number
  operator_name: string
  remark: string
  created_at: string
}

// 终端统计
export interface TerminalStats {
  total: number
  stock: number
  distributed: number
  activated: number
  returned: number
  yesterday_activated: number
  today_activated: number
  month_activated: number
}

// 终端下发参数
export interface TerminalDispatchParams {
  terminal_ids: number[]
  to_agent_id: number
  cargo_deduction?: {
    unit_price: number
    wallet_sources: string[]
  }
}

// 终端查询参数
export interface TerminalQueryParams {
  channel_id?: number
  status?: TerminalStatus
  owner_agent_id?: number
  sn?: string
  page?: number
  page_size?: number
}
