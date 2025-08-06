/* eslint-disable no-useless-catch */
import type { OssComponent, OssComponentUpdateRequest } from '@/api'
import { defineStore } from 'pinia'
import { OssService } from '@/api'

export const useOssStore = defineStore('oss', {
  state: () => ({
    list: [] as OssComponent[],
    page: 1,
    size: 50,
    total: 0,
    loading: false,
    filters: {
      name: '' as string,
      layers: undefined as string | undefined,
      tag: undefined as string | undefined,
      inScopeOnly: false,
    },
  }),
  actions: {
    async fetchList () {
      this.loading = true
      try {
        const res = await OssService.listOssComponents({
          page: this.page,
          size: this.size,
          name: this.filters.name || undefined,
          layers: this.filters.layers,
          tag: this.filters.tag,
          inScopeOnly: this.filters.inScopeOnly || undefined,
        })
        this.list = res.items ?? []
        this.total = res.total ?? 0
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },
    async update (id: string, payload: OssComponentUpdateRequest) {
      try {
        const updated = await OssService.updateOssComponent({
          ossId: id,
          requestBody: payload,
        })
        const idx = this.list.findIndex(item => item.id === updated.id)
        if (idx !== -1) {
          this.list[idx] = updated
        }
        return updated
      } catch (error) {
        throw error
      }
    },
  },
})
