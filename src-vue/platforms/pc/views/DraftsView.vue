<template>
  <div class="drafts-container">
    <aside class="drafts-sidebar">
      <div class="panel-header drafts-header" @click="draftsExpanded = !draftsExpanded">
        <span :class="['collapse-icon', { expanded: draftsExpanded }]"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg></span>
        <h2>用户保存内容</h2>
        <button class="btn-add-circle" @click.stop="createNewDraft" title="新建草稿">+</button>
      </div>
      <div v-show="draftsExpanded" class="drafts-list">
        <template v-for="item in draftTree" :key="item.id">
          <div
            v-if="item.isFolder"
            class="draft-folder"
            @dragover.prevent="onDraftDragOver($event, item)"
            @dragleave="onDraftDragLeave($event, item)"
            @drop="onDraftDrop($event, item)"
            :class="{ 'drag-over': dragOverFolderId === item.id }"
          >
            <div
              :class="['draft-folder-header', { expanded: expandedDraftFolders[item.id] }]"
              @click="toggleDraftFolder(item.id)"
            >
              <span class="folder-icon">{{ expandedDraftFolders[item.id] ? '📂' : '📁' }}</span>
              <span class="folder-name">{{ item.name || '未命名文件夹' }}</span>
              <div class="folder-actions" @click.stop>
                <button class="btn-icon" @click="startRenameDraft(item)" title="重命名">✎</button>
                <button class="btn-icon btn-danger" @click="deleteDraftItem(item)" title="删除">×</button>
              </div>
            </div>
            <div v-show="expandedDraftFolders[item.id]" class="folder-children">
              <div
                v-for="child in item.children"
                :key="child.id"
                :class="['draft-item', { active: selectedDraft?.id === child.id, folder: child.isFolder, 'has-pending': !child.isFolder && hasPendingDraftSave(child.id) }]"
                :draggable="!child.isFolder"
                @dragstart="onDraftDragStart($event, child)"
                @dragend="onDraftDragEnd"
                @click="child.isFolder ? toggleDraftFolder(child.id) : selectDraft(child)"
              >
                <span v-if="child.isFolder" class="folder-icon">{{ expandedDraftFolders[child.id] ? '📂' : '📁' }}</span>
                <span class="draft-name">{{ child.name || '未命名' }}</span>
                <div class="draft-actions" @click.stop>
                  <button class="btn-icon" @click="startRenameDraft(child)" title="重命名">✎</button>
                  <button class="btn-icon btn-danger" @click="deleteDraftItem(child)" title="删除">×</button>
                </div>
              </div>
            </div>
          </div>
          <div
            v-else
            :class="['draft-item', { active: selectedDraft?.id === item.id, 'has-pending': hasPendingDraftSave(item.id) }]"
            draggable="true"
            @dragstart="onDraftDragStart($event, item)"
            @dragend="onDraftDragEnd"
            @click="selectDraft(item)"
          >
            <span class="draft-name">{{ item.name || '未命名' }}</span>
            <div class="draft-actions" @click.stop>
              <button class="btn-icon" @click="startRenameDraft(item)" title="重命名">✎</button>
              <button class="btn-icon btn-danger" @click="deleteDraftItem(item)" title="删除">×</button>
            </div>
          </div>
        </template>
        <div v-if="draftTree.length === 0" class="drafts-empty">暂无保存的草稿</div>
      </div>
      <div class="panel-header clipboard-header" @click="clipboardExpanded = !clipboardExpanded">
        <span :class="['collapse-icon', { expanded: clipboardExpanded }]"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg></span>
        <h2>捕获剪贴板</h2>
        <div class="clipboard-header-actions" @click.stop>
          <button v-if="clipboardCaptures.length > 0" class="btn-clear" @click="clearAllClipboardCaptures" title="清空">清空</button>
          <label class="clipboard-switch">
            <input type="checkbox" v-model="clipboardMonitorEnabled" @change="toggleClipboardMonitor">
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>
      <div v-show="clipboardExpanded" class="clipboard-list">
        <div
          v-for="capture in clipboardCaptures"
          :key="capture.id"
          :class="['clipboard-item', { active: selectedClipboardCapture?.id === capture.id, 'has-pending': hasPendingClipboardSave(capture.id) }]"
          @click="selectClipboardCapture(capture)"
        >
          <span class="clipboard-preview">{{ truncateText(capture.content, 50) }}</span>
          <span class="clipboard-time-badge">{{ formatRelativeTime(capture.createdAt) }}</span>
          <button class="btn-save-to-draft" @click.stop="saveClipboardToDraft(capture)" title="保存到草稿">+</button>
        </div>
        <div v-if="clipboardCaptures.length === 0" class="clipboard-empty">
          {{ clipboardMonitorEnabled ? '暂无捕获内容' : '剪贴板监听已关闭' }}
        </div>
      </div>
    </aside>
    <div class="drafts-content">
      <div v-if="!selectedDraft && !selectedClipboardCapture" class="drafts-placeholder">
        <span>请选择一个草稿或剪贴板内容</span>
      </div>
      <div v-else class="draft-editor">
        <div class="draft-editor-header">
          <input
            v-model="draftEditorName"
            class="draft-name-input"
            placeholder="输入名称..."
          />
          <div class="draft-editor-actions">
            <button class="btn btn-primary" @click="saveDraft">保存</button>
            <button class="btn btn-secondary" @click="saveDraftAndCopy">保存并复制</button>
            <button class="btn-close-editor" @click="closeDraftEditor" title="关闭">×</button>
          </div>
        </div>
        <textarea
          v-model="draftEditorContent"
          class="draft-editor-textarea"
          placeholder="输入内容..."
        ></textarea>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, reactive } from 'vue'
