import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Project, ProjectInfo, TimelineNode } from '@/types'
import { projectApi } from '@/api/modules/project'

export const useWorldbookProjectStore = defineStore('worldbook-project', () => {
  const projects = ref<ProjectInfo[]>([])
  const currentProject = ref<Project | null>(null)
  const currentProjectName = ref('')
  const isLoading = ref(false)
  const isSaving = ref(false)

  const hasProject = computed(() => !!currentProject.value)

  const rootNodes = computed(() => {
    if (!currentProject.value?.timeline) return []
    return currentProject.value.timeline.filter(n => !n.parentId)
  })

  async function loadProjects() {
    isLoading.value = true
    try {
      projects.value = await projectApi.getAll() || []
    } finally {
      isLoading.value = false
    }
  }

  async function loadProject(fileName: string) {
    console.log('[Store.loadProject] 开始加载:', fileName)
    isLoading.value = true
    try {
      const loaded = await projectApi.load(fileName)
      console.log('[Store.loadProject] API返回:', loaded)
      if (loaded) {
        loaded.timeline = loaded.timeline || []
        loaded.worldBook = loaded.worldBook || []
      }
      currentProject.value = loaded
      currentProjectName.value = fileName
      console.log('[Store.loadProject] 加载完成')
    } catch (e) {
      console.error('[Store.loadProject] 加载失败:', e)
      throw e
    } finally {
      isLoading.value = false
    }
  }

  async function saveProject() {
    if (!currentProject.value || !currentProjectName.value) return
    isSaving.value = true
    try {
      await projectApi.save(currentProjectName.value, currentProject.value)
    } finally {
      isSaving.value = false
    }
  }

  async function deleteProject(fileName: string) {
    await projectApi.delete(fileName)
    projects.value = projects.value.filter(p => p.fileName !== fileName)
    if (currentProjectName.value === fileName) {
      currentProject.value = null
      currentProjectName.value = ''
    }
  }

  async function renameProject(oldName: string, newName: string) {
    await projectApi.rename(oldName, newName)
    const project = projects.value.find(p => p.fileName === oldName)
    if (project) {
      project.fileName = newName
      project.name = newName.replace(/\.[^/.]+$/, '')
    }
    if (currentProjectName.value === oldName) {
      currentProjectName.value = newName
    }
  }

  function setProject(project: Project | null) {
    currentProject.value = project
  }

  function updateTimeline(timeline: TimelineNode[]) {
    if (currentProject.value) {
      currentProject.value.timeline = timeline
    }
  }

  function clear() {
    currentProject.value = null
    currentProjectName.value = ''
  }

  return {
    projects,
    currentProject,
    currentProjectName,
    isLoading,
    isSaving,
    hasProject,
    rootNodes,
    loadProjects,
    loadProject,
    saveProject,
    deleteProject,
    renameProject,
    setProject,
    updateTimeline,
    clear
  }
})
