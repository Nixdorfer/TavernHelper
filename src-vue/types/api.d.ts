export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
}
export interface AuthLoginResponse {
  access_token: string
  token_type: string
  expires_in: number
}
export interface UserProfile {
  id: string
  name: string
  email: string
  avatar?: string
}
export interface ProjectInfo {
  id: number
  fileName: string
  name: string
  type: 'play' | 'create'
  updatedAt: string
  isEncrypted?: boolean
}
export interface Project {
  id: number
  name: string
  fileName: string
  type: 'play' | 'create'
  createTime: string
  updateTime: string
  timeline: TimelineNode[]
  worldBook?: WorldBookEntry[]
  currentNode?: number
  branchData?: Record<string, string>
}
export interface TimelineNode {
  id: number
  name: string
  parentId?: number
  note?: string
  tags?: string[]
  createdAt?: string
  pre_text?: PromptEntry[]
  post_text?: PromptEntry[]
  pre_prompt?: PromptEntry[]
  world_book?: WorldBookEntry[]
  children?: TimelineNode[]
}
export interface PromptEntry {
  id: string
  name?: string
  content?: string
  isFolder?: boolean
  collapsed?: boolean
}
export interface WorldBookEntry {
  id: string
  name?: string
  key: string
  keywords?: string[]
  matchMode?: 'or' | 'and'
  contentItems?: ContentItem[]
  isFolder?: boolean
  parentId?: string
}
export interface ContentItem {
  id: string
  title?: string
  content: string
  keyRegion?: number
  valueRegion?: number
}
export interface Conversation {
  id: string
  name: string
  appId: string
  createdAt: string
  updatedAt: string
}
export interface Message {
  id: string
  query: string
  answer: string
  createdAt: string
  files?: string[]
}
export interface GalleryImage {
  id: string
  hash: string
  localPath: string
  remoteUrl?: string
  fileName: string
  fileSize: number
  createdAt: string
  folderPath?: string
  isValid?: boolean
}
export interface Draft {
  id: string
  name: string
  content: string
  parentId?: string
  isFolder?: boolean
  sortOrder?: number
  createdAt: string
  updatedAt: string
}
export interface TriggerTemplate {
  mode: 'and' | 'or'
  words: string[]
  system: boolean
  user: boolean
  ai: boolean
}
export interface LineData {
  content: string | null
  syncDot: string
}
export interface BlockTemplate {
  title: string
  lines: Record<string, LineData>
}
export interface CardTemplate {
  name: string
  desc?: string
  trigger?: TriggerTemplate
  blocks: Record<string, BlockTemplate>
  image?: string[]
}
export interface FolderTemplate {
  name: string
  desc?: string
  folders: Record<string, FolderTemplate>
  cards: Record<string, CardTemplate>
}
export interface NodeTemplate {
  name: string
  desc?: string
  folders: Record<string, FolderTemplate>
}
export interface ImmediateChange {
  name: string
  action: 'add' | 'del'
  level: 'folder' | 'card' | 'block'
  target: number | null
}
export interface SaveFolderChange {
  name?: string
  desc?: string
}
export interface SaveCardChange {
  name?: string
  desc?: string
  trigger?: TriggerTemplate
  image?: string[]
}
export interface SaveBlockChange {
  name?: string
  content?: string
}
export interface SaveChanges {
  folder?: Record<string, SaveFolderChange>
  card?: Record<string, SaveCardChange>
  block?: Record<string, SaveBlockChange>
}
export interface SaveNodeChangesResult {
  blocks: Record<string, Record<string, LineData>>
}
