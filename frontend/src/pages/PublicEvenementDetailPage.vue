<template>
  <div class="bg-[#f7f9ff] pb-16">
    <div class="max-w-[1280px] mx-auto px-6 py-10">

      <div class="flex items-center gap-2 text-sm text-[#40617f] mb-8">
        <RouterLink to="/evenements" class="hover:text-[#006d35] transition">{{ t('public.eventDetail.breadcrumbEvents') }}</RouterLink>
        <ChevronRightIcon class="w-3.5 h-3.5" />
        <span class="text-[#001d32] font-medium truncate max-w-[300px]">{{ event?.title || '…' }}</span>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-32">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <!-- Not found -->
      <div v-else-if="!event" class="flex flex-col items-center justify-center py-32 text-center">
        <CalendarDaysIcon class="w-16 h-16 text-[#40617f]/30 mb-4" />
        <p class="font-jakarta font-bold text-[#001d32] text-xl mb-2">{{ t('public.eventDetail.notFound') }}</p>
        <RouterLink to="/evenements" class="text-[#006d35] font-semibold hover:underline text-sm mt-2">
          {{ t('public.eventDetail.back') }}
        </RouterLink>
      </div>

      <template v-else>
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">

          <!-- Main content -->
          <div class="lg:col-span-8 flex flex-col gap-8">

            <!-- Hero image -->
            <div class="relative h-[360px] sm:h-[440px] rounded-[32px] overflow-hidden flex items-center justify-center shadow-sm"
              :style="`background: ${cardGradient(event.id)};`">
              <CalendarDaysIcon class="w-40 h-40 text-white/20" />
              <div class="absolute top-6 left-6 flex gap-2 flex-wrap">
                <span v-if="event.category"
                  class="bg-white/90 backdrop-blur-sm font-bold text-sm px-4 py-2 rounded-full shadow-sm"
                  :class="categoryColor(event.category.name)">
                  {{ event.category.name }}
                </span>
                <span v-if="isPast"
                  class="bg-gray-700/80 text-white text-sm font-bold px-4 py-2 rounded-full">
                  {{ t('public.eventDetail.ctaPast') }}
                </span>
              </div>
              <div class="absolute top-6 right-6">
                <span v-if="isFull && !isPast"
                  class="bg-red-500 text-white text-sm font-bold px-4 py-2 rounded-full">
                  {{ t('public.eventDetail.badgeFull') }}
                </span>
                <span v-else-if="!isPast && spotsLeft !== null"
                  class="bg-[#006d35] text-white text-sm font-bold px-4 py-2 rounded-full">
                  {{ spotsLeft }} {{ spotsLeft === 1 ? 'place' : 'places' }}
                </span>
              </div>
            </div>

            <!-- Title + meta -->
            <div class="flex flex-col gap-3">
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl leading-tight tracking-tight">
                {{ event.title }}
              </h1>
              <div class="flex flex-wrap gap-5 text-[#40617f] text-base mt-1">
                <div class="flex items-center gap-2">
                  <CalendarDaysIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ formatDate(event.start_date) }}</span>
                </div>
                <div v-if="sameDay" class="flex items-center gap-2">
                  <ClockIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ formatTime(event.start_date) }} – {{ formatTime(event.end_date) }}</span>
                </div>
                <div v-else class="flex items-center gap-2">
                  <ClockIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ formatDate(event.start_date) }} – {{ formatDate(event.end_date) }}</span>
                </div>
                <div v-if="event.location" class="flex items-center gap-2">
                  <MapPinIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ event.location }}</span>
                </div>
              </div>
            </div>

            <!-- Description -->
            <div v-if="event.description" class="bg-white rounded-[32px] p-10 shadow-sm flex flex-col gap-4">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.eventDetail.aboutTitle') }}</h2>
              <p class="text-[#40617f] text-lg leading-relaxed whitespace-pre-wrap">{{ event.description }}</p>
            </div>

            <!-- Organizer -->
            <div v-if="event.creator" class="bg-[#edf4ff] rounded-[32px] p-8 flex gap-6 items-start">
              <div class="w-20 h-20 rounded-2xl overflow-hidden bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] flex items-center justify-center shrink-0 shadow-md">
                <img v-if="event.creator.avatar_url" :src="event.creator.avatar_url" class="w-full h-full object-cover" alt="" />
                <UserCircleIcon v-else class="w-12 h-12 text-[#006d35]/40" />
              </div>
              <div class="flex flex-col gap-1">
                <p class="text-[#006d35] text-xs font-bold uppercase tracking-widest">{{ t('public.eventDetail.organizerTag') }}</p>
                <h3 class="font-jakarta font-bold text-[#001d32] text-xl">
                  {{ event.creator.first_name }} {{ event.creator.last_name }}
                </h3>
              </div>
            </div>

          </div>

          <!-- Sidebar -->
          <div class="lg:col-span-4 flex flex-col gap-6">

            <!-- Registration card -->
            <div class="bg-white border border-[#d8eaff] rounded-[40px] p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] relative overflow-hidden sticky top-6">
              <div class="absolute -right-8 -top-8 w-32 h-32 rounded-full bg-[rgba(0,109,53,0.05)]" />

              <!-- Price -->
              <div class="flex items-baseline justify-between mb-2">
                <span class="font-jakarta font-extrabold text-[#001d32] text-4xl">{{ t('public.eventDetail.priceFree') }}</span>
              </div>
              <p class="text-[#40617f] text-sm mb-6">{{ t('public.eventDetail.freeLabel') }}</p>

              <!-- Spots info -->
              <div v-if="!isPast" class="mb-6">
                <div v-if="event.max_participants" class="flex items-center gap-2">
                  <div class="flex-1 bg-gray-100 rounded-full h-2 overflow-hidden">
                    <div class="h-full bg-[#006d35] rounded-full transition-all"
                      :style="`width: ${Math.min(100, (event.registrations_count / event.max_participants) * 100)}%`" />
                  </div>
                  <span class="text-sm font-semibold text-[#40617f] shrink-0">
                    {{ event.registrations_count }}/{{ event.max_participants }}
                  </span>
                </div>
                <p v-else class="text-sm text-[#40617f]">{{ t('public.eventDetail.spotsUnlimited') }}</p>
              </div>

              <!-- Registered badge -->
              <div v-if="isRegistered && !isPast" class="flex items-center gap-2 bg-[#d1fae5] text-[#065f46] px-4 py-2 rounded-xl mb-4 text-sm font-semibold">
                <CheckCircleIcon class="w-4 h-4 shrink-0" />
                {{ t('public.eventDetail.badgeRegistered') }}
              </div>

              <!-- CTA button -->
              <button v-if="isPast" disabled
                class="w-full py-4 rounded-2xl text-[#40617f] font-bold text-base bg-gray-100 cursor-default">
                {{ t('public.eventDetail.ctaPast') }}
              </button>
              <button v-else-if="!userAuth.isLoggedIn"
                @click="router.push(`/connexion?redirect=/evenements/${event.id}`)"
                class="w-full py-4 rounded-2xl text-white font-bold text-base flex items-center justify-center gap-2 transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,109,53,0.2)]"
                style="background: linear-gradient(168deg, #006d35 0%, #1b8848 100%);">
                {{ t('public.eventDetail.ctaLogin') }}
                <ArrowRightIcon class="w-4 h-4" />
              </button>
              <button v-else-if="isRegistered"
                @click="unregister"
                :disabled="actionLoading"
                class="w-full py-4 rounded-2xl font-bold text-base transition border-2 border-red-300 text-red-600 hover:bg-red-50 disabled:opacity-60">
                {{ actionLoading ? '…' : t('public.eventDetail.ctaUnregister') }}
              </button>
              <button v-else-if="isFull" disabled
                class="w-full py-4 rounded-2xl text-[#40617f] font-bold text-base bg-gray-100 cursor-default">
                {{ t('public.eventDetail.ctaFull') }}
              </button>
              <button v-else
                @click="register"
                :disabled="actionLoading"
                class="w-full py-4 rounded-2xl text-white font-bold text-base flex items-center justify-center gap-2 transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,109,53,0.2)] disabled:opacity-60"
                style="background: linear-gradient(168deg, #006d35 0%, #1b8848 100%);">
                {{ actionLoading ? t('public.eventDetail.registering') : t('public.eventDetail.ctaReserve') }}
                <ArrowRightIcon v-if="!actionLoading" class="w-4 h-4" />
              </button>

              <!-- Feedback message -->
              <p v-if="feedback" class="mt-3 text-center text-sm font-medium"
                :class="feedbackError ? 'text-red-600' : 'text-[#006d35]'">
                {{ feedback }}
              </p>
            </div>

            <!-- Location card -->
            <div v-if="event.location" class="bg-white border border-[#d8eaff] rounded-[32px] p-6 flex items-center justify-between gap-4">
              <div>
                <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ event.location }}</p>
              </div>
              <div class="w-12 h-12 bg-[#d8eaff] rounded-xl flex items-center justify-center shrink-0">
                <MapPinIcon class="w-5 h-5 text-[#40617f]" />
              </div>
            </div>

          </div>
        </div>
      </template>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ChevronRightIcon,
  CalendarDaysIcon,
  ClockIcon,
  MapPinIcon,
  UserCircleIcon,
  ArrowRightIcon,
  CheckCircleIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const userAuth = useUserAuthStore()

