import { computed } from 'vue'
import {
  useWorldbookProjectStore,
  useWorldbookEditorStore,
  useWorldbookTabsStore
} from '@/stores'
import { useNotification } from '../core/useNotification'
import { nodeApi } from '@/api/modules/project'

export function useEditor() {
  const projectStore = useWorldbookProjectStore()
  const editorStore = useWorldbookEditorStore()
  const tabsStore = useWorldbookTabsStore()
  const { success, error: showError } = useNotification()

  const hasUnsavedChanges = computed(() => {
    return editorStore.isDirty || tabsStore.hasDirtyTabs
  })

  const currentNodeId = computed(() => editorStore.selectedNodeId)
  const currentNode = computed(() => editorStore.selectedNode)

  function selectNode(nodeId: string, node: any) {
    editorStore.selectNode(nodeId, node)
  }

  async function saveCurrentNode() {
    if (!projectStore.currentProjectName || !editorStore.selectedNodeId) return false
    try {
      await projectStore.saveProject()
      editorStore.markClean()
      if (tabsStore.activeTabId) {
        tabsStore.markClean(tabsStore.activeTabId)
      }
      success('保存成功')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function createChildNode(parentId: string, name: string) {
    if (!projectStore.currentProjectName) return null
    try {
      const node = await nodeApi.createChild(projectStore.currentProjectName, parentId, name)
      return node
    } catch (e) {
      showError((e as Error).message)
      return null
    }
  }

  async function createSiblingNode(siblingId: string, name: string) {
    if (!projectStore.currentProjectName) return null
    try {
      const node = await nodeApi.createBrother(projectStore.currentProjectName, siblingId, name)
      return node
    } catch (e) {
      showError((e as Error).message)
      return null
    }
  }

  async function renameNode(nodeId: string, newName: string) {
    try {
      await nodeApi.rename(nodeId, newName)
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function deleteNode(nodeId: string) {
    if (!projectStore.currentProjectName) return false
    try {
      await nodeApi.delete(projectStore.currentProjectName, nodeId)
      if (editorStore.selectedNodeId === nodeId) {
        editorStore.selectNode(null, null)
      }
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  return {
    hasUnsavedChanges,
    currentNodeId,
    currentNode,
    selectNode,
    saveCurrentNode,
    createChildNode,
    createSiblingNode,
    renameNode,
    deleteNode
  }
}
