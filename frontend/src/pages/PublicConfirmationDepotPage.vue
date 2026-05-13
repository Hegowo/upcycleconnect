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
              <div class="w-full sm:w-40 shrink-0 flex flex-col items-center gap-2">
                <div class="w-40 h-40 rounded-[20px] bg-[#edf4ff] border-2 border-[#d8eaff] flex flex-col items-center justify-center">
                  <img v-if="qrDataUrl" :src="qrDataUrl" class="w-32 h-32" />
                  <QrCodeIcon v-else class="w-20 h-20 text-[#001d32] opacity-30" />
                </div>
                <span v-if="deposit.qr_code" class="text-[#001d32] font-mono font-bold text-xs tracking-widest text-center">
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
              <!-- Photo(s) or placeholder -->
              <div v-if="deposit.photo1 || deposit.photo2 || deposit.photo3" class="flex gap-2 overflow-x-auto">
                <img
                  v-for="src in [deposit.photo1, deposit.photo2, deposit.photo3].filter(Boolean)"
                  :key="src.slice(0,30)"
                  :src="src"
                  class="h-32 rounded-xl object-cover flex-shrink-0 cursor-pointer hover:opacity-90 transition"
                  :class="[deposit.photo1, deposit.photo2, deposit.photo3].filter(Boolean).length === 1 ? 'w-full' : 'w-32'"
                  @click="lightboxSrc = src"
                />
              </div>
              <div v-else class="h-32 rounded-xl bg-gradient-to-br from-[#fef3c7] to-[#fde68a] flex items-center justify-center">
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

            <div v-if="deposit.collection_point" class="bg-[#edf4ff] rounded-[24px] p-6 flex flex-col gap-4">
              <h3 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.depotConfirmation.collectPointTitle') }}</h3>
              <div class="flex items-start gap-3">
                <div class="w-10 h-10 bg-white rounded-xl flex items-center justify-center shrink-0">
                  <MapPinIcon class="w-5 h-5 text-[#006d35]" />
                </div>
                <div>
                  <p class="font-jakarta font-bold text-[#001d32] text-base">{{ deposit.collection_point.name }}</p>
                  <p class="text-[#40617f] text-sm mt-0.5">{{ deposit.collection_point.address }}</p>
                  <p class="text-[#40617f] text-sm">{{ deposit.collection_point.postal_code }} {{ deposit.collection_point.city }}</p>
                  <p v-if="deposit.collection_point.opening_hours" class="text-[#40617f] text-xs mt-1">{{ deposit.collection_point.opening_hours }}</p>
                  <p v-if="deposit.collection_point.phone" class="text-[#40617f] text-xs mt-0.5">{{ deposit.collection_point.phone }}</p>
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

  <!-- Photo lightbox -->
  <Teleport to="body">
    <div v-if="lightboxSrc" class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center" @click="lightboxSrc = null">
      <img :src="lightboxSrc" class="max-w-[90vw] max-h-[90vh] rounded-2xl object-contain shadow-2xl" @click.stop />
      <button @click="lightboxSrc = null" class="absolute top-4 right-4 p-2 rounded-full bg-white/20 hover:bg-white/40 transition text-white">
        <XMarkIcon class="w-6 h-6" />
      </button>
    </div>
  </Teleport>
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
  XMarkIcon,
} from '@heroicons/vue/24/outline'
import QRCode from 'qrcode'
import { userApi } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const userAuth = useUserAuthStore()

const loading = ref(true)
const deposit = ref(null)
const qrDataUrl = ref(null)
const lightboxSrc = ref(null)

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

const steps = computed(() => {
  const cp = deposit.value?.collection_point
  const qrCode = deposit.value?.qr_code || '—'
  return [
    { num: 1, title: t('public.depotConfirmation.step1Title'), desc: t('public.depotConfirmation.step1Desc') },
    {
      num: 2,
      title: t('public.depotConfirmation.step2Title'),
      desc: cp
        ? t('public.depotConfirmation.step2Desc', { name: cp.name, address: cp.address, postalCode: cp.postal_code, city: cp.city })
        : t('public.depotConfirmation.step2Desc', { name: '—', address: '', postalCode: '', city: '' }),
    },
    {
      num: 3,
      title: t('public.depotConfirmation.step3Title'),
      desc: t('public.depotConfirmation.step3Desc', { qrCode }),
    },
  ]
})

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

