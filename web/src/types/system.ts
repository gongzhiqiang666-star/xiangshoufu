// 系统管理相关类型定义

// 系统用户
export interface SystemUser {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  role: UserRole
  status: number
  last_login_at: string
  last_login_ip: string
  created_at: string
  updated_at: string
}

// 用户角色
export type UserRole = 'admin' | 'finance' | 'operation' | 'readonly'

// 用户角色配置
export const USER_ROLE_CONFIG: Record<UserRole, { label: string; color: string }> = {
  admin: { label: '管理员', color: '#409eff' },
  finance: { label: '财务', color: '#67c23a' },
  operation: { label: '运营', color: '#e6a23c' },
  readonly: { label: '只读用户', color: '#909399' },
}

// 创建用户请求
export interface CreateUserRequest {
  username: string
  password: string
  nickname: string
  email?: string
  phone?: string
  role: UserRole
}

// 更新用户请求
export interface UpdateUserRequest {
  nickname?: string
  email?: string
  phone?: string
  role?: UserRole
  status?: number
}

// 操作日志
export interface OperationLog {
  id: number
  user_id: number
  username: string
  nickname: string
  module: string
  action: string
  method: string
  path: string
  ip: string
  user_agent: string
  request_body: string
  response_code: number
  response_time: number // 毫秒
  created_at: string
}

// 操作模块
export type LogModule = 'auth' | 'agent' | 'merchant' | 'terminal' | 'transaction' | 'profit' | 'wallet' | 'policy' | 'system'

// 模块配置
export const LOG_MODULE_CONFIG: Record<LogModule, { label: string; color: string }> = {
  auth: { label: '认证', color: '#409eff' },
  agent: { label: '代理管理', color: '#67c23a' },
  merchant: { label: '商户管理', color: '#e6a23c' },
  terminal: { label: '终端管理', color: '#f56c6c' },
  transaction: { label: '交易管理', color: '#909399' },
  profit: { label: '分润管理', color: '#409eff' },
  wallet: { label: '钱包管理', color: '#67c23a' },
  policy: { label: '政策管理', color: '#e6a23c' },
  system: { label: '系统管理', color: '#f56c6c' },
}
