<template>
  <div class="space-y-6">
    
    <div class="flex items-start justify-between gap-3">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Demandes de Dépôt</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Validez les demandes de dépôt d'objets des utilisateurs</p>
      </div>
      <span class="text-xs sm:text-sm font-semibold px-3 py-1 rounded-full shrink-0" style="background:#fef9c3; color:#854d0e;">
        {{ pendingCount }} en attente
      </span>
    </div>

    
    <div class="flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
      <div class="relative flex-1">
        <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          v-model="search"
          @input="onSearchInput"
          type="text"
          placeholder="Rechercher..."
          class="w-full pl-9 pr-3 py-2 text-sm bg-white border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
        />
      </div>
      <div class="flex gap-1 bg-white border border-[#e5e7eb] rounded-xl p-1 overflow-x-auto">
        <button
          v-for="f in statusFilters"
          :key="f.value"
          @click="setFilter(f.value)"
          class="px-3 py-1.5 rounded-lg text-xs font-semibold transition whitespace-nowrap shrink-0"
          :class="statusFilter === f.value ? 'bg-[#006d35] text-white' : 'text-[#40617f] hover:bg-gray-50'"
        >
          {{ f.label }}
        </button>
      </div>
    </div>

    
    <div class="grid grid-cols-1 lg:grid-cols-5 gap-6">

      
      <div class="lg:col-span-2" :class="selectedRequest ? 'hidden lg:block' : ''">
        <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          <div class="px-5 py-4 border-b border-[#f1f5f9]">
            <h3 class="font-semibold text-[#001d32]">File d'attente</h3>
            <p class="text-xs text-gray-400 mt-0.5">{{ meta.total }} demandes au total</p>
          </div>

          
          <div v-if="loading" class="flex items-center justify-center py-12">
            <div class="w-6 h-6 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin"></div>
          </div>

          
          <div v-else-if="requests.length === 0" class="flex flex-col items-center justify-center py-12 text-gray-400">
            <InboxIcon class="w-10 h-10 mb-2 text-gray-300" />
            <p class="text-sm">Aucune demande</p>
          </div>

          
          <div v-else class="divide-y divide-[#f8fafc] max-h-[600px] overflow-y-auto">
            <div
              v-for="req in requests"
              :key="req.id"
              @click="selectRequest(req)"
              class="flex items-start gap-3 p-4 cursor-pointer transition-colors"
              :class="selectedRequest?.id === req.id ? 'bg-[#f0fdf4] border-l-2 border-[#006d35]' : 'hover:bg-[#f8fafc] border-l-2 border-transparent'"
            >
              <div class="w-10 h-10 rounded-xl overflow-hidden shrink-0 bg-gray-100 flex items-center justify-center">
                <img v-if="req.photo1" :src="req.photo1" class="w-full h-full object-cover" />
                <PhotoIcon v-else class="w-5 h-5 text-gray-400" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-semibold text-[#001d32] truncate">{{ req.title }}</p>
                <p class="text-xs text-gray-400 mt-0.5">
                  {{ req.user?.name || '—' }} · {{ req.category?.name || 'Sans catégorie' }}
                </p>
                <div class="flex items-center justify-between mt-1.5">
                  <span class="text-[10px] text-gray-400">{{ formatDate(req.created_at) }}</span>
                  <span class="text-xs font-semibold px-2 py-0.5 rounded-full" :class="statusBadge(req.status)">
                    {{ statusText(req.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          
          <div v-if="meta.last_page > 1" class="flex items-center justify-between px-5 py-3 border-t border-[#f1f5f9]">
            <button @click="changePage(meta.current_page - 1)" :disabled="meta.current_page <= 1" class="text-xs text-[#40617f] disabled:opacity-40 hover:text-[#006d35] transition">← Précédent</button>
            <span class="text-xs text-gray-400">{{ meta.current_page }} / {{ meta.last_page }}</span>
            <button @click="changePage(meta.current_page + 1)" :disabled="meta.current_page >= meta.last_page" class="text-xs text-[#40617f] disabled:opacity-40 hover:text-[#006d35] transition">Suivant →</button>
          </div>
        </div>
      </div>

      
      <div class="lg:col-span-3" ref="detailPanel">
        <div v-if="!selectedRequest" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm flex items-center justify-center h-96">
          <div class="text-center text-gray-400">
            <InboxIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p class="text-sm font-medium">Sélectionnez une demande</p>
            <p class="text-xs mt-1">Cliquez sur une demande pour voir les détails</p>
          </div>
        </div>

        <div v-else class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          
          <div class="lg:hidden px-4 py-2 border-b border-[#f1f5f9]">
            <button @click="selectedRequest = null" class="flex items-center gap-1.5 text-sm text-[#40617f] font-medium">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/></svg>
              Retour à la liste
            </button>
          </div>
          
          <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
            <div>
              <h3 class="font-semibold text-[#001d32]">{{ selectedRequest.title }}</h3>
              <p class="text-xs text-gray-400 mt-0.5">Demande #{{ selectedRequest.id }} · {{ formatDate(selectedRequest.created_at) }}</p>
            </div>
            <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="statusBadge(selectedRequest.status)">
              {{ statusText(selectedRequest.status) }}
            </span>
          </div>

          <div class="p-6 space-y-5">
            
            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-3">Informations</p>
              <div class="grid grid-cols-2 sm:grid-cols-2 gap-3">
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Utilisateur</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.user?.name || '—' }}</p>
                  <p class="text-xs text-gray-400">{{ selectedRequest.user?.email || '' }}</p>
                </div>
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Catégorie</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.category?.name || 'Sans catégorie' }}</p>
                </div>
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">État de l'objet</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ conditionText(selectedRequest.condition) }}</p>
                </div>
                <div v-if="selectedRequest.estimated_weight" class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Poids estimé</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.estimated_weight }} kg</p>
                </div>
              </div>
            </div>

            <!-- Photos -->
            <div v-if="selectedRequest.photo1 || selectedRequest.photo2 || selectedRequest.photo3">
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Photos</p>
              <div class="flex gap-2">
                <img
                  v-for="src in [selectedRequest.photo1, selectedRequest.photo2, selectedRequest.photo3].filter(Boolean)"
                  :key="src.slice(0,30)"
                  :src="src"
                  class="w-24 h-24 rounded-xl object-cover border border-[#e5e7eb] cursor-pointer hover:opacity-90 transition"
                  @click="lightboxSrc = src"
                />
              </div>
            </div>

            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Description</p>
              <p class="text-sm text-gray-600 bg-[#f8fafc] rounded-xl p-3 leading-relaxed">{{ selectedRequest.description }}</p>
            </div>

            
            <div v-if="selectedRequest.history">
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Historique de l'objet</p>
              <p class="text-sm text-gray-600 bg-[#f8fafc] rounded-xl p-3 leading-relaxed">{{ selectedRequest.history }}</p>
            </div>

            
            <div v-if="selectedRequest.carbon_savings">
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-3">Impact Environnemental</p>
              <div class="rounded-xl p-3 text-center" style="background:#f0fdf4;">
                <p class="text-lg font-bold" style="color:#006d35;">{{ selectedRequest.carbon_savings }} kg</p>
                <p class="text-[10px] text-gray-400 mt-0.5">CO₂ économisé</p>
              </div>
            </div>

            <div v-if="selectedRequest.qr_code" class="bg-[#f0fdf4] rounded-xl p-3 flex items-center gap-4">
              <img v-if="qrDataUrl" :src="qrDataUrl" class="w-16 h-16 rounded-lg border border-[#bbf7d0]" />
              <div>
                <p class="text-xs font-semibold text-[#006d35]">Code QR généré</p>
                <p class="text-sm font-mono text-[#001d32]">{{ selectedRequest.qr_code }}</p>
              </div>
            </div>

            
            <div v-if="selectedRequest.validation_note && selectedRequest.status !== 'pending'">
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Note de validation</p>
              <p class="text-sm text-gray-600 bg-[#f8fafc] rounded-xl p-3">{{ selectedRequest.validation_note }}</p>
            </div>

            
            <template v-if="selectedRequest.status === 'pending'">
              <div>
                <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Note de validation</p>
                <textarea
                  v-model="validationNote"
                  class="w-full px-3 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition resize-none"
                  rows="3"
                  placeholder="Ajouter une note (optionnel)..."
                ></textarea>
              </div>

              
              <p v-if="actionError" class="text-xs text-red-500">{{ actionError }}</p>

              <div class="flex items-center gap-3 pt-2">
                <button
                  @click="updateStatus('rejected')"
                  :disabled="actionLoading"
                  class="flex-1 py-3 rounded-xl text-sm font-semibold border-2 border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition flex items-center justify-center gap-2 disabled:opacity-50"
                >
                  <XCircleIcon class="w-4 h-4 text-red-400" />
                  Rejeter
                </button>
                <button
                  @click="updateStatus('approved')"
                  :disabled="actionLoading"
                  class="flex-1 py-3 rounded-xl text-sm font-semibold text-white transition hover:opacity-90 flex items-center justify-center gap-2 disabled:opacity-50"
                  style="background-color:#006d35;"
                >
                  <span v-if="actionLoading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                  <CheckCircleIcon v-else class="w-4 h-4" />
                  Approuver et Générer Code
                </button>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Photo lightbox -->
  <Teleport to="body">
    <div v-if="lightboxSrc" class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center" @click="lightboxSrc = null">
      <img :src="lightboxSrc" class="max-w-[90vw] max-h-[90vh] rounded-2xl object-contain shadow-2xl" @click.stop />
      <button @click="lightboxSrc = null" class="absolute top-4 right-4 p-2 rounded-full bg-white/20 hover:bg-white/40 transition text-white">
        <XMarkIcon class="w-6 h-6" />
      </button>
    </div>
  </Teleport>
