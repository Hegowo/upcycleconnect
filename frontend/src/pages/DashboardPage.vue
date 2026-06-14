<template>
  <div class="space-y-4 sm:space-y-6">

    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ t('dashboard.title') }}</h2>
        <p class="text-xs sm:text-sm text-[#40617f] mt-0.5">{{ todayFormatted }}</p>
      </div>
      <div class="flex items-center gap-2 shrink-0">

        <button @click="customizing = true" class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition" :title="t('dashboard.customize')">
          <AdjustmentsHorizontalIcon class="w-4 h-4 shrink-0" />
          <span class="hidden sm:inline">{{ t('dashboard.customize') }}</span>
        </button>

        <button @click="exportDashboard" class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition">
          <ArrowDownTrayIcon class="w-4 h-4 shrink-0" />
          <span class="hidden sm:inline">{{ t('common.exportCsv') }}</span>
        </button>

        <button @click="refresh" class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg text-white transition hover:opacity-90" style="background-color:#006d35;">
          <ArrowPathIcon class="w-4 h-4" :class="loading ? 'animate-spin' : ''" />
          <span class="hidden sm:inline">{{ t('common.refresh') }}</span>
        </button>
      </div>
    </div>

    <div v-if="visibleIds.length === 0" class="text-center py-16 text-gray-400 text-sm border border-dashed border-[#e2e8f0] rounded-2xl">
      {{ t('dashboard.emptyLayout') }}
      <button @click="customizing = true" class="block mx-auto mt-2 text-[#006d35] font-semibold hover:underline">{{ t('dashboard.customize') }}</button>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-3 sm:gap-4 items-start">
      <template v-for="id in visibleIds" :key="id">
        <div :class="spanOf(id)">

          <StatWidget v-if="id === 'users'" :icon="UsersIcon" iconBg="#dcfce7" iconColor="#006d35"
            :value="stats?.users_count ?? 0" :label="t('dashboard.statsUsers')" :cornerLabel="t('dashboard.members')" to="/admin/users" />

          <StatWidget v-else-if="id === 'providers'" :icon="BriefcaseIcon" iconBg="#dbeafe" iconColor="#2563eb"
            :value="stats?.providers_count ?? 0" :label="t('dashboard.statsProviders')" :cornerLabel="t('dashboard.statsActive')" to="/admin/providers?status=approved" />

          <StatWidget v-else-if="id === 'env'" :icon="GlobeAltIcon" iconBg="#f0fdf4" iconColor="#006d35"
            :value="envScore" suffix="/100" :label="t('dashboard.envScore')" :cornerLabel="t('dashboard.envScoreImpact')" :tooltip="t('dashboard.envScoreFormula')" />

          <StatWidget v-else-if="id === 'deposits'" :icon="CubeIcon" iconBg="#fef9c3" iconColor="#ca8a04"
            :value="stats?.deposits_accepted ?? 0" :label="t('dashboard.objectsCollected')" :cornerLabel="t('dashboard.totalLabel')" />

          <StatWidget v-else-if="id === 'prestations'" :icon="TagIcon" iconBg="#ede9fe" iconColor="#7c3aed"
            :value="stats?.prestations_count ?? 0" :label="t('dashboard.statsPrestations')" to="/admin/prestations" />

          <StatWidget v-else-if="id === 'events'" :icon="CalendarDaysIcon" iconBg="#e0f2fe" iconColor="#0891b2"
            :value="stats?.events_count ?? 0" :label="t('dashboard.statsEvents')" />

          <TrendsWidget v-else-if="id === 'trends'" />

          <ActivityWidget v-else-if="id === 'activity'" :logs="logs" :loading="logsLoading" />

          <PendingProvidersWidget v-else-if="id === 'pending'" :pending="stats?.providers_pending ?? 0" />

        </div>
      </template>
    </div>

    <DashboardCustomizer
      v-model:open="customizing"
      :items="customizerItems"
      @update:items="onUpdateItems"
      @reset="resetLayout"
    />

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import {
  UsersIcon, BriefcaseIcon, GlobeAltIcon, CubeIcon, TagIcon, CalendarDaysIcon,
  AdjustmentsHorizontalIcon, ArrowDownTrayIcon, ArrowPathIcon,
} from '@heroicons/vue/24/outline'
import StatWidget from '@/components/dashboard/StatWidget.vue'
import TrendsWidget from '@/components/dashboard/TrendsWidget.vue'
import ActivityWidget from '@/components/dashboard/ActivityWidget.vue'
import PendingProvidersWidget from '@/components/dashboard/PendingProvidersWidget.vue'
import DashboardCustomizer from '@/components/dashboard/DashboardCustomizer.vue'

