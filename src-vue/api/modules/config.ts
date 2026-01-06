import { api } from '../index'
export interface AppConfig {
  theme: string
  language: string
  lastOpenedProject: string
  colorScheme: number
  colorMode: string
  systemPrompt: string
  systemPromptType: string
  debugMode: boolean
  safeMode: boolean
  safeModeAction: string
  safeModeTemplate: string
  debugTestReply: string
  bytedanceApiKey: string
  autoGenerateImage: boolean
  noImageMode: boolean
  strictMode: boolean
  claudeApiKey: string
  geminiApiKey: string
  grokApiKey: string
}
export const configApi = {
  load(): Promise<AppConfig> {
    return api.config.load()
  },
  save(config: Partial<AppConfig>) {
    return api.config.save(config)
  },
  update(key: keyof AppConfig, value: any) {
    return api.config.update(key, value)
  }
}
export const sessionApi = {
  load() {
    return api.session.load()
  },
  save(state: any) {
    return api.session.save(state)
  }
}
