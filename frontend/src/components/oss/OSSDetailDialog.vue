<template>
  <v-dialog v-model="modelOpen" max-width="600">
    <v-card>
      <v-card-title>{{ detail?.name }}</v-card-title>
      <v-card-text>
        <div v-if="detail">
          <div>
            <strong>Name:</strong>
            <span>{{ detail.name }}</span>
          </div>
          <div v-if="detail.homepageUrl" class="mt-2">
            <strong>Homepage URL:</strong>
            <a
              :href="detail.homepageUrl"
              rel="noopener"
              target="_blank"
            >{{ detail.homepageUrl }}</a>
          </div>
          <div v-if="detail.repositoryUrl" class="mt-2">
            <strong>Repository URL:</strong>
            <a
              :href="detail.repositoryUrl"
              rel="noopener"
              target="_blank"
            >{{ detail.repositoryUrl }}</a>
          </div>
          <div v-if="detail.description" class="mt-2">
            <strong>Description:</strong>
            <p>{{ detail.description }}</p>
          </div>
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="close">{{ t('oss.detail.cancel') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { OssComponent } from '@/api'
  import { computed, ref, watch } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { OssService } from '@/api'

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

  const detail = ref<OssComponent | null>(null)

  watch(
    () => props.open,
    val => {
      if (val) {
        loadDetail()
      }
    },
  )

  watch(
    () => props.ossId,
    () => {
      if (props.open) {
        loadDetail()
      }
    },
  )

  async function loadDetail () {
    if (!props.ossId) {
      detail.value = null
      return
    }
    try {
      detail.value = await OssService.getOssComponent({ ossId: props.ossId })
    } catch (error) {
      console.error(error)
    }
  }

  function close () {
    emit('update:open', false)
  }
</script>
