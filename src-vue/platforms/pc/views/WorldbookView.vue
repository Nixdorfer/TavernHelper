<template>
  <div class="worldbook-content">
    <div v-if="showProjectSelector" class="project-selector-overlay">
      <div class="project-selector">
        <h1>项目</h1>
        <p class="project-selector-hint">选择或创建一个项目开始</p>
        <div v-if="loadingProjects" class="loading-hint">
          <span class="loading-spinner"></span>
          <span>加载项目列表...</span>
        </div>
        <div v-else-if="projects && projects.length > 0" class="project-list">
          <div
            v-for="project in projects"
            :key="project.fileName"
            :class="['project-card', { encrypted: project.isEncrypted }]"
            @click="selectProject(project)"
          >
            <span :class="['status-dot', getProjectStatusClass(project)]"></span>
            <span class="project-card-time">{{ formatProjectTime(project.updatedAt) }}</span>
            <div class="project-card-name">{{ project.name }}</div>
            <div class="project-card-actions">
              <button class="btn btn-sm btn-secondary" @click.stop="showRenameProject(project)">重命名</button>
              <button class="btn btn-sm btn-danger" @click.stop="showDeleteProject(project)">删除</button>
            </div>
          </div>
        </div>
        <div v-else class="empty-hint">暂无项目</div>
        <div class="project-selector-actions">
          <button class="btn btn-primary" @click="openCreateProjectModal">+ 新建项目</button>
        </div>
      </div>
    </div>
    <div v-if="showCreateProjectModal" class="modal-overlay" @mousedown.self="showCreateProjectModal = false">
      <div class="modal modal-md">
        <div class="modal-header">
          <h3>{{ createProjectStep === 1 ? '新建项目' : '输入项目名称' }}</h3>
          <button class="btn btn-icon" @click="showCreateProjectModal = false">×</button>
        </div>
        <div class="modal-body">
          <template v-if="createProjectStep === 1">
            <div class="project-type-cards">
              <div class="project-type-card" @click="selectProjectType('play')">
                <div class="type-card-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/><path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/></svg>
                </div>
                <div class="type-card-title">游玩</div>
                <div class="type-card-desc">供您游玩他人剧本时提供动态记忆</div>
              </div>
              <div class="project-type-card" @click="selectProjectType('create')">
                <div class="type-card-icon">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/></svg>
                </div>
                <div class="type-card-title">创作</div>
                <div class="type-card-desc">供您从零开始编辑自己的剧本</div>
              </div>
            </div>
          </template>
          <template v-else>
            <div class="form-group">
              <label>项目名称</label>
              <input type="text" class="form-control" v-model="newProjectName" placeholder="输入项目名称" @keyup.enter="createProject" ref="inputCreateProject">
            </div>
          </template>
        </div>
        <div class="modal-footer" v-if="createProjectStep === 2">
          <button class="btn btn-secondary" @click="createProjectStep = 1">返回</button>
          <button class="btn btn-primary" @click="createProject" :disabled="!newProjectName.trim()">创建</button>
        </div>
      </div>
    </div>
    <div v-if="showRenameProjectModal" class="modal-overlay" @mousedown.self="showRenameProjectModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>重命名项目</h3>
          <button class="btn btn-icon" @click="showRenameProjectModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>新名称</label>
            <input type="text" class="form-control" v-model="renameProjectName" placeholder="输入新名称" @keyup.enter="confirmRenameProject" ref="inputRenameProject">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showRenameProjectModal = false">取消</button>
          <button class="btn btn-primary" @click="confirmRenameProject" :disabled="!renameProjectName.trim()">确定</button>
        </div>
      </div>
    </div>
    <div v-if="showCreateNodeModal" class="modal-overlay" @mousedown.self="showCreateNodeModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>{{ createNodeType === 'child' ? '创建子节点' : createNodeType === 'logic' ? '创建逻辑节点' : createNodeType === 'branch' ? '创建分支' : '创建根节点' }}</h3>
          <button class="btn btn-icon" @click="showCreateNodeModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>节点名称</label>
            <input type="text" class="form-control" v-model="newNodeName" placeholder="输入节点名称" @keyup.enter="confirmCreateNode">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showCreateNodeModal = false">取消</button>
          <button class="btn btn-primary" @click="confirmCreateNode" :disabled="!newNodeName.trim()">创建</button>
        </div>
      </div>
    </div>
    <div v-if="showRenameNodeModal" class="modal-overlay" @mousedown.self="showRenameNodeModal = false">
      <div class="modal modal-sm">
        <div class="modal-header">
          <h3>重命名节点</h3>
          <button class="btn btn-icon" @click="showRenameNodeModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>新名称</label>
            <input type="text" class="form-control" v-model="renameNodeName" placeholder="输入新名称" @keyup.enter="confirmRenameNode">
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showRenameNodeModal = false">取消</button>
          <button class="btn btn-primary" @click="confirmRenameNode" :disabled="!renameNodeName.trim()">确定</button>
        </div>
      </div>
    </div>
    <div v-show="!showProjectSelector" class="project-tabs">
      <div class="tabs-container">
        <div
          v-for="(tab, index) in openTabs"
          :key="tab.projectInfo.fileName"
          :class="['project-tab', { active: activeTabIndex === index }]"
          @click="switchTab(index)"
        >
          <span class="tab-project-name">{{ truncateProjectName(tab.projectInfo.name) }}</span>
          <span class="tab-node-name" v-if="getTabCurrentNodeName(tab)">{{ getTabCurrentNodeName(tab) }}</span>
          <button class="tab-close" @click.stop="closeTab(index)" title="关闭">×</button>
        </div>
      </div>
      <button class="tab-add" @click="showProjectSelector = true" title="打开项目">+</button>
    </div>
    <div class="main-container" v-show="!showProjectSelector">
      <aside class="panel timeline-panel collapsed" style="width: 40px; min-width: 40px;">
        <div class="panel-header" @click="toggleFullWorldTree">
          <span :class="['collapse-icon', { expanded: fullWorldTreeActive }]">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
          </span>
        </div>
      </aside>
      <div class="content-area">
        <Transition name="world-tree-slide">
          <div v-if="fullWorldTreeActive && currentProject" class="full-world-tree-overlay">
            <div class="full-world-tree-container" ref="fullWorldTreeContainer">
              <div class="full-tree-scale-wrapper" :style="{ transform: 'scale(' + fullTreeScale + ')' }">
                <TimelineGraph
                  ref="fullWorldTreeGraph"
                  :nodes="currentProject.timeline || []"
                  :current-node-id="currentProject.currentNode"
                  :branch-names="currentProject.branchData || {}"
                  :full-mode="true"
                  @select="handleSelectNode"
                  @set-current="handleSetCurrentNode"
                  @create-child="showCreateChild"
                  @create-logic="showCreateLogic"
                  @create-branch="showCreateBranch"
                  @rename="showRenameNode"
                  @delete="handleDeleteNode"
                  @create-root="handleCreateRootNode"
                  @full-mode-edit="handleFullModeEdit"
                  @node-click="handleNodeClick"
                />
              </div>
            </div>
          </div>
        </Transition>
        <template v-if="currentProject && currentNode">
          <main class="editor-panel">
            <div v-if="editorTabs.length > 0" class="editor-tabs">
              <div
                v-for="(tab, idx) in editorTabs"
                :key="tab.nodeId"
                :class="['editor-tab', { active: idx === activeEditorTabIndex, preview: tab.isPreview }]"
                @click="switchEditorTab(idx)"
                @dblclick="pinEditorTab(idx)"
              >
                <span v-if="isNodeTabUnsaved(tab)" class="tab-status-dot unsaved"></span>
                <span class="editor-tab-name">{{ tab.nodeName || '未命名节点' }}</span>
                <button class="editor-tab-close" @click.stop="closeEditorTab(idx)">×</button>
              </div>
            </div>
            <template v-if="selectedPromptEntry || selectedWorldBookItem">
              <template v-if="selectedPromptEntry && editingPromptEntry">
                <div class="editor-body">
                  <div class="editor-section-tabs">
                    <div v-if="!selectedPromptEntry.isSystem" :class="['section-tab', { active: activePromptSection === 'basic' }]" @click="activePromptSection = 'basic'">基础</div>
                    <div v-if="!selectedPromptEntry.isSystem" :class="['section-tab', { active: activePromptSection === 'triggers' }]" @click="activePromptSection = 'triggers'">触发词</div>
                    <div :class="['section-tab', { active: activePromptSection === 'contentItems' }]" @click="activePromptSection = 'contentItems'">词条</div>
                    <div v-if="!selectedPromptEntry.isSystem" :class="['section-tab', { active: activePromptSection === 't2i' }]" @click="activePromptSection = 't2i'">生图</div>
                  </div>
                  <div class="editor-content">
                    <div v-show="activePromptSection === 'basic'" class="editor-section">
                      <div class="form-group">
                        <label>条目名称</label>
                        <input type="text" class="form-control" v-model="editingPromptEntry.name" placeholder="条目名称">
                      </div>
                      <div class="form-group flex-grow">
                        <label>描述</label>
                        <MarkdownEditor v-model="editingPromptEntry.desc" placeholder="输入描述..." />
                      </div>
                    </div>
                    <div v-show="activePromptSection === 'triggers'" class="editor-section">
                      <div class="form-group">
                        <label>触发词</label>
                        <div class="keywords-container">
                          <div class="match-mode-toggle">
                            <span :class="['match-label', { active: editingPromptEntry?.matchMode === 'or' }]" @click="editingPromptEntry && (editingPromptEntry.matchMode = 'or')">或</span>
                            <div class="match-switch" :class="{ checked: editingPromptEntry?.matchMode !== 'or' }" @click="togglePromptMatchMode"></div>
                            <span :class="['match-label', { active: editingPromptEntry?.matchMode !== 'or' }]" @click="editingPromptEntry && (editingPromptEntry.matchMode = 'and')">和</span>
                          </div>
                          <div class="keyword-input-row">
                            <input
                              type="text"
                              class="form-control"
                              v-model="promptKeywordInput"
                              @keyup.enter="addPromptKeyword"
                              placeholder="输入触发词后按回车添加"
                            >
                          </div>
                        </div>
                      </div>
                      <div class="form-group keywords-list-group">
                        <div class="keywords-list">
                          <span v-for="(kw, idx) in (editingPromptEntry.keywords || [])" :key="idx" class="keyword-tag">
                            {{ kw }}
                            <button @click="removePromptKeyword(idx)" class="keyword-remove">×</button>
                          </span>
                        </div>
                      </div>
                      <div class="form-group trigger-toggles">
                        <label>触发来源</label>
                        <div class="trigger-toggle-row">
                          <span class="trigger-toggle-label">系统</span>
                          <div class="trigger-switch" :class="{ checked: editingPromptEntry?.triggerSystem }" @click="editingPromptEntry && (editingPromptEntry.triggerSystem = !editingPromptEntry.triggerSystem)"></div>
                        </div>
                        <div class="trigger-toggle-row">
                          <span class="trigger-toggle-label">用户</span>
                          <div class="trigger-switch" :class="{ checked: editingPromptEntry?.triggerUser }" @click="editingPromptEntry && (editingPromptEntry.triggerUser = !editingPromptEntry.triggerUser)"></div>
                        </div>
                        <div class="trigger-toggle-row">
                          <span class="trigger-toggle-label">AI</span>
                          <div class="trigger-switch" :class="{ checked: editingPromptEntry?.triggerAI }" @click="editingPromptEntry && (editingPromptEntry.triggerAI = !editingPromptEntry.triggerAI)"></div>
                        </div>
                      </div>
                    </div>
                    <div v-show="activePromptSection === 'contentItems'" class="editor-section">
                      <div v-for="(item, idx) in promptContentItemsWithEmpty" :key="'prompt-content-' + idx" :class="['content-item-card', { collapsed: item.collapsed, 'empty-placeholder': item.isPlaceholder }]">
                        <div class="content-item-header">
                          <span :class="['collapse-toggle', { expanded: !item.collapsed }]" @click="togglePromptContentItemCollapsed(idx)">
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                          </span>
                          <input type="text" class="content-title-input" v-model="item.title" placeholder="点击输入标题..." @click.stop @input="onPromptContentItemInput(idx)" />
                          <button v-if="!item.isPlaceholder" class="btn-remove-content" @click.stop="removePromptContentItem(idx)" title="删除">×</button>
                        </div>
                        <div v-show="!item.collapsed" class="content-item-body">
                          <SyncEditor v-model="item.content" :rows="6" :lines="item.lines" :is-dirty="hasUnsavedChanges" placeholder="输入内容..." @copy-serial="handleCopySerial" @update:modelValue="onPromptContentItemInput(idx)" />
                        </div>
                      </div>
                    </div>
                    <div v-show="activePromptSection === 't2i'" class="editor-section">
                      <div class="form-group">
                        <label>生图提示词</label>
                        <div class="keyword-input-row">
                          <input type="text" class="form-control" v-model="promptT2iInput" @keyup.enter="addPromptT2iKeyword" placeholder="输入生图提示词后按回车添加">
                        </div>
                      </div>
                      <div class="form-group keywords-list-group">
                        <div class="keywords-list">
                          <span v-for="(kw, idx) in (editingPromptEntry.t2iKeywords || [])" :key="idx" class="keyword-tag">
                            {{ kw }}
                            <button @click="removePromptT2iKeyword(idx)" class="keyword-remove">×</button>
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </template>
              <template v-else-if="selectedWorldBookItem && editingEntry">
                <div class="editor-body">
                  <template v-if="!editingEntry.isFolder">
                    <div class="editor-section-tabs">
                      <div :class="['section-tab', { active: activeEditorSection === 'basic' }]" @click="activeEditorSection = 'basic'">基础设置</div>
                      <div :class="['section-tab', { active: activeEditorSection === 'contentItems' }]" @click="activeEditorSection = 'contentItems'">词条内容</div>
                      <div :class="['section-tab', { active: activeEditorSection === 't2iKeywords' }]" @click="activeEditorSection = 't2iKeywords'">生图基础词</div>
                    </div>
                    <div class="editor-content">
                      <div v-show="activeEditorSection === 'basic'" class="editor-section">
                        <div class="form-group">
                          <label>名称</label>
                          <input type="text" class="form-control" v-model="editingEntry.name" placeholder="条目名称">
                        </div>
                        <div class="form-group">
                          <label>触发词</label>
                          <div class="keywords-container">
                            <div class="match-mode-toggle">
                              <span :class="['match-label', { active: editingEntry.matchMode === 'or' }]" @click="editingEntry.matchMode = 'or'">或</span>
                              <div class="match-switch" :class="{ checked: editingEntry.matchMode !== 'or' }" @click="toggleMatchMode"></div>
                              <span :class="['match-label', { active: editingEntry.matchMode !== 'or' }]" @click="editingEntry.matchMode = 'and'">和</span>
                            </div>
                            <div class="keyword-input-row">
                              <input type="text" class="form-control" v-model="newKeywordInput" @keyup.enter="addKeyword" placeholder="输入触发词后按回车添加">
                            </div>
                          </div>
                        </div>
                        <div class="form-group keywords-list-group">
                          <div class="keywords-list">
                            <span v-for="(kw, idx) in editingEntry.keywords" :key="idx" class="keyword-tag">
                              {{ kw }}
                              <button @click="removeKeyword(idx)" class="keyword-remove">×</button>
                            </span>
                          </div>
                        </div>
                      </div>
                      <div v-show="activeEditorSection === 'contentItems'" class="editor-section">
                        <div v-for="(item, idx) in contentItemsWithEmpty" :key="'content-' + item.id" :class="['content-item-card', { collapsed: item.collapsed, 'empty-placeholder': item.isPlaceholder }]">
                          <div class="content-item-header">
                            <span :class="['collapse-toggle', { expanded: !item.collapsed }]" @click="toggleContentItemCollapsed(item)">
                              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                            </span>
                            <span v-if="item.serial" class="content-serial-badge">{{ item.serial }}</span>
                            <input type="text" class="content-title-input" v-model="item.title" :placeholder="getContentFirstLine(item.content) || '点击输入标题...'" @click.stop @input="onContentItemInput(idx)" />
                            <button v-if="!item.isPlaceholder" class="btn-remove-content" @click.stop="confirmRemoveContentItem(idx)" title="删除">×</button>
                          </div>
                          <div v-show="!item.collapsed" class="content-item-body">
                            <div class="content-item-controls">
                              <div class="region-toggles">
                                <label class="region-toggle-switch">
                                  <input type="checkbox" v-model="item.keySystem">
                                  <span class="slider"></span>
                                  <span class="toggle-label">系统</span>
                                </label>
                                <label class="region-toggle-switch">
                                  <input type="checkbox" v-model="item.keyUser">
                                  <span class="slider"></span>
                                  <span class="toggle-label">用户</span>
                                </label>
                                <label class="region-toggle-switch">
                                  <input type="checkbox" v-model="item.keyAI">
                                  <span class="slider"></span>
                                  <span class="toggle-label">AI</span>
                                </label>
                              </div>
                              <div class="value-region-select">
                                <select v-model="item.valueRegion" class="form-control">
                                  <option :value="1">全局词</option>
                                  <option :value="2">前置词</option>
                                  <option :value="4">后置词</option>
                                </select>
                              </div>
                            </div>
                            <SyncEditor v-model="item.content" :rows="8" :lines="item.lines" :is-dirty="hasUnsavedChanges" placeholder="输入词条内容..." @copy-serial="handleCopySerial" @update:modelValue="onContentItemInput(idx)" />
                          </div>
                        </div>
                      </div>
                      <div v-show="activeEditorSection === 't2iKeywords'" class="editor-section">
                        <input type="text" class="t2i-keyword-input" v-model="t2iKeywordInput" placeholder="输入后按回车添加" @keydown.enter.prevent="addT2iKeyword">
                        <div v-if="editingEntry.t2iKeywords && editingEntry.t2iKeywords.length > 0" class="keywords-list t2i-keywords-list">
                          <span v-for="(kw, idx) in editingEntry.t2iKeywords" :key="'t2i-' + idx" class="keyword-tag t2i-tag">
                            {{ kw }}
                            <button @click="removeT2iKeyword(idx)" class="keyword-remove">×</button>
                          </span>
                        </div>
                        <div v-else class="position-empty">暂无生图基础词</div>
                      </div>
                    </div>
                  </template>
                  <template v-else>
                    <div class="folder-edit-content">
                      <div class="form-group">
                        <label>目录名称</label>
                        <input type="text" class="form-control" v-model="editingEntry.name">
                      </div>
                    </div>
                  </template>
                </div>
              </template>
            </template>
            <template v-else-if="editingNodeInfo">
              <div class="editor-header">
                <h3>事件详情</h3>
              </div>
              <div class="editor-body">
                <div class="form-group">
                  <label>事件名称</label>
                  <input type="text" class="form-control" v-model="editingNodeName" placeholder="事件名称">
                </div>
                <div class="form-group flex-grow">
                  <label>事件描述</label>
                  <MarkdownEditor v-model="editingNodeDescription" placeholder="输入事件描述..." />
                </div>
              </div>
            </template>
          </main>
          <aside :class="['panel', 'right-panel', { 'right-panel-collapsed': rightPanelCollapsed }]">
            <div class="right-panel-header" @click="rightPanelCollapsed = !rightPanelCollapsed">
              <span v-if="!rightPanelCollapsed" class="right-panel-header-text" :style="{ background: currentNodeColor, color: currentNodeColor === '#ffffff' ? '#333' : '#fff' }">{{ currentNode?.name || '事件详情' }}</span>
              <span :class="['right-panel-collapse-icon', { expanded: !rightPanelCollapsed }]">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
              </span>
            </div>
            <div v-if="rightPanelCollapsed" class="right-panel-collapsed-footer" @click="rightPanelCollapsed = false">
              <span :class="['collapsed-status-dot', saveStatusDotClass]"></span>
            </div>
            <div class="panel-content" v-show="!rightPanelCollapsed">
              <div class="entry-section">
                <div class="section-header" @click="openNodeInfoEditor">
                  <span class="section-title">事件详情</span>
                  <span class="section-summary">{{ currentNode?.name }}</span>
                </div>
              </div>
              <div v-if="currentNode?.isLogic" class="logic-node-hint">
                <div class="logic-hint-icon">?</div>
                <div class="logic-hint-text">逻辑节点仅用于梳理剧情分支，不包含提示词和世界书内容。</div>
              </div>
              <template v-if="!currentNode?.isLogic">
                <div class="entry-section">
                  <div class="section-header" @click="toggleSection('preText')" @dragover.prevent @drop="onDropToRoot($event, 'preText')">
                    <span :class="['collapse-icon', { expanded: expandedSections.preText }]">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                    </span>
                    <span class="section-title">前置提示词</span>
                    <span class="section-count">({{ preTextCards.length + 1 }})</span>
                    <div class="add-dropdown" @click.stop>
                      <button class="btn-add-circle" @click="togglePromptAddMenu('preText')" title="添加">+</button>
                      <div v-if="promptAddMenu === 'preText'" class="dropdown-menu">
                        <div class="dropdown-item" @click="addPromptEntry('preText')">新建卡片</div>
                        <div class="dropdown-item" @click="addPromptFolder('preText')">新建目录</div>
                      </div>
                    </div>
                  </div>
                  <div v-show="expandedSections.preText" class="prompt-entries">
                    <div :class="['section-item', { active: selectedPromptEntry?.type === 'preText' && selectedPromptEntry?.index === -1 }]" @click="selectSystemCard('preText')" @dragover.prevent @drop="onDrop($event, 'preText')">
                      <span class="badge-pill badge-system">系统</span>
                      <span class="item-name">前置提示词</span>
                    </div>
                    <div
                      v-for="(item, idx) in preTextCards"
                      :key="item.cardId"
                      :class="['section-item', { active: selectedPromptEntry?.type === 'preText' && selectedPromptEntry?.index === idx }]"
                      @click="selectPromptEntryByCard('preText', idx, item)"
                      draggable="true"
                      @dragstart="onDragStart($event, 'card', item.cardId, 'preText')"
                      @dragover.prevent
                      @drop="onDrop($event, 'preText')"
                    >
                      <span class="badge-pill badge-card">卡片</span>
                      <span class="item-name">{{ item.card.name || '条目 ' + (idx + 1) }}</span>
                      <button class="btn-delete-sm" @click.stop="deletePromptEntryByCard('preText', item.cardId)">×</button>
                    </div>
                    <div v-for="folderItem in preTextFolders" :key="'folder-' + folderItem.folderId" :class="['directory-container', { expanded: expandedDirectories.has(folderItem.folderId) }]">
                      <div
                        class="section-item folder-item"
                        @click="toggleDirectory(folderItem.folderId)"
                        draggable="true"
                        @dragstart="onDragStart($event, 'folder', folderItem.folderId, 'preText')"
                        @dragover.prevent
                        @drop="onDropToFolder($event, folderItem.folderId, 'preText')"
                      >
                        <span :class="['collapse-icon', { expanded: expandedDirectories.has(folderItem.folderId) }]">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                        </span>
                        <span class="badge-pill badge-folder">目录</span>
                        <span class="item-name">{{ folderItem.folder.name }}</span>
                      </div>
                      <div v-if="expandedDirectories.has(folderItem.folderId)" class="directory-content">
                        <template v-for="subFolder in getSortedFolderContent(folderItem.folder).folders" :key="'subfolder-' + subFolder.folderId">
                          <div
                            class="section-item folder-item"
                            draggable="true"
                            @dragstart="onDragStart($event, 'folder', subFolder.folderId, 'preText')"
                            @dragover.prevent
                            @drop="onDropToFolder($event, subFolder.folderId, 'preText')"
                          >
                            <span class="badge-pill badge-folder">目录</span>
                            <span class="item-name">{{ subFolder.folder.name }}</span>
                          </div>
                        </template>
                        <template v-for="subCard in getSortedFolderContent(folderItem.folder).cards" :key="'subcard-' + subCard.cardId">
                          <div
                            class="section-item"
                            draggable="true"
                            @dragstart="onDragStart($event, 'card', subCard.cardId, 'preText')"
                            @dragover.prevent
                            @drop="onDrop($event, 'preText')"
                          >
                            <span class="badge-pill badge-card">卡片</span>
                            <span class="item-name">{{ subCard.card.name }}</span>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="entry-section">
                  <div class="section-header" @click="toggleSection('postText')" @dragover.prevent @drop="onDropToRoot($event, 'postText')">
                    <span :class="['collapse-icon', { expanded: expandedSections.postText }]">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                    </span>
                    <span class="section-title">后置提示词</span>
                    <span class="section-count">({{ postTextCards.length + 1 }})</span>
                    <div class="add-dropdown" @click.stop>
                      <button class="btn-add-circle" @click="togglePromptAddMenu('postText')" title="添加">+</button>
                      <div v-if="promptAddMenu === 'postText'" class="dropdown-menu">
                        <div class="dropdown-item" @click="addPromptEntry('postText')">新建卡片</div>
                        <div class="dropdown-item" @click="addPromptFolder('postText')">新建目录</div>
                      </div>
                    </div>
                  </div>
                  <div v-show="expandedSections.postText" class="prompt-entries">
                    <div :class="['section-item', { active: selectedPromptEntry?.type === 'postText' && selectedPromptEntry?.index === -1 }]" @click="selectSystemCard('postText')" @dragover.prevent @drop="onDrop($event, 'postText')">
                      <span class="badge-pill badge-system">系统</span>
                      <span class="item-name">后置提示词</span>
                    </div>
                    <div
                      v-for="(item, idx) in postTextCards"
                      :key="item.cardId"
                      :class="['section-item', { active: selectedPromptEntry?.type === 'postText' && selectedPromptEntry?.index === idx }]"
                      @click="selectPromptEntryByCard('postText', idx, item)"
                      draggable="true"
                      @dragstart="onDragStart($event, 'card', item.cardId, 'postText')"
                      @dragover.prevent
                      @drop="onDrop($event, 'postText')"
                    >
                      <span class="badge-pill badge-card">卡片</span>
                      <span class="item-name">{{ item.card.name || '条目 ' + (idx + 1) }}</span>
                      <button class="btn-delete-sm" @click.stop="deletePromptEntryByCard('postText', item.cardId)">×</button>
                    </div>
                    <div v-for="folderItem in postTextFolders" :key="'folder-' + folderItem.folderId" :class="['directory-container', { expanded: expandedDirectories.has(folderItem.folderId) }]">
                      <div
                        class="section-item folder-item"
                        @click="toggleDirectory(folderItem.folderId)"
                        draggable="true"
                        @dragstart="onDragStart($event, 'folder', folderItem.folderId, 'postText')"
                        @dragover.prevent
                        @drop="onDropToFolder($event, folderItem.folderId, 'postText')"
                      >
                        <span :class="['collapse-icon', { expanded: expandedDirectories.has(folderItem.folderId) }]">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                        </span>
                        <span class="badge-pill badge-folder">目录</span>
                        <span class="item-name">{{ folderItem.folder.name }}</span>
                      </div>
                      <div v-if="expandedDirectories.has(folderItem.folderId)" class="directory-content">
                        <template v-for="subFolder in getSortedFolderContent(folderItem.folder).folders" :key="'subfolder-' + subFolder.folderId">
                          <div
                            class="section-item folder-item"
                            draggable="true"
                            @dragstart="onDragStart($event, 'folder', subFolder.folderId, 'postText')"
                            @dragover.prevent
                            @drop="onDropToFolder($event, subFolder.folderId, 'postText')"
                          >
                            <span class="badge-pill badge-folder">目录</span>
                            <span class="item-name">{{ subFolder.folder.name }}</span>
                          </div>
                        </template>
                        <template v-for="subCard in getSortedFolderContent(folderItem.folder).cards" :key="'subcard-' + subCard.cardId">
                          <div
                            class="section-item"
                            draggable="true"
                            @dragstart="onDragStart($event, 'card', subCard.cardId, 'postText')"
                            @dragover.prevent
                            @drop="onDrop($event, 'postText')"
                          >
                            <span class="badge-pill badge-card">卡片</span>
                            <span class="item-name">{{ subCard.card.name }}</span>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </div>
                <div class="entry-section">
                  <div class="section-header" @click="toggleSection('prePrompt')" @dragover.prevent @drop="onDropToRoot($event, 'prePrompt')">
                    <span :class="['collapse-icon', { expanded: expandedSections.prePrompt }]">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                    </span>
                    <span class="section-title">全局提示词</span>
                    <span class="section-count">({{ prePromptCards.length + 1 }})</span>
                    <div class="add-dropdown" @click.stop>
                      <button class="btn-add-circle" @click="togglePromptAddMenu('prePrompt')" title="添加">+</button>
                      <div v-if="promptAddMenu === 'prePrompt'" class="dropdown-menu">
                        <div class="dropdown-item" @click="addPromptEntry('prePrompt')">新建卡片</div>
                        <div class="dropdown-item" @click="addPromptFolder('prePrompt')">新建目录</div>
                      </div>
                    </div>
                  </div>
                  <div v-show="expandedSections.prePrompt" class="prompt-entries">
                    <div :class="['section-item', { active: selectedPromptEntry?.type === 'prePrompt' && selectedPromptEntry?.index === -1 }]" @click="selectSystemCard('prePrompt')" @dragover.prevent @drop="onDrop($event, 'prePrompt')">
                      <span class="badge-pill badge-system">系统</span>
                      <span class="item-name">全局提示词</span>
                    </div>
                    <div
                      v-for="(item, idx) in prePromptCards"
                      :key="item.cardId"
                      :class="['section-item', { active: selectedPromptEntry?.type === 'prePrompt' && selectedPromptEntry?.index === idx }]"
                      @click="selectPromptEntryByCard('prePrompt', idx, item)"
                      draggable="true"
                      @dragstart="onDragStart($event, 'card', item.cardId, 'prePrompt')"
                      @dragover.prevent
                      @drop="onDrop($event, 'prePrompt')"
                    >
                      <span class="badge-pill badge-card">卡片</span>
                      <span class="item-name">{{ item.card.name || '条目 ' + (idx + 1) }}</span>
                      <button class="btn-delete-sm" @click.stop="deletePromptEntryByCard('prePrompt', item.cardId)">×</button>
                    </div>
                    <div v-for="folderItem in prePromptFolders" :key="'folder-' + folderItem.folderId" :class="['directory-container', { expanded: expandedDirectories.has(folderItem.folderId) }]">
                      <div
                        class="section-item folder-item"
                        @click="toggleDirectory(folderItem.folderId)"
                        draggable="true"
                        @dragstart="onDragStart($event, 'folder', folderItem.folderId, 'prePrompt')"
                        @dragover.prevent
                        @drop="onDropToFolder($event, folderItem.folderId, 'prePrompt')"
                      >
                        <span :class="['collapse-icon', { expanded: expandedDirectories.has(folderItem.folderId) }]">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M9 5l7 7-7 7"/></svg>
                        </span>
                        <span class="badge-pill badge-folder">目录</span>
                        <span class="item-name">{{ folderItem.folder.name }}</span>
                      </div>
                      <div v-if="expandedDirectories.has(folderItem.folderId)" class="directory-content">
                        <template v-for="subFolder in getSortedFolderContent(folderItem.folder).folders" :key="'subfolder-' + subFolder.folderId">
                          <div
                            class="section-item folder-item"
                            draggable="true"
                            @dragstart="onDragStart($event, 'folder', subFolder.folderId, 'prePrompt')"
                            @dragover.prevent
                            @drop="onDropToFolder($event, subFolder.folderId, 'prePrompt')"
                          >
                            <span class="badge-pill badge-folder">目录</span>
                            <span class="item-name">{{ subFolder.folder.name }}</span>
                          </div>
                        </template>
                        <template v-for="subCard in getSortedFolderContent(folderItem.folder).cards" :key="'subcard-' + subCard.cardId">
                          <div
                            class="section-item"
                            draggable="true"
                            @dragstart="onDragStart($event, 'card', subCard.cardId, 'prePrompt')"
                            @dragover.prevent
                            @drop="onDrop($event, 'prePrompt')"
                          >
                            <span class="badge-pill badge-card">卡片</span>
                            <span class="item-name">{{ subCard.card.name }}</span>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </div>
              </template>
            </div>
            <div v-show="!rightPanelCollapsed" class="right-panel-footer">
              <div v-if="justSaved" class="save-status-pill saved" @click="saveProject">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12l5 5L20 7"/></svg>
                <span>已保存</span>
              </div>
              <div v-else-if="hasUnsavedChanges" class="save-status-pill unsaved" @click="saveProject">
                <span :class="['save-status-dot', saveStatusDotClass]"></span>
                <span>{{ lastSaveTimeAgo }}</span>
              </div>
              <div v-else class="save-status-pill saved" @click="saveProject">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M5 12l5 5L20 7"/></svg>
                <span>已保存</span>
              </div>
            </div>
          </aside>
        </template>
        <div v-if="!currentProject || !currentNode" class="empty-state full-width">
          <p>{{ currentProject ? '请选择一个节点' : '请选择一个项目' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, reactive, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useNotificationStore, useConfirmStore } from '@/stores'
import { useWorldbookProjectStore, useWorldbookEditorStore, useWorldbookTabsStore, useWorldbookTreeStore } from '@/stores/modules/worldbook'
import { projectApi, nodeApi, worldTreeApi } from '@/api/modules/project'
import type { NodeTemplate, FolderTemplate, CardTemplate } from '@/types'
import TimelineGraph from '@/components/business/TimelineGraph.vue'
import MarkdownEditor from '@/components/business/MarkdownEditor.vue'
import SyncEditor from '@/components/business/SyncEditor.vue'
interface ProjectInfo {
  fileName: string
  name: string
  updatedAt?: string
  isEncrypted?: boolean
  decryptFailed?: boolean
}
interface OpenTab {
  projectInfo: ProjectInfo
  project: any
}
interface EditorTab {
  nodeId: string
  nodeName: string
  isPreview?: boolean
}
interface PromptEntry {
  type: string
  index: number
  isSystem?: boolean
}
const notificationStore = useNotificationStore()
const confirmStore = useConfirmStore()
const projectStore = useWorldbookProjectStore()
const editorStore = useWorldbookEditorStore()
const tabsStore = useWorldbookTabsStore()
const treeStore = useWorldbookTreeStore()
const showProjectSelector = ref(true)
const loadingProjects = ref(false)
const fullWorldTreeActive = ref(false)
const fullTreeScale = ref(1)
const fullWorldTreeContainer = ref<HTMLElement | null>(null)
const rightPanelCollapsed = ref(false)
const lastSavedAt = ref<number | null>(null)
const justSaved = ref(false)
const hasUnsavedChanges = ref(false)
const systemFolderChangeIds = ref<Record<string, number>>({})
async function loadSystemFolderChangeIds() {
  if (!currentProject.value?.id) return
  try {
    systemFolderChangeIds.value = await worldTreeApi.getSystemFolderChangeIds(currentProject.value.id)
    console.log('[loadSystemFolderChangeIds] 系统文件夹changeId:', systemFolderChangeIds.value)
  } catch (e) {
    console.error('[loadSystemFolderChangeIds] 获取失败:', e)
  }
}
const nodeDetail = ref<NodeTemplate | null>(null)
const isLoadingNodeDetail = ref(false)
async function loadNodeDetail(nodeId: number) {
  if (isLoadingNodeDetail.value) return
  isLoadingNodeDetail.value = true
  try {
    console.log('[loadNodeDetail] 加载节点详情:', nodeId)
    nodeDetail.value = await worldTreeApi.getNodeDetail(nodeId)
    console.log('[loadNodeDetail] 节点详情加载完成:', nodeDetail.value)
  } catch (e) {
    console.error('[loadNodeDetail] 加载失败:', e)
    nodeDetail.value = null
  } finally {
    isLoadingNodeDetail.value = false
  }
}
function findFolderByName(name: string): { folderId: string; folder: FolderTemplate } | null {
  if (!nodeDetail.value?.folders) return null
  for (const [folderId, folder] of Object.entries(nodeDetail.value.folders)) {
    if (folder.name === name) {
      return { folderId, folder }
    }
  }
  return null
}
function getCardsForSection(sectionName: string): Array<{ cardId: string; card: CardTemplate }> {
  const folderData = findFolderByName(sectionName)
  if (!folderData) return []
  return Object.entries(folderData.folder.cards || {})
    .filter(([_, card]) => !card.name.endsWith('_BASE'))
    .map(([cardId, card]) => ({ cardId, card }))
}
function getSystemCard(sectionName: string): { cardId: string; card: CardTemplate } | null {
  const folderData = findFolderByName(sectionName)
  if (!folderData) return null
  for (const [cardId, card] of Object.entries(folderData.folder.cards || {})) {
    if (card.name.endsWith('_BASE')) {
      return { cardId, card }
    }
  }
  return null
}
function getFoldersForSection(sectionName: string): Array<{ folderId: string; folder: FolderTemplate }> {
  const folderData = findFolderByName(sectionName)
  if (!folderData) return []
  return Object.entries(folderData.folder.folders || {}).map(([folderId, folder]) => ({ folderId, folder }))
}
const preTextCards = computed(() => getCardsForSection('SYS_PRE'))
const preTextFolders = computed(() => getFoldersForSection('SYS_PRE'))
const postTextFolders = computed(() => getFoldersForSection('SYS_POST'))
const prePromptFolders = computed(() => getFoldersForSection('SYS_GLOBAL'))
const postTextCards = computed(() => getCardsForSection('SYS_POST'))
const prePromptCards = computed(() => getCardsForSection('SYS_GLOBAL'))
const saveStatusDotClass = computed(() => {
  if (!lastSavedAt.value) return 'red'
  const diff = Date.now() - lastSavedAt.value
  if (diff < 60000) return 'green'
  if (diff < 300000) return 'yellow'
  return 'red'
})
const lastSaveTimeAgo = computed(() => {
  if (!lastSavedAt.value) return '未保存'
  const diff = Date.now() - lastSavedAt.value
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚保存'
  if (minutes < 60) return `${minutes}分钟前保存`
  const hours = Math.floor(minutes / 60)
  return `${hours}小时前保存`
})
const editingNodeInfo = ref(false)
const editingNodeName = ref('')
const editingNodeDescription = ref('')
const activePromptSection = ref('basic')
const activeEditorSection = ref('basic')
const promptKeywordInput = ref('')
const promptT2iInput = ref('')
const newKeywordInput = ref('')
const t2iKeywordInput = ref('')
const promptAddMenu = ref('')
const expandedSections = reactive({
  preText: true,
  postText: true,
  prePrompt: true
})
const expandedDirectories = ref<Set<string>>(new Set())
function toggleDirectory(folderId: string) {
  if (expandedDirectories.value.has(folderId)) {
    expandedDirectories.value.delete(folderId)
  } else {
    expandedDirectories.value.add(folderId)
  }
  expandedDirectories.value = new Set(expandedDirectories.value)
}
function getSortedFolderContent(folder: FolderTemplate): { folders: Array<{ folderId: string; folder: FolderTemplate }>; cards: Array<{ cardId: string; card: CardTemplate }> } {
  const folders = Object.entries(folder.folders || {}).map(([folderId, f]) => ({ folderId, folder: f })).sort((a, b) => (a.folder.name || '').localeCompare(b.folder.name || ''))
  const cards = Object.entries(folder.cards || {}).map(([cardId, c]) => ({ cardId, card: c })).sort((a, b) => (a.card.name || '').localeCompare(b.card.name || ''))
  return { folders, cards }
}
const selectedPromptEntry = ref<PromptEntry | null>(null)
const selectedWorldBookItem = ref<any>(null)
const editingPromptEntry = ref<any>(null)
const editingEntry = ref<any>(null)
const effectiveContentItems = ref<any[]>([])
const editorTabs = ref<EditorTab[]>([])
const activeEditorTabIndex = ref(-1)
const openTabs = ref<OpenTab[]>([])
const activeTabIndex = ref(-1)
const showCreateProjectModal = ref(false)
const createProjectStep = ref(1)
const newProjectType = ref<'play' | 'create'>('create')
const newProjectName = ref('')
const showRenameProjectModal = ref(false)
const renameProjectName = ref('')
const renamingProject = ref<ProjectInfo | null>(null)
const inputCreateProject = ref<HTMLInputElement | null>(null)
const inputRenameProject = ref<HTMLInputElement | null>(null)
const showCreateNodeModal = ref(false)
const createNodeType = ref<'child' | 'logic' | 'branch' | 'root'>('child')
const createNodeParentId = ref<number | null>(null)
const newNodeName = ref('')
const showRenameNodeModal = ref(false)
const renamingNodeId = ref<number | null>(null)
const renameNodeName = ref('')
const projects = computed(() => projectStore.projects)
const currentProject = computed(() => projectStore.currentProject)
const currentNode = computed(() => {
  if (!currentProject.value?.timeline || !currentProject.value?.currentNode) return null
  return currentProject.value.timeline?.find(n => n.id === currentProject.value!.currentNode) || null
})
const BRANCH_COLORS = ['#3b82f6', '#f59e0b', '#8b5cf6', '#ec4899', '#06b6d4', '#84cc16', '#f97316', '#6366f1', '#14b8a6', '#eab308']
const currentNodeColor = computed(() => {
  if (!currentNode.value) return '#3b82f6'
  if (!currentNode.value.parentId) return '#ffffff'
  const branchId = currentNode.value.branchId
  if (!branchId) return BRANCH_COLORS[0]
  const branchData = currentProject.value?.branchData?.[branchId]
  if (branchData?.color) return branchData.color
  const branchIndex = Object.keys(currentProject.value?.branchData || {}).indexOf(branchId)
  return BRANCH_COLORS[branchIndex >= 0 ? branchIndex % BRANCH_COLORS.length : 0]
})
const emptyPromptContentItem = { title: '', content: '', collapsed: true, isPlaceholder: true }
const promptContentItemsWithEmpty = computed(() => {
  const items = editingPromptEntry.value?.contentItems || []
  return [...items, { ...emptyPromptContentItem }]
})
const emptyContentItem = { id: 'placeholder', title: '', content: '', collapsed: true, isPlaceholder: true, keySystem: false, keyUser: false, keyAI: true, valueRegion: 1 }
const contentItemsWithEmpty = computed(() => {
  const items = effectiveContentItems.value || []
  return [...items, { ...emptyContentItem }]
})
function onPromptContentItemInput(idx: number) {
  if (!editingPromptEntry.value) return
  const items = editingPromptEntry.value.contentItems || []
  if (idx === items.length) {
    const placeholder = promptContentItemsWithEmpty.value[idx]
    if (!editingPromptEntry.value.contentItems) {
      editingPromptEntry.value.contentItems = []
    }
    editingPromptEntry.value.contentItems.push({ title: placeholder?.title || '', content: placeholder?.content || '', collapsed: false })
  }
}
function onContentItemInput(idx: number) {
  if (!editingEntry.value) return
  const items = editingEntry.value.contentItems || []
  if (idx === items.length) {
    const placeholder = contentItemsWithEmpty.value[idx]
    if (!editingEntry.value.contentItems) {
      editingEntry.value.contentItems = []
    }
    const newItem = { id: Date.now().toString(), title: placeholder?.title || '', content: placeholder?.content || '', collapsed: false, keySystem: false, keyUser: false, keyAI: true, valueRegion: 1 }
    editingEntry.value.contentItems.push(newItem)
    effectiveContentItems.value = [...editingEntry.value.contentItems]
  }
}
function getProjectStatusClass(project: ProjectInfo): string {
  if (project.isEncrypted && !project.decryptFailed) return 'status-encrypted'
  if (project.decryptFailed) return 'status-error'
  return 'status-normal'
}
function getTabStatusClass(tab: OpenTab, index: number): string {
  return index === activeTabIndex.value ? 'status-active' : 'status-normal'
}
function formatProjectTime(time?: string): string {
  if (!time) return ''
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return date.toLocaleDateString()
}
function truncateProjectName(name: string): string {
  if (!name) return ''
  return name.length > 5 ? name.slice(0, 5) + '…' : name
}
function getTabCurrentNodeName(tab: OpenTab): string {
  const nodeId = tab.project?.currentNode
  if (!nodeId) return ''
  const node = tab.project?.timeline?.find((n: any) => n.id === nodeId)
  return node?.name || ''
}
async function loadProjects() {
  loadingProjects.value = true
  try {
    await projectStore.loadProjects()
  } finally {
    loadingProjects.value = false
  }
}
const isLoadingProject = ref(false)
async function selectProject(project: ProjectInfo) {
  if (isLoadingProject.value) {
    console.log('[SelectProject] 已在加载中，跳过:', project.fileName)
    return
  }
  isLoadingProject.value = true
  try {
    console.log('[SelectProject] 开始加载项目:', project.fileName)
    const result = await projectStore.loadProject(project.fileName)
    console.log('[SelectProject] 项目加载完成:', currentProject.value)
    await loadSystemFolderChangeIds()
    showProjectSelector.value = false
    openTabs.value.push({
      projectInfo: project,
      project: currentProject.value
    })
    activeTabIndex.value = openTabs.value.length - 1
    if (!currentProject.value?.currentNode && currentProject.value?.timeline?.length) {
      const rootNode = currentProject.value.timeline.find(n => !n.parentId)
      if (rootNode) {
        currentProject.value.currentNode = rootNode.id
      }
    }
    if (currentProject.value?.currentNode) {
      fullWorldTreeActive.value = false
      openNodeInfoEditor()
    } else {
      fullWorldTreeActive.value = true
    }
  } catch (e: any) {
    console.error('[SelectProject] 加载失败:', e)
    notificationStore.showNotification('加载项目失败: ' + e.message, 'error')
  } finally {
    isLoadingProject.value = false
  }
}
function showRenameProject(project: ProjectInfo) {
  renamingProject.value = project
  renameProjectName.value = project.name
  showRenameProjectModal.value = true
}
async function showDeleteProject(project: ProjectInfo) {
  const confirmed = await confirmStore.show({
    title: '删除项目',
    message: `确定要删除项目 "${project.name}" 吗？此操作不可撤销。`,
    confirmText: '删除',
    cancelText: '取消',
    type: 'danger'
  })
  if (confirmed) {
    try {
      await projectStore.deleteProject(project.fileName)
      notificationStore.showNotification('项目已删除', 'success')
    } catch (e: any) {
      const msg = e?.message || String(e) || '未知错误'
      notificationStore.showNotification('删除失败: ' + msg, 'error')
    }
  }
}
function openCreateProjectModal() {
  newProjectName.value = ''
  createProjectStep.value = 1
  newProjectType.value = 'create'
  showCreateProjectModal.value = true
}
function selectProjectType(type: 'play' | 'create') {
  newProjectType.value = type
  createProjectStep.value = 2
}
async function createProject() {
  if (!newProjectName.value.trim()) return
  try {
    console.log('[CreateProject] 开始创建项目:', newProjectName.value.trim(), newProjectType.value)
    await projectApi.create(newProjectName.value.trim(), newProjectType.value)
    console.log('[CreateProject] 项目创建成功，重新加载项目列表')
    await projectStore.loadProjects()
    showCreateProjectModal.value = false
    notificationStore.showNotification('项目创建成功', 'success')
    const created = projects.value?.find(p => p.name === newProjectName.value.trim())
    console.log('[CreateProject] 查找创建的项目:', created)
    if (created) {
      await selectProject(created)
    }
  } catch (e: any) {
    console.error('[CreateProject] 创建失败:', e)
    notificationStore.showNotification('创建失败: ' + e.message, 'error')
  }
}
async function confirmRenameProject() {
  if (!renameProjectName.value.trim() || !renamingProject.value) return
  try {
    const oldName = renamingProject.value.fileName
    const newName = renameProjectName.value.trim() + '.json'
    await projectStore.renameProject(oldName, newName)
    showRenameProjectModal.value = false
    notificationStore.showNotification('重命名成功', 'success')
  } catch (e: any) {
    notificationStore.showNotification('重命名失败: ' + e.message, 'error')
  }
}
function switchTab(index: number) {
  activeTabIndex.value = index
  const tab = openTabs.value[index]
  if (tab) {
    projectStore.setProject(tab.project)
  }
}
function closeTab(index: number) {
  openTabs.value.splice(index, 1)
  if (openTabs.value.length === 0) {
    showProjectSelector.value = true
    projectStore.clear()
    activeTabIndex.value = -1
  } else if (index <= activeTabIndex.value) {
    activeTabIndex.value = Math.max(0, activeTabIndex.value - 1)
    switchTab(activeTabIndex.value)
  }
}
function toggleFullWorldTree() {
  fullWorldTreeActive.value = !fullWorldTreeActive.value
}
async function saveCurrentNode(nodeId: number) {
  const projectName = openTabs.value[activeTabIndex.value]?.projectInfo.name
  if (!projectName) return
  try {
    await projectApi.updateCurrentNode(projectName, nodeId)
  } catch (e) {
    console.error('Failed to save current node:', e)
  }
}
async function handleSelectNode(node: any) {
  if (!currentProject.value) return
  currentProject.value.currentNode = node.id
  await saveCurrentNode(node.id)
}
function handleSetCurrentNode(node: any) {
  handleSelectNode(node)
}
async function handleFullModeEdit(node: any) {
  if (!currentProject.value || !node) return
  currentProject.value.currentNode = node.id
  await saveCurrentNode(node.id)
  fullWorldTreeActive.value = false
  openNodeInfoEditor()
}
async function handleNodeClick(node: any) {
  if (!currentProject.value || !node) return
  currentProject.value.currentNode = node.id
  await saveCurrentNode(node.id)
}
function showCreateChild(node: any) {
  createNodeType.value = 'child'
  createNodeParentId.value = node?.id || currentProject.value?.currentNode || null
  newNodeName.value = ''
  showCreateNodeModal.value = true
}
function showCreateLogic(node: any) {
  createNodeType.value = 'logic'
  createNodeParentId.value = node?.id || currentProject.value?.currentNode || null
  newNodeName.value = ''
  showCreateNodeModal.value = true
}
function showCreateBranch(node: any) {
  createNodeType.value = 'branch'
  createNodeParentId.value = node?.id || currentProject.value?.currentNode || null
  newNodeName.value = ''
  showCreateNodeModal.value = true
}
function showRenameNode(node: any) {
  if (!node) return
  renamingNodeId.value = node.id
  renameNodeName.value = node.name || ''
  showRenameNodeModal.value = true
}
async function handleDeleteNode(node: any) {
  if (!node || !currentProject.value) return
  const confirmed = await confirmStore.show({
    title: '删除节点',
    message: `确定要删除节点 "${node.name}" 吗？此操作不可撤销。`,
    confirmText: '删除',
    cancelText: '取消',
    type: 'danger'
  })
  if (confirmed) {
    try {
      const projectName = openTabs.value[activeTabIndex.value]?.projectInfo.fileName
      if (!projectName) return
      await nodeApi.delete(projectName, node.id)
      currentProject.value.timeline = currentProject.value.timeline?.filter(n => n.id !== node.id && n.parentId !== node.id) || []
      if (currentProject.value.currentNode === node.id) {
        currentProject.value.currentNode = undefined
      }
      notificationStore.showNotification('节点已删除', 'success')
    } catch (e: any) {
      const msg = e?.message || String(e) || '未知错误'
      notificationStore.showNotification('删除失败: ' + msg, 'error')
    }
  }
}
function handleCreateRootNode() {
  createNodeType.value = 'root'
  createNodeParentId.value = null
  newNodeName.value = ''
  showCreateNodeModal.value = true
}
async function confirmCreateNode() {
  if (!newNodeName.value.trim() || !currentProject.value) return
  try {
    const projectName = openTabs.value[activeTabIndex.value]?.projectInfo.name
    if (!projectName) return
    let result: any
    if (createNodeType.value === 'root') {
      result = await nodeApi.createChild(projectName, 0, newNodeName.value.trim())
    } else if (createNodeType.value === 'branch') {
      result = await nodeApi.createBrother(projectName, createNodeParentId.value || 0, newNodeName.value.trim())
    } else {
      result = await nodeApi.createChild(projectName, createNodeParentId.value || 0, newNodeName.value.trim())
    }
    if (result) {
      if (!currentProject.value.timeline) currentProject.value.timeline = []
      currentProject.value.timeline.push(result)
      currentProject.value.currentNode = result.id
      await saveCurrentNode(result.id)
    }
    showCreateNodeModal.value = false
    notificationStore.showNotification('节点创建成功', 'success')
  } catch (e: any) {
    notificationStore.showNotification('创建失败: ' + e.message, 'error')
  }
}
async function confirmRenameNode() {
  if (!renameNodeName.value.trim() || !renamingNodeId.value) return
  try {
    await nodeApi.rename(renamingNodeId.value, renameNodeName.value.trim())
    const node = currentProject.value?.timeline?.find(n => n.id === renamingNodeId.value)
    if (node) node.name = renameNodeName.value.trim()
    showRenameNodeModal.value = false
    notificationStore.showNotification('重命名成功', 'success')
  } catch (e: any) {
    notificationStore.showNotification('重命名失败: ' + e.message, 'error')
  }
}
function switchEditorTab(idx: number) {
  activeEditorTabIndex.value = idx
}
function pinEditorTab(idx: number) {
  const tab = editorTabs.value[idx]
  if (tab) tab.isPreview = false
}
function closeEditorTab(idx: number) {
  editorTabs.value.splice(idx, 1)
  if (idx <= activeEditorTabIndex.value) {
    activeEditorTabIndex.value = Math.max(0, activeEditorTabIndex.value - 1)
  }
}
function isNodeTabUnsaved(tab: EditorTab): boolean {
  return false
}
function toggleSection(section: 'preText' | 'postText' | 'prePrompt') {
  expandedSections[section] = !expandedSections[section]
}
const dragData = ref<{ type: 'card' | 'folder'; id: string; source: string } | null>(null)
function onDragStart(e: DragEvent, itemType: 'card' | 'folder', id: string, source: string) {
  dragData.value = { type: itemType, id, source }
  if (e.dataTransfer) {
    e.dataTransfer.effectAllowed = 'move'
  }
}
function onDrop(e: DragEvent, target: string) {
  if (!dragData.value) return
  const { type, id, source } = dragData.value
  if (source === target) {
    dragData.value = null
    return
  }
  console.log(`[Drag] Moving ${type} ${id} from ${source} to ${target}`)
  dragData.value = null
}
function onDropToFolder(e: DragEvent, folderId: string, section: string) {
  e.stopPropagation()
  if (!dragData.value) return
  const { type, id } = dragData.value
  if (type === 'card' || type === 'folder') {
    console.log(`[Drag] Moving ${type} ${id} into folder ${folderId} in section ${section}`)
  }
  dragData.value = null
}
function onDropToRoot(e: DragEvent, section: string) {
  e.stopPropagation()
  if (!dragData.value) return
  const { type, id } = dragData.value
  if (type === 'card' || type === 'folder') {
    console.log(`[Drag] Moving ${type} ${id} to root of section ${section}`)
  }
  dragData.value = null
}
function togglePromptAddMenu(type: string) {
  promptAddMenu.value = promptAddMenu.value === type ? '' : type
}
async function addPromptFolder(type: string) {
  if (!currentNode.value) return
  const sysFolderName = type === 'preText' ? 'SYS_PRE' : type === 'postText' ? 'SYS_POST' : 'SYS_GLOBAL'
  const parentFolderChangeId = systemFolderChangeIds.value[sysFolderName]
  if (parentFolderChangeId && currentNode.value.id) {
    try {
      await worldTreeApi.addFolderWithParent(currentNode.value.id, parentFolderChangeId, '新目录')
      await loadNodeDetail(currentNode.value.id)
    } catch (e) {
      console.error('[addPromptFolder] 创建文件夹失败:', e)
    }
  }
  promptAddMenu.value = ''
  notificationStore.showNotification('目录已添加', 'success')
}
async function addPromptEntry(type: string) {
  console.log('[addPromptEntry] 添加条目, type:', type)
  console.log('[addPromptEntry] currentNode:', currentNode.value)
  if (!currentNode.value) return
  const fieldName = type === 'preText' ? 'pre_text' : type === 'postText' ? 'post_text' : 'pre_prompt'
  if (!currentNode.value[fieldName]) currentNode.value[fieldName] = []
  const sysFolderName = type === 'preText' ? 'SYS_PRE' : type === 'postText' ? 'SYS_POST' : 'SYS_GLOBAL'
  const folderChangeId = systemFolderChangeIds.value[sysFolderName]
  console.log('[addPromptEntry] folderChangeId:', folderChangeId, 'for', sysFolderName)
  if (folderChangeId && currentNode.value.id) {
    try {
      const newCardChangeId = await worldTreeApi.addCard(currentNode.value.id, folderChangeId, '新条目', '')
      console.log('[addPromptEntry] 创建card成功, changeId:', newCardChangeId)
      await loadNodeDetail(currentNode.value.id)
      const cards = type === 'preText' ? preTextCards.value : type === 'postText' ? postTextCards.value : prePromptCards.value
      const newCardIndex = cards.findIndex(c => c.cardId === String(newCardChangeId))
      if (newCardIndex >= 0) {
        const newCard = cards[newCardIndex]
        selectPromptEntryByCard(type, newCardIndex, newCard)
      }
    } catch (e) {
      console.error('[addPromptEntry] 创建card失败:', e)
    }
  }
  promptAddMenu.value = ''
  notificationStore.showNotification('条目已添加', 'success')
}
function selectSystemCard(type: string) {
  selectedPromptEntry.value = { type, index: -1, isSystem: true }
  selectedWorldBookItem.value = null
  activePromptSection.value = 'contentItems'
  const sysFolder = type === 'preText' ? 'SYS_PRE' : type === 'postText' ? 'SYS_POST' : 'SYS_GLOBAL'
  const systemCard = getSystemCard(sysFolder)
  isUpdatingSyncDots.value = true
  if (systemCard) {
    const contentItems: any[] = []
    for (const [blockId, block] of Object.entries(systemCard.card.blocks || {})) {
      const lines = Object.entries(block.lines || {}).map(([sn, lineData]) => ({
        sn,
        content: lineData.content || '',
        syncDot: lineData.syncDot
      }))
      contentItems.push({
        id: blockId,
        blockChangeId: Number(blockId),
        title: block.title,
        content: lines.map(l => l.content).join('\n'),
        lines
      })
    }
    editingPromptEntry.value = {
      cardId: systemCard.cardId,
      cardChangeId: Number(systemCard.cardId),
      name: systemCard.card.name,
      desc: systemCard.card.desc || '',
      keywords: [],
      contentItems,
      t2iKeywords: [],
      isSystem: true
    }
  } else {
    editingPromptEntry.value = {
      name: `系统${type === 'preText' ? '前置' : type === 'postText' ? '后置' : '全局'}提示词`,
      desc: '',
      keywords: [],
      contentItems: [],
      t2iKeywords: [],
      isSystem: true
    }
  }
  editingEntry.value = null
  nextTick(() => { isUpdatingSyncDots.value = false })
}
function selectPromptEntry(data: { type: string; idx: number }) {
  selectedPromptEntry.value = { type: data.type, index: data.idx }
  selectedWorldBookItem.value = null
  const entries = currentNode.value?.[data.type === 'preText' ? 'pre_text' : data.type === 'postText' ? 'post_text' : 'pre_prompt'] || []
  if (entries[data.idx]) {
    editingPromptEntry.value = { ...entries[data.idx], desc: entries[data.idx].desc || '', matchMode: entries[data.idx].matchMode || 'or' }
  } else {
    editingPromptEntry.value = null
  }
  editingEntry.value = null
}
function selectPromptEntryByCard(type: string, idx: number, item: { cardId: string; card: CardTemplate }) {
  selectedPromptEntry.value = { type, index: idx }
  selectedWorldBookItem.value = null
  isUpdatingSyncDots.value = true
  const contentItems: any[] = []
  for (const [blockId, block] of Object.entries(item.card.blocks || {})) {
    const lines = Object.entries(block.lines || {}).map(([sn, lineData]) => ({
      sn,
      content: lineData.content || '',
      syncDot: lineData.syncDot
    }))
    contentItems.push({
      id: blockId,
      blockChangeId: Number(blockId),
      title: block.title,
      content: lines.map(l => l.content).join('\n'),
      lines
    })
  }
  editingPromptEntry.value = {
    cardId: item.cardId,
    cardChangeId: Number(item.cardId),
    name: item.card.name,
    desc: item.card.desc || '',
    matchMode: item.card.trigger?.mode || 'or',
    keywords: item.card.trigger?.words || [],
    triggerSystem: item.card.trigger?.system ?? false,
    triggerUser: item.card.trigger?.user ?? true,
    triggerAI: item.card.trigger?.ai ?? true,
    t2iKeywords: item.card.image || [],
    contentItems
  }
  editingEntry.value = null
  editingNodeInfo.value = false
  nextTick(() => { isUpdatingSyncDots.value = false })
}
async function deletePromptEntryByCard(type: string, cardId: string) {
  const confirmed = await confirmStore.confirmDelete('此卡片')
  if (!confirmed) return
  console.log('[deletePromptEntryByCard] 删除卡片:', cardId)
  notificationStore.showNotification('删除功能待实现', 'info')
}
async function deletePromptEntry(data: { type: string; idx: number }) {
  if (!currentNode.value) return
  const fieldName = data.type === 'preText' ? 'pre_text' : data.type === 'postText' ? 'post_text' : 'pre_prompt'
  const entries = currentNode.value[fieldName] || []
  const entry = entries[data.idx]
  if (!entry) return
  const confirmed = await confirmStore.confirmDelete(entry.name || `条目 ${data.idx + 1}`)
  if (!confirmed) return
  currentNode.value[fieldName].splice(data.idx, 1)
  if (selectedPromptEntry.value?.type === data.type && selectedPromptEntry.value?.index === data.idx) {
    selectedPromptEntry.value = null
    editingPromptEntry.value = null
  }
  notificationStore.showNotification('条目已删除', 'success')
}
async function openNodeInfoEditor() {
  editingNodeInfo.value = true
  if (currentNode.value) {
    editingNodeName.value = currentNode.value.name || ''
    editingNodeDescription.value = currentNode.value.note || ''
    await loadNodeDetail(currentNode.value.id)
  }
  selectedPromptEntry.value = null
  selectedWorldBookItem.value = null
  editingPromptEntry.value = null
  editingEntry.value = null
}
function addPromptContentItem() {
  if (!editingPromptEntry.value) return
  if (!editingPromptEntry.value.contentItems) {
    editingPromptEntry.value.contentItems = []
  }
  editingPromptEntry.value.contentItems.push({
    title: '',
    content: '',
    collapsed: false
  })
  activePromptSection.value = 'contentItems'
}
function togglePromptContentItemCollapsed(idx: number) {
  if (editingPromptEntry.value?.contentItems?.[idx]) {
    editingPromptEntry.value.contentItems[idx].collapsed = !editingPromptEntry.value.contentItems[idx].collapsed
  }
}
function removePromptContentItem(idx: number) {
  editingPromptEntry.value?.contentItems?.splice(idx, 1)
}
function addPromptKeyword() {
  if (!promptKeywordInput.value.trim() || !editingPromptEntry.value) return
  if (!editingPromptEntry.value.keywords) {
    editingPromptEntry.value.keywords = []
  }
  editingPromptEntry.value.keywords.push(promptKeywordInput.value.trim())
  promptKeywordInput.value = ''
}
function removePromptKeyword(idx: number) {
  editingPromptEntry.value?.keywords?.splice(idx, 1)
}
function addPromptT2iKeyword() {
  if (!promptT2iInput.value.trim() || !editingPromptEntry.value) return
  if (!editingPromptEntry.value.t2iKeywords) {
    editingPromptEntry.value.t2iKeywords = []
  }
  editingPromptEntry.value.t2iKeywords.push(promptT2iInput.value.trim())
  promptT2iInput.value = ''
}
function removePromptT2iKeyword(idx: number) {
  editingPromptEntry.value?.t2iKeywords?.splice(idx, 1)
}
function addContentItem() {
  console.log('[addContentItem] 添加词条内容')
  console.log('[addContentItem] editingEntry:', editingEntry.value)
  if (!editingEntry.value) return
  if (!editingEntry.value.contentItems) {
    editingEntry.value.contentItems = []
  }
  editingEntry.value.contentItems.push({
    id: Date.now().toString(),
    title: '',
    content: '',
    collapsed: false,
    keySystem: false,
    keyUser: false,
    keyAI: true,
    valueRegion: 1
  })
  effectiveContentItems.value = [...editingEntry.value.contentItems]
  console.log('[addContentItem] 添加后contentItems:', editingEntry.value.contentItems)
}
function toggleContentItemCollapsed(item: any) {
  item.collapsed = !item.collapsed
}
function confirmRemoveContentItem(idx: number) {
  effectiveContentItems.value.splice(idx, 1)
  if (editingEntry.value?.contentItems) {
    editingEntry.value.contentItems.splice(idx, 1)
  }
}
function getContentFirstLine(content: string): string {
  if (!content) return ''
  const firstLine = content.split('\n')[0]
  return firstLine.length > 30 ? firstLine.slice(0, 30) + '...' : firstLine
}
function toggleMatchMode() {
  if (!editingEntry.value) return
  editingEntry.value.matchMode = editingEntry.value.matchMode === 'or' ? 'and' : 'or'
}
function togglePromptMatchMode() {
  if (!editingPromptEntry.value) return
  editingPromptEntry.value.matchMode = editingPromptEntry.value.matchMode === 'or' ? 'and' : 'or'
}
function addKeyword() {
  if (!newKeywordInput.value.trim() || !editingEntry.value) return
  if (!editingEntry.value.keywords) {
    editingEntry.value.keywords = []
  }
  editingEntry.value.keywords.push(newKeywordInput.value.trim())
  newKeywordInput.value = ''
}
function removeKeyword(idx: number) {
  editingEntry.value?.keywords?.splice(idx, 1)
}
function addT2iKeyword() {
  if (!t2iKeywordInput.value.trim() || !editingEntry.value) return
  if (!editingEntry.value.t2iKeywords) {
    editingEntry.value.t2iKeywords = []
  }
  editingEntry.value.t2iKeywords.push(t2iKeywordInput.value.trim())
  t2iKeywordInput.value = ''
}
function removeT2iKeyword(idx: number) {
  editingEntry.value?.t2iKeywords?.splice(idx, 1)
}
const isUpdatingSyncDots = ref(false)
function markDirty() {
  if (isUpdatingSyncDots.value) return
  hasUnsavedChanges.value = true
  justSaved.value = false
  isUpdatingSyncDots.value = true
  clearSyncDots()
  isUpdatingSyncDots.value = false
}
function clearSyncDots() {
  if (editingPromptEntry.value?.contentItems) {
    for (const item of editingPromptEntry.value.contentItems) {
      if (item.lines) {
        for (const l of item.lines) {
          l.syncDot = ''
        }
      }
    }
  }
  if (editingEntry.value?.contentItems) {
    for (const item of editingEntry.value.contentItems) {
      if (item.lines) {
        for (const l of item.lines) {
          l.syncDot = ''
        }
      }
    }
  }
}
function updateSyncDotsFromResult(blocks: Record<string, Record<string, { content: string | null; syncDot: string }>>) {
  console.log('[updateSyncDotsFromResult] 收到blocks:', JSON.stringify(blocks, null, 2))
  isUpdatingSyncDots.value = true
  const updateItems = (contentItems: any[], label: string) => {
    console.log(`[updateSyncDotsFromResult] 更新${label}, contentItems数量:`, contentItems.length)
    for (const item of contentItems) {
      const blockKey = String(item.blockChangeId)
      console.log(`[updateSyncDotsFromResult] item blockChangeId:`, item.blockChangeId, 'blockKey:', blockKey, 'blocks中是否存在:', !!blocks[blockKey])
      if (blocks[blockKey]) {
        const blockLines = blocks[blockKey]
        console.log(`[updateSyncDotsFromResult] blockLines:`, blockLines)
        item.lines = Object.entries(blockLines).map(([sn, lineData]) => ({
          sn,
          content: lineData.content || '',
          syncDot: lineData.syncDot
        }))
        console.log(`[updateSyncDotsFromResult] 更新后item.lines:`, item.lines)
      }
    }
  }
  if (editingPromptEntry.value?.contentItems) {
    updateItems(editingPromptEntry.value.contentItems, 'editingPromptEntry')
  }
  if (editingEntry.value?.contentItems) {
    updateItems(editingEntry.value.contentItems, 'editingEntry')
  }
  nextTick(() => {
    isUpdatingSyncDots.value = false
    console.log('[updateSyncDotsFromResult] 完成, isUpdatingSyncDots已重置')
  })
}
function handleCopySerial(serial: string) {
  notificationStore.showNotification(`已复制行序列号: ${serial}`, 'success')
}
function markSaved() {
  lastSavedAt.value = Date.now()
  hasUnsavedChanges.value = false
  justSaved.value = true
  setTimeout(() => {
    justSaved.value = false
  }, 2000)
}
async function saveProject() {
  console.log('[saveProject] 开始保存')
  if (!currentProject.value || !currentNode.value) {
    console.log('[saveProject] 无currentProject或currentNode, 跳过')
    return
  }
  try {
    const nodeId = currentNode.value.id
    const projectId = currentProject.value.id
    const changes: any = { card: {}, block: {} }
    if (editingPromptEntry.value && selectedPromptEntry.value) {
      const zone = selectedPromptEntry.value.type === 'preText' ? 'pre'
        : selectedPromptEntry.value.type === 'postText' ? 'post' : 'global'
      const cardChangeId = editingPromptEntry.value.cardChangeId
      if (cardChangeId) {
        if (selectedPromptEntry.value.index === -1) {
          changes.card[cardChangeId] = {
            name: editingPromptEntry.value.name,
            desc: editingPromptEntry.value.desc || ''
          }
        } else {
          changes.card[cardChangeId] = {
            name: editingPromptEntry.value.name,
            desc: editingPromptEntry.value.desc || '',
            trigger: {
              mode: editingPromptEntry.value.matchMode || 'or',
              words: editingPromptEntry.value.keywords || [],
              system: editingPromptEntry.value.triggerSystem || false,
              user: editingPromptEntry.value.triggerUser !== false,
              ai: editingPromptEntry.value.triggerAI !== false
            },
            image: editingPromptEntry.value.t2iKeywords || []
          }
        }
      }
      const contentItems = (editingPromptEntry.value.contentItems || []).filter(item => item.title?.trim() || item.content?.trim())
      for (const item of contentItems) {
        if (!item.blockChangeId && cardChangeId) {
          const blockChangeId = await worldTreeApi.addBlock(nodeId, cardChangeId, item.title || '词条', zone)
          item.blockChangeId = blockChangeId
        }
        if (item.blockChangeId) {
          changes.block[item.blockChangeId] = { content: item.content ?? '' }
        }
      }
      editingPromptEntry.value.contentItems = contentItems
      if (selectedPromptEntry.value.index >= 0) {
        const fieldName = selectedPromptEntry.value.type === 'preText' ? 'pre_text'
          : selectedPromptEntry.value.type === 'postText' ? 'post_text' : 'pre_prompt'
        if (currentNode.value[fieldName]) {
          currentNode.value[fieldName][selectedPromptEntry.value.index] = { ...editingPromptEntry.value }
        }
      }
    }
    if (editingEntry.value && selectedWorldBookItem.value) {
      const cardChangeId = editingEntry.value.changeId
      if (cardChangeId) {
        changes.card[cardChangeId] = {
          name: editingEntry.value.name,
          desc: editingEntry.value.desc || '',
          trigger: {
            mode: editingEntry.value.matchMode || 'or',
            words: editingEntry.value.keywords || [],
            system: editingEntry.value.triggerSystem || false,
            user: editingEntry.value.triggerUser !== false,
            ai: editingEntry.value.triggerAI !== false
          },
          image: editingEntry.value.t2iKeywords || []
        }
      }
      const contentItems = (editingEntry.value.contentItems || []).filter(item => item.title?.trim() || item.content?.trim())
      for (const item of contentItems) {
        if (!item.blockChangeId && cardChangeId) {
          const blockChangeId = await worldTreeApi.addBlock(nodeId, cardChangeId, item.title || '词条', 'card')
          item.blockChangeId = blockChangeId
        }
        if (item.blockChangeId) {
          changes.block[item.blockChangeId] = { content: item.content ?? '' }
        }
      }
      editingEntry.value.contentItems = contentItems
      effectiveContentItems.value = [...contentItems]
    }
    const hasChanges = Object.keys(changes.card).length > 0 || Object.keys(changes.block).length > 0
    console.log('[saveProject] hasChanges:', hasChanges, 'changes.block keys:', Object.keys(changes.block))
    if (hasChanges) {
      console.log('[saveProject] 调用saveNodeChanges:', JSON.stringify(changes, null, 2))
      const result = await worldTreeApi.saveNodeChanges(nodeId, projectId, changes)
      console.log('[saveProject] saveNodeChanges返回:', JSON.stringify(result, null, 2))
      if (result?.blocks) {
        console.log('[saveProject] 调用updateSyncDotsFromResult前, editingPromptEntry.contentItems:', JSON.stringify(editingPromptEntry.value?.contentItems?.map(i => ({ blockChangeId: i.blockChangeId, lines: i.lines })), null, 2))
        updateSyncDotsFromResult(result.blocks)
        console.log('[saveProject] 调用updateSyncDotsFromResult后, editingPromptEntry.contentItems:', JSON.stringify(editingPromptEntry.value?.contentItems?.map(i => ({ blockChangeId: i.blockChangeId, lines: i.lines })), null, 2))
      } else {
        console.log('[saveProject] result.blocks为空或undefined')
      }
    }
    if (editingNodeInfo.value) {
      currentNode.value.name = editingNodeName.value
      currentNode.value.note = editingNodeDescription.value
      await nodeApi.rename(currentNode.value.id, editingNodeName.value)
      await nodeApi.updateNote(currentProject.value.name, currentNode.value.id, editingNodeDescription.value || '')
    }
    const projectName = openTabs.value[activeTabIndex.value]?.projectInfo.name
    if (projectName) {
      await projectApi.save(projectName, currentProject.value)
    }
    markSaved()
    notificationStore.showNotification('保存成功', 'success')
  } catch (e: any) {
    console.error('[saveProject] 保存失败:', e)
    notificationStore.showNotification('保存失败: ' + e.message, 'error')
  }
}
function handleKeyDown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 's') {
    e.preventDefault()
    saveProject()
  }
}
watch([editingPromptEntry, editingEntry], () => {
  if (editingPromptEntry.value || editingEntry.value) {
    markDirty()
  }
}, { deep: true })
watch([editingNodeName, editingNodeDescription], () => {
  if (editingNodeInfo.value) {
    markDirty()
  }
})
let saveTimeUpdateInterval: ReturnType<typeof setInterval> | null = null
onMounted(() => {
  loadProjects()
  saveTimeUpdateInterval = setInterval(() => {
    if (lastSavedAt.value) {
      lastSavedAt.value = lastSavedAt.value
    }
  }, 60000)
  window.addEventListener('keydown', handleKeyDown)
})
onUnmounted(() => {
  if (saveTimeUpdateInterval) {
    clearInterval(saveTimeUpdateInterval)
  }
  window.removeEventListener('keydown', handleKeyDown)
})
</script>
<style scoped>
.worldbook-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
.project-selector-overlay {
  position: absolute;
  inset: 0;
  background: var(--bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
.project-selector {
  text-align: center;
  max-width: 600px;
  width: 90%;
}
.project-selector h1 {
  font-size: 32px;
  margin-bottom: 12px;
}
.project-selector-hint {
  color: var(--text-secondary);
  margin-bottom: 24px;
}
.project-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 400px;
  overflow-y: auto;
  margin-bottom: 20px;
}
.project-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}
.project-card:hover {
  background: var(--bg-hover);
}
.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-dot.status-normal {
  background: #22c55e;
}
.status-dot.status-encrypted {
  background: #f59e0b;
}
.status-dot.status-error {
  background: #ef4444;
}
.status-dot.status-active {
  background: var(--primary-color);
}
.project-card-time {
  font-size: 11px;
  color: var(--text-tertiary);
  flex-shrink: 0;
}
.project-card-name {
  flex: 1;
  text-align: left;
  font-weight: 500;
}
.project-card-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s;
}
.project-card:hover .project-card-actions {
  opacity: 1;
}
.project-selector-actions {
  display: flex;
  justify-content: center;
}
.loading-hint,
.empty-hint {
  padding: 20px;
  color: var(--text-secondary);
  text-align: center;
}
.loading-spinner {
  display: inline-block;
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin-right: 8px;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}
.project-tabs {
  display: flex;
  align-items: center;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 0 8px;
}
.tabs-container {
  display: flex;
  flex: 1;
  overflow-x: auto;
  gap: 2px;
}
.project-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border-radius: 6px 6px 0 0;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
  white-space: nowrap;
  max-width: 200px;
}
.project-tab:hover {
  background: var(--bg-hover);
}
.project-tab.active {
  background: var(--bg-primary);
}
.tab-project-name {
  padding: 2px 10px;
  background: var(--primary-color);
  color: white;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  flex-shrink: 0;
}
.tab-node-name {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  color: var(--text-secondary);
  font-size: 12px;
}
.tab-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.tab-name {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
}
.tab-close {
  padding: 0 4px;
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0.6;
  font-size: 14px;
  flex-shrink: 0;
  color: var(--text-primary);
}
.tab-close:hover {
  opacity: 1;
}
.tab-add {
  padding: 6px 12px;
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 18px;
  color: var(--text-secondary);
}
.tab-add:hover {
  color: var(--text-primary);
}
.main-container {
  flex: 1;
  display: flex;
  overflow: hidden;
}
.panel {
  background: var(--bg-secondary);
  display: flex;
  flex-direction: column;
}
.timeline-panel {
  border-right: 1px solid var(--border-color);
}
.panel-header {
  padding: 12px;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}