const { t } = useI18n()
const auth = useAuthStore()

const stats = ref(null)
const logs  = ref([])
const loading     = ref(true)
const logsLoading = ref(true)
const customizing = ref(false)

const todayFormatted = computed(() =>
  new Date().toLocaleDateString('fr-FR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }))

const envScore = computed(() => {
  if (!stats.value) return 0
  const accepted = stats.value.deposits_accepted ?? 0
  const co2 = stats.value.co2_total ?? 0
  return Math.min(100, Math.round(accepted * 5 + co2 * 1.5))
})

// ---- Widget registry: order here = default order ----
const WIDGET_META = [
  { id: 'users',       defaultVisible: true,  span: 'xl:col-span-1' },
  { id: 'providers',   defaultVisible: true,  span: 'xl:col-span-1' },
  { id: 'env',         defaultVisible: true,  span: 'xl:col-span-1' },
  { id: 'deposits',    defaultVisible: true,  span: 'xl:col-span-1' },
  { id: 'prestations', defaultVisible: false, span: 'xl:col-span-1' },
  { id: 'events',      defaultVisible: false, span: 'xl:col-span-1' },
  { id: 'trends',      defaultVisible: true,  span: 'sm:col-span-2 xl:col-span-2' },
  { id: 'activity',    defaultVisible: true,  span: 'sm:col-span-2 xl:col-span-2' },
  { id: 'pending',     defaultVisible: true,  span: 'sm:col-span-2 xl:col-span-2' },
]
const spanMap = Object.fromEntries(WIDGET_META.map(w => [w.id, w.span]))
function spanOf(id) { return spanMap[id] || 'xl:col-span-1' }

const WIDGET_LABELS = {
  users: 'dashboard.widgetUsers', providers: 'dashboard.widgetProviders', env: 'dashboard.widgetEnv',
  deposits: 'dashboard.widgetDeposits', prestations: 'dashboard.widgetPrestations', events: 'dashboard.widgetEvents',
  trends: 'dashboard.widgetTrends', activity: 'dashboard.widgetActivity', pending: 'dashboard.widgetPending',
}
function widgetLabel(id) { return t(WIDGET_LABELS[id] || id) }

// ---- Persisted layout (per admin, localStorage) ----
const storageKey = computed(() => `admin_dashboard_layout_v1_${auth.user?.id || 'default'}`)
const layout = ref([]) // [{ id, visible }]

function defaultLayout() { return WIDGET_META.map(w => ({ id: w.id, visible: w.defaultVisible })) }

function loadLayout() {
  try {
    const raw = JSON.parse(localStorage.getItem(storageKey.value) || 'null')
    if (Array.isArray(raw)) {
      const known = new Set(WIDGET_META.map(w => w.id))
      const merged = raw.filter(x => x && known.has(x.id)).map(x => ({ id: x.id, visible: !!x.visible }))
      const present = new Set(merged.map(x => x.id))
      for (const w of WIDGET_META) if (!present.has(w.id)) merged.push({ id: w.id, visible: w.defaultVisible })
      layout.value = merged
      return
    }
  } catch { /* fall through to default */ }
  layout.value = defaultLayout()
}
function saveLayout() { localStorage.setItem(storageKey.value, JSON.stringify(layout.value)) }
function resetLayout() { layout.value = defaultLayout(); saveLayout() }

const visibleIds = computed(() => layout.value.filter(x => x.visible).map(x => x.id))
const customizerItems = computed(() => layout.value.map(x => ({ id: x.id, visible: x.visible, label: widgetLabel(x.id) })))
function onUpdateItems(items) { layout.value = items; saveLayout() }

// ---- Data ----
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
  const csv = rows.map(r => r.map(v => `"${String(v).replace(/"/g, '""')}"`).join(',')).join('\n')
  const blob = new Blob(['﻿' + csv], { type: 'text/csv;charset=utf-8;' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `dashboard_${new Date().toISOString().slice(0, 10)}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

onMounted(() => {
  loadLayout()
  loadData()
})
</script>
