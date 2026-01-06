<template>
  <button
    :class="[
      'base-button',
      `base-button--${variant}`,
      `base-button--${size}`,
      { 'base-button--loading': loading, 'base-button--disabled': disabled }
    ]"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <span v-if="loading" class="base-button__spinner"></span>
    <span v-else-if="icon" class="base-button__icon">{{ icon }}</span>
    <span v-if="$slots.default" class="base-button__text"><slot /></span>
  </button>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  variant?: 'primary' | 'secondary' | 'danger' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  icon?: string
}>(), {
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false
})
const emit = defineEmits<{
  click: [e: MouseEvent]
}>()
function handleClick(e: MouseEvent) {
  emit('click', e)
}
</script>

<style scoped>
.base-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border: none;
  border-radius: 999px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
  white-space: nowrap;
}
.base-button--primary {
  background: var(--primary-color);
  color: white;
}
.base-button--primary:hover:not(:disabled) {
  filter: brightness(1.1);
}
.base-button--secondary {
  background: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}
.base-button--secondary:hover:not(:disabled) {
  background: var(--bg-hover);
}
.base-button--danger {
  background: #ef4444;
  color: white;
}
.base-button--danger:hover:not(:disabled) {
  background: #dc2626;
}
.base-button--ghost {
  background: transparent;
  color: var(--text-secondary);
}
.base-button--ghost:hover:not(:disabled) {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.base-button--sm {
  padding: 4px 12px;
  font-size: 12px;
}
.base-button--md {
  padding: 8px 16px;
  font-size: 14px;
}
.base-button--lg {
  padding: 12px 24px;
  font-size: 16px;
}
.base-button--disabled,
.base-button--loading {
  opacity: 0.6;
  cursor: not-allowed;
}
.base-button__spinner {
  width: 14px;
  height: 14px;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
