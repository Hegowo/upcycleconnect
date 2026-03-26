<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-[#001d32]">Gestion Catalogue</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Prestations & Services — Gérez l'offre de la plateforme</p>
      </div>
    </div>

    <!-- 3-column layout: sidebar filters + grid -->
    <div class="flex gap-6 items-start">

      <!-- Sidebar filters -->
      <div class="w-64 shrink-0">
        <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-5 sticky top-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="font-semibold text-[#001d32] text-sm">Filtres Catalogue</h3>
            <button @click="resetFilters" class="text-xs font-medium hover:underline" style="color:#006d35;">
              Réinitialiser
            </button>
          </div>

          <!-- Search -->
          <div class="mb-4">
            <div class="relative">
              <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                v-model="filters.search"
                @input="debouncedFetch"
                type="text"
                class="w-full pl-9 pr-3 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
                placeholder="Rechercher..."
              />
            </div>
          </div>

          <!-- Price type (content type) -->
          <div class="mb-4">
            <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Type de Contenu</p>
            <div class="space-y-2">
              <label v-for="type in priceTypes" :key="type.value" class="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  :value="type.value"
                  v-model="selectedTypes"
                  @change="fetchPrestations"
                  class="w-4 h-4 rounded border-gray-300 accent-green-700"
                />
                <span class="text-sm text-gray-600">{{ type.label }}</span>
              </label>
            </div>
          </div>

          <!-- Status -->
          <div class="mb-4">
            <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Statut</p>
            <div class="space-y-2">
              <label v-for="s in statusOptions" :key="s.value" class="flex items-center gap-2 cursor-pointer">
                <input
                  type="radio"
                  :value="s.value"
                  v-model="filters.status"
                  @change="fetchPrestations"
                  class="w-4 h-4 border-gray-300 accent-green-700"
                />
                <span class="text-sm text-gray-600">{{ s.label }}</span>
              </label>
            </div>
          </div>

          <!-- Category -->
          <div>
            <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Catégorie</p>
            <select v-model="filters.category_id" @change="fetchPrestations" class="w-full text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
              <option value="">Toutes</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Grid -->
      <div class="flex-1 min-w-0">
        <!-- Loading -->
        <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div v-for="n in 4" :key="n" class="bg-white rounded-2xl border border-[#f1f5f9] h-52 animate-pulse"></div>
        </div>

        <!-- Empty -->
        <div v-else-if="!prestations.length" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm py-16 text-center">
          <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
            <MagnifyingGlassIcon class="w-6 h-6 text-gray-400" />
          </div>
          <p class="text-gray-500 font-medium">Aucune prestation trouvée</p>
          <p class="text-gray-400 text-sm mt-1">Modifiez vos critères de filtrage</p>
        </div>

        <!-- Cards grid -->
        <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <div
            v-for="p in prestations"
            :key="p.id"
            class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden hover:shadow-md transition-shadow"
          >
            <div class="p-5">
              <!-- Badges top -->
              <div class="flex items-center gap-2 mb-3 flex-wrap">
                <span class="text-xs font-semibold px-2 py-0.5 rounded-full" :class="prestationStatusBadge(p.status)">
                  {{ prestationStatusText(p.status) }}
                </span>
                <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-[#f1f5f9] text-[#475569]">
                  {{ priceTypeLabel(p.price_type) }}
                </span>
                <span v-if="p.category?.name" class="text-xs font-medium px-2 py-0.5 rounded-full bg-[#dbeafe] text-blue-700">
                  {{ p.category.name }}
                </span>
              </div>

              <!-- Title -->
              <h4 class="text-sm font-bold text-[#001d32] mb-2 line-clamp-2">{{ p.title }}</h4>

              <!-- Details -->
              <div class="space-y-1.5 mb-4">
                <div class="flex items-center gap-2 text-xs text-gray-500">
                  <UserIcon class="w-3.5 h-3.5 text-gray-400 shrink-0" />
                  <span>{{ p.provider?.first_name }} {{ p.provider?.last_name }}</span>
                </div>
                <div class="flex items-center gap-2 text-xs text-gray-500">
                  <CurrencyEuroIcon class="w-3.5 h-3.5 text-gray-400 shrink-0" />
                  <span class="font-semibold text-[#001d32]">{{ p.price ? `${p.price} €` : 'Sur devis' }}</span>
                  <span class="text-gray-400">({{ priceTypeLabel(p.price_type) }})</span>
                </div>
              </div>

              <!-- Actions -->
              <div class="flex items-center gap-2 pt-3 border-t border-[#f8fafc]">
                <button
                  v-if="p.status === 'draft'"
                  @click="changeStatus(p, 'published')"
                  class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
                  style="background:#dcfce7; color:#166534;"
                >
                  Publier
                </button>
                <button
                  v-if="p.status === 'published'"
                  @click="changeStatus(p, 'suspended')"
                  class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
                  style="background:#fff7ed; color:#c2410c;"
                >
                  Suspendre
                </button>
                <RouterLink
                  :to="`/admin/prestations/${p.id}/edit`"
                  class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
                  style="background:#f1f5f9; color:#475569;"
                >
                  Modifier
                </RouterLink>
                <button @click="openDelete(p)" class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="prestations.length" class="mt-4 flex justify-end">
          <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchPrestations" />
        </div>
      </div>
    </div>

    <!-- Footer stats bar -->
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-4">
      <div class="grid grid-cols-2 lg:grid-cols-4 gap-6">
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">Prévision Revenus</p>
          <p class="text-lg font-bold text-[#001d32]">12 400 €</p>
        </div>
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">Apprenants Actifs</p>
          <p class="text-lg font-bold text-[#001d32]">{{ meta.total ?? 0 }}</p>
        </div>
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">Santé Système</p>
          <div class="flex items-center justify-center gap-1.5">
            <span class="w-2 h-2 rounded-full bg-green-500"></span>
            <p class="text-lg font-bold text-green-600">99.9% Up</p>
          </div>
        </div>
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">Total Prestations</p>
          <p class="text-lg font-bold text-[#001d32]">{{ meta.total ?? 0 }}</p>
        </div>
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      title="Supprimer cette prestation ?"
      :message="`Voulez-vous supprimer : ${deleteConfirm.item?.title} ?`"
      confirm-label="Supprimer"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />

    <!-- Floating create button -->
    <RouterLink
      to="/admin/prestations/create"
      class="fixed bottom-8 right-8 w-14 h-14 rounded-full flex items-center justify-center text-white shadow-lg hover:shadow-xl transition-shadow z-40"
      style="background-color:#006d35;"
      title="Nouvelle prestation"
    >
      <PlusIcon class="w-6 h-6" />
    </RouterLink>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { prestationService } from '@/services/prestationService'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/utils/useToast'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'
import {
  MagnifyingGlassIcon, PlusIcon, TrashIcon,
  UserIcon, CurrencyEuroIcon,
} from '@heroicons/vue/24/outline'

const toast       = useToast()
const prestations = ref([])
const categories  = ref([])
const loading     = ref(false)
const meta        = ref({ current_page: 1, last_page: 1, total: 0 })
const filters     = reactive({ status: '', category_id: '', search: '' })
const selectedTypes = ref([])
const deleteConfirm = reactive({ show: false, item: null, loading: false })

const priceTypes = [
  { value: 'fixed',  label: 'Training' },
  { value: 'hourly', label: 'Atelier' },
  { value: 'quote',  label: 'Conférence' },
]

const statusOptions = [
  { value: '',          label: 'Tous' },
  { value: 'published', label: 'Publié' },
  { value: 'draft',     label: 'Brouillon' },
  { value: 'suspended', label: 'Suspendu' },
  { value: 'archived',  label: 'Archivé' },
]

function resetFilters() {
  filters.status = ''
  filters.category_id = ''
  filters.search = ''
  selectedTypes.value = []
  fetchPrestations()
}

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchPrestations, 400)
}

