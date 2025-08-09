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
              <th class="text-left">{{ t('project.usageList.modified') }}</th>
              <th class="text-left">{{ t('project.usageList.actions') }}</th>
            </tr>
          </template>
          <template #item="{ item }">
            <tr>
              <td>{{ item.ossName }}</td>
              <td>{{ item.ossVersion }}</td>
              <td>{{ item.usageRole }}</td>
              <td>{{ item.scopeStatus }}</td>
              <td>
                <v-chip
                  v-if="item.modified"
                  color="warning"
                  size="small"
                >
                  {{ t('project.usageList.modifiedYes') }}
                </v-chip>
                <span v-else>{{ t('project.usageList.modifiedNo') }}</span>
              </td>
              <td>
                <v-btn
                  icon="mdi-pencil"
                  size="small"
                  variant="text"
                  @click="editItem(item)"
                />
              </td>
            </tr>
          </template>
        </v-data-table>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="close">{{ t('project.usageList.close') }}</v-btn>
      </v-card-actions>
    </v-card>

    <!-- Edit Project Usage Dialog -->
    <v-dialog v-model="editDialog" max-width="600">
      <v-card>
        <v-card-title>{{ t('project.usageEdit.title') }}</v-card-title>
        <v-card-text>
          <v-form ref="editFormRef">
            <h3 class="mb-3">{{ t('project.usage.modificationInfo') }}</h3>
            <v-switch
              v-model="editForm.modified"
              :label="t('project.usage.modified')"
            />
            <v-textarea
              v-if="editForm.modified"
              v-model="editForm.modificationDescription"
              :label="t('project.usage.modificationDescription')"
              rows="3"
            />
            <v-select
              v-if="editForm.modified"
              v-model="editForm.supplierType"
              clearable
              item-title="title"
              item-value="value"
              :items="supplierTypeOptions"
              :label="t('project.usage.supplierType')"
            />
            <v-text-field
              v-if="editForm.modified"
              v-model="editForm.forkOriginUrl"
              :label="t('project.usage.forkOriginUrl')"
              :rules="[urlRule]"
            />
          </v-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" :loading="saving" @click="onSave">
            {{ t('project.usageEdit.save') }}
          </v-btn>
          <v-btn text @click="closeEdit">{{ t('project.usageEdit.cancel') }}</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { ProjectUsage, SupplierType } from '@/api'
  import { computed, reactive, ref, watch } from 'vue'
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
  const editDialog = ref(false)
  const editFormRef = ref()
  const saving = ref(false)
  const selectedUsage = ref<ProjectUsage>()

  const ossNameCache = new Map<string, string>()
  const ossVersionCache = new Map<string, string>()

  const urlRule = (v?: string) => {
    if (!v) return true
    try {
      new URL(v)
      return true
    } catch {
      return t('error.invalidUrl')
    }
  }

  interface EditForm {
    modified: boolean
    modificationDescription?: string
    supplierType?: SupplierType
    forkOriginUrl?: string
  }

  const editForm = reactive<EditForm>({
    modified: false,
    modificationDescription: undefined,
    supplierType: undefined,
    forkOriginUrl: undefined,
  })

  const supplierTypeOptions = ref<{ title: string, value: SupplierType }[]>([
    { title: 'UPSTREAM', value: 'UPSTREAM' },
    { title: 'INTERNAL_FORK', value: 'INTERNAL_FORK' },
    { title: 'REPACKAGE', value: 'REPACKAGE' },
  ])

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

  function editItem (item: UsageItem) {
    selectedUsage.value = item
    editForm.modified = item.modified ?? false
    editForm.modificationDescription = item.modificationDescription ?? undefined
    editForm.supplierType = item.supplierType ?? undefined
    editForm.forkOriginUrl = item.forkOriginUrl ?? undefined
    editDialog.value = true
  }

  function closeEdit () {
    editDialog.value = false
    selectedUsage.value = undefined
    resetEditForm()
  }

  function resetEditForm () {
    editForm.modified = false
    editForm.modificationDescription = undefined
    editForm.supplierType = undefined
    editForm.forkOriginUrl = undefined
  }

  async function onSave () {
    if (!selectedUsage.value || !props.projectId) return
    const result = await editFormRef.value?.validate()
    if (!result?.valid) return

    saving.value = true
    try {
      await ProjectUsagesService.updateProjectUsage({
        projectId: props.projectId,
        usageId: selectedUsage.value.id,
        requestBody: {
          modified: editForm.modified,
          modificationDescription: editForm.modified ? editForm.modificationDescription : undefined,
          supplierType: editForm.modified ? editForm.supplierType : undefined,
          forkOriginUrl: editForm.modified ? editForm.forkOriginUrl : undefined,
        },
      })
      await load()
      closeEdit()
    } catch (error) {
      console.error(error)
    } finally {
      saving.value = false
    }
  }
</script>
