<template>
  <div class="space-y-6">
    
    <div class="flex items-center justify-between gap-3">
      <div class="min-w-0">
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32] truncate">{{ t('providers.title') }}</h2>
        <p class="text-sm text-[#40617f] mt-0.5 hidden sm:block">{{ t('providers.headerSubtitle') }}</p>
      </div>
      <span class="text-xs sm:text-sm font-semibold px-3 py-1 rounded-full shrink-0" style="background:#dcfce7; color:#166634;">
        {{ meta.total ?? providers.length }} {{ t('providers.registered') }}
      </span>
    </div>

    
    <div class="bg-white rounded-2xl p-3 sm:p-4 border border-[#f1f5f9] shadow-sm space-y-2 sm:space-y-0 sm:flex sm:flex-wrap sm:gap-3 sm:items-center">
      <div class="relative flex-1 min-w-0">
        <MagnifyingGlassIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
        <input
          v-model="filters.search"
          @input="debouncedFetch"
          type="text"
          class="w-full pl-9 pr-4 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
          :placeholder="t('providers.searchSiret')"
        />
      </div>
      <div class="flex items-center gap-2">
        <select v-model="filters.status" @change="fetchProviders" class="flex-1 text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
          <option value="">{{ t('providers.allStatuses') }}</option>
          <option value="pending">{{ t('providers.statusPending') }}</option>
          <option value="approved">{{ t('providers.statusApproved') }}</option>
          <option value="suspended">{{ t('providers.statusSuspended') }}</option>
        </select>
        <button class="p-2 sm:px-4 sm:py-2 text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition flex items-center gap-2 shrink-0">
          <ArrowDownTrayIcon class="w-4 h-4" />
          <span class="hidden sm:inline">{{ t('providers.exportCsv') }}</span>
        </button>
        <button class="p-2 sm:px-4 sm:py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2 shrink-0" style="background-color:#006d35;">
          <PlusIcon class="w-4 h-4" />
          <span class="hidden sm:inline">{{ t('providers.addProvider') }}</span>
        </button>
      </div>
    </div>

    
