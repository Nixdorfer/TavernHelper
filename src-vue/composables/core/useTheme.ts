import { computed } from 'vue'
import { useConfigStore } from '@/stores'

const COLOR_SCHEMES = [
  { name: 'dark-green', primary: '#10b981', bg: '#1a1a2e' },
  { name: 'dark-blue', primary: '#3b82f6', bg: '#1e293b' },
  { name: 'dark-purple', primary: '#8b5cf6', bg: '#1e1b2e' },
  { name: 'light-blue', primary: '#2563eb', bg: '#f8fafc' },
  { name: 'light-green', primary: '#059669', bg: '#f0fdf4' }
]

export function useTheme() {
  const configStore = useConfigStore()

  const currentScheme = computed(() => {
    return COLOR_SCHEMES[configStore.config.colorScheme] || COLOR_SCHEMES[0]
  })

  const isDark = computed(() => {
    return configStore.config.colorMode === 'dark'
  })

  const primaryColor = computed(() => currentScheme.value.primary)
  const backgroundColor = computed(() => currentScheme.value.bg)

  function setColorScheme(index: number) {
    if (index >= 0 && index < COLOR_SCHEMES.length) {
      configStore.update('colorScheme', index)
    }
  }

  function setColorMode(mode: 'dark' | 'light' | 'auto') {
    configStore.update('colorMode', mode)
  }

  function toggleDarkMode() {
    configStore.update('colorMode', isDark.value ? 'light' : 'dark')
  }

  function applyTheme() {
    const root = document.documentElement
    root.style.setProperty('--primary-color', primaryColor.value)
    root.style.setProperty('--bg-color', backgroundColor.value)
    root.classList.toggle('dark', isDark.value)
  }

  return {
    currentScheme,
    isDark,
    primaryColor,
    backgroundColor,
    colorSchemes: COLOR_SCHEMES,
    setColorScheme,
    setColorMode,
    toggleDarkMode,
    applyTheme
  }
}
