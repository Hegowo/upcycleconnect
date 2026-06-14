<template>
  <header :class="[
    'h-14 lg:h-16 border-b flex items-center justify-between px-4 md:px-6 shrink-0 transition-colors duration-200',
    theme.isDark
      ? 'bg-[#1e293b] border-[#334155]'
      : 'bg-white border-[#e5e7eb]'
  ]">

    <h1 :class="['text-sm lg:text-base font-semibold truncate', theme.isDark ? 'text-[#f1f5f9]' : 'text-gray-900']">
      {{ title }}
    </h1>

    <div class="hidden md:flex items-center flex-1 max-w-md mx-8">
      <div class="relative w-full">
        <MagnifyingGlassIcon :class="['absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4', theme.isDark ? 'text-[#475569]' : 'text-gray-400']" />
        <input
          type="text"
          :placeholder="t('common.search')"
          :class="[
            'w-full pl-9 pr-4 py-2 text-sm border rounded-lg focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition',
            theme.isDark
              ? 'bg-[#0f172a] border-[#334155] text-[#f1f5f9] placeholder-[#475569]'
              : 'bg-[#f8fafc] border-[#e5e7eb] text-gray-900'
          ]"
        />
      </div>
    </div>

    <div class="flex items-center gap-2 md:gap-3">

      <button
        @click="theme.toggle()"
        :class="[
          'p-2 rounded-lg transition',
          theme.isDark
            ? 'text-yellow-400 hover:text-yellow-300 hover:bg-[#263244]'
            : 'text-gray-400 hover:text-gray-600 hover:bg-gray-50'
        ]"
        :title="theme.isDark ? 'Passer en mode clair' : 'Passer en mode sombre'"
      >
        <SunIcon v-if="theme.isDark" class="w-5 h-5" />
        <MoonIcon v-else class="w-5 h-5" />
      </button>

      <NotificationsDropdown variant="admin" :is-dark="theme.isDark" />

      <div :class="['h-6 w-px hidden sm:block', theme.isDark ? 'bg-[#334155]' : 'bg-[#e5e7eb]']"></div>

      <div class="flex items-center gap-2.5">
        <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#001d32;">
          {{ initials }}
        </div>
        <div class="hidden md:block">
          <p :class="['text-sm font-semibold leading-none', theme.isDark ? 'text-[#f1f5f9]' : 'text-gray-800']">{{ auth.fullName }}</p>
          <p class="text-xs capitalize mt-0.5 font-medium" style="color:#006d35;">{{ auth.role?.replace('_', ' ') }}</p>
        </div>
      </div>

      <button
        @click="handleLogout"
        :disabled="loading"
        :class="[
          'p-2 rounded-lg transition',
          theme.isDark
            ? 'text-[#64748b] hover:text-red-400 hover:bg-red-900/20'
            : 'text-gray-400 hover:text-red-500 hover:bg-red-50'
        ]"
        :title="t('nav.logout')"
      >
        <ArrowRightOnRectangleIcon class="w-5 h-5" />
      </button>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
import {
  MagnifyingGlassIcon,
  ArrowRightOnRectangleIcon,
  SunIcon,
  MoonIcon,
} from '@heroicons/vue/24/outline'
import NotificationsDropdown from '@/components/NotificationsDropdown.vue'

const { t } = useI18n()
const route   = useRoute()
const router  = useRouter()
const auth    = useAuthStore()
const theme   = useThemeStore()
const loading = ref(false)

const title = computed(() => route.meta.title || t('nav.administration'))

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
