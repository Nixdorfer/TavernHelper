<template>
  <div id="app" class="no-select" @contextmenu.prevent>
    <GlobalSidebar
      :token="authStore.token"
      :active-panel="appStore.activePanel"
      :plaza-expanded="appStore.plazaExpanded"
      :plaza-current-tab="appStore.plazaCurrentTab"
      :plaza-tabs="plazaTabs"
      :show-debug="appStore.isDevMode"
      @switch-panel="appStore.switchPanel"
      @click-plaza-main="appStore.clickPlazaMain"
      @select-plaza-tab="appStore.selectPlazaTab"
      @settings="showSettings = true"
    />
    <main class="app-container">
      <KeepAlive>
        <component :is="currentPage" v-if="currentPage" :key="appStore.activePanel" />
      </KeepAlive>
    </main>
    <LoginModal v-model="showLoginModal" @success="onLoginSuccess" />
    <SettingsModal v-model="showSettings" />
    <NotificationToast />
    <ConfirmDialog />
  </div>
</template>
<script setup>
import '@/styles/global.css'
import { ref, computed, defineAsyncComponent, onMounted } from 'vue'
import GlobalSidebar from './components/GlobalSidebar.vue'
import LoginModal from '@/components/modals/LoginModal.vue'
import NotificationToast from '@/components/common/NotificationToast.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { useAppStore, useAuthStore, useConfigStore } from '@/stores'
import { api } from '@/api'
import { logger } from '@/utils/logger'
const SettingsModal = defineAsyncComponent(() => import('@/components/modals/SettingsModal.vue'))
const ConversationView = defineAsyncComponent(() => import('./views/ConversationView.vue'))
const PlazaView = defineAsyncComponent(() => import('./views/PlazaView.vue'))
const DraftsView = defineAsyncComponent(() => import('./views/DraftsView.vue'))
const CreationView = defineAsyncComponent(() => import('./views/CreationView.vue'))
const GalleryView = defineAsyncComponent(() => import('./views/GalleryView.vue'))
const WorldbookView = defineAsyncComponent(() => import('./views/WorldbookView.vue'))
const DebugPanel = import.meta.env.DEV ? defineAsyncComponent(() => import('@/components/business/DebugPanel.vue')) : null
const appStore = useAppStore()
const authStore = useAuthStore()
const configStore = useConfigStore()
const showLoginModal = ref(false)
const showSettings = ref(false)
const pageComponents = {
  conversation: ConversationView,
  plaza: PlazaView,
  drafts: DraftsView,
  creation: CreationView,
  gallery: GalleryView,
  worldbook: WorldbookView,
  debug: DebugPanel
}
const currentPage = computed(() => {
  const panel = appStore.activePanel
  if (panel === 'debug' && !appStore.isDevMode) return null
  return pageComponents[panel] || null
})
const plazaTabs = computed(() => {
  const tabs = [
    { key: 'favorite', label: '我的收藏', requiresAuth: true },
    { key: 'recommended', label: '推荐榜' },
    { key: 'daily_rank', label: '日榜' },
    { key: 'weekly_rank', label: '周榜' },
    { key: 'monthly_rank', label: '月榜' },
    { key: 'overall_rank', label: '总榜' }
  ]
  return authStore.token ? tabs : tabs.filter(t => !t.requiresAuth)
})
const onLoginSuccess = async () => {
  showLoginModal.value = false
}
onMounted(async () => {
  try {
    const cfg = await api.config.load()
    if (cfg) {
      configStore.$patch(cfg)
    }
  } catch (e) {
    logger.error('加载配置失败:', e)
  }
})
</script>
<style scoped>
#app {
  display: flex;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-primary);
}
.app-container {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}
</style>
