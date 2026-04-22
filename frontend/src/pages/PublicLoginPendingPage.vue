<template>
  <div class="min-h-[70vh] flex items-center justify-center py-16 px-4">
    <div class="w-full max-w-md text-center">

      <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto mb-6" style="background:#fef9c3;">
        <EnvelopeIcon class="w-10 h-10" style="color:#854d0e;" />
      </div>

      <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">
        Vérifiez votre boîte mail
      </h1>

      <p class="text-[#40617f] mt-3 leading-relaxed">
        Pour votre sécurité, une connexion depuis
        <strong class="text-[#001d32]">un nouvel emplacement</strong>
        a été détectée.
      </p>

      <div class="mt-6 bg-white rounded-2xl border border-[#edf4ff] p-6 text-left shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)]">
        <div class="flex gap-3 items-start">
          <ShieldCheckIcon class="w-5 h-5 text-[#006d35] shrink-0 mt-0.5" />
          <div class="text-sm text-[#40617f] leading-relaxed">
            Un email contenant un lien de validation vient d'être envoyé
            <span v-if="email" class="text-[#001d32] font-semibold">à {{ email }}</span>.
            <br />
            Cliquez sur le lien pour finaliser votre connexion.
          </div>
        </div>
        <div class="mt-4 flex gap-3 items-start">
          <ClockIcon class="w-5 h-5 text-[#854d0e] shrink-0 mt-0.5" />
          <div class="text-sm text-[#40617f] leading-relaxed">
            Le lien est valable <strong class="text-[#001d32]">15 minutes</strong>.
            Pensez à vérifier vos spams si vous ne le recevez pas.
          </div>
        </div>
      </div>

      <div class="mt-6 flex flex-col gap-3">
        <RouterLink
          to="/connexion"
          class="inline-flex items-center justify-center gap-2 px-8 py-3 rounded-xl font-bold text-[#006d35] border-2 border-[#006d35]/30 hover:bg-[#f0fdf4] transition"
        >
          <ArrowLeftIcon class="w-4 h-4" />
          Retour à la connexion
        </RouterLink>

        <RouterLink
          to="/"
          class="text-xs text-gray-400 hover:text-[#006d35] transition underline-offset-2 hover:underline"
        >
          Revenir à l'accueil
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { EnvelopeIcon, ShieldCheckIcon, ClockIcon, ArrowLeftIcon } from '@heroicons/vue/24/outline'

const route = useRoute()
const email = computed(() => route.query.email || '')
</script>
