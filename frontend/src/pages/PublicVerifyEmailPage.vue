<template>
  <div class="min-h-[70vh] flex items-center justify-center py-16 px-4">
    <div class="w-full max-w-md text-center">

      
      <div v-if="state === 'loading'" class="space-y-4">
        <div class="w-14 h-14 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto"></div>
        <p class="text-[#40617f] text-sm">Vérification en cours...</p>
      </div>

      
      <div v-else-if="state === 'success'" class="space-y-6">
        <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#dcfce7;">
          <CheckCircleIcon class="w-10 h-10 text-[#006d35]" />
        </div>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Email vérifié !</h1>
          <p class="text-[#40617f] mt-2">Votre compte est maintenant actif. Vous pouvez vous connecter.</p>
        </div>
        <RouterLink
          to="/connexion"
          class="inline-flex items-center gap-2 px-8 py-3.5 rounded-xl font-jakarta font-bold text-white text-base transition hover:opacity-90"
          style="background: linear-gradient(135deg, #006d35, #1b8848);"
        >
          <ArrowRightCircleIcon class="w-5 h-5" />
          Se connecter
        </RouterLink>
      </div>

      
      <div v-else-if="state === 'used'" class="space-y-6">
        <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#dcfce7;">
          <CheckCircleIcon class="w-10 h-10 text-[#006d35]" />
        </div>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl">Compte déjà activé</h1>
          <p class="text-[#40617f] mt-2">Votre compte a déjà été validé. Vous pouvez vous connecter.</p>
        </div>
        <RouterLink to="/connexion" class="inline-flex items-center gap-2 px-8 py-3 rounded-xl font-bold text-white transition hover:opacity-90" style="background:#006d35;">
          Se connecter
        </RouterLink>
      </div>

      
      <div v-else class="space-y-6">
        <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#fee2e2;">
          <XCircleIcon class="w-10 h-10 text-red-500" />
        </div>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl">Lien invalide</h1>
          <p class="text-[#40617f] mt-2">{{ errorMsg }}</p>
        </div>
        <RouterLink to="/inscription" class="inline-flex items-center gap-2 px-8 py-3 rounded-xl font-bold text-[#006d35] border-2 border-[#006d35]/30 hover:bg-[#f0fdf4] transition">
          Recommencer l'inscription
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<script setup>import { ref, onMounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { CheckCircleIcon, XCircleIcon, ArrowRightCircleIcon } from '@heroicons/vue/24/outline'

const route    = useRoute()
const state    = ref('loading')
const errorMsg = ref('Ce lien est invalide ou a expiré.')

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    state.value = 'error'
    return
  }
  try {
    const res = await fetch(`/api/v1/auth/verify-email?token=${encodeURIComponent(token)}`)
    if (res.ok) {
      state.value = 'success'
    } else if (res.status === 410) {
      state.value = 'used'
    } else {
      const json = await res.json().catch(() => ({}))
      errorMsg.value = json.message || 'Ce lien est invalide ou a expiré.'
      state.value = 'error'
    }
  } catch {
    errorMsg.value = 'Une erreur réseau est survenue.'
    state.value = 'error'
  }
})
</script>
