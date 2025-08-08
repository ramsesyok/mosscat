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
          <v-text-field v-model="form.licenseExpressionRaw" label="License Expression" />
          <v-text-field v-model="form.purl" label="PURL" :rules="[urlRule]" />
          <v-text-field v-model="form.cpeListInput" label="CPE List" />
          <v-text-field v-model="form.hashSha256" label="SHA256" />
          <v-checkbox v-model="form.modified" label="Modified" />
          <v-text-field v-model="form.modificationDescription" label="Modification Description" />
          <v-select
            v-model="form.supplierType"
            clearable
            :items="supplierTypeOptions"
            label="Supplier Type"
          />
          <v-text-field
            v-model="form.forkOriginUrl"
            label="Fork Origin URL"
            :rules="[urlRule]"
          />
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
    SupplierType,
  } from '@/api'
  import { computed, reactive, ref, watch } from 'vue'
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
    modified: boolean
    modificationDescription?: string
    supplierType?: SupplierType
    forkOriginUrl?: string
  }

  const form = reactive<Form>({
    version: '',
    releaseDate: '',
    licenseExpressionRaw: undefined,
    purl: undefined,
    cpeListInput: '',
    hashSha256: undefined,
    modified: false,
    modificationDescription: undefined,
    supplierType: undefined,
    forkOriginUrl: undefined,
  })

  const supplierTypeOptions: SupplierType[] = ['UPSTREAM', 'INTERNAL_FORK', 'REPACKAGE']

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
      form.modified = detail.value.modified ?? false
      form.modificationDescription = detail.value.modificationDescription ?? undefined
      form.supplierType = detail.value.supplierType ?? undefined
      form.forkOriginUrl = detail.value.forkOriginUrl ?? undefined
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
    form.modified = false
    form.modificationDescription = undefined
    form.supplierType = undefined
    form.forkOriginUrl = undefined
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
          modified: form.modified,
          modificationDescription: form.modificationDescription,
          supplierType: form.supplierType,
          forkOriginUrl: form.forkOriginUrl,
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
          modified: form.modified,
          modificationDescription: form.modificationDescription,
          supplierType: form.supplierType,
          forkOriginUrl: form.forkOriginUrl,
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
