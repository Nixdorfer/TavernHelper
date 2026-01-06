<template>
  <aside class="global-sidebar">
    <div class="sidebar-menu">
      <div
        :class="['sidebar-item', { active: activePanel === 'worldbook' }]"
        @click="$emit('switch-panel', 'worldbook')"
      >
        <span class="sidebar-text">项目</span>
      </div>
      <div
        v-if="token"
        :class="['sidebar-item', { active: activePanel === 'conversation' }]"
        @click="$emit('switch-panel', 'conversation')"
      >
        <span class="sidebar-text">聊天</span>
      </div>
      <div
        :class="['sidebar-item', { active: activePanel === 'drafts' }]"
        @click="$emit('switch-panel', 'drafts')"
      >
        <span class="sidebar-text">草稿箱</span>
      </div>
      <div
        :class="['sidebar-item', { active: activePanel === 'creation' }]"
        @click="$emit('switch-panel', 'creation')"
      >
        <span class="sidebar-text">创作</span>
      </div>
      <div
        :class="['sidebar-item', { active: activePanel === 'gallery' }]"
        @click="$emit('switch-panel', 'gallery')"
      >
        <span class="sidebar-text">图库</span>
      </div>
      <div :class="['sidebar-item-group', { expanded: plazaExpanded }]">
        <div
          :class="['sidebar-item', 'plaza-item', { active: activePanel === 'plaza' }]"
          @click="$emit('click-plaza-main')"
        >
          <span class="sidebar-text">广场</span>
        </div>
        <transition name="plaza-sub">
          <div v-if="plazaExpanded" class="sidebar-sub-items">
            <transition-group name="plaza-item" tag="div">
              <div
                v-for="(tab, index) in plazaTabs"
                :key="tab.key"
                :class="['sidebar-sub-item', { active: activePanel === 'plaza' && plazaCurrentTab === tab.key }]"
                :style="{ transitionDelay: index * 30 + 'ms' }"
                @click="$emit('select-plaza-tab', tab.key)"
              >
                {{ tab.label }}
              </div>
            </transition-group>
          </div>
        </transition>
      </div>
      <div
        v-if="showDebug"
        :class="['sidebar-item debug-item', { active: activePanel === 'debug' }]"
        @click="$emit('switch-panel', 'debug')"
      >
        <span class="sidebar-text">调试</span>
      </div>
    </div>
    <div class="sidebar-footer">
      <div class="sidebar-item" @click="$emit('settings')">
        <span class="sidebar-text">设置</span>
      </div>
    </div>
  </aside>
</template>
<script>
export default {
  name: 'GlobalSidebar',
  props: {
    token: { type: String, default: '' },
    activePanel: { type: String, default: 'worldbook' },
    plazaExpanded: { type: Boolean, default: false },
    plazaCurrentTab: { type: String, default: '' },
    plazaTabs: { type: Array, default: () => [] },
    showDebug: { type: Boolean, default: false }
  },
  emits: ['switch-panel', 'click-plaza-main', 'select-plaza-tab', 'settings']
}
</script>
<style scoped>
.global-sidebar {
  width: 72px;
  min-width: 72px;
  flex-shrink: 0;
  background: var(--sidebar-bg);
  display: flex;
  flex-direction: column;
  border-right: 1px solid var(--border-color);
  padding: 12px 0;
  z-index: 100;
}
.sidebar-menu {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 0 8px;
}
.sidebar-item {
  padding: 12px 8px;
  text-align: center;
  cursor: pointer;
  border-radius: 8px;
  transition: background 0.2s;
}
.sidebar-item:hover {
  background: var(--hover-bg);
}
.sidebar-item.active {
  background: var(--primary-color);
  color: #fff;
}
.sidebar-text {
  font-size: 12px;
  writing-mode: horizontal-tb;
}
.sidebar-item-group {
  display: flex;
  flex-direction: column;
  border-radius: 8px;
  transition: background 0.2s;
}
.sidebar-item-group.expanded {
  background: transparent;
}
.sidebar-sub-items {
  padding-left: 0;
  overflow: hidden;
}
.sidebar-sub-item {
  padding: 8px 6px;
  font-size: 11px;
  text-align: center;
  cursor: pointer;
  border-radius: 6px;
  color: var(--text-secondary);
  transition: all 0.2s;
}
.sidebar-sub-item:hover {
  background: var(--hover-bg);
  color: var(--text-color);
}
.sidebar-sub-item.active {
  background: var(--primary-light);
  color: var(--primary-color);
}
.plaza-item.active {
  background: var(--primary-light);
  color: var(--primary-color);
}
.sidebar-footer {
  padding: 8px;
  border-top: 1px solid var(--border-color);
  margin-top: 8px;
}
.plaza-sub-enter-active,
.plaza-sub-leave-active {
  transition: all 0.3s ease;
}
.plaza-sub-enter-from,
.plaza-sub-leave-to {
  opacity: 0;
  max-height: 0;
}
.plaza-sub-enter-to,
.plaza-sub-leave-from {
  opacity: 1;
  max-height: 200px;
}
.plaza-item-enter-active,
.plaza-item-leave-active {
  transition: all 0.2s ease;
}
.plaza-item-enter-from {
  opacity: 0;
  transform: translateX(-10px);
}
.plaza-item-leave-to {
  opacity: 0;
  transform: translateX(10px);
}
.debug-item {
  margin-top: 8px;
  border-top: 1px dashed var(--border-color);
  padding-top: 12px;
}
.debug-item.active {
  background: #f59e0b;
}
</style>
