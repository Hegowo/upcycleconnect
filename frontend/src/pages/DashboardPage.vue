<template>
  <div class="space-y-4 sm:space-y-6">

    
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ t('dashboard.title') }}</h2>
        <p class="text-xs sm:text-sm text-[#40617f] mt-0.5">{{ todayFormatted }}</p>
      </div>
      <div class="flex items-center gap-2 shrink-0">
        
        <button @click="exportDashboard" class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition">
          <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
          </svg>
          <span class="hidden sm:inline">{{ t('common.exportCsv') }}</span>
        </button>
        
        <button @click="refresh" class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg text-white transition hover:opacity-90" style="background-color:#006d35;">
          <span v-if="loading">
            <svg class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.4 0 0 5.4 0 12h4z"/></svg>
          </span>
          <span v-else>
            <svg class="w-4 h-4 sm:hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            <span class="hidden sm:inline">{{ t('common.refresh') }}</span>
          </span>
        </button>
      </div>
    </div>

    
    <div v-if="loading" class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
      <div v-for="n in 4" :key="n" class="bg-white rounded-2xl p-4 sm:p-5 animate-pulse h-24 sm:h-28 border border-[#f1f5f9]"></div>
    </div>
    <div v-else class="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">

      <RouterLink to="/admin/users" class="bg-white rounded-2xl p-4 sm:p-5 border border-[#f1f5f9] shadow-sm hover:shadow-md transition-shadow block">
        <div class="flex items-center justify-between mb-2 sm:mb-3">
          <div class="w-8 h-8 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center" style="background-color:#dcfce7;">
            <UsersIcon class="w-4 h-4 sm:w-5 sm:h-5" style="color:#006d35;" />
          </div>
          <span class="text-[10px] sm:text-xs font-medium text-gray-400 text-right leading-tight">{{ t('dashboard.members') }}</span>
        </div>
        <p class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ stats?.users_count ?? 0 }}</p>
        <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5 sm:mt-1 leading-tight">{{ t('dashboard.statsUsers') }}</p>
      </RouterLink>

      <RouterLink to="/admin/providers?status=approved" class="bg-white rounded-2xl p-4 sm:p-5 border border-[#f1f5f9] shadow-sm hover:shadow-md transition-shadow block">
        <div class="flex items-center justify-between mb-2 sm:mb-3">
          <div class="w-8 h-8 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center" style="background-color:#dbeafe;">
            <BriefcaseIcon class="w-4 h-4 sm:w-5 sm:h-5 text-blue-600" />
          </div>
          <span class="text-[10px] sm:text-xs font-medium text-gray-400 text-right leading-tight">{{ t('dashboard.statsActive') }}</span>
        </div>
        <p class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ stats?.providers_count ?? 0 }}</p>
        <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5 sm:mt-1 leading-tight">{{ t('dashboard.statsProviders') }}</p>
      </RouterLink>

      <div class="bg-white rounded-2xl p-4 sm:p-5 border border-[#f1f5f9] shadow-sm relative">
        <div class="flex items-center justify-between mb-2 sm:mb-3">
          <div class="w-8 h-8 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center" style="background-color:#f0fdf4;">
            <svg class="w-4 h-4 sm:w-5 sm:h-5" style="color:#006d35;" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064"/>
            </svg>
          </div>
          <div class="flex items-center gap-1">
            <span class="hidden sm:inline text-xs font-semibold px-2 py-0.5 rounded-full" style="background:#dcfce7; color:#166534;">{{ t('dashboard.envScoreImpact') }}</span>
            <div class="relative group">
              <InformationCircleIcon class="w-4 h-4 text-gray-400 cursor-help" />
              <div class="absolute right-0 top-5 w-52 bg-[#001d32] text-white text-xs rounded-xl p-3 z-10 hidden group-hover:block shadow-lg leading-relaxed">
                {{ t('dashboard.envScoreFormula') }}
              </div>
            </div>
          </div>
        </div>
        <p class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ envScore }}<span class="text-xs sm:text-sm font-normal text-gray-400">/100</span></p>
        <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5 sm:mt-1 leading-tight">{{ t('dashboard.envScore') }}</p>
      </div>

      <div class="bg-white rounded-2xl p-4 sm:p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-2 sm:mb-3">
          <div class="w-8 h-8 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center" style="background-color:#fef9c3;">
            <svg class="w-4 h-4 sm:w-5 sm:h-5 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"/>
            </svg>
          </div>
          <span class="text-[10px] sm:text-xs font-medium text-gray-400 text-right leading-tight">{{ t('dashboard.totalLabel') }}</span>
        </div>
        <p class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ stats?.deposits_accepted ?? 0 }}</p>
        <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5 sm:mt-1 leading-tight">{{ t('dashboard.objectsCollected') }}</p>
      </div>

    </div>

    
    <div class="grid grid-cols-2 gap-3 xl:hidden">
      <RouterLink to="/admin/providers"
        class="flex items-center gap-3 bg-gradient-to-r from-[#006d35] to-[#1b8848] rounded-2xl p-4 text-white">
        <div class="w-9 h-9 rounded-xl flex items-center justify-center bg-white/20 shrink-0">
          <BriefcaseIcon class="w-5 h-5 text-white" />
        </div>
        <div class="min-w-0">
          <p class="text-xs font-bold leading-tight">{{ t('dashboard.validateProviders') }}</p>
          <p v-if="(stats?.providers_pending ?? 0) > 0" class="text-[10px] text-white/70 mt-0.5">
            {{ stats.providers_pending }} {{ t('dashboard.statsPending').toLowerCase() }}
          </p>
        </div>
      </RouterLink>
      <RouterLink to="/admin/categories"
        class="flex items-center gap-3 bg-white rounded-2xl p-4 border border-[#f1f5f9] shadow-sm" style="border-left: 3px solid #40617f;">
        <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0" style="background-color:#f0f4f8;">
          <TagIcon class="w-5 h-5 text-[#40617f]" />
        </div>
        <div class="min-w-0">
          <p class="text-xs font-bold text-[#001d32] leading-tight">{{ t('dashboard.manageCategories') }}</p>
          <p class="text-[10px] text-gray-400 mt-0.5">{{ t('dashboard.manageCategoriesDesc') }}</p>
        </div>
      </RouterLink>
    </div>

    
    <div class="grid grid-cols-1 xl:grid-cols-3 gap-4 sm:gap-6">

      
      <div class="xl:col-span-2 space-y-4 sm:space-y-6">

        
        <div class="bg-white rounded-2xl p-4 sm:p-6 border border-[#f1f5f9] shadow-sm">
          <div class="flex items-center justify-between mb-4 sm:mb-5">
            <div>
              <h3 class="font-semibold text-[#001d32] text-sm sm:text-base">{{ t('dashboard.trends') }}</h3>
              <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5">{{ t('dashboard.trendsSubtitle') }} — 2026</p>
            </div>
            <span class="text-xs border border-[#e5e7eb] rounded-lg px-2 py-1 text-gray-500 bg-[#f8fafc]">2026</span>
          </div>
          <div class="flex items-end gap-1 sm:gap-2 h-24 sm:h-32">
            <div v-for="(bar, i) in monthlyData" :key="i" class="flex-1 flex flex-col items-center gap-0.5 sm:gap-1">
              <div
                class="w-full rounded-t-md transition-all duration-500"
                :style="{
                  height: (bar.value / maxMonthly * 100) + '%',
                  backgroundColor: i === currentMonth ? '#006d35' : '#dcfce7',
                  minHeight: '4px'
                }"
              ></div>
              <span class="text-[8px] sm:text-[9px] text-gray-400">{{ bar.label }}</span>
            </div>
          </div>
        </div>

        
        <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          <div class="flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4 border-b border-[#f1f5f9]">
            <h3 class="font-semibold text-[#001d32] text-sm sm:text-base">{{ t('dashboard.recentActivity') }}</h3>
            <RouterLink to="/admin/logs" class="text-xs font-medium hover:underline" style="color:#006d35;">{{ t('common.seeAll') }}</RouterLink>
          </div>

          <div v-if="logsLoading" class="p-4 space-y-3">
            <div v-for="n in 5" :key="n" class="h-10 bg-gray-50 rounded-lg animate-pulse"></div>
          </div>
          <div v-else-if="logs.length === 0" class="text-center py-12 text-gray-400 text-sm">
            {{ t('logs.noActivity') }}
          </div>
          <div v-else class="divide-y divide-[#f8fafc]">
            <div
              v-for="log in logs.slice(0, 8)"
              :key="log.id"
              class="flex items-start gap-3 px-4 sm:px-6 py-3 sm:py-3.5 hover:bg-[#f8fafc] transition-colors cursor-pointer"
              @click="openModal(log)"
            >
              
              <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full flex items-center justify-center shrink-0 text-xs font-bold text-white mt-0.5"
                :style="{ backgroundColor: avatarColor(log.admin_name) }">
                {{ (log.admin_name?.trim() || 'S').charAt(0).toUpperCase() }}
              </div>

              
              <div class="flex-1 min-w-0">
                
                <div class="hidden sm:flex items-center gap-3">
                  <p class="text-sm text-gray-800 flex-1 min-w-0 truncate">
                    <span class="font-medium">{{ log.admin_name?.trim() || 'Système' }}</span>
                    — {{ actionLabel(log.action) }}
                    <span v-if="log.resource_id" class="text-gray-400"> #{{ log.resource_id }}</span>
                  </p>
                  <div class="flex items-center gap-2 shrink-0">
                    <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
                      {{ log.resource_type || '—' }}
                    </span>
                    <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="statusBadgeClass(log.action)">
                      {{ statusLabel(log.action) }}
                    </span>
                  </div>
                </div>
                
                <div class="sm:hidden">
                  <p class="text-xs text-gray-800 leading-snug">
                    <span class="font-semibold">{{ log.admin_name?.trim() || 'Système' }}</span>
                    — {{ actionLabel(log.action) }}
                  </p>
                  <div class="flex items-center gap-1.5 mt-1 flex-wrap">
                    <span class="text-[10px] text-gray-400">{{ formatDate(log.created_at) }}</span>
                    <span class="text-[10px] px-1.5 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
                      {{ log.resource_type || '—' }}
                    </span>
                    <span class="text-[10px] px-1.5 py-0.5 rounded-full font-medium" :class="statusBadgeClass(log.action)">
                      {{ statusLabel(log.action) }}
                    </span>
                  </div>
                </div>
                
                <p class="hidden sm:block text-xs text-gray-400 mt-0.5">{{ formatDate(log.created_at) }}</p>
              </div>
            </div>
          </div>
        </div>

      </div>

      
      <div class="hidden xl:block space-y-4">
        <div class="rounded-2xl p-5 text-white" style="background: linear-gradient(135deg, #006d35, #1b8848);">
          <div class="flex items-start justify-between mb-4">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center bg-white/20">
              <BriefcaseIcon class="w-5 h-5 text-white" />
            </div>
            <span v-if="(stats?.providers_pending ?? 0) > 0" class="text-xs font-semibold px-2 py-0.5 rounded-full bg-white/20 text-white">
              {{ stats.providers_pending }} {{ t('dashboard.statsPending').toLowerCase() }}
            </span>
          </div>
          <h3 class="font-bold text-base mb-1">{{ t('dashboard.validateProviders') }}</h3>
          <p class="text-xs text-white/70 mb-4">{{ t('dashboard.validateProvidersDesc') }}</p>
          <RouterLink to="/admin/providers" class="inline-flex items-center gap-1.5 text-sm font-semibold text-white/90 hover:text-white transition">
            {{ t('dashboard.seeProviders') }}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </RouterLink>
        </div>

        <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm" style="border-top: 3px solid #40617f;">
          <div class="flex items-start justify-between mb-4">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background-color:#f0f4f8;">
              <TagIcon class="w-5 h-5 text-[#40617f]" />
            </div>
          </div>
          <h3 class="font-bold text-base text-[#001d32] mb-1">{{ t('dashboard.manageCategories') }}</h3>
          <p class="text-xs text-gray-500 mb-4">{{ t('dashboard.manageCategoriesDesc') }}</p>
          <RouterLink to="/admin/categories" class="inline-flex items-center gap-1.5 text-sm font-semibold transition" style="color:#40617f;">
            {{ t('dashboard.manageCategoriasLink') }}
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
          </RouterLink>
        </div>
      </div>

    </div>

    
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="selectedLog" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center sm:p-4" @click.self="closeModal">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
          
          <div class="relative w-full sm:max-w-2xl rounded-t-2xl sm:rounded-xl overflow-hidden shadow-2xl border border-gray-700">
            <div class="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-gray-700">
              <button @click="closeModal" class="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors"></button>
              <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
              <span class="ml-3 text-xs font-mono text-gray-400 truncate">upcycleconnect — log#{{ selectedLog.id }}</span>
              <button @click="closeModal" class="ml-auto text-gray-500 hover:text-gray-300 transition-colors shrink-0">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>
            <div class="bg-gray-900 p-4 sm:p-5 font-mono text-xs sm:text-sm leading-relaxed max-h-[70vh] overflow-y-auto">
              <div class="flex items-center gap-2 mb-4 pb-3 border-b border-gray-700 flex-wrap">
                <span class="text-green-400">root@upcycleconnect</span>
                <span class="text-gray-500">:</span>
                <span class="text-blue-400">~/logs</span>
                <span class="text-gray-500">$</span>
                <span class="text-white">cat log.{{ selectedLog.id }}.json</span>
              </div>
              <div class="space-y-1.5">
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">id</span><span class="text-yellow-300">{{ selectedLog.id }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">timestamp</span><span class="text-cyan-300 break-all">{{ formatDateFull(selectedLog.created_at) }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">ip_address</span><span class="text-green-300">{{ selectedLog.ip_address || 'unknown' }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">admin</span><span class="text-white">{{ selectedLog.admin_name?.trim() || 'system' }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">admin_id</span><span class="text-yellow-300">{{ selectedLog.admin_id ?? '—' }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">action</span><span :class="actionTerminalColor(selectedLog.action)" class="font-bold break-all">{{ selectedLog.action }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">resource_type</span><span class="text-purple-300">{{ selectedLog.resource_type || '—' }}</span></div>
                <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">resource_id</span><span class="text-yellow-300">{{ selectedLog.resource_id ?? '—' }}</span></div>
                <template v-if="selectedLog.old_values">
                  <div class="flex gap-2 sm:gap-3 mt-2 pt-2 border-t border-gray-700">
                    <span class="text-gray-500 w-24 sm:w-28 shrink-0">old_values</span>
                    <span class="text-red-300 break-all text-xs font-mono whitespace-pre">{{ formatJson(selectedLog.old_values) }}</span>
                  </div>
                </template>
                <template v-if="selectedLog.new_values">
                  <div class="flex gap-2 sm:gap-3">
                    <span class="text-gray-500 w-24 sm:w-28 shrink-0">new_values</span>
                    <span class="text-green-300 break-all text-xs font-mono whitespace-pre">{{ formatJson(selectedLog.new_values) }}</span>
                  </div>
                </template>
                <div v-if="selectedLog.details" class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">details</span><span class="text-gray-300 break-all">{{ selectedLog.details }}</span></div>
              </div>
              <div class="mt-4 pt-3 border-t border-gray-700 flex items-center gap-1 flex-wrap">
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

<script setup>import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import { UsersIcon, BriefcaseIcon, TagIcon, InformationCircleIcon } from '@heroicons/vue/24/outline'

const { t } = useI18n()

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

const envScore = computed(() => {
  if (!stats.value) return 0
  const accepted = stats.value.deposits_accepted ?? 0
  const co2 = stats.value.co2_total ?? 0
  return Math.min(100, Math.round(accepted * 5 + co2 * 1.5))
})

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

function exportDashboard() {
  if (!stats.value) return
  const rows = [
    ['Indicateur', 'Valeur'],
    ['Total utilisateurs', stats.value.users_count ?? 0],
    ['Prestataires approuvés', stats.value.providers_count ?? 0],
    ['Prestataires en attente', stats.value.providers_pending ?? 0],
    ['Total prestations', stats.value.prestations_count ?? 0],
    ['Total événements', stats.value.events_count ?? 0],
    ...Object.entries(stats.value.prestations_by_status || {}).map(([s, c]) => [`Prestations ${s}`, c]),
  ]
  downloadCsv(rows, `dashboard_${new Date().toISOString().slice(0,10)}.csv`)
}

function downloadCsv(rows, filename) {
  const csv = rows.map(r => r.map(v => `"${String(v).replace(/"/g, '""')}"`).join(',')).join('\n')
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = filename
  a.click()
  URL.revokeObjectURL(a.href)
}

onMounted(loadData)

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
function formatJson(raw) {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
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
