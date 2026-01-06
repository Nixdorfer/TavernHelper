import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Message, Conversation } from '@/types'

export const useMessagesStore = defineStore('conversation-messages', () => {
  const selectedConversation = ref<Conversation | null>(null)
  const selectedConversationApp = ref<string | null>(null)
  const messages = ref<Message[]>([])
  const isLoading = ref(false)
  const hasMore = ref(false)
  const currentPage = ref(1)
  const total = ref(0)
  const input = ref('')
  const isSending = ref(false)
  const sendButtonState = ref<'ready' | 'sending' | 'streaming'>('ready')
  const canScrollDown = ref(false)
  const autoScrollLocked = ref(false)
  const editingMessageId = ref<string | null>(null)
  const editingMessageContent = ref('')
  const isSendingEdit = ref(false)
  const regeneratingMessageId = ref<string | null>(null)
  const generatingImageMessageId = ref<string | null>(null)
  const messageImages = ref<Record<string, string>>({})
  const pendingImageTasks = ref<Record<string, any>>({})
  const currentBranchId = ref<string | null>(null)

  const filteredMessages = computed(() => {
    return messages.value.filter(msg => !msg.isHidden)
  })

  const lastMessageHasError = computed(() => {
    if (messages.value.length === 0) return false
    const lastMsg = messages.value[messages.value.length - 1]
    return lastMsg && lastMsg.answer && lastMsg.answer.includes('[ERROR]')
  })

  function selectConversation(conversation: Conversation | null, app: string | null = null) {
    selectedConversation.value = conversation
    selectedConversationApp.value = app
    messages.value = []
    currentPage.value = 1
    hasMore.value = false
  }

  function clearConversation() {
    selectedConversation.value = null
    selectedConversationApp.value = null
    messages.value = []
  }

  function setMessages(newMessages: Message[], newTotal: number, newHasMore: boolean) {
    messages.value = newMessages
    total.value = newTotal
    hasMore.value = newHasMore
  }

  function appendMessages(newMessages: Message[]) {
    messages.value = [...newMessages, ...messages.value]
  }

  function addMessage(message: Message) {
    messages.value.push(message)
  }

  function updateMessage(id: string, updates: Partial<Message>) {
    const index = messages.value.findIndex(m => m.id === id)
    if (index > -1) {
      messages.value[index] = { ...messages.value[index], ...updates }
    }
  }

  function removeMessage(id: string) {
    const index = messages.value.findIndex(m => m.id === id)
    if (index > -1) {
      messages.value.splice(index, 1)
    }
  }

  function setInput(value: string) {
    input.value = value
  }

  function clearInput() {
    input.value = ''
  }

  function startEditing(messageId: string, content: string) {
    editingMessageId.value = messageId
    editingMessageContent.value = content
  }

  function cancelEditing() {
    editingMessageId.value = null
    editingMessageContent.value = ''
  }

  function setMessageImage(messageId: string, imageUrl: string) {
    messageImages.value[messageId] = imageUrl
  }

  return {
    selectedConversation,
    selectedConversationApp,
    messages,
    isLoading,
    hasMore,
    currentPage,
    total,
    input,
    isSending,
    sendButtonState,
    canScrollDown,
    autoScrollLocked,
    editingMessageId,
    editingMessageContent,
    isSendingEdit,
    regeneratingMessageId,
    generatingImageMessageId,
    messageImages,
    pendingImageTasks,
    currentBranchId,
    filteredMessages,
    lastMessageHasError,
    selectConversation,
    clearConversation,
    setMessages,
    appendMessages,
    addMessage,
    updateMessage,
    removeMessage,
    setInput,
    clearInput,
    startEditing,
    cancelEditing,
    setMessageImage
  }
})