async function printVoucher() {
  const d = deposit.value
  if (!d) return
  const qrImg = qrDataUrl.value || await QRCode.toDataURL(d.qr_code || '', { width: 200, margin: 2 })
  const cp = d.collection_point
  const photos = [d.photo1, d.photo2, d.photo3].filter(Boolean)
  const photoHtml = photos.map(p => `<img src="${p}" style="width:120px;height:120px;object-fit:cover;border-radius:10px;" />`).join('')
  const win = window.open('', '_blank', 'width=794,height=1123')
  win.document.write(`<!DOCTYPE html><html lang="fr"><head>
<meta charset="utf-8" />
<title>Bon de Dépôt — ${d.qr_code || ''}</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: 'Segoe UI', Arial, sans-serif; background: #fff; color: #001d32; padding: 40px; }
  .header { display: flex; align-items: center; justify-content: space-between; border-bottom: 2px solid #006d35; padding-bottom: 20px; margin-bottom: 28px; }
  .logo { font-size: 22px; font-weight: 900; color: #006d35; letter-spacing: -0.5px; }
  .badge { background: #d1fae5; color: #006d35; font-size: 11px; font-weight: 700; padding: 5px 14px; border-radius: 20px; text-transform: uppercase; letter-spacing: 1px; }
  h1 { font-size: 26px; font-weight: 900; margin-bottom: 6px; }
  .subtitle { color: #40617f; font-size: 13px; margin-bottom: 28px; }
  .qr-row { display: flex; gap: 28px; align-items: flex-start; margin-bottom: 28px; }
  .qr-box { background: #edf4ff; border-radius: 16px; padding: 16px; display: flex; flex-direction: column; align-items: center; gap: 8px; min-width: 168px; }
  .qr-code-str { font-family: monospace; font-size: 11px; font-weight: 700; color: #001d32; letter-spacing: 2px; }
  .info-grid { flex: 1; display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
  .info-cell { background: #f7f9ff; border-radius: 10px; padding: 10px 14px; }
  .info-label { font-size: 10px; text-transform: uppercase; letter-spacing: 1px; color: #40617f; margin-bottom: 4px; }
  .info-value { font-size: 15px; font-weight: 700; }
  .section { margin-bottom: 22px; }
  .section-title { font-size: 11px; text-transform: uppercase; letter-spacing: 1px; color: #6b7280; font-weight: 600; margin-bottom: 10px; border-bottom: 1px solid #f1f5f9; padding-bottom: 6px; }
  .cp-box { background: #edf4ff; border-radius: 12px; padding: 14px 18px; }
  .cp-name { font-weight: 700; font-size: 15px; margin-bottom: 4px; }
  .cp-addr { color: #40617f; font-size: 13px; line-height: 1.6; }
  .photos { display: flex; gap: 10px; flex-wrap: wrap; }
  .instructions { background: #f0fdf4; border-radius: 12px; padding: 16px 20px; margin-top: 28px; }
  .instructions p { font-size: 13px; color: #166534; line-height: 1.7; }
  .instructions strong { display: block; margin-bottom: 6px; font-size: 14px; }
  .footer { margin-top: 32px; padding-top: 16px; border-top: 1px solid #e5e7eb; text-align: center; color: #94a3b8; font-size: 11px; }
</style>
</head><body>
<div class="header">
  <div class="logo">UpcycleConnect</div>
  <div class="badge">Bon de Dépôt Validé</div>
</div>
<h1>Bon de Dépôt</h1>
<p class="subtitle">Présentez ce document au point de collecte pour déposer votre objet.</p>
<div class="qr-row">
  <div class="qr-box">
    <img src="${qrImg}" width="136" height="136" />
    <span class="qr-code-str">${d.qr_code || ''}</span>
  </div>
  <div class="info-grid">
    <div class="info-cell"><div class="info-label">Référence</div><div class="info-value">#UP-${d.id}</div></div>
    <div class="info-cell"><div class="info-label">Statut</div><div class="info-value" style="color:#006d35;">Validé</div></div>
    <div class="info-cell"><div class="info-label">Validé le</div><div class="info-value">${d.validated_at ? new Date(d.validated_at).toLocaleDateString('fr-FR') : '—'}</div></div>
    <div class="info-cell"><div class="info-label">Expire le</div><div class="info-value">${(() => { const base = d.validated_at || d.created_at; if (!base) return '—'; const dt = new Date(base); dt.setMonth(dt.getMonth() + 1); return dt.toLocaleDateString('fr-FR'); })()}</div></div>
  </div>
</div>
${cp ? `<div class="section"><div class="section-title">Point de collecte</div><div class="cp-box"><div class="cp-name">${cp.name}</div><div class="cp-addr">${cp.address}<br/>${cp.postal_code} ${cp.city}${cp.opening_hours ? '<br/>' + cp.opening_hours : ''}</div></div></div>` : ''}
<div class="section"><div class="section-title">Objet</div>
<div class="info-grid">
  <div class="info-cell"><div class="info-label">Nom</div><div class="info-value">${d.title}</div></div>
  <div class="info-cell"><div class="info-label">Catégorie</div><div class="info-value">${d.category?.name || '—'}</div></div>
  ${d.estimated_weight ? `<div class="info-cell"><div class="info-label">Poids estimé</div><div class="info-value">${d.estimated_weight} kg</div></div>` : ''}
  ${d.carbon_savings ? `<div class="info-cell"><div class="info-label">CO₂ économisé</div><div class="info-value" style="color:#006d35;">${d.carbon_savings} kg</div></div>` : ''}
</div>
</div>
${photos.length ? `<div class="section"><div class="section-title">Photos de l'objet</div><div class="photos">${photoHtml}</div></div>` : ''}
<div class="instructions"><p><strong>Comment procéder ?</strong>1. Présentez ce bon (imprimé ou sur écran) au point de collecte.<br/>2. Un agent scannera le QR code pour valider votre dépôt.<br/>3. Vous recevrez une confirmation par e-mail.</p></div>
<div class="footer">UpcycleConnect — Ensemble pour une économie circulaire · Bon valable 1 mois à compter de la validation</div>
</body></html>`)
  win.document.close()
  win.focus()
  setTimeout(() => { win.print(); }, 500)
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
    if (deposit.value?.qr_code) {
      qrDataUrl.value = await QRCode.toDataURL(deposit.value.qr_code, { width: 160, margin: 1, errorCorrectionLevel: 'M' })
    }
  } catch {
    deposit.value = null
  } finally {
    loading.value = false
  }
})
</script>
