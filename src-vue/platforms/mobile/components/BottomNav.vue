<template>
  <nav class="bottom-nav">
    <button
      v-for="item in navItems"
      :key="item.id"
      :class="['bottom-nav__item', { 'bottom-nav__item--active': currentView === item.id }]"
      @click="emit('navigate', item.id)"
    >
      <span class="bottom-nav__icon">{{ item.icon }}</span>
      <span class="bottom-nav__label">{{ item.label }}</span>
    </button>
  </nav>
</template>

<script setup lang="ts">
import type { ViewType } from '@/stores/modules/app'
defineProps<{
  currentView: ViewType
}>()
const emit = defineEmits<{
  navigate: [view: ViewType]
}>()
const navItems = [
  { id: 'conversation', label: '对话', icon: '💬' },
  { id: 'worldbook', label: '项目', icon: '📖' },
  { id: 'plaza', label: '广场', icon: '🏛️' },
  { id: 'gallery', label: '图库', icon: '🖼️' }
]
</script>

<style scoped>
.bottom-nav {
  display: flex;
  align-items: center;
  justify-content: space-around;
  height: 56px;
  background: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  padding-bottom: env(safe-area-inset-bottom);
}
.bottom-nav__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 6px 16px;
  border: none;
  background: transparent;
  cursor: pointer;
  transition: all var(--transition-fast);
}
.bottom-nav__item--active {
  color: var(--primary-color);
}
.bottom-nav__icon {
  font-size: 20px;
}
.bottom-nav__label {
  font-size: 11px;
  color: var(--text-secondary);
}
.bottom-nav__item--active .bottom-nav__label {
  color: var(--primary-color);
}
</style>
