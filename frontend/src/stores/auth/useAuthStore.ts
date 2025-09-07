import { defineStore } from 'pinia'
import { AuthService, OpenAPI, type Role, UsersService } from '@/api'

interface AuthState {
  token: string
  roles: Role[]
}

const storedToken = typeof localStorage === 'undefined' ? '' : localStorage.getItem('jwtToken') || ''
if (storedToken) {
  OpenAPI.TOKEN = storedToken
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: storedToken,
    roles: [],
  }),
  actions: {
    async login (username: string, password: string) {
      const res = await AuthService.login({
        requestBody: { username, password },
      })
      this.token = res.accessToken
      OpenAPI.TOKEN = res.accessToken
      localStorage.setItem('jwtToken', res.accessToken)
      await this.fetchCurrentUser()
    },
    async fetchCurrentUser () {
      if (!this.token) {
        return
      }
      try {
        const user = await UsersService.getCurrentUser()
        this.roles = user.roles
      } catch (error) {
        console.error(error)
      }
    },
    async logout () {
      try {
        await AuthService.logout()
      } catch {
        // ignore logout errors
      }
      this.token = ''
      this.roles = []
      OpenAPI.TOKEN = undefined
      localStorage.removeItem('jwtToken')
    },
  },
})
