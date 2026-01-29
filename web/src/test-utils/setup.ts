/**
 * 全局测试设置
 * 在 vitest.config.ts 中通过 setupFiles 引入
 */

import { config } from '@vue/test-utils'
import { vi } from 'vitest'

// Mock Element Plus 组件
config.global.stubs = {
  'el-icon': true,
  'el-button': true,
  'el-input': true,
  'el-form': true,
  'el-form-item': true,
  'el-table': true,
  'el-table-column': true,
  'el-pagination': true,
  'el-dialog': true,
  'el-select': true,
  'el-option': true,
  'el-date-picker': true,
  'el-message': true,
  'el-message-box': true,
  // Element Plus Icons
  'CaretTop': true,
  'CaretBottom': true,
}

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
}
Object.defineProperty(window, 'localStorage', { value: localStorageMock })

// Mock console.error 避免测试输出过多噪音
vi.spyOn(console, 'error').mockImplementation(() => {})
