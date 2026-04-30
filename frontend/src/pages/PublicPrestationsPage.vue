<template>
  <div class="bg-[#f7f9ff]">

    <section class="px-4 sm:px-6 pt-10 sm:pt-16 pb-8 sm:pb-12 max-w-[1280px] mx-auto">
      <div class="flex flex-col md:flex-row items-start gap-8 md:gap-16">
        <div class="flex-1 flex flex-col gap-6">
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl md:text-[64px] md:leading-[72px] leading-tight tracking-[-0.025em]">
            {{ t('public.prestations.heroTitlePart1') }}<br />
            <span class="text-[#006d35]">{{ t('public.prestations.heroTitleHighlight') }}</span>
          </h1>
          <p class="text-[#40617f] text-xl leading-[32px] max-w-[672px]">
            {{ t('public.prestations.heroSubtitle') }}
          </p>
        </div>

        <div class="hidden md:block shrink-0 relative w-72 h-72 lg:w-80 lg:h-80">
          <div class="absolute inset-0 bg-[rgba(0,109,53,0.05)] blur-[32px] rounded-full" />
          <div class="relative w-full h-full rounded-[40px] overflow-hidden shadow-2xl rotate-3 bg-gradient-to-br from-[#d1fae5] to-[#bbf7d0] flex items-center justify-center">
            <HomeModernIcon class="w-32 h-32 text-[#006d35]/30" />
          </div>
        </div>
      </div>
    </section>

    <section class="px-4 sm:px-6 pb-8 sm:pb-12 max-w-[1280px] mx-auto">
      <div class="bg-[#edf4ff] rounded-[32px] p-6 sm:p-12 flex flex-col gap-6 sm:gap-8">

        <div class="flex gap-4">
          <div class="flex-1 relative">
            <input
              v-model="searchQuery"
              @keyup.enter="refresh"
              type="text"
              :placeholder="t('public.prestations.searchPlaceholder')"
              class="w-full bg-white pl-12 pr-4 py-[18px] rounded-xl text-base text-gray-600 outline-none shadow-sm focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400"
            />
            <MagnifyingGlassIcon class="w-[18px] h-[18px] absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
          </div>
          <button
            @click="refresh"
            class="flex items-center gap-2 px-8 py-4 rounded-xl text-white font-bold text-base transition hover:opacity-90"
            style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
          >
            <AdjustmentsHorizontalIcon class="w-[18px] h-[18px]" />
            {{ t('public.prestations.btnFilter') }}
          </button>
        </div>

        <div class="flex flex-wrap gap-3">
          <button
            @click="selectCategory(null)"
            class="px-6 py-2 rounded-full text-base font-semibold transition-all duration-200"
            :class="activeCategoryId === null
              ? 'bg-[#006d35] text-white'
              : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
          >
            {{ t('public.prestations.filterAll') }}
          </button>
          <button
            v-for="cat in categories"
            :key="cat.id"
            @click="selectCategory(cat.id)"
            class="px-6 py-2 rounded-full text-base font-semibold transition-all duration-200"
            :class="activeCategoryId === cat.id
              ? 'bg-[#006d35] text-white'
              : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
          >
            {{ cat.name }}
          </button>
        </div>

      </div>
    </section>

    <section class="px-4 sm:px-6 pb-12 sm:pb-16 max-w-[1280px] mx-auto">
      <div v-if="loading" class="text-center py-16 text-[#40617f]">{{ t('public.prestations.loading') }}</div>
      <div v-else-if="error" class="text-center py-16 text-red-600">{{ error }}</div>
      <div v-else-if="prestations.length === 0" class="text-center py-16 text-[#40617f]">
        {{ t('public.prestations.emptyResults') }}
      </div>
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">

        <article
          v-for="(prestation, idx) in prestations"
          :key="prestation.id"
          @click="goToDetail(prestation.id)"
          class="bg-white rounded-[32px] overflow-hidden flex flex-col hover:shadow-xl transition-shadow duration-300 cursor-pointer"
        >
          <div class="relative h-64 overflow-hidden flex-shrink-0">
            <div class="absolute inset-0 flex items-center justify-center" :style="`background: linear-gradient(135deg, ${palette(idx).from}, ${palette(idx).to});`">
              <component :is="palette(idx).icon" class="w-24 h-24 text-white/30" />
            </div>
            <div class="absolute top-4 left-4">
              <span
                class="bg-white/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-bold"
                :class="priceBadgeClass(prestation)"
              >
                {{ priceBadgeLabel(prestation) }}
              </span>
            </div>
          </div>

          <div class="p-8 flex flex-col flex-1">
            <h3 class="font-jakarta font-extrabold text-[#001d32] text-2xl leading-[30px] mb-4">
              {{ prestation.title }}
            </h3>
            <p class="text-[#40617f] text-sm leading-[22px] flex-1 line-clamp-4">
              {{ prestation.description || t('public.prestations.fallbackDescription') }}
            </p>

            <div class="border-t border-gray-100/50 mt-6 pt-5 flex items-center justify-between">
              <div>
                <p class="text-[#40617f] text-xs uppercase tracking-wider">{{ priceLabel(prestation) }}</p>
                <p class="text-[#006d35] font-bold text-xl">{{ priceValue(prestation) }}</p>
              </div>
              <button
                @click.stop="goToDetail(prestation.id)"
                class="bg-[#cee5ff] text-[#001d32] font-bold text-base px-5 py-3 rounded-xl hover:bg-[#b8d8ff] transition-colors"
              >
                {{ prestation.price_type === 'quote' ? t('public.prestations.ctaQuote') : t('public.prestations.ctaReserve') }}
              </button>
            </div>
          </div>
        </article>

      </div>
    </section>

    <section class="bg-[#edf4ff] px-4 sm:px-8 py-12 sm:py-20">
      <div class="max-w-[896px] mx-auto flex flex-col items-center gap-8 text-center px-6">
        <span class="flex items-center gap-2 bg-[rgba(0,109,53,0.1)] text-[#006d35] font-bold text-sm px-4 py-2 rounded-full">
          {{ t('public.prestations.didYouKnowTag') }}
        </span>
        <h2 class="font-jakarta font-extrabold text-[#001d32] text-[36px] leading-[40px]">
          {{ t('public.prestations.didYouKnowTitle') }}
        </h2>
        <p class="text-[#40617f] text-base leading-[24px]">
          {{ t('public.prestations.didYouKnowSubtitle') }}
        </p>
      </div>
    </section>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  MagnifyingGlassIcon,
  AdjustmentsHorizontalIcon,
  HomeModernIcon,
  ScissorsIcon,
  SparklesIcon,
  LightBulbIcon,
  ArchiveBoxIcon,
} from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

