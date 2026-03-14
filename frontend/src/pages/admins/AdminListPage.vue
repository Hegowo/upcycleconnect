<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h2 class="text-base font-medium text-gray-700">Comptes administrateurs</h2>
      <RouterLink to="/admin/admins/create" class="btn-primary text-sm px-4 py-2 rounded-md">
        + Nouvel administrateur
      </RouterLink>
    </div>

    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nom</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Rôle</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Statut</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="5" class="px-4 py-8 text-center">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-green-600 mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="!admins.length">
              <td colspan="5" class="px-4 py-8 text-center text-gray-400">Aucun administrateur.</td>
            </tr>
            <tr v-for="admin in admins" :key="admin.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900">
                {{ admin.first_name }} {{ admin.last_name }}
                <span v-if="admin.id === currentUserId" class="text-xs text-gray-400">(moi)</span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ admin.email }}</td>
              <td class="px-4 py-3 text-sm">
                <AppBadge :label="admin.role" :variant="admin.role === 'super_admin' ? 'info' : 'neutral'" />
              </td>
              <td class="px-4 py-3">
                <AppBadge :label="admin.status" />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <RouterLink :to="`/admin/admins/${admin.id}/edit`" class="text-xs text-blue-600 hover:text-blue-800 font-medium">
                    Modifier
                  </RouterLink>
                  <button
                    v-if="admin.id !== currentUserId"
                    @click="openDelete(admin)"
                    class="text-xs text-red-600 hover:text-red-800 font-medium"
                  >
                    Désactiver
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      title="Désactiver l'administrateur"
      :message="`Désactiver le compte de ${deleteConfirm.admin?.first_name} ${deleteConfirm.admin?.last_name} ?`"
      confirm-label="Désactiver"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { adminUserService } from '@/services/userService'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/common/AppBadge.vue'
import AppConfirmDialog from '@/components/common/AppConfirmDialog.vue'

const toast  = useToast()
const auth   = useAuthStore()
const admins = ref([])
const loading = ref(false)
const currentUserId = computed(() => auth.user?.id)
const deleteConfirm = reactive({ show: false, admin: null, loading: false })

async function fetchAdmins() {
  loading.value = true
  try {
    const data = await adminUserService.list()
    admins.value = data.data
  } catch {
    toast.showError('Impossible de charger les administrateurs.')
  } finally {
    loading.value = false
  }
}

function openDelete(admin) {
  deleteConfirm.admin = admin
  deleteConfirm.show  = true
}

async function executeDelete() {
  deleteConfirm.loading = true
  try {
    await adminUserService.remove(deleteConfirm.admin.id)
    admins.value = admins.value.filter((a) => a.id !== deleteConfirm.admin.id)
    toast.showSuccess('Administrateur désactivé.')
    deleteConfirm.show = false
  } catch {
    toast.showError('Impossible de désactiver ce compte.')
  } finally {
    deleteConfirm.loading = false
  }
}

onMounted(fetchAdmins)
</script>
