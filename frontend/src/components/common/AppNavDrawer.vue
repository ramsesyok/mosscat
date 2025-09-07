<template>
  <v-list density="compact" nav>
    <v-list-item
      v-for="item in items"
      :key="item.titleKey"
      :prepend-icon="item.icon"
      router
      :title="t(item.titleKey)"
      :to="item.to"
      @click="onNavigate"
    />
  </v-list>
</template>

<script setup lang="ts">
  import type { RouteLocationRaw } from 'vue-router'
  import { computed } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useAuthStore } from '@/stores/auth/useAuthStore'

  interface NavItem {
    titleKey: string
    icon: string
    to: RouteLocationRaw
  }

  const emit = defineEmits<{ (e: 'navigate'): void }>()

  const { t } = useI18n()
  const auth = useAuthStore()

  const items = computed<NavItem[]>(() => {
    const list: NavItem[] = [
      { titleKey: 'oss.listTitle', icon: 'mdi-package-variant', to: '/' },
      { titleKey: 'project.listTitle', icon: 'mdi-briefcase', to: '/projects' },
    ]
    if (auth.roles.includes('ADMIN')) {
      list.push({ titleKey: 'user.listTitle', icon: 'mdi-account-multiple', to: '/users' })
    }
    list.push({ titleKey: 'settings.title', icon: 'mdi-cog', to: '/settings' })
    return list
  })

  function onNavigate () {
    emit('navigate')
  }
</script>
