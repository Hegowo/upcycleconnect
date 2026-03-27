<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between gap-3">
      <div class="min-w-0">
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32] truncate">{{ t('events.listTitle') }}</h2>
        <p class="text-sm text-[#40617f] mt-0.5 hidden sm:block">{{ t('events.listSubtitle') }}</p>
      </div>
      <RouterLink
        to="/admin/events/create"
        class="p-2 sm:px-4 sm:py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2 shrink-0"
        style="background-color:#006d35;"
      >
        <PlusIcon class="w-4 h-4" />
        <span class="hidden sm:inline">{{ t('events.planEvent') }}</span>
      </RouterLink>
    </div>

    <!-- Next event + health cards -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Next event -->
      <div class="lg:col-span-2 bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-[#001d32]">{{ t('events.nextWorkshop') }}</h3>
          <span class="text-xs text-gray-400">{{ t('events.comingSoon') }}</span>
        </div>
        <div v-if="nextEvent">
          <div class="flex items-start gap-4">
            <div class="w-14 h-14 rounded-xl flex items-center justify-center shrink-0" style="background:#dcfce7;">
              <CalendarIcon class="w-7 h-7" style="color:#006d35;" />
            </div>
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1 flex-wrap">
                <h4 class="text-base font-bold text-[#001d32]">{{ nextEvent.title }}</h4>
                <span class="text-xs font-medium px-2 py-0.5 rounded-full bg-[#dcfce7] text-[#166534]">
                  {{ nextEvent.status }}
                </span>
              </div>
              <div class="flex items-center gap-4 text-xs text-gray-500 mb-3 flex-wrap">
                <span class="flex items-center gap-1">
                  <MapPinIcon class="w-3.5 h-3.5" />
                  {{ nextEvent.location || t('events.online') }}
                </span>
                <span class="flex items-center gap-1">
                  <ClockIcon class="w-3.5 h-3.5" />
                  {{ formatDateTime(nextEvent.start_date) }}
                </span>
                <span class="flex items-center gap-1">
                  <UsersIcon class="w-3.5 h-3.5" />
                  {{ nextEvent.max_participants ?? '∞' }} {{ t('events.places') }}
                </span>
              </div>
              <div v-if="nextEvent.max_participants">
                <div class="flex justify-between text-xs text-gray-400 mb-1">
                  <span>{{ t('events.capacity') }}</span>
                  <span>{{ nextEvent.registrations_count ?? 0 }} / {{ nextEvent.max_participants }}</span>
                </div>
                <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                  <div class="h-full rounded-full" :style="{ width: capacityPct(nextEvent) + '%', background: '#006d35' }"></div>
                </div>
              </div>
              <div v-else class="text-xs text-gray-400">{{ t('events.noLimit') }}</div>
            </div>
          </div>
        </div>
        <div v-else class="py-6 text-center text-gray-400 text-sm">
          {{ t('events.noEvent') }}
        </div>
      </div>

      <!-- Event health -->
      <div class="bg-white rounded-2xl p-5 border border-[#f1f5f9] shadow-sm">
        <h3 class="font-semibold text-[#001d32] mb-4">{{ t('events.healthTitle') }}</h3>
        <div class="space-y-4">
          <div>
            <div class="flex justify-between text-xs text-gray-500 mb-1">
              <span>{{ t('events.fillRate') }}</span>
              <span class="font-bold text-[#001d32]">{{ globalFillRate }}%</span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div class="h-full rounded-full" :style="{ width: globalFillRate + '%', background: '#006d35' }"></div>
            </div>
          </div>
          <div class="flex items-center justify-between py-2 border-t border-[#f8fafc]">
            <span class="text-xs text-gray-500">{{ t('events.totalRegistrations') }}</span>
            <span class="text-sm font-bold text-[#001d32]">{{ totalRegistrations }}</span>
          </div>
          <div class="flex items-center justify-between py-2 border-t border-[#f8fafc]">
            <span class="text-xs text-gray-500">{{ t('events.publishedSlots') }}</span>
            <span class="text-sm font-bold text-[#001d32]">{{ events.filter(e => e.status === 'published').length }}</span>
          </div>
          <div class="flex items-center justify-between py-2 border-t border-[#f8fafc]">
            <span class="text-xs text-gray-500">{{ t('events.totalEvents') }}</span>
            <span class="text-sm font-bold text-[#001d32]">{{ meta.total ?? 0 }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-2xl p-4 border border-[#f1f5f9] shadow-sm flex flex-wrap gap-3 items-center">
      <select v-model="filters.status" @change="fetchEvents" class="text-sm border border-[#e5e7eb] rounded-lg px-3 py-2 bg-[#f8fafc] text-gray-600 focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
        <option value="">{{ t('events.allStatuses') }}</option>
        <option value="draft">{{ t('events.statusDraft') }}</option>
        <option value="published">{{ t('events.statusPublished') }}</option>
        <option value="cancelled">{{ t('events.statusCancelled') }}</option>
      </select>
    </div>

    <!-- Mobile event cards (lg:hidden) -->
<div class="lg:hidden">
  <!-- Loading -->
  <div v-if="loading" class="space-y-3">
    <div v-for="n in 5" :key="n" class="bg-white rounded-2xl border border-[#f1f5f9] h-28 animate-pulse"></div>
  </div>
  <!-- Empty -->
  <div v-else-if="!events.length" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm py-12 text-center">
    <CalendarIcon class="w-8 h-8 text-gray-300 mx-auto mb-2" />
    <p class="text-gray-500 font-medium text-sm">{{ t('events.noFound') }}</p>
  </div>
  <!-- Cards -->
  <div v-else class="space-y-3">
    <div v-for="e in events" :key="e.id" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-4">
      <!-- Top row: title + status badge -->
      <div class="flex items-start gap-2 mb-3">
        <div class="flex-1 min-w-0">
          <h4 class="text-sm font-semibold text-[#001d32] leading-tight truncate">{{ e.title }}</h4>
          <div class="flex items-center gap-3 mt-1 text-xs text-gray-400">
            <span class="flex items-center gap-1">
              <MapPinIcon class="w-3 h-3" />
              {{ e.location || t('events.online') }}
            </span>
          </div>
        </div>
        <span class="text-xs font-semibold px-2 py-0.5 rounded-full shrink-0" :class="eventStatusBadge(e.status)">
          {{ eventStatusText(e.status) }}
        </span>
      </div>
      <!-- Info row -->
      <div class="flex items-center gap-3 text-xs text-gray-500 mb-3">
        <span class="flex items-center gap-1">
          <ClockIcon class="w-3 h-3" />
          {{ formatDateShort(e.start_date) }} {{ formatTime(e.start_date) }}
        </span>
        <span class="flex items-center gap-1">
          <UsersIcon class="w-3 h-3" />
          {{ e.max_participants ?? '∞' }} {{ t('events.places') }}
        </span>
      </div>
      <!-- Actions -->
      <div class="flex items-center gap-2 pt-2 border-t border-[#f8fafc]">
        <button
          v-if="e.status === 'draft'"
          @click="changeStatus(e, 'published')"
          class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
          style="background:#dcfce7; color:#166534;"
        >
          {{ t('events.actionPublish') }}
        </button>
        <button
          v-if="e.status === 'published'"
          @click="changeStatus(e, 'cancelled')"
          class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
          style="background:#fff7ed; color:#c2410c;"
        >
          {{ t('events.actionCancel') }}
        </button>
        <RouterLink
          :to="`/admin/events/${e.id}/edit`"
          class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
          style="background:#f1f5f9; color:#475569;"
        >
          {{ t('common.edit') }}
        </RouterLink>
        <button @click="openDelete(e)" class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition">
          <TrashIcon class="w-4 h-4" />
        </button>
      </div>
    </div>
    <!-- Mobile pagination -->
    <div class="flex justify-center pt-2">
      <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchEvents" />
    </div>
  </div>
</div>

    <!-- Events table -->
    <div class="hidden lg:block bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <span class="font-semibold text-sm text-[#001d32]">{{ t('events.allEvents') }}</span>
        <span class="text-xs px-2.5 py-1 rounded-full font-semibold" style="background:#dcfce7; color:#166534;">
          {{ meta.total ?? 0 }} total
        </span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full">
          <thead>
            <tr class="bg-[#f8fafc]">
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('events.colTitle') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('events.colLocation') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('events.colDate') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('events.colParticipantsHeader') }}</th>
              <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('events.colStatus') }}</th>
              <th class="px-6 py-3 text-right text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('events.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#f8fafc]">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-40"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-32"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-12"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-20"></div></td>
              <td class="px-6 py-4"><div class="h-4 bg-gray-100 rounded animate-pulse w-24 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!events.length">
              <td :colspan="6" class="px-6 py-16 text-center">
                <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
                  <CalendarIcon class="w-6 h-6 text-gray-400" />
                </div>
                <p class="text-gray-500 font-medium">{{ t('events.noFound') }}</p>
              </td>
            </tr>
            <tr v-else v-for="e in events" :key="e.id" class="hover:bg-[#f8fafc] transition-colors">
              <td class="px-6 py-4 text-sm font-semibold text-[#001d32] max-w-xs truncate">{{ e.title }}</td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ e.location || '—' }}</td>
              <td class="px-6 py-4">
                <span class="text-sm font-medium text-[#001d32]">{{ formatDateShort(e.start_date) }}</span>
                <span class="text-xs text-gray-400 ml-1">{{ formatTime(e.start_date) }}</span>
              </td>
              <td class="px-6 py-4 text-sm text-gray-500">{{ e.max_participants ?? '∞' }}</td>
              <td class="px-6 py-4">
                <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="eventStatusBadge(e.status)">
                  {{ eventStatusText(e.status) }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex justify-end gap-1.5">
                  <button
                    v-if="e.status === 'draft'"
                    @click="changeStatus(e, 'published')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#dcfce7; color:#166534;"
                  >
                    {{ t('events.actionPublish') }}
                  </button>
                  <button
                    v-if="e.status === 'published'"
                    @click="changeStatus(e, 'cancelled')"
                    class="px-3 py-1 rounded-lg text-xs font-semibold transition"
                    style="background:#fff7ed; color:#c2410c;"
                  >
                    {{ t('events.actionCancel') }}
                  </button>
                  <RouterLink
                    :to="`/admin/events/${e.id}/edit`"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-[#40617f] hover:bg-blue-50 transition"
                    :title="t('common.edit')"
                  >
                    <PencilIcon class="w-4 h-4" />
                  </RouterLink>
                  <button @click="openDelete(e)" class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition" :title="t('events.actionDelete')">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-6 py-4 border-t border-[#f1f5f9]">
        <AppPagination :current-page="meta.current_page" :last-page="meta.last_page" @page-change="fetchEvents" />
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      :title="t('events.confirmDeleteTitle')"
      :message="t('events.confirmDeleteMsg', { name: deleteConfirm.item?.title })"
      :confirm-label="t('events.actionDelete')"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { eventService } from '@/services/eventService'
import { useToast } from '@/utils/useToast'
import AppPagination from '@/components/AppPagination.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'
import {
  CalendarIcon, PlusIcon, PencilIcon, TrashIcon,
  MapPinIcon, ClockIcon, UsersIcon,
} from '@heroicons/vue/24/outline'

const auth  = useAuthStore()
const toast = useToast()
const { t } = useI18n()
const events  = ref([])
const loading = ref(false)
const meta    = ref({ current_page: 1, last_page: 1, total: 0 })
const filters = reactive({ status: '' })
const deleteConfirm = reactive({ show: false, item: null, loading: false })

const nextEvent = computed(() => events.value.find(e => e.status === 'published') || events.value[0] || null)

function capacityPct(event) {
  if (!event.max_participants || event.max_participants <= 0) return 0
  return Math.min(100, Math.round((event.registrations_count ?? 0) / event.max_participants * 100))
}

const totalRegistrations = computed(() =>
  events.value.reduce((sum, e) => sum + (e.registrations_count ?? 0), 0)
)

const globalFillRate = computed(() => {
  const withCap = events.value.filter(e => e.max_participants > 0)
  if (!withCap.length) return 0
  const total = withCap.reduce((s, e) => s + e.max_participants, 0)
  const filled = withCap.reduce((s, e) => s + (e.registrations_count ?? 0), 0)
  return Math.min(100, Math.round(filled / total * 100))
})

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
    toast.showSuccess(t('events.toastStatusUpdatedMsg'))
  } catch {
    toast.showError(t('events.toastErrorOccurred'))
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
    toast.showError(t('events.toastErrorOccurred'))
  } finally {
    deleteConfirm.loading = false
  }
}

function eventStatusBadge(status) {
  if (status === 'published') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'cancelled') return 'bg-[#fee2e2] text-[#991b1b]'
  if (status === 'draft')     return 'bg-[#f1f5f9] text-[#475569]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function eventStatusText(status) {
  const map = {
    published: t('events.statusPublished'),
    draft:     t('events.statusDraft'),
    cancelled: t('events.statusCancelled'),
  }
  return map[status] || status
}
function formatDateShort(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' })
}
function formatDateTime(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}
function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  fetchEvents()
})
</script>
