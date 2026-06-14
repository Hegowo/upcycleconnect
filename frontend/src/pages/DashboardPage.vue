<template>
  <div class="space-y-4 sm:space-y-6">

    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">{{ t('dashboard.title') }}</h2>
        <p class="text-xs sm:text-sm text-[#40617f] mt-0.5">{{ todayFormatted }}</p>
      </div>

      <div class="flex items-center gap-2 shrink-0">
        <template v-if="editMode">
          <div class="relative">
            <button @click.stop="addMenu = !addMenu" :disabled="hiddenWidgets.length === 0"
              class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition disabled:opacity-40">
              <PlusIcon class="w-4 h-4 shrink-0" />
              <span class="hidden sm:inline">{{ t('dashboard.addWidget') }}</span>
            </button>
            <div v-if="addMenu" class="absolute right-0 mt-1 w-56 bg-white rounded-xl border border-[#e5e7eb] shadow-lg z-30 p-1.5">
              <button v-for="h in hiddenWidgets" :key="h.id" @click.stop="showWidget(h.id)"
                class="w-full text-left px-3 py-2 rounded-lg text-sm text-[#001d32] hover:bg-[#f0fdf4] transition">
                {{ widgetLabel(h.id) }}
              </button>
              <p v-if="hiddenWidgets.length === 0" class="px-3 py-2 text-xs text-gray-400">{{ t('dashboard.allWidgetsShown') }}</p>
            </div>
          </div>
          <button @click="resetLayout" class="px-2 sm:px-4 py-2 text-xs sm:text-sm font-medium rounded-lg text-gray-500 hover:bg-gray-50 transition">
            {{ t('dashboard.layoutReset') }}
          </button>
          <button @click="exitEdit" class="flex items-center gap-1.5 px-3 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg text-white transition hover:opacity-90" style="background-color:#006d35;">
            <CheckIcon class="w-4 h-4" />
            <span class="hidden sm:inline">{{ t('dashboard.doneEditing') }}</span>
          </button>
        </template>

        <template v-else>
          <button @click="editMode = true" class="flex items-center gap-1.5 px-2 sm:px-4 py-2 text-xs sm:text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition">
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
        </template>
      </div>
    </div>

    <div v-if="editMode" class="flex items-center gap-2 text-xs text-[#40617f] bg-[#edf4ff] border border-[#dbeafe] rounded-xl px-3 py-2">
      <InformationCircleIcon class="w-4 h-4 shrink-0" />
      {{ t('dashboard.editHint') }}
    </div>

    <div v-if="visibleTiles.length === 0" class="text-center py-16 text-gray-400 text-sm border border-dashed border-[#e2e8f0] rounded-2xl">
      {{ t('dashboard.emptyLayout') }}
      <button @click="editMode = true" class="block mx-auto mt-2 text-[#006d35] font-semibold hover:underline">{{ t('dashboard.customize') }}</button>
    </div>

    <div v-else class="dash-grid" :class="{ 'is-editing': editMode }">
      <div
        v-for="tile in visibleTiles"
        :key="tile.id"
        class="dash-tile"
        :class="{ dragging: dragId === tile.id, 'drag-over': dragOverId === tile.id && dragId && dragId !== tile.id }"
        :style="{ '--w': tile.w, '--h': tile.h }"
        :draggable="editMode"
        @dragstart="onDragStart(tile)"
        @dragover.prevent="dragOverId = tile.id"
        @drop="onDrop(tile)"
        @dragend="onDragEnd"
      >
        <div class="tile-content" :class="{ 'pointer-events-none select-none': editMode }">
          <StatWidget v-if="tile.id === 'users'" :icon="UsersIcon" iconBg="#dcfce7" iconColor="#006d35"
            :value="stats?.users_count ?? 0" :label="t('dashboard.statsUsers')" :cornerLabel="t('dashboard.members')" to="/admin/users" />
          <StatWidget v-else-if="tile.id === 'providers'" :icon="BriefcaseIcon" iconBg="#dbeafe" iconColor="#2563eb"
            :value="stats?.providers_count ?? 0" :label="t('dashboard.statsProviders')" :cornerLabel="t('dashboard.statsActive')" to="/admin/providers?status=approved" />
          <StatWidget v-else-if="tile.id === 'env'" :icon="GlobeAltIcon" iconBg="#f0fdf4" iconColor="#006d35"
            :value="envScore" suffix="/100" :label="t('dashboard.envScore')" :cornerLabel="t('dashboard.envScoreImpact')" :tooltip="t('dashboard.envScoreFormula')" />
          <StatWidget v-else-if="tile.id === 'deposits'" :icon="CubeIcon" iconBg="#fef9c3" iconColor="#ca8a04"
            :value="stats?.deposits_accepted ?? 0" :label="t('dashboard.objectsCollected')" :cornerLabel="t('dashboard.totalLabel')" />
          <StatWidget v-else-if="tile.id === 'prestations'" :icon="TagIcon" iconBg="#ede9fe" iconColor="#7c3aed"
            :value="stats?.prestations_count ?? 0" :label="t('dashboard.statsPrestations')" to="/admin/prestations" />
          <StatWidget v-else-if="tile.id === 'events'" :icon="CalendarDaysIcon" iconBg="#e0f2fe" iconColor="#0891b2"
            :value="stats?.events_count ?? 0" :label="t('dashboard.statsEvents')" />
          <TrendsWidget v-else-if="tile.id === 'trends'" />
          <ActivityWidget v-else-if="tile.id === 'activity'" :logs="logs" :loading="logsLoading" />
          <PendingProvidersWidget v-else-if="tile.id === 'pending'" :pending="stats?.providers_pending ?? 0" />
        </div>

        <div v-if="editMode" class="tile-tools">
          <span class="tile-btn cursor-grab active:cursor-grabbing" :title="t('dashboard.dragHint')"><Bars3Icon class="w-4 h-4" /></span>
          <div class="relative">
            <button class="tile-btn font-semibold tabular-nums" @click.stop="sizeMenu = sizeMenu === tile.id ? null : tile.id" :title="t('dashboard.size')">
              {{ tile.w }}×{{ tile.h }}
            </button>
            <div v-if="sizeMenu === tile.id" class="absolute right-0 mt-1 bg-white rounded-xl border border-[#e5e7eb] shadow-lg z-30 p-1.5 grid grid-cols-2 gap-1 w-28">
              <button v-for="s in sizesOf(tile.id)" :key="s.join('x')" @click.stop="setSize(tile, s)"
                class="px-2 py-1.5 rounded-lg text-xs font-medium tabular-nums transition"
                :class="s[0] === tile.w && s[1] === tile.h ? 'bg-[#006d35] text-white' : 'text-[#001d32] hover:bg-[#f0fdf4]'">
                {{ s[0] }}×{{ s[1] }}
              </button>
            </div>
          </div>
          <button class="tile-btn" @click.stop="hideWidget(tile.id)" :title="t('dashboard.hideWidget')"><EyeSlashIcon class="w-4 h-4" /></button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import {
  UsersIcon, BriefcaseIcon, GlobeAltIcon, CubeIcon, TagIcon, CalendarDaysIcon,
  AdjustmentsHorizontalIcon, ArrowDownTrayIcon, ArrowPathIcon,
  PlusIcon, CheckIcon, Bars3Icon, EyeSlashIcon, InformationCircleIcon,
} from '@heroicons/vue/24/outline'
import StatWidget from '@/components/dashboard/StatWidget.vue'
import TrendsWidget from '@/components/dashboard/TrendsWidget.vue'
import ActivityWidget from '@/components/dashboard/ActivityWidget.vue'
import PendingProvidersWidget from '@/components/dashboard/PendingProvidersWidget.vue'

