<template>
  <BaseModal
    :modelValue="modelValue"
    title="登录"
    size="sm"
    @update:modelValue="emit('update:modelValue', $event)"
    @confirm="handleSubmit"
  >
    <form class="login-form" @submit.prevent="handleSubmit">
      <BaseInput
        v-model="email"
        type="email"
        label="邮箱"
        placeholder="请输入邮箱"
        :error="errors.email"
      />
      <BaseInput
        v-model="password"
        type="password"
        label="密码"
        placeholder="请输入密码"
        :error="errors.password"
      />
      <label class="login-form__remember">
        <input type="checkbox" v-model="remember" />
        <span>记住登录状态</span>
      </label>
    </form>
    <template #footer>
      <BaseButton variant="secondary" @click="emit('update:modelValue', false)">取消</BaseButton>
      <BaseButton variant="primary" :loading="isLoading" @click="handleSubmit">登录</BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores'
import { useNotification } from '@/composables'
import BaseModal from '../common/BaseModal.vue'
import BaseInput from '../common/BaseInput.vue'
import BaseButton from '../common/BaseButton.vue'
defineProps<{
  modelValue: boolean
}>()
const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()
const authStore = useAuthStore()
const { success, error: showError } = useNotification()
const email = ref('')
const password = ref('')
const remember = ref(true)
const isLoading = ref(false)
const errors = reactive({ email: '', password: '' })
function validate(): boolean {
  errors.email = ''
  errors.password = ''
  if (!email.value) {
    errors.email = '请输入邮箱'
    return false
  }
  if (!password.value) {
    errors.password = '请输入密码'
    return false
  }
  return true
}
async function handleSubmit() {
  if (!validate()) return
  isLoading.value = true
  try {
    await authStore.login(email.value, password.value, remember.value)
    success('登录成功')
    emit('update:modelValue', false)
    emit('success')
  } catch (e) {
    showError((e as Error).message || '登录失败')
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.login-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.login-form__remember {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text-secondary);
  cursor: pointer;
}
.login-form__remember input {
  width: 16px;
  height: 16px;
  accent-color: var(--primary-color);
}
</style>
