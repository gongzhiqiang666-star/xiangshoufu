// 上传相关类型定义

// 上传响应
export interface UploadImageResponse {
  id: number
  original_name: string
  file_url: string
  thumbnail_url?: string
  file_size: number
  width: number
  height: number
  mime_type: string
}

// 上传模块类型
export type UploadModule = 'banner' | 'poster' | ''
