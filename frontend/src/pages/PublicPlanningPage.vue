<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[900px] mx-auto px-4 sm:px-6 py-10">

      <!-- Header -->
      <div class="mb-8">
        <button @click="router.push('/profil')"
          class="text-sm text-[#40617f] hover:text-[#006d35] flex items-center gap-1.5 mb-4 transition">
          <ArrowLeftIcon class="w-4 h-4" />
          {{ t('common.back') }}
        </button>
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">
          {{ t('public.planning.title') }}
        </h1>
        <p class="text-[#40617f] text-sm mt-1">{{ t('public.planning.subtitle') }}</p>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else>
        <!-- Filter tabs -->
        <div class="flex gap-2 mb-6">
          <button v-for="f in filters" :key="f.value"
            @click="activeFilter = f.value"
            :class="[
              'px-4 py-1.5 rounded-full text-sm font-medium transition',
              activeFilter === f.value
                ? 'bg-[#006d35] text-white'
                : 'bg-white text-[#40617f] border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35]'
            ]">
            {{ f.label }}
          </button>
        </div>

        <!-- Events list -->
        <div v-if="filteredItems.length" class="space-y-4 mb-10">
          <div v-for="item in filteredItems" :key="`${item.type}-${item.id}`"
            class="bg-white rounded-2xl p-5 flex gap-4 items-start shadow-sm hover:shadow-md transition">

            <!-- Date badge -->
            <div class="shrink-0 w-14 h-14 rounded-xl flex flex-col items-center justify-center text-white font-jakarta font-bold"
              :class="item.past ? 'bg-[#9ca3af]' : 'bg-gradient-to-br from-[#006d35] to-[#1b8848]'">
              <span class="text-xs uppercase leading-none">{{ monthShort(item.startDate) }}</span>
              <span class="text-xl leading-tight">{{ dayNum(item.startDate) }}</span>
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap mb-1">
                <span :class="[
                  'text-xs font-semibold px-2 py-0.5 rounded-full',
                  item.type === 'event'
                    ? 'bg-[#edf4ff] text-[#1a73e8]'
                    : 'bg-[#fef3c7] text-[#92400e]'
                ]">
                  {{ item.type === 'event' ? t('public.planning.typeEvent') : t('public.planning.typePrestation') }}
                </span>
                <span v-if="item.organizer"
                  class="text-xs font-bold px-2 py-0.5 rounded-full bg-[#d1fae5] text-[#065f46] flex items-center gap-1">
                  <CheckBadgeIcon class="w-3 h-3" />
                  {{ t('public.planning.organizerBadge') }}
                </span>
                <span v-if="item.type === 'prestation' && item.status"
                  :class="['text-xs px-2 py-0.5 rounded-full font-medium', statusClass(item.status)]">
                  {{ statusLabel(item.status) }}
                </span>
                <span v-if="item.past"
                  class="text-xs text-[#9ca3af] ml-auto shrink-0">
                  {{ t('public.planning.filterPast') }}
                </span>
              </div>
              <p class="font-jakarta font-bold text-[#001d32] text-base truncate">{{ item.title }}</p>
              <p v-if="item.location" class="text-[#40617f] text-sm flex items-center gap-1 mt-0.5">
                <MapPinIcon class="w-3.5 h-3.5 shrink-0" />
                {{ item.location }}
              </p>
              <p v-if="item.type === 'event' && item.endDate" class="text-[#40617f] text-xs mt-1">
                {{ formatDateRange(item.startDate, item.endDate) }}
              </p>
              <p v-else class="text-[#40617f] text-xs mt-1">
                {{ formatDate(item.startDate) }}
              </p>
            </div>

            <!-- Action -->
            <div class="shrink-0">
              <button v-if="item.type === 'event'"
                @click="router.push(`/evenements/${item.eventId}`)"
                class="text-xs text-[#006d35] font-medium hover:underline">
                {{ t('public.planning.seeDetail') }}
              </button>
              <button v-else
                @click="router.push(`/profil/reservations/${item.id}`)"
                class="text-xs text-[#006d35] font-medium hover:underline">
                {{ t('public.planning.seeDetail') }}
              </button>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else class="text-center py-16 text-[#40617f]">
          <CalendarDaysIcon class="w-12 h-12 mx-auto mb-3 text-[#d1d5db]" />
          <p class="text-sm">{{ emptyMessage }}</p>
        </div>

        <!-- iCal subscription card -->
        <div class="bg-white rounded-2xl p-6 shadow-sm mt-4">
          <div class="flex items-center gap-3 mb-3">
            <div class="w-10 h-10 rounded-xl bg-[#edf4ff] flex items-center justify-center">
              <CalendarDaysIcon class="w-5 h-5 text-[#1a73e8]" />
            </div>
            <div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.planning.calendarTitle') }}</h2>
              <p class="text-[#40617f] text-sm">{{ t('public.planning.calendarDesc') }}</p>
            </div>
          </div>

          <!-- Token loading -->
          <div v-if="tokenLoading" class="flex items-center gap-2 text-sm text-[#40617f] py-2">
            <div class="w-4 h-4 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin" />
            {{ t('common.loading') }}
          </div>

          <template v-else-if="calendarToken">
            <!-- URL field -->
            <div class="mt-4">
              <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">
                {{ t('public.planning.calendarUrlLabel') }}
              </label>
              <div class="flex items-center gap-2">
                <input readonly :value="httpsUrl"
                  class="flex-1 text-xs bg-[#f7f9ff] border border-[#e5e7eb] rounded-lg px-3 py-2 font-mono text-[#001d32] select-all truncate" />
                <button @click="copyUrl"
                  :class="[
                    'shrink-0 px-3 py-2 rounded-lg text-xs font-semibold transition',
                    copied ? 'bg-[#006d35] text-white' : 'bg-[#edf4ff] text-[#1a73e8] hover:bg-[#d0e8ff]'
                  ]">
                  {{ copied ? t('public.planning.calendarCopied') : t('public.planning.calendarCopy') }}
                </button>
              </div>
            </div>

            <!-- Apple Calendar button -->
            <a :href="webcalUrl"
              class="mt-3 inline-flex items-center gap-2 px-4 py-2 bg-[#001d32] text-white text-sm font-semibold rounded-xl hover:bg-[#003060] transition">
              <svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
                <path d="M19 3h-1V1h-2v2H8V1H6v2H5a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2V5a2 2 0 00-2-2zm0 16H5V9h14v10zm0-12H5V5h14v2z"/>
              </svg>
              {{ t('public.planning.calendarOpenApple') }}
            </a>

            <!-- Instructions -->
            <details class="mt-4">
              <summary class="text-sm font-semibold text-[#40617f] cursor-pointer hover:text-[#006d35] transition">
                {{ t('public.planning.calendarHowTitle') }}
              </summary>
              <ul class="mt-3 space-y-2 text-sm text-[#40617f] list-disc list-inside">
                <li>{{ t('public.planning.calendarHowApple') }}</li>
                <li>{{ t('public.planning.calendarHowGoogle') }}</li>
                <li>{{ t('public.planning.calendarHowOutlook') }}</li>
              </ul>
            </details>

            <!-- Regenerate -->
            <button @click="regenToken"
              class="mt-4 text-xs text-[#ef4444] hover:underline">
              {{ t('public.planning.calendarRegen') }}
            </button>
          </template>
        </div>
      </template>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ArrowLeftIcon,
  MapPinIcon,
  CalendarDaysIcon,
  CheckBadgeIcon,
} from '@heroicons/vue/24/outline'

