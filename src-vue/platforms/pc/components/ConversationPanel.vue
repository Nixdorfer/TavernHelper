<template>
  <aside :class="['panel', 'conversation-panel', { collapsed: collapsed, embedded: embedded }]" :style="embedded ? {} : panelStyle">
    <div v-if="!embedded" class="panel-header">
      <div class="header-left" @click="emit('toggle')">
        <span :class="['collapse-icon', 'horizontal', { expanded: !collapsed }]">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
        </span>
        <h2>对话记录</h2>
      </div>
      <button v-show="!collapsed" class="btn-refresh" @click.stop="refreshAll" :disabled="refreshing" title="刷新应用和对话列表">
        <span :class="{ spinning: refreshing }">🔄</span>
      </button>
    </div>
    <div v-show="!collapsed" class="panel-content">
      <div v-if="loadingApps" class="loading-hint">加载应用列表...</div>
      <div v-else-if="installedApps.length > 0" class="apps-list">
        <div v-for="app in installedApps" :key="app.id" class="app-item">
          <div class="app-header" @click="toggleApp(app)">
            <span :class="['collapse-icon', { expanded: app.expanded }]">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
            </span>
            <span class="app-name">{{ app.app?.name || app.name }}</span>
            <button class="btn-new-conv" @click.stop="createNewConversation(app)" title="新建对话">+</button>
          </div>
          <div v-if="app.expanded" class="conversations-list">
            <div v-if="app.loadingConversations" class="loading-hint">加载对话...</div>
            <template v-else>
              <div
                v-if="pendingNewConv && pendingNewConv.appId === (app.app?.id || app.id)"
                class="conversation-item pending-conv"
              >
                <span class="conversation-name">正在发起新对话...</span>
              </div>
              <div
                v-for="conv in app.conversations"
                :key="conv.id"
                :class="['conversation-item', { active: selectedConversation?.id === conv.id }]"
                @click="selectConversation(app, conv)"
              >
                <span class="conversation-name">{{ conv.name }}</span>
              </div>
              <div v-if="app.conversations.length === 0 && !pendingNewConv" class="empty-hint">暂无对话</div>
            </template>
          </div>
        </div>
      </div>
      <div v-else class="empty-hint">暂无已安装应用</div>
    </div>
    <div v-if="confirmingDelete" class="confirm-overlay" @click.self="cancelDelete">
      <div class="confirm-dialog">
        <div class="confirm-title">删除对话</div>
        <div class="confirm-message">确定要删除对话 "{{ confirmingDelete.conv.name }}" 吗？此操作不可撤销。</div>
        <div class="confirm-actions">
          <button class="btn btn-secondary" @click="cancelDelete">取消</button>
          <button class="btn btn-danger" @click="executeDelete">删除</button>
        </div>
      </div>
    </div>
  </aside>
