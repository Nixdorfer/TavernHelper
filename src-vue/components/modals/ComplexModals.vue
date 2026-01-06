<template>
  <div class="complex-modals">
    <div v-if="showRenameProject" class="modal-overlay" @mousedown.self="$emit('close-rename-project')">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>重命名项目</h3>
          <button class="btn btn-icon" @click="$emit('close-rename-project')">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>项目名称</label>
            <input type="text" class="form-control" :value="renameProjectValue" @input="handleRenameInput" @keyup.enter="$emit('confirm-rename-project')">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-rename-project')">取消</button>
          <button class="btn btn-primary" @click="$emit('confirm-rename-project')">确认</button>
        </div>
      </div>
    </div>
    <div v-if="showKeyManager" class="modal-overlay" @mousedown.self="$emit('close-key-manager')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>密钥管理</h3>
          <button class="btn btn-icon" @click="$emit('close-key-manager')">×</button>
        </div>
        <div class="modal-body">
          <div v-if="!currentProject" class="empty-hint">请先选择项目</div>
          <template v-else>
            <div class="form-group">
              <label>当前项目: {{ currentProject.name }}</label>
            </div>
            <div v-if="projectKeys" class="key-display">
              <div class="form-group">
                <label>公钥</label>
                <div class="key-box">
                  <code class="key-content">{{ projectKeys.publicKey }}</code>
                  <button class="btn btn-secondary btn-sm" @click="$emit('copy-to-clipboard', projectKeys.publicKey, '公钥')">复制</button>
                </div>
              </div>
              <div class="form-group">
                <label>私钥</label>
                <div class="key-box">
                  <code class="key-content">{{ projectKeys.privateKey }}</code>
                  <button class="btn btn-secondary btn-sm" @click="$emit('copy-to-clipboard', projectKeys.privateKey, '私钥')">复制</button>
                </div>
              </div>
              <button class="btn btn-danger" @click="$emit('delete-project-keys')">删除密钥对</button>
            </div>
            <div v-else class="no-keys">
              <p>该项目尚未创建密钥对</p>
              <button class="btn btn-primary" @click="$emit('create-project-keys')">生成密钥对</button>
            </div>
          </template>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-key-manager')">关闭</button>
        </div>
      </div>
    </div>
    <div v-if="showBindNodeModal" class="modal-overlay" @mousedown.self="$emit('close-bind-node')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>绑定节点到远程对话</h3>
          <button class="btn btn-icon" @click="$emit('close-bind-node')">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>当前节点: {{ currentNode?.name }}</label>
          </div>
          <div class="form-group">
            <label>选择目标应用</label>
            <div v-if="!conversationPanelApps.length" class="empty-hint">加载应用列表中...</div>
            <div v-else class="app-select-list">
              <div
                v-for="app in conversationPanelApps"
                :key="app.id"
                :class="['app-select-item', { active: bindTargetApp?.id === app.id }]"
                @click="$emit('select-bind-target-app', app)"
              >
                <span class="app-select-name">{{ app.app.name }}</span>
              </div>
            </div>
          </div>
          <div v-if="bindTargetApp" class="form-group">
            <label>选择目标对话</label>
            <div v-if="bindConversationsLoading" class="empty-hint">加载对话列表中...</div>
            <div v-else-if="!bindConversationsList.length" class="empty-hint">暂无对话</div>
            <div v-else class="app-select-list">
              <div
                v-for="conv in bindConversationsList"
                :key="conv.id"
                :class="['app-select-item', { active: bindTargetConversation?.id === conv.id }]"
                @click="$emit('select-bind-target-conversation', conv)"
              >
                <span class="app-select-name">{{ conv.name || '未命名对话' }}</span>
              </div>
            </div>
          </div>
          <div v-if="bindTargetConversation" class="upload-preview">
            <div class="form-group">
              <label>绑定内容预览</label>
              <div class="preview-stats">
                <div class="stat-item">
                  <span class="stat-label">前置词条目</span>
                  <span class="stat-value">{{ currentNode?.pre_text?.length || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">后置词条目</span>
                  <span class="stat-value">{{ currentNode?.post_text?.length || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">全局词条目</span>
                  <span class="stat-value">{{ currentNode?.pre_prompt?.length || 0 }}</span>
                </div>
                <div class="stat-item">
                  <span class="stat-label">世界书条目</span>
                  <span class="stat-value">{{ worldBookEntriesCount }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-bind-node')">取消</button>
          <button
            class="btn btn-primary"
            :disabled="!bindTargetConversation || bindLoading"
            @click="$emit('bind-node-to-conversation')"
          >
            {{ bindLoading ? '绑定中...' : '绑定' }}
          </button>
        </div>
      </div>
    </div>
    <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="$emit('cancel-delete')">
      <div class="modal modal-sm delete-confirm-modal">
        <div class="modal-header">
          <h3>确认删除</h3>
          <button class="btn btn-icon" @click="$emit('cancel-delete')">×</button>
        </div>
        <div class="modal-body">
          <p class="delete-confirm-text">
            确定要删除对话 <strong>{{ selectedConversation?.name || '未命名对话' }}</strong> 吗？
          </p>
          <p class="delete-confirm-warning">此操作不可恢复！</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('cancel-delete')">取消</button>
          <button class="btn btn-danger" @click="$emit('confirm-delete')">删除</button>
        </div>
      </div>
    </div>
    <div v-if="showWorldTreeCommandPreview" class="modal-overlay" @click.self="$emit('cancel-world-tree-command')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>{{ pendingWorldTreeCommand?.type === 'new_node' ? '新增/修改节点' : '修订设定' }}</h3>
          <button class="btn btn-icon" @click="$emit('cancel-world-tree-command')">×</button>
        </div>
        <div class="modal-body preview-modal-body">
          <div v-if="pendingWorldTreeCommand?.parsed?.desc" class="preview-desc">
            {{ pendingWorldTreeCommand.parsed.desc }}
          </div>
          <div v-if="previewDiff.add.length > 0" class="preview-diff-section">
            <div class="preview-section-header preview-add-header">
              <span class="preview-icon">+</span>
              新增 ({{ selectedPreviewCount('add') }}/{{ previewDiff.add.length }})
            </div>
            <div v-for="(item, idx) in previewDiff.add" :key="'add-'+item.path" :class="['preview-diff-item', { 'preview-item-deselected': item.deselected }]">
              <div class="preview-item-header">
                <label class="preview-toggle" @click.stop>
                  <input type="checkbox" :checked="!item.deselected" @change="$emit('toggle-preview-item', 'add', idx)">
                  <span class="preview-toggle-slider"></span>
                </label>
                <span class="preview-location">{{ item.targetEntryName }}</span>
              </div>
              <textarea class="preview-edit-input" :value="item.content" @input="handlePreviewContentInput('add', idx, $event)" @click.stop rows="3"></textarea>
            </div>
          </div>
          <div v-if="previewDiff.change.length > 0" class="preview-diff-section">
            <div class="preview-section-header preview-change-header">
              <span class="preview-icon">~</span>
              修改 ({{ selectedPreviewCount('change') }}/{{ previewDiff.change.length }})
            </div>
            <div v-for="(item, idx) in previewDiff.change" :key="'change-'+item.serial" :class="['preview-diff-item', { 'preview-item-deselected': item.deselected }]">
              <div class="preview-item-header">
                <label class="preview-toggle" @click.stop>
                  <input type="checkbox" :checked="!item.deselected" @change="$emit('toggle-preview-item', 'change', idx)">
                  <span class="preview-toggle-slider"></span>
                </label>
                <span class="preview-location">{{ item.entryName }}</span>
                <span class="preview-serial-capsule">{{ item.serial }}</span>
              </div>
              <div class="preview-diff-compare">
                <div class="preview-old-row">
                  <span class="diff-label-tag diff-old-tag">旧</span>
                  <span class="diff-content-text">{{ item.oldContent || '(空)' }}</span>
                </div>
                <div class="preview-new-row">
                  <span class="diff-label-tag diff-new-tag">新</span>
                  <textarea class="preview-edit-input preview-edit-inline" :value="item.newContent" @input="handlePreviewNewContentInput('change', idx, $event)" @click.stop rows="2"></textarea>
                </div>
              </div>
            </div>
          </div>
          <div v-if="previewDiff.del.length > 0" class="preview-diff-section">
            <div class="preview-section-header preview-del-header">
              <span class="preview-icon">-</span>
              删除 ({{ selectedPreviewCount('del') }}/{{ previewDiff.del.length }})
            </div>
            <div v-for="(item, idx) in previewDiff.del" :key="'del-'+item.serial" :class="['preview-diff-item', { 'preview-item-deselected': item.deselected }]">
              <div class="preview-item-header">
                <label class="preview-toggle" @click.stop>
                  <input type="checkbox" :checked="!item.deselected" @change="$emit('toggle-preview-item', 'del', idx)">
                  <span class="preview-toggle-slider"></span>
                </label>
                <span class="preview-location">{{ item.entryName }}</span>
                <span class="preview-serial-capsule">{{ item.serial }}</span>
              </div>
              <div class="preview-content-box preview-deleted-box">{{ item.content || '(未找到)' }}</div>
            </div>
          </div>
          <div v-if="!previewDiff.add.length && !previewDiff.change.length && !previewDiff.del.length" class="preview-empty">
            无法识别的变更或内容未找到
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('cancel-world-tree-command')">取消</button>
          <button class="btn btn-danger" @click="$emit('reject-world-tree-command')">拒绝</button>
          <button class="btn btn-primary" @click="$emit('apply-world-tree-command')">应用</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { watch, onMounted, onUnmounted, nextTick } from 'vue'
interface Project {
  id: string
  name: string
}
interface ProjectKeys {
  publicKey: string
  privateKey: string
}
interface Node {
  name: string
  pre_text?: string[]
  post_text?: string[]
  pre_prompt?: string[]
}
interface App {
  id: string
  app: { name: string }
}
interface Conversation {
  id: string
  name?: string
}
interface PreviewItem {
  path?: string
  serial?: string
  content?: string
  oldContent?: string
  newContent?: string
  targetEntryName?: string
  entryName?: string
  deselected?: boolean
}
interface PreviewDiff {
  add: PreviewItem[]
  change: PreviewItem[]
  del: PreviewItem[]
}
interface WorldTreeCommand {
  type: string
  parsed?: { desc?: string }
}
const props = withDefaults(defineProps<{
  showRenameProject?: boolean
  renameProjectValue?: string
  showKeyManager?: boolean
  currentProject?: Project | null
  projectKeys?: ProjectKeys | null
  showBindNodeModal?: boolean
  currentNode?: Node | null
  conversationPanelApps?: App[]
  bindTargetApp?: App | null
  bindConversationsLoading?: boolean
  bindConversationsList?: Conversation[]
  bindTargetConversation?: Conversation | null
  worldBookEntriesCount?: number
  bindLoading?: boolean
  showDeleteConfirm?: boolean
  selectedConversation?: Conversation | null
  showWorldTreeCommandPreview?: boolean
  pendingWorldTreeCommand?: WorldTreeCommand | null
  previewDiff?: PreviewDiff
}>(), {
  showRenameProject: false,
  renameProjectValue: '',
  showKeyManager: false,
  currentProject: null,
  projectKeys: null,
  showBindNodeModal: false,
  currentNode: null,
  conversationPanelApps: () => [],
  bindTargetApp: null,
  bindConversationsLoading: false,
  bindConversationsList: () => [],
  bindTargetConversation: null,
  worldBookEntriesCount: 0,
  bindLoading: false,
  showDeleteConfirm: false,
  selectedConversation: null,
  showWorldTreeCommandPreview: false,
  pendingWorldTreeCommand: null,
  previewDiff: () => ({ add: [], change: [], del: [] })
})
const emit = defineEmits<{
  'close-rename-project': []
  'update:renameProjectValue': [value: string]
  'confirm-rename-project': []
  'close-key-manager': []
  'copy-to-clipboard': [content: string, label: string]
  'delete-project-keys': []
  'create-project-keys': []
  'close-bind-node': []
  'select-bind-target-app': [app: App]
  'select-bind-target-conversation': [conv: Conversation]
  'bind-node-to-conversation': []
  'cancel-delete': []
  'confirm-delete': []
  'cancel-world-tree-command': []
  'reject-world-tree-command': []
  'apply-world-tree-command': []
  'toggle-preview-item': [type: string, idx: number]
  'update-preview-item-content': [type: string, idx: number, value: string]
  'update-preview-item-new-content': [type: string, idx: number, value: string]
}>()
function selectedPreviewCount(type: 'add' | 'change' | 'del'): number {
  if (!props.previewDiff[type]) return 0
  return props.previewDiff[type].filter(item => !item.deselected).length
}
function handleRenameInput(e: Event) {
  emit('update:renameProjectValue', (e.target as HTMLInputElement).value)
}
function handlePreviewContentInput(type: string, idx: number, e: Event) {
  emit('update-preview-item-content', type, idx, (e.target as HTMLTextAreaElement).value)
}
function handlePreviewNewContentInput(type: string, idx: number, e: Event) {
  emit('update-preview-item-new-content', type, idx, (e.target as HTMLTextAreaElement).value)
}
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (props.showRenameProject) emit('close-rename-project')
    if (props.showKeyManager) emit('close-key-manager')
    if (props.showBindNodeModal) emit('close-bind-node')
    if (props.showDeleteConfirm) emit('cancel-delete')
    if (props.showWorldTreeCommandPreview) emit('cancel-world-tree-command')
  }
  if (e.key === 'Enter') {
    const target = e.target as HTMLElement
    if (target.tagName === 'TEXTAREA') return
    if (props.showRenameProject) {
      e.preventDefault()
      emit('confirm-rename-project')
    }
    if (props.showDeleteConfirm) {
      e.preventDefault()
      emit('confirm-delete')
    }
  }
}
function focusFirstInput(selector: string) {
  nextTick(() => {
    const modal = document.querySelector(selector)
    if (!modal) return
    const input = modal.querySelector('input:not([type="hidden"]):not([disabled]), textarea:not([disabled])') as HTMLElement
    if (input) input.focus()
  })
}
watch(() => props.showRenameProject, (val) => { if (val) focusFirstInput('.modal-overlay') })
watch(() => props.showKeyManager, (val) => { if (val) focusFirstInput('.modal-overlay') })
watch(() => props.showBindNodeModal, (val) => { if (val) focusFirstInput('.modal-overlay') })
onMounted(() => window.addEventListener('keydown', handleKeydown))
onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
</script>
<style scoped>
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
.modal {
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.modal-sm { width: 400px; }
.modal-md { width: 560px; }
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}
.modal-header h3 { margin: 0; font-size: 16px; }
.modal-body { padding: 20px; overflow-y: auto; }
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}
.form-group { margin-bottom: 16px; }
.form-group:last-child { margin-bottom: 0; }
.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  color: var(--text-secondary);
}
.form-control {
  width: 100%;
  padding: 10px 12px;
  font-size: 14px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--input-bg);
  color: var(--text-color);
  box-sizing: border-box;
  resize: none;
}
.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
}
.btn {
  padding: 8px 16px;
  font-size: 14px;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-sm { padding: 6px 12px; font-size: 12px; }
.btn-primary { background: var(--primary-color); color: #fff; }
.btn-primary:hover:not(:disabled) { opacity: 0.9; }
.btn-primary:disabled { opacity: 0.6; cursor: not-allowed; }
.btn-secondary { background: var(--secondary-bg); color: var(--text-color); }
.btn-secondary:hover { background: var(--hover-bg); }
.btn-danger { background: #ef4444; color: #fff; }
.btn-danger:hover { background: #dc2626; }
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
  font-size: 18px;
  color: var(--text-secondary);
}
.btn-icon:hover { background: var(--hover-bg); }
.empty-hint {
  text-align: center;
  color: var(--text-muted);
  padding: 20px;
}
.key-display { margin-top: 12px; }
.key-box {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--input-bg);
  padding: 8px 12px;
  border-radius: 8px;
}
.key-content {
  flex: 1;
  font-size: 12px;
  word-break: break-all;
}
.no-keys { text-align: center; padding: 20px 0; }
.app-select-list {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid var(--border-color);
  border-radius: 8px;
}
.app-select-item {
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.2s;
}
.app-select-item:hover { background: var(--hover-bg); }
.app-select-item.active {
  background: var(--primary-light);
  color: var(--primary-color);
}
.app-select-name { font-size: 14px; }
.upload-preview { margin-top: 16px; }
.preview-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  background: var(--input-bg);
  padding: 12px;
  border-radius: 8px;
}
.stat-item {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
}
.stat-label { color: var(--text-secondary); }
.stat-value { font-weight: 500; }
.delete-confirm-text { margin-bottom: 8px; }
.delete-confirm-warning {
  color: #ef4444;
  font-size: 13px;
}
.preview-modal-body {
  max-height: 60vh;
  overflow-y: auto;
}
.preview-desc {
  padding: 12px;
  background: var(--input-bg);
  border-radius: 8px;
  margin-bottom: 16px;
  font-size: 14px;
}
.preview-diff-section { margin-bottom: 16px; }
.preview-section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px 8px 0 0;
  font-size: 13px;
  font-weight: 500;
}
.preview-add-header { background: rgba(34, 197, 94, 0.1); color: #22c55e; }
.preview-change-header { background: rgba(251, 191, 36, 0.1); color: #f59e0b; }
.preview-del-header { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.preview-icon { font-weight: bold; }
.preview-diff-item {
  border: 1px solid var(--border-color);
  border-top: none;
  padding: 12px;
}
.preview-diff-item:last-child { border-radius: 0 0 8px 8px; }
.preview-item-deselected { opacity: 0.5; }
.preview-item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}
.preview-toggle {
  position: relative;
  display: inline-block;
  width: 36px;
  height: 20px;
  cursor: pointer;
}
.preview-toggle input { opacity: 0; width: 0; height: 0; }
.preview-toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--border-color);
  transition: 0.3s;
  border-radius: 20px;
}
.preview-toggle-slider::before {
  position: absolute;
  content: "";
  height: 14px;
  width: 14px;
  left: 3px;
  bottom: 3px;
  background: #fff;
  transition: 0.3s;
  border-radius: 50%;
}
.preview-toggle input:checked + .preview-toggle-slider { background: var(--primary-color); }
.preview-toggle input:checked + .preview-toggle-slider::before { transform: translateX(16px); }
.preview-location { flex: 1; font-size: 13px; }
.preview-serial-capsule {
  padding: 2px 8px;
  background: var(--hover-bg);
  border-radius: 10px;
  font-size: 11px;
}
.preview-edit-input {
  width: 100%;
  padding: 8px 10px;
  font-size: 13px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--input-bg);
  color: var(--text-color);
  box-sizing: border-box;
  resize: none;
}
.preview-diff-compare { margin-top: 8px; }
.preview-old-row, .preview-new-row {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  margin-bottom: 6px;
}
.diff-label-tag {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}
.diff-old-tag { background: rgba(239, 68, 68, 0.1); color: #ef4444; }
.diff-new-tag { background: rgba(34, 197, 94, 0.1); color: #22c55e; }
.diff-content-text {
  flex: 1;
  font-size: 13px;
  color: var(--text-secondary);
}
.preview-edit-inline { flex: 1; }
.preview-content-box {
  padding: 8px 10px;
  background: var(--input-bg);
  border-radius: 6px;
  font-size: 13px;
}
.preview-deleted-box {
  background: rgba(239, 68, 68, 0.05);
  border: 1px solid rgba(239, 68, 68, 0.2);
}
.preview-empty {
  text-align: center;
  color: var(--text-muted);
  padding: 20px;
}
</style>
