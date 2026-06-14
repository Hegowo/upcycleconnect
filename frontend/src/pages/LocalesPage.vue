<template>
  <div class="max-w-5xl space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Langues</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Ajoutez et traduisez des langues sans modifier le code du site.</p>
      </div>
      <button @click="openAdd" class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
        <PlusIcon class="w-4 h-4" /> Ajouter une langue
      </button>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>

    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-[#f8fafc] text-[#64748b] text-xs uppercase">
          <tr>
            <th class="text-left px-5 py-3">Langue</th>
            <th class="text-left px-5 py-3">Code</th>
            <th class="text-left px-5 py-3">Type</th>
            <th class="text-left px-5 py-3">Statut</th>
            <th class="text-right px-5 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="l in locales" :key="l.code" class="border-t border-[#f1f5f9]">
            <td class="px-5 py-3 font-semibold text-[#001d32]">{{ l.name }}</td>
            <td class="px-5 py-3 font-mono text-[#40617f]">{{ l.code }}</td>
            <td class="px-5 py-3">
              <span :class="l.builtin ? 'bg-blue-100 text-blue-700' : 'bg-[#dcfce7] text-[#166534]'" class="text-xs font-bold px-2 py-0.5 rounded-full">
                {{ l.builtin ? 'Intégrée' : 'Personnalisée' }}
              </span>
            </td>
            <td class="px-5 py-3">
              <span :class="l.enabled ? 'bg-[#dcfce7] text-[#166534]' : 'bg-gray-100 text-gray-500'" class="text-xs font-bold px-2 py-0.5 rounded-full">
                {{ l.enabled ? 'Active' : 'Désactivée' }}
              </span>
            </td>
            <td class="px-5 py-3 text-right space-x-1.5">
              <button v-if="!l.builtin" @click="openEdit(l)" class="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs rounded-lg border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35] transition">
                <PencilSquareIcon class="w-3.5 h-3.5" /> Traduire
              </button>
              <button v-if="!l.builtin" @click="toggle(l)" class="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs rounded-lg border border-[#e5e7eb] hover:border-[#006d35] transition">
                {{ l.enabled ? 'Désactiver' : 'Activer' }}
              </button>
              <button v-if="!l.builtin" @click="destroy(l)" class="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs rounded-lg border border-red-200 text-red-600 hover:bg-red-50 transition">
                <TrashIcon class="w-3.5 h-3.5" />
              </button>
              <span v-if="l.builtin" class="text-xs text-[#94a3b8]">—</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Teleport to="body">
      <div v-if="showAdd" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showAdd = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-1">Nouvelle langue</h3>
          <p class="text-xs text-[#64748b] mb-4">La langue sera pré-remplie avec les textes français, prêts à traduire.</p>
          <div class="space-y-3">
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Code (ex: es, de, pt-BR)</label>
              <input v-model="addForm.code" type="text" placeholder="es" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm font-mono focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Nom affiché</label>
              <input v-model="addForm.name" type="text" placeholder="Español" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
            <p v-if="addError" class="text-red-600 text-sm">{{ addError }}</p>
          </div>
          <div class="flex justify-end gap-3 mt-5">
            <button @click="showAdd = false" class="px-4 py-2 text-sm text-[#40617f] font-medium">Annuler</button>
            <button @click="createLocale" :disabled="adding" class="px-5 py-2 text-white text-sm font-semibold rounded-xl hover:opacity-90 transition disabled:opacity-50" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              {{ adding ? 'Création...' : 'Créer' }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="showEdit" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showEdit = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-6xl h-[92vh] flex flex-col overflow-hidden">

          <div class="px-6 py-4 border-b flex items-center justify-between shrink-0 gap-4" style="background:linear-gradient(135deg,#006d35,#1b8848);">
            <div class="min-w-0">
              <h3 class="text-white font-bold text-lg leading-tight">Traduction — {{ editLocale?.name }}</h3>
              <p class="text-white/70 text-xs mt-0.5">Référence : Français · {{ totalKeys }} textes</p>
            </div>
            <div class="flex items-center gap-4 shrink-0">
              <div class="text-right">
                <div class="flex items-center gap-2">
                  <div class="w-32 h-2 bg-white/25 rounded-full overflow-hidden">
                    <div class="h-full bg-white rounded-full transition-all" :style="{ width: globalProgress + '%' }" />
                  </div>
                  <span class="text-white font-bold text-sm">{{ globalProgress }}%</span>
                </div>
                <p class="text-white/70 text-[11px] mt-0.5">{{ translatedCount }} / {{ totalKeys }} traduits</p>
              </div>
              <button @click="showEdit = false" class="text-white/70 hover:text-white p-2 rounded-lg hover:bg-white/10"><XMarkIcon class="w-5 h-5" /></button>
            </div>
          </div>

          <div class="px-6 py-3 border-b border-gray-100 shrink-0 flex flex-wrap items-center gap-3">
            <div class="relative flex-1 min-w-[200px]">
              <MagnifyingGlassIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#94a3b8]" />
              <input v-model="search" type="search" placeholder="Rechercher dans les textes..." class="w-full pl-9 pr-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
            <div class="flex gap-1.5">
              <button v-for="f in filters" :key="f.value" @click="filterMode = f.value"
                :class="['px-3 py-1.5 rounded-full text-xs font-semibold border transition',
                  filterMode === f.value ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]']">
                {{ f.label }}
              </button>
            </div>
          </div>

          <div class="flex-1 flex min-h-0">

            <div class="w-56 border-r border-gray-100 overflow-y-auto bg-[#fafbfc] shrink-0">
              <button v-for="s in sections" :key="s.name" @click="activeSection = s.name"
                :class="['w-full text-left px-4 py-2.5 border-b border-gray-50 transition',
                  activeSection === s.name ? 'bg-[#f0fdf4]' : 'hover:bg-white']">
                <div class="flex items-center justify-between gap-2">
                  <span class="text-sm font-medium truncate" :class="activeSection === s.name ? 'text-[#006d35]' : 'text-[#001d32]'">{{ sectionLabel(s.name) }}</span>
                  <span class="text-[10px] font-bold shrink-0" :class="s.done === s.total ? 'text-[#16a34a]' : 'text-[#94a3b8]'">{{ s.done }}/{{ s.total }}</span>
                </div>
                <div class="w-full h-1 bg-gray-200 rounded-full mt-1.5 overflow-hidden">
                  <div class="h-full rounded-full transition-all" :class="s.done === s.total ? 'bg-[#16a34a]' : 'bg-[#006d35]'" :style="{ width: (s.total ? (s.done/s.total*100) : 0) + '%' }" />
                </div>
              </button>
            </div>

            <div class="flex-1 overflow-y-auto p-5 space-y-3">
              <div v-if="!visibleKeys.length" class="text-center py-16 text-[#94a3b8] text-sm">
                Aucun texte ne correspond à ce filtre.
              </div>
              <div v-for="k in visibleKeys" :key="k"
                class="rounded-xl border p-3.5 transition"
                :class="isTranslated(k) ? 'border-[#e5e7eb] bg-white' : 'border-amber-200 bg-amber-50/40'">
                <div class="flex items-center justify-between gap-2 mb-2">
                  <span class="font-mono text-[10px] text-[#94a3b8] truncate">{{ k }}</span>
                  <span v-if="!isTranslated(k)" class="shrink-0 text-[10px] font-bold uppercase tracking-wide text-amber-600 bg-amber-100 px-2 py-0.5 rounded-full">À traduire</span>
                  <span v-else class="shrink-0 text-[10px] font-bold uppercase tracking-wide text-[#16a34a] bg-[#dcfce7] px-2 py-0.5 rounded-full">Traduit</span>
                </div>
                <div class="grid sm:grid-cols-2 gap-3">
                  <div class="bg-[#f8fafc] rounded-lg px-3 py-2 text-sm text-[#64748b] leading-relaxed">
                    <span class="block text-[9px] font-bold uppercase tracking-wider text-[#cbd5e1] mb-0.5">🇫🇷 Français</span>
                    {{ frFlat[k] }}
                  </div>
                  <div>
                    <span class="block text-[9px] font-bold uppercase tracking-wider text-[#bbf7d0] mb-0.5">{{ editLocale?.name }}</span>
                    <textarea v-model="editFlat[k]" :placeholder="frFlat[k]" rows="2"
                      class="w-full px-3 py-2 bg-white border border-[#e5e7eb] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-y leading-relaxed" />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="px-6 py-3.5 border-t bg-[#f8fafc] flex items-center justify-between gap-3 shrink-0">
            <p class="text-xs text-[#64748b]">
              <span v-if="dirtyCount > 0" class="text-[#006d35] font-semibold">{{ dirtyCount }} modification(s) non enregistrée(s)</span>
              <span v-else>Tout est enregistré</span>
            </p>
            <div class="flex gap-2">
              <button @click="showEdit = false" class="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700">Fermer</button>
              <button @click="saveTranslations" :disabled="savingEdit" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm flex items-center gap-2 hover:opacity-90 disabled:opacity-50 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
                <div v-if="savingEdit" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                {{ savingEdit ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onMounted } from 'vue'
import { PlusIcon, PencilSquareIcon, TrashIcon, XMarkIcon, MagnifyingGlassIcon } from '@heroicons/vue/24/outline'
import frMessages from '@/locales/fr.json'

const BASE = '/api/admin/v1'
function authHeaders() {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const locales = ref([])
const loading = ref(true)

function flatten(obj, prefix = '', out = {}) {
  for (const [k, v] of Object.entries(obj || {})) {
    const key = prefix ? `${prefix}.${k}` : k
    if (v && typeof v === 'object' && !Array.isArray(v)) flatten(v, key, out)
    else out[key] = v
  }
  return out
}
function unflatten(flat) {
  const out = {}
  for (const [key, val] of Object.entries(flat)) {
    const parts = key.split('.')
    let cur = out
    parts.forEach((p, i) => {
      if (i === parts.length - 1) cur[p] = val
      else { cur[p] = cur[p] || {}; cur = cur[p] }
    })
  }
  return out
}

const frFlat = flatten(frMessages)

async function fetchLocales() {
  loading.value = true
  try {
    const res = await fetch(`${BASE}/locales`, { headers: authHeaders() })
    const j = await res.json()
    locales.value = j.data || []
  } finally { loading.value = false }
}

const showAdd = ref(false)
const adding = ref(false)
const addError = ref('')
const addForm = ref({ code: '', name: '' })
function openAdd() { addForm.value = { code: '', name: '' }; addError.value = ''; showAdd.value = true }
async function createLocale() {
  if (!addForm.value.code || !addForm.value.name) { addError.value = 'Code et nom requis'; return }
  adding.value = true
  addError.value = ''
  try {
    const res = await fetch(`${BASE}/locales`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({ code: addForm.value.code.trim(), name: addForm.value.name.trim(), messages: frMessages, enabled: true }),
    })
    if (!res.ok) { const e = await res.json().catch(() => ({})); throw new Error(e.message || 'Erreur') }
    showAdd.value = false
    await fetchLocales()
  } catch (e) { addError.value = e.message } finally { adding.value = false }
}

const showEdit = ref(false)
const savingEdit = ref(false)
const editLocale = ref(null)
const editFlat = ref({})
const savedSnapshot = ref({})
const search = ref('')
const filterMode = ref('all')
const activeSection = ref('')

const filters = [
  { value: 'all',  label: 'Tous' },
  { value: 'todo', label: 'À traduire' },
  { value: 'done', label: 'Traduits' },
]

const SECTION_LABELS = {
  common: 'Commun', nav: 'Navigation', login: 'Connexion', dashboard: 'Tableau de bord',
  users: 'Utilisateurs', providers: 'Prestataires', categories: 'Catégories',
  prestations: 'Prestations', events: 'Événements', admins: 'Administrateurs',
  logs: 'Journaux', settings: 'Paramètres', public: 'Site public',
}
function sectionLabel(name) { return SECTION_LABELS[name] || name }

function isTranslated(k) {
  const v = (editFlat.value[k] || '').trim()
  return v !== '' && v !== String(frFlat[k]).trim()
}

const allKeys = Object.keys(frFlat)
const totalKeys = computed(() => allKeys.length)
const translatedCount = computed(() => allKeys.filter(isTranslated).length)
const globalProgress = computed(() => totalKeys.value ? Math.round(translatedCount.value / totalKeys.value * 100) : 0)
const dirtyCount = computed(() => allKeys.filter(k => (editFlat.value[k] || '') !== (savedSnapshot.value[k] || '')).length)

const sections = computed(() => {
  const map = {}
  for (const k of allKeys) {
    const top = k.split('.')[0]
    if (!map[top]) map[top] = { name: top, total: 0, done: 0 }
    map[top].total++
    if (isTranslated(k)) map[top].done++
  }
  return Object.values(map).sort((a, b) => a.name.localeCompare(b.name))
})

const visibleKeys = computed(() => {
  const q = search.value.trim().toLowerCase()
  return allKeys.filter(k => {
    if (!q && activeSection.value && k.split('.')[0] !== activeSection.value) return false
    if (q) {
      const hit = k.toLowerCase().includes(q) || String(frFlat[k]).toLowerCase().includes(q) || String(editFlat.value[k] || '').toLowerCase().includes(q)
      if (!hit) return false
    }
    if (filterMode.value === 'todo' && isTranslated(k)) return false
    if (filterMode.value === 'done' && !isTranslated(k)) return false
    return true
  })
})

async function openEdit(l) {
  const res = await fetch(`${BASE}/locales/${l.code}`, { headers: authHeaders() })
  const data = await res.json()
  editLocale.value = data
  editFlat.value = { ...frFlat, ...flatten(data.messages || {}) }
  savedSnapshot.value = { ...editFlat.value }
  search.value = ''
  filterMode.value = 'all'
  activeSection.value = sections.value[0]?.name || ''
  showEdit.value = true
}

async function saveTranslations() {
  savingEdit.value = true
  try {
    await fetch(`${BASE}/locales/${editLocale.value.code}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({ messages: unflatten(editFlat.value) }),
    })
    savedSnapshot.value = { ...editFlat.value }
    showEdit.value = false
    fetchLocales()
  } finally { savingEdit.value = false }
}

async function toggle(l) {
  await fetch(`${BASE}/locales/${l.code}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify({ enabled: !l.enabled }),
  })
  fetchLocales()
}
async function destroy(l) {
  if (!confirm(`Supprimer la langue "${l.name}" ?`)) return
  await fetch(`${BASE}/locales/${l.code}`, { method: 'DELETE', headers: authHeaders() })
  fetchLocales()
}

onMounted(fetchLocales)
</script>
