import { computed } from 'vue'
import { useGalleryStore } from '@/stores'
import { useNotification } from '../core/useNotification'
import { useConfirm } from '../core/useConfirm'

export function useGallery() {
  const store = useGalleryStore()
  const { success, error: showError } = useNotification()
  const { confirmDelete } = useConfirm()

  const images = computed(() => store.filteredImages)
  const folders = computed(() => store.folders)
  const selectedImage = computed(() => store.selectedImage)
  const isLoading = computed(() => store.isLoading)

  async function loadImages() {
    await store.loadImages()
  }

  async function loadFolders() {
    await store.loadFolders()
  }

  async function createFolder(name: string) {
    try {
      await store.createFolder(name)
      success('文件夹创建成功')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function deleteFolder(name: string) {
    const confirmed = await confirmDelete(name)
    if (!confirmed) return false
    try {
      await store.deleteFolder(name)
      success('文件夹已删除')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function deleteImage(id: string) {
    const confirmed = await confirmDelete('此图片')
    if (!confirmed) return false
    try {
      await store.deleteImage(id)
      success('图片已删除')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  async function moveToFolder(imageId: string, folderPath: string) {
    try {
      await store.moveToFolder(imageId, folderPath)
      success('图片已移动')
      return true
    } catch (e) {
      showError((e as Error).message)
      return false
    }
  }

  function selectImage(imageId: string | null) {
    store.selectImage(imageId)
  }

  function setViewMode(mode: 'grid' | 'list') {
    store.setViewMode(mode)
  }

  function setCurrentFolder(folder: string | null) {
    store.setCurrentFolder(folder)
  }

  return {
    images,
    folders,
    selectedImage,
    isLoading,
    currentFolder: store.currentFolder,
    viewMode: store.viewMode,
    loadImages,
    loadFolders,
    createFolder,
    deleteFolder,
    deleteImage,
    moveToFolder,
    selectImage,
    setViewMode,
    setCurrentFolder
  }
}
