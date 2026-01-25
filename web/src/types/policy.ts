// 政策模板
export interface PolicyTemplate {
  id: number
  name: string
  channel_id: number
  channel_name: string
  is_default: boolean
  status: number

  // 动态费率配置
  rate_configs: RateConfigs

  created_at: string
  updated_at: string
}

// 费率类型定义（从通道配置读取）
export interface RateTypeDefinition {
  code: string       // 费率类型编码
  name: string       // 费率类型名称
  sort_order: number // 排序
  min_rate: string   // 费率下限
  max_rate: string   // 费率上限
}

// 费率配置值
export interface RateConfigValue {
  rate: string // 费率值（百分比）
}

// 费率配置集合
export type RateConfigs = Record<string, RateConfigValue>

// 费率阶梯调整值
export type RateDeltas = Record<string, string>

// 通道
export interface Channel {
  id: number
  channel_code: string
  channel_name: string
  description: string
  status: number
  priority: number
  config: string
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
  // 动态费率配置
  rate_configs: RateConfigs
  rate_stages: Omit<RateStage, 'id' | 'template_id'>[]
  activation_rewards: Omit<ActivationReward, 'id' | 'template_id'>[]
  deposit_cashbacks: Omit<DepositCashback, 'id' | 'template_id'>[]
  sim_cashbacks: Omit<SimCashback, 'id' | 'template_id'>[]
}
