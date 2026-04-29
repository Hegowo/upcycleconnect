<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-4 sm:px-8 pt-8 sm:pt-12 pb-16">

      <div v-if="loading" class="flex flex-col items-center justify-center py-24">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
        <p class="text-[#40617f] mt-4 text-sm">{{ t('public.depotConfirmation.loadingTitle') }}</p>
      </div>

      <div v-else-if="!deposit" class="bg-white rounded-[24px] p-12 text-center">
        <ExclamationCircleIcon class="w-16 h-16 mx-auto text-amber-500 mb-4" />
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight mb-2">
          {{ t('public.depotConfirmation.notFoundTitle') }}
        </h1>
        <p class="text-[#40617f] max-w-md mx-auto">{{ t('public.depotConfirmation.notFoundDesc') }}</p>
        <RouterLink to="/profil" class="inline-block mt-6 px-6 py-3 rounded-xl bg-[#006d35] text-white font-semibold hover:opacity-90 transition">
          {{ t('public.depotConfirmation.backToProfile') }}
        </RouterLink>
      </div>

      <template v-else>
        <div class="flex flex-col items-center text-center mb-12">
          <div
            class="w-20 h-20 rounded-full flex items-center justify-center mb-6 shadow-sm"
            :class="statusBg"
          >
            <CheckCircleIcon v-if="deposit.status === 'approved'" class="w-10 h-10 text-[#006d35]" />
            <ClockIcon v-else-if="deposit.status === 'pending'" class="w-10 h-10 text-amber-600" />
            <XCircleIcon v-else class="w-10 h-10 text-red-500" />
          </div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl sm:text-5xl tracking-tight mb-3">
            {{ headerTitle }}
          </h1>
          <p class="text-[#40617f] text-lg max-w-xl leading-relaxed">
            {{ headerSubtitle }}
          </p>
        </div>

        <div class="grid grid-cols-12 gap-4 sm:gap-6">

          <div class="col-span-12 lg:col-span-8 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6 sm:gap-8">

            <div class="flex items-center justify-between flex-wrap gap-3">
              <div>
                <h2 class="font-jakarta font-bold text-[#001d32] text-xl sm:text-2xl">{{ t('public.depotConfirmation.voucherTitle') }}</h2>
                <p class="text-[#40617f] text-sm mt-1">{{ t('public.depotConfirmation.validUntil', { date: expiryDate }) }}</p>
              </div>
              <span
                class="text-xs sm:text-sm font-bold uppercase tracking-wide px-4 py-2 rounded-full"
                :class="statusBadgeClass"
              >
                {{ statusBadgeLabel }}
              </span>
            </div>

            <div class="flex flex-col sm:flex-row gap-6 sm:gap-8 items-start">
              <div class="w-full sm:w-40 h-40 rounded-[20px] bg-[#edf4ff] border-2 border-[#d8eaff] flex flex-col items-center justify-center shrink-0 gap-3">
                <QrCodeIcon class="w-20 h-20 text-[#001d32]" :class="!deposit.qr_code ? 'opacity-30' : ''" />
                <span v-if="deposit.qr_code" class="text-[#001d32] font-mono font-bold text-xs tracking-widest">
                  {{ deposit.qr_code }}
                </span>
                <span v-else class="text-[#40617f] text-[10px] uppercase tracking-wider">{{ t('public.depotConfirmation.qrPending') }}</span>
              </div>

              <div class="flex-1 flex flex-col gap-4 w-full">
                <div class="grid grid-cols-2 gap-3 sm:gap-4">
                  <div class="bg-[#f7f9ff] rounded-xl p-3 sm:p-4">
                    <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.depotConfirmation.labelCode') }}</p>
                    <p class="font-jakarta font-bold text-[#001d32] text-sm sm:text-lg font-mono break-all">
                      {{ deposit.qr_code || '—' }}
                    </p>
                  </div>
                  <div class="bg-[#f7f9ff] rounded-xl p-3 sm:p-4">
                    <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.depotConfirmation.labelObjectRef') }}</p>
                    <p class="font-jakarta font-bold text-[#001d32] text-sm sm:text-lg">#UP-{{ deposit.id }}</p>
                  </div>
                  <div class="bg-[#f7f9ff] rounded-xl p-3 sm:p-4">
                    <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.depotConfirmation.labelValidatedAt') }}</p>
                    <p class="font-jakarta font-bold text-[#001d32] text-sm">
                      {{ deposit.validated_at ? formatDate(deposit.validated_at) : '—' }}
                    </p>
                  </div>
                  <div class="bg-[#f7f9ff] rounded-xl p-3 sm:p-4">
                    <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.depotConfirmation.labelExpiresAt') }}</p>
                    <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ expiryDate }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="deposit.status === 'rejected' && deposit.validation_note" class="bg-red-50 border border-red-100 rounded-xl p-4">
              <p class="text-red-700 text-xs font-bold uppercase tracking-wide mb-1">{{ t('public.depotConfirmation.rejectionNoteLabel') }}</p>
              <p class="text-red-800 text-sm">{{ deposit.validation_note }}</p>
            </div>

            <button
              v-if="deposit.status === 'approved'"
              @click="printVoucher"
              class="flex items-center justify-center gap-2 w-full py-4 rounded-xl text-white font-bold text-base transition hover:opacity-90"
              style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
            >
              <PrinterIcon class="w-5 h-5" />
              {{ t('public.depotConfirmation.printVoucher') }}
            </button>
          </div>

          <div class="col-span-12 lg:col-span-4 flex flex-col gap-4 sm:gap-6">
            <div class="bg-white rounded-[24px] p-6 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-4">
              <h3 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.depotConfirmation.objectTitle') }}</h3>
              <div class="h-32 rounded-xl bg-gradient-to-br from-[#fef3c7] to-[#fde68a] flex items-center justify-center">
                <HomeModernIcon class="w-16 h-16 text-white/40" />
              </div>
              <div>
                <p class="font-jakarta font-bold text-[#001d32] text-base">{{ deposit.title }}</p>
                <p class="text-[#40617f] text-sm mt-1">
                  {{ deposit.category?.name || t('public.depotConfirmation.categoryUnknown') }}
                </p>
                <p class="text-[#40617f] text-xs mt-1.5">
                  {{ t('public.depotConfirmation.conditionLabel') }} : {{ conditionLabel(deposit.condition) }}
                </p>
              </div>
              <div class="grid grid-cols-2 gap-2">
                <div v-if="deposit.estimated_weight" class="bg-[#f7f9ff] rounded-lg p-2.5">
                  <p class="text-[#40617f] text-[10px] uppercase tracking-wide">{{ t('public.depotConfirmation.weightLabel') }}</p>
                  <p class="font-bold text-[#001d32] text-sm mt-0.5">{{ deposit.estimated_weight }} kg</p>
                </div>
                <div v-if="deposit.carbon_savings" class="bg-[#f0fdf4] rounded-lg p-2.5">
                  <p class="text-[#006d35] text-[10px] uppercase tracking-wide">{{ t('public.depotConfirmation.co2Label') }}</p>
                  <p class="font-bold text-[#006d35] text-sm mt-0.5">{{ deposit.carbon_savings }} kg</p>
                </div>
              </div>
            </div>

            <div class="bg-[#edf4ff] rounded-[24px] p-6 flex flex-col gap-4">
              <h3 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.depotConfirmation.collectPointTitle') }}</h3>
              <div class="flex items-start gap-3">
                <div class="w-10 h-10 bg-white rounded-xl flex items-center justify-center shrink-0">
                  <MapPinIcon class="w-5 h-5 text-[#006d35]" />
                </div>
                <div>
                  <p class="font-jakarta font-bold text-[#001d32] text-base">{{ t('public.depotConfirmation.demoCollectPointName') }}</p>
                  <p class="text-[#40617f] text-sm mt-0.5">{{ t('public.depotConfirmation.demoCollectPointAddress1') }}</p>
                  <p class="text-[#40617f] text-sm">{{ t('public.depotConfirmation.demoCollectPointAddress2') }}</p>
                  <span class="inline-block mt-2 text-[#006d35] text-xs font-bold bg-[#d1fae5] px-2 py-0.5 rounded-full">
                    {{ t('public.depotConfirmation.openingBadge') }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="col-span-12 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)]">
            <h3 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">{{ t('public.depotConfirmation.stepsTitle') }}</h3>
            <div class="grid grid-cols-1 sm:grid-cols-3 gap-6">
              <div v-for="step in steps" :key="step.num" class="flex gap-4 items-start">
                <div class="w-10 h-10 rounded-full bg-[#006d35] text-white font-bold text-lg flex items-center justify-center shrink-0">
                  {{ step.num }}
                </div>
                <div>
                  <p class="font-jakarta font-bold text-[#001d32] text-base">{{ step.title }}</p>
                  <p class="text-[#40617f] text-sm leading-relaxed mt-1">{{ step.desc }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="flex items-center justify-center gap-4 mt-10 flex-wrap">
          <RouterLink to="/profil" class="flex items-center gap-2 text-[#40617f] hover:text-[#006d35] font-medium text-sm transition">
            <ArrowLeftIcon class="w-4 h-4" />
            {{ t('public.depotConfirmation.backToProfile') }}
          </RouterLink>
          <RouterLink to="/depot" class="flex items-center gap-2 px-4 py-2 rounded-lg bg-[#006d35] text-white font-medium text-sm hover:opacity-90 transition">
            {{ t('public.depotConfirmation.newDepositCta') }}
          </RouterLink>
        </div>
      </template>

    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  CheckCircleIcon,
  ClockIcon,
  XCircleIcon,
  QrCodeIcon,
  PrinterIcon,
  MapPinIcon,
  HomeModernIcon,
  ArrowLeftIcon,
  ExclamationCircleIcon,
} from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const userAuth = useUserAuthStore()

const loading = ref(true)
const deposit = ref(null)

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

const headerTitle = computed(() => {
  if (!deposit.value) return ''
  if (deposit.value.status === 'approved') return t('public.depotConfirmation.title')
  if (deposit.value.status === 'rejected') return t('public.depotConfirmation.rejectedTitle')
  return t('public.depotConfirmation.pendingTitle')
})

const headerSubtitle = computed(() => {
  if (!deposit.value) return ''
  if (deposit.value.status === 'approved') return t('public.depotConfirmation.subtitle')
  if (deposit.value.status === 'rejected') return t('public.depotConfirmation.rejectedSubtitle')
  return t('public.depotConfirmation.pendingSubtitle')
})

const statusBg = computed(() => {
  if (!deposit.value) return 'bg-gray-100'
  if (deposit.value.status === 'approved') return 'bg-[#d1fae5]'
  if (deposit.value.status === 'rejected') return 'bg-red-100'
  return 'bg-amber-100'
})

const statusBadgeClass = computed(() => {
  if (!deposit.value) return 'bg-gray-100 text-gray-700'
  if (deposit.value.status === 'approved') return 'bg-[#d1fae5] text-[#006d35]'
  if (deposit.value.status === 'rejected') return 'bg-red-100 text-red-700'
  return 'bg-amber-100 text-amber-700'
})

const statusBadgeLabel = computed(() => {
  if (!deposit.value) return ''
  if (deposit.value.status === 'approved') return t('public.depotConfirmation.validatedBadge')
  if (deposit.value.status === 'rejected') return t('public.depotConfirmation.rejectedBadge')
  return t('public.depotConfirmation.pendingBadge')
})

const expiryDate = computed(() => {
  if (!deposit.value) return ''
  const base = deposit.value.validated_at || deposit.value.created_at
  if (!base) return ''
  const d = new Date(base)
  d.setMonth(d.getMonth() + 1)
  return formatDate(d.toISOString())
})

const steps = computed(() => [
  { num: 1, title: t('public.depotConfirmation.step1Title'), desc: t('public.depotConfirmation.step1Desc') },
  { num: 2, title: t('public.depotConfirmation.step2Title'), desc: t('public.depotConfirmation.step2Desc') },
  { num: 3, title: t('public.depotConfirmation.step3Title'), desc: t('public.depotConfirmation.step3Desc') },
])

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: 'numeric', month: 'long', year: 'numeric' })
}

function conditionLabel(c) {
  return {
    good: t('public.depotConfirmation.conditionGood'),
    fair: t('public.depotConfirmation.conditionFair'),
    poor: t('public.depotConfirmation.conditionPoor'),
  }[c] || c || '—'
}

function printVoucher() {
  window.print()
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil')
    return
  }
  const id = route.query.id
  if (!id) {
    loading.value = false
    return
  }
  try {
    deposit.value = await userApi(`/deposits/${id}`)
  } catch {
    deposit.value = null
  } finally {
    loading.value = false
  }
})
</script>
