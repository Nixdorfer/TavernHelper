import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Revision {
  id: string
  messageId: string
  content: string
  type: 'edit' | 'regenerate' | 'branch'
  timestamp: number
}

export const useRevisionsStore = defineStore('conversation-revisions', () => {
  const list = ref<Revision[]>([])
  const emptyIds = ref<string[]>([])
  const showPending = ref(true)

  function add(revision: Revision) {
    list.value.push(revision)
  }

  function remove(index: number) {
    list.value.splice(index, 1)
  }

  function removeById(id: string) {
    const index = list.value.findIndex(r => r.id === id)
    if (index > -1) {
      list.value.splice(index, 1)
    }
  }

  function clear() {
    list.value = []
    emptyIds.value = []
  }

  function addEmptyId(id: string) {
    if (!emptyIds.value.includes(id)) {
      emptyIds.value.push(id)
    }
  }

  function removeEmptyId(id: string) {
    const index = emptyIds.value.indexOf(id)
    if (index > -1) {
      emptyIds.value.splice(index, 1)
    }
  }

  function setShowPending(show: boolean) {
    showPending.value = show
  }

  function getByMessageId(messageId: string): Revision[] {
    return list.value.filter(r => r.messageId === messageId)
  }

  return {
    list,
    emptyIds,
    showPending,
    add,
    remove,
    removeById,
    clear,
    addEmptyId,
    removeEmptyId,
    setShowPending,
    getByMessageId
  }
})
