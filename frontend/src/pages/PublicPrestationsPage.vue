<template>
  <div class="bg-[#f7f9ff]">

    
    <section class="px-4 sm:px-6 pt-10 sm:pt-16 pb-8 sm:pb-12 max-w-[1280px] mx-auto">
      <div class="flex flex-col md:flex-row items-start gap-8 md:gap-16">
        <div class="flex-1 flex flex-col gap-6">
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl md:text-[64px] md:leading-[72px] leading-tight tracking-[-0.025em]">
            Nos Prestations<br />
            <span class="text-[#006d35]">d'Upcycling</span>
          </h1>
          <p class="text-[#40617f] text-xl leading-[32px] max-w-[672px]">
            Transformez vos objets du quotidien en pièces uniques avec l'aide de nos artisans locaux.
            Donnez une seconde vie à votre intérieur tout en préservant la planète.
          </p>
        </div>
        
        <div class="hidden md:block shrink-0 relative w-72 h-72 lg:w-80 lg:h-80">
          <div class="absolute inset-0 bg-[rgba(0,109,53,0.05)] blur-[32px] rounded-full" />
          <div class="relative w-full h-full rounded-[40px] overflow-hidden shadow-2xl rotate-3 bg-gradient-to-br from-[#d1fae5] to-[#bbf7d0] flex items-center justify-center">
            <HomeModernIcon class="w-32 h-32 text-[#006d35]/30" />
          </div>
        </div>
      </div>
    </section>

    
    <section class="px-4 sm:px-6 pb-8 sm:pb-12 max-w-[1280px] mx-auto">
      <div class="bg-[#edf4ff] rounded-[32px] p-6 sm:p-12 flex flex-col gap-6 sm:gap-8">

        
        <div class="flex gap-4">
          <div class="flex-1 relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Rechercher une prestation (ex: canapé, robe, bijoux...)"
              class="w-full bg-white pl-12 pr-4 py-[18px] rounded-xl text-base text-gray-600 outline-none shadow-sm focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400"
            />
            <MagnifyingGlassIcon class="w-[18px] h-[18px] absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
          </div>
          <button
            class="flex items-center gap-2 px-8 py-4 rounded-xl text-white font-bold text-base transition hover:opacity-90"
            style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
          >
            <AdjustmentsHorizontalIcon class="w-[18px] h-[18px]" />
            Filtrer
          </button>
        </div>

        
        <div class="flex flex-wrap gap-3">
          <button
            v-for="cat in categories"
            :key="cat"
            @click="activeCategory = cat"
            class="px-6 py-2 rounded-full text-base font-semibold transition-all duration-200"
            :class="activeCategory === cat
              ? 'bg-[#006d35] text-white'
              : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
          >
            {{ cat }}
          </button>
        </div>

      </div>
    </section>

    
    <section class="px-4 sm:px-6 pb-12 sm:pb-16 max-w-[1280px] mx-auto">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">

        <article
          v-for="prestation in filteredPrestations"
          :key="prestation.id"
          class="bg-white rounded-[32px] overflow-hidden flex flex-col hover:shadow-xl transition-shadow duration-300"
        >
          
          <div class="relative h-64 overflow-hidden flex-shrink-0">
            <div class="absolute inset-0 flex items-center justify-center" :style="`background: linear-gradient(135deg, ${prestation.colorFrom}, ${prestation.colorTo});`">
              <component :is="prestation.icon" class="w-24 h-24 text-white/30" />
            </div>
            
            <div class="absolute top-4 left-4">
              <div class="flex items-center gap-1.5 bg-white/90 backdrop-blur-sm px-3 py-1 rounded-full text-xs font-bold text-[#006d35]">
                <span>🌿</span> Score {{ prestation.score }}
              </div>
            </div>
          </div>

          
          <div class="p-8 flex flex-col flex-1">
            <h3 class="font-jakarta font-extrabold text-[#001d32] text-2xl leading-[30px] mb-4">
              {{ prestation.title }}
            </h3>
            <p class="text-[#40617f] text-sm leading-[22px] flex-1">
              {{ prestation.description }}
            </p>

            <div class="border-t border-gray-100/50 mt-6 pt-5 flex items-center justify-between">
              <div>
                <p class="text-[#40617f] text-xs uppercase tracking-wider">À partir de</p>
                <p class="text-[#006d35] font-bold text-xl">{{ prestation.price }}</p>
              </div>
              <button class="bg-[#cee5ff] text-[#001d32] font-bold text-base px-5 py-3 rounded-xl hover:bg-[#b8d8ff] transition-colors">
                Réserver
              </button>
            </div>
          </div>
        </article>

        
        <article class="rounded-[32px] overflow-hidden flex flex-col p-8 justify-between" style="background: #1b8848;">
          <div class="flex flex-col gap-4">
            <span class="self-start text-xs font-bold uppercase tracking-[0.6px] text-white bg-white/20 backdrop-blur-sm px-4 py-1.5 rounded-full">
              Sur-mesure
            </span>
            <h3 class="font-jakarta font-bold text-white text-3xl leading-[37px]">
              Vous avez un projet spécifique ?
            </h3>
            <p class="text-white/80 text-base leading-[26px]">
              Nos artisans partenaires étudient vos demandes personnalisées pour des projets d'upcycling hors normes.
            </p>
          </div>
          <button
            class="mt-8 w-full bg-white text-[#006d35] font-bold text-base py-4 rounded-xl hover:opacity-90 transition shadow-lg"
          >
            Faire une demande spéciale
          </button>
        </article>

      </div>
    </section>

    
    <section class="bg-[#edf4ff] px-4 sm:px-8 py-12 sm:py-20">
      <div class="max-w-[896px] mx-auto flex flex-col items-center gap-8 text-center px-6">
        <span class="flex items-center gap-2 bg-[rgba(0,109,53,0.1)] text-[#006d35] font-bold text-sm px-4 py-2 rounded-full">
          🌿 LE SAVIEZ-VOUS ?
        </span>
        <h2 class="font-jakarta font-extrabold text-[#001d32] text-[36px] leading-[40px]">
          L'upcycling réduit de 70% l'empreinte carbone par rapport au recyclage traditionnel.
        </h2>
        <p class="text-[#40617f] text-base leading-[24px]">
          Rejoignez 12 000 passionnés et recevez nos meilleurs conseils chaque mois.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 w-full max-w-[512px]">
          <input
            v-model="newsletterEmail"
            type="email"
            placeholder="votre@email.com"
            class="flex-1 bg-white px-6 py-[18px] rounded-xl text-base outline-none shadow-sm focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400"
          />
          <button
            class="px-8 py-4 rounded-xl text-white font-bold text-base transition hover:opacity-90"
            style="background: linear-gradient(135deg, #006d35, #1b8848);"
          >
            S'abonner
          </button>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>import { ref, computed } from 'vue'
