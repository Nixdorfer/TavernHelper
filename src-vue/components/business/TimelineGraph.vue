<template>
  <div :class="['timeline-graph', { 'readonly-mode': readonly, 'full-mode': fullMode }]">
    <div v-if="!readonly && !fullMode && visibleBranchList.length > 0" class="branch-headers" :style="{ width: svgWidth + 'px' }">
      <div class="branch-headers-list">
        <div
          v-for="branch in visibleBranchList"
          :key="branch.id"
          class="branch-header"
          :style="{ width: actualColumnWidth + 'px', borderTopColor: branch.color }"
          @click="startEditBranchName(branch)"
        >
          <span class="branch-color-dot" :style="{ background: branch.color }"></span>
          <span class="branch-name">{{ getBranchDisplayName(branch) }}</span>
        </div>
      </div>
      <button class="btn-branch-selector" @click.stop="openBranchSelector" title="选择显示的世界线">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M4 6h16M4 12h16M4 18h16"/>
        </svg>
      </button>
    </div>
    <div class="graph-container" ref="graphContainer" :key="graphRenderKey" :style="{ width: svgWidth + 'px' }" @contextmenu.prevent="!readonly && showGraphContextMenu($event)" @mousedown="handlePanStart" @mousemove="handlePanMove" @mouseup="handlePanEnd" @mouseleave="handlePanEnd" @wheel="handleWheel">
      <div class="graph-content" :style="{ minHeight: svgHeight + 'px', width: svgWidth + 'px', transform: 'scale(' + scale + ')', transformOrigin: 'top left' }">
        <svg class="graph-lines" :style="{ width: svgWidth + 'px', height: svgHeight + 'px' }">
          <path
            v-for="line in lines"
            :key="line.id"
            :d="line.path"
            class="connection-line"
            :class="{
              highlighted: isNodeHighlighted(line.fromId) && isNodeHighlighted(line.toId),
              'logic-line': line.isLogicConnection
            }"
            :style="{ stroke: line.isLogicConnection ? '#f97316' : line.color }"
          />
        </svg>
        <div
          v-for="capsule in logicCapsules"
          :key="capsule.id"
          class="logic-capsule"
          :style="{ left: capsule.x + 'px', top: capsule.y + 'px' }"
          @click="showLogicCapsuleDetail(capsule)"
        >
          {{ capsule.label }}
        </div>
        <div
          v-for="node in displayNodes"
          :key="node.id"
          class="node-wrapper"
          :style="{ left: node.x + 'px', top: (node.y + getNodeOffset(node.id)) + 'px', width: actualColumnWidth + 'px' }"
          @mouseenter="!readonly && setHoveredNode(node.id, $event)"
          @mouseleave="!readonly && clearHoveredNode()"
        >
          <div class="node-shape-container">
            <div
              :class="['node-shape', {
                'node-capsule': !node.isLogic && !node.isHiddenBranch,
                'node-hexagon': node.isLogic,
                'node-wide-hexagon': node.isHiddenBranch,
                current: node.id === currentNodeId && !node.isHiddenBranch,
                selected: !readonly && node.id === selectedNodeId && !node.isHiddenBranch,
                filtered: searchTag && !isNodeHighlighted(node.id),
                hovered: !readonly && hoveredNodeId === node.id && !node.isHiddenBranch,
                'has-conversation': nodeConversationMap?.[node.id] && !node.isHiddenBranch,
                'has-error': nodeErrorMap?.[node.id] && !node.isHiddenBranch,
                'hidden-branch-node': node.isHiddenBranch
              }]"
              :style="!node.isLogic ? { borderColor: node.branchColor, '--node-border-color': node.branchColor } : {}"
              :draggable="!readonly && !node.isHiddenBranch"
              @click="!readonly && !node.isHiddenBranch && selectAndOpenNode(node)"
              @dblclick="!readonly && !node.isHiddenBranch && handleNodeDblClick(node)"
              @contextmenu.prevent="((!readonly || allowParallelWorld) && !node.isHiddenBranch) && showContextMenu($event, node)"
              @dragstart="!readonly && !node.isHiddenBranch && handleNodeDragStart($event, node)"
              @drag="!readonly && !node.isHiddenBranch && handleNodeDrag($event, node)"
              @dragend="!readonly && !node.isHiddenBranch && handleNodeDragEnd($event, node)"
            >
              <span v-if="nodeErrorMap?.[node.id] && !node.isHiddenBranch" class="error-name-capsule" :data-node-id="node.id">{{ node.name }}</span>
              <span v-else class="node-shape-name">{{ node.name }}</span>
            </div>
            <transition v-if="!readonly && !node.isHiddenBranch && !selectOnly" name="fade">
              <div v-if="hoveredNodeId === node.id" class="node-action-buttons">
                <button class="add-btn add-child" @click.stop="emit('create-child', node)" title="新建后续事件">+</button>
                <button v-if="fullMode" class="add-btn edit-btn" @click.stop="handleFullModeEdit(node)" title="编辑节点">✎</button>
              </div>
            </transition>
          </div>
          <div class="node-tags" v-if="node.tags && node.tags.length && !node.isHiddenBranch">
            <span v-for="tag in node.tags" :key="tag" class="tag">#{{ tag }}</span>
          </div>
          <div
            v-if="draggingNodeId === node.id && snapPreviewOffset !== null && snapPreviewOffset !== getNodeOffset(node.id)"
            class="node-snap-ghost"
            :style="{ top: (snapPreviewOffset - getNodeOffset(node.id)) + 'px' }"
          >
            <div class="ghost-shape"></div>
          </div>
          <div
            v-if="!readonly && hoveredNodeId === node.id && nodeConversationMap?.[node.id] && nodeConversationNames?.[node.id]?.length"
            class="node-conversations-tooltip"
          >
            <div class="tooltip-title">绑定的对话</div>
            <div v-for="conv in (nodeConversationNames?.[node.id] || [])" :key="conv.id" class="tooltip-item">{{ conv.name || conv.id }}</div>
          </div>
        </div>
      </div>
    </div>
    <div v-if="contextMenu.show" class="context-menu" :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }">
      <div v-if="contextMenu.isGraphMenu && !readonly" class="menu-item" @click="createRootNode">新建根节点</div>
      <template v-else-if="allowParallelWorld && readonly">
        <div class="menu-item" @click="createParallelWorld">平行世界对话</div>
      </template>
      <template v-else-if="!readonly">
        <div class="menu-item danger" @click="deleteNodeAndRelink">删除事件节点</div>
      </template>
    </div>
    <div v-if="contextMenu.show" class="menu-overlay" @click="closeContextMenu"></div>
    <div v-if="operationMode" class="operation-hint">
      <span>{{ operationMode === 'merge' ? '选择要合并到的目标事件' : '选择新的父事件' }}</span>
      <button class="btn btn-secondary btn-sm" @click="cancelOperation">取消</button>
    </div>
    <div v-if="showBranchSelector" class="modal-overlay" @click.self="showBranchSelector = false">
      <div class="branch-selector-modal">
        <div class="modal-header">
          <h3>选择显示的世界线（最多3条）</h3>
          <button class="btn btn-icon" @click="showBranchSelector = false">×</button>
        </div>
        <div class="modal-body">
          <div class="branch-options-grid">
            <div
              v-for="(branch, index) in selectableBranches"
              :key="branch.id"
              :class="['branch-option-capsule', { selected: pendingBranchIds.includes(branch.id), required: branch.id === requiredBranchId }]"
              @click="togglePendingBranchVisibility(branch.id)"
            >
              <span class="branch-color-dot" :style="{ background: branch.color }"></span>
              <span class="branch-name">{{ branch.name || '世界线 ' + (index + 1) }}</span>
              <span v-if="branch.id === requiredBranchId" class="required-badge">当前</span>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button :class="['btn', pendingVisibleBranchCount > 3 ? 'btn-danger' : 'btn-primary']" :disabled="pendingVisibleBranchCount < 1" @click="confirmBranchSelection">确定</button>
        </div>
      </div>
    </div>
    <div v-if="showBranchEditModal" class="modal-overlay" @click.self="showBranchEditModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>编辑世界线</h3>
          <button class="btn btn-icon" @click="showBranchEditModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>世界线名称</label>
            <input type="text" class="form-control" v-model="editingBranchName" placeholder="输入世界线名称">
          </div>
          <div class="form-group" style="margin-top: 12px;">
            <label>世界线颜色</label>
            <div class="color-picker-grid">
              <div
                v-for="(color, index) in BRANCH_COLORS"
                :key="index"
                class="color-option"
                :class="{ selected: editingBranchColor === color }"
                :style="{ background: color }"
                @click="editingBranchColor = color"
              ></div>
            </div>
          </div>
          <div class="form-group" style="margin-top: 12px;">
            <label>世界线描述</label>
            <textarea class="form-control" v-model="editingBranchDescription" rows="3" placeholder="输入世界线描述（可选）"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showBranchEditModal = false">取消</button>
          <button class="btn btn-primary" @click="saveBranchEdit">保存</button>
        </div>
      </div>
    </div>
    <div v-if="showNoteModal" class="modal-overlay" @click.self="showNoteModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>{{ editingNoteNode?.name }}</h3>
          <button class="btn btn-icon" @click="showNoteModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>事件备注</label>
            <textarea class="form-control note-textarea" v-model="editingNote" placeholder="添加节点备注..." rows="6"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showNoteModal = false">取消</button>
          <button class="btn btn-primary" @click="saveNote">保存</button>
        </div>
      </div>
    </div>
    <div v-if="showLogicCapsuleModal" class="modal-overlay" @click.self="showLogicCapsuleModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>编辑分支标记</h3>
          <button class="btn btn-icon" @click="showLogicCapsuleModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>标记名称</label>
            <input type="text" class="form-control" v-model="editingCapsuleLabel" placeholder="输入分支标记名称">
          </div>
          <div class="form-group">
            <label>描述</label>
            <textarea class="form-control" v-model="editingCapsuleDescription" rows="3" placeholder="输入分支描述（可选）"></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showLogicCapsuleModal = false">取消</button>
          <button class="btn btn-primary" @click="saveCapsuleEdit">保存</button>
        </div>
      </div>
    </div>
    <Teleport to="body">
      <div v-if="hoveredNodeId && nodeErrorMap?.[hoveredNodeId]" class="node-error-tooltip-fixed" :style="errorTooltipStyle">
        <div v-for="(err, idx) in (nodeErrorMap?.[hoveredNodeId] || [])" :key="idx" class="error-item">
          <span class="error-entry">{{ err.entryName }}</span>: {{ err.error }}
        </div>
      </div>
    </Teleport>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
