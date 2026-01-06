import { invoke } from '@tauri-apps/api/core'

export const callTauri = async <T = any>(command: string, args?: Record<string, any>): Promise<T> => {
  try {
    return await invoke(command, args ?? {})
  } catch (e: any) {
    if (typeof e === 'string') {
      throw new Error(e)
    }
    throw e
  }
}

export const tauriAdapter = {
  config: {
    load: () => callTauri('load_config'),
    save: (config: any) => callTauri('save_config', { config }),
    update: (key: string, value: any) => callTauri('update_config', { key, value }),
  },
  session: {
    load: () => callTauri('load_session_state'),
    save: (state: any) => callTauri('save_session_state', { state }),
  },
  project: {
    getAll: () => callTauri('get_projects'),
    load: (name: string) => callTauri('load_project_by_name', { name }),
    save: (fileName: string, data: any) => callTauri('save_project', { fileName, data }),
    delete: (fileName: string) => callTauri('delete_project', { fileName }),
    rename: (oldName: string, newName: string) => callTauri('rename_project', { oldName, newName }),
    create: (name: string, projectType: string) => callTauri('create_project', { name, projectType }),
    updateCurrentNode: (projectName: string, nodeId: number) => callTauri('update_project_current_node', { projectName, nodeId }),
  },
  node: {
    createChild: (projectName: string, parentId: number, name: string) => callTauri('create_child_node', { projectName, parentId, name }),
    createBrother: (projectName: string, siblingId: number, name: string) => callTauri('create_brother_node', { projectName, siblingId, name }),
    rename: (nodeId: number, newName: string) => callTauri('rename_node', { nodeId, newName }),
    updateNote: (projectName: string, nodeId: number, note: string) => callTauri('update_node_note', { projectName, nodeId, note }),
    updateDesc: (nodeId: number, desc: string) => callTauri('update_node_desc', { nodeId, desc }),
    delete: (projectName: string, nodeId: number) => callTauri('delete_node', { projectName, nodeId }),
  },
  worldTree: {
    getNodeContent: (nodeId: number) => callTauri('get_node_content', { nodeId }),
    addFolder: (nodeId: number, name: string) => callTauri('add_folder', { nodeId, name }),
    addFolderWithParent: (nodeId: number, parentChangeId: number, name: string) => callTauri('add_folder_with_parent', { nodeId, parentChangeId, name }),
    addCard: (nodeId: number, folderChangeId: number, name: string, keyWord: string) => callTauri('add_card', { nodeId, folderChangeId, name, keyWord }),
    addBlock: (nodeId: number, cardChangeId: number, title: string, zone: string) => callTauri('add_block', { nodeId, cardChangeId, title, zone }),
    addLine: (nodeId: number, projectId: number, blockChangeId: number, content: string) => callTauri('add_line', { nodeId, projectId, blockChangeId, content }),
    updateLineContent: (lineId: number, content: string) => callTauri('update_line_content', { lineId, content }),
    deleteFolder: (nodeId: number, folderChangeId: number) => callTauri('delete_folder', { nodeId, folderChangeId }),
    deleteCard: (nodeId: number, cardChangeId: number) => callTauri('delete_card', { nodeId, cardChangeId }),
    deleteBlock: (nodeId: number, blockChangeId: number) => callTauri('delete_block', { nodeId, blockChangeId }),
    deleteLine: (nodeId: number, lineChangeId: number) => callTauri('delete_line', { nodeId, lineChangeId }),
    getNodeDetail: (nodeId: number) => callTauri('get_node_detail', { nodeId }),
    immediateChange: (nodeId: number, change: any) => callTauri('immediate_change', { nodeId, change }),
    saveNodeChanges: (nodeId: number, projectId: number, changes: any) => callTauri('save_node_changes', { nodeId, projectId, changes }),
    getSystemFolderChangeIds: (projectId: number) => callTauri('get_system_folder_change_ids', { projectId }),
  },
  gallery: {
    getImages: () => callTauri('get_gallery_images'),
    getFolders: () => callTauri('get_gallery_folders'),
    createFolder: (name: string) => callTauri('create_gallery_folder', { name }),
    deleteFolder: (name: string) => callTauri('delete_gallery_folder', { name }),
    deleteImage: (id: string) => callTauri('delete_gallery_image', { id }),
    deleteImages: (ids: string[]) => callTauri('delete_gallery_images', { ids }),
    moveToFolder: (imageId: string, folderPath: string) => callTauri('move_gallery_image_to_folder', { imageId, folderPath }),
    readAsBase64: (id: string) => callTauri('read_gallery_image_as_base64', { id }),
    selectAndAdd: () => callTauri('select_and_add_gallery_image'),
    selectAndAddToFolder: (folderPath: string) => callTauri('select_and_add_to_folder', { folderPath }),
    addFromBase64: (base64Data: string, fileName: string) => callTauri('add_gallery_image_from_base64', { base64Data, fileName }),
    addFromBase64ToFolder: (base64Data: string, fileName: string, folderPath: string) => callTauri('add_gallery_image_from_base64_to_folder', { base64Data, fileName, folderPath }),
    renameImage: (id: string, newName: string) => callTauri('rename_gallery_image', { id, newName }),
    renameFolder: (oldName: string, newName: string) => callTauri('rename_gallery_folder', { oldName, newName }),
    updateImageUrl: (id: string, url: string) => callTauri('update_gallery_image_url', { id, url }),
  },
  drafts: {
    getAll: () => callTauri('get_all_drafts'),
    create: (draft: any) => callTauri('create_draft', { draft }),
    update: (draft: any) => callTauri('update_draft', { draft }),
    delete: (id: number) => callTauri('delete_draft', { id }),
    getClipboard: () => callTauri('get_clipboard_captures'),
    moveClipboardToDraft: (captureId: number, name: string, parentId: number | null) => callTauri('move_clipboard_to_draft', { captureId, name, parentId }),
    startClipboardMonitor: () => callTauri('start_clipboard_monitor'),
    stopClipboardMonitor: () => callTauri('stop_clipboard_monitor'),
    clearAllClipboard: () => callTauri('clear_all_clipboard_captures'),
    copyToClipboard: (content: string) => callTauri('copy_to_clipboard', { content }),
  },
  auth: {
    login: (email: string, password: string, rememberMe: boolean) => callTauri('auth_login', { email, password, rememberMe }),
    getProfile: (token: string) => callTauri('auth_get_profile', { token }),
    logout: (token: string) => callTauri('auth_logout', { token }),
    getPoints: (token: string, userId: string) => callTauri('auth_get_points', { token, userId }),
  },
  conversation: {
    getList: (token: string, appId: string, page: number, limit: number) => callTauri('get_conversations', { token, appId, page, limit }),
    getDetail: (token: string, appId: string, conversationId: string) => callTauri('get_conversation_detail', { token, appId, conversationId }),
    delete: (token: string, appId: string, conversationId: string) => callTauri('delete_conversation', { token, appId, conversationId }),
    rename: (token: string, appId: string, conversationId: string, newName: string) => callTauri('rename_conversation', { token, appId, conversationId, newName }),
    create: (token: string, appId: string, query: string, name: string) => callTauri('create_new_conversation', { token, appId, query, name }),
  },
  creation: {
    getLocalCreations: () => callTauri('get_local_creations'),
    getLocalCreationConfig: (folderName: string) => callTauri('get_local_creation_config', { folderName }),
    saveLocalCreation: (folderName: string, config: any) => callTauri('save_local_creation', { folderName, config }),
    deleteLocalCreation: (folderName: string) => callTauri('delete_local_creation', { folderName }),
    fetchWithAuth: (token: string, url: string, method: string, body?: string) => callTauri('fetch_with_auth', { token, url, method, body }),
    sendTestChat: (request: any) => callTauri('send_test_chat', { request }),
  },
  system: {
    hideWindow: () => callTauri('hide_window'),
    showWindow: () => callTauri('show_window'),
    isDebugMode: () => callTauri('is_debug_mode'),
    openFileDialog: (title: string, filters: [string, string[]][]) => callTauri('open_file_dialog', { title, filters }),
    saveFileDialog: (title: string, defaultName: string, filters: [string, string[]][]) => callTauri('save_file_dialog', { title, defaultName, filters }),
  },
}

export default tauriAdapter
