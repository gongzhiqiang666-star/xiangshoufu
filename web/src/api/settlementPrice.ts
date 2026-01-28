import { get, post, put } from './request'

// ============================================================
// 结算价类型定义
// ============================================================

/** 费率配置 */
export interface RateConfig {
  rate: string
}

/** 费率配置集合 */
export interface RateConfigs {
  [key: string]: RateConfig
}

/** 押金返现配置项 */
export interface DepositCashbackItem {
  deposit_amount: number
  cashback_amount: number
}

/** 激活奖励配置项 */
export interface ActivationRewardItem {
  reward_name: string
  min_register_days: number
  max_register_days: number
  target_amount: number
  reward_amount: number
  priority: number
}

/** 结算价 */
export interface SettlementPrice {
  id: number
  agent_id: number
  agent_name?: string
  channel_id: number
  channel_name?: string
  template_id?: number
  brand_code: string
  rate_configs: RateConfigs
  credit_rate?: string
  debit_rate?: string
  debit_cap?: string
  unionpay_rate?: string
  wechat_rate?: string
  alipay_rate?: string
  deposit_cashbacks: DepositCashbackItem[]
  sim_first_cashback: number
  sim_second_cashback: number
  sim_third_plus_cashback: number
  high_rate_configs: HighRateConfigs
  d0_extra_configs: D0ExtraConfigs
  version: number
  status: number
  effective_at: string
  created_at: string
  updated_at: string
}

/** 结算价列表项 */
export interface SettlementPriceItem {
  id: number
  agent_id: number
  agent_name: string
  channel_id: number
  channel_name: string
  brand_code: string
  rate_configs: RateConfigs
  deposit_cashbacks: DepositCashbackItem[]
  sim_first_cashback: number
  sim_second_cashback: number
  sim_third_plus_cashback: number
  high_rate_configs: HighRateConfigs
  d0_extra_configs: D0ExtraConfigs
  version: number
  status: number
  effective_at: string
  updated_at: string
}

/** 结算价列表响应 */
export interface SettlementPriceListResponse {
  list: SettlementPriceItem[]
  total: number
  page: number
  size: number
}

/** 结算价列表请求参数 */
export interface SettlementPriceListParams {
  agent_id?: number
  channel_id?: number
  status?: number
  page?: number
  page_size?: number
}

/** 创建结算价请求 */
export interface CreateSettlementPriceRequest {
  agent_id: number
  channel_id: number
  template_id?: number
  brand_code?: string
}

/** 更新费率请求 */
export interface UpdateRateRequest {
  rate_configs?: RateConfigs
  credit_rate?: string
  debit_rate?: string
  debit_cap?: string
  unionpay_rate?: string
  wechat_rate?: string
  alipay_rate?: string
}

/** 更新押金返现请求 */
export interface UpdateDepositCashbackRequest {
  deposit_cashbacks: DepositCashbackItem[]
}

/** 更新流量费返现请求 */
export interface UpdateSimCashbackRequest {
  sim_first_cashback: number
  sim_second_cashback: number
  sim_third_plus_cashback: number
}

/** 高调费率配置项 */
export interface HighRateConfig {
  rate: string  // 高调费率（%）
}

/** 高调费率配置集合 */
export interface HighRateConfigs {
  [rateType: string]: HighRateConfig
}

/** P+0加价配置项 */
export interface D0ExtraConfig {
  extra_fee: number  // 加价金额（分）
}

/** P+0加价配置集合 */
export interface D0ExtraConfigs {
  [rateType: string]: D0ExtraConfig
}

/** 更新高调费率请求 */
export interface UpdateHighRateRequest {
  high_rate_configs: HighRateConfigs
}

/** 更新P+0加价请求 */
export interface UpdateD0ExtraRequest {
  d0_extra_configs: D0ExtraConfigs
}

// ============================================================
// 代理商奖励配置类型定义
// ============================================================

/** 代理商奖励配置 */
export interface AgentRewardSetting {
  id: number
  agent_id: number
  agent_name?: string
  template_id?: number
  template_name?: string
  reward_amount: number
  activation_rewards: ActivationRewardItem[]
  version: number
  status: number
  effective_at: string
  created_at: string
  updated_at: string
}

/** 创建代理商奖励配置请求 */
export interface CreateAgentRewardSettingRequest {
  agent_id: number
  template_id?: number
  reward_amount?: number
}

/** 更新激活奖励请求 */
export interface UpdateActivationRewardRequest {
  activation_rewards: ActivationRewardItem[]
}

// ============================================================
// 调价记录类型定义
// ============================================================

/** 变更类型 */
export enum ChangeType {
  Init = 1,       // 初始化
  Rate = 2,       // 费率调整
  Deposit = 3,    // 押金返现调整
  Sim = 4,        // 流量费返现调整
  Activation = 5, // 激活奖励调整
  Batch = 6,      // 批量调整
  Sync = 7,       // 模板同步
}

/** 配置类型 */
export enum ConfigType {
  Settlement = 1, // 结算价
  Reward = 2,     // 奖励配置
}

