<template>
  <div class="debug-panel">
    <div class="debug-header">
      <h2>API 调试面板</h2>
      <input v-model="searchQuery" class="search-input" placeholder="搜索接口..." />
    </div>
    <div class="debug-content">
      <div class="api-list">
        <div
          v-for="api in filteredApis"
          :key="api.name"
          :class="['api-item', { active: selectedApi?.name === api.name }]"
          @click="selectApi(api)"
        >
          <span class="api-name">{{ api.name }}</span>
          <span class="api-file">{{ api.file }}</span>
        </div>
      </div>
      <div class="api-detail">
        <template v-if="selectedApi">
          <div class="api-detail-header">
            <h3>{{ selectedApi.name }}</h3>
            <span class="api-signature">{{ selectedApi.signature }}</span>
          </div>
          <div class="api-params">
            <div class="param-label">请求参数 (JSON)</div>
            <textarea v-model="requestParams" class="param-input" placeholder='例如: ["arg1", 123, {"key": "value"}]'></textarea>
          </div>
          <div class="api-actions">
            <button class="btn-call" @click="callApi" :disabled="calling">
              {{ calling ? '调用中...' : '调用接口' }}
            </button>
            <button class="btn-clear" @click="clearResult">清空结果</button>
          </div>
          <div class="api-result">
            <div class="result-label">返回结果</div>
            <div v-if="callError" class="result-error">{{ callError }}</div>
            <pre v-else class="result-content">{{ resultDisplay }}</pre>
          </div>
        </template>
        <div v-else class="no-selection">
          请选择一个接口进行测试
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { invoke } from '@tauri-apps/api/core'
interface ApiItem {
  name: string
  command: string
  file: string
  signature: string
  params: string[]
}
const API_LIST: ApiItem[] = [
  { name: 'IsDebugMode', command: 'is_debug_mode', file: 'system.rs', signature: '() -> bool', params: [] },
  { name: 'LoadConfig', command: 'load_config', file: 'config.rs', signature: '() -> Result<Config>', params: [] },
  { name: 'SaveConfig', command: 'save_config', file: 'config.rs', signature: '(config: Config) -> Result<()>', params: ['config'] },
  { name: 'UpdateConfig', command: 'update_config', file: 'config.rs', signature: '(key: &str, value: Value) -> Result<()>', params: ['key', 'value'] },
  { name: 'SaveSessionState', command: 'save_session_state', file: 'config.rs', signature: '(state: SessionState) -> Result<()>', params: ['state'] },
  { name: 'LoadSessionState', command: 'load_session_state', file: 'config.rs', signature: '() -> Result<SessionState>', params: [] },
  { name: 'GetProjects', command: 'get_projects', file: 'project.rs', signature: '() -> Result<Vec<ProjectInfo>>', params: [] },
  { name: 'LoadProjectByName', command: 'load_project_by_name', file: 'project.rs', signature: '(name: &str) -> Result<Project>', params: ['name'] },
  { name: 'CreateProject', command: 'create_project', file: 'project.rs', signature: '(name: &str, projectType: &str) -> Result<i32>', params: ['name', 'projectType'] },
  { name: 'DeleteProject', command: 'delete_project', file: 'project.rs', signature: '(fileName: &str) -> Result<()>', params: ['fileName'] },
  { name: 'RenameProject', command: 'rename_project', file: 'project.rs', signature: '(oldName: &str, newName: &str) -> Result<()>', params: ['oldName', 'newName'] },
  { name: 'CreateChildNode', command: 'create_child_node', file: 'node.rs', signature: '(projectName: &str, parentId: i32, name: &str) -> Result<NodeInfo>', params: ['projectName', 'parentId', 'name'] },
  { name: 'CreateBrotherNode', command: 'create_brother_node', file: 'node.rs', signature: '(projectName: &str, siblingId: i32, name: &str) -> Result<NodeInfo>', params: ['projectName', 'siblingId', 'name'] },
  { name: 'DeleteNode', command: 'delete_node', file: 'node.rs', signature: '(projectName: &str, nodeId: i32) -> Result<()>', params: ['projectName', 'nodeId'] },
  { name: 'RenameNode', command: 'rename_node', file: 'node.rs', signature: '(nodeId: i32, newName: &str) -> Result<()>', params: ['nodeId', 'newName'] },
  { name: 'GetGalleryImages', command: 'get_gallery_images', file: 'gallery.rs', signature: '() -> Result<Vec<GalleryImage>>', params: [] },
  { name: 'GetGalleryFolders', command: 'get_gallery_folders', file: 'gallery.rs', signature: '() -> Result<Vec<GalleryFolder>>', params: [] },
  { name: 'DeleteGalleryImage', command: 'delete_gallery_image', file: 'gallery.rs', signature: '(id: &str) -> Result<()>', params: ['id'] },
  { name: 'GetAllDrafts', command: 'get_all_drafts', file: 'drafts.rs', signature: '() -> Result<Vec<Draft>>', params: [] },
  { name: 'CreateDraft', command: 'create_draft', file: 'drafts.rs', signature: '(draft: Draft) -> Result<()>', params: ['draft'] },
  { name: 'DeleteDraft', command: 'delete_draft', file: 'drafts.rs', signature: '(id: i32) -> Result<()>', params: ['id'] },
  { name: 'GetClipboardCaptures', command: 'get_clipboard_captures', file: 'drafts.rs', signature: '() -> Result<Vec<ClipboardCapture>>', params: [] },
  { name: 'CopyToClipboard', command: 'copy_to_clipboard', file: 'clipboard.rs', signature: '(content: &str) -> Result<()>', params: ['content'] },
  { name: 'GetLocalCreations', command: 'get_local_creations', file: 'creation.rs', signature: '() -> Result<Vec<Value>>', params: [] },
  { name: 'SendTestChat', command: 'send_test_chat', file: 'creation.rs', signature: '(request: TestChatRequest) -> TestChatResponse', params: ['request'] },
  { name: 'GetNodeContent', command: 'get_node_content', file: 'worldtree.rs', signature: '(nodeId: i32) -> Result<NodeContent>', params: ['nodeId'] },
  { name: 'AddFolder', command: 'add_folder', file: 'worldtree.rs', signature: '(nodeId: i32, name: &str) -> Result<i32>', params: ['nodeId', 'name'] },
  { name: 'AddCard', command: 'add_card', file: 'worldtree.rs', signature: '(nodeId: i32, folderChangeId: i32, name: &str, keyWord: &str) -> Result<i32>', params: ['nodeId', 'folderChangeId', 'name', 'keyWord'] },
  { name: 'AuthLogin', command: 'auth_login', file: 'auth.rs', signature: '(email: &str, password: &str, rememberMe: bool) -> Result<Value>', params: ['email', 'password', 'rememberMe'] },
  { name: 'AuthGetProfile', command: 'auth_get_profile', file: 'auth.rs', signature: '(token: &str) -> Result<Value>', params: ['token'] },
  { name: 'AuthLogout', command: 'auth_logout', file: 'auth.rs', signature: '(token: &str) -> Result<()>', params: ['token'] },
  { name: 'FetchWithAuth', command: 'fetch_with_auth', file: 'creation.rs', signature: '(token: &str, url: &str, method: &str, body: Option<&str>) -> Result<String>', params: ['token', 'url', 'method', 'body'] },
]
const searchQuery = ref('')
const apis = ref(API_LIST)
const selectedApi = ref<ApiItem | null>(null)
const requestParams = ref('[]')
const callResult = ref<any>(null)
const callError = ref<string | null>(null)
const calling = ref(false)
const filteredApis = computed(() => {
  if (!searchQuery.value) return apis.value
  const q = searchQuery.value.toLowerCase()
  return apis.value.filter(api =>
    api.name.toLowerCase().includes(q) ||
    api.file.toLowerCase().includes(q)
  )
})
const resultDisplay = computed(() => {
  if (callResult.value === null) return '等待调用...'
  try {
    return JSON.stringify(callResult.value, null, 2)
  } catch {
    return String(callResult.value)
  }
})
function selectApi(api: ApiItem) {
  selectedApi.value = api
  requestParams.value = '[]'
  callResult.value = null
  callError.value = null
}
async function callApi() {
  if (!selectedApi.value || calling.value) return
  calling.value = true
  callError.value = null
  callResult.value = null
  try {
    let paramsObj: Record<string, any> = {}
    try {
      const parsed = JSON.parse(requestParams.value)
      if (Array.isArray(parsed)) {
        selectedApi.value.params.forEach((paramName, index) => {
          if (index < parsed.length) {
            paramsObj[paramName] = parsed[index]
          }
        })
      } else if (typeof parsed === 'object') {
        paramsObj = parsed
      }
    } catch (e: any) {
      throw new Error('参数格式错误: ' + e.message)
    }
    const result = await invoke(selectedApi.value.command, paramsObj)
    callResult.value = result
  } catch (e: any) {
    callError.value = e.message || String(e)
  } finally {
    calling.value = false
  }
}
function clearResult() {
  callResult.value = null
  callError.value = null
}
</script>
<style scoped>
.debug-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-main);
  overflow: hidden;
}
.debug-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}
.debug-header h2 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}
.search-input {
  flex: 1;
  max-width: 300px;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 20px;
  color: var(--text-primary);
  font-size: 13px;
}
.search-input:focus {
  outline: none;
  border-color: var(--primary);
}
.debug-content {
  flex: 1;
  display: flex;
  overflow: hidden;
}
.api-list {
  width: 280px;
  min-width: 280px;
  border-right: 1px solid var(--border);
  overflow-y: auto;
  padding: 8px;
}
.api-item {
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 4px;
  transition: background 0.15s;
}
.api-item:hover {
  background: var(--bg-tertiary);
}
.api-item.active {
  background: var(--primary);
  color: white;
}
.api-name {
  display: block;
  font-size: 13px;
  font-weight: 500;
}
.api-file {
  display: block;
  font-size: 11px;
  opacity: 0.7;
  margin-top: 2px;
}
.api-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  padding: 16px;
}
.no-selection {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 14px;
}
.api-detail-header {
  margin-bottom: 16px;
}
.api-detail-header h3 {
  margin: 0 0 8px;
  font-size: 18px;
}
.api-signature {
  font-family: monospace;
  font-size: 12px;
  color: var(--text-secondary);
  background: var(--bg-tertiary);
  padding: 4px 8px;
  border-radius: 4px;
}
.api-params {
  margin-bottom: 16px;
}
.param-label,
.result-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}
.param-input {
  width: 100%;
  height: 100px;
  padding: 10px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-primary);
  font-family: monospace;
  font-size: 13px;
  resize: none;
}
.param-input:focus {
  outline: none;
  border-color: var(--primary);
}
.api-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
.btn-call {
  padding: 8px 20px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
}
.btn-call:hover:not(:disabled) {
  opacity: 0.9;
}
.btn-call:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn-clear {
  padding: 8px 16px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border: 1px solid var(--border);
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
}
.btn-clear:hover {
  background: var(--bg-secondary);
}
.api-result {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
}
.result-content {
  flex: 1;
  margin: 0;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  font-family: monospace;
  font-size: 12px;
  line-height: 1.5;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-word;
}
.result-error {
  padding: 12px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid #ef4444;
  border-radius: 8px;
  color: #ef4444;
  font-size: 13px;
}
</style>
