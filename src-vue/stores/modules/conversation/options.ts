import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface ConversationOption {
  id: string
  label: string
  value: any
}

export const useOptionsStore = defineStore('conversation-options', () => {
  const showOverlay = ref(false)
  const list = ref<ConversationOption[]>([])
  const selected = ref<Record<string, boolean>>({})
  const remarks = ref<Record<string, string>>({})
  const isReadonly = ref(false)

  const hasSelectedOptions = computed(() => {
    return Object.values(selected.value).some(v => v)
  })

  const selectedCount = computed(() => {
    return Object.values(selected.value).filter(v => v).length
  })

  const selectedDisplay = computed(() => {
    return `${selectedCount.value} 个选项`
  })

  function setList(options: ConversationOption[]) {
    list.value = options
  }

  function toggle(optionId: string) {
    selected.value[optionId] = !selected.value[optionId]
  }

  function select(optionId: string) {
    selected.value[optionId] = true
  }

  function deselect(optionId: string) {
    selected.value[optionId] = false
  }

  function setRemark(optionId: string, remark: string) {
    remarks.value[optionId] = remark
  }

  function clear() {
    selected.value = {}
    remarks.value = {}
  }

  function setReadonly(readonly: boolean) {
    isReadonly.value = readonly
  }

  function showOptions() {
    showOverlay.value = true
  }

  function hideOptions() {
    showOverlay.value = false
  }

  return {
    showOverlay,
    list,
    selected,
    remarks,
    isReadonly,
    hasSelectedOptions,
    selectedCount,
    selectedDisplay,
    setList,
    toggle,
    select,
    deselect,
    setRemark,
    clear,
    setReadonly,
    showOptions,
    hideOptions
  }
})
