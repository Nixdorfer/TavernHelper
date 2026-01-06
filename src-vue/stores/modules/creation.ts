import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CreationProject {
  id: string
  name: string
  type: string
  content: string
  createdAt: string
  updatedAt: string
}

export const useCreationStore = defineStore('creation', () => {
  const projects = ref<CreationProject[]>([])
  const currentProjectId = ref<string | null>(null)
  const isLoading = ref(false)
  const isSaving = ref(false)
  const isDirty = ref(false)
  const editorContent = ref('')

  const currentProject = computed(() => {
    return projects.value.find(p => p.id === currentProjectId.value) || null
  })

  function setProjects(newProjects: CreationProject[]) {
    projects.value = newProjects
  }

  function addProject(project: CreationProject) {
    projects.value.push(project)
  }

  function updateProject(id: string, updates: Partial<CreationProject>) {
    const index = projects.value.findIndex(p => p.id === id)
    if (index > -1) {
      projects.value[index] = { ...projects.value[index], ...updates }
    }
  }

  function removeProject(id: string) {
    projects.value = projects.value.filter(p => p.id !== id)
    if (currentProjectId.value === id) {
      currentProjectId.value = null
    }
  }

  function selectProject(projectId: string | null) {
    currentProjectId.value = projectId
    if (projectId) {
      const project = projects.value.find(p => p.id === projectId)
      if (project) {
        editorContent.value = project.content
      }
    } else {
      editorContent.value = ''
    }
    isDirty.value = false
  }

  function setEditorContent(content: string) {
    editorContent.value = content
    isDirty.value = true
  }

  function markClean() {
    isDirty.value = false
  }

  function clear() {
    projects.value = []
    currentProjectId.value = null
    editorContent.value = ''
    isDirty.value = false
  }

  return {
    projects,
    currentProjectId,
    isLoading,
    isSaving,
    isDirty,
    editorContent,
    currentProject,
    setProjects,
    addProject,
    updateProject,
    removeProject,
    selectProject,
    setEditorContent,
    markClean,
    clear
  }
})
