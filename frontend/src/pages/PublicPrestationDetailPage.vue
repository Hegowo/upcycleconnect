<template>
  <div class="bg-[#f7f9ff] min-h-screen">
    <section class="px-4 sm:px-6 pt-10 pb-12 max-w-[1024px] mx-auto">

      <button
        @click="$router.back()"
        class="flex items-center gap-2 text-[#40617f] hover:text-[#006d35] mb-6"
      >
        <ArrowLeftIcon class="w-5 h-5" />
        {{ t('public.prestationDetail.back') }}
      </button>

      <div v-if="loading" class="text-center py-16 text-[#40617f]">{{ t('public.prestationDetail.loading') }}</div>
      <div v-else-if="error" class="text-center py-16 text-red-600">{{ error }}</div>
      <div v-else-if="prestation" class="bg-white rounded-[32px] overflow-hidden shadow-sm">

        <div class="relative h-72 bg-gradient-to-br from-[#d1fae5] to-[#bbf7d0] flex items-center justify-center">
          <SparklesIcon class="w-32 h-32 text-white/40" />
          <div class="absolute top-6 left-6">
            <span class="bg-white/90 px-4 py-1.5 rounded-full text-xs font-bold text-[#006d35]">
              {{ prestation.category?.name || t('public.prestationDetail.categoryFallback') }}
            </span>
          </div>
        </div>

        <div class="p-8 md:p-12">
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl md:text-4xl mb-4">
            {{ prestation.title }}
          </h1>

          <div v-if="prestation.provider" class="text-[#40617f] mb-6 text-sm">
            {{ t('public.prestationDetail.providerPrefix') }}
            <span class="font-semibold text-[#001d32]">
              {{ prestation.provider.first_name }} {{ prestation.provider.last_name }}
            </span>
          </div>

          <p class="text-[#40617f] text-base leading-7 whitespace-pre-wrap mb-8">
            {{ prestation.description || t('public.prestationDetail.noDescription') }}
          </p>

          <div class="bg-[#edf4ff] rounded-[24px] p-6 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
            <div>
              <p class="text-[#40617f] text-xs uppercase tracking-wider">{{ priceLabel }}</p>
              <p class="text-[#006d35] font-extrabold text-3xl">{{ priceValue }}</p>
              <p v-if="prestation.price_type === 'quote'" class="text-xs text-[#40617f] mt-1">
                {{ t('public.prestationDetail.quoteNote') }}
              </p>
            </div>

            <div class="flex flex-col gap-3 w-full sm:w-auto">
              <textarea
                v-model="notes"
                :placeholder="t('public.prestationDetail.notesPlaceholder')"
                rows="3"
                class="w-full sm:w-80 p-3 rounded-xl text-sm text-gray-700 outline-none border border-[#cee5ff] focus:ring-2 focus:ring-[#006d35]/20"
              />
              <button
                @click="reserve"
                :disabled="reserving"
                class="px-8 py-4 rounded-xl text-white font-bold text-base transition hover:opacity-90 disabled:opacity-60"
                style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
              >
                <span v-if="reserving">{{ t('public.prestationDetail.ctaProcessing') }}</span>
                <span v-else-if="prestation.price_type === 'quote'">{{ t('public.prestationDetail.ctaQuote') }}</span>
                <span v-else>{{ t('public.prestationDetail.ctaReserve') }}</span>
              </button>
            </div>
          </div>

          <div v-if="feedback" class="mt-6 p-4 rounded-xl bg-[#d1fae5] text-[#006d35] text-sm">
            {{ feedback }}
          </div>
          <div v-if="reserveError" class="mt-6 p-4 rounded-xl bg-red-50 text-red-600 text-sm">
            {{ reserveError }}
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowLeftIcon, SparklesIcon } from '@heroicons/vue/24/outline'
import { publicGet, userApi } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userAuth = useUserAuthStore()

const prestation = ref(null)
const loading = ref(true)
const error = ref('')
const notes = ref('')
const reserving = ref(false)
const feedback = ref('')
const reserveError = ref('')

const priceValue = computed(() => {
  const p = prestation.value
  if (!p) return ''
  if (p.price_type === 'quote') return t('public.prestationDetail.priceQuote')
  if (!p.price) return t('public.prestationDetail.priceMissing')
  const amount = parseFloat(p.price)
  if (Number.isNaN(amount)) return p.price
  return `${amount.toFixed(2)} €${p.price_type === 'hourly' ? ' / h' : ''}`
})
const priceLabel = computed(() => {
  const p = prestation.value
  if (!p) return ''
  if (p.price_type === 'quote') return t('public.prestationDetail.labelQuote')
  if (p.price_type === 'hourly') return t('public.prestationDetail.labelHourly')
  return t('public.prestationDetail.labelFixed')
})

async function loadDetail() {
  loading.value = true
  try {
    prestation.value = await publicGet(`/prestations/${route.params.id}`)
  } catch (e) {
    error.value = e.message || t('public.prestationDetail.errorNotFound')
  } finally {
    loading.value = false
  }
}

async function reserve() {
  if (!userAuth.isLoggedIn) {
    router.push(`/connexion?redirect=/prestations/${route.params.id}`)
    return
  }
  reserving.value = true
  feedback.value = ''
  reserveError.value = ''
  try {
    const res = await userApi(`/prestations/${route.params.id}/reserve`, {
      method: 'POST',
      body: JSON.stringify({ notes: notes.value || null }),
    })
    if (res.type === 'quote') {
      feedback.value = res.message || t('public.prestationDetail.feedbackQuoteSent')
      setTimeout(() => router.push('/profil/factures'), 2000)
    } else if (res.checkout_url) {
      window.location.href = res.checkout_url
    } else {
      reserveError.value = t('public.prestationDetail.errorUnexpected')
    }
  } catch (e) {
    reserveError.value = e.message || t('public.prestationDetail.errorReserve')
  } finally {
    reserving.value = false
  }
}

onMounted(loadDetail)
</script>
