import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, UserRole } from '@/types'
import { login as loginApi, logout as logoutApi, getProfile } from '@/api/auth'
import { getToken, setToken, setUser, setExpires, clearAuth, getUser } from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string | null>(getToken())
  const userInfo = ref<UserInfo | null>(getUser<UserInfo>())
  const loading = ref(false)

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const role = computed(() => userInfo.value?.role || '')
  const realName = computed(() => userInfo.value?.real_name || '')

  // 登录
  async function login(username: string, password: string) {
    loading.value = true
    try {
      const res = await loginApi({ username, password })
      token.value = res.access_token

      // 转换用户信息格式
      const user: UserInfo = {
        id: res.user.id,
        username: res.user.username,
        role: getRoleName(res.user.role_type),
        real_name: res.agent?.agent_name || res.user.username,
        phone: '',
        email: '',
        status: 1,
        last_login_at: new Date().toISOString(),
        created_at: '',
      }
      userInfo.value = user

      // 计算过期时间
      const expiresAt = new Date(Date.now() + res.expires_in * 1000).toISOString()

      // 持久化存储
      setToken(res.access_token)
      setUser(user)
      setExpires(expiresAt)

      return res
    } finally {
      loading.value = false
    }
  }

  // 根据 role_type 获取角色名称
  function getRoleName(roleType: number): UserRole {
    switch (roleType) {
      case 2:
        return 'admin'
      case 1:
        return 'readonly' // 代理商默认只读
      default:
        return 'readonly'
    }
  }

  // 登出
  async function logout() {
    try {
      await logoutApi()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // 无论API调用是否成功，都清除本地状态
      token.value = null
      userInfo.value = null
      clearAuth()
    }
  }

  // 获取用户信息
  async function fetchUserInfo() {
    if (!token.value) return null

    try {
      const user = await getProfile()
      userInfo.value = user
      setUser(user)
      return user
    } catch (error) {
      console.error('Fetch user info error:', error)
      // Token无效时清除登录状态
      token.value = null
      userInfo.value = null
      clearAuth()
      return null
    }
  }

  // 检查是否有权限
  function hasRole(roles: string | string[]): boolean {
    if (!userInfo.value) return false
    const roleList = Array.isArray(roles) ? roles : [roles]
    return roleList.includes(userInfo.value.role)
  }

  return {
    // 状态
    token,
    userInfo,
    loading,
    // 计算属性
    isLoggedIn,
    username,
    role,
    realName,
    // 方法
    login,
    logout,
    fetchUserInfo,
    hasRole,
  }
})
