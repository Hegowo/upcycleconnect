<template>
  <div class="bg-[#f7f9ff]">

    <section class="px-4 sm:px-6 pt-10 sm:pt-16 pb-8 sm:pb-12 max-w-[1280px] mx-auto">
      <div class="flex flex-col md:flex-row items-start gap-8 md:gap-16">
        <div class="flex-1 flex flex-col gap-6">
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl md:text-[64px] md:leading-[72px] leading-tight tracking-[-0.025em]">
            {{ t('public.events.heroTitlePart1') }}<br />
            <span class="text-[#006d35]">{{ t('public.events.heroTitleHighlight') }}</span>
          </h1>
          <p class="text-[#40617f] text-xl leading-[32px] max-w-[672px]">
            {{ t('public.events.heroSubtitle') }}
          </p>
        </div>
        <div class="hidden md:block shrink-0 relative w-72 h-72 lg:w-80 lg:h-80">
          <div class="absolute inset-0 bg-[rgba(0,109,53,0.05)] blur-[32px] rounded-full" />
          <div class="relative w-full h-full rounded-[40px] overflow-hidden shadow-2xl rotate-3 bg-gradient-to-br from-[#dbeafe] to-[#bfdbfe] flex items-center justify-center">
            <CalendarDaysIcon class="w-32 h-32 text-[#1d4ed8]/30" />
          </div>
        </div>
      </div>
    </section>

    <section class="px-6 pb-12 max-w-[1280px] mx-auto">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-6">

        <div class="md:col-span-8 bg-[#edf4ff] rounded-[32px] p-8 flex flex-col gap-4">
          <p class="text-[#001d32] font-jakarta font-bold text-lg">{{ t('public.events.searchSectionTitle') }}</p>
          <div class="relative">
            <input v-model="searchQuery" type="text"
              :placeholder="t('public.events.searchPlaceholder')"
              class="w-full bg-white pl-12 pr-4 py-[18px] rounded-xl text-base text-gray-600 outline-none shadow-sm focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400" />
            <MagnifyingGlassIcon class="w-[18px] h-[18px] absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
          </div>
          <div class="flex gap-3">
            <button @click="searchQuery = ''; activeCategory = 'all'"
              class="px-6 py-3 rounded-xl text-[#40617f] font-semibold text-sm bg-white border border-gray-200 hover:bg-gray-50 transition">
              {{ t('public.events.btnReset') }}
            </button>
          </div>
        </div>

        <div class="md:col-span-4 bg-[#edf4ff] rounded-[32px] p-8 flex flex-col gap-4">
          <p class="text-[#001d32] font-jakarta font-bold text-lg">{{ t('public.events.typeSectionTitle') }}</p>
          <div class="flex flex-col gap-2">
            <button @click="activeCategory = 'all'"
              class="w-full text-left px-5 py-3 rounded-xl text-sm font-semibold transition-all duration-200"
              :class="activeCategory === 'all' ? 'bg-[#006d35] text-white' : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'">
              {{ t('public.events.typeAll') }}
            </button>
            <button v-for="cat in categories" :key="cat.id"
              @click="activeCategory = cat.id"
              class="w-full text-left px-5 py-3 rounded-xl text-sm font-semibold transition-all duration-200"
              :class="activeCategory === cat.id ? 'bg-[#006d35] text-white' : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'">
              {{ cat.name }}
            </button>
          </div>
        </div>
      </div>
    </section>

    <section class="px-6 pb-16 max-w-[1280px] mx-auto">

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else>
        <div class="flex items-center justify-between mb-8">
          <p class="text-[#40617f] text-sm">
            {{ filteredEvents.length > 1
              ? t('public.events.resultsCountPlural', { count: filteredEvents.length })
              : t('public.events.resultsCount', { count: filteredEvents.length }) }}
          </p>
        </div>

        <div v-if="filteredEvents.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <article v-for="event in displayedEvents" :key="event.id"
            class="bg-white rounded-[32px] overflow-hidden flex flex-col hover:shadow-xl transition-shadow duration-300 cursor-pointer"
            @click="router.push(`/evenements/${event.id}`)">

            <div class="relative h-56 overflow-hidden flex-shrink-0">
              <div class="absolute inset-0 flex items-center justify-center"
                :style="`background: ${cardGradient(event.id)};`">
                <CalendarDaysIcon class="w-20 h-20 text-white/20" />
              </div>
              <div class="absolute top-4 left-4">
                <span v-if="event.category" class="text-xs font-bold uppercase tracking-[0.6px] px-3 py-1 rounded-full"
                  :class="categoryBadgeClass(event.category?.name)">
                  {{ event.category.name }}
                </span>
              </div>
              <div class="absolute top-4 right-4">
                <span class="text-xs font-bold px-3 py-1 rounded-full"
                  :class="isFull(event) ? 'bg-red-100 text-red-700' : 'bg-emerald-100 text-emerald-700'">
                  {{ isFull(event)
                    ? t('public.events.badgeFull')
                    : spotsLabel(event) }}
                </span>
              </div>
            </div>

            <div class="p-8 flex flex-col flex-1">
              <div class="flex items-center gap-2 text-[#006d35] text-sm font-semibold mb-3">
                <CalendarDaysIcon class="w-4 h-4" />
                {{ formatDate(event.start_date) }}
              </div>
              <h3 class="font-jakarta font-extrabold text-[#001d32] text-xl leading-[26px] mb-3">{{ event.title }}</h3>
              <p class="text-[#40617f] text-sm leading-[22px] flex-1 line-clamp-3">{{ event.description || '' }}</p>
              <div v-if="event.location" class="flex items-center gap-3 mt-4 mb-6">
                <MapPinIcon class="w-4 h-4 text-[#40617f] shrink-0" />
                <span class="text-[#40617f] text-sm">{{ event.location }}</span>
              </div>
              <div class="border-t border-gray-100 pt-5 flex items-center justify-between">
                <div>
                  <p class="text-[#40617f] text-xs uppercase tracking-wider">{{ t('public.events.freeLabel') }}</p>
                  <p class="text-[#006d35] font-bold text-lg">{{ t('public.events.priceFree') }}</p>
                </div>
                <button v-if="!isFull(event)"
                  class="bg-[#006d35] text-white font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#1b8848] transition-colors">
                  {{ t('public.events.ctaReserve') }}
                </button>
                <button v-else
                  class="bg-[#edf4ff] text-[#40617f] font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#d8eaff] transition-colors">
                  {{ t('public.events.ctaWaitingList') }}
                </button>
              </div>
            </div>
          </article>
        </div>

        <div v-if="filteredEvents.length > displayCount" class="flex justify-center mt-12">
          <button @click="displayCount += 6"
            class="flex items-center gap-2 px-10 py-4 rounded-xl border-2 border-[#006d35] text-[#006d35] font-bold text-base hover:bg-[#006d35]/5 transition">
            {{ t('public.events.loadMore') }}
            <ChevronDownIcon class="w-5 h-5" />
          </button>
        </div>

        <div v-if="!filteredEvents.length" class="flex flex-col items-center justify-center py-20 text-center">
          <CalendarDaysIcon class="w-16 h-16 text-[#40617f]/30 mb-4" />
          <p class="text-[#001d32] font-bold text-xl mb-2">{{ t('public.events.emptyTitle') }}</p>
          <p class="text-[#40617f] text-base">{{ t('public.events.emptyHint') }}</p>
        </div>
      </template>
    </section>

    <section class="px-6 pb-16 max-w-[1280px] mx-auto">
      <div class="rounded-[40px] p-12 flex flex-col md:flex-row items-center justify-between gap-8" style="background: #195687;">
        <div class="flex flex-col gap-4 flex-1">
          <span class="self-start text-xs font-bold uppercase tracking-[0.6px] text-white bg-white/20 px-4 py-1.5 rounded-full">
            {{ t('public.events.providerCtaTag') }}
          </span>
          <h2 class="font-jakarta font-bold text-white text-3xl leading-[38px]">
            {{ t('public.events.providerCtaTitle1') }}<br />{{ t('public.events.providerCtaTitle2') }}
          </h2>
          <p class="text-white/80 text-base leading-[26px] max-w-md">{{ t('public.events.providerCtaText') }}</p>
        </div>
        <div class="shrink-0">
          <button @click="router.push('/profil/pro')"
            class="bg-white text-[#006d35] font-bold text-base px-10 py-4 rounded-xl hover:opacity-90 transition shadow-lg">
            {{ t('public.events.providerCtaButton') }}
          </button>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  MagnifyingGlassIcon,
  CalendarDaysIcon,
  MapPinIcon,
  ChevronDownIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const router = useRouter()

