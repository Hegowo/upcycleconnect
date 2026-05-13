<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-6 py-10">

      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else class="grid grid-cols-12 gap-4 sm:gap-6">

        <div class="col-span-12 lg:col-span-8 bg-white rounded-[24px] p-6 sm:p-8 flex gap-4 sm:gap-8 overflow-hidden relative">
          <div class="relative shrink-0">
            <div class="w-32 h-32 rounded-[24px] overflow-hidden shadow-sm bg-gradient-to-br from-[#006d35] to-[#1b8848] flex items-center justify-center">
              <img v-if="userAuth.user?.avatar_url" :src="userAuth.user.avatar_url" class="w-full h-full object-cover" />
              <span v-else class="text-white font-jakarta font-extrabold text-4xl">{{ userAuth.initials }}</span>
            </div>
            <div class="absolute bottom-0 right-0 w-6 h-6 bg-[#006d35] rounded-full flex items-center justify-center">
              <CheckIcon class="w-3.5 h-3.5 text-white" />
            </div>
          </div>

          <div class="flex-1 flex flex-col gap-2">
            <div class="flex items-center gap-3">
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">{{ userAuth.fullName }}</h1>
              <span class="bg-[rgba(27,136,72,0.1)] text-[#006d35] text-xs font-bold uppercase tracking-wide px-3 py-1 rounded-full">
                {{ memberLabel }}
              </span>
              <button @click="router.push('/profil/modifier')"
                class="ml-auto flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-[#e5e7eb] text-xs font-medium text-[#40617f] hover:bg-gray-50 transition shrink-0">
                <PencilIcon class="w-3.5 h-3.5" /> {{ t('public.profile.edit') }}
              </button>
            </div>
            <div class="flex items-center gap-2 text-[#40617f] text-sm">
              <EnvelopeIcon class="w-3.5 h-3.5" />
              <span>{{ userAuth.user?.email }}</span>
            </div>
            <p class="text-[#3f4a3f] text-base leading-relaxed mt-2 max-w-md text-sm text-gray-500">
              {{ t('public.profile.memberSince', { date: memberSince }) }}
            </p>
          </div>

          <div class="absolute -right-10 -top-10 w-36 h-36 rounded-full bg-[rgba(0,109,53,0.05)]" />
        </div>

        <div class="col-span-12 lg:col-span-4 bg-[#edf4ff] rounded-[24px] p-6 sm:p-8 flex flex-col justify-between">
          <div class="flex items-center justify-between mb-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.profile.scoreTitle') }}</h2>
            <ArrowTrendingUpIcon class="w-5 h-5 text-[#006d35]" />
          </div>

          <div class="mb-4">
            <div class="flex items-baseline gap-1">
              <span class="font-jakarta font-extrabold text-[#006d35] text-5xl">{{ profile.score }}</span>
              <span class="text-[#40617f] text-lg">{{ t('public.profile.scoreOutOf') }}</span>
            </div>
            <div class="mt-3 h-3 bg-[#cee5ff] rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-gradient-to-r from-[#006d35] to-[#1b8848] transition-all duration-700"
                :style="`width: ${profile.score}%;`" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="bg-white rounded-xl p-3">
              <p class="text-[#40617f] text-xs uppercase tracking-wide">{{ t('public.profile.scoreCo2') }}</p>
              <p class="font-semibold text-[#001d32] text-lg mt-0.5">{{ co2Display }} kg</p>
            </div>
            <div class="bg-white rounded-xl p-3">
              <p class="text-[#40617f] text-xs uppercase tracking-wide">{{ t('public.profile.scoreObjects') }}</p>
              <p class="font-semibold text-[#001d32] text-lg mt-0.5">{{ profile.objects_saved }}</p>
            </div>
          </div>
        </div>

        <div class="col-span-12 bg-white rounded-[24px] p-8">
          <div class="flex items-center justify-between mb-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.profile.depositsTitle') }}</h2>
            <span class="text-xs text-[#40617f]">{{ t('public.profile.depositsTotal', { n: profile.deposits?.length || 0 }) }}</span>
          </div>

          <div v-if="!profile.deposits?.length" class="text-center py-12 text-[#40617f] text-sm">
            {{ t('public.profile.depositsEmpty') }}
          </div>
          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            <button
              v-for="dep in depositsToShow"
              :key="dep.id"
              class="flex flex-col text-left focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 rounded-xl"
              @click="router.push(`/confirmation-depot?id=${dep.id}`)"
            >
              <div class="relative rounded-xl overflow-hidden h-40 w-full flex items-center justify-center group"
                :style="`background: ${statusGradient(dep.status)};`"
              >
                <img v-if="dep.photo1" :src="dep.photo1" class="w-full h-full object-cover" />
                <ArchiveBoxIcon v-else class="w-16 h-16 text-white/20" />
                <div class="absolute top-3 left-3">
                  <span class="text-xs font-bold uppercase tracking-tight px-2 py-1 rounded-lg" :class="statusClass(dep.status)">
                    {{ statusLabel(dep.status) }}
                  </span>
                </div>
                <div class="absolute inset-0 bg-black/0 group-hover:bg-black/10 transition rounded-xl" />
              </div>
              <h3 class="font-semibold text-[#001d32] text-sm mt-3 truncate">{{ dep.title }}</h3>
              <p class="text-[#40617f] text-xs mt-0.5">{{ formatDate(dep.created_at) }}</p>
            </button>
          </div>
        </div>

        <div class="col-span-12 lg:col-span-7 bg-[#edf4ff] rounded-[24px] p-6 sm:p-8">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">{{ t('public.profile.reservationsTitle') }}</h2>

          <div v-if="!profile.reservations?.length" class="text-center py-8 text-[#40617f] text-sm">
            {{ t('public.profile.reservationsEmpty') }}
          </div>
          <div v-else class="flex flex-col gap-4">
            <button
              v-for="resa in profile.reservations"
              :key="`${resa.type}-${resa.id}`"
              @click="openReservation(resa)"
              class="bg-white rounded-xl p-4 flex items-center gap-4 text-left w-full hover:shadow-md hover:-translate-y-0.5 transition cursor-pointer focus:outline-none focus:ring-2 focus:ring-[#006d35]/30"
              :class="resa.past ? 'opacity-70 hover:opacity-100' : ''"
            >
              <div
                class="w-14 h-14 rounded-lg flex items-center justify-center shrink-0"
                :class="reservationIconBg(resa)"
              >
                <component :is="reservationIcon(resa)" class="w-6 h-6" :class="reservationIconColor(resa)" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 flex-wrap">
                  <p class="font-semibold text-[#001d32] text-base truncate">{{ resa.title }}</p>
                  <span v-if="resa.type === 'prestation' && resa.status"
                    class="text-[10px] font-bold uppercase tracking-wide px-2 py-0.5 rounded-full"
                    :class="prestaStatusClass(resa.status)"
                  >
                    {{ prestaStatusLabel(resa.status) }}
                  </span>
                </div>
                <p class="text-[#40617f] text-sm">
                  <template v-if="resa.type === 'prestation'">
                    {{ formatDate(resa.start_date) }}
                    <span v-if="resa.amount_cents"> · {{ formatAmount(resa.amount_cents, resa.currency) }}</span>
                  </template>
                  <template v-else>
                    {{ formatDate(resa.start_date) }}{{ resa.location ? ` · ${resa.location}` : '' }}
                  </template>
                </p>
              </div>
              <ChevronRightIcon class="w-4 h-4 text-[#94a3b8] shrink-0" />
            </button>
          </div>
        </div>

        <div class="col-span-12 lg:col-span-5 bg-white rounded-[24px] p-6 sm:p-8 flex flex-col justify-between">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">{{ t('public.profile.badgesTitle') }}</h2>

          <div class="flex gap-4 flex-wrap">
            <div
              v-for="badge in profile.badges"
              :key="badge.key"
              class="relative group"
            >
              <div
                class="w-16 h-16 rounded-full flex items-center justify-center border-2 transition-opacity"
                :class="badge.earned ? '' : 'opacity-30'"
                :style="badge.earned ? `background: ${badgeBg(badge.key)}; border-color: ${badgeBorder(badge.key)};` : 'background:#f1f5f9; border-color:#cbd5e1;'"
              >
                <component :is="badgeIcon(badge.key)" class="w-6 h-6"
                  :style="badge.earned ? `color: ${badgeIconColor(badge.key)};` : 'color:#94a3b8;'"
                />
              </div>
              <div class="absolute -bottom-8 left-1/2 -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-[#001d32] text-white text-[10px] px-2 py-1 rounded whitespace-nowrap z-10 pointer-events-none">
                {{ badge.label }}<span v-if="!badge.earned" class="opacity-60"> — {{ badge.desc }}</span>
              </div>
            </div>
          </div>

          <div class="border-t border-[#edf4ff] pt-6 mt-6 flex items-center justify-between">
            <span class="text-[#40617f] text-xs font-semibold uppercase tracking-wide">
              {{ t('public.profile.badgesProgress', { earned: earnedBadgesCount, total: profile.badges?.length || 0 }) }}
            </span>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  CheckIcon,
  ChevronRightIcon,
  CalendarDaysIcon,
  ArrowTrendingUpIcon,
  EnvelopeIcon,
  ArchiveBoxIcon,
  GlobeAltIcon,
  HeartIcon,
  BoltIcon,
  SparklesIcon,
  TrophyIcon,
  AcademicCapIcon,
  PencilIcon,
  WrenchScrewdriverIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const router   = useRouter()
