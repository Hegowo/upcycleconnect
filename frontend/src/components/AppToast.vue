<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 w-80">
      <TransitionGroup name="toast">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['flex items-start gap-3 p-4 rounded-lg shadow-lg text-sm', typeClass(toast.type)]"
        >
          <span class="flex-1">{{ toast.message }}</span>
          <button @click="remove(toast.id)" class="text-current opacity-60 hover:opacity-100">✕</button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup>import { useToast } from '@/utils/useToast'

const { toasts, remove } = useToast()

function typeClass(type) {
  return {
    success: 'bg-green-600 text-white',
    error:   'bg-red-600 text-white',
    info:    'bg-blue-600 text-white',
  }[type] || 'bg-gray-800 text-white'
}
</script>

<style scoped>
.toast-enter-active, .toast-leave-active { transition: all 0.3s ease; }
.toast-enter-from { opacity: 0; transform: translateX(100%); }
.toast-leave-to   { opacity: 0; transform: translateX(100%); }
</style>
