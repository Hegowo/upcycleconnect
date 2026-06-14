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

    <Teleport to="body">
      <div v-if="showBroadcast" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="showBroadcast = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-4">Diffuser une notification</h3>
          <div class="space-y-3">
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Audience</label>
              <select v-model="bcast.audience" @change="onAudienceChange" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                <option value="all">Tous les utilisateurs</option>
                <option value="particuliers">Particuliers</option>
                <option value="pros">Professionnels / artisans</option>
                <option value="user">Un utilisateur précis</option>
              </select>
            </div>

            <div v-if="bcast.audience === 'user'">
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Destinataire</label>
              <div v-if="selectedUser" class="flex items-center justify-between gap-2 px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl">
                <div class="min-w-0">
                  <p class="text-sm font-medium text-[#001d32] truncate">{{ selectedUser.name }}<span v-if="selectedUser.company" class="text-gray-400"> · {{ selectedUser.company }}</span></p>
                  <p class="text-xs text-gray-400 truncate">{{ selectedUser.email }}</p>
                </div>
                <button @click="clearUser" class="text-xs font-semibold text-[#006d35] shrink-0">Changer</button>
              </div>
              <div v-else class="relative">
                <input v-model="userQuery" @input="searchUsers" type="text" placeholder="Email, nom, prénom ou société…" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
                <div v-if="userResults.length" class="absolute z-10 left-0 right-0 mt-1 bg-white border border-[#e5e7eb] rounded-xl shadow-lg max-h-52 overflow-y-auto">
                  <button v-for="u in userResults" :key="u.id" @click="pickUser(u)" class="w-full text-left px-3 py-2 hover:bg-[#f0fdf4] transition">
                    <p class="text-sm text-[#001d32] truncate">{{ u.name }}<span v-if="u.company" class="text-gray-400"> · {{ u.company }}</span></p>
                    <p class="text-xs text-gray-400 truncate">{{ u.email }}</p>
                  </button>
                </div>
                <p v-else-if="userQuery.length >= 2 && !userSearching" class="text-xs text-gray-400 mt-1">Aucun utilisateur trouvé.</p>
              </div>
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
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Destination au clic</label>
              <select v-model="linkMode" @change="onLinkModeChange" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                <option value="none">Aucune (pas de lien)</option>
                <option value="/">Accueil</option>
                <option value="/prestations">Liste des prestations</option>
                <option value="prestation">Une prestation précise…</option>
                <option value="/evenements">Liste des événements</option>
                <option value="event">Un événement précis…</option>
                <option value="/communaute">Communauté / Forum</option>
                <option value="/depot">Déposer un objet</option>
                <option value="/profil">Profil utilisateur</option>
                <option value="custom">URL personnalisée…</option>
              </select>
            </div>

            <div v-if="linkMode === 'prestation' || linkMode === 'event'">
              <div v-if="selectedItem" class="flex items-center justify-between gap-2 px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl">
                <p class="text-sm text-[#001d32] truncate">{{ selectedItem.label }}</p>
                <button @click="selectedItem = null" class="text-xs font-semibold text-[#006d35] shrink-0">Changer</button>
              </div>
              <div v-else class="relative">
                <input v-model="itemQuery" @input="searchItems" type="text" :placeholder="linkMode === 'prestation' ? 'Rechercher une prestation…' : 'Rechercher un événement…'" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
                <div v-if="itemResults.length" class="absolute z-10 left-0 right-0 mt-1 bg-white border border-[#e5e7eb] rounded-xl shadow-lg max-h-52 overflow-y-auto">
                  <button v-for="(it, i) in itemResults" :key="i" @click="pickItem(it)" class="block w-full text-left px-3 py-2 text-sm text-[#001d32] hover:bg-[#f0fdf4] transition truncate">{{ it.label }}</button>
                </div>
                <p v-else-if="itemQuery.length >= 1 && !itemSearching" class="text-xs text-gray-400 mt-1">Aucun résultat.</p>
              </div>
            </div>

            <div v-if="linkMode === 'custom'">
              <input v-model="customLink" type="text" placeholder="/prestations/12" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
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
const bcast = ref({ audience: 'all', title: '', body: '' })