<div class="lg:hidden">
  
  <div v-if="loading" class="space-y-3">
    <div v-for="n in 5" :key="n" class="bg-white rounded-2xl border border-[#f1f5f9] h-28 animate-pulse"></div>
  </div>
  
  <div v-else-if="!providers.length" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm py-12 text-center">
    <MagnifyingGlassIcon class="w-8 h-8 text-gray-300 mx-auto mb-2" />
    <p class="text-gray-500 font-medium text-sm">{{ t('providers.noFound') }}</p>
  </div>
  
  <div v-else class="space-y-3">
    <div
      v-for="p in providers"
      :key="p.id"
      class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-4"
      :class="p.profile?.status === 'pending' ? 'border-l-4 border-l-orange-400' : ''"
    >
      <div class="flex items-center gap-3 mb-3">
        <div class="w-10 h-10 rounded-xl flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#006d35;">
          {{ providerInitials(p) }}
        </div>
        <div class="flex-1 min-w-0">
          <RouterLink :to="`/admin/providers/${p.id}`" class="text-sm font-semibold text-[#001d32] truncate block" @click.stop>
            {{ p.profile?.company_name || p.company_name || '—' }}
          </RouterLink>
          <p class="text-xs text-gray-400 truncate">{{ p.user?.email || p.email }}</p>
        </div>
        <span class="text-xs font-semibold px-2 py-0.5 rounded-full shrink-0" :class="providerStatusBadge(p.profile?.status || p.status)">
          {{ providerStatusText(p.profile?.status || p.status) }}
        </span>
      </div>
      <div class="flex items-center justify-between pt-2 border-t border-[#f8fafc]">
        <div class="flex items-center gap-2">
          <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-[#f1f5f9] text-[#475569]">
            {{ p.profile?.category || p.category || '—' }}
          </span>
          <span class="text-xs font-mono text-gray-400">{{ (p.profile?.siret || p.siret || '').slice(0, 9) }}...</span>
        </div>
        <div class="flex items-center gap-1.5">
          <button
            v-if="(p.profile?.status || p.status) !== 'approved'"
            @click="changeStatus(p, 'approved')"
            class="px-2.5 py-1 rounded-lg text-xs font-semibold transition"
            style="background:#dcfce7; color:#166534;"
          >
            {{ t('providers.actionApprove') }}
          </button>
          <button
            v-if="(p.profile?.status || p.status) !== 'suspended'"
            @click="changeStatus(p, 'suspended')"
            class="px-2.5 py-1 rounded-lg text-xs font-semibold transition"
            style="background:#fee2e2; color:#991b1b;"
          >
            {{ t('providers.actionSuspend') }}
          </button>
        </div>
      </div>
    </div>
    
    <div class="flex justify-center pt-2">
      <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchProviders" />
    </div>
  </div>
</div>

    
    <div class="hidden lg:block bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <span class="font-semibold text-sm text-[#001d32]">{{ t('providers.listTitle') }}</span>
        <span class="text-xs px-2.5 py-1 rounded-full font-semibold" style="background:#dcfce7; color:#166534;">
          {{ meta.total ?? providers.length }} total
        </span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-[#f8fafc]">
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('providers.colRaison') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('providers.colSiret') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('providers.colContact') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('providers.colCategory') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('providers.colStatus') }}</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('providers.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#f8fafc]">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-36"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-44"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!providers.length">
              <td :colspan="6" class="px-6 py-16 text-center">
                <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
                  <MagnifyingGlassIcon class="w-6 h-6 text-gray-400" />
                </div>
                <p class="text-gray-500 font-medium">{{ t('providers.noFound') }}</p>
                <p class="text-gray-400 text-sm mt-1">{{ t('common.searchCriteria') }}</p>
              </td>
            </tr>
            <tr
              v-else
              v-for="p in providers"
              :key="p.id"
              class="hover:bg-[#f8fafc] transition-colors"
              :class="p.profile?.status === 'pending' ? 'border-l-2 border-orange-400' : ''"
            >
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#006d35;">
                    {{ providerInitials(p) }}
                  </div>
                  <div>
                    <RouterLink :to="`/admin/providers/${p.id}`" class="text-sm font-semibold text-[#001d32] hover:text-[#006d35] transition">
                      {{ p.profile?.company_name || p.company_name || '—' }}
                    </RouterLink>
                    <p class="text-xs text-gray-400">{{ p.profile?.city || p.city || '' }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 text-xs font-mono text-gray-500">
                {{ p.profile?.siret || p.siret || '—' }}
              </td>
              <td class="px-6 py-4">
                <p class="text-sm text-[#001d32] font-medium">{{ p.user?.first_name || p.first_name }} {{ p.user?.last_name || p.last_name }}</p>
                <p class="text-xs text-gray-400">{{ p.user?.email || p.email }}</p>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-medium px-2.5 py-1 rounded-full bg-[#f1f5f9] text-[#475569]">
                  {{ p.profile?.category || p.category || '—' }}
                </span>
              </td>
              <td class="px-6 py-4">
                <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="providerStatusBadge(p.profile?.status || p.status)">
                  {{ providerStatusText(p.profile?.status || p.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-right" @click.stop>
                <div class="flex justify-end gap-1.5">
                  <button class="p-1.5 rounded-lg text-gray-400 hover:text-[#40617f] hover:bg-blue-50 transition" :title="t('common.edit')">
                    <PencilIcon class="w-4 h-4" />
                  </button>
                  <button
                    v-if="(p.profile?.status || p.status) !== 'approved'"
                    @click="changeStatus(p, 'approved')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#dcfce7; color:#166534;"
                  >
                    {{ t('providers.actionApprove') }}
                  </button>
                  <button
                    v-if="(p.profile?.status || p.status) !== 'suspended'"
                    @click="changeStatus(p, 'suspended')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#fee2e2; color:#991b1b;"
                  >
                    {{ t('providers.actionSuspend') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-6 py-4 border-t border-[#f1f5f9]">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchProviders" />
      </div>
    </div>

    
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
      <div class="rounded-2xl p-5 text-white" style="background: linear-gradient(135deg, #006d35, #1b8848);">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center bg-white/20">
            <ArrowTrendingUpIcon class="w-5 h-5 text-white" />
          </div>
          <div>
            <p class="text-sm font-bold">{{ t('providers.growth') }}</p>
            <p class="text-xs text-white/70">{{ t('providers.growthSub') }}</p>
          </div>
        </div>
        <p class="text-3xl font-bold">+12%</p>
        <p class="text-xs text-white/70 mt-1">{{ t('providers.growthDesc') }}</p>
      </div>

      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-9 h-9 rounded-xl flex items-center justify-center bg-orange-50">
            <ClockIcon class="w-5 h-5 text-orange-500" />
          </div>
          <div>
            <p class="text-sm font-bold text-[#001d32]">{{ t('providers.verificationQueue') }}</p>
            <p class="text-xs text-gray-400">{{ t('providers.verificationQueueSub') }}</p>
          </div>
        </div>
        <p class="text-3xl font-bold" :class="pendingProviders > 0 ? 'text-orange-500' : 'text-[#001d32]'">{{ pendingProviders }}</p>
        <p class="text-xs text-gray-400 mt-1">{{ t('providers.pendingValidation') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { providerService } from '@/services/userService'
import { useToast } from '@/utils/useToast'
import AppPagination from '@/components/AppPagination.vue'
import {
  MagnifyingGlassIcon, PlusIcon, ArrowDownTrayIcon,
  PencilIcon, ClockIcon, ArrowTrendingUpIcon,
} from '@heroicons/vue/24/outline'

const toast     = useToast()
const { t }     = useI18n()
const providers = ref([])
const loading   = ref(false)
const meta      = ref({ current_page: 1, last_page: 1, total: 0 })
const filters   = reactive({ status: '', search: '' })

const pendingProviders = computed(() => providers.value.filter(p => (p.profile?.status || p.status) === 'pending').length)

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
    if (provider.profile) provider.profile.status = status
    else provider.status = status
    toast.showSuccess(t('providers.toastStatusUpdated', { status }))
  } catch {
    toast.showError(t('providers.toastErrorOccurred'))
  }
}

function providerInitials(p) {
  const name = p.profile?.company_name || p.company_name || p.first_name || ''
  const words = name.split(' ')
  if (words.length >= 2) return (words[0][0] + words[1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
}

function providerStatusBadge(status) {
  if (status === 'approved') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'pending')  return 'bg-[#fef9c3] text-[#854d0e]'
  if (status === 'suspended') return 'bg-[#fee2e2] text-[#991b1b]'
  return 'bg-[#f1f5f9] text-[#475569]'
}

function providerStatusText(status) {
  const map = {
    approved:  t('providers.statusApproved'),
    pending:   t('providers.statusPending'),
    suspended: t('providers.statusSuspended'),
  }
  return map[status] || status || '—'
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR')
}

onMounted(() => fetchProviders())
</script>
