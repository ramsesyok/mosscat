<template>
  <v-dialog v-model="modelOpen" max-width="600">
    <v-card>
      <v-card-title>
        {{ isNew ? t('oss.versionDetail.titleNew') : t('oss.versionDetail.titleEdit') }}
      </v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-text-field v-model="form.version" :disabled="!isNew" label="Version" required />
          <v-text-field v-model="form.releaseDate" label="Release Date" type="date" />
          <v-select
            v-model="form.licenseExpressionRaw"
            clearable
            item-props
            item-title="title"
            item-value="value"
            :items="licenseOptions"
            label="License"
          />
          <v-text-field v-model="form.purl" label="PURL" :rules="[urlRule]" />
          <v-text-field v-model="form.cpeListInput" label="CPE List" />
          <v-text-field v-model="form.hashSha256" label="SHA256" />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="primary" :loading="saving" @click="onSave">
          {{ t('oss.versionDetail.save') }}
        </v-btn>
        <v-btn text @click="close">{{ t('oss.versionDetail.cancel') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type {
    OssVersion,
    OssVersionCreateRequest,
    OssVersionUpdateRequest,
  } from '@/api'
  import { computed, onMounted, reactive, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { OssVersionsService } from '@/api'

  interface Props {
    open: boolean
    ossId?: string
    versionId?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{ (e: 'update:open', v: boolean): void, (e: 'saved'): void }>()

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

  const modelOpen = computed({
    get: () => props.open,
    set: (val: boolean) => emit('update:open', val),
  })

  const isNew = computed(() => !props.versionId)

  interface Form {
    version: string
    releaseDate: string
    licenseExpressionRaw?: string
    purl?: string
    cpeListInput: string
    hashSha256?: string
  }

  const form = reactive<Form>({
    version: '',
    releaseDate: '',
    licenseExpressionRaw: undefined,
    purl: undefined,
    cpeListInput: '',
    hashSha256: undefined,
  })

  const licenseOptions = ref<{ title: string, value: string, subtitle: string }[]>([])

  onMounted(async () => {
    try {
      const res = await fetch('/licenses.json')
      const data: Array<{ key: string, name: string, spdx_id: string, url: string }> = await res.json()
      licenseOptions.value = data.map(l => ({
        title: l.key,
        value: l.spdx_id,
        subtitle: l.name,
      }))
    } catch (error) {
      console.error(error)
    }
  })

  const formRef = ref()
  const saving = ref(false)
  const detail = ref<OssVersion | null>(null)

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
    () => props.versionId,
    () => {
      if (props.open) {
        loadDetail()
      }
    },
  )

  async function loadDetail () {
    if (!props.ossId || !props.versionId) {
      resetForm()
      detail.value = null
      return
    }
    try {
      detail.value = await OssVersionsService.getOssVersion({
        ossId: props.ossId,
        versionId: props.versionId,
      })
      form.version = detail.value.version
      form.releaseDate = detail.value.releaseDate ?? ''
      form.licenseExpressionRaw = detail.value.licenseExpressionRaw ?? undefined
      form.purl = detail.value.purl ?? undefined
      form.cpeListInput = (detail.value.cpeList ?? []).join(', ')
      form.hashSha256 = detail.value.hashSha256 ?? undefined
    } catch (error) {
      console.error(error)
    }
  }

  function resetForm () {
    form.version = ''
    form.releaseDate = ''
    form.licenseExpressionRaw = undefined
    form.purl = undefined
    form.cpeListInput = ''
    form.hashSha256 = undefined
  }

  function close () {
    emit('update:open', false)
  }

  async function onSave () {
    if (!props.ossId) return
    const result = await formRef.value?.validate()
    if (!result?.valid) return
    saving.value = true
    try {
      const cpeList = form.cpeListInput
        .split(',')
        .map(c => c.trim())
        .filter(Boolean)
      if (props.versionId) {
        const payload: OssVersionUpdateRequest = {
          releaseDate: form.releaseDate || undefined,
          licenseExpressionRaw: form.licenseExpressionRaw,
          purl: form.purl,
          cpeList,
          hashSha256: form.hashSha256,
        }
        await OssVersionsService.updateOssVersion({
          ossId: props.ossId,
          versionId: props.versionId,
          requestBody: payload,
        })
      } else {
        const payload: OssVersionCreateRequest = {
          version: form.version,
          releaseDate: form.releaseDate || undefined,
          licenseExpressionRaw: form.licenseExpressionRaw,
          purl: form.purl,
          cpeList,
          hashSha256: form.hashSha256,
        }
        await OssVersionsService.createOssVersion({
          ossId: props.ossId,
          requestBody: payload,
        })
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
