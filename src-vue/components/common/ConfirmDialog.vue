<template>
  <BaseModal
    :modelValue="confirmStore.isVisible"
    :title="confirmStore.options.title"
    size="sm"
    :closable="true"
    @update:modelValue="handleClose"
    @confirm="confirmStore.confirm"
  >
    <p class="confirm-dialog__message">{{ confirmStore.options.message }}</p>
    <template #footer>
      <BaseButton variant="secondary" @click="confirmStore.cancel">
        {{ confirmStore.options.cancelText }}
      </BaseButton>
      <BaseButton
        :variant="confirmStore.options.type === 'danger' ? 'danger' : 'primary'"
        @click="confirmStore.confirm"
      >
        {{ confirmStore.options.confirmText }}
      </BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { useConfirmStore } from '@/stores'
import BaseModal from './BaseModal.vue'
import BaseButton from './BaseButton.vue'
const confirmStore = useConfirmStore()
function handleClose(value: boolean) {
  if (!value) {
    confirmStore.cancel()
  }
}
</script>

<style scoped>
.confirm-dialog__message {
  margin: 0;
  font-size: 15px;
  color: var(--text-primary);
  line-height: 1.6;
}
</style>