const loading = ref(true)
const event = ref(null)
const isRegistered = ref(false)
const actionLoading = ref(false)
const feedback = ref('')
const feedbackError = ref(false)

const GRADIENTS = [
  ['#dbeafe', '#bfdbfe'],
  ['#f3e8ff', '#e9d5ff'],
  ['#d1fae5', '#a7f3d0'],
  ['#fef3c7', '#fde68a'],
  ['#fee2e2', '#fecaca'],
  ['#e0e7ff', '#c7d2fe'],
]

function cardGradient(id) {
  const [from, to] = GRADIENTS[id % GRADIENTS.length]
  return `linear-gradient(135deg, ${from}, ${to})`
}

function categoryColor(name = '') {
  const n = name.toLowerCase()
  if (n.includes('atelier') || n.includes('workshop')) return 'text-[#1d4ed8]'
  if (n.includes('formation') || n.includes('training')) return 'text-[#065f46]'
  if (n.includes('conf')) return 'text-[#7c3aed]'
  return 'text-[#40617f]'
}

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatTime(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleTimeString(localeCode.value, { hour: '2-digit', minute: '2-digit' })
}

const sameDay = computed(() => {
  if (!event.value) return false
  const s = new Date(event.value.start_date)
  const e = new Date(event.value.end_date)
  return s.toDateString() === e.toDateString()
})

