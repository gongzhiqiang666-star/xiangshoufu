// 用户信息
export interface UserInfo {
  id: number
  username: string
  role: UserRole
  real_name: string
  phone: string
  email: string
  status: number
  last_login_at: string
  created_at: string
}

// 用户角色
export type UserRole = 'admin' | 'finance' | 'operation' | 'readonly'

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录用户信息 (从后端返回)
export interface LoginUserInfo {
  id: number
  username: string
  role_type: number // 1=代理商, 2=管理员
}

// 登录代理商信息
export interface LoginAgentInfo {
  id: number
  agent_no: string
  agent_name: string
  level: number
}

// 登录响应
export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number // 秒数
  token_type: string
  user: LoginUserInfo
  agent?: LoginAgentInfo
}

// Token刷新响应
export interface RefreshTokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
  user: LoginUserInfo
  agent?: LoginAgentInfo
}
