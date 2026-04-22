<template>
  <div class="bg-[#f7f9ff] min-h-screen flex items-center justify-center px-4">
    <div class="bg-white rounded-[32px] p-10 max-w-[520px] w-full text-center shadow-sm">
      <div class="w-20 h-20 mx-auto rounded-full bg-orange-50 flex items-center justify-center mb-6">
        <XCircleIcon class="w-12 h-12 text-orange-500" />
      </div>
      <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl mb-3">
        Paiement annulé
      </h1>
      <p class="text-[#40617f] text-base leading-7 mb-6">
        Votre paiement a été annulé. Aucun montant ne vous a été facturé. Vous pouvez réessayer
        à tout moment depuis la page de la prestation.
      </p>
      <router-link
        to="/prestations"
        class="px-6 py-3 rounded-xl text-white font-bold transition hover:opacity-90 inline-block"
        style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
      >
        Retour aux prestations
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { XCircleIcon } from '@heroicons/vue/24/outline'
</script>
