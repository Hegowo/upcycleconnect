<template>
  <div class="bg-[#f7f9ff]">

    
    <section class="px-4 sm:px-6 pt-10 sm:pt-16 pb-8 sm:pb-12 max-w-[1280px] mx-auto">
      <div class="flex flex-col md:flex-row items-start gap-8 md:gap-16">
        <div class="flex-1 flex flex-col gap-6">
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl md:text-[64px] md:leading-[72px] leading-tight tracking-[-0.025em]">
            Nos Événements<br />
            <span class="text-[#006d35]">& Formations</span>
          </h1>
          <p class="text-[#40617f] text-xl leading-[32px] max-w-[672px]">
            Découvrez nos ateliers pratiques, formations certifiantes et conférences inspirantes
            autour de l'upcycling et de l'économie circulaire.
          </p>
        </div>
        
        <div class="hidden md:block shrink-0 relative w-72 h-72 lg:w-80 lg:h-80">
          <div class="absolute inset-0 bg-[rgba(0,109,53,0.05)] blur-[32px] rounded-full" />
          <div class="relative w-full h-full rounded-[40px] overflow-hidden shadow-2xl rotate-3 bg-gradient-to-br from-[#dbeafe] to-[#bfdbfe] flex items-center justify-center">
            <CalendarDaysIcon class="w-32 h-32 text-[#1d4ed8]/30" />
          </div>
        </div>
      </div>
    </section>

    
    <section class="px-6 pb-12 max-w-[1280px] mx-auto">
      <div class="grid grid-cols-1 md:grid-cols-12 gap-6">

        
        <div class="md:col-span-8 bg-[#edf4ff] rounded-[32px] p-8 flex flex-col gap-4">
          <p class="text-[#001d32] font-jakarta font-bold text-lg">Trouver un événement</p>
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Rechercher un atelier, une formation... (ex: textile, mobilier)"
              class="w-full bg-white pl-12 pr-4 py-[18px] rounded-xl text-base text-gray-600 outline-none shadow-sm focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400"
            />
            <MagnifyingGlassIcon class="w-[18px] h-[18px] absolute left-4 top-1/2 -translate-y-1/2 text-gray-400" />
          </div>
          <div class="flex gap-3">
            <button
              class="flex items-center gap-2 px-6 py-3 rounded-xl text-white font-bold text-sm transition hover:opacity-90"
              style="background: linear-gradient(134deg, #006d35 0%, #1b8848 100%);"
            >
              <MagnifyingGlassIcon class="w-4 h-4" />
              Rechercher
            </button>
            <button
              @click="searchQuery = ''; activeType = 'Tout'"
              class="px-6 py-3 rounded-xl text-[#40617f] font-semibold text-sm bg-white border border-gray-200 hover:bg-gray-50 transition"
            >
              Réinitialiser
            </button>
          </div>
        </div>

        
        <div class="md:col-span-4 bg-[#edf4ff] rounded-[32px] p-8 flex flex-col gap-4">
          <p class="text-[#001d32] font-jakarta font-bold text-lg">Type d'événement</p>
          <div class="flex flex-col gap-2">
            <button
              v-for="type in eventTypes"
              :key="type"
              @click="activeType = type"
              class="w-full text-left px-5 py-3 rounded-xl text-sm font-semibold transition-all duration-200"
              :class="activeType === type
                ? 'bg-[#006d35] text-white'
                : 'bg-white text-[#40617f] hover:bg-[#006d35]/5'"
            >
              {{ type }}
            </button>
          </div>
        </div>

      </div>
    </section>

    
    <section class="px-6 pb-16 max-w-[1280px] mx-auto">

      
      <div class="flex items-center justify-between mb-8">
        <p class="text-[#40617f] text-sm">
          <span class="font-bold text-[#001d32]">{{ filteredEvents.length }}</span> événement{{ filteredEvents.length > 1 ? 's' : '' }} trouvé{{ filteredEvents.length > 1 ? 's' : '' }}
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">

        <article
          v-for="event in displayedEvents"
          :key="event.id"
          class="bg-white rounded-[32px] overflow-hidden flex flex-col hover:shadow-xl transition-shadow duration-300"
        >
          
          <div class="relative h-56 overflow-hidden flex-shrink-0">
            <div
              class="absolute inset-0 flex items-center justify-center"
              :style="`background: linear-gradient(135deg, ${event.colorFrom}, ${event.colorTo});`"
            >
              <component :is="event.icon" class="w-20 h-20 text-white/25" />
            </div>
            
            <div class="absolute top-4 left-4">
              <span
                class="text-xs font-bold uppercase tracking-[0.6px] px-3 py-1 rounded-full"
                :class="typeBadgeClass(event.type)"
              >
                {{ event.type }}
              </span>
            </div>
            
            <div class="absolute top-4 right-4">
              <span
                class="text-xs font-bold px-3 py-1 rounded-full"
                :class="event.complet
                  ? 'bg-red-100 text-red-700'
                  : 'bg-emerald-100 text-emerald-700'"
              >
                {{ event.complet ? 'Complet' : `${event.places} place${event.places > 1 ? 's' : ''}` }}
              </span>
            </div>
          </div>

          
          <div class="p-8 flex flex-col flex-1">
            
            <div class="flex items-center gap-2 text-[#006d35] text-sm font-semibold mb-3">
              <CalendarDaysIcon class="w-4 h-4" />
              {{ event.date }}
            </div>

            <h3 class="font-jakarta font-extrabold text-[#001d32] text-xl leading-[26px] mb-3">
              {{ event.title }}
            </h3>
            <p class="text-[#40617f] text-sm leading-[22px] flex-1">
              {{ event.description }}
            </p>

            <div class="flex items-center gap-3 mt-4 mb-6">
              <MapPinIcon class="w-4 h-4 text-[#40617f] shrink-0" />
              <span class="text-[#40617f] text-sm">{{ event.location }}</span>
            </div>

            <div class="border-t border-gray-100 pt-5 flex items-center justify-between">
              <div>
                <p class="text-[#40617f] text-xs uppercase tracking-wider">Prix</p>
                <p class="text-[#006d35] font-bold text-xl">{{ event.price }}</p>
              </div>
              <button
                v-if="!event.complet"
                class="bg-[#006d35] text-white font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#1b8848] transition-colors"
              >
                Réserver
              </button>
              <button
                v-else
                class="bg-[#edf4ff] text-[#40617f] font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#d8eaff] transition-colors"
              >
                Liste d'attente
              </button>
            </div>
          </div>
        </article>

      </div>

      
      <div v-if="filteredEvents.length > displayCount" class="flex justify-center mt-12">
        <button
          @click="displayCount += 6"
          class="flex items-center gap-2 px-10 py-4 rounded-xl border-2 border-[#006d35] text-[#006d35] font-bold text-base hover:bg-[#006d35]/5 transition"
        >
          Afficher plus d'événements
          <ChevronDownIcon class="w-5 h-5" />
        </button>
      </div>

      
      <div v-if="filteredEvents.length === 0" class="flex flex-col items-center justify-center py-20 text-center">
        <CalendarDaysIcon class="w-16 h-16 text-[#40617f]/30 mb-4" />
        <p class="text-[#001d32] font-bold text-xl mb-2">Aucun événement trouvé</p>
        <p class="text-[#40617f] text-base">Essayez de modifier vos critères de recherche.</p>
      </div>

    </section>

    
    <section class="px-6 pb-16 max-w-[1280px] mx-auto">
      <div class="rounded-[40px] p-12 flex flex-col md:flex-row items-center justify-between gap-8" style="background: #195687;">
        <div class="flex flex-col gap-4 flex-1">
          <span class="self-start text-xs font-bold uppercase tracking-[0.6px] text-white bg-white/20 px-4 py-1.5 rounded-full">
            Vous êtes formateur ou artisan ?
          </span>
          <h2 class="font-jakarta font-bold text-white text-3xl leading-[38px]">
            Proposez votre propre<br />atelier ou formation
          </h2>
          <p class="text-white/80 text-base leading-[26px] max-w-md">
            Rejoignez notre réseau de professionnels et partagez votre passion pour l'upcycling avec notre communauté.
          </p>
        </div>
        <div class="shrink-0">
          <button class="bg-white text-[#006d35] font-bold text-base px-10 py-4 rounded-xl hover:opacity-90 transition shadow-lg">
            Proposer un événement
          </button>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>import { ref, computed } from 'vue'
