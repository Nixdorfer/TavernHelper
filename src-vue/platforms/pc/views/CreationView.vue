<template>
  <div class="creation-view">
    <aside
      :class="['creation-sidebar', { collapsed: creationPanelCollapsed }]"
      :style="{ width: creationPanelCollapsed ? '40px' : creationPanelWidth + 'px' }">
    <div class="panel-header creation-header" @click="toggleCreationPanel">
      <span :class="['collapse-icon', { expanded: !creationPanelCollapsed }]"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg></span>
      <h2 v-show="!creationPanelCollapsed">创作项目</h2>
      <button v-show="!creationPanelCollapsed" class="btn-refresh" @click.stop="loadCreationApps" :disabled="creationLoading" title="刷新">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M23 4v6h-6M1 20v-6h6"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
        </svg>
      </button>
    </div>
    <div v-show="!creationPanelCollapsed" class="creation-sidebar-content">
      <div class="creation-list">
        <div v-if="creationLoading" class="creation-loading">
          <span class="loading-spinner"></span>
          <span>加载中...</span>
        </div>
        <div v-else-if="creationApps.length === 0" class="creation-empty">
          暂无创作项目
        </div>
        <div
          v-else
          v-for="app in creationApps"
          :key="app.id"
          :class="['creation-item', { active: selectedCreationApp?.id === app.id }]"
          @click="selectCreationApp(app)"
        >
          <div class="creation-item-icon">
            <img v-if="app.cover" :src="app.cover" />
            <span v-else class="creation-item-icon-placeholder">{{ app.name?.charAt(0) || '?' }}</span>
          </div>
          <div class="creation-item-info">
            <div class="creation-item-name">{{ app.name }}</div>
            <div class="creation-item-mode">{{ app.mode === 'advanced-chat' ? '高级对话' : app.mode }}</div>
          </div>
          <button class="btn-delete-creation" @click.stop="confirmDeleteCreation(app)" title="删除">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
    <div v-show="creationPanelCollapsed && creationEditorTab === 'description'" class="editor-tools-sidebar">
      <div class="editor-tool-group">
        <button :class="['editor-tool-btn', { active: editorShowBlocks }]" @click="toggleEditorBlocks" title="组件">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/></svg>
        </button>
        <button :class="['editor-tool-btn', { active: editorShowStyles }]" @click="toggleEditorStyles" title="样式">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
        </button>
      </div>
      <div class="editor-tool-divider"></div>
      <button :class="['editor-tool-btn', { active: editorShowCode }]" @click="toggleEditorCode" title="代码">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16,18 22,12 16,6"/><polyline points="8,6 2,12 8,18"/></svg>
      </button>
    </div>
  </aside>
  <div class="creation-content">
    <div v-if="!selectedCreationApp" class="creation-placeholder">
      <span>请选择一个创作项目</span>
    </div>
    <div v-else-if="creationDetailLoading" class="creation-loading-center">
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>
    <div v-else class="creation-editor-wrapper">
      <div class="creation-main-editor">
        <div class="creation-editor-body">
          <div v-show="creationEditorTab === 'basic'" class="creation-tab-content">
            <div class="form-section">
              <label class="form-label">标题</label>
              <input v-model="creationConfig.app.name" class="form-input" placeholder="作品名称..." @input="markCreationUnsaved" />
            </div>
            <div class="form-section">
              <label class="form-label">简介</label>
              <textarea v-model="creationConfig.app.summary" class="form-textarea" rows="3" placeholder="简介..." @input="markCreationUnsaved"></textarea>
            </div>
            <div class="form-section">
              <label class="form-label">图片设置</label>
              <div class="image-columns">
                <div class="image-column">
                  <div class="column-label">封面</div>
                  <div class="column-preview">
                    <img v-if="creationConfig.app.cover" :src="creationConfig.app.cover" />
                    <div v-else class="preview-placeholder">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                        <circle cx="8.5" cy="8.5" r="1.5"/>
                        <path d="M21 15l-5-5L5 21"/>
                      </svg>
                    </div>
                  </div>
                  <div class="column-buttons">
                    <button class="btn-column" @click="openImagePicker('cover')">选择</button>
                    <button class="btn-column btn-danger" :class="{ disabled: !creationConfig.app.cover }" :disabled="!creationConfig.app.cover" @click="clearCreationImage('cover')">删除</button>
                  </div>
                </div>
                <div class="image-column">
                  <div class="column-label">手机背景图</div>
                  <div class="column-preview">
                    <img v-if="creationConfig.bg_mobile" :src="creationConfig.bg_mobile" />
                    <div v-else class="preview-placeholder">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                        <circle cx="8.5" cy="8.5" r="1.5"/>
                        <path d="M21 15l-5-5L5 21"/>
                      </svg>
                    </div>
                  </div>
                  <div class="column-buttons">
                    <button class="btn-column" @click="openImagePicker('bg_mobile')">选择</button>
                    <button class="btn-column btn-danger" :class="{ disabled: !creationConfig.bg_mobile }" :disabled="!creationConfig.bg_mobile" @click="clearCreationImage('bg_mobile')">删除</button>
                  </div>
                </div>
                <div class="image-column">
                  <div class="column-label">PC背景图</div>
                  <div class="column-preview">
                    <img v-if="creationConfig.bg_image" :src="creationConfig.bg_image" />
                    <div v-else class="preview-placeholder">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                        <rect x="3" y="3" width="18" height="18" rx="2"/>
                        <circle cx="8.5" cy="8.5" r="1.5"/>
                        <path d="M21 15l-5-5L5 21"/>
                      </svg>
                    </div>
                  </div>
                  <div class="column-buttons">
                    <button class="btn-column" @click="openImagePicker('bg_image')">选择</button>
                    <button class="btn-column btn-danger" :class="{ disabled: !creationConfig.bg_image }" :disabled="!creationConfig.bg_image" @click="clearCreationImage('bg_image')">删除</button>
                  </div>
                </div>
              </div>
            </div>
            <div class="form-row">
              <div class="form-section third">
                <label class="form-label">语言</label>
                <select v-model="creationConfig.app.language" class="form-select" @change="markCreationUnsaved">
                  <option value="zh-Hans">简体中文</option>
                  <option value="zh-Hant">繁体中文</option>
                  <option value="en">English</option>
                  <option value="ja">日本語</option>
                  <option value="ko">한국어</option>
                </select>
              </div>
              <div class="form-section third">
                <label class="form-label">面向群体</label>
                <select v-model="creationConfig.app.gender" class="form-select" @change="markCreationUnsaved">
                  <option :value="1">全部</option>
                  <option :value="2">男性向</option>
                  <option :value="3">女性向</option>
                </select>
              </div>
              <div class="form-section third">
                <label class="form-label">MOD权限</label>
                <select v-model="creationConfig.app.mod_permission" class="form-select" @change="markCreationUnsaved">
                  <option :value="0">无</option>
                  <option :value="1">仅查看</option>
                  <option :value="2">评论</option>
                  <option :value="3">编辑</option>
                  <option :value="4">完全控制</option>
                </select>
              </div>
            </div>
            <div class="form-section">
              <label class="form-label">屏蔽词</label>
              <div class="banned-words-tags">
                <span v-for="(word, idx) in creationConfig.banned_words" :key="idx" class="banned-word-tag">
                  {{ word }}
                  <button class="tag-remove" @click="removeBannedWord(idx)">×</button>
                </span>
                <input v-model="bannedWordInput" class="banned-word-input" placeholder="输入后回车..." @keydown.enter.prevent="addBannedWord" />
              </div>
            </div>
          </div>
          <div v-show="creationEditorTab === 'interaction'" class="creation-tab-content">
            <div class="form-section">
              <label class="form-label">开场白</label>
              <textarea v-model="creationConfig.opening_statement" class="form-textarea" rows="4" placeholder="开场白..." @input="markCreationUnsaved"></textarea>
            </div>
            <div class="form-section">
              <label class="form-label">建议问题</label>
              <div class="list-items">
                <div v-for="(q, idx) in creationConfig.suggested_questions" :key="idx" class="list-item">
                  <input v-model="creationConfig.suggested_questions[idx]" class="form-input" placeholder="建议问题..." @input="markCreationUnsaved" />
                  <button class="btn-remove-sm" @click="removeSuggestedQuestion(idx)">×</button>
                </div>
                <button class="btn-add-item" @click="addSuggestedQuestion">+ 添加问题</button>
              </div>
            </div>
            <div class="form-section">
              <label class="form-label">快捷回复</label>
              <div class="list-items">
                <div v-for="(cmd, idx) in creationConfig.shortcut_commands" :key="idx" class="list-item shortcut-item">
                  <input v-model="cmd.label" class="form-input" placeholder="标签..." />
                  <input v-model="cmd.command" class="form-input" placeholder="命令..." />
                  <button class="btn-remove-sm" @click="removeShortcutCommand(idx)">×</button>
                </div>
                <button class="btn-add-item" @click="addShortcutCommand">+ 添加快捷回复</button>
              </div>
            </div>
          </div>
          <div v-show="creationEditorTab === 'description'" class="creation-tab-content creation-description-area">
            <GrapesEditor
              ref="grapesEditorRef"
              v-model="creationConfig.app.description"
              :view-mode="descriptionViewMode"
              :hide-toolbar="true"
              external-views-container="#gjs-views-container"
              @save="saveCreationDescription"
              @change="markCreationUnsaved"
            />
          </div>
          <div v-show="creationEditorTab === 'prompt'" class="creation-tab-content">
            <div class="form-section">
              <label class="form-label">全局词 (pre_prompt)</label>
              <textarea v-model="creationConfig.pre_prompt" class="form-textarea" rows="8" placeholder="全局提示词..." @input="markCreationUnsaved"></textarea>
            </div>
            <div class="form-section">
              <label class="form-label">前置词 (pre_text)</label>
              <textarea v-model="creationConfig.pre_text" class="form-textarea" rows="6" placeholder="前置词..." @input="markCreationUnsaved"></textarea>
            </div>
            <div class="form-section">
              <label class="form-label">后置词 (post_text)</label>
              <textarea v-model="creationConfig.post_text" class="form-textarea" rows="6" placeholder="后置词..." @input="markCreationUnsaved"></textarea>
            </div>
          </div>
          <div v-show="creationEditorTab === 'worldbook'" class="creation-tab-content creation-worldbook-tab">
            <div v-if="creationSelectedWorldBookEntry" class="worldbook-entry-editor">
              <div class="entry-editor-header">
                <div class="form-section trigger-section">
                  <div class="trigger-mode-switch">
                    <span :class="['mode-option', { active: !creationSelectedWorldBookEntry.isOrMode }]" @click="creationSelectedWorldBookEntry.isOrMode = false">和</span>
                    <span :class="['mode-option', { active: creationSelectedWorldBookEntry.isOrMode }]" @click="creationSelectedWorldBookEntry.isOrMode = true">或</span>
                  </div>
                  <label class="form-label">{{ creationSelectedWorldBookEntry.isOrMode ? '或触发词' : '和触发词' }}</label>
                  <div class="trigger-input-row">
                    <input v-model="creationWbKeywordInput" class="form-input" placeholder="输入关键词后按回车添加..." @keydown.enter.prevent="addCreationWbKeyword" />
                  </div>
                  <div class="trigger-tags" v-if="getCreationWbKeywords().length">
                    <span v-for="(kw, idx) in getCreationWbKeywords()" :key="idx" class="trigger-tag">
                      {{ kw }}
                      <button class="tag-remove" @click="removeCreationWbKeyword(idx)">×</button>
                    </span>
                  </div>
                </div>
              </div>
              <div class="entry-editor-controls">
                <div class="control-group">
                  <label class="control-label">触发区域</label>
                  <div class="region-toggles">
                    <label class="region-toggle-switch">
                      <input type="checkbox" v-model="creationSelectedWorldBookEntry.keySystem">
                      <span class="slider"></span>
                      <span class="toggle-label">系统</span>
                    </label>
                    <label class="region-toggle-switch">
                      <input type="checkbox" v-model="creationSelectedWorldBookEntry.keyUser">
                      <span class="slider"></span>
                      <span class="toggle-label">用户</span>
                    </label>
                    <label class="region-toggle-switch">
                      <input type="checkbox" v-model="creationSelectedWorldBookEntry.keyAI">
                      <span class="slider"></span>
                      <span class="toggle-label">AI</span>
                    </label>
                  </div>
                </div>
                <div class="control-group">
                  <label class="control-label">插入位置</label>
                  <div class="region-toggles">
                    <label class="region-toggle-switch">
                      <input type="checkbox" v-model="creationSelectedWorldBookEntry.posGlobal">
                      <span class="slider"></span>
                      <span class="toggle-label">全局词</span>
                    </label>
                    <label class="region-toggle-switch">
                      <input type="checkbox" v-model="creationSelectedWorldBookEntry.posPre">
                      <span class="slider"></span>
                      <span class="toggle-label">前置词</span>
                    </label>
                    <label class="region-toggle-switch">
                      <input type="checkbox" v-model="creationSelectedWorldBookEntry.posPost">
                      <span class="slider"></span>
                      <span class="toggle-label">后置词</span>
                    </label>
                  </div>
                </div>
              </div>
              <div class="entry-editor-body">
                <label class="form-label">内容</label>
                <textarea v-model="creationSelectedWorldBookEntry.value" class="form-textarea entry-value-textarea" placeholder="世界书条目内容..."></textarea>
              </div>
            </div>
            <div v-else class="worldbook-empty-hint">
              请在右侧世界书列表中选择或添加条目
            </div>
          </div>
        </div>
      </div>
      <aside class="creation-right-panel">
        <div class="right-panel-tabs">
          <div :class="['right-panel-tab', { active: creationEditorTab === 'basic' }]" @click="creationEditorTab = 'basic'">
            <span class="tab-label">基础信息</span>
          </div>
          <div :class="['right-panel-tab', { active: creationEditorTab === 'interaction' }]" @click="creationEditorTab = 'interaction'">
            <span class="tab-label">交互信息</span>
          </div>
          <div :class="['right-panel-tab', { active: creationEditorTab === 'description' }]" @click="switchToDescriptionTab">
            <span class="tab-label">对话首页</span>
          </div>
          <div :class="['right-panel-tab', { active: creationEditorTab === 'prompt' }]" @click="creationEditorTab = 'prompt'">
            <span class="tab-label">提示词</span>
          </div>
        </div>
        <div class="right-panel-worldbook">
          <div class="worldbook-section-header" @click="toggleWorldBookExpand">
            <span :class="['section-collapse-icon', { expanded: creationWorldBookExpanded }]"><svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg></span>
            <span class="section-title">世界书</span>
            <span class="section-count">({{ creationConfig.world_book.length }})</span>
            <button class="btn-add-wb-circle" @click.stop="addCreationWorldBookCard" title="添加条目">+</button>
          </div>
          <div :class="['worldbook-section-content', { collapsed: !creationWorldBookExpanded }]">
            <div class="worldbook-panel-list">
              <div v-for="(entry, idx) in creationConfig.world_book" :key="entry.id || idx"
                   class="worldbook-panel-item"
                   :class="{ active: creationSelectedWorldBookEntry && creationSelectedWorldBookEntry.id === entry.id }"
                   @click="creationEditorTab = 'worldbook'; selectCreationWorldBookEntry(entry)">
                <span class="item-name">{{ entry.key || '未命名条目' }}</span>
                <button class="btn-delete-wb-item" @click.stop="removeCreationWorldBookEntry(entry, idx)">×</button>
              </div>
            </div>
          </div>
        </div>
        <div v-show="creationEditorTab === 'description'" class="right-panel-divider"></div>
        <div v-show="creationEditorTab === 'description'" class="right-panel-gjs-tools">
          <div id="gjs-views-container"></div>
        </div>
        <div class="creation-save-status-footer">
          <button class="btn-test-chat" @click="openTestDrawer">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <span>测试对话</span>
          </button>
          <div v-if="creationJustSaved" class="save-status-pill saved">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12l5 5L20 7"/></svg>
            <span>已保存</span>
          </div>
          <div v-else-if="creationHasUnsavedChanges" class="save-status-pill unsaved clickable" @click="saveLocalCreationConfig">
            <span :class="['save-status-dot', saveStatusDotClass]"></span>
            <span>{{ creationLastSaveTimeAgo }}</span>
          </div>
          <div v-else class="save-status-pill saved">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12l5 5L20 7"/></svg>
            <span>已保存</span>
          </div>
        </div>
      </aside>
    </div>
  </div>
  <ImagePickerModal
    v-if="showImagePicker"
    @close="showImagePicker = false"
    @confirm="handleImagePickerConfirm"
  />
  <Transition name="drawer-slide">
    <div v-if="showTestDrawer" class="test-drawer-overlay" @mousedown.self="closeTestDrawer">
      <div class="test-drawer">
        <div class="test-drawer-header">
          <span class="test-drawer-title">测试对话</span>
          <select v-model="selectedTestProvider" class="test-provider-select">
            <option value="claude">Claude</option>
            <option value="gemini">Gemini</option>
            <option value="grok">Grok</option>
          </select>
          <button class="btn-close-drawer" @click="closeTestDrawer">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M18 6L6 18M6 6l12 12"/>
            </svg>
          </button>
        </div>
        <div class="test-drawer-body">
          <div class="test-preview-area">
            <div class="preview-label">首页预览</div>
            <div class="preview-content" v-html="creationConfig.app.description || '暂无内容'"></div>
          </div>
          <div class="test-messages-area">
            <div v-if="testMessages.length === 0" class="test-messages-empty">
              发送消息开始测试对话
            </div>
            <div v-for="(msg, idx) in testMessages" :key="idx" :class="['test-message', msg.role]">
              <div class="message-bubble">{{ msg.content }}</div>
            </div>
            <div v-if="sendingTestMessage" class="test-message assistant loading">
              <div class="message-bubble">
                <span class="typing-indicator">
                  <span></span><span></span><span></span>
                </span>
              </div>
            </div>
          </div>
        </div>
        <div class="test-drawer-footer">
          <input
            v-model="testMessageInput"
            class="test-input"
            placeholder="输入消息..."
            @keydown.enter.prevent="sendTestMessage"
            :disabled="sendingTestMessage"
          />
          <button class="btn-send-test" @click="sendTestMessage" :disabled="sendingTestMessage || !testMessageInput.trim()">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </Transition>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useAuthStore, useNotificationStore, useConfirmStore } from '@/stores'
