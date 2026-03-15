<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('prestations.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ t('prestations.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/prestations/create" class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90" style="background-color:#1B8848;">
        {{ t('prestations.createBtn') }}
      </RouterLink>
    </div>

    <div class="card p-4 flex flex-wrap gap-3 items-center">
      <select v-model="filters.status" @change="fetchPrestations" class="input w-44">
        <option value="">{{ t('common.allStatuses') }}</option>
        <option value="draft">{{ t('prestations.statusDraft') }}</option>
        <option value="published">{{ t('prestations.statusPublished') }}</option>
        <option value="suspended">{{ t('prestations.statusSuspended') }}</option>
        <option value="archived">{{ t('prestations.statusArchived') }}</option>
      </select>
      <select v-model="filters.category_id" @change="fetchPrestations" class="input w-48">
        <option value="">{{ t('prestations.allCategories') }}</option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
      </select>
    </div>

    <div class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
        <span class="font-medium text-sm text-gray-700">{{ t('prestations.title') }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 font-semibold">{{ meta.total ?? 0 }} {{ t('common.total') }}</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('prestations.colTitle') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('prestations.colCategory') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('prestations.colProvider') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('prestations.colPrice') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('prestations.colStatus') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">{{ t('prestations.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-40"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-32"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-16"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!prestations.length">
              <td :colspan="6" class="px-4 py-16 text-center">
                <div class="text-4xl mb-2">🔍</div>
                <p class="text-gray-500 font-medium">{{ t('common.noResults') }}</p>
                <p class="text-gray-400 text-sm mt-1">{{ t('common.noResultsHint') }}</p>
              </td>
            </tr>
            <tr v-else v-for="p in prestations" :key="p.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900 max-w-xs truncate">{{ p.title }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ p.category?.name || '—' }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ p.provider?.first_name }} {{ p.provider?.last_name }}</td>
              <td class="px-4 py-3 text-sm">
                <span class="font-medium text-gray-900">{{ p.price ? `${p.price} €` : t('prestations.priceOnQuote') }}</span>
                <span class="text-xs text-gray-400 ml-1">({{ priceTypeLabel(p.price_type) }})</span>
              </td>
              <td class="px-4 py-3">
                <AppBadge :label="p.status" />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <button
                    v-if="p.status === 'draft'"
                    @click="changeStatus(p, 'published')"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#DCFCE7; color:#15803d;"
                  >
                    {{ t('prestations.actionPublish') }}
                  </button>
                  <button
                    v-if="p.status === 'published'"
                    @click="changeStatus(p, 'suspended')"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#FFF7ED; color:#c2410c;"
                  >
                    {{ t('prestations.actionSuspend') }}
                  </button>
                  <RouterLink :to="`/admin/prestations/${p.id}/edit`" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#DBEAFE; color:#1d4ed8;">
                    {{ t('common.edit') }}
                  </RouterLink>
                  <button @click="openDelete(p)" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#FEE2E2; color:#dc2626;">
                    {{ t('common.delete') }}
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
      :title="t('prestations.confirmDeleteTitle')"
      :message="t('prestations.confirmDeleteMsg', { name: deleteConfirm.item?.title })"
      :confirm-label="t('common.delete')"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { prestationService } from '@/services/prestationService'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'

const { t } = useI18n()
const toast       = useToast()
const prestations = ref([])
const categories  = ref([])
const loading     = ref(false)
const meta        = ref({ current_page: 1, last_page: 1, total: 0 })
const filters     = reactive({ status: '', category_id: '' })
const deleteConfirm = reactive({ show: false, item: null, loading: false })

const priceTypeLabel = (type) => ({
  fixed:  t('prestations.priceTypeFixed'),
  hourly: t('prestations.priceTypeHourly'),
  quote:  t('prestations.priceTypeQuote'),
}[type] || type)

async function fetchPrestations(page = 1) {
  loading.value = true
  try {
    const res = await prestationService.list({ page, ...filters })
    prestations.value = res.data
    meta.value = res.meta
  } catch {
    toast.showError(t('prestations.toastLoadError'))
  } finally {
    loading.value = false
  }
}

async function changeStatus(prestation, status) {
  try {
    const updated = await prestationService.updateStatus(prestation.id, status)
    prestation.status = updated.status
    toast.showSuccess(t('prestations.toastStatusUpdated'))
  } catch {
    toast.showError(t('common.actionFailed'))
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
    toast.showSuccess(t('prestations.toastDeleted'))
    deleteConfirm.show = false
  } catch {
    toast.showError(t('common.actionFailed'))
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