/** 调价记录 */
export interface PriceChangeLog {
  id: number
  agent_id: number
  agent_name: string
  channel_id?: number
  channel_name: string
  change_type: ChangeType
  change_type_name: string
  config_type: ConfigType
  config_type_name: string
  field_name: string
  old_value: string
  new_value: string
  change_summary: string
  operator_name: string
  source: string
  created_at: string
}

/** 调价记录列表响应 */
export interface PriceChangeLogListResponse {
  list: PriceChangeLog[]
  total: number
  page: number
  size: number
}

/** 调价记录列表请求参数 */
export interface PriceChangeLogListParams {
  agent_id?: number
  channel_id?: number
  change_type?: ChangeType
  config_type?: ConfigType
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

// ============================================================
// 结算价API
// ============================================================

/**
 * 获取结算价列表
 */
export function getSettlementPrices(
  params: SettlementPriceListParams
): Promise<SettlementPriceListResponse> {
  return get<SettlementPriceListResponse>('/v1/settlement-prices', params)
}

/**
 * 获取结算价详情
 */
export function getSettlementPrice(id: number): Promise<SettlementPrice> {
  return get<SettlementPrice>(`/v1/settlement-prices/${id}`)
}

/**
 * 创建结算价
 */
export function createSettlementPrice(
  data: CreateSettlementPriceRequest
): Promise<SettlementPrice> {
  return post<SettlementPrice>('/v1/settlement-prices', data)
}

/**
 * 更新费率
 */
export function updateSettlementPriceRate(
  id: number,
  data: UpdateRateRequest
): Promise<SettlementPrice> {
  return put<SettlementPrice>(`/v1/settlement-prices/${id}/rate`, data)
}

/**
 * 更新押金返现
 */
export function updateSettlementPriceDeposit(
  id: number,
  data: UpdateDepositCashbackRequest
): Promise<SettlementPrice> {
  return put<SettlementPrice>(`/v1/settlement-prices/${id}/deposit`, data)
}

/**
 * 更新流量费返现
 */
export function updateSettlementPriceSim(
  id: number,
  data: UpdateSimCashbackRequest
): Promise<SettlementPrice> {
  return put<SettlementPrice>(`/v1/settlement-prices/${id}/sim`, data)
}

/**
 * 获取结算价调价记录
 */
export function getSettlementPriceChangeLogs(
  id: number,
  params?: { page?: number; page_size?: number }
): Promise<PriceChangeLogListResponse> {
  return get<PriceChangeLogListResponse>(`/v1/settlement-prices/${id}/change-logs`, params)
}

/**
 * 更新高调费率配置
 */
export function updateSettlementPriceHighRate(
  id: number,
  data: UpdateHighRateRequest
): Promise<SettlementPrice> {
  return put<SettlementPrice>(`/v1/settlement-prices/${id}/high-rate`, data)
}

/**
 * 更新P+0加价配置
 */
export function updateSettlementPriceD0Extra(
  id: number,
  data: UpdateD0ExtraRequest
): Promise<SettlementPrice> {
  return put<SettlementPrice>(`/v1/settlement-prices/${id}/d0-extra`, data)
}

// ============================================================
// 代理商奖励配置API
// ============================================================

/**
 * 获取奖励配置列表
 */
export function getRewardSettings(params?: {
  page?: number
  page_size?: number
}): Promise<{ list: AgentRewardSetting[]; total: number; page: number; size: number }> {
  return get('/v1/reward-settings', params)
}

/**
 * 获取奖励配置详情
 */
export function getRewardSetting(id: number): Promise<AgentRewardSetting> {
  return get<AgentRewardSetting>(`/v1/reward-settings/${id}`)
}

/**
 * 创建奖励配置
 */
export function createRewardSetting(
  data: CreateAgentRewardSettingRequest
): Promise<AgentRewardSetting> {
  return post<AgentRewardSetting>('/v1/reward-settings', data)
}

/**
 * 更新激活奖励
 */
export function updateRewardSettingActivation(
  id: number,
  data: UpdateActivationRewardRequest
): Promise<AgentRewardSetting> {
  return put<AgentRewardSetting>(`/v1/reward-settings/${id}/activation`, data)
}

/**
 * 获取奖励配置调价记录
 */
export function getRewardSettingChangeLogs(
  id: number,
  params?: { page?: number; page_size?: number }
): Promise<PriceChangeLogListResponse> {
  return get<PriceChangeLogListResponse>(`/v1/reward-settings/${id}/change-logs`, params)
}

// ============================================================
// 调价记录API
// ============================================================

/**
 * 获取调价记录列表
 */
export function getPriceChangeLogs(
  params: PriceChangeLogListParams
): Promise<PriceChangeLogListResponse> {
  return get<PriceChangeLogListResponse>('/v1/price-change-logs', params)
}

/**
 * 获取调价记录详情
 */
export function getPriceChangeLog(id: number): Promise<PriceChangeLog> {
  return get<PriceChangeLog>(`/v1/price-change-logs/${id}`)
}

/**
 * 按代理商获取调价记录
 */
export function getAgentPriceChangeLogs(
  agentId: number,
  params?: { page?: number; page_size?: number }
): Promise<PriceChangeLogListResponse> {
  return get<PriceChangeLogListResponse>(`/v1/agents/${agentId}/price-change-logs`, params)
}
