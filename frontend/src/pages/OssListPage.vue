<template>
  <v-container fluid>
    <OssTable ref="tableRef" @delete="onDelete" />
  </v-container>
</template>

<script setup lang="ts">
  import type { OssComponent } from '@/api'
  import { ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { OssService } from '@/api'

  const { t } = useI18n()
  const tableRef = ref<{ reload: () => void }>()

  async function onDelete (item: OssComponent) {
    if (!confirm(t('common.confirmDelete'))) return
    try {
      await OssService.deprecateOssComponent({ ossId: item.id })
      tableRef.value?.reload()
    } catch (error) {
      console.error(error)
    }
  }
</script>
