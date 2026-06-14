<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-[#001d32]/85 backdrop-blur-sm" />

      <div class="relative bg-white rounded-[28px] shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col max-h-[92vh]">

        <div class="px-6 pt-6 pb-4 shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">
          <div class="flex items-center justify-between mb-3">
            <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-white/20 text-white text-[10px] font-bold uppercase tracking-wider">
              <SparklesIcon class="w-3.5 h-3.5" />
              Configuration de votre profil pro
            </span>
            <span class="text-white/70 text-xs font-semibold">{{ step + 1 }} / {{ steps.length }}</span>
          </div>
          <h2 class="font-jakarta font-extrabold text-white text-2xl leading-tight">{{ current.title }}</h2>
          <p v-if="current.subtitle" class="text-white/80 text-sm mt-1">{{ current.subtitle }}</p>
        </div>

        <div class="flex items-center gap-1.5 px-6 pt-4 shrink-0">
          <div
            v-for="(_, i) in steps"
            :key="i"
            class="h-1 flex-1 rounded-full transition-all"
            :class="i <= step ? 'bg-[#006d35]' : 'bg-[#e5e7eb]'"
          />
        </div>

        <div class="flex-1 overflow-y-auto px-6 py-6">

          <div v-if="current.key === 'welcome'" class="text-center py-4">
            <div class="w-20 h-20 rounded-2xl mx-auto mb-4 flex items-center justify-center" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <BuildingStorefrontIcon class="w-10 h-10 text-white" />
            </div>
            <p class="text-[#001d32] text-base leading-relaxed mb-4">
              Bienvenue <strong>{{ companyName }}</strong> sur UpcycleConnect !<br/>
              Votre compte vient d'être validé par notre équipe.
            </p>
            <p class="text-[#40617f] text-sm mb-6">
              Avant de commencer, prenons 2 minutes pour personnaliser votre profil professionnel afin qu'il soit attractif pour la communauté.
            </p>
            <div class="bg-[#edf4ff] rounded-xl px-4 py-3 inline-flex items-start gap-2.5 text-left">
              <LightBulbIcon class="w-5 h-5 text-[#006d35] shrink-0 mt-0.5" />
              <div>
                <p class="font-semibold text-[#001d32] text-sm">Astuce</p>
                <p class="text-[#40617f] text-xs mt-0.5">Vous pourrez modifier toutes ces informations plus tard depuis votre profil.</p>
              </div>
            </div>
          </div>

          <div v-else-if="current.key === 'description'" class="space-y-4">
            <p class="text-[#40617f] text-sm">
              Décrivez votre activité, vos spécialités et ce qui rend votre savoir-faire unique. Cette description sera visible publiquement sur votre fiche.
            </p>
            <textarea
              v-model="form.description"
              rows="7"
              placeholder="Ex : Spécialiste de la rénovation de meubles anciens depuis 15 ans, je redonne vie au mobilier en bois massif grâce à des techniques traditionnelles d'ébénisterie..."
              class="w-full px-4 py-3 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition resize-none"
            />
            <p class="text-xs text-[#94a3b8]">{{ form.description.length }} caractères — minimum 50 recommandé.</p>
          </div>

          <div v-else-if="current.key === 'website'" class="space-y-4">
            <p class="text-[#40617f] text-sm">
              Si vous avez un site web, un portfolio ou une page Instagram, partagez-le pour que la communauté puisse découvrir votre univers en détail.
            </p>
            <div class="relative">
              <GlobeAltIcon class="w-5 h-5 absolute left-4 top-1/2 -translate-y-1/2 text-[#94a3b8]" />
              <input
                v-model="form.website"
                type="url"
                placeholder="https://mon-site.com"
                class="w-full pl-12 pr-4 py-3 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
              />
            </div>
            <p class="text-xs text-[#94a3b8]">Optionnel — vous pouvez passer cette étape.</p>
          </div>

          <div v-else-if="current.key === 'prestation'" class="space-y-4">
            <p class="text-[#40617f] text-sm">
              Créez votre première prestation pour montrer ce que vous proposez. Vous pourrez en ajouter d'autres plus tard.
            </p>

            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wider mb-1.5">Titre de la prestation</label>
              <input
                v-model="form.prestationTitle"
                type="text"
                placeholder="Ex : Restauration de meubles en bois massif"
                class="w-full px-3 py-2.5 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wider mb-1.5">Description</label>
              <textarea
                v-model="form.prestationDescription"
                rows="3"
                placeholder="Décrivez en quelques lignes ce que comprend cette prestation..."
                class="w-full px-3 py-2.5 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition resize-none"
              />
            </div>

            <div class="grid grid-cols-2 gap-3">
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wider mb-1.5">Type de tarification</label>
                <select
                  v-model="form.priceType"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
                >
                  <option value="fixed">Prix fixe</option>
                  <option value="hourly">Par heure</option>
                  <option value="quote">Sur devis</option>
                </select>
              </div>
              <div v-if="form.priceType !== 'quote'">
                <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wider mb-1.5">Prix (€)</label>
                <input
                  v-model="form.price"
                  type="number"
                  step="0.01"
                  placeholder="0.00"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wider mb-1.5">Catégorie</label>
              <select
                v-model="form.categoryId"
                class="w-full px-3 py-2.5 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
              >
                <option :value="null">— Sans catégorie —</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
            <p class="text-xs text-[#94a3b8]">Optionnel — vous pouvez passer cette étape pour le faire plus tard.</p>
          </div>

          <div v-else-if="current.key === 'done'" class="text-center py-4">
            <div class="w-20 h-20 rounded-full mx-auto mb-4 flex items-center justify-center" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <CheckCircleIcon class="w-12 h-12 text-white" />
            </div>
            <p class="text-[#001d32] text-base leading-relaxed mb-3">
              Votre profil professionnel est configuré !
            </p>
            <p class="text-[#40617f] text-sm mb-6">
              Vous pouvez dès maintenant gérer vos projets, vos prestations, vos campagnes publicitaires et collecter les objets déposés par les particuliers.
            </p>
            <div class="grid grid-cols-2 gap-3 max-w-sm mx-auto">
              <div v-for="item in nextSteps" :key="item.label" class="bg-[#f7f9ff] rounded-xl p-3 text-left">
                <component :is="item.icon" class="w-5 h-5 text-[#006d35] mb-1.5" />
                <p class="text-xs font-semibold text-[#001d32]">{{ item.label }}</p>
                <p class="text-[10px] text-[#40617f] mt-0.5">{{ item.desc }}</p>
              </div>
            </div>
          </div>

          <p v-if="error" class="text-red-600 text-sm mt-4 text-center">{{ error }}</p>
        </div>

        <div class="px-6 py-4 bg-[#f8fafc] border-t border-gray-100 flex items-center justify-between gap-3 shrink-0">
          <button
            @click="prev"
            :disabled="step === 0 || saving"
            class="px-4 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] disabled:opacity-30 disabled:cursor-not-allowed hover:bg-white transition"
          >
            ← Précédent
          </button>

          <div class="flex items-center gap-2">
            <button
              v-if="canSkip"
              @click="next"
              :disabled="saving"
              class="px-4 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] hover:bg-white transition disabled:opacity-50"
            >
              Passer
            </button>
            <button
              @click="saveAndNext"
              :disabled="saving"
              class="px-6 py-2.5 rounded-xl text-white font-jakarta font-bold text-sm transition hover:opacity-90 flex items-center gap-2 disabled:opacity-60"
              style="background:linear-gradient(135deg,#006d35,#1b8848);"
            >
              <div v-if="saving" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
              <span v-else>{{ isLast ? "Terminer" : "Suivant" }}</span>
              <ArrowRightIcon v-if="!saving && !isLast" class="w-4 h-4" />
              <CheckCircleIcon v-else-if="!saving && isLast" class="w-4 h-4" />
            </button>
          </div>
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
  BuildingStorefrontIcon,
  GlobeAltIcon,
  CalendarDaysIcon,
  WrenchScrewdriverIcon,
  QrCodeIcon,
  MegaphoneIcon,
} from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'
import { userApi, publicGet } from '@/services/publicApi'

