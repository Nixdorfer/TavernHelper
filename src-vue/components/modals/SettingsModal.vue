<template>
  <BaseModal
    :modelValue="modelValue"
    title="设置"
    size="lg"
    :confirmOnEnter="false"
    @update:modelValue="emit('update:modelValue', $event)"
  >
    <div class="settings-modal">
      <div class="settings-modal__tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          :class="['settings-modal__tab', { 'settings-modal__tab--active': activeTab === tab.id }]"
          @click="activeTab = tab.id"
        >
          {{ tab.label }}
        </button>
      </div>
      <div class="settings-modal__content">
        <div v-if="activeTab === 'general'" class="settings-section">
          <h4 class="settings-section__title">通用设置</h4>
          <div class="settings-item">
            <span class="settings-item__label">无图模式</span>
            <label class="settings-switch">
              <input type="checkbox" :checked="configStore.config.noImageMode" @change="toggleNoImageMode" />
              <span class="settings-switch__slider"></span>
            </label>
          </div>
        </div>
        <div v-if="activeTab === 'api'" class="settings-section">
          <h4 class="settings-section__title">API配置</h4>
          <BaseInput
            :modelValue="configStore.config.bytedanceApiKey"
            label="Bytedance API Key"
            type="password"
            placeholder="请输入API Key"
            @update:modelValue="configStore.update('bytedanceApiKey', $event)"
          />
          <BaseInput
            :modelValue="configStore.config.grokApiKey"
            label="Grok API Key"
            type="password"
            placeholder="请输入API Key"
            @update:modelValue="configStore.update('grokApiKey', $event)"
          />
        </div>
        <div v-if="activeTab === 'connection'" class="settings-section">
          <h4 class="settings-section__title">平台连接</h4>
          <div class="platform-cards">
            <div
              :class="['platform-card', { 'platform-card--selected': selectedPlatform === 'fengyue' }]"
              @click="selectedPlatform = 'fengyue'"
            >
              <div class="platform-card__header">
                <span :class="['platform-dot', authStore.isLoggedIn ? 'platform-dot--online' : 'platform-dot--offline']"></span>
                <span class="platform-card__name">风月</span>
              </div>
              <div v-if="authStore.isLoggedIn" class="platform-card__points">
                {{ authStore.points }} 积分
              </div>
            </div>
            <div class="platform-card platform-card--disabled">
              <div class="platform-card__header">
                <span class="platform-dot platform-dot--offline"></span>
                <span class="platform-card__name">魅魔岛</span>
              </div>
              <div class="platform-card__hint">暂不支持</div>
            </div>
          </div>
          <template v-if="selectedPlatform === 'fengyue'">
            <div v-if="!authStore.isLoggedIn" class="login-form">
              <BaseInput
                v-model="loginEmail"
                label="账号"
                type="text"
                placeholder="请输入邮箱"
              />
              <BaseInput
                v-model="loginPassword"
                label="密码"
                type="password"
                placeholder="请输入密码"
              />
              <label class="remember-checkbox">
                <input type="checkbox" v-model="rememberMe" />
                <span>记住登录状态</span>
              </label>
              <BaseButton variant="primary" :loading="isLoggingIn" @click="handleLogin">登录</BaseButton>
            </div>
            <div v-else class="logged-info">
              <div class="logged-info__user">已登录：{{ authStore.userName }}</div>
              <BaseButton variant="secondary" :loading="isLoggingOut" @click="handleLogout">登出</BaseButton>
            </div>
          </template>
        </div>
      </div>
    </div>
    <template #footer>
      <BaseButton variant="secondary" @click="emit('update:modelValue', false)">关闭</BaseButton>
      <BaseButton variant="primary" :loading="isSaving" @click="saveSettings">保存</BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useConfigStore, useAuthStore } from '@/stores'
