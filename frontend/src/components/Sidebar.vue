<template>
  <aside class="w-64 text-white flex flex-col shrink-0" style="background-color: #103652;">
    <div class="flex items-center justify-center px-3 py-4 border-b" style="border-color: #174d72;">
      <img
        src="/logoentier.png"
        alt="UpcycleConnect"
        class="w-full h-auto object-contain"
        style="max-height: 80px;"
        @error="onLogoError"
        ref="logoImg"
      />
    </div>

    <nav class="flex-1 px-3 py-4 space-y-1">
      <NavItem to="/admin/dashboard" icon="chart-bar">{{ t('nav.dashboard') }}</NavItem>
      <NavItem to="/admin/users" icon="users">{{ t('nav.users') }}</NavItem>
      <NavItem to="/admin/providers" icon="briefcase">{{ t('nav.providers') }}</NavItem>
      <NavItem to="/admin/categories" icon="tag">{{ t('nav.categories') }}</NavItem>
      <NavItem to="/admin/prestations" icon="clipboard-list">{{ t('nav.prestations') }}</NavItem>
      <NavItem to="/admin/events" icon="calendar">{{ t('nav.events') }}</NavItem>
      <NavItem to="/admin/logs" icon="clipboard-list">{{ t('nav.logs') }}</NavItem>
      <NavItem v-if="auth.isSuperAdmin" to="/admin/admins" icon="shield-check">{{ t('nav.admins') }}</NavItem>
      <NavItem v-if="auth.isSuperAdmin" to="/admin/database" icon="database">{{ t('nav.database') }}</NavItem>
      <NavItem to="/admin/settings" icon="cog">{{ t('nav.settings') }}</NavItem>
    </nav>

    <div class="px-4 py-3 text-sm" style="border-top: 1px solid #174d72;">
      <p class="text-xs" style="color: #8ab4c8;">{{ t('nav.connectedAs') }}</p>
      <p class="text-white font-medium truncate">{{ auth.fullName }}</p>
      <p class="text-xs capitalize" style="color: #1B8848;">{{ auth.role?.replace('_', ' ') }}</p>
    </div>
  </aside>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import NavItem from './NavItem.vue'

const { t } = useI18n()
const auth = useAuthStore()
const logoImg = ref(null)

function onLogoError() {
  if (logoImg.value) logoImg.value.src = '/logoentier.svg'
}
</script>