const props = defineProps({
  show: { type: Boolean, default: false },
  companyName: { type: String, default: '' },
})
const emit = defineEmits(['done'])

const userAuth = useUserAuthStore()
const step    = ref(0)
const saving  = ref(false)
const error   = ref('')
const categories = ref([])

const form = ref({
  description: '',
  website: '',
  prestationTitle: '',
  prestationDescription: '',
  priceType: 'fixed',
  price: '',
  categoryId: null,
})

const steps = [
  { key: 'welcome',     title: 'Configurons votre profil professionnel', subtitle: 'Quelques étapes pour bien démarrer.', canSkip: false },
  { key: 'description', title: 'Présentez votre activité',             subtitle: 'Une bonne description attire plus de clients.', canSkip: true  },
  { key: 'website',     title: 'Votre site ou portfolio',              subtitle: 'Renforcez votre crédibilité avec un lien externe.', canSkip: true  },
  { key: 'prestation',  title: 'Créez votre première prestation',     subtitle: 'Présentez vos services à la communauté.', canSkip: true  },
  { key: 'done',        title: 'Tout est prêt !',                      subtitle: 'Votre espace pro est désormais opérationnel.', canSkip: false },
]

const current = computed(() => steps[step.value])
const isLast  = computed(() => step.value === steps.length - 1)
const canSkip = computed(() => current.value.canSkip)

