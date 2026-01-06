export const callWails = async <T = any>(method: string, ...args: any[]): Promise<T> => {
  if (!window.go?.main?.App?.[method]) {
    throw new Error(`Wails method not found: ${method}`)
  }
  try {
    return await window.go.main.App[method](...args)
  } catch (e: any) {
    if (typeof e === 'string') {
      throw new Error(e)
    }
    throw e
  }
}
export const wailsAdapter = {
  auth: {
    login: (email: string, password: string, rememberMe: boolean) =>
      callWails('AuthLogin', email, password, rememberMe),
    logout: (token: string) =>
      callWails('AuthLogout', token),
    getProfile: (token: string) =>
      callWails('AuthGetProfile', token),
    getPoints: (token: string, userId: string) =>
      callWails('AuthGetPoints', token, userId)
  },
  project: {
    getAll: () => callWails('GetProjects'),
    load: (fileName: string) => callWails('LoadProjectByName', fileName),
    save: (fileName: string, data: any) => callWails('SaveProject', fileName, data),
    delete: (fileName: string) => callWails('DeleteProject', fileName),
    rename: (oldName: string, newName: string) => callWails('RenameProject', oldName, newName),
    create: (name: string, type: string) => callWails('CreateProject', name, type),
    updateCurrentNode: (projectName: string, nodeId: number | null) =>
      callWails('UpdateProjectCurrentNode', projectName, nodeId)
  },
  node: {
    createChild: (projectName: string, parentId: number, name: string) =>
      callWails('CreateChildNode', projectName, parentId, name),
    createBrother: (projectName: string, siblingId: number, name: string) =>
      callWails('CreateBrotherNode', projectName, siblingId, name),
    rename: (nodeId: number, newName: string) =>
      callWails('RenameNode', nodeId, newName),
    updateNote: (projectName: string, nodeId: number, note: string) =>
      callWails('UpdateNodeNote', projectName, nodeId, note),
    delete: (projectName: string, nodeId: number) =>
      callWails('DeleteNode', projectName, nodeId)
  },
  worldTree: {
    getNodeContent: (nodeId: number) => callWails('GetNodeContent', nodeId),
    addFolder: (nodeId: number, name: string) => callWails('AddFolder', nodeId, name),
    addFolderWithParent: (nodeId: number, parentChangeId: number | null, name: string) =>
      callWails('AddFolderWithParent', nodeId, parentChangeId, name),
    addCard: (nodeId: number, folderChangeId: number | null, name: string, keyWord: string) =>
      callWails('AddCard', nodeId, folderChangeId, name, keyWord),
    addBlock: (nodeId: number, cardChangeId: number | null, title: string, zone: string) =>
      callWails('AddBlock', nodeId, cardChangeId, title, zone),
    addLine: (nodeId: number, projectId: number, blockChangeId: number, content: string) =>
      callWails('AddLine', nodeId, projectId, blockChangeId, content),
    updateLineContent: (lineId: number, content: string) =>
      callWails('UpdateLineContent', lineId, content),
    getSystemFolderChangeIds: (projectId: number) =>
      callWails('GetSystemFolderChangeIds', projectId),
    getNodeDetail: (nodeId: number) => callWails('GetNodeDetail', nodeId),
    immediateChange: (nodeId: number, change: any) =>
      callWails('ImmediateChange', nodeId, change),
    saveNodeChanges: (nodeId: number, projectId: number, changes: any) =>
      callWails('SaveNodeChanges', nodeId, projectId, changes)
  },
  conversation: {
    getList: (token: string, appId: string, page: number, limit: number) =>
      callWails('GetConversations', token, appId, page, limit),
    getDetail: (token: string, appId: string, conversationId: string) =>
      callWails('GetConversationDetail', token, appId, conversationId),
    delete: (token: string, appId: string, conversationId: string) =>
      callWails('DeleteConversation', token, appId, conversationId),
    rename: (token: string, appId: string, conversationId: string, newName: string) =>
      callWails('RenameConversation', token, appId, conversationId, newName),
    create: (token: string, appId: string, query: string, name: string) =>
      callWails('CreateNewConversation', token, appId, query, name)
  },
  gallery: {
    getImages: () => callWails('GetGalleryImages'),
    getFolders: () => callWails('GetGalleryFolders'),
    createFolder: (name: string) => callWails('CreateGalleryFolder', name),
    deleteFolder: (name: string) => callWails('DeleteGalleryFolder', name),
    deleteImage: (id: string) => callWails('DeleteGalleryImage', id),
    deleteImages: (ids: string[]) => callWails('DeleteGalleryImages', ids),
    moveToFolder: (imageId: string, folderPath: string) =>
      callWails('MoveGalleryImageToFolder', imageId, folderPath),
    moveImageToFolder: (imageId: string, folderPath: string) =>
      callWails('MoveGalleryImageToFolder', imageId, folderPath),
    readAsBase64: (id: string) => callWails('ReadGalleryImageAsBase64', id),
    selectAndAdd: () => callWails('SelectAndAddGalleryImage'),
    selectAndAddToFolder: (folderPath: string) =>
      callWails('SelectAndAddGalleryImageToFolder', folderPath),
    addFromBase64: (base64: string, fileName: string) =>
      callWails('AddGalleryImageFromBase64', base64, fileName),
    addFromBase64ToFolder: (base64: string, fileName: string, folderPath: string) =>
      callWails('AddGalleryImageFromBase64ToFolder', base64, fileName, folderPath),
    renameImage: (id: string, newName: string) =>
      callWails('RenameGalleryImage', id, newName),
    renameFolder: (oldName: string, newName: string) =>
      callWails('RenameGalleryFolder', oldName, newName),
    updateImageUrl: (id: string, url: string) =>
      callWails('UpdateGalleryImageURL', id, url)
  },
  drafts: {
    getAll: () => callWails('GetAllDrafts'),
    create: (draft: any) => callWails('CreateDraft', draft),
    update: (draft: any) => callWails('UpdateDraft', draft),
    delete: (id: string) => callWails('DeleteDraft', id),
    getClipboard: () => callWails('GetClipboardCaptures'),
    moveClipboardToDraft: (captureId: string, name: string, parentId: string) =>
      callWails('MoveClipboardToDraft', captureId, name, parentId),
    startClipboardMonitor: () => callWails('StartClipboardMonitor'),
    stopClipboardMonitor: () => callWails('StopClipboardMonitor'),
    clearAllClipboard: () => callWails('ClearAllClipboardCaptures')
  },
  config: {
    load: () => callWails('LoadConfig'),
    save: (config: any) => callWails('SaveConfig', config),
    update: (key: string, value: any) => callWails('UpdateConfig', key, String(value))
  },
  session: {
    load: () => callWails('LoadSessionState'),
    save: (state: any) => callWails('SaveSessionState', state)
  },
  file: {
    openDialog: (title: string, filters: string) => callWails('OpenFileDialog', title, filters),
    saveDialog: (title: string, defaultName: string, filters: string) =>
      callWails('SaveFileDialog', title, defaultName, filters),
    copyToClipboard: (content: string) => callWails('CopyToClipboard', content)
  },
  system: {
    isDebugMode: () => callWails('IsDebugMode'),
    hideWindow: () => callWails('HideWindow'),
    showWindow: () => callWails('ShowWindow')
  }
}
