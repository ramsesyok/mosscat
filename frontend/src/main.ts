/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Composables
import { createApp } from 'vue'

import { loadConfig } from '@/config'
// Plugins
import { registerPlugins } from '@/plugins'
import { configureHttp } from '@/plugins/http'
import { i18n } from '@/plugins/i18n'
import router from '@/router'
import pinia from '@/stores'

// Components
import App from './App.vue'

// Styles
import 'unfonts.css'

async function bootstrap () {
  const cfg = await loadConfig()
  configureHttp(cfg)

  const app = createApp(App)

  app.use(pinia)
  registerPlugins(app)
  app.use(router)
  app.use(i18n)
  app.mount('#app')
}

bootstrap()
