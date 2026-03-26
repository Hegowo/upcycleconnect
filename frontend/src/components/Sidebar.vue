<template>
  <aside class="w-64 bg-white border-r border-[#e5e7eb] flex flex-col shrink-0 h-screen sticky top-0">
    <!-- Logo -->
    <div class="flex items-center gap-3 px-5 py-5 border-b border-[#e5e7eb]">
      <div class="w-9 h-9 rounded-lg flex items-center justify-center shrink-0" style="background-color:#006d35;">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M12 2C10.5 2 9.2 2.6 8.2 3.6L5.4 6.4C4.5 7.3 4 8.6 4 10c0 2.8 2.2 5 5 5h.5l-2 3.5H9l2-3.5h2l2 3.5h1.5l-2-3.5H15c2.8 0 5-2.2 5-5 0-1.4-.5-2.7-1.4-3.6L15.8 3.6C14.8 2.6 13.5 2 12 2z" fill="white" opacity="0.9"/>
          <path d="M7 16.5C7 18.4 8.6 20 10.5 20h3c1.9 0 3.5-1.6 3.5-3.5" stroke="white" stroke-width="1.5" stroke-linecap="round" fill="none"/>
        </svg>
      </div>
      <div>
        <p class="text-sm font-bold leading-none" style="color:#006d35;">UpcycleConnect</p>
        <p class="text-[10px] font-medium tracking-widest uppercase mt-0.5 text-gray-400">Admin Terminal</p>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-3 py-4 space-y-0.5 overflow-y-auto">
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        custom
        v-slot="{ isActive, navigate }"
      >
        <button
          @click="navigate"
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
            @click="navigate"
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

    <!-- Nouveau Rapport button -->
    <div class="px-3 pb-3">
      <button class="w-full flex items-center justify-center gap-2 px-4 py-2.5 rounded-lg text-white text-sm font-semibold transition hover:opacity-90" style="background-color:#006d35;">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Nouveau Rapport
      </button>
    </div>

    <!-- User footer -->
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

<script setup>
import { computed } from 'vue'
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
} from '@heroicons/vue/24/outline'

const auth = useAuthStore()

const initials = computed(() => {
  if (!auth.user) return '?'
  return ((auth.user.first_name?.[0] || '') + (auth.user.last_name?.[0] || '')).toUpperCase()
})

const navItems = [
  { to: '/admin/dashboard',    icon: Squares2X2Icon,           label: 'Tableau de bord' },
  { to: '/admin/users',        icon: UsersIcon,                label: 'Utilisateurs' },
  { to: '/admin/providers',    icon: BriefcaseIcon,            label: 'Prestataires' },
  { to: '/admin/categories',   icon: TagIcon,                  label: 'Catégories' },
  { to: '/admin/prestations',  icon: ClipboardDocumentListIcon, label: 'Prestations' },
  { to: '/admin/events',       icon: CalendarIcon,             label: 'Événements' },
  { to: '/admin/logs',         icon: ClockIcon,                label: 'Journaux' },
  { to: '/admin/box-requests', icon: InboxIcon,                label: 'Demandes Dépôt' },
  { to: '/admin/settings',     icon: Cog6ToothIcon,            label: 'Paramètres' },
]

const superAdminItems = [
  { to: '/admin/admins',   icon: ShieldCheckIcon, label: 'Administrateurs' },
  { to: '/admin/database', icon: CircleStackIcon,  label: 'Base de données' },
]
</script>
