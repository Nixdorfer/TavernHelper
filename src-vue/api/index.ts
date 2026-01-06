import { wailsAdapter } from './adapter/wails'
const BACKEND = import.meta.env.VITE_BACKEND || 'wails'
const adapters: Record<string, typeof wailsAdapter> = {
  wails: wailsAdapter
}
export const api = adapters[BACKEND] || wailsAdapter
export * from './modules/auth'
export * from './modules/project'
export * from './modules/conversation'
export * from './modules/config'
export * from './modules/gallery'
export * from './modules/drafts'
export * from './modules/creation'
