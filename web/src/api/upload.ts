import request from './request'
import type { UploadImageResponse, UploadModule } from '@/types/upload'

// 上传图片
export function uploadImage(file: File, module?: UploadModule) {
  const formData = new FormData()
  formData.append('file', file)
  if (module) {
    formData.append('module', module)
  }

  return request.post('/api/v1/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }).then(res => res.data.data as UploadImageResponse)
}
