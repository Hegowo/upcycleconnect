<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/profil/pro" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" /> Retour profil pro
      </RouterLink>

      <div class="mb-8">
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Abonnement professionnel</h1>
        <p class="text-[#40617f] text-sm mt-1">Choisissez un plan pour accéder aux fonctionnalités pro d'UpcycleConnect.</p>
      </div>

      <div v-if="subscription" class="bg-white rounded-2xl border border-[#edf4ff] p-5 mb-8 flex items-center justify-between gap-4 flex-wrap">
        <div class="flex items-center gap-4">
          <div class="w-12 h-12 rounded-xl flex items-center justify-center" style="background: linear-gradient(135deg,#006d35,#1b8848);">
            <StarIcon class="w-6 h-6 text-white" />
          </div>
          <div>
            <p class="font-jakarta font-bold text-[#001d32]">Abonnement {{ subscription.plan_label }} actif</p>
            <p class="text-[#40617f] text-xs mt-0.5" v-if="subscription.current_period_end">
              Renouvellement le {{ formatDate(subscription.current_period_end) }}
            </p>
          </div>
        </div>
        <button @click="cancelSub" :disabled="cancelling"
          class="px-4 py-2 rounded-xl border border-red-200 text-red-600 text-sm font-semibold hover:bg-red-50 transition disabled:opacity-50">
          {{ cancelling ? 'Annulation...' : 'Résilier' }}
        </button>
      </div>

      <div v-if="success" class="bg-green-50 border border-green-200 rounded-2xl p-4 mb-6 flex items-center gap-3">
        <CheckCircleIcon class="w-5 h-5 text-green-600 shrink-0" />
        <p class="text-green-800 text-sm font-medium">Abonnement activé avec succès ! Merci pour votre confiance.</p>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
        <div
          v-for="plan in plans"
          :key="plan.key"
          class="bg-white rounded-[24px] border-2 flex flex-col overflow-hidden"
          :class="plan.key === 'premium' ? 'border-[#006d35] shadow-lg' : 'border-[#edf4ff]'"
        >
          <div v-if="plan.key === 'premium'" class="bg-[#006d35] text-white text-center text-xs font-bold py-1.5 tracking-widest uppercase">
            Recommandé
          </div>
          <div class="p-7 flex flex-col gap-5 flex-1">
            <div>
              <p class="text-xs font-bold uppercase tracking-widest text-[#40617f]">{{ plan.label }}</p>
              <p class="font-jakarta font-extrabold text-[#001d32] text-4xl mt-1">
                {{ formatEUR(plan.amount_cents) }}
                <span class="text-[#40617f] font-normal text-base">/mois</span>
              </p>
            </div>
            <ul class="flex flex-col gap-2.5">
              <li v-for="f in plan.features" :key="f" class="flex items-start gap-2.5 text-sm text-[#001d32]">
                <CheckCircleIcon class="w-4 h-4 text-[#006d35] shrink-0 mt-0.5" />
                {{ f }}
              </li>
            </ul>
            <button
              @click="subscribe(plan.key)"
              :disabled="loading || (subscription && subscription.status === 'active')"
              class="mt-auto w-full py-3 rounded-xl font-jakarta font-bold text-sm transition hover:opacity-90 disabled:opacity-50"
              :style="plan.key === 'premium'
                ? 'background: linear-gradient(135deg,#006d35,#1b8848); color: white;'
                : 'background: white; border: 2px solid #006d35; color: #006d35;'"
            >
              {{ loading === plan.key ? 'Chargement...' : subscription ? 'Abonnement actif' : 'Choisir ce plan' }}
            </button>
          </div>
        </div>
      </div>

      <p v-if="error" class="text-red-600 text-sm mt-6 text-center">{{ error }}</p>
      <p class="text-[#94a3b8] text-xs text-center mt-8">Paiement sécurisé par Stripe. Résiliation possible à tout moment.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeftIcon, CheckCircleIcon, StarIcon } from '@heroicons/vue/24/outline'
import { publicGet, userApi } from '@/services/publicApi'

const route = useRoute()
const plans = ref([])
const subscription = ref(null)
const loading = ref(null)
const cancelling = ref(false)
const error = ref('')
const success = ref(route.query.success === '1')

async function fetchData() {
  const [plansData, subData] = await Promise.all([
    publicGet('/subscription/plans'),
    userApi('/subscription').catch(() => ({ subscription: null })),
  ])
  plans.value = (plansData.data || []).sort((a, b) => a.amount_cents - b.amount_cents)
  subscription.value = subData.subscription
}

async function subscribe(plan) {
  loading.value = plan
  error.value = ''
  try {
    const res = await userApi('/subscription/checkout', { method: 'POST', body: JSON.stringify({ plan }) })
    if (res.checkout_url) window.location.href = res.checkout_url
    else error.value = res.message || 'Erreur lors du démarrage du paiement.'
  } catch (e) {
    error.value = e.message || 'Erreur inattendue.'
  } finally {
    loading.value = null
  }
}

async function cancelSub() {
  if (!confirm('Confirmer la résiliation de votre abonnement ?')) return
  cancelling.value = true
  error.value = ''
  try {
    await userApi('/subscription/cancel', { method: 'POST' })
    await fetchData()
  } catch (e) {
    error.value = e.message || 'Erreur lors de la résiliation.'
  } finally {
    cancelling.value = false
  }
}

function formatEUR(cents) {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100)
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}

onMounted(fetchData)
</script>
