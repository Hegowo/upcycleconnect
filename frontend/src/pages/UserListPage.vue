<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('users.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ t('users.subtitle') }}</p>
      </div>
    </div>

    <div class="card p-4 flex flex-wrap gap-3 items-center">
      <input
        v-model="filters.search"
        @input="debouncedFetch"
        type="text"
        class="input w-64"
        :placeholder="t('users.searchPlaceholder')"
      />
      <select v-model="filters.status" @change="fetchUsers" class="input w-40">
        <option value="">{{ t('common.allStatuses') }}</option>
        <option value="active">{{ t('users.statusActive') }}</option>
        <option value="inactive">{{ t('users.statusInactive') }}</option>
        <option value="banned">{{ t('users.statusBanned') }}</option>
      </select>
    </div>

    <div class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
        <span class="font-medium text-sm text-gray-700">{{ t('users.title') }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 font-semibold">{{ meta.total }} {{ t('common.total') }}</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('users.colName') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('users.colEmail') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('users.colPhone') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('users.colStatus') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('users.colJoined') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">{{ t('users.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-32"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-44"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-16"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-16 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!users.length">
              <td :colspan="6" class="px-4 py-16 text-center">
                <div class="text-4xl mb-2">🔍</div>
                <p class="text-gray-500 font-medium">{{ t('common.noResults') }}</p>
                <p class="text-gray-400 text-sm mt-1">{{ t('common.noResultsHint') }}</p>
              </td>
            </tr>
            <tr
              v-else
              v-for="user in users"
              :key="user.id"
              class="hover:bg-gray-50 cursor-pointer"
              @click="router.push(`/admin/users/${user.id}`)"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#1B5B88;">
                    {{ (user.first_name?.[0] || '') + (user.last_name?.[0] || '') }}
                  </div>
                  <p class="text-sm font-medium text-gray-900">{{ user.first_name }} {{ user.last_name }}</p>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ user.email }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ user.phone || '—' }}</td>
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
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#FEE2E2; color:#dc2626;"
                  >
                    {{ t('users.actionBan') }}
                  </button>
                  <button
                    v-else
                    @click="openConfirm('activate', user)"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#DCFCE7; color:#15803d;"
                  >
                    {{ t('users.actionActivate') }}
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
      :title="confirm.action === 'ban' ? t('users.confirmBanTitle') : t('users.confirmActivateTitle')"
      :message="confirm.action === 'ban'
        ? t('users.confirmBanMsg', { name: `${confirm.user?.first_name} ${confirm.user?.last_name}` })
        : t('users.confirmActivateMsg', { name: `${confirm.user?.first_name} ${confirm.user?.last_name}` })"
      :confirm-label="confirm.action === 'ban' ? t('users.actionBan') : t('users.actionActivate')"
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
import { useI18n } from 'vue-i18n'
import { userService } from '@/services/userService'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'

const { t } = useI18n()
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
    toast.showError(t('users.toastLoadError'))
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
      toast.showSuccess(t('users.toastBanned', { name: `${confirm.user.first_name} ${confirm.user.last_name}` }))
    } else {
      await userService.activate(confirm.user.id)
      toast.showSuccess(t('users.toastActivated', { name: `${confirm.user.first_name} ${confirm.user.last_name}` }))
    }
    confirm.show = false
    fetchUsers(meta.value.current_page)
  } catch {
    toast.showError(t('users.toastError'))
  } finally {
    confirm.loading = false
  }
}

onMounted(() => fetchUsers())
</script>
