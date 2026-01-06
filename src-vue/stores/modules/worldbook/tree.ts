import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { TimelineNode } from '@/types'

export const useWorldbookTreeStore = defineStore('worldbook-tree', () => {
  const nodes = ref<TimelineNode[]>([])
  const expandedNodeIds = ref<Set<string>>(new Set())
  const selectedNodeId = ref<string | null>(null)
  const draggingNodeId = ref<string | null>(null)
  const dropTargetId = ref<string | null>(null)
  const searchQuery = ref('')
  const isRenaming = ref(false)
  const renamingNodeId = ref<string | null>(null)

  const filteredNodes = computed(() => {
    if (!searchQuery.value) return nodes.value
    const query = searchQuery.value.toLowerCase()
    return nodes.value.filter(node =>
      node.name?.toLowerCase().includes(query)
    )
  })

  const selectedNode = computed(() => {
    return nodes.value.find(n => n.id === selectedNodeId.value) || null
  })

  function setNodes(newNodes: TimelineNode[]) {
    nodes.value = newNodes
  }

  function addNode(node: TimelineNode) {
    nodes.value.push(node)
  }

  function updateNode(id: string, updates: Partial<TimelineNode>) {
    const index = nodes.value.findIndex(n => n.id === id)
    if (index > -1) {
      nodes.value[index] = { ...nodes.value[index], ...updates }
    }
  }

  function removeNode(id: string) {
    nodes.value = nodes.value.filter(n => n.id !== id && n.parentId !== id)
  }

  function selectNode(nodeId: string | null) {
    selectedNodeId.value = nodeId
  }

  function toggleExpanded(nodeId: string) {
    if (expandedNodeIds.value.has(nodeId)) {
      expandedNodeIds.value.delete(nodeId)
    } else {
      expandedNodeIds.value.add(nodeId)
    }
  }

  function expandNode(nodeId: string) {
    expandedNodeIds.value.add(nodeId)
  }

  function collapseNode(nodeId: string) {
    expandedNodeIds.value.delete(nodeId)
  }

  function expandAll() {
    nodes.value.forEach(n => expandedNodeIds.value.add(n.id))
  }

  function collapseAll() {
    expandedNodeIds.value.clear()
  }

  function startDragging(nodeId: string) {
    draggingNodeId.value = nodeId
  }

  function stopDragging() {
    draggingNodeId.value = null
    dropTargetId.value = null
  }

  function setDropTarget(nodeId: string | null) {
    dropTargetId.value = nodeId
  }

  function startRenaming(nodeId: string) {
    renamingNodeId.value = nodeId
    isRenaming.value = true
  }

  function stopRenaming() {
    renamingNodeId.value = null
    isRenaming.value = false
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  function getChildren(parentId: string): TimelineNode[] {
    return nodes.value.filter(n => n.parentId === parentId)
  }

  function clear() {
    nodes.value = []
    expandedNodeIds.value.clear()
    selectedNodeId.value = null
  }

  return {
    nodes,
    expandedNodeIds,
    selectedNodeId,
    draggingNodeId,
    dropTargetId,
    searchQuery,
    isRenaming,
    renamingNodeId,
    filteredNodes,
    selectedNode,
    setNodes,
    addNode,
    updateNode,
    removeNode,
    selectNode,
    toggleExpanded,
    expandNode,
    collapseNode,
    expandAll,
    collapseAll,
    startDragging,
    stopDragging,
    setDropTarget,
    startRenaming,
    stopRenaming,
    setSearchQuery,
    getChildren,
    clear
  }
})
