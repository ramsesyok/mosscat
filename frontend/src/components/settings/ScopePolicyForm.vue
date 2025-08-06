<template>
  <div>
    <v-switch v-model="form.runtimeRequiredDefaultInScope" :label="$t('settings.policy.runtimeRequiredDefaultInScope')" />
    <v-switch v-model="form.serverEnvIncluded" :label="$t('settings.policy.serverEnvIncluded')" />
    <v-switch v-model="form.autoMarkForksInScope" :label="$t('settings.policy.autoMarkForksInScope')" />
    <v-btn class="mt-4" color="primary" :loading="saving" @click="onSave">
      {{ $t('settings.policy.save') }}
    </v-btn>
  </div>
</template>

<script setup lang="ts">
  import type { ScopePolicy } from '@/api'
  import { storeToRefs } from 'pinia'
  import { reactive, ref, watch } from 'vue'
  import { useScopePolicyStore } from '@/stores/useScopePolicyStore'

  const store = useScopePolicyStore()
  const { policy } = storeToRefs(store)
  const emit = defineEmits<{ (e: 'saved', policy: ScopePolicy): void }>()

  const form = reactive({
    runtimeRequiredDefaultInScope: false,
    serverEnvIncluded: false,
    autoMarkForksInScope: false,
  })

  watch(policy, p => {
    if (p) {
      form.runtimeRequiredDefaultInScope = p.runtimeRequiredDefaultInScope ?? false
      form.serverEnvIncluded = p.serverEnvIncluded ?? false
      form.autoMarkForksInScope = p.autoMarkForksInScope ?? false
    }
  }, { immediate: true })

  const saving = ref(false)

  async function onSave () {
    saving.value = true
    try {
      const updated = await store.update({
        runtimeRequiredDefaultInScope: form.runtimeRequiredDefaultInScope,
        serverEnvIncluded: form.serverEnvIncluded,
        autoMarkForksInScope: form.autoMarkForksInScope,
      })
      emit('saved', updated)
    } catch (error) {
      console.error(error)
    } finally {
      saving.value = false
    }
  }
</script>
