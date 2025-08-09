<template>
  <div>
    <div class="d-flex justify-end mb-2">
      <v-btn color="primary" @click="openCreate">
        {{ $t('oss.table.add') }}
      </v-btn>
    </div>
    <v-data-table
      item-value="id"
      :items="items"
      :items-length="totalItems"
      :items-per-page="itemsPerPage"
      :loading="loading"
      :page="page"
      @update:items-per-page="onUpdateItemsPerPage"
      @update:page="onUpdatePage"
    >
      <template #headers>
        <tr>
          <th class="text-left">{{ $t('oss.table.name') }}</th>
          <th v-if="display.xl.value" class="text-left" style="width: 25%">{{ $t('oss.table.homepage') }}</th>
          <th v-if="display.smAndUp.value" class="text-left" style="width: 25%">{{ $t('oss.table.repository') }}</th>
          <th v-if="display.mdAndUp.value" class="text-left" style="width: auto">{{ $t('oss.table.primaryLanguage') }}</th>
          <th class="text-left" style="width: auto">{{ $t('oss.table.layers') }}</th>
          <th v-if="display.xl.value" class="text-left" style="width: 20%">{{ $t('oss.table.tags') }}</th>
          <th class="text-left" style="width: auto">{{ $t('oss.table.actions') }}</th>
        </tr>
      </template>
      <template #item="{ item }">
        <tr>
          <td>{{ item.name }}</td>
          <td v-if="display.xl.value" style="width: 25%">
            <a
              v-if="item.homepageUrl"
              :href="item.homepageUrl"
              rel="noopener"
              target="_blank"
            >{{ item.homepageUrl }}</a>
          </td>
          <td v-if="display.smAndUp.value" style="width: 25%">
            <a
              v-if="item.repositoryUrl"
              :href="item.repositoryUrl"
              rel="noopener"
              target="_blank"
            >{{ item.repositoryUrl }}</a>
          </td>
          <td v-if="display.mdAndUp.value" style="width: auto">{{ item.primaryLanguage }}</td>
          <td style="width: auto">{{ (item.layers || []).join(', ') }}</td>
          <td v-if="display.xl.value" style="width: 20%">
            <v-chip
              v-for="tag in item.tags"
              :key="tag.id"
              class="ma-1"
              size="small"
              variant="tonal"
            >{{ tag.name }}</v-chip>
          </td>
          <td class="text-right" style="width: auto">
            <v-icon
              class="mr-2"
              size="small"
              @click.stop="openEdit(item.id)"
            >mdi-pencil</v-icon>
            <v-icon
              class="mr-2"
              size="small"
              @click.stop="openVersionList(item.id)"
            >mdi-format-list-bulleted</v-icon>
            <v-menu location="bottom">
              <template #activator="{ props }">
                <v-btn icon="mdi-dots-vertical" variant="text" v-bind="props" />
              </template>
              <v-list density="compact">
                <v-list-item @click="emit('delete', item)">
                  <v-list-item-title>{{ $t('common.delete') }}</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </td>
        </tr>
      </template>
    </v-data-table>
    <OssEditDialog
      v-model:open="editOpen"
      :oss-id="selectedId"
      @saved="onSaved"
    />
    <OssVersionListDialog
      v-model:open="versionListOpen"
      :oss-id="selectedId"
    />
  </div>
</template>

<script setup lang="ts">
  import type { OssComponent } from '@/api'
  import { onMounted, ref } from 'vue'
  import { useDisplay } from 'vuetify'
  import { OssService } from '@/api'
  import OssEditDialog from './OssEditDialog.vue'
  import OssVersionListDialog from './OssVersionListDialog.vue'

  const items = ref<OssComponent[]>([])
  const loading = ref(false)
  const page = ref(1)
  const itemsPerPage = ref(50)
  const totalItems = ref(0)

  const editOpen = ref(false)
  const versionListOpen = ref(false)
  const selectedId = ref<string>()
  const display = useDisplay()

  const emit = defineEmits<{
    (e: 'delete', item: OssComponent): void
  }>()

  async function fetchList () {
    loading.value = true
    try {
      const res = await OssService.listOssComponents({
        page: page.value,
        size: itemsPerPage.value,
      })
      items.value = res.items ?? []
      totalItems.value = res.total ?? 0
    } catch (error) {
      console.error(error)
    } finally {
      loading.value = false
    }
  }

  function onUpdatePage (p: number) {
    page.value = p
    fetchList()
  }

  function onUpdateItemsPerPage (s: number) {
    page.value = 1
    itemsPerPage.value = s
    fetchList()
  }

  function reload () {
    page.value = 1
    fetchList()
  }

  function openCreate () {
    selectedId.value = undefined
    editOpen.value = true
  }

  function openEdit (id: string) {
    selectedId.value = id
    editOpen.value = true
  }

  function openVersionList (id: string) {
    selectedId.value = id
    versionListOpen.value = true
  }

  async function onSaved () {
    await fetchList()
  }

  defineExpose({ reload })

  onMounted(fetchList)
</script>
