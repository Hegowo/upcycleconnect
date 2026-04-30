<template>
  <div class="bg-[#f7f9ff] min-h-screen">
    <section class="px-4 sm:px-6 pt-10 pb-12 max-w-[1024px] mx-auto">
      <button
        @click="$router.back()"
        class="flex items-center gap-2 text-[#40617f] hover:text-[#006d35] mb-6"
      >
        <ArrowLeftIcon class="w-5 h-5" />
        {{ t('public.invoices.back') }}
      </button>

      <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl md:text-4xl mb-8">
        {{ t('public.invoices.title') }}
      </h1>

      <div class="flex gap-3 mb-6">
        <button
          @click="filter = ''"
          class="px-5 py-2 rounded-full text-sm font-semibold transition"
          :class="filter === '' ? 'bg-[#006d35] text-white' : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
        >
          {{ t('public.invoices.filterAll') }}
        </button>
        <button
          @click="filter = 'invoice'"
          class="px-5 py-2 rounded-full text-sm font-semibold transition"
          :class="filter === 'invoice' ? 'bg-[#006d35] text-white' : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
        >
          {{ t('public.invoices.filterInvoices') }}
        </button>
        <button
          @click="filter = 'quote'"
          class="px-5 py-2 rounded-full text-sm font-semibold transition"
          :class="filter === 'quote' ? 'bg-[#006d35] text-white' : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
        >
          {{ t('public.invoices.filterQuotes') }}
        </button>
      </div>

      <div v-if="loading" class="text-center py-16 text-[#40617f]">{{ t('public.invoices.loading') }}</div>
      <div v-else-if="error" class="text-center py-16 text-red-600">{{ error }}</div>
      <div v-else-if="items.length === 0" class="bg-white rounded-[24px] p-10 text-center text-[#40617f]">
        {{ t('public.invoices.empty') }}
      </div>
      <div v-else class="flex flex-col gap-3">
        <div
          v-for="inv in items"
          :key="inv.id"
          class="bg-white rounded-2xl p-5 flex items-center justify-between gap-4"
        >
          <div class="flex items-center gap-4 min-w-0">
            <div
              class="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0"
              :class="inv.type === 'quote' ? 'bg-orange-50 text-orange-500' : 'bg-[#d1fae5] text-[#006d35]'"
            >
              <DocumentTextIcon class="w-6 h-6" />
            </div>
            <div class="min-w-0">
              <p class="font-bold text-[#001d32] truncate">{{ inv.prestation_title }}</p>
              <p class="text-xs text-[#40617f]">
                {{ inv.type === 'quote' ? t('public.invoices.labelQuote') : t('public.invoices.labelInvoice') }} {{ t('public.invoices.numberPrefix') }} {{ inv.number }} •
                {{ formatDate(inv.issued_at) }}
              </p>
            </div>
          </div>
          <div class="flex items-center gap-4 flex-shrink-0">
            <span class="font-bold text-[#006d35] hidden sm:inline">
              {{ formatAmount(inv.amount_cents) }}
            </span>
            <button
              v-if="inv.has_pdf"
              @click="download(inv)"
              class="flex items-center gap-1 px-4 py-2 rounded-xl bg-[#cee5ff] text-[#001d32] font-semibold text-sm hover:bg-[#b8d8ff]"
            >
              <ArrowDownTrayIcon class="w-4 h-4" />
              {{ t('public.invoices.downloadButton') }}
            </button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ArrowLeftIcon, DocumentTextIcon, ArrowDownTrayIcon } from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const { t, locale } = useI18n()
const items = ref([])
const loading = ref(false)
const error = ref('')
const filter = ref('')

function formatDate(s) {
  if (!s) return ''
  try {
    const localeCode = locale.value === 'en' ? 'en-US' : 'fr-FR'
    return new Date(s).toLocaleDateString(localeCode, { day: '2-digit', month: 'long', year: 'numeric' })
  } catch {
    return s
  }
}
function formatAmount(cents) {
  if (!cents && cents !== 0) return '—'
  return `${(cents / 100).toFixed(2)} €`
}

async function load() {
  loading.value = true
  error.value = ''
  try {
    const path = filter.value ? `/invoices?type=${filter.value}` : '/invoices'
    const res = await userApi(path)
    items.value = res.data || []
  } catch (e) {
    error.value = e.message || t('public.invoices.loadError')
  } finally {
    loading.value = false
  }
}

async function download(inv) {
  const token = localStorage.getItem('user_token')
  const res = await fetch(`/api/v1/invoices/${inv.id}/download`, {
    headers: { Authorization: `Bearer ${token}` },
  })
  if (!res.ok) {
    alert(t('public.invoices.downloadError'))
    return
  }
  const blob = await res.blob()
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${inv.number}.pdf`
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

watch(filter, load)
onMounted(load)
</script>
