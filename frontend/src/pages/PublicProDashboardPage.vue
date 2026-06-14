<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/profil/pro" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" /> Retour profil pro
      </RouterLink>

      <div class="flex items-center justify-between flex-wrap gap-3 mb-8">
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Tableau de bord</h1>
          <p class="text-[#40617f] text-sm mt-1">Suivez votre activité et l'impact de votre travail d'upcycling.</p>
        </div>
        <span v-if="data?.subscription?.active" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-bold text-white" style="background:linear-gradient(135deg,#006d35,#1b8848);">
          <StarIcon class="w-3.5 h-3.5" /> Abonnement {{ planLabel }}
        </span>
      </div>

      <div v-if="loading" class="py-20 text-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>

      <template v-else-if="data">

        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          <div class="bg-white rounded-2xl border border-[#edf4ff] p-5">
            <ClipboardDocumentListIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ data.basic.prestations_count }}</p>
            <p class="text-[#40617f] text-xs">Prestations ({{ data.basic.active_prestations }} publiées)</p>
          </div>
          <div class="bg-white rounded-2xl border border-[#edf4ff] p-5">
            <WrenchScrewdriverIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ data.basic.projects_count }}</p>
            <p class="text-[#40617f] text-xs">Projets d'upcycling</p>
          </div>
          <div class="bg-white rounded-2xl border border-[#edf4ff] p-5">
            <MegaphoneIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ data.basic.campaigns_count }}</p>
            <p class="text-[#40617f] text-xs">Campagnes</p>
          </div>
          <div class="bg-white rounded-2xl border border-[#edf4ff] p-5">
            <ChartBarIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ data.advanced ? formatEUR(data.advanced.revenue_total_cents) : '—' }}</p>
            <p class="text-[#40617f] text-xs">Chiffre d'affaires</p>
          </div>
        </div>

        <div class="relative mb-8">
          <div v-if="!data.advanced" class="absolute inset-0 z-10 rounded-[24px] backdrop-blur-sm bg-white/60 flex items-center justify-center">
            <PaywallCard
              title="Tableaux de bord avancés"
              desc="Suivez votre chiffre d'affaires, vos réservations et les statistiques sur les matériaux disponibles."
              tier="Basic"
            />
          </div>
          <div :class="!data.advanced ? 'pointer-events-none select-none opacity-60' : ''">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

              <div class="bg-white rounded-[24px] border border-[#edf4ff] p-6">
                <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-4">Chiffre d'affaires (6 mois)</h2>
                <div v-if="advancedMonthly.length" class="flex items-end gap-2 h-40">
                  <div v-for="m in advancedMonthly" :key="m.month" class="flex-1 flex flex-col items-center gap-1">
                    <div class="w-full rounded-t-lg transition-all" style="background:linear-gradient(180deg,#006d35,#1b8848);"
                      :style="{ height: barHeight(m.cents) }" :title="formatEUR(m.cents)" />
                    <span class="text-[10px] text-[#94a3b8]">{{ shortMonth(m.month) }}</span>
                  </div>
                </div>
                <p v-else class="text-[#94a3b8] text-sm py-10 text-center">Pas encore de revenus enregistrés.</p>
                <div class="mt-4 flex items-center justify-between text-sm pt-3 border-t border-[#f1f5f9]">
                  <span class="text-[#40617f]">Réservations payées</span>
                  <span class="font-semibold text-[#001d32]">{{ data.advanced?.reservations_paid || 0 }} / {{ data.advanced?.reservations_total || 0 }}</span>
                </div>
              </div>

              <div class="bg-white rounded-[24px] border border-[#edf4ff] p-6">
                <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-4">Matériaux disponibles</h2>
                <div v-if="advancedMaterials.length" class="space-y-3">
                  <div v-for="m in advancedMaterials" :key="m.category" class="flex items-center justify-between gap-3">
                    <span class="text-sm text-[#001d32] truncate">{{ m.category }}</span>
                    <div class="flex items-center gap-2 shrink-0">
                      <span class="text-xs text-[#94a3b8]">{{ m.weight_kg }} kg</span>
                      <span class="text-sm font-bold text-[#006d35] bg-[#f0fdf4] px-2 py-0.5 rounded-full">{{ m.count }}</span>
                    </div>
                  </div>
                </div>
                <p v-else class="text-[#94a3b8] text-sm py-10 text-center">Aucun matériau disponible actuellement.</p>
                <p class="text-[10px] text-[#94a3b8] mt-4 pt-3 border-t border-[#f1f5f9]">Objets validés en conteneur, prêts à être récupérés.</p>
              </div>
            </div>
          </div>
        </div>

        <div class="relative">
          <div v-if="!data.premium" class="absolute inset-0 z-10 rounded-[24px] backdrop-blur-sm bg-white/60 flex items-center justify-center">
            <PaywallCard
              title="Analyse d'impact & alertes collecte"
              desc="Mesurez votre impact écologique détaillé et recevez des alertes priorisées sur les objets à collecter."
              tier="Premium"
            />
          </div>
          <div :class="!data.premium ? 'pointer-events-none select-none opacity-60' : ''">
            <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

              <div class="bg-gradient-to-br from-[#006d35] to-[#1b8848] rounded-[24px] p-6 text-white">
                <h2 class="font-jakarta font-bold text-lg mb-4 flex items-center gap-2">
                  <GlobeAltIcon class="w-5 h-5" /> Impact écologique
                </h2>
                <div class="grid grid-cols-3 gap-3">
                  <div class="text-center">
                    <p class="font-jakarta font-extrabold text-3xl">{{ premiumEco.co2 }}</p>
                    <p class="text-white/80 text-xs mt-1">kg CO₂ évités</p>
                  </div>
                  <div class="text-center">
                    <p class="font-jakarta font-extrabold text-3xl">{{ premiumEco.weight }}</p>
                    <p class="text-white/80 text-xs mt-1">kg revalorisés</p>
                  </div>
                  <div class="text-center">
                    <p class="font-jakarta font-extrabold text-3xl">{{ premiumEco.items }}</p>
                    <p class="text-white/80 text-xs mt-1">objets collectés</p>
                  </div>
                </div>
              </div>

              <div class="bg-white rounded-[24px] border border-[#edf4ff] p-6">
                <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-4 flex items-center gap-2">
                  <BellAlertIcon class="w-5 h-5 text-[#006d35]" /> Alertes collecte
                </h2>
                <div v-if="premiumAlerts.length" class="space-y-2 max-h-52 overflow-y-auto">
                  <div v-for="a in premiumAlerts" :key="a.id" class="flex items-center justify-between gap-3 p-2 rounded-lg hover:bg-[#f8fafc]">
                    <div class="min-w-0">
                      <p class="text-sm font-medium text-[#001d32] truncate">{{ a.title }}</p>
                      <p class="text-xs text-[#94a3b8]">{{ a.category || 'Non catégorisé' }} · {{ a.collection_point || 'Point inconnu' }}</p>
                    </div>
                    <RouterLink to="/profil/pro/collecte" class="text-xs font-bold text-[#006d35] hover:underline shrink-0">Collecter</RouterLink>
                  </div>
                </div>
                <p v-else class="text-[#94a3b8] text-sm py-8 text-center">Aucun objet à collecter pour l'instant.</p>
              </div>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, h } from 'vue'
