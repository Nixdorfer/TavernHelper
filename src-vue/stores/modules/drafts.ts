import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Draft } from '@/types'
import { draftsApi } from '@/api/modules/drafts'
import type { ClipboardCapture } from '@/api/modules/drafts'

export const useDraftsStore = defineStore('drafts', () => {
  const drafts = ref<Draft[]>([])
  const clipboard = ref<ClipboardCapture[]>([])
  const isLoading = ref(false)
  const selectedDraftId = ref<string | null>(null)
  const searchQuery = ref('')
  const editingDraftId = ref<string | null>(null)

  const filteredDrafts = computed(() => {
    if (!searchQuery.value) return drafts.value
    const query = searchQuery.value.toLowerCase()
    return drafts.value.filter(draft =>
      draft.name?.toLowerCase().includes(query) ||
      draft.content?.toLowerCase().includes(query)
    )
  })

  const selectedDraft = computed(() => {
    return drafts.value.find(d => d.id === selectedDraftId.value) || null
  })

  async function loadDrafts() {
    isLoading.value = true
    try {
      drafts.value = await draftsApi.getAll()
    } finally {
      isLoading.value = false
    }
  }

  async function createDraft(draft: Omit<Draft, 'id' | 'createdAt' | 'updatedAt'>) {
    const newDraft = await draftsApi.create(draft)
    drafts.value.push(newDraft)
    return newDraft
  }

  async function updateDraft(draft: Draft) {
    await draftsApi.update(draft)
    const index = drafts.value.findIndex(d => d.id === draft.id)
    if (index > -1) {
      drafts.value[index] = draft
    }
  }

  async function deleteDraft(id: string) {
    await draftsApi.delete(id)
    drafts.value = drafts.value.filter(d => d.id !== id)
    if (selectedDraftId.value === id) {
      selectedDraftId.value = null
    }
  }

  async function loadClipboard() {
    clipboard.value = await draftsApi.getClipboard()
  }

  async function moveClipboardToDraft(captureId: string, name: string, parentId?: string) {
    const draft = await draftsApi.moveClipboardToDraft(captureId, name, parentId)
    drafts.value.push(draft)
    clipboard.value = clipboard.value.filter(c => c.id !== captureId)
    return draft
  }

  function selectDraft(draftId: string | null) {
    selectedDraftId.value = draftId
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function startEditing(draftId: string) {
    editingDraftId.value = draftId
  }

  function stopEditing() {
    editingDraftId.value = null
  }

  function clear() {
    drafts.value = []
    clipboard.value = []
    selectedDraftId.value = null
    searchQuery.value = ''
  }

  return {
    drafts,
    clipboard,
    isLoading,
    selectedDraftId,
    searchQuery,
    editingDraftId,
    filteredDrafts,
    selectedDraft,
    loadDrafts,
    createDraft,
    updateDraft,
    deleteDraft,
    loadClipboard,
    moveClipboardToDraft,
    selectDraft,
    setSearchQuery,
    startEditing,
    stopEditing,
    clear
  }
})
