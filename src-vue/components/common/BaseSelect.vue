<template>
  <div :class="['base-select', { 'base-select--open': isOpen, 'base-select--disabled': disabled }]" ref="selectRef">
    <label v-if="label" class="base-select__label">{{ label }}</label>
    <div class="base-select__trigger" @click="toggle">
      <span class="base-select__value">{{ displayValue }}</span>
      <span class="base-select__arrow">▼</span>
    </div>
    <Transition name="dropdown">
      <div v-if="isOpen" class="base-select__dropdown">
        <div
          v-for="option in options"
          :key="option.value"
          :class="['base-select__option', { 'base-select__option--selected': option.value === modelValue }]"
          @click="selectOption(option)"
        >
          {{ option.label }}
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
export interface SelectOption {
  label: string
  value: string | number
}
const props = defineProps<{
  modelValue?: string | number
  options: SelectOption[]
  label?: string
  placeholder?: string
  disabled?: boolean
}>()
const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  change: [value: string | number]
}>()
const isOpen = ref(false)
const selectRef = ref<HTMLElement | null>(null)
const displayValue = computed(() => {
  const selected = props.options.find(o => o.value === props.modelValue)
  return selected?.label || props.placeholder || '请选择'
})
function toggle() {
  if (!props.disabled) {
    isOpen.value = !isOpen.value
  }
}
function selectOption(option: SelectOption) {
  emit('update:modelValue', option.value)
  emit('change', option.value)
  isOpen.value = false
}
function handleClickOutside(e: MouseEvent) {
  if (selectRef.value && !selectRef.value.contains(e.target as Node)) {
    isOpen.value = false
  }
}
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})
onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.base-select {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.base-select__label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}
.base-select__trigger {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  cursor: pointer;
  transition: border-color 0.2s;
}
.base-select--open .base-select__trigger,
.base-select__trigger:hover {
  border-color: var(--primary-color);
}
.base-select--disabled .base-select__trigger {
  opacity: 0.6;
  cursor: not-allowed;
}
.base-select__value {
  font-size: 14px;
  color: var(--text-primary);
}
.base-select__arrow {
  font-size: 10px;
  color: var(--text-tertiary);
  transition: transform 0.2s;
}
.base-select--open .base-select__arrow {
  transform: rotate(180deg);
}
.base-select__dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 4px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  max-height: 200px;
  overflow-y: auto;
}
.base-select__option {
  padding: 10px 12px;
  font-size: 14px;
  color: var(--text-primary);
  cursor: pointer;
  transition: background 0.15s;
}
.base-select__option:hover {
  background: var(--bg-hover);
}
.base-select__option--selected {
  color: var(--primary-color);
  font-weight: 500;
}
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.15s ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
