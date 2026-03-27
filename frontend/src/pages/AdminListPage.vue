<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between gap-3">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('admins.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ t('admins.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/admins/create" class="p-2 sm:px-4 sm:py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2 shrink-0" style="background-color:#006d35;">
        <PlusIcon class="w-4 h-4" />
        <span class="hidden sm:inline">{{ t('admins.createBtn') }}</span>
      </RouterLink>
    </div>

    
    <div class="lg:hidden">
      
      <div v-if="loading" class="space-y-3">
        <div v-for="n in 4" :key="n" class="bg-white rounded-2xl border border-[#f1f5f9] h-20 animate-pulse"></div>
      </div>
      
      <div v-else-if="!admins.length" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm py-12 text-center">
        <p class="text-gray-500 font-medium text-sm">{{ t('common.noResults') }}</p>
      </div>
      
      <div v-else class="space-y-3">
        <div v-for="admin in admins" :key="admin.id" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-4">
          <div class="flex items-center gap-3 mb-3">
            <div
              class="w-10 h-10 rounded-full flex items-center justify-center text-white text-sm font-bold shrink-0"
              :style="admin.role === 'super_admin' ? 'background-color:#103652;' : 'background-color:#1B5B88;'"
            >
              {{ (admin.first_name?.[0] || '') + (admin.last_name?.[0] || '') }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-[#001d32] truncate">
                {{ admin.first_name }} {{ admin.last_name }}
                <span v-if="admin.id === currentUserId" class="text-xs text-gray-400 ml-1">(moi)</span>
              </p>
              <p class="text-xs text-gray-400 truncate">{{ admin.email }}</p>
            </div>
            <AppBadge :label="admin.status" />
          </div>
          <div class="flex items-center justify-between pt-2 border-t border-[#f8fafc]">
            <span
              class="text-xs px-2 py-0.5 rounded-full font-semibold"
              :class="admin.role === 'super_admin' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'"
            >
              {{ admin.role === 'super_admin' ? t('admins.roleSuperAdmin') : t('admins.roleAdmin') }}
            </span>
            <div class="flex items-center gap-2">
              <RouterLink :to="`/admin/admins/${admin.id}/edit`" class="text-xs px-3 py-1.5 rounded-full font-medium transition" style="background:#DBEAFE; color:#1d4ed8;">
                {{ t('admins.actionEdit') }}
              </RouterLink>
              <button
                v-if="admin.id !== currentUserId"
                @click="openDelete(admin)"
                class="text-xs px-3 py-1.5 rounded-full font-medium transition"
                style="background:#FEE2E2; color:#dc2626;"
              >
                {{ t('admins.actionDelete') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="hidden lg:block card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
        <span class="font-medium text-sm text-gray-700">{{ t('admins.title') }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 font-semibold">{{ admins.length }} {{ t('common.total') }}</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('admins.colName') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('admins.colEmail') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('admins.colRole') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('common.status') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">{{ t('admins.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" v-for="n in 4" :key="n">
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-36"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-44"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-16"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!admins.length">
              <td :colspan="5" class="px-4 py-16 text-center">
                <div class="text-4xl mb-2">🔍</div>
                <p class="text-gray-500 font-medium">{{ t('common.noResults') }}</p>
              </td>
            </tr>
            <tr v-else v-for="admin in admins" :key="admin.id" class="hover:bg-gray-50">
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div
                    class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0"
                    :style="admin.role === 'super_admin' ? 'background-color:#103652;' : 'background-color:#1B5B88;'"
                  >
                    {{ (admin.first_name?.[0] || '') + (admin.last_name?.[0] || '') }}
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">
                      {{ admin.first_name }} {{ admin.last_name }}
                      <span v-if="admin.id === currentUserId" class="text-xs text-gray-400 ml-1">(moi)</span>
                    </p>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ admin.email }}</td>
              <td class="px-4 py-3">
                <span
                  class="text-xs px-2 py-0.5 rounded-full font-semibold"
                  :class="admin.role === 'super_admin' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'"
                >
                  {{ admin.role === 'super_admin' ? t('admins.roleSuperAdmin') : t('admins.roleAdmin') }}
                </span>
              </td>
              <td class="px-4 py-3">
                <AppBadge :label="admin.status" />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <RouterLink :to="`/admin/admins/${admin.id}/edit`" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#DBEAFE; color:#1d4ed8;">
                    {{ t('admins.actionEdit') }}
                  </RouterLink>
                  <button
                    v-if="admin.id !== currentUserId"
                    @click="openDelete(admin)"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#FEE2E2; color:#dc2626;"
                  >
                    {{ t('admins.actionDelete') }}
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
      :title="t('admins.confirmDeleteTitle')"
      :message="t('admins.confirmDeleteMsg', { name: `${deleteConfirm.admin?.first_name} ${deleteConfirm.admin?.last_name}` })"
      :confirm-label="t('common.delete')"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { adminUserService } from '@/services/userService'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'
import { PlusIcon } from '@heroicons/vue/24/outline'

const { t } = useI18n()
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
    toast.showError(t('admins.toastLoadError'))
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
    toast.showSuccess(t('admins.toastDeleted', { name: `${deleteConfirm.admin.first_name} ${deleteConfirm.admin.last_name}` }))
    deleteConfirm.show = false
  } catch {
    toast.showError(t('admins.toastDeleteError'))
  } finally {
    deleteConfirm.loading = false
  }
}

onMounted(fetchAdmins)
</script>
