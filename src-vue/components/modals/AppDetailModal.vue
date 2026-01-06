<template>
  <Teleport to="body">
    <div class="modal-overlay" @click.self="emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <h3>{{ app?.name }}</h3>
          <button class="btn-close" @click="emit('close')">×</button>
        </div>
        <div class="modal-body">
          <div v-if="loading" class="loading-center">
            <span class="loading-spinner"></span>
          </div>
          <template v-else>
            <div v-if="app?.icon && !noImageMode" class="app-cover">
              <img :src="app.icon" />
            </div>
            <div class="app-info">
              <div v-if="app?.account_name" class="app-author">作者: {{ app.account_name }}</div>
              <div v-if="app?.summary" class="app-summary">{{ app.summary }}</div>
              <div class="app-stats">
                <span v-if="app?.players_count">{{ formatCount(app.players_count) }} 人玩过</span>
                <span v-if="app?.like_count">{{ formatCount(app.like_count) }} 收藏</span>
              </div>
            </div>
          </template>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="emit('close')">关闭</button>
          <button class="btn btn-primary" @click="emit('start-conversation')">开始对话</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { useConfigStore } from '@/stores'
interface PlazaApp {
  id: string
  name: string
  icon?: string
  summary?: string
  account_name?: string
  like_count?: number
  players_count?: number
}
const props = defineProps<{
  app: PlazaApp | null
  loading?: boolean
}>()
const emit = defineEmits<{
  close: []
  'start-conversation': []
}>()
const configStore = useConfigStore()
const noImageMode = computed(() => configStore.noImageMode)
const formatCount = (num: number): string => {
  if (!num) return ''
  if (num >= 10000) return (num / 10000).toFixed(1) + 'w'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'k'
  return num.toString()
}
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    emit('close')
  }
  if (e.key === 'Enter') {
    e.preventDefault()
    emit('start-conversation')
  }
}
onMounted(() => window.addEventListener('keydown', handleKeydown))
onUnmounted(() => window.removeEventListener('keydown', handleKeydown))
</script>
<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}
.modal-container {
  background: var(--bg-primary);
  border-radius: 16px;
  width: 90%;
  max-width: 480px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}
.modal-header h3 {
  margin: 0;
  font-size: 18px;
  color: var(--text-primary);
}
.btn-close {
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  font-size: 24px;
  cursor: pointer;
  color: var(--text-secondary);
  border-radius: 8px;
}
.btn-close:hover {
  background: var(--bg-hover);
}
.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}
.loading-center {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
}
.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.app-cover {
  width: 100%;
  aspect-ratio: 16 / 9;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 16px;
}
.app-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.app-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.app-author {
  font-size: 14px;
  color: var(--primary-color);
}
.app-summary {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}
.app-stats {
  display: flex;
  gap: 16px;
  font-size: 13px;
  color: var(--text-tertiary);
}
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}
.btn {
  padding: 10px 24px;
  border-radius: 8px;
  font-size: 14px;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}
.btn-secondary {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.btn-secondary:hover {
  background: var(--border-color);
}
.btn-primary {
  background: var(--primary-color);
  color: #fff;
}
.btn-primary:hover {
  background: var(--primary-hover);
}
</style>
