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

    <div v-if="loading" class="flex justify-center py-16">
      <div class="w-6 h-6 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin" />
    </div>

    <div v-else-if="points.length === 0" class="bg-white rounded-2xl border border-[#f1f5f9] p-12 text-center text-gray-400">
      <MapPinIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
      <p class="font-medium">Aucun point de collecte</p>
      <p class="text-sm mt-1">Ajoutez votre premier point de collecte</p>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
      <div
        v-for="pt in points"
        :key="pt.id"
        class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-5 flex flex-col gap-3"
      >
        <div class="flex items-start justify-between gap-2">
          <div class="flex items-start gap-3">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0" :style="pt.is_active ? 'background:#f0fdf4' : 'background:#f1f5f9'">
              <MapPinIcon class="w-5 h-5" :style="pt.is_active ? 'color:#006d35' : 'color:#9ca3af'" />
            </div>
            <div>
              <p class="font-semibold text-[#001d32] text-sm">{{ pt.name }}</p>
              <p class="text-xs text-gray-400 mt-0.5">{{ pt.address }}, {{ pt.postal_code }} {{ pt.city }}</p>
            </div>
          </div>
          <span
            class="text-[10px] font-bold px-2 py-0.5 rounded-full shrink-0"
            :class="pt.is_active ? 'bg-[#dcfce7] text-[#166534]' : 'bg-[#f1f5f9] text-[#6b7280]'"
          >
            {{ pt.is_active ? 'Actif' : 'Inactif' }}
          </span>
        </div>

        <div v-if="pt.opening_hours" class="text-xs text-[#40617f] bg-[#f8fafc] rounded-lg px-3 py-2 flex items-start gap-1.5">
          <ClockIcon class="w-3.5 h-3.5 shrink-0 mt-0.5" />
          <span>{{ pt.opening_hours }}</span>
        </div>
        <div v-if="pt.phone" class="text-xs text-[#40617f] flex items-center gap-1.5">
          <PhoneIcon class="w-3.5 h-3.5 shrink-0" />
          {{ pt.phone }}
        </div>

        <div class="flex gap-2 pt-1 border-t border-[#f1f5f9]">
          <button @click="openEdit(pt)" class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-[#40617f] bg-[#f8fafc] hover:bg-[#edf4ff] transition">
            Modifier
          </button>
          <button @click="confirmDelete(pt)" class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-red-500 bg-red-50 hover:bg-red-100 transition">
            Supprimer
          </button>
        </div>
      </div>
    </div>

    <!-- Modal create/edit -->
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

    <!-- Delete confirm -->
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { PlusIcon, MapPinIcon, ClockIcon, PhoneIcon, XMarkIcon } from '@heroicons/vue/24/outline'

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
