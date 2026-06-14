<template>
  <aside :class="[
    'hidden lg:flex w-64 flex-col shrink-0 h-screen border-r transition-colors duration-200',
    theme.isDark
      ? 'bg-[#1e293b] border-[#334155]'
      : 'bg-white border-[#e5e7eb]'
  ]">

    <div :class="['flex items-center gap-3 px-4 py-4 border-b', theme.isDark ? 'border-[#334155]' : 'border-[#e5e7eb]']">
      <img src="/logoentier.png" alt="UpcycleConnect" class="h-9 w-auto object-contain" />
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
            'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 text-left border-l-2',
            isActive
              ? (theme.isDark
                  ? 'bg-[#006d35]/20 text-[#4ade80] border-[#4ade80]'
                  : 'bg-[#f0fdf4] text-[#006d35] border-[#006d35]')
              : (theme.isDark
                  ? 'text-[#94a3b8] hover:bg-[#263244] border-transparent hover:text-[#f1f5f9]'
                  : 'text-[#374151] hover:bg-[#f9fafb] border-transparent')
          ]"
        >
          <component
            :is="item.icon"
            class="w-5 h-5 shrink-0"
            :class="isActive
              ? (theme.isDark ? 'text-[#4ade80]' : 'text-[#006d35]')
              : (theme.isDark ? 'text-[#64748b]' : 'text-[#9ca3af]')"
          />
          <span>{{ item.label }}</span>
        </button>
      </RouterLink>

      <div v-if="auth.isSuperAdmin" class="pt-3 pb-1">
        <p :class="['px-3 text-[10px] font-semibold uppercase tracking-widest mb-1', theme.isDark ? 'text-[#475569]' : 'text-gray-400']">
          Super Admin
        </p>
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
              'w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-all duration-150 text-left border-l-2',
              isActive
                ? (theme.isDark
                    ? 'bg-[#006d35]/20 text-[#4ade80] border-[#4ade80]'
                    : 'bg-[#f0fdf4] text-[#006d35] border-[#006d35]')
                : (theme.isDark
                    ? 'text-[#94a3b8] hover:bg-[#263244] border-transparent hover:text-[#f1f5f9]'
                    : 'text-[#374151] hover:bg-[#f9fafb] border-transparent')
            ]"
          >
            <component
              :is="item.icon"
              class="w-5 h-5 shrink-0"
              :class="isActive
                ? (theme.isDark ? 'text-[#4ade80]' : 'text-[#006d35]')
                : (theme.isDark ? 'text-[#64748b]' : 'text-[#9ca3af]')"
            />
            <span>{{ item.label }}</span>
          </button>
        </RouterLink>
      </template>
    </nav>

    <div :class="['px-4 py-4 border-t', theme.isDark ? 'border-[#334155]' : 'border-[#e5e7eb]']">
      <div class="flex items-center gap-3">
        <div class="w-9 h-9 rounded-full flex items-center justify-center text-white text-sm font-bold shrink-0" style="background-color:#001d32;">
          {{ initials }}
        </div>
        <div class="flex-1 min-w-0">
          <p :class="['text-sm font-semibold truncate leading-none', theme.isDark ? 'text-[#f1f5f9]' : 'text-gray-900']">
            {{ auth.fullName }}
          </p>
          <p class="text-xs capitalize mt-0.5 font-medium" style="color:#006d35;">{{ auth.role?.replace('_', ' ') }}</p>
        </div>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useThemeStore } from '@/stores/theme'
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
  MapPinIcon,
  ChatBubbleLeftEllipsisIcon,
  BanknotesIcon,
  LanguageIcon,
  BellAlertIcon,
  ArrowDownTrayIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const auth  = useAuthStore()
const theme = useThemeStore()

const initials = computed(() => {
  if (!auth.user) return '?'
  return ((auth.user.first_name?.[0] || '') + (auth.user.last_name?.[0] || '')).toUpperCase()
})

const navItems = computed(() => [
  { to: '/admin/dashboard',         icon: Squares2X2Icon,             label: t('nav.dashboard') },
  { to: '/admin/users',             icon: UsersIcon,                  label: t('nav.users') },
  { to: '/admin/providers',         icon: BriefcaseIcon,              label: t('nav.providers') },
  { to: '/admin/categories',        icon: TagIcon,                    label: t('nav.categories') },
  { to: '/admin/prestations',       icon: ClipboardDocumentListIcon,  label: t('nav.prestations') },
  { to: '/admin/events',            icon: CalendarIcon,               label: t('nav.events') },
  { to: '/admin/forum',             icon: ChatBubbleLeftEllipsisIcon, label: t('nav.forum') },
  { to: '/admin/monetization',      icon: BanknotesIcon,              label: 'Monétisation' },
  { to: '/admin/locales',           icon: LanguageIcon,               label: 'Langues' },
  { to: '/admin/notifications',     icon: BellAlertIcon,              label: 'Notifications' },
  { to: '/admin/export',            icon: ArrowDownTrayIcon,          label: 'Export & Audit' },
  { to: '/admin/logs',              icon: ClockIcon,                  label: t('nav.logs') },
  { to: '/admin/box-requests',      icon: InboxIcon,                  label: t('nav.boxRequests') },
  { to: '/admin/collection-points', icon: MapPinIcon,                 label: t('nav.collectionPoints') },
  { to: '/admin/settings',          icon: Cog6ToothIcon,              label: t('nav.settings') },
  { to: '/admin/docs',              icon: CodeBracketIcon,            label: t('nav.apiDocs') },
])

const superAdminItems = computed(() => [
  { to: '/admin/admins',   icon: ShieldCheckIcon, label: t('nav.admins') },
  { to: '/admin/database', icon: CircleStackIcon,  label: t('nav.database') },
])
</script>