const userAuth = useUserAuthStore()

const loading = ref(true)
const profile = ref({
  score: 0,
  co2_saved: 0,
  objects_saved: 0,
  deposits: [],
  reservations: [],
  badges: [],
})

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil')
    return
  }
  try {
    const res = await fetch('/api/v1/profile', {
      headers: {
        Authorization: `Bearer ${userAuth.token}`,
        Accept: 'application/json',
      },
    })
    if (res.ok) {
      profile.value = await res.json()
    }
  } catch {}
  loading.value = false
})

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

const memberSince = computed(() => {
  if (!userAuth.user?.created_at) return ''
  return new Date(userAuth.user.created_at).toLocaleDateString(localeCode.value, { year: 'numeric', month: 'long' })
})

const memberLabel = computed(() => {
  const score = profile.value.score
  if (score >= 75) return t('public.profile.memberLabelChampion')
  if (score >= 50) return t('public.profile.memberLabelActive')
  if (score >= 20) return t('public.profile.memberLabelEcoActive')
  return t('public.profile.memberLabelNew')
})

const co2Display = computed(() => {
  const v = profile.value.co2_saved || 0
  return v % 1 === 0 ? v.toFixed(0) : v.toFixed(1)
})

const depositsToShow = computed(() => (profile.value.deposits || []).slice(0, 6))
const earnedBadgesCount = computed(() => (profile.value.badges || []).filter(b => b.earned).length)

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: 'numeric', month: 'long', year: 'numeric' })
}

