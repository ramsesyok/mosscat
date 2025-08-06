/* eslint-disable no-useless-catch */
import type { Tag } from '@/api'
import { defineStore } from 'pinia'
import { TagsService } from '@/api'

export const useTagStore = defineStore('tag', {
  state: () => ({
    items: [] as Tag[],
    loading: false,
  }),
  actions: {
    async fetch () {
      this.loading = true
      try {
        this.items = await TagsService.listTags()
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },
    async create (name: string) {
      try {
        await TagsService.createTag({ requestBody: { name } })
      } catch (error) {
        throw error
      }
    },
    async remove (id: string) {
      try {
        await TagsService.deleteTag({ tagId: id })
      } catch (error) {
        throw error
      }
    },
  },
})
