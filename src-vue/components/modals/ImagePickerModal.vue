<template>
  <Teleport to="body">
    <Transition name="modal">
      <div class="modal-overlay" @mousedown.self="close">
        <div
          class="image-picker-modal"
          @dragenter.prevent="handleDragEnter"
          @dragover.prevent="handleDragOver"
          @dragleave.prevent="handleDragLeave"
          @drop.prevent="handleDrop"
        >
          <div v-if="isDragging" class="drag-overlay">
            <div class="drag-hint">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M21 15v4a2 2 0 01-2 2H5a2 2 0 01-2-2v-4M17 8l-5-5-5 5M12 3v12"/>
              </svg>
              <span>拖放图片到此处</span>
            </div>
          </div>
          <div class="picker-header">
            <h3>选择图片</h3>
          </div>
          <div class="picker-body">
            <div v-if="loading" class="picker-loading">
              <span class="loading-spinner"></span>
              <span>加载中...</span>
            </div>
            <div v-else-if="images.length === 0" class="picker-empty">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="18" height="18" rx="2"/>
                <circle cx="8.5" cy="8.5" r="1.5"/>
                <path d="M21 15l-5-5L5 21"/>
              </svg>
              <span>暂无图片</span>
              <span class="empty-hint">点击"新建"添加图片或直接拖入</span>
            </div>
            <div v-else class="picker-grid">
              <div
                v-for="img in images"
                :key="img.id"
                :class="['picker-item', { selected: selectedImage?.id === img.id }]"
                @click="selectImage(img)"
              >
                <img :src="getImageSrc(img)" :alt="img.fileName" @error="handleImageError" />
                <div class="item-check">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
                    <path d="M5 12l5 5L20 7"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>
          <div class="picker-footer">
            <button class="btn btn-secondary" @click="close">关闭</button>
            <button class="btn btn-secondary" @click="openFilePicker">新建</button>
            <button class="btn btn-primary" :disabled="!selectedImage" @click="confirm">确定</button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { galleryApi } from '@/api'
