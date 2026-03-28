<template>
  <aside class="hidden lg:flex w-64 bg-white border-r border-[#e5e7eb] flex-col shrink-0 h-screen">
    
    <div class="flex items-center gap-3 px-4 py-4 border-b border-[#e5e7eb]">
      <img src="/logoentier.png" alt="UpcycleConnect" class="h-9 w-auto object-contain" />
      <p class="text-[10px] font-semibold tracking-widest uppercase text-gray-400 whitespace-nowrap">Admin Terminal</p>
    </div>

    
    <nav class="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        custom
        v-slot="{ isActive, navigate }"
      >
        <button
          @click="navigate()"
          :class="[
            'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 text-left',
            isActive
              ? 'bg-[#f0fdf4] text-[#006d35] border-l-2 border-[#006d35]'
              : 'text-[#374151] hover:bg-[#f9fafb] border-l-2 border-transparent'
          ]"
        >
          <component :is="item.icon" class="w-5 h-5 shrink-0" :class="isActive ? 'text-[#006d35]' : 'text-[#9ca3af]'" />
          <span>{{ item.label }}</span>
        </button>
      </RouterLink>

      <div v-if="auth.isSuperAdmin" class="pt-3 pb-1">
        <p class="px-3 text-[10px] font-semibold uppercase tracking-widest text-gray-400 mb-1">Super Admin</p>
      </div>

      <template v-if="auth.isSuperAdmin">
        <RouterLink
          v-for="item in superAdminItems"
          :key="item.to"
          :to="item.to"
          custom
          v-slot="{ isActive, navigate }"
        >
          <button
            @click="navigate()"
            :class="[
              'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 text-left',
              isActive
                ? 'bg-[#f0fdf4] text-[#006d35] border-l-2 border-[#006d35]'
                : 'text-[#374151] hover:bg-[#f9fafb] border-l-2 border-transparent'
            ]"
          >
            <component :is="item.icon" class="w-5 h-5 shrink-0" :class="isActive ? 'text-[#006d35]' : 'text-[#9ca3af]'" />
            <span>{{ item.label }}</span>
          </button>
        </RouterLink>
      </template>
    </nav>

    
    <div class="px-3 pb-3">
      <button class="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg text-white text-sm font-semibold transition hover:opacity-90" style="background-color:#006d35;">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('nav.newReport') }}
      </button>
    </div>

    
    <div class="px-4 py-4 border-t border-[#e5e7eb]">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-full flex items-center justify-center text-white text-sm font-bold shrink-0" style="background-color:#001d32;">
          {{ initials }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold text-gray-900 truncate leading-none">{{ auth.fullName }}</p>
          <p class="text-xs capitalize mt-0.5 font-medium" style="color:#006d35;">{{ auth.role?.replace('_', ' ') }}</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import {
  Squares2X2Icon,
  UsersIcon,
  BriefcaseIcon,
  TagIcon,
  ClipboardDocumentListIcon,
  CalendarIcon,
  ClockIcon,
  ShieldCheckIcon,
  CircleStackIcon,
  Cog6ToothIcon,
  InboxIcon,
  CodeBracketIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const auth = useAuthStore()

const initials = computed(() => {
  if (!auth.user) return '?'
  return ((auth.user.first_name?.[0] || '') + (auth.user.last_name?.[0] || '')).toUpperCase()
})

const navItems = computed(() => [
  { to: '/admin/dashboard',    icon: Squares2X2Icon,            label: t('nav.dashboard') },
  { to: '/admin/users',        icon: UsersIcon,                 label: t('nav.users') },
  { to: '/admin/providers',    icon: BriefcaseIcon,             label: t('nav.providers') },
  { to: '/admin/categories',   icon: TagIcon,                   label: t('nav.categories') },
  { to: '/admin/prestations',  icon: ClipboardDocumentListIcon, label: t('nav.prestations') },
  { to: '/admin/events',       icon: CalendarIcon,              label: t('nav.events') },
  { to: '/admin/logs',         icon: ClockIcon,                 label: t('nav.logs') },
  { to: '/admin/box-requests', icon: InboxIcon,                 label: t('nav.boxRequests') },
  { to: '/admin/settings',     icon: Cog6ToothIcon,             label: t('nav.settings') },
  { to: '/admin/docs',         icon: CodeBracketIcon,           label: t('nav.apiDocs') },
])

const superAdminItems = computed(() => [
  { to: '/admin/admins',   icon: ShieldCheckIcon, label: t('nav.admins') },
  { to: '/admin/database', icon: CircleStackIcon,  label: t('nav.database') },
])
</script>
