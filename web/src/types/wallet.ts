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

// ========== 充值钱包相关类型 ==========

// 代理商钱包配置
export interface AgentWalletConfig {
  agent_id: number
  charging_wallet_enabled: boolean
  charging_wallet_limit: number
  charging_wallet_limit_yuan: number
  settlement_wallet_enabled: boolean
  settlement_ratio: number
  enabled_at?: string
}

// 开通充值钱包请求
export interface EnableChargingWalletParams {
  agent_id: number
  limit: number // 限额(分)
}

// 充值钱包汇总
export interface ChargingWalletSummary {
  balance: number
  balance_yuan: number
  total_issued: number
  total_issued_yuan: number
}

// 充值记录状态
export type ChargingDepositStatus = 0 | 1 | 2 // 0=待确认 1=已确认 2=已拒绝

// 充值记录
export interface ChargingWalletDeposit {
  id: number
  deposit_no: string
  agent_id: number
  amount: number
  amount_yuan: number
  payment_method: number
  payment_method_name: string
  payment_ref: string
  status: ChargingDepositStatus
  status_name: string
  reject_reason?: string
  remark?: string
  created_at: string
  confirmed_at?: string
}

// 申请充值请求
export interface CreateChargingDepositParams {
  amount: number // 分
  payment_method: number // 1=银行转账 2=微信 3=支付宝
  payment_ref?: string
  remark?: string
}

// 奖励发放记录
export interface ChargingWalletReward {
  id: number
  reward_no: string
  from_agent_id: number
  from_agent_name: string
  to_agent_id: number
  to_agent_name: string
  amount: number
  amount_yuan: number
  reward_type: number
  reward_type_name: string
  status: number
  status_name: string
  remark?: string
  created_at: string
  revoked_at?: string
}

// 发放奖励请求
export interface IssueChargingRewardParams {
  to_agent_id: number
  amount: number // 分
  remark?: string
}

// 充值记录查询参数
export interface ChargingDepositQueryParams {
  status?: ChargingDepositStatus
  page?: number
  page_size?: number
}

// 奖励记录查询参数
export interface ChargingRewardQueryParams {
  direction: 'from' | 'to' // from=我发放的, to=我收到的
  page?: number
  page_size?: number
}

// ========== 沉淀钱包相关类型 ==========

// 开通沉淀钱包请求
export interface EnableSettlementWalletParams {
  agent_id: number
  ratio: number // 沉淀比例(1-100)
}

// 沉淀钱包汇总
export interface SettlementWalletSummary {
  subordinate_total_balance: number
  subordinate_total_balance_yuan: number
  settlement_ratio: number
  available_amount: number
  available_amount_yuan: number
  used_amount: number
  used_amount_yuan: number
  pending_return_amount: number
  pending_return_amount_yuan: number
}

// 下级余额信息
export interface SubordinateBalance {
  agent_id: number
  agent_name: string
  available_balance: number
  available_balance_yuan: number
}

// 沉淀使用类型
export type SettlementUsageType = 1 | 2 // 1=使用 2=归还

// 沉淀使用状态
export type SettlementUsageStatus = 1 | 2 // 1=正常 2=待归还

// 沉淀使用记录
export interface SettlementWalletUsage {
  id: number
  usage_no: string
  agent_id: number
  amount: number
  amount_yuan: number
  usage_type: SettlementUsageType
  usage_type_name: string
  status: SettlementUsageStatus
  status_name: string
  return_deadline?: string
  returned_at?: string
  remark?: string
  created_at: string
}

// 使用沉淀款请求
export interface UseSettlementParams {
  amount: number // 分
  remark?: string
}

// 归还沉淀款请求
export interface ReturnSettlementParams {
  amount: number // 分
  remark?: string
}

// 沉淀使用记录查询参数
export interface SettlementUsageQueryParams {
  usage_type?: SettlementUsageType
  page?: number
  page_size?: number
}

// ========== 钱包拆分相关类型 ==========

// 代理商钱包拆分配置
export interface AgentWalletSplitConfig {
  id?: number
  agent_id: number
  split_by_channel: boolean
  configured_by?: number
  configured_at?: string
  created_at?: string
  updated_at?: string
}

// 政策模版提现门槛配置
export interface PolicyWithdrawThreshold {
  id?: number
  template_id: number
  wallet_type: number // 1=分润 2=服务费 3=奖励
  channel_id: number // 0=通用
  threshold_amount: number // 分
  created_at?: string
  updated_at?: string
}

// 子钱包（按通道拆分时的二级钱包）
export interface SubWallet {
  channel_id: number
  channel_name: string
  balance: number
  frozen_amount: number
  withdraw_threshold: number
  can_withdraw: boolean
}

// 钱包展示（支持拆分模式）
export interface WalletDisplay {
  wallet_type: number
  wallet_type_name: string
  balance: number
  frozen_amount: number
  total_income: number
  total_withdraw: number
  withdraw_threshold: number
  can_withdraw: boolean
  sub_wallets?: SubWallet[]
}

// 钱包列表响应（支持拆分模式）
export interface WalletListWithSplitResponse {
  split_by_channel: boolean
  wallets: WalletDisplay[]
}

// 设置拆分配置请求
export interface SetSplitConfigRequest {
  split_by_channel: boolean
}

// 设置提现门槛请求
export interface SetWithdrawThresholdRequest {
  wallet_type: number
  channel_id: number
  threshold_amount: number
}

// 带通道的提现请求
export interface WithdrawWithChannelRequest {
  wallet_id: number
  channel_id?: number // 拆分模式下必填
  amount: number
}

