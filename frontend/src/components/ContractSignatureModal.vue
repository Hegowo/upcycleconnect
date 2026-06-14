<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4">

      <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="cancel" />

      <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-4xl max-h-[92vh] flex flex-col overflow-hidden">

        <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between shrink-0" style="background: linear-gradient(135deg, #006d35, #1b8848);">
          <div>
            <h2 class="text-white font-jakarta font-bold text-lg">Contrat de prestation</h2>
            <p class="text-white/70 text-xs mt-0.5">À lire attentivement et à signer avant le paiement</p>
          </div>
          <button @click="cancel" class="text-white/70 hover:text-white transition p-2 rounded-lg hover:bg-white/10">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div v-if="loading" class="flex-1 flex items-center justify-center py-20">
          <div class="w-12 h-12 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin"></div>
        </div>

        <div v-else-if="errorMsg" class="flex-1 flex flex-col items-center justify-center gap-3 py-20 px-6 text-center">
          <ExclamationTriangleIcon class="w-12 h-12 text-red-500" />
          <p class="text-red-600 font-medium">{{ errorMsg }}</p>
          <button @click="cancel" class="mt-2 px-4 py-2 rounded-lg border border-gray-200 text-sm">Fermer</button>
        </div>

        <div v-else class="flex-1 overflow-y-auto px-6 py-5 contract-body">

          <div class="grid grid-cols-2 gap-4 text-xs mb-6 pb-4 border-b border-gray-100">
            <div>
              <p class="uppercase text-gray-400 font-semibold mb-1 tracking-wide">Numéro</p>
              <p class="text-gray-800 font-mono">CONTRAT-{{ year }}-{{ String(data.prestation.id).padStart(6, '0') }}-* </p>
            </div>
            <div>
              <p class="uppercase text-gray-400 font-semibold mb-1 tracking-wide">Date</p>
              <p class="text-gray-800">{{ today }}</p>
            </div>
          </div>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 1 — Entre les soussignés</h3>

          <div class="grid sm:grid-cols-2 gap-4 mb-6">
            <div class="p-4 rounded-xl bg-[#f0fdf4] border border-[#bbf7d0]">
              <p class="text-xs uppercase text-[#006d35] font-semibold mb-2">Le Prestataire</p>
              <p class="font-semibold text-gray-900 text-sm">UpcycleConnect SAS</p>
              <p class="text-xs text-gray-600 mt-1">174, rue La Fayette<br/>75010 Paris</p>
              <p class="text-xs text-gray-500 mt-2">contact@upcycleconnect.xyz</p>
            </div>
            <div class="p-4 rounded-xl bg-[#f8fafc] border border-gray-200">
              <p class="text-xs uppercase text-gray-500 font-semibold mb-2">Le Client</p>
              <p class="font-semibold text-gray-900 text-sm">{{ data.customer.name }}</p>
              <p class="text-xs text-gray-600 mt-1">{{ data.customer.email }}</p>
              <p v-if="data.customer.phone" class="text-xs text-gray-500 mt-1">Tél : {{ data.customer.phone }}</p>
              <p v-if="data.customer.address" class="text-xs text-gray-500 mt-1">{{ data.customer.address }}</p>
            </div>
          </div>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 2 — Objet du contrat</h3>
          <p class="text-sm text-gray-700 leading-relaxed mb-3">
            Le présent contrat a pour objet la prestation de service suivante, commandée par le Client auprès du Prestataire via la plateforme UpcycleConnect :
          </p>
          <div class="pl-4 border-l-4 border-[#006d35] mb-4 py-1">
            <p class="font-semibold text-gray-900">{{ data.prestation.title }}</p>
            <p v-if="data.prestation.description" class="text-sm text-gray-600 mt-1 whitespace-pre-line">{{ data.prestation.description }}</p>
          </div>
          <p v-if="data.provider?.name && data.provider.name !== 'UpcycleConnect'" class="text-xs text-gray-500 mb-6">
            Prestation fournie par : <span class="font-medium text-gray-700">{{ data.provider.name }}</span>
          </p>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 3 — Prix et modalités de paiement</h3>
          <p class="text-xs text-gray-500 mb-1">Montant total TTC :</p>
          <p class="text-3xl font-bold text-[#006d35] mb-3">{{ formatEUR(data.amount_cents) }}</p>
          <p class="text-sm text-gray-700 leading-relaxed mb-6">
            Le paiement est effectué en une seule fois, en ligne, via Stripe au moment de la signature du présent contrat. La TVA française au taux normal de 20 % est incluse dans le montant total. Une facture définitive est émise après réception du paiement.
          </p>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 4 — Exécution de la prestation</h3>
          <p class="text-sm text-gray-700 leading-relaxed mb-4">
            Le Prestataire s'engage à exécuter la prestation avec soin et dans le respect des règles de l'art. Les modalités précises de réalisation (date, lieu, durée) sont convenues entre les parties via la messagerie et l'agenda intégrés à la plateforme.
          </p>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 5 — Droit de rétractation</h3>
          <p class="text-sm text-gray-700 leading-relaxed mb-4">
            Conformément aux articles L.221-18 et suivants du Code de la consommation, le Client dispose d'un délai de quatorze (14) jours à compter de la conclusion du présent contrat pour exercer son droit de rétractation. Toutefois, si l'exécution de la prestation commence, à la demande expresse du Client, avant l'expiration de ce délai, le Client renonce expressément à son droit de rétractation pour la part de prestation déjà exécutée.
          </p>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 6 — Protection des données personnelles</h3>
          <p class="text-sm text-gray-700 leading-relaxed mb-4">
            Les données personnelles du Client sont traitées conformément au RGPD. Elles sont collectées pour l'exécution du présent contrat. Le Client dispose d'un droit d'accès, de rectification, d'effacement et de portabilité de ses données via dpo@upcycleconnect.xyz.
          </p>

          <h3 class="font-jakarta font-bold text-[#006d35] text-sm mb-2">Article 7 — Loi applicable et règlement des litiges</h3>
          <p class="text-sm text-gray-700 leading-relaxed mb-6">
            Le présent contrat est soumis au droit français. En cas de litige, les parties s'engagent à rechercher en priorité une solution amiable. À défaut, compétence est donnée aux tribunaux du ressort de Paris.
          </p>

          <div class="border-t border-gray-100 pt-5 mt-2">
            <h3 class="font-jakarta font-bold text-[#001d32] text-sm mb-3 flex items-center gap-2">
              <PencilIcon class="w-4 h-4 text-[#006d35]" />
              Signature du Client
            </h3>
            <p class="text-xs text-gray-500 mb-3">Signez ci-dessous avec votre souris ou votre doigt (mobile). Conformément au règlement eIDAS, cette signature a la même valeur juridique qu'une signature manuscrite.</p>

            <div class="relative rounded-xl border-2 border-dashed border-gray-300 bg-[#fafbfd] overflow-hidden" style="aspect-ratio: 4 / 1; min-height: 140px;">
              <canvas
                ref="canvas"
                class="absolute inset-0 w-full h-full cursor-crosshair touch-none"
                @mousedown="startDraw"
                @mousemove="draw"
                @mouseup="endDraw"
                @mouseleave="endDraw"
                @touchstart.prevent="startDraw"
                @touchmove.prevent="draw"
                @touchend.prevent="endDraw"
              />
              <div v-if="!hasSignature" class="absolute inset-0 flex items-center justify-center pointer-events-none text-gray-300 text-sm italic">
                Signez ici
              </div>
            </div>

            <div class="flex items-center justify-between mt-2">
              <button
                @click="clearSig"
                :disabled="!hasSignature"
                class="text-xs text-gray-500 hover:text-gray-700 underline underline-offset-2 disabled:opacity-40 disabled:no-underline"
              >
                Effacer la signature
              </button>
              <p v-if="hasSignature" class="text-xs text-[#006d35] font-medium flex items-center gap-1">
                <CheckCircleIcon class="w-4 h-4" />
                Signature enregistrée
              </p>
            </div>

            <label class="flex items-start gap-3 mt-5 cursor-pointer select-none">
              <input
                v-model="accepted"
                type="checkbox"
                class="mt-0.5 w-4 h-4 rounded border-gray-300 text-[#006d35] focus:ring-[#006d35]/30"
              />
              <span class="text-xs text-gray-700 leading-relaxed">
                Je reconnais avoir lu et accepté l'intégralité des conditions du présent contrat, et je consens à sa conclusion par voie électronique avec signature manuscrite numérisée.
              </span>
            </label>
          </div>
        </div>

        <div v-if="!loading && !errorMsg" class="px-6 py-4 border-t border-gray-100 flex items-center justify-between gap-3 shrink-0 bg-[#f8fafc]">
          <button
            @click="cancel"
            :disabled="submitting"
            class="px-5 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700 hover:bg-gray-50 transition disabled:opacity-50"
          >
            Annuler
          </button>
          <button
            @click="submit"
            :disabled="!canSubmit"
            class="flex-1 sm:flex-none px-6 py-2.5 rounded-xl font-jakarta font-bold text-white text-sm transition-all hover:opacity-90 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            style="background: linear-gradient(135deg, #006d35, #1b8848);"
          >
            <div v-if="submitting" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            <LockClosedIcon v-else class="w-4 h-4" />
            {{ submitting ? 'Préparation du paiement…' : `Signer et payer ${formatEUR(data?.amount_cents || 0)}` }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onUnmounted } from 'vue'
import {
  XMarkIcon,
  PencilIcon,
  CheckCircleIcon,
  LockClosedIcon,
  ExclamationTriangleIcon,
} from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const props = defineProps({
  show: { type: Boolean, default: false },
  prestationId: { type: [Number, String], default: null },

  reservationId: { type: [Number, String], default: null },
  notes: { type: String, default: '' },
})

const mode = computed(() => props.reservationId ? 'accept-quote' : 'reserve')

const emit = defineEmits(['close', 'signed'])

const canvas        = ref(null)
const loading       = ref(false)
const submitting    = ref(false)
const errorMsg      = ref('')
const data          = ref(null)
const accepted      = ref(false)
const hasSignature  = ref(false)

let ctx = null
let drawing = false
let lastX = 0
let lastY = 0

const today = new Date().toLocaleDateString('fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
const year  = new Date().getFullYear()

const canSubmit = computed(() => hasSignature.value && accepted.value && !submitting.value && data.value)

function formatEUR(cents) {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format((cents || 0) / 100)
}

function setupCanvas() {
  const c = canvas.value
  if (!c) return
  const rect = c.getBoundingClientRect()
  const dpr  = window.devicePixelRatio || 1
  c.width  = rect.width * dpr
  c.height = rect.height * dpr
  ctx = c.getContext('2d')
  ctx.scale(dpr, dpr)
  ctx.strokeStyle = '#001d32'
  ctx.lineWidth = 2
  ctx.lineCap   = 'round'
  ctx.lineJoin  = 'round'
}

function getPos(e) {
  const rect = canvas.value.getBoundingClientRect()
  const point = e.touches ? e.touches[0] : e
  return { x: point.clientX - rect.left, y: point.clientY - rect.top }
}

function startDraw(e) {
  if (!ctx) return
  drawing = true
  const p = getPos(e)
  lastX = p.x
  lastY = p.y
}
function draw(e) {
  if (!drawing || !ctx) return
  const p = getPos(e)
  ctx.beginPath()
  ctx.moveTo(lastX, lastY)
  ctx.lineTo(p.x, p.y)
  ctx.stroke()
  lastX = p.x
  lastY = p.y
  hasSignature.value = true
}
function endDraw() {
  drawing = false
}

function clearSig() {
  if (!ctx) return
  ctx.clearRect(0, 0, canvas.value.width, canvas.value.height)
  hasSignature.value = false
}

watch(() => props.show, async (val) => {
  if (val) {
    await fetchPreview()
    await nextTick()
    setupCanvas()
  } else {
    reset()
  }
})

function reset() {
  data.value = null
  loading.value = false
  submitting.value = false
  errorMsg.value = ''
  accepted.value = false
  hasSignature.value = false
}

async function fetchPreview() {
  if (!props.prestationId && !props.reservationId) return
  loading.value = true
  errorMsg.value = ''
  try {
    const url = mode.value === 'accept-quote'
      ? `/reservations/${props.reservationId}/quote-contract-preview`
      : `/prestations/${props.prestationId}/contract-preview`
    data.value = await userApi(url)
  } catch (e) {
    errorMsg.value = e.message || 'Impossible de charger le contrat.'
  } finally {
    loading.value = false
  }
}

function cancel() {
  if (submitting.value) return
  emit('close')
}

async function submit() {
  if (!canSubmit.value) return
  submitting.value = true
  try {
    const signature = canvas.value.toDataURL('image/png')
    const url = mode.value === 'accept-quote'
      ? `/reservations/${props.reservationId}/accept-quote`
      : `/prestations/${props.prestationId}/reserve`
    const res = await userApi(url, {
      method: 'POST',
      body: JSON.stringify({
        notes: props.notes || null,
        signature,
      }),
    })
    emit('signed', res)
  } catch (e) {
    errorMsg.value = e.message || 'Erreur lors de la signature.'
    submitting.value = false
  }
}

function onResize() {
  if (!props.show || !canvas.value) return
  setupCanvas()
}
window.addEventListener('resize', onResize)
onUnmounted(() => window.removeEventListener('resize', onResize))
</script>

<style scoped>
.contract-body { scroll-behavior: smooth; }
.contract-body::-webkit-scrollbar { width: 8px; }
.contract-body::-webkit-scrollbar-thumb { background: #cbd5e1; border-radius: 4px; }
</style>
