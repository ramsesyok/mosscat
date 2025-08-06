/// <reference types="vite/client" />
/// <reference types="unplugin-vue-router/client" />
/// <reference types="vite-plugin-vue-layouts-next/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_BASIC_USER: string
  readonly VITE_BASIC_PASS: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
