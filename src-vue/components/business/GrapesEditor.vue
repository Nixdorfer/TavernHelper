<template>
  <div class="grapes-editor-container" :class="{ 'code-fullscreen': showCodeFullscreen }">
    <div v-if="showClearConfirm" class="modal-overlay" @click.self="showClearConfirm = false">
      <div class="confirm-modal">
        <div class="confirm-title">清空画布</div>
        <div class="confirm-content">确定要清空画布吗？此操作不可撤销。</div>
        <div class="confirm-actions">
          <button class="btn-cancel" @click="showClearConfirm = false">取消</button>
          <button class="btn-confirm" @click="confirmClearCanvas">确定</button>
        </div>
      </div>
    </div>
    <div v-if="showRatioInput" class="modal-overlay" @click.self="showRatioInput = false">
      <div class="ratio-modal">
        <div class="ratio-title">自定义比例</div>
        <div class="ratio-inputs">
          <input type="number" v-model.number="customRatioW" min="1" max="100" placeholder="宽" class="ratio-input">
          <span class="ratio-separator">:</span>
          <input type="number" v-model.number="customRatioH" min="1" max="100" placeholder="高" class="ratio-input">
        </div>
        <div class="ratio-actions">
          <button class="btn-cancel" @click="showRatioInput = false">取消</button>
          <button class="btn-confirm" @click="applyCustomRatio">应用</button>
        </div>
      </div>
    </div>
    <div v-if="!hideToolbar" class="editor-toolbar">
      <div class="toolbar-left">
        <button :class="['toolbar-btn', { active: currentViewMode === 'visual' }]" @click="setViewMode('visual')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M3 9h18M9 21V9"/></svg>
          <span>可视化</span>
        </button>
        <button :class="['toolbar-btn', { active: currentViewMode === 'code' }]" @click="setViewMode('code')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="16,18 22,12 16,6"/><polyline points="8,6 2,12 8,18"/></svg>
          <span>代码</span>
        </button>
        <button :class="['toolbar-btn', { active: currentViewMode === 'split' }]" @click="setViewMode('split')">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="12" y1="3" x2="12" y2="21"/></svg>
          <span>分屏</span>
        </button>
      </div>
      <div class="toolbar-center">
        <button class="toolbar-btn" @click="undo" title="撤销">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 7v6h6"/><path d="M21 17a9 9 0 00-9-9 9 9 0 00-6 2.3L3 13"/></svg>
        </button>
        <button class="toolbar-btn" @click="redo" title="重做">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 7v6h-6"/><path d="M3 17a9 9 0 019-9 9 9 0 016 2.3L21 13"/></svg>
        </button>
        <div class="toolbar-divider"></div>
        <button class="toolbar-btn" @click="toggleFullscreen" :title="isFullscreen ? '退出全屏' : '全屏'">
          <svg v-if="!isFullscreen" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8 3H5a2 2 0 00-2 2v3m18 0V5a2 2 0 00-2-2h-3m0 18h3a2 2 0 002-2v-3M3 16v3a2 2 0 002 2h3"/></svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8 3v3a2 2 0 01-2 2H3m18 0h-3a2 2 0 01-2-2V3m0 18v-3a2 2 0 012-2h3M3 16h3a2 2 0 012 2v3"/></svg>
        </button>
      </div>
      <div class="toolbar-right">
        <button class="toolbar-btn primary" @click="save">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 01-2-2V5a2 2 0 012-2h11l5 5v11a2 2 0 01-2 2z"/><polyline points="17,21 17,13 7,13 7,21"/><polyline points="7,3 7,8 15,8"/></svg>
          <span>保存</span>
        </button>
      </div>
    </div>
    <div class="editor-main" :class="currentViewMode">
      <div class="visual-panel">
        <div class="canvas-area">
          <div ref="editorContainer" class="gjs-editor-wrapper" :class="'device-' + currentDevice"></div>
        </div>
        <Transition name="panel-slide-left">
          <div class="blocks-panel floating" v-show="showBlocks && !showCodeFullscreen" :style="{ left: blocksPanelPos.x + 'px', top: blocksPanelPos.y + 'px' }">
            <div class="panel-header draggable" @mousedown="startDragPanel($event, 'blocks')">
              <span>组件</span>
              <button class="btn-close-panel" @click.stop="showBlocks = false">×</button>
            </div>
            <div class="panel-options">
              <label class="toggle-option">
                <input type="checkbox" v-model="componentGuideEnabled">
                <span class="toggle-switch"></span>
                <span class="toggle-label">组件指引</span>
              </label>
            </div>
            <div ref="blocksContainer" class="panel-content"></div>
          </div>
        </Transition>
        <Transition name="panel-slide-left">
          <div class="styles-panel floating" v-show="showStyles && !showCodeFullscreen" :style="{ left: stylesPanelPos.x + 'px', top: stylesPanelPos.y + 'px' }">
            <div class="panel-header draggable" @mousedown="startDragPanel($event, 'styles')">
              <span>样式</span>
              <button class="btn-close-panel" @click.stop="showStyles = false">×</button>
            </div>
            <div class="panel-options">
              <label class="toggle-option">
                <input type="checkbox" v-model="editModeEnabled">
                <span class="toggle-switch"></span>
                <span class="toggle-label">编辑模式</span>
              </label>
            </div>
            <div ref="stylesContainer" class="panel-content"></div>
            <div ref="traitsContainer" class="panel-content"></div>
          </div>
        </Transition>
        <Transition name="panel-slide-left">
          <div class="code-edit-panel floating" v-show="showStyles && editModeEnabled && selectedComponentCode && !showCodeFullscreen" :style="{ left: (stylesPanelPos.x + 340) + 'px', top: stylesPanelPos.y + 'px' }">
            <div class="panel-header">
              <span>组件代码</span>
            </div>
            <div class="code-edit-tabs">
              <button :class="['code-edit-tab', { active: editCodeTab === 'html' }]" @click="editCodeTab = 'html'">HTML</button>
              <button :class="['code-edit-tab', { active: editCodeTab === 'css' }]" @click="editCodeTab = 'css'">CSS</button>
            </div>
            <div class="code-edit-content">
              <CodeEditor
                v-show="editCodeTab === 'html'"
                :modelValue="selectedComponentCode?.html || ''"
                language="html"
                @update:modelValue="updateComponentHtml"
              />
              <CodeEditor
                v-show="editCodeTab === 'css'"
                :modelValue="selectedComponentCode?.css || ''"
                language="css"
                @update:modelValue="updateComponentCss"
              />
            </div>
          </div>
        </Transition>
        <div v-if="componentGuideEnabled && hoveredBlockInfo" class="block-tooltip" :style="{ left: blockTooltipPos.x + 'px', top: blockTooltipPos.y + 'px' }">
          <div class="tooltip-title">{{ hoveredBlockInfo.label }}</div>
          <div class="tooltip-desc">{{ hoveredBlockInfo.description }}</div>
          <div class="tooltip-usage">适用场景: {{ hoveredBlockInfo.usage }}</div>
        </div>
      </div>
      <Transition name="code-fade">
        <div v-if="showCodeFullscreen" class="code-panel fullscreen">
          <div class="code-panel-header">
            <div class="code-tabs">
              <button :class="['code-tab', { active: codeTab === 'html' }]" @click="codeTab = 'html'">HTML</button>
              <button :class="['code-tab', { active: codeTab === 'css' }]" @click="codeTab = 'css'">CSS</button>
            </div>
            <div class="code-actions">
              <button class="btn-format" @click="formatCurrentCode" title="格式化代码">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M4 7h16M4 12h10M4 17h16"/></svg>
              </button>
            </div>
            <button class="btn-close-code" @click="closeCodeFullscreen">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M18 6L6 18M6 6l12 12"/></svg>
            </button>
          </div>
          <div class="code-editor-wrapper">
            <CodeEditor
              v-show="codeTab === 'html'"
              :modelValue="htmlCode"
              language="html"
              @update:modelValue="updateHtmlCode"
            />
            <CodeEditor
              v-show="codeTab === 'css'"
              :modelValue="cssCode"
              language="css"
              @update:modelValue="updateCssCode"
            />
          </div>
        </div>
      </Transition>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import grapesjs from 'grapesjs'