async function fetchPrestations(page = 1) {
  loading.value = true
  try {
    const res = await prestationService.list({ page, ...filters })
    prestations.value = res.data
    meta.value = res.meta
  } catch {
    toast.showError('Impossible de charger les prestations')
  } finally {
    loading.value = false
  }
}

async function changeStatus(prestation, status) {
  try {
    const updated = await prestationService.updateStatus(prestation.id, status)
    prestation.status = updated.status
    toast.showSuccess('Statut mis à jour')
  } catch {
    toast.showError('Une erreur est survenue')
  }
}

function openDelete(item) {
  deleteConfirm.item = item
  deleteConfirm.show = true
}

async function executeDelete() {
  deleteConfirm.loading = true
  try {
    await prestationService.remove(deleteConfirm.item.id)
    prestations.value = prestations.value.filter((p) => p.id !== deleteConfirm.item.id)
    toast.showSuccess('Prestation supprimée')
    deleteConfirm.show = false
  } catch {
    toast.showError('Une erreur est survenue')
  } finally {
    deleteConfirm.loading = false
  }
}

const priceTypeLabel = (type) => ({
  fixed:  'Formation',
  hourly: 'Atelier',
  quote:  'Sur devis',
}[type] || type)

function prestationStatusBadge(status) {
  if (status === 'published') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'suspended') return 'bg-[#fee2e2] text-[#991b1b]'
  if (status === 'archived')  return 'bg-[#fee2e2] text-[#991b1b]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function prestationStatusText(status) {
  const map = { published: 'Publié', draft: 'Brouillon', suspended: 'Suspendu', archived: 'Archivé' }
  return map[status] || status
}

onMounted(async () => {
  try {
    const catData = await categoryService.list()
    categories.value = catData.data
  } catch {}
  fetchPrestations()
})
</script>
