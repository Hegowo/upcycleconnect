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

      <div v-else-if="reservation" class="grid grid-cols-12 gap-6 items-start">

        <div class="col-span-12 lg:col-span-7 space-y-6">

        <div class="bg-white rounded-[24px] p-6 sm:p-8 relative overflow-hidden">
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

        <div class="bg-white rounded-[24px] p-6 sm:p-8">
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

        <div class="space-y-6">

          <div v-if="!isProviderView" class="bg-[#edf4ff] rounded-[24px] p-6">
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

          <div v-if="isProviderView && canEditQuote" class="bg-white rounded-[24px] p-6 border border-[#edf4ff]">
            <h3 class="font-jakarta font-bold text-[#001d32] text-base mb-1">Établir le devis</h3>
            <p class="text-[#40617f] text-xs mb-4">Ajoutez une ligne par élément demandé. Le total TTC est calculé automatiquement.</p>

            <div class="space-y-2 mb-2">
              <div v-for="(line, i) in quoteLines" :key="i" class="flex items-center gap-2">
                <input v-model="line.label" type="text" placeholder="Désignation"
                  class="flex-1 min-w-0 px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
                <input v-model.number="line.amount" type="number" min="0" step="0.01" placeholder="€"
                  class="w-24 px-3 py-2 bg-[#f8fafc] border border-[#e5e7eb] rounded-lg text-sm text-right focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
                <button @click="removeLine(i)" class="text-[#94a3b8] hover:text-red-500 shrink-0 p-1" title="Supprimer">
                  <XMarkIcon class="w-4 h-4" />
                </button>
              </div>
            </div>
            <button @click="addLine" class="text-sm font-semibold text-[#006d35] hover:underline mb-3">+ Ajouter une ligne</button>

            <div class="flex items-center justify-between border-t border-[#edf4ff] pt-3 mb-3">
              <span class="text-sm font-semibold text-[#001d32]">Total TTC</span>
              <span class="text-lg font-bold text-[#006d35]">{{ formatAmount(Math.round(quoteTotal * 100), 'eur') }}</span>
            </div>

            <textarea v-model="quoteMessage" rows="2" placeholder="Note pour le client (optionnel)…"
              class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm mb-3 resize-none focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            <button @click="sendQuote" :disabled="quoteSending"
              class="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm transition hover:opacity-90 disabled:opacity-60"
              style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);">
              {{ quoteSending ? 'Envoi…' : 'Envoyer le devis au client' }}
            </button>
            <p v-if="quoteFeedback" class="text-sm mt-2" :class="quoteFeedback.includes('✓') ? 'text-[#006d35]' : 'text-red-600'">{{ quoteFeedback }}</p>
          </div>

          <div v-else-if="isProviderView && contract" class="bg-[#f0fdf4] rounded-[24px] p-6 border border-[#bbf7d0]">
            <p class="font-semibold text-[#166534] text-sm">Devis signé par le client ✓</p>
            <p class="text-[#40617f] text-xs mt-1">Montant accepté : {{ formatAmount(reservation.amount_cents, reservation.currency) }}</p>
          </div>

          <div v-if="quote && quote.lines && quote.lines.length && !canEditQuote" class="bg-white rounded-[24px] p-6 border border-[#edf4ff]">
            <h3 class="font-jakarta font-bold text-[#001d32] text-base mb-3">Détail du devis</h3>
            <div class="divide-y divide-[#edf4ff]">
              <div v-for="(l, i) in quote.lines" :key="i" class="flex items-center justify-between py-2 gap-3">
                <span class="text-sm text-[#001d32]">{{ l.label }}</span>
                <span class="text-sm font-medium text-[#001d32] shrink-0">{{ formatAmount(l.amount_cents, 'eur') }}</span>
              </div>
            </div>
            <div class="flex items-center justify-between border-t-2 border-[#edf4ff] pt-3 mt-2">
              <span class="font-bold text-[#001d32]">Total TTC</span>
              <span class="font-bold text-[#006d35] text-lg">{{ formatAmount(quote.amount_cents, 'eur') }}</span>
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

            <button
              v-if="canAcceptQuote"
              @click="showQuoteSign = true"
              class="w-full mt-3 inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm transition hover:opacity-90"
              style="background: linear-gradient(134deg, #001d32 0%, #003e6b 100%);"
            >
              <PencilSquareIcon class="w-4 h-4" />
              Accepter et signer ce devis
            </button>

            <router-link
              to="/profil/factures"
              class="block text-center text-xs text-[#40617f] hover:text-[#006d35] mt-3"
            >
              {{ t('public.reservation.seeAllInvoices') }}
            </router-link>
          </div>

          <ContractSignatureModal
            v-if="canAcceptQuote"
            :show="showQuoteSign"
            :reservation-id="reservation.id"
            @close="showQuoteSign = false"
            @signed="onQuoteSigned"
          />

          <div v-if="contract" class="bg-white rounded-[24px] p-6 border border-[#edf4ff]">
            <h3 class="font-jakarta font-bold text-[#001d32] text-base mb-3 flex items-center gap-2">
              <PencilSquareIcon class="w-5 h-5 text-[#006d35]" />
              Contrat de prestation
            </h3>
            <p class="text-[#40617f] text-xs mb-1">Numéro de contrat</p>
            <p class="font-semibold text-[#001d32] text-sm mb-1 font-mono">{{ contract.number }}</p>
            <p class="text-[#40617f] text-xs">Signé le {{ formatContractDate(contract.signed_at) }}</p>
            <button
              v-if="contract.has_pdf"
              @click="downloadContract"
              :disabled="downloadingContract"
              class="w-full mt-4 inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl border-2 border-[#006d35] text-[#006d35] font-bold text-sm transition hover:bg-[#f0fdf4] disabled:opacity-60"
            >
              <ArrowDownTrayIcon class="w-4 h-4" />
              {{ downloadingContract ? 'Téléchargement…' : 'Télécharger le contrat' }}
            </button>
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

        <div class="col-span-12 lg:col-span-5">
        <div class="lg:sticky lg:top-6 bg-white rounded-[24px] p-6 flex flex-col" style="height: calc(100vh - 6rem); max-height: 720px;">
          <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-1 flex items-center gap-2 shrink-0">
            <ChatBubbleLeftRightIcon class="w-5 h-5 text-[#006d35]" />
            Messagerie
          </h2>
          <p class="text-[#40617f] text-xs mb-4 shrink-0">Échangez directement avec {{ isProviderView ? 'le client' : 'le prestataire' }}.</p>
          <div ref="chatBox" class="flex-1 space-y-3 overflow-y-auto px-1 py-2">
            <p v-if="!messages.length" class="text-center text-[#94a3b8] text-sm py-8">Aucun message pour l'instant. Démarrez la conversation.</p>
            <div v-for="m in messages" :key="m.id" class="flex" :class="m.sender_id === userAuth.user?.id ? 'justify-end' : 'justify-start'">
              <div class="max-w-[75%] rounded-2xl px-4 py-2"
                :class="m.sender_id === userAuth.user?.id ? 'bg-[#006d35] text-white' : 'bg-[#f1f5f9] text-[#001d32]'">
                <p v-if="m.sender_id !== userAuth.user?.id" class="text-[10px] font-semibold opacity-70 mb-0.5">{{ m.sender_name }}</p>
                <p class="text-sm whitespace-pre-wrap break-words">{{ m.content }}</p>
                <p class="text-[10px] mt-0.5" :class="m.sender_id === userAuth.user?.id ? 'text-white/70' : 'text-[#94a3b8]'">{{ formatChatTime(m.created_at) }}</p>
              </div>
            </div>
          </div>
          <form @submit.prevent="sendMessage" class="flex items-end gap-2 mt-4">
            <textarea v-model="newMessage" rows="1" placeholder="Votre message…"
              @keydown.enter.exact.prevent="sendMessage"
              class="flex-1 px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm resize-none focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            <button type="submit" :disabled="chatSending || !newMessage.trim()"
              class="px-4 py-2.5 rounded-xl text-white font-bold text-sm transition hover:opacity-90 disabled:opacity-40 shrink-0"
              style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);">
              Envoyer
            </button>
          </form>
        </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import ContractSignatureModal from '@/components/ContractSignatureModal.vue'
