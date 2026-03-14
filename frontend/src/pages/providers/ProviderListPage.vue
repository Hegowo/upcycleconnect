<template>
  <div class="space-y-4">
    <div class="card p-4 flex gap-3">
      <select v-model="filters.status" @change="fetchProviders" class="input w-48">
        <option value="">Tous les statuts</option>
        <option value="pending">En attente</option>
        <option value="approved">Approuvé</option>
        <option value="suspended">Suspendu</option>
      </select>
      <input
        v-model="filters.search"
        @input="debouncedFetch"
        type="text"
        class="input w-64"
        placeholder="Rechercher..."
      />
    </div>

    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Entreprise</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Contact</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Inscription</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="5" class="px-4 py-8 text-center">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-green-600 mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="!providers.length">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400">Aucun prestataire trouvé.</td>
            </tr>
            <tr
              v-for="p in providers"
              :key="p.id"
              class="hover:bg-gray-50"
            >
              <td class="px-4 py-3 text-sm font-medium text-gray-900">
                <RouterLink :to="`/admin/providers/${p.id}`" class="hover:text-green-600">
                  {{ p.profile?.company_name || '—' }}
                </RouterLink>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">
                {{ p.first_name }} {{ p.last_name }}<br />
                <span class="text-xs">{{ p.email }}</span>
              </td>
              <td class="px-4 py-3">
                <AppBadge :label="p.profile?.status || '—'" />
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(p.created_at) }}</td>
              <td class="px-4 py-3 text-right" @click.stop>
                <div class="flex justify-end gap-2">
                  <button
                    v-if="p.profile?.status !== 'approved'"
                    @click="changeStatus(p, 'approved')"
                    class="text-xs text-green-600 hover:text-green-800 font-medium"
                  >
                    Approuver
                  </button>
                  <button
                    v-if="p.profile?.status !== 'suspended'"
                    @click="changeStatus(p, 'suspended')"
                    class="text-xs text-red-600 hover:text-red-800 font-medium"
                  >
                    Suspendre
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-4 py-3 border-t border-gray-200">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchProviders" />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { providerService } from '@/services/userService'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/common/AppBadge.vue'
import AppPagination from '@/components/common/AppPagination.vue'

const toast     = useToast()
const providers = ref([])
const loading   = ref(false)
const meta      = ref({ current_page: 1, last_page: 1 })
const filters   = reactive({ status: '', search: '' })
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
    toast.showError('Impossible de charger les prestataires.')
  } finally {
    loading.value = false
  }
}

async function changeStatus(provider, status) {
  try {
    await providerService.updateStatus(provider.id, status)
    provider.profile.status = status
    toast.showSuccess(`Statut mis à jour : ${status}`)
  } catch {
    toast.showError('Impossible de modifier le statut.')
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR')
}

onMounted(() => fetchProviders())
</script>
