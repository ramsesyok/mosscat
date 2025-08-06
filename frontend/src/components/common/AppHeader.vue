<template>
  <v-app-bar app>
    <v-app-bar-nav-icon @click="emit('toggle-nav')" />
    <v-toolbar-title class="text-h6">{{ title }}</v-toolbar-title>
    <div class="flex-grow-1" />
  </v-app-bar>
</template>

<script setup lang="ts">
  import { computed } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRoute } from 'vue-router'

  const emit = defineEmits<{ (e: 'toggle-nav'): void }>()

  const route = useRoute()
  const { t } = useI18n()

  const title = computed(() => {
    const key = route.meta.titleKey as string | undefined
    return key ? t(key) : 'OSSカタログ管理'
  })
</script>
