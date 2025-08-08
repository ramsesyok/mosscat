import type { AppConfig } from '@/config'
import type { ApiError, Problem } from '@/types/problem'
import axios, { type AxiosError } from 'axios'
import { OpenAPI } from '@/api'
import router from '@/router'
import { useAuthStore } from '@/stores/auth/useAuthStore'

export function configureHttp (config: AppConfig) {
  axios.defaults.baseURL = config.apiBaseUrl
  axios.defaults.timeout = 30_000
  OpenAPI.BASE = config.apiBaseUrl
}

axios.interceptors.request.use(config => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers = config.headers || {}
    if (!('Authorization' in config.headers)) {
      config.headers['Authorization'] = `Bearer ${auth.token}`
    }
  }
  return config
})

axios.interceptors.response.use(
  response => response,
  (error: AxiosError<Problem>) => {
    const problem = error.response?.data
    const apiError: ApiError = new Error(problem?.title || 'ネットワークエラー') as ApiError
    apiError.status = error.response?.status
    if (problem && problem.title) {
      apiError.problem = problem
    }
    if (apiError.status === 401) {
      const url = error.config?.url || ''
      if (!url.endsWith('/auth/login') && !url.endsWith('/auth/logout')) {
        const auth = useAuthStore()
        auth.logout()
        if (router.currentRoute.value.path !== '/') {
          router.push('/')
        }
      }
    }
    return Promise.reject(apiError)
  },
)

export { default } from 'axios'
