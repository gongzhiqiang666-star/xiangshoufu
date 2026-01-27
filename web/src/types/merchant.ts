// 商户类型（5档分类）
export type MerchantType = 'quality' | 'medium' | 'normal' | 'warning' | 'churned'

// 商户信息
export interface Merchant {
  id: number
  merchant_code: string
  name: string
  phone: string
  phone_masked: string
  agent_id: number
  agent_name: string
  is_direct: boolean
  merchant_type: MerchantType
  channel_id: number
  channel_name: string
  terminal_sn: string
  activated_at: string
  created_at: string
}

// 商户详情
export interface MerchantDetail extends Merchant {
  // 流量费信息
  first_sim_fee_time: string
  first_sim_fee_amount: number
  last_renewal_time: string
  renewal_count: number

  // 费率信息
  credit_rate: number
  debit_rate: number
  debit_cap: number
  t0_fee: number

  // 交易统计
  total_transaction_amount: number
  month_transaction_amount: number
  month_credit_amount: number
  month_debit_amount: number
  month_wechat_amount: number
  month_alipay_amount: number

  // 登记信息
  registered_phone: string
  register_remark: string
}

// 商户统计（5档分类）
export interface MerchantStats {
  total: number
  direct: number
  team: number
  quality: number   // 优质
  medium: number    // 中等
  normal: number    // 普通
  warning: number   // 预警
  churned: number   // 流失
}

// 商户类型配置（5档分类）
export const MERCHANT_TYPE_CONFIG: Record<MerchantType, { label: string; color: string }> = {
  quality: { label: '优质商户', color: '#67c23a' },
  medium: { label: '中等商户', color: '#409eff' },
  normal: { label: '普通商户', color: '#909399' },
  warning: { label: '预警商户', color: '#e6a23c' },
  churned: { label: '流失商户', color: '#f56c6c' },
}

// 商户查询参数
export interface MerchantQueryParams {
  channel_id?: number
  merchant_type?: MerchantType
  owner_type?: 'all' | 'direct' | 'team'
  keyword?: string
  page?: number
  page_size?: number
}