const nextSteps = [
  { icon: WrenchScrewdriverIcon, label: 'Projets',  desc: 'Documentez vos transformations' },
  { icon: CalendarDaysIcon,      label: 'Événements', desc: 'Organisez des ateliers' },
  { icon: QrCodeIcon,            label: 'Collecte', desc: 'Scannez les codes barres' },
  { icon: MegaphoneIcon,         label: 'Pub',      desc: 'Lancez des campagnes' },
]

function prev() {
  if (step.value > 0) step.value--
  error.value = ''
}

function next() {
  if (step.value < steps.length - 1) {
    step.value++
    error.value = ''
  }
}

async function saveAndNext() {
  if (saving.value) return
  error.value = ''

  if (current.value.key === 'welcome') {
    next()
    return
  }

  saving.value = true
  try {
    if (current.value.key === 'description' && form.value.description.trim()) {
      await userApi('/provider/profile', {
        method: 'PUT',
        body: JSON.stringify({ description: form.value.description.trim() }),
      })
    }
    if (current.value.key === 'website' && form.value.website.trim()) {
      await userApi('/provider/profile', {
        method: 'PUT',
        body: JSON.stringify({ website: form.value.website.trim() }),
      })
    }

    if (current.value.key === 'prestation' && form.value.prestationTitle.trim()) {
      await userApi('/provider/prestations', {
        method: 'POST',
        body: JSON.stringify({
          title: form.value.prestationTitle.trim(),
          description: form.value.prestationDescription.trim() || null,
          price_type: form.value.priceType,
          price: form.value.priceType === 'quote' ? null : (form.value.price || null),
          category_id: form.value.categoryId,
        }),
      })
    }

    if (isLast.value) {
      await userApi('/provider/onboarding/complete', { method: 'POST' })
      await userAuth.fetchMe()
      emit('done')
      return
    }
    next()
  } catch (e) {
    error.value = e.message || 'Erreur lors de l\'enregistrement.'
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

onMounted(async () => {
  window.addEventListener('keydown', blockEscape, true)
  try {
    const data = await publicGet('/categories')
    categories.value = data.data || data || []
  } catch {}
})
onUnmounted(() => window.removeEventListener('keydown', blockEscape, true))

watch(() => props.show, (val) => {
  if (val) {
    step.value = 0
    error.value = ''
    document.body.style.overflow = 'hidden'
  } else {
    document.body.style.overflow = ''
  }
})
</script>
