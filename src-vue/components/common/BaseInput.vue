<template>
  <div :class="['base-input', { 'base-input--error': error, 'base-input--disabled': disabled }]">
    <label v-if="label" class="base-input__label">{{ label }}</label>
    <div class="base-input__wrapper">
      <span v-if="prefix" class="base-input__prefix">{{ prefix }}</span>
      <input
        ref="inputRef"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        class="base-input__field"
        @input="handleInput"
        @focus="emit('focus', $event)"
        @blur="emit('blur', $event)"
        @keydown.enter="emit('enter', $event)"
      />
      <span v-if="suffix" class="base-input__suffix">{{ suffix }}</span>
      <button v-if="clearable && modelValue" class="base-input__clear" @click="handleClear">×</button>
    </div>
    <p v-if="error" class="base-input__error">{{ error }}</p>
    <p v-else-if="hint" class="base-input__hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
const props = defineProps<{
  modelValue?: string | number
  type?: string
  placeholder?: string
  label?: string
  hint?: string
  error?: string
  prefix?: string
  suffix?: string
  disabled?: boolean
  readonly?: boolean
  clearable?: boolean
}>()
const emit = defineEmits<{
  'update:modelValue': [value: string]
  focus: [e: FocusEvent]
  blur: [e: FocusEvent]
  enter: [e: KeyboardEvent]
  clear: []
}>()
const inputRef = ref<HTMLInputElement | null>(null)
function handleInput(e: Event) {
  emit('update:modelValue', (e.target as HTMLInputElement).value)
}
function handleClear() {
  emit('update:modelValue', '')
  emit('clear')
  inputRef.value?.focus()
}
function focus() {
  inputRef.value?.focus()
}
function select() {
  inputRef.value?.select()
}
defineExpose({ inputRef, focus, select })
</script>

<style scoped>
.base-input {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.base-input__label {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}
.base-input__wrapper {
  display: flex;
  align-items: center;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 0 12px;
  transition: border-color 0.2s;
}
.base-input__wrapper:focus-within {
  border-color: var(--primary-color);
}
.base-input--error .base-input__wrapper {
  border-color: #ef4444;
}
.base-input--disabled .base-input__wrapper {
  opacity: 0.6;
  cursor: not-allowed;
}
.base-input__field {
  flex: 1;
  border: none;
  background: transparent;
  padding: 10px 0;
  font-size: 14px;
  color: var(--text-primary);
  outline: none;
  min-width: 0;
}
.base-input__field::placeholder {
  color: var(--text-tertiary);
}
.base-input__prefix,
.base-input__suffix {
  color: var(--text-tertiary);
  font-size: 14px;
}
.base-input__prefix {
  margin-right: 8px;
}
.base-input__suffix {
  margin-left: 8px;
}
.base-input__clear {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  margin-left: 4px;
  border: none;
  background: var(--bg-tertiary);
  border-radius: 50%;
  color: var(--text-tertiary);
  cursor: pointer;
  font-size: 14px;
  line-height: 1;
}
.base-input__clear:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.base-input__error {
  font-size: 12px;
  color: #ef4444;
  margin: 0;
}
.base-input__hint {
  font-size: 12px;
  color: var(--text-tertiary);
  margin: 0;
}
</style>
