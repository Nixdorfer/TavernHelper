<template>
  <div class="fullscreen-overlays">
    <transition name="fade-overlay">
      <div v-if="showNodeSelectDialog && pendingProjectForNodeSelect" class="full-worldtree-overlay" @keydown.esc="$emit('cancel-node-selection')" tabindex="0" ref="nodeSelectOverlay">
        <div class="full-worldtree-container">
          <TimelineGraph
            :nodes="pendingProjectForNodeSelect.projectData.timeline"
            :current-node-id="nodeSelectDialogSelectedId"
            :branch-names="{}"
            :readonly="false"
            :full-mode="true"
            :select-only="true"
            @node-click="$emit('node-select-dialog-click', $event)"
          />
        </div>
        <div class="full-tree-bottom-buttons">
          <button class="btn-full-tree-action btn-confirm-select" :class="{ disabled: !nodeSelectDialogSelectedId }" @click="nodeSelectDialogSelectedId && $emit('confirm-node-selection', nodeSelectDialogSelectedId)" title="确认选择">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
          </button>
          <button class="btn-full-tree-action btn-close-action" @click="$emit('cancel-node-selection')" title="取消">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>
        </div>
      </div>
    </transition>
    <div v-if="capsLockWarning" class="caps-lock-warning">
      大写锁定已开启，安全模式可能无法立刻生效
    </div>
    <div v-if="safeModeActive && safeModeAction === 'randomChars'" class="safe-mode-decrypt-error">
      程序解密失败，请重新订阅
    </div>
    <div v-if="showFullWorldTree" class="full-worldtree-overlay" @mousedown="startDrag" @mousemove="onDrag" @mouseup="stopDrag" @mouseleave="stopDrag" @wheel.prevent="onZoom" @keydown.esc="$emit('close-full-world-tree')" tabindex="0" ref="fullWorldTreeOverlay">
      <div class="full-worldtree-container" :style="{ transform: `translate(${fullTreeDragOffset.x}px, ${fullTreeDragOffset.y}px) scale(${fullTreeZoom})` }">
        <TimelineGraph
          :nodes="fullWorldTreeData.timeline"
          :current-node-id="fullTreeSelectedNodeId || fullWorldTreeData.currentNode"
          :branch-names="{}"
          :readonly="false"
          :full-mode="true"
          :select-only="true"
          :node-conversation-map="fullTreeNodeConversationMap"
          :node-conversation-names="fullTreeNodeConversationNames"
          @node-click="$emit('full-tree-node-click', $event)"
        />
      </div>
      <div class="full-tree-bottom-actions">
        <button :class="['btn-full-tree-bottom', 'btn-add', { disabled: !fullTreeSelectedNodeId }]" @click="fullTreeSelectedNodeId && $emit('show-create-parallel-world', { nodeId: fullTreeSelectedNodeId })" :disabled="!fullTreeSelectedNodeId" title="添加平行世界线">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"></line>
            <line x1="5" y1="12" x2="19" y2="12"></line>
          </svg>
        </button>
        <button :class="['btn-full-tree-bottom', 'btn-relocate', { disabled: !fullTreeSelectedNodeId }]" @click="fullTreeSelectedNodeId && $emit('relocate-conversation-to-node')" :disabled="!fullTreeSelectedNodeId" title="定位对话到选中节点">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="23 4 23 10 17 10"></polyline>
            <path d="M20.49 15a9 9 0 1 1-2.12-9.36L23 10"></path>
          </svg>
        </button>
        <button class="btn-full-tree-bottom btn-close" @click="$emit('close-full-world-tree')" title="关闭">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="18" y1="6" x2="6" y2="18"></line>
            <line x1="6" y1="6" x2="18" y2="18"></line>
          </svg>
        </button>
      </div>
    </div>
    <Transition name="modal">
      <div v-if="galleryPreviewImage" class="modal-overlay" @mousedown.self="$emit('close-gallery-preview')">
        <div class="gallery-preview-modal">
          <div class="preview-modal-header">
            <h3>{{ galleryPreviewImage.fileName }}</h3>
            <button class="btn btn-icon" @click="$emit('close-gallery-preview')">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="preview-modal-body">
            <img :src="galleryImageSrc" :alt="galleryPreviewImage.fileName" />
          </div>
          <div class="preview-modal-footer">
            <div class="preview-info-row">
              <span class="preview-label">URL:</span>
              <span class="preview-value">{{ galleryPreviewImage.remoteUrl || '无远程URL' }}</span>
              <button v-if="galleryPreviewImage.remoteUrl" class="btn-copy" @click="$emit('copy-to-clipboard', galleryPreviewImage.remoteUrl)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
              </button>
            </div>
            <div class="preview-info-row">
              <span class="preview-label">Hash:</span>
              <span class="preview-value hash-value">{{ galleryPreviewImage.hash }}</span>
              <button class="btn-copy" @click="$emit('copy-to-clipboard', galleryPreviewImage.hash)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import TimelineGraph from '@/components/business/TimelineGraph.vue'
