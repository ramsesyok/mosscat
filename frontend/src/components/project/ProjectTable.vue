<template>
  <div>
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
      <template #top>
        <v-toolbar flat>
          <v-toolbar-title>{{ $t('project.listTitle') }}</v-toolbar-title>
          <v-spacer></v-spacer>
          <v-btn color="primary" @click="emit('create')">
            {{ $t('project.table.add') }}
          </v-btn>
        </v-toolbar>
      </template>
      <template #headers>
        <tr>
          <th class="text-left" style="width: 80px">{{ $t('project.table.code') }}</th>
          <th class="text-left">{{ $t('project.table.name') }}</th>
          <th class="text-left d-none d-md-table-cell" style="width: 140px">{{ $t('project.table.department') }}</th>
          <th class="text-left" style="width: 120px">{{ $t('project.table.deliveryDate') }}</th>
          <th class="text-left d-none d-md-table-cell" style="width: 140px">{{ $t('project.table.manager') }}</th>
          <th class="text-left" style="width: 120px">{{ $t('project.table.actions') }}</th>
        </tr>
      </template>
      <template #item="{ item }">
        <tr>
          <td style="width: 80px">{{ item.projectCode }}</td>
          <td>{{ item.name }}</td>
          <td class="d-none d-md-table-cell" style="width: 140px">{{ item.department }}</td>
          <td style="width: 120px">{{ item.deliveryDate }}</td>
          <td class="d-none d-md-table-cell" style="width: 140px">{{ item.manager }}</td>
          <td style="width: 120px">
            <v-btn icon size="small" variant="text" @click.stop="emit('detail', item)">
              <v-icon size="small">mdi-pencil</v-icon>
            </v-btn>
            <v-btn icon size="small" variant="text" @click.stop="emit('usage', item)">
              <v-icon size="small">mdi-view-list</v-icon>
            </v-btn>
            <v-menu location="bottom">
              <template #activator="{ props }">
                <v-btn icon size="small" variant="text" v-bind="props">
                  <v-icon size="small">mdi-dots-vertical</v-icon>
                </v-btn>
              </template>
              <v-list density="compact">
                <v-list-item @click="emitExport(item.id, 'csv')">
                  <v-list-item-title>{{ $t('export.csv') }}</v-list-item-title>
                </v-list-item>
                <v-list-item @click="emitExport(item.id, 'spdx-json')">
                  <v-list-item-title>{{ $t('export.spdx') }}</v-list-item-title>
                </v-list-item>
                <v-list-item @click="emit('delete', item)">
                  <v-list-item-title>{{ $t('common.delete') }}</v-list-item-title>
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
    (e: 'row-click' | 'detail' | 'usage' | 'delete', item: Project): void
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
