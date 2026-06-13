<template>
  <div class="max-w-5xl space-y-6">
    <div class="flex items-center justify-between flex-wrap gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Notifications</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Historique des notifications envoyées et diffusion à une audience.</p>
      </div>
      <button @click="showBroadcast = true" class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
        <PaperAirplaneIcon class="w-4 h-4" /> Diffuser une notification
      </button>
    </div>

    <!-- Filters -->
    <div class="card p-4 flex flex-wrap gap-3 items-center">
      <input v-model="search" @input="debouncedFetch" type="search" placeholder="Rechercher..." class="flex-1 min-w-[160px] text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
      <select v-model="typeFilter" @change="fetchSent(1)" class="text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
        <option value="">Tous les types</option>
        <option v-for="t in knownTypes" :key="t" :value="t">{{ t }}</option>
      </select>
    </div>

    <div v-if="loading" class="py-16 text-center">
      <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>

    <div v-else-if="!items.length" class="card p-10 text-center text-[#64748b] text-sm">Aucune notification.</div>

    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-[#f8fafc] text-[#64748b] text-xs uppercase">
          <tr>
            <th class="text-left px-5 py-3">Notification</th>
            <th class="text-left px-5 py-3 hidden sm:table-cell">Destinataire</th>
            <th class="text-left px-5 py-3 hidden md:table-cell">Type</th>
            <th class="text-left px-5 py-3">Date</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="n in items" :key="n.id" class="border-t border-[#f1f5f9]">
            <td class="px-5 py-3">
              <p class="font-semibold text-[#001d32]">{{ n.title }}</p>
              <p class="text-[#64748b] text-xs line-clamp-1">{{ n.body }}</p>
            </td>
            <td class="px-5 py-3 hidden sm:table-cell">
              <p class="text-[#001d32]">{{ n.recipient }}</p>
              <p class="text-[#94a3b8] text-xs">{{ n.email }}</p>
            </td>
            <td class="px-5 py-3 hidden md:table-cell">
              <span class="text-xs font-mono text-[#40617f] bg-[#f1f5f9] px-2 py-0.5 rounded">{{ n.type }}</span>
            </td>
            <td class="px-5 py-3 text-[#64748b] text-xs whitespace-nowrap">{{ formatDate(n.created_at) }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="meta.last_page > 1" class="flex items-center justify-between px-5 py-3 border-t border-[#f1f5f9]">
        <button @click="fetchSent(meta.current_page - 1)" :disabled="meta.current_page <= 1" class="text-sm text-[#40617f] disabled:opacity-40">← Précédent</button>
        <span class="text-xs text-[#94a3b8]">Page {{ meta.current_page }} / {{ meta.last_page }}</span>
        <button @click="fetchSent(meta.current_page + 1)" :disabled="meta.current_page >= meta.last_page" class="text-sm text-[#40617f] disabled:opacity-40">Suivant →</button>
      </div>
    </div>

    <!-- Broadcast modal -->
    <Teleport to="body">
      <div v-if="showBroadcast" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showBroadcast = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-4">Diffuser une notification</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Audience</label>
              <select v-model="bcast.audience" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                <option value="all">Tous les utilisateurs</option>
                <option value="particuliers">Particuliers</option>
                <option value="pros">Professionnels / artisans</option>
              </select>
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Titre</label>
              <input v-model="bcast.title" type="text" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Message</label>
              <textarea v-model="bcast.body" rows="3" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Lien (optionnel)</label>
              <input v-model="bcast.link" type="text" placeholder="/prestations" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
            <p v-if="bcastMsg" :class="bcastError ? 'text-red-600' : 'text-green-600'" class="text-sm">{{ bcastMsg }}</p>
          </div>
          <div class="flex justify-end gap-3 mt-5">
            <button @click="showBroadcast = false" class="px-4 py-2 text-sm text-[#40617f] font-medium">Annuler</button>
            <button @click="sendBroadcast" :disabled="sending" class="px-5 py-2 text-white text-sm font-semibold rounded-xl hover:opacity-90 transition disabled:opacity-50 flex items-center gap-2" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <div v-if="sending" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              {{ sending ? 'Envoi...' : 'Envoyer' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { PaperAirplaneIcon } from '@heroicons/vue/24/outline'

const BASE = '/api/admin/v1'
function authHeaders() {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const items = ref([])
const meta = ref({ current_page: 1, last_page: 1 })
const loading = ref(true)
const search = ref('')
const typeFilter = ref('')
const knownTypes = ref([])
let searchTimer = null

async function fetchSent(page = 1) {
  loading.value = true
  try {
    const params = new URLSearchParams({ page, per_page: 30 })
    if (typeFilter.value) params.set('type', typeFilter.value)
    if (search.value) params.set('search', search.value)
    const res = await fetch(`${BASE}/notifications/sent?${params}`, { headers: authHeaders() })
    const j = await res.json()
    items.value = j.data || []
    meta.value = j.meta || meta.value
    // collect distinct types for the filter
    const set = new Set(knownTypes.value)
    items.value.forEach(n => set.add(n.type))
    knownTypes.value = [...set].sort()
  } finally { loading.value = false }
}
function debouncedFetch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchSent(1), 350)
}

const showBroadcast = ref(false)
const sending = ref(false)
const bcastMsg = ref('')
const bcastError = ref(false)
const bcast = ref({ audience: 'all', title: '', body: '', link: '' })
async function sendBroadcast() {
  if (!bcast.value.title || !bcast.value.body) { bcastMsg.value = 'Titre et message requis'; bcastError.value = true; return }
  sending.value = true
  bcastMsg.value = ''
  try {
    const res = await fetch(`${BASE}/notifications/broadcast`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify(bcast.value),
    })
    const j = await res.json()
    if (!res.ok) throw new Error(j.message || 'Erreur')
    bcastError.value = false
    bcastMsg.value = `Envoyé à ${j.sent} destinataire(s).`
    setTimeout(() => { showBroadcast.value = false; bcast.value = { audience: 'all', title: '', body: '', link: '' }; bcastMsg.value = ''; fetchSent(1) }, 1200)
  } catch (e) { bcastError.value = true; bcastMsg.value = e.message } finally { sending.value = false }
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('fr-FR', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => fetchSent(1))
</script>