import 'grapesjs/dist/css/grapes.min.css'
import gjsPresetWebpage from 'grapesjs-preset-webpage'
import CodeEditor from '@/components/common/CodeEditor.vue'
import { formatHtml, formatCss } from '@/utils/codeFormatter'
const props = withDefaults(defineProps<{
  modelValue?: string
  viewMode?: string
  hideToolbar?: boolean
  externalViewsContainer?: string
}>(), {
  modelValue: '',
  viewMode: 'visual',
  hideToolbar: false,
  externalViewsContainer: ''
})
const emit = defineEmits<{
  'update:modelValue': [value: string]
  save: [content: string]
  fullscreen: [value: boolean]
  change: []
}>()
const editorContainer = ref<HTMLElement | null>(null)
const blocksContainer = ref<HTMLElement | null>(null)
const stylesContainer = ref<HTMLElement | null>(null)
const traitsContainer = ref<HTMLElement | null>(null)
let editor: any = null
const internalViewMode = ref('visual')
const codeTab = ref('html')
const htmlCode = ref('')
const cssCode = ref('')
const showBlocks = ref(false)
const showStyles = ref(false)
const isFullscreen = ref(false)
const showCodeFullscreen = ref(false)
let codeUpdateTimer: ReturnType<typeof setTimeout> | null = null
const blocksPanelPos = ref({ x: 12, y: 12 })
const stylesPanelPos = ref({ x: 12, y: 12 })
const currentDevice = ref('fullscreen')
const showClearConfirm = ref(false)
const showRatioInput = ref(false)
const customRatioW = ref(16)
const customRatioH = ref(9)
const componentGuideEnabled = ref(false)
const editModeEnabled = ref(false)
const editCodeTab = ref('html')
const selectedComponentCode = ref<{ html: string; css: string } | null>(null)
const hoveredBlockInfo = ref<{ label: string; description: string; usage: string } | null>(null)
const blockTooltipPos = ref({ x: 0, y: 0 })
let isDraggingPanel = false
let dragPanelType = ''
let dragStartX = 0
let dragStartY = 0
let dragStartPosX = 0
let dragStartPosY = 0
const currentViewMode = computed(() => {
  return props.hideToolbar ? props.viewMode : internalViewMode.value
})
function initEditor() {
  if (!editorContainer.value) return
  editor = grapesjs.init({
    container: editorContainer.value,
    fromElement: false,
    height: '100%',
    width: 'auto',
    storageManager: false,
    panels: { defaults: [] },
    blockManager: {
      appendTo: blocksContainer.value || undefined
    },
    styleManager: {
      appendTo: stylesContainer.value || undefined
    },
    traitManager: {
      appendTo: traitsContainer.value || undefined
    },
    deviceManager: {
      devices: [
        { name: '默认', width: '' }
      ]
    },
    i18n: {
      locale: 'zh',
      messages: {
        zh: {
          styleManager: {
            empty: '选择一个元素后进行样式编辑',
            layer: '图层',
            fileButton: '图片',
            sectors: {
              general: '通用',
              layout: '布局',
              typography: '排版',
              decorations: '装饰',
              extra: '其他',
              flex: '弹性布局',
              dimension: '尺寸'
            },
            properties: {
              'float': '浮动',
              'display': '显示',
              'position': '定位',
              'top': '上',
              'right': '右',
              'left': '左',
              'bottom': '下',
              'width': '宽度',
              'height': '高度',
              'max-width': '最大宽度',
              'max-height': '最大高度',
              'margin': '外边距',
              'margin-top': '上外边距',
              'margin-right': '右外边距',
              'margin-bottom': '下外边距',
              'margin-left': '左外边距',
              'padding': '内边距',
              'padding-top': '上内边距',
              'padding-right': '右内边距',
              'padding-bottom': '下内边距',
              'padding-left': '左内边距',
              'font-family': '字体',
              'font-size': '字号',
              'font-weight': '字重',
              'letter-spacing': '字间距',
              'color': '颜色',
              'line-height': '行高',
              'text-align': '对齐',
              'text-decoration': '装饰',
              'text-shadow': '文字阴影',
              'border-radius': '圆角',
              'border': '边框',
              'box-shadow': '阴影',
              'background': '背景',
              'background-color': '背景色',
              'opacity': '透明度',
              'transition': '过渡',
              'perspective': '透视',
              'transform': '变换'
            }
          },
          traitManager: {
            empty: '选择一个元素后编辑属性',
            label: '组件属性',
            traits: {
              labels: {
                id: 'ID',
                alt: '替代文本',
                title: '标题',
                href: '链接'
              }
            }
          },
          blockManager: {
            labels: {
              'column1': '一列',
              'column2': '两列',
              'column3': '三列',
              'column3-7': '两列 3/7',
              'text': '文本',
              'link': '链接',
              'image': '图片',
              'video': '视频',
              'map': '地图'
            },
            categories: {
              'Basic': '基础',
              'Layout': '布局',
              'Extra': '其他'
            }
          },
          panels: {
            buttons: {
              titles: {
                'preview': '预览',
                'fullscreen': '全屏',
                'export-template': '查看代码',
                'open-sm': '样式管理器',
                'open-tm': '属性',
                'open-layers': '图层',
                'open-blocks': '组件'
              }
            }
          },
          domComponents: {
            names: {
              '': '元素',
              wrapper: '主体',
              text: '文本',
              comment: '注释',
              image: '图片',
              video: '视频',
              label: '标签',
              link: '链接',
              map: '地图',
              tfoot: '表格底部',
              tbody: '表格主体',
              thead: '表格头部',
              table: '表格',
              row: '表格行',
              cell: '表格单元格',
              default: '组件',
              section: '区块',
              container: '容器',
              column: '列'
            }
          },
          assetManager: {
            addButton: '添加图片',
            inputPlh: '输入图片地址',
            modalTitle: '选择图片',
            uploadTitle: '拖放图片到此处或点击上传'
          },
          selectorManager: {
            label: '选择器',
            selected: '已选中',
            emptyState: '选择一个元素',
            states: {
              hover: '悬停',
              active: '激活',
              focus: '聚焦',
              'nth-of-type(2n)': '偶数'
            }
          },
          deviceManager: {
            device: '设备',
            devices: {
              desktop: '桌面端',
              tablet: '平板',
              mobileLandscape: '手机横屏',
              mobilePortrait: '手机竖屏'
            }
          },
          commands: {
            'core:undo': '撤销',
            'core:redo': '重做',
            'core:canvas-clear': '清空画布',
            'core:component-delete': '删除组件',
            'core:copy': '复制',
            'core:paste': '粘贴',
            'core:component-outline': '组件轮廓',
            'core:fullscreen': '全屏',
            'core:preview': '预览',
            'core:open-code': '查看代码'
          },
          richTextEditor: {
            bold: '粗体',
            italic: '斜体',
            underline: '下划线',
            strikethrough: '删除线',
            link: '链接',
            wrap: '换行'
          },
          layerManager: {
            layers: '图层',
            noElement: '无元素'
          },
          storageManager: {
            recover: '恢复'
          },
          modal: {
            title: '标题',
            close: '关闭'
          }
        }
      }
    },
    plugins: [gjsPresetWebpage],
    pluginsOpts: {
      'grapesjs-preset-webpage': {
        blocksBasicOpts: {
          blocks: ['column1', 'column2', 'column3', 'column3-7', 'text', 'link', 'image', 'video', 'map'],
          flexGrid: true
        },
        formsOpts: false,
        navbarOpts: false,
        countdownOpts: false,
        exportOpts: false,
        aviaryOpts: false,
        filestackOpts: false
      }
    },
    canvas: {
      styles: [
        'https://cdn.jsdelivr.net/npm/github-markdown-css@5.5.1/github-markdown-light.min.css'
      ]
    }
  })
  addCustomBlocks()
  if (props.modelValue) {
    setContent(props.modelValue)
  }
  editor.on('change:changesCount', () => {
    syncCode()
    emit('change')
  })
  editor.on('component:selected', () => {
    showBlocks.value = false
    showStyles.value = true
    updateSelectedComponentCode()
  })
  editor.on('component:deselected', () => {
    updateSelectedComponentCode()
  })
  editor.Commands.add('core:canvas-clear', {
    run() {
      openClearConfirm()
    }
  })
  syncCode()
  window.addEventListener('resize', updateCanvasSize)
  initDeviceButtons()
  setupBlockHoverListeners()
}
function addCustomBlocks() {
  const bm = editor.BlockManager
  bm.add('heading', {
    label: '标题',
    category: '基础',
    content: '<h2>标题文字</h2>',
    attributes: { class: 'gjs-block-heading' }
  })
  bm.add('paragraph', {
    label: '段落',
    category: '基础',
    content: '<p>这是一段文字内容，点击编辑...</p>',
    attributes: { class: 'gjs-block-paragraph' }
  })
  bm.add('divider', {
    label: '分隔线',
    category: '基础',
    content: '<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">',
    attributes: { class: 'gjs-block-divider' }
  })
  bm.add('quote', {
    label: '引用',
    category: '基础',
    content: '<blockquote style="border-left: 4px solid #ddd; padding-left: 16px; margin: 16px 0; color: #666;">引用内容...</blockquote>',
    attributes: { class: 'gjs-block-quote' }
  })
  bm.add('card', {
    label: '卡片',
    category: '布局',
    content: `
      <div style="border: 1px solid #e0e0e0; border-radius: 8px; padding: 16px; margin: 10px 0; background: #fff;">
        <h3 style="margin: 0 0 8px 0;">卡片标题</h3>
        <p style="margin: 0; color: #666;">卡片内容描述...</p>
      </div>
    `,
    attributes: { class: 'gjs-block-card' }
  })
  bm.add('info-box', {
    label: '信息框',
    category: '布局',
    content: `
      <div style="background: #e3f2fd; border-left: 4px solid #2196f3; padding: 12px 16px; margin: 10px 0; border-radius: 0 4px 4px 0;">
        <strong style="color: #1565c0;">提示</strong>
        <p style="margin: 8px 0 0 0; color: #1976d2;">这是一条提示信息...</p>
      </div>
    `,
    attributes: { class: 'gjs-block-info' }
  })
  bm.add('warning-box', {
    label: '警告框',
    category: '布局',
    content: `
      <div style="background: #fff3e0; border-left: 4px solid #ff9800; padding: 12px 16px; margin: 10px 0; border-radius: 0 4px 4px 0;">
        <strong style="color: #e65100;">警告</strong>
        <p style="margin: 8px 0 0 0; color: #ef6c00;">这是一条警告信息...</p>
      </div>
    `,
    attributes: { class: 'gjs-block-warning' }
  })
  bm.add('image-text', {
    label: '图文混排',
    category: '布局',
    content: `
      <div style="display: flex; gap: 16px; align-items: flex-start; margin: 16px 0;">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='150' height='150'%3E%3Crect fill='%23ddd' width='150' height='150'/%3E%3Ctext x='50%25' y='50%25' dominant-baseline='middle' text-anchor='middle' fill='%23999'%3E150×150%3C/text%3E%3C/svg%3E" style="width: 150px; height: 150px; object-fit: cover; border-radius: 8px; flex-shrink: 0;">
        <div>
          <h3 style="margin: 0 0 8px 0;">标题</h3>
          <p style="margin: 0; color: #666; line-height: 1.6;">这里是描述文字，可以详细介绍图片相关的内容...</p>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-image-text' }
  })
  bm.add('button', {
    label: '按钮',
    category: '基础',
    content: '<button style="display: inline-block; padding: 10px 24px; background: #4f46e5; color: #fff; border: none; border-radius: 6px; font-size: 14px; cursor: pointer;">按钮文字</button>',
    attributes: { class: 'gjs-block-button' }
  })
  bm.add('button-outline', {
    label: '边框按钮',
    category: '基础',
    content: '<button style="display: inline-block; padding: 10px 24px; background: transparent; color: #4f46e5; border: 1px solid #4f46e5; border-radius: 6px; font-size: 14px; cursor: pointer;">边框按钮</button>',
    attributes: { class: 'gjs-block-button-outline' }
  })
  bm.add('button-group', {
    label: '按钮组',
    category: '基础',
    content: `
      <div style="display: flex; gap: 8px; margin: 10px 0;">
        <button style="padding: 10px 24px; background: #4f46e5; color: #fff; border: none; border-radius: 6px; font-size: 14px; cursor: pointer;">主按钮</button>
        <button style="padding: 10px 24px; background: transparent; color: #4f46e5; border: 1px solid #4f46e5; border-radius: 6px; font-size: 14px; cursor: pointer;">次按钮</button>
      </div>
    `,
    attributes: { class: 'gjs-block-button-group' }
  })
  bm.add('list-ul', {
    label: '无序列表',
    category: '基础',
    content: `
      <ul style="margin: 10px 0; padding-left: 24px; line-height: 1.8;">
        <li>列表项一</li>
        <li>列表项二</li>
        <li>列表项三</li>
      </ul>
    `,
    attributes: { class: 'gjs-block-list-ul' }
  })
  bm.add('list-ol', {
    label: '有序列表',
    category: '基础',
    content: `
      <ol style="margin: 10px 0; padding-left: 24px; line-height: 1.8;">
        <li>第一步</li>
        <li>第二步</li>
        <li>第三步</li>
      </ol>
    `,
    attributes: { class: 'gjs-block-list-ol' }
  })
  bm.add('badge', {
    label: '徽章',
    category: '基础',
    content: '<span style="display: inline-block; padding: 4px 12px; background: #4f46e5; color: #fff; border-radius: 999px; font-size: 12px;">徽章</span>',
    attributes: { class: 'gjs-block-badge' }
  })
  bm.add('tag-group', {
    label: '标签组',
    category: '基础',
    content: `
      <div style="display: flex; flex-wrap: wrap; gap: 8px; margin: 10px 0;">
        <span style="padding: 4px 12px; background: #e0e7ff; color: #4f46e5; border-radius: 999px; font-size: 12px;">标签一</span>
        <span style="padding: 4px 12px; background: #dcfce7; color: #16a34a; border-radius: 999px; font-size: 12px;">标签二</span>
        <span style="padding: 4px 12px; background: #fef3c7; color: #d97706; border-radius: 999px; font-size: 12px;">标签三</span>
      </div>
    `,
    attributes: { class: 'gjs-block-tag-group' }
  })
  bm.add('avatar', {
    label: '头像',
    category: '基础',
    content: '<div style="width: 48px; height: 48px; border-radius: 50%; background: #e5e7eb; display: flex; align-items: center; justify-content: center; color: #6b7280; font-size: 18px;">A</div>',
    attributes: { class: 'gjs-block-avatar' }
  })
  bm.add('grid-2', {
    label: '两列网格',
    category: '布局',
    content: `
      <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin: 10px 0;">
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px;">列 1</div>
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px;">列 2</div>
      </div>
    `,
    attributes: { class: 'gjs-block-grid-2' }
  })
  bm.add('grid-3', {
    label: '三列网格',
    category: '布局',
    content: `
      <div style="display: grid; grid-template-columns: 1fr 1fr 1fr; gap: 16px; margin: 10px 0;">
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px;">列 1</div>
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px;">列 2</div>
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px;">列 3</div>
      </div>
    `,
    attributes: { class: 'gjs-block-grid-3' }
  })
  bm.add('flex-row', {
    label: '弹性行',
    category: '布局',
    content: `
      <div style="display: flex; gap: 16px; align-items: center; margin: 10px 0;">
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px; flex: 1;">项目 1</div>
        <div style="padding: 16px; background: #f3f4f6; border-radius: 8px; flex: 1;">项目 2</div>
      </div>
    `,
    attributes: { class: 'gjs-block-flex-row' }
  })
  bm.add('section', {
    label: '区块',
    category: '布局',
    content: `
      <section style="padding: 32px; margin: 20px 0; background: #f9fafb; border-radius: 12px;">
        <h2 style="margin: 0 0 16px 0;">区块标题</h2>
        <p style="margin: 0; color: #6b7280;">区块内容描述文字...</p>
      </section>
    `,
    attributes: { class: 'gjs-block-section' }
  })
  bm.add('hero', {
    label: '主视觉',
    category: '布局',
    content: `
      <div style="text-align: center; padding: 48px 24px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px; color: #fff; margin: 20px 0;">
        <h1 style="margin: 0 0 16px 0; font-size: 32px;">欢迎标题</h1>
        <p style="margin: 0 0 24px 0; opacity: 0.9;">这里是副标题或描述文字</p>
        <button style="padding: 12px 32px; background: #fff; color: #667eea; border: none; border-radius: 6px; font-size: 16px; cursor: pointer;">立即开始</button>
      </div>
    `,
    attributes: { class: 'gjs-block-hero' }
  })
  bm.add('tabs', {
    label: '选项卡',
    category: '交互',
    content: `
      <div class="tabs-container" style="margin: 20px 0;">
        <div class="tabs-header" style="display: flex; border-bottom: 1px solid #e5e7eb;">
          <button class="tab-btn active" style="padding: 12px 24px; background: none; border: none; border-bottom: 2px solid #4f46e5; color: #4f46e5; font-size: 14px; cursor: pointer;" onclick="this.parentElement.querySelectorAll('.tab-btn').forEach(b=>{b.classList.remove('active');b.style.borderBottomColor='transparent';b.style.color='#6b7280'});this.classList.add('active');this.style.borderBottomColor='#4f46e5';this.style.color='#4f46e5';this.parentElement.parentElement.querySelectorAll('.tab-panel').forEach((p,i)=>p.style.display=i===0?'block':'none')">选项一</button>
          <button class="tab-btn" style="padding: 12px 24px; background: none; border: none; border-bottom: 2px solid transparent; color: #6b7280; font-size: 14px; cursor: pointer;" onclick="this.parentElement.querySelectorAll('.tab-btn').forEach(b=>{b.classList.remove('active');b.style.borderBottomColor='transparent';b.style.color='#6b7280'});this.classList.add('active');this.style.borderBottomColor='#4f46e5';this.style.color='#4f46e5';this.parentElement.parentElement.querySelectorAll('.tab-panel').forEach((p,i)=>p.style.display=i===1?'block':'none')">选项二</button>
          <button class="tab-btn" style="padding: 12px 24px; background: none; border: none; border-bottom: 2px solid transparent; color: #6b7280; font-size: 14px; cursor: pointer;" onclick="this.parentElement.querySelectorAll('.tab-btn').forEach(b=>{b.classList.remove('active');b.style.borderBottomColor='transparent';b.style.color='#6b7280'});this.classList.add('active');this.style.borderBottomColor='#4f46e5';this.style.color='#4f46e5';this.parentElement.parentElement.querySelectorAll('.tab-panel').forEach((p,i)=>p.style.display=i===2?'block':'none')">选项三</button>
        </div>
        <div class="tabs-content">
          <div class="tab-panel" style="padding: 20px 0; display: block;">选项一的内容区域</div>
          <div class="tab-panel" style="padding: 20px 0; display: none;">选项二的内容区域</div>
          <div class="tab-panel" style="padding: 20px 0; display: none;">选项三的内容区域</div>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-tabs' }
  })
  bm.add('accordion', {
    label: '折叠卡片',
    category: '交互',
    content: `
      <div class="accordion" style="margin: 20px 0; border: 1px solid #e5e7eb; border-radius: 8px; overflow: hidden;">
        <div class="accordion-item" style="border-bottom: 1px solid #e5e7eb;">
          <div class="accordion-header" style="padding: 16px; background: #f9fafb; cursor: pointer; display: flex; justify-content: space-between; align-items: center;" onclick="const content=this.nextElementSibling;const arrow=this.querySelector('.arrow');if(content.style.maxHeight){content.style.maxHeight=null;arrow.style.transform='rotate(0deg)'}else{content.style.maxHeight=content.scrollHeight+'px';arrow.style.transform='rotate(180deg)'}">
            <span style="font-weight: 500;">折叠项标题一</span>
            <span class="arrow" style="transition: transform 0.3s;">▼</span>
          </div>
          <div class="accordion-content" style="max-height: 0; overflow: hidden; transition: max-height 0.3s ease;">
            <div style="padding: 16px;">这是折叠项一的内容，点击标题可以展开或收起。</div>
          </div>
        </div>
        <div class="accordion-item" style="border-bottom: 1px solid #e5e7eb;">
          <div class="accordion-header" style="padding: 16px; background: #f9fafb; cursor: pointer; display: flex; justify-content: space-between; align-items: center;" onclick="const content=this.nextElementSibling;const arrow=this.querySelector('.arrow');if(content.style.maxHeight){content.style.maxHeight=null;arrow.style.transform='rotate(0deg)'}else{content.style.maxHeight=content.scrollHeight+'px';arrow.style.transform='rotate(180deg)'}">
            <span style="font-weight: 500;">折叠项标题二</span>
            <span class="arrow" style="transition: transform 0.3s;">▼</span>
          </div>
          <div class="accordion-content" style="max-height: 0; overflow: hidden; transition: max-height 0.3s ease;">
            <div style="padding: 16px;">这是折叠项二的内容，点击标题可以展开或收起。</div>
          </div>
        </div>
        <div class="accordion-item">
          <div class="accordion-header" style="padding: 16px; background: #f9fafb; cursor: pointer; display: flex; justify-content: space-between; align-items: center;" onclick="const content=this.nextElementSibling;const arrow=this.querySelector('.arrow');if(content.style.maxHeight){content.style.maxHeight=null;arrow.style.transform='rotate(0deg)'}else{content.style.maxHeight=content.scrollHeight+'px';arrow.style.transform='rotate(180deg)'}">
            <span style="font-weight: 500;">折叠项标题三</span>
            <span class="arrow" style="transition: transform 0.3s;">▼</span>
          </div>
          <div class="accordion-content" style="max-height: 0; overflow: hidden; transition: max-height 0.3s ease;">
            <div style="padding: 16px;">这是折叠项三的内容，点击标题可以展开或收起。</div>
          </div>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-accordion' }
  })
  bm.add('collapse', {
    label: '单个折叠',
    category: '交互',
    content: `
      <div class="collapse-item" style="margin: 10px 0; border: 1px solid #e5e7eb; border-radius: 8px; overflow: hidden;">
        <div class="collapse-header" style="padding: 16px; background: #f9fafb; cursor: pointer; display: flex; justify-content: space-between; align-items: center;" onclick="const content=this.nextElementSibling;const arrow=this.querySelector('.arrow');if(content.style.maxHeight){content.style.maxHeight=null;arrow.style.transform='rotate(0deg)'}else{content.style.maxHeight=content.scrollHeight+'px';arrow.style.transform='rotate(180deg)'}">
          <span style="font-weight: 500;">点击展开/收起</span>
          <span class="arrow" style="transition: transform 0.3s;">▼</span>
        </div>
        <div class="collapse-content" style="max-height: 0; overflow: hidden; transition: max-height 0.3s ease;">
          <div style="padding: 16px;">这里是可折叠的内容区域。</div>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-collapse' }
  })
  bm.add('tooltip', {
    label: '提示气泡',
    category: '交互',
    content: `
      <span style="position: relative; display: inline-block; cursor: help; border-bottom: 1px dashed #6b7280;" onmouseenter="this.querySelector('.tooltip-text').style.visibility='visible';this.querySelector('.tooltip-text').style.opacity='1'" onmouseleave="this.querySelector('.tooltip-text').style.visibility='hidden';this.querySelector('.tooltip-text').style.opacity='0'">
        悬停查看提示
        <span class="tooltip-text" style="visibility: hidden; opacity: 0; position: absolute; bottom: 125%; left: 50%; transform: translateX(-50%); background: #1f2937; color: #fff; padding: 8px 12px; border-radius: 6px; font-size: 12px; white-space: nowrap; transition: opacity 0.3s; z-index: 100;">这是提示内容</span>
      </span>
    `,
    attributes: { class: 'gjs-block-tooltip' }
  })
  bm.add('progress', {
    label: '进度条',
    category: '交互',
    content: `
      <div style="margin: 10px 0;">
        <div style="display: flex; justify-content: space-between; margin-bottom: 4px; font-size: 14px;">
          <span>进度</span>
          <span>75%</span>
        </div>
        <div style="height: 8px; background: #e5e7eb; border-radius: 999px; overflow: hidden;">
          <div style="width: 75%; height: 100%; background: #4f46e5; border-radius: 999px;"></div>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-progress' }
  })
  bm.add('table', {
    label: '表格',
    category: '数据',
    content: `
      <table style="width: 100%; border-collapse: collapse; margin: 10px 0;">
        <thead>
          <tr style="background: #f9fafb;">
            <th style="padding: 12px; text-align: left; border-bottom: 2px solid #e5e7eb;">列标题1</th>
            <th style="padding: 12px; text-align: left; border-bottom: 2px solid #e5e7eb;">列标题2</th>
            <th style="padding: 12px; text-align: left; border-bottom: 2px solid #e5e7eb;">列标题3</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">数据1</td>
            <td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">数据2</td>
            <td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">数据3</td>
          </tr>
          <tr>
            <td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">数据4</td>
            <td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">数据5</td>
            <td style="padding: 12px; border-bottom: 1px solid #e5e7eb;">数据6</td>
          </tr>
        </tbody>
      </table>
    `,
    attributes: { class: 'gjs-block-table' }
  })
  bm.add('code-block', {
    label: '代码块',
    category: '数据',
    content: `
      <pre style="background: #1f2937; color: #e5e7eb; padding: 16px; border-radius: 8px; overflow-x: auto; margin: 10px 0; font-family: 'Consolas', 'Monaco', monospace; font-size: 13px; line-height: 1.5;"><code>function hello() {
  console.log("Hello World!");
}</code></pre>
    `,
    attributes: { class: 'gjs-block-code' }
  })
  bm.add('keyboard', {
    label: '键盘按键',
    category: '数据',
    content: '<kbd style="display: inline-block; padding: 4px 8px; background: #f3f4f6; border: 1px solid #d1d5db; border-radius: 4px; font-family: monospace; font-size: 12px; box-shadow: 0 2px 0 #d1d5db;">Ctrl</kbd> + <kbd style="display: inline-block; padding: 4px 8px; background: #f3f4f6; border: 1px solid #d1d5db; border-radius: 4px; font-family: monospace; font-size: 12px; box-shadow: 0 2px 0 #d1d5db;">C</kbd>',
    attributes: { class: 'gjs-block-kbd' }
  })
  bm.add('mark', {
    label: '高亮文本',
    category: '基础',
    content: '<p>这是一段包含<mark style="background: #fef08a; padding: 2px 4px; border-radius: 2px;">高亮标记</mark>的文本。</p>',
    attributes: { class: 'gjs-block-mark' }
  })
  bm.add('abbr', {
    label: '缩写',
    category: '基础',
    content: '<p>使用<abbr title="超文本标记语言" style="text-decoration: underline dotted; cursor: help;">HTML</abbr>创建网页。</p>',
    attributes: { class: 'gjs-block-abbr' }
  })
  bm.add('details', {
    label: '详情摘要',
    category: '交互',
    content: `
      <details style="margin: 10px 0; border: 1px solid #e5e7eb; border-radius: 8px;">
        <summary style="padding: 12px 16px; cursor: pointer; background: #f9fafb; border-radius: 8px; font-weight: 500;">点击查看详情</summary>
        <div style="padding: 16px;">这里是详细内容，点击上方摘要可以展开或收起。</div>
      </details>
    `,
    attributes: { class: 'gjs-block-details' }
  })
  bm.add('figure', {
    label: '图片说明',
    category: '媒体',
    content: `
      <figure style="margin: 20px 0; text-align: center;">
        <img src="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='300' height='200'%3E%3Crect fill='%23ddd' width='300' height='200'/%3E%3Ctext x='50%25' y='50%25' dominant-baseline='middle' text-anchor='middle' fill='%23999'%3E300×200%3C/text%3E%3C/svg%3E" style="max-width: 100%; border-radius: 8px;" alt="图片描述">
        <figcaption style="margin-top: 8px; color: #6b7280; font-size: 14px;">图片说明文字</figcaption>
      </figure>
    `,
    attributes: { class: 'gjs-block-figure' }
  })
  bm.add('video-embed', {
    label: '视频容器',
    category: '媒体',
    content: `
      <div style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden; border-radius: 8px; margin: 10px 0; background: #000;">
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); color: #fff; text-align: center;">
          <div style="width: 60px; height: 60px; border-radius: 50%; background: rgba(255,255,255,0.2); display: flex; align-items: center; justify-content: center; margin: 0 auto 10px;">▶</div>
          <span>视频占位符</span>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-video-embed' }
  })
  bm.add('audio', {
    label: '音频播放器',
    category: '媒体',
    content: `
      <div style="padding: 16px; background: #f9fafb; border-radius: 8px; margin: 10px 0;">
        <div style="display: flex; align-items: center; gap: 12px;">
          <button style="width: 40px; height: 40px; border-radius: 50%; background: #4f46e5; color: #fff; border: none; cursor: pointer; display: flex; align-items: center; justify-content: center;">▶</button>
          <div style="flex: 1;">
            <div style="height: 4px; background: #e5e7eb; border-radius: 2px;">
              <div style="width: 30%; height: 100%; background: #4f46e5; border-radius: 2px;"></div>
            </div>
          </div>
          <span style="font-size: 12px; color: #6b7280;">1:23 / 4:56</span>
        </div>
      </div>
    `,
    attributes: { class: 'gjs-block-audio' }
  })
  bm.add('iframe', {
    label: '嵌入框架',
    category: '媒体',
    content: `
      <div style="border: 2px dashed #e5e7eb; border-radius: 8px; padding: 40px; text-align: center; margin: 10px 0; background: #f9fafb;">
        <div style="color: #6b7280;">iframe 嵌入区域</div>
        <div style="font-size: 12px; color: #9ca3af; margin-top: 4px;">在代码视图中编辑 src 属性</div>
      </div>
    `,
    attributes: { class: 'gjs-block-iframe' }
  })
  bm.add('input-text', {
    label: '输入框',
    category: '表单',
    content: `
      <div style="margin: 10px 0;">
        <label style="display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500;">标签文字</label>
        <input type="text" placeholder="请输入..." style="width: 100%; padding: 10px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; box-sizing: border-box;">
      </div>
    `,
    attributes: { class: 'gjs-block-input' }
  })
  bm.add('textarea', {
    label: '文本域',
    category: '表单',
    content: `
      <div style="margin: 10px 0;">
        <label style="display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500;">标签文字</label>
        <textarea placeholder="请输入内容..." style="width: 100%; padding: 10px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; min-height: 100px; box-sizing: border-box; resize: vertical;"></textarea>
      </div>
    `,
    attributes: { class: 'gjs-block-textarea' }
  })
  bm.add('select', {
    label: '下拉框',
    category: '表单',
    content: `
      <div style="margin: 10px 0;">
        <label style="display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500;">选择选项</label>
        <select style="width: 100%; padding: 10px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; background: #fff;">
          <option>选项一</option>
          <option>选项二</option>
          <option>选项三</option>
        </select>
      </div>
    `,
    attributes: { class: 'gjs-block-select' }
  })
  bm.add('checkbox', {
    label: '复选框',
    category: '表单',
    content: `
      <div style="margin: 10px 0;">
        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
          <input type="checkbox" style="width: 18px; height: 18px; accent-color: #4f46e5;">
          <span style="font-size: 14px;">复选框选项</span>
        </label>
      </div>
    `,
    attributes: { class: 'gjs-block-checkbox' }
  })
  bm.add('radio', {
    label: '单选框',
    category: '表单',
    content: `
      <div style="margin: 10px 0; display: flex; flex-direction: column; gap: 8px;">
        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
          <input type="radio" name="radio-group" style="width: 18px; height: 18px; accent-color: #4f46e5;">
          <span style="font-size: 14px;">选项一</span>
        </label>
        <label style="display: flex; align-items: center; gap: 8px; cursor: pointer;">
          <input type="radio" name="radio-group" style="width: 18px; height: 18px; accent-color: #4f46e5;">
          <span style="font-size: 14px;">选项二</span>
        </label>
      </div>
    `,
    attributes: { class: 'gjs-block-radio' }
  })
  bm.add('switch', {
    label: '开关',
    category: '表单',
    content: `
      <label style="display: inline-flex; align-items: center; gap: 8px; cursor: pointer;">
        <span style="position: relative; width: 44px; height: 24px;">
          <input type="checkbox" style="opacity: 0; width: 0; height: 0;" onchange="this.nextElementSibling.style.background=this.checked?'#4f46e5':'#d1d5db';this.nextElementSibling.querySelector('span').style.transform=this.checked?'translateX(20px)':'translateX(0)'">
          <span style="position: absolute; inset: 0; background: #d1d5db; border-radius: 999px; transition: 0.3s;">
            <span style="position: absolute; top: 2px; left: 2px; width: 20px; height: 20px; background: #fff; border-radius: 50%; transition: 0.3s; box-shadow: 0 1px 3px rgba(0,0,0,0.2);"></span>
          </span>
        </span>
        <span style="font-size: 14px;">开关选项</span>
      </label>
    `,
    attributes: { class: 'gjs-block-switch' }
  })
  bm.add('form', {
    label: '表单容器',
    category: '表单',
    content: `
      <form style="padding: 24px; background: #f9fafb; border-radius: 12px; margin: 10px 0;">
        <h3 style="margin: 0 0 20px 0;">表单标题</h3>
        <div style="margin-bottom: 16px;">
          <label style="display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500;">姓名</label>
          <input type="text" placeholder="请输入姓名" style="width: 100%; padding: 10px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; box-sizing: border-box;">
        </div>
        <div style="margin-bottom: 16px;">
          <label style="display: block; margin-bottom: 6px; font-size: 14px; font-weight: 500;">邮箱</label>
          <input type="email" placeholder="请输入邮箱" style="width: 100%; padding: 10px 12px; border: 1px solid #d1d5db; border-radius: 6px; font-size: 14px; box-sizing: border-box;">
        </div>
        <button type="submit" style="width: 100%; padding: 12px; background: #4f46e5; color: #fff; border: none; border-radius: 6px; font-size: 14px; cursor: pointer;">提交</button>
      </form>
    `,
    attributes: { class: 'gjs-block-form' }
  })
  bm.add('address', {
    label: '地址',
    category: '基础',
    content: `
      <address style="font-style: normal; padding: 16px; background: #f9fafb; border-radius: 8px; margin: 10px 0;">
        <strong>联系地址</strong><br>
        北京市朝阳区某某街道123号<br>
        邮编: 100000<br>
        电话: 010-12345678
      </address>
    `,
    attributes: { class: 'gjs-block-address' }
  })
  bm.add('time', {
    label: '时间日期',
    category: '基础',
    content: '<p>发布于 <time datetime="2024-01-01" style="color: #6b7280;">2024年1月1日</time></p>',
    attributes: { class: 'gjs-block-time' }
  })
  bm.add('definition-list', {
    label: '定义列表',
    category: '基础',
    content: `
      <dl style="margin: 10px 0;">
        <dt style="font-weight: 600; margin-bottom: 4px;">术语一</dt>
        <dd style="margin: 0 0 16px 20px; color: #6b7280;">术语一的定义和解释。</dd>
        <dt style="font-weight: 600; margin-bottom: 4px;">术语二</dt>
        <dd style="margin: 0 0 16px 20px; color: #6b7280;">术语二的定义和解释。</dd>
      </dl>
    `,
    attributes: { class: 'gjs-block-dl' }
  })
  bm.add('sub-sup', {
    label: '上下标',
    category: '基础',
    content: '<p>H<sub>2</sub>O 是水的化学式，E=mc<sup>2</sup> 是质能方程。</p>',
    attributes: { class: 'gjs-block-subsup' }
  })
  bm.add('meter', {
    label: '度量',
    category: '数据',
    content: `
      <div style="margin: 10px 0;">
        <label style="display: block; margin-bottom: 6px; font-size: 14px;">存储空间使用情况</label>
        <div style="height: 20px; background: #e5e7eb; border-radius: 4px; overflow: hidden;">
          <div style="width: 60%; height: 100%; background: linear-gradient(90deg, #22c55e 0%, #22c55e 50%, #eab308 80%, #ef4444 100%); background-size: 166.67% 100%;"></div>
        </div>
        <div style="font-size: 12px; color: #6b7280; margin-top: 4px;">已使用 6GB / 10GB</div>
      </div>
    `,
    attributes: { class: 'gjs-block-meter' }
  })
  bm.add('blockquote-cite', {
    label: '引用带出处',
    category: '基础',
    content: `
      <figure style="margin: 20px 0;">
        <blockquote style="border-left: 4px solid #4f46e5; padding: 16px 20px; margin: 0; background: #f5f3ff; border-radius: 0 8px 8px 0; font-style: italic; color: #4b5563;">
          "生活不是等待暴风雨过去，而是学会在雨中跳舞。"
        </blockquote>
        <figcaption style="margin-top: 8px; text-align: right; color: #6b7280; font-size: 14px;">— <cite>维维安·格林</cite></figcaption>
      </figure>
    `,
    attributes: { class: 'gjs-block-blockquote-cite' }
  })
}
function setContent(html: string) {
  if (!editor) return
  const styleMatch = html.match(/<style[^>]*>([\s\S]*?)<\/style>/i)
  const css = styleMatch ? styleMatch[1] : ''
  const bodyMatch = html.match(/<body[^>]*>([\s\S]*?)<\/body>/i)
  const bodyContent = bodyMatch ? bodyMatch[1] : html.replace(/<style[^>]*>[\s\S]*?<\/style>/gi, '')
  editor.setComponents(bodyContent)
  editor.setStyle(css)
  syncCode()
}
function getFullHtml() {
  if (!editor) return ''
  const html = editor.getHtml()
  const css = editor.getCss()
  return css ? `<style>${css}</style>\n${html}` : html
}
function syncCode() {
  if (!editor) return
  htmlCode.value = editor.getHtml()
  cssCode.value = editor.getCss()
}
function onCodeChange() {
  if (codeUpdateTimer) {
    clearTimeout(codeUpdateTimer)
  }
  codeUpdateTimer = setTimeout(() => {
    if (editor) {
      editor.setComponents(htmlCode.value)
      editor.setStyle(cssCode.value)
    }
  }, 500)
}
function updateHtmlCode(value: string) {
  htmlCode.value = value
  onCodeChange()
}
function updateCssCode(value: string) {
  cssCode.value = value
  onCodeChange()
}
function setViewMode(mode: string) {
  internalViewMode.value = mode
  if (mode === 'code' || mode === 'split') {
    syncCode()
  }
  if (editor) {
    nextTick(() => {
      editor.refresh()
    })
  }
}
function undo() {
  if (editor) {
    editor.UndoManager.undo()
  }
}
function redo() {
  if (editor) {
    editor.UndoManager.redo()
  }
}
function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value
  emit('fullscreen', isFullscreen.value)
}
function save() {
  const content = getFullHtml()
  emit('update:modelValue', content)
  emit('save', content)
}
function openCodeFullscreen() {
  syncCode()
  showCodeFullscreen.value = true
}
function closeCodeFullscreen() {
  showCodeFullscreen.value = false
}
function formatCurrentCode() {
  if (codeTab.value === 'html') {
    htmlCode.value = formatHtml(htmlCode.value)
  } else {
    cssCode.value = formatCss(cssCode.value)
  }
  onCodeChange()
}
function toggleBlocksPanel() {
  showBlocks.value = !showBlocks.value
}
function toggleStylesPanel() {
  showStyles.value = !showStyles.value
}
function setDevice(device: string) {
  currentDevice.value = device
  updateCanvasSize()
  updateDeviceButtonsState()
}
function updateCanvasSize() {
  if (!editor) return
  const canvasEl = document.querySelector('.canvas-area') as HTMLElement
  if (!canvasEl) return
  const wrapperEl = canvasEl.querySelector('.gjs-editor-wrapper') as HTMLElement
  if (!wrapperEl) return
  const containerWidth = canvasEl.clientWidth - 24
  const containerHeight = canvasEl.clientHeight - 24
  let canvasWidth = containerWidth
  let canvasHeight = containerHeight
  if (currentDevice.value === 'desktop') {
    const ratio = 16 / 9
    if (containerWidth / containerHeight > ratio) {
      canvasHeight = containerHeight
      canvasWidth = canvasHeight * ratio
    } else {
      canvasWidth = containerWidth
      canvasHeight = canvasWidth / ratio
    }
  } else if (currentDevice.value === 'mobile') {
    const ratio = 9 / 20
    if (containerWidth / containerHeight > ratio) {
      canvasHeight = containerHeight
      canvasWidth = canvasHeight * ratio
    } else {
      canvasWidth = containerWidth
      canvasHeight = canvasWidth / ratio
    }
  } else if (currentDevice.value === 'custom') {
    const ratio = customRatioW.value / customRatioH.value
    if (containerWidth / containerHeight > ratio) {
      canvasHeight = containerHeight
      canvasWidth = canvasHeight * ratio
    } else {
      canvasWidth = containerWidth
      canvasHeight = canvasWidth / ratio
    }
  }
  wrapperEl.style.width = currentDevice.value === 'fullscreen' ? '100%' : canvasWidth + 'px'
  wrapperEl.style.height = currentDevice.value === 'fullscreen' ? '100%' : canvasHeight + 'px'
  nextTick(() => {
    editor.refresh()
  })
}
function applyCustomRatio() {
  if (customRatioW.value > 0 && customRatioH.value > 0) {
    currentDevice.value = 'custom'
    showRatioInput.value = false
    updateCanvasSize()
    updateDeviceButtonsState()
  }
}
function openClearConfirm() {
  showClearConfirm.value = true
}
function confirmClearCanvas() {
  if (editor) {
    editor.DomComponents.clear()
    editor.CssComposer.clear()
  }
  showClearConfirm.value = false
}
function updateComponentHtml(value: string) {
  if (!editor || !selectedComponentCode.value) return
  const selected = editor.getSelected()
  if (selected) {
    selected.set('content', '')
    selected.components(value)
  }
}
function updateComponentCss(value: string) {
  if (!editor || !selectedComponentCode.value) return
  const selected = editor.getSelected()
  if (selected) {
    const rule = editor.CssComposer.getRule(`#${selected.getId()}`) || editor.CssComposer.setRule(`#${selected.getId()}`, {})
    if (rule) {
      const styleObj: Record<string, string> = {}
      value.replace(/([a-z-]+)\s*:\s*([^;]+)/gi, (_, prop, val) => {
        styleObj[prop.trim()] = val.trim()
        return ''
      })
      rule.setStyle(styleObj)
    }
  }
}
function updateSelectedComponentCode() {
  if (!editor || !editModeEnabled.value) {
    selectedComponentCode.value = null
    return
  }
  const selected = editor.getSelected()
  if (selected) {
    const html = selected.toHTML()
    const id = selected.getId()
    const rule = editor.CssComposer.getRule(`#${id}`)
    const css = rule ? rule.toCSS() : ''
    selectedComponentCode.value = { html, css }
  } else {
    selectedComponentCode.value = null
  }
}
const blockDescriptions: Record<string, { description: string; usage: string }> = {
  'heading': { description: '用于显示标题文字，支持h1-h6不同级别', usage: '页面标题、章节标题、文章标题' },
  'paragraph': { description: '用于显示普通段落文字', usage: '正文内容、描述说明、介绍文字' },
  'divider': { description: '水平分隔线，用于分隔内容区域', usage: '区分不同内容区块、视觉分割' },
  'quote': { description: '引用样式，带有左边框装饰', usage: '引用名言、重要提示、摘要' },
  'card': { description: '卡片容器，带边框和圆角', usage: '信息展示、产品卡片、个人简介' },
  'info-box': { description: '蓝色提示信息框', usage: '提示信息、帮助说明、重要通知' },
  'warning-box': { description: '橙色警告信息框', usage: '警告提醒、注意事项、风险提示' },
  'image-text': { description: '左图右文的混排布局', usage: '产品介绍、文章配图、个人简介' },
  'button': { description: '实心按钮，主色调背景', usage: '主要操作、表单提交、链接跳转' },
  'button-outline': { description: '边框按钮，透明背景', usage: '次要操作、取消按钮、辅助操作' },
  'button-group': { description: '按钮组合，多个按钮并排', usage: '操作组合、选择器、工具栏' },
  'list-ul': { description: '无序列表，圆点标记', usage: '功能列表、特点介绍、注意事项' },
  'list-ol': { description: '有序列表，数字标记', usage: '步骤说明、排名列表、操作流程' },
  'badge': { description: '徽章标签，圆角胶囊形状', usage: '状态标识、标签分类、计数显示' },
  'tag-group': { description: '标签组，多个不同颜色标签', usage: '分类标签、技能标签、属性展示' },
  'avatar': { description: '头像容器，圆形', usage: '用户头像、联系人、团队成员' },
  'grid-2': { description: '两列等宽网格布局', usage: '对比展示、双列内容、图文并排' },
  'grid-3': { description: '三列等宽网格布局', usage: '特性展示、产品列表、团队成员' },
  'flex-row': { description: '弹性行布局，自适应宽度', usage: '导航栏、工具条、内容分布' },
  'section': { description: '区块容器，带背景和内边距', usage: '页面分区、功能模块、内容分组' },
  'hero': { description: '主视觉区域，渐变背景居中布局', usage: '首屏展示、活动banner、欢迎页' },
  'tabs': { description: '选项卡切换，点击显示不同内容', usage: '内容分类、功能切换、设置面板' },
  'accordion': { description: '折叠卡片组，可展开收起', usage: 'FAQ列表、目录导航、详情折叠' },
  'collapse': { description: '单个折叠面板', usage: '详情展开、高级选项、附加信息' },
  'tooltip': { description: '悬停提示气泡', usage: '术语解释、操作提示、帮助信息' },
  'progress': { description: '进度条显示', usage: '完成进度、加载状态、数据占比' },
  'table': { description: '数据表格，带表头', usage: '数据展示、对比列表、价格表' },
  'code-block': { description: '代码块，等宽字体深色背景', usage: '代码示例、命令行、配置信息' },
  'keyboard': { description: '键盘按键样式', usage: '快捷键说明、操作指引' },
  'mark': { description: '高亮标记文本', usage: '重点强调、关键词高亮' },
  'abbr': { description: '缩写词，悬停显示全称', usage: '术语解释、缩写说明' },
  'details': { description: '原生详情摘要元素', usage: '可展开详情、FAQ、隐藏内容' },
  'figure': { description: '图片带说明文字', usage: '插图配文、作品展示' },
  'video-embed': { description: '视频容器占位', usage: '视频嵌入、播放器占位' },
  'audio': { description: '音频播放器样式', usage: '音乐播放、语音消息' },
  'iframe': { description: '嵌入框架占位', usage: '第三方内容嵌入、地图嵌入' },
  'input-text': { description: '文本输入框带标签', usage: '表单输入、搜索框' },
  'textarea': { description: '多行文本输入域', usage: '评论输入、内容编辑' },
  'select': { description: '下拉选择框', usage: '选项选择、筛选条件' },
  'checkbox': { description: '复选框选项', usage: '多选设置、同意条款' },
  'radio': { description: '单选框选项组', usage: '单选设置、类型选择' },
  'switch': { description: '开关切换按钮', usage: '开关设置、状态切换' },
  'form': { description: '表单容器，包含输入和提交', usage: '用户注册、信息收集' },
  'address': { description: '地址信息展示', usage: '联系地址、公司信息' },
  'time': { description: '时间日期显示', usage: '发布时间、事件日期' },
  'definition-list': { description: '定义列表，术语和解释配对', usage: '词汇解释、参数说明' },
  'sub-sup': { description: '上标下标文本', usage: '化学公式、数学公式、注释' },
  'meter': { description: '度量显示，带渐变颜色', usage: '用量显示、健康指数' },
  'blockquote-cite': { description: '引用带出处署名', usage: '名人名言、文献引用' }
}
function setupBlockHoverListeners() {
  nextTick(() => {
    const blocks = document.querySelectorAll('.blocks-panel .gjs-block')
    blocks.forEach(block => {
      block.addEventListener('mouseenter', (e: Event) => {
        if (!componentGuideEnabled.value) return
        const target = e.currentTarget as HTMLElement
        const label = target.querySelector('.gjs-block-label')?.textContent || ''
        const blockId = target.getAttribute('data-block') || target.getAttribute('title') || ''
        const info = Object.entries(blockDescriptions).find(([key, val]) => label.includes(key) || blockId.includes(key))
        if (info) {
          hoveredBlockInfo.value = { label, ...info[1] }
          const rect = target.getBoundingClientRect()
          blockTooltipPos.value = { x: rect.right + 10, y: rect.top }
        }
      })
      block.addEventListener('mouseleave', () => {
        hoveredBlockInfo.value = null
      })
    })
  })
}
function initDeviceButtons() {
  nextTick(() => {
    const commandsPanel = document.querySelector('.gjs-pn-panel.gjs-pn-commands')
    if (!commandsPanel) return
    const existingToolbar = commandsPanel.querySelector('.custom-device-toolbar')
    if (existingToolbar) existingToolbar.remove()
    const deviceToolbar = document.createElement('div')
    deviceToolbar.className = 'custom-device-toolbar'
    deviceToolbar.innerHTML = `
      <button class="device-btn ${currentDevice.value === 'desktop' ? 'active' : ''}" data-device="desktop" title="桌面端 16:9">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2"/><path d="M8 21h8M12 17v4"/></svg>
      </button>
      <button class="device-btn ${currentDevice.value === 'mobile' ? 'active' : ''}" data-device="mobile" title="手机 9:20">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="5" y="2" width="14" height="20" rx="2"/><path d="M12 18h.01"/></svg>
      </button>
      <button class="device-btn ${currentDevice.value === 'fullscreen' ? 'active' : ''}" data-device="fullscreen" title="铺满">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M8 3H5a2 2 0 00-2 2v3m18 0V5a2 2 0 00-2-2h-3m0 18h3a2 2 0 002-2v-3M3 16v3a2 2 0 002 2h3"/></svg>
      </button>
      <button class="device-btn ${currentDevice.value === 'custom' ? 'active' : ''}" data-device="custom" title="自定义比例">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><path d="M9 3v18M3 9h18"/></svg>
      </button>
    `
    commandsPanel.insertBefore(deviceToolbar, commandsPanel.firstChild)
    deviceToolbar.querySelectorAll('.device-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const device = (e.currentTarget as HTMLElement).getAttribute('data-device')
        if (device === 'custom') {
          showRatioInput.value = true
        } else if (device) {
          setDevice(device)
          updateDeviceButtonsState()
        }
      })
    })
  })
}
function updateDeviceButtonsState() {
  const buttons = document.querySelectorAll('.custom-device-toolbar .device-btn')
  buttons.forEach(btn => {
    const device = btn.getAttribute('data-device')
    btn.classList.toggle('active', device === currentDevice.value)
  })
}
function startDragPanel(e: MouseEvent, panelType: string) {
  if ((e.target as HTMLElement).closest('.btn-close-panel')) return
  isDraggingPanel = true
  dragPanelType = panelType
  dragStartX = e.clientX
  dragStartY = e.clientY
  if (panelType === 'blocks') {
    dragStartPosX = blocksPanelPos.value.x
    dragStartPosY = blocksPanelPos.value.y
  } else {
    dragStartPosX = stylesPanelPos.value.x
    dragStartPosY = stylesPanelPos.value.y
  }
  document.addEventListener('mousemove', onDragPanel)
  document.addEventListener('mouseup', stopDragPanel)
  e.preventDefault()
}
function onDragPanel(e: MouseEvent) {
  if (!isDraggingPanel) return
  const deltaX = e.clientX - dragStartX
  const deltaY = e.clientY - dragStartY
  const canvas = document.querySelector('.visual-panel') as HTMLElement
  if (!canvas) return
  const canvasRect = canvas.getBoundingClientRect()
  const panelEl = dragPanelType === 'blocks'
    ? document.querySelector('.blocks-panel.floating') as HTMLElement
    : document.querySelector('.styles-panel.floating') as HTMLElement
  if (!panelEl) return
  const panelWidth = panelEl.offsetWidth
  const panelHeight = panelEl.offsetHeight
  const maxX = canvasRect.width - panelWidth
  const maxY = canvasRect.height - panelHeight
  const newX = Math.max(0, Math.min(maxX, dragStartPosX + deltaX))
  const newY = Math.max(0, Math.min(maxY, dragStartPosY + deltaY))
  if (dragPanelType === 'blocks') {
    blocksPanelPos.value = { x: newX, y: newY }
  } else {
    stylesPanelPos.value = { x: newX, y: newY }
  }
}
function stopDragPanel() {
  isDraggingPanel = false
  dragPanelType = ''
  document.removeEventListener('mousemove', onDragPanel)
  document.removeEventListener('mouseup', stopDragPanel)
}
watch(() => props.modelValue, (newVal) => {
  if (editor && newVal !== getFullHtml()) {
    setContent(newVal)
  }
})
watch(() => props.viewMode, (newVal) => {
  if (props.hideToolbar && newVal !== internalViewMode.value) {
    setViewMode(newVal)
  }
})
watch(editModeEnabled, () => {
  updateSelectedComponentCode()
})
onMounted(() => {
  nextTick(() => {
    initEditor()
  })
})
onBeforeUnmount(() => {
  if (editor) {
    editor.destroy()
  }
  if (codeUpdateTimer) {
    clearTimeout(codeUpdateTimer)
  }
  window.removeEventListener('resize', updateCanvasSize)
})
defineExpose({
  showBlocks,
  showStyles,
  openCodeFullscreen,
  closeCodeFullscreen,
  toggleBlocksPanel,
  toggleStylesPanel,
  getFullHtml,
  setContent
})
</script>
<style scoped>
.grapes-editor-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-secondary);
  border-radius: 8px;
  overflow: hidden;
}
.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg-tertiary);
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}
.toolbar-left,
.toolbar-center,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}
.toolbar-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 10px;
  background: transparent;
  border: 1px solid transparent;
  border-radius: 6px;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}
