import { api } from '../index'
import type { ProjectInfo, Project, NodeTemplate, ImmediateChange, SaveChanges, SaveNodeChangesResult } from '@/types'
export const projectApi = {
  getAll(): Promise<ProjectInfo[]> {
    return api.project.getAll()
  },
  load(fileName: string): Promise<Project> {
    return api.project.load(fileName)
  },
  save(fileName: string, data: Project) {
    return api.project.save(fileName, data)
  },
  delete(fileName: string) {
    return api.project.delete(fileName)
  },
  rename(oldName: string, newName: string) {
    return api.project.rename(oldName, newName)
  },
  create(name: string, type: string): Promise<number> {
    return api.project.create(name, type)
  },
  updateCurrentNode(projectName: string, nodeId: number | null) {
    return api.project.updateCurrentNode(projectName, nodeId)
  }
}
export const nodeApi = {
  createChild(projectName: string, parentId: number, name: string) {
    return api.node.createChild(projectName, parentId, name)
  },
  createBrother(projectName: string, siblingId: number, name: string) {
    return api.node.createBrother(projectName, siblingId, name)
  },
  rename(nodeId: number, newName: string) {
    return api.node.rename(nodeId, newName)
  },
  updateNote(projectName: string, nodeId: number, note: string) {
    return api.node.updateNote(projectName, nodeId, note)
  },
  delete(projectName: string, nodeId: number) {
    return api.node.delete(projectName, nodeId)
  }
}
export const worldTreeApi = {
  getNodeContent(nodeId: number) {
    return api.worldTree.getNodeContent(nodeId)
  },
  addFolder(nodeId: number, name: string) {
    return api.worldTree.addFolder(nodeId, name)
  },
  addFolderWithParent(nodeId: number, parentChangeId: number | null, name: string) {
    return api.worldTree.addFolderWithParent(nodeId, parentChangeId, name)
  },
  addCard(nodeId: number, folderChangeId: number | null, name: string, keyWord: string) {
    return api.worldTree.addCard(nodeId, folderChangeId, name, keyWord)
  },
  addBlock(nodeId: number, cardChangeId: number | null, title: string, zone: string) {
    return api.worldTree.addBlock(nodeId, cardChangeId, title, zone)
  },
  addLine(nodeId: number, projectId: number, blockChangeId: number, content: string) {
    return api.worldTree.addLine(nodeId, projectId, blockChangeId, content)
  },
  updateLineContent(lineId: number, content: string) {
    return api.worldTree.updateLineContent(lineId, content)
  },
  getSystemFolderChangeIds(projectId: number): Promise<Record<string, number>> {
    return api.worldTree.getSystemFolderChangeIds(projectId)
  },
  getNodeDetail(nodeId: number): Promise<NodeTemplate> {
    return api.worldTree.getNodeDetail(nodeId)
  },
  immediateChange(nodeId: number, change: ImmediateChange): Promise<number> {
    return api.worldTree.immediateChange(nodeId, change)
  },
  saveNodeChanges(nodeId: number, projectId: number, changes: SaveChanges): Promise<SaveNodeChangesResult> {
    return api.worldTree.saveNodeChanges(nodeId, projectId, changes)
  }
}
