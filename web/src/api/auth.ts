import { post, get } from './request'
import type { LoginRequest, LoginResponse, RefreshTokenResponse, UserInfo } from '@/types'
import { encryptPassword } from '@/utils/crypto'

/**
 * 用户登录（密码加密传输）
 */
export async function login(data: LoginRequest): Promise<LoginResponse> {
  // 加密密码
  const encryptedPassword = await encryptPassword(data.password)

  return post<LoginResponse>('/v1/auth/login', {
    username: data.username,
    password: encryptedPassword,
    encrypted: true, // 标识密码已加密
  })
}

/**
 * 用户登出
 */
export function logout(): Promise<void> {
  return post<void>('/v1/auth/logout')
}

/**
 * 刷新Token
 */
export function refreshToken(): Promise<RefreshTokenResponse> {
  return post<RefreshTokenResponse>('/v1/auth/refresh')
}

/**
 * 获取当前用户信息
 */
export function getProfile(): Promise<UserInfo> {
  return get<UserInfo>('/v1/auth/profile')
}

/**
 * 修改密码
 */
export function changePassword(data: {
  old_password: string
  new_password: string
}): Promise<void> {
  return post<void>('/v1/auth/change-password', data)
}
