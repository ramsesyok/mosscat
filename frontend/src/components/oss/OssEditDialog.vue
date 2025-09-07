<template>
  <v-dialog v-model="modelOpen" max-width="600">
    <v-card>
      <v-card-title>
        {{ isNew ? t('oss.detail.titleNew') : t('oss.detail.titleEdit') }}
        <span v-if="!isNew" class="ml-2">{{ detail?.name }}</span>
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-text-field v-model="form.name" label="Name" required />
          <v-text-field
            v-model="form.homepageUrl"
            label="Homepage URL"
            :rules="[urlRule]"
          />
          <v-text-field
            v-model="form.repositoryUrl"
            label="Repository URL"
            :rules="[urlRule]"
          />
          <v-select
            v-model="form.primaryLanguage"
            clearable
            :items="languageOptions"
            :label="t('oss.table.primaryLanguage')"
          />
          <v-select
            v-model="form.layer"
            clearable
            :items="layerOptions"
            :label="t('oss.table.layers')"
          />
          <v-select
            v-model="form.tagIds"
            chips
            clearable
            item-title="name"
            item-value="id"
            :items="tags"
            :label="t('oss.table.tags')"
            multiple
          />
          <v-textarea v-model="form.description" label="Description" />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="primary" :loading="saving" @click="onSave">{{ t('oss.detail.save') }}</v-btn>
        <v-btn text @click="close">{{ t('oss.detail.cancel') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { Layer, OssComponent, OssComponentCreateRequest, OssComponentUpdateRequest } from '@/api'
  import { storeToRefs } from 'pinia'
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { OssService } from '@/api'
  import { useTagStore } from '@/stores/useTagStore'

  interface Props {
    open: boolean
    ossId?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:open', value: boolean): void
    (e: 'saved'): void
  }>()

  const { t } = useI18n()

  const urlRule = (v?: string) => {
    if (!v) return true
    try {
      new URL(v)
      return true
    } catch {
      return t('error.invalidUrl')
    }
  }

  const languageOptions = ['C/C++', 'Java', 'Python', 'Go', 'Rust', 'PHP', 'C#']
  const layerOptions: Layer[] = [
    'OS',
    'RUNTIME',
    'DB',
    'MIDDLEWARE',
    'LIB',
    'FRAMEWORK',
    'TOOL_BUILD',
    'TOOL_DEV',
    'TOOL_TEST',
    'DOC',
    'OTHER',
  ]

  const tagStore = useTagStore()
  const { items: tags } = storeToRefs(tagStore)

  onMounted(async () => {
    try {
      await tagStore.fetch()
    } catch (error) {
      console.error(error)
    }
  })

  const modelOpen = computed({
    get: () => props.open,
    set: (val: boolean) => emit('update:open', val),
  })

  const isNew = computed(() => !props.ossId)

  const detail = ref<OssComponent | null>(null)

  interface OssComponentForm {
    name: string
    homepageUrl?: string
    repositoryUrl?: string
    description?: string
    primaryLanguage?: string
    layer?: Layer
    tagIds: string[]
  }

  const form = reactive<OssComponentForm>({
    name: '',
    homepageUrl: undefined,
    repositoryUrl: undefined,
    description: undefined,
    primaryLanguage: undefined,
    layer: undefined,
    tagIds: [],
  })

  const formRef = ref()
  const saving = ref(false)

  watch(
    () => props.open,
    val => {
      if (val) {
        loadDetail()
      }
    },
    { immediate: true },
  )

  watch(
    () => props.ossId,
    () => {
      if (props.open) {
        loadDetail()
      }
    },
    { immediate: true },
  )

  async function loadDetail () {
    if (!props.ossId) {
      detail.value = null
      resetForm()
      return
    }
    try {
      detail.value = await OssService.getOssComponent({ ossId: props.ossId })
      form.name = detail.value.name
      form.homepageUrl = detail.value.homepageUrl ?? undefined
      form.repositoryUrl = detail.value.repositoryUrl ?? undefined
      form.description = detail.value.description ?? undefined
      form.primaryLanguage = detail.value.primaryLanguage ?? undefined
      form.layer = detail.value.layers?.[0] ?? undefined
      form.tagIds = (detail.value.tags ?? []).map(tag => tag.id)
    } catch (error) {
      console.error(error)
    }
  }

  function resetForm () {
    form.name = ''
    form.homepageUrl = undefined
    form.repositoryUrl = undefined
    form.description = undefined
    form.primaryLanguage = undefined
    form.layer = undefined
    form.tagIds = []
  }

  function close () {
    emit('update:open', false)
  }

  async function onSave () {
    const result = await formRef.value?.validate()
    if (!result?.valid) return
    saving.value = true
    try {
      if (props.ossId) {
        const payload: OssComponentUpdateRequest = {
          name: form.name,
          homepageUrl: form.homepageUrl,
          repositoryUrl: form.repositoryUrl,
          description: form.description,
          primaryLanguage: form.primaryLanguage,
          layers: form.layer ? [form.layer] : [],
          tagIds: form.tagIds,
        }
        await OssService.updateOssComponent({ ossId: props.ossId, requestBody: payload })
      } else {
        const payload: OssComponentCreateRequest = {
          name: form.name,
          homepageUrl: form.homepageUrl,
          repositoryUrl: form.repositoryUrl,
          description: form.description,
          primaryLanguage: form.primaryLanguage,
          layers: form.layer ? [form.layer] : [],
          tagIds: form.tagIds,
        }
        await OssService.createOssComponent({ requestBody: payload })
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
