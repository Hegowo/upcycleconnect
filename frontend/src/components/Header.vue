<template>
  <header class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-6 shrink-0">
    <div class="flex items-center gap-3">
      <h1 class="text-base font-semibold text-gray-900">{{ title }}</h1>
    </div>

    <div class="flex items-center gap-4">
      <RouterLink
        v-if="auth.isSuperAdmin"
        to="/admin/providers?status=pending"
        class="relative text-gray-400 hover:text-gray-600 transition"
        :title="t('providers.statusPending')"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6 6 0 10-12 0v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
      </RouterLink>

      <div class="flex items-center gap-2">
        <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-bold" style="background-color: #1B5B88;">
          {{ initials }}
        </div>
        <div class="hidden md:block text-right">
          <p class="text-sm font-medium text-gray-800 leading-none">{{ auth.fullName }}</p>
          <p class="text-xs capitalize mt-0.5" style="color: #1B8848;">{{ auth.role?.replace('_', ' ') }}</p>
        </div>
      </div>

      <button
        @click="handleLogout"
        :disabled="loading"
        class="flex items-center gap-1.5 text-sm text-gray-400 hover:text-red-500 transition-colors px-2 py-1 rounded-lg hover:bg-red-50"
        :title="t('nav.logout')"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
        </svg>
        <span class="hidden md:inline">{{ loading ? t('nav.loggingOut') : t('nav.logout') }}</span>
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()
const loading = ref(false)

const title    = computed(() => route.meta.title || t('nav.administration'))
const initials = computed(() => {
  if (!auth.user) return '?'
  return ((auth.user.first_name?.[0] || '') + (auth.user.last_name?.[0] || '')).toUpperCase()
})

async function handleLogout() {
  loading.value = true
  await auth.logout()
  router.push('/admin/login')
}
</script>
