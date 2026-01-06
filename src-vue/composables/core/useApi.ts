import { ref } from 'vue'
import { useNotificationStore } from '@/stores'

export function useApi<T = any>() {
  const isLoading = ref(false)
  const error = ref<Error | null>(null)
  const data = ref<T | null>(null)
  const notificationStore = useNotificationStore()

  async function execute<R = T>(
    fn: () => Promise<R>,
    options: {
      showError?: boolean
      errorMessage?: string
    } = {}
  ): Promise<R | null> {
    const { showError = true, errorMessage } = options
    isLoading.value = true
    error.value = null
    try {
      const result = await fn()
      data.value = result as any
      return result
    } catch (e) {
      error.value = e as Error
      if (showError) {
        notificationStore.error(errorMessage || (e as Error).message)
      }
      return null
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    error,
    data,
    execute
  }
}
