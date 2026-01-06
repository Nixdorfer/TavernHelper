<template>
  <div class="code-editor">
    <div class="code-editor-container" ref="containerRef">
      <pre class="code-highlight" :class="language" aria-hidden="true"><code v-html="highlightedCode"></code></pre>
      <textarea
        ref="textareaRef"
        class="code-input"
        :value="modelValue"
        @input="handleInput"
        @scroll="syncScroll"
        @keydown="handleKeydown"
        spellcheck="false"
        autocomplete="off"
        autocorrect="off"
        autocapitalize="off"
      ></textarea>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
const props = withDefaults(defineProps<{
  modelValue: string
  language: 'html' | 'css'
}>(), {
  modelValue: '',
  language: 'html'
})
const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()
const containerRef = ref<HTMLElement | null>(null)
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const highlightedCode = computed(() => {
  const code = props.modelValue || ''
  if (props.language === 'html') {
    return highlightHtml(code)
  } else {
    return highlightCss(code)
  }
})
function highlightHtml(code: string): string {
  let result = escapeHtml(code)
  result = result.replace(/(&lt;\/?)([\w-]+)/g, '<span class="hl-tag">$1$2</span>')
  result = result.replace(/([\w-]+)(=)(&quot;[^&]*&quot;)/g, '<span class="hl-attr">$1</span>$2<span class="hl-string">$3</span>')
  result = result.replace(/(&lt;!--[\s\S]*?--&gt;)/g, '<span class="hl-comment">$1</span>')
  return result + '\n'
}
function highlightCss(code: string): string {
  let result = escapeHtml(code)
  result = result.replace(/(\/\*[\s\S]*?\*\/)/g, '<span class="hl-comment">$1</span>')
  result = result.replace(/([.#]?[\w-]+)(\s*\{)/g, '<span class="hl-selector">$1</span>$2')
  result = result.replace(/([\w-]+)(\s*:)/g, '<span class="hl-property">$1</span>$2')
  result = result.replace(/:\s*([^;{}]+)(;)/g, ': <span class="hl-value">$1</span>$2')
  result = result.replace(/(\d+)(px|em|rem|%|vh|vw|s|ms)/g, '<span class="hl-number">$1$2</span>')
  result = result.replace(/(#[0-9a-fA-F]{3,8})/g, '<span class="hl-color">$1</span>')
  return result + '\n'
}
function escapeHtml(str: string): string {
  return str
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}
function handleInput(e: Event) {
  const target = e.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
}
function syncScroll() {
  if (!textareaRef.value || !containerRef.value) return
  const pre = containerRef.value.querySelector('pre')
  if (pre) {
    pre.scrollTop = textareaRef.value.scrollTop
    pre.scrollLeft = textareaRef.value.scrollLeft
  }
}
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Tab') {
    e.preventDefault()
    const textarea = textareaRef.value
    if (!textarea) return
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const value = textarea.value
    const newValue = value.substring(0, start) + '  ' + value.substring(end)
    emit('update:modelValue', newValue)
    nextTick(() => {
      textarea.selectionStart = textarea.selectionEnd = start + 2
    })
  }
}
defineExpose({
  focus: () => textareaRef.value?.focus()
})
</script>
<style scoped>
.code-editor {
  width: 100%;
  height: 100%;
  position: relative;
}
.code-editor-container {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
}
.code-highlight,
.code-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  padding: 12px;
  margin: 0;
  border: none;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  white-space: pre-wrap;
  word-wrap: break-word;
  overflow: auto;
  box-sizing: border-box;
}
.code-highlight {
  pointer-events: none;
  background: var(--bg-secondary);
  color: var(--text-primary);
  z-index: 1;
}
.code-highlight code {
  display: block;
  font-family: inherit;
}
.code-input {
  background: transparent;
  color: transparent;
  caret-color: var(--text-primary);
  resize: none;
  outline: none;
  z-index: 2;
  -webkit-text-fill-color: transparent;
}
.code-input::selection {
  background: rgba(var(--primary-rgb), 0.3);
  -webkit-text-fill-color: transparent;
}
:deep(.hl-tag) {
  color: #e06c75;
}
:deep(.hl-attr) {
  color: #d19a66;
}
:deep(.hl-string) {
  color: #98c379;
}
:deep(.hl-comment) {
  color: #5c6370;
  font-style: italic;
}
:deep(.hl-selector) {
  color: #e06c75;
}
:deep(.hl-property) {
  color: #56b6c2;
}
:deep(.hl-value) {
  color: #d19a66;
}
:deep(.hl-number) {
  color: #d19a66;
}
:deep(.hl-color) {
  color: #98c379;
}
</style>