const { t } = useI18n()
const router = useRouter()

const categories = ref([])
const prestations = ref([])
const activeCategoryId = ref(null)
const searchQuery = ref('')
const loading = ref(false)
const error = ref('')

const palettes = [
  { from: '#d1fae5', to: '#bbf7d0', icon: HomeModernIcon },
  { from: '#dbeafe', to: '#bfdbfe', icon: ScissorsIcon },
  { from: '#f3e8ff', to: '#e9d5ff', icon: SparklesIcon },
  { from: '#fef3c7', to: '#fde68a', icon: LightBulbIcon },
  { from: '#d1fae5', to: '#a7f3d0', icon: ArchiveBoxIcon },
]
function palette(idx) { return palettes[idx % palettes.length] }

function priceValue(p) {
  if (p.price_type === 'quote') return t('public.prestations.priceQuoteValue')
  if (!p.price) return t('public.prestations.priceMissing')
  const amount = parseFloat(p.price)
  if (Number.isNaN(amount)) return p.price
  return `${amount.toFixed(2)} €${p.price_type === 'hourly' ? ' / h' : ''}`
}
function priceLabel(p) {
  if (p.price_type === 'quote') return t('public.prestations.priceQuoteLabel')
  if (p.price_type === 'hourly') return t('public.prestations.priceHourlyLabel')
  return t('public.prestations.priceFixedLabel')
}
function priceBadgeLabel(p) {
  if (p.price_type === 'quote') return t('public.prestations.badgeQuote')
  if (p.price_type === 'hourly') return t('public.prestations.badgeHourly')
  return t('public.prestations.badgeFixed')
}
function priceBadgeClass(p) {
  if (p.price_type === 'quote') return 'text-orange-600'
  return 'text-[#006d35]'
}

async function loadCategories() {
  try {
    const res = await publicGet('/categories')
    categories.value = res.data || []
  } catch {
    categories.value = []
  }
}

async function refresh() {
  loading.value = true
  error.value = ''
  try {
    const params = {}
    if (activeCategoryId.value) params.category_id = activeCategoryId.value
    if (searchQuery.value) params.search = searchQuery.value
    const res = await publicGet('/prestations', params)
    prestations.value = res.data || []
  } catch (e) {
    error.value = e.message || t('public.common.loadError')
  } finally {
    loading.value = false
  }
}

function selectCategory(id) {
  activeCategoryId.value = id
  refresh()
}

function goToDetail(id) {
  router.push(`/prestations/${id}`)
}

onMounted(async () => {
  await loadCategories()
  await refresh()
})
</script>
