<template>
  <div class="max-w-4xl space-y-6">
    <div>
      <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Export & Audit des données</h2>
      <p class="text-sm text-[#40617f] mt-0.5">Exportez n'importe quelle donnée de la plateforme en CSV ou PDF, sur la période de votre choix.</p>
    </div>

    <div class="card p-5 space-y-4">
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div>
          <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Du (optionnel)</label>
          <input v-model="from" type="date" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Au (optionnel)</label>
          <input v-model="to" type="date" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
        </div>
        <div>
          <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Raccourcis</label>
          <div class="flex gap-1.5 flex-wrap">
            <button v-for="p in presets" :key="p.label" @click="applyPreset(p)" class="px-2.5 py-1 text-xs rounded-lg border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35] transition">{{ p.label }}</button>
          </div>
        </div>
      </div>
      <p class="text-xs text-[#94a3b8]">Sans dates, l'export contient toutes les données. Le filtre s'applique à la date de création.</p>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <div v-for="ds in datasets" :key="ds.key" class="card p-5 flex items-center justify-between gap-4">
        <div class="min-w-0">
          <p class="font-jakarta font-bold text-[#001d32]">{{ ds.label }}</p>
          <p class="text-xs text-[#94a3b8] mt-0.5">{{ ds.columns }} colonnes</p>
        </div>
        <div class="flex gap-2 shrink-0">
          <button @click="download(ds.key, 'csv')" :disabled="busy === ds.key+'csv'"
            class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg border border-[#006d35] text-[#006d35] text-xs font-bold hover:bg-[#f0fdf4] transition disabled:opacity-50">
            <TableCellsIcon class="w-4 h-4" /> CSV
          </button>
          <button @click="download(ds.key, 'pdf')" :disabled="busy === ds.key+'pdf'"
            class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg text-white text-xs font-bold hover:opacity-90 transition disabled:opacity-50" style="background:linear-gradient(135deg,#006d35,#1b8848);">
            <DocumentTextIcon class="w-4 h-4" /> PDF
          </button>
        </div>
      </div>
    </div>

    <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { TableCellsIcon, DocumentTextIcon } from '@heroicons/vue/24/outline'

const BASE = '/api/admin/v1'
function token() { return localStorage.getItem('admin_token') }
function authHeaders() { const t = token(); return t ? { Authorization: `Bearer ${t}` } : {} }

const datasets = ref([])
const loading = ref(true)
const from = ref('')
const to = ref('')
const busy = ref('')
const error = ref('')

const presets = [
  { label: '7 jours',  days: 7 },
  { label: '30 jours', days: 30 },
  { label: 'Cette année', year: true },
  { label: 'Tout', clear: true },
]
function applyPreset(p) {
  const now = new Date()
  if (p.clear) { from.value = ''; to.value = ''; return }
  to.value = now.toISOString().slice(0, 10)
  if (p.year) { from.value = `${now.getFullYear()}-01-01`; return }
  const d = new Date(now); d.setDate(d.getDate() - p.days)
  from.value = d.toISOString().slice(0, 10)
}

async function fetchDatasets() {
  loading.value = true
  try {
    const res = await fetch(`${BASE}/export/datasets`, { headers: authHeaders() })
    const j = await res.json()
    datasets.value = (j.data || []).sort((a, b) => a.label.localeCompare(b.label))
  } finally { loading.value = false }
}

async function download(dataset, format) {
  busy.value = dataset + format
  error.value = ''
  try {
    const params = new URLSearchParams({ dataset, format })
    if (from.value) params.set('from', from.value)
    if (to.value) params.set('to', to.value)
    const res = await fetch(`${BASE}/export?${params}`, { headers: authHeaders() })
    if (!res.ok) { const j = await res.json().catch(() => ({})); throw new Error(j.message || 'Erreur export') }
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `upcycleconnect-${dataset}-${new Date().toISOString().slice(0,10)}.${format}`
    a.click()
    URL.revokeObjectURL(url)
  } catch (e) { error.value = e.message } finally { busy.value = '' }
}

onMounted(fetchDatasets)
</script>