</template>

<script setup>import { ref, computed, onMounted, watch, nextTick } from 'vue'
import {
  InboxIcon, PhotoIcon, CheckCircleIcon, XCircleIcon,
  MagnifyingGlassIcon, QrCodeIcon, XMarkIcon,
} from '@heroicons/vue/24/outline'
import QRCode from 'qrcode'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const BASE = import.meta.env.VITE_API_BASE_URL || '/api/admin/v1'

function authHeaders() {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const requests = ref([])
const meta = ref({ current_page: 1, last_page: 1, per_page: 20, total: 0 })
const loading = ref(false)
const selectedRequest = ref(null)
const validationNote = ref('')
const actionLoading = ref(false)
const actionError = ref('')
const search = ref('')
const statusFilter = ref('')
const detailPanel = ref(null)
const lightboxSrc = ref(null)
const qrDataUrl = ref(null)
let searchTimer = null

const statusFilters = [
  { value: '', label: 'Tous' },
  { value: 'pending', label: 'En attente' },
  { value: 'approved', label: 'Approuvés' },
  { value: 'rejected', label: 'Rejetés' },
]

const pendingCount = computed(() => {
  if (statusFilter.value === 'pending') return meta.value.total
  return requests.value.filter(r => r.status === 'pending').length
})

async function fetchRequests(page = 1) {
  loading.value = true
  try {
    const params = new URLSearchParams({ page, per_page: 20 })
    if (statusFilter.value) params.set('status', statusFilter.value)
    if (search.value) params.set('search', search.value)

    const res = await fetch(`${BASE}/deposits?${params}`, { headers: authHeaders() })
    if (!res.ok) throw new Error('Erreur réseau')
    const json = await res.json()
    requests.value = json.data || []
    meta.value = json.meta || meta.value
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function selectRequest(req) {
  validationNote.value = ''
  actionError.value = ''
  try {
    const res = await fetch(`${BASE}/deposits/${req.id}`, { headers: authHeaders() })
    if (!res.ok) throw new Error()
    selectedRequest.value = await res.json()
  } catch {
    selectedRequest.value = req
  }
}

async function updateStatus(status) {
  if (!selectedRequest.value) return
  actionLoading.value = true
  actionError.value = ''
  try {
    const res = await fetch(`${BASE}/deposits/${selectedRequest.value.id}/status`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify({ status, validation_note: validationNote.value || null }),
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({}))
      throw new Error(err.message || 'Erreur lors de la mise à jour')
    }
    const updated = await res.json()
    selectedRequest.value = updated
    const idx = requests.value.findIndex(r => r.id === updated.id)
    if (idx !== -1) requests.value[idx] = { ...requests.value[idx], status: updated.status }
    validationNote.value = ''
  } catch (e) {
    actionError.value = e.message
  } finally {
    actionLoading.value = false
  }
}

function setFilter(val) {
  statusFilter.value = val
  fetchRequests(1)
}

function onSearchInput() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(() => fetchRequests(1), 400)
}

function changePage(page) {
  if (page < 1 || page > meta.value.last_page) return
  fetchRequests(page)
}

function formatDate(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
}

function statusBadge(status) {
  if (status === 'pending')  return 'bg-[#fef9c3] text-[#854d0e]'
  if (status === 'approved') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'rejected') return 'bg-[#fee2e2] text-[#991b1b]'
  return 'bg-[#f1f5f9] text-[#475569]'
}

function statusText(status) {
  const map = { pending: 'En attente', approved: 'Approuvé', rejected: 'Rejeté' }
  return map[status] || status
}

function conditionText(condition) {
  const map = { new: 'Neuf', good: 'Bon état', fair: 'État correct', poor: 'Usé' }
  return map[condition] || condition || '—'
}

watch(() => selectedRequest.value, async (val) => {
  qrDataUrl.value = null
  if (val?.qr_code) {
    qrDataUrl.value = await QRCode.toDataURL(val.qr_code, { width: 160, margin: 1, errorCorrectionLevel: 'M' })
  }
  if (val && window.innerWidth < 1024) {
    await nextTick()
    detailPanel.value?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
})

onMounted(() => fetchRequests())
</script>
