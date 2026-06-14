<template>
  <Teleport to="body">
    <div
      v-if="show"
      class="fixed inset-0 z-[100] flex items-center justify-center p-4"
      @keydown.capture.prevent
      tabindex="-1"
    >
      <div class="absolute inset-0 bg-[#001d32]/85 backdrop-blur-sm" />

      <div class="relative bg-white rounded-[28px] shadow-2xl w-full max-w-lg overflow-hidden">
        <div class="px-6 pt-6 pb-4" style="background: linear-gradient(135deg, #006d35, #1b8848);">
          <div class="flex items-center justify-between mb-3">
            <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-white/20 text-white text-[10px] font-bold uppercase tracking-wider">
              <SparklesIcon class="w-3.5 h-3.5" />
              Bienvenue sur UpcycleConnect
            </span>
            <span class="text-white/70 text-xs font-semibold">{{ step + 1 }} / {{ steps.length }}</span>
          </div>
          <h2 class="font-jakarta font-extrabold text-white text-2xl leading-tight">{{ current.title }}</h2>
        </div>

        <div class="flex items-center gap-1.5 px-6 pt-4">
          <div
            v-for="(_, i) in steps"
            :key="i"
            class="h-1 flex-1 rounded-full transition-all"
            :class="i <= step ? 'bg-[#006d35]' : 'bg-[#e5e7eb]'"
          />
        </div>

        <div class="px-6 py-6">
          <div class="flex items-start gap-4 mb-4">
            <div class="w-12 h-12 rounded-2xl flex items-center justify-center shrink-0" style="background: linear-gradient(135deg, #006d35, #1b8848);">
              <component :is="current.icon" class="w-6 h-6 text-white" />
            </div>
            <p class="text-[#001d32] text-sm leading-relaxed pt-1">{{ current.desc }}</p>
          </div>

          <div v-if="current.tip" class="bg-[#edf4ff] rounded-xl px-4 py-3 flex items-start gap-2.5">
            <LightBulbIcon class="w-5 h-5 text-[#006d35] shrink-0 mt-0.5" />
            <p class="text-xs text-[#001d32] leading-relaxed"><span class="font-semibold">Astuce :</span> {{ current.tip }}</p>
          </div>
        </div>

        <div class="px-6 py-4 bg-[#f8fafc] border-t border-gray-100 flex items-center justify-between gap-3">
          <button
            @click="prev"
            :disabled="step === 0 || saving"
            class="px-4 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] disabled:opacity-30 disabled:cursor-not-allowed hover:bg-white transition"
          >
            ← Précédent
          </button>
          <button
            @click="next"
            :disabled="saving"
            class="px-6 py-2.5 rounded-xl text-white font-jakarta font-bold text-sm transition hover:opacity-90 flex items-center gap-2 disabled:opacity-60"
            style="background: linear-gradient(135deg, #006d35, #1b8848);"
          >
            <div v-if="saving" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            <span v-else>{{ isLast ? "C'est parti !" : 'Suivant' }}</span>
            <ArrowRightIcon v-if="!saving && !isLast" class="w-4 h-4" />
            <CheckCircleIcon v-else-if="!saving && isLast" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import {
  SparklesIcon,
  LightBulbIcon,
  ArrowRightIcon,
  CheckCircleIcon,
  ArchiveBoxIcon,
  ArchiveBoxArrowDownIcon,
  ShoppingBagIcon,
  GlobeAltIcon,
  CalendarDaysIcon,
  BellAlertIcon,
} from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'

const props = defineProps({
  show: { type: Boolean, default: false },
})
const emit = defineEmits(['done'])

const userAuth = useUserAuthStore()
const step = ref(0)
const saving = ref(false)

const steps = [
  {
    icon: SparklesIcon,
    title: 'Bienvenue sur UpcycleConnect !',
    desc: "La plateforme de l'upcycling intelligent : donnez une seconde vie aux objets, trouvez des matériaux, participez à des ateliers et suivez votre impact écologique.",
    tip: 'Ce tour guidé prend moins d\'une minute et vous montre l\'essentiel.',
  },
  {
    icon: ArchiveBoxIcon,
    title: 'Publiez vos annonces',
    desc: "Vendez ou donnez les objets dont vous n'avez plus besoin. Les artisans et particuliers de la communauté pourront les récupérer pour leur donner une nouvelle vie.",
    tip: 'Vos annonces sont accessibles depuis « Mes annonces » dans votre profil.',
  },
  {
    icon: ArchiveBoxArrowDownIcon,
    title: 'Déposez vos objets en conteneur',
    desc: 'Demandez un dépôt dans un conteneur près de chez vous. Une fois validé, vous recevez un code pour ouvrir la box et un code-barres que le professionnel scannera pour le récupérer.',
    tip: 'Rendez-vous dans « Déposer un objet » depuis l\'accueil ou votre profil.',
  },
  {
    icon: ShoppingBagIcon,
    title: 'Services, formations & événements',
    desc: 'Réservez et payez en ligne nos prestations, formations et ateliers proposés par nos artisans et animateurs partenaires. Vos factures et contrats signés sont automatiquement archivés.',
    tip: 'Vous retrouvez tout dans « Mes factures » et « Mes réservations ».',
  },
  {
    icon: GlobeAltIcon,
    title: 'Suivez votre Upcycling Score',
    desc: "Chaque objet sauvé et chaque participation à un atelier vous fait gagner des points. Suivez vos économies en CO₂ et l'impact écologique de vos actions.",
    tip: 'Votre score est visible en haut de votre profil.',
  },
  {
    icon: CalendarDaysIcon,
    title: 'Planning, notifications & communauté',
    desc: "Retrouvez tous vos rendez-vous dans votre planning, recevez des notifications pour ne rien manquer, et échangez avec la communauté via le forum dédié.",
    tip: 'Vous pouvez activer les notifications push dans Paramètres pour être averti en temps réel.',
  },
]

const current = computed(() => steps[step.value])
const isLast = computed(() => step.value === steps.length - 1)

function prev() {
  if (step.value > 0) step.value--
}

async function next() {
  if (saving.value) return
  if (!isLast.value) {
    step.value++
    return
  }
  saving.value = true
  try {
    await userAuth.completeOnboarding()
    emit('done')
  } catch (e) {
    console.error('[onboarding] complete failed', e)
  } finally {
    saving.value = false
  }
}

function blockEscape(e) {
  if (!props.show) return
  if (e.key === 'Escape') {
    e.preventDefault()
    e.stopPropagation()
  }
}

onMounted(() => window.addEventListener('keydown', blockEscape, true))
onUnmounted(() => window.removeEventListener('keydown', blockEscape, true))

watch(() => props.show, (val) => {
  if (val) {
    step.value = 0
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>