.toolbar-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.toolbar-btn.active {
  background: rgba(var(--primary-rgb), 0.15);
  color: var(--primary-color);
  border-color: var(--primary-color);
}
.toolbar-btn.primary {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}
.toolbar-btn.primary:hover {
  opacity: 0.9;
}
.toolbar-btn svg {
  width: 16px;
  height: 16px;
}
.toolbar-divider {
  width: 1px;
  height: 20px;
  background: var(--border-color);
  margin: 0 8px;
}
.editor-main {
  flex: 1;
  display: flex;
  min-height: 0;
  position: relative;
}
.visual-panel {
  flex: 1;
  display: flex;
  min-width: 0;
  position: relative;
}
.blocks-panel.floating,
.styles-panel.floating {
  position: absolute;
  width: 330px;
  max-height: 80%;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  z-index: 20;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}
.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
}
.panel-header.draggable {
  cursor: move;
  user-select: none;
}
.btn-close-panel {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 14px;
}
.btn-close-panel:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}
.canvas-area {
  flex: 1;
  min-width: 0;
  background: #1a1a1a;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
}
.panel-options {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border-color);
}
.toggle-option {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 12px;
  color: var(--text-secondary);
}
.toggle-option input {
  display: none;
}
.toggle-switch {
  position: relative;
  width: 32px;
  height: 18px;
  background: var(--border-color);
  border-radius: 999px;
  transition: all 0.2s;
}
.toggle-switch::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 14px;
  height: 14px;
  background: #fff;
  border-radius: 50%;
  transition: all 0.2s;
}
.toggle-option input:checked + .toggle-switch {
  background: var(--primary-color);
}
.toggle-option input:checked + .toggle-switch::after {
  left: 16px;
}
.toggle-label {
  color: var(--text-primary);
}
.code-edit-panel.floating {
  width: 350px;
  height: 400px;
}
.code-edit-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-color);
}
.code-edit-tab {
  flex: 1;
  padding: 8px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: all 0.15s;
}
.code-edit-tab:hover {
  color: var(--text-primary);
}
.code-edit-tab.active {
  color: var(--primary-color);
  background: rgba(var(--primary-rgb), 0.1);
}
.code-edit-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}
.block-tooltip {
  position: fixed;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 12px;
  max-width: 250px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  z-index: 1000;
  pointer-events: none;
}
.tooltip-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 8px;
}
.tooltip-desc {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  line-height: 1.5;
}
.tooltip-usage {
  font-size: 11px;
  color: var(--primary-color);
}
:deep(.custom-device-toolbar) {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-right: 8px;
  padding-right: 8px;
  border-right: 1px solid var(--border-color);
}
:deep(.custom-device-toolbar .device-btn) {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 4px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
:deep(.custom-device-toolbar .device-btn:hover) {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}
:deep(.custom-device-toolbar .device-btn.active) {
  background: rgba(var(--primary-rgb), 0.2);
  color: var(--primary-color);
}
:deep(.custom-device-toolbar .device-btn svg) {
  width: 16px;
  height: 16px;
}
:deep(.gjs-pn-panel.gjs-pn-commands) {
  width: 100% !important;
  box-sizing: border-box;
}
.gjs-editor-wrapper {
  height: 100%;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  transition: width 0.3s, height 0.3s;
}
.gjs-editor-wrapper.device-fullscreen {
  width: 100% !important;
  height: 100% !important;
}
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
.confirm-modal,
.ratio-modal {
  background: var(--bg-primary);
  border-radius: 12px;
  padding: 20px;
  min-width: 300px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}
.confirm-title,
.ratio-title {
  font-size: 16px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 12px;
}
.confirm-content {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 20px;
}
.confirm-actions,
.ratio-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
.btn-cancel,
.btn-confirm {
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.15s;
}
.btn-cancel {
  background: transparent;
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}
.btn-cancel:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.btn-confirm {
  background: var(--primary-color);
  border: 1px solid var(--primary-color);
  color: white;
}
.btn-confirm:hover {
  opacity: 0.9;
}
.ratio-inputs {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
}
.ratio-input {
  width: 80px;
  padding: 8px 12px;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
  text-align: center;
}
.ratio-input:focus {
  outline: none;
  border-color: var(--primary-color);
}
.ratio-separator {
  font-size: 16px;
  color: var(--text-secondary);
}
.code-panel {
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  border-left: 1px solid var(--border-color);
}
.code-panel.fullscreen {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 50;
  border-left: none;
  border-radius: 0;
}
.code-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 12px;
  border-bottom: 1px solid var(--border-color);
  flex-shrink: 0;
}
.btn-close-code {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 50%;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
.btn-close-code:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.btn-close-code svg {
  width: 18px;
  height: 18px;
}
.code-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
  margin-right: 8px;
}
.btn-format {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: 6px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
.btn-format:hover {
  background: var(--bg-hover);
  color: var(--primary-color);
}
.btn-format svg {
  width: 16px;
  height: 16px;
}
.code-tabs {
  display: flex;
  flex-shrink: 0;
}
.code-tab {
  padding: 8px 16px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  font-size: 12px;
  cursor: pointer;
  border-bottom: 2px solid transparent;
  transition: all 0.15s;
}
.code-tab:hover {
  color: var(--text-primary);
}
.code-tab.active {
  color: var(--primary-color);
  border-bottom-color: var(--primary-color);
}
.code-editor-wrapper {
  flex: 1;
  min-height: 0;
}
:deep(.gjs-one-bg) {
  background-color: var(--bg-secondary);
}
:deep(.gjs-two-color) {
  color: var(--text-primary);
}
:deep(.gjs-three-bg) {
  background-color: var(--bg-tertiary);
}
:deep(.gjs-four-color),
:deep(.gjs-four-color-h:hover) {
  color: var(--primary-color);
}
:deep(.gjs-block) {
  width: auto;
  min-height: auto;
  padding: 8px;
  margin: 4px;
  border-radius: 6px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
}
:deep(.gjs-block:hover) {
  border-color: var(--primary-color);
}
:deep(.gjs-block__media) {
  display: none;
}
:deep(.gjs-block-label) {
  font-size: 12px;
}
:deep(.gjs-block-categories) {
  display: flex;
  flex-wrap: wrap;
}
:deep(.gjs-block-category) {
  width: 100%;
}
:deep(.gjs-block-category:first-child .gjs-title) {
  display: none;
}
:deep(.gjs-block-category:first-child .gjs-blocks-c) {
  display: flex;
  flex-wrap: wrap;
}
:deep(.gjs-category-title),
:deep(.gjs-title) {
  font-size: 12px;
  padding: 8px;
  border-bottom: 1px solid var(--border-color);
}
:deep(.gjs-sm-sector-title) {
  font-size: 12px;
  padding: 8px;
  border-bottom: 1px solid var(--border-color);
}
:deep(.gjs-field) {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  border-radius: 4px;
}
:deep(.gjs-field input) {
  color: var(--text-primary);
}
:deep(.gjs-pn-panel.gjs-pn-views-container) {
  display: none !important;
}
:deep(.gjs-pn-panel.gjs-pn-views) {
  display: none !important;
}
:deep(.gjs-pn-buttons) {
  justify-content: flex-start !important;
}
:deep(.gjs-pn-panel.gjs-pn-options) {
  position: absolute !important;
  right: 0 !important;
  left: auto !important;
}
:deep(.gjs-pn-btn[title="查看代码"]),
:deep(.gjs-pn-btn.fa-code),
:deep(.gjs-pn-btn[title="Import"]),
:deep(.gjs-pn-btn[title="Full-screen"]),
:deep(.gjs-pn-btn[title="Preview"]),
:deep(.gjs-pn-btn[title="View code"]),
:deep(.gjs-pn-btn[title="View components"]),
:deep(.gjs-pn-btn[title="导入"]),
:deep(.gjs-pn-btn[title="全屏"]),
:deep(.gjs-pn-btn[title="预览"]),
:deep(.gjs-pn-btn[title="代码"]),
:deep(.gjs-pn-btn.fa-download),
:deep(.gjs-pn-btn.fa-arrows-alt),
:deep(.gjs-pn-btn.fa-eye) {
  display: none !important;
}
:deep(.gjs-pn-panel.gjs-pn-devices-c) {
  display: none !important;
}
:deep(.gjs-cv-canvas) {
  width: 100% !important;
}
:deep(.gjs-toolbar) {
  background: var(--bg-secondary);
  border-radius: 6px;
  border: 1px solid var(--border-color);
}
:deep(.gjs-toolbar-item) {
  font-size: 12px;
}
:deep(.gjs-toolbar-item[title="Move"]) {
  font-size: 0;
}
:deep(.gjs-toolbar-item[title="Move"])::after {
  content: '移动';
  font-size: 12px;
}
:deep(.gjs-toolbar-item[title="Select parent"]) {
  font-size: 0;
}
:deep(.gjs-toolbar-item[title="Select parent"])::after {
  content: '选择父级';
  font-size: 12px;
}
:deep(.gjs-toolbar-item[title="Delete"]) {
  font-size: 0;
}
:deep(.gjs-toolbar-item[title="Delete"])::after {
  content: '删除';
  font-size: 12px;
}
:deep(.gjs-toolbar-item[title="Copy"]) {
  font-size: 0;
}
:deep(.gjs-toolbar-item[title="Copy"])::after {
  content: '复制';
  font-size: 12px;
}
:deep(.gjs-toolbar-item[title="Drag"]) {
  font-size: 0;
}
:deep(.gjs-toolbar-item[title="Drag"])::after {
  content: '拖动';
  font-size: 12px;
}
:deep(.gjs-rte-action[title="Bold"]) {
  font-size: 0;
}
:deep(.gjs-rte-action[title="Bold"])::after {
  content: 'B';
  font-size: 14px;
  font-weight: bold;
}
:deep(.gjs-rte-action[title="Italic"]) {
  font-size: 0;
}
:deep(.gjs-rte-action[title="Italic"])::after {
  content: 'I';
  font-size: 14px;
  font-style: italic;
}
:deep(.gjs-rte-action[title="Underline"]) {
  font-size: 0;
}
:deep(.gjs-rte-action[title="Underline"])::after {
  content: 'U';
  font-size: 14px;
  text-decoration: underline;
}
:deep(.gjs-rte-action[title="Strikethrough"]) {
  font-size: 0;
}
:deep(.gjs-rte-action[title="Strikethrough"])::after {
  content: 'S';
  font-size: 14px;
  text-decoration: line-through;
}
:deep(.gjs-rte-action[title="Link"]) {
  font-size: 0;
}
:deep(.gjs-rte-action[title="Link"])::after {
  content: '链接';
  font-size: 12px;
}
:deep(.gjs-rte-action[title="Wrap for style"]) {
  font-size: 0;
}
:deep(.gjs-rte-action[title="Wrap for style"])::after {
  content: '样式';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Preview"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Preview"])::after {
  content: '预览';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="View code"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="View code"])::after {
  content: '代码';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Full-screen"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Full-screen"])::after {
  content: '全屏';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Undo"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Undo"])::after {
  content: '撤销';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Redo"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Redo"])::after {
  content: '重做';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Import"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Import"])::after {
  content: '导入';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Clear canvas"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Clear canvas"])::after {
  content: '清空';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Open Style Manager"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Open Style Manager"])::after {
  content: '样式';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Settings"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Settings"])::after {
  content: '设置';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Open Layer Manager"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Open Layer Manager"])::after {
  content: '图层';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Open Blocks"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Open Blocks"])::after {
  content: '组件';
  font-size: 12px;
}
:deep(.gjs-pn-btn[title="Component outline"]) {
  font-size: 0;
}
:deep(.gjs-pn-btn[title="Component outline"])::after {
  content: '轮廓';
  font-size: 12px;
}
:deep(.gjs-clm-tags #gjs-clm-new) {
  font-size: 0;
}
:deep(.gjs-clm-tags #gjs-clm-new::placeholder) {
  font-size: 12px;
}
:deep(.gjs-am-add-asset .gjs-am-add-field) {
  font-size: 0;
}
:deep(.gjs-am-add-asset .gjs-am-add-field::placeholder) {
  font-size: 12px;
}
.panel-slide-left-enter-active,
.panel-slide-left-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.panel-slide-left-enter-from,
.panel-slide-left-leave-to {
  transform: translateX(-20px);
  opacity: 0;
}
.panel-slide-right-enter-active,
.panel-slide-right-leave-active {
  transition: transform 0.2s ease, opacity 0.2s ease;
}
.panel-slide-right-enter-from,
.panel-slide-right-leave-to {
  transform: translateX(20px);
  opacity: 0;
}
.code-fade-enter-active,
.code-fade-leave-active {
  transition: opacity 0.2s ease;
}
.code-fade-enter-from,
.code-fade-leave-to {
  opacity: 0;
}
</style>
