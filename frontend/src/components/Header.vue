<template>
  <header class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6 shrink-0">
    <h1 class="text-lg font-semibold text-gray-900">{{ title }}</h1>
    <button
      @click="handleLogout"
      :disabled="loading"
      class="flex items-center gap-2 text-sm text-gray-500 hover:text-gray-700 transition-colors"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
      </svg>
      {{ loading ? 'Déconnexion...' : 'Se déconnecter' }}
    </button>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()
const loading = ref(false)

const title = computed(() => route.meta.title || 'Administration')

async function handleLogout() {
  loading.value = true
  await auth.logout()
  router.push('/login')
}
</script>
