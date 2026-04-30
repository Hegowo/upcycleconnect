<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-10">
    <div class="max-w-[960px] mx-auto px-6 py-10">

      <button
        @click="$router.push('/profil')"
        class="inline-flex items-center gap-2 text-[#40617f] hover:text-[#001d32] text-sm font-medium mb-6 transition"
      >
        <ArrowLeftIcon class="w-4 h-4" />
        {{ t('public.reservation.back') }}
      </button>

      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else-if="error" class="bg-white rounded-[24px] p-10 text-center">
        <p class="text-[#40617f]">{{ error }}</p>
      </div>

      <div v-else-if="reservation" class="grid grid-cols-12 gap-6">

        <div class="col-span-12 bg-white rounded-[24px] p-6 sm:p-8 relative overflow-hidden">
          <div class="absolute -right-10 -top-10 w-40 h-40 rounded-full bg-[rgba(0,109,53,0.05)]" />

          <div class="flex flex-wrap items-start gap-4 mb-4">
            <div class="w-14 h-14 rounded-xl flex items-center justify-center shrink-0"
              :class="statusBg(reservation.status)"
            >
              <WrenchScrewdriverIcon class="w-7 h-7" :class="statusText(reservation.status)" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap mb-1">
                <span
                  class="text-[11px] font-bold uppercase tracking-wide px-2 py-0.5 rounded-full"
                  :class="statusBadge(reservation.status)"
                >
                  {{ statusLabel(reservation.status) }}
                </span>
                <span class="text-xs text-[#40617f]">
                  {{ t('public.reservation.idLabel', { id: reservation.id }) }}
                </span>
              </div>
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl tracking-tight">
                {{ reservation.prestation?.title || t('public.reservation.fallbackTitle') }}
              </h1>
              <p v-if="reservation.prestation?.category?.name" class="text-[#40617f] text-sm mt-0.5">
                {{ reservation.prestation.category.name }}
              </p>
            </div>
          </div>

          <p v-if="reservation.prestation?.description" class="text-[#3f4a3f] text-sm leading-relaxed mt-4">
            {{ reservation.prestation.description }}
          </p>
        </div>

        <div class="col-span-12 lg:col-span-7 bg-white rounded-[24px] p-6 sm:p-8">
          <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-5">{{ t('public.reservation.detailsTitle') }}</h2>

          <dl class="divide-y divide-[#edf4ff]">
            <div class="flex items-center justify-between py-3">
              <dt class="text-[#40617f] text-sm flex items-center gap-2">
                <CalendarDaysIcon class="w-4 h-4" />
                {{ t('public.reservation.createdAt') }}
              </dt>
              <dd class="text-[#001d32] text-sm font-medium">{{ formatDateTime(reservation.created_at) }}</dd>
            </div>

            <div class="flex items-center justify-between py-3">
              <dt class="text-[#40617f] text-sm flex items-center gap-2">
                <ClockIcon class="w-4 h-4" />
                {{ t('public.reservation.updatedAt') }}
              </dt>
              <dd class="text-[#001d32] text-sm font-medium">{{ formatDateTime(reservation.updated_at) }}</dd>
            </div>

            <div v-if="reservation.amount_cents" class="flex items-center justify-between py-3">
              <dt class="text-[#40617f] text-sm flex items-center gap-2">
                <BanknotesIcon class="w-4 h-4" />
                {{ t('public.reservation.amountLabel') }}
              </dt>
              <dd class="text-[#006d35] text-base font-bold">{{ formatAmount(reservation.amount_cents, reservation.currency) }}</dd>
            </div>

            <div v-if="reservation.prestation?.provider" class="flex items-center justify-between py-3">
              <dt class="text-[#40617f] text-sm flex items-center gap-2">
                <UserIcon class="w-4 h-4" />
                {{ t('public.reservation.providerLabel') }}
              </dt>
              <dd class="text-[#001d32] text-sm font-medium">
                {{ providerName(reservation.prestation.provider) }}
              </dd>
            </div>

            <div v-if="reservation.notes" class="py-3">
              <dt class="text-[#40617f] text-sm flex items-center gap-2 mb-2">
                <DocumentTextIcon class="w-4 h-4" />
                {{ t('public.reservation.notesLabel') }}
              </dt>
              <dd class="text-[#001d32] text-sm bg-[#f7f9ff] rounded-lg p-3">{{ reservation.notes }}</dd>
            </div>
          </dl>
        </div>

        <div class="col-span-12 lg:col-span-5 flex flex-col gap-6">

          <div class="bg-[#edf4ff] rounded-[24px] p-6">
            <h3 class="font-jakarta font-bold text-[#001d32] text-base mb-4">{{ t('public.reservation.paymentStatusTitle') }}</h3>
            <div class="flex items-start gap-3">
              <div class="w-10 h-10 rounded-full flex items-center justify-center shrink-0"
                :class="statusBg(reservation.status)"
              >
                <component :is="statusIcon(reservation.status)" class="w-5 h-5" :class="statusText(reservation.status)" />
              </div>
              <div class="flex-1">
                <p class="font-semibold text-[#001d32] text-sm">{{ statusHeadline(reservation.status) }}</p>
                <p class="text-[#40617f] text-xs mt-1 leading-relaxed">{{ statusDescription(reservation.status) }}</p>
              </div>
            </div>
          </div>

          <div v-if="invoice" class="bg-white rounded-[24px] p-6 border border-[#edf4ff]">
            <h3 class="font-jakarta font-bold text-[#001d32] text-base mb-3">
              {{ invoice.type === 'quote' ? t('public.reservation.quoteTitle') : t('public.reservation.invoiceTitle') }}
            </h3>
            <p class="text-[#40617f] text-xs mb-1">{{ t('public.reservation.invoiceNumber') }}</p>
            <p class="font-semibold text-[#001d32] text-sm mb-4">{{ invoice.number }}</p>
            <button
              v-if="invoice.has_pdf"
              @click="downloadInvoice"
              :disabled="downloading"
              class="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm transition hover:opacity-90 disabled:opacity-60"
              style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
            >
              <ArrowDownTrayIcon class="w-4 h-4" />
              {{ downloading ? t('public.reservation.downloading') : t('public.reservation.downloadPdf') }}
            </button>
            <router-link
              to="/profil/factures"
              class="block text-center text-xs text-[#40617f] hover:text-[#006d35] mt-3"
            >
              {{ t('public.reservation.seeAllInvoices') }}
            </router-link>
          </div>

          <router-link
            v-if="reservation.prestation?.id"
            :to="`/prestations/${reservation.prestation.id}`"
            class="bg-white rounded-[24px] p-6 border border-[#edf4ff] flex items-center justify-between hover:border-[#006d35] transition"
          >
            <div>
              <p class="text-[#40617f] text-xs">{{ t('public.reservation.viewPrestationLabel') }}</p>
              <p class="font-semibold text-[#001d32] text-sm">{{ t('public.reservation.viewPrestationCta') }}</p>
            </div>
            <ChevronRightIcon class="w-5 h-5 text-[#40617f]" />
          </router-link>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ArrowLeftIcon,
  ArrowDownTrayIcon,
  BanknotesIcon,
  CalendarDaysIcon,
  CheckCircleIcon,
  ChevronRightIcon,
  ClockIcon,
  DocumentTextIcon,
  ExclamationTriangleIcon,
  UserIcon,
  WrenchScrewdriverIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const route    = useRoute()
