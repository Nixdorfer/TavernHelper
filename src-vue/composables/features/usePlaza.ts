import { computed } from 'vue'
import { usePlazaStore } from '@/stores'

export function usePlaza() {
  const store = usePlazaStore()

  const apps = computed(() => store.filteredApps)
  const categories = computed(() => store.categories)
  const selectedApp = computed(() => store.selectedApp)
  const isLoading = computed(() => store.isLoading)

  function setApps(apps: any[]) {
    store.setApps(apps)
  }

  function search(query: string) {
    store.setSearchQuery(query)
  }

  function filterByCategory(category: string | null) {
    store.setCategory(category)
  }

  function selectApp(appId: string | null) {
    store.selectApp(appId)
  }

  function clear() {
    store.clear()
  }

  return {
    apps,
    categories,
    selectedApp,
    isLoading,
    searchQuery: store.searchQuery,
    selectedCategory: store.selectedCategory,
    setApps,
    search,
    filterByCategory,
    selectApp,
    clear
  }
}
