/**
 * 测试工具函数集合
 */

import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import type { Component } from 'vue'

/**
 * 创建带 Pinia 的组件挂载器
 * @param component Vue组件
 * @param options 挂载选项
 */
export function mountWithPinia<T extends Component>(
  component: T,
  options: Record<string, unknown> = {}
) {
  const pinia = createPinia()
  setActivePinia(pinia)

  const globalOptions = options.global as Record<string, unknown> | undefined

  return mount(component, {
    global: {
      plugins: [pinia],
      stubs: {
        'el-icon': true,
        'CaretTop': true,
        'CaretBottom': true,
      },
      ...globalOptions,
    },
    ...options,
  })
}

/**
 * 创建 Pinia 实例用于 Store 测试
 */
export function createTestPinia() {
  const pinia = createPinia()
  setActivePinia(pinia)
  return pinia
}

/**
 * 模拟 API 响应
 * @param data 响应数据
 * @param delay 延迟时间（毫秒）
 */
export function mockApiResponse<T>(data: T, delay = 0): Promise<T> {
  return new Promise((resolve) => {
    setTimeout(() => resolve(data), delay)
  })
}

/**
 * 模拟 API 错误
 * @param message 错误消息
 * @param code 错误码
 */
export function mockApiError(message: string, code = 500): Promise<never> {
  return Promise.reject({
    response: {
      status: code,
      data: { message },
    },
  })
}

/**
 * 等待组件更新
 * @param wrapper Vue组件包装器
 * @param times 更新次数
 */
export async function flushPromises(times = 1): Promise<void> {
  for (let i = 0; i < times; i++) {
    await new Promise((resolve) => setTimeout(resolve, 0))
  }
}

/**
 * 创建模拟用户信息
 */
export function createMockUser(overrides: Record<string, unknown> = {}) {
  return {
    id: 1,
    username: 'testuser',
    role: 'admin',
    real_name: '测试用户',
    phone: '13800138000',
    email: 'test@example.com',
    status: 1,
    last_login_at: '2024-01-01T00:00:00Z',
    created_at: '2024-01-01T00:00:00Z',
    ...overrides,
  }
}

/**
 * 创建模拟登录响应
 */
export function createMockLoginResponse(overrides: Record<string, unknown> = {}) {
  return {
    access_token: 'mock-token-12345',
    expires_in: 86400,
    user: {
      id: 1,
      username: 'testuser',
      role_type: 2,
    },
    agent: {
      agent_name: '测试代理商',
    },
    ...overrides,
  }
}
