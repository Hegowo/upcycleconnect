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
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="t('public.events.searchPlaceholder')"
              class="w-full bg-white pl-12 pr-4 py-[18px] rounded-xl text-base text-gray-600 outline-none shadow-sm focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400"
            />
            <MagnifyingGlassIcon class="w-[18px] h-[18px] absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
          </div>
          <div class="flex gap-3">
            <button
              class="flex items-center gap-2 px-6 py-3 rounded-xl text-white font-bold text-sm transition hover:opacity-90"
              style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
            >
              <MagnifyingGlassIcon class="w-4 h-4" />
              {{ t('public.events.btnSearch') }}
            </button>
            <button
              @click="searchQuery = ''; activeType = 'all'"
              class="px-6 py-3 rounded-xl text-[#40617f] font-semibold text-sm bg-white border border-gray-200 hover:bg-gray-50 transition"
            >
              {{ t('public.events.btnReset') }}
            </button>
          </div>
        </div>

        <div class="md:col-span-4 bg-[#edf4ff] rounded-[32px] p-8 flex flex-col gap-4">
          <p class="text-[#001d32] font-jakarta font-bold text-lg">{{ t('public.events.typeSectionTitle') }}</p>
          <div class="flex flex-col gap-2">
            <button
              v-for="type in eventTypes"
              :key="type.key"
              @click="activeType = type.key"
              class="w-full text-left px-5 py-3 rounded-xl text-sm font-semibold transition-all duration-200"
              :class="activeType === type.key
                ? 'bg-[#006d35] text-white'
                : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
            >
              {{ type.label }}
            </button>
          </div>
        </div>

      </div>
    </section>

    <section class="px-6 pb-16 max-w-[1280px] mx-auto">

      <div class="flex items-center justify-between mb-8">
        <p class="text-[#40617f] text-sm">
          {{ filteredEvents.length > 1
            ? t('public.events.resultsCountPlural', { count: filteredEvents.length })
            : t('public.events.resultsCount', { count: filteredEvents.length }) }}
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">

        <article
          v-for="event in displayedEvents"
          :key="event.id"
          class="bg-white rounded-[32px] overflow-hidden flex flex-col hover:shadow-xl transition-shadow duration-300"
        >
          <div class="relative h-56 overflow-hidden flex-shrink-0">
            <div
              class="absolute inset-0 flex items-center justify-center"
              :style="`background: linear-gradient(135deg, ${event.colorFrom}, ${event.colorTo});`"
            >
              <component :is="event.icon" class="w-20 h-20 text-white/25" />
            </div>

            <div class="absolute top-4 left-4">
              <span
                class="text-xs font-bold uppercase tracking-[0.6px] px-3 py-1 rounded-full"
                :class="typeBadgeClass(event.typeKey)"
              >
                {{ typeLabel(event.typeKey) }}
              </span>
            </div>

            <div class="absolute top-4 right-4">
              <span
                class="text-xs font-bold px-3 py-1 rounded-full"
                :class="event.complet
                  ? 'bg-red-100 text-red-700'
                  : 'bg-emerald-100 text-emerald-700'"
              >
                {{ event.complet
                  ? t('public.events.badgeFull')
                  : (event.places > 1
                    ? t('public.events.badgePlaces', { n: event.places })
                    : t('public.events.badgePlace', { n: event.places })) }}
              </span>
            </div>
          </div>

          <div class="p-8 flex flex-col flex-1">

            <div class="flex items-center gap-2 text-[#006d35] text-sm font-semibold mb-3">
              <CalendarDaysIcon class="w-4 h-4" />
              {{ event.date }}
            </div>

            <h3 class="font-jakarta font-extrabold text-[#001d32] text-xl leading-[26px] mb-3">
              {{ event.title }}
            </h3>
            <p class="text-[#40617f] text-sm leading-[22px] flex-1">
              {{ event.description }}
            </p>

            <div class="flex items-center gap-3 mt-4 mb-6">
              <MapPinIcon class="w-4 h-4 text-[#40617f] shrink-0" />
              <span class="text-[#40617f] text-sm">{{ event.location }}</span>
            </div>

            <div class="border-t border-gray-100 pt-5 flex items-center justify-between">
              <div>
                <p class="text-[#40617f] text-xs uppercase tracking-wider">{{ t('public.events.priceLabel') }}</p>
                <p class="text-[#006d35] font-bold text-xl">{{ event.price }}</p>
              </div>
              <button
                v-if="!event.complet"
                class="bg-[#006d35] text-white font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#1b8848] transition-colors"
              >
                {{ t('public.events.ctaReserve') }}
              </button>
              <button
                v-else
                class="bg-[#edf4ff] text-[#40617f] font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#d8eaff] transition-colors"
              >
                {{ t('public.events.ctaWaitingList') }}
              </button>
            </div>
          </div>
        </article>

      </div>

      <div v-if="filteredEvents.length > displayCount" class="flex justify-center mt-12">
        <button
          @click="displayCount += 6"
          class="flex items-center gap-2 px-10 py-4 rounded-xl border-2 border-[#006d35] text-[#006d35] font-bold text-base hover:bg-[#006d35]/5 transition"
        >
          {{ t('public.events.loadMore') }}
          <ChevronDownIcon class="w-5 h-5" />
        </button>
      </div>

      <div v-if="filteredEvents.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
        <CalendarDaysIcon class="w-16 h-16 text-[#40617f]/30 mb-4" />
        <p class="text-[#001d32] font-bold text-xl mb-2">{{ t('public.events.emptyTitle') }}</p>
        <p class="text-[#40617f] text-base">{{ t('public.events.emptyHint') }}</p>
      </div>

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
          <p class="text-white/80 text-base leading-[26px] max-w-md">
            {{ t('public.events.providerCtaText') }}
          </p>
        </div>
        <div class="shrink-0">
          <button class="bg-white text-[#006d35] font-bold text-base px-10 py-4 rounded-xl hover:opacity-90 transition shadow-lg">
            {{ t('public.events.providerCtaButton') }}
          </button>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  MagnifyingGlassIcon,
  CalendarDaysIcon,
  MapPinIcon,
  ChevronDownIcon,
  ScissorsIcon,
  SparklesIcon,
  LightBulbIcon,
  AcademicCapIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()

