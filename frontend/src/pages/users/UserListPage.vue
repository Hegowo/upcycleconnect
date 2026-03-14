<template>
  <div class="space-y-4">
    <!-- Filtres -->
    <div class="card p-4 flex flex-wrap gap-3">
      <input
        v-model="filters.search"
        @input="debouncedFetch"
        type="text"
        class="input w-64"
        placeholder="Rechercher par nom ou email..."
      />
      <select v-model="filters.status" @change="fetchUsers" class="input w-40">
        <option value="">Tous les statuts</option>
        <option value="active">Actif</option>
        <option value="inactive">Inactif</option>
        <option value="banned">Banni</option>
      </select>
    </div>

    <!-- Table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nom</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Inscription</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-green-600 mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="!users.length">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400">Aucun utilisateur trouvé.</td>
            </tr>
            <tr
              v-for="user in users"
              :key="user.id"
              class="hover:bg-gray-50 cursor-pointer"
              @click="router.push(`/admin/users/${user.id}`)"
            >
              <td class="px-4 py-3 text-sm font-medium text-gray-900">
                {{ user.first_name }} {{ user.last_name }}
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ user.email }}</td>
              <td class="px-4 py-3">
                <AppBadge :label="user.status" />
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">
                {{ formatDate(user.created_at) }}
              </td>
              <td class="px-4 py-3 text-right" @click.stop>
                <div class="flex justify-end gap-2">
                  <button
                    v-if="user.status !== 'banned'"
                    @click="openConfirm('ban', user)"
                    class="text-xs text-red-600 hover:text-red-800 font-medium"
                  >
                    Bannir
                  </button>
                  <button
                    v-else
                    @click="openConfirm('activate', user)"
                    class="text-xs text-green-600 hover:text-green-800 font-medium"
                  >
                    Activer
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-4 py-3 border-t border-gray-200">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" :total="meta.total" @page-change="changePage" />
      </div>
    </div>

    <AppConfirmDialog
      :show="confirm.show"
      :title="confirm.action === 'ban' ? 'Bannir l\'utilisateur' : 'Activer l\'utilisateur'"
      :message="confirm.action === 'ban'
        ? `Êtes-vous sûr de vouloir bannir ${confirm.user?.first_name} ${confirm.user?.last_name} ?`
        : `Êtes-vous sûr de vouloir activer ${confirm.user?.first_name} ${confirm.user?.last_name} ?`"
      :confirm-label="confirm.action === 'ban' ? 'Bannir' : 'Activer'"
      :confirm-variant="confirm.action === 'ban' ? 'danger' : 'primary'"
      :loading="confirm.loading"
      @confirm="executeAction"
      @cancel="confirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { userService } from '@/services/userService'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/common/AppBadge.vue'
import AppPagination from '@/components/common/AppPagination.vue'
import AppConfirmDialog from '@/components/common/AppConfirmDialog.vue'

const router = useRouter()
const toast  = useToast()

const users   = ref([])
const loading = ref(false)
const meta    = ref({ current_page: 1, last_page: 1, total: 0 })
const filters = reactive({ search: '', status: '' })
const confirm = reactive({ show: false, action: '', user: null, loading: false })

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
    toast.showError('Impossible de charger les utilisateurs.')
  } finally {
    loading.value = false
  }
}

function changePage(page) { fetchUsers(page) }

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
      toast.showSuccess(`${confirm.user.first_name} ${confirm.user.last_name} a été banni.`)
    } else {
      await userService.activate(confirm.user.id)
      toast.showSuccess(`${confirm.user.first_name} ${confirm.user.last_name} a été activé.`)
    }
    confirm.show = false
    fetchUsers(meta.value.current_page)
  } catch {
    toast.showError('L\'action a échoué.')
  } finally {
    confirm.loading = false
  }
}

onMounted(() => fetchUsers())
</script>