import { useDraftsStore, useNotificationStore, useConfirmStore } from '@/stores'
import { draftsApi, type ClipboardCapture } from '@/api/modules/drafts'
import type { Draft } from '@/types'
interface DraftTreeItem extends Draft {
  children?: DraftTreeItem[]
}
const draftsStore = useDraftsStore()
const notificationStore = useNotificationStore()
const confirmStore = useConfirmStore()
const drafts = computed(() => draftsStore.drafts)
const clipboardCaptures = ref<ClipboardCapture[]>([])
const selectedDraft = ref<(Draft & { isNew?: boolean }) | null>(null)
const selectedClipboardCapture = ref<ClipboardCapture | null>(null)
const draftEditorContent = ref('')
const draftEditorName = ref('')
const draftsExpanded = ref(true)
const clipboardExpanded = ref(true)
const expandedDraftFolders = reactive<Record<string, boolean>>({})
const clipboardMonitorEnabled = ref(true)
const draggingDraft = ref<Draft | null>(null)
const dragOverFolderId = ref<string | null>(null)
const pendingSaves = reactive<Record<string, { name: string; content: string; isClipboard?: boolean; captureId?: string }>>({})
const draftTree = computed<DraftTreeItem[]>(() => {
  if (!drafts.value || drafts.value.length === 0) return []
  const rootItems = drafts.value.filter(d => !d.parentId)
  const getChildren = (parentId: string): DraftTreeItem[] => {
    const children = drafts.value.filter(d => d.parentId === parentId)
    return children.map(child => ({
      ...child,
      children: child.isFolder ? getChildren(child.id) : []
    }))
  }
  return rootItems.map(item => ({
    ...item,
    children: item.isFolder ? getChildren(item.id) : []
  }))
})
function truncateText(text: string, maxLen: number): string {
  if (!text) return ''
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
}
function formatRelativeTime(time: string): string {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString()
}
async function loadDrafts() {
  await draftsStore.loadDrafts()
}
async function loadClipboardCaptures() {
  try {
    const result = await draftsApi.getClipboard()
    clipboardCaptures.value = result || []
  } catch (e) {
    console.error('加载剪贴板捕获失败:', e)
    clipboardCaptures.value = []
  }
}
function toggleDraftFolder(folderId: string) {
  expandedDraftFolders[folderId] = !expandedDraftFolders[folderId]
}
function shallowSaveCurrentDraft() {
  if (selectedDraft.value) {
    const hasChanges = draftEditorName.value !== (selectedDraft.value.name || '') ||
                      draftEditorContent.value !== (selectedDraft.value.content || '')
    if (hasChanges) {
      pendingSaves[selectedDraft.value.id] = {
        name: draftEditorName.value,
        content: draftEditorContent.value
      }
    }
  } else if (selectedClipboardCapture.value) {
    const hasChanges = draftEditorContent.value !== (selectedClipboardCapture.value.content || '')
    if (hasChanges) {
      pendingSaves['clip_' + selectedClipboardCapture.value.id] = {
        name: draftEditorName.value,
        content: draftEditorContent.value,
        isClipboard: true,
        captureId: selectedClipboardCapture.value.id
      }
    }
  }
}
function hasPendingDraftSave(id: string): boolean {
  return !!pendingSaves[id]
}
function hasPendingClipboardSave(id: string): boolean {
  return !!pendingSaves['clip_' + id]
}
function selectDraft(draft: Draft) {
  shallowSaveCurrentDraft()
  selectedDraft.value = draft
  selectedClipboardCapture.value = null
  const pending = pendingSaves[draft.id]
  if (pending) {
    draftEditorName.value = pending.name
    draftEditorContent.value = pending.content
  } else {
    draftEditorName.value = draft.name || ''
    draftEditorContent.value = draft.content || ''
  }
}
function selectClipboardCapture(capture: ClipboardCapture) {
  shallowSaveCurrentDraft()
  selectedClipboardCapture.value = capture
  selectedDraft.value = null
  const pending = pendingSaves['clip_' + capture.id]
  if (pending) {
    draftEditorName.value = pending.name
    draftEditorContent.value = pending.content
  } else {
    draftEditorName.value = ''
    draftEditorContent.value = capture.content || ''
  }
}
function getDraftDisplayName(): string {
  if (draftEditorName.value && draftEditorName.value.trim()) {
    return draftEditorName.value.trim()
  }
  const content = draftEditorContent.value || ''
  const firstLine = content.split('\n')[0] || ''
  return firstLine.slice(0, 10).trim() || '未命名'
}
async function saveDraft() {
  try {
    const now = new Date().toISOString()
    const displayName = getDraftDisplayName()
    if (selectedClipboardCapture.value) {
      const captureId = selectedClipboardCapture.value.id
      const draft = await draftsStore.moveClipboardToDraft(captureId, displayName, '')
      if (draft) {
        draft.content = draftEditorContent.value
        draft.name = displayName
        await draftsStore.updateDraft(draft)
        delete pendingSaves['clip_' + captureId]
        selectedDraft.value = draft
        selectedClipboardCapture.value = null
        notificationStore.showNotification('已保存到草稿', 'success')
      }
    } else if (selectedDraft.value) {
      const draftId = selectedDraft.value.id
      if (selectedDraft.value.isNew) {
        const newDraft = await draftsStore.createDraft({
          name: displayName,
          content: draftEditorContent.value,
          parentId: '',
          isFolder: false
        })
        delete pendingSaves[draftId]
        selectedDraft.value = newDraft
        notificationStore.showNotification('草稿已创建', 'success')
      } else {
        const updatedDraft: Draft = {
          ...selectedDraft.value,
          name: displayName,
          content: draftEditorContent.value,
          updatedAt: now
        }
        await draftsStore.updateDraft(updatedDraft)
        delete pendingSaves[draftId]
        selectedDraft.value = updatedDraft
        notificationStore.showNotification('草稿已保存', 'success')
      }
    }
  } catch (e: any) {
    console.error('保存草稿失败:', e)
    notificationStore.showNotification('保存失败: ' + e.message, 'error')
  }
}
async function saveDraftAndCopy() {
  await saveDraft()
  try {
    await draftsApi.copyToClipboard(draftEditorContent.value)
    notificationStore.showNotification('已保存并复制到剪贴板', 'success')
  } catch (e: any) {
    console.error('复制到剪贴板失败:', e)
    notificationStore.showNotification('复制失败: ' + e.message, 'error')
  }
}
function startRenameDraft(item: Draft) {
  const newName = prompt('输入新名称', item.name || '')
  if (newName !== null && newName.trim()) {
    confirmRenameDraft(item, newName.trim())
  }
}
async function confirmRenameDraft(item: Draft, newName: string) {
  try {
    const updated: Draft = {
      ...item,
      name: newName,
      updatedAt: new Date().toISOString()
    }
    await draftsStore.updateDraft(updated)
    if (selectedDraft.value?.id === item.id) {
      selectedDraft.value = updated
      draftEditorName.value = updated.name || ''
    }
    notificationStore.showNotification('重命名成功', 'success')
  } catch (e: any) {
    console.error('重命名失败:', e)
    notificationStore.showNotification('重命名失败: ' + e.message, 'error')
  }
}
async function deleteDraftItem(item: Draft) {
  const itemName = item.name || (item.isFolder ? '未命名文件夹' : '未命名草稿')
  const confirmed = confirm(`确定要删除"${itemName}"吗？${item.isFolder ? '文件夹内的所有内容也将被删除。' : ''}`)
  if (!confirmed) return
  try {
    await draftsStore.deleteDraft(item.id)
    if (selectedDraft.value?.id === item.id) {
      selectedDraft.value = null
      draftEditorName.value = ''
      draftEditorContent.value = ''
    }
    notificationStore.showNotification('已删除', 'success')
  } catch (e: any) {
    console.error('删除失败:', e)
    const msg = e?.message || String(e) || '未知错误'
    notificationStore.showNotification('删除失败: ' + msg, 'error')
  }
}
async function saveClipboardToDraft(capture: ClipboardCapture) {
  try {
    const draft = await draftsStore.moveClipboardToDraft(
      capture.id,
      '剪贴板内容 ' + formatRelativeTime(capture.createdAt),
      ''
    )
    if (draft) {
      await loadClipboardCaptures()
      notificationStore.showNotification('已保存到草稿', 'success')
    }
  } catch (e: any) {
    console.error('保存到草稿失败:', e)
    notificationStore.showNotification('保存失败: ' + e.message, 'error')
  }
}
async function toggleClipboardMonitor() {
  try {
    if (clipboardMonitorEnabled.value) {
      await draftsApi.startClipboardMonitor()
    } else {
      await draftsApi.stopClipboardMonitor()
    }
  } catch (e) {
    console.error('切换剪贴板监听失败:', e)
  }
}
function handleClipboardCaptured(capture: ClipboardCapture) {
  clipboardCaptures.value.unshift(capture)
  if (clipboardCaptures.value.length > 50) {
    clipboardCaptures.value = clipboardCaptures.value.slice(0, 50)
  }
}
function createNewDraft() {
  const now = new Date().toISOString()
  selectedDraft.value = {
    id: 'draft_' + Date.now(),
    name: '',
    content: '',
    parentId: '',
    isFolder: false,
    sortOrder: 0,
    createdAt: now,
    updatedAt: now,
    isNew: true
  }
  selectedClipboardCapture.value = null
  draftEditorName.value = ''
  draftEditorContent.value = ''
}
function onDraftDragStart(event: DragEvent, item: Draft) {
  draggingDraft.value = item
  event.dataTransfer!.effectAllowed = 'move'
  event.dataTransfer!.setData('text/plain', item.id)
}
function onDraftDragOver(event: DragEvent, folder: Draft) {
  if (!draggingDraft.value || draggingDraft.value.id === folder.id) return
  event.preventDefault()
  dragOverFolderId.value = folder.id
}
function onDraftDragLeave(event: DragEvent, folder: Draft) {
  if ((event.currentTarget as HTMLElement).contains(event.relatedTarget as HTMLElement)) return
  if (dragOverFolderId.value === folder.id) {
    dragOverFolderId.value = null
  }
}
async function onDraftDrop(event: DragEvent, folder: Draft) {
  event.preventDefault()
  dragOverFolderId.value = null
  if (!draggingDraft.value || draggingDraft.value.id === folder.id) return
  if (draggingDraft.value.isFolder) return
  try {
    const updated: Draft = {
      ...draggingDraft.value,
      parentId: folder.id,
      updatedAt: new Date().toISOString()
    }
    await draftsStore.updateDraft(updated)
    notificationStore.showNotification('已移动到文件夹', 'success')
  } catch (e: any) {
    console.error('移动草稿失败:', e)
    notificationStore.showNotification('移动失败: ' + e.message, 'error')
  }
  draggingDraft.value = null
}
function onDraftDragEnd() {
  draggingDraft.value = null
  dragOverFolderId.value = null
}
async function clearAllClipboardCaptures() {
  if (clipboardCaptures.value.length === 0) return
  const confirmed = await confirmStore.show({
    title: '清空剪贴板',
    message: `确定要清空所有捕获的剪贴板内容吗？共${clipboardCaptures.value.length}条记录将被删除。`,
    confirmText: '清空',
    cancelText: '取消',
    type: 'danger'
  })
  if (!confirmed) return
  try {
    await draftsApi.clearAllClipboard()
    clipboardCaptures.value = []
    if (selectedClipboardCapture.value) {
      selectedClipboardCapture.value = null
      draftEditorContent.value = ''
    }
    notificationStore.showNotification('已清空剪贴板', 'success')
  } catch (e: any) {
    console.error('清空剪贴板失败:', e)
    notificationStore.showNotification('清空失败: ' + e.message, 'error')
  }
}
function closeDraftEditor() {
  shallowSaveCurrentDraft()
  selectedDraft.value = null
  selectedClipboardCapture.value = null
  draftEditorName.value = ''
  draftEditorContent.value = ''
}
onMounted(async () => {
  await loadDrafts()
  await loadClipboardCaptures()
  if (clipboardMonitorEnabled.value) {
    await draftsApi.startClipboardMonitor()
  }
  if (typeof window !== 'undefined' && (window as any).runtime) {
    (window as any).runtime.EventsOn('clipboard-captured', handleClipboardCaptured)
  }
})
onBeforeUnmount(() => {
  if (typeof window !== 'undefined' && (window as any).runtime) {
    (window as any).runtime.EventsOff('clipboard-captured')
  }
})
</script>
<style scoped>
.drafts-container {
  display: flex;
  flex: 1;
  overflow: hidden;
}
.drafts-sidebar {
  width: 320px;
  min-width: 320px;
  min-height: 400px;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
}
.panel-header {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  gap: 8px;
}
.panel-header h2 {
  flex: 1;
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}
.collapse-icon {
  width: 16px;
  height: 16px;
  transition: transform 0.2s;
}
.collapse-icon.expanded {
  transform: rotate(90deg);
}
.collapse-icon svg {
  width: 100%;
  height: 100%;
}
.drafts-list,
.clipboard-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.draft-item,
.clipboard-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  gap: 8px;
  margin-bottom: 4px;
}
.draft-item:hover,
.clipboard-item:hover {
  background: var(--bg-hover);
}
.draft-item.active,
.clipboard-item.active {
  background: var(--primary-color);
  color: white;
}
.draft-item.has-pending,
.clipboard-item.has-pending {
  background: rgba(255, 193, 7, 0.25);
}
.draft-item.has-pending:not(.active),
.clipboard-item.has-pending:not(.active) {
  border-left: 3px solid #ffc107;
}
.draft-name,
.clipboard-preview {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 13px;
}
.draft-actions,
.folder-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.2s;
}
.draft-item:hover .draft-actions,
.draft-folder-header:hover .folder-actions {
  opacity: 1;
}
.btn-icon {
  padding: 2px 6px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 4px;
  font-size: 12px;
}
.btn-icon:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.btn-icon.btn-danger:hover {
  background: #ef4444;
  color: white;
}
.draft-folder {
  margin-bottom: 4px;
}
.draft-folder-header {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  gap: 8px;
}
.draft-folder-header:hover {
  background: var(--bg-hover);
}
.folder-icon {
  font-size: 14px;
}
.folder-name {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
}
.folder-children {
  padding-left: 20px;
}
.drafts-empty,
.clipboard-empty {
  padding: 20px;
  text-align: center;
  color: var(--text-tertiary);
  font-size: 13px;
}
.clipboard-switch {
  position: relative;
  width: 36px;
  height: 20px;
}
.clipboard-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.clipboard-switch .switch-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--bg-tertiary);
  transition: 0.3s;
  border-radius: 20px;
}
.clipboard-switch .switch-slider:before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}
.clipboard-switch input:checked + .switch-slider {
  background-color: var(--primary-color);
}
.clipboard-switch input:checked + .switch-slider:before {
  transform: translateX(16px);
}
.btn-add-circle {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
.btn-add-circle:hover {
  filter: brightness(1.2);
}
.clipboard-header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
.btn-clear {
  padding: 4px 8px;
  font-size: 11px;
  background: transparent;
  color: var(--text-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
  cursor: pointer;
}
.btn-clear:hover {
  background: #ef4444;
  border-color: #ef4444;
  color: white;
}
.clipboard-time-badge {
  font-size: 11px;
  color: var(--text-tertiary);
  background: var(--bg-tertiary);
  padding: 2px 8px;
  border-radius: 10px;
  white-space: nowrap;
}
.btn-save-to-draft {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
  flex-shrink: 0;
}
.btn-save-to-draft:hover {
  filter: brightness(1.2);
}
.draft-folder.drag-over {
  background: rgba(34, 197, 94, 0.15);
  border-radius: 8px;
}
.draft-item[draggable="true"]:active {
  cursor: grabbing;
}
.drafts-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
}
.drafts-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}
.draft-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
  gap: 16px;
  overflow: hidden;
}
.draft-editor-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}
.draft-name-input {
  flex: 1;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 14px;
  resize: none;
}
.draft-name-input:focus {
  outline: none;
  border-color: var(--primary-color);
}
.draft-editor-actions {
  display: flex;
  gap: 8px;
}
.draft-editor-textarea {
  flex: 1;
  padding: 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.6;
  resize: none;
}
.draft-editor-textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}
.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  border: none;
}
.btn-primary {
  background: var(--primary-color);
  color: white;
}
.btn-primary:hover {
  filter: brightness(1.1);
}
.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}
.btn-secondary:hover {
  background: var(--bg-secondary);
}
.btn-close-editor {
  width: 28px;
  height: 28px;
  border-radius: 4px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border: none;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-close-editor:hover {
  background: #ef4444;
  color: white;
}
</style>