const router   = useRouter()
const userAuth = useUserAuthStore()

const loading     = ref(true)
const error       = ref('')
const reservation = ref(null)
const invoice     = ref(null)
const downloading = ref(false)

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push(`/connexion?redirect=/profil/reservations/${route.params.id}`)
    return
  }
  try {
    const res = await fetch(`/api/v1/reservations/${route.params.id}`, {
      headers: {
        Authorization: `Bearer ${userAuth.token}`,
        Accept: 'application/json',
      },
    })
    if (res.status === 404) {
      error.value = t('public.reservation.notFound')
    } else if (!res.ok) {
      error.value = t('public.reservation.loadError')
    } else {
      const data = await res.json()
      reservation.value = data.reservation
      invoice.value     = data.invoice
    }
  } catch (e) {
    error.value = t('public.reservation.networkError')
  }
  loading.value = false
})

async function downloadInvoice() {
  if (!invoice.value?.id) return
  downloading.value = true
  try {
    const res = await fetch(`/api/v1/invoices/${invoice.value.id}/download`, {
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (!res.ok) throw new Error('download failed')
    const blob = await res.blob()
    const url  = URL.createObjectURL(blob)
    const a    = document.createElement('a')
    a.href     = url
    a.download = `${invoice.value.number}.pdf`
    document.body.appendChild(a)
    a.click()
    a.remove()
    URL.revokeObjectURL(url)
  } catch {
  } finally {
    downloading.value = false
  }
}

function formatDateTime(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleString(localeCode.value, {
    day: 'numeric', month: 'long', year: 'numeric',
    hour: '2-digit', minute: '2-digit',
  })
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

function providerName(p) {
  if (!p) return ''
  const parts = [p.first_name, p.last_name].filter(Boolean).join(' ')
  return parts || p.company_name || p.email || ''
}

function statusLabel(s) {
  return {
    pending:         t('public.reservation.statusLabelPending'),
    paid:            t('public.reservation.statusLabelPaid'),
    failed:          t('public.reservation.statusLabelFailed'),
    quote_requested: t('public.reservation.statusLabelQuote'),
  }[s] ?? s
}
function statusHeadline(s) {
  return {
    pending:         t('public.reservation.statusHeadlinePending'),
    paid:            t('public.reservation.statusHeadlinePaid'),
    failed:          t('public.reservation.statusHeadlineFailed'),
    quote_requested: t('public.reservation.statusHeadlineQuote'),
  }[s] ?? t('public.reservation.statusHeadlineDefault')
}
function statusDescription(s) {
  return {
    pending:         t('public.reservation.statusDescPending'),
    paid:            t('public.reservation.statusDescPaid'),
    failed:          t('public.reservation.statusDescFailed'),
    quote_requested: t('public.reservation.statusDescQuote'),
  }[s] ?? ''
}
function statusIcon(s) {
  if (s === 'paid') return CheckCircleIcon
  if (s === 'failed') return ExclamationTriangleIcon
  if (s === 'quote_requested') return DocumentTextIcon
  return ClockIcon
}
function statusBg(s) {
  return {
    pending:         'bg-amber-100',
    paid:            'bg-emerald-100',
    failed:          'bg-red-100',
    quote_requested: 'bg-blue-100',
  }[s] ?? 'bg-gray-100'
}
function statusText(s) {
  return {
    pending:         'text-amber-600',
    paid:            'text-emerald-600',
    failed:          'text-red-600',
    quote_requested: 'text-blue-600',
  }[s] ?? 'text-gray-600'
}
function statusBadge(s) {
  return {
    pending:         'bg-amber-100 text-amber-700',
    paid:            'bg-emerald-100 text-emerald-700',
    failed:          'bg-red-100 text-red-700',
    quote_requested: 'bg-blue-100 text-blue-700',
  }[s] ?? 'bg-gray-100 text-gray-700'
}
</script>