function statusLabel(s) {
  return {
    pending:   t('public.profile.depositStatusPending'),
    approved:  t('public.profile.depositStatusAccepted'),
    accepted:  t('public.profile.depositStatusAccepted'),
    validated: t('public.profile.depositStatusValidated'),
    rejected:  t('public.profile.depositStatusRejected'),
  }[s] ?? s
}

function statusGradient(s) {
  return {
    pending:   'linear-gradient(135deg, #fef3c7, #fde68a)',
    approved:  'linear-gradient(135deg, #d1fae5, #bbf7d0)',
    accepted:  'linear-gradient(135deg, #d1fae5, #bbf7d0)',
    validated: 'linear-gradient(135deg, #d1fae5, #86efac)',
    rejected:  'linear-gradient(135deg, #fee2e2, #fecaca)',
  }[s] ?? 'linear-gradient(135deg, #e2e8f0, #cbd5e1)'
}

function openReservation(r) {
  if (r.type === 'prestation') {
    router.push(`/profil/reservations/${r.id}`)
  } else if (r.type === 'event' && r.event_id) {
    router.push(`/evenements/${r.event_id}`)
  }
}

function formatAmount(cents, currency) {
  const val = (cents || 0) / 100
  const cur = (currency || 'eur').toUpperCase()
  try {
    return new Intl.NumberFormat(localeCode.value, { style: 'currency', currency: cur }).format(val)
  } catch {
    return `${val.toFixed(2)} €`
  }
}

