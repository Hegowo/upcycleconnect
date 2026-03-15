<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('providers.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ t('providers.subtitle') }}</p>
      </div>
    </div>

    <div class="card p-4 flex flex-wrap gap-3 items-center">
      <select v-model="filters.status" @change="fetchProviders" class="input w-48">
        <option value="">{{ t('common.allStatuses') }}</option>
        <option value="pending">{{ t('providers.statusPending') }}</option>
        <option value="approved">{{ t('providers.statusApproved') }}</option>
        <option value="suspended">{{ t('providers.statusSuspended') }}</option>
      </select>
      <input
        v-model="filters.search"
        @input="debouncedFetch"
        type="text"
        class="input w-64"
        :placeholder="t('providers.searchPlaceholder')"
      />
    </div>

    <div class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
        <span class="font-medium text-sm text-gray-700">{{ t('providers.title') }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 font-semibold">{{ meta.total ?? providers.length }} {{ t('common.total') }}</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('providers.colCompany') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('providers.colContact') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('providers.colStatus') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('providers.colJoined') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">{{ t('providers.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-36"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-44"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!providers.length">
              <td :colspan="5" class="px-4 py-16 text-center">
                <div class="text-4xl mb-2">🔍</div>
                <p class="text-gray-500 font-medium">{{ t('common.noResults') }}</p>
                <p class="text-gray-400 text-sm mt-1">{{ t('common.noResultsHint') }}</p>
              </td>
            </tr>
            <tr
              v-else
              v-for="p in providers"
              :key="p.id"
              class="hover:bg-gray-50 border-l-4"
              :style="p.profile?.status === 'pending' ? 'border-left-color: #f97316;' : 'border-left-color: transparent;'"
            >
              <td class="px-4 py-3">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#1B5B88;">
                    {{ (p.profile?.company_name?.[0] || p.first_name?.[0] || '').toUpperCase() }}{{ (p.profile?.company_name?.split(' ')?.[1]?.[0] || p.last_name?.[0] || '').toUpperCase() }}
                  </div>
                  <RouterLink :to="`/admin/providers/${p.id}`" class="text-sm font-medium text-gray-900 hover:text-green-600">
                    {{ p.profile?.company_name || '—' }}
                  </RouterLink>
                </div>
              </td>
              <td class="px-4 py-3">
                <p class="text-sm text-gray-700">{{ p.first_name }} {{ p.last_name }}</p>
                <p class="text-xs text-gray-400">{{ p.email }}</p>
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
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#DCFCE7; color:#15803d;"
                  >
                    {{ t('providers.actionApprove') }}
                  </button>
                  <button
                    v-if="p.profile?.status !== 'suspended'"
                    @click="changeStatus(p, 'suspended')"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#FEE2E2; color:#dc2626;"
                  >
                    {{ t('providers.actionSuspend') }}
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
import { useI18n } from 'vue-i18n'
import { providerService } from '@/services/userService'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppPagination from '@/components/AppPagination.vue'

const { t } = useI18n()
const toast     = useToast()
const providers = ref([])
const loading   = ref(false)
const meta      = ref({ current_page: 1, last_page: 1, total: 0 })
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
    toast.showError(t('providers.toastLoadError'))
  } finally {
    loading.value = false
  }
}

async function changeStatus(provider, status) {
  try {
    await providerService.updateStatus(provider.id, status)
    provider.profile.status = status
    toast.showSuccess(t('providers.toastStatusUpdated', { status }))
  } catch {
    toast.showError(t('providers.toastStatusError'))
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR')
}

onMounted(() => fetchProviders())
</script>