import { RouterLink } from 'vue-router'
import {
  ArrowLeftIcon, StarIcon, ClipboardDocumentListIcon, WrenchScrewdriverIcon,
  MegaphoneIcon, ChartBarIcon, GlobeAltIcon, BellAlertIcon, LockClosedIcon,
} from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const PaywallCard = {
  props: ['title', 'desc', 'tier'],
  setup(props) {
    return () => h('div', { class: 'bg-white rounded-2xl shadow-lg border-2 border-[#006d35] p-6 max-w-sm text-center mx-4' }, [
      h(LockClosedIcon, { class: 'w-10 h-10 text-[#006d35] mx-auto mb-3' }),
      h('p', { class: 'font-jakarta font-bold text-[#001d32] text-base mb-1' }, props.title),
      h('p', { class: 'text-[#40617f] text-sm mb-4' }, props.desc),
      h(RouterLink, {
        to: '/profil/pro/abonnement',
        class: 'inline-block px-5 py-2.5 rounded-xl text-white font-bold text-sm',
        style: 'background:linear-gradient(135deg,#006d35,#1b8848);',
      }, () => `Débloquer avec ${props.tier} →`),
    ])
  },
}

const data    = ref(null)
const loading  = ref(true)

const planLabel = computed(() => {
  const p = data.value?.subscription?.plan
  return p === 'premium' ? 'Premium' : p === 'basic' ? 'Basic' : ''
})
const advancedMonthly   = computed(() => data.value?.advanced?.monthly_revenue || [])
const advancedMaterials = computed(() => data.value?.advanced?.material_stats || [])
const premiumEco = computed(() => {
  const e = data.value?.premium?.eco_impact
  return {
    co2: e ? Math.round(e.co2_saved_kg) : 0,
    weight: e ? Math.round(e.weight_saved_kg) : 0,
    items: e ? e.items_collected : 0,
  }
})
const premiumAlerts = computed(() => data.value?.premium?.collection_alerts || [])

const maxRevenue = computed(() => Math.max(...advancedMonthly.value.map(m => m.cents), 1))
function barHeight(cents) {
  const pct = Math.max((cents / maxRevenue.value) * 100, 4)
  return `${pct}%`
}
function shortMonth(ym) {
  const [, m] = ym.split('-')
  return ['Jan','Fév','Mar','Avr','Mai','Jun','Jul','Aoû','Sep','Oct','Nov','Déc'][parseInt(m) - 1] || ym
}
function formatEUR(cents) {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR', maximumFractionDigits: 0 }).format((cents || 0) / 100)
}

onMounted(async () => {
  try {
    data.value = await userApi('/provider/dashboard')
  } catch {}
  loading.value = false
})
</script>
