import { ref, onMounted, nextTick } from 'vue'
import type { Ref } from 'vue'

export function useAutoFocus(delay = 0) {
  const inputRef = ref<HTMLInputElement | null>(null)

  function focus() {
    nextTick(() => {
      setTimeout(() => {
        inputRef.value?.focus()
        inputRef.value?.select()
      }, delay)
    })
  }

  onMounted(() => {
    if (delay >= 0) {
      focus()
    }
  })

  return {
    inputRef,
    focus
  }
}

export function useFocusOnShow(showRef: Ref<boolean>, delay = 50) {
  const inputRef = ref<HTMLInputElement | null>(null)

  function focus() {
    if (showRef.value && inputRef.value) {
      setTimeout(() => {
        inputRef.value?.focus()
        inputRef.value?.select()
      }, delay)
    }
  }

  return {
    inputRef,
    focus
  }
}
