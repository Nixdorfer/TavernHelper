import { computed } from 'vue'
import { useWorldTreeStore } from '@/stores'

export function useWorldTree() {
  const store = useWorldTreeStore()

  const hasTree = computed(() => !!store.tree)
  const isEnabled = computed(() => store.isEnabled)
  const showInSidebar = computed(() => store.showInSidebar)

  function setTree(tree: any) {
    store.setTree(tree)
  }

  function clearTree() {
    store.clearTree()
  }

  function toggleSidebar() {
    store.toggleSidebar()
  }

  function setEnabled(enabled: boolean) {
    store.setEnabled(enabled)
  }

  function setSystemPrompt(prompt: string) {
    store.setSystemPrompt(prompt)
  }

  function acceptCommand(raw: string) {
    store.acceptCommand(raw)
  }

  function isCommandAccepted(raw: string): boolean {
    return store.acceptedCommandRaws.has(raw)
  }

  function setPendingCommand(command: any) {
    store.setPendingCommand(command)
  }

  function clearPreview() {
    store.clearPreview()
  }

  return {
    tree: store.tree,
    hasTree,
    isEnabled,
    showInSidebar,
    systemPrompt: store.systemPrompt,
    autoAcceptCommands: store.autoAcceptCommands,
    showCommandPreview: store.showCommandPreview,
    pendingCommand: store.pendingCommand,
    previewDiff: store.previewDiff,
    setTree,
    clearTree,
    toggleSidebar,
    setEnabled,
    setSystemPrompt,
    acceptCommand,
    isCommandAccepted,
    setPendingCommand,
    clearPreview
  }
}
