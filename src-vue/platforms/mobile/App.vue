<template>
  <div class="app app--mobile" :class="{ 'app--dark': isDark }">
    <main class="app__main">
      <component :is="currentViewComponent" />
    </main>
    <BottomNav
      :currentView="appStore.currentView"
      @navigate="appStore.setView"
    />
    <ModalsContainer ref="modalsRef" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, defineAsyncComponent } from 'vue'
import { useAppStore, useAuthStore, useConfigStore } from '@/stores'
import { useTheme } from '@/composables'
import BottomNav from './components/BottomNav.vue'
import { ModalsContainer } from '@/components'
const ConversationView = defineAsyncComponent(() => import('./views/ConversationView.vue'))
const WorldbookView = defineAsyncComponent(() => import('./views/WorldbookView.vue'))
const PlazaView = defineAsyncComponent(() => import('./views/PlazaView.vue'))
const GalleryView = defineAsyncComponent(() => import('./views/GalleryView.vue'))
const appStore = useAppStore()
const authStore = useAuthStore()
const configStore = useConfigStore()
const { isDark, applyTheme } = useTheme()
const modalsRef = ref<InstanceType<typeof ModalsContainer> | null>(null)
const viewComponents = {
  conversation: ConversationView,
  worldbook: WorldbookView,
  plaza: PlazaView,
  gallery: GalleryView,
  creation: ConversationView,
  drafts: ConversationView
}
const currentViewComponent = computed(() => {
  return viewComponents[appStore.currentView] || ConversationView
})
onMounted(async () => {
  await configStore.load()
  applyTheme()
  appStore.setInitialized(true)
})
</script>

<style scoped>
.app--mobile {
  display: flex;
  flex-direction: column;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
  background: var(--bg-primary);
}
.app__main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>