import {
  MagnifyingGlassIcon,
  CalendarDaysIcon,
  MapPinIcon,
  ChevronDownIcon,
  ScissorsIcon,
  SparklesIcon,
  LightBulbIcon,
  AcademicCapIcon,
} from '@heroicons/vue/24/outline'

const searchQuery = ref('')
const activeType = ref('Tout')
const displayCount = ref(6)

const eventTypes = ['Tout', 'Atelier', 'Formation', 'Conférence']

const typeBadgeClass = (type) => {
  const map = {
    'Atelier':     'bg-[#d8eaff] text-[#1d4ed8]',
    'Formation':   'bg-[#d1fae5] text-[#065f46]',
    'Conférence':  'bg-[#f3e8ff] text-[#7c3aed]',
  }
  return map[type] || 'bg-gray-100 text-gray-600'
}

const events = [
  {
    id: 1,
    title: 'Atelier Upcycling Textile',
    description: 'Apprenez à transformer vos vieux vêtements en créations uniques. Techniques de teinture naturelle, couture créative et détournement de matières.',
    type: 'Atelier',
    date: '12 avril 2026',
    location: 'Paris 11e',
    price: '45€',
    places: 4,
    complet: false,
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: ScissorsIcon,
  },
  {
    id: 2,
    title: 'Conférence : Économie Circulaire & Design',
    description: 'Réflexion sur l\'impact environnemental du design. Intervenants experts, tables rondes et échanges avec des acteurs du secteur.',
    type: 'Conférence',
    date: '19 avril 2026',
    location: 'Lyon 2e',
    price: 'Gratuit',
    places: 12,
    complet: false,
    colorFrom: '#f3e8ff',
    colorTo: '#e9d5ff',
    icon: SparklesIcon,
  },
  {
    id: 3,
    title: 'Formation Création de Bijoux',
    description: 'Formation certifiante sur 2 jours : refonte de vieux métaux, travail du verre recyclé et création de pièces uniques.',
    type: 'Formation',
    date: '26-27 avril 2026',
    location: 'Marseille',
    price: '180€',
    places: 0,
    complet: true,
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: SparklesIcon,
  },
  {
    id: 4,
    title: 'Atelier Restauration de Mobilier',
    description: 'Ponçage, vernis écologique, réparation structurelle. Donnez un second souffle à vos meubles avec l\'aide de nos ébénistes.',
    type: 'Atelier',
    date: '3 mai 2026',
    location: 'Bordeaux',
    price: '65€',
    places: 2,
    complet: false,
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: LightBulbIcon,
  },
  {
    id: 5,
    title: 'Formation Bee-Wraps & Emballages',
    description: 'Créez vos propres alternatives aux emballages plastiques : bee-wraps, sacs à vrac, boites en carton recyclé.',
    type: 'Formation',
    date: '10 mai 2026',
    location: 'Toulouse',
    price: '55€',
    places: 8,
    complet: false,
    colorFrom: '#d1fae5',
    colorTo: '#bbf7d0',
    icon: AcademicCapIcon,
  },
  {
    id: 6,
    title: 'Conférence : Zéro Déchet au Quotidien',
    description: 'Découvrez des gestes simples et efficaces pour réduire vos déchets et adopter un mode de vie plus durable.',
    type: 'Conférence',
    date: '17 mai 2026',
    location: 'Nantes',
    price: 'Gratuit',
    places: 30,
    complet: false,
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: SparklesIcon,
  },
  {
    id: 7,
    title: 'Atelier Luminaires Industriels',
    description: 'Créez des objets déco et luminaires originaux à partir de pièces industrielles récupérées. Résultats garantis.',
    type: 'Atelier',
    date: '24 mai 2026',
    location: 'Strasbourg',
    price: '75€',
    places: 6,
    complet: false,
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: LightBulbIcon,
  },
  {
    id: 8,
    title: 'Formation Artisan Certifié Upcycling',
    description: 'Formation complète sur 5 jours pour les professionnels souhaitant se certifier dans les métiers de l\'upcycling.',
    type: 'Formation',
    date: '2-6 juin 2026',
    location: 'Paris 3e',
    price: '490€',
    places: 3,
    complet: false,
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: AcademicCapIcon,
  },
]

const filteredEvents = computed(() => {
  return events.filter(e => {
    const matchType = activeType.value === 'Tout' || e.type === activeType.value
    const matchSearch = !searchQuery.value ||
      e.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      e.description.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      e.location.toLowerCase().includes(searchQuery.value.toLowerCase())
    return matchType && matchSearch
  })
})

const displayedEvents = computed(() => filteredEvents.value.slice(0, displayCount.value))
</script>
