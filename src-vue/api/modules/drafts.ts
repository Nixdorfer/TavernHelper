import { api } from '../index'
import type { Draft } from '@/types'
export interface ClipboardCapture {
  id: string
  content: string
  createdAt: string
}
export const draftsApi = {
  getAll(): Promise<Draft[]> {
    return api.drafts.getAll()
  },
  create(draft: Omit<Draft, 'id' | 'createdAt' | 'updatedAt'>): Promise<Draft> {
    return api.drafts.create(draft)
  },
  update(draft: Draft) {
    return api.drafts.update(draft)
  },
  delete(id: string) {
    return api.drafts.delete(id)
  },
  getClipboard(): Promise<ClipboardCapture[]> {
    return api.drafts.getClipboard()
  },
  moveClipboardToDraft(captureId: string, name: string, parentId?: string): Promise<Draft> {
    return api.drafts.moveClipboardToDraft(captureId, name, parentId || '')
  },
  startClipboardMonitor() {
    return api.drafts.startClipboardMonitor()
  },
  stopClipboardMonitor() {
    return api.drafts.stopClipboardMonitor()
  },
  clearAllClipboard() {
    return api.drafts.clearAllClipboard()
  },
  copyToClipboard(content: string) {
    return api.file.copyToClipboard(content)
  }
}
