// 钱包类型
export type WalletType = 'profit' | 'service' | 'reward' | 'recharge' | 'deposit'

// 钱包信息
export interface Wallet {
  id: number
  agent_id: number
  agent_name: string
  channel_id: number
  channel_name: string
  wallet_type: WalletType
  balance: number
  frozen: number
  available: number
  total_income: number
  total_withdraw: number
  updated_at: string
}

// 钱包流水
export interface WalletLog {
  id: number
  wallet_id: number
  log_no: string
  type: 'income' | 'expense' | 'freeze' | 'unfreeze' | 'withdraw'
  amount: number
  balance_before: number
  balance_after: number
  related_no: string
  remark: string
  created_at: string
}

// 钱包汇总
export interface WalletSummary {
  profit_balance: number
  service_balance: number
  reward_balance: number
  recharge_balance: number
  deposit_balance: number
  total_balance: number
  total_available: number
  total_frozen: number
}

// 钱包类型配置
export const WALLET_TYPE_CONFIG: Record<WalletType, { label: string; color: string; description: string }> = {
  profit: { label: '分润钱包', color: '#409eff', description: '交易分润收入' },
  service: { label: '服务费钱包', color: '#67c23a', description: '流量费+押金返现' },
  reward: { label: '奖励钱包', color: '#e6a23c', description: '激活奖励收入' },
  recharge: { label: '充值钱包', color: '#f56c6c', description: '上级额外奖励' },
  deposit: { label: '沉淀钱包', color: '#909399', description: '下级未提现余额' },
}

// 提现申请
export interface WithdrawRequest {
  wallet_id: number
  amount: number
}

// 提现记录
export interface Withdrawal {
  id: number
  withdrawal_no: string
  agent_id: number
  agent_name: string
  channel_id: number
  channel_name: string
  wallet_type: WalletType
  amount: number
  fee: number
  actual_amount: number
  bank_name: string
  bank_card_no: string
  status: 'pending' | 'approved' | 'rejected' | 'paid' | 'failed'
  remark: string
  created_at: string
  audited_at: string
  paid_at: string
}

// 钱包流水查询参数
export interface WalletLogQueryParams {
  wallet_id?: number
  type?: string
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}
