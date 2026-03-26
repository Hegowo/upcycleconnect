<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-8 pt-12 pb-16">

      <!-- Success Hero -->
      <div class="flex flex-col items-center text-center mb-12">
        <div class="w-20 h-20 rounded-full bg-[#d1fae5] flex items-center justify-center mb-6 shadow-sm">
          <CheckCircleIcon class="w-10 h-10 text-[#006d35]" />
        </div>
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-5xl tracking-tight mb-3">
          Demande de dépôt acceptée !
        </h1>
        <p class="text-[#40617f] text-lg max-w-xl leading-relaxed">
          Votre objet a été validé. Voici votre bon de dépôt à présenter au point de collecte.
        </p>
      </div>

      <!-- Main Bento Grid -->
      <div class="grid grid-cols-12 gap-6">

        <!-- Depot Pass (col-span-8) -->
        <div class="col-span-8 bg-white rounded-[24px] p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-8">

          <!-- Pass header -->
          <div class="flex items-center justify-between">
            <div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Bon de Dépôt</h2>
              <p class="text-[#40617f] text-sm mt-1">Valable jusqu'au 24 Mai 2024</p>
            </div>
            <span class="bg-[#d1fae5] text-[#006d35] font-bold text-sm px-4 py-2 rounded-full uppercase tracking-wide">
              Validé
            </span>
          </div>

          <!-- QR + Code -->
          <div class="flex gap-8 items-start">
            <!-- QR Code placeholder -->
            <div class="w-40 h-40 rounded-[20px] bg-[#edf4ff] border-2 border-[#d8eaff] flex flex-col items-center justify-center shrink-0 gap-3">
              <QrCodeIcon class="w-20 h-20 text-[#001d32]" />
              <span class="text-[#001d32] font-mono font-bold text-xs tracking-widest">UP-789-QZ</span>
            </div>

            <!-- Details -->
            <div class="flex-1 flex flex-col gap-5">
              <div class="grid grid-cols-2 gap-4">
                <div class="bg-[#f7f9ff] rounded-xl p-4">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">Code de dépôt</p>
                  <p class="font-jakarta font-bold text-[#001d32] text-xl font-mono">UP-789-QZ</p>
                </div>
                <div class="bg-[#f7f9ff] rounded-xl p-4">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">Référence objet</p>
                  <p class="font-jakarta font-bold text-[#001d32] text-xl">#UP-2024-9182</p>
                </div>
                <div class="bg-[#f7f9ff] rounded-xl p-4">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">Validé le</p>
                  <p class="font-jakarta font-bold text-[#001d32] text-base">24 Avril 2024</p>
                </div>
                <div class="bg-[#f7f9ff] rounded-xl p-4">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">Expire le</p>
                  <p class="font-jakarta font-bold text-[#001d32] text-base">24 Mai 2024</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Print button -->
          <button
            class="flex items-center justify-center gap-2 w-full py-4 rounded-xl text-white font-bold text-base transition hover:opacity-90"
            style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
          >
            <PrinterIcon class="w-5 h-5" />
            Imprimer le bon de dépôt
          </button>
        </div>

        <!-- Right column: Object + Location -->
        <div class="col-span-4 flex flex-col gap-6">

          <!-- Object summary -->
          <div class="bg-white rounded-[24px] p-6 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-4">
            <h3 class="font-jakarta font-bold text-[#001d32] text-lg">Votre objet</h3>
            <div class="h-32 rounded-xl bg-gradient-to-br from-[#fef3c7] to-[#fde68a] flex items-center justify-center">
              <HomeModernIcon class="w-16 h-16 text-white/40" />
            </div>
            <div>
              <p class="font-jakarta font-bold text-[#001d32] text-base">Chaise Scandinave</p>
              <p class="text-[#40617f] text-sm mt-1">Mobilier • Bois</p>
            </div>
            <div class="flex items-center gap-2">
              <span class="bg-[#d1fae5] text-[#006d35] text-xs font-bold uppercase tracking-wide px-3 py-1 rounded-full">
                Approuvé
              </span>
            </div>
          </div>

          <!-- Location card -->
          <div class="bg-[#edf4ff] rounded-[24px] p-6 flex flex-col gap-4">
            <h3 class="font-jakarta font-bold text-[#001d32] text-lg">Point de collecte</h3>
            <div class="flex items-start gap-3">
              <div class="w-10 h-10 bg-white rounded-xl flex items-center justify-center shrink-0">
                <MapPinIcon class="w-5 h-5 text-[#006d35]" />
              </div>
              <div>
                <p class="font-jakarta font-bold text-[#001d32] text-base">Eco-Box Bastille</p>
                <p class="text-[#40617f] text-sm mt-0.5">12 Rue de la Roquette</p>
                <p class="text-[#40617f] text-sm">75011 Paris</p>
                <span class="inline-block mt-2 text-[#006d35] text-xs font-bold bg-[#d1fae5] px-2 py-0.5 rounded-full">
                  Ouvert 24h/24
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- Instructions (full width) -->
        <div class="col-span-12 bg-white rounded-[24px] p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)]">
          <h3 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">Comment procéder ?</h3>
          <div class="grid grid-cols-3 gap-6">
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

      <!-- Back link -->
      <div class="flex items-center justify-center mt-10">
        <RouterLink to="/profil" class="flex items-center gap-2 text-[#40617f] hover:text-[#006d35] font-medium text-sm transition">
          <ArrowLeftIcon class="w-4 h-4" />
          Retour au profil
        </RouterLink>
      </div>

    </div>

    <!-- Footer -->
    <footer class="bg-[#edf4ff] border-t border-[rgba(226,232,240,0.5)]">
      <div class="max-w-[1280px] mx-auto px-8 py-12 flex items-center justify-between">
        <div>
          <p class="font-jakarta font-bold text-[#006d35] text-xl">UpcycleConnect</p>
          <p class="text-[#40617f] text-sm mt-1">© 2024 UpcycleConnect. The Curated Ecosystem.</p>
        </div>
        <div class="flex gap-8">
          <span v-for="link in footerLinks" :key="link" class="text-[#40617f] text-sm opacity-80 cursor-pointer hover:opacity-100">
            {{ link }}
          </span>
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import {
  CheckCircleIcon,
  QrCodeIcon,
  PrinterIcon,
  MapPinIcon,
  HomeModernIcon,
  ArrowLeftIcon,
} from '@heroicons/vue/24/outline'

const steps = [
  {
    num: 1,
    title: 'Imprimez votre bon',
    desc: 'Téléchargez et imprimez le bon de dépôt ou enregistrez-le sur votre téléphone.',
  },
  {
    num: 2,
    title: 'Rendez-vous au point de collecte',
    desc: 'Apportez votre objet à l\'Eco-Box Bastille, 12 Rue de la Roquette, 75011 Paris.',
  },
  {
    num: 3,
    title: 'Scannez et déposez',
    desc: 'Présentez le QR code UP-789-QZ au terminal de l\'Eco-Box pour finaliser le dépôt.',
  },
]

const footerLinks = ['Privacy Policy', 'Terms of Service', 'Sustainability Report', 'Contact']
</script>