const isPast = computed(() => {
  if (!event.value) return false
  return new Date(event.value.end_date) < new Date()
})

const isFull = computed(() => {
  if (!event.value?.max_participants) return false
  return event.value.registrations_count >= event.value.max_participants
})

const spotsLeft = computed(() => {
  if (!event.value?.max_participants) return null
  return Math.max(0, event.value.max_participants - event.value.registrations_count)
})

function showFeedback(msg, isError = false) {
  feedback.value = msg
  feedbackError.value = isError
  setTimeout(() => { feedback.value = '' }, 4000)
}

async function register() {
  actionLoading.value = true
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/register`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      isRegistered.value = true
      event.value.registrations_count++
      showFeedback(t('public.eventDetail.registerSuccess'))
    } else {
      const data = await res.json().catch(() => ({}))
      showFeedback(data.message || t('public.eventDetail.registerError'), true)
    }
  } catch {
    showFeedback(t('public.eventDetail.registerError'), true)
  }
  actionLoading.value = false
}

async function unregister() {
  if (!confirm(t('public.eventDetail.unregisterConfirm'))) return
  actionLoading.value = true
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/register`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      isRegistered.value = false
      event.value.registrations_count = Math.max(0, event.value.registrations_count - 1)
      showFeedback(t('public.eventDetail.unregisterSuccess'))
    }
  } catch {}
  actionLoading.value = false
}

onMounted(async () => {
  const id = route.params.id
  try {
    const res = await fetch(`/api/public/v1/events/${id}`)
    if (res.ok) {
      event.value = await res.json()
    }
  } catch {}

  if (userAuth.isLoggedIn && event.value) {
    try {
      const res = await fetch(`/api/v1/events/${id}/registration`, {
        headers: { Authorization: `Bearer ${userAuth.token}` },
      })
      if (res.ok) {
        const data = await res.json()
        isRegistered.value = data.registered
      }
    } catch {}
  }

  loading.value = false
})
</script>
