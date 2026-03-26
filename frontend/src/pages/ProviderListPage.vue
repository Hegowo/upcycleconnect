<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-3">
        <div>
          <h2 class="text-2xl font-bold text-[#001d32]">Prestataires</h2>
          <p class="text-sm text-[#40617f] mt-0.5">Service Providers — Gestion des prestataires inscrits</p>
        </div>
        <span class="text-sm font-semibold px-3 py-1 rounded-full" style="background:#dcfce7; color:#166534;">
          {{ meta.total ?? providers.length }} inscrits
        </span>
      </div>
    </div>

    <!-- Filters bar -->
    <div class="bg-white rounded-2xl p-4 border border-[#f1f5f9] shadow-sm flex flex-wrap gap-3 items-center">
      <div class="relative flex-1 min-w-[240px]">
        <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          v-model="filters.search"
          @input="debouncedFetch"
          type="text"
          class="w-full pl-9 pr-4 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
          placeholder="Rechercher par nom, SIRET..."
        />
      </div>
      <select v-model="filters.status" @change="fetchProviders" class="text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
        <option value="">Tous Statuts</option>
        <option value="pending">En attente</option>
        <option value="approved">Approuvé</option>
        <option value="suspended">Suspendu</option>
      </select>
      <div class="ml-auto flex items-center gap-2">
        <button class="px-4 py-2 text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition flex items-center gap-2">
          <ArrowDownTrayIcon class="w-4 h-4" />
          Exporter CSV
        </button>
        <button class="px-4 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2" style="background-color:#006d35;">
          <PlusIcon class="w-4 h-4" />
          Ajouter Prestataire
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <span class="font-semibold text-sm text-[#001d32]">Liste des prestataires</span>
        <span class="text-xs px-2.5 py-1 rounded-full font-semibold" style="background:#dcfce7; color:#166534;">
          {{ meta.total ?? providers.length }} total
        </span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-[#f8fafc]">
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Raison Sociale</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">SIRET</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Contact</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Catégorie</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Statut</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#f8fafc]">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-36"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-44"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!providers.length">
              <td :colspan="6" class="px-6 py-16 text-center">
                <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
                  <MagnifyingGlassIcon class="w-6 h-6 text-gray-400" />
                </div>
                <p class="text-gray-500 font-medium">Aucun prestataire trouvé</p>
                <p class="text-gray-400 text-sm mt-1">Modifiez vos critères de recherche</p>
              </td>
            </tr>
            <tr
              v-else
              v-for="p in providers"
              :key="p.id"
              class="hover:bg-[#f8fafc] transition-colors"
              :class="p.profile?.status === 'pending' ? 'border-l-2 border-orange-400' : ''"
            >
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#006d35;">
                    {{ providerInitials(p) }}
                  </div>
                  <div>
                    <RouterLink :to="`/admin/providers/${p.id}`" class="text-sm font-semibold text-[#001d32] hover:text-[#006d35] transition">
                      {{ p.profile?.company_name || p.company_name || '—' }}
                    </RouterLink>
                    <p class="text-xs text-gray-400">{{ p.profile?.city || p.city || '' }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-xs font-mono text-gray-500">
                {{ p.profile?.siret || p.siret || '—' }}
              </td>
              <td class="px-6 py-4">
                <p class="text-sm text-[#001d32] font-medium">{{ p.user?.first_name || p.first_name }} {{ p.user?.last_name || p.last_name }}</p>
                <p class="text-xs text-gray-400">{{ p.user?.email || p.email }}</p>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-medium px-2.5 py-1 rounded-full bg-[#f1f5f9] text-[#475569]">
                  {{ p.profile?.category || p.category || '—' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="providerStatusBadge(p.profile?.status || p.status)">
                  {{ providerStatusText(p.profile?.status || p.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-right" @click.stop>
                <div class="flex justify-end gap-1.5">
                  <button class="p-1.5 rounded-lg text-gray-400 hover:text-[#40617f] hover:bg-blue-50 transition" title="Modifier">
                    <PencilIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="(p.profile?.status || p.status) !== 'approved'"
                    @click="changeStatus(p, 'approved')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#dcfce7; color:#166534;"
                  >
                    Approuver
                  </button>
                  <button
                    v-if="(p.profile?.status || p.status) !== 'suspended'"
                    @click="changeStatus(p, 'suspended')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#fee2e2; color:#991b1b;"
                  >
                    Suspendre
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-6 py-4 border-t border-[#f1f5f9]">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchProviders" />
      </div>
    </div>

    <!-- Bottom cards -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div class="rounded-2xl p-5 text-white" style="background: linear-gradient(135deg, #006d35, #1b8848);">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center bg-white/20">
            <ArrowTrendingUpIcon class="w-5 h-5 text-white" />
          </div>
          <div>
            <p class="text-sm font-bold">Tendance de Croissance</p>
            <p class="text-xs text-white/70">Prestataires — 30 derniers jours</p>
          </div>
        </div>
        <p class="text-3xl font-bold">+12%</p>
        <p class="text-xs text-white/70 mt-1">Hausse des nouvelles inscriptions</p>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center bg-orange-50">
            <ClockIcon class="w-5 h-5 text-orange-500" />
          </div>
          <div>
            <p class="text-sm font-bold text-[#001d32]">File de Vérification</p>
            <p class="text-xs text-gray-400">Prestataires en attente d'examen</p>
          </div>
        </div>
        <p class="text-3xl font-bold" :class="pendingProviders > 0 ? 'text-orange-500' : 'text-[#001d32]'">{{ pendingProviders }}</p>
        <p class="text-xs text-gray-400 mt-1">En attente de validation</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { providerService } from '@/services/userService'
import { useToast } from '@/utils/useToast'
import AppPagination from '@/components/AppPagination.vue'
import {
  MagnifyingGlassIcon, PlusIcon, ArrowDownTrayIcon,
  PencilIcon, ClockIcon, ArrowTrendingUpIcon,
} from '@heroicons/vue/24/outline'

const toast     = useToast()
const providers = ref([])
const loading   = ref(false)
const meta      = ref({ current_page: 1, last_page: 1, total: 0 })
const filters   = reactive({ status: '', search: '' })

const pendingProviders = computed(() => providers.value.filter(p => (p.profile?.status || p.status) === 'pending').length)

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => fetchProviders(), 400)
}

async function fetchProviders(page = 1) {
  loading.value = true
  try {
    const res = await providerService.list({ page, ...filters })
    providers.value = res.data
    meta.value      = res.meta
  } catch {
    toast.showError('Impossible de charger les prestataires')
  } finally {
    loading.value = false
  }
}

async function changeStatus(provider, status) {
  try {
    await providerService.updateStatus(provider.id, status)
    if (provider.profile) provider.profile.status = status
    else provider.status = status
    toast.showSuccess(`Statut mis à jour : ${status}`)
  } catch {
    toast.showError('Une erreur est survenue')
  }
}

function providerInitials(p) {
  const name = p.profile?.company_name || p.company_name || p.first_name || ''
  const words = name.split(' ')
  if (words.length >= 2) return (words[0][0] + words[1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
}

function providerStatusBadge(status) {
  if (status === 'approved') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'pending')  return 'bg-[#fef9c3] text-[#854d0e]'
  if (status === 'suspended') return 'bg-[#fee2e2] text-[#991b1b]'
  return 'bg-[#f1f5f9] text-[#475569]'
}

function providerStatusText(status) {
  const map = { approved: 'Approuvé', pending: 'En attente', suspended: 'Suspendu' }
  return map[status] || status || '—'
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR')
}

onMounted(() => fetchProviders())
</script>
