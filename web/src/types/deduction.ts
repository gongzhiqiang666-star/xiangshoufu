/**
 * 代扣管理类型定义
 */

// 代扣计划状态
export type DeductionPlanStatus = 1 | 2 | 3 | 4
export const DEDUCTION_PLAN_STATUS = {
  ACTIVE: 1 as const,      // 进行中
  COMPLETED: 2 as const,   // 已完成
  PAUSED: 3 as const,      // 已暂停
  CANCELLED: 4 as const,   // 已取消
}

// 代扣计划类型
export type DeductionPlanType = 1 | 2 | 3
export const DEDUCTION_PLAN_TYPE = {
  GOODS: 1 as const,       // 货款代扣
  PARTNER: 2 as const,     // 伙伴代扣
  DEPOSIT: 3 as const,     // 押金代扣
}

// 代扣记录状态
export type DeductionRecordStatus = 0 | 1 | 2 | 3
export const DEDUCTION_RECORD_STATUS = {
  PENDING: 0 as const,         // 待扣款
  SUCCESS: 1 as const,         // 成功
  PARTIAL_SUCCESS: 2 as const, // 部分成功
  FAILED: 3 as const,          // 失败
}

// 货款代扣状态
export type GoodsDeductionStatus = 1 | 2 | 3 | 4
export const GOODS_DEDUCTION_STATUS = {
  PENDING_ACCEPT: 1 as const,  // 待接收
  IN_PROGRESS: 2 as const,     // 进行中
  COMPLETED: 3 as const,       // 已完成
  REJECTED: 4 as const,        // 已拒绝
}

// 扣款来源
export type DeductionSource = 1 | 2 | 3
export const DEDUCTION_SOURCE = {
  PROFIT: 1 as const,          // 分润钱包
  SERVICE_FEE: 2 as const,     // 服务费钱包
  BOTH: 3 as const,            // 两者都扣
}

// 状态配置
export type TagType = 'primary' | 'success' | 'warning' | 'info' | 'danger'

export const DEDUCTION_PLAN_STATUS_CONFIG: Record<DeductionPlanStatus, { label: string; color: string; type: TagType }> = {
  1: { label: '进行中', color: '#409eff', type: 'primary' },
  2: { label: '已完成', color: '#67c23a', type: 'success' },
  3: { label: '已暂停', color: '#e6a23c', type: 'warning' },
  4: { label: '已取消', color: '#909399', type: 'info' },
}

export const DEDUCTION_PLAN_TYPE_CONFIG: Record<DeductionPlanType, { label: string; color: string }> = {
  1: { label: '货款代扣', color: '#409eff' },
  2: { label: '伙伴代扣', color: '#67c23a' },
  3: { label: '押金代扣', color: '#e6a23c' },
}

export const DEDUCTION_RECORD_STATUS_CONFIG: Record<DeductionRecordStatus, { label: string; color: string; type: TagType }> = {
  0: { label: '待扣款', color: '#909399', type: 'info' },
  1: { label: '成功', color: '#67c23a', type: 'success' },
  2: { label: '部分成功', color: '#e6a23c', type: 'warning' },
  3: { label: '失败', color: '#f56c6c', type: 'danger' },
}

export const GOODS_DEDUCTION_STATUS_CONFIG: Record<GoodsDeductionStatus, { label: string; color: string; type: TagType }> = {
  1: { label: '待接收', color: '#e6a23c', type: 'warning' },
  2: { label: '进行中', color: '#409eff', type: 'primary' },
  3: { label: '已完成', color: '#67c23a', type: 'success' },
  4: { label: '已拒绝', color: '#f56c6c', type: 'danger' },
}

export const DEDUCTION_SOURCE_CONFIG: Record<DeductionSource, { label: string; color: string }> = {
  1: { label: '分润钱包', color: '#409eff' },
  2: { label: '服务费钱包', color: '#67c23a' },
  3: { label: '分润+服务费', color: '#e6a23c' },
}

// ==================== 代扣计划接口 ====================

// 代扣计划
export interface DeductionPlan {
  id: number
  plan_no: string
  deductor_id: number
  deductor_name: string
  deductee_id: number
  deductee_name: string
  plan_type: DeductionPlanType
  total_amount: number        // 总金额（分）
  deducted_amount: number     // 已扣金额（分）
  remaining_amount: number    // 剩余金额（分）
  total_periods: number       // 总期数
  current_period: number      // 当前期数
  period_amount: number       // 每期金额（分）
  status: DeductionPlanStatus
  related_type?: string
  related_id?: number
  remark?: string
  created_by: number
  created_at: string
  updated_at: string
  completed_at?: string
}

