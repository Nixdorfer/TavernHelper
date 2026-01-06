<template>
  <div class="sync-editor">
    <div class="sync-legend">
      <span class="legend-item"><span class="dot dot-circle dot-green"></span>继承且保留</span>
      <span class="legend-item"><span class="dot dot-circle dot-yellow"></span>继承但部分删除</span>
      <span class="legend-item"><span class="dot dot-circle dot-red"></span>继承但全部删除</span>
      <span class="legend-item"><span class="dot dot-square dot-green"></span>新增且保留</span>
      <span class="legend-item"><span class="dot dot-square dot-yellow"></span>新增但部分删除</span>
      <span class="legend-item"><span class="dot dot-square dot-red"></span>新增但删除</span>
    </div>
    <div class="editor-container">
      <div class="sync-gutter" ref="gutterRef">
        <div
          v-for="(line, index) in computedLines"
          :key="index"
          class="gutter-line"
        >
          <div
            v-if="line.syncDot"
            :class="getSyncDotClasses(line.syncDot)"
            @click="copySerial(line.sn)"
            :title="line.sn ? '点击复制: ' + line.sn : ''"
          ></div>
          <div v-else-if="isDirty && line.sn" class="sync-placeholder"></div>
        </div>
      </div>
      <div class="text-area-wrapper">
        <textarea
          ref="textareaRef"
          :value="modelValue"
          @input="handleInput"
          @scroll="syncScroll"
          :placeholder="placeholder"
          :rows="rows"
          class="sync-textarea"
        ></textarea>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed } from 'vue'
import { logger } from '@/utils/logger'
interface LineInfo {
  sn: string
  content: string
  syncDot?: string
}
const props = withDefaults(defineProps<{
  modelValue: string
  placeholder?: string
  rows?: number
  lines?: LineInfo[]
  isDirty?: boolean
}>(), {
  placeholder: '',
  rows: 12,
  lines: () => [],
  isDirty: false
})
const emit = defineEmits<{
  'update:modelValue': [value: string]
  'copy-serial': [serial: string]
}>()
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const gutterRef = ref<HTMLElement | null>(null)
const syncDotMap = computed(() => {
  const map: Record<number, { sn: string; syncDot: string }> = {}
  if (!props.modelValue || !props.lines.length) return map
  const contentLines = props.modelValue.split('\n')
  props.lines.forEach((line, index) => {
    if (index < contentLines.length && (line.syncDot || line.sn)) {
      map[index] = { sn: line.sn || '', syncDot: line.syncDot || '' }
    }
  })
  return map
})
const computedLines = computed(() => {
  const contentLines = props.modelValue.split('\n')
  return contentLines.map((content, index) => {
    const dotInfo = syncDotMap.value[index]
    return {
      content,
      sn: dotInfo?.sn || '',
      syncDot: dotInfo?.syncDot || ''
    }
  })
})
function getSyncDotClasses(syncDot: string) {
  const [color, shape] = syncDot.split('-')
  const classes = ['sync-dot']
  if (shape === 'dot') {
    classes.push('dot-circle')
  } else if (shape === 'square') {
    classes.push('dot-square')
  }
  if (color === 'green') {
    classes.push('dot-green')
  } else if (color === 'yellow') {
    classes.push('dot-yellow')
  } else if (color === 'red') {
    classes.push('dot-red')
  }
  return classes
}
async function copySerial(sn: string) {
  if (!sn) return
  try {
    await navigator.clipboard.writeText(sn)
    emit('copy-serial', sn)
  } catch (err) {
    logger.error('复制失败:', err)
  }
}
function handleInput(e: Event) {
  const target = e.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
}
function syncScroll() {
  if (gutterRef.value && textareaRef.value) {
    gutterRef.value.scrollTop = textareaRef.value.scrollTop
  }
}
</script>
<style scoped>
.sync-editor {
  position: relative;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  background: var(--bg-primary);
}
.sync-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-color);
  font-size: 11px;
  color: var(--text-secondary);
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
}
.dot {
  width: 10px;
  height: 10px;
  flex-shrink: 0;
}
.dot-circle {
  border-radius: 50%;
}
.dot-square {
  border-radius: 2px;
}
.dot-green {
  background-color: #22c55e;
}
.dot-yellow {
  background-color: #f59e0b;
}
.dot-red {
  background-color: #ef4444;
}
.editor-container {
  display: flex;
  position: relative;
}
.sync-gutter {
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  overflow-y: hidden;
  overflow-x: hidden;
  flex-shrink: 0;
  width: 24px;
}
.gutter-line {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 24px;
  line-height: 24px;
}
.sync-dot {
  width: 8px;
  height: 8px;
  flex-shrink: 0;
  cursor: pointer;
  transition: transform 0.15s;
}
.sync-dot:hover {
  transform: scale(1.4);
}
.sync-placeholder {
  width: 10px;
  height: 2px;
  background-color: rgba(255, 255, 255, 0.4);
  border-radius: 1px;
}
.text-area-wrapper {
  flex: 1;
  position: relative;
  display: flex;
}
.sync-textarea {
  flex: 1;
  width: 100%;
  min-height: 200px;
  padding: 0 12px;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  line-height: 24px;
  resize: none;
  outline: none;
}
.sync-textarea::placeholder {
  color: var(--text-secondary);
  opacity: 0.6;
}
</style>
