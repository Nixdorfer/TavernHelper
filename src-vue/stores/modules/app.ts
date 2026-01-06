import { defineStore } from 'pinia'
import { ref } from 'vue'
export type PanelType = 'conversation' | 'worldbook' | 'plaza' | 'gallery' | 'creation' | 'drafts' | 'debug'
export const useAppStore = defineStore('app', () => {
  const activePanel = ref<PanelType>('worldbook')
  const isDevMode = ref(import.meta.env.DEV)
  const safeModeActive = ref(false)
  const plazaExpanded = ref(false)
  const plazaCurrentTab = ref<string | null>('favorite')
  const switchPanel = (panel: PanelType) => {
    activePanel.value = panel
    if (panel !== 'plaza') {
      plazaExpanded.value = false
    }
  }
  const togglePlazaExpand = () => {
    plazaExpanded.value = !plazaExpanded.value
    if (plazaExpanded.value) {
      activePanel.value = 'plaza'
      plazaCurrentTab.value = null
    }
  }
  const clickPlazaMain = () => {
    if (!plazaExpanded.value) plazaExpanded.value = true
    activePanel.value = 'plaza'
    plazaCurrentTab.value = null
  }
  const selectPlazaTab = (tabKey: string) => {
    plazaCurrentTab.value = tabKey
    activePanel.value = 'plaza'
  }
  return {
    activePanel,
    isDevMode,
    safeModeActive,
    plazaExpanded,
    plazaCurrentTab,
    switchPanel,
    togglePlazaExpand,
    clickPlazaMain,
    selectPlazaTab
  }
})
