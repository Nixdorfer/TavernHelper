<template>
  <div class="custom-select" :class="{ open: isOpen }" ref="selectRef">
    <div class="select-trigger" @click="toggle">
      <span class="select-value">{{ currentLabel }}</span>
      <span class="select-arrow">▾</span>
    </div>
    <transition name="dropdown">
      <div v-if="isOpen" class="select-dropdown">
        <div
          v-for="option in options"
          :key="option.value"
          class="select-option"
          :class="{ selected: option.value === modelValue }"
          @click="selectOption(option)"
        >
          {{ option.label }}
        </div>
      </div>
    </transition>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
interface SelectOption {
  value: string | number
  label: string
}
const props = withDefaults(defineProps<{
  modelValue?: string | number
  options: SelectOption[]
  placeholder?: string
}>(), {
  modelValue: '',
  placeholder: '请选择'
})
const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  change: [value: string | number]
}>()
const isOpen = ref(false)
const selectRef = ref<HTMLElement | null>(null)
const currentLabel = computed(() => {
  const option = props.options.find(o => o.value === props.modelValue)
  return option ? option.label : props.placeholder
})
function toggle() {
  isOpen.value = !isOpen.value
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
onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
<style scoped>
.custom-select {
  position: relative;
  width: 100%;
}
.select-trigger {
  width: 100%;
  padding: 10px 12px;
  background: var(--input-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  color: var(--text-color);
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  transition: border-color 0.2s;
}
.select-trigger:hover {
  border-color: var(--primary-color);
}
.custom-select.open .select-trigger {
  border-color: var(--primary-color);
}
.select-arrow {
  font-size: 12px;
  color: var(--text-secondary);
  transition: transform 0.2s;
}
.custom-select.open .select-arrow {
  transform: rotate(180deg);
}
.select-dropdown {
  position: absolute;
  top: calc(100% + 4px);
  left: 0;
  right: 0;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  z-index: 100;
  overflow: hidden;
  transform-origin: top center;
}
.select-option {
  padding: 10px 12px;
  font-size: 14px;
  color: var(--text-color);
  cursor: pointer;
  transition: background 0.15s;
}
.select-option:hover {
  background: var(--hover-bg);
}
.select-option.selected {
  background: var(--primary-light);
  color: var(--primary-color);
}
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scaleY(0.8) translateY(-4px);
}
.dropdown-enter-to,
.dropdown-leave-from {
  opacity: 1;
  transform: scaleY(1) translateY(0);
}
</style>