.collapse-icon {
  width: 16px;
  height: 16px;
  transition: transform 0.2s;
}
.collapse-icon.expanded {
  transform: rotate(90deg);
}
.collapse-icon svg {
  width: 100%;
  height: 100%;
  color: var(--text-secondary);
}
.content-area {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
}
.editor-panel {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-primary);
}
.full-world-tree-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--bg-primary);
  z-index: 10;
  overflow: auto;
}
.full-world-tree-container {
  padding: 20px;
  min-height: 100%;
}
.full-tree-scale-wrapper {
  transform-origin: top left;
}
.editor-tabs {
  display: flex;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  overflow-x: auto;
}
.editor-tab {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 12px;
  cursor: pointer;
  font-size: 12px;
  border-right: 1px solid var(--border-color);
  transition: background 0.2s;
}
.editor-tab:hover {
  background: var(--bg-hover);
}
.editor-tab.active {
  background: var(--bg-primary);
}
.editor-tab.preview .editor-tab-name {
  font-style: italic;
  opacity: 0.7;
}
.tab-status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.tab-status-dot.unsaved {
  background: #f59e0b;
}
.editor-tab-close {
  padding: 0 2px;
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0.6;
}
.editor-header {
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}
.editor-header h3 {
  margin: 0;
  font-size: 14px;
}
.editor-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  overflow-y: auto;
  gap: 12px;
  padding: 12px 16px;
}
.editor-section-tabs {
  display: flex;
  gap: 4px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-color);
}
.section-tab {
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  background: var(--bg-tertiary);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 4px;
}
.section-tab:hover {
  background: var(--bg-hover);
}
.section-tab.active {
  background: var(--primary-color);
  color: white;
}
.btn-add-content-item {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  border: none;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.editor-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}
.editor-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.form-group.flex-grow {
  flex: 1;
}
.form-group label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-secondary);
}
.form-control {
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
}
.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
}
.keywords-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.match-mode-toggle {
  display: flex;
  align-items: center;
  gap: 8px;
}
.match-label {
  font-size: 12px;
  color: var(--text-secondary);
  cursor: pointer;
}
.match-label.active {
  color: var(--primary-color);
  font-weight: 500;
}
.match-switch {
  width: 32px;
  height: 18px;
  background: var(--bg-tertiary);
  border-radius: 9px;
  cursor: pointer;
  position: relative;
}
.match-switch::after {
  content: '';
  position: absolute;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: white;
  top: 2px;
  left: 2px;
  transition: transform 0.2s;
}
.match-switch.checked::after {
  transform: translateX(14px);
}
.match-switch.checked {
  background: var(--primary-color);
}
.trigger-toggles {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.trigger-toggle-row {
  display: flex;
  align-items: center;
  gap: 12px;
}
.trigger-toggle-label {
  font-size: 13px;
  color: var(--text-primary);
  min-width: 40px;
}
.trigger-switch {
  width: 32px;
  height: 18px;
  background: var(--bg-tertiary);
  border-radius: 9px;
  cursor: pointer;
  position: relative;
}
.trigger-switch::after {
  content: '';
  position: absolute;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  background: white;
  top: 2px;
  left: 2px;
  transition: transform 0.2s;
}
.trigger-switch.checked {
  background: var(--primary-color);
}
.trigger-switch.checked::after {
  transform: translateX(14px);
}
.keyword-input-row {
  display: flex;
  gap: 8px;
}
.keyword-input-row .form-control {
  flex: 1;
}
.keywords-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.keyword-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: var(--bg-tertiary);
  border-radius: 4px;
  font-size: 12px;
}
.keyword-remove {
  padding: 0 2px;
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0.6;
}
.keyword-remove:hover {
  opacity: 1;
}
.content-item-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  margin-bottom: 12px;
}
.content-item-card.empty-placeholder {
  border-style: dashed;
  opacity: 0.6;
}
.content-item-card.empty-placeholder:hover {
  opacity: 1;
}
.content-item-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 12px;
  cursor: pointer;
}
.collapse-toggle {
  width: 16px;
  height: 16px;
  transition: transform 0.2s;
}
.collapse-toggle.expanded {
  transform: rotate(90deg);
}
.collapse-toggle svg {
  width: 100%;
  height: 100%;
  color: var(--text-secondary);
}
.content-serial-badge {
  font-size: 10px;
  padding: 2px 6px;
  background: var(--primary-color);
  color: white;
  border-radius: 4px;
}
.content-title-input {
  flex: 1;
  padding: 4px 8px;
  border: 1px solid transparent;
  border-radius: 4px;
  background: transparent;
  font-size: 13px;
  line-height: 1.4;
  color: var(--text-primary);
}
.content-title-input:focus {
  outline: none;
  border-color: var(--border-color);
  background: var(--bg-primary);
}
.btn-remove-content {
  padding: 2px 6px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--text-secondary);
  opacity: 0;
  transition: opacity 0.2s;
}
.content-item-header:hover .btn-remove-content {
  opacity: 1;
}
.btn-remove-content:hover {
  color: #ef4444;
}
.content-item-body {
  padding: 0 12px 12px;
}
.content-item-controls {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 12px;
}
.region-toggles {
  display: flex;
  gap: 12px;
}
.region-toggle-switch {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  font-size: 12px;
}
.region-toggle-switch input {
  width: 14px;
  height: 14px;
}
.toggle-label {
  color: var(--text-secondary);
}
.value-region-select .form-control {
  padding: 4px 8px;
  font-size: 12px;
}
.content-textarea {
  min-height: 120px;
  resize: none;
  margin: 8px;
  width: calc(100% - 16px);
}
.position-empty {
  padding: 30px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 13px;
}
.t2i-keyword-input {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
  margin-bottom: 12px;
}
.t2i-keywords-list {
  margin-top: 8px;
}
.right-panel {
  width: 280px;
  min-width: 280px;
  border-left: 1px solid var(--border-color);
  overflow: hidden;
  transition: width 0.2s ease, min-width 0.2s ease;
}
.right-panel.right-panel-collapsed {
  width: 40px;
  min-width: 40px;
  position: relative;
}
.right-panel > .panel-content,
.right-panel > .right-panel-footer {
  width: 280px;
  min-width: 280px;
  flex-shrink: 0;
}
.right-panel-header {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px;
  min-height: 47px;
  box-sizing: border-box;
  border-bottom: 1px solid var(--border-color);
  cursor: pointer;
  position: relative;
}
.right-panel-header-text {
  font-size: 13px;
  font-weight: 500;
  padding: 4px 12px;
  border-radius: 6px;
}
.right-panel-collapsed-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 12px 0;
  min-height: 50px;
  box-sizing: border-box;
  border-top: 1px solid var(--border-color);
  cursor: pointer;
}
.right-panel-collapse-icon {
  position: absolute;
  right: 12px;
  width: 16px;
  height: 16px;
  transition: transform 0.2s;
}
.right-panel-collapse-icon.expanded {
  transform: rotate(180deg);
}
.right-panel-collapse-icon svg {
  width: 100%;
  height: 100%;
  color: var(--text-secondary);
}
.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.entry-section {
  margin-bottom: 8px;
}
.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.2s;
}
.section-header:hover {
  background: var(--bg-hover);
}
.section-header.active {
  background: var(--primary-color);
  color: white;
}
.section-title {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
}
.section-summary {
  font-size: 12px;
  color: var(--text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100px;
}
.section-count {
  font-size: 11px;
  color: var(--text-tertiary);
}
.add-dropdown {
  position: relative;
}
.btn-add-circle {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  border: none;
  cursor: pointer;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.dropdown-container {
  position: relative;
  display: inline-block;
}
.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  z-index: 100;
  min-width: 120px;
}
.dropdown-item {
  padding: 8px 12px;
  font-size: 12px;
  cursor: pointer;
  transition: background 0.2s;
}
.dropdown-item:hover {
  background: var(--bg-hover);
}
.prompt-entries {
  padding-left: 24px;
}
.section-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 12px;
  transition: background 0.2s;
}
.section-item:hover {
  background: var(--bg-hover);
}
.section-item.active {
  background: rgba(var(--primary-rgb), 0.15);
  color: var(--primary-color);
}
.badge-pill {
  font-size: 10px;
  padding: 2px 6px;
  color: white;
  border-radius: 10px;
  flex-shrink: 0;
}
.badge-system {
  background: #8b5cf6;
}
.badge-card {
  background: #16a34a;
}
.badge-folder {
  background: #c2410c;
}
.item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.folder-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: var(--text-secondary);
}
.directory-container {
  position: relative;
}
.directory-container.expanded > .section-item.folder-item {
  position: relative;
}
.directory-container.expanded > .section-item.folder-item::before {
  content: '';
  position: absolute;
  left: 4px;
  top: 50%;
  bottom: -4px;
  width: 2px;
  background: #ea580c;
}
.directory-content {
  position: relative;
  padding-left: 12px;
}
.directory-content::before {
  content: '';
  position: absolute;
  left: 4px;
  top: 0;
  bottom: 4px;
  width: 2px;
  background: #ea580c;
}
.directory-content .section-item {
  position: relative;
}
.btn-delete-sm {
  padding: 0 4px;
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--text-secondary);
  opacity: 0;
}
.section-item:hover .btn-delete-sm {
  opacity: 1;
}
.btn-delete-sm:hover {
  color: #ef4444;
}
.logic-node-hint {
  display: flex;
  gap: 8px;
  padding: 12px;
  background: #fef3c7;
  border-radius: 6px;
  margin: 8px;
}
.logic-hint-icon {
  width: 20px;
  height: 20px;
  background: #f59e0b;
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  flex-shrink: 0;
}
.logic-hint-text {
  font-size: 12px;
  color: #92400e;
  line-height: 1.4;
}
.folder-edit-content {
  padding: 16px;
}
.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-secondary);
}
.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
}
.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}
.btn-primary {
  background: var(--primary-color);
  color: white;
}
.btn-primary:hover {
  background: var(--primary-hover);
}
.btn-secondary {
  background: var(--bg-tertiary);
  color: var(--text-primary);
}
.btn-secondary:hover {
  background: var(--bg-hover);
}
.btn-danger {
  background: #ef4444;
  color: white;
}
.btn-danger:hover {
  background: #dc2626;
}
.modal-md {
  width: 480px;
}
.project-type-cards {
  display: flex;
  gap: 16px;
  justify-content: center;
}
.project-type-card {
  flex: 1;
  max-width: 200px;
  padding: 24px 16px;
  background: var(--bg-tertiary);
  border: 2px solid var(--border-color);
  border-radius: 12px;
  cursor: pointer;
  text-align: center;
  transition: all 0.2s;
}
.project-type-card:hover {
  border-color: var(--primary-color);
  background: var(--bg-hover);
}
.type-card-icon {
  width: 48px;
  height: 48px;
  margin: 0 auto 12px;
  color: var(--primary-color);
}
.type-card-icon svg {
  width: 100%;
  height: 100%;
}
.type-card-title {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
}
.type-card-desc {
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
}
.world-tree-slide-enter-active {
  animation: world-tree-slide-down 0.3s ease-out;
}
.world-tree-slide-leave-active {
  animation: world-tree-fade-out 0.2s ease-out;
}
@keyframes world-tree-slide-down {
  0% {
    transform: translateY(-100%);
    opacity: 0;
  }
  100% {
    transform: translateY(0);
    opacity: 1;
  }
}
@keyframes world-tree-fade-out {
  0% {
    opacity: 1;
  }
  100% {
    opacity: 0;
  }
}
.right-panel-footer {
  padding: 12px;
  border-top: 1px solid var(--border-color);
  flex-shrink: 0;
}
.save-status-pill {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 16px;
  font-size: 12px;
  cursor: pointer;
  transition: opacity 0.2s;
}
.save-status-pill:hover {
  opacity: 0.8;
}
.save-status-pill.saved {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}
.save-status-pill.saved svg {
  width: 14px;
  height: 14px;
}
.save-status-pill.unsaved {
  background: var(--bg-tertiary);
  color: var(--text-secondary);
}
.save-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}
.save-status-dot.green {
  background: #22c55e;
}
.save-status-dot.yellow {
  background: #eab308;
}
.save-status-dot.red {
  background: #ef4444;
}
.save-status-dot.red-flash {
  background: #ef4444;
  animation: flash 1s infinite;
}
@keyframes flash {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}
.collapsed-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.collapsed-status-dot.green {
  background: #22c55e;
}
.collapsed-status-dot.yellow {
  background: #eab308;
}
.collapsed-status-dot.red {
  background: #ef4444;
  animation: flash 1s infinite;
}
</style>
