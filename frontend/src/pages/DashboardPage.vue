<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-[#001d32]">Tableau de bord</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Terminal Overview — {{ todayFormatted }}</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="px-4 py-2 text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition">
          Exporter
        </button>
        <button @click="refresh" class="px-4 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90" style="background-color:#006d35;">
          <span v-if="loading" class="flex items-center gap-1.5">
            <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.4 0 0 5.4 0 12h4z"/></svg>
            Actualisation...
          </span>
          <span v-else>Actualiser</span>
        </button>
      </div>
    </div>

    <!-- Alert providers pending -->
    <div v-if="(stats?.providers_pending ?? 0) > 0"
      class="flex items-center justify-between rounded-xl p-4 border"
      style="background:#FFF7ED; border-color:#fed7aa;"
    >
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-full flex items-center justify-center" style="background:#fed7aa;">
          <svg class="w-4 h-4 text-orange-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
          </svg>
        </div>
        <div>
          <p class="text-sm font-semibold text-orange-800">{{ stats.providers_pending }} prestataire(s) en attente de validation</p>
          <p class="text-xs text-orange-600">Action requise — cliquez pour traiter</p>
        </div>
      </div>
      <RouterLink to="/admin/providers?status=pending"
        class="text-sm font-semibold px-4 py-2 rounded-lg text-white shrink-0"
        style="background-color:#c2410c;"
      >
        Traiter maintenant
      </RouterLink>
    </div>

    <!-- Stat cards -->
    <div v-if="loading" class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="n in 4" :key="n" class="bg-white rounded-2xl p-5 animate-pulse h-28 border border-[#f1f5f9]"></div>
    </div>
    <div v-else class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Total Utilisateurs -->
      <RouterLink to="/admin/users" class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm hover:shadow-md transition-shadow block">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color:#dcfce7;">
            <UsersIcon class="w-5 h-5" style="color:#006d35;" />
          </div>
          <span class="text-xs font-medium text-gray-400">Membres</span>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">{{ stats?.users_count ?? 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">Total Utilisateurs</p>
      </RouterLink>

      <!-- Prestataires -->
      <RouterLink to="/admin/providers?status=approved" class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm hover:shadow-md transition-shadow block">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color:#dbeafe;">
            <BriefcaseIcon class="w-5 h-5 text-blue-600" />
          </div>
          <span class="text-xs font-medium text-gray-400">Actifs</span>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">{{ stats?.providers_count ?? 0 }}</p>
        <p class="text-xs text-gray-400 mt-1">Prestataires</p>
      </RouterLink>

      <!-- Score Environnemental (mock) -->
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color:#f0fdf4;">
            <svg class="w-5 h-5" style="color:#006d35;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064"/>
            </svg>
          </div>
          <span class="text-xs font-semibold px-2 py-0.5 rounded-full" style="background:#dcfce7; color:#166534;">Impact Élevé</span>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">94<span class="text-sm font-normal text-gray-400">/100</span></p>
        <p class="text-xs text-gray-400 mt-1">Score Environnemental</p>
      </div>

      <!-- Objets Collectés (mock) -->
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color:#fef9c3;">
            <svg class="w-5 h-5 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
            </svg>
          </div>
          <span class="text-xs font-medium text-gray-400">Ce mois</span>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">1 247</p>
        <p class="text-xs text-gray-400 mt-1">Objets Collectés</p>
      </div>
    </div>

    <!-- Main content 2 columns -->
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">

      <!-- Left col (2/3): chart + activity -->
      <div class="xl:col-span-2 space-y-6">

        <!-- Monthly trend chart -->
        <div class="bg-white rounded-2xl p-6 border border-[#f1f5f9] shadow-sm">
          <div class="flex items-center justify-between mb-5">
            <div>
              <h3 class="font-semibold text-[#001d32]">Tendances Mensuelles</h3>
              <p class="text-xs text-gray-400 mt-0.5">Inscriptions par mois — 2025</p>
            </div>
            <select class="text-xs border border-[#e5e7eb] rounded-lg px-2 py-1 text-gray-500 bg-[#f8fafc]">
              <option>2025</option>
              <option>2024</option>
            </select>
          </div>
          <div class="flex items-end gap-2 h-32">
            <div v-for="(bar, i) in monthlyData" :key="i" class="flex-1 flex flex-col items-center gap-1">
              <div
                class="w-full rounded-t-md transition-all duration-500"
                :style="{
                  height: (bar.value / maxMonthly * 100) + '%',
                  backgroundColor: i === currentMonth ? '#006d35' : '#dcfce7',
                  minHeight: '4px'
                }"
              ></div>
              <span class="text-[9px] text-gray-400">{{ bar.label }}</span>
            </div>
          </div>
        </div>

        <!-- Activity feed -->
        <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          <div class="flex items-center justify-between px-6 py-4 border-b border-[#f1f5f9]">
            <h3 class="font-semibold text-[#001d32]">Flux d'Activité Récente</h3>
            <RouterLink to="/admin/logs" class="text-xs font-medium hover:underline" style="color:#006d35;">Voir tout</RouterLink>
          </div>

          <div v-if="logsLoading" class="p-4 space-y-3">
            <div v-for="n in 5" :key="n" class="h-10 bg-gray-50 rounded-lg animate-pulse"></div>
          </div>
          <div v-else-if="logs.length === 0" class="text-center py-12 text-gray-400 text-sm">
            Aucune activité récente
          </div>
          <div v-else class="divide-y divide-[#f8fafc]">
            <div
              v-for="log in logs.slice(0, 8)"
              :key="log.id"
              class="flex items-center gap-4 px-6 py-3.5 hover:bg-[#f8fafc] transition-colors cursor-pointer"
              @click="openModal(log)"
            >
              <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 text-xs font-bold text-white"
                :style="{ backgroundColor: avatarColor(log.admin_name) }">
                {{ (log.admin_name?.trim() || 'S').charAt(0).toUpperCase() }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-800">
                  <span class="font-medium">{{ log.admin_name?.trim() || 'Système' }}</span>
                  — {{ actionLabel(log.action) }}
                  <span v-if="log.resource_id" class="text-gray-400"> #{{ log.resource_id }}</span>
                </p>
                <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(log.created_at) }}</p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
                  {{ log.resource_type || '—' }}
                </span>
                <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="statusBadgeClass(log.action)">
                  {{ statusLabel(log.action) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right col (1/3): CTA cards -->
      <div class="space-y-4">
        <!-- Pending providers CTA -->
        <div class="rounded-2xl p-5 text-white" style="background: linear-gradient(135deg, #006d35, #1b8848);">
          <div class="flex items-start justify-between mb-4">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center bg-white/20">
              <BriefcaseIcon class="w-5 h-5 text-white" />
            </div>
            <span v-if="(stats?.providers_pending ?? 0) > 0" class="text-xs font-semibold px-2 py-0.5 rounded-full bg-white/20 text-white">
              {{ stats.providers_pending }} en attente
            </span>
          </div>
          <h3 class="font-bold text-base mb-1">Valider Prestataires</h3>
          <p class="text-xs text-white/70 mb-4">Examinez et approuvez les nouveaux prestataires inscrits sur la plateforme.</p>
          <RouterLink to="/admin/providers" class="inline-flex items-center gap-1.5 text-sm font-semibold text-white/90 hover:text-white transition">
            Voir les prestataires
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </RouterLink>
        </div>

        <!-- Categories CTA -->
        <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm" style="border-top: 3px solid #40617f;">
          <div class="flex items-start justify-between mb-4">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color:#f0f4f8;">
              <TagIcon class="w-5 h-5 text-[#40617f]" />
            </div>
          </div>
          <h3 class="font-bold text-base text-[#001d32] mb-1">Gestion Catégories</h3>
          <p class="text-xs text-gray-500 mb-4">Organisez et gérez les catégories d'objets et de prestations.</p>
          <RouterLink to="/admin/categories" class="inline-flex items-center gap-1.5 text-sm font-semibold transition" style="color:#40617f;">
            Gérer les catégories
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </RouterLink>
        </div>

        <!-- Prestations stats card -->
        <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
          <h3 class="font-semibold text-[#001d32] mb-4">Prestations par statut</h3>
          <div v-if="!stats || !Object.keys(stats.prestations_by_status || {}).length" class="text-sm text-gray-400">
            Aucune donnée
          </div>
          <div v-else class="space-y-3">
            <div v-for="(count, status) in stats.prestations_by_status" :key="status">
              <div class="flex justify-between items-center mb-1">
                <span class="text-xs font-medium capitalize text-gray-600">{{ status }}</span>
                <span class="text-xs font-bold text-gray-800">{{ count }}</span>
              </div>
              <div class="h-1.5 bg-gray-100 rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :style="{ width: barWidth(count) + '%', backgroundColor: statusColor(status) }"
                ></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Terminal Modal -->
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="selectedLog" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="closeModal">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
          <div class="relative w-full max-w-2xl rounded-xl overflow-hidden shadow-2xl border border-gray-700">
            <div class="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-gray-700">
              <button @click="closeModal" class="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors"></button>
              <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
              <span class="ml-3 text-xs font-mono text-gray-400">upcycleconnect — log#{{ selectedLog.id }}</span>
              <button @click="closeModal" class="ml-auto text-gray-500 hover:text-gray-300 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>
            <div class="bg-gray-900 p-5 font-mono text-sm leading-relaxed">
              <div class="flex items-center gap-2 mb-4 pb-3 border-b border-gray-700">
                <span class="text-green-400">root@upcycleconnect</span>
                <span class="text-gray-500">:</span>
                <span class="text-blue-400">~/logs</span>
                <span class="text-gray-500">$</span>
                <span class="text-white">cat log.{{ selectedLog.id }}.json</span>
              </div>
              <div class="space-y-1.5">
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">id</span><span class="text-yellow-300">{{ selectedLog.id }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">timestamp</span><span class="text-cyan-300">{{ formatDateFull(selectedLog.created_at) }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">ip_address</span><span class="text-green-300">{{ selectedLog.ip_address || 'unknown' }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">admin</span><span class="text-white">{{ selectedLog.admin_name?.trim() || 'system' }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">admin_id</span><span class="text-yellow-300">{{ selectedLog.admin_id ?? '—' }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">action</span><span :class="actionTerminalColor(selectedLog.action)" class="font-bold">{{ selectedLog.action }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">resource_type</span><span class="text-purple-300">{{ selectedLog.resource_type || '—' }}</span></div>
                <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">resource_id</span><span class="text-yellow-300">{{ selectedLog.resource_id ?? '—' }}</span></div>
                <div v-if="selectedLog.details" class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">details</span><span class="text-gray-300 break-all">{{ selectedLog.details }}</span></div>
              </div>
              <div class="mt-4 pt-3 border-t border-gray-700 flex items-center gap-1">
                <span class="text-green-400">root@upcycleconnect</span>
                <span class="text-gray-500">:</span>
                <span class="text-blue-400">~/logs</span>
                <span class="text-gray-500">$</span>
                <span class="inline-block w-2 h-4 bg-green-400 ml-1 animate-pulse"></span>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import { UsersIcon, BriefcaseIcon, TagIcon } from '@heroicons/vue/24/outline'

const auth = useAuthStore()
const stats = ref(null)
const logs  = ref([])
const loading     = ref(true)
const logsLoading = ref(true)
const selectedLog = ref(null)

function openModal(log) { selectedLog.value = log }
function closeModal() { selectedLog.value = null }
function onKeydown(e) { if (e.key === 'Escape') closeModal() }
onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

const todayFormatted = computed(() => {
  return new Date().toLocaleDateString('fr-FR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

const currentMonth = new Date().getMonth()

const monthlyData = [
  { label: 'Jan', value: 42 }, { label: 'Fév', value: 58 }, { label: 'Mar', value: 71 },
  { label: 'Avr', value: 65 }, { label: 'Mai', value: 89 }, { label: 'Jun', value: 103 },
  { label: 'Jul', value: 94 }, { label: 'Aoû', value: 78 }, { label: 'Sep', value: 112 },
  { label: 'Oct', value: 87 }, { label: 'Nov', value: 134 }, { label: 'Déc', value: 98 },
]
const maxMonthly = computed(() => Math.max(...monthlyData.map(d => d.value)))

async function loadData() {
  loading.value = true
  logsLoading.value = true
  try {
    const { data } = await api.get('/dashboard/stats')
    stats.value = data
  } finally {
    loading.value = false
  }
  try {
    const { data } = await api.get('/logs')
    logs.value = data.data || []
  } finally {
    logsLoading.value = false
  }
}

function refresh() { loadData() }

onMounted(loadData)

const totalPrestations = computed(() => {
  if (!stats.value) return 1
  return Object.values(stats.value.prestations_by_status || {}).reduce((a, b) => a + b, 0) || 1
})
function barWidth(count) {
  return Math.round((count / totalPrestations.value) * 100)
}
function statusColor(status) {
  return { published: '#006d35', draft: '#94a3b8', suspended: '#f97316', archived: '#e11d48' }[status] || '#94a3b8'
}

const avatarColors = ['#006d35', '#1b8848', '#40617f', '#001d32', '#7c3aed', '#0891b2', '#b45309']
function avatarColor(name) {
  if (!name) return '#94a3b8'
  const i = name.charCodeAt(0) % avatarColors.length
  return avatarColors[i]
}

function actionLabel(action) {
  const labels = {
    'user.banned':        'a banni un utilisateur',
    'user.activated':     'a activé un utilisateur',
    'user.deleted':       'a supprimé un utilisateur',
    'provider.approved':  'a approuvé un prestataire',
    'provider.suspended': 'a suspendu un prestataire',
    'admin.created':      'a créé un administrateur',
    'admin.deleted':      'a supprimé un administrateur',
    'category.created':   'a créé une catégorie',
    'category.deleted':   'a supprimé une catégorie',
    'category.toggled':   'a modifié une catégorie',
  }
  return labels[action] || action
}
function actionBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'bg-[#fee2e2] text-[#991b1b]'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'bg-[#dcfce7] text-[#166534]'
  if (action.includes('suspended')) return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function statusBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#dcfce7] text-[#166534]'
}
function statusLabel(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'Warning'
  return 'Success'
}
function formatDate(dateStr) {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now - d) / 1000)
  if (diff < 60) return "à l'instant"
  if (diff < 3600) return `il y a ${Math.floor(diff / 60)}min`
  if (diff < 86400) return `il y a ${Math.floor(diff / 3600)}h`
  return d.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
function formatDateFull(dateStr) {
  const d = new Date(dateStr)
  return d.toLocaleDateString('fr-FR', {
    weekday: 'short', day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  }) + ` (UTC${d.getTimezoneOffset() <= 0 ? '+' : '-'}${Math.abs(d.getTimezoneOffset() / 60)})`
}
function actionTerminalColor(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'text-red-400'
  if (action.includes('suspended')) return 'text-orange-400'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'text-green-400'
  return 'text-gray-300'
}
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
