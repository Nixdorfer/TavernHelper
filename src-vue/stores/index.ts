import { createPinia } from 'pinia'

export const pinia = createPinia()

export { useAuthStore } from './modules/auth'
export { useConfigStore } from './modules/config'
export { useAppStore } from './modules/app'
export { useNotificationStore } from './modules/notification'
export { useConfirmStore } from './modules/confirm'
export {
  useConversationStore,
  useMessagesStore,
  useStreamingStore,
  useWorldTreeStore,
  useOptionsStore,
  useRevisionsStore
} from './modules/conversation'
export {
  useWorldbookProjectStore,
  useWorldbookEditorStore,
  useWorldbookTabsStore,
  useWorldbookTreeStore
} from './modules/worldbook'
export { usePlazaStore } from './modules/plaza'
export { useGalleryStore } from './modules/gallery'
export { useDraftsStore } from './modules/drafts'
export { useCreationStore } from './modules/creation'
