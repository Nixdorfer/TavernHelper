<template>
  <div class="markdown-editor">
    <textarea
      ref="textareaRef"
      :value="modelValue"
      @input="handleInput"
      :placeholder="placeholder"
      :rows="rows"
      class="editor-textarea"
    ></textarea>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
const props = withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
  rows?: number
}>(), {
  placeholder: '',
  rows: 6
})
const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
const textareaRef = ref<HTMLTextAreaElement | null>(null)
function handleInput(e: Event) {
  const target = e.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
}
</script>
<style scoped>
.markdown-editor {
  width: 100%;
}
.editor-textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  line-height: 1.6;
  resize: vertical;
  font-family: inherit;
}
.editor-textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}
</style>
