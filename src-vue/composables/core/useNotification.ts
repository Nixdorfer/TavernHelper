import { useNotificationStore } from '@/stores'

export function useNotification() {
  const store = useNotificationStore()
  return {
    show: store.show,
    success: store.success,
    error: store.error,
    warning: store.warning,
    info: store.info,
    remove: store.remove,
    clear: store.clear,
    notifications: store.notifications
  }
}