const selectedUser = ref(null)
const userQuery = ref('')
const userResults = ref([])
const userSearching = ref(false)
let userTimer = null
function searchUsers() {
  clearTimeout(userTimer)
  const q = userQuery.value.trim()
  if (q.length < 2) { userResults.value = []; return }
  userSearching.value = true
  userTimer = setTimeout(async () => {
    try {
      const res = await fetch(`${BASE}/notifications/recipients?q=${encodeURIComponent(q)}`, { headers: authHeaders() })
      const j = await res.json()
      userResults.value = j.data || []
    } catch { userResults.value = [] } finally { userSearching.value = false }
  }, 300)
}
function pickUser(u) { selectedUser.value = u; userResults.value = []; userQuery.value = '' }
function clearUser() { selectedUser.value = null }
function onAudienceChange() { if (bcast.value.audience !== 'user') clearUser() }

const linkMode = ref('none')
const customLink = ref('')
const selectedItem = ref(null)
const itemQuery = ref('')
const itemResults = ref([])
const itemSearching = ref(false)
let itemTimer = null
function onLinkModeChange() { selectedItem.value = null; itemResults.value = []; itemQuery.value = '' }
function searchItems() {
  clearTimeout(itemTimer)
  const q = itemQuery.value.trim()
  itemSearching.value = true
  itemTimer = setTimeout(async () => {
    try {
      const res = await fetch(`${BASE}/notifications/link-targets?type=${linkMode.value}&q=${encodeURIComponent(q)}`, { headers: authHeaders() })
      const j = await res.json()
      itemResults.value = j.data || []
    } catch { itemResults.value = [] } finally { itemSearching.value = false }
  }, 300)
}
function pickItem(it) { selectedItem.value = it; itemResults.value = []; itemQuery.value = '' }
function resolveLink() {
  if (linkMode.value === 'none') return ''
  if (linkMode.value === 'custom') return customLink.value.trim()
  if (linkMode.value === 'prestation' || linkMode.value === 'event') return selectedItem.value?.path || ''
  return linkMode.value
}

function resetBroadcast() {
  bcast.value = { audience: 'all', title: '', body: '' }
  selectedUser.value = null; userQuery.value = ''; userResults.value = []
  linkMode.value = 'none'; customLink.value = ''; selectedItem.value = null; itemQuery.value = ''; itemResults.value = []
  bcastMsg.value = ''
}

async function sendBroadcast() {
  if (!bcast.value.title || !bcast.value.body) { bcastMsg.value = 'Titre et message requis'; bcastError.value = true; return }
  if (bcast.value.audience === 'user' && !selectedUser.value) { bcastMsg.value = 'Choisis un destinataire'; bcastError.value = true; return }
  if ((linkMode.value === 'prestation' || linkMode.value === 'event') && !selectedItem.value) { bcastMsg.value = 'Choisis la cible du lien'; bcastError.value = true; return }
  sending.value = true
  bcastMsg.value = ''
  try {
    const payload = {
      audience: bcast.value.audience,
      title: bcast.value.title,
      body: bcast.value.body,
      link: resolveLink(),
      user_id: bcast.value.audience === 'user' ? selectedUser.value?.id : null,
    }
    const res = await fetch(`${BASE}/notifications/broadcast`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify(payload),
    })
    const j = await res.json()
    if (!res.ok) throw new Error(j.message || 'Erreur')
    bcastError.value = false
    bcastMsg.value = `Envoyé à ${j.sent} destinataire(s).`
    setTimeout(() => { showBroadcast.value = false; resetBroadcast(); fetchSent(1) }, 1200)
  } catch (e) { bcastError.value = true; bcastMsg.value = e.message } finally { sending.value = false }
}

function formatDate(iso) {
  return new Date(iso).toLocaleString('fr-FR', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
}

onMounted(() => fetchSent(1))
</script>