const searchQuery = ref('')
const activeType = ref('all')
const displayCount = ref(6)

const eventTypes = computed(() => [
  { key: 'all',        label: t('public.events.typeAll') },
  { key: 'workshop',   label: t('public.events.typeWorkshop') },
  { key: 'training',   label: t('public.events.typeTraining') },
  { key: 'conference', label: t('public.events.typeConference') },
])

function typeLabel(key) {
  return {
    workshop:   t('public.events.typeWorkshop'),
    training:   t('public.events.typeTraining'),
    conference: t('public.events.typeConference'),
  }[key] || key
}

const typeBadgeClass = (key) => {
  const map = {
    workshop:   'bg-[#d8eaff] text-[#1d4ed8]',
    training:   'bg-[#d1fae5] text-[#065f46]',
    conference: 'bg-[#f3e8ff] text-[#7c3aed]',
  }
  return map[key] || 'bg-gray-100 text-gray-600'
}

const events = computed(() => [
  {
    id: 1,
    title: 'Atelier Upcycling Textile',
    description: 'Apprenez à transformer vos vieux vêtements en créations uniques. Techniques de teinture naturelle, couture créative et détournement de matières.',
    typeKey: 'workshop',
    date: '12 avril 2026',
    location: 'Paris 11e',
    price: '45€',
    places: 4,
    complet: false,
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: ScissorsIcon,
  },
  {
    id: 2,
    title: 'Conférence : Économie Circulaire & Design',
    description: "Réflexion sur l'impact environnemental du design. Intervenants experts, tables rondes et échanges avec des acteurs du secteur.",
    typeKey: 'conference',
    date: '19 avril 2026',
    location: 'Lyon 2e',
    price: t('public.events.priceFree'),
    places: 12,
    complet: false,
    colorFrom: '#f3e8ff',
    colorTo: '#e9d5ff',
    icon: SparklesIcon,
  },
  {
    id: 3,
    title: 'Formation Création de Bijoux',
    description: 'Formation certifiante sur 2 jours : refonte de vieux métaux, travail du verre recyclé et création de pièces uniques.',
    typeKey: 'training',
    date: '26-27 avril 2026',
    location: 'Marseille',
    price: '180€',
    places: 0,
    complet: true,
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: SparklesIcon,
  },
  {
    id: 4,
    title: 'Atelier Restauration de Mobilier',
    description: "Ponçage, vernis écologique, réparation structurelle. Donnez un second souffle à vos meubles avec l'aide de nos ébénistes.",
    typeKey: 'workshop',
    date: '3 mai 2026',
    location: 'Bordeaux',
    price: '65€',
    places: 2,
    complet: false,
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: LightBulbIcon,
  },
  {
    id: 5,
    title: 'Formation Bee-Wraps & Emballages',
    description: 'Créez vos propres alternatives aux emballages plastiques : bee-wraps, sacs à vrac, boites en carton recyclé.',
    typeKey: 'training',
    date: '10 mai 2026',
    location: 'Toulouse',
    price: '55€',
    places: 8,
    complet: false,
    colorFrom: '#d1fae5',
    colorTo: '#bbf7d0',
    icon: AcademicCapIcon,
  },
  {
    id: 6,
    title: 'Conférence : Zéro Déchet au Quotidien',
    description: 'Découvrez des gestes simples et efficaces pour réduire vos déchets et adopter un mode de vie plus durable.',
    typeKey: 'conference',
    date: '17 mai 2026',
    location: 'Nantes',
    price: t('public.events.priceFree'),
    places: 30,
    complet: false,
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: SparklesIcon,
  },
  {
    id: 7,
    title: 'Atelier Luminaires Industriels',
    description: 'Créez des objets déco et luminaires originaux à partir de pièces industrielles récupérées. Résultats garantis.',
    typeKey: 'workshop',
    date: '24 mai 2026',
    location: 'Strasbourg',
    price: '75€',
    places: 6,
    complet: false,
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: LightBulbIcon,
  },
  {
    id: 8,
    title: 'Formation Artisan Certifié Upcycling',
    description: 'Formation complète sur 5 jours pour les professionnels souhaitant se certifier dans les métiers de l\'upcycling.',
    typeKey: 'training',
    date: '2-6 juin 2026',
    location: 'Paris 3e',
    price: '490€',
    places: 3,
    complet: false,
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: AcademicCapIcon,
  },
])

const filteredEvents = computed(() => {
  return events.value.filter(e => {
    const matchType = activeType.value === 'all' || e.typeKey === activeType.value
    const matchSearch = !searchQuery.value ||
      e.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      e.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      e.location.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchType && matchSearch
  })
})

const displayedEvents = computed(() => filteredEvents.value.slice(0, displayCount.value))
</script>
