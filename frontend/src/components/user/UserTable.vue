<template>
  <div>
    <div class="d-flex justify-end mb-2">
      <v-btn color="primary" @click="emit('create')">
        {{ $t('user.table.add') }}
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
          <th class="text-left">{{ $t('user.table.username') }}</th>
          <th class="text-left">{{ $t('user.table.displayName') }}</th>
          <th class="text-left">{{ $t('user.table.email') }}</th>
          <th class="text-left">{{ $t('user.table.roles') }}</th>
          <th class="text-left">{{ $t('user.table.active') }}</th>
          <th class="text-left" style="width: 120px">{{ $t('user.table.actions') }}</th>
        </tr>
      </template>
      <template #item="{ item }">
        <tr>
          <td>{{ item.username }}</td>
          <td>{{ item.displayName }}</td>
          <td>{{ item.email }}</td>
          <td>{{ item.roles?.join(', ') }}</td>
          <td>{{ item.active ? '✔' : '' }}</td>
          <td class="text-right" style="width: 120px">
            <v-btn icon="mdi-pencil" variant="text" @click.stop="emit('detail', item)" />
            <v-btn icon="mdi-delete" variant="text" @click.stop="emit('delete', item)" />
          </td>
        </tr>
      </template>
    </v-data-table>
  </div>
</template>

<script setup lang="ts">
  import type { User } from '@/api'

  interface Props {
    items: User[]
    loading: boolean
    page: number
    itemsPerPage: number
    totalItems: number
  }

  defineProps<Props>()
  const emit = defineEmits<{
    (e: 'update:page' | 'update:items-per-page', value: number): void
    (e: 'detail' | 'delete', item: User): void
    (e: 'create'): void
  }>()

  function onRowClick (_: unknown, row: { item: User }) {
    emit('detail', row.item)
  }
</script>
