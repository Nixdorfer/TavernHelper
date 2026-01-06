import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './styles/global.css'
import { api } from './api'
import { initLogger } from './utils/logger'
const PLATFORM = import.meta.env.VITE_PLATFORM || 'pc'
const loadApp = async () => {
  const logEnabled = await api.system.isLogEnabled().catch(() => false)
  initLogger(logEnabled)
  const App = PLATFORM === 'mobile'
    ? (await import('./platforms/mobile/App.vue')).default
    : (await import('./platforms/pc/App.vue')).default
  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)
  app.mount('#app')
}
loadApp()
