<template>
  <Suspense>
    <component :is="PlatformApp" />
  </Suspense>
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue'
const PLATFORM = import.meta.env.VITE_PLATFORM || 'pc'
const PlatformApp = PLATFORM === 'mobile'
  ? defineAsyncComponent(() => import('./platforms/mobile/App.vue'))
  : defineAsyncComponent(() => import('./platforms/pc/App.vue'))
</script>
