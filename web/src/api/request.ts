import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { ElMessage } from 'element-plus'
import { getToken, clearAuth } from '@/utils/storage'
import type { ApiResponse } from '@/types'
import router from '@/router'

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    // 添加Token到请求头
    const token = getToken()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { data } = response

    // 业务错误处理
    if (data.code !== 0) {
      ElMessage.error(data.message || '请求失败')

      // Token过期或无效
      if (data.code === 401) {
        clearAuth()
        router.push('/login')
      }

      return Promise.reject(new Error(data.message || '请求失败'))
    }

    return response
  },
  (error) => {
    console.error('Response error:', error)

    // HTTP错误处理
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          ElMessage.error('登录已过期，请重新登录')
          clearAuth()
          router.push('/login')
          break
        case 403:
          ElMessage.error('没有权限访问该资源')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(data?.message || '网络请求失败')
      }
    } else if (error.request) {
      ElMessage.error('网络连接失败，请检查网络')
    } else {
      ElMessage.error('请求配置错误')
    }

    return Promise.reject(error)
  }
)

// 封装GET请求
export function get<T>(url: string, params?: object, config?: AxiosRequestConfig): Promise<T> {
  return request.get(url, { params, ...config }).then((res) => res.data.data)
}

// 封装POST请求
export function post<T>(url: string, data?: object, config?: AxiosRequestConfig): Promise<T> {
  return request.post(url, data, config).then((res) => res.data.data)
}

// 封装PUT请求
export function put<T>(url: string, data?: object, config?: AxiosRequestConfig): Promise<T> {
  return request.put(url, data, config).then((res) => res.data.data)
}

// 封装DELETE请求
export function del<T>(url: string, config?: AxiosRequestConfig): Promise<T> {
  return request.delete(url, config).then((res) => res.data.data)
}

export default request
