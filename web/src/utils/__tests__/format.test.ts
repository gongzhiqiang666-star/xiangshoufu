/**
 * 格式化工具函数测试
 * 覆盖: 正常流程、边界情况、错误处理、特殊输入
 */

import { describe, it, expect } from 'vitest'
import {
  formatAmount,
  formatNumber,
  formatPercent,
  formatDateTime,
  formatDate,
  calculateTrend,
} from '../format'

describe('formatAmount', () => {
  // ✅ 正常流程 (Happy Path)
  it('should format 100 cents to "1.00"', () => {
    expect(formatAmount(100)).toBe('1.00')
  })

  it('should format 12345 cents to "123.45"', () => {
    expect(formatAmount(12345)).toBe('123.45')
  })

  // ✅ 边界情况 (Edge Cases)
  it('should handle zero', () => {
    expect(formatAmount(0)).toBe('0.00')
  })

  it('should handle small values (1 cent)', () => {
    expect(formatAmount(1)).toBe('0.01')
  })

  // ✅ 错误处理 (Error Handling)
  it('should handle negative values', () => {
    expect(formatAmount(-100)).toBe('-1.00')
  })

  // ✅ 特殊输入 (Special Inputs)
  it('should format large numbers with thousand separators', () => {
    expect(formatAmount(10000000)).toBe('100,000.00')
  })

  it('should respect custom decimals parameter', () => {
    expect(formatAmount(12345, 0)).toBe('123')
    expect(formatAmount(12345, 1)).toBe('123.5')
    expect(formatAmount(12345, 3)).toBe('123.450')
  })

  it('should handle decimal precision correctly', () => {
    expect(formatAmount(199)).toBe('1.99')
    expect(formatAmount(1)).toBe('0.01')
  })
})

describe('formatNumber', () => {
  // ✅ 正常流程
  it('should format regular numbers with locale', () => {
    expect(formatNumber(1234)).toBe('1,234')
  })

  // ✅ 边界情况
  it('should handle zero', () => {
    expect(formatNumber(0)).toBe('0')
  })

  it('should format numbers >= 10000 with 万 unit', () => {
    expect(formatNumber(10000)).toBe('1.00万')
    expect(formatNumber(12345)).toBe('1.23万')
    expect(formatNumber(99999)).toBe('10.00万')
  })

  it('should format numbers >= 100000000 with 亿 unit', () => {
    expect(formatNumber(100000000)).toBe('1.00亿')
    expect(formatNumber(123456789)).toBe('1.23亿')
  })

  // ✅ 特殊输入
  it('should handle boundary between units', () => {
    expect(formatNumber(9999)).toBe('9,999')
    expect(formatNumber(10000)).toBe('1.00万')
    expect(formatNumber(99999999)).toBe('10000.00万')
    expect(formatNumber(100000000)).toBe('1.00亿')
  })
})

describe('formatPercent', () => {
  // ✅ 正常流程
  it('should format 0.5 to "50.00%"', () => {
    expect(formatPercent(0.5)).toBe('50.00%')
  })

  it('should format 0.1234 to "12.34%"', () => {
    expect(formatPercent(0.1234)).toBe('12.34%')
  })

  // ✅ 边界情况
  it('should handle zero', () => {
    expect(formatPercent(0)).toBe('0.00%')
  })

  it('should handle 100%', () => {
    expect(formatPercent(1)).toBe('100.00%')
  })

  // ✅ 特殊输入
  it('should handle values > 100%', () => {
    expect(formatPercent(1.5)).toBe('150.00%')
  })

  it('should respect custom decimals', () => {
    expect(formatPercent(0.12345, 0)).toBe('12%')
    expect(formatPercent(0.12345, 1)).toBe('12.3%')
    expect(formatPercent(0.12345, 3)).toBe('12.345%')
  })
})

describe('formatDateTime', () => {
  const testDate = new Date('2024-03-15T10:30:45')

  // ✅ 正常流程
  it('should format date to datetime by default', () => {
    expect(formatDateTime(testDate)).toBe('2024-03-15 10:30:45')
  })

  it('should format date string input', () => {
    expect(formatDateTime('2024-03-15T10:30:45')).toBe('2024-03-15 10:30:45')
  })

  // ✅ 不同格式
  it('should format to date only', () => {
    expect(formatDateTime(testDate, 'date')).toBe('2024-03-15')
  })

  it('should format to time only', () => {
    expect(formatDateTime(testDate, 'time')).toBe('10:30:45')
  })

  // ✅ 边界情况
  it('should pad single digit months and days', () => {
    const date = new Date('2024-01-05T09:05:05')
    expect(formatDateTime(date, 'date')).toBe('2024-01-05')
    expect(formatDateTime(date, 'time')).toBe('09:05:05')
  })
})

describe('formatDate', () => {
  // ✅ 正常流程
  it('should format valid date', () => {
    expect(formatDate('2024-03-15T10:30:45')).toBe('2024-03-15 10:30:45')
  })

  // ✅ 空值处理
  it('should return "-" for null', () => {
    expect(formatDate(null)).toBe('-')
  })

  it('should return "-" for undefined', () => {
    expect(formatDate(undefined)).toBe('-')
  })

  it('should return "-" for empty string', () => {
    expect(formatDate('')).toBe('-')
  })
})

describe('calculateTrend', () => {
  // ✅ 正常流程
  it('should calculate positive trend', () => {
    expect(calculateTrend(150, 100)).toBe(50)
  })

  it('should calculate negative trend', () => {
    expect(calculateTrend(50, 100)).toBe(-50)
  })

  // ✅ 边界情况
  it('should return 0 when both values are 0', () => {
    expect(calculateTrend(0, 0)).toBe(0)
  })

  it('should return 100 when previous is 0 and current is positive', () => {
    expect(calculateTrend(100, 0)).toBe(100)
  })

  it('should return 0 when previous is 0 and current is 0', () => {
    expect(calculateTrend(0, 0)).toBe(0)
  })

  // ✅ 特殊输入
  it('should handle same values (0% change)', () => {
    expect(calculateTrend(100, 100)).toBe(0)
  })

  it('should handle decimal results', () => {
    expect(calculateTrend(110, 100)).toBe(10)
    expect(calculateTrend(133, 100)).toBe(33)
  })
})
