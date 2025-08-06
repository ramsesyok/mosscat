<template>
  <div>
    <v-row class="align-center" dense>
      <v-col cols="8" sm="6">
        <v-text-field
          v-model="newTag"
          density="comfortable"
          hide-details
          :label="$t('settings.tags.newTag')"
          @keyup.enter="onAdd"
        />
      </v-col>
      <v-col cols="4" sm="2">
        <v-btn class="ma-0" color="primary" :loading="adding" @click="onAdd">
          {{ $t('settings.tags.add') }}
        </v-btn>
      </v-col>
    </v-row>
    <div class="my-2">
      <v-chip
        v-for="tag in items"
        :key="tag.id"
        class="ma-1"
        closable
        size="small"
        variant="tonal"
        @click:close="openConfirm(tag)"
      >{{ tag.name }}</v-chip>
    </div>
    <v-dialog v-model="confirmOpen" max-width="400">
      <v-card>
        <v-card-title>{{ $t('settings.tags.confirmDelete') }}</v-card-title>
        <v-card-actions>
          <v-spacer />
          <v-btn color="primary" :loading="deleting" @click="onDelete">
            {{ $t('settings.policy.save') }}
          </v-btn>
          <v-btn text @click="confirmOpen = false">
            {{ $t('project.detail.cancel') }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    <v-snackbar v-model="snackbar" timeout="3000">{{ snackbarMessage }}</v-snackbar>
  </div>
</template>

<script setup lang="ts">
  import type { Tag } from '@/api'
  import { storeToRefs } from 'pinia'
  import { onMounted, ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useTagStore } from '@/stores/useTagStore'

  const store = useTagStore()
  const { items } = storeToRefs(store)
  const { t } = useI18n()

  const newTag = ref('')
  const adding = ref(false)
  const deleting = ref(false)
  const confirmOpen = ref(false)
  const target = ref<Tag>()
  const snackbar = ref(false)
  const snackbarMessage = ref('')

  onMounted(() => {
    store.fetch()
  })

  async function onAdd () {
    if (!newTag.value) return
    adding.value = true
    try {
      await store.create(newTag.value)
      await store.fetch()
      newTag.value = ''
      showToast(t('settings.tags.added'))
    } catch (error) {
      console.error(error)
    } finally {
      adding.value = false
    }
  }

  function openConfirm (tag: Tag) {
    target.value = tag
    confirmOpen.value = true
  }

  async function onDelete () {
    if (!target.value) return
    deleting.value = true
    try {
      await store.remove(target.value.id)
      await store.fetch()
      showToast(t('settings.tags.deleted'))
    } catch (error) {
      console.error(error)
    } finally {
      deleting.value = false
      confirmOpen.value = false
      target.value = undefined
    }
  }

  function showToast (msg: string) {
    snackbarMessage.value = msg
    snackbar.value = true
  }
</script>
