import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

const platform = process.env.VITE_PLATFORM || 'pc'

export default defineConfig({
  plugins: [vue()],
  define: {
    'import.meta.env.VITE_PLATFORM': JSON.stringify(platform)
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, '.'),
      '@api': path.resolve(__dirname, 'api'),
      '@stores': path.resolve(__dirname, 'stores'),
      '@composables': path.resolve(__dirname, 'composables'),
      '@components': path.resolve(__dirname, 'components'),
      '@platform': path.resolve(__dirname, `platforms/${platform}`),
      '@types': path.resolve(__dirname, 'types')
    }
  },
  server: {
    port: platform === 'mobile' ? 5174 : 5173,
    hmr: {
      host: 'localhost',
      protocol: 'ws'
    }
  },
  build: {
    outDir: path.resolve(__dirname, `dist/${platform}`),
    emptyOutDir: true,
    chunkSizeWarningLimit: 1500
  }
})
