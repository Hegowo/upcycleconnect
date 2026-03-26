<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-6 py-10">
      <div class="grid grid-cols-12 gap-6">

        <!-- ── Profile Header (col-span-8) ── -->
        <div class="col-span-8 bg-white rounded-[24px] p-8 flex gap-8 overflow-hidden relative">
          <!-- Avatar -->
          <div class="relative shrink-0">
            <div class="w-32 h-32 rounded-[24px] overflow-hidden shadow-sm bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] flex items-center justify-center">
              <UserCircleIcon class="w-20 h-20 text-[#006d35]/40" />
            </div>
            <div class="absolute bottom-0 right-0 w-6 h-6 bg-[#006d35] rounded-full flex items-center justify-center">
              <CheckIcon class="w-3.5 h-3.5 text-white" />
            </div>
          </div>

          <!-- Info -->
          <div class="flex-1 flex flex-col gap-2">
            <div class="flex items-center gap-3">
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Clara Durand</h1>
              <span class="bg-[rgba(27,136,72,0.1)] text-[#006d35] text-xs font-bold uppercase tracking-wide px-3 py-1 rounded-full">
                Membre Éco-Actif
              </span>
            </div>
            <div class="flex items-center gap-2 text-[#40617f] text-sm">
              <MapPinIcon class="w-3.5 h-3.5" />
              <span>Lyon, France</span>
            </div>
            <p class="text-[#3f4a3f] text-base leading-relaxed mt-2 max-w-md">
              Passionnée par l'upcycling textile et la restauration de mobilier vintage. Je crois en une économie circulaire et créative.
            </p>
          </div>

          <!-- Decorative circle -->
          <div class="absolute -right-10 -top-10 w-36 h-36 rounded-full bg-[rgba(0,109,53,0.05)]" />
        </div>

        <!-- ── Upcycling Score (col-span-4) ── -->
        <div class="col-span-4 bg-[#edf4ff] rounded-[24px] p-8 flex flex-col justify-between">
          <div class="flex items-center justify-between mb-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-lg">Upcycling Score</h2>
            <ArrowTrendingUpIcon class="w-5 h-5 text-[#006d35]" />
          </div>

          <div class="mb-4">
            <div class="flex items-baseline gap-1">
              <span class="font-jakarta font-extrabold text-[#006d35] text-5xl">75</span>
              <span class="text-[#40617f] text-lg">/100</span>
            </div>
            <div class="mt-3 h-3 bg-[#cee5ff] rounded-full overflow-hidden">
              <div class="h-full rounded-full bg-gradient-to-r from-[#006d35] to-[#1b8848]" style="width: 75%;" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div class="bg-white rounded-xl p-3">
              <p class="text-[#40617f] text-xs uppercase tracking-wide">CO2 Économisé</p>
              <p class="font-semibold text-[#001d32] text-lg mt-0.5">12.4 kg</p>
            </div>
            <div class="bg-white rounded-xl p-3">
              <p class="text-[#40617f] text-xs uppercase tracking-wide">Objets Sauvés</p>
              <p class="font-semibold text-[#001d32] text-lg mt-0.5">24</p>
            </div>
          </div>
        </div>

        <!-- ── Mes Annonces (col-span-12) ── -->
        <div class="col-span-12 bg-white rounded-[24px] p-8">
          <div class="flex items-center justify-between mb-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">Mes Annonces</h2>
            <button class="flex items-center gap-1 text-[#006d35] font-semibold text-sm hover:underline">
              Tout voir
              <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <div class="grid grid-cols-3 gap-6">
            <div
              v-for="annonce in annonces"
              :key="annonce.id"
              class="flex flex-col"
            >
              <div class="relative rounded-xl overflow-hidden h-52 bg-gradient-to-br"
                :style="`background: linear-gradient(135deg, ${annonce.colorFrom}, ${annonce.colorTo});`"
              >
                <div class="absolute inset-0 flex items-center justify-center">
                  <component :is="annonce.icon" class="w-20 h-20 text-white/20" />
                </div>
                <div class="absolute top-3 left-3">
                  <span
                    class="text-xs font-bold uppercase tracking-tight px-2 py-1 rounded-lg backdrop-blur-sm"
                    :class="statusClass(annonce.status)"
                  >
                    {{ annonce.status }}
                  </span>
                </div>
              </div>
              <h3 class="font-semibold text-[#001d32] text-base mt-3">{{ annonce.title }}</h3>
              <p class="text-[#40617f] text-xs mt-0.5">{{ annonce.date }}</p>
            </div>
          </div>
        </div>

        <!-- ── Mes Réservations (col-span-7) ── -->
        <div class="col-span-7 bg-[#edf4ff] rounded-[24px] p-8">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">Mes Réservations</h2>

          <div class="flex flex-col gap-4">
            <div
              v-for="resa in reservations"
              :key="resa.id"
              class="bg-white rounded-xl p-4 flex items-center gap-4"
              :class="resa.past ? 'opacity-70' : ''"
            >
              <div
                class="w-14 h-14 rounded-lg flex items-center justify-center shrink-0"
                :class="resa.past ? 'bg-[rgba(185,219,254,0.2)]' : 'bg-[rgba(27,136,72,0.1)]'"
              >
                <CalendarDaysIcon class="w-6 h-6" :class="resa.past ? 'text-[#40617f]' : 'text-[#006d35]'" />
              </div>
              <div class="flex-1">
                <p class="font-semibold text-[#001d32] text-base">{{ resa.title }}</p>
                <p class="text-[#40617f] text-sm">{{ resa.date }}</p>
              </div>
              <component
                :is="resa.past ? CheckIcon : ChevronRightIcon"
                class="w-4 h-4 text-[#94a3b8] shrink-0"
              />
            </div>
          </div>
        </div>

        <!-- ── Mes Badges (col-span-5) ── -->
        <div class="col-span-5 bg-white rounded-[24px] p-8 flex flex-col justify-between">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">Mes Badges</h2>

          <div class="flex gap-4 flex-wrap">
            <div
              v-for="badge in badges"
              :key="badge.id"
              class="relative group"
            >
              <div
                class="w-16 h-16 rounded-full flex items-center justify-center border-2"
                :style="`background: ${badge.bg}; border-color: ${badge.border};`"
              >
                <component :is="badge.icon" class="w-6 h-6" :style="`color: ${badge.iconColor};`" />
              </div>
              <div class="absolute -bottom-6 left-1/2 -translate-x-1/2 opacity-0 group-hover:opacity-100 transition-opacity bg-[#001d32] text-white text-[10px] px-2 py-1 rounded whitespace-nowrap z-10">
                {{ badge.label }}
              </div>
            </div>
            <!-- Locked badge -->
            <div class="w-16 h-16 rounded-full flex items-center justify-center border-2 border-dashed border-[#becabc] opacity-40">
              <LockClosedIcon class="w-5 h-5 text-[#64748b]" />
            </div>
          </div>

          <div class="border-t border-[#edf4ff] pt-6 mt-6 flex items-center justify-between">
            <span class="text-[#40617f] text-xs font-semibold uppercase tracking-wide">3 / 12 Badges</span>
            <button class="bg-[#edf4ff] p-2 rounded-lg hover:bg-[#d8eaff] transition">
              <PlusIcon class="w-4 h-4 text-[#40617f]" />
            </button>
          </div>
        </div>

      </div>
    </div>

    <!-- Footer -->
    <footer class="bg-[#edf4ff] border-t border-[rgba(226,232,240,0.5)] mt-10">
      <div class="max-w-[1280px] mx-auto px-10 py-12 flex items-center justify-between">
        <div>
          <p class="font-jakarta font-bold text-[#40617f] text-lg">UpcycleConnect</p>
          <p class="text-[#40617f] text-xs uppercase tracking-wide mt-1">© 2024 UpcycleConnect. L'écosystème durable.</p>
        </div>
        <div class="flex gap-8">
          <span v-for="link in footerLinks" :key="link" class="text-[#40617f] text-xs uppercase tracking-wide opacity-80 cursor-pointer hover:opacity-100">
            {{ link }}
          </span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import {
  UserCircleIcon,
  MapPinIcon,
  ChevronRightIcon,
  CheckIcon,
  CalendarDaysIcon,
  ArrowTrendingUpIcon,
  LockClosedIcon,
  PlusIcon,
  HomeModernIcon,
  ScissorsIcon,
  SparklesIcon,
  BoltIcon,
  GlobeAltIcon,
  HeartIcon,
} from '@heroicons/vue/24/outline'