function reservationIcon(r) {
  return r.type === 'prestation' ? WrenchScrewdriverIcon : CalendarDaysIcon
}
function reservationIconBg(r) {
  if (r.type === 'prestation') {
    if (r.status === 'paid') return 'bg-[rgba(27,136,72,0.1)]'
    if (r.status === 'failed') return 'bg-[rgba(239,68,68,0.1)]'
    if (r.status === 'quote_requested') return 'bg-[rgba(59,130,246,0.1)]'
    return 'bg-[rgba(245,158,11,0.1)]'
  }
  return r.past ? 'bg-[rgba(185,219,254,0.2)]' : 'bg-[rgba(27,136,72,0.1)]'
}
function reservationIconColor(r) {
  if (r.type === 'prestation') {
    if (r.status === 'paid') return 'text-[#006d35]'
    if (r.status === 'failed') return 'text-red-500'
    if (r.status === 'quote_requested') return 'text-blue-500'
    return 'text-amber-500'
  }
  return r.past ? 'text-[#40617f]' : 'text-[#006d35]'
}
function prestaStatusLabel(s) {
  return {
    pending:         t('public.profile.reservationStatusPending'),
    paid:            t('public.profile.reservationStatusPaid'),
    failed:          t('public.profile.reservationStatusFailed'),
    quote_requested: t('public.profile.reservationStatusQuote'),
  }[s] ?? s
}
function prestaStatusClass(s) {
  return {
    pending: 'bg-amber-100 text-amber-700',
    paid: 'bg-emerald-100 text-emerald-700',
    failed: 'bg-red-100 text-red-700',
    quote_requested: 'bg-blue-100 text-blue-700',
  }[s] ?? 'bg-gray-100 text-gray-700'
}

function statusClass(s) {
  return {
    pending:   'bg-yellow-100 text-yellow-700',
    approved:  'bg-white/90 text-[#006d35]',
    accepted:  'bg-white/90 text-[#006d35]',
    validated: 'bg-white/90 text-[#006d35]',
    rejected:  'bg-white/90 text-red-500',
  }[s] ?? 'bg-white/90 text-gray-600'
}

const BADGE_META = {
  premier_pas:  { bg: '#fef9c3', border: 'rgba(202,138,4,0.3)', iconColor: '#92400e', icon: SparklesIcon },
  donneur:      { bg: '#fce7f3', border: 'rgba(219,39,119,0.2)', iconColor: '#be185d', icon: HeartIcon },
  reboiseur:    { bg: '#92f8ab', border: 'rgba(0,109,53,0.2)',   iconColor: '#0f522b', icon: GlobeAltIcon },
  activateur:   { bg: '#b0f2bd', border: 'rgba(44,106,64,0.2)', iconColor: '#006d35', icon: BoltIcon },
  expert:       { bg: '#dbeafe', border: 'rgba(59,130,246,0.2)', iconColor: '#1d4ed8', icon: AcademicCapIcon },
  eco_champion: { bg: '#fef3c7', border: 'rgba(217,119,6,0.3)', iconColor: '#b45309', icon: TrophyIcon },
}
const badgeBg        = k => BADGE_META[k]?.bg ?? '#f1f5f9'
const badgeBorder    = k => BADGE_META[k]?.border ?? '#cbd5e1'
const badgeIconColor = k => BADGE_META[k]?.iconColor ?? '#94a3b8'
const badgeIcon      = k => BADGE_META[k]?.icon ?? SparklesIcon
</script>
