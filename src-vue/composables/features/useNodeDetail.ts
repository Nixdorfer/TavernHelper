import { ref, computed } from 'vue'
import { worldTreeApi } from '@/api/modules/project'
import type {
  NodeTemplate,
  FolderTemplate,
  CardTemplate,
  BlockTemplate,
  ImmediateChange,
  SaveChanges,
  SaveFolderChange,
  SaveCardChange,
  SaveBlockChange
} from '@/types'
export interface PendingChanges {
  folder: Record<string, SaveFolderChange>
  card: Record<string, SaveCardChange>
  block: Record<string, SaveBlockChange>
}
function createEmptyPendingChanges(): PendingChanges {
  return {
    folder: {},
    card: {},
    block: {}
  }
}
export function useNodeDetail() {
  const nodeDetail = ref<NodeTemplate | null>(null)
  const isLoading = ref(false)
  const isSaving = ref(false)
  const pendingChanges = ref<PendingChanges>(createEmptyPendingChanges())
  const currentNodeId = ref<number | null>(null)
  const currentProjectId = ref<number | null>(null)
  const hasChanges = computed(() => {
    const p = pendingChanges.value
    return Object.keys(p.folder).length > 0 ||
           Object.keys(p.card).length > 0 ||
           Object.keys(p.block).length > 0
  })
  async function loadNodeDetail(nodeId: number) {
    if (isLoading.value) return
    isLoading.value = true
    try {
      nodeDetail.value = await worldTreeApi.getNodeDetail(nodeId)
      currentNodeId.value = nodeId
      pendingChanges.value = createEmptyPendingChanges()
    } finally {
      isLoading.value = false
    }
  }
  async function immediateAddFolder(parentFolderId: number | null, name: string): Promise<number> {
    if (!currentNodeId.value) throw new Error('No node selected')
    const change: ImmediateChange = {
      name,
      action: 'add',
      level: 'folder',
      target: parentFolderId
    }
    return await worldTreeApi.immediateChange(currentNodeId.value, change)
  }
  async function immediateDeleteFolder(folderId: number): Promise<number> {
    if (!currentNodeId.value) throw new Error('No node selected')
    const change: ImmediateChange = {
      name: '',
      action: 'del',
      level: 'folder',
      target: folderId
    }
    return await worldTreeApi.immediateChange(currentNodeId.value, change)
  }
  async function immediateAddCard(folderChangeId: number | null, name: string, keyWord: string = ''): Promise<number> {
    if (!currentNodeId.value) throw new Error('No node selected')
    const change: ImmediateChange = {
      name: JSON.stringify({ name, keyWord }),
      action: 'add',
      level: 'card',
      target: folderChangeId
    }
    return await worldTreeApi.immediateChange(currentNodeId.value, change)
  }
  async function immediateDeleteCard(cardId: number): Promise<number> {
    if (!currentNodeId.value) throw new Error('No node selected')
    const change: ImmediateChange = {
      name: '',
      action: 'del',
      level: 'card',
      target: cardId
    }
    return await worldTreeApi.immediateChange(currentNodeId.value, change)
  }
  async function immediateAddBlock(cardChangeId: number | null, title: string, zone: string): Promise<number> {
    if (!currentNodeId.value) throw new Error('No node selected')
    const change: ImmediateChange = {
      name: JSON.stringify({ title, zone }),
      action: 'add',
      level: 'block',
      target: cardChangeId
    }
    return await worldTreeApi.immediateChange(currentNodeId.value, change)
  }
  async function immediateDeleteBlock(blockId: number): Promise<number> {
    if (!currentNodeId.value) throw new Error('No node selected')
    const change: ImmediateChange = {
      name: '',
      action: 'del',
      level: 'block',
      target: blockId
    }
    return await worldTreeApi.immediateChange(currentNodeId.value, change)
  }
  function trackFolderChange(folderId: string, change: SaveFolderChange) {
    const existing = pendingChanges.value.folder[folderId] || {}
    pendingChanges.value.folder[folderId] = { ...existing, ...change }
  }
  function trackCardChange(cardId: string, change: SaveCardChange) {
    const existing = pendingChanges.value.card[cardId] || {}
    pendingChanges.value.card[cardId] = { ...existing, ...change }
  }
  function trackBlockChange(blockId: string, change: SaveBlockChange) {
    const existing = pendingChanges.value.block[blockId] || {}
    pendingChanges.value.block[blockId] = { ...existing, ...change }
  }
  async function saveChanges(): Promise<void> {
    if (!currentNodeId.value || !currentProjectId.value) {
      throw new Error('No node or project selected')
    }
    if (!hasChanges.value) return
    isSaving.value = true
    try {
      const changes: SaveChanges = {}
      if (Object.keys(pendingChanges.value.folder).length > 0) {
        changes.folder = pendingChanges.value.folder
      }
      if (Object.keys(pendingChanges.value.card).length > 0) {
        changes.card = pendingChanges.value.card
      }
      if (Object.keys(pendingChanges.value.block).length > 0) {
        changes.block = pendingChanges.value.block
      }
      await worldTreeApi.saveNodeChanges(currentNodeId.value, currentProjectId.value, changes)
      pendingChanges.value = createEmptyPendingChanges()
    } finally {
      isSaving.value = false
    }
  }
  function setProjectId(projectId: number) {
    currentProjectId.value = projectId
  }
  function clear() {
    nodeDetail.value = null
    currentNodeId.value = null
    currentProjectId.value = null
    pendingChanges.value = createEmptyPendingChanges()
  }
  function getFolder(folderId: string): FolderTemplate | null {
    if (!nodeDetail.value?.folders) return null
    return findFolderInStruct(nodeDetail.value.folders, folderId)
  }
  function findFolderInStruct(folders: Record<string, FolderTemplate>, folderId: string): FolderTemplate | null {
    if (folders[folderId]) return folders[folderId]
    for (const folder of Object.values(folders)) {
      if (folder.folders) {
        const found = findFolderInStruct(folder.folders, folderId)
        if (found) return found
      }
    }
    return null
  }
  function getCard(cardId: string): CardTemplate | null {
    if (!nodeDetail.value?.folders) return null
    return findCardInStruct(nodeDetail.value.folders, cardId)
  }
  function findCardInStruct(folders: Record<string, FolderTemplate>, cardId: string): CardTemplate | null {
    for (const folder of Object.values(folders)) {
      if (folder.cards && folder.cards[cardId]) {
        return folder.cards[cardId]
      }
      if (folder.folders) {
        const found = findCardInStruct(folder.folders, cardId)
        if (found) return found
      }
    }
    return null
  }
  function getBlock(blockId: string): BlockTemplate | null {
    if (!nodeDetail.value?.folders) return null
    return findBlockInStruct(nodeDetail.value.folders, blockId)
  }
  function findBlockInStruct(folders: Record<string, FolderTemplate>, blockId: string): BlockTemplate | null {
    for (const folder of Object.values(folders)) {
      if (folder.cards) {
        for (const card of Object.values(folder.cards)) {
          if (card.blocks && card.blocks[blockId]) {
            return card.blocks[blockId]
          }
        }
      }
      if (folder.folders) {
        const found = findBlockInStruct(folder.folders, blockId)
        if (found) return found
      }
    }
    return null
  }
  return {
    nodeDetail,
    isLoading,
    isSaving,
    hasChanges,
    pendingChanges,
    currentNodeId,
    currentProjectId,
    loadNodeDetail,
    immediateAddFolder,
    immediateDeleteFolder,
    immediateAddCard,
    immediateDeleteCard,
    immediateAddBlock,
    immediateDeleteBlock,
    trackFolderChange,
    trackCardChange,
    trackBlockChange,
    saveChanges,
    setProjectId,
    clear,
    getFolder,
    getCard,
    getBlock
  }
}
