<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h2 class="text-base font-medium text-gray-700">{{ meta.total ?? 0 }} événement(s)</h2>
      <RouterLink to="/admin/events/create" class="btn-primary text-sm px-4 py-2 rounded-md">
        + Nouvel événement
      </RouterLink>
    </div>

    <div class="card p-4 flex gap-3">
      <select v-model="filters.status" @change="fetchEvents" class="input w-44">
        <option value="">Tous statuts</option>
        <option value="draft">Brouillon</option>
        <option value="published">Publié</option>
        <option value="cancelled">Annulé</option>
      </select>
    </div>

    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Titre</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Lieu</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Places</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="6" class="px-4 py-8 text-center">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-green-600 mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="!events.length">
              <td colspan="6" class="px-4 py-8 text-center text-gray-400">Aucun événement trouvé.</td>
            </tr>
            <tr v-for="e in events" :key="e.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900 max-w-xs truncate">{{ e.title }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ e.location || '—' }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(e.start_date) }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ e.max_participants ?? '∞' }}</td>
              <td class="px-4 py-3">
                <AppBadge :label="e.status" />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <button
                    v-if="e.status === 'draft'"
                    @click="changeStatus(e, 'published')"
                    class="text-xs text-green-600 hover:text-green-800 font-medium"
                  >
                    Publier
                  </button>
                  <button
                    v-if="e.status === 'published'"
                    @click="changeStatus(e, 'cancelled')"
                    class="text-xs text-orange-600 hover:text-orange-800 font-medium"
                  >
                    Annuler
                  </button>
                  <RouterLink :to="`/admin/events/${e.id}/edit`" class="text-xs text-blue-600 hover:text-blue-800 font-medium">
                    Modifier
                  </RouterLink>
                  <button @click="openDelete(e)" class="text-xs text-red-600 hover:text-red-800 font-medium">
                    Suppr.
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-4 py-3 border-t border-gray-200">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchEvents" />
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      title="Supprimer l'événement"
      :message="`Supprimer &quot;${deleteConfirm.item?.title}&quot; ?`"
      confirm-label="Supprimer"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { eventService } from '@/services/eventService'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/common/AppBadge.vue'
import AppPagination from '@/components/common/AppPagination.vue'
import AppConfirmDialog from '@/components/common/AppConfirmDialog.vue'

const toast   = useToast()
const events  = ref([])
const loading = ref(false)
const meta    = ref({ current_page: 1, last_page: 1, total: 0 })
const filters = reactive({ status: '' })
const deleteConfirm = reactive({ show: false, item: null, loading: false })

async function fetchEvents(page = 1) {
  loading.value = true
  try {
    const res = await eventService.list({ page, ...filters })
    events.value = res.data
    meta.value   = res.meta
  } catch {
    toast.showError('Impossible de charger les événements.')
  } finally {
    loading.value = false
  }
}

async function changeStatus(event, status) {
  try {
    const updated = await eventService.updateStatus(event.id, status)
    event.status = updated.status
    toast.showSuccess(`Statut mis à jour : ${status}`)
  } catch {
    toast.showError('Impossible de modifier le statut.')
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
    toast.showSuccess('Événement supprimé.')
    deleteConfirm.show = false
  } catch {
    toast.showError('Impossible de supprimer.')
  } finally {
    deleteConfirm.loading = false
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

onMounted(fetchEvents)
</script>
