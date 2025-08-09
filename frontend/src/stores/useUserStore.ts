import type { Role, User, UserCreateRequest, UserUpdateRequest } from '@/api'
import { defineStore } from 'pinia'
import { UsersService } from '@/api'

export const useUserStore = defineStore('user', {
  state: () => ({
    items: [] as User[],
    total: 0,
    page: 1,
    size: 50,
    filters: { username: '', role: undefined as Role | undefined },
    loading: false,
  }),
  actions: {
    async fetchList () {
      this.loading = true
      try {
        const res = await UsersService.listUsers({
          page: this.page,
          size: this.size,
          username: this.filters.username || undefined,
          role: this.filters.role,
        })
        this.items = res.items ?? []
        this.total = res.total ?? 0
      } catch (error) {
        console.error(error)
      } finally {
        this.loading = false
      }
    },
    async get (id: string) {
      return await UsersService.getUser({ userId: id })
    },
    async create (payload: UserCreateRequest) {
      return await UsersService.createUser({ requestBody: payload })
    },
    async update (id: string, payload: UserUpdateRequest) {
      const updated = await UsersService.updateUser({ userId: id, requestBody: payload })
      const idx = this.items.findIndex(i => i.id === updated.id)
      if (idx !== -1) {
        this.items[idx] = updated
      }
      return updated
    },
    async delete (id: string) {
      await UsersService.deleteUser({ userId: id })
      this.items = this.items.filter(i => i.id !== id)
    },
  },
})
