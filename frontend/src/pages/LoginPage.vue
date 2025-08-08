<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" md="4" sm="8">
        <v-card>
          <v-card-title>{{ t('login.title') }}</v-card-title>
          <v-card-text>
            <v-form @submit.prevent="onSubmit">
              <v-text-field
                v-model="username"
                class="mb-2"
                :label="t('login.username')"
                required
              />
              <v-text-field
                v-model="password"
                class="mb-4"
                :label="t('login.password')"
                required
                type="password"
              />
              <v-btn color="primary" :loading="loading" type="submit">
                {{ t('login.submit') }}
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import { useI18n } from 'vue-i18n'
  import { useRouter } from 'vue-router'
  import { useAuthStore } from '@/stores/auth/useAuthStore'
  import { RouteName } from '@/types/routes'
  import { useToast } from '@/composables/useToast'

  const username = ref('')
  const password = ref('')
  const loading = ref(false)
  const router = useRouter()
  const { t } = useI18n()
  const auth = useAuthStore()
  const toast = useToast()

  async function onSubmit () {
    loading.value = true
    try {
      await auth.login(username.value, password.value)
      router.push({ name: RouteName.OssList as any })
    } catch {
      toast.error(t('login.failed'))
    } finally {
      loading.value = false
    }
  }
</script>
