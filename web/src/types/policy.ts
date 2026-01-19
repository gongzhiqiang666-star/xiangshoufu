// 政策模板
export interface PolicyTemplate {
  id: number
  name: string
  channel_id: number
  channel_name: string
  is_default: boolean
  status: number

  // 费率设置
  credit_rate: number
  debit_rate: number
  debit_cap: number
  unionpay_qr_rate: number
  wechat_rate: number
  alipay_rate: number
  t0_fee: number

  created_at: string
  updated_at: string
}

// 政策模板详情
export interface PolicyTemplateDetail extends PolicyTemplate {
  // 费率阶梯
  rate_stages: RateStage[]

  // 激活奖励
  activation_rewards: ActivationReward[]

  // 押金返现
  deposit_cashbacks: DepositCashback[]

  // 流量返现
  sim_cashbacks: SimCashback[]
}

// 费率阶梯
export interface RateStage {
  id: number
  template_id: number
  start_day: number
  end_day: number
  rate_adjustment: number
  base_type: 'merchant' | 'agent'
}

// 激活奖励
export interface ActivationReward {
  id: number
  template_id: number
  start_day: number
  end_day: number
  target_amount: number
  reward_amount: number
  is_multi_level: boolean
  effective_date: string
}

// 押金返现
export interface DepositCashback {
  id: number
  template_id: number
  deposit_amount: number
  cashback_amount: number
}

// 流量返现
export interface SimCashback {
  id: number
  template_id: number
  type: 'first' | 'second' | 'renewal'
  sim_fee: number
  cashback_amount: number
  is_multi_level: boolean
}

// 代理商政策
export interface AgentPolicy {
  id: number
  agent_id: number
  agent_name: string
  channel_id: number
  channel_name: string
  template_id: number
  template_name: string

  // 实际费率
  credit_rate: number
  debit_rate: number
  debit_cap: number

  created_at: string
}

// 政策模板查询参数
export interface PolicyQueryParams {
  channel_id?: number
  status?: number
  keyword?: string
  page?: number
  page_size?: number
}

// 创建/更新政策模板
export interface PolicyTemplateForm {
  name: string
  channel_id: number
  is_default: boolean
  credit_rate: number
  debit_rate: number
  debit_cap: number
  unionpay_qr_rate: number
  wechat_rate: number
  alipay_rate: number
  t0_fee: number
  rate_stages: Omit<RateStage, 'id' | 'template_id'>[]
  activation_rewards: Omit<ActivationReward, 'id' | 'template_id'>[]
  deposit_cashbacks: Omit<DepositCashback, 'id' | 'template_id'>[]
  sim_cashbacks: Omit<SimCashback, 'id' | 'template_id'>[]
}
