import { ref, onMounted, onUnmounted } from 'vue'

export interface ResizableOptions {
  minWidth?: number
  maxWidth?: number
  minHeight?: number
  maxHeight?: number
  direction?: 'horizontal' | 'vertical' | 'both'
}

export function useResizable(options: ResizableOptions = {}) {
  const {
    minWidth = 100,
    maxWidth = 800,
    minHeight = 100,
    maxHeight = 600,
    direction = 'horizontal'
  } = options
  const width = ref(200)
  const height = ref(200)
  const isResizing = ref(false)
  const startX = ref(0)
  const startY = ref(0)
  const startWidth = ref(0)
  const startHeight = ref(0)

  function startResize(e: MouseEvent) {
    isResizing.value = true
    startX.value = e.clientX
    startY.value = e.clientY
    startWidth.value = width.value
    startHeight.value = height.value
    e.preventDefault()
  }

  function onMouseMove(e: MouseEvent) {
    if (!isResizing.value) return
    if (direction === 'horizontal' || direction === 'both') {
      const deltaX = e.clientX - startX.value
      const newWidth = startWidth.value + deltaX
      width.value = Math.min(Math.max(newWidth, minWidth), maxWidth)
    }
    if (direction === 'vertical' || direction === 'both') {
      const deltaY = e.clientY - startY.value
      const newHeight = startHeight.value + deltaY
      height.value = Math.min(Math.max(newHeight, minHeight), maxHeight)
    }
  }

  function onMouseUp() {
    isResizing.value = false
  }

  onMounted(() => {
    window.addEventListener('mousemove', onMouseMove)
    window.addEventListener('mouseup', onMouseUp)
  })

  onUnmounted(() => {
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
  })

  return {
    width,
    height,
    isResizing,
    startResize
  }
}
