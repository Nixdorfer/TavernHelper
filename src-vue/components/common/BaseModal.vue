<template>
  <Teleport to="body">
    <Transition name="modal" @after-enter="handleAfterEnter">
      <div v-if="modelValue" class="base-modal" @click.self="handleBackdropClick">
        <div ref="contentRef" :class="['base-modal__content', `base-modal--${size}`]" @click.stop>
          <div v-if="title || $slots.header || closable" class="base-modal__header">
            <slot name="header">
              <h3 class="base-modal__title">{{ title }}</h3>
            </slot>
            <button v-if="closable" class="base-modal__close" @click="close">×</button>
          </div>
          <div class="base-modal__body">
            <slot />
          </div>
          <div v-if="$slots.footer" class="base-modal__footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
const props = withDefaults(defineProps<{
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closable?: boolean
  closeOnBackdrop?: boolean
  closeOnEscape?: boolean
  confirmOnEnter?: boolean
}>(), {
  size: 'md',
  closable: true,
  closeOnBackdrop: true,
  closeOnEscape: true,
  confirmOnEnter: true
})
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
  confirm: []
}>()
const contentRef = ref<HTMLElement | null>(null)
function close() {
  emit('update:modelValue', false)
  emit('close')
}
function handleBackdropClick() {
  if (props.closeOnBackdrop) {
    close()
  }
}
function handleAfterEnter() {
  if (!contentRef.value) return
  const firstInput = contentRef.value.querySelector('input:not([type="hidden"]):not([disabled]), textarea:not([disabled])') as HTMLElement
  if (firstInput) {
    firstInput.focus()
  }
}
function handleKeydown(e: KeyboardEvent) {
  if (!props.modelValue) return
  if (e.key === 'Escape' && props.closeOnEscape) {
    close()
  }
  if (e.key === 'Enter' && props.confirmOnEnter) {
    const target = e.target as HTMLElement
    if (target.tagName === 'TEXTAREA') return
    e.preventDefault()
    emit('confirm')
  }
}
watch(() => props.modelValue, (val) => {
  document.body.style.overflow = val ? 'hidden' : ''
})
onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.base-modal {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
  padding: 20px;
}
.base-modal__content {
  background: var(--bg-primary);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 40px);
  overflow: hidden;
}
.base-modal--sm { width: 320px; }
.base-modal--md { width: 480px; }
.base-modal--lg { width: 640px; }
.base-modal--xl { width: 800px; }
.base-modal--full { width: calc(100vw - 40px); height: calc(100vh - 40px); }
.base-modal__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}
.base-modal__title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}
.base-modal__close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 50%;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 24px;
  line-height: 1;
  transition: all 0.2s;
}
.base-modal__close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.base-modal__body {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}
.base-modal__footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-active .base-modal__content,
.modal-leave-active .base-modal__content {
  transition: transform 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .base-modal__content,
.modal-leave-to .base-modal__content {
  transform: scale(0.95);
}
</style>
