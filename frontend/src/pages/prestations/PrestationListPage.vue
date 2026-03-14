<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h2 class="text-base font-medium text-gray-700">{{ meta.total ?? 0 }} prestation(s)</h2>
      <RouterLink to="/admin/prestations/create" class="btn-primary text-sm px-4 py-2 rounded-md">
        + Nouvelle prestation
      </RouterLink>
    </div>

    <div class="card p-4 flex flex-wrap gap-3">
      <select v-model="filters.status" @change="fetchPrestations" class="input w-44">
        <option value="">Tous statuts</option>
        <option value="draft">Brouillon</option>
        <option value="published">Publié</option>
        <option value="suspended">Suspendu</option>
        <option value="archived">Archivé</option>
      </select>
      <select v-model="filters.category_id" @change="fetchPrestations" class="input w-48">
        <option value="">Toutes catégories</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
    </div>

    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Titre</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Catégorie</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prestataire</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prix</th>
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
            <tr v-else-if="!prestations.length">
              <td colspan="6" class="px-4 py-8 text-center text-gray-400">Aucune prestation trouvée.</td>
            </tr>
            <tr v-for="p in prestations" :key="p.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900 max-w-xs truncate">{{ p.title }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ p.category?.name || '—' }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ p.provider?.first_name }} {{ p.provider?.last_name }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">
                {{ p.price ? `${p.price} €` : 'Sur devis' }}
                <span class="text-xs text-gray-400">({{ priceTypeLabel(p.price_type) }})</span>
              </td>
              <td class="px-4 py-3">
                <AppBadge :label="p.status" />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <button
                    v-if="p.status === 'draft'"
                    @click="changeStatus(p, 'published')"
                    class="text-xs text-green-600 hover:text-green-800 font-medium"
                  >
                    Publier
                  </button>
                  <button
                    v-if="p.status === 'published'"
                    @click="changeStatus(p, 'suspended')"
                    class="text-xs text-orange-600 hover:text-orange-800 font-medium"
                  >
                    Suspendre
                  </button>
                  <RouterLink :to="`/admin/prestations/${p.id}/edit`" class="text-xs text-blue-600 hover:text-blue-800 font-medium">
                    Modifier
                  </RouterLink>
                  <button @click="openDelete(p)" class="text-xs text-red-600 hover:text-red-800 font-medium">
                    Suppr.
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-4 py-3 border-t border-gray-200">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchPrestations" />
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      title="Supprimer la prestation"
      :message="`Êtes-vous sûr de vouloir supprimer &quot;${deleteConfirm.item?.title}&quot; ?`"
      confirm-label="Supprimer"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { prestationService } from '@/services/prestationService'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'

const toast       = useToast()
const prestations = ref([])
const categories  = ref([])
const loading     = ref(false)
const meta        = ref({ current_page: 1, last_page: 1, total: 0 })
const filters     = reactive({ status: '', category_id: '' })
const deleteConfirm = reactive({ show: false, item: null, loading: false })

const priceTypeLabel = (t) => ({ fixed: 'fixe', hourly: '/h', quote: 'devis' }[t] || t)

async function fetchPrestations(page = 1) {
  loading.value = true
  try {
    const res = await prestationService.list({ page, ...filters })
    prestations.value = res.data
    meta.value = res.meta
  } catch {
    toast.showError('Impossible de charger les prestations.')
  } finally {
    loading.value = false
  }
}

async function changeStatus(prestation, status) {
  try {
    const updated = await prestationService.updateStatus(prestation.id, status)
    prestation.status = updated.status
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
    await prestationService.remove(deleteConfirm.item.id)
    prestations.value = prestations.value.filter((p) => p.id !== deleteConfirm.item.id)
    toast.showSuccess('Prestation supprimée.')
    deleteConfirm.show = false
  } catch {
    toast.showError('Impossible de supprimer.')
  } finally {
    deleteConfirm.loading = false
  }
}

onMounted(async () => {
  const catData = await categoryService.list()
  categories.value = catData.data
  fetchPrestations()
})
</script>
