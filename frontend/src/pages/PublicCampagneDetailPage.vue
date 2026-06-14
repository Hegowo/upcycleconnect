<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" /> Retour à l'accueil
      </RouterLink>

      <div v-if="loading" class="py-20 text-center">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>

      <div v-else-if="!campaign" class="bg-white rounded-2xl p-12 text-center border border-[#edf4ff]">
        <MegaphoneIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <p class="text-[#001d32] font-semibold">Campagne introuvable</p>
        <p class="text-[#40617f] text-sm mt-1">Cette campagne n'existe plus ou n'est plus active.</p>
      </div>

      <template v-else>

        <div class="bg-white rounded-[24px] overflow-hidden border border-[#edf4ff] mb-8">
          <div class="aspect-[21/9] bg-gradient-to-br from-[#006d35] to-[#1b8848] relative overflow-hidden">
            <img v-if="campaign.image_url" :src="campaign.image_url" :alt="campaign.title" class="w-full h-full object-cover" />
            <MegaphoneIcon v-else class="absolute inset-0 m-auto w-16 h-16 text-white/30" />
            <span class="absolute top-4 left-4 bg-[#006d35] text-white text-xs font-bold uppercase tracking-wider px-3 py-1.5 rounded-full">
              Mis en avant
            </span>
          </div>
          <div class="p-6 sm:p-8">
            <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl sm:text-3xl tracking-tight mb-2">{{ campaign.title }}</h1>
            <p class="text-[#40617f] text-sm flex items-center gap-1.5 mb-4">
              <UserIcon class="w-4 h-4" /> Proposé par {{ campaign.provider_name || 'un artisan UpcycleConnect' }}
            </p>
            <p v-if="campaign.description" class="text-[#001d32] text-base leading-relaxed whitespace-pre-line">{{ campaign.description }}</p>
          </div>
        </div>

        <div v-if="prestations.length" class="mb-8">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-4">Prestations mises en avant</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <RouterLink
              v-for="p in prestations"
              :key="p.id"
              :to="`/prestations/${p.id}`"
              class="group bg-white rounded-2xl border border-[#edf4ff] p-5 hover:border-[#006d35] hover:shadow-md transition flex flex-col"
            >
              <span v-if="p.category" class="text-[10px] font-bold uppercase tracking-wider text-[#006d35] mb-1">{{ p.category.name }}</span>
              <h3 class="font-jakarta font-bold text-[#001d32] text-base group-hover:text-[#006d35] transition line-clamp-2">{{ p.title }}</h3>
              <p v-if="p.description" class="text-[#40617f] text-xs mt-1 line-clamp-2">{{ p.description }}</p>
              <p class="text-[#006d35] font-semibold text-sm mt-3">{{ p.price ? formatEUR(p.price) : 'Sur devis' }}</p>
            </RouterLink>
          </div>
        </div>

        <div v-if="events.length" class="mb-8">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-4">Événements mis en avant</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <RouterLink
              v-for="e in events"
              :key="e.id"
              :to="`/evenements/${e.id}`"
              class="group bg-white rounded-2xl border border-[#edf4ff] p-5 hover:border-[#006d35] hover:shadow-md transition flex flex-col"
            >
              <span class="text-[10px] font-bold uppercase tracking-wider text-[#006d35] mb-1 flex items-center gap-1">
                <CalendarDaysIcon class="w-3.5 h-3.5" /> {{ formatDate(e.start_date) }}
              </span>
              <h3 class="font-jakarta font-bold text-[#001d32] text-base group-hover:text-[#006d35] transition line-clamp-2">{{ e.title }}</h3>
              <p v-if="e.description" class="text-[#40617f] text-xs mt-1 line-clamp-2">{{ e.description }}</p>
            </RouterLink>
          </div>
        </div>

        <div v-if="!prestations.length && !events.length" class="bg-white rounded-2xl border border-[#edf4ff] p-8 text-center">
          <p class="text-[#40617f] text-sm">Cette campagne ne référence pas encore de prestations ou d'événements.</p>
          <RouterLink to="/prestations" class="inline-block mt-3 text-[#006d35] font-semibold text-sm hover:underline">
            Découvrir toutes les prestations →
          </RouterLink>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeftIcon, MegaphoneIcon, UserIcon, CalendarDaysIcon } from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

const route = useRoute()
const campaign = ref(null)
const prestations = ref([])
const events = ref([])
const loading = ref(true)

async function load(id) {
  loading.value = true
  try {
    const data = await publicGet(`/campaigns/${id}`)
    campaign.value = data.campaign
    prestations.value = data.prestations || []
    events.value = data.events || []
  } catch {
    campaign.value = null
  } finally {
    loading.value = false
  }
}

function formatEUR(price) {
  const n = parseFloat(price)
  return isNaN(n) ? price : new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(n)
}
function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(() => load(route.params.id))
watch(() => route.params.id, (id) => id && load(id))
</script>
