/* eslint-disable no-useless-catch */
import type { ScopePolicy, ScopePolicyUpdateRequest } from '@/api'
import { defineStore } from 'pinia'
import { ScopePolicyService } from '@/api'

export const useScopePolicyStore = defineStore('scopePolicy', {
  state: () => ({
    policy: null as ScopePolicy | null,
    loading: false,
  }),
  actions: {
    async fetch () {
      this.loading = true
      try {
        this.policy = await ScopePolicyService.getScopePolicy()
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },
    async update (payload: ScopePolicyUpdateRequest) {
      this.loading = true
      try {
        const updated = await ScopePolicyService.updateScopePolicy({ requestBody: payload })
        this.policy = updated
        return updated
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },
  },
})
