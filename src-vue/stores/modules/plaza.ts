import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface PlazaApp {
  id: string
  name: string
  description: string
  icon?: string
  category?: string
}

export const usePlazaStore = defineStore('plaza', () => {
  const apps = ref<PlazaApp[]>([])
  const isLoading = ref(false)
  const searchQuery = ref('')
  const selectedCategory = ref<string | null>(null)
  const selectedAppId = ref<string | null>(null)

  const filteredApps = computed(() => {
    let result = apps.value
    if (selectedCategory.value) {
      result = result.filter(app => app.category === selectedCategory.value)
    }
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      result = result.filter(app =>
        app.name.toLowerCase().includes(query) ||
        app.description?.toLowerCase().includes(query)
      )
    }
    return result
  })

  const selectedApp = computed(() => {
    return apps.value.find(app => app.id === selectedAppId.value) || null
  })

  const categories = computed(() => {
    const cats = new Set<string>()
    apps.value.forEach(app => {
      if (app.category) cats.add(app.category)
    })
    return Array.from(cats)
  })

  function setApps(newApps: PlazaApp[]) {
    apps.value = newApps
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function setCategory(category: string | null) {
    selectedCategory.value = category
  }

  function selectApp(appId: string | null) {
    selectedAppId.value = appId
  }

  function clear() {
    apps.value = []
    searchQuery.value = ''
    selectedCategory.value = null
    selectedAppId.value = null
  }

  return {
    apps,
    isLoading,
    searchQuery,
    selectedCategory,
    selectedAppId,
    filteredApps,
    selectedApp,
    categories,
    setApps,
    setSearchQuery,
    setCategory,
    selectApp,
    clear
  }
})
