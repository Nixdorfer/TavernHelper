<template>
  <Teleport to="body">
    <Transition name="context-menu">
      <div
        v-if="visible"
        class="context-menu"
        :style="{ left: position.x + 'px', top: position.y + 'px' }"
        @click.stop
      >
        <div
          v-for="(item, index) in items"
          :key="index"
          class="context-menu-item"
          :class="{ danger: item.danger }"
          @click="handleClick(item)"
        >
          <span v-if="item.icon" class="context-menu-icon" v-html="item.icon"></span>
          <span>{{ item.label }}</span>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
export interface ContextMenuItem {
  label: string
  icon?: string
  danger?: boolean
  action: () => void
}
const props = defineProps<{
  items: ContextMenuItem[]
}>()
const visible = ref(false)
const position = ref({ x: 0, y: 0 })
function show(x: number, y: number) {
  const menuWidth = 140
  const menuHeight = props.items.length * 36 + 8
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  let finalX = x
  let finalY = y
  if (x + menuWidth > viewportWidth) {
    finalX = viewportWidth - menuWidth - 8
  }
  if (y + menuHeight > viewportHeight) {
    finalY = viewportHeight - menuHeight - 8
  }
  position.value = { x: finalX, y: finalY }
  visible.value = true
}
function hide() {
  visible.value = false
}
function handleClick(item: ContextMenuItem) {
  item.action()
  hide()
}
function handleOutsideClick() {
  hide()
}
onMounted(() => {
  document.addEventListener('click', handleOutsideClick)
  document.addEventListener('contextmenu', handleOutsideClick)
})
onBeforeUnmount(() => {
  document.removeEventListener('click', handleOutsideClick)
  document.removeEventListener('contextmenu', handleOutsideClick)
})
defineExpose({ show, hide })
</script>
<style scoped>
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 4px;
  min-width: 120px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
  z-index: 9999;
}
.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-primary);
  transition: background 0.15s;
}
.context-menu-item:hover {
  background: var(--bg-tertiary);
}
.context-menu-item.danger {
  color: #ef4444;
}
.context-menu-item.danger:hover {
  background: rgba(239, 68, 68, 0.15);
}
.context-menu-icon {
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.context-menu-icon :deep(svg) {
  width: 16px;
  height: 16px;
}
.context-menu-enter-active,
.context-menu-leave-active {
  transition: opacity 0.15s, transform 0.15s;
}
.context-menu-enter-from,
.context-menu-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
