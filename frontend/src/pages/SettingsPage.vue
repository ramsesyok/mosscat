<template>
  <v-container fluid>
    <v-row>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>{{ $t('settings.tags.title') }}</v-card-title>
          <v-card-text>
            <TagManager />
          </v-card-text>
        </v-card>
      </v-col>
      <v-col cols="12" md="6">
        <v-card>
          <v-card-title>{{ $t('settings.policy.title') }}</v-card-title>
          <v-card-text>
            <ScopePolicyForm @saved="onPolicySaved" />
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    <v-snackbar v-model="snackbar" timeout="3000">{{ snackbarMessage }}</v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
  import { onMounted, ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import ScopePolicyForm from '@/components/settings/ScopePolicyForm.vue'
  import TagManager from '@/components/settings/TagManager.vue'
  import { useScopePolicyStore } from '@/stores/useScopePolicyStore'
  import { useTagStore } from '@/stores/useTagStore'

  const tagStore = useTagStore()
  const policyStore = useScopePolicyStore()
  const { t } = useI18n()

  const snackbar = ref(false)
  const snackbarMessage = ref('')

  onMounted(() => {
    tagStore.fetch()
    policyStore.fetch()
  })

  function onPolicySaved () {
    snackbarMessage.value = t('settings.policy.saved')
    snackbar.value = true
  }
</script>
