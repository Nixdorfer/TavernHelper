/// <reference types="vite/client" />
declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}
interface ImportMetaEnv {
  readonly VITE_PLATFORM: 'pc' | 'mobile'
  readonly VITE_BACKEND: 'wails' | 'http' | 'mock'
}
interface ImportMeta {
  readonly env: ImportMetaEnv
}
interface Window {
  go?: {
    main: {
      App: Record<string, (...args: any[]) => Promise<any>>
    }
  }
}
