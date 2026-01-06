import { api } from '../index'
import type { GalleryImage } from '@/types'
export interface GalleryFolder {
  name: string
  path: string
  count: number
}
export const galleryApi = {
  getImages(): Promise<GalleryImage[]> {
    return api.gallery.getImages()
  },
  getFolders(): Promise<GalleryFolder[]> {
    return api.gallery.getFolders()
  },
  createFolder(name: string): Promise<GalleryFolder> {
    return api.gallery.createFolder(name)
  },
  deleteFolder(name: string) {
    return api.gallery.deleteFolder(name)
  },
  deleteImage(id: string) {
    return api.gallery.deleteImage(id)
  },
  deleteImages(ids: string[]) {
    return api.gallery.deleteImages(ids)
  },
  moveToFolder(imageId: string, folderPath: string) {
    return api.gallery.moveToFolder(imageId, folderPath)
  },
  moveImageToFolder(imageId: string, folderPath: string) {
    return api.gallery.moveImageToFolder(imageId, folderPath)
  },
  readAsBase64(id: string): Promise<string> {
    return api.gallery.readAsBase64(id)
  },
  readImageAsBase64(id: string): Promise<string> {
    return api.gallery.readAsBase64(id)
  },
  selectAndAddImage(): Promise<GalleryImage | null> {
    return api.gallery.selectAndAdd()
  },
  selectAndAdd(): Promise<GalleryImage | null> {
    return api.gallery.selectAndAdd()
  },
  selectAndAddToFolder(folderPath: string): Promise<GalleryImage | null> {
    return api.gallery.selectAndAddToFolder(folderPath)
  },
  addImageFromBase64(base64: string, fileName: string): Promise<GalleryImage | null> {
    return api.gallery.addFromBase64(base64, fileName)
  },
  addFromBase64ToFolder(base64: string, fileName: string, folderPath: string): Promise<GalleryImage | null> {
    return api.gallery.addFromBase64ToFolder(base64, fileName, folderPath)
  },
  renameImage(id: string, newName: string) {
    return api.gallery.renameImage(id, newName)
  },
  renameFolder(oldName: string, newName: string) {
    return api.gallery.renameFolder(oldName, newName)
  },
  updateImageUrl(id: string, url: string) {
    return api.gallery.updateImageUrl(id, url)
  }
}
