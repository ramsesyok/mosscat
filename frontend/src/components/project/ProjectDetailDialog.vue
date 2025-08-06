<template>
  <v-dialog v-model="modelOpen" max-width="600">
    <v-card>
      <v-card-title>{{ isNew ? $t('project.detail.titleNew') : $t('project.detail.titleEdit') }}</v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-text-field v-model="form.projectCode" :disabled="!isNew" label="Code" required />
          <v-text-field v-model="form.name" label="Name" required />
          <v-text-field v-model="form.department" label="Department" />
          <v-text-field v-model="form.manager" label="Manager" />
          <v-text-field v-model="form.deliveryDate" label="Delivery Date" type="date" />
          <v-textarea v-model="form.description" label="Description" />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="primary" :loading="saving" @click="onSave">{{ $t('project.detail.save') }}</v-btn>
        <v-btn text @click="close">{{ $t('project.detail.cancel') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { ProjectCreateRequest, ProjectUpdateRequest } from '@/api'
  import { computed, reactive, ref, watch } from 'vue'
  import { useProjectStore } from '@/stores/useProjectStore'

  interface Props {
    open: boolean
    projectId?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:open', v: boolean): void
    (e: 'saved'): void
  }>()

  const store = useProjectStore()

  const modelOpen = computed({
    get: () => props.open,
    set: (v: boolean) => emit('update:open', v),
  })

  const isNew = computed(() => !props.projectId)

  const form = reactive<ProjectCreateRequest>({
    projectCode: '',
    name: '',
    department: undefined,
    manager: undefined,
    deliveryDate: undefined,
    description: undefined,
  })

  const formRef = ref()
  const saving = ref(false)

  watch(
    () => props.open,
    val => {
      if (val) loadDetail()
    },
  )
  watch(
    () => props.projectId,
    () => {
      if (props.open) loadDetail()
    },
  )

  async function loadDetail () {
    if (!props.projectId) {
      resetForm()
      return
    }
    try {
      const detail = await store.get(props.projectId)
      form.projectCode = detail.projectCode
      form.name = detail.name
      form.department = detail.department ?? undefined
      form.manager = detail.manager ?? undefined
      form.deliveryDate = detail.deliveryDate ?? undefined
      form.description = detail.description ?? undefined
    } catch (error) {
      console.error(error)
    }
  }

  function resetForm () {
    form.projectCode = ''
    form.name = ''
    form.department = undefined
    form.manager = undefined
    form.deliveryDate = undefined
    form.description = undefined
  }

  function close () {
    emit('update:open', false)
  }

  async function onSave () {
    saving.value = true
    try {
      if (isNew.value) {
        const payload: ProjectCreateRequest = { ...form }
        await store.create(payload)
      } else if (props.projectId) {
        const payload: ProjectUpdateRequest = {
          name: form.name,
          department: form.department,
          manager: form.manager,
          deliveryDate: form.deliveryDate,
          description: form.description,
        }
        await store.update(props.projectId, payload)
      }
      emit('saved')
      close()
    } catch (error) {
      console.error(error)
    } finally {
      saving.value = false
    }
  }
</script>