const router = useRouter()
const { t, locale } = useI18n()
const userAuth = useUserAuthStore()

const loading = ref(true)
const tokenLoading = ref(true)
const items = ref([])
const calendarToken = ref(null)
const copied = ref(false)

const filters = computed(() => [
  { value: 'upcoming', label: t('public.planning.filterUpcoming') },
  { value: 'all', label: t('public.planning.filterAll') },
  { value: 'past', label: t('public.planning.filterPast') },
])
const activeFilter = ref('upcoming')

const filteredItems = computed(() => {
  if (activeFilter.value === 'upcoming') return items.value.filter(i => !i.past)
  if (activeFilter.value === 'past') return items.value.filter(i => i.past)
  return items.value
})

const emptyMessage = computed(() => {
  if (activeFilter.value === 'upcoming') return t('public.planning.emptyUpcoming')
  if (activeFilter.value === 'past') return t('public.planning.emptyPast')
  return t('public.planning.emptyAll')
})

const httpsUrl = computed(() => {
  if (!calendarToken.value) return ''
  const base = `${window.location.protocol}//${window.location.host}`
  return `${base}/api/public/v1/calendar.ics?token=${calendarToken.value}`
})

const webcalUrl = computed(() => {
  if (!calendarToken.value) return '#'
  const base = `${window.location.host}`
  return `webcal://${base}/api/public/v1/calendar.ics?token=${calendarToken.value}`
})

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatDateRange(startIso, endIso) {
  if (!startIso) return ''
  const start = new Date(startIso)
  const end = endIso ? new Date(endIso) : null
  const opts = { hour: '2-digit', minute: '2-digit' }
  if (!end || start.toDateString() === end.toDateString()) {
    return `${start.toLocaleTimeString(localeCode.value, opts)}${end ? ' – ' + end.toLocaleTimeString(localeCode.value, opts) : ''}`
  }
  return `${formatDate(startIso)} – ${formatDate(endIso)}`
}