const BRANCH_COLORS = [
  '#3b82f6', '#f59e0b', '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16',
  '#f97316', '#6366f1', '#14b8a6', '#eab308', '#a855f7', '#0ea5e9', '#65a30d',
  '#fb923c', '#818cf8', '#2dd4bf', '#fbbf24', '#c084fc', '#38bdf8', '#a3e635',
  '#fdba74', '#a5b4fc', '#5eead4', '#fde047', '#d8b4fe', '#7dd3fc', '#bef264'
]
const ROOT_NODE_COLOR = '#ffffff'
const COLUMN_WIDTH = 80
const LEVEL_HEIGHT = 80
const PADDING = 20
interface TimelineNode {
  id: number
  name: string
  parentId?: number | null
  branchId?: string
  isLogic?: boolean
  tags?: string[]
  note?: string
  optionLabel?: string
}
interface Branch {
  id: number
  name: string
  color: string
  description?: string
  nodeCount: number
  columnIndex: number
}
interface ContextMenuState {
  show: boolean
  x: number
  y: number
  node: TimelineNode | null
  isGraphMenu: boolean
}
interface LogicCapsule {
  id: string
  x: number
  y: number
  label: string
  description: string
  nodeId: number
}
const props = defineProps<{
  nodes: TimelineNode[]
  currentNodeId?: number | null
  branchNames?: Record<string, { name?: string; color?: string; description?: string }>
  readonly?: boolean
  nodeConversationMap?: Record<number, boolean>
  nodeConversationNames?: Record<number, { id: string; name: string }[]>
  requiredBranchId?: number | null
  nodeErrorMap?: Record<number, { entryName: string; error: string }[]>
  fullMode?: boolean
  pathMode?: boolean
  selectOnly?: boolean
  allowParallelWorld?: boolean
}>()
const emit = defineEmits<{
  select: [node: TimelineNode]
  'set-current': [node: TimelineNode]
  'create-child': [node: TimelineNode]
  'create-logic': []
  'create-branch': []
  'create-root': []
  rename: [node: TimelineNode]
  'edit-tags': [node: TimelineNode]
  merge: [data: { source: TimelineNode; target: TimelineNode }]
  rebase: [data: { node: TimelineNode; newParent: TimelineNode }]
  delete: [node: TimelineNode]
  'delete-and-relink': [node: TimelineNode]
  'update-branch': [data: { branchId: number; name: string; color: string; description: string }]
  'update-note': [nodeId: number, label: string, description: string]
  'width-change': [width: number]
  'context-menu-open': []
  'create-parallel-world': [node: TimelineNode]
  'node-click': [node: TimelineNode]
  'full-mode-edit': [node: TimelineNode]
}>()
const graphContainer = ref<HTMLElement | null>(null)
const selectedNodeId = ref<number | null>(null)
const hoveredNodeId = ref<number | null>(null)
const searchTag = ref('')
const filteredNodeIds = ref<Set<number>>(new Set())
const displayBranchIds = ref<number[]>([])
const pendingBranchIds = ref<number[]>([])
const showBranchSelector = ref(false)
const showBranchEditModal = ref(false)
const showNoteModal = ref(false)
const showLogicCapsuleModal = ref(false)
const selectedLogicCapsule = ref<LogicCapsule | null>(null)
const editingCapsuleLabel = ref('')
const editingCapsuleDescription = ref('')
const editingBranchId = ref<number | null>(null)
const editingBranchName = ref('')
const editingBranchColor = ref('')
const editingBranchDescription = ref('')
const editingNoteNode = ref<TimelineNode | null>(null)
const editingNote = ref('')
const contextMenu = ref<ContextMenuState>({ show: false, x: 0, y: 0, node: null, isGraphMenu: false })
const operationMode = ref<'merge' | 'rebase' | null>(null)
const operationSourceNode = ref<TimelineNode | null>(null)
const nodePositions = ref<Map<number, { x: number; y: number }>>(new Map())
const nodeBranchMap = ref<Map<number, number | null>>(new Map())
const hoverTimeout = ref<ReturnType<typeof setTimeout> | null>(null)
const nodeOffsets = ref<Record<number, number>>({})
const draggingNodeId = ref<number | null>(null)
const dragStartY = ref(0)
const dragStartOffset = ref(0)
const snapPreviewOffset = ref<number | null>(null)
const renderVersion = ref(0)
const errorTooltipPosition = ref({ x: 0, y: 0 })
const isPanning = ref(false)
const panStartX = ref(0)
const panStartY = ref(0)
const panStartScrollLeft = ref(0)
const panStartScrollTop = ref(0)
const scale = ref(1)
const allBranches = computed<Branch[]>(() => {
  const branches: Branch[] = []
  const rootNodes = props.nodes.filter(n => !n.parentId)
  rootNodes.forEach(root => {
    const rootChildren = props.nodes.filter(n => n.parentId === root.id)
    rootChildren.forEach(child => {
      const branchIndex = branches.length
      const branchData = props.branchNames?.[child.id] || {}
      branches.push({
        id: child.id,
        name: branchData.name || '',
        color: branchData.color || BRANCH_COLORS[branchIndex % BRANCH_COLORS.length],
        description: branchData.description || '',
        nodeCount: 1 + countDescendants(child.id),
        columnIndex: branchIndex
      })
    })
  })
  const processed = new Set<number>()
  rootNodes.forEach(root => processed.add(root.id))
  props.nodes.forEach(node => {
    if (node.parentId && !processed.has(node.parentId)) {
      const siblings = props.nodes.filter(n => n.parentId === node.parentId)
      if (siblings.length > 1) {
        siblings.forEach((sibling, idx) => {
          if (idx > 0) {
            const branchIndex = branches.length
            const branchData = props.branchNames?.[sibling.id] || {}
            branches.push({
              id: sibling.id,
              name: branchData.name || '',
              color: branchData.color || BRANCH_COLORS[branchIndex % BRANCH_COLORS.length],
              description: branchData.description || '',
              nodeCount: 1 + countDescendants(sibling.id),
              columnIndex: branchIndex
            })
          }
        })
      }
      processed.add(node.parentId)
    }
  })
  branches.sort((a, b) => a.id - b.id)
  return branches
})
const selectableBranches = computed(() => allBranches.value)
const pendingVisibleBranchCount = computed(() => pendingBranchIds.value.length)
const visibleBranchList = computed(() => {
  if (props.readonly || props.fullMode || allBranches.value.length <= 3) {
    return allBranches.value
  }
  const idSet = new Set(displayBranchIds.value)
  return allBranches.value.filter(b => idSet.has(b.id))
})
const graphRenderKey = computed(() => {
  return renderVersion.value + '-' + displayBranchIds.value.join('-') + '-' + allBranches.value.length
})
const pathNodes = computed(() => {
  if (!props.pathMode || !props.currentNodeId) return new Set<number>()
  const pathSet = new Set<number>()
  let current = props.nodes.find(n => n.id === props.currentNodeId)
  while (current) {
    pathSet.add(current.id)
    if (!current.parentId) break
    current = props.nodes.find(n => n.id === current!.parentId)
  }
  return pathSet
})
const displayNodes = computed(() => {
  nodePositions.value.clear()
  nodeBranchMap.value.clear()
  assignNodeBranches()
  if (props.pathMode && props.currentNodeId) {
    const pathNodeIds = pathNodes.value
    let levelIndex = 0
    const pathNodesOrdered: TimelineNode[] = []
    let current = props.nodes.find(n => !n.parentId && pathNodeIds.has(n.id))
    while (current) {
      pathNodesOrdered.push(current)
      const next = props.nodes.find(n => n.parentId === current!.id && pathNodeIds.has(n.id))
      current = next
    }
    const nodeWidth = 160
    const centerX = (320 - nodeWidth) / 2
    return pathNodesOrdered.map(node => {
      const branchId = nodeBranchMap.value.get(node.id)
      const branch = allBranches.value.find(b => b.id === branchId)
      const isRootNode = !node.parentId
      const pos = { x: centerX, y: PADDING + levelIndex * LEVEL_HEIGHT }
      nodePositions.value.set(node.id, pos)
      levelIndex++
      return {
        ...node,
        x: pos.x,
        y: pos.y,
        branchColor: isRootNode ? ROOT_NODE_COLOR : (branch?.color || BRANCH_COLORS[0]),
        isHiddenBranch: false
      }
    })
  }
  const visibleBranchIdSet = props.fullMode || props.readonly
    ? new Set(allBranches.value.map(b => b.id))
    : new Set(displayBranchIds.value)
  const visibleNodes = new Set<number>()
  const neededHiddenNodes = new Set<number>()
  props.nodes.forEach(node => {
    const branchId = nodeBranchMap.value.get(node.id)
    const isRootNode = !node.parentId
    const isHiddenBranch = !props.fullMode && !isRootNode && branchId && !visibleBranchIdSet.has(branchId)
    if (!isHiddenBranch) {
      visibleNodes.add(node.id)
    }
  })
  visibleNodes.forEach(nodeId => {
    let current = props.nodes.find(n => n.id === nodeId)
    while (current && current.parentId) {
      const parent = props.nodes.find(n => n.id === current!.parentId)
      if (parent && !visibleNodes.has(parent.id)) {
        neededHiddenNodes.add(parent.id)
      }
      current = parent
    }
  })
  const nodesToDisplay = props.nodes.filter(node => {
    return visibleNodes.has(node.id) || neededHiddenNodes.has(node.id)
  })
  return nodesToDisplay.map(node => {
    const branchId = nodeBranchMap.value.get(node.id)
    const branch = allBranches.value.find(b => b.id === branchId)
    const isRootNode = !node.parentId
    const isHiddenBranch = neededHiddenNodes.has(node.id)
    const pos = calculateNodePosition(node, isHiddenBranch)
    nodePositions.value.set(node.id, pos)
    return {
      ...node,
      x: pos.x,
      y: pos.y,
      branchColor: isRootNode ? ROOT_NODE_COLOR : (branch?.color || BRANCH_COLORS[0]),
      isHiddenBranch
    }
  })
})
const lines = computed(() => {
  const result: { id: string; fromId: number; toId: number; path: string; color: string; isLogicConnection: boolean }[] = []
  for (const node of displayNodes.value) {
    if (node.parentId) {
      const parentNode = props.nodes.find(n => n.id === node.parentId)
      const parentPos = nodePositions.value.get(node.parentId)
      if (parentPos) {
        const parentOffset = getNodeOffset(node.parentId)
        const childOffset = getNodeOffset(node.id)
        const adjustedParentPos = { x: parentPos.x, y: parentPos.y + parentOffset }
        const childPos = { x: node.x, y: node.y + childOffset }
        const isLogicConnection = parentNode?.isLogic === true
        const isChildLogic = node.isLogic === true
        result.push({
          id: `${node.parentId}-${node.id}`,
          fromId: node.parentId,
          toId: node.id,
          path: createPath(adjustedParentPos, childPos, isLogicConnection, isChildLogic),
          color: node.branchColor,
          isLogicConnection
        })
      }
    }
  }
  return result
})
const logicCapsules = computed<LogicCapsule[]>(() => {
  const capsules: LogicCapsule[] = []
  for (const line of lines.value) {
    if (line.isLogicConnection) {
      const childNode = displayNodes.value.find(n => n.id === line.toId)
      if (childNode) {
        const childOffset = getNodeOffset(line.toId)
        const centerOffset = actualColumnWidth.value / 2
        const x = childNode.x + centerOffset
        const y = childNode.y + childOffset - 8
        capsules.push({
          id: `capsule-${line.id}`,
          x,
          y,
          label: childNode.optionLabel || childNode.name || '选项',
          description: childNode.note || '',
          nodeId: childNode.id
        })
      }
    }
  }
  return capsules
})
const svgWidth = computed(() => {
  if (props.pathMode) return 320
  return Math.max(requiredPanelWidth.value, 300)
})
const svgHeight = computed(() => {
  if (displayNodes.value.length === 0) return 300
  return Math.max(...displayNodes.value.map(n => n.y)) + 100
})
const requiredPanelWidth = computed(() => {
  const cols = visibleBranchList.value.length
  if (cols <= 1) return 200
  return cols * COLUMN_WIDTH + PADDING * 2
})
const actualColumnWidth = computed(() => {
  const cols = visibleBranchList.value.length
  if (cols <= 1) return requiredPanelWidth.value - PADDING * 2
  return COLUMN_WIDTH
})
const errorTooltipStyle = computed(() => ({
  left: errorTooltipPosition.value.x + 'px',
  top: errorTooltipPosition.value.y + 'px'
}))
watch(() => props.nodes, () => {
  initVisibleBranches()
  initSelectedNode()
}, { immediate: true, deep: true })
watch(() => props.currentNodeId, (newId) => {
  if (newId) selectedNodeId.value = newId
}, { immediate: true })
watch(requiredPanelWidth, (newWidth) => {
  emit('width-change', newWidth)
}, { immediate: true })
function countDescendants(nodeId: number): number {
  let count = 0
  const children = props.nodes.filter(n => n.parentId === nodeId)
  for (const child of children) {
    count += 1 + countDescendants(child.id)
  }
  return count
}
function assignNodeBranches() {
  const assignBranch = (nodeId: number, branchId: number | null) => {
    nodeBranchMap.value.set(nodeId, branchId)
    const children = props.nodes.filter(n => n.parentId === nodeId)
    children.forEach((child, index) => {
      if (index === 0) {
        assignBranch(child.id, branchId)
      } else {
        assignBranch(child.id, child.id)
      }
    })
  }
  const rootNodes = props.nodes.filter(n => !n.parentId)
  rootNodes.forEach(root => {
    nodeBranchMap.value.set(root.id, null)
    const rootChildren = props.nodes.filter(n => n.parentId === root.id)
    rootChildren.forEach(child => {
      assignBranch(child.id, child.id)
    })
  })
}
function calculateNodePosition(node: TimelineNode, isHiddenBranch = false): { x: number; y: number } {
  let level = 0
  let currentId = node.parentId
  while (currentId) {
    level++
    const parent = props.nodes.find(n => n.id === currentId)
    currentId = parent?.parentId
  }
  const branchId = nodeBranchMap.value.get(node.id)
  if (branchId === null) {
    const totalCols = visibleBranchList.value.length || 1
    const centerX = (totalCols * actualColumnWidth.value) / 2 - actualColumnWidth.value / 2
    return { x: PADDING + centerX, y: PADDING + level * LEVEL_HEIGHT }
  }
  if (isHiddenBranch) {
    const parentNode = props.nodes.find(n => n.id === node.parentId)
    if (parentNode) {
      const parentPos = nodePositions.value.get(parentNode.id)
      if (parentPos) {
        return { x: parentPos.x, y: PADDING + level * LEVEL_HEIGHT }
      }
    }
  }
  let colIndex = 0
  const branchIdx = visibleBranchList.value.findIndex(b => b.id === branchId)
  if (branchIdx >= 0) colIndex = branchIdx
  return { x: PADDING + colIndex * actualColumnWidth.value, y: PADDING + level * LEVEL_HEIGHT }
}
function createPath(from: { x: number; y: number }, to: { x: number; y: number }, isParentLogic = false, isChildLogic = false): string {
  const centerOffset = actualColumnWidth.value / 2
  const startX = from.x + centerOffset
  const startY = from.y + (isParentLogic ? 30 : 36)
  const endX = to.x + centerOffset
  const endY = to.y + (isChildLogic ? 10 : 10)
  if (Math.abs(startX - endX) < 1) {
    return `M ${startX} ${startY} L ${endX} ${endY}`
  } else {
    const midY = (startY + endY) / 2
    return `M ${startX} ${startY} L ${startX} ${midY} L ${endX} ${midY} L ${endX} ${endY}`
  }
}
function getNodeOffset(nodeId: number): number {
  return nodeOffsets.value[nodeId] || 0
}
function getBranchDisplayName(branch: Branch): string {
  if (branch.name) return branch.name
  const idx = allBranches.value.findIndex(b => b.id === branch.id)
  return idx === 0 ? '主世界线' : '世界线 ' + idx
}
function setHoveredNode(nodeId: number, event: MouseEvent) {
  if (hoverTimeout.value) {
    clearTimeout(hoverTimeout.value)
    hoverTimeout.value = null
  }
  hoveredNodeId.value = nodeId
  if (event && props.nodeErrorMap?.[nodeId]) {
    nextTick(() => {
      const target = event.currentTarget as HTMLElement
      if (target) {
        const rect = target.getBoundingClientRect()
        errorTooltipPosition.value = { x: rect.left + rect.width / 2, y: rect.top - 8 }
      }
    })
  }
}
function clearHoveredNode() {
  if (hoverTimeout.value) clearTimeout(hoverTimeout.value)
  hoverTimeout.value = setTimeout(() => {
    hoveredNodeId.value = null
    hoverTimeout.value = null
  }, 50)
}
function initVisibleBranches() {
  const allBranchIds = allBranches.value.map(b => b.id)
  if (allBranches.value.length <= 3) {
    displayBranchIds.value = allBranchIds
  } else {
    const currentValid = displayBranchIds.value.filter(id => allBranchIds.includes(id))
    if (currentValid.length === 3) {
      const newBranchIds = allBranchIds.filter(id => !currentValid.includes(id))
      if (newBranchIds.length > 0) {
        const newestBranchId = newBranchIds[newBranchIds.length - 1]
        const sortedCurrent = [...currentValid].sort((a, b) => a - b)
        sortedCurrent.shift()
        const newIds = [...sortedCurrent, newestBranchId]
        newIds.sort((a, b) => a - b)
        displayBranchIds.value = newIds
      } else {
        displayBranchIds.value = currentValid
      }
    } else {
      const sortedDesc = [...allBranches.value].sort((a, b) => b.id - a.id)
      let top3Ids = sortedDesc.slice(0, 3).map(b => b.id)
      if (props.requiredBranchId && !top3Ids.includes(props.requiredBranchId)) {
        top3Ids = [props.requiredBranchId, ...top3Ids.slice(0, 2)]
      }
      top3Ids.sort((a, b) => a - b)
      displayBranchIds.value = top3Ids
    }
  }
}
function initSelectedNode() {
  if (!selectedNodeId.value && props.nodes.length > 0) {
    const root = props.nodes.find(n => !n.parentId)
    if (root) selectedNodeId.value = root.id
  }
}
function isNodeHighlighted(nodeId: number): boolean {
  if (!searchTag.value) return true
  return filteredNodeIds.value.has(nodeId)
}
function selectAndOpenNode(node: TimelineNode) {
  if (props.fullMode) {
    selectedNodeId.value = node.id
    emit('node-click', node)
    return
  }
  if (operationMode.value) {
    completeOperation(node)
    return
  }
  if (node.isLogic) {
    selectedLogicCapsule.value = {
      id: 'capsule-' + node.id,
      nodeId: node.id,
      label: node.optionLabel || node.name || '选项',
      description: node.note || '',
      x: 0,
      y: 0
    }
    editingCapsuleLabel.value = selectedLogicCapsule.value.label
    editingCapsuleDescription.value = selectedLogicCapsule.value.description
    showLogicCapsuleModal.value = true
    return
  }
  selectedNodeId.value = node.id
  emit('select', node)
  emit('set-current', node)
}
function showContextMenu(event: MouseEvent, node: TimelineNode) {
  event.stopPropagation()
  emit('context-menu-open')
  contextMenu.value = { show: true, x: event.clientX, y: event.clientY, node, isGraphMenu: false }
}
function showGraphContextMenu(event: MouseEvent) {
  if ((event.target as HTMLElement).closest('.node-wrapper')) return
  emit('context-menu-open')
  contextMenu.value = { show: true, x: event.clientX, y: event.clientY, node: null, isGraphMenu: true }
}
function closeContextMenu() {
  contextMenu.value.show = false
}
function handleFullModeEdit(node: TimelineNode) {
  emit('full-mode-edit', node)
}
function createRootNode() {
  emit('create-root')
  closeContextMenu()
}
function createParallelWorld() {
  if (contextMenu.value.node) {
    emit('create-parallel-world', contextMenu.value.node)
  }
  closeContextMenu()
}
function deleteNodeAndRelink() {
  if (contextMenu.value.node) {
    emit('delete-and-relink', contextMenu.value.node)
  }
  closeContextMenu()
}
function completeOperation(targetNode: TimelineNode) {
  if (operationMode.value === 'merge' && operationSourceNode.value) {
    emit('merge', { source: operationSourceNode.value, target: targetNode })
  } else if (operationMode.value === 'rebase' && operationSourceNode.value) {
    emit('rebase', { node: operationSourceNode.value, newParent: targetNode })
  }
  cancelOperation()
}
function cancelOperation() {
  operationMode.value = null
  operationSourceNode.value = null
}
function openBranchSelector() {
  pendingBranchIds.value = [...displayBranchIds.value]
  showBranchSelector.value = true
}
function togglePendingBranchVisibility(branchId: number) {
  if (branchId === props.requiredBranchId) return
  const index = pendingBranchIds.value.indexOf(branchId)
  if (index >= 0) {
    pendingBranchIds.value.splice(index, 1)
  } else {
    pendingBranchIds.value.push(branchId)
  }
}
function confirmBranchSelection() {
  displayBranchIds.value = [...pendingBranchIds.value].sort((a, b) => a - b)
  showBranchSelector.value = false
}
function startEditBranchName(branch: Branch) {
  editingBranchId.value = branch.id
  editingBranchName.value = branch.name || ''
  editingBranchColor.value = branch.color && branch.color !== ROOT_NODE_COLOR ? branch.color : BRANCH_COLORS[0]
  editingBranchDescription.value = branch.description || ''
  showBranchEditModal.value = true
}
function saveBranchEdit() {
  if (editingBranchId.value) {
    emit('update-branch', {
      branchId: editingBranchId.value,
      name: editingBranchName.value,
      color: editingBranchColor.value,
      description: editingBranchDescription.value
    })
  }
  showBranchEditModal.value = false
  editingBranchId.value = null
  editingBranchName.value = ''
  editingBranchColor.value = ''
  editingBranchDescription.value = ''
}
function saveNote() {
  if (editingNoteNode.value) {
    emit('update-note', editingNoteNode.value.id, '', editingNote.value)
  }
  showNoteModal.value = false
  editingNoteNode.value = null
  editingNote.value = ''
}
function showLogicCapsuleDetail(capsule: LogicCapsule) {
  selectedLogicCapsule.value = capsule
  editingCapsuleLabel.value = capsule.label || ''
  editingCapsuleDescription.value = capsule.description || ''
  showLogicCapsuleModal.value = true
}
function saveCapsuleEdit() {
  if (selectedLogicCapsule.value && selectedLogicCapsule.value.nodeId) {
    emit('update-note', selectedLogicCapsule.value.nodeId, editingCapsuleLabel.value, editingCapsuleDescription.value)
  }
  showLogicCapsuleModal.value = false
}
function handleNodeDragStart(e: DragEvent, node: TimelineNode) {
  draggingNodeId.value = node.id
  dragStartY.value = e.clientY
  dragStartOffset.value = getNodeOffset(node.id)
  const img = new Image()
  img.src = 'data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7'
  e.dataTransfer?.setDragImage(img, 0, 0)
  if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move'
}
function handleNodeDrag(e: DragEvent, node: TimelineNode) {
  if (!draggingNodeId.value || e.clientY === 0) return
  const deltaY = e.clientY - dragStartY.value
  const newOffset = dragStartOffset.value + deltaY
  const clampedOffset = Math.max(0, newOffset)
  nodeOffsets.value = { ...nodeOffsets.value, [node.id]: clampedOffset }
  snapPreviewOffset.value = Math.round(clampedOffset / LEVEL_HEIGHT) * LEVEL_HEIGHT
}
function handleNodeDragEnd(e: DragEvent, node: TimelineNode) {
  if (draggingNodeId.value) {
    const currentOffset = getNodeOffset(node.id)
    const snappedOffset = Math.round(currentOffset / LEVEL_HEIGHT) * LEVEL_HEIGHT
    nodeOffsets.value = { ...nodeOffsets.value, [node.id]: snappedOffset }
  }
  draggingNodeId.value = null
  snapPreviewOffset.value = null
}
function handleNodeDblClick(node: TimelineNode) {
  emit('full-mode-edit', node)
}
function handlePanStart(e: MouseEvent) {
  if ((e.target as HTMLElement).closest('.node-shape') || (e.target as HTMLElement).closest('.add-btn') || (e.target as HTMLElement).closest('.logic-capsule')) return
  isPanning.value = true
  panStartX.value = e.clientX
  panStartY.value = e.clientY
  const container = graphContainer.value
  if (container) {
    panStartScrollLeft.value = container.scrollLeft
    panStartScrollTop.value = container.scrollTop
    container.style.cursor = 'grabbing'
  }
}
function handlePanMove(e: MouseEvent) {
  if (!isPanning.value) return
  const container = graphContainer.value
  if (!container) return
  const dx = e.clientX - panStartX.value
  const dy = e.clientY - panStartY.value
  container.scrollLeft = panStartScrollLeft.value - dx
  container.scrollTop = panStartScrollTop.value - dy
}
function handlePanEnd() {
  isPanning.value = false
  const container = graphContainer.value
  if (container) container.style.cursor = ''
}
function handleWheel(e: WheelEvent) {
  if (!e.ctrlKey && !e.metaKey) {
    e.preventDefault()
    const delta = e.deltaY > 0 ? -0.1 : 0.1
    const newScale = Math.max(0.5, Math.min(2, scale.value + delta))
    scale.value = newScale
  }
}
</script>
<style scoped>
.timeline-graph {
  display: flex;
  flex-direction: column;
  height: 100%;
  position: relative;
}
.timeline-graph.readonly-mode {
  pointer-events: none;
  cursor: default;
}
.timeline-graph.full-mode {
  height: auto;
  min-height: 100%;
}
.timeline-graph.full-mode .graph-container {
  overflow: visible;
}
.timeline-graph.full-mode .node-shape.selected {
  box-shadow: 0 0 0 3px var(--accent);
}
.timeline-graph.readonly-mode .node-shape {
  cursor: default;
}
.branch-headers {
  display: flex;
  border-bottom: 1px solid var(--border);
  background: var(--bg-tertiary);
  overflow-x: auto;
  align-items: center;
  box-sizing: border-box;
  justify-content: space-between;
}
.branch-headers-list {
  display: flex;
  align-items: center;
  flex: 1;
  max-width: calc(100% - 40px);
}
.btn-branch-selector {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin: 0 6px;
  transition: all 0.15s;
}
.btn-branch-selector:hover {
  background: var(--accent);
  color: white;
}
.branch-header {
  padding: 6px 10px;
  display: flex;
  align-items: center;
  gap: 6px;
  border-top: 3px solid;
  cursor: pointer;
  transition: background 0.15s;
  min-width: 0;
  box-sizing: border-box;
}
.branch-header:hover {
  background: var(--bg-secondary);
}
.branch-color-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.branch-name {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-primary);
}
.graph-container {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  position: relative;
}
.graph-content {
  position: relative;
  padding: 20px;
}
.graph-lines {
  position: absolute;
  top: 0;
  left: 0;
  pointer-events: none;
  z-index: 1;
}
.connection-line {
  fill: none;
  stroke-width: 2;
  transition: stroke 0.2s;
}
.connection-line.highlighted {
  stroke-width: 3;
}
.connection-line.logic-line {
  stroke: #f97316 !important;
  stroke-width: 2.5;
}
.logic-capsule {
  position: absolute;
  transform: translate(-50%, -50%);
  background: #f97316;
  color: white;
  font-size: 10px;
  padding: 3px 8px;
  border-radius: 10px;
  white-space: nowrap;
  cursor: pointer;
  z-index: 5;
  transition: all 0.2s;
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.logic-capsule:hover {
  background: #ea580c;
  transform: translate(-50%, -50%) scale(1.1);
}
.node-wrapper {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 20px 40px 20px;
  z-index: 2;
}
.node-snap-ghost {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  pointer-events: none;
  z-index: 100;
}
.ghost-shape {
  width: 80px;
  height: 28px;
  background: var(--accent);
  opacity: 0.3;
  border-radius: 20px;
  border: 2px dashed var(--accent);
}
.node-shape-container {
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
}
.node-shape {
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 6px 16px;
  min-width: 60px;
  max-width: 120px;
}
.node-shape-name {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.node-capsule {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 2px solid;
  border-radius: 20px;
}
.node-capsule:hover,
.node-capsule.hovered {
  background: var(--bg-secondary);
}
.node-capsule.selected {
  background: var(--accent);
  color: white;
  box-shadow: 0 0 0 3px rgba(233, 69, 96, 0.3);
}
.node-capsule.filtered {
  opacity: 0.3;
}
.node-capsule.has-conversation {
  background: #3b82f6;
  color: white;
}
.node-capsule.has-conversation:hover,
.node-capsule.has-conversation.hovered {
  background: #2563eb;
}
.node-capsule.has-error {
  background: var(--bg-tertiary);
}
.node-capsule.has-error:hover,
.node-capsule.has-error.hovered {
  background: var(--bg-secondary);
}
.node-capsule.current {
  background: #22c55e;
  color: #000;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.3);
}
.node-capsule.current:hover,
.node-capsule.current.hovered {
  background: #16a34a;
}
.node-wide-hexagon {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: none;
  clip-path: polygon(8% 50%, 15% 0%, 85% 0%, 92% 50%, 85% 100%, 15% 100%);
  padding: 6px 20px;
  position: relative;
}
.node-wide-hexagon::before {
  content: '';
  position: absolute;
  top: -2px;
  left: -2px;
  right: -2px;
  bottom: -2px;
  background: var(--node-border-color, #6b7280);
  clip-path: polygon(8% 50%, 15% 0%, 85% 0%, 92% 50%, 85% 100%, 15% 100%);
  z-index: -1;
}
.node-shape.hidden-branch-node {
  cursor: default;
  pointer-events: none;
}
.error-name-capsule {
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  background: #ef4444;
  color: #000;
  padding: 2px 8px;
  border-radius: 10px;
}
.node-hexagon {
  background: var(--warning);
  color: white;
  border: none;
  clip-path: polygon(10% 0%, 90% 0%, 100% 50%, 90% 100%, 10% 100%, 0% 50%);
  padding: 6px 20px;
}
.node-hexagon:hover,
.node-hexagon.hovered {
  background: #e6a700;
  box-shadow: 0 0 8px rgba(255, 152, 0, 0.5);
}
.node-hexagon.current {
  background: var(--warning);
  box-shadow: 0 0 0 3px rgba(255, 152, 0, 0.4);
}
.node-hexagon.selected {
  background: var(--warning);
  box-shadow: 0 0 0 3px rgba(255, 152, 0, 0.6);
}
.node-hexagon.filtered {
  opacity: 0.3;
}
.add-btn {
  position: absolute;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--accent);
  color: white;
  border: none;
  font-size: 18px;
  line-height: 1;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  z-index: 10;
}
.add-btn:hover {
  transform: scale(1.2);
  background: var(--accent-hover, #d63850);
}
.node-action-buttons {
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-top: 4px;
  display: flex;
  gap: 6px;
  z-index: 10;
}
.node-action-buttons .add-btn {
  position: relative;
  top: auto;
  left: auto;
  margin: 0;
  transform: none;
}
.node-action-buttons .add-btn:hover {
  transform: scale(1.2);
}
.edit-btn {
  background: var(--accent);
  color: white;
}
.edit-btn:hover {
  background: var(--accent-hover, #2563eb);
  color: white;
}
.add-child {
  top: 100%;
  left: 50%;
  margin-top: 4px;
  transform: translateX(-50%);
}
.add-child:hover {
  transform: translateX(-50%) scale(1.2);
}
.fade-enter-active {
  animation: slide-from-center-down 0.2s ease-out forwards;
}
.fade-leave-active {
  transition: opacity 0.15s ease-out;
}
.fade-leave-to {
  opacity: 0;
}
@keyframes slide-from-center-down {
  from {
    opacity: 0;
    top: 50%;
    transform: translateX(-50%) translateY(-50%);
  }
  to {
    opacity: 1;
    top: 100%;
    margin-top: 4px;
    transform: translateX(-50%);
  }
}
.node-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-top: 6px;
  justify-content: center;
}
.tag {
  font-size: 9px;
  color: var(--accent);
  background: rgba(233, 69, 96, 0.15);
  padding: 1px 5px;
  border-radius: 8px;
}
.note-textarea {
  min-height: 100px;
}
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 6px;
  padding: 4px 0;
  min-width: 140px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}
.menu-item {
  padding: 8px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.15s;
}
.menu-item:hover {
  background: var(--bg-tertiary);
}
.menu-item.danger {
  color: var(--danger);
}
.menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
}
.operation-hint {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--accent);
  color: white;
  padding: 8px 16px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}
