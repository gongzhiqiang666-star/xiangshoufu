import request from '@/utils/request'
import type { UploadImageResponse, UploadModule } from '@/types/upload'

// 上传图片
export function uploadImage(file: File, module?: UploadModule) {
  const formData = new FormData()
  formData.append('file', file)
  if (module) {
    formData.append('module', module)
  }

  return request<UploadImageResponse>({
    url: '/api/v1/upload/image',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}
