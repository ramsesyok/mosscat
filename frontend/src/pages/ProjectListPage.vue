<template>
  <v-container fluid>
    <ProjectFilterBar
      v-model:code="filters.code"
      v-model:name="filters.name"
      class="mb-4"
      @search="onSearch"
    />
    <ProjectTable
      :items="items"
      :items-per-page="size"
      :loading="loading || exporting"
      :page="page"
      :total-items="total"
      @create="onCreate"
      @detail="onDetail"
      @export="onExport"
      @row-click="onRowClick"
      @update:items-per-page="onItemsPerPageChange"
      @update:page="onPageChange"
      @usage="onUsage"
    />
    <ProjectDetailDialog
      v-model:open="detailOpen"
      :project-id="selectedId"
      @saved="onSaved"
    />
    <ProjectUsageListDialog
      v-model:open="usageOpen"
      :project-id="selectedId"
    />
    <v-snackbar v-model="snackbar" timeout="3000">{{ snackbarMessage }}</v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
  import { useI18n } from 'vue-i18n'
  import { ExportService } from '@/api'
  import { useProjectStore } from '@/stores/useProjectStore'

  const store = useProjectStore()
  const { items, page, size, total, filters, loading } = storeToRefs(store)
  const { fetchList } = store

  const { t } = useI18n()

  const detailOpen = ref(false)
  const selectedId = ref<string>()
  const usageOpen = ref(false)
  const exporting = ref(false)
  const snackbar = ref(false)
  const snackbarMessage = ref('')

  onMounted(() => {
    fetchList()
  })

  function onSearch () {
    store.page = 1
    fetchList()
  }
  function onPageChange (p: number) {
    store.page = p
    fetchList()
  }
  function onItemsPerPageChange (s: number) {
    store.page = 1
    store.size = s
    fetchList()
  }
  function onRowClick (item: any) {
    openDetail(item.id)
  }
  function onDetail (item: any) {
    openDetail(item.id)
  }
  function onUsage (item: any) {
    selectedId.value = item.id
    usageOpen.value = true
  }
  function onCreate () {
    openDetail()
  }
  function openDetail (id?: string) {
    selectedId.value = id
    detailOpen.value = true
  }
  async function onSaved () {
    await fetchList()
    showToast(t('project.detail.save'))
  }
  function showToast (msg: string) {
    snackbarMessage.value = msg
    snackbar.value = true
  }
  async function onExport (payload: { id: string, format: 'csv' | 'spdx-json' }) {
    exporting.value = true
    try {
      const data = await ExportService.exportProjectArtifacts({
        projectId: payload.id,
        format: payload.format,
      })
      const text = typeof data === 'string' ? data : JSON.stringify(data)
      const blob = new Blob([text], { type: payload.format === 'csv' ? 'text/csv' : 'application/json' })
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = payload.format === 'csv' ? 'project.csv' : 'project.spdx.json'
      a.click()
      URL.revokeObjectURL(url)
      showToast(t('export.toast'))
    } catch (error) {
      console.error(error)
    } finally {
      exporting.value = false
    }
  }
</script>
