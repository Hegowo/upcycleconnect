<template>
  <nav class="fixed bottom-0 left-0 right-0 z-40 lg:hidden bg-white border-t border-[#e5e7eb] shadow-[0_-2px_12px_rgba(0,0,0,0.1)] nav-safe">
    <div class="flex overflow-x-auto nav-scroll">
      <RouterLink
        v-for="item in visibleItems"
        :key="item.to"
        :to="item.to"
        class="flex flex-col items-center justify-center gap-1 px-4 py-2.5 min-w-[68px] shrink-0 transition-colors relative"
        :class="isActive(item.to) ? 'text-[#006d35]' : 'text-gray-400'"
      >
        <div
          v-if="isActive(item.to)"
          class="absolute top-0 left-1/2 -translate-x-1/2 w-8 h-[3px] rounded-full bg-[#006d35]"
        />
        <component :is="item.icon" class="w-5 h-5 shrink-0" :class="isActive(item.to) ? 'text-[#006d35]' : 'text-gray-400'" />
        <span class="text-[9px] font-semibold leading-tight text-center whitespace-nowrap">{{ item.label }}</span>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup>import { computed } from 'vue'
import { useRoute } from 'vue-router'
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
const route = useRoute()
const auth  = useAuthStore()

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

const superItems = computed(() => auth.isSuperAdmin ? [
  { to: '/admin/admins',   icon: ShieldCheckIcon, label: t('nav.admins') },
  { to: '/admin/database', icon: CircleStackIcon,  label: t('nav.database') },
] : [])

const visibleItems = computed(() => [...navItems.value, ...superItems.value])

function isActive(to) {
  return route.path === to || route.path.startsWith(to + '/')
}
</script>

<style scoped>
.nav-scroll {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
.nav-scroll::-webkit-scrollbar {
  display: none;
}
.nav-safe {
  padding-bottom: env(safe-area-inset-bottom, 0px);
}
</style>
