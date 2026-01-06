import { api } from '../index'
import type { Conversation, Message } from '@/types'
export const conversationApi = {
  getList(token: string, appId: string, page = 1, limit = 20): Promise<{
    items: Conversation[]
    total: number
    hasMore: boolean
  }> {
    return api.conversation.getList(token, appId, page, limit)
  },
  getDetail(token: string, appId: string, conversationId: string): Promise<{
    conversation: Conversation
    messages: Message[]
  }> {
    return api.conversation.getDetail(token, appId, conversationId)
  },
  delete(token: string, appId: string, conversationId: string) {
    return api.conversation.delete(token, appId, conversationId)
  },
  rename(token: string, appId: string, conversationId: string, newName: string) {
    return api.conversation.rename(token, appId, conversationId, newName)
  },
  create(token: string, appId: string, query: string, name: string) {
    return api.conversation.create(token, appId, query, name)
  }
}
