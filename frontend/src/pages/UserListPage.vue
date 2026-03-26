<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-[#001d32]">Gestion Utilisateurs</h2>
        <p class="text-sm text-[#40617f] mt-0.5">User Management — Gérez les membres de la plateforme</p>
      </div>
      <div class="flex items-center gap-2">
        <button class="px-4 py-2 text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition flex items-center gap-2">
          <AdjustmentsHorizontalIcon class="w-4 h-4" />
          Filtres Avancés
        </button>
        <button class="px-4 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2" style="background-color:#006d35;">
          <PlusIcon class="w-4 h-4" />
          Ajouter Membre
        </button>
      </div>
    </div>

    <!-- Stat cards -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center" style="background-color:#dcfce7;">
            <UsersIcon class="w-5 h-5" style="color:#006d35;" />
          </div>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">{{ meta.total }}</p>
        <p class="text-xs text-gray-400 mt-1">Total Membres</p>
      </div>
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center" style="background-color:#dbeafe;">
            <CheckCircleIcon class="w-5 h-5 text-blue-600" />
          </div>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">{{ activeCount }}</p>
        <p class="text-xs text-gray-400 mt-1">Utilisateurs Actifs</p>
      </div>
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center" style="background-color:#f0fdf4;">
            <BriefcaseIcon class="w-5 h-5" style="color:#006d35;" />
          </div>
        </div>
        <p class="text-2xl font-bold text-[#001d32]">—</p>
        <p class="text-xs text-gray-400 mt-1">Prestataires Vérifiés</p>
      </div>
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm" :class="pendingCount > 0 ? 'border-red-200' : ''">
        <div class="flex items-center justify-between mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center" :class="pendingCount > 0 ? 'bg-red-50' : 'bg-gray-50'">
            <ClockIcon class="w-5 h-5" :class="pendingCount > 0 ? 'text-red-500' : 'text-gray-400'" />
          </div>
          <span v-if="pendingCount > 0" class="text-xs font-semibold px-2 py-0.5 rounded-full bg-red-100 text-red-700">Action Requise</span>
        </div>
        <p class="text-2xl font-bold" :class="pendingCount > 0 ? 'text-red-600' : 'text-[#001d32]'">{{ pendingCount }}</p>
        <p class="text-xs text-gray-400 mt-1">En Attente d'Approbation</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 p-1 bg-[#f8fafc] rounded-xl w-fit border border-[#e5e7eb]">
      <button
        @click="activeTab = 'all'; fetchUsers()"
        :class="activeTab === 'all' ? 'bg-white shadow text-[#001d32] font-semibold' : 'text-gray-500 hover:text-gray-700'"
        class="px-4 py-2 rounded-lg text-sm transition"
      >
        Tous les Membres
      </button>
      <button
        @click="activeTab = 'professionals'; filters.status = ''; fetchUsers()"
        :class="activeTab === 'professionals' ? 'bg-white shadow text-[#001d32] font-semibold' : 'text-gray-500 hover:text-gray-700'"
        class="px-4 py-2 rounded-lg text-sm transition"
      >
        Professionnels
      </button>
      <button
        @click="activeTab = 'pending'; filters.status = 'pending'; fetchUsers()"
        :class="activeTab === 'pending' ? 'bg-white shadow text-[#001d32] font-semibold' : 'text-gray-500 hover:text-gray-700'"
        class="px-4 py-2 rounded-lg text-sm transition"
      >
        En Attente de Validation
      </button>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-2xl p-4 border border-[#f1f5f9] shadow-sm flex flex-wrap gap-3 items-center">
      <div class="relative flex-1 max-w-sm">
        <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          v-model="filters.search"
          @input="debouncedFetch"
          type="text"
          class="w-full pl-9 pr-4 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
          placeholder="Rechercher un utilisateur..."
        />
      </div>
      <select v-model="filters.status" @change="fetchUsers" class="text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
        <option value="">Tous les statuts</option>
        <option value="active">Actif</option>
        <option value="inactive">Inactif</option>
        <option value="banned">Banni</option>
      </select>
    </div>

    <!-- Table -->
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <span class="font-semibold text-sm text-[#001d32]">Liste des membres</span>
        <span class="text-xs px-2.5 py-1 rounded-full font-semibold" style="background:#dcfce7; color:#166534;">{{ meta.total }} total</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-[#f8fafc]">
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Nom / Email</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Rôle</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Statut</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Siret / ID</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Inscription</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#f8fafc]">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-40"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-24"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-16"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-16 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!users.length">
              <td :colspan="6" class="px-6 py-16 text-center">
                <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
                  <MagnifyingGlassIcon class="w-6 h-6 text-gray-400" />
                </div>
                <p class="text-gray-500 font-medium">Aucun résultat</p>
                <p class="text-gray-400 text-sm mt-1">Modifiez vos critères de recherche</p>
              </td>
            </tr>
            <tr
              v-else
              v-for="user in users"
              :key="user.id"
              class="hover:bg-[#f8fafc] cursor-pointer transition-colors"
              @click="router.push(`/admin/users/${user.id}`)"
            >
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#001d32;">
                    {{ (user.first_name?.[0] || '') + (user.last_name?.[0] || '') }}
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-[#001d32]">{{ user.first_name }} {{ user.last_name }}</p>
                    <p class="text-xs text-gray-400">{{ user.email }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-medium px-2.5 py-1 rounded-full"
                  :class="hasProviderRole(user) ? 'bg-[#dbeafe] text-blue-700' : 'bg-[#f1f5f9] text-[#475569]'">
                  {{ hasProviderRole(user) ? 'Professionnel' : 'Particulier' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="statusBadge(user.status)">
                  {{ statusText(user.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-xs font-mono text-gray-400">
                #{{ user.id }}
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">
                {{ formatDate(user.created_at) }}
              </td>
              <td class="px-6 py-4 text-right" @click.stop>
                <div class="flex justify-end gap-1.5">
                  <button @click="router.push(`/admin/users/${user.id}`)" class="p-1.5 rounded-lg text-gray-400 hover:text-[#40617f] hover:bg-blue-50 transition" title="Voir">
                    <EyeIcon class="w-4 h-4" />
                  </button>
                  <button class="p-1.5 rounded-lg text-gray-400 hover:text-[#006d35] hover:bg-green-50 transition" title="Modifier">
                    <PencilIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="user.status !== 'banned'"
                    @click="openConfirm('ban', user)"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition"
                    title="Bannir"
                  >
                    <NoSymbolIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-else
                    @click="openConfirm('activate', user)"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-green-600 hover:bg-green-50 transition"
                    title="Activer"
                  >
                    <CheckCircleIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-6 py-4 border-t border-[#f1f5f9] flex items-center justify-between">
        <p class="text-xs text-gray-400">{{ meta.total }} membres au total</p>
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" :total="meta.total" @page-change="changePage" />
      </div>
    </div>

    <!-- Bottom cards -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center bg-gray-100">
            <ShieldCheckIcon class="w-5 h-5 text-gray-500" />
          </div>
          <div>
            <p class="text-sm font-semibold text-[#001d32]">Vérification d'Intégrité des Données</p>
            <p class="text-xs text-gray-400">Analyse automatique hebdomadaire</p>
          </div>
        </div>
        <div class="flex items-center gap-2 mt-3">
          <div class="h-1.5 flex-1 bg-gray-100 rounded-full overflow-hidden">
            <div class="h-full rounded-full bg-green-500 w-[97%]"></div>
          </div>
          <span class="text-xs font-semibold text-green-600">97%</span>
        </div>
        <p class="text-xs text-gray-400 mt-1">Données cohérentes</p>
      </div>

      <RouterLink to="/admin/providers" class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm hover:shadow-md transition-shadow block" style="border-top: 3px solid #006d35;">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center" style="background-color:#dcfce7;">
            <BriefcaseIcon class="w-5 h-5" style="color:#006d35;" />
          </div>
          <div>
            <p class="text-sm font-semibold text-[#001d32]">Vérification Prestataires</p>
            <p class="text-xs text-gray-400">File d'attente de validation</p>
          </div>
        </div>
        <p class="text-xs font-medium" style="color:#006d35;">Accéder à la liste des prestataires →</p>
      </RouterLink>
    </div>

    <AppConfirmDialog
      :show="confirm.show"
      :title="confirm.action === 'ban' ? 'Bannir cet utilisateur ?' : 'Activer cet utilisateur ?'"
      :message="confirm.action === 'ban'
        ? `Voulez-vous bannir ${confirm.user?.first_name} ${confirm.user?.last_name} ?`
        : `Voulez-vous activer ${confirm.user?.first_name} ${confirm.user?.last_name} ?`"
      :confirm-label="confirm.action === 'ban' ? 'Bannir' : 'Activer'"
      :confirm-variant="confirm.action === 'ban' ? 'danger' : 'primary'"
      :loading="confirm.loading"
      @confirm="executeAction"
      @cancel="confirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { userService } from '@/services/userService'
import { useToast } from '@/utils/useToast'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'
import {
  UsersIcon, BriefcaseIcon, CheckCircleIcon, ClockIcon,
  MagnifyingGlassIcon, AdjustmentsHorizontalIcon, PlusIcon,
  EyeIcon, PencilIcon, NoSymbolIcon, ShieldCheckIcon,
} from '@heroicons/vue/24/outline'

const router = useRouter()
const toast  = useToast()

const users      = ref([])
const loading    = ref(false)
const activeTab  = ref('all')
const meta       = ref({ current_page: 1, last_page: 1, total: 0 })
const filters    = reactive({ search: '', status: '' })
const confirm    = reactive({ show: false, action: '', user: null, loading: false })

const activeCount  = computed(() => users.value.filter(u => u.status === 'active').length)
const pendingCount = computed(() => users.value.filter(u => u.status === 'pending').length)

let debounceTimer = null
function debouncedFetch() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(fetchUsers, 400)
}

async function fetchUsers(page = 1) {
  loading.value = true
  try {
    const res = await userService.list({ page, per_page: 20, ...filters })
    users.value = res.data
    meta.value  = res.meta
  } catch {
    toast.showError('Impossible de charger les utilisateurs')
  } finally {
    loading.value = false
  }
}

function changePage(page) { fetchUsers(page) }

function hasProviderRole(user) {
  return user.roles?.some(r => r === 'provider' || r?.name === 'provider')
}

function statusBadge(status) {
  if (status === 'active') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'banned') return 'bg-[#fee2e2] text-[#991b1b]'
  if (status === 'pending') return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#f1f5f9] text-[#475569]'
}

function statusText(status) {
  const map = { active: 'Actif', banned: 'Banni', inactive: 'Inactif', pending: 'En attente' }
  return map[status] || status
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function openConfirm(action, user) {
  confirm.action = action
  confirm.user   = user
  confirm.show   = true
}

async function executeAction() {
  confirm.loading = true
  try {
    if (confirm.action === 'ban') {
      await userService.ban(confirm.user.id)
      toast.showSuccess(`${confirm.user.first_name} ${confirm.user.last_name} a été banni`)
    } else {
      await userService.activate(confirm.user.id)
      toast.showSuccess(`${confirm.user.first_name} ${confirm.user.last_name} a été activé`)
    }
    confirm.show = false
    fetchUsers(meta.value.current_page)
  } catch {
    toast.showError('Une erreur est survenue')
  } finally {
    confirm.loading = false
  }
}

onMounted(() => fetchUsers())
</script>
