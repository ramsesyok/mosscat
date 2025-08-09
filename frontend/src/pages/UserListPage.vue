<template>
  <v-container fluid>
    <UserTable
      :items="items"
      :items-per-page="size"
      :loading="loading"
      :page="page"
      :total-items="total"
      @create="onCreate"
      @delete="onDelete"
      @detail="onDetail"
      @update:items-per-page="onItemsPerPageChange"
      @update:page="onPageChange"
    />
    <UserDetailDialog
      v-model:open="detailOpen"
      :user-id="selectedId"
      @saved="onSaved"
    />
  </v-container>
</template>

<script setup lang="ts">
  import type { User } from '@/api'
  import { storeToRefs } from 'pinia'
  import { onMounted, ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import UserDetailDialog from '@/components/user/UserDetailDialog.vue'
  import UserTable from '@/components/user/UserTable.vue'
  import { useUserStore } from '@/stores/useUserStore'

  const store = useUserStore()
  const { items, page, size, total, loading } = storeToRefs(store)
  const { fetchList } = store
  const { t } = useI18n()

  const detailOpen = ref(false)
  const selectedId = ref<string>()

  onMounted(() => {
    fetchList()
  })

  function onPageChange (p: number) {
    store.page = p
    fetchList()
  }
  function onItemsPerPageChange (s: number) {
    store.page = 1
    store.size = s
    fetchList()
  }
  function onDetail (item: User) {
    openDetail(item.id)
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
  }
  async function onDelete (item: User) {
    if (!confirm(t('common.confirmDelete'))) return
    await store.delete(item.id)
    await fetchList()
  }
</script>
