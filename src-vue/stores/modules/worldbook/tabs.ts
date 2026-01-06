import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface EditorTab {
  id: string
  nodeId: string
  name: string
  type: 'worldbook' | 'content' | 'note'
  isDirty: boolean
}

export const useWorldbookTabsStore = defineStore('worldbook-tabs', () => {
  const tabs = ref<EditorTab[]>([])
  const activeTabId = ref<string | null>(null)

  const activeTab = computed(() => {
    return tabs.value.find(t => t.id === activeTabId.value) || null
  })

  const hasDirtyTabs = computed(() => {
    return tabs.value.some(t => t.isDirty)
  })

  function addTab(tab: EditorTab) {
    const existing = tabs.value.find(t => t.id === tab.id)
    if (!existing) {
      tabs.value.push(tab)
    }
    activeTabId.value = tab.id
  }

  function removeTab(tabId: string) {
    const index = tabs.value.findIndex(t => t.id === tabId)
    if (index > -1) {
      tabs.value.splice(index, 1)
      if (activeTabId.value === tabId) {
        activeTabId.value = tabs.value[Math.max(0, index - 1)]?.id || null
      }
    }
  }

  function setActiveTab(tabId: string) {
    activeTabId.value = tabId
  }

  function markDirty(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.isDirty = true
    }
  }

  function markClean(tabId: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.isDirty = false
    }
  }

  function updateTabName(tabId: string, name: string) {
    const tab = tabs.value.find(t => t.id === tabId)
    if (tab) {
      tab.name = name
    }
  }

  function clear() {
    tabs.value = []
    activeTabId.value = null
  }

  return {
    tabs,
    activeTabId,
    activeTab,
    hasDirtyTabs,
    addTab,
    removeTab,
    setActiveTab,
    markDirty,
    markClean,
    updateTabName,
    clear
  }
})
