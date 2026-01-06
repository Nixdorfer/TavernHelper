import { computed } from 'vue'
import { useDraftsStore } from '@/stores'
import { useNotification } from '../core/useNotification'
import { useConfirm } from '../core/useConfirm'
import type { Draft } from '@/types'

export function useDrafts() {
  const store = useDraftsStore()
  const { success, error: showError } = useNotification()
  const { confirmDelete } = useConfirm()

  const drafts = computed(() => store.filteredDrafts)
  const clipboard = computed(() => store.clipboard)
  const selectedDraft = computed(() => store.selectedDraft)
  const isLoading = computed(() => store.isLoading)

  async function loadDrafts() {
    await store.loadDrafts()
  }

  async function loadClipboard() {
    await store.loadClipboard()
  }

  async function createDraft(draft: Omit<Draft, 'id' | 'createdAt' | 'updatedAt'>) {
    try {
      const newDraft = await store.createDraft(draft)
      success('草稿已创建')
      return newDraft
    } catch (e) {
      showError((e as Error).message)
      return null
    }
  }

  async function updateDraft(draft: Draft) {
    try {
      await store.updateDraft(draft)
      success('草稿已更新')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function deleteDraft(id: string) {
    const draft = store.drafts.find(d => d.id === id)
    const confirmed = await confirmDelete(draft?.name || '此草稿')
    if (!confirmed) return false
    try {
      await store.deleteDraft(id)
      success('草稿已删除')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function moveClipboardToDraft(captureId: string, name: string, parentId?: string) {
    try {
      const draft = await store.moveClipboardToDraft(captureId, name, parentId)
      success('已移至草稿')
      return draft
    } catch (e) {
      showError((e as Error).message)
      return null
    }
  }

  function selectDraft(draftId: string | null) {
    store.selectDraft(draftId)
  }

  function search(query: string) {
    store.setSearchQuery(query)
  }

  return {
    drafts,
    clipboard,
    selectedDraft,
    isLoading,
    searchQuery: store.searchQuery,
    loadDrafts,
    loadClipboard,
    createDraft,
    updateDraft,
    deleteDraft,
    moveClipboardToDraft,
    selectDraft,
    search
  }
}
