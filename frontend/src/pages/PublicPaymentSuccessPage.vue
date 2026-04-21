<template>
  <div class="bg-[#f7f9ff] min-h-screen flex items-center justify-center px-4">
    <div class="bg-white rounded-[32px] p-10 max-w-[560px] w-full text-center shadow-sm">
      <div class="w-20 h-20 mx-auto rounded-full bg-[#d1fae5] flex items-center justify-center mb-6">
        <CheckCircleIcon class="w-12 h-12 text-[#006d35]" />
      </div>
      <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl mb-3">
        Paiement confirmé !
      </h1>
      <p class="text-[#40617f] text-base leading-7 mb-6">
        Merci pour votre confiance. Votre facture vous a été envoyée par email et est également
        disponible dans votre espace personnel.
      </p>

      <div v-if="loading" class="text-[#40617f] text-sm mb-6">Chargement...</div>
      <div v-else-if="reservation" class="bg-[#edf4ff] rounded-2xl p-5 text-left text-sm mb-6">
        <p class="text-[#40617f] mb-1">Prestation</p>
        <p class="font-bold text-[#001d32] mb-3">{{ reservation.prestation?.title }}</p>
        <p class="text-[#40617f] mb-1">Montant réglé</p>
        <p class="font-bold text-[#006d35]">{{ formatAmount(reservation.amount_cents) }}</p>
      </div>

      <div class="flex flex-col sm:flex-row gap-3 justify-center">
        <router-link
          to="/profil/factures"
          class="px-6 py-3 rounded-xl text-white font-bold transition hover:opacity-90"
          style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
        >
          Voir mes factures
        </router-link>
        <router-link
          to="/prestations"
          class="px-6 py-3 rounded-xl bg-[#cee5ff] text-[#001d32] font-bold hover:bg-[#b8d8ff]"
        >
          Retour aux prestations
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { CheckCircleIcon } from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const route = useRoute()
const reservation = ref(null)
const loading = ref(false)

function formatAmount(cents) {
  if (!cents && cents !== 0) return '—'
  return `${(cents / 100).toFixed(2)} €`
}

onMounted(async () => {
  const sessionId = route.query.session_id
  if (!sessionId) return
  loading.value = true
  for (let i = 0; i < 5; i++) {
    try {
      const r = await userApi(`/payments/session?session_id=${encodeURIComponent(sessionId)}`)
      reservation.value = r
      if (r.status === 'paid') break
    } catch {}
    await new Promise((r) => setTimeout(r, 1500))
  }
  loading.value = false
})
</script>
