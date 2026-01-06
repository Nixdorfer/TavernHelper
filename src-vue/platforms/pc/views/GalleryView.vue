<template>
  <div class="gallery-content"
       @dragover.prevent="handleGalleryDragOver"
       @dragleave="handleGalleryDragLeave"
       @drop.prevent="handleGalleryDrop"
       @mousedown="handleSelectionStart"
       @mousemove="handleSelectionMove"
       @mouseup="handleSelectionEnd"
       :class="{ 'drag-over': galleryDragOver }">
    <div class="gallery-header-bar">
      <div class="gallery-breadcrumb">
        <span
          class="breadcrumb-item"
          :class="{ active: !currentFolder, 'drag-over': breadcrumbDragOver === 'root' }"
          @click="navigateToFolder('')"
          @dragover.prevent="onBreadcrumbDragOver($event, 'root')"
          @dragleave="onBreadcrumbDragLeave"
          @drop.prevent="onDropToBreadcrumb($event, '')"
        >全部</span>
        <template v-if="currentFolder">
          <span class="breadcrumb-sep">/</span>
          <span class="breadcrumb-item active" @dblclick="startRenamingFolder(currentFolder)">
            <template v-if="renamingFolder === currentFolder">
              <input
                v-model="renamingFolderName"
                class="inline-edit-input"
                @blur="finishRenamingFolder"
                @keyup.enter="finishRenamingFolder"
                @keyup.escape="cancelRenamingFolder"
                ref="folderRenameInputRef"
              />
            </template>
            <template v-else>{{ currentFolder }}</template>
          </span>
        </template>
      </div>
      <div class="gallery-actions">
        <button class="btn-gallery-action" @click="createFolder" title="新建文件夹">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/>
            <path d="M12 11v6M9 14h6"/>
          </svg>
        </button>
        <button class="btn-gallery-action" @click="addImage" title="添加图片">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 5v14M5 12h14"/>
          </svg>
        </button>
        <button v-if="selectedImages.length > 0" class="btn-gallery-action btn-danger" @click="deleteSelectedImages" title="删除选中">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
          </svg>
          <span>({{ selectedImages.length }})</span>
        </button>
      </div>
    </div>
    <div v-if="loading" class="gallery-loading-area">
      <span class="loading-spinner"></span>
      <span>加载中...</span>
    </div>
    <div v-else class="gallery-main-area" @contextmenu.prevent="showEmptyContextMenu($event)">
      <div v-if="displayImages.length === 0 && folders.length === 0 && !currentFolder" class="gallery-empty-hint">
        <span>暂无图片，点击右上角按钮添加或拖拽图片到此处</span>
      </div>
      <div v-else class="gallery-images-grid">
        <template v-if="!currentFolder">
          <div
            v-for="folder in folders"
            :key="'folder-' + folder.name"
            :class="['gallery-item-container', { selected: selectedFolders.includes(folder.name), clicked: clickedFolder === folder.name }]"
            :data-folder="folder.name"
          >
            <div
              :class="['gallery-folder-card', { 'drag-over': dragOverFolder === folder.name }]"
              @click="navigateToFolder(folder.name)"
              @contextmenu.prevent.stop="showFolderContextMenu($event, folder)"
              @dragover.prevent="onFolderDragOver($event, folder.name)"
              @dragleave="onFolderDragLeave"
              @drop.prevent="onDropToFolder($event, folder.name)"
            >
              <div class="folder-icon-wrapper">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/>
                </svg>
              </div>
            </div>
            <div :class="['gallery-item-name', { selected: selectedFolders.includes(folder.name) }]" @dblclick.stop="startRenamingFolder(folder.name)">
              <template v-if="renamingFolder === folder.name">
                <input
                  v-model="renamingFolderName"
                  class="inline-edit-input"
                  @click.stop
                  @blur="finishRenamingFolder"
                  @keyup.enter="finishRenamingFolder"
                  @keyup.escape="cancelRenamingFolder"
                  ref="folderRenameInputRef"
                />
              </template>
              <template v-else>{{ folder.name }}</template>
            </div>
          </div>
        </template>
        <div
          v-for="img in displayImages"
          :key="img.id"
          :class="['gallery-item-container', { selected: selectedImages.includes(img.id) }]"
          :data-image="img.id"
        >
          <div
            :class="['gallery-image-card', { 'no-url': !img.remoteUrl }]"
            @click="handleImageClick(img, $event)"
            @contextmenu.prevent.stop="showImageContextMenu($event, img)"
            draggable="true"
            @dragstart="onImageDragStart($event, img)"
            @dragend="onImageDragEnd"
          >
            <img :src="getImageSrc(img)" :alt="img.fileName" @error="handleImageError" draggable="false" />
            <div v-if="manageMode" class="select-overlay" :class="{ checked: selectedImages.includes(img.id) }">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                <path d="M5 12l5 5L20 7"/>
              </svg>
            </div>
          </div>
          <div :class="['gallery-item-name', { selected: selectedImages.includes(img.id) }]" @dblclick.stop="startRenamingImage(img)">
            <template v-if="renamingImage === img.id">
              <input
                v-model="renamingImageName"
                class="inline-edit-input"
                @blur="finishRenamingImage"
                @keyup.enter="finishRenamingImage"
                @keyup.escape="cancelRenamingImage"
                ref="imageRenameInputRef"
              />
            </template>
            <template v-else>{{ getFileNameWithoutExt(img.fileName) }}</template>
          </div>
        </div>
      </div>
    </div>
    <div v-if="galleryDragOver" class="gallery-drop-overlay">
      <div class="drop-hint">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12"/>
        </svg>
        <span>释放以添加图片</span>
      </div>
    </div>
    <Teleport to="body">
      <Transition name="gallery-modal">
        <div v-if="previewImage" class="modal-overlay gallery-modal-overlay" @mousedown.self="closePreview">
          <div class="gallery-preview-modal">
            <div class="preview-modal-header">
              <input
                v-model="previewImage.fileName"
                class="preview-name-input"
                @blur="savePreviewImageName"
                @keyup.enter="($event.target as HTMLInputElement).blur()"
              />
              <button class="btn btn-icon" @click="closePreview">×</button>
            </div>
            <div class="preview-modal-body">
              <img :src="getImageSrc(previewImage)" :alt="previewImage.fileName" />
            </div>
            <div class="preview-modal-footer">
              <div class="preview-actions">
                <button v-if="previewImage.remoteUrl" class="btn-preview btn-preview-success" @click="copyText(previewImage.remoteUrl)">复制</button>
                <button class="btn-preview btn-preview-danger" @click="deleteSingleImage(previewImage.id)">删除</button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="showCreateFolderModal" class="modal-overlay" @mousedown.self="cancelCreateFolder">
          <div class="modal modal-sm">
            <div class="modal-header">
              <h3>新建文件夹</h3>
              <button class="btn btn-icon" @click="cancelCreateFolder">×</button>
            </div>
            <div class="modal-body">
              <div class="form-group">
                <label class="form-label">文件夹名称</label>
                <input
                  ref="newFolderInputRef"
                  v-model="newFolderName"
                  class="form-input"
                  placeholder="请输入文件夹名称"
                  @keyup.enter="confirmCreateFolder"
                  @keyup.escape="cancelCreateFolder"
                />
              </div>
            </div>
            <div class="modal-footer">
              <button class="btn btn-secondary" @click="cancelCreateFolder">取消</button>
              <button class="btn btn-primary" @click="confirmCreateFolder" :disabled="!newFolderName.trim()">创建</button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
    <ContextMenu ref="contextMenuRef" :items="contextMenuItems" />
    <div
      v-if="isSelecting"
      class="selection-box"
      :style="selectionBoxStyle"
    ></div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, nextTick, reactive } from 'vue'