const loading = ref(true)
const events = ref([])
const categories = ref([])
const searchQuery = ref('')
const activeCategory = ref('all')
const displayCount = ref(6)

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

function categoryBadgeClass(name = '') {
  const n = name.toLowerCase()
  if (n.includes('atelier') || n.includes('workshop')) return 'bg-[#d8eaff] text-[#1d4ed8]'
  if (n.includes('formation') || n.includes('training')) return 'bg-[#d1fae5] text-[#065f46]'
  if (n.includes('conf')) return 'bg-[#f3e8ff] text-[#7c3aed]'
  return 'bg-gray-100 text-gray-600'
}

function isFull(event) {
  return event.max_participants !== null && event.max_participants !== undefined &&
    event.registrations_count >= event.max_participants
}

function spotsLabel(event) {
  if (!event.max_participants) return t('public.events.freeLabel')
  const left = event.max_participants - (event.registrations_count || 0)
  if (left <= 0) return t('public.events.badgeFull')
  return left === 1
    ? t('public.events.badgePlace', { n: left })
    : t('public.events.badgePlaces', { n: left })
}

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: 'numeric', month: 'long', year: 'numeric' })
}

const filteredEvents = computed(() => {
  return events.value.filter(e => {
    const matchCat = activeCategory.value === 'all' || e.category?.id === activeCategory.value
    const q = searchQuery.value.toLowerCase()
    const matchSearch = !q ||
      e.title?.toLowerCase().includes(q) ||
      (e.description || '').toLowerCase().includes(q) ||
      (e.location || '').toLowerCase().includes(q)
    return matchCat && matchSearch
  })
})

const displayedEvents = computed(() => filteredEvents.value.slice(0, displayCount.value))

onMounted(async () => {
  try {
    const res = await fetch('/api/public/v1/events?upcoming=true&per_page=50')
    if (res.ok) {
      const data = await res.json()
      events.value = data.data || []
      const seen = new Set()
      categories.value = events.value
        .filter(e => e.category && !seen.has(e.category.id) && seen.add(e.category.id))
        .map(e => e.category)
    }
  } catch {}
  loading.value = false
})
</script>
