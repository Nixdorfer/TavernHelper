import { computed } from 'vue'
import { useConversationStore, useMessagesStore } from '@/stores'

export function useImageGeneration() {
  const conversationStore = useConversationStore()
  const messagesStore = useMessagesStore()

  const isAutoGenerateEnabled = computed(() => conversationStore.autoGenerateImage)
  const t2iPrompt = computed(() => conversationStore.t2iPrompt)
  const isGenerating = computed(() => !!messagesStore.generatingImageMessageId)

  function setAutoGenerate(enabled: boolean) {
    conversationStore.autoGenerateImage = enabled
  }

  function setT2iPrompt(prompt: string) {
    conversationStore.t2iPrompt = prompt
  }

  function showT2iSelector(messageId: string) {
    conversationStore.t2iMessageId = messageId
    conversationStore.showT2iSelectorModal = true
  }

  function hideT2iSelector() {
    conversationStore.showT2iSelectorModal = false
    conversationStore.t2iMessageId = null
  }

  function showImagePromptModal(messageId: string) {
    conversationStore.imagePromptMessageId = messageId
    conversationStore.showImagePromptModal = true
    conversationStore.imagePromptText = ''
  }

  function hideImagePromptModal() {
    conversationStore.showImagePromptModal = false
    conversationStore.imagePromptMessageId = null
    conversationStore.imagePromptText = ''
  }

  function getMessageImage(messageId: string): string | undefined {
    return messagesStore.messageImages[messageId]
  }

  return {
    isAutoGenerateEnabled,
    t2iPrompt,
    isGenerating,
    setAutoGenerate,
    setT2iPrompt,
    showT2iSelector,
    hideT2iSelector,
    showImagePromptModal,
    hideImagePromptModal,
    getMessageImage
  }
}
