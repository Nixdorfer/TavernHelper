<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="visible" class="modal-overlay" @mousedown.self="close">
        <div class="image-preview-modal">
          <div class="preview-header">
            <h3>图片预览</h3>
            <button class="btn-close" @click="close">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>
          <div class="preview-body">
            <img v-if="imageUrl" :src="imageUrl" alt="预览图片" />
            <div v-else class="preview-empty">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="18" height="18" rx="2"/>
                <circle cx="8.5" cy="8.5" r="1.5"/>
                <path d="M21 15l-5-5L5 21"/>
              </svg>
              <span>暂无图片</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
<script setup lang="ts">
import { watch, onMounted, onUnmounted } from 'vue'
const props = withDefaults(defineProps<{
  visible?: boolean
  imageUrl?: string
}>(), {
  visible: false,
  imageUrl: ''
})
const emit = defineEmits<{
  'update:visible': [value: boolean]
}>()
function close() {
  emit('update:visible', false)
}
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && props.visible) {
    close()
  }
}
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
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.image-preview-modal {
  position: relative;
  max-width: 90vw;
  max-height: 90vh;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.preview-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}
.preview-header h3 {
  margin: 0;
  font-size: 14px;
  font-weight: 500;
}
.btn-close {
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
  color: var(--text-secondary);
  transition: all 0.2s;
}
.btn-close:hover {
  background: var(--hover-bg);
  color: var(--text-color);
}
.btn-close svg {
  width: 16px;
  height: 16px;
}
.preview-body {
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
}
.preview-body img {
  max-width: 80vw;
  max-height: 70vh;
  object-fit: contain;
  border-radius: 8px;
}
.preview-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 48px;
  color: var(--text-secondary);
}
.preview-empty svg {
  width: 64px;
  height: 64px;
  opacity: 0.5;
}
.modal-enter-active {
  transition: opacity 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-leave-active {
  transition: opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-enter-active .image-preview-modal {
  transition: transform 0.25s cubic-bezier(0.34, 1.56, 0.64, 1),
              opacity 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-leave-active .image-preview-modal {
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1),
              opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .image-preview-modal {
  transform: scale(0.9) translateY(10px);
  opacity: 0;
}
.modal-leave-to .image-preview-modal {
  transform: scale(0.95);
  opacity: 0;
}
</style>