interface ProjectData {
  projectData: {
    timeline: any[]
  }
}
interface WorldTreeData {
  timeline: any[]
  currentNode: string
}
interface GalleryImage {
  fileName: string
  remoteUrl?: string
  hash: string
}
interface DragOffset {
  x: number
  y: number
}
const props = withDefaults(defineProps<{
  showNodeSelectDialog?: boolean
  pendingProjectForNodeSelect?: ProjectData | null
  nodeSelectDialogSelectedId?: string
  capsLockWarning?: boolean
  safeModeActive?: boolean
  safeModeAction?: string
  showFullWorldTree?: boolean
  fullWorldTreeData?: WorldTreeData
  fullTreeSelectedNodeId?: string
  fullTreeNodeConversationMap?: Record<string, any>
  fullTreeNodeConversationNames?: Record<string, string>
  fullTreeDragOffset?: DragOffset
  fullTreeZoom?: number
  galleryPreviewImage?: GalleryImage | null
  galleryImageSrc?: string
}>(), {
  showNodeSelectDialog: false,
  pendingProjectForNodeSelect: null,
  nodeSelectDialogSelectedId: '',
  capsLockWarning: false,
  safeModeActive: false,
  safeModeAction: '',
  showFullWorldTree: false,
  fullWorldTreeData: () => ({ timeline: [], currentNode: '' }),
  fullTreeSelectedNodeId: '',
  fullTreeNodeConversationMap: () => ({}),
  fullTreeNodeConversationNames: () => ({}),
  fullTreeDragOffset: () => ({ x: 0, y: 0 }),
  fullTreeZoom: 1,
  galleryPreviewImage: null,
  galleryImageSrc: ''
})
const emit = defineEmits<{
  'cancel-node-selection': []
  'confirm-node-selection': [nodeId: string]
  'node-select-dialog-click': [event: any]
  'start-full-tree-drag': [event: MouseEvent]
  'on-full-tree-drag': [event: MouseEvent]
  'stop-full-tree-drag': []
  'on-full-tree-zoom': [event: WheelEvent]
  'close-full-world-tree': []
  'full-tree-node-click': [event: any]
  'show-create-parallel-world': [data: { nodeId: string }]
  'relocate-conversation-to-node': []
  'close-gallery-preview': []
  'copy-to-clipboard': [content: string]
}>()
const nodeSelectOverlay = ref<HTMLElement | null>(null)
const fullWorldTreeOverlay = ref<HTMLElement | null>(null)
function startDrag(e: MouseEvent) {
  emit('start-full-tree-drag', e)
}
function onDrag(e: MouseEvent) {
  emit('on-full-tree-drag', e)
}
function stopDrag() {
  emit('stop-full-tree-drag')
}
function onZoom(e: WheelEvent) {
  emit('on-full-tree-zoom', e)
}
</script>
<style scoped>
.full-worldtree-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--overlay-bg, rgba(0, 0, 0, 0.85));
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.full-worldtree-container {
  transform-origin: center center;
  transition: transform 0.1s ease-out;
}
.full-tree-bottom-buttons {
  position: fixed;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 16px;
  z-index: 1001;
}
.btn-full-tree-action {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-confirm-select {
  background: var(--primary-color);
  color: #fff;
}
.btn-confirm-select:hover:not(.disabled) {
  background: var(--primary-hover);
  transform: scale(1.1);
}
.btn-confirm-select.disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.btn-close-action {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}
.btn-close-action:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.1);
}
.full-tree-bottom-actions {
  position: fixed;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 12px;
  z-index: 1001;
}
.btn-full-tree-bottom {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}
.btn-full-tree-bottom:hover:not(.disabled) {
  transform: scale(1.1);
}
.btn-full-tree-bottom.disabled {
  opacity: 0.4;
  cursor: not-allowed;
}
.btn-add:hover:not(.disabled) {
  background: var(--primary-color);
}
.btn-relocate:hover:not(.disabled) {
  background: #3b82f6;
}
.btn-close:hover {
  background: rgba(255, 255, 255, 0.2);
}
.caps-lock-warning {
  position: fixed;
  top: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: #fbbf24;
  color: #000;
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  z-index: 9999;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}
.safe-mode-decrypt-error {
  position: fixed;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #ef4444;
  color: #fff;
  padding: 16px 24px;
  border-radius: 12px;
  font-size: 14px;
  z-index: 9999;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.gallery-preview-modal {
  background: var(--card-bg);
  border-radius: 12px;
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.preview-modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}
.preview-modal-header h3 {
  margin: 0;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.preview-modal-body {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
  overflow: auto;
}
.preview-modal-body img {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}
.preview-modal-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
}
.preview-info-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  font-size: 13px;
}
.preview-info-row:last-child { margin-bottom: 0; }
.preview-label {
  color: var(--text-secondary);
  min-width: 40px;
}
.preview-value {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.hash-value { font-family: monospace; }
.btn-copy {
  padding: 4px;
  background: transparent;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-copy:hover {
  background: var(--hover-bg);
  color: var(--text-color);
}
.btn-icon {
  width: 28px;
  height: 28px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  color: var(--text-secondary);
}
.btn-icon:hover { background: var(--hover-bg); }
.fade-overlay-enter-active, .fade-overlay-leave-active {
  transition: opacity 0.3s ease;
}
.fade-overlay-enter-from, .fade-overlay-leave-to {
  opacity: 0;
}
.modal-enter-active, .modal-leave-active {
  transition: all 0.3s ease;
}
.modal-enter-from, .modal-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
