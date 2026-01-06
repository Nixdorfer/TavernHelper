import { useConfirmStore } from '@/stores'

export function useConfirm() {
  const store = useConfirmStore()
  return {
    show: store.show,
    confirm: store.confirm,
    cancel: store.cancel,
    confirmDelete: store.confirmDelete,
    isVisible: store.isVisible,
    options: store.options
  }
}
