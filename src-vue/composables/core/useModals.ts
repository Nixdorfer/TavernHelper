import { ref, computed } from 'vue'

export type ModalName = 'login' | 'settings' | 'imagePreview' | 'imagePicker' | 'nodeForm' | 'entryForm'

export function useModals() {
  const activeModals = ref<Set<ModalName>>(new Set())
  const modalData = ref<Record<string, any>>({})

  const hasActiveModal = computed(() => activeModals.value.size > 0)

  function open(name: ModalName, data?: any) {
    activeModals.value.add(name)
    if (data) {
      modalData.value[name] = data
    }
  }

  function close(name: ModalName) {
    activeModals.value.delete(name)
    delete modalData.value[name]
  }

  function isOpen(name: ModalName): boolean {
    return activeModals.value.has(name)
  }

  function getData<T = any>(name: ModalName): T | undefined {
    return modalData.value[name]
  }

  function closeAll() {
    activeModals.value.clear()
    modalData.value = {}
  }

  return {
    activeModals,
    modalData,
    hasActiveModal,
    open,
    close,
    isOpen,
    getData,
    closeAll
  }
}