import {
  ArrowLeftIcon,
  ArrowDownTrayIcon,
  BanknotesIcon,
  CalendarDaysIcon,
  ChatBubbleLeftRightIcon,
  CheckCircleIcon,
  ChevronRightIcon,
  ClockIcon,
  DocumentTextIcon,
  ExclamationTriangleIcon,
  PencilSquareIcon,
  UserIcon,
  WrenchScrewdriverIcon,
  XMarkIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const route    = useRoute()
const router   = useRouter()
const userAuth = useUserAuthStore()

const loading     = ref(true)
const error       = ref('')
const reservation = ref(null)
const invoice     = ref(null)
const contract    = ref(null)
const downloading = ref(false)
const downloadingContract = ref(false)
const showQuoteSign = ref(false)
const isProviderView = ref(false)

const quote = ref(null)
const quoteLines = ref([{ label: '', amount: null }])
const quoteMessage = ref('')
const quoteSending = ref(false)
const quoteFeedback = ref('')
const quoteTotal = computed(() =>
  quoteLines.value.reduce((s, l) => s + (Number(l.amount) || 0), 0)
)
const canEditQuote = computed(() =>
  isProviderView.value && !contract.value &&
  ['quote_requested', 'quote_issued'].includes(reservation.value?.status)
)
function addLine() { quoteLines.value.push({ label: '', amount: null }) }
function removeLine(i) { quoteLines.value.splice(i, 1); if (!quoteLines.value.length) addLine() }
function prefillQuoteLines() {
  if (quote.value?.lines?.length) {
    quoteLines.value = quote.value.lines.map(l => ({ label: l.label, amount: (l.amount_cents || 0) / 100 }))
  }
}
async function sendQuote() {
  const lines = quoteLines.value
    .map(l => ({ label: (l.label || '').trim(), amount: Number(l.amount) || 0 }))
    .filter(l => l.label && l.amount > 0)
  if (!lines.length) { quoteFeedback.value = 'Ajoute au moins une ligne avec un montant.'; return }
  quoteSending.value = true; quoteFeedback.value = ''
  try {
    const res = await fetch(`/api/v1/reservations/${route.params.id}/quote`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify({ lines, message: quoteMessage.value || null }),
    })
    const j = await res.json()
    if (!res.ok) throw new Error(j.message || 'Erreur')
    quoteFeedback.value = 'Devis envoyé au client ✓'
    if (reservation.value) { reservation.value.status = 'quote_issued'; reservation.value.amount_cents = j.amount_cents }
    quote.value = { ...(quote.value || {}), amount_cents: j.amount_cents, lines: lines.map(l => ({ label: l.label, amount_cents: Math.round(l.amount * 100) })) }
  } catch (e) { quoteFeedback.value = e.message } finally { quoteSending.value = false }
}

