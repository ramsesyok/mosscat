<template>
  <v-row class="align-center" dense>
    <v-col cols="12" sm="4">
      <v-text-field
        v-model="modelCode"
        clearable
        density="comfortable"
        :label="$t('project.filter.code')"
      />
    </v-col>
    <v-col cols="12" sm="4">
      <v-text-field
        v-model="modelName"
        clearable
        density="comfortable"
        :label="$t('project.filter.name')"
      />
    </v-col>
    <v-col cols="12" sm="2">
      <v-btn class="ma-0" color="primary" @click="emit('search')">
        {{ $t('project.filter.search') }}
      </v-btn>
    </v-col>
  </v-row>
</template>

<script setup lang="ts">
  interface Props {
    code: string | undefined
    name: string | undefined
  }

  const props = defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:code' | 'update:name', v: string | undefined): void
    (e: 'search'): void
  }>()

  const modelCode = computed({
    get: () => props.code,
    set: val => emit('update:code', val),
  })
  const modelName = computed({
    get: () => props.name,
    set: val => emit('update:name', val),
  })
</script>
