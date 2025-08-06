<template>
  <v-text-field
    v-model="search"
    append-inner-icon="mdi-magnify"
    clearable
    density="comfortable"
    :label="label"
    @click:append-inner="emitSearch"
    @keyup.enter="emitSearch"
  />
</template>

<script setup lang="ts">
  interface Props {
    name: string
    label?: string
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:name', value: string): void
    (e: 'search'): void
  }>()

  const search = computed({
    get: () => props.name,
    set: val => emit('update:name', val),
  })

  const label = computed(() => props.label ?? 'Search')

  function emitSearch () {
    emit('search')
  }
</script>
