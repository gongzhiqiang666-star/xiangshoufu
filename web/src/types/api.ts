// API响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页响应
export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 分页请求参数
export interface PaginationParams {
  page?: number
  page_size?: number
}
