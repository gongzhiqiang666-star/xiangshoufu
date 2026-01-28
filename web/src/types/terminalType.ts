// 终端类型
export interface TerminalType {
  id: number
  channel_id: number
  channel_code: string
  channel_name: string
  brand_code: string
  brand_name: string
  model_code: string
  model_name: string
  full_name: string
  description: string
  status: number
  created_at: string
  updated_at: string
}

// 终端类型查询参数
export interface TerminalTypeQueryParams {
  channel_id?: number
  status?: number
  keyword?: string
  page?: number
  page_size?: number
}

// 创建终端类型参数
export interface CreateTerminalTypeParams {
  channel_id: number
  brand_code: string
  brand_name: string
  model_code: string
  model_name?: string
  description?: string
}

// 更新终端类型参数
export interface UpdateTerminalTypeParams {
  brand_code?: string
  brand_name?: string
  model_code?: string
  model_name?: string
  description?: string
}

// 终端类型状态
export const TerminalTypeStatus = {
  DISABLED: 0,
  ENABLED: 1,
} as const

export const TerminalTypeStatusText: Record<number, string> = {
  [TerminalTypeStatus.DISABLED]: '禁用',
  [TerminalTypeStatus.ENABLED]: '启用',
}