import { creationApi } from '@/api'
import GrapesEditor from '@/components/business/GrapesEditor.vue'
import ImagePickerModal from '@/components/modals/ImagePickerModal.vue'
interface CreationApp {
  id: string
  name: string
  cover?: string
  mode?: string
  folder_name?: string
  config_updated_at?: string
}
interface WorldBookEntry {
  id: string
  key: string
  keywords: string[]
  isOrMode: boolean
  value: string
  keySystem: boolean
  keyUser: boolean
  keyAI: boolean
  posGlobal: boolean
  posPre: boolean
  posPost: boolean
  isFolder: boolean
  parentId?: string
}
interface ShortcutCommand {
  label: string
  command: string
}
interface TestMessage {
  role: 'user' | 'assistant'
  content: string
}
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const confirmStore = useConfirmStore()
const token = computed(() => authStore.token)
const grapesEditorRef = ref<InstanceType<typeof GrapesEditor> | null>(null)
const creationPanelWidth = ref(280)
const creationPanelCollapsed = ref(false)
const creationApps = ref<CreationApp[]>([])
const creationLoading = ref(false)
const selectedCreationApp = ref<CreationApp | null>(null)
const creationDetailLoading = ref(false)
const creationConfig = reactive({
  pre_prompt: '',
  pre_text: '',
  post_text: '',
  opening_statement: '',
  suggested_questions: [] as string[],
  world_book: [] as WorldBookEntry[],
  cg_book: [] as any[],
  banned_words: [] as string[],
  shortcut_commands: [] as ShortcutCommand[],
  preset_type: 1,
  preset_chats: [] as any[],
  bg_image: '',
  bg_mobile: '',
  builtInCss: '',
  app: {
    name: '',
    description: '',
    summary: '',
    language: 'zh-Hans',
    gender: 1,
    cover: '',
    cover_tiny: '',
    mod_permission: 4,
    is_anonymous: false,
    update_content: '',
    is_available_not_public: true
  }
})
const creationLocalFolderName = ref<string | null>(null)
const creationLastSavedAt = ref<number | null>(null)
const creationHasUnsavedChanges = ref(false)
const creationJustSaved = ref(false)
const creationEditorTab = ref('basic')
const descriptionViewMode = ref('visual')
const editorShowBlocks = ref(false)
const editorShowStyles = ref(false)
const editorShowCode = ref(false)
const creationWorldBookExpanded = ref(true)
const creationSelectedWorldBookEntry = ref<WorldBookEntry | null>(null)
const creationWbKeywordInput = ref('')
const bannedWordInput = ref('')
const showImagePicker = ref(false)
const imagePickerTarget = ref('')
let creationTempAutoSaveTimer: ReturnType<typeof setInterval> | null = null
const showTestDrawer = ref(false)
const testMessages = ref<TestMessage[]>([])
const testMessageInput = ref('')
const sendingTestMessage = ref(false)
const selectedTestProvider = ref('claude')
const saveStatusDotClass = computed(() => {
  if (!creationLastSavedAt.value) return 'red'
  const diff = Date.now() - creationLastSavedAt.value
  if (diff < 60000) return 'green'
  if (diff < 300000) return 'yellow'
  return 'red'
})
const creationLastSaveTimeAgo = computed(() => {
  if (!creationLastSavedAt.value) return '未保存'
  const diff = Date.now() - creationLastSavedAt.value
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  return `${hours}小时前`
})
function toggleCreationPanel() {
  creationPanelCollapsed.value = !creationPanelCollapsed.value
}
async function loadCreationApps() {
  if (creationLoading.value) return
  creationLoading.value = true
  try {
    if (!token.value) {
      const localCreations = await creationApi.getLocalCreations()
      creationApps.value = (localCreations || []).map((c: any) => ({
        id: c.remote_id || c.folder_name,
        name: c.name || c.folder_name,
        cover: c.cover || '',
        mode: 'advanced-chat',
        folder_name: c.folder_name,
        config_updated_at: c.config_updated_at || c.updated_at
      }))
      return
    }
    const apps = await creationApi.getApps(token.value)
    creationApps.value = apps || []
  } catch (e: any) {
    notificationStore.showNotification('加载作品列表失败: ' + e.message, 'error')
  } finally {
    creationLoading.value = false
  }
}
async function selectCreationApp(app: CreationApp) {
  clearCreationTempAutoSave()
  selectedCreationApp.value = app
  creationPanelCollapsed.value = true
  startCreationTempAutoSave()
  if (!token.value) {
    const localData = await creationApi.getLocalCreationConfig(app.folder_name!).catch(() => null)
    if (localData) {
      localData.folder_name = app.folder_name
      loadCreationDetailFromLocal(localData)
    }
    return
  }
  const localData = await creationApi.getLocalCreationByRemoteId(app.id).catch(() => null)
  if (localData && localData.config_updated_at) {
    const localUpdatedAt = localData.config_updated_at
    const remoteUpdatedAt = app.config_updated_at
    if (remoteUpdatedAt && remoteUpdatedAt > localUpdatedAt) {
      const confirmed = await confirmStore.showConfirmDialog({
        title: '远程更新检测',
        message: `检测到项目"${app.name}"在远程有更新，是否覆盖本地？`,
        confirmText: '覆盖本地',
        cancelText: '使用本地'
      })
      if (confirmed) {
        await loadCreationDetailFromRemote(app.id, localData.folder_name)
      } else {
        loadCreationDetailFromLocal(localData)
      }
      return
    }
    loadCreationDetailFromLocal(localData)
    return
  }
  await loadCreationDetailFromRemote(app.id, localData ? localData.folder_name : null)
}
function loadCreationDetailFromLocal(localData: any) {
  creationDetailLoading.value = true
  try {
    creationConfig.opening_statement = localData.opening_statement || ''
    creationConfig.suggested_questions = localData.suggested_questions || []
    creationConfig.pre_prompt = localData.pre_prompt || ''
    creationConfig.pre_text = localData.pre_text || ''
    creationConfig.post_text = localData.post_text || ''
    creationConfig.world_book = localData.world_book || []
    creationConfig.cg_book = localData.cg_book || []
    creationConfig.banned_words = localData.banned_words || []
    creationConfig.shortcut_commands = localData.shortcut_commands || []
    creationConfig.preset_type = localData.preset_type || 1
    creationConfig.preset_chats = localData.preset_chats || []
    creationConfig.bg_image = localData.bg_image || ''
    creationConfig.bg_mobile = localData.bg_mobile || ''
    creationConfig.builtInCss = localData.builtInCss || ''
    creationConfig.app.name = localData.name || ''
    creationConfig.app.description = localData.description || ''
    creationConfig.app.summary = localData.summary || ''
    creationConfig.app.language = localData.language || 'zh-Hans'
    creationConfig.app.gender = localData.gender || 1
    creationConfig.app.cover = localData.cover || ''
    creationConfig.app.cover_tiny = localData.cover_tiny || ''
    creationConfig.app.mod_permission = localData.mod_permission !== undefined ? localData.mod_permission : 4
    creationConfig.app.is_anonymous = localData.is_anonymous || false
    creationConfig.app.update_content = localData.update_content || ''
    creationConfig.app.is_available_not_public = localData.is_available_not_public !== undefined ? localData.is_available_not_public : true
    creationLocalFolderName.value = localData.folder_name
    creationLastSavedAt.value = Date.now()
    creationHasUnsavedChanges.value = false
    creationJustSaved.value = true
  } finally {
    creationDetailLoading.value = false
  }
}
async function loadCreationDetailFromRemote(appId: string, existingFolderName: string | null) {
  creationDetailLoading.value = true
  try {
    const mc = await creationApi.getModelConfig(token.value, appId)
    if (mc) {
      loadCreationConfigFromCache(mc)
    }
    const appResp = await creationApi.getAppDetail(token.value, appId)
    if (appResp) {
      creationConfig.app.name = appResp.name || ''
      creationConfig.app.description = appResp.description || ''
      creationConfig.app.summary = appResp.summary || ''
      creationConfig.app.language = appResp.language || 'zh-Hans'
      creationConfig.app.gender = appResp.gender || 1
      creationConfig.app.cover = appResp.cover || ''
      creationConfig.app.cover_tiny = appResp.cover_tiny || ''
      creationConfig.app.mod_permission = appResp.mod_permission !== undefined ? appResp.mod_permission : 4
      creationConfig.app.is_anonymous = appResp.is_anonymous || false
      creationConfig.app.update_content = appResp.update_content || ''
      creationConfig.app.is_available_not_public = appResp.is_available_not_public !== undefined ? appResp.is_available_not_public : true
    }
    let folderName = existingFolderName
    if (!folderName) {
      const safeName = (appResp?.name || 'project').replace(/[\\/:*?"<>|]/g, '_')
      folderName = `${safeName}_${appId}`
    }
    creationLocalFolderName.value = folderName
    await saveLocalCreationConfig()
  } catch (e: any) {
    notificationStore.showNotification('加载作品详情失败: ' + e.message, 'error')
  } finally {
    creationDetailLoading.value = false
  }
}
function loadCreationConfigFromCache(mc: any) {
  creationConfig.opening_statement = mc.opening_statement || ''
  creationConfig.suggested_questions = mc.suggested_questions || []
  creationConfig.pre_prompt = mc.pre_prompt || ''
  creationConfig.pre_text = mc.pre_text || ''
  creationConfig.post_text = mc.post_text || ''
  creationConfig.world_book = mc.world_book || []
  creationConfig.cg_book = mc.cg_book || []
  creationConfig.banned_words = mc.banned_words || []
  creationConfig.shortcut_commands = mc.shortcut_commands || []
  creationConfig.preset_type = mc.preset_type || 1
  creationConfig.preset_chats = mc.preset_chats || []
  creationConfig.bg_image = mc.bg_image || ''
  creationConfig.bg_mobile = mc.bg_mobile || ''
  creationConfig.builtInCss = mc.builtInCss || ''
}
async function saveLocalCreationConfig() {
  const folderName = creationLocalFolderName.value
  if (!folderName) return
  if (grapesEditorRef.value) {
    const html = (grapesEditorRef.value as any).getFullHtml()
    if (html) {
      creationConfig.app.description = html
    }
  }
  try {
    const config = {
      name: creationConfig.app.name,
      description: '',
      summary: creationConfig.app.summary || '',
      language: creationConfig.app.language || '',
      gender: creationConfig.app.gender || 1,
      pre_prompt: creationConfig.pre_prompt || null,
      pre_text: creationConfig.pre_text || '',
      post_text: creationConfig.post_text || '',
      cover: creationConfig.app.cover || null,
      cover_tiny: creationConfig.app.cover_tiny || null,
      bg_image: creationConfig.bg_image || null,
      bg_mobile: creationConfig.bg_mobile || null,
      builtInCss: creationConfig.builtInCss || null,
      opening_statement: creationConfig.opening_statement || null,
      suggested_questions: creationConfig.suggested_questions.filter(q => q.trim()),
      world_book: creationConfig.world_book || [],
      cg_book: creationConfig.cg_book || [],
      banned_words: creationConfig.banned_words || [],
      shortcut_commands: creationConfig.shortcut_commands || [],
      preset_type: creationConfig.preset_type || 1,
      preset_chats: creationConfig.preset_chats || [],
      mod_permission: creationConfig.app.mod_permission || 4,
      is_anonymous: creationConfig.app.is_anonymous || false,
      is_available_not_public: creationConfig.app.is_available_not_public !== undefined ? creationConfig.app.is_available_not_public : true,
      remote_id: selectedCreationApp.value?.id || null,
      config_updated_at: selectedCreationApp.value?.config_updated_at || null,
      updated_at: new Date().toISOString()
    }
    await creationApi.saveLocalCreation(folderName, config)
    if (creationConfig.app.description) {
      await creationApi.saveLocalCreationPage(folderName, 'description', creationConfig.app.description)
    }
    creationLastSavedAt.value = Date.now()
    creationHasUnsavedChanges.value = false
    creationJustSaved.value = true
    setTimeout(() => {
      creationJustSaved.value = false
    }, 2000)
  } catch (e: any) {
    notificationStore.showNotification('保存失败: ' + e.message, 'error')
  }
}
function handleSaveKeydown(e: KeyboardEvent) {
  if (!selectedCreationApp.value) return
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    saveLocalCreationConfig()
  }
}
function markCreationUnsaved() {
  if (creationLocalFolderName.value) {
    creationHasUnsavedChanges.value = true
    creationJustSaved.value = false
  }
}
function startCreationTempAutoSave() {
  clearCreationTempAutoSave()
  creationTempAutoSaveTimer = setInterval(() => {
    if (creationHasUnsavedChanges.value) {
      saveLocalCreationConfig()
    }
  }, 30000)
}
function clearCreationTempAutoSave() {
  if (creationTempAutoSaveTimer) {
    clearInterval(creationTempAutoSaveTimer)
    creationTempAutoSaveTimer = null
  }
}
async function confirmDeleteCreation(app: CreationApp) {
  const confirmed = await confirmStore.showConfirmDialog({
    title: '删除创作项目',
    message: '确定要删除"' + app.name + '"吗？此操作不可恢复。',
    confirmText: '删除',
    type: 'danger'
  })
  if (confirmed) {
    try {
      await creationApi.deleteLocalCreation(app.folder_name || app.id)
      creationApps.value = creationApps.value.filter(a => a.id !== app.id)
      if (selectedCreationApp.value?.id === app.id) {
        selectedCreationApp.value = null
        clearCreationTempAutoSave()
      }
      notificationStore.showNotification('删除成功', 'success')
    } catch (e: any) {
      const msg = e?.message || String(e) || '未知错误'
      notificationStore.showNotification('删除失败: ' + msg, 'error')
    }
  }
}
function addSuggestedQuestion() {
  creationConfig.suggested_questions.push('')
}
function removeSuggestedQuestion(idx: number) {
  creationConfig.suggested_questions.splice(idx, 1)
}
function addShortcutCommand() {
  creationConfig.shortcut_commands.push({ label: '', command: '' })
}
function removeShortcutCommand(idx: number) {
  creationConfig.shortcut_commands.splice(idx, 1)
}
function addBannedWord() {
  const word = bannedWordInput.value.trim()
  if (word && !creationConfig.banned_words.includes(word)) {
    creationConfig.banned_words.push(word)
  }
  bannedWordInput.value = ''
}
function removeBannedWord(idx: number) {
  creationConfig.banned_words.splice(idx, 1)
}
function openImagePicker(target: string) {
  imagePickerTarget.value = target
  showImagePicker.value = true
}
function handleImagePickerConfirm(img: any) {
  if (!img) return
  const url = img.remoteUrl || ''
  if (!url) {
    notificationStore.showNotification('请选择已上传到远程的图片', 'error')
    return
  }
  if (imagePickerTarget.value === 'cover') {
    creationConfig.app.cover = url
  } else if (imagePickerTarget.value === 'bg_mobile') {
    creationConfig.bg_mobile = url
  } else if (imagePickerTarget.value === 'bg_image') {
    creationConfig.bg_image = url
  }
  showImagePicker.value = false
  markCreationUnsaved()
}
function clearCreationImage(type: string) {
  if (type === 'cover') {
    creationConfig.app.cover = ''
    creationConfig.app.cover_tiny = ''
  } else if (type === 'bg_mobile') {
    creationConfig.bg_mobile = ''
  } else if (type === 'bg_image') {
    creationConfig.bg_image = ''
  }
  markCreationUnsaved()
}
function switchToDescriptionTab() {
  creationEditorTab.value = 'description'
  creationWorldBookExpanded.value = false
}
function toggleWorldBookExpand() {
  if (creationEditorTab.value === 'description') return
  creationWorldBookExpanded.value = !creationWorldBookExpanded.value
}
function addCreationWorldBookCard() {
  const newEntry: WorldBookEntry = {
    id: 'wb_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9),
    key: '',
    keywords: [],
    isOrMode: false,
    value: '',
    keySystem: true,
    keyUser: true,
    keyAI: true,
    posGlobal: true,
    posPre: false,
    posPost: false,
    isFolder: false
  }
  creationConfig.world_book.push(newEntry)
  creationSelectedWorldBookEntry.value = newEntry
  creationWbKeywordInput.value = ''
  creationEditorTab.value = 'worldbook'
}
function selectCreationWorldBookEntry(entry: WorldBookEntry) {
  creationSelectedWorldBookEntry.value = entry
}
function removeCreationWorldBookEntry(entry: WorldBookEntry, idx: number) {
  if (entry.isFolder) {
    creationConfig.world_book = creationConfig.world_book.filter(e => e.id !== entry.id && e.parentId !== entry.id)
  } else {
    creationConfig.world_book.splice(idx, 1)
  }
  if (creationSelectedWorldBookEntry.value && creationSelectedWorldBookEntry.value.id === entry.id) {
    creationSelectedWorldBookEntry.value = null
  }
}
function getCreationWbKeywords(): string[] {
  if (!creationSelectedWorldBookEntry.value) return []
  return creationSelectedWorldBookEntry.value.keywords || []
}
function addCreationWbKeyword() {
  if (!creationSelectedWorldBookEntry.value || !creationWbKeywordInput.value.trim()) return
  if (!creationSelectedWorldBookEntry.value.keywords) {
    creationSelectedWorldBookEntry.value.keywords = []
  }
  creationSelectedWorldBookEntry.value.keywords.push(creationWbKeywordInput.value.trim())
  creationWbKeywordInput.value = ''
}
function removeCreationWbKeyword(idx: number) {
  if (!creationSelectedWorldBookEntry.value || !creationSelectedWorldBookEntry.value.keywords) return
  creationSelectedWorldBookEntry.value.keywords.splice(idx, 1)
}
function saveCreationDescription(content: string) {
  creationConfig.app.description = content
  markCreationUnsaved()
}
function toggleEditorBlocks() {
  editorShowBlocks.value = !editorShowBlocks.value
  if (editorShowBlocks.value) {
    editorShowStyles.value = false
  }
  if (grapesEditorRef.value) {
    (grapesEditorRef.value as any).showBlocks = editorShowBlocks.value;
    (grapesEditorRef.value as any).showStyles = editorShowStyles.value
  }
}
function toggleEditorStyles() {
  editorShowStyles.value = !editorShowStyles.value
  if (editorShowStyles.value) {
    editorShowBlocks.value = false
  }
  if (grapesEditorRef.value) {
    (grapesEditorRef.value as any).showStyles = editorShowStyles.value;
    (grapesEditorRef.value as any).showBlocks = editorShowBlocks.value
  }
}
function toggleEditorCode() {
  editorShowCode.value = !editorShowCode.value
  if (grapesEditorRef.value) {
    if (editorShowCode.value) {
      (grapesEditorRef.value as any).openCodeFullscreen()
    } else {
      (grapesEditorRef.value as any).closeCodeFullscreen()
    }
  }
}
function openTestDrawer() {
  showTestDrawer.value = true
  testMessages.value = []
  testMessageInput.value = ''
}
function closeTestDrawer() {
  showTestDrawer.value = false
}
function buildTestSystemPrompt(): string {
  const parts: string[] = []
  if (creationConfig.pre_prompt) parts.push(creationConfig.pre_prompt)
  if (creationConfig.pre_text) parts.push(creationConfig.pre_text)
  if (creationConfig.post_text) parts.push(creationConfig.post_text)
  return parts.join('\n\n')
}
async function sendTestMessage() {
  if (!testMessageInput.value.trim() || sendingTestMessage.value) return
  const userMessage = testMessageInput.value.trim()
  testMessageInput.value = ''
  testMessages.value.push({ role: 'user', content: userMessage })
  sendingTestMessage.value = true
  try {
    const request = {
      provider: selectedTestProvider.value,
      messages: testMessages.value.map(m => ({ role: m.role, content: m.content })),
      systemPrompt: buildTestSystemPrompt(),
      worldBook: creationConfig.world_book || []
    }
    const response = await creationApi.sendTestChat(request)
    if (response.error) {
      notificationStore.showNotification(response.error, 'error')
    } else {
      testMessages.value.push({ role: 'assistant', content: response.content })
    }
  } catch (e: any) {
    notificationStore.showNotification('发送失败: ' + e.message, 'error')
  } finally {
    sendingTestMessage.value = false
  }
}
onMounted(() => {
  loadCreationApps()
  document.addEventListener('keydown', handleSaveKeydown)
})
onBeforeUnmount(() => {
  clearCreationTempAutoSave()
  document.removeEventListener('keydown', handleSaveKeydown)
})
</script>
<style scoped>
.creation-view {
  display: flex;
  flex: 1;
  height: 100%;
  overflow: hidden;
}
.creation-sidebar {
  background: var(--bg-panel);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  flex-shrink: 0;
  transition: width 0.2s ease;
}
.creation-sidebar.collapsed {
  width: 40px !important;
}
.creation-header {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  gap: 8px;
}
.creation-header h2 {
  flex: 1;
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}
.collapse-icon {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}
.collapse-icon.expanded {
  transform: rotate(90deg);
}
.collapse-icon svg {
  width: 12px;
  height: 12px;
}
.btn-refresh {
  padding: 4px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-refresh:hover {
  background: var(--bg-tertiary);
  color: var(--primary);
}
.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.creation-sidebar-content {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.creation-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.creation-loading,
.creation-loading-center {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 20px;
  color: var(--text-secondary);
  font-size: 13px;
}
.creation-loading-center {
  flex: 1;
}
.creation-empty {
  padding: 20px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 13px;
}
.creation-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  gap: 10px;
  margin-bottom: 4px;
}
.creation-item:hover {
  background: var(--bg-tertiary);
}
.creation-item.active {
  background: var(--primary);
}
.creation-item-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  overflow: hidden;
  flex-shrink: 0;
}
.creation-item-icon img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.creation-item-icon-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  font-size: 16px;
  font-weight: 500;
}
.creation-item-info {
  flex: 1;
  overflow: hidden;
}
.creation-item-name {
  font-size: 13px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.creation-item-mode {
  font-size: 11px;
  color: var(--text-secondary);
  margin-top: 2px;
}
.btn-delete-creation {
  opacity: 0;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: transparent;
  color: var(--text-secondary);
  border: none;
  cursor: pointer;
  transition: opacity 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.creation-item:hover .btn-delete-creation {
  opacity: 1;
}
.btn-delete-creation:hover {
  background: #ef4444;
  color: white;
}
.creation-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-main);
  overflow: hidden;
  position: relative;
}
.creation-placeholder {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 14px;
}
.creation-editor-wrapper {
  flex: 1;
  display: flex;
  overflow: hidden;
}
.creation-main-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.creation-editor-body {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}
.creation-tab-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  padding: 16px;
}
.creation-description-area {
  padding: 0;
  overflow: hidden;
}
.form-section {
  margin-bottom: 0;
}
.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}
.form-input {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 13px;
  resize: none;
}
.form-input:focus {
  outline: none;
  border-color: var(--primary);
}
.form-textarea {
  width: 100%;
  padding: 10px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 8px;
  color: var(--text-primary);
  font-size: 13px;
  line-height: 1.5;
  resize: none;
}
.form-textarea:focus {
  outline: none;
  border-color: var(--primary);
}
.form-select {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 13px;
}
.form-row {
  display: flex;
  gap: 16px;
}
.form-section.third {
  flex: 1;
}
.image-columns {
  display: flex;
  gap: 16px;
}
.image-column {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: var(--bg-secondary);
  border-radius: 8px;
  padding: 12px;
}
.column-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
  text-align: center;
}
.column-preview {
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  background: var(--bg-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
}
.column-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.preview-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: var(--text-secondary);
  opacity: 0.4;
}
.preview-placeholder svg {
  width: 48px;
  height: 48px;
}
.column-buttons {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.btn-column {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 12px;
  font-size: 12px;
  border: none;
  border-radius: 16px;
  cursor: pointer;
  transition: all 0.15s;
  background: var(--primary);
  color: white;
}
.btn-column:hover:not(:disabled) {
  opacity: 0.9;
}
.btn-column.disabled,
.btn-column:disabled {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  opacity: 0.5;
  cursor: not-allowed;
}
.btn-column.btn-danger:not(.disabled):not(:disabled) {
  background: #ef4444;
}
.banned-words-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 10px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  min-height: 42px;
}
.banned-word-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: var(--bg-secondary);
  border-radius: 14px;
  font-size: 12px;
}
.tag-remove {
  width: 14px;
  height: 14px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 12px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}
