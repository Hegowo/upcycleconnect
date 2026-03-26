<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-[#001d32]">Événements & Équipe</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Gestion des événements et de l'équipe administrative</p>
      </div>
      <RouterLink
        to="/admin/events/create"
        class="px-4 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2"
        style="background-color:#006d35;"
      >
        <PlusIcon class="w-4 h-4" />
        Planifier un Événement
      </RouterLink>
    </div>

    <!-- Next event + health cards -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Next event -->
      <div class="lg:col-span-2 bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-[#001d32]">Prochain Atelier</h3>
          <span class="text-xs text-gray-400">Prochainement</span>
        </div>
        <div v-if="nextEvent">
          <div class="flex items-start gap-4">
            <div class="w-14 h-14 rounded-xl flex items-center justify-center shrink-0" style="background:#dcfce7;">
              <CalendarIcon class="w-7 h-7" style="color:#006d35;" />
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1 flex-wrap">
                <h4 class="text-base font-bold text-[#001d32]">{{ nextEvent.title }}</h4>
                <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-[#dcfce7] text-[#166534]">
                  {{ nextEvent.status }}
                </span>
              </div>
              <div class="flex items-center gap-4 text-xs text-gray-500 mb-3 flex-wrap">
                <span class="flex items-center gap-1">
                  <MapPinIcon class="w-3.5 h-3.5" />
                  {{ nextEvent.location || 'En ligne' }}
                </span>
                <span class="flex items-center gap-1">
                  <ClockIcon class="w-3.5 h-3.5" />
                  {{ formatDateTime(nextEvent.start_date) }}
                </span>
                <span class="flex items-center gap-1">
                  <UsersIcon class="w-3.5 h-3.5" />
                  {{ nextEvent.max_participants ?? '∞' }} places
                </span>
              </div>
              <div>
                <div class="flex justify-between text-xs text-gray-400 mb-1">
                  <span>Capacité</span>
                  <span>{{ Math.floor(Math.random() * 70 + 20) }}%</span>
                </div>
                <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div class="h-full rounded-full" style="width:65%; background:#006d35;"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="py-6 text-center text-gray-400 text-sm">
          Aucun événement à venir
        </div>
      </div>

      <!-- Event health -->
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <h3 class="font-semibold text-[#001d32] mb-4">Santé des Événements</h3>
        <div class="space-y-4">
          <div>
            <div class="flex justify-between text-xs text-gray-500 mb-1">
              <span>Taux de complétion</span>
              <span class="font-bold text-[#001d32]">94%</span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full" style="width:94%; background:#006d35;"></div>
            </div>
          </div>
          <div class="flex items-center justify-between py-2 border-t border-[#f8fafc]">
            <span class="text-xs text-gray-500">Note moyenne</span>
            <div class="flex items-center gap-1">
              <svg class="w-4 h-4 text-yellow-400" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/></svg>
              <span class="text-sm font-bold text-[#001d32]">4.9</span>
            </div>
          </div>
          <div class="flex items-center justify-between py-2 border-t border-[#f8fafc]">
            <span class="text-xs text-gray-500">Créneaux disponibles</span>
            <span class="text-sm font-bold text-[#001d32]">{{ events.filter(e => e.status === 'published').length }}</span>
          </div>
          <div class="flex items-center justify-between py-2 border-t border-[#f8fafc]">
            <span class="text-xs text-gray-500">Total événements</span>
            <span class="text-sm font-bold text-[#001d32]">{{ meta.total ?? 0 }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-2xl p-4 border border-[#f1f5f9] shadow-sm flex flex-wrap gap-3 items-center">
      <select v-model="filters.status" @change="fetchEvents" class="text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
        <option value="">Tous les statuts</option>
        <option value="draft">Brouillon</option>
        <option value="published">Publié</option>
        <option value="cancelled">Annulé</option>
      </select>
    </div>

    <!-- Events table -->
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <span class="font-semibold text-sm text-[#001d32]">Tous les événements</span>
        <span class="text-xs px-2.5 py-1 rounded-full font-semibold" style="background:#dcfce7; color:#166534;">
          {{ meta.total ?? 0 }} total
        </span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-[#f8fafc]">
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Titre</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Lieu</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Date</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Participants</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Statut</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#f8fafc]">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-40"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-32"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-12"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!events.length">
              <td :colspan="6" class="px-6 py-16 text-center">
                <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
                  <CalendarIcon class="w-6 h-6 text-gray-400" />
                </div>
                <p class="text-gray-500 font-medium">Aucun événement trouvé</p>
              </td>
            </tr>
            <tr v-else v-for="e in events" :key="e.id" class="hover:bg-[#f8fafc] transition-colors">
              <td class="px-6 py-4 text-sm font-semibold text-[#001d32] max-w-xs truncate">{{ e.title }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ e.location || '—' }}</td>
              <td class="px-6 py-4">
                <span class="text-sm font-medium text-[#001d32]">{{ formatDateShort(e.start_date) }}</span>
                <span class="text-xs text-gray-400 ml-1">{{ formatTime(e.start_date) }}</span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ e.max_participants ?? '∞' }}</td>
              <td class="px-6 py-4">
                <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="eventStatusBadge(e.status)">
                  {{ eventStatusText(e.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-1.5">
                  <button
                    v-if="e.status === 'draft'"
                    @click="changeStatus(e, 'published')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#dcfce7; color:#166534;"
                  >
                    Publier
                  </button>
                  <button
                    v-if="e.status === 'published'"
                    @click="changeStatus(e, 'cancelled')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#fff7ed; color:#c2410c;"
                  >
                    Annuler
                  </button>
                  <RouterLink
                    :to="`/admin/events/${e.id}/edit`"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-[#40617f] hover:bg-blue-50 transition"
                    title="Modifier"
                  >
                    <PencilIcon class="w-4 h-4" />
                  </RouterLink>
                  <button @click="openDelete(e)" class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition" title="Supprimer">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-6 py-4 border-t border-[#f1f5f9]">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchEvents" />
      </div>
    </div>

    <!-- Admin team section -->
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <div>
          <h3 class="font-semibold text-[#001d32]">Gestion de l'Équipe Admin</h3>
          <p class="text-xs text-gray-400 mt-0.5">Administrateurs de la plateforme</p>
        </div>
        <RouterLink
          v-if="auth.isSuperAdmin"
          to="/admin/admins/create"
          class="px-4 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2"
          style="background-color:#006d35;"
        >
          <PlusIcon class="w-4 h-4" />
          Ajouter Administrateur
        </RouterLink>
      </div>

      <div v-if="!auth.isSuperAdmin" class="px-6 py-8 text-center text-gray-400 text-sm">
        <ShieldCheckIcon class="w-8 h-8 mx-auto mb-2 text-gray-300" />
        Accès réservé aux super-administrateurs
      </div>

      <template v-else>
        <div v-if="adminsLoading" class="p-4 space-y-3">
          <div v-for="n in 3" :key="n" class="h-12 bg-gray-50 rounded-lg animate-pulse"></div>
        </div>
        <div v-else-if="!admins.length" class="px-6 py-8 text-center text-gray-400 text-sm">
          Aucun administrateur
        </div>
        <table v-else class="min-w-full">
          <thead>
            <tr class="bg-[#f8fafc]">
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Admin User</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Rôle</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Dernière Activité</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-[#6b7280] uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#f8fafc]">
            <tr v-for="admin in admins" :key="admin.id" class="hover:bg-[#f8fafc] transition-colors">
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-bold" style="background-color:#001d32;">
                    {{ ((admin.first_name?.[0] || '') + (admin.last_name?.[0] || '')).toUpperCase() }}
                  </div>
                  <div>
                    <p class="text-sm font-semibold text-[#001d32]">{{ admin.first_name }} {{ admin.last_name }}</p>
                    <p class="text-xs text-gray-400">{{ admin.email }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-semibold px-2.5 py-1 rounded-full capitalize"
                  :class="admin.role === 'super_admin' ? 'bg-[#fef9c3] text-[#854d0e]' : 'bg-[#dcfce7] text-[#166534]'">
                  {{ admin.role?.replace('_', ' ') }}
                </span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-400">
                {{ formatDateShort(admin.created_at) }}
              </td>
              <td class="px-6 py-4 text-right">
                <RouterLink :to="`/admin/admins/${admin.id}/edit`" class="p-1.5 rounded-lg text-gray-400 hover:text-[#40617f] hover:bg-blue-50 transition inline-flex">
                  <PencilIcon class="w-4 h-4" />
                </RouterLink>
              </td>
            </tr>
          </tbody>
        </table>
      </template>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      title="Supprimer cet événement ?"
      :message="`Voulez-vous supprimer l'événement : ${deleteConfirm.item?.title} ?`"
      confirm-label="Supprimer"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { eventService } from '@/services/eventService'
import api from '@/services/api'
import { useToast } from '@/utils/useToast'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'
import {
  CalendarIcon, PlusIcon, PencilIcon, TrashIcon,
  MapPinIcon, ClockIcon, UsersIcon, ShieldCheckIcon,
} from '@heroicons/vue/24/outline'

const auth  = useAuthStore()
const toast = useToast()
const events  = ref([])
const admins  = ref([])
const loading = ref(false)
const adminsLoading = ref(false)
const meta    = ref({ current_page: 1, last_page: 1, total: 0 })
const filters = reactive({ status: '' })
const deleteConfirm = reactive({ show: false, item: null, loading: false })

const nextEvent = computed(() => events.value.find(e => e.status === 'published') || events.value[0] || null)

async function fetchEvents(page = 1) {
  loading.value = true
  try {
    const res = await eventService.list({ page, ...filters })
    events.value = res.data
    meta.value   = res.meta
  } catch {
    toast.showError('Impossible de charger les événements')
  } finally {
    loading.value = false
  }
}

async function fetchAdmins() {
  if (!auth.isSuperAdmin) return
  adminsLoading.value = true
  try {
    const { data } = await api.get('/admins')
    admins.value = data.data || []
  } catch {
    // silent
  } finally {
    adminsLoading.value = false
  }
}

async function changeStatus(event, status) {
  try {
    const updated = await eventService.updateStatus(event.id, status)
    event.status = updated.status
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
    await eventService.remove(deleteConfirm.item.id)
    events.value = events.value.filter((e) => e.id !== deleteConfirm.item.id)
    toast.showSuccess('Événement supprimé')
    deleteConfirm.show = false
  } catch {
    toast.showError('Une erreur est survenue')
  } finally {
    deleteConfirm.loading = false
  }
}

function eventStatusBadge(status) {
  if (status === 'published') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'cancelled') return 'bg-[#fee2e2] text-[#991b1b]'
  if (status === 'draft')     return 'bg-[#f1f5f9] text-[#475569]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function eventStatusText(status) {
  const map = { published: 'Publié', draft: 'Brouillon', cancelled: 'Annulé' }
  return map[status] || status
}
function formatDateShort(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}
function formatDateTime(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}
function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchEvents()
  fetchAdmins()
})
</script>
