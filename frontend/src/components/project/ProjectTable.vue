<template>
  <div>
    <div class="d-flex justify-end mb-2">
      <v-btn color="primary" @click="emit('create')">
        {{ $t('project.table.add') }}
      </v-btn>
    </div>
    <v-data-table
      item-value="id"
      :items="items"
      :items-length="totalItems"
      :items-per-page="itemsPerPage"
      :loading="loading"
      :page="page"
      @click:row="onRowClick"
      @update:items-per-page="s => emit('update:items-per-page', s)"
      @update:page="p => emit('update:page', p)"
    >
      <template #headers>
        <tr>
          <th class="text-left">{{ $t('project.table.code') }}</th>
          <th class="text-left">{{ $t('project.table.name') }}</th>
          <th class="text-left">{{ $t('project.table.department') }}</th>
          <th class="text-left">{{ $t('project.table.deliveryDate') }}</th>
          <th class="text-left">{{ $t('project.table.manager') }}</th>
          <th class="text-left" style="width: 120px">{{ $t('project.table.actions') }}</th>
        </tr>
      </template>
      <template #item="{ item }">
        <tr>
          <td>{{ item.projectCode }}</td>
          <td>{{ item.name }}</td>
          <td>{{ item.department }}</td>
          <td>{{ item.deliveryDate }}</td>
          <td>{{ item.manager }}</td>
          <td class="text-right" style="width: 120px">
            <v-btn icon="mdi-pencil" variant="text" @click.stop="emit('detail', item)" />
            <v-btn icon="mdi-view-list" variant="text" @click.stop="emit('usage', item)" />
            <v-menu location="bottom">
              <template #activator="{ props }">
                <v-btn icon="mdi-dots-vertical" variant="text" v-bind="props" />
              </template>
              <v-list density="compact">
                <v-list-item @click="emitExport(item.id, 'csv')">
                  <v-list-item-title>{{ $t('export.csv') }}</v-list-item-title>
                </v-list-item>
                <v-list-item @click="emitExport(item.id, 'spdx-json')">
                  <v-list-item-title>{{ $t('export.spdx') }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script setup lang="ts">
  import type { Project } from '@/api'

  interface Props {
    items: Project[]
    loading: boolean
    page: number
    itemsPerPage: number
    totalItems: number
  }

  defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:page' | 'update:items-per-page', value: number): void
    (e: 'row-click' | 'detail' | 'usage', item: Project): void
    (e: 'export', payload: { id: string, format: 'csv' | 'spdx-json' }): void
    (e: 'create'): void
  }>()

  function onRowClick (_: unknown, row: { item: Project }) {
    emit('row-click', row.item)
  }

  function emitExport (id: string, format: 'csv' | 'spdx-json') {
    emit('export', { id, format })
  }
</script>
