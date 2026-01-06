import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { defaultSystemPrompt, defaultT2iPrompt } from '@/constants/defaults'
import { useMessagesStore } from './messages'
import { useStreamingStore } from './streaming'
import { useWorldTreeStore } from './worldtree'
import { useOptionsStore } from './options'
import { useRevisionsStore } from './revisions'

export interface SelectionPopup {
  show: boolean
  x: number
  y: number
  text: string
  messageId: string | null
}

export interface ContextMenu {
  show: boolean
  x: number
  y: number
  fullAnswer: string
  hasSelection: boolean
}

export interface CapsuleDetail {
  show: boolean
  type: string
  text: string
  instruction: string
}

export const useConversationStore = defineStore('conversation', () => {
  const pendingConversationRestore = ref<any>(null)
  const newConversationStartPage = ref<any>(null)
  const newConvWorldTree = ref<any>(null)
  const newConvSelectingNode = ref(false)
  const newConvSelectedNodeId = ref<string | null>(null)
  const pendingNewConvProject = ref<string | null>(null)
  const newConvMessageInput = ref('')
  const sendingNewConvMessage = ref(false)
  const conversationPanelWidth = ref(320)
  const conversationPanelCollapsed = ref(false)
  const systemPrompt = ref(defaultSystemPrompt)
  const selectionPopup = ref<SelectionPopup>({
    show: false,
    x: 0,
    y: 0,
    text: '',
    messageId: null
  })
  const justShowedSelectionPopup = ref(false)
  const contextMenu = ref<ContextMenu>({
    show: false,
    x: 0,
    y: 0,
    fullAnswer: '',
    hasSelection: false
  })
  const capsuleDetail = ref<CapsuleDetail>({
    show: false,
    type: '',
    text: '',
    instruction: ''
  })
  const showMoreActions = ref(false)
  const showModSelector = ref(false)
  const modList = ref<any[]>([])
  const loadingMods = ref(false)
  const selectedModIds = ref<string[]>([])
  const savingMods = ref(false)
  const autoGenerateImage = ref(false)
  const t2iPrompt = ref(defaultT2iPrompt)
  const showT2iSelectorModal = ref(false)
  const t2iScenes = ref<any[]>([])
  const t2iMessageId = ref<string | null>(null)
  const showImagePromptModal = ref(false)
  const imagePromptText = ref('')
  const imagePromptMessageId = ref<string | null>(null)
  const submittingImagePrompt = ref(false)
  const imageProviders = ref<any[]>([])
  const selectedImageProviderId = ref('')
  const loadingImageProviders = ref(false)
  const pendingNewConv = ref<any>(null)
  const showParallelWorldConv = ref(false)
  const parallelWorldNewConv = ref<any>(null)
  const parallelWorldConvInput = ref('')
  const sendingParallelWorldMsg = ref(false)
  const showDeleteConfirm = ref(false)
  const isDebugTestConversation = ref(false)
  const showDebugReplyModal = ref(false)
  const syncingConfig = ref(false)

  const hasActiveSystemPrompts = computed(() => {
    const worldTreeStore = useWorldTreeStore()
    return !!(
      systemPrompt.value ||
      (worldTreeStore.isEnabled && worldTreeStore.systemPrompt) ||
      (autoGenerateImage.value && t2iPrompt.value)
    )
  })

  function setSystemPrompt(prompt: string) {
    systemPrompt.value = prompt
  }

  function showSelectionPopup(x: number, y: number, text: string, messageId: string | null) {
    selectionPopup.value = { show: true, x, y, text, messageId }
    justShowedSelectionPopup.value = true
  }

  function hideSelectionPopup() {
    selectionPopup.value.show = false
  }

  function showContextMenu(x: number, y: number, fullAnswer: string, hasSelection: boolean) {
    contextMenu.value = { show: true, x, y, fullAnswer, hasSelection }
  }

  function hideContextMenu() {
    contextMenu.value.show = false
  }

  function showCapsule(type: string, text: string, instruction: string) {
    capsuleDetail.value = { show: true, type, text, instruction }
  }

  function hideCapsule() {
    capsuleDetail.value.show = false
  }

  function setPanelWidth(width: number) {
    conversationPanelWidth.value = width
  }

  function togglePanelCollapsed() {
    conversationPanelCollapsed.value = !conversationPanelCollapsed.value
  }

  return {
    pendingConversationRestore,
    newConversationStartPage,
    newConvWorldTree,
    newConvSelectingNode,
    newConvSelectedNodeId,
    pendingNewConvProject,
    newConvMessageInput,
    sendingNewConvMessage,
    conversationPanelWidth,
    conversationPanelCollapsed,
    systemPrompt,
    selectionPopup,
    justShowedSelectionPopup,
    contextMenu,
    capsuleDetail,
    showMoreActions,
    showModSelector,
    modList,
    loadingMods,
    selectedModIds,
    savingMods,
    autoGenerateImage,
    t2iPrompt,
    showT2iSelectorModal,
    t2iScenes,
    t2iMessageId,
    showImagePromptModal,
    imagePromptText,
    imagePromptMessageId,
    submittingImagePrompt,
    imageProviders,
    selectedImageProviderId,
    loadingImageProviders,
    pendingNewConv,
    showParallelWorldConv,
    parallelWorldNewConv,
    parallelWorldConvInput,
    sendingParallelWorldMsg,
    showDeleteConfirm,
    isDebugTestConversation,
    showDebugReplyModal,
    syncingConfig,
    hasActiveSystemPrompts,
    setSystemPrompt,
    showSelectionPopup,
    hideSelectionPopup,
    showContextMenu,
    hideContextMenu,
    showCapsule,
    hideCapsule,
    setPanelWidth,
    togglePanelCollapsed
  }
})

export { useMessagesStore } from './messages'
export { useStreamingStore } from './streaming'
export { useWorldTreeStore } from './worldtree'
export { useOptionsStore } from './options'
export { useRevisionsStore } from './revisions'
