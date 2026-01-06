import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { WorldBookEntry, ContentItem } from '@/types'

export const useWorldbookEditorStore = defineStore('worldbook-editor', () => {
  const selectedNodeId = ref<string | null>(null)
  const selectedNode = ref<any>(null)
  const worldBook = ref<WorldBookEntry[]>([])
  const contents = ref<ContentItem[]>([])
  const isDirty = ref(false)
  const isEditing = ref(false)
  const editingEntryId = ref<string | null>(null)
  const searchQuery = ref('')
  const expandedEntryIds = ref<Set<string>>(new Set())

  const filteredWorldBook = computed(() => {
    if (!searchQuery.value) return worldBook.value
    const query = searchQuery.value.toLowerCase()
    return worldBook.value.filter(entry =>
      entry.name?.toLowerCase().includes(query) ||
      entry.content?.toLowerCase().includes(query)
    )
  })

  const filteredContents = computed(() => {
    if (!searchQuery.value) return contents.value
    const query = searchQuery.value.toLowerCase()
    return contents.value.filter(item =>
      item.name?.toLowerCase().includes(query) ||
      item.content?.toLowerCase().includes(query)
    )
  })

  function selectNode(nodeId: string | null, node: any = null) {
    selectedNodeId.value = nodeId
    selectedNode.value = node
  }

  function setWorldBook(entries: WorldBookEntry[]) {
    worldBook.value = entries
  }

  function setContents(items: ContentItem[]) {
    contents.value = items
  }

  function addEntry(entry: WorldBookEntry) {
    worldBook.value.push(entry)
    isDirty.value = true
  }

  function updateEntry(id: string, updates: Partial<WorldBookEntry>) {
    const index = worldBook.value.findIndex(e => e.id === id)
    if (index > -1) {
      worldBook.value[index] = { ...worldBook.value[index], ...updates }
      isDirty.value = true
    }
  }

  function removeEntry(id: string) {
    const index = worldBook.value.findIndex(e => e.id === id)
    if (index > -1) {
      worldBook.value.splice(index, 1)
      isDirty.value = true
    }
  }

  function addContent(item: ContentItem) {
    contents.value.push(item)
    isDirty.value = true
  }

  function updateContent(id: string, updates: Partial<ContentItem>) {
    const index = contents.value.findIndex(c => c.id === id)
    if (index > -1) {
      contents.value[index] = { ...contents.value[index], ...updates }
      isDirty.value = true
    }
  }

  function removeContent(id: string) {
    const index = contents.value.findIndex(c => c.id === id)
    if (index > -1) {
      contents.value.splice(index, 1)
      isDirty.value = true
    }
  }

  function toggleExpanded(entryId: string) {
    if (expandedEntryIds.value.has(entryId)) {
      expandedEntryIds.value.delete(entryId)
    } else {
      expandedEntryIds.value.add(entryId)
    }
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function startEditing(entryId: string) {
    editingEntryId.value = entryId
    isEditing.value = true
  }

  function stopEditing() {
    editingEntryId.value = null
    isEditing.value = false
  }

  function markClean() {
    isDirty.value = false
  }

  function clear() {
    selectedNodeId.value = null
    selectedNode.value = null
    worldBook.value = []
    contents.value = []
    isDirty.value = false
  }

  return {
    selectedNodeId,
    selectedNode,
    worldBook,
    contents,
    isDirty,
    isEditing,
    editingEntryId,
    searchQuery,
    expandedEntryIds,
    filteredWorldBook,
    filteredContents,
    selectNode,
    setWorldBook,
    setContents,
    addEntry,
    updateEntry,
    removeEntry,
    addContent,
    updateContent,
    removeContent,
    toggleExpanded,
    setSearchQuery,
    startEditing,
    stopEditing,
    markClean,
    clear
  }
})