const messages = ref([])
const newMessage = ref('')
const chatSending = ref(false)
const chatBox = ref(null)
function scrollChat() { if (chatBox.value) chatBox.value.scrollTop = chatBox.value.scrollHeight }
async function fetchMessages() {
  try {
    const res = await fetch(`/api/v1/reservations/${route.params.id}/messages`, {
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) { const j = await res.json(); messages.value = j.data || []; await nextTick(); scrollChat() }
  } catch {}
}
async function sendMessage() {
  const content = newMessage.value.trim()
  if (!content) return
  chatSending.value = true
  try {
    const res = await fetch(`/api/v1/reservations/${route.params.id}/messages`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify({ content }),
    })
    const j = await res.json()
    if (res.ok && j.data) { messages.value.push(j.data); newMessage.value = ''; await nextTick(); scrollChat() }
  } catch {} finally { chatSending.value = false }
}

const canAcceptQuote = computed(() => {
  if (!reservation.value || !invoice.value) return false
  if (invoice.value.type !== 'quote') return false
  if (contract.value) return false
  const s = reservation.value.status
  return s === 'quote_requested' || s === 'quote_issued'
})

function onQuoteSigned(res) {
  if (res?.checkout_url) {
    window.location.href = res.checkout_url
  }
}

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
      reservation.value    = data.reservation
      invoice.value        = data.invoice
      quote.value          = data.quote || null
      isProviderView.value = !!data.is_provider_view
      prefillQuoteLines()

      try {
        const cRes = await fetch(`/api/v1/reservations/${route.params.id}/contract`, {
          headers: { Authorization: `Bearer ${userAuth.token}` },
        })
        if (cRes.ok) contract.value = await cRes.json()
      } catch {}
      await fetchMessages()
    }
  } catch (e) {
    error.value = t('public.reservation.networkError')
  }
  loading.value = false
})

async function downloadContract() {
  if (!contract.value?.id) return
  downloadingContract.value = true
  try {
    const res = await fetch(`/api/v1/contracts/${contract.value.id}/download`, {
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (!res.ok) throw new Error('download failed')
    const blob = await res.blob()
    const url  = URL.createObjectURL(blob)
    const a    = document.createElement('a')
    a.href     = url
    a.download = `${contract.value.number}.pdf`
    a.click()
    URL.revokeObjectURL(url)
  } catch {
  } finally {
    downloadingContract.value = false
  }
}

function formatContractDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: '2-digit', month: 'long', year: 'numeric' })
}

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

function formatChatTime(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleString(localeCode.value, { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' })
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
