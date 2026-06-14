<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Points de collecte</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Gérez les lieux où les utilisateurs peuvent déposer leurs objets</p>
      </div>
      <button
        @click="openCreate"
        class="flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90"
        style="background:#006d35;"
      >
        <PlusIcon class="w-4 h-4" />
        Ajouter
      </button>
    </div>

    <div class="card p-4 flex flex-wrap items-center gap-3">
      <div class="relative flex-1 min-w-[200px]">
        <MagnifyingGlassIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#94a3b8]" />
        <input v-model="search" type="search" placeholder="Rechercher par nom, ville, code postal…"
          class="w-full pl-9 pr-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
      </div>
      <select v-model="deptFilter" class="text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
        <option value="">Tous les départements</option>
        <option v-for="d in departments" :key="d" :value="d">{{ d }} — {{ deptName(d) }}</option>
      </select>
      <div class="flex gap-1.5">
        <button v-for="f in statusFilters" :key="f.value" @click="statusFilter = f.value; page = 1"
          :class="['px-3 py-1.5 rounded-full text-xs font-semibold border transition',
            statusFilter === f.value ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]']">
          {{ f.label }}
        </button>
      </div>
    </div>

    <div class="flex items-center gap-4 text-xs text-[#64748b] px-1">
      <span><strong class="text-[#001d32]">{{ filtered.length }}</strong> résultat(s)</span>
      <span><strong class="text-[#16a34a]">{{ activeCount }}</strong> actifs</span>
      <span><strong class="text-[#6b7280]">{{ points.length - activeCount }}</strong> inactifs</span>
    </div>

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-6 h-6 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="!filtered.length" class="bg-white rounded-2xl border border-[#f1f5f9] p-12 text-center text-gray-400">
      <MapPinIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
      <p class="font-medium">Aucun point de collecte</p>
      <p class="text-sm mt-1">Aucun résultat pour ces filtres</p>
    </div>

    <div v-else class="card overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-[#f8fafc] text-[#64748b] text-xs uppercase">
          <tr>
            <th class="text-left px-4 py-3">Nom</th>
            <th class="text-left px-4 py-3 hidden md:table-cell">Adresse</th>
            <th class="text-left px-4 py-3">Ville</th>
            <th class="text-left px-4 py-3 hidden lg:table-cell">Horaires</th>
            <th class="text-left px-4 py-3">Statut</th>
            <th class="text-right px-4 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="pt in paginated" :key="pt.id" class="border-t border-[#f1f5f9] hover:bg-[#f8fafc]">
            <td class="px-4 py-2.5">
              <div class="flex items-center gap-2.5">
                <MapPinIcon class="w-4 h-4 shrink-0" :style="pt.is_active ? 'color:#006d35' : 'color:#cbd5e1'" />
                <span class="font-medium text-[#001d32]">{{ pt.name }}</span>
              </div>
            </td>
            <td class="px-4 py-2.5 text-[#40617f] hidden md:table-cell">{{ pt.address }}</td>
            <td class="px-4 py-2.5 text-[#40617f]"><span class="font-mono text-xs text-[#94a3b8]">{{ pt.postal_code }}</span> {{ pt.city }}</td>
            <td class="px-4 py-2.5 text-[#64748b] text-xs hidden lg:table-cell">{{ pt.opening_hours || '—' }}</td>
            <td class="px-4 py-2.5">
              <div class="flex flex-col gap-1 items-start">
                <span v-if="pt.out_of_service" class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-[#fef2f2] text-[#b91c1c]">
                  Hors service
                </span>
                <button v-else @click="toggleActive(pt)"
                  class="text-[10px] font-bold px-2 py-0.5 rounded-full transition"
                  :class="pt.is_active ? 'bg-[#dcfce7] text-[#166534] hover:bg-[#bbf7d0]' : 'bg-[#f1f5f9] text-[#6b7280] hover:bg-[#e5e7eb]'">
                  {{ pt.is_active ? 'Actif' : 'Inactif' }}
                </button>
                <span v-if="pt.out_of_service && pt.out_of_service_until" class="text-[10px] text-[#94a3b8]">
                  jusqu'au {{ formatDate(pt.out_of_service_until) }}
                </span>
                <span v-else-if="pt.out_of_service" class="text-[10px] text-[#94a3b8]">durée indéterminée</span>
              </div>
            </td>
            <td class="px-4 py-2.5 text-right whitespace-nowrap">
              <button v-if="!pt.out_of_service" @click="openOOS(pt)" class="p-1.5 rounded-lg text-gray-400 hover:text-amber-600 hover:bg-amber-50 transition" title="Déclarer hors service">
                <NoSymbolIcon class="w-4 h-4" />
              </button>
              <button v-else @click="restore(pt)" class="p-1.5 rounded-lg text-gray-400 hover:text-[#006d35] hover:bg-[#f0fdf4] transition" title="Remettre en service">
                <ArrowPathIcon class="w-4 h-4" />
              </button>
              <button @click="openEdit(pt)" class="p-1.5 rounded-lg text-gray-400 hover:text-[#006d35] hover:bg-[#f0fdf4] transition" title="Modifier">
                <PencilSquareIcon class="w-4 h-4" />
              </button>
              <button @click="confirmDelete(pt)" class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition" title="Supprimer">
                <TrashIcon class="w-4 h-4" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-if="lastPage > 1" class="flex items-center justify-between px-4 py-3 border-t border-[#f1f5f9]">
        <button @click="page--" :disabled="page <= 1" class="text-sm text-[#40617f] disabled:opacity-40 hover:text-[#006d35]">← Précédent</button>
        <span class="text-xs text-[#94a3b8]">Page {{ page }} / {{ lastPage }} · {{ filtered.length }} points</span>
        <button @click="page++" :disabled="page >= lastPage" class="text-sm text-[#40617f] disabled:opacity-40 hover:text-[#006d35]">Suivant →</button>
      </div>
    </div>

    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-[#001d32]/60" @click.self="closeModal">
          <div class="bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[90vh] overflow-y-auto" @click.stop>
            <div class="sticky top-0 bg-white border-b border-[#f1f5f9] px-6 py-4 flex items-center justify-between rounded-t-2xl">
              <h3 class="font-bold text-[#001d32]">{{ editTarget ? 'Modifier le point' : 'Nouveau point de collecte' }}</h3>
              <button @click="closeModal" class="p-1 rounded-lg hover:bg-gray-100 transition">
                <XMarkIcon class="w-5 h-5 text-gray-400" />
              </button>
            </div>

            <form @submit.prevent="submitForm" class="p-6 space-y-4">
              <div v-if="formError" class="text-sm text-red-600 bg-red-50 rounded-xl px-3 py-2">{{ formError }}</div>

              <div>
                <label class="block text-xs font-semibold text-[#001d32] mb-1.5">Nom du point <span class="text-red-500">*</span></label>
                <input v-model="form.name" type="text" required placeholder="Ex: Recyclerie du 11e" class="w-full px-3 py-2.5 text-sm border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>

              <div>
                <label class="block text-xs font-semibold text-[#001d32] mb-1.5">Adresse complète <span class="text-red-500">*</span></label>
                <div class="relative">
                  <input
                    v-model="addressQuery"
                    @input="onAddressInput"
                    type="text"
                    placeholder="Rechercher une adresse..."
                    autocomplete="off"
                    class="w-full px-3 py-2.5 text-sm border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30"
                  />
                  <ul v-if="addressSuggestions.length" class="absolute z-10 mt-1 w-full bg-white border border-[#e5e7eb] rounded-xl shadow-lg overflow-hidden">
                    <li
                      v-for="s in addressSuggestions"
                      :key="s.properties.id"
                      @mousedown.prevent="selectAddress(s)"
                      class="px-3 py-2.5 text-sm cursor-pointer hover:bg-[#f0fdf4] flex items-start gap-2"
                    >
                      <MapPinIcon class="w-4 h-4 text-[#006d35] shrink-0 mt-0.5" />
                      <span>{{ s.properties.label }}</span>
                    </li>
                  </ul>
                </div>
                <div v-if="form.address" class="mt-2 grid grid-cols-2 gap-2">
                  <input v-model="form.address" type="text" placeholder="Adresse" class="px-3 py-2 text-xs border border-[#e5e7eb] rounded-lg focus:outline-none" />
                  <input v-model="form.postal_code" type="text" placeholder="Code postal" class="px-3 py-2 text-xs border border-[#e5e7eb] rounded-lg focus:outline-none" />
                  <input v-model="form.city" type="text" placeholder="Ville" class="px-3 py-2 text-xs border border-[#e5e7eb] rounded-lg col-span-2 focus:outline-none" />
                  <div class="col-span-2 text-xs text-gray-400 flex gap-4">
                    <span>Lat: {{ form.latitude.toFixed(6) }}</span>
                    <span>Lng: {{ form.longitude.toFixed(6) }}</span>
                  </div>
                </div>
              </div>

              <div>
                <label class="block text-xs font-semibold text-[#001d32] mb-1.5">Téléphone</label>
                <input v-model="form.phone" type="text" placeholder="01 23 45 67 89" class="w-full px-3 py-2.5 text-sm border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>

              <div>
                <label class="block text-xs font-semibold text-[#001d32] mb-1.5">Horaires d'ouverture</label>
                <textarea v-model="form.opening_hours" rows="2" placeholder="Ex: Lun-Ven 9h-18h, Sam 10h-16h" class="w-full px-3 py-2.5 text-sm border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" />
              </div>

              <div class="flex items-center gap-3">
                <button
                  type="button"
                  @click="form.is_active = !form.is_active"
                  class="relative w-11 h-6 rounded-full transition-colors duration-200 focus:outline-none"
                  :class="form.is_active ? 'bg-[#006d35]' : 'bg-gray-300'"
                >
                  <span class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow transition-transform duration-200" :class="form.is_active ? 'translate-x-5' : 'translate-x-0'" />
                </button>
                <span class="text-sm font-medium text-[#001d32]">Point actif</span>
              </div>

              <div class="flex gap-3 pt-2 border-t border-[#f1f5f9]">
                <button type="button" @click="closeModal" class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] bg-[#f8fafc] hover:bg-[#edf4ff] transition">Annuler</button>
                <button type="submit" :disabled="submitting || !form.address" class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-white transition hover:opacity-90 disabled:opacity-50" style="background:#006d35;">
                  <span v-if="submitting" class="inline-block w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                  <span v-else>{{ editTarget ? 'Enregistrer' : 'Créer' }}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="deleteTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-[#001d32]/60" @click.self="deleteTarget = null">
          <div class="bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6 space-y-4">
            <h3 class="font-bold text-[#001d32]">Supprimer ce point ?</h3>
            <p class="text-sm text-[#40617f]">Le point <strong>{{ deleteTarget.name }}</strong> sera définitivement supprimé.</p>
            <div class="flex gap-3">
              <button @click="deleteTarget = null" class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] bg-[#f8fafc] hover:bg-[#edf4ff] transition">Annuler</button>
              <button @click="doDelete" :disabled="submitting" class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-white bg-red-500 hover:bg-red-600 transition disabled:opacity-50">
                {{ submitting ? '...' : 'Supprimer' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <Transition name="modal-fade">
        <div v-if="oosTarget" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-[#001d32]/60" @click.self="oosTarget = null">
          <div class="bg-white rounded-2xl shadow-2xl w-full max-w-sm p-6 space-y-4">
            <h3 class="font-bold text-[#001d32]">Déclarer hors service</h3>
            <p class="text-sm text-[#40617f]">Le box <strong>{{ oosTarget.name }}</strong> ne pourra plus être choisi par les utilisateurs pour un dépôt.</p>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Jusqu'au (facultatif)</label>
              <input v-model="oosUntil" type="date" :min="todayStr" class="w-full px-3 py-2.5 text-sm border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              <p class="text-xs text-[#94a3b8] mt-1.5">Laissez vide pour une durée indéterminée. Le box redeviendra disponible automatiquement à la date choisie.</p>
            </div>
            <div class="flex gap-3">
              <button @click="oosTarget = null" class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] bg-[#f8fafc] hover:bg-[#edf4ff] transition">Annuler</button>
              <button @click="doOOS" :disabled="submitting" class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-white transition disabled:opacity-50" style="background-color:#d97706;">
                {{ submitting ? '...' : 'Mettre hors service' }}
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { PlusIcon, MapPinIcon, ClockIcon, PhoneIcon, XMarkIcon, MagnifyingGlassIcon, PencilSquareIcon, TrashIcon, NoSymbolIcon, ArrowPathIcon } from '@heroicons/vue/24/outline'

const BASE = import.meta.env.VITE_API_BASE_URL || '/api/admin/v1'

function authHeaders(json = false) {
  const token = localStorage.getItem('admin_token')
  const h = {}
  if (token) h['Authorization'] = `Bearer ${token}`
  if (json) h['Content-Type'] = 'application/json'
  return h
}

const points     = ref([])
const loading    = ref(false)
const showModal  = ref(false)
const editTarget = ref(null)
const deleteTarget = ref(null)
const submitting = ref(false)
const formError  = ref('')

const addressQuery       = ref('')
const addressSuggestions = ref([])
let addressTimer = null

const form = ref(emptyForm())

const search = ref('')
const statusFilter = ref('')
const deptFilter = ref('')
const page = ref(1)
const perPage = 20

const statusFilters = [
  { value: '', label: 'Tous' },
  { value: 'active', label: 'Actifs' },
  { value: 'inactive', label: 'Inactifs' },
]

const DEPT_NAMES = {
  '01':'Ain','02':'Aisne','03':'Allier','04':'Alpes-de-Haute-Provence','05':'Hautes-Alpes',
  '06':'Alpes-Maritimes','07':'Ardèche','08':'Ardennes','09':'Ariège','10':'Aube','11':'Aude',
  '12':'Aveyron','13':'Bouches-du-Rhône','14':'Calvados','15':'Cantal','16':'Charente',
  '17':'Charente-Maritime','18':'Cher','19':'Corrèze','21':"Côte-d'Or",'22':"Côtes-d'Armor",
  '23':'Creuse','24':'Dordogne','25':'Doubs','26':'Drôme','27':'Eure','28':'Eure-et-Loir',
  '29':'Finistère','30':'Gard','31':'Haute-Garonne','32':'Gers','33':'Gironde','34':'Hérault',
  '35':'Ille-et-Vilaine','36':'Indre','37':'Indre-et-Loire','38':'Isère','39':'Jura','40':'Landes',
  '41':'Loir-et-Cher','42':'Loire','43':'Haute-Loire','44':'Loire-Atlantique','45':'Loiret',
  '46':'Lot','47':'Lot-et-Garonne','48':'Lozère','49':'Maine-et-Loire','50':'Manche','51':'Marne',
  '52':'Haute-Marne','53':'Mayenne','54':'Meurthe-et-Moselle','55':'Meuse','56':'Morbihan',
  '57':'Moselle','58':'Nièvre','59':'Nord','60':'Oise','61':'Orne','62':'Pas-de-Calais',
  '63':'Puy-de-Dôme','64':'Pyrénées-Atlantiques','65':'Hautes-Pyrénées','66':'Pyrénées-Orientales',
  '67':'Bas-Rhin','68':'Haut-Rhin','69':'Rhône','70':'Haute-Saône','71':'Saône-et-Loire',
  '72':'Sarthe','73':'Savoie','74':'Haute-Savoie','75':'Paris','76':'Seine-Maritime',
  '77':'Seine-et-Marne','78':'Yvelines','79':'Deux-Sèvres','80':'Somme','81':'Tarn',
  '82':'Tarn-et-Garonne','83':'Var','84':'Vaucluse','85':'Vendée','86':'Vienne','87':'Haute-Vienne',
  '88':'Vosges','89':'Yonne','90':'Territoire de Belfort','91':'Essonne','92':'Hauts-de-Seine',
  '93':'Seine-Saint-Denis','94':'Val-de-Marne','95':"Val-d'Oise",
}
function deptName(d) { return DEPT_NAMES[d] || 'Autre' }

const departments = computed(() => {
  const set = new Set(points.value.map(p => (p.postal_code || '').slice(0, 2)).filter(Boolean))
  return [...set].sort()
})

const activeCount = computed(() => points.value.filter(p => p.is_active).length)

const filtered = computed(() => {
  const q = search.value.trim().toLowerCase()
  return points.value.filter(p => {
    if (statusFilter.value === 'active' && !p.is_active) return false
    if (statusFilter.value === 'inactive' && p.is_active) return false
    if (deptFilter.value && !(p.postal_code || '').startsWith(deptFilter.value)) return false
    if (q) {
      const hay = `${p.name} ${p.city} ${p.postal_code} ${p.address}`.toLowerCase()
      if (!hay.includes(q)) return false
    }
    return true
  })
})

const lastPage = computed(() => Math.max(1, Math.ceil(filtered.value.length / perPage)))
const paginated = computed(() => filtered.value.slice((page.value - 1) * perPage, page.value * perPage))

watch([search, deptFilter, statusFilter], () => { page.value = 1 })

async function toggleActive(pt) {
  const updated = { ...pt, is_active: !pt.is_active }
  await fetch(`${BASE}/collection-points/${pt.id}`, { method: 'PUT', headers: authHeaders(true), body: JSON.stringify(updated) })
  pt.is_active = !pt.is_active
}

const oosTarget = ref(null)
const oosUntil = ref('')
const todayStr = new Date().toISOString().slice(0, 10)
function openOOS(pt) { oosTarget.value = pt; oosUntil.value = ''; }
async function doOOS() {
  submitting.value = true
  try {
    const body = { out_of_service: true, until: oosUntil.value || null }
    const res = await fetch(`${BASE}/collection-points/${oosTarget.value.id}/out-of-service`, { method: 'PUT', headers: authHeaders(true), body: JSON.stringify(body) })
    if (res.ok) { const upd = await res.json(); Object.assign(oosTarget.value, upd) }
    oosTarget.value = null
  } finally { submitting.value = false }
}
async function restore(pt) {
  const res = await fetch(`${BASE}/collection-points/${pt.id}/out-of-service`, { method: 'PUT', headers: authHeaders(true), body: JSON.stringify({ out_of_service: false }) })
  if (res.ok) { const upd = await res.json(); Object.assign(pt, upd) }
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}

function emptyForm() {
  return { name: '', address: '', city: '', postal_code: '', latitude: 0, longitude: 0, phone: '', opening_hours: '', is_active: true }
}

async function fetchPoints() {
  loading.value = true
  try {
    const res = await fetch(`${BASE}/collection-points`, { headers: authHeaders() })
    points.value = res.ok ? await res.json() : []
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editTarget.value = null
  form.value = emptyForm()
  addressQuery.value = ''
  addressSuggestions.value = []
  formError.value = ''
  showModal.value = true
}

function openEdit(pt) {
  editTarget.value = pt
  form.value = { name: pt.name, address: pt.address, city: pt.city, postal_code: pt.postal_code, latitude: pt.latitude, longitude: pt.longitude, phone: pt.phone || '', opening_hours: pt.opening_hours || '', is_active: pt.is_active }
  addressQuery.value = `${pt.address}, ${pt.postal_code} ${pt.city}`
  addressSuggestions.value = []
  formError.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editTarget.value = null
}

function onAddressInput() {
  clearTimeout(addressTimer)
  const q = addressQuery.value.trim()
  if (q.length < 3) { addressSuggestions.value = []; return }
  addressTimer = setTimeout(async () => {
    const res = await fetch(`https://api-adresse.data.gouv.fr/search/?q=${encodeURIComponent(q)}&limit=5&autocomplete=1`)
    const data = await res.json()
    addressSuggestions.value = data.features || []
  }, 300)
}

function selectAddress(s) {
  const p = s.properties
  form.value.address     = [p.housenumber, p.street || p.name].filter(Boolean).join(' ')
  form.value.city        = p.city || ''
  form.value.postal_code = p.postcode || ''
  form.value.latitude    = s.geometry.coordinates[1]
  form.value.longitude   = s.geometry.coordinates[0]
  addressQuery.value     = p.label
  addressSuggestions.value = []
}

async function submitForm() {
  if (!form.value.address || !form.value.latitude) {
    formError.value = 'Veuillez sélectionner une adresse dans la liste.'
    return
  }
  formError.value = ''
  submitting.value = true
  try {
    const url = editTarget.value
      ? `${BASE}/collection-points/${editTarget.value.id}`
      : `${BASE}/collection-points`
    const method = editTarget.value ? 'PUT' : 'POST'
    const res = await fetch(url, { method, headers: authHeaders(true), body: JSON.stringify(form.value) })
    if (!res.ok) {
      const e = await res.json().catch(() => ({}))
      throw new Error(e.message || 'Erreur')
    }
    await fetchPoints()
    closeModal()
  } catch (e) {
    formError.value = e.message
  } finally {
    submitting.value = false
  }
}

function confirmDelete(pt) {
  deleteTarget.value = pt
}

async function doDelete() {
  if (!deleteTarget.value) return
  submitting.value = true
  try {
    await fetch(`${BASE}/collection-points/${deleteTarget.value.id}`, { method: 'DELETE', headers: authHeaders() })
    deleteTarget.value = null
    await fetchPoints()
  } finally {
    submitting.value = false
  }
}

onMounted(fetchPoints)
</script>

<style scoped>
.modal-fade-enter-active { transition: opacity 0.2s ease; }
.modal-fade-leave-active { transition: opacity 0.15s ease; }
.modal-fade-enter-from,
.modal-fade-leave-to { opacity: 0; }
</style>