const { t } = useI18n()
const auth = useAuthStore()

const stats = ref(null)
const logs  = ref([])
const loading     = ref(true)
const logsLoading = ref(true)

const editMode  = ref(false)
const addMenu   = ref(false)
const sizeMenu  = ref(null)
const dragId    = ref(null)
const dragOverId = ref(null)

const todayFormatted = computed(() =>
  new Date().toLocaleDateString('fr-FR', { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' }))

const envScore = computed(() => {
  if (!stats.value) return 0
  const accepted = stats.value.deposits_accepted ?? 0
  const co2 = stats.value.co2_total ?? 0
  return Math.min(100, Math.round(accepted * 5 + co2 * 1.5))
})

const WIDGET_META = [
  { id: 'users',       defaultVisible: true,  w: 1, h: 1, sizes: [[1, 1], [2, 1]] },
  { id: 'providers',   defaultVisible: true,  w: 1, h: 1, sizes: [[1, 1], [2, 1]] },
  { id: 'env',         defaultVisible: true,  w: 1, h: 1, sizes: [[1, 1], [2, 1]] },
  { id: 'deposits',    defaultVisible: true,  w: 1, h: 1, sizes: [[1, 1], [2, 1]] },
  { id: 'prestations', defaultVisible: false, w: 1, h: 1, sizes: [[1, 1], [2, 1]] },
  { id: 'events',      defaultVisible: false, w: 1, h: 1, sizes: [[1, 1], [2, 1]] },
  { id: 'trends',      defaultVisible: true,  w: 2, h: 2, sizes: [[2, 2], [4, 2], [2, 3], [3, 2]] },
  { id: 'activity',    defaultVisible: true,  w: 2, h: 3, sizes: [[2, 3], [2, 4], [4, 3], [2, 2]] },
  { id: 'pending',     defaultVisible: true,  w: 2, h: 1, sizes: [[2, 1], [2, 2], [1, 1], [4, 1]] },
]
const metaMap = Object.fromEntries(WIDGET_META.map(w => [w.id, w]))
function sizesOf(id) { return metaMap[id]?.sizes || [[1, 1]] }

const WIDGET_LABELS = {
  users: 'dashboard.widgetUsers', providers: 'dashboard.widgetProviders', env: 'dashboard.widgetEnv',
  deposits: 'dashboard.widgetDeposits', prestations: 'dashboard.widgetPrestations', events: 'dashboard.widgetEvents',
  trends: 'dashboard.widgetTrends', activity: 'dashboard.widgetActivity', pending: 'dashboard.widgetPending',
}
function widgetLabel(id) { return t(WIDGET_LABELS[id] || id) }

const storageKey = computed(() => `admin_dashboard_layout_v4_${auth.user?.id || 'default'}`)
const layout = ref([])

function defaultLayout() {
  return WIDGET_META.map(w => ({ id: w.id, visible: w.defaultVisible, w: w.w, h: w.h }))
}
function loadLayout() {
  try {
    const raw = JSON.parse(localStorage.getItem(storageKey.value) || 'null')
    if (Array.isArray(raw)) {
      const merged = raw
        .filter(x => x && metaMap[x.id])
        .map(x => ({ id: x.id, visible: !!x.visible, w: clampW(x.w, x.id), h: clampH(x.h, x.id) }))
      const present = new Set(merged.map(x => x.id))
      for (const w of WIDGET_META) if (!present.has(w.id)) merged.push({ id: w.id, visible: w.defaultVisible, w: w.w, h: w.h })
      layout.value = merged
      return
    }
  } catch {  }
  layout.value = defaultLayout()
}
function clampW(v, id) { const def = metaMap[id]; return Math.min(4, Math.max(1, Number(v) || def.w)) }
function clampH(v, id) { return Math.min(4, Math.max(1, Number(v) || metaMap[id].h)) }
function saveLayout() { localStorage.setItem(storageKey.value, JSON.stringify(layout.value)) }
function resetLayout() { layout.value = defaultLayout(); saveLayout() }

const visibleTiles  = computed(() => layout.value.filter(x => x.visible))
const hiddenWidgets = computed(() => layout.value.filter(x => !x.visible))

function setSize(tile, s) {
  const item = layout.value.find(x => x.id === tile.id)
  if (item) { item.w = s[0]; item.h = s[1]; saveLayout() }
  sizeMenu.value = null
}
function hideWidget(id) {
  const item = layout.value.find(x => x.id === id)
  if (item) { item.visible = false; saveLayout() }
  sizeMenu.value = null
}
function showWidget(id) {
  const item = layout.value.find(x => x.id === id)
  if (item) { item.visible = true; saveLayout() }
  addMenu.value = false
}
function exitEdit() { editMode.value = false; addMenu.value = false; sizeMenu.value = null }

function onDragStart(tile) { dragId.value = tile.id; sizeMenu.value = null }
function onDrop(target) {
  if (!dragId.value || dragId.value === target.id) return
  const from = layout.value.findIndex(x => x.id === dragId.value)
  const to   = layout.value.findIndex(x => x.id === target.id)
  if (from < 0 || to < 0) return
  const [moved] = layout.value.splice(from, 1)
  layout.value.splice(to, 0, moved)
  saveLayout()
}
function onDragEnd() { dragId.value = null; dragOverId.value = null }

function onDocClick() { addMenu.value = false; sizeMenu.value = null }

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
  window.addEventListener('click', onDocClick)
})
onUnmounted(() => window.removeEventListener('click', onDocClick))
</script>

<style scoped>
.dash-grid {
  display: grid;
  grid-template-columns: 1fr;
  grid-auto-rows: 9rem;
  grid-auto-flow: row dense;
  gap: 0.75rem;
}
.dash-tile {
  position: relative;
  grid-column: 1 / -1;
  grid-row: span var(--h, 1);
  min-width: 0;
  border-radius: 1rem;
}
.dash-tile > .tile-content { height: 100%; overflow: hidden; border-radius: 1rem; }

@media (min-width: 640px) {
  .dash-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .dash-tile { grid-column: span min(var(--w, 1), 2); }
}
@media (min-width: 1280px) {
  .dash-grid { grid-template-columns: repeat(4, minmax(0, 1fr)); }
  .dash-tile { grid-column: span min(var(--w, 1), 4); }
}

.is-editing .dash-tile {
  cursor: grab;
  outline: 2px dashed #cbd5e1;
  outline-offset: -2px;
}
.is-editing .dash-tile:active { cursor: grabbing; }
.dash-tile.dragging { opacity: 0.45; }
.dash-tile.drag-over { outline: 2px solid #006d35; outline-offset: -2px; }

.tile-tools {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  z-index: 10;
}
.tile-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 1.75rem;
  height: 1.75rem;
  padding: 0 0.4rem;
  border-radius: 0.5rem;
  background: #ffffff;
  border: 1px solid #e5e7eb;
  color: #40617f;
  font-size: 0.7rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.06);
  transition: background 0.15s;
}
.tile-btn:hover { background: #f0fdf4; color: #006d35; }
</style>
