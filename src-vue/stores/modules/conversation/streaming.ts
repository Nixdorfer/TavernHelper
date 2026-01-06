import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useStreamingStore = defineStore('conversation-streaming', () => {
  const isEnabled = ref(true)
  const messageId = ref<string | null>(null)
  const taskId = ref<string | null>(null)
  const reader = ref<ReadableStreamDefaultReader | null>(null)
  const isActive = ref(false)
  const buffer = ref('')

  function start(msgId: string, tId: string) {
    messageId.value = msgId
    taskId.value = tId
    isActive.value = true
    buffer.value = ''
  }

  function setReader(r: ReadableStreamDefaultReader | null) {
    reader.value = r
  }

  function appendBuffer(text: string) {
    buffer.value += text
  }

  function stop() {
    if (reader.value) {
      try {
        reader.value.cancel()
      } catch (e) {}
    }
    messageId.value = null
    taskId.value = null
    reader.value = null
    isActive.value = false
    buffer.value = ''
  }

  function setEnabled(enabled: boolean) {
    isEnabled.value = enabled
  }

  return {
    isEnabled,
    messageId,
    taskId,
    reader,
    isActive,
    buffer,
    start,
    setReader,
    appendBuffer,
    stop,
    setEnabled
  }
})