.operation-hint .btn {
  background: rgba(255, 255, 255, 0.2);
}
.operation-hint .btn:hover {
  background: rgba(255, 255, 255, 0.3);
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.branch-selector-modal {
  background: var(--bg-secondary);
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.branch-selector-modal .modal-header {
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.branch-selector-modal .modal-header h3 {
  font-size: 14px;
  font-weight: 500;
}
.branch-selector-modal .modal-body {
  padding: 8px;
  overflow-y: auto;
  flex: 1;
}
.branch-selector-modal .modal-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: flex-end;
}
.branch-options-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px;
}
.branch-option-capsule {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.15s;
  border: 2px solid var(--border);
  background: transparent;
  font-size: 13px;
}
.branch-option-capsule:hover {
  border-color: var(--accent);
}
.branch-option-capsule.selected {
  background: var(--accent);
  border-color: var(--accent);
  color: white;
}
.branch-option-capsule.selected .branch-name {
  color: white;
}
.branch-option-capsule.required {
  cursor: default;
}
.branch-option-capsule.required:not(.selected) {
  border-color: var(--accent);
  background: rgba(233, 69, 96, 0.1);
}
.required-badge {
  font-size: 10px;
  background: var(--accent);
  color: white;
  padding: 1px 6px;
  border-radius: 8px;
  margin-left: 4px;
}
.branch-option-capsule .branch-color-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}
.branch-option-capsule .branch-name {
  color: var(--text-primary);
  white-space: nowrap;
}
.modal {
  background: var(--bg-secondary);
  border-radius: 8px;
  width: 90%;
  max-width: 400px;
  overflow: hidden;
}
.modal-sm {
  max-width: 360px;
}
.modal .modal-header {
  padding: 14px 16px;
  border-bottom: 1px solid var(--border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.modal .modal-header h3 {
  font-size: 14px;
}
.modal .modal-body {
  padding: 16px;
}
.modal .modal-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
.form-group {
  margin-bottom: 0;
}
.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}
.form-control {
  width: 100%;
  padding: 10px;
  background: var(--bg-primary);
  border: 1px solid var(--border);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 13px;
}
.form-control:focus {
  outline: none;
  border-color: var(--accent);
}
.color-picker-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: 6px;
  padding: 8px;
  background: var(--bg-primary);
  border-radius: 6px;
  border: 1px solid var(--border);
}
.color-option {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid transparent;
  transition: all 0.15s;
}
.color-option:hover {
  transform: scale(1.15);
}
.color-option.selected {
  border-color: var(--text-primary);
  box-shadow: 0 0 0 2px var(--bg-secondary);
}
.node-conversations-tooltip {
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 12px;
  min-width: 120px;
  max-width: 200px;
  z-index: 1000;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  pointer-events: none;
}
.node-conversations-tooltip .tooltip-title {
  font-size: 10px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  padding-bottom: 4px;
  border-bottom: 1px solid var(--border);
}
.node-conversations-tooltip .tooltip-item {
  font-size: 11px;
  color: var(--text-primary);
  padding: 2px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  border: none;
  transition: all 0.15s;
}
.btn-primary {
  background: var(--accent);
  color: white;
}
.btn-primary:hover {
  background: var(--accent-hover, #d63850);
}
.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.btn-secondary:hover {
  background: var(--bg-hover);
}
.btn-danger {
  background: var(--danger);
  color: white;
}
.btn-icon {
  padding: 4px 8px;
  background: transparent;
  font-size: 18px;
}
.btn-icon:hover {
  background: var(--bg-tertiary);
}
.btn-sm {
  padding: 4px 12px;
  font-size: 12px;
}
</style>
<style>
.node-error-tooltip-fixed {
  position: fixed;
  transform: translateX(-50%) translateY(-100%);
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 8px 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 10000;
  min-width: 150px;
  max-width: 280px;
  pointer-events: none;
}
.node-error-tooltip-fixed .error-item {
  font-size: 12px;
  color: var(--text-secondary);
  padding: 2px 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.node-error-tooltip-fixed .error-entry {
  color: #ef4444;
  font-weight: 500;
}
</style>