const annonces = [
  {
    id: 1,
    title: 'Chaise vintage scandinave',
    date: 'Donné le 12 Mars 2024',
    status: 'Donné',
    colorFrom: '#d1fae5',
    colorTo: '#bbf7d0',
    icon: HomeModernIcon,
  },
  {
    id: 2,
    title: 'Lot de 12 bocaux en verre',
    date: 'Vendu le 05 Mars 2024',
    status: 'Vendu',
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: SparklesIcon,
  },
  {
    id: 3,
    title: 'Miroir biseauté doré',
    date: 'Ajouté il y a 2 jours',
    status: 'Disponible',
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: ScissorsIcon,
  },
]

const statusClass = (status) => {
  if (status === 'Donné' || status === 'Disponible') return 'bg-white/90 text-[#006d35]'
  if (status === 'Vendu') return 'bg-white/90 text-[#40617f]'
  return 'bg-white/90 text-gray-600'
}

const reservations = [
  { id: 1, title: 'Atelier Couture Denim', date: 'Samedi 24 Avril • 14h30', past: false },
  { id: 2, title: 'Conférence Zero Déchet', date: 'Passé le 15 Mars', past: true },
]

const badges = [
  { id: 1, icon: GlobeAltIcon, label: 'Reboiseur', bg: '#92f8ab', border: 'rgba(0,109,53,0.2)', iconColor: '#0f522b' },
  { id: 2, icon: HeartIcon, label: 'Donneur Or', bg: '#b9dbfe', border: 'rgba(64,97,127,0.2)', iconColor: '#1d4ed8' },
  { id: 3, icon: BoltIcon, label: 'Activateur', bg: '#b0f2bd', border: 'rgba(44,106,64,0.2)', iconColor: '#006d35' },
]

const footerLinks = ['À propos', 'Confidentialité', 'Aide', 'Partenaires']
</script>
