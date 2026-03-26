<template>
  <header class="h-16 bg-white border-b border-[#e5e7eb] flex items-center justify-between px-6 shrink-0">
    <!-- Page title -->
    <div class="flex items-center gap-3">
      <h1 class="text-base font-semibold text-gray-900">{{ title }}</h1>
    </div>

    <!-- Center search bar -->
    <div class="hidden md:flex items-center flex-1 max-w-md mx-8">
      <div class="relative w-full">
        <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          type="text"
          placeholder="Rechercher..."
          class="w-full pl-9 pr-4 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
        />
      </div>
    </div>

    <!-- Right actions -->
    <div class="flex items-center gap-3">
      <!-- Notification bell -->
      <button
        class="relative p-2 rounded-lg text-gray-400 hover:text-gray-600 hover:bg-gray-50 transition"
        title="Notifications"
      >
        <BellIcon class="w-5 h-5" />
        <span
          v-if="pendingCount > 0"
          class="absolute top-1 right-1 w-2 h-2 rounded-full bg-red-500"
        ></span>
      </button>

      <!-- Separator -->
      <div class="h-6 w-px bg-[#e5e7eb]"></div>

      <!-- User info -->
      <div class="flex items-center gap-2.5">
        <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#001d32;">
          {{ initials }}
        </div>
        <div class="hidden md:block">
          <p class="text-sm font-semibold text-gray-800 leading-none">{{ auth.fullName }}</p>
          <p class="text-xs capitalize mt-0.5 font-medium" style="color:#006d35;">{{ auth.role?.replace('_', ' ') }}</p>
        </div>
      </div>

      <!-- Logout -->
      <button
        @click="handleLogout"
        :disabled="loading"
        class="p-2 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition"
        title="Se déconnecter"
      >
        <ArrowRightOnRectangleIcon class="w-5 h-5" />
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  BellIcon,
  MagnifyingGlassIcon,
  ArrowRightOnRectangleIcon,
} from '@heroicons/vue/24/outline'

const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()
const loading = ref(false)
const pendingCount = ref(0)

const title = computed(() => route.meta.title || 'Administration')

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
