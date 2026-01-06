import { computed } from 'vue'
import {
  useConversationStore,
  useMessagesStore,
  useStreamingStore,
  useAuthStore
} from '@/stores'
import { conversationApi } from '@/api/modules/conversation'
import { useNotification } from '../core/useNotification'

export function useConversation() {
  const conversationStore = useConversationStore()
  const messagesStore = useMessagesStore()
  const streamingStore = useStreamingStore()
  const authStore = useAuthStore()
  const { error: showError } = useNotification()

  const hasConversation = computed(() => !!messagesStore.selectedConversation)
  const canSend = computed(() => {
    return !messagesStore.isSending && !streamingStore.isActive && messagesStore.input.trim().length > 0
  })

  async function loadConversations(appId: string, page = 1, limit = 20) {
    if (!authStore.token) return null
    messagesStore.isLoading = true
    try {
      return await conversationApi.getList(authStore.token, appId, page, limit)
    } catch (e) {
      showError((e as Error).message)
      return null
    } finally {
      messagesStore.isLoading = false
    }
  }

  async function loadMessages(conversationId: string, appId: string) {
    if (!authStore.token) return
    messagesStore.isLoading = true
    try {
      const result = await conversationApi.getDetail(authStore.token, appId, conversationId)
      messagesStore.selectConversation(result.conversation, appId)
      messagesStore.setMessages(result.messages, result.messages.length, false)
    } catch (e) {
      showError((e as Error).message)
    } finally {
      messagesStore.isLoading = false
    }
  }

  async function deleteConversation(conversationId: string, appId: string) {
    if (!authStore.token) return false
    try {
      await conversationApi.delete(authStore.token, appId, conversationId)
      if (messagesStore.selectedConversation?.id === conversationId) {
        messagesStore.clearConversation()
      }
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function renameConversation(conversationId: string, appId: string, newName: string) {
    if (!authStore.token) return false
    try {
      await conversationApi.rename(authStore.token, appId, conversationId, newName)
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  function stopStreaming() {
    streamingStore.stop()
  }

  return {
    hasConversation,
    canSend,
    loadConversations,
    loadMessages,
    deleteConversation,
    renameConversation,
    stopStreaming
  }
}
