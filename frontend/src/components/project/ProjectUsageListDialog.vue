<template>
  <v-dialog v-model="modelOpen" width="800">
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-btn icon="mdi-close" variant="text" @click="close" />
        <span class="ml-2">{{ t('project.usageList.title') }}</span>
      </v-card-title>
      <v-card-text>
        <v-data-table
          item-value="id"
          :items="usages"
          :loading="loading"
        >
          <template #headers>
            <tr>
              <th class="text-left">{{ t('project.usageList.oss') }}</th>
              <th class="text-left">{{ t('project.usageList.version') }}</th>
              <th class="text-left">{{ t('project.usageList.role') }}</th>
              <th class="text-left">{{ t('project.usageList.scope') }}</th>
            </tr>
          </template>
          <template #item="{ item }">
            <tr>
              <td>{{ item.ossName }}</td>
              <td>{{ item.ossVersion }}</td>
              <td>{{ item.usageRole }}</td>
              <td>{{ item.scopeStatus }}</td>
            </tr>
          </template>
        </v-data-table>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="close">{{ t('project.usageList.close') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { ProjectUsage } from '@/api'
  import { computed, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { OssService, OssVersionsService, ProjectUsagesService } from '@/api'

  interface Props {
    open: boolean
    projectId?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{ (e: 'update:open', value: boolean): void }>()
  const { t } = useI18n()

  const modelOpen = computed({
    get: () => props.open,
    set: (val: boolean) => emit('update:open', val),
  })

  interface UsageItem extends ProjectUsage {
    ossName: string
    ossVersion: string
  }

  const usages = ref<UsageItem[]>([])
  const loading = ref(false)

  const ossNameCache = new Map<string, string>()
  const ossVersionCache = new Map<string, string>()

  watch(
    () => props.projectId,
    () => {
      if (modelOpen.value && props.projectId) {
        load()
      }
    },
  )

  watch(modelOpen, val => {
    if (val && props.projectId) {
      load()
    }
  })

  async function load () {
    if (!props.projectId) return
    loading.value = true
    try {
      const res = await ProjectUsagesService.listProjectUsages({ projectId: props.projectId })
      const items = res.items ?? []
      usages.value = await Promise.all(items.map(async u => {
        const ossName = await fetchOssName(u.ossId)
        const ossVersion = await fetchOssVersion(u.ossId, u.ossVersionId)
        return { ...u, ossName, ossVersion }
      }))
    } finally {
      loading.value = false
    }
  }

  async function fetchOssName (ossId: string) {
    let name = ossNameCache.get(ossId)
    if (!name) {
      const res = await OssService.getOssComponent({ ossId })
      name = res.name
      ossNameCache.set(ossId, name)
    }
    return name
  }

  async function fetchOssVersion (ossId: string, versionId: string) {
    let version = ossVersionCache.get(versionId)
    if (!version) {
      const res = await OssVersionsService.getOssVersion({ ossId, versionId })
      version = res.version
      ossVersionCache.set(versionId, version)
    }
    return version
  }

  function close () {
    emit('update:open', false)
  }
</script>
