<template>
  <div class="conversation-view">
    <aside :class="['conversation-sidebar', { collapsed: conversationPanelCollapsed }]" :style="sidebarStyle">
      <div class="panel-header conversation-header" @click="togglePanel">
        <span :class="['collapse-icon', { expanded: !sidebarShowWorldTree }]">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
        </span>
        <h2>对话记录</h2>
      </div>
      <div v-show="!conversationPanelCollapsed && !sidebarShowWorldTree" class="sidebar-main-content">
        <ConversationPanel
          ref="conversationPanelRef"
          :collapsed="false"
          :width="conversationPanelWidth"
          :token="token"
          :user-id="userId"
          :embedded="true"
          :restore-conversation="pendingConversationRestore"
          :debug-mode="isDebugMode"
          :pending-new-conv="pendingNewConv"
          @toggle="handlePanelToggle"
          @select-conversation="handleSelectConversation"
          @show-notification="handleShowNotification"
          @conversation-restored="handleConversationRestored"
          @new-conversation="handleNewConversation"
        />
      </div>
      <div v-show="!conversationPanelCollapsed" :class="['worldtree-section', { expanded: sidebarShowWorldTree }]">
        <div class="worldtree-header" @click="toggleWorldTree">
          <span :class="['collapse-icon', { expanded: sidebarShowWorldTree }]">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
          </span>
          <label :class="['worldtree-switch', { disabled: !selectedConversation && !newConversationStartPage }]" @click.stop>
            <input type="checkbox" :checked="worldTreeEnabled" @change="handleWorldTreeToggle" :disabled="!selectedConversation && !newConversationStartPage">
            <span class="switch-slider"></span>
          </label>
          <h2>世界树</h2>
          <button v-if="hasWorldTree" class="btn-expand-tree" @click.stop="openFullWorldtree" title="全屏查看">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M15 3h6v6M9 21H3v-6M21 3l-7 7M3 21l7-7"/>
            </svg>
          </button>
          <button v-if="hasWorldTree" class="btn-unbind" @click.stop="unbindWorldtree" title="解绑世界树">×</button>
        </div>
        <div v-show="sidebarShowWorldTree" :class="['world-tree-content', { disabled: !worldTreeEnabled }]">
          <template v-if="conversationWorldTree?.timeline?.length > 0">
            <div class="readonly-timeline-graph">
              <TimelineGraph
                :nodes="conversationWorldTree.timeline"
                :current-node-id="conversationWorldTree.currentNode"
                :branch-names="{}"
                :readonly="true"
                :path-mode="true"
                :required-branch-id="conversationCurrentBranchId"
                :allow-parallel-world="true"
                @create-parallel-world="handleCreateParallelWorld"
              />
            </div>
          </template>
          <template v-else-if="newConvWorldTree">
            <div class="readonly-timeline-graph new-conv-graph">
              <TimelineGraph
                :nodes="newConvWorldTree.timeline"
                :current-node-id="newConvWorldTree.currentNode"
                :branch-names="{}"
                :readonly="!newConvSelectingNode"
                :full-mode="newConvSelectingNode"
                @node-click="handleNewConvNodeClick"
              />
            </div>
            <div v-if="newConvSelectingNode" class="node-select-actions">
              <button class="btn btn-sm btn-primary" @click="confirmNodeSelection" :disabled="!newConvSelectedNodeId">确定</button>
              <button class="btn btn-sm btn-secondary" @click="cancelNodeSelection">取消</button>
            </div>
          </template>
          <div v-else-if="selectedConversation || newConversationStartPage" class="bind-worldtree-area">
            <div class="bind-worldtree-hint">未绑定世界树</div>
            <div class="bind-worldtree-actions">
              <button v-if="selectedConversation" class="btn btn-sm btn-primary" @click="initWorldtreeFromRemote">从远端同步</button>
              <button class="btn btn-sm btn-secondary" @click="showProjectSelect">从项目加载</button>
            </div>
          </div>
          <div v-else class="outline-empty">
            <span>请选择对话</span>
          </div>
        </div>
      </div>
    </aside>
    <div class="conversation-content">
      <div v-if="showParallelWorldConv && parallelWorldNewConv" class="parallel-world-conv-area">
        <div class="new-conv-header">
          <h3>{{ parallelWorldNewConv.nodeName }}</h3>
          <button class="btn-header" @click="closeParallelWorld">关闭</button>
        </div>
        <div class="empty-conv-hint">
          <span>发送第一条消息开始对话</span>
        </div>
        <div class="new-conv-input-area">
          <textarea
            v-model="parallelWorldConvInput"
            class="message-input"
            placeholder="输入第一条消息..."
            rows="1"
            @keydown.enter.exact="handleParallelWorldEnterKey"
            @keydown.enter.shift.exact.stop
            :disabled="sendingParallelWorldMsg"
          ></textarea>
          <button
            class="send-btn"
            @click="sendParallelWorldMessage"
            :disabled="!parallelWorldConvInput.trim() || sendingParallelWorldMsg"
          >
            <span v-if="sendingParallelWorldMsg" class="loading-spinner-sm"></span>
            <span v-else>发送</span>
          </button>
        </div>
      </div>
      <div v-else-if="newConversationStartPage" class="new-conversation-area">
        <div class="new-conv-header">
          <h3>{{ newConversationStartPage.appName }}</h3>
          <button class="btn-header" @click="closeNewConversation">关闭</button>
        </div>
        <div class="start-page-container">
          <iframe
            ref="startPageIframe"
            :srcdoc="newConversationStartPage.html"
            class="start-page-iframe"
            sandbox="allow-scripts allow-same-origin"
          ></iframe>
        </div>
        <div class="message-input-area">
          <textarea
            v-model="newConvMessageInput"
            class="message-input"
            placeholder="输入第一条消息开始对话..."
            rows="1"
            @keydown.enter.exact="handleNewConvEnterKey"
            @keydown.enter.shift.exact.stop
            :disabled="sendingNewConvMessage"
          ></textarea>
          <div class="input-buttons">
            <button
              :class="['send-btn', { 'send-btn-streaming': streamingEnabled }]"
              @click="sendNewConversationMessage"
              :disabled="!newConvMessageInput.trim() || sendingNewConvMessage"
            >
              <span v-if="sendingNewConvMessage" class="loading-spinner-sm"></span>
              <span v-else>发送</span>
            </button>
          </div>
        </div>
      </div>
      <div v-else-if="!selectedConversation" class="conversation-placeholder">
        <span>请选择一个对话</span>
      </div>
      <div v-else class="messages-area">
        <div class="messages-header">
          <h3>{{ selectedConversation.name }}</h3>
          <span v-if="loadingMessages" class="loading-spinner"></span>
          <div class="header-actions">
            <button v-if="isDebugTestConversation" class="btn-header" @click="showDebugReplyModal">设定回复</button>
            <button class="btn-header" @click="closeConversation">关闭</button>
            <button v-if="!isDebugTestConversation" class="btn-header btn-danger" @click="showDeleteConfirm">删除</button>
          </div>
        </div>
        <div class="messages-container" ref="messagesContainer" @scroll="onMessagesScroll">
          <div v-if="hasMoreMessages && !loadingMessages" class="load-more">
            <button class="btn btn-secondary btn-sm" @click="loadMoreMessages">加载更早的消息</button>
          </div>
          <div v-for="(msg, msgIndex) in filteredConversationMessages" :key="msg.id" :class="['message-item', { editing: editingMessageId === msg.id }]">
            <div class="message-query">
              <div class="message-role">用户</div>
              <div v-if="editingMessageId === msg.id" class="message-edit-area">
                <textarea
                  v-model="editingMessageContent"
                  class="message-edit-textarea"
                  ref="editTextarea"
                ></textarea>
                <div class="message-edit-actions">
                  <button class="btn-edit-save" @click="saveEditedMessage(msg)" :disabled="sendingEditMessage">
                    {{ sendingEditMessage ? '发送中...' : '保存并重新生成' }}
                  </button>
                  <button class="btn-edit-cancel" @click="cancelEditMessage">取消</button>
                </div>
              </div>
              <div v-else-if="getCleanUserQuery(msg.query)" class="message-content">{{ getCleanUserQuery(msg.query) }}</div>
            </div>
            <div v-if="editingMessageId !== msg.id && lastMessageHasError && msgIndex === filteredConversationMessages.length - 1" class="message-error-retry">
              <button class="btn-error-retry" @click="regenerateMessage(msg)" :disabled="regeneratingMessageId === msg.id">
                {{ regeneratingMessageId === msg.id ? '重试中...' : '生成错误 点击重试' }}
              </button>
            </div>
            <div v-else-if="editingMessageId !== msg.id" class="message-answer" @contextmenu="handleAnswerContextMenu($event, msg)">
              <div class="message-role">
                <span>AI</span>
                <span class="model-badge">{{ msg.model_id }}</span>
              </div>
              <div class="message-content selectable-text markdown-body" v-html="renderMarkdown(getMessageContent(msg.answer))"></div>
            </div>
            <div class="message-footer">
              <div class="footer-info">
                <div class="footer-line">{{ formatMessageTime(msg.created_at) }}</div>
                <div class="footer-line">发送: {{ msg.message_points }} 积分 / {{ msg.message_tokens }} tokens</div>
                <div class="footer-line">回复: {{ msg.answer_points }} 积分 / {{ msg.answer_tokens }} tokens</div>
              </div>
              <div v-if="editingMessageId !== msg.id" class="footer-actions">
                <button class="btn-msg-action" @click="regenerateMessage(msg)" :disabled="regeneratingMessageId === msg.id">
                  {{ regeneratingMessageId === msg.id ? '生成中...' : '重新生成' }}
                </button>
                <button class="btn-msg-action" @click="startEditMessage(msg)">编辑</button>
                <button class="btn-msg-action btn-danger" @click="deleteMessage(msg)">删除</button>
              </div>
            </div>
          </div>
          <div v-if="filteredConversationMessages.length === 0 && !loadingMessages" class="empty-hint">暂无消息</div>
        </div>
        <div class="message-input-area">
          <textarea
            ref="messageInputRef"
            v-model="messageInput"
            class="message-input"
            placeholder="输入消息..."
            rows="1"
            @keydown.enter.exact="handleEnterKey"
            @keydown.enter.shift.exact.stop
            :disabled="sendingMessage"
          ></textarea>
          <div class="input-buttons">
            <button :class="['scroll-bottom-btn', { active: canScrollDown, locked: autoScrollLocked }]" @click="toggleAutoScroll" :title="autoScrollLocked ? '自动滚动已开启' : '滚动到底部'">
              <svg v-if="autoScrollLocked" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
              </svg>
              <span v-else>↓</span>
            </button>
            <button
              :class="['send-btn', 'send-btn-' + sendButtonState, { 'send-btn-streaming': streamingEnabled && sendButtonState === 'ready' }]"
              @click="sendButtonState === 'streaming' ? stopStreaming() : sendMessage()"
              :disabled="sendButtonState !== 'ready' && sendButtonState !== 'streaming' ? true : (sendButtonState === 'ready' && !messageInput.trim())"
            >
              <span v-if="sendButtonState === 'sending'" class="loading-spinner-sm"></span>
              <span v-else-if="sendButtonState === 'streaming'">终止</span>
              <span v-else>发送</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useAuthStore, useConfigStore, useNotificationStore, useConfirmStore } from '@/stores'
