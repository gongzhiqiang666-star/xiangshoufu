// 交易类型
export type TransactionType = 'credit' | 'debit' | 'wechat' | 'alipay' | 'unionpay_qr'

// 交易状态
export type TransactionStatus = 'success' | 'failed' | 'pending' | 'refunded'

// 交易信息
export interface Transaction {
  id: number
  transaction_no: string
  channel_transaction_no: string
  channel_id: number
  channel_name: string
  merchant_id: number
  merchant_name: string
  terminal_sn: string
  transaction_type: TransactionType
  amount: number
  fee: number
  status: TransactionStatus
  transaction_time: string
  created_at: string
}

// 交易详情
export interface TransactionDetail extends Transaction {
  agent_id: number
  agent_name: string

  // 费率信息
  rate: number
  t0_fee: number

  // 分润信息
  profit_distributed: boolean
  profit_records: {
    agent_id: number
    agent_name: string
    profit_amount: number
  }[]
}

// 交易统计
export interface TransactionStats {
  total_amount: number
  total_count: number
  today_amount: number
  today_count: number
  month_amount: number
  month_count: number
}

// 交易趋势
export interface TransactionTrend {
  dates: string[]
  amounts: number[]
  counts: number[]
}

// 交易类型配置
export const TRANSACTION_TYPE_CONFIG: Record<TransactionType, { label: string; color: string }> = {
  credit: { label: '贷记卡', color: '#409eff' },
  debit: { label: '借记卡', color: '#67c23a' },
  wechat: { label: '微信', color: '#07c160' },
  alipay: { label: '支付宝', color: '#1677ff' },
  unionpay_qr: { label: '云闪付', color: '#e60012' },
}

// 交易查询参数
export interface TransactionQueryParams {
  channel_id?: number
  transaction_type?: TransactionType
  status?: TransactionStatus
  merchant_id?: number
  terminal_sn?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}
