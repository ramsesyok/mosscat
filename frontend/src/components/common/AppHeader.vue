<template>
  <v-app-bar app>
    <v-app-bar-nav-icon @click="emit('toggle-nav')" />
    <v-toolbar-title class="text-h6">{{ title }}</v-toolbar-title>
    <div class="flex-grow-1" />
    <v-btn
      v-if="auth.token"
      icon
      :title="t('common.logout')"
      @click="onLogout"
    >
      <v-icon icon="mdi-logout" />
    </v-btn>
  </v-app-bar>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRoute } from 'vue-router'
  import router from '@/router'
  import { useAuthStore } from '@/stores/auth/useAuthStore'
  import { RouteName } from '@/types/routes'

  const emit = defineEmits<{ (e: 'toggle-nav'): void }>()

  const route = useRoute()
  const { t } = useI18n()
  const auth = useAuthStore()

  const title = computed(() => {
    const key = route.meta.titleKey as string | undefined
    return key ? t(key) : 'OSSカタログ管理'
  })

  async function onLogout () {
    await auth.logout()
    router.push({ name: RouteName.Login })
  }
</script>