.tag-remove:hover {
  color: #ef4444;
}
.banned-word-input {
  flex: 1;
  min-width: 100px;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: 12px;
  outline: none;
}
.list-items {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.list-item {
  display: flex;
  gap: 8px;
  align-items: center;
}
.list-item .form-input {
  flex: 1;
}
.shortcut-item .form-input {
  flex: 1;
}
.btn-remove-sm {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border: none;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.btn-remove-sm:hover {
  background: #ef4444;
  color: white;
}
.btn-add-item {
  padding: 8px 16px;
  background: var(--bg-tertiary);
  border: 1px dashed var(--border);
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 13px;
  text-align: center;
}
.btn-add-item:hover {
  color: var(--primary);
  border-color: var(--primary);
}
.creation-worldbook-tab {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.worldbook-entry-editor {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.entry-editor-header {
  margin-bottom: 16px;
}
.trigger-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.trigger-mode-switch {
  display: inline-flex;
  background: var(--bg-tertiary);
  border-radius: 16px;
  padding: 3px;
  gap: 2px;
}
.trigger-mode-switch .mode-option {
  padding: 4px 16px;
  font-size: 13px;
  cursor: pointer;
  border-radius: 14px;
  color: var(--text-secondary);
  transition: all 0.15s;
}
.trigger-mode-switch .mode-option.active {
  background: var(--primary);
  color: white;
}
.trigger-input-row {
  display: flex;
  gap: 8px;
}
.trigger-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 8px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  min-height: 36px;
}
.trigger-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: var(--bg-secondary);
  border-radius: 14px;
  font-size: 12px;
  color: var(--text-primary);
}
.entry-editor-controls {
  display: flex;
  gap: 24px;
  margin-bottom: 16px;
  padding: 12px;
  background: var(--bg-tertiary);
  border-radius: 8px;
}
.control-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.control-label {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
}
.region-toggles {
  display: flex;
  gap: 12px;
}
.region-toggle-switch {
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
}
.region-toggle-switch input {
  width: 16px;
  height: 16px;
}
.toggle-label {
  font-size: 12px;
  color: var(--text-primary);
}
.entry-editor-body {
  flex: 1;
  display: flex;
  flex-direction: column;
}
.entry-value-textarea {
  flex: 1;
  min-height: 300px;
}
.worldbook-empty-hint {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 14px;
}
.creation-right-panel {
  width: 180px;
  min-width: 180px;
  border-left: 1px solid var(--border);
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
}
.right-panel-tabs {
  display: flex;
  flex-direction: column;
}
.right-panel-tab {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  cursor: pointer;
  border-bottom: 1px solid var(--border);
  transition: background 0.15s;
}
.right-panel-tab:hover {
  background: var(--bg-tertiary);
}
.right-panel-tab.active {
  background: var(--primary);
  color: white;
}
.tab-label {
  font-size: 13px;
  font-weight: 500;
}
.right-panel-worldbook {
  border-bottom: 1px solid var(--border);
}
.worldbook-section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  cursor: pointer;
  transition: background 0.15s;
}
.worldbook-section-header:hover {
  background: var(--bg-tertiary);
}
.section-collapse-icon {
  width: 14px;
  height: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s;
}
.section-collapse-icon.expanded {
  transform: rotate(90deg);
}
.section-collapse-icon svg {
  width: 12px;
  height: 12px;
}
.section-title {
  font-size: 13px;
  font-weight: 500;
}
.section-count {
  font-size: 12px;
  color: var(--text-secondary);
}
.btn-add-wb-circle {
  margin-left: auto;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--primary);
  color: white;
  border: none;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.worldbook-section-content {
  max-height: 300px;
  overflow-y: auto;
  transition: max-height 0.2s;
}
.worldbook-section-content.collapsed {
  max-height: 0;
  overflow: hidden;
}
.worldbook-panel-list {
  padding: 0 8px 8px;
}
.worldbook-panel-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  margin-bottom: 2px;
  transition: background 0.15s;
}
.worldbook-panel-item:hover {
  background: var(--bg-tertiary);
}
.worldbook-panel-item.active {
  background: var(--primary);
  color: white;
}
.item-name {
  flex: 1;
  font-size: 12px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.btn-delete-wb-item {
  opacity: 0;
  width: 18px;
  height: 18px;
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 12px;
  transition: opacity 0.15s;
}
.worldbook-panel-item:hover .btn-delete-wb-item {
  opacity: 1;
}
.btn-delete-wb-item:hover {
  color: #ef4444;
}
.right-panel-divider {
  height: 1px;
  background: var(--border);
}
.right-panel-gjs-tools {
  flex: 1;
  overflow-y: auto;
}
.creation-save-status-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border);
  margin-top: auto;
}
.save-status-pill {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 12px;
}
.save-status-pill.saved {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}
.save-status-pill.saved svg {
  width: 14px;
  height: 14px;
}
.save-status-pill.unsaved {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}
.save-status-pill.clickable {
  cursor: pointer;
  transition: background 0.15s;
}
.save-status-pill.clickable:hover {
  background: var(--bg-hover);
}
.save-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.save-status-dot.green {
  background: #22c55e;
}
.save-status-dot.yellow {
  background: #eab308;
}
.save-status-dot.red {
  background: #ef4444;
}
.editor-tools-sidebar {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px 4px;
  gap: 6px;
}
.editor-tool-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
  background: var(--bg-tertiary);
  border-radius: 8px;
  padding: 4px;
}
.editor-tool-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
.editor-tool-btn:hover {
  background: var(--hover-bg);
  color: var(--text-color);
}
.editor-tool-btn.active {
  background: var(--primary-light);
  color: var(--primary);
}
.editor-tool-btn svg {
  width: 18px;
  height: 18px;
}
.editor-tool-divider {
  width: 24px;
  height: 1px;
  background: var(--border);
  margin: 4px 0;
}
.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid var(--border);
  border-top-color: var(--primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.btn-test-chat {
  width: 100%;
  margin-bottom: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px 16px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 13px;
  transition: opacity 0.15s;
}
.btn-test-chat:hover {
  opacity: 0.9;
}
.test-drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 1000;
  display: flex;
  justify-content: flex-end;
}
.test-drawer {
  width: 450px;
  max-width: 100%;
  height: 100%;
  background: var(--bg-panel);
  display: flex;
  flex-direction: column;
  box-shadow: -4px 0 24px rgba(0, 0, 0, 0.2);
}
.test-drawer-header {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}
.test-drawer-title {
  font-size: 16px;
  font-weight: 600;
}
.test-provider-select {
  padding: 6px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 16px;
  color: var(--text-primary);
  font-size: 13px;
  cursor: pointer;
}
.btn-close-drawer {
  margin-left: auto;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--bg-tertiary);
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.btn-close-drawer:hover {
  background: var(--bg-secondary);
  color: var(--text-primary);
}
.test-drawer-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.test-preview-area {
  max-height: 200px;
  overflow-y: auto;
  border-bottom: 1px solid var(--border);
  padding: 12px 16px;
}
.preview-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
}
.preview-content {
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-primary);
}
.test-messages-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.test-messages-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
  font-size: 14px;
}
.test-message {
  display: flex;
}
.test-message.user {
  justify-content: flex-end;
}
.test-message.assistant {
  justify-content: flex-start;
}
.message-bubble {
  max-width: 85%;
  padding: 10px 14px;
  border-radius: 16px;
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
  white-space: pre-wrap;
}
.test-message.user .message-bubble {
  background: var(--primary);
  color: white;
  border-bottom-right-radius: 4px;
}
.test-message.assistant .message-bubble {
  background: var(--bg-tertiary);
  color: var(--text-primary);
  border-bottom-left-radius: 4px;
}
.typing-indicator {
  display: flex;
  gap: 4px;
}
.typing-indicator span {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-secondary);
  animation: typing-bounce 1.2s infinite;
}
.typing-indicator span:nth-child(2) {
  animation-delay: 0.2s;
}
.typing-indicator span:nth-child(3) {
  animation-delay: 0.4s;
}
@keyframes typing-bounce {
  0%, 60%, 100% { transform: translateY(0); }
  30% { transform: translateY(-4px); }
}
.test-drawer-footer {
  display: flex;
  gap: 10px;
  padding: 16px;
  border-top: 1px solid var(--border);
}
.test-input {
  flex: 1;
  padding: 10px 16px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border);
  border-radius: 20px;
  color: var(--text-primary);
  font-size: 13px;
  resize: none;
}
.test-input:focus {
  outline: none;
  border-color: var(--primary);
}
.btn-send-test {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: var(--primary);
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.btn-send-test:hover:not(:disabled) {
  opacity: 0.9;
}
.btn-send-test:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: all 0.3s ease;
}
.drawer-slide-enter-active .test-drawer,
.drawer-slide-leave-active .test-drawer {
  transition: transform 0.3s ease;
}
.drawer-slide-enter-from,
.drawer-slide-leave-to {
  background: rgba(0, 0, 0, 0);
}
.drawer-slide-enter-from .test-drawer,
.drawer-slide-leave-to .test-drawer {
  transform: translateX(100%);
}
</style>
