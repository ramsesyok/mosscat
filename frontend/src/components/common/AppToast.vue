<template>
  <v-snackbar
    v-model="visible"
    :color="color"
    location="bottom right"
    :timeout="timeout"
  >
    {{ message }}
  </v-snackbar>
</template>

<script setup lang="ts">
  import type { ToastEvent } from '@/types/toast'
  import { onMounted, onUnmounted, ref } from 'vue'
  import { offToast, onToast } from '@/composables/useToast'

  const visible = ref(false)
  const message = ref('')
  const color = ref('primary')
  const timeout = ref(3000)

  function handleToast (event: ToastEvent) {
    message.value = event.message
    timeout.value = event.timeout
    color.value = event.type === 'success'
      ? 'green'
      : (event.type === 'error'
        ? 'red'
        : 'primary')
    visible.value = true
  }

  onMounted(() => {
    onToast(handleToast)
  })

  onUnmounted(() => {
    offToast(handleToast)
  })
</script>
