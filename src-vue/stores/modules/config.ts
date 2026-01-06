import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { configApi, sessionApi } from '@/api/modules/config'
import type { AppConfig } from '@/api/modules/config'

export const useConfigStore = defineStore('config', () => {
  const config = ref<AppConfig>({
    theme: 'default',
    language: 'zh-CN',
    lastOpenedProject: '',
    colorScheme: 1,
    colorMode: 'dark',
    systemPrompt: '',
    systemPromptType: 'system',
    debugMode: false,
    safeMode: false,
    safeModeAction: 'randomChars',
    safeModeTemplate: '',
    debugTestReply: '',
    bytedanceApiKey: '',
    autoGenerateImage: false,
    noImageMode: false,
    strictMode: false,
    claudeApiKey: '',
    geminiApiKey: '',
    grokApiKey: ''
  })
  const isLoading = ref(false)
  const noImageMode = computed(() => config.value.noImageMode)

  async function load() {
    isLoading.value = true
    try {
      const loaded = await configApi.load()
      Object.assign(config.value, loaded)
    } finally {
      isLoading.value = false
    }
  }

  async function save() {
    await configApi.save(config.value)
  }

  async function update<K extends keyof AppConfig>(key: K, value: AppConfig[K]) {
    config.value[key] = value
    await configApi.update(key, value)
  }

  async function loadSession() {
    return await sessionApi.load()
  }

  async function saveSession(state: any) {
    await sessionApi.save(state)
  }

  return {
    config,
    isLoading,
    noImageMode,
    load,
    save,
    update,
    loadSession,
    saveSession
  }
})