interface GalleryImage {
  id: string
  fileName: string
  localPath?: string
  remoteUrl?: string
}
const emit = defineEmits<{
  close: []
  confirm: [image: GalleryImage]
}>()
const images = ref<GalleryImage[]>([])
const loading = ref(false)
const selectedImage = ref<GalleryImage | null>(null)
const imageCache = reactive<Record<string, string>>({})
const isDragging = ref(false)
const dragCounter = ref(0)
function close() {
  emit('close')
}
async function loadImages() {
  loading.value = true
  try {
    const result = await galleryApi.getImages()
    images.value = result || []
    for (const img of images.value) {
      if (!imageCache[img.id]) {
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
function handleImageError(event: Event) {
  const target = event.target as HTMLImageElement
  target.src = 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect fill="%23333" width="100" height="100"/><text fill="%23666" x="50" y="55" text-anchor="middle" font-size="12">加载失败</text></svg>'
}
function selectImage(img: GalleryImage) {
  if (selectedImage.value?.id === img.id) {
    selectedImage.value = null
  } else {
    selectedImage.value = img
  }
}
async function openFilePicker() {
  try {
    const result = await galleryApi.selectAndAddImage()
    if (result) {
      const exists = images.value.find(i => i.id === result.id)
      if (!exists) {
        images.value.unshift(result)
      }
      await loadImageBase64(result)
      selectedImage.value = result
    }
  } catch (err) {
    console.error('添加图片失败:', err)
  }
}
function confirm() {
  if (selectedImage.value) {
    emit('confirm', selectedImage.value)
    close()
  }
}
function handleDragEnter(e: DragEvent) {
  dragCounter.value++
  if (e.dataTransfer?.types.includes('Files')) {
    isDragging.value = true
  }
}
function handleDragOver(e: DragEvent) {
  if (e.dataTransfer) {
    e.dataTransfer.dropEffect = 'copy'
  }
}
function handleDragLeave() {
  dragCounter.value--
  if (dragCounter.value === 0) {
    isDragging.value = false
  }
}
async function handleDrop(e: DragEvent) {
  isDragging.value = false
  dragCounter.value = 0
  const files = Array.from(e.dataTransfer?.files || []).filter(f => f.type.startsWith('image/'))
  if (files.length === 0) return
  for (const file of files) {
    await addImageFromFile(file)
  }
}
async function addImageFromFile(file: File) {
  try {
    const reader = new FileReader()
    const base64 = await new Promise<string>((resolve, reject) => {
      reader.onload = () => resolve(reader.result as string)
      reader.onerror = reject
      reader.readAsDataURL(file)
    })
    const result = await galleryApi.addImageFromBase64(base64, file.name)
    if (result) {
      const exists = images.value.find(i => i.id === result.id)
      if (!exists) {
        images.value.unshift(result)
      }
      imageCache[result.id] = base64
      selectedImage.value = result
    }
  } catch (err) {
    console.error('添加图片失败:', err)
  }
}
function handleKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    close()
  }
  if (e.key === 'Enter' && selectedImage.value) {
    e.preventDefault()
    confirm()
  }
}
onMounted(() => {
  loadImages()
  window.addEventListener('keydown', handleKeydown)
})
onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>
<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}
.image-picker-modal {
  position: relative;
  width: 600px;
  max-width: 90vw;
  height: 500px;
  max-height: 80vh;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.drag-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(var(--primary-rgb, 59, 130, 246), 0.15);
  border: 3px dashed var(--primary-color);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  animation: dragPulse 1s ease-in-out infinite;
}
@keyframes dragPulse {
  0%, 100% { background: rgba(var(--primary-rgb, 59, 130, 246), 0.15); }
  50% { background: rgba(var(--primary-rgb, 59, 130, 246), 0.25); }
}
.drag-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: var(--primary-color);
}
.drag-hint svg {
  width: 48px;
  height: 48px;
  animation: dragBounce 0.6s ease-in-out infinite;
}
@keyframes dragBounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-8px); }
}
.drag-hint span {
  font-size: 16px;
  font-weight: 500;
}
.picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
  background: var(--card-bg);
}
.picker-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}
.picker-body {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--bg-tertiary);
}
.picker-loading,
.picker-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  gap: 12px;
  color: var(--text-secondary);
}
.picker-empty svg {
  width: 64px;
  height: 64px;
  opacity: 0.5;
}
.picker-empty span {
  font-size: 14px;
}
.empty-hint {
  font-size: 12px;
  opacity: 0.7;
}
.picker-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 12px;
  background: var(--bg-tertiary);
  padding: 12px;
  border-radius: 8px;
}
.picker-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  background: var(--bg-secondary);
  transition: transform 0.15s, box-shadow 0.15s;
}
.picker-item:hover {
  transform: scale(1.03);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
.picker-item.selected {
  outline: 3px solid var(--primary-color);
  outline-offset: 2px;
}
.picker-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.item-check {
  position: absolute;
  top: 6px;
  right: 6px;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--primary-color);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transform: scale(0.8);
  transition: all 0.15s;
}
.picker-item.selected .item-check {
  opacity: 1;
  transform: scale(1);
}
.item-check svg {
  width: 14px;
  height: 14px;
  color: white;
}
.picker-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
  flex-shrink: 0;
  background: var(--card-bg);
}
.btn {
  padding: 8px 20px;
  font-size: 14px;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-secondary {
  background: var(--hover-bg);
  color: var(--text-color);
}
.btn-secondary:hover {
  background: var(--border-color);
}
.btn-primary {
  background: var(--primary-color);
  color: white;
}
.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}
.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.loading-spinner {
  width: 24px;
  height: 24px;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.modal-enter-active {
  transition: opacity 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-leave-active {
  transition: opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-enter-active .image-picker-modal {
  transition: transform 0.25s cubic-bezier(0.34, 1.56, 0.64, 1),
              opacity 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-leave-active .image-picker-modal {
  transition: transform 0.2s cubic-bezier(0.4, 0, 0.2, 1),
              opacity 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .image-picker-modal {
  transform: scale(0.9) translateY(10px);
  opacity: 0;
}
.modal-leave-to .image-picker-modal {
  transform: scale(0.95);
  opacity: 0;
}
</style>
