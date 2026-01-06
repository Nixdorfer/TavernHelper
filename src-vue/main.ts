import { createApp } from 'vue'
import { createPinia } from 'pinia'
import './styles/global.css'
const PLATFORM = import.meta.env.VITE_PLATFORM || 'pc'
const loadApp = async () => {
  const App = PLATFORM === 'mobile'
    ? (await import('./platforms/mobile/App.vue')).default
    : (await import('./platforms/pc/App.vue')).default
  const app = createApp(App)
  const pinia = createPinia()
  app.use(pinia)
  app.mount('#app')
}
loadApp()