import { useNotification } from '@/composables'
import BaseModal from '../common/BaseModal.vue'
import BaseButton from '../common/BaseButton.vue'
import BaseInput from '../common/BaseInput.vue'
defineProps<{
  modelValue: boolean
}>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()
const configStore = useConfigStore()
const authStore = useAuthStore()
const { success, error: showError } = useNotification()
const activeTab = ref('general')
const isSaving = ref(false)
const isLoggingIn = ref(false)
const isLoggingOut = ref(false)
const loginEmail = ref('')
const loginPassword = ref('')
const rememberMe = ref(true)
const selectedPlatform = ref('fengyue')
const tabs = [
  { id: 'general', label: '通用' },
  { id: 'api', label: 'API' },
  { id: 'connection', label: '连接' }
]
async function handleLogin() {
  if (!loginEmail.value || !loginPassword.value) {
    showError('请输入账号和密码')
    return
  }
  isLoggingIn.value = true
  try {
    await authStore.login(loginEmail.value, loginPassword.value, rememberMe.value)
    success('登录成功')
    loginEmail.value = ''
    loginPassword.value = ''
  } catch (e) {
    showError((e as Error).message || '登录失败')
  } finally {
    isLoggingIn.value = false
  }
}
async function handleLogout() {
  isLoggingOut.value = true
  try {
    await authStore.logout()
    success('已登出')
  } catch (e) {
    showError((e as Error).message || '登出失败')
  } finally {
    isLoggingOut.value = false
  }
}
function toggleNoImageMode(e: Event) {
  configStore.update('noImageMode', (e.target as HTMLInputElement).checked)
}
async function saveSettings() {
  isSaving.value = true
  try {
    await configStore.save()
    success('设置已保存')
    emit('update:modelValue', false)
  } catch (e) {
    showError((e as Error).message || '保存失败')
  } finally {
    isSaving.value = false
  }
}
</script>

<style scoped>
.settings-modal {
  display: flex;
  gap: 24px;
  min-height: 400px;
}
.settings-modal__tabs {
  display: flex;
  flex-direction: column;
  gap: 4px;
  width: 120px;
  flex-shrink: 0;
}
.settings-modal__tab {
  padding: 10px 16px;
  border: none;
  background: transparent;
  border-radius: 8px;
  text-align: left;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.15s;
}
.settings-modal__tab:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}
.settings-modal__tab--active {
  background: var(--primary-color);
  color: white;
}
.settings-modal__content {
  flex: 1;
}
.settings-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.settings-section__title {
  margin: 0 0 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}
.settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}
.settings-item__label {
  font-size: 14px;
  color: var(--text-primary);
}
.settings-switch {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
}
.settings-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}
.settings-switch__slider {
  position: absolute;
  cursor: pointer;
  inset: 0;
  background: var(--bg-tertiary);
  border-radius: 24px;
  transition: 0.2s;
}
.settings-switch__slider::before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background: white;
  border-radius: 50%;
  transition: 0.2s;
}
.settings-switch input:checked + .settings-switch__slider {
  background: var(--primary-color);
}
.settings-switch input:checked + .settings-switch__slider::before {
  transform: translateX(20px);
}
.platform-cards {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}
.platform-card {
  flex: 1;
  padding: 16px;
  border-radius: 12px;
  border: 2px solid var(--border-color);
  background: var(--bg-secondary);
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
.platform-card:hover:not(.platform-card--disabled) {
  border-color: var(--primary-color);
}
.platform-card--selected {
  background: rgba(var(--primary-rgb, 76, 175, 80), 0.15);
  border-color: var(--primary-color);
}
.platform-card--disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--bg-tertiary);
}
.platform-card__header {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.platform-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.platform-dot--online {
  background: #4caf50;
}
.platform-dot--offline {
  background: #f44336;
}
.platform-card__name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-primary);
}
.platform-card__points {
  margin-top: 8px;
  display: inline-block;
  padding: 4px 12px;
  background: #4caf50;
  color: white;
  font-size: 12px;
  border-radius: 20px;
}
.platform-card__hint {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-tertiary);
}
.login-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.remember-checkbox {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}
.remember-checkbox input {
  width: 16px;
  height: 16px;
  accent-color: var(--primary-color);
}
.logged-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: var(--bg-secondary);
  border-radius: 8px;
}
.logged-info__user {
  font-size: 14px;
  color: var(--text-primary);
}
</style>
