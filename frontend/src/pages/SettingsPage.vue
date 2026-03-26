<template>
  <div class="space-y-6 max-w-3xl">
    <!-- Header -->
    <div>
      <h2 class="text-2xl font-bold text-[#001d32]">Paramètres Système</h2>
      <p class="text-sm text-[#40617f] mt-0.5">Configuration Globale — Gérez les paramètres de la plateforme</p>
    </div>

    <!-- Global settings card -->
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-6 space-y-6">
      <h3 class="font-semibold text-[#001d32] text-base">Paramètres Globaux</h3>

      <!-- Toggle 1: Auto validation -->
      <div class="flex items-center justify-between py-3 border-b border-[#f8fafc]">
        <div>
          <p class="text-sm font-medium text-[#001d32]">Validation automatique des petits objets</p>
          <p class="text-xs text-gray-400 mt-0.5">Approuver automatiquement les dépôts de moins de 5kg</p>
        </div>
        <button
          @click="toggleAutoValidation"
          :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors', autoValidation ? 'bg-[#006d35]' : 'bg-gray-300']"
        >
          <span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform shadow', autoValidation ? 'translate-x-6' : 'translate-x-1']" />
        </button>
      </div>

      <!-- Languages -->
      <div class="py-3 border-b border-[#f8fafc]">
        <p class="text-sm font-medium text-[#001d32] mb-1">Gestion des langues</p>
        <p class="text-xs text-gray-400 mb-3">Langues actives sur la plateforme</p>
        <div class="flex items-center gap-2 flex-wrap">
          <span class="px-3 py-1.5 rounded-lg text-xs font-semibold border-2 flex items-center gap-1.5" style="border-color:#006d35; color:#006d35; background:#f0fdf4;">
            🇫🇷 FR
            <span class="text-[10px] text-green-600">actif</span>
          </span>
          <span class="px-3 py-1.5 rounded-lg text-xs font-semibold border-2 flex items-center gap-1.5" style="border-color:#006d35; color:#006d35; background:#f0fdf4;">
            🇬🇧 EN
            <span class="text-[10px] text-green-600">actif</span>
          </span>
          <span class="px-3 py-1.5 rounded-lg text-xs font-medium border border-gray-200 text-gray-400 flex items-center gap-1.5">
            🇩🇪 DE
          </span>
          <button class="px-3 py-1.5 rounded-lg text-xs font-medium border border-dashed border-[#006d35] text-[#006d35] hover:bg-[#f0fdf4] transition">
            + Ajouter
          </button>
        </div>
      </div>

      <!-- Toggle 2: Email notifications -->
      <div class="flex items-center justify-between py-3 border-b border-[#f8fafc]">
        <div>
          <p class="text-sm font-medium text-[#001d32]">Emails de notification activés</p>
          <p class="text-xs text-gray-400 mt-0.5">Envoyer des notifications par email aux administrateurs</p>
        </div>
        <button
          @click="emailNotifications = !emailNotifications"
          :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors', emailNotifications ? 'bg-[#006d35]' : 'bg-gray-300']"
        >
          <span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform shadow', emailNotifications ? 'translate-x-6' : 'translate-x-1']" />
        </button>
      </div>

      <!-- Language selector (functional) -->
      <div class="py-3">
        <p class="text-sm font-medium text-[#001d32] mb-3">Langue de l'interface admin</p>
        <div class="flex gap-3">
          <button
            @click="setLang('fr')"
            class="flex items-center gap-2 px-5 py-3 rounded-xl transition text-sm font-medium border-2"
            :style="currentLang === 'fr' ? 'border-color:#006d35; color:#006d35; background:#f0fdf4;' : 'border-color:#e5e7eb; color:#6b7280;'"
          >
            🇫🇷 Français
          </button>
          <button
            @click="setLang('en')"
            class="flex items-center gap-2 px-5 py-3 rounded-xl transition text-sm font-medium border-2"
            :style="currentLang === 'en' ? 'border-color:#006d35; color:#006d35; background:#f0fdf4;' : 'border-color:#e5e7eb; color:#6b7280;'"
          >
            🇬🇧 English
          </button>
        </div>
      </div>

      <!-- Profile info (read-only) -->
      <div class="py-3 border-t border-[#f8fafc]">
        <p class="text-sm font-medium text-[#001d32] mb-3">Informations du profil</p>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">Prénom</label>
            <p class="text-sm text-[#001d32] mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg">{{ auth.user?.first_name || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">Nom</label>
            <p class="text-sm text-[#001d32] mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg">{{ auth.user?.last_name || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">Email</label>
            <p class="text-sm text-[#001d32] mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg">{{ auth.user?.email || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">Rôle</label>
            <p class="text-sm capitalize mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg font-medium" style="color:#006d35;">{{ auth.user?.role?.replace('_', ' ') || '—' }}</p>
          </div>
        </div>
      </div>

      <!-- Action buttons -->
      <div class="flex items-center gap-3 pt-2">
        <button class="px-6 py-2.5 text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition">
          Annuler
        </button>
        <button @click="saveSettings" class="px-6 py-2.5 text-sm font-semibold rounded-lg text-white transition hover:opacity-90" style="background-color:#006d35;">
          Sauvegarder Configuration
        </button>
      </div>
    </div>

    <!-- API Access dark card -->
    <div class="rounded-2xl p-6" style="background: linear-gradient(135deg, #001d32, #0a2e4a);">
      <div class="flex items-start justify-between">
        <div>
          <div class="flex items-center gap-2 mb-2">
            <div class="w-8 h-8 rounded-lg flex items-center justify-center bg-white/10">
              <CodeBracketIcon class="w-4 h-4 text-white" />
            </div>
            <h3 class="font-bold text-white">Accès API Avancé</h3>
          </div>
          <p class="text-sm text-white/60 max-w-md">
            Accédez à la documentation complète de l'API REST, gérez vos tokens d'authentification et explorez les endpoints disponibles.
          </p>
          <div class="flex items-center gap-3 mt-4">
            <div class="flex items-center gap-1.5 text-xs text-white/40">
              <span class="w-1.5 h-1.5 rounded-full bg-green-400"></span>
              API v1 active
            </div>
            <div class="text-xs text-white/40">•</div>
            <div class="text-xs text-white/40">Go / Gin backend</div>
          </div>
        </div>
        <button class="px-5 py-2.5 rounded-lg text-sm font-semibold text-[#001d32] bg-white hover:bg-gray-100 transition shrink-0">
          Portail Développeur
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/utils/useToast'
import { CodeBracketIcon } from '@heroicons/vue/24/outline'

const { locale } = useI18n()
const auth = useAuthStore()
const toast = useToast()

const currentLang = computed(() => locale.value)
const autoValidation   = ref(false)
const emailNotifications = ref(true)

function toggleAutoValidation() {
  autoValidation.value = !autoValidation.value
}

function setLang(lang) {
  locale.value = lang
  localStorage.setItem('admin_lang', lang)
  toast.showSuccess('Langue mise à jour')
}

function saveSettings() {
  toast.showSuccess('Configuration sauvegardée')
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