// 代扣计划详情（包含记录列表）
export interface DeductionPlanDetail extends DeductionPlan {
  records?: DeductionRecord[]
  deductor_info?: AgentBasicInfo
  deductee_info?: AgentBasicInfo
}

// 代扣记录
export interface DeductionRecord {
  id: number
  plan_id: number
  plan_no: string
  deductor_id: number
  deductee_id: number
  period_num: number
  amount: number              // 应扣金额（分）
  actual_amount: number       // 实扣金额（分）
  status: DeductionRecordStatus
  wallet_details?: WalletDeductDetail[]
  fail_reason?: string
  scheduled_at: string
  deducted_at?: string
  created_at: string
}

// 钱包扣款明细
export interface WalletDeductDetail {
  wallet_id: number
  wallet_type: number
  wallet_name: string
  balance_before: number
  deduct_amount: number
  balance_after: number
}

// 代理商基本信息
export interface AgentBasicInfo {
  id: number
  agent_no: string
  agent_name: string
  phone?: string
}

// 创建代扣计划请求
export interface CreateDeductionPlanRequest {
  deductee_id: number           // 被扣款方代理商ID
  plan_type: DeductionPlanType  // 计划类型
  total_amount: number          // 总金额（分）
  total_periods: number         // 总期数
  remark?: string               // 备注
}

// 代扣计划查询参数
export interface DeductionPlanQueryParams {
  plan_type?: DeductionPlanType
  status?: DeductionPlanStatus
  deductee_id?: number
  start_date?: string
  end_date?: string
  page?: number
  page_size?: number
}

// 代扣计划统计
export interface DeductionPlanStats {
  total_count: number
  active_count: number
  completed_count: number
  paused_count: number
  total_amount: number
  deducted_amount: number
  remaining_amount: number
}

// ==================== 货款代扣接口 ====================

// 货款代扣
export interface GoodsDeduction {
  id: number
  deduction_no: string
  from_agent_id: number
  from_agent_name: string
  to_agent_id: number
  to_agent_name: string
  total_amount: number          // 总金额（分）
  total_amount_yuan: number     // 总金额（元）
  deducted_amount: number       // 已扣金额（分）
  remaining_amount: number      // 剩余金额（分）
  deduction_source: DeductionSource
  source_name: string
  terminal_count: number
  unit_price: number            // 单价（分）
  status: GoodsDeductionStatus
  status_name: string
  progress: number              // 进度百分比
  agreement_signed: boolean
  agreement_url?: string
  distribute_id?: number
  remark?: string
  created_by: number
  created_at: string
  accepted_at?: string
  completed_at?: string
}

// 货款代扣详情
export interface GoodsDeductionDetail extends GoodsDeduction {
  terminals?: GoodsDeductionTerminal[]
  details?: GoodsDeductionDetailRecord[]
  from_agent_info?: AgentBasicInfo
  to_agent_info?: AgentBasicInfo
}

// 货款代扣终端
export interface GoodsDeductionTerminal {
  id: number
  deduction_id: number
  terminal_id: number
  terminal_sn: string
  unit_price: number
  created_at: string
}

// 货款代扣扣款明细
export interface GoodsDeductionDetailRecord {
  id: number
  deduction_id: number
  deduction_no: string
  amount: number                // 本次扣款金额（分）
  wallet_type: number
  wallet_type_name: string
  channel_id?: number
  channel_name?: string
  wallet_balance_before: number
  wallet_balance_after: number
  cumulative_deducted: number
  remaining_after: number
  trigger_type: string
  trigger_transaction_id?: number
  trigger_profit_id?: number
  created_at: string
}

// 货款代扣查询参数
export interface GoodsDeductionQueryParams {
  status?: GoodsDeductionStatus
  deduction_source?: DeductionSource
  start_date?: string
  end_date?: string
  keyword?: string
  page?: number
  page_size?: number
}

// 货款代扣统计
export interface GoodsDeductionSummary {
  total_count: number
  pending_count: number
  in_progress_count: number
  completed_count: number
  total_amount: number
  deducted_amount: number
  remaining_amount: number
}

// 创建货款代扣请求
export interface CreateGoodsDeductionRequest {
  to_agent_id: number
  unit_price: number
  deduction_source: DeductionSource
  terminals: CreateGoodsDeductionTerminal[]
  agreement_url?: string
  remark?: string
  distribute_id?: number
}

// 创建货款代扣终端
export interface CreateGoodsDeductionTerminal {
  terminal_id: number
  terminal_sn?: string
  unit_price?: number
}

// 拒绝货款代扣请求
export interface RejectGoodsDeductionRequest {
  reason: string
}
