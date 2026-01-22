// 定时任务相关类型定义

// 任务执行状态
export enum JobStatus {
  Success = 1,
  Failed = 2,
  Running = 3
}

// 任务触发类型
export enum JobTriggerType {
  Auto = 1,
  Manual = 2
}

// 告警通道类型
export enum AlertChannelType {
  DingTalk = 1,
  WeChatWork = 2,
  Email = 3
}

// 告警类型
export enum AlertType {
  JobFailed = 1,
  ConsecutiveFail = 2,
  JobTimeout = 3
}

// 告警发送状态
export enum AlertSendStatus {
  Pending = 0,
  Sent = 1,
  Failed = 2
}

// 状态名称映射
export const JobStatusName: Record<number, string> = {
  1: '成功',
  2: '失败',
  3: '运行中'
}

export const JobTriggerTypeName: Record<number, string> = {
  1: '自动触发',
  2: '手动触发'
}

export const AlertChannelTypeName: Record<number, string> = {
  1: '钉钉',
  2: '企业微信',
  3: '邮件'
}

export const AlertTypeName: Record<number, string> = {
  1: '任务失败',
  2: '连续失败',
  3: '任务超时'
}

export const AlertSendStatusName: Record<number, string> = {
  0: '待发送',
  1: '已发送',
  2: '发送失败'
}

// 任务配置
export interface JobConfig {
  id: number
  job_name: string
  job_desc: string
  cron_expr: string
  interval_seconds: number
  is_enabled: boolean
  max_retries: number
  retry_interval: number
  alert_threshold: number
  timeout_seconds: number
  created_at: string
  updated_at: string
}

// 任务列表项
export interface JobListItem {
  id: number
  job_name: string
  job_desc: string
  interval_seconds: number
  is_enabled: boolean
  max_retries: number
  alert_threshold: number
  is_running: boolean
  last_run_at: string | null
  last_status: number | null
  updated_at: string
}

// 任务详情
export interface JobDetail {
  config: JobConfig
  is_running: boolean
  latest_logs: JobExecutionLog[]
}

// 任务执行日志
export interface JobExecutionLog {
  id: number
  job_name: string
  started_at: string
  ended_at: string | null
  duration_ms: number
  status: number
  processed_count: number
  success_count: number
  fail_count: number
  error_message: string
  error_stack: string
  retry_count: number
  trigger_type: number
  created_at: string
}

// 任务执行统计
export interface JobExecutionStats {
  job_name: string
  total_count: number
  success_count: number
  fail_count: number
  avg_duration_ms: number
  max_duration_ms: number
  min_duration_ms: number
}

// 告警配置
export interface AlertConfig {
  id: number
  name: string
  channel_type: number
  channel_type_name?: string
  webhook_url: string
  webhook_secret?: string
  email_addresses: string
  email_smtp_host?: string
  email_smtp_port?: number
  email_username?: string
  is_enabled: boolean
  created_by: number | null
  created_at: string
  updated_at: string
}

// 告警记录
export interface AlertLog {
  id: number
  job_name: string
  alert_type: number
  channel_type: number
  config_id: number | null
  title: string
  message: string
  send_status: number
  send_at: string | null
  error_message: string
  created_at: string
}

// 更新任务配置请求
export interface UpdateJobConfigRequest {
  job_desc?: string
  interval_seconds?: number
  max_retries?: number
  retry_interval?: number
  alert_threshold?: number
  timeout_seconds?: number
}

// 创建告警配置请求
export interface CreateAlertConfigRequest {
  name: string
  channel_type: number
  webhook_url?: string
  webhook_secret?: string
  email_addresses?: string
  email_smtp_host?: string
  email_smtp_port?: number
  email_username?: string
  email_password?: string
}

// 更新告警配置请求
export interface UpdateAlertConfigRequest {
  name?: string
  webhook_url?: string
  webhook_secret?: string
  email_addresses?: string
  email_smtp_host?: string
  email_smtp_port?: number
  email_username?: string
  email_password?: string
}
