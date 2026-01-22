// 税筹通道扣费类型
export type TaxFeeType = 1 | 2 // 1=付款扣 2=出款扣

// 税筹通道状态
export type TaxChannelStatus = 0 | 1 // 0=禁用 1=启用

// 税筹通道
export interface TaxChannel {
  id: number
  channel_code: string
  channel_name: string
  fee_type: TaxFeeType
  fee_type_name: string
  tax_rate: number
  tax_rate_percent: number // 百分比显示
  fixed_fee: number
  fixed_fee_yuan: number
  status: TaxChannelStatus
  status_name: string
  remark: string
  created_at: string
  updated_at: string
}

// 创建税筹通道请求
export interface CreateTaxChannelParams {
  channel_code: string
  channel_name: string
  fee_type: TaxFeeType
  tax_rate: number // 小数形式，如0.09
  fixed_fee?: number // 分
  api_url?: string
  api_key?: string
  api_secret?: string
  remark?: string
}

// 更新税筹通道请求
export interface UpdateTaxChannelParams {
  channel_name?: string
  fee_type?: TaxFeeType
  tax_rate?: number
  fixed_fee?: number
  api_url?: string
  api_key?: string
  api_secret?: string
  status?: TaxChannelStatus
  remark?: string
}

// 通道-税筹通道映射
export interface ChannelTaxMapping {
  id: number
  channel_id: number
  wallet_type: number
  wallet_type_name: string
  tax_channel_id: number
  tax_channel_name: string
  tax_rate: number
  fixed_fee: number
  created_at: string
}

// 设置通道税筹映射请求
export interface SetChannelTaxMappingParams {
  channel_id: number
  wallet_type: number // 1=分润 2=服务费 3=奖励 4=充值 5=沉淀
  tax_channel_id: number
}

// 计算税费请求
export interface CalculateTaxParams {
  channel_id: number
  wallet_type: number
  amount: number // 分
}

// 税费计算结果
export interface TaxCalculationResult {
  original_amount: number // 原金额(分)
  tax_rate: number // 税率
  tax_fee: number // 税费(分)
  fixed_fee: number // 固定费用(分)
  total_fee: number // 总费用(分)
  actual_amount: number // 实际到账(分)
  tax_channel_id: number
  tax_channel_name: string
}

// 税筹通道查询参数
export interface TaxChannelQueryParams {
  status?: TaxChannelStatus
}

// 扣费类型配置
export const TAX_FEE_TYPE_CONFIG: Record<TaxFeeType, { label: string; description: string }> = {
  1: { label: '付款扣', description: '充值时扣除税费' },
  2: { label: '出款扣', description: '提现时扣除税费' },
}

// 税筹通道状态配置
export const TAX_CHANNEL_STATUS_CONFIG: Record<TaxChannelStatus, { label: string; color: string }> = {
  0: { label: '禁用', color: '#909399' },
  1: { label: '启用', color: '#67c23a' },
}
