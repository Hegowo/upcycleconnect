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

    <!-- Add language modal -->
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

      <!-- Translation editor -->
      <div v-if="showEdit" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showEdit = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-3xl max-h-[92vh] flex flex-col overflow-hidden">
          <div class="px-6 py-4 border-b flex items-center justify-between shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">
            <h3 class="text-white font-bold text-lg">Traduire — {{ editLocale?.name }}</h3>
            <button @click="showEdit = false" class="text-white/70 hover:text-white p-2 rounded-lg hover:bg-white/10"><XMarkIcon class="w-5 h-5" /></button>
          </div>
          <div class="px-6 py-3 border-b border-gray-100 shrink-0">
            <input v-model="search" type="search" placeholder="Rechercher une clé ou un texte..." class="w-full px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            <p class="text-xs text-[#94a3b8] mt-1.5">{{ filteredKeys.length }} clé(s) — colonne gauche = français (référence), droite = traduction.</p>
          </div>
          <div class="flex-1 overflow-y-auto p-4 space-y-2">
            <div v-for="k in filteredKeys" :key="k" class="grid grid-cols-2 gap-3 items-start">
              <div class="bg-[#f8fafc] rounded-lg px-3 py-2 text-xs text-[#64748b]">
                <p class="font-mono text-[10px] text-[#94a3b8] mb-0.5 truncate">{{ k }}</p>
                {{ frFlat[k] }}
              </div>
              <input v-model="editFlat[k]" type="text" :placeholder="frFlat[k]" class="w-full px-3 py-2 bg-white border border-[#e5e7eb] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
          </div>
          <div class="px-6 py-4 border-t bg-[#f8fafc] flex items-center justify-between gap-3 shrink-0">
            <button @click="showEdit = false" class="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700">Fermer</button>
            <button @click="saveTranslations" :disabled="savingEdit" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm flex items-center gap-2 hover:opacity-90 disabled:opacity-50 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <div v-if="savingEdit" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              {{ savingEdit ? 'Enregistrement...' : 'Enregistrer les traductions' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onMounted } from 'vue'
import { PlusIcon, PencilSquareIcon, TrashIcon, XMarkIcon } from '@heroicons/vue/24/outline'
import frMessages from '@/locales/fr.json'

const BASE = '/api/admin/v1'
function authHeaders() {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const locales = ref([])
const loading = ref(true)

// Flatten/unflatten helpers — keys like "public.home.title"
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

// ── Add ──
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
    // Seed the new language with the FR message tree so the admin only has to translate.
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

// ── Edit translations ──
const showEdit = ref(false)
const savingEdit = ref(false)
const editLocale = ref(null)
const editFlat = ref({})
const search = ref('')
async function openEdit(l) {
  const res = await fetch(`${BASE}/locales/${l.code}`, { headers: authHeaders() })
  const data = await res.json()
  editLocale.value = data
  editFlat.value = { ...frFlat, ...flatten(data.messages || {}) }
  search.value = ''
  showEdit.value = true
}
const filteredKeys = computed(() => {
  const keys = Object.keys(frFlat)
  const q = search.value.trim().toLowerCase()
  if (!q) return keys
  return keys.filter(k => k.toLowerCase().includes(q) || String(frFlat[k]).toLowerCase().includes(q) || String(editFlat.value[k] || '').toLowerCase().includes(q))
})
async function saveTranslations() {
  savingEdit.value = true
  try {
    await fetch(`${BASE}/locales/${editLocale.value.code}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({ messages: unflatten(editFlat.value) }),
    })
    showEdit.value = false
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
