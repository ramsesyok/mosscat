/* eslint-disable no-useless-catch */
import type {
  Project,
  ProjectCreateRequest,
  ProjectUpdateRequest,
} from '@/api'
import { defineStore } from 'pinia'
import { ProjectsService } from '@/api'

export const useProjectStore = defineStore('project', {
  state: () => ({
    items: [] as Project[],
    total: 0,
    page: 1,
    size: 50,
    filters: { code: '', name: '' } as { code?: string, name?: string },
    loading: false,
  }),
  actions: {
    async fetchList () {
      this.loading = true
      try {
        const res = await ProjectsService.listProjects({
          page: this.page,
          size: this.size,
          code: this.filters.code || undefined,
          name: this.filters.name || undefined,
        })
        this.items = res.items ?? []
        this.total = res.total ?? 0
      } catch (error) {
        throw error
      } finally {
        this.loading = false
      }
    },
    async get (id: string) {
      return await ProjectsService.getProject({ projectId: id })
    },
    async create (payload: ProjectCreateRequest) {
      return await ProjectsService.createProject({ requestBody: payload })
    },
    async update (id: string, payload: ProjectUpdateRequest) {
      const updated = await ProjectsService.updateProject({ projectId: id, requestBody: payload })
      const idx = this.items.findIndex(i => i.id === updated.id)
      if (idx !== -1) {
        this.items[idx] = updated
      }
      return updated
    },
    async delete (id: string) {
      await ProjectsService.deleteProject({ projectId: id })
      this.items = this.items.filter(i => i.id !== id)
    },
  },
})