function monthShort(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { month: 'short' }).replace('.', '')
}

function dayNum(iso) {
  if (!iso) return ''
  return new Date(iso).getDate()
}

function statusLabel(s) {
  const map = {
    pending: t('public.planning.statusPending'),
    paid: t('public.planning.statusPaid'),
    failed: t('public.planning.statusFailed'),
    quote: t('public.planning.statusQuote'),
  }
  return map[s] || s
}

function statusClass(s) {
  const map = {
    pending: 'bg-yellow-50 text-yellow-700',
    paid: 'bg-green-50 text-green-700',
    failed: 'bg-red-50 text-red-700',
    quote: 'bg-blue-50 text-blue-700',
  }
  return map[s] || 'bg-gray-100 text-gray-600'
}

async function copyUrl() {
  try {
    await navigator.clipboard.writeText(httpsUrl.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 2500)
  } catch {}
}

async function regenToken() {
  if (!confirm(t('public.planning.calendarRegenConfirm'))) return
  try {
    const res = await fetch('/api/v1/calendar/token/regenerate', {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      const data = await res.json()
      calendarToken.value = data.calendar_token
    }
  } catch {}
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/planning')
    return
  }

  const headers = { Authorization: `Bearer ${userAuth.token}`, Accept: 'application/json' }

  const [profileRes, tokenRes] = await Promise.all([
    fetch('/api/v1/profile', { headers }).catch(() => null),
    fetch('/api/v1/calendar/token', { headers }).catch(() => null),
  ])

  if (profileRes?.ok) {
    const data = await profileRes.json()
    const now = new Date()

    const events = (data.reservations || [])
      .filter(r => r.type === 'event')
      .map(r => ({
        id: r.id,
        eventId: r.event_id,
        type: 'event',
        title: r.title,
        startDate: r.start_date,
        endDate: null,
        location: r.location || null,
        past: r.past,
        organizer: r.organizer || false,
      }))

    const prestations = (data.reservations || [])
      .filter(r => r.type === 'prestation')
      .map(r => ({
        id: r.id,
        type: 'prestation',
        title: r.title,
        startDate: r.start_date,
        endDate: null,
        location: null,
        status: r.status,
        past: r.past,
      }))

    items.value = [...events, ...prestations].sort((a, b) => {
      const da = a.startDate ? new Date(a.startDate) : new Date(0)
      const db = b.startDate ? new Date(b.startDate) : new Date(0)
      return db - da
    })
  }

  if (tokenRes?.ok) {
    const data = await tokenRes.json()
    calendarToken.value = data.calendar_token
  }

  loading.value = false
  tokenLoading.value = false
})
</script>
