<template>
  <v-dialog v-model="modelOpen" width="800">
    <v-card>
      <v-card-title>{{ t('project.listTitle') }}</v-card-title>
      <v-card-text>
        <div v-if="store.loading" class="d-flex justify-center pa-6">
          <v-progress-circular color="primary" indeterminate />
        </div>
        <v-data-table
          v-else
          item-value="id"
          :items="store.items"
          :items-length="store.total"
          :items-per-page="store.size"
          :page="store.page"
          @update:items-per-page="s => { store.size = s; store.fetchList() }"
          @update:page="p => { store.page = p; store.fetchList() }"
        >
          <template #headers>
            <tr>
              <th class="text-left">{{ t('project.table.code') }}</th>
              <th class="text-left">{{ t('project.table.name') }}</th>
              <th class="text-left">{{ t('project.table.deliveryDate') }}</th>
            </tr>
          </template>
          <template #item="{ item }">
            <tr
              :class="{ 'bg-grey-lighten-3': selected?.id === item.id }"
              @click="selected = item"
            >
              <td>{{ item.projectCode }}</td>
              <td>{{ item.name }}</td>
              <td>{{ item.deliveryDate }}</td>
            </tr>
          </template>
        </v-data-table>

        <v-divider class="my-4" />

        <v-form ref="formRef">
          <v-switch
            v-model="form.modified"
            :label="t('project.usage.modified')"
          />
          <v-textarea
            v-if="form.modified"
            v-model="form.modificationDescription"
            :label="t('project.usage.modificationDescription')"
            rows="3"
          />
          <v-select
            v-if="form.modified"
            v-model="form.supplierType"
            clearable
            item-title="title"
            item-value="value"
            :items="supplierTypeOptions"
            :label="t('project.usage.supplierType')"
          />
          <v-text-field
            v-if="form.modified"
            v-model="form.forkOriginUrl"
            :label="t('project.usage.forkOriginUrl')"
            :rules="[urlRule]"
          />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn
          color="primary"
          :disabled="!selected"
          @click="registerUsage"
        >{{ t('project.registerUsage') }}</v-btn>
        <v-btn text @click="close">{{ t('common.cancel') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { Project, SupplierType } from '@/api'
  import { computed, reactive, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { ProjectUsagesService } from '@/api'
  import { useProjectStore } from '@/stores/useProjectStore'

  interface Props {
    open: boolean
    ossId: string
    ossVersionId: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{ (e: 'update:open', value: boolean): void, (e: 'added'): void }>()
  const { t } = useI18n()
  const store = useProjectStore()

  const urlRule = (v?: string) => {
    if (!v) return true
    try {
      new URL(v)
      return true
    } catch {
      return t('error.invalidUrl')
    }
  }

  const modelOpen = computed({
    get: () => props.open,
    set: (val: boolean) => emit('update:open', val),
  })

  const selected = ref<Project>()
  const formRef = ref()

  interface Form {
    modified: boolean
    modificationDescription?: string
    supplierType?: SupplierType
    forkOriginUrl?: string
  }

  const form = reactive<Form>({
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

  watch(modelOpen, val => {
    if (val) {
      store.fetchList()
      selected.value = undefined
      resetForm()
    }
  })

  function resetForm () {
    form.modified = false
    form.modificationDescription = undefined
    form.supplierType = undefined
    form.forkOriginUrl = undefined
  }

  async function registerUsage () {
    if (!selected.value) return
    const result = await formRef.value?.validate()
    if (!result?.valid) return

    await ProjectUsagesService.createProjectUsage({
      projectId: selected.value.id,
      requestBody: {
        ossId: props.ossId,
        ossVersionId: props.ossVersionId,
        usageRole: 'RUNTIME_REQUIRED',
        modified: form.modified,
        modificationDescription: form.modified ? form.modificationDescription : undefined,
        supplierType: form.modified ? form.supplierType : undefined,
        forkOriginUrl: form.modified ? form.forkOriginUrl : undefined,
      },
    })
    emit('added')
    emit('update:open', false)
  }

  function close () {
    emit('update:open', false)
  }
</script>
