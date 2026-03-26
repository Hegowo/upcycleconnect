<template>
  <div class="bg-[#f7f9ff] min-h-screen flex">

    <!-- Sidebar -->
    <aside class="w-64 shrink-0 bg-[#f8fafc] border-r border-[#e2e8f0] flex flex-col py-6 px-4 sticky top-[65px] h-[calc(100vh-65px)] overflow-y-auto">
      <div class="px-4 py-6">
        <p class="text-[#15803d] font-semibold text-sm">Navigation Forum</p>
        <p class="text-[#64748b] text-xs mt-0.5">Écosystème Curé</p>
      </div>

      <nav class="flex flex-col gap-1">
        <RouterLink
          v-for="item in navItems"
          :key="item.label"
          :to="item.to"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition"
          :class="item.active
            ? 'bg-[#f0fdf4] text-[#15803d]'
            : 'text-[#64748b] hover:bg-[#f0fdf4] hover:text-[#15803d]'"
        >
          <component :is="item.icon" class="w-5 h-5 shrink-0" />
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="mt-8 px-4">
        <button
          class="w-full py-3 rounded-xl text-white font-semibold text-base text-center transition hover:opacity-90"
          style="background: linear-gradient(165deg, #006d35 0%, #1b8848 100%);"
        >
          Nouveau Sujet
        </button>
      </div>

      <!-- Bottom links -->
      <div class="mt-auto flex flex-col gap-1 pb-4">
        <button class="flex items-center gap-3 px-4 py-3 rounded-xl text-[#64748b] text-sm font-semibold hover:bg-[#f0fdf4] hover:text-[#15803d] transition">
          <Cog6ToothIcon class="w-5 h-5 shrink-0" />
          Paramètres
        </button>
        <button class="flex items-center gap-3 px-4 py-3 rounded-xl text-[#64748b] text-sm font-semibold hover:bg-[#f0fdf4] hover:text-[#15803d] transition">
          <QuestionMarkCircleIcon class="w-5 h-5 shrink-0" />
          Aide
        </button>
      </div>
    </aside>

    <!-- Main content -->
    <main class="flex-1 min-w-0 p-8">
      <div class="max-w-[1152px] flex flex-col gap-12">

        <!-- Header -->
        <div class="flex flex-col gap-2">
          <p class="text-[#006d35] text-xs font-semibold uppercase tracking-[1.2px]">Forum Communautaire</p>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-5xl tracking-tight">Catégories de Discussion</h1>
          <p class="text-[#40617f] text-lg mt-2 max-w-2xl leading-relaxed">
            Explorez nos différents espaces thématiques et rejoignez une communauté engagée pour transformer le futur, un objet à la fois.
          </p>
        </div>

        <!-- Bento Grid -->
        <div class="grid grid-cols-12 gap-6">

          <!-- Matériaux & Ressources (col-7) -->
          <div class="col-span-7 bg-white rounded-[24px] p-10 relative flex flex-col justify-between min-h-[378px]">
            <div class="flex items-start justify-between">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
                <CubeIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <span class="bg-[#edf4ff] text-[#40617f] text-sm font-semibold px-3 py-1 rounded-full">1.2k Sujets</span>
            </div>
            <div class="mt-12">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Matériaux & Ressources</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3 max-w-lg">
                Découvrez comment traiter et réutiliser le bois, le métal, le plastique et les textiles pour vos prochaines créations.
              </p>
              <div class="flex flex-wrap gap-2 mt-4">
                <span v-for="tag in ['Bois', 'Métal', 'Plastique', 'Textile']" :key="tag"
                  class="bg-[#e3efff] text-[#001d32] text-xs font-semibold px-3 py-1 rounded-lg">
                  {{ tag }}
                </span>
              </div>
            </div>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
              Explorer <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Techniques d'Upcycling (col-5) -->
          <div class="col-span-5 bg-white rounded-[24px] p-10 relative flex flex-col justify-between min-h-[378px]">
            <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
              <WrenchScrewdriverIcon class="w-7 h-7 text-[#40617f]" />
            </div>
            <div class="mt-12">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Techniques d'Upcycling</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3">
                Apprenez la transformation, la peinture et l'assemblage avec nos experts.
              </p>
              <p class="text-[#40617f] text-sm font-medium mt-4">850 discussions actives</p>
            </div>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
              Explorer <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Projets & Inspirations (col-4) -->
          <div class="col-span-4 bg-white rounded-[24px] p-10 relative flex flex-col justify-between min-h-[374px] overflow-hidden">
            <!-- Decorative corner element -->
            <div class="absolute -right-4 -top-4 w-24 h-32 opacity-5 pointer-events-none">
              <div class="w-full h-full bg-[#006d35] rounded-xl" />
            </div>
            <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
              <LightBulbIcon class="w-7 h-7 text-[#40617f]" />
            </div>
            <div class="mt-12">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Projets &<br>Inspirations</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3">
                Partagez vos réussites et inspirez-vous des créations de la communauté.
              </p>
            </div>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
              Explorer <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Éco-conception & Impact (col-8) -->
          <div class="col-span-8 bg-white rounded-[24px] p-10 flex gap-8 items-center min-h-[374px]">
            <div class="flex-1">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
                <GlobeAltIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl mt-12">Éco-conception &<br>Impact</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3">
                Calculez votre Upcycling Score et apprenez à réduire votre empreinte carbone par le design.
              </p>
              <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
                Explorer <ChevronRightIcon class="w-3.5 h-3.5" />
              </button>
            </div>

            <!-- A+ Score ring -->
            <div class="bg-[#edf4ff] rounded-xl w-64 h-64 flex items-center justify-center shrink-0 relative">
              <!-- SVG ring -->
              <svg class="absolute inset-0 w-full h-full -rotate-90" viewBox="0 0 256 256">
                <circle cx="128" cy="128" r="100" fill="none" stroke="#d8eaff" stroke-width="12" />
                <circle cx="128" cy="128" r="100" fill="none" stroke="#006d35" stroke-width="12"
                  stroke-dasharray="628" stroke-dashoffset="38" stroke-linecap="round" />
              </svg>
              <div class="relative text-center">
                <p class="font-semibold text-[#006d35] text-4xl">A+</p>
                <p class="text-[#40617f] text-xs uppercase tracking-[1.2px] font-semibold mt-1">Score Moyen</p>
              </div>
            </div>
          </div>

          <!-- Entraide & Questions (col-6) -->
          <div class="col-span-6 bg-white rounded-[24px] p-10 flex flex-col gap-6 self-start">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center shrink-0">
                <ChatBubbleLeftEllipsisIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <div>
                <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Entraide & Questions</h2>
                <p class="text-[#40617f] text-sm mt-0.5">Un problème ? La communauté a la solution.</p>
              </div>
            </div>
            <p class="text-[#3f4a3f] text-base leading-relaxed">
              Besoin d'un conseil technique ou d'un coup de main sur un projet spécifique ? Postez vos questions ici.
            </p>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base pt-2">
              Explorer <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Vente & Business (col-6) -->
          <div class="col-span-6 bg-white rounded-[24px] p-10 flex flex-col gap-6 self-start">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center shrink-0">
                <ShoppingBagIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <div>
                <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Vente & Business</h2>
                <p class="text-[#40617f] text-sm mt-0.5">Pour les professionnels et artisans.</p>
              </div>
            </div>
            <p class="text-[#3f4a3f] text-base leading-relaxed">
              Monétisez vos créations, trouvez des partenaires locaux et développez votre activité artisanale circulaire.
            </p>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base pt-2">
              Explorer <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <!-- Dark CTA (full width) -->
          <div class="col-span-12 bg-[#0f172a] rounded-[32px] px-12 py-16 flex items-center justify-between relative overflow-hidden">
            <!-- Decorative glow -->
            <div class="absolute -right-20 -top-16 w-64 h-64 rounded-full bg-[rgba(0,109,53,0.2)] blur-[32px] pointer-events-none" />
            <div>
              <h3 class="font-jakarta font-bold text-white text-3xl">Vous ne trouvez pas votre bonheur ?</h3>
              <p class="text-[#94a3b8] text-base mt-2">Proposez une nouvelle catégorie ou créez un sujet général.</p>
            </div>
            <button
              class="bg-[#006d35] text-white font-semibold text-base px-8 py-4 rounded-xl hover:bg-[#1b8848] transition shrink-0"
            >
              Proposer un thème
            </button>
          </div>

        </div>
      </div>
    </main>

  </div>
</template>

<script setup>
import { RouterLink } from 'vue-router'
import {
  ChatBubbleLeftEllipsisIcon,
  LightBulbIcon,
  WrenchScrewdriverIcon,
  CubeIcon,
  GlobeAltIcon,
  ShoppingBagIcon,
  ChevronRightIcon,
  Cog6ToothIcon,
  QuestionMarkCircleIcon,
} from '@heroicons/vue/24/outline'

const navItems = [
  { label: 'Forum',            icon: ChatBubbleLeftEllipsisIcon, to: '/communaute/categories-forum', active: true },
  { label: 'Projets Partagés', icon: LightBulbIcon,              to: '/communaute',                 active: false },
  { label: 'Actualités',       icon: GlobeAltIcon,               to: '/communaute',                 active: false },
  { label: 'Conseils',         icon: WrenchScrewdriverIcon,      to: '/communaute',                 active: false },
]

const footerLinks = ['Confidentialité', 'Conditions', 'Contact']
</script>
