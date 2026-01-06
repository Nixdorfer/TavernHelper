import { api } from '../index'

const BASE_URL = 'https://aipornhub.ltd/console/api'

export const creationApi = {
  async getLocalCreations(): Promise<any[]> {
    return await api.creation.getLocalCreations()
  },
  async getApps(token: string): Promise<any[]> {
    const url = `${BASE_URL}/apps?page=1&limit=100`
    const response = await api.creation.fetchWithAuth(token, url, 'GET', '')
    const resp = JSON.parse(response)
    return resp?.data || []
  },
  async getLocalCreationConfig(folderName: string): Promise<any> {
    return await api.creation.getLocalCreationConfig(folderName)
  },
  async getLocalCreationByRemoteId(remoteId: string): Promise<any> {
    return await api.creation.getLocalCreationByRemoteId(remoteId)
  },
  async getModelConfig(token: string, appId: string): Promise<any> {
    const url = `${BASE_URL}/apps/${appId}/model-config`
    const response = await api.creation.fetchWithAuth(token, url, 'GET', '')
    return JSON.parse(response)
  },
  async getAppDetail(token: string, appId: string): Promise<any> {
    const url = `${BASE_URL}/apps/${appId}`
    const response = await api.creation.fetchWithAuth(token, url, 'GET', '')
    return JSON.parse(response)
  },
  async saveLocalCreation(folderName: string, config: any): Promise<void> {
    await api.creation.saveLocalCreation(folderName, config)
  },
  async saveLocalCreationPage(folderName: string, pageName: string, content: string): Promise<void> {
    await api.creation.saveLocalCreationPage(folderName, pageName, content)
  },
  async deleteLocalCreation(folderName: string): Promise<void> {
    await api.creation.deleteLocalCreation(folderName)
  },
  async sendTestChat(request: {
    provider: string
    messages: { role: string; content: string }[]
    systemPrompt: string
    worldBook: any[]
  }): Promise<{ content?: string; error?: string }> {
    return await api.creation.sendTestChat(request)
  }
}
