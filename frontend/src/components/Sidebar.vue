<template>
  <aside class="w-64 text-white flex flex-col shrink-0" style="background-color: #103652;">
    <!-- Logo -->
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

    <!-- Navigation -->
    <nav class="flex-1 px-3 py-4 space-y-1">
      <NavItem to="/admin/dashboard" icon="chart-bar">Tableau de bord</NavItem>
      <NavItem to="/admin/users" icon="users">Utilisateurs</NavItem>
      <NavItem to="/admin/providers" icon="briefcase">Prestataires</NavItem>
      <NavItem to="/admin/categories" icon="tag">Catégories</NavItem>
      <NavItem to="/admin/prestations" icon="clipboard-list">Prestations</NavItem>
      <NavItem to="/admin/events" icon="calendar">Événements</NavItem>
      <NavItem v-if="auth.isSuperAdmin" to="/admin/admins" icon="shield-check">Administrateurs</NavItem>
      <NavItem v-if="auth.isSuperAdmin" to="/admin/database" icon="database">Base de données</NavItem>
    </nav>

    <!-- User info -->
    <div class="px-4 py-3 text-sm" style="border-top: 1px solid #174d72;">
      <p class="text-xs" style="color: #8ab4c8;">Connecté en tant que</p>
      <p class="text-white font-medium truncate">{{ auth.fullName }}</p>
      <p class="text-xs capitalize" style="color: #1B8848;">{{ auth.role?.replace('_', ' ') }}</p>
    </div>
  </aside>
</template>

<script setup>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import NavItem from './NavItem.vue'

const auth = useAuthStore()
const logoImg = ref(null)

function onLogoError() {
  if (logoImg.value) logoImg.value.src = '/logoentier.svg'
}
</script>
