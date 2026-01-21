/**
 * 格式化金额（分转元）
 * @param amount 金额（分）
 * @param decimals 小数位数
 */
export function formatAmount(amount: number, decimals = 2): string {
  const yuan = amount / 100
  return yuan.toLocaleString('zh-CN', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

/**
 * 格式化大数字（带单位）
 * @param num 数字
 */
export function formatNumber(num: number): string {
  if (num >= 100000000) {
    return (num / 100000000).toFixed(2) + '亿'
  }
  if (num >= 10000) {
    return (num / 10000).toFixed(2) + '万'
  }
  return num.toLocaleString('zh-CN')
}

/**
 * 格式化百分比
 * @param value 值（0-1之间）
 * @param decimals 小数位数
 */
export function formatPercent(value: number, decimals = 2): string {
  return (value * 100).toFixed(decimals) + '%'
}

/**
 * 格式化日期时间
 * @param date 日期字符串或Date对象
 * @param format 格式类型
 */
export function formatDateTime(
  date: string | Date,
  format: 'date' | 'datetime' | 'time' = 'datetime'
): string {
  const d = typeof date === 'string' ? new Date(date) : date

  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  const hours = String(d.getHours()).padStart(2, '0')
  const minutes = String(d.getMinutes()).padStart(2, '0')
  const seconds = String(d.getSeconds()).padStart(2, '0')

  switch (format) {
    case 'date':
      return `${year}-${month}-${day}`
    case 'time':
      return `${hours}:${minutes}:${seconds}`
    case 'datetime':
    default:
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  }
}

/**
 * 格式化日期（简写，等同于 formatDateTime 默认 datetime 格式）
 * @param date 日期字符串或Date对象
 */
export function formatDate(date: string | Date | null | undefined): string {
  if (!date) return '-'
  return formatDateTime(date, 'datetime')
}

/**
 * 计算趋势百分比
 * @param current 当前值
 * @param previous 前期值
 */
export function calculateTrend(current: number, previous: number): number {
  if (previous === 0) return current > 0 ? 100 : 0
  return ((current - previous) / previous) * 100
}
