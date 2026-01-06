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
interface ApiItem {
  name: string
  file: string
  signature: string
  params: string[]
}
const API_LIST: ApiItem[] = [
  { name: 'IsDebugMode', file: 'app.go', signature: '() bool', params: [] },
  { name: 'LogDebug', file: 'app.go', signature: '(tag string, message string)', params: ['tag', 'message'] },
  { name: 'LoadConfig', file: 'config.go', signature: '() (*Config, error)', params: [] },
  { name: 'SaveConfig', file: 'config.go', signature: '(config Config) error', params: ['config'] },
  { name: 'UpdateConfig', file: 'config.go', signature: '(key string, value string) error', params: ['key', 'value'] },
  { name: 'SaveSessionState', file: 'config.go', signature: '(state SessionState) error', params: ['state'] },
  { name: 'LoadSessionState', file: 'config.go', signature: '() (*SessionState, error)', params: [] },
  { name: 'GetProjects', file: 'project.go', signature: '() ([]ProjectInfo, error)', params: [] },
  { name: 'LoadProjectByName', file: 'project.go', signature: '(name string) (*Project, error)', params: ['name'] },
  { name: 'LoadProject', file: 'project.go', signature: '(projectID int) (*Project, error)', params: ['projectID'] },
  { name: 'CreateProject', file: 'project.go', signature: '(name string) (int, error)', params: ['name'] },
  { name: 'DeleteProject', file: 'project.go', signature: '(projectName string) error', params: ['projectName'] },
  { name: 'RenameProject', file: 'project.go', signature: '(oldName string, newName string) error', params: ['oldName', 'newName'] },
  { name: 'UpdateProjectTime', file: 'project.go', signature: '(projectID int) error', params: ['projectID'] },
  { name: 'CreateChildNode', file: 'node.go', signature: '(projectName string, parentNodeID any, name string) (*NodeInfo, error)', params: ['projectName', 'parentNodeID', 'name'] },
  { name: 'CreateBrotherNode', file: 'node.go', signature: '(projectName string, siblingNodeID any, name string) (*NodeInfo, error)', params: ['projectName', 'siblingNodeID', 'name'] },
  { name: 'UpdateNode', file: 'node.go', signature: '(projectName string, nodeData map[string]any) error', params: ['projectName', 'nodeData'] },
  { name: 'DeleteNode', file: 'node.go', signature: '(projectName string, nodeID any) error', params: ['projectName', 'nodeID'] },
  { name: 'RenameNode', file: 'node.go', signature: '(nodeID int, newName string) error', params: ['nodeID', 'newName'] },
  { name: 'RebaseNode', file: 'node.go', signature: '(nodeID int, newParentID *int) error', params: ['nodeID', 'newParentID'] },
  { name: 'GetNodePath', file: 'node.go', signature: '(nodeID int) ([]int, error)', params: ['nodeID'] },
  { name: 'GetApps', file: 'database.go', signature: '() ([]WTApp, error)', params: [] },
  { name: 'CreateApp', file: 'database.go', signature: '(name string) (int, error)', params: ['name'] },
  { name: 'DeleteApp', file: 'database.go', signature: '(appID int) error', params: ['appID'] },
  { name: 'RenameApp', file: 'database.go', signature: '(appID int, newName string) error', params: ['appID', 'newName'] },
  { name: 'GetLocalConversations', file: 'database.go', signature: '(appID int) ([]WTConversation, error)', params: ['appID'] },
  { name: 'CreateConversation', file: 'database.go', signature: '(appID int, name string) (int, error)', params: ['appID', 'name'] },
  { name: 'GetDialogues', file: 'database.go', signature: '(conversationID int, page, limit int) ([]WTDialogue, int, error)', params: ['conversationID', 'page', 'limit'] },
  { name: 'CreateDialogue', file: 'database.go', signature: '(conversationID int, requestContent, responseContent string) (int, error)', params: ['conversationID', 'requestContent', 'responseContent'] },
  { name: 'GetGalleryImages', file: 'gallery.go', signature: '() ([]GalleryImage, error)', params: [] },
  { name: 'GetGalleryFolders', file: 'gallery.go', signature: '() ([]GalleryFolder, error)', params: [] },
  { name: 'DeleteGalleryImage', file: 'gallery.go', signature: '(id string) error', params: ['id'] },
  { name: 'GetAllDrafts', file: 'drafts.go', signature: '() ([]Draft, error)', params: [] },
  { name: 'CreateDraft', file: 'drafts.go', signature: '(draft Draft) error', params: ['draft'] },
  { name: 'DeleteDraft', file: 'drafts.go', signature: '(id string) error', params: ['id'] },
  { name: 'GetClipboardCaptures', file: 'drafts.go', signature: '() ([]ClipboardCapture, error)', params: [] },
  { name: 'CopyToClipboard', file: 'clipboard.go', signature: '(content string) error', params: ['content'] },
  { name: 'GetSafeModeTxtFiles', file: 'safemode.go', signature: '() ([]string, error)', params: [] },
  { name: 'IsCapsLockOn', file: 'safemode.go', signature: '() bool', params: [] },
  { name: 'GetTempDir', file: 'temp_manager.go', signature: '() (string, error)', params: [] },
  { name: 'ListTempFiles', file: 'temp_manager.go', signature: '() ([]TempFileInfo, error)', params: [] },
  { name: 'GetUnsavedSessions', file: 'temp_manager.go', signature: '() ([]SessionRecoveryInfo, error)', params: [] },
  { name: 'GetLocalCreations', file: 'api.go', signature: '() ([]map[string]any, error)', params: [] },
  { name: 'GetCreationsDir', file: 'api.go', signature: '() (string, error)', params: [] },
  { name: 'SendTestChat', file: 'llm_api.go', signature: '(request TestChatRequest) TestChatResponse', params: ['request'] },
  { name: 'GetNodeContent', file: 'worldtree.go', signature: '(nodeID int) (*NodeContent, error)', params: ['nodeID'] },
  { name: 'AddFolder', file: 'worldtree.go', signature: '(nodeID int, name string) (int, error)', params: ['nodeID', 'name'] },
  { name: 'AddCard', file: 'worldtree.go', signature: '(nodeID int, folderChangeID *int, name string, keyWord string) (int, error)', params: ['nodeID', 'folderChangeID', 'name', 'keyWord'] },
  { name: 'AuthLogin', file: 'auth.go', signature: '(email, password string, rememberMe bool) (map[string]any, error)', params: ['email', 'password', 'rememberMe'] },
  { name: 'AuthGetProfile', file: 'auth.go', signature: '(token string) (map[string]any, error)', params: ['token'] },
  { name: 'AuthLogout', file: 'auth.go', signature: '(token string) error', params: ['token'] },
  { name: 'FetchWithAuth', file: 'api.go', signature: '(token, url, method, body string) (string, error)', params: ['token', 'url', 'method', 'body'] },
  { name: 'GetDirtyState', file: 'dirty_manager.go', signature: '(id string) *DirtyState', params: ['id'] },
  { name: 'HasAnyDirtyContent', file: 'dirty_manager.go', signature: '() bool', params: [] },
  { name: 'GetAllDirtyStates', file: 'dirty_manager.go', signature: '() []DirtyState', params: [] }
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
    let params: any[] = []
    try {
      params = JSON.parse(requestParams.value)
      if (!Array.isArray(params)) params = [params]
    } catch (e: any) {
      throw new Error('参数格式错误: ' + e.message)
    }
    const fn = (window as any).go?.main?.App?.[selectedApi.value.name]
    if (!fn) {
      throw new Error('接口不存在: ' + selectedApi.value.name)
    }
    const result = await fn(...params)
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
