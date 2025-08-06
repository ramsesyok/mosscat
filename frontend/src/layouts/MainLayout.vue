<template>
  <v-app>
    <v-navigation-drawer
      v-model="drawer"
      app
      :permanent="!display.mobile.value"
      :temporary="display.mobile.value"
    >
      <AppNavDrawer @navigate="onNavigate" />
    </v-navigation-drawer>

    <AppHeader @toggle-nav="toggleDrawer" />

    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
  import { ref, watch } from 'vue'
  import { useDisplay } from 'vuetify'

  import AppHeader from '@/components/common/AppHeader.vue'
  import AppNavDrawer from '@/components/common/AppNavDrawer.vue'

  const drawer = ref(true)
  const display = useDisplay()

  watch(() => display.mobile.value, val => {
    if (!val) drawer.value = true
  })

  function toggleDrawer () {
    drawer.value = !drawer.value
  }

  function onNavigate () {
    if (display.mobile.value) drawer.value = false
  }
</script>
