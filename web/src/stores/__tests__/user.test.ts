/**
 * User Store 测试
 * 覆盖: 初始化状态、登录流程、登出流程、权限检查
 */

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useUserStore } from '../user'

// Mock API 模块
vi.mock('@/api/auth', () => ({
  login: vi.fn(),
  logout: vi.fn(),
  getProfile: vi.fn(),
}))

// Mock storage 模块
vi.mock('@/utils/storage', () => ({
  getToken: vi.fn(() => null),
  setToken: vi.fn(),
  getUser: vi.fn(() => null),
  setUser: vi.fn(),
  setExpires: vi.fn(),
  clearAuth: vi.fn(),
}))

// 导入 mocked 模块
import { login as loginApi, logout as logoutApi, getProfile } from '@/api/auth'
import { getToken, setToken, getUser, setUser, setExpires, clearAuth } from '@/utils/storage'

describe('useUserStore', () => {
  beforeEach(() => {
    // 每个测试前创建新的 Pinia 实例
    setActivePinia(createPinia())
    // 重置所有 mock
    vi.clearAllMocks()
  })

  describe('初始化状态', () => {
    it('should initialize with null token and userInfo', () => {
      const store = useUserStore()
      expect(store.token).toBeNull()
      expect(store.userInfo).toBeNull()
      expect(store.loading).toBe(false)
    })

    it('should have isLoggedIn as false initially', () => {
      const store = useUserStore()
      expect(store.isLoggedIn).toBe(false)
    })

    it('should have empty username and role initially', () => {
      const store = useUserStore()
      expect(store.username).toBe('')
      expect(store.role).toBe('')
      expect(store.realName).toBe('')
    })
  })

  describe('登录流程', () => {
    it('should login successfully and update state', async () => {
      const mockResponse = {
        access_token: 'test-token-123',
        expires_in: 86400,
        user: {
          id: 1,
          username: 'testuser',
          role_type: 2,
        },
        agent: {
          agent_name: '测试代理商',
        },
      }

      vi.mocked(loginApi).mockResolvedValue(mockResponse)

      const store = useUserStore()
      const result = await store.login('testuser', 'password123')

      // 验证 API 调用
      expect(loginApi).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password123',
      })

      // 验证状态更新
      expect(store.token).toBe('test-token-123')
      expect(store.userInfo).not.toBeNull()
      expect(store.userInfo?.username).toBe('testuser')
      expect(store.userInfo?.role).toBe('admin')
      expect(store.isLoggedIn).toBe(true)

      // 验证存储调用
      expect(setToken).toHaveBeenCalledWith('test-token-123')
      expect(setUser).toHaveBeenCalled()
      expect(setExpires).toHaveBeenCalled()

      // 验证返回值
      expect(result).toEqual(mockResponse)
    })

    it('should set loading state during login', async () => {
      let resolveLogin: (value: unknown) => void
      const loginPromise = new Promise((resolve) => {
        resolveLogin = resolve
      })

      vi.mocked(loginApi).mockReturnValue(loginPromise as Promise<never>)

      const store = useUserStore()
      const loginTask = store.login('user', 'pass')

      // 登录过程中 loading 应为 true
      expect(store.loading).toBe(true)

      // 完成登录
      resolveLogin!({
        access_token: 'token',
        expires_in: 86400,
        user: { id: 1, username: 'user', role_type: 1 },
        agent: { agent_name: 'Agent' },
      })

      await loginTask

      // 登录完成后 loading 应为 false
      expect(store.loading).toBe(false)
    })

    it('should handle login failure', async () => {
      vi.mocked(loginApi).mockRejectedValue(new Error('Invalid credentials'))

      const store = useUserStore()

      await expect(store.login('user', 'wrong')).rejects.toThrow('Invalid credentials')
      expect(store.loading).toBe(false)
      expect(store.token).toBeNull()
    })

    it('should map role_type correctly', async () => {
      // role_type 2 -> admin
      vi.mocked(loginApi).mockResolvedValue({
        access_token: 'token',
        expires_in: 86400,
        user: { id: 1, username: 'admin', role_type: 2 },
        agent: { agent_name: 'Admin' },
      })

      const store = useUserStore()
      await store.login('admin', 'pass')
      expect(store.role).toBe('admin')

      // 重置并测试 role_type 1 -> readonly
      setActivePinia(createPinia())
      const store2 = useUserStore()

      vi.mocked(loginApi).mockResolvedValue({
        access_token: 'token',
        expires_in: 86400,
        user: { id: 2, username: 'agent', role_type: 1 },
        agent: { agent_name: 'Agent' },
      })

      await store2.login('agent', 'pass')
      expect(store2.role).toBe('readonly')
    })
  })

  describe('登出流程', () => {
    it('should logout and clear state', async () => {
      vi.mocked(logoutApi).mockResolvedValue(undefined)

      const store = useUserStore()
      // 模拟已登录状态
      store.token = 'existing-token'
      store.userInfo = {
        id: 1,
        username: 'user',
        role: 'admin',
        real_name: '用户',
        phone: '',
        email: '',
        status: 1,
        last_login_at: '',
        created_at: '',
      }

      await store.logout()

      expect(store.token).toBeNull()
      expect(store.userInfo).toBeNull()
      expect(clearAuth).toHaveBeenCalled()
    })

    it('should clear state even if logout API fails', async () => {
      vi.mocked(logoutApi).mockRejectedValue(new Error('Network error'))

      const store = useUserStore()
      store.token = 'token'
      store.userInfo = {
        id: 1,
        username: 'user',
        role: 'admin',
        real_name: '用户',
        phone: '',
        email: '',
        status: 1,
        last_login_at: '',
        created_at: '',
      }

      await store.logout()

      // 即使 API 失败，本地状态也应清除
      expect(store.token).toBeNull()
      expect(store.userInfo).toBeNull()
      expect(clearAuth).toHaveBeenCalled()
    })
  })

  describe('获取用户信息', () => {
    it('should fetch user info when token exists', async () => {
      const mockUser = {
        id: 1,
        username: 'testuser',
        role: 'admin',
        real_name: '测试用户',
        phone: '13800138000',
        email: 'test@example.com',
        status: 1,
        last_login_at: '2024-01-01T00:00:00Z',
        created_at: '2024-01-01T00:00:00Z',
      }

      vi.mocked(getProfile).mockResolvedValue(mockUser)

      const store = useUserStore()
      store.token = 'valid-token'

      const result = await store.fetchUserInfo()

      expect(getProfile).toHaveBeenCalled()
      expect(store.userInfo).toEqual(mockUser)
      expect(setUser).toHaveBeenCalledWith(mockUser)
      expect(result).toEqual(mockUser)
    })

    it('should return null when no token', async () => {
      const store = useUserStore()
      store.token = null

      const result = await store.fetchUserInfo()

      expect(getProfile).not.toHaveBeenCalled()
      expect(result).toBeNull()
    })

    it('should clear auth when fetch fails', async () => {
      vi.mocked(getProfile).mockRejectedValue(new Error('Unauthorized'))

      const store = useUserStore()
      store.token = 'invalid-token'

      const result = await store.fetchUserInfo()

      expect(store.token).toBeNull()
      expect(store.userInfo).toBeNull()
      expect(clearAuth).toHaveBeenCalled()
      expect(result).toBeNull()
    })
  })

  describe('权限检查', () => {
    it('should return false when not logged in', () => {
      const store = useUserStore()
      expect(store.hasRole('admin')).toBe(false)
    })

    it('should check single role correctly', () => {
      const store = useUserStore()
      store.userInfo = {
        id: 1,
        username: 'admin',
        role: 'admin',
        real_name: '管理员',
        phone: '',
        email: '',
        status: 1,
        last_login_at: '',
        created_at: '',
      }

      expect(store.hasRole('admin')).toBe(true)
      expect(store.hasRole('readonly')).toBe(false)
    })

    it('should check multiple roles correctly', () => {
      const store = useUserStore()
      store.userInfo = {
        id: 1,
        username: 'user',
        role: 'readonly',
        real_name: '用户',
        phone: '',
        email: '',
        status: 1,
        last_login_at: '',
        created_at: '',
      }

      expect(store.hasRole(['admin', 'readonly'])).toBe(true)
      expect(store.hasRole(['admin', 'superadmin'])).toBe(false)
    })
  })
})
