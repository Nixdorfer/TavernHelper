import { callWails } from '../adapter/wails'
const BASE_URL = 'https://aipornhub.ltd/console/api'
export const creationApi = {
  async getLocalCreations(): Promise<any[]> {
    return await callWails('GetLocalCreations')
  },
  async getApps(token: string): Promise<any[]> {
    const url = `${BASE_URL}/apps?page=1&limit=100`
    const response = await callWails('FetchWithAuth', token, url, 'GET', '')
    const resp = JSON.parse(response)
    return resp?.data || []
  },
  async getLocalCreationConfig(folderName: string): Promise<any> {
    return await callWails('GetLocalCreationConfig', folderName)
  },
  async getLocalCreationByRemoteId(remoteId: string): Promise<any> {
    return await callWails('GetLocalCreationByRemoteId', remoteId)
  },
  async getModelConfig(token: string, appId: string): Promise<any> {
    const url = `${BASE_URL}/apps/${appId}/model-config`
    const response = await callWails('FetchWithAuth', token, url, 'GET', '')
    return JSON.parse(response)
  },
  async getAppDetail(token: string, appId: string): Promise<any> {
    const url = `${BASE_URL}/apps/${appId}`
    const response = await callWails('FetchWithAuth', token, url, 'GET', '')
    return JSON.parse(response)
  },
  async saveLocalCreation(folderName: string, config: any): Promise<void> {
    await callWails('SaveLocalCreation', folderName, config)
  },
  async saveLocalCreationPage(folderName: string, pageName: string, content: string): Promise<void> {
    await callWails('SaveLocalCreationPage', folderName, pageName, content)
  },
  async deleteLocalCreation(folderName: string): Promise<void> {
    await callWails('DeleteLocalCreation', folderName)
  },
  async sendTestChat(request: {
    provider: string
    messages: { role: string; content: string }[]
    systemPrompt: string
    worldBook: any[]
  }): Promise<{ content?: string; error?: string }> {
    return await callWails('SendTestChat', request)
  }
}
