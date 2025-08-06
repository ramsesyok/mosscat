import { defineStore } from 'pinia'
import { AuthService, OpenAPI } from '@/api'

interface AuthState {
  token: string
}

const storedToken = typeof localStorage === 'undefined' ? '' : localStorage.getItem('jwtToken') || ''
if (storedToken) {
  OpenAPI.TOKEN = storedToken
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: storedToken,
  }),
  actions: {
    async login (username: string, password: string) {
      const res = await AuthService.login({
        requestBody: { username, password },
      })
      this.token = res.accessToken
      OpenAPI.TOKEN = res.accessToken
      localStorage.setItem('jwtToken', res.accessToken)
    },
    logout () {
      this.token = ''
      OpenAPI.TOKEN = undefined
      localStorage.removeItem('jwtToken')
    },
  },
})