</template>
<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { conversationApi } from '@/api/modules/conversation'
import { api } from '@/api'
import { logger } from '@/utils/logger'
interface AppItem {
  id: string
  app?: {
    id: string
    name: string
    icon?: string
    description?: string
  }
  name?: string
  expanded: boolean
  loadingConversations: boolean
  conversations: ConvItem[]
}
interface ConvItem {
  id: string
  name: string
  expanded?: boolean
}
const props = withDefaults(defineProps<{
  collapsed?: boolean
  width?: number
  token?: string
  embedded?: boolean
  userId?: string
  restoreConversation?: any
  debugMode?: boolean
  pendingNewConv?: any
}>(), {
  collapsed: false,
  width: 350,
  token: '',
  embedded: false,
  userId: '',
  debugMode: false
})
const emit = defineEmits<{
  toggle: []
  'select-conversation': [data: { app: AppItem; conversation: ConvItem }]
  'show-notification': [data: { message: string; type: string }]
  'conversation-restored': []
  'new-conversation': [data: { app: AppItem; appId: string }]
}>()
const loadingApps = ref(false)
const installedApps = ref<AppItem[]>([])
const selectedConversation = ref<ConvItem | null>(null)
const selectedApp = ref<AppItem | null>(null)
const refreshing = ref(false)
const confirmingDelete = ref<{ app: AppItem; conv: ConvItem } | null>(null)
const hasRestoredConversation = ref(false)
const panelStyle = computed(() => ({
  width: props.collapsed ? '40px' : props.width + 'px',
  minWidth: props.collapsed ? '40px' : props.width + 'px'
}))
watch(() => props.token, (newToken) => {
  if (newToken) {
    loadInstalledApps()
  } else {
    installedApps.value = []
    selectedConversation.value = null
    if (props.debugMode) {
      addTestApp()
    }
  }
}, { immediate: true })
watch(() => props.debugMode, (newVal) => {
  if (newVal) {
    addTestApp()
  }
}, { immediate: true })
watch(() => props.restoreConversation, () => {
  if (props.restoreConversation && !hasRestoredConversation.value && installedApps.value.length > 0) {
    tryRestoreConversation()
  }
}, { immediate: true })
watch(installedApps, () => {
  if (props.restoreConversation && !hasRestoredConversation.value && installedApps.value.length > 0) {
    tryRestoreConversation()
  }
})
async function loadInstalledApps() {
  if (props.debugMode) {
    addTestApp()
  }
  if (!props.token) {
    loadingApps.value = false
    return
  }
  loadingApps.value = true
  try {
    await fetchAppsFromAPI()
  } catch (e) {
    logger.error('加载应用失败:', e)
  }
  loadingApps.value = false
}
async function fetchAppsFromAPI() {
  try {
    const response = await fetch('https://aipornhub.ltd/console/api/used-installed-apps?page=1&limit=100', {
      headers: {
        'accept': '*/*',
        'authorization': `Bearer ${props.token}`,
        'content-type': 'application/json',
        'x-language': 'zh-Hans'
      }
    })
    const data = await response.json()
    if (data.installed_apps) {
      const expandedIds = new Set(installedApps.value.filter(a => a.expanded).map(a => a.app?.id || a.id))
      installedApps.value = data.installed_apps.map((app: any) => ({
        ...app,
        expanded: expandedIds.has(app.app?.id),
        loadingConversations: false,
        conversations: []
      }))
      if (props.debugMode) {
        addTestApp()
      }
    }
  } catch (e) {
    logger.error('加载应用失败:', e)
  }
}
async function toggleApp(app: AppItem) {
  app.expanded = !app.expanded
  if (app.expanded && app.conversations.length === 0) {
    await loadConversations(app)
  }
}
async function loadConversations(app: AppItem) {
  app.loadingConversations = true
  const appId = app.app?.id || app.id
  if (appId === '__debug_test_app__') {
    app.loadingConversations = false
    return
  }
  try {
    const response = await fetch(`https://aipornhub.ltd/console/api/installed-apps/${appId}/conversations?limit=500&pinned=false`, {
      headers: {
        'accept': '*/*',
        'authorization': `Bearer ${props.token}`,
        'content-type': 'application/json',
        'x-language': 'zh-Hans'
      }
    })
    const data = await response.json()
    if (data.data) {
      const expandedId = selectedConversation.value?.id
      app.conversations = data.data.map((conv: any) => ({
        ...conv,
        expanded: conv.id === expandedId
      }))
    }
  } catch (e) {
    logger.error('加载对话失败:', e)
  }
  app.loadingConversations = false
}
function selectConversation(app: AppItem, conv: ConvItem) {
  selectedConversation.value = conv
  selectedApp.value = app
  emit('select-conversation', {
    app: app,
    conversation: conv
  })
}
function cancelDelete() {
  confirmingDelete.value = null
}
async function executeDelete() {
  if (!confirmingDelete.value) return
  const { app, conv } = confirmingDelete.value
  confirmingDelete.value = null
  const appId = app.app?.id || app.id
  try {
    await api.conversation.delete(props.token, appId, conv.id)
    const idx = app.conversations.findIndex(c => c.id === conv.id)
    if (idx > -1) {
      app.conversations.splice(idx, 1)
    }
    if (selectedConversation.value?.id === conv.id) {
      selectedConversation.value = null
      selectedApp.value = null
      emit('select-conversation', { app: null as any, conversation: null as any })
    }
  } catch (e) {
    logger.error('删除对话失败:', e)
    emit('show-notification', { message: '删除对话失败: ' + e, type: 'error' })
  }
}
async function refreshAll() {
  if (refreshing.value) return
  refreshing.value = true
  try {
    await fetchAppsFromAPI()
    const expandedApps = installedApps.value.filter(app => app.expanded)
    await Promise.all(expandedApps.map(app => loadConversations(app)))
  } catch (e) {
    logger.error('刷新失败:', e)
  } finally {
    refreshing.value = false
  }
}
async function tryRestoreConversation() {
  if (!props.restoreConversation || hasRestoredConversation.value) return
  const { appId, conversationId } = props.restoreConversation
  if (!appId || !conversationId) return
  const app = installedApps.value.find(a =>
    (a.app?.id === appId) || (a.id === appId)
  )
  if (!app) {
    return
  }
  if (!app.expanded) {
    app.expanded = true
    await loadConversations(app)
  }
  const conv = app.conversations?.find(c => c.id === conversationId)
  if (conv) {
    hasRestoredConversation.value = true
    selectConversation(app, conv)
    emit('conversation-restored')
  }
}
function addTestApp() {
  const existingTestApp = installedApps.value.find(a => a.id === '__debug_test_app__')
  if (existingTestApp) return
  const testApp: AppItem = {
    id: '__debug_test_app__',
    app: {
      id: '__debug_test_app__',
      name: '测试应用'
    },
    expanded: false,
    loadingConversations: false,
    conversations: [
      { id: '__debug_test_conversation__', name: '测试对话' }
    ]
  }
  installedApps.value.unshift(testApp)
}
function createNewConversation(app: AppItem) {
  const appId = app.app?.id || app.id
  app.expanded = true
  emit('new-conversation', { app, appId })
}
function addNewConversation(appId: string, conversation: ConvItem) {
  const app = installedApps.value.find(a =>
    (a.app?.id === appId) || (a.id === appId)
  )
  if (app) {
    if (!app.conversations) {
      app.conversations = []
    }
    const exists = app.conversations.find(c => c.id === conversation.id)
    if (!exists) {
      app.conversations.unshift(conversation)
    }
  }
}
function addAppToList(appData: any) {
  const appId = appData.app?.id || appData.id
  const existing = installedApps.value.find(a => (a.app?.id || a.id) === appId)
  if (!existing) {
    const newApp: AppItem = {
      id: appId,
      app: appData.app || { id: appId, name: appData.name, icon: appData.icon },
      name: appData.name,
      conversations: [],
      expanded: true,
      loadingConversations: false
    }
    installedApps.value.unshift(newApp)
    return newApp
  }
  existing.expanded = true
  return existing
}
function removeAppIfEmpty(appId: string) {
  const appIndex = installedApps.value.findIndex(a => (a.app?.id || a.id) === appId)
  if (appIndex >= 0) {
    const app = installedApps.value[appIndex]
    if (!app.conversations || app.conversations.length === 0) {
      installedApps.value.splice(appIndex, 1)
    }
  }
}
defineExpose({
  addNewConversation,
  addAppToList,
  removeAppIfEmpty,
  refreshAll
})
</script>
<style scoped>
.conversation-panel {
  border-left: 1px solid var(--border-color);
  border-right: none;
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  overflow: hidden;
}
.conversation-panel.embedded {
  border-left: none;
  border-right: none;
  width: 100%;
  height: 100%;
}
.conversation-panel.embedded .panel-content {
  flex: 1;
  max-height: none;
}
.panel-header {
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}
.panel-header h2 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}
.btn-refresh {
  padding: 2px 8px;
  border-radius: 4px;
  border: none;
  background: var(--primary-color);
  color: white;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-refresh:hover:not(:disabled) {
  background: var(--primary-hover);
}
.btn-refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
.btn-refresh .spinning {
  display: inline-block;
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.collapse-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.collapse-icon svg {
  width: 12px;
  height: 12px;
  color: var(--text-secondary);
}
.collapse-icon.expanded {
  transform: rotate(90deg);
}
.collapse-icon.horizontal {
  transform: rotate(180deg);
}
.collapse-icon.horizontal.expanded {
  transform: rotate(0deg);
}
.panel-content {
  flex: 0 0 auto;
  max-height: 40%;
  overflow-y: auto;
  padding: 8px;
}
.loading-hint,
.empty-hint {
  padding: 12px;
  color: var(--text-secondary);
  font-size: 12px;
  text-align: center;
}
.apps-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.app-item {
  background: var(--bg-tertiary);
  border-radius: 6px;
  overflow: hidden;
}
.app-header {
  padding: 8px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: background 0.2s;
}
.app-header:hover {
  background: var(--bg-hover);
}
.app-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}
.btn-new-conv {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: none;
  background: var(--primary-color);
  color: white;
  font-size: 14px;
  font-weight: bold;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  opacity: 0;
  transition: all 0.2s;
}
.app-header:hover .btn-new-conv {
  opacity: 1;
}
.btn-new-conv:hover {
  transform: scale(1.1);
}
.conversations-list {
  padding: 4px 8px 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.conversation-item {
  padding: 6px 10px;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  transition: background 0.2s;
  font-size: 12px;
  color: var(--text-secondary);
}
.conversation-item:hover {
  background: var(--bg-hover);
}
.conversation-item.active {
  background: var(--primary-color);
  color: white;
}
.conversation-item.pending-conv {
  background: #fef3c7;
  color: #92400e;
  animation: pendingPulse 1.5s ease-in-out infinite;
}
@keyframes pendingPulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}
.conversation-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.confirm-overlay {
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
.confirm-dialog {
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 20px;
  min-width: 300px;
  max-width: 400px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}
.confirm-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 12px;
}
.confirm-message {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 20px;
  line-height: 1.5;
}
.confirm-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
.confirm-actions .btn {
  padding: 6px 16px;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  font-size: 14px;
}
.confirm-actions .btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.confirm-actions .btn-danger {
  background: #ef4444;
  color: white;
}
.confirm-actions .btn:hover {
  opacity: 0.9;
}
</style>
