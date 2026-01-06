import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ConfirmOptions {
  title?: string
  message: string
  confirmText?: string
  cancelText?: string
  type?: 'info' | 'warning' | 'danger'
}

export const useConfirmStore = defineStore('confirm', () => {
  const isVisible = ref(false)
  const options = ref<ConfirmOptions>({ message: '' })
  let resolvePromise: ((value: boolean) => void) | null = null

  function show(opts: ConfirmOptions): Promise<boolean> {
    options.value = {
      title: opts.title || '确认',
      message: opts.message,
      confirmText: opts.confirmText || '确定',
      cancelText: opts.cancelText || '取消',
      type: opts.type || 'info'
    }
    isVisible.value = true
    return new Promise(resolve => {
      resolvePromise = resolve
    })
  }

  function confirm() {
    isVisible.value = false
    resolvePromise?.(true)
    resolvePromise = null
  }

  function cancel() {
    isVisible.value = false
    resolvePromise?.(false)
    resolvePromise = null
  }

  async function confirmDelete(itemName: string): Promise<boolean> {
    return show({
      title: '确认删除',
      message: `确定要删除"${itemName}"吗？此操作不可撤销。`,
      confirmText: '删除',
      cancelText: '取消',
      type: 'danger'
    })
  }

  return {
    isVisible,
    options,
    show,
    confirm,
    cancel,
    confirmDelete
  }
})
