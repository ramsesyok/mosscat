import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

import LoginPage from '@/pages/LoginPage.vue'
import NotFoundPage from '@/pages/NotFoundPage.vue'
import OssListPage from '@/pages/OssListPage.vue'
import ProjectListPage from '@/pages/ProjectListPage.vue'
import SettingsPage from '@/pages/SettingsPage.vue'

import { RouteName } from '@/types/routes'
import { useAuthStore } from '@/stores/auth/useAuthStore'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: RouteName.Login,
    component: LoginPage,
    meta: { titleKey: 'login.title' },
  },
  {
    path: '/oss',
    name: RouteName.OssList,
    component: OssListPage,
    meta: { titleKey: 'oss.listTitle' },
  },
  {
    path: '/projects',
    name: RouteName.ProjectList,
    component: ProjectListPage,
    meta: { titleKey: 'project.listTitle' },
  },
  {
    path: '/settings',
    name: RouteName.Settings,
    component: SettingsPage,
    meta: { titleKey: 'settings.title' },
  },
  {
    path: '/:pathMatch(.*)*',
    name: RouteName.NotFound,
    component: NotFoundPage,
  },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior () {
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  const name = to.name as string | undefined
  if (name !== RouteName.Login && !auth.token) {
    return { name: RouteName.Login as any }
  }
  if (name === RouteName.Login && auth.token) {
    return { name: RouteName.OssList as any }
  }
})

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.('Failed to fetch dynamically imported module')) {
    if (localStorage.getItem('vuetify:dynamic-reload')) {
      console.error('Dynamic import error, reloading page did not fix it', err)
    } else {
      console.log('Reloading page to fix dynamic import error')
      localStorage.setItem('vuetify:dynamic-reload', 'true')
      location.assign(to.fullPath)
    }
  } else {
    console.error(err)
  }
})

router.isReady().then(() => {
  localStorage.removeItem('vuetify:dynamic-reload')
})

export default router
