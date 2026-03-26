<template>
  <div class="bg-[#f8fafc]">

    <!-- ── Hero ── -->
    <section class="relative overflow-hidden px-8 pt-32 pb-24">
      <!-- Cercles décoratifs blur -->
      <div class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] rounded-full bg-[rgba(28,130,65,0.05)] blur-[32px] pointer-events-none" />
      <div class="absolute right-[-96px] top-[-96px] w-[400px] h-[400px] rounded-full bg-[rgba(25,86,135,0.05)] blur-[32px] pointer-events-none" />

      <div class="max-w-[1280px] mx-auto relative">
        <div class="flex flex-col items-center gap-6 max-w-[896px] mx-auto text-center">
          <!-- Titre -->
          <h1 class="font-jakarta font-bold text-[#0f172a] tracking-[-0.025em] text-5xl md:text-[72px] leading-[1] md:leading-[72px]">
            Donnez une
            <span class="text-[#1c8241]"> seconde vie</span>
            <br />à vos objets
          </h1>

          <!-- Sous-titre -->
          <p class="text-[#334155] text-xl md:text-2xl leading-[32px] max-w-[672px]">
            Connectez-vous avec des artisans passionnés par l'upcycling pour transformer,
            réparer ou réutiliser vos biens de manière durable.
          </p>

          <!-- CTA -->
          <div class="flex flex-col sm:flex-row gap-6 items-center pt-4">
            <RouterLink
              to="/prestations"
              class="inline-flex items-center justify-center px-8 py-[17px] rounded-full bg-[#1c8241] text-white text-base font-medium transition hover:opacity-90 shadow-lg"
            >
              Voir les prestations
            </RouterLink>
            <RouterLink
              to="/inscription"
              class="inline-flex items-center justify-center px-8 py-[17px] rounded-full border-2 border-[#1c8241] text-[#1c8241] text-base font-medium bg-transparent transition hover:bg-[#1c8241]/5 shadow-sm"
            >
              Devenir prestataire
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Catégories ── -->
    <section class="bg-white border-t border-b border-[#e2e8f0] py-24 px-8">
      <div class="max-w-[1280px] mx-auto flex flex-col gap-12">

        <!-- En-tête -->
        <div class="flex flex-col items-center gap-4 text-center">
          <h2 class="font-jakarta font-bold text-[#0f172a] text-[36px] leading-[40px]">
            Nos catégories d'upcycling
          </h2>
          <p class="text-[#334155] text-lg leading-[28px] max-w-[672px]">
            Explorez les différents domaines où nos artisans exercent leur magie pour revaloriser vos objets.
          </p>
        </div>

        <!-- Bento grid -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-6" style="grid-template-rows: 198px 218px;">

          <!-- Grande carte : Mobilier & Décoration (col 1-2, row 1-2) -->
          <RouterLink
            to="/prestations"
            class="col-span-2 row-span-2 relative flex flex-col justify-between p-8 rounded-[32px] overflow-hidden border border-[#bbf7d0] cursor-pointer hover:shadow-lg transition-shadow"
            style="background: linear-gradient(143deg, #f0fdf4 0%, #d1fae5 100%);"
          >
            <div class="absolute top-0 right-0 w-40 h-36 opacity-20">
              <HomeModernIcon class="w-full h-full text-[#006d35]" />
            </div>

            <div class="relative z-10">
              <div class="w-12 h-12 bg-white rounded-full shadow-sm flex items-center justify-center mb-12">
                <HomeModernIcon class="w-5 h-5 text-[#006d35]" />
              </div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-2xl leading-[32px]">
                Mobilier & Décoration
              </h3>
              <p class="text-[#334155] text-base leading-[24px] mt-4 max-w-xs">
                Restauration de meubles anciens, création d'objets déco à partir de matériaux de récupération.
              </p>
            </div>

            <div class="flex items-center gap-1 text-[#1c8241] font-medium text-base relative z-10">
              Découvrir
              <ChevronRightIcon class="w-2.5 h-2.5" />
            </div>
          </RouterLink>

          <!-- Textile & Mode -->
          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col gap-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="w-10 h-10 bg-[#eff6ff] rounded-full flex items-center justify-center">
              <ScissorsIcon class="w-4 h-4 text-[#3b82f6]" />
            </div>
            <div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">Textile & Mode</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">Retouches, transformation de vêtements, création d'accessoires.</p>
            </div>
          </RouterLink>

          <!-- Emballages -->
          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col gap-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="w-10 h-10 bg-[#fff7ed] rounded-full flex items-center justify-center">
              <ArchiveBoxIcon class="w-4 h-4 text-[#f97316]" />
            </div>
            <div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">Emballages</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">Solutions d'emballages réutilisables et upcyclés.</p>
            </div>
          </RouterLink>

          <!-- Électronique -->
          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col gap-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="w-10 h-10 bg-[#faf5ff] rounded-full flex items-center justify-center">
              <ComputerDesktopIcon class="w-4 h-4 text-[#a855f7]" />
            </div>
            <div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">Électronique</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">Réparation et réemploi de petits appareils électroniques.</p>
            </div>
          </RouterLink>

          <!-- Autre -->
          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col justify-between shadow-sm hover:shadow-md transition-shadow">
            <div>
              <div class="w-10 h-10 bg-[#f3f4f6] rounded-full flex items-center justify-center mb-4">
                <EllipsisHorizontalIcon class="w-4 h-4 text-[#6b7280]" />
              </div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">Autre</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">Jouets, livres, articles de sport, et bien plus encore.</p>
            </div>
            <div class="flex items-center gap-1 text-[#4b5563] font-medium text-sm">
              Voir tout <ChevronRightIcon class="w-2 h-2" />
            </div>
          </RouterLink>

        </div>
      </div>
    </section>

    <!-- ── CTA Prestataire ── -->
    <section class="py-20 px-8">
      <div class="max-w-[1280px] mx-auto">
        <div
          class="relative overflow-hidden rounded-[32px] px-16 py-20 flex flex-col md:flex-row items-center justify-between gap-10 shadow-xl"
          style="background: #195687;"
        >
          <div class="flex flex-col gap-4 max-w-lg">
            <h2 class="font-jakarta font-bold text-white text-[36px] leading-[40px]">
              Vous êtes artisan ou prestataire ?
            </h2>
            <p class="text-[#dbeafe] text-lg leading-[28px]">
              Rejoignez notre communauté croissante. Proposez vos services, gagnez en visibilité et participez à l'économie circulaire.
            </p>
          </div>
          <RouterLink
            to="/inscription"
            class="shrink-0 inline-flex items-center justify-center bg-white font-jakarta font-bold text-[#195687] text-base px-8 py-[17px] rounded-full shadow-lg transition hover:opacity-90"
          >
            Créer mon profil prestataire
          </RouterLink>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import {
  HomeModernIcon,
  ScissorsIcon,
  ArchiveBoxIcon,
  ComputerDesktopIcon,
  EllipsisHorizontalIcon,
  ChevronRightIcon,
} from '@heroicons/vue/24/outline'
</script>