import { useAuthStore, useNotificationStore, useGalleryStore, useConfirmStore } from '@/stores'
import { galleryApi } from '@/api/modules/gallery'
import ContextMenu, { type ContextMenuItem } from '@/components/common/ContextMenu.vue'
interface GalleryImage {
  id: string
  fileName: string
  remoteUrl?: string
  folderPath?: string
  isValid?: boolean
}
interface GalleryFolder {
  name: string
  path: string
  createdAt: string
}
const authStore = useAuthStore()
const notificationStore = useNotificationStore()
const galleryStore = useGalleryStore()
const confirmStore = useConfirmStore()
const token = computed(() => authStore.token)
const contextMenuRef = ref<InstanceType<typeof ContextMenu> | null>(null)
const contextMenuItems = ref<ContextMenuItem[]>([])
const isSelecting = ref(false)
const selectionStart = ref({ x: 0, y: 0 })
const selectionCurrent = ref({ x: 0, y: 0 })
const selectionBoxStyle = computed(() => {
  const left = Math.min(selectionStart.value.x, selectionCurrent.value.x)
  const top = Math.min(selectionStart.value.y, selectionCurrent.value.y)
  const width = Math.abs(selectionCurrent.value.x - selectionStart.value.x)
  const height = Math.abs(selectionCurrent.value.y - selectionStart.value.y)
  return {
    left: left + 'px',
    top: top + 'px',
    width: width + 'px',
    height: height + 'px'
  }
})
const images = ref<GalleryImage[]>([])
const folders = ref<GalleryFolder[]>([])
const loading = ref(false)
const manageMode = ref(false)
const selectedImages = ref<string[]>([])
const selectedFolders = ref<string[]>([])
const clickedFolder = ref('')
const imageCache = reactive<Record<string, string>>({})
const previewImage = ref<GalleryImage | null>(null)
const currentFolder = ref('')
const galleryDragOver = ref(false)
const imageDragging = ref(false)
const dragOverFolder = ref<string | null>(null)
const breadcrumbDragOver = ref<string | null>(null)
const renamingImage = ref<string | null>(null)
const renamingImageName = ref('')
const _renamingImageExt = ref('')
const renamingFolder = ref<string | null>(null)
const renamingFolderName = ref('')
const showCreateFolderModal = ref(false)
const newFolderName = ref('')
const newFolderInputRef = ref<HTMLInputElement | null>(null)
const folderRenameInputRef = ref<HTMLInputElement | HTMLInputElement[] | null>(null)
const imageRenameInputRef = ref<HTMLInputElement | HTMLInputElement[] | null>(null)
const displayImages = computed(() => {
  return images.value.filter(img => {
    if (img.isValid === false) return false
    if (currentFolder.value) {
      return img.folderPath === currentFolder.value
    }
    return !img.folderPath
  })
})
async function loadData() {
  loading.value = true
  try {
    const [imgs, flds] = await Promise.all([
      galleryApi.getImages(),
      galleryApi.getFolders()
    ])
    images.value = imgs || []
    folders.value = flds || []
    for (const img of images.value) {
      if (img.isValid !== false && !imageCache[img.id]) {
        loadImageBase64(img)
      }
    }
  } catch (err) {
    console.error('加载图库失败:', err)
  } finally {
    loading.value = false
  }
}
async function loadImageBase64(img: GalleryImage) {
  try {
    const base64 = await galleryApi.readImageAsBase64(img.id)
    imageCache[img.id] = base64
  } catch (err) {
    console.error('加载图片失败:', img.id, err)
  }
}
function getImageSrc(img: GalleryImage): string {
  if (imageCache[img.id]) {
    return imageCache[img.id]
  }
  if (img.remoteUrl) {
    return img.remoteUrl
  }
  return ''
}
function handleImageError(e: Event) {
  const target = e.target as HTMLImageElement
  target.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect fill="%23333" width="100" height="100"/><text fill="%23666" x="50" y="55" text-anchor="middle" font-size="12">加载失败</text></svg>'
}
function handleImageClick(img: GalleryImage, event: MouseEvent) {
  if (event.ctrlKey || event.metaKey || manageMode.value) {
    const idx = selectedImages.value.indexOf(img.id)
    if (idx >= 0) {
      selectedImages.value.splice(idx, 1)
    } else {
      selectedImages.value.push(img.id)
    }
  } else {
    selectedImages.value = []
    selectedFolders.value = []
    previewImage.value = { ...img }
  }
}
function handleFolderClick(folder: GalleryFolder, event: MouseEvent) {
  if (event.ctrlKey || event.metaKey) {
    const idx = selectedFolders.value.indexOf(folder.name)
    if (idx >= 0) {
      selectedFolders.value.splice(idx, 1)
    } else {
      selectedFolders.value.push(folder.name)
    }
  } else {
    clickedFolder.value = clickedFolder.value === folder.name ? '' : folder.name
  }
}
function onFolderDragOver(e: DragEvent, folderName: string) {
  if (!imageDragging.value) return
  dragOverFolder.value = folderName
}
function onFolderDragLeave() {
  dragOverFolder.value = null
}
async function onDropToFolder(e: DragEvent, folderName: string) {
  dragOverFolder.value = null
  if (!imageDragging.value) return
  const imageId = e.dataTransfer?.getData('text/plain')
  if (!imageId) return
  try {
    await galleryApi.moveImageToFolder(imageId, folderName)
    const img = images.value.find(i => i.id === imageId)
    if (img) {
      img.folderPath = folderName
    }
    notificationStore.showNotification('移动成功', 'success')
  } catch (err) {
    console.error('移动图片失败:', err)
    notificationStore.showNotification('移动失败', 'error')
  }
}
function onBreadcrumbDragOver(e: DragEvent, target: string) {
  if (!imageDragging.value) return
  breadcrumbDragOver.value = target
}
function onBreadcrumbDragLeave() {
  breadcrumbDragOver.value = null
}
async function onDropToBreadcrumb(e: DragEvent, folderPath: string) {
  breadcrumbDragOver.value = null
  if (!imageDragging.value) return
  const imageId = e.dataTransfer?.getData('text/plain')
  if (!imageId) return
  try {
    await galleryApi.moveImageToFolder(imageId, folderPath)
    const img = images.value.find(i => i.id === imageId)
    if (img) {
      img.folderPath = folderPath
    }
    notificationStore.showNotification('移动成功', 'success')
  } catch (err) {
    console.error('移动图片失败:', err)
    notificationStore.showNotification('移动失败', 'error')
  }
}
function closePreview() {
  previewImage.value = null
}
async function savePreviewImageName() {
  if (!previewImage.value) return
  try {
    await galleryApi.renameImage(previewImage.value.id, previewImage.value.fileName)
    const img = images.value.find(i => i.id === previewImage.value!.id)
    if (img) img.fileName = previewImage.value.fileName
  } catch (err) {
    console.error('重命名失败:', err)
  }
}
async function copyText(text: string) {
  try {
    await navigator.clipboard.writeText(text)
    notificationStore.showNotification('已复制到剪贴板', 'success')
  } catch {
    notificationStore.showNotification('复制失败', 'error')
  }
}
async function addImage() {
  try {
    const result = currentFolder.value
      ? await galleryApi.selectAndAddToFolder(currentFolder.value)
      : await galleryApi.selectAndAdd()
    if (result) {
      const exists = images.value.find(i => i.id === result.id)
      if (!exists) {
        images.value.unshift(result)
      }
      loadImageBase64(result)
    }
  } catch (err) {
    console.error('添加图片失败:', err)
  }
}
async function deleteSelectedImages() {
  if (selectedImages.value.length === 0) return
  const confirmed = await confirmStore.show({
    title: '确认删除',
    message: `确定要删除选中的 ${selectedImages.value.length} 张图片吗？此操作不可撤销。`,
    confirmText: '删除',
    type: 'danger'
  })
  if (!confirmed) return
  try {
    await galleryApi.deleteImages(selectedImages.value)
    images.value = images.value.filter(img => !selectedImages.value.includes(img.id))
    for (const id of selectedImages.value) {
      delete imageCache[id]
    }
    selectedImages.value = []
    notificationStore.showNotification('删除成功', 'success')
  } catch (err) {
    console.error('删除图片失败:', err)
    notificationStore.showNotification('删除失败', 'error')
  }
}
async function deleteSingleImage(id: string) {
  const img = images.value.find(i => i.id === id)
  const confirmed = await confirmStore.confirmDelete(img?.fileName || '图片')
  if (!confirmed) return
  try {
    await galleryApi.deleteImage(id)
    images.value = images.value.filter(img => img.id !== id)
    delete imageCache[id]
    previewImage.value = null
    notificationStore.showNotification('删除成功', 'success')
  } catch (err) {
    console.error('删除图片失败:', err)
    notificationStore.showNotification('删除失败', 'error')
  }
}
function navigateToFolder(folderName: string) {
  currentFolder.value = folderName
  selectedImages.value = []
}
function createFolder() {
  newFolderName.value = ''
  showCreateFolderModal.value = true
  nextTick(() => {
    if (newFolderInputRef.value) newFolderInputRef.value.focus()
  })
}
async function confirmCreateFolder() {
  if (!newFolderName.value || !newFolderName.value.trim()) {
    showCreateFolderModal.value = false
    return
  }
  try {
    await galleryApi.createFolder(newFolderName.value.trim())
    folders.value.push({ name: newFolderName.value.trim(), path: newFolderName.value.trim(), createdAt: new Date().toISOString() })
    notificationStore.showNotification('文件夹创建成功', 'success')
  } catch (err) {
    console.error('创建文件夹失败:', err)
    notificationStore.showNotification('创建文件夹失败', 'error')
  }
  showCreateFolderModal.value = false
  newFolderName.value = ''
}
function cancelCreateFolder() {
  showCreateFolderModal.value = false
  newFolderName.value = ''
}
function startRenamingFolder(name: string) {
  renamingFolder.value = name
  renamingFolderName.value = name
  nextTick(() => {
    const input = folderRenameInputRef.value
    if (input) {
      const el = Array.isArray(input) ? input[0] : input
      if (el) {
        el.focus()
        el.select()
      }
    }
  })
}
async function finishRenamingFolder() {
  if (!renamingFolder.value || !renamingFolderName.value.trim()) {
    cancelRenamingFolder()
    return
  }
  if (renamingFolderName.value === renamingFolder.value) {
    cancelRenamingFolder()
    return
  }
  try {
    await galleryApi.renameFolder(renamingFolder.value, renamingFolderName.value.trim())
    const folder = folders.value.find(f => f.name === renamingFolder.value)
    if (folder) {
      folder.name = renamingFolderName.value.trim()
      folder.path = renamingFolderName.value.trim()
    }
    if (currentFolder.value === renamingFolder.value) {
      currentFolder.value = renamingFolderName.value.trim()
    }
    images.value.forEach(img => {
      if (img.folderPath === renamingFolder.value) {
        img.folderPath = renamingFolderName.value.trim()
      }
    })
  } catch (err) {
    console.error('重命名文件夹失败:', err)
    notificationStore.showNotification('重命名文件夹失败', 'error')
  }
  renamingFolder.value = null
  renamingFolderName.value = ''
}
function cancelRenamingFolder() {
  renamingFolder.value = null
  renamingFolderName.value = ''
}
function startRenamingImage(img: GalleryImage) {
  renamingImage.value = img.id
  renamingImageName.value = getFileNameWithoutExt(img.fileName)
  _renamingImageExt.value = getFileExt(img.fileName)
  nextTick(() => {
    const input = imageRenameInputRef.value
    if (input) {
      const el = Array.isArray(input) ? input[0] : input
      if (el) {
        el.focus()
        el.select()
      }
    }
  })
}
function getFileExt(fileName: string): string {
  if (!fileName) return ''
  const lastDot = fileName.lastIndexOf('.')
  if (lastDot === -1) return ''
  return fileName.substring(lastDot)
}
async function finishRenamingImage() {
  if (!renamingImage.value || !renamingImageName.value.trim()) {
    cancelRenamingImage()
    return
  }
  try {
    const newName = renamingImageName.value.trim() + (_renamingImageExt.value || '')
    await galleryApi.renameImage(renamingImage.value, newName)
    const img = images.value.find(i => i.id === renamingImage.value)
    if (img) img.fileName = newName
  } catch (err) {
    console.error('重命名图片失败:', err)
    notificationStore.showNotification('重命名图片失败', 'error')
  }
  renamingImage.value = null
  renamingImageName.value = ''
  _renamingImageExt.value = ''
}
function cancelRenamingImage() {
  renamingImage.value = null
  renamingImageName.value = ''
}
function handleGalleryDragOver(e: DragEvent) {
  if (imageDragging.value) return
  const types = e.dataTransfer?.types
  if (!types || !Array.from(types).includes('Files')) return
  galleryDragOver.value = true
}
function handleGalleryDragLeave(e: DragEvent) {
  if ((e.currentTarget as HTMLElement).contains(e.relatedTarget as HTMLElement)) return
  galleryDragOver.value = false
}
function onImageDragStart(e: DragEvent, img: GalleryImage) {
  imageDragging.value = true
  e.dataTransfer!.setData('text/plain', img.id)
  e.dataTransfer!.effectAllowed = 'move'
}
function onImageDragEnd() {
  imageDragging.value = false
}
function getFileNameWithoutExt(fileName: string): string {
  if (!fileName) return ''
  const lastDot = fileName.lastIndexOf('.')
  if (lastDot === -1) return fileName
  return fileName.substring(0, lastDot)
}
async function handleGalleryDrop(e: DragEvent) {
  galleryDragOver.value = false
  const types = e.dataTransfer?.types
  if (!types || !Array.from(types).includes('Files')) return
  const files = e.dataTransfer?.files
  if (!files || files.length === 0) return
  for (const file of files) {
    if (!file.type.startsWith('image/')) continue
    try {
      const base64 = await readFileAsBase64(file)
      const result = await galleryApi.addFromBase64ToFolder(base64, file.name, currentFolder.value || '')
      if (result) {
        const exists = images.value.find(i => i.id === result.id)
        if (!exists) {
          images.value.unshift(result)
          imageCache[result.id] = base64
        } else if (result.folderPath && result.folderPath !== currentFolder.value) {
          navigateToFolder(result.folderPath)
        }
        loadImageBase64(result)
      }
    } catch (err) {
      console.error('添加图片失败:', err)
    }
  }
}
function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}
onMounted(() => {
  loadData()
})
function showEmptyContextMenu(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (target.closest('.gallery-image-card') || target.closest('.gallery-folder-card')) return
  contextMenuItems.value = [
    {
      label: '导入图片',
      icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4"/><polyline points="17 8 12 3 7 8"/><line x1="12" y1="3" x2="12" y2="15"/></svg>',
      action: () => addImage()
    },
    {
      label: '新建文件夹',
      icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 19a2 2 0 01-2 2H4a2 2 0 01-2-2V5a2 2 0 012-2h5l2 3h9a2 2 0 012 2z"/><path d="M12 11v6M9 14h6"/></svg>',
      action: () => createFolder()
    }
  ]
  contextMenuRef.value?.show(e.clientX, e.clientY)
}
function showFolderContextMenu(e: MouseEvent, folder: GalleryFolder) {
  contextMenuItems.value = [
    {
      label: '重命名',
      icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>',
      action: () => startRenamingFolder(folder.name)
    },
    {
      label: '删除',
      icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>',
      danger: true,
      action: () => deleteFolder(folder.name)
    }
  ]
  contextMenuRef.value?.show(e.clientX, e.clientY)
}
function showImageContextMenu(e: MouseEvent, img: GalleryImage) {
  contextMenuItems.value = [
    {
      label: '重命名',
      icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>',
      action: () => startRenamingImage(img)
    },
    {
      label: '删除',
      icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 6h18M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/></svg>',
      danger: true,
      action: () => deleteSingleImage(img.id)
    }
  ]
  contextMenuRef.value?.show(e.clientX, e.clientY)
}
async function deleteFolder(name: string) {
  const confirmed = await confirmStore.confirmDelete(name)
  if (!confirmed) return
  try {
    await galleryApi.deleteFolder(name)
    folders.value = folders.value.filter(f => f.name !== name)
    images.value = images.value.filter(img => img.folderPath !== name)
    notificationStore.showNotification('删除成功', 'success')
  } catch (err) {
    console.error('删除文件夹失败:', err)
    notificationStore.showNotification('删除失败', 'error')
  }
}
function handleSelectionStart(e: MouseEvent) {
  const target = e.target as HTMLElement
  if (target.closest('.gallery-image-card') || target.closest('.gallery-folder-card') || target.closest('.gallery-header-bar') || target.closest('.btn')) return
  if (e.button !== 0) return
  isSelecting.value = true
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  selectionStart.value = { x: e.clientX - rect.left, y: e.clientY - rect.top }
  selectionCurrent.value = { ...selectionStart.value }
  selectedImages.value = []
  selectedFolders.value = []
}
function handleSelectionMove(e: MouseEvent) {
  if (!isSelecting.value) return
  const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
  selectionCurrent.value = { x: e.clientX - rect.left, y: e.clientY - rect.top }
  updateSelectionFromBox(e.currentTarget as HTMLElement)
}
function handleSelectionEnd() {
  isSelecting.value = false
}
function updateSelectionFromBox(container: HTMLElement) {
  const boxLeft = Math.min(selectionStart.value.x, selectionCurrent.value.x)
  const boxTop = Math.min(selectionStart.value.y, selectionCurrent.value.y)
  const boxRight = Math.max(selectionStart.value.x, selectionCurrent.value.x)
  const boxBottom = Math.max(selectionStart.value.y, selectionCurrent.value.y)
  const newSelectedImages: string[] = []
  const newSelectedFolders: string[] = []
  const items = container.querySelectorAll('.gallery-item-container')
  items.forEach(item => {
    const rect = item.getBoundingClientRect()
    const containerRect = container.getBoundingClientRect()
    const itemLeft = rect.left - containerRect.left
    const itemTop = rect.top - containerRect.top
    const itemRight = itemLeft + rect.width
    const itemBottom = itemTop + rect.height
    const intersects = !(itemRight < boxLeft || itemLeft > boxRight || itemBottom < boxTop || itemTop > boxBottom)
    if (intersects) {
      const folderId = item.querySelector('.gallery-folder-card') ? item.getAttribute('data-folder') : null
      const imageId = item.querySelector('.gallery-image-card') ? item.getAttribute('data-image') : null
      if (folderId) {
        newSelectedFolders.push(folderId)
      } else if (imageId) {
        newSelectedImages.push(imageId)
      }
    }
  })
  selectedImages.value = newSelectedImages
  selectedFolders.value = newSelectedFolders
}
</script>
<style scoped>
.gallery-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  overflow: hidden;
  position: relative;
}
.gallery-content.drag-over {
  background: rgba(34, 197, 94, 0.05);
}
.gallery-header-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
}
.gallery-breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
}
.breadcrumb-item {
  padding: 4px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-secondary);
  transition: all 0.2s;
}
.breadcrumb-item:hover {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.breadcrumb-item.active {
  color: var(--text-primary);
  font-weight: 500;
}
.breadcrumb-item.drag-over {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.4);
}
.breadcrumb-sep {
  color: var(--text-secondary);
}
.gallery-actions {
  display: flex;
  gap: 8px;
}
.btn-gallery-action {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 10px;
  background: var(--bg-tertiary);
  border: none;
  border-radius: 6px;
  cursor: pointer;
  color: var(--text-primary);
  transition: all 0.2s;
}
.btn-gallery-action:hover {
  background: var(--bg-primary);
}
.btn-gallery-action.btn-danger {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}
.btn-gallery-action.btn-danger:hover {
  background: #ef4444;
  color: white;
}
.btn-gallery-action svg {
  width: 18px;
  height: 18px;
}
.gallery-loading-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: var(--text-secondary);
}
.gallery-main-area {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}
.gallery-empty-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-secondary);
}
.gallery-images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;
}
.gallery-item-container {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.gallery-item-container.selected .gallery-image-card,
.gallery-item-container.selected .gallery-folder-card {
  box-shadow: 0 0 0 3px var(--primary-color);
}
.gallery-item-container.clicked .gallery-folder-card {
  background: var(--bg-hover);
}
.gallery-folder-card,
.gallery-image-card {
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  position: relative;
}
.gallery-folder-card:hover,
.gallery-image-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}
.gallery-folder-card.drag-over {
  background: rgba(34, 197, 94, 0.15);
  border-color: #22c55e;
}
.folder-icon-wrapper {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}
.folder-icon-wrapper svg {
  width: 48px;
  height: 48px;
  color: var(--text-secondary);
}
.gallery-image-card img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.gallery-image-card.no-url {
  opacity: 0.7;
}
.select-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s;
}
.gallery-image-card:hover .select-overlay {
  opacity: 1;
}
.select-overlay.checked {
  opacity: 1;
  background: rgba(34, 197, 94, 0.5);
}
.select-overlay svg {
  width: 32px;
  height: 32px;
  color: white;
}
.gallery-item-name {
  font-size: 12px;
  text-align: center;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 4px;
}
.gallery-item-name.selected {
  color: var(--primary-color);
  font-weight: 500;
}
.inline-edit-input {
  width: 100%;
  padding: 2px 6px;
  font-size: 12px;
  border: 1px solid var(--primary-color);
  border-radius: 4px;
  background: var(--bg-tertiary);
  color: var(--text-primary);
  text-align: center;
  resize: none;
}
.inline-edit-input:focus {
  outline: none;
}
.gallery-drop-overlay {
  position: absolute;
  inset: 0;
  background: rgba(34, 197, 94, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3px dashed #22c55e;
  border-radius: 8px;
  z-index: 10;
}
.drop-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #22c55e;
}
.drop-hint svg {
  width: 48px;
  height: 48px;
}
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.gallery-modal-overlay {
  background: rgba(0, 0, 0, 0.85);
}
.gallery-preview-modal {
  max-width: 90vw;
  max-height: 90vh;
  background: var(--bg-secondary);
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
.preview-modal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
  gap: 12px;
}
.preview-name-input {
  flex: 1;
  padding: 6px 10px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 14px;
  resize: none;
}
.preview-name-input:focus {
  outline: none;
  border-color: var(--primary-color);
}
.preview-modal-body {
  flex: 1;
  overflow: auto;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px;
}
.preview-modal-body img {
  max-width: 100%;
  max-height: 70vh;
  object-fit: contain;
}
.preview-modal-footer {
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
}
.preview-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
}
.btn-preview {
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 13px;
  cursor: pointer;
  border: none;
  transition: all 0.2s;
}
.btn-preview-primary {
  background: var(--primary-color);
  color: white;
}
.btn-preview-success {
  background: #22c55e;
  color: white;
}
.btn-preview-danger {
  background: transparent;
  color: #ef4444;
  border: 1px solid #ef4444;
}
.btn-preview-danger:hover {
  background: #ef4444;
  color: white;
}
.modal {
  background: var(--bg-secondary);
  border-radius: 12px;
  min-width: 320px;
  max-width: 90vw;
}
.modal-sm {
  min-width: 280px;
  max-width: 320px;
  width: 320px;
}
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
}
.modal-header h3 {
  margin: 0;
  font-size: 16px;
}
.modal-body {
  padding: 16px;
}
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  border-top: 1px solid var(--border-color);
}
.form-group {
  margin-bottom: 12px;
}
.form-label {
  display: block;
  font-size: 13px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}
.form-input {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  color: var(--text-primary);
  font-size: 14px;
  resize: none;
}
.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
}
.btn {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  border: none;
}
.btn-icon {
  padding: 4px 8px;
  background: transparent;
}
.btn-primary {
  background: var(--primary-color);
  color: white;
}
.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
.loading-spinner-xs {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.gallery-modal-enter-active,
.gallery-modal-leave-active {
  transition: opacity 0.2s;
}
.gallery-modal-enter-from,
.gallery-modal-leave-to {
  opacity: 0;
}
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.selection-box {
  position: absolute;
  border: 2px solid var(--primary-color);
  background: rgba(var(--primary-color-rgb, 34, 197, 94), 0.15);
  border-radius: 4px;
  pointer-events: none;
  z-index: 100;
}
</style>
