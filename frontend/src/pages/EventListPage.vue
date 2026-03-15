<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('events.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ t('events.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/events/create" class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90" style="background-color:#1B8848;">
        {{ t('events.createBtn') }}
      </RouterLink>
    </div>

    <div class="card p-4 flex flex-wrap gap-3 items-center">
      <select v-model="filters.status" @change="fetchEvents" class="input w-44">
        <option value="">{{ t('common.allStatuses') }}</option>
        <option value="draft">{{ t('events.statusDraft') }}</option>
        <option value="published">{{ t('events.statusPublished') }}</option>
        <option value="cancelled">{{ t('events.statusCancelled') }}</option>
      </select>
    </div>

    <div class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
        <span class="font-medium text-sm text-gray-700">{{ t('events.title') }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 font-semibold">{{ meta.total ?? 0 }} {{ t('common.total') }}</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('events.colTitle') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('events.colLocation') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('events.colDates') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('events.colParticipants') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('events.colStatus') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">{{ t('events.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-40"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-32"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-12"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!events.length">
              <td :colspan="6" class="px-4 py-16 text-center">
                <div class="text-4xl mb-2">🔍</div>
                <p class="text-gray-500 font-medium">{{ t('common.noResults') }}</p>
                <p class="text-gray-400 text-sm mt-1">{{ t('common.noResultsHint') }}</p>
              </td>
            </tr>
            <tr v-else v-for="e in events" :key="e.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900 max-w-xs truncate">{{ e.title }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ e.location || '—' }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">
                <span class="font-medium text-gray-700">{{ formatDateShort(e.start_date) }}</span>
                <span class="text-xs text-gray-400 ml-1">{{ formatTime(e.start_date) }}</span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ e.max_participants ?? '∞' }}</td>
              <td class="px-4 py-3">
                <AppBadge :label="e.status" />
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <button
                    v-if="e.status === 'draft'"
                    @click="changeStatus(e, 'published')"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#DCFCE7; color:#15803d;"
                  >
                    {{ t('events.actionPublish') }}
                  </button>
                  <button
                    v-if="e.status === 'published'"
                    @click="changeStatus(e, 'cancelled')"
                    class="text-xs px-3 py-1 rounded-full font-medium transition"
                    style="background:#FFF7ED; color:#c2410c;"
                  >
                    {{ t('events.actionCancel') }}
                  </button>
                  <RouterLink :to="`/admin/events/${e.id}/edit`" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#DBEAFE; color:#1d4ed8;">
                    {{ t('common.edit') }}
                  </RouterLink>
                  <button @click="openDelete(e)" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#FEE2E2; color:#dc2626;">
                    {{ t('common.delete') }}
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
      :title="t('events.confirmDeleteTitle')"
      :message="t('events.confirmDeleteMsg', { name: deleteConfirm.item?.title })"
      :confirm-label="t('common.delete')"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { eventService } from '@/services/eventService'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'

const { t, locale } = useI18n()
const dateLocale = computed(() => locale.value === 'en' ? 'en-GB' : 'fr-FR')
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
    toast.showError(t('events.toastLoadError'))
  } finally {
    loading.value = false
  }
}

async function changeStatus(event, status) {
  try {
    const updated = await eventService.updateStatus(event.id, status)
    event.status = updated.status
    toast.showSuccess(t('events.toastStatusUpdated'))
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
    await eventService.remove(deleteConfirm.item.id)
    events.value = events.value.filter((e) => e.id !== deleteConfirm.item.id)
    toast.showSuccess(t('events.toastDeleted'))
    deleteConfirm.show = false
  } catch {
    toast.showError(t('common.actionFailed'))
  } finally {
    deleteConfirm.loading = false
  }
}

function formatDateShort(iso) {
  return new Date(iso).toLocaleDateString(dateLocale.value, { day: '2-digit', month: 'short', year: 'numeric' })
}

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString(dateLocale.value, { hour: '2-digit', minute: '2-digit' })
}

onMounted(fetchEvents)
</script>
