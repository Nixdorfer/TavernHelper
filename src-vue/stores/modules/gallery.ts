import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { GalleryImage } from '@/types'
import { galleryApi } from '@/api/modules/gallery'
import type { GalleryFolder } from '@/api/modules/gallery'

export const useGalleryStore = defineStore('gallery', () => {
  const images = ref<GalleryImage[]>([])
  const folders = ref<GalleryFolder[]>([])
  const isLoading = ref(false)
  const currentFolder = ref<string | null>(null)
  const selectedImageId = ref<string | null>(null)
  const viewMode = ref<'grid' | 'list'>('grid')
  const sortBy = ref<'date' | 'name'>('date')
  const sortOrder = ref<'asc' | 'desc'>('desc')

  const filteredImages = computed(() => {
    let result = images.value
    if (currentFolder.value) {
      result = result.filter(img => img.folder === currentFolder.value)
    }
    result = [...result].sort((a, b) => {
      let cmp = 0
      if (sortBy.value === 'date') {
        cmp = new Date(a.createdAt).getTime() - new Date(b.createdAt).getTime()
      } else {
        cmp = (a.name || '').localeCompare(b.name || '')
      }
      return sortOrder.value === 'desc' ? -cmp : cmp
    })
    return result
  })

  const selectedImage = computed(() => {
    return images.value.find(img => img.id === selectedImageId.value) || null
  })

  async function loadImages() {
    isLoading.value = true
    try {
      images.value = await galleryApi.getImages()
    } finally {
      isLoading.value = false
    }
  }

  async function loadFolders() {
    folders.value = await galleryApi.getFolders()
  }

  async function createFolder(name: string) {
    const folder = await galleryApi.createFolder(name)
    folders.value.push(folder)
    return folder
  }

  async function deleteFolder(name: string) {
    await galleryApi.deleteFolder(name)
    folders.value = folders.value.filter(f => f.name !== name)
    if (currentFolder.value === name) {
      currentFolder.value = null
    }
  }

  async function deleteImage(id: string) {
    await galleryApi.deleteImage(id)
    images.value = images.value.filter(img => img.id !== id)
    if (selectedImageId.value === id) {
      selectedImageId.value = null
    }
  }

  async function moveToFolder(imageId: string, folderPath: string) {
    await galleryApi.moveToFolder(imageId, folderPath)
    const image = images.value.find(img => img.id === imageId)
    if (image) {
      image.folder = folderPath
    }
  }

  function setCurrentFolder(folder: string | null) {
    currentFolder.value = folder
  }

  function selectImage(imageId: string | null) {
    selectedImageId.value = imageId
  }

  function setViewMode(mode: 'grid' | 'list') {
    viewMode.value = mode
  }

  function setSortBy(by: 'date' | 'name') {
    sortBy.value = by
  }

  function toggleSortOrder() {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  }

  function clear() {
    images.value = []
    folders.value = []
    currentFolder.value = null
    selectedImageId.value = null
  }

  return {
    images,
    folders,
    isLoading,
    currentFolder,
    selectedImageId,
    viewMode,
    sortBy,
    sortOrder,
    filteredImages,
    selectedImage,
    loadImages,
    loadFolders,
    createFolder,
    deleteFolder,
    deleteImage,
    moveToFolder,
    setCurrentFolder,
    selectImage,
    setViewMode,
    setSortBy,
    toggleSortOrder,
    clear
  }
})
