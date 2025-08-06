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
  import { useI18n } from 'vue-i18n'

  interface NavItem {
    titleKey: string
    icon: string
    to: RouteLocationRaw
  }

  const emit = defineEmits<{ (e: 'navigate'): void }>()

  const { t } = useI18n()

  const items: NavItem[] = [
    { titleKey: 'oss.listTitle', icon: 'mdi-package-variant', to: '/' },
    { titleKey: 'project.listTitle', icon: 'mdi-briefcase', to: '/projects' },
    { titleKey: 'settings.title', icon: 'mdi-cog', to: '/settings' },
  ]

  function onNavigate () {
    emit('navigate')
  }
</script>
