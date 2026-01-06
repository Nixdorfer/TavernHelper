import { defineStore } from 'pinia'
import { ref } from 'vue'
import { defaultWorldTreeSystemPrompt } from '@/constants/defaults'

export interface WorldTreeCommand {
  type: string
  target?: string
  content?: string
  raw: string
}

export interface PreviewDiff {
  add: any[]
  change: any[]
  del: any[]
}

export const useWorldTreeStore = defineStore('conversation-worldtree', () => {
  const tree = ref<any>(null)
  const isEnabled = ref(true)
  const systemPrompt = ref(defaultWorldTreeSystemPrompt)
  const showInSidebar = ref(false)
  const autoAcceptCommands = ref(false)
  const acceptedCommandRaws = ref<Set<string>>(new Set())
  const showProjectSelect = ref(false)
  const showCommandPreview = ref(false)
  const pendingCommand = ref<WorldTreeCommand | null>(null)
  const previewDiff = ref<PreviewDiff>({ add: [], change: [], del: [] })
  const previewAffectedEntryIds = ref<Set<string>>(new Set())

  function setTree(newTree: any) {
    tree.value = newTree
  }

  function clearTree() {
    tree.value = null
  }

  function toggleSidebar() {
    showInSidebar.value = !showInSidebar.value
  }

  function setEnabled(enabled: boolean) {
    isEnabled.value = enabled
  }

  function setSystemPrompt(prompt: string) {
    systemPrompt.value = prompt
  }

  function setAutoAccept(auto: boolean) {
    autoAcceptCommands.value = auto
  }

  function acceptCommand(raw: string) {
    acceptedCommandRaws.value.add(raw)
  }

  function clearAcceptedCommands() {
    acceptedCommandRaws.value.clear()
  }

  function setPendingCommand(command: WorldTreeCommand | null) {
    pendingCommand.value = command
    showCommandPreview.value = !!command
  }

  function setPreviewDiff(diff: PreviewDiff) {
    previewDiff.value = diff
  }

  function setAffectedEntryIds(ids: Set<string>) {
    previewAffectedEntryIds.value = ids
  }

  function clearPreview() {
    pendingCommand.value = null
    showCommandPreview.value = false
    previewDiff.value = { add: [], change: [], del: [] }
    previewAffectedEntryIds.value.clear()
  }

  return {
    tree,
    isEnabled,
    systemPrompt,
    showInSidebar,
    autoAcceptCommands,
    acceptedCommandRaws,
    showProjectSelect,
    showCommandPreview,
    pendingCommand,
    previewDiff,
    previewAffectedEntryIds,
    setTree,
    clearTree,
    toggleSidebar,
    setEnabled,
    setSystemPrompt,
    setAutoAccept,
    acceptCommand,
    clearAcceptedCommands,
    setPendingCommand,
    setPreviewDiff,
    setAffectedEntryIds,
    clearPreview
  }
})
