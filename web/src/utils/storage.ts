const TOKEN_KEY = 'xsf_token'
const USER_KEY = 'xsf_user'
const EXPIRES_KEY = 'xsf_expires'

/**
 * 获取Token
 */
export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

/**
 * 设置Token
 */
export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

/**
 * 移除Token
 */
export function removeToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

/**
 * 获取用户信息
 */
export function getUser<T>(): T | null {
  const user = localStorage.getItem(USER_KEY)
  if (user) {
    try {
      return JSON.parse(user) as T
    } catch {
      return null
    }
  }
  return null
}

/**
 * 设置用户信息
 */
export function setUser<T>(user: T): void {
  localStorage.setItem(USER_KEY, JSON.stringify(user))
}

/**
 * 移除用户信息
 */
export function removeUser(): void {
  localStorage.removeItem(USER_KEY)
}

/**
 * 获取Token过期时间
 */
export function getExpires(): string | null {
  return localStorage.getItem(EXPIRES_KEY)
}

/**
 * 设置Token过期时间
 */
export function setExpires(expires: string): void {
  localStorage.setItem(EXPIRES_KEY, expires)
}

/**
 * 移除Token过期时间
 */
export function removeExpires(): void {
  localStorage.removeItem(EXPIRES_KEY)
}

/**
 * 清除所有认证信息
 */
export function clearAuth(): void {
  removeToken()
  removeUser()
  removeExpires()
}

/**
 * 检查Token是否过期
 */
export function isTokenExpired(): boolean {
  const expires = getExpires()
  if (!expires) return true
  return new Date(expires).getTime() < Date.now()
}
