<template>
  <v-container class="text-center" fluid>
    <v-row justify="center">
      <v-col cols="12">
        <h1>ページが見つかりません</h1>
        <v-btn color="primary" @click="goHome">トップに戻る</v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
  import { useRouter } from 'vue-router'

  const router = useRouter()

  function goHome () {
    router.push('/')
  }
</script>