import {
  MagnifyingGlassIcon,
  AdjustmentsHorizontalIcon,
  HomeModernIcon,
  ScissorsIcon,
  SparklesIcon,
  LightBulbIcon,
  ArchiveBoxIcon,
} from '@heroicons/vue/24/outline'

const searchQuery = ref('')
const newsletterEmail = ref('')
const activeCategory = ref('Tous')

const categories = ['Tous', 'Mobilier', 'Textile', 'Bijoux', 'Décoration', 'Emballages']

const prestations = [
  {
    id: 1,
    title: 'Restauration de Mobilier Ancien',
    description: 'Donnez un second souffle à vos meubles de famille. Ponçage, vernis écologique et réparation structurelle par nos ébénistes.',
    price: '85€',
    score: 94,
    category: 'Mobilier',
    colorFrom: '#d1fae5',
    colorTo: '#bbf7d0',
    icon: HomeModernIcon,
  },
  {
    id: 2,
    title: 'Upcycling Textile sur-mesure',
    description: 'Transformez vos vieux jeans ou rideaux en accessoires de mode uniques ou en nouveaux vêtements design.',
    price: '45€',
    score: 88,
    category: 'Textile',
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: ScissorsIcon,
  },
  {
    id: 3,
    title: 'Création de Bijoux Éco-responsables',
    description: 'Refonte de vieux métaux ou création à partir de matériaux recyclés (verre, plastique océanique, bois).',
    price: '60€',
    score: 92,
    category: 'Bijoux',
    colorFrom: '#f3e8ff',
    colorTo: '#e9d5ff',
    icon: SparklesIcon,
  },
  {
    id: 4,
    title: 'Luminaires & Décoration Upcyclée',
    description: "Création d'objets déco et luminaires à partir de pièces industrielles ou d'objets détournés.",
    price: '75€',
    score: 96,
    category: 'Décoration',
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: LightBulbIcon,
  },
  {
    id: 5,
    title: 'Emballages Durables & Ateliers',
    description: 'Ateliers de création de bee-wraps et sacs à vrac personnalisés à partir de chutes de tissus.',
    price: '25€',
    score: 99,
    category: 'Emballages',
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: ArchiveBoxIcon,
  },
]

const filteredPrestations = computed(() => {
  return prestations.filter(p => {
    const matchCat = activeCategory.value === 'Tous' || p.category === activeCategory.value
    const matchSearch = !searchQuery.value ||
      p.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      p.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchCat && matchSearch
  })
})
</script>
