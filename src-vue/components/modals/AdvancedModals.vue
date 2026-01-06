<template>
  <div class="advanced-modals">
    <div v-if="showModSelector" class="modal-overlay" @click.self="$emit('close-mod-selector')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>选择 Mod</h3>
          <button class="btn btn-icon" @click="$emit('close-mod-selector')">×</button>
        </div>
        <div class="modal-body">
          <div v-if="loadingMods" class="loading-state">
            <div class="loading-spinner"></div>
            <span>加载 Mod 列表...</span>
          </div>
          <div v-else-if="modList.length === 0" class="empty-state">
            <p>暂无可用的 Mod</p>
          </div>
          <div v-else class="mod-list">
            <div
              v-for="mod in modList"
              :key="mod.id"
              :class="['mod-item', { selected: isModSelected(mod) }]"
              @click="$emit('toggle-mod-selection', mod)"
            >
              <div class="mod-info">
                <div class="mod-name">{{ mod.name }}</div>
                <div class="mod-description" v-if="mod.description">{{ mod.description }}</div>
                <div class="mod-version" v-if="mod.versions && mod.versions.length > 0">
                  <select :value="mod.selectedVersion" @change="handleModVersionChange(mod, $event)" @click.stop class="version-select">
                    <option v-for="ver in mod.versions" :key="ver.id" :value="ver.id">
                      {{ ver.version }}
                    </option>
                  </select>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-mod-selector')">取消</button>
          <button class="btn btn-primary" @click="$emit('save-mods')" :disabled="savingMods">
            {{ savingMods ? '保存中...' : '保存' }}
          </button>
        </div>
      </div>
    </div>
    <div v-if="showT2iSelectorModal" class="modal-overlay" @click.self="$emit('close-t2i-selector')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>选择要生成的图片</h3>
          <button class="btn btn-icon" @click="$emit('close-t2i-selector')">×</button>
        </div>
        <div class="modal-body t2i-modal-body">
          <div v-if="t2iScenes.length === 0" class="empty-state">
            <p>没有可用的场景</p>
          </div>
          <div v-else class="t2i-scene-list">
            <div
              v-for="(scene, index) in t2iScenes"
              :key="index"
              :class="['t2i-scene-item', { selected: scene.selected, invalid: !scene.valid }]"
              @click="scene.valid && $emit('toggle-t2i-scene', index)"
            >
              <div class="t2i-scene-content">
                <textarea
                  :value="scene.prompt"
                  @input="handleT2iSceneInput(index, $event)"
                  class="t2i-scene-textarea"
                  rows="3"
                  @click.stop
                ></textarea>
              </div>
              <div v-if="!scene.valid" class="t2i-scene-invalid">缺少角色生图关键字</div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-t2i-selector')">取消</button>
          <button class="btn btn-primary" @click="$emit('submit-t2i-scenes')" :disabled="!hasSelectedT2iScenes">
            生成选中的图片 ({{ selectedT2iScenesCount }})
          </button>
        </div>
      </div>
    </div>
    <div v-if="showImagePromptModal" class="modal-overlay" @click.self="$emit('close-image-prompt')">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>生成图片</h3>
          <button class="btn btn-icon" @click="$emit('close-image-prompt')">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>图片描述</label>
            <textarea
              :value="imagePromptText"
              @input="handleImagePromptInput"
              class="form-control image-prompt-textarea"
              rows="6"
              placeholder="请输入图片描述..."
              @keydown.enter.ctrl="$emit('submit-image-prompt')"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-image-prompt')">取消</button>
          <button class="btn btn-primary" @click="$emit('submit-image-prompt')" :disabled="!imagePromptText.trim() || submittingImagePrompt">
            {{ submittingImagePrompt ? '生成中...' : '生成' }}
          </button>
        </div>
      </div>
    </div>
    <div v-if="showUnbindRemote" class="modal-overlay" @mousedown.self="$emit('close-unbind-remote')">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>解除远程绑定</h3>
          <button class="btn btn-icon" @click="$emit('close-unbind-remote')">×</button>
        </div>
        <div class="modal-body">
          <p class="text-secondary">确定要解除与远程应用的绑定吗？</p>
          <div class="form-group" style="margin-top: 16px;">
            <div class="toggle-row">
              <label class="toggle-label">
                <input type="checkbox" :checked="unbindClearRemoteConfig" @change="handleUnbindCheckbox">
                <span class="toggle-slider"></span>
                <span>同时清空远程配置</span>
              </label>
            </div>
            <p class="setting-hint">开启后将向远程发送空配置，关闭则仅删除本地绑定关系</p>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-unbind-remote')">取消</button>
          <button class="btn btn-danger" @click="$emit('confirm-unbind-remote')">确认解绑</button>
        </div>
      </div>
    </div>
    <div v-if="showEditNodeName" class="modal-overlay" @mousedown.self="$emit('close-edit-node-name')">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>重命名</h3>
          <button class="btn btn-icon" @click="$emit('close-edit-node-name')">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>节点名称</label>
            <input type="text" class="form-control" :value="editingNodeName" @input="handleNodeNameInput" @keyup.enter="$emit('save-node-name')">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-edit-node-name')">取消</button>
          <button class="btn btn-primary" @click="$emit('save-node-name')">保存</button>
        </div>
      </div>
    </div>
    <div v-if="showEditNodeDescription" class="modal-overlay" @mousedown.self="$emit('close-edit-node-description')">
      <div class="modal">
        <div class="modal-header">
          <h3>描述</h3>
          <button class="btn btn-icon" @click="$emit('close-edit-node-description')">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <textarea
              class="form-control"
              :value="editingNodeDescription"
              @input="handleNodeDescriptionInput"
              rows="6"
              placeholder="暂无描述"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-edit-node-description')">取消</button>
          <button class="btn btn-primary" @click="$emit('save-node-description')">保存</button>
        </div>
      </div>
    </div>
    <div v-if="showDebugReplyModal" class="modal-overlay" @mousedown.self="$emit('close-debug-reply')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>设定测试回复</h3>
          <button class="btn btn-icon" @click="$emit('close-debug-reply')">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>回复内容</label>
            <textarea
              class="form-control"
              :value="debugTestReply"
              @input="handleDebugReplyInput"
              rows="10"
              placeholder="输入测试回复内容..."
            ></textarea>
          </div>
          <div class="form-group toggle-group">
            <label>返回收到的消息</label>
            <label class="toggle-switch">
              <input type="checkbox" :checked="debugEchoMessage" @change="handleDebugEchoChange">
              <span class="toggle-slider"></span>
            </label>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-debug-reply')">关闭</button>
        </div>
      </div>
    </div>
    <div v-if="showCreateProjectFromRemote" class="modal-overlay" @click.self="$emit('cancel-create-from-remote')">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>新建项目</h3>
          <button class="btn btn-icon" @click="$emit('cancel-create-from-remote')">×</button>
        </div>
        <div class="modal-body">
          <p style="margin-bottom: 12px; color: var(--text-secondary);">远端数据已获取，请输入项目名称来创建新项目：</p>
          <input
            type="text"
            class="form-control"
            :value="newProjectNameFromRemote"
            @input="handleNewProjectInput"
            placeholder="项目名称"
            @keydown.enter="$emit('confirm-create-from-remote')"
          >
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('cancel-create-from-remote')">取消</button>
          <button class="btn btn-primary" @click="$emit('confirm-create-from-remote')">创建</button>
        </div>
      </div>
    </div>
    <div v-if="showT2iModal" class="modal-overlay t2i-overlay" @mousedown.self="$emit('close-t2i-modal')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>生成图片</h3>
          <button class="btn btn-icon" @click="$emit('close-t2i-modal')">×</button>
        </div>
        <div class="modal-body t2i-modal-body">
          <div v-if="t2iPrompts.length === 0" class="t2i-empty">暂无生图提示词</div>
          <div v-for="(prompt, idx) in t2iPrompts" :key="idx" class="t2i-prompt-item">
            <label class="t2i-checkbox">
              <input type="checkbox" :checked="t2iSelectedPrompts.includes(idx)" @change="$emit('toggle-t2i-prompt', idx)">
              <span class="t2i-checkbox-mark"></span>
            </label>
            <textarea
              class="form-control t2i-prompt-input"
              :value="prompt"
              @input="handleT2iPromptInput(idx, $event)"
              rows="2"
            ></textarea>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="$emit('close-t2i-modal')">取消</button>
          <button class="btn btn-primary t2i-generate-btn" @click="$emit('generate-t2i-images')" :disabled="t2iSelectedPrompts.length === 0">生成</button>
        </div>
      </div>
    </div>
    <div v-if="showWorldTreeSystemPrompt" class="modal-overlay" @mousedown.self="$emit('close-world-tree-system-prompt')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>世界树系统词</h3>
          <button class="btn btn-icon" @click="$emit('close-world-tree-system-prompt')">×</button>
        </div>
        <div class="modal-body">
          <p style="margin-bottom: 12px; color: var(--text-secondary);">开启世界树时，此内容将添加在用户输入前发送</p>
          <textarea
            class="form-control"
            :value="worldTreeSystemPrompt"
            @input="handleWorldTreePromptInput"
            placeholder="输入世界树系统提示词..."
            rows="8"
            style="resize: none;"
          ></textarea>
        </div>
        <div class="modal-footer">
          <button class="btn btn-primary" @click="$emit('close-world-tree-system-prompt')">确定</button>
        </div>
      </div>
    </div>
    <div v-if="showProjectSelectForWorldTree" class="modal-overlay" @mousedown.self="$emit('close-project-select-for-world-tree')">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>选择项目</h3>
          <button class="btn btn-icon" @click="$emit('close-project-select-for-world-tree')">×</button>
        </div>
        <div class="modal-body">
          <div v-if="!projects || projects.length === 0" class="empty-hint">暂无项目，请先在项目管理中创建项目</div>
          <div v-else class="project-list-bind">
            <div
              v-for="project in projects"
              :key="project.id"
              class="project-item-bind"
              @click="$emit('select-project-for-world-tree', project)"
            >
              <span class="project-name">{{ project.name }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { watch, onMounted, onUnmounted, nextTick } from 'vue'
interface Mod {
  id: string
  name: string
  description?: string
  versions?: { id: string; version: string }[]
  selectedVersion?: string
}
interface T2iScene {
  prompt: string
  selected: boolean
  valid: boolean
}
interface Project {
  id: string
  name: string
}
const props = withDefaults(defineProps<{
  showModSelector?: boolean
  showT2iSelectorModal?: boolean
  showImagePromptModal?: boolean
  showUnbindRemote?: boolean
  showEditNodeName?: boolean
  showEditNodeDescription?: boolean
  showDebugReplyModal?: boolean
  showCreateProjectFromRemote?: boolean
  showT2iModal?: boolean
  showWorldTreeSystemPrompt?: boolean
  showProjectSelectForWorldTree?: boolean
  loadingMods?: boolean
  modList?: Mod[]
  savingMods?: boolean
  t2iScenes?: T2iScene[]
  hasSelectedT2iScenes?: boolean
  selectedT2iScenesCount?: number
  imagePromptText?: string
  submittingImagePrompt?: boolean
  unbindClearRemoteConfig?: boolean
  editingNodeName?: string
  editingNodeDescription?: string
  debugTestReply?: string
  debugEchoMessage?: boolean
  newProjectNameFromRemote?: string
  t2iPrompts?: string[]
  t2iSelectedPrompts?: number[]
  worldTreeSystemPrompt?: string
  projects?: Project[]
  isModSelected?: (mod: Mod) => boolean
}>(), {
  showModSelector: false,
  showT2iSelectorModal: false,
  showImagePromptModal: false,
  showUnbindRemote: false,
  showEditNodeName: false,
  showEditNodeDescription: false,
  showDebugReplyModal: false,
  showCreateProjectFromRemote: false,
  showT2iModal: false,
  showWorldTreeSystemPrompt: false,
  showProjectSelectForWorldTree: false,
  loadingMods: false,
  modList: () => [],
  savingMods: false,
  t2iScenes: () => [],
  hasSelectedT2iScenes: false,
  selectedT2iScenesCount: 0,
  imagePromptText: '',
  submittingImagePrompt: false,
  unbindClearRemoteConfig: false,
  editingNodeName: '',
  editingNodeDescription: '',
  debugTestReply: '',
  debugEchoMessage: false,
  newProjectNameFromRemote: '',
  t2iPrompts: () => [],
  t2iSelectedPrompts: () => [],
  worldTreeSystemPrompt: '',
  projects: () => [],
  isModSelected: () => false
})
const emit = defineEmits<{
  'close-mod-selector': []
  'toggle-mod-selection': [mod: Mod]
  'update-mod-version': [data: { mod: Mod; version: string }]
  'save-mods': []
  'close-t2i-selector': []
  'toggle-t2i-scene': [index: number]
  'update-t2i-scene-prompt': [data: { index: number; value: string }]
  'submit-t2i-scenes': []
  'close-image-prompt': []
  'update:imagePromptText': [value: string]
  'submit-image-prompt': []
  'close-unbind-remote': []
  'update:unbindClearRemoteConfig': [value: boolean]
  'confirm-unbind-remote': []
  'close-edit-node-name': []
  'update:editingNodeName': [value: string]
  'save-node-name': []
  'close-edit-node-description': []
  'update:editingNodeDescription': [value: string]
  'save-node-description': []
  'close-debug-reply': []
  'update:debugTestReply': [value: string]
  'update:debugEchoMessage': [value: boolean]
  'cancel-create-from-remote': []
  'update:newProjectNameFromRemote': [value: string]
  'confirm-create-from-remote': []
  'close-t2i-modal': []
  'toggle-t2i-prompt': [idx: number]
  'update-t2i-prompt': [data: { idx: number; value: string }]
  'generate-t2i-images': []
  'close-world-tree-system-prompt': []
  'update:worldTreeSystemPrompt': [value: string]
  'close-project-select-for-world-tree': []
  'select-project-for-world-tree': [project: Project]
}>()
function handleModVersionChange(mod: Mod, e: Event) {
  emit('update-mod-version', { mod, version: (e.target as HTMLSelectElement).value })
}
function handleT2iSceneInput(index: number, e: Event) {
  emit('update-t2i-scene-prompt', { index, value: (e.target as HTMLTextAreaElement).value })
}
function handleImagePromptInput(e: Event) {
  emit('update:imagePromptText', (e.target as HTMLTextAreaElement).value)
}
function handleUnbindCheckbox(e: Event) {
  emit('update:unbindClearRemoteConfig', (e.target as HTMLInputElement).checked)
}
function handleNodeNameInput(e: Event) {
  emit('update:editingNodeName', (e.target as HTMLInputElement).value)
}
function handleNodeDescriptionInput(e: Event) {
  emit('update:editingNodeDescription', (e.target as HTMLTextAreaElement).value)
}
function handleDebugReplyInput(e: Event) {
  emit('update:debugTestReply', (e.target as HTMLTextAreaElement).value)
}
function handleDebugEchoChange(e: Event) {
  emit('update:debugEchoMessage', (e.target as HTMLInputElement).checked)
}
function handleNewProjectInput(e: Event) {
  emit('update:newProjectNameFromRemote', (e.target as HTMLInputElement).value)
}
function handleT2iPromptInput(idx: number, e: Event) {
  emit('update-t2i-prompt', { idx, value: (e.target as HTMLTextAreaElement).value })
}
function handleWorldTreePromptInput(e: Event) {
  emit('update:worldTreeSystemPrompt', (e.target as HTMLTextAreaElement).value)
}
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (props.showModSelector) emit('close-mod-selector')
    if (props.showT2iSelectorModal) emit('close-t2i-selector')
    if (props.showImagePromptModal) emit('close-image-prompt')
    if (props.showUnbindRemote) emit('close-unbind-remote')
    if (props.showEditNodeName) emit('close-edit-node-name')
    if (props.showEditNodeDescription) emit('close-edit-node-description')
    if (props.showDebugReplyModal) emit('close-debug-reply')
    if (props.showCreateProjectFromRemote) emit('cancel-create-from-remote')
    if (props.showT2iModal) emit('close-t2i-modal')
    if (props.showWorldTreeSystemPrompt) emit('close-world-tree-system-prompt')
    if (props.showProjectSelectForWorldTree) emit('close-project-select-for-world-tree')
  }
  if (e.key === 'Enter') {
    const target = e.target as HTMLElement
    if (target.tagName === 'TEXTAREA') return
    if (props.showEditNodeName) {
      e.preventDefault()
      emit('save-node-name')
    }
    if (props.showCreateProjectFromRemote) {
      e.preventDefault()
      emit('confirm-create-from-remote')
    }
    if (props.showUnbindRemote) {
      e.preventDefault()
      emit('confirm-unbind-remote')
    }
  }
}
function focusFirstInput(selector: string) {
  nextTick(() => {
    const modal = document.querySelector(selector)
    if (!modal) return
    const input = modal.querySelector('input:not([type="hidden"]):not([type="checkbox"]):not([disabled]), textarea:not([disabled])') as HTMLElement
    if (input) input.focus()
  })
}
watch(() => props.showEditNodeName, (val) => { if (val) focusFirstInput('.modal-overlay') })
watch(() => props.showEditNodeDescription, (val) => { if (val) focusFirstInput('.modal-overlay') })
watch(() => props.showCreateProjectFromRemote, (val) => { if (val) focusFirstInput('.modal-overlay') })
watch(() => props.showImagePromptModal, (val) => { if (val) focusFirstInput('.modal-overlay') })
watch(() => props.showWorldTreeSystemPrompt, (val) => { if (val) focusFirstInput('.modal-overlay') })
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
.loading-state, .empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px;
  color: var(--text-secondary);
}
.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-bottom: 12px;
}
@keyframes spin { to { transform: rotate(360deg); } }
.mod-list { max-height: 400px; overflow-y: auto; }
.mod-item {
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.mod-item:hover { background: var(--hover-bg); }
.mod-item.selected {
  border-color: var(--primary-color);
  background: var(--primary-light);
}
.mod-name { font-weight: 500; margin-bottom: 4px; }
.mod-description { font-size: 13px; color: var(--text-secondary); }
.mod-version { margin-top: 8px; }
.version-select {
  padding: 4px 8px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--input-bg);
  color: var(--text-color);
  font-size: 12px;
}
.t2i-modal-body { max-height: 400px; overflow-y: auto; }
.t2i-scene-list { display: flex; flex-direction: column; gap: 8px; }
.t2i-scene-item {
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
}
.t2i-scene-item.selected { border-color: var(--primary-color); }
.t2i-scene-item.invalid { opacity: 0.6; cursor: not-allowed; }
.t2i-scene-textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--input-bg);
  color: var(--text-color);
  font-size: 13px;
  resize: none;
}
.t2i-scene-invalid {
  margin-top: 8px;
  font-size: 12px;
  color: #ef4444;
}
.toggle-row { display: flex; align-items: center; gap: 8px; }
.toggle-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
.toggle-group {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
}
.toggle-switch input { opacity: 0; width: 0; height: 0; }
.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--border-color);
  transition: 0.3s;
  border-radius: 22px;
}
.toggle-slider::before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 3px;
  bottom: 3px;
  background: #fff;
  transition: 0.3s;
  border-radius: 50%;
}
input:checked + .toggle-slider { background: var(--primary-color); }
input:checked + .toggle-slider::before { transform: translateX(18px); }
.setting-hint { font-size: 12px; color: var(--text-muted); margin-top: 8px; }
.text-secondary { color: var(--text-secondary); }
.t2i-empty { text-align: center; padding: 20px; color: var(--text-secondary); }
.t2i-prompt-item { display: flex; gap: 12px; margin-bottom: 12px; }
.t2i-checkbox {
  display: flex;
  align-items: center;
  cursor: pointer;
}
.t2i-prompt-input { flex: 1; }
.empty-hint { text-align: center; color: var(--text-muted); padding: 20px; }
.project-list-bind { display: flex; flex-direction: column; gap: 8px; }
.project-item-bind {
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.project-item-bind:hover { background: var(--hover-bg); }
.project-name { font-weight: 500; }
</style>
