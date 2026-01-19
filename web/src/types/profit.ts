// 分润类型
export type ProfitType = 'transaction' | 'deposit_cashback' | 'sim_cashback' | 'activation_reward'

// 分润记录
export interface Profit {
  id: number
  profit_no: string
  agent_id: number
  agent_name: string
  channel_id: number
  channel_name: string
  profit_type: ProfitType
  related_no: string
  transaction_amount: number
  profit_amount: number
  wallet_type: string
  status: number
  created_at: string
}

// 分润统计
export interface ProfitStats {
  transaction_profit: number
  deposit_cashback: number
  sim_cashback: number
  activation_reward: number
  total: number
}

// 分润日汇总
export interface ProfitDailySummary {
  date: string
  transaction_profit: number
  deposit_cashback: number
  sim_cashback: number
  activation_reward: number
  total: number
}

// 分润类型配置
export const PROFIT_TYPE_CONFIG: Record<ProfitType, { label: string; color: string }> = {
  transaction: { label: '交易分润', color: '#409eff' },
  deposit_cashback: { label: '押金返现', color: '#67c23a' },
  sim_cashback: { label: '流量返现', color: '#e6a23c' },
  activation_reward: { label: '激活奖励', color: '#f56c6c' },
}

// 分润查询参数
export interface ProfitQueryParams {
  channel_id?: number
  profit_type?: ProfitType
  agent_id?: number
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

// 手动调账
export interface AdjustmentRequest {
  agent_id: number
  adjust_type: 'increase' | 'decrease'
  wallet_type: string
  amount: number
  reason: string
}

// 调账记录
export interface Adjustment {
  id: number
  adjustment_no: string
  agent_id: number
  agent_name: string
  adjust_type: 'increase' | 'decrease'
  wallet_type: string
  amount: number
  reason: string
  status: 'pending' | 'approved' | 'rejected'
  applicant_id: number
  applicant_name: string
  auditor_id: number
  auditor_name: string
  audit_remark: string
  created_at: string
  audited_at: string
}