import { useConversationStore, useMessagesStore, useWorldTreeStore, useStreamingStore, useRevisionsStore } from '@/stores/modules/conversation'
import { conversationApi } from '@/api/modules/conversation'
import ConversationPanel from '@/platforms/pc/components/ConversationPanel.vue'
import TimelineGraph from '@/components/business/TimelineGraph.vue'
const props = defineProps<{
  isDebugMode?: boolean
}>()
const emit = defineEmits<{
  'toggle-panel': []
  'switch-to-worldbook': []
  'show-notification': [event: any]
}>()
const authStore = useAuthStore()
const configStore = useConfigStore()
const notificationStore = useNotificationStore()
const confirmStore = useConfirmStore()
const conversationStore = useConversationStore()
const messagesStore = useMessagesStore()
const worldTreeStore = useWorldTreeStore()
const streamingStore = useStreamingStore()
const revisionsStore = useRevisionsStore()
const conversationPanelRef = ref<any>(null)
const messagesContainer = ref<HTMLElement | null>(null)
const messageInputRef = ref<HTMLTextAreaElement | null>(null)
const startPageIframe = ref<HTMLIFrameElement | null>(null)
const sidebarShowWorldTree = ref(false)
const token = computed(() => authStore.token)
const userId = computed(() => authStore.userId)
const conversationPanelWidth = computed(() => conversationStore.conversationPanelWidth)
const conversationPanelCollapsed = computed(() => conversationStore.conversationPanelCollapsed)
const pendingConversationRestore = computed(() => conversationStore.pendingConversationRestore)
const newConversationStartPage = computed(() => conversationStore.newConversationStartPage)
const newConvWorldTree = computed(() => conversationStore.newConvWorldTree)
const newConvSelectingNode = computed(() => conversationStore.newConvSelectingNode)
const newConvSelectedNodeId = computed(() => conversationStore.newConvSelectedNodeId)
const newConvMessageInput = computed({
  get: () => conversationStore.newConvMessageInput,
  set: (v) => { conversationStore.newConvMessageInput = v }
})
const sendingNewConvMessage = computed(() => conversationStore.sendingNewConvMessage)
const pendingNewConv = computed(() => conversationStore.pendingNewConv)
const showParallelWorldConv = computed(() => conversationStore.showParallelWorldConv)
const parallelWorldNewConv = computed(() => conversationStore.parallelWorldNewConv)
const parallelWorldConvInput = computed({
  get: () => conversationStore.parallelWorldConvInput,
  set: (v) => { conversationStore.parallelWorldConvInput = v }
})
const sendingParallelWorldMsg = computed(() => conversationStore.sendingParallelWorldMsg)
const selectedConversation = computed(() => messagesStore.selectedConversation)
const conversationMessages = computed(() => messagesStore.messages)
const loadingMessages = computed(() => messagesStore.isLoading)
const hasMoreMessages = computed(() => messagesStore.hasMore)
const messageInput = computed({
  get: () => messagesStore.input,
  set: (v) => messagesStore.setInput(v)
})
const sendingMessage = computed(() => messagesStore.isSending)
const sendButtonState = computed(() => messagesStore.sendButtonState)
const canScrollDown = computed(() => messagesStore.canScrollDown)
const autoScrollLocked = computed(() => messagesStore.autoScrollLocked)
const editingMessageId = computed(() => messagesStore.editingMessageId)
const editingMessageContent = computed({
  get: () => messagesStore.editingMessageContent,
  set: (v) => { messagesStore.editingMessageContent = v }
})
const sendingEditMessage = computed(() => messagesStore.isSendingEdit)
const regeneratingMessageId = computed(() => messagesStore.regeneratingMessageId)
const conversationCurrentBranchId = computed(() => messagesStore.currentBranchId)
const worldTreeEnabled = computed(() => worldTreeStore.isEnabled)
const conversationWorldTree = computed(() => worldTreeStore.conversationTree)
const streamingEnabled = computed(() => streamingStore.isEnabled)
const sidebarStyle = computed(() => ({
  width: conversationPanelCollapsed.value ? '40px' : conversationPanelWidth.value + 'px',
  minWidth: conversationPanelCollapsed.value ? '40px' : conversationPanelWidth.value + 'px'
}))
const isDebugTestConversation = computed(() => selectedConversation.value?.id === '__debug_test_conversation__')
const filteredConversationMessages = computed(() => conversationMessages.value.filter(msg => !msg.isHidden))
const lastMessageHasError = computed(() => messagesStore.lastMessageHasError)
const hasWorldTree = computed(() => {
  return (conversationWorldTree.value?.timeline?.length > 0) || newConvWorldTree.value
})
function togglePanel() {
  conversationStore.togglePanelCollapsed()
}
function toggleWorldTree() {
  sidebarShowWorldTree.value = !sidebarShowWorldTree.value
}
function handleWorldTreeToggle(e: Event) {
  const target = e.target as HTMLInputElement
  worldTreeStore.setEnabled(target.checked)
}
function handlePanelToggle() {
  emit('switch-to-worldbook')
}
function handleSelectConversation(data: any) {
  if (data.conversation) {
    messagesStore.selectConversation(data.conversation, data.app?.app?.id)
  }
}
function handleShowNotification(event: any) {
  notificationStore.showNotification(event.message, event.type)
}
function handleConversationRestored() {
  conversationStore.pendingConversationRestore = null
}
function handleNewConversation(data: any) {
  conversationStore.pendingNewConv = data
}
function openFullWorldtree() {
  emit('switch-to-worldbook')
}
function unbindWorldtree() {
  worldTreeStore.clearConversationTree()
}
function handleCreateParallelWorld(node: any) {
  notificationStore.showNotification('平行世界功能开发中', 'info')
}
function handleNewConvNodeClick(node: any) {
  conversationStore.newConvSelectedNodeId = node.id
}
function confirmNodeSelection() {
  conversationStore.newConvSelectingNode = false
}
function cancelNodeSelection() {
  conversationStore.newConvSelectingNode = false
  conversationStore.newConvSelectedNodeId = null
}
function initWorldtreeFromRemote() {
  notificationStore.showNotification('从远端同步功能开发中', 'info')
}
function showProjectSelect() {
  notificationStore.showNotification('从项目加载功能开发中', 'info')
}
function closeParallelWorld() {
  conversationStore.showParallelWorldConv = false
  conversationStore.parallelWorldNewConv = null
}
function handleParallelWorldEnterKey(e: KeyboardEvent) {
  if (!e.shiftKey && parallelWorldConvInput.value.trim() && !sendingParallelWorldMsg.value) {
    e.preventDefault()
    sendParallelWorldMessage()
  }
}
function sendParallelWorldMessage() {
  notificationStore.showNotification('平行世界对话功能开发中', 'info')
}
function closeNewConversation() {
  conversationStore.newConversationStartPage = null
  conversationStore.newConvWorldTree = null
}
function handleNewConvEnterKey(e: KeyboardEvent) {
  if (!e.shiftKey && newConvMessageInput.value.trim() && !sendingNewConvMessage.value) {
    e.preventDefault()
    sendNewConversationMessage()
  }
}
async function sendNewConversationMessage() {
  if (!newConvMessageInput.value.trim()) return
  const startPage = newConversationStartPage.value
  if (!startPage?.appId) {
    notificationStore.showNotification('未找到应用信息', 'error')
    return
  }
  conversationStore.sendingNewConvMessage = true
  try {
    const query = newConvMessageInput.value.trim()
    const result = await conversationApi.create(token.value, startPage.appId, query, query.slice(0, 20))
    if (result?.conversation) {
      messagesStore.selectConversation(result.conversation, startPage.appId)
      if (result.messages) {
        messagesStore.setMessages(result.messages, result.messages.length, false)
      }
      conversationStore.newConversationStartPage = null
      conversationStore.newConvWorldTree = null
      conversationStore.newConvMessageInput = ''
      notificationStore.showNotification('对话已创建', 'success')
    }
  } catch (e: any) {
    notificationStore.showNotification('创建失败: ' + e.message, 'error')
  } finally {
    conversationStore.sendingNewConvMessage = false
  }
}
function closeConversation() {
  messagesStore.clearConversation()
}
async function showDeleteConfirm() {
  if (!selectedConversation.value) return
  const confirmed = await confirmStore.confirmDelete(selectedConversation.value.name || '此对话')
  if (!confirmed) return
  try {
    const appId = messagesStore.selectedConversationApp
    if (!appId) return
    await conversationApi.delete(token.value, appId, selectedConversation.value.id)
    messagesStore.clearConversation()
    notificationStore.showNotification('对话已删除', 'success')
    conversationPanelRef.value?.refreshList?.()
  } catch (e: any) {
    const msg = e?.message || String(e) || '未知错误'
    notificationStore.showNotification('删除失败: ' + msg, 'error')
  }
}
function showDebugReplyModal() {
  conversationStore.showDebugReplyModal = true
}
async function loadMoreMessages() {
  if (!selectedConversation.value || messagesStore.isLoading) return
  const appId = messagesStore.selectedConversationApp
  if (!appId) return
  messagesStore.isLoading = true
  try {
    const page = messagesStore.currentPage + 1
    const result = await conversationApi.getDetail(token.value, appId, selectedConversation.value.id)
    if (result?.messages) {
      messagesStore.appendMessages(result.messages)
      messagesStore.currentPage = page
      messagesStore.hasMore = false
    }
  } catch (e: any) {
    notificationStore.showNotification('加载失败: ' + e.message, 'error')
  } finally {
    messagesStore.isLoading = false
  }
}
function onMessagesScroll() {
  if (!messagesContainer.value) return
  const el = messagesContainer.value
  messagesStore.canScrollDown = el.scrollHeight - el.scrollTop - el.clientHeight > 100
}
function toggleAutoScroll() {
  messagesStore.autoScrollLocked = !messagesStore.autoScrollLocked
  if (messagesStore.autoScrollLocked && messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}
function handleEnterKey(e: KeyboardEvent) {
  if (!e.shiftKey && messageInput.value.trim() && !sendingMessage.value) {
    e.preventDefault()
    sendMessage()
  }
}
async function sendMessage() {
  if (!messageInput.value.trim() || !selectedConversation.value) return
  const query = messageInput.value.trim()
  messagesStore.isSending = true
  messagesStore.sendButtonState = 'sending'
  try {
    const newMessage = {
      id: Date.now().toString(),
      query: query,
      answer: '',
      created_at: new Date().toISOString(),
      model_id: '',
      message_points: 0,
      message_tokens: 0,
      answer_points: 0,
      answer_tokens: 0
    }
    messagesStore.addMessage(newMessage)
    messagesStore.clearInput()
    await nextTick()
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
    }
    notificationStore.showNotification('消息发送需要流式API支持', 'info')
  } catch (e: any) {
    notificationStore.showNotification('发送失败: ' + e.message, 'error')
  } finally {
    messagesStore.isSending = false
    messagesStore.sendButtonState = 'ready'
  }
}
function stopStreaming() {
  streamingStore.stop()
}
function getCleanUserQuery(query: string): string {
  if (!query) return ''
  return query.replace(/^<HIDE>[\s\S]*?<\/HIDE>\s*/i, '').trim()
}
function getMessageContent(answer: string): string {
  if (!answer) return ''
  return answer
    .replace(/<WORLD_TREE_COMMAND>[\s\S]*?<\/WORLD_TREE_COMMAND>/gi, '')
    .replace(/<T2I>[\s\S]*?<\/T2I>/gi, '')
    .replace(/<OPTIONS>[\s\S]*?<\/OPTIONS>/gi, '')
    .trim()
}
function renderMarkdown(text: string): string {
  if (!text) return ''
  let html = text
  const codeBlocks: string[] = []
  html = html.replace(/```\w*\n?([\s\S]*?)```/g, (match, code) => {
    const placeholder = `\x00CODEBLOCK${codeBlocks.length}\x00`
    const escapedCode = code.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').trim()
    codeBlocks.push(`<pre><code>${escapedCode}</code></pre>`)
    return placeholder
  })
  const inlineCodes: string[] = []
  html = html.replace(/`([^`]+)`/g, (match, code) => {
    const placeholder = `\x00INLINECODE${inlineCodes.length}\x00`
    const escapedCode = code.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
    inlineCodes.push(`<code>${escapedCode}</code>`)
    return placeholder
  })
  html = html.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  html = html.replace(/__([^_]+)__/g, '<strong>$1</strong>')
  html = html.replace(/(?<!\*)\*([^*]+)\*(?!\*)/g, '<em>$1</em>')
  html = html.replace(/(?<!_)_([^_]+)_(?!_)/g, '<em>$1</em>')
  html = html.replace(/~~([^~]+)~~/g, '<del>$1</del>')
  html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>')
  html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>')
  html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>')
  html = html.replace(/^\* (.+)$/gm, '<li>$1</li>')
  html = html.replace(/^- (.+)$/gm, '<li>$1</li>')
  html = html.replace(/^\d+\. (.+)$/gm, '<li>$1</li>')
  html = html.replace(/^> (.+)$/gm, '<blockquote>$1</blockquote>')
  html = html.replace(/\n/g, '<br>')
  inlineCodes.forEach((code, i) => {
    html = html.replace(`\x00INLINECODE${i}\x00`, code)
  })
  codeBlocks.forEach((block, i) => {
    html = html.replace(`\x00CODEBLOCK${i}\x00`, block)
  })
  return html
}
function formatMessageTime(time: string): string {
  if (!time) return ''
  const date = new Date(time)
  return date.toLocaleString()
}
function handleAnswerContextMenu(e: MouseEvent, msg: any) {
  e.preventDefault()
  conversationStore.showContextMenu(e.clientX, e.clientY, msg.answer, !!window.getSelection()?.toString())
}
async function regenerateMessage(msg: any) {
  if (!selectedConversation.value || !msg) return
  messagesStore.regeneratingMessageId = msg.id
  notificationStore.showNotification('重新生成功能需要流式API支持', 'info')
  messagesStore.regeneratingMessageId = null
}
function startEditMessage(msg: any) {
  messagesStore.startEditing(msg.id, msg.query)
}
async function saveEditedMessage(msg: any) {
  if (!selectedConversation.value || !msg) return
  messagesStore.isSendingEdit = true
  try {
    messagesStore.updateMessage(msg.id, { query: messagesStore.editingMessageContent })
    messagesStore.cancelEditing()
    notificationStore.showNotification('消息已更新', 'success')
  } catch (e: any) {
    notificationStore.showNotification('保存失败: ' + e.message, 'error')
  } finally {
    messagesStore.isSendingEdit = false
  }
}
function cancelEditMessage() {
  messagesStore.cancelEditing()
}
async function deleteMessage(msg: any) {
  if (!msg) return
  const confirmed = await confirmStore.confirmDelete('此消息')
  if (!confirmed) return
  messagesStore.removeMessage(msg.id)
  notificationStore.showNotification('消息已删除', 'success')
}
onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
})
onUnmounted(() => {
  document.removeEventListener('click', handleDocumentClick)
})
function handleDocumentClick() {
  conversationStore.hideContextMenu()
  conversationStore.hideSelectionPopup()
}
</script>
<style scoped>
.conversation-view {
  display: flex;
  flex: 1;
  height: 100%;
  overflow: hidden;
}
.conversation-sidebar {
  display: flex;
  flex-direction: column;
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  transition: width 0.2s, min-width 0.2s;
  overflow: hidden;
}
.conversation-sidebar.collapsed {
  width: 40px !important;
  min-width: 40px !important;
}
.panel-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
}
.panel-header h2 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
  white-space: nowrap;
}
.collapse-icon {
  display: inline-flex;
  width: 16px;
  height: 16px;
  transition: transform 0.2s;
}
.collapse-icon svg {
  width: 100%;
  height: 100%;
  color: var(--text-secondary);
}
.collapse-icon.expanded {
  transform: rotate(90deg);
}
.sidebar-main-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.worldtree-section {
  border-top: 1px solid var(--border-color);
}
.worldtree-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
}
.worldtree-header h2 {
  margin: 0;
  font-size: 13px;
  font-weight: 500;
  flex: 1;
}
.worldtree-switch {
  position: relative;
  display: inline-flex;
  align-items: center;
}
.worldtree-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.switch-slider {
  width: 32px;
  height: 18px;
  background: var(--bg-tertiary);
  border-radius: 9px;
  position: relative;
  transition: background 0.2s;
  cursor: pointer;
}
.switch-slider::after {
  content: '';
  position: absolute;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: white;
  top: 2px;
  left: 2px;
  transition: transform 0.2s;
}
.worldtree-switch input:checked + .switch-slider {
  background: var(--primary-color);
}
.worldtree-switch input:checked + .switch-slider::after {
  transform: translateX(14px);
}
.worldtree-switch.disabled .switch-slider {
  opacity: 0.5;
  cursor: not-allowed;
}
.btn-expand-tree,
.btn-unbind {
  padding: 4px 8px;
  border: none;
  background: transparent;
  cursor: pointer;
  border-radius: 4px;
  color: var(--text-secondary);
}
.btn-expand-tree:hover,
.btn-unbind:hover {
  background: var(--bg-hover);
}
.world-tree-content {
  padding: 12px;
  overflow-y: auto;
}
.world-tree-content.disabled {
  opacity: 0.5;
  pointer-events: none;
}
.readonly-timeline-graph {
  background: var(--bg-tertiary);
  border-radius: 8px;
  padding: 12px;
}
.bind-worldtree-area {
  text-align: center;
  padding: 20px;
}
.bind-worldtree-hint {
  color: var(--text-secondary);
  margin-bottom: 12px;
}
.bind-worldtree-actions {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.outline-empty {
  text-align: center;
  padding: 20px;
  color: var(--text-secondary);
}
.node-select-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  margin-top: 12px;
}
.conversation-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-primary);
}
.conversation-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}
.messages-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.messages-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  gap: 12px;
}
.messages-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.header-actions {
  display: flex;
  gap: 8px;
}
.btn-header {
  padding: 6px 12px;
  border: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}
.btn-header:hover {
  background: var(--bg-hover);
}
.btn-header.btn-danger {
  background: #ef4444;
  color: white;
}
.btn-header.btn-danger:hover {
  background: #dc2626;
}
.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.load-more {
  text-align: center;
  padding: 12px;
}
.message-item {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.message-query,
.message-answer {
  padding: 12px 16px;
  border-radius: 12px;
  background: var(--bg-secondary);
}
.message-query {
  border-left: 3px solid var(--primary-color);
}
.message-answer {
  border-left: 3px solid #10b981;
}
.message-role {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.model-badge {
  font-size: 10px;
  padding: 2px 6px;
  background: var(--bg-tertiary);
  border-radius: 4px;
  font-weight: normal;
}
.message-content {
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-primary);
  word-break: break-word;
}
.message-edit-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.message-edit-textarea {
  width: 100%;
  min-height: 100px;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  resize: none;
}
.message-edit-actions {
  display: flex;
  gap: 8px;
}
.btn-edit-save,
.btn-edit-cancel {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}
.btn-edit-save {
  background: var(--primary-color);
  color: white;
}
.btn-edit-cancel {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.message-error-retry {
  text-align: center;
}
.btn-error-retry {
  padding: 8px 16px;
  background: #ef4444;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}
.message-footer {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 8px 16px;
  font-size: 11px;
  color: var(--text-tertiary);
}
.footer-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.footer-actions {
  display: flex;
  gap: 8px;
}
.btn-msg-action {
  padding: 4px 8px;
  border: none;
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border-radius: 4px;
  cursor: pointer;
  font-size: 11px;
}
.btn-msg-action:hover {
  background: var(--bg-hover);
}
.btn-msg-action.btn-danger {
  color: #ef4444;
}
.empty-hint {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}
.message-input-area {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid var(--border-color);
  background: var(--bg-secondary);
}
.message-input {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid var(--border-color);
  border-radius: 12px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  resize: none;
  max-height: 200px;
}
.input-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}
.scroll-bottom-btn {
  width: 36px;
  height: 36px;
  border: none;
  background: var(--bg-tertiary);
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}
.scroll-bottom-btn.locked {
  background: var(--primary-color);
  color: white;
}
.send-btn {
  padding: 10px 20px;
  background: var(--primary-color);
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.2s;
}
.send-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.send-btn:not(:disabled):hover {
  background: var(--primary-hover);
}
.send-btn-streaming {
  background: #22c55e;
}
.send-btn-streaming:hover {
  background: #16a34a;
}
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
.loading-spinner-sm {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.new-conversation-area,
.parallel-world-conv-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.new-conv-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}
.new-conv-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}
.start-page-container {
  flex: 1;
  overflow: hidden;
}
.start-page-iframe {
  width: 100%;
  height: 100%;
  border: none;
}
.empty-conv-hint {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}
.new-conv-input-area {
  display: flex;
  gap: 12px;
  padding: 16px;
  border-top: 1px solid var(--border-color);
}
.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
}
.btn-sm {
  padding: 6px 12px;
  font-size: 12px;
}
.btn-primary {
  background: var(--primary-color);
  color: white;
}
.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.btn-danger {
  background: #ef4444;
  color: white;
}
</style>
