<template>
  <Teleport to="body">
    <div class="notification-container">
      <TransitionGroup name="notification">
        <div
          v-for="notification in notificationStore.notifications"
          :key="notification.id"
          :class="['notification-toast', `notification-toast--${notification.type}`]"
        >
          <span class="notification-toast__icon">{{ icons[notification.type] }}</span>
          <span class="notification-toast__message">{{ notification.message }}</span>
          <button class="notification-toast__close" @click="notificationStore.remove(notification.id)">×</button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useNotificationStore } from '@/stores'
const notificationStore = useNotificationStore()
const icons = {
  success: '✓',
  error: '✕',
  warning: '⚠',
  info: 'ℹ'
}
</script>

<style scoped>
.notification-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;
}
.notification-toast {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  pointer-events: auto;
  min-width: 280px;
  max-width: 400px;
}
.notification-toast--success {
  border-left: 4px solid #10b981;
}
.notification-toast--error {
  border-left: 4px solid #ef4444;
}
.notification-toast--warning {
  border-left: 4px solid #f59e0b;
}
.notification-toast--info {
  border-left: 4px solid #3b82f6;
}
.notification-toast__icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  font-size: 12px;
  font-weight: bold;
}
.notification-toast--success .notification-toast__icon {
  background: #10b981;
  color: white;
}
.notification-toast--error .notification-toast__icon {
  background: #ef4444;
  color: white;
}
.notification-toast--warning .notification-toast__icon {
  background: #f59e0b;
  color: white;
}
.notification-toast--info .notification-toast__icon {
  background: #3b82f6;
  color: white;
}
.notification-toast__message {
  flex: 1;
  font-size: 14px;
  color: var(--text-primary);
  line-height: 1.4;
}
.notification-toast__close {
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  border-radius: 50%;
  color: var(--text-tertiary);
  cursor: pointer;
  font-size: 18px;
  line-height: 1;
  transition: all 0.15s;
}
.notification-toast__close:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}
.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}
.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
.notification-move {
  transition: transform 0.3s ease;
}
</style>
