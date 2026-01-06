<template>
  <div class="plaza-content-full">
    <div v-if="!currentTab" class="plaza-welcome">
      <div class="plaza-welcome-title">欢迎来到广场</div>
      <div class="plaza-welcome-subtitle">在这里您可以自由探索新鲜事物</div>
    </div>
    <div v-else-if="loading" class="plaza-loading-center">
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>
    <div v-else-if="apps.length === 0" class="plaza-empty-center">
      <span>暂无内容</span>
    </div>
    <div v-else :class="['plaza-cards-grid', { 'no-image-mode': noImageMode }]">
      <div
        v-for="app in apps"
        :key="app.id"
        class="plaza-card"
        @click="handleAppClick(app)"
      >
        <div v-if="!noImageMode" class="plaza-card-cover">
          <img v-if="app.icon" :src="app.icon" />
          <div v-else class="plaza-card-cover-placeholder">
            <span>{{ app.name?.charAt(0) || '?' }}</span>
          </div>
        </div>
        <div class="plaza-card-info">
          <div class="plaza-card-title">{{ app.name }}</div>
          <div v-if="app.isAuthor" class="plaza-card-meta">
            <span>{{ app.appCount }} 个作品</span>
            <span>{{ app.followerCount }} 粉丝</span>
          </div>
          <div v-else class="plaza-card-meta">
            <span v-if="app.account_name" class="card-author">{{ app.account_name }}</span>
            <span v-if="app.players_count" class="card-stat">{{ formatCount(app.players_count) }}人玩过</span>
          </div>
        </div>
      </div>
    </div>
    <AppDetailModal
      v-if="showAppDetail"
      :app="selectedApp"
      :loading="appDetailLoading"
      @close="closeAppDetail"
      @start-conversation="startConversation"
    />
  </div>
</template>
<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useAuthStore, useConfigStore, useAppStore, useNotificationStore } from '@/stores'
import { useConversationStore } from '@/stores/modules/conversation'
import AppDetailModal from '@/components/modals/AppDetailModal.vue'
interface PlazaApp {
  id: string
  name: string
  icon?: string
  summary?: string
  account_name?: string
  like_count?: number
  players_count?: number
  appCount?: number
  followerCount?: number
  isAuthor: boolean
}
const emit = defineEmits<{
  'switch-to-conversation': []
}>()
const authStore = useAuthStore()
const configStore = useConfigStore()
const appStore = useAppStore()
const notificationStore = useNotificationStore()
const conversationStore = useConversationStore()
const apps = ref<PlazaApp[]>([])
const loading = ref(false)
const showAppDetail = ref(false)
const selectedApp = ref<PlazaApp | null>(null)
const appDetailLoading = ref(false)
const token = computed(() => authStore.token)
const currentTab = computed(() => appStore.plazaCurrentTab)
const noImageMode = computed(() => configStore.noImageMode)
const formatCount = (num: number): string => {
  if (!num) return ''
  if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return num.toString()
}
const loadApps = async () => {
  loading.value = true
  try {
    const tabKey = currentTab.value
    let result: PlazaApp[] = []
    if (tabKey === 'starauthor') {
      const url = `https://aipornhub.ltd/console/api/account/author-ranking?page=1&page_size=30&sort_by=total_overall_rank&sort_order=desc`
      const response = await (window as any).go.main.App.FetchWithAuth(token.value, url, 'GET', '')
      const data = JSON.parse(response)
      const items = data.data || []
      if (Array.isArray(items)) {
        result = items.map((item: any) => ({
          id: item.account_id,
          name: item.name,
          icon: item.avatar,
          appCount: item.app_count,
          followerCount: item.follower_count,
          isAuthor: true
        }))
      }
    } else {
      const url = `https://aipornhub.ltd/go/api/explore/search?keywords=&ranking=${tabKey}&page=1&limit=30&order=default&lang=zh-Hans`
      const response = await (window as any).go.main.App.FetchWithAuth(token.value, url, 'GET', '')
      const data = JSON.parse(response)
      const items = data.data?.apps || data.data?.items || data.items || []
      if (Array.isArray(items)) {
        result = items.map((item: any) => ({
          id: item.id,
          name: item.name,
          icon: item.cover || item.icon,
          summary: item.summary,
          account_name: item.account_name,
          like_count: item.like_count,
          players_count: item.players_count,
          isAuthor: false
        }))
      }
    }
    apps.value = result
  } catch (e) {
    console.error('加载广场列表失败:', e)
    apps.value = []
  }
  loading.value = false
}
const handleAppClick = (app: PlazaApp) => {
  if (app.isAuthor) return
  selectedApp.value = app
  showAppDetail.value = true
}
const closeAppDetail = () => {
  showAppDetail.value = false
  selectedApp.value = null
}
const startConversation = async () => {
  if (!selectedApp.value) return
  const app = selectedApp.value
  conversationStore.newConversationStartPage = {
    appId: app.id,
    appName: app.name,
    html: `<div style="text-align:center;padding:40px;color:#666;"><h2>${app.name}</h2><p>${app.summary || ''}</p></div>`
  }
  closeAppDetail()
  emit('switch-to-conversation')
}
watch(currentTab, (val) => {
  if (val) {
    loadApps()
  } else {
    apps.value = []
  }
}, { immediate: true })
</script>
<style scoped>
.plaza-content-full {
  flex: 1;
  background: var(--bg-primary);
  overflow-y: auto;
  padding: 20px;
}
.plaza-welcome {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 16px;
}
.plaza-welcome-title {
  font-size: 36px;
  font-weight: 600;
  color: var(--text-primary);
}
.plaza-welcome-subtitle {
  font-size: 16px;
  color: var(--text-secondary);
}
.plaza-loading-center,
.plaza-empty-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 12px;
  color: var(--text-secondary);
}
.plaza-cards-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
}
.plaza-cards-grid.no-image-mode {
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
}
.plaza-card {
  background: var(--bg-secondary);
  border-radius: 12px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid var(--border-color);
}
.plaza-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.3);
}
.plaza-card-cover {
  width: 100%;
  aspect-ratio: 3 / 4;
  background: var(--bg-tertiary);
  overflow: hidden;
}
.plaza-card-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.plaza-card-cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  font-weight: bold;
  color: var(--text-secondary);
  background: linear-gradient(135deg, var(--bg-tertiary), var(--bg-secondary));
}
.plaza-card-info {
  padding: 12px;
}
.plaza-card-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-bottom: 6px;
}
.plaza-card-meta {
  display: flex;
  gap: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  flex-wrap: wrap;
}
.plaza-card-meta .card-author {
  color: var(--primary-color);
}
.plaza-card-meta .card-stat {
  color: var(--text-tertiary);
}
.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
