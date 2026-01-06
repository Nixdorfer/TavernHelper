<template>
  <div class="modals-container">
    <ConfirmDialog />
    <NotificationToast />
    <LoginModal v-model="showLogin" @success="onLoginSuccess" />
    <SettingsModal v-model="showSettings" />
  </div>
</template>

<script setup lang="ts">
import { ref, provide } from 'vue'
import ConfirmDialog from '../common/ConfirmDialog.vue'
import NotificationToast from '../common/NotificationToast.vue'
import LoginModal from './LoginModal.vue'
import SettingsModal from './SettingsModal.vue'
const showLogin = ref(false)
const showSettings = ref(false)
function openLogin() {
  showLogin.value = true
}
function openSettings() {
  showSettings.value = true
}
function onLoginSuccess() {
  showLogin.value = false
}
provide('modals', {
  openLogin,
  openSettings
})
defineExpose({
  showLogin,
  showSettings,
  openLogin,
  openSettings
})
</script>

<style scoped>
.modals-container {
  position: relative;
  z-index: 1000;
}
</style>
