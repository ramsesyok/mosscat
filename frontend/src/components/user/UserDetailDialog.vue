<template>
  <v-dialog v-model="modelOpen" max-width="600">
    <v-card>
      <v-card-title>{{ isNew ? $t('user.detail.titleNew') : $t('user.detail.titleEdit') }}</v-card-title>
      <v-card-text>
        <v-form ref="formRef">
          <v-text-field
            v-model="form.username"
            :disabled="!isNew"
            :label="$t('user.detail.username')"
            required
          />
          <v-text-field
            v-model="form.displayName"
            :label="$t('user.detail.displayName')"
          />
          <v-text-field
            v-model="form.email"
            :label="$t('user.detail.email')"
          />
          <v-text-field
            v-model="form.password"
            :label="$t('user.detail.password')"
            type="password"
          />
          <v-select
            v-model="form.roles"
            :items="roles"
            :label="$t('user.detail.roles')"
            multiple
          />
          <v-switch
            v-model="form.active"
            :label="$t('user.detail.active')"
          />
        </v-form>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn color="primary" :loading="saving" @click="onSave">{{ $t('user.detail.save') }}</v-btn>
        <v-btn text @click="close">{{ $t('user.detail.cancel') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
  import type { UserCreateRequest, UserUpdateRequest } from '@/api'
  import { computed, reactive, ref, watch } from 'vue'
  import { useUserStore } from '@/stores/useUserStore'

  interface Props {
    open: boolean
    userId?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{ (e: 'update:open', v: boolean): void, (e: 'saved'): void }>()

  const store = useUserStore()

  const modelOpen = computed({
    get: () => props.open,
    set: (v: boolean) => emit('update:open', v),
  })

  const isNew = computed(() => !props.userId)

  const roles = ['ADMIN', 'EDITOR', 'VIEWER']

  const form = reactive<UserCreateRequest>({
    username: '',
    displayName: undefined,
    email: undefined,
    password: undefined,
    roles: [],
    active: true,
  })

  const formRef = ref()
  const saving = ref(false)

  watch(() => props.open, val => {
    if (val) loadDetail()
  })
  watch(() => props.userId, () => {
    if (props.open) loadDetail()
  })

  async function loadDetail () {
    if (!props.userId) {
      resetForm()
      return
    }
    try {
      const detail = await store.get(props.userId)
      form.username = detail.username
      form.displayName = detail.displayName ?? undefined
      form.email = detail.email ?? undefined
      form.password = undefined
      form.roles = detail.roles.slice()
      form.active = detail.active
    } catch (error) {
      console.error(error)
    }
  }

  function resetForm () {
    form.username = ''
    form.displayName = undefined
    form.email = undefined
    form.password = undefined
    form.roles = []
    form.active = true
  }

  function close () {
    emit('update:open', false)
  }

  async function onSave () {
    saving.value = true
    try {
      if (isNew.value) {
        const payload: UserCreateRequest = { ...form }
        await store.create(payload)
      } else if (props.userId) {
        const payload: UserUpdateRequest = {
          displayName: form.displayName,
          email: form.email,
          password: form.password,
          roles: form.roles,
          active: form.active,
        }
        await store.update(props.userId, payload)
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
