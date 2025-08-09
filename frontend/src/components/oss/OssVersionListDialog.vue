<template>
  <v-dialog v-model="modelOpen">
    <v-card>
      <v-card-title class="d-flex align-center">
        <v-btn icon="mdi-close" variant="text" @click="close" />
        <span class="ml-2">{{ t('oss.versionList.title') }}</span>
        <v-spacer />
        <v-btn color="primary" @click="openAdd">{{ t('oss.versionList.add') }}</v-btn>
      </v-card-title>
      <v-card-text>
        <v-data-table
          item-value="id"
          :items="versions"
          :loading="loading"
        >
          <template #headers>
            <tr>
              <th class="text-left">{{ t('oss.versionList.version') }}</th>
              <th class="text-left">{{ t('oss.versionList.releaseDate') }}</th>
              <th class="text-left" style="width: 120px">{{ t('oss.table.actions') }}</th>
            </tr>
          </template>
          <template #item="{ item }">
            <tr>
              <td>{{ item.version }}</td>
              <td>{{ item.releaseDate }}</td>
              <td class="text-right" style="width: 120px">
                <v-btn icon="mdi-pencil" variant="text" @click.stop="openEdit(item.id)" />
                <v-btn icon="mdi-package-variant-plus" variant="text" @click.stop="openProjectSelect(item.id)" />
                <v-btn icon="mdi-view-list" variant="text" @click.stop="openUsageList(item.id)" />
                <v-menu location="bottom">
                  <template #activator="{ props: activatorProps }">
                    <v-btn icon="mdi-dots-vertical" variant="text" v-bind="activatorProps" />
                  </template>
                  <v-list density="compact">
                    <v-list-item @click="remove(item.id)">
                      <v-list-item-title>{{ t('common.delete') }}</v-list-item-title>
                    </v-list-item>
                  </v-list>
                </v-menu>
              </td>
            </tr>
          </template>
        </v-data-table>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="close">{{ t('oss.versionList.close') }}</v-btn>
      </v-card-actions>
    </v-card>
    <OssVersionEditDialog
      v-model:open="editOpen"
      :oss-id="props.ossId"
      :version-id="selectedId"
      @saved="load"
    />
    <ProjectSelectDialog
      v-if="props.ossId && selectedId"
      v-model:open="projectOpen"
      :oss-id="props.ossId"
      :oss-version-id="selectedId"
      @added="load"
    />
  </v-dialog>
</template>

<script setup lang="ts">
  import type { OssVersion } from '@/api'
  import { computed, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { OssVersionsService } from '@/api'
  import ProjectSelectDialog from '../project/ProjectSelectDialog.vue'
  import OssVersionEditDialog from './OssVersionEditDialog.vue'

  interface Props {
    open: boolean
    ossId?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{ (e: 'update:open', value: boolean): void }>()
  const { t } = useI18n()

  const modelOpen = computed({
    get: () => props.open,
    set: (val: boolean) => emit('update:open', val),
  })

  const versions = ref<OssVersion[]>([])
  const loading = ref(false)
  const editOpen = ref(false)
  const projectOpen = ref(false)
  const selectedId = ref<string>()

  watch(
    () => props.open,
    val => {
      if (val) {
        load()
      }
    },
  )

  watch(
    () => props.ossId,
    () => {
      if (props.open) {
        load()
      }
    },
  )

  async function load () {
    if (!props.ossId) {
      versions.value = []
      return
    }
    loading.value = true
    try {
      const res = await OssVersionsService.listOssVersions({ ossId: props.ossId })
      versions.value = res.items ?? []
    } finally {
      loading.value = false
    }
  }

  function close () {
    emit('update:open', false)
  }

  function openAdd () {
    selectedId.value = undefined
    editOpen.value = true
  }

  function openEdit (id: string) {
    selectedId.value = id
    editOpen.value = true
  }

  function openProjectSelect (id: string) {
    selectedId.value = id
    projectOpen.value = true
  }

  function openUsageList (id: string) {
    console.log('open usage list', id)
  }

  async function remove (id: string) {
    if (!props.ossId) return
    await OssVersionsService.deleteOssVersion({ ossId: props.ossId, versionId: id })
    await load()
  }
</script>
