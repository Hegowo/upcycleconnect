<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-4 sm:px-8 pt-8 sm:pt-12 pb-12 sm:pb-16">

      <div class="mb-8 flex items-start justify-between gap-4 flex-wrap">
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl sm:text-5xl tracking-tight">{{ t('public.depot.title') }}</h1>
          <p class="text-[#40617f] text-lg mt-4 max-w-2xl leading-relaxed">
            {{ t('public.depot.subtitle') }}
          </p>
        </div>
        <button
          @click="openTutorial"
          class="flex items-center gap-2 px-4 py-2.5 rounded-xl bg-white text-[#006d35] border-2 border-[#006d35]/15 hover:border-[#006d35]/40 font-semibold text-sm transition shrink-0"
        >
          <QuestionMarkCircleIcon class="w-5 h-5" />
          {{ t('public.depot.tutorialBtn') }}
        </button>
      </div>

      <div v-if="!isLoggedIn" class="mb-8 bg-amber-50 border border-amber-200 rounded-2xl p-6 flex items-start gap-4">
        <ExclamationCircleIcon class="w-6 h-6 text-amber-600 shrink-0 mt-0.5" />
        <div class="flex-1">
          <p class="font-semibold text-amber-900">{{ t('public.depot.authRequiredTitle') }}</p>
          <p class="text-amber-800 text-sm mt-1">{{ t('public.depot.authRequiredDesc') }}</p>
          <RouterLink to="/connexion?redirect=/depot" class="inline-block mt-3 px-4 py-2 rounded-lg bg-[#006d35] text-white text-sm font-semibold hover:opacity-90 transition">
            {{ t('public.depot.authRequiredCta') }}
          </RouterLink>
        </div>
      </div>

      <div class="grid grid-cols-12 gap-4 sm:gap-6">

        <div ref="photosCard" class="col-span-12 lg:col-span-8 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <div class="flex items-center justify-between">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.photosTitle') }}</h2>
            <span class="text-[#40617f] text-xs uppercase tracking-wide">{{ t('public.depot.photosMin') }}</span>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
            <div
              v-for="slot in photoSlots"
              :key="slot.id"
              class="bg-[#edf4ff] border-2 border-dashed border-[#becabc] rounded-xl flex flex-col items-center justify-center cursor-pointer hover:bg-[#d8eaff] transition h-40 sm:h-60 group"
              @click="triggerUpload(slot.id)"
            >
              <div v-if="!slot.file">
                <PhotoIcon class="w-10 h-10 text-[#40617f]/40 mx-auto mb-2 group-hover:text-[#40617f]/60" />
                <p class="text-[#40617f] text-sm text-center font-medium">{{ slotLabel(slot.id) }}</p>
              </div>
              <div v-else class="text-center">
                <CheckCircleIcon class="w-10 h-10 text-[#006d35] mx-auto mb-2" />
                <p class="text-[#006d35] text-sm font-semibold">{{ t('public.depot.photoAdded') }}</p>
              </div>
            </div>
          </div>

          <p class="text-[#40617f] text-sm">
            {{ t('public.depot.photoTip') }}
          </p>
        </div>

        <div ref="descCard" class="col-span-12 lg:col-span-4 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.descTitle') }}</h2>

          <div class="flex flex-col gap-5">
            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldName') }}</label>
              <input
                v-model="form.name"
                type="text"
                :placeholder="t('public.depot.fieldNamePlaceholder')"
                class="bg-[#edf4ff] px-3 py-3.5 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280]"
              />
            </div>

            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldDetails') }}</label>
              <textarea
                v-model="form.details"
                rows="6"
                :placeholder="t('public.depot.fieldDetailsPlaceholder')"
                class="bg-[#edf4ff] px-3 py-3 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280] resize-none"
              />
            </div>

            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldWeight') }}</label>
              <input
                v-model.number="form.weight"
                type="number"
                step="0.1"
                min="0"
                :placeholder="t('public.depot.fieldWeightPlaceholder')"
                class="bg-[#edf4ff] px-3 py-3 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280]"
              />
            </div>
          </div>
        </div>

        <div ref="categoryCard" class="col-span-12 lg:col-span-5 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.classificationTitle') }}</h2>

          <div class="flex flex-col gap-3">
            <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldCategory') }}</label>
            <div v-if="categoriesLoading" class="text-[#40617f] text-sm">{{ t('common.loading') }}</div>
            <div v-else-if="categories.length === 0" class="text-[#40617f] text-sm italic">{{ t('public.depot.noCategories') }}</div>
            <div v-else class="grid grid-cols-2 gap-3">
              <button
                v-for="cat in categories"
                :key="cat.id"
                @click="form.categoryId = cat.id"
                class="px-5 py-3 rounded-xl text-sm font-medium border-2 transition-all"
                :class="form.categoryId === cat.id
                  ? 'bg-[#006d35] text-white border-[#006d35]'
                  : 'bg-[#edf4ff] text-[#001d32] border-transparent hover:border-[#006d35]/20'"
              >
                {{ cat.name }}
              </button>
            </div>
          </div>
        </div>

        <div ref="conditionCard" class="col-span-12 lg:col-span-7 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.conditionTitle') }}</h2>

          <div class="flex flex-wrap gap-3">
            <button
              v-for="cond in conditions"
              :key="cond.key"
              @click="form.condition = cond.key"
              class="px-5 py-2.5 rounded-full text-sm font-medium transition-all border-2"
              :class="form.condition === cond.key
                ? 'bg-[#006d35] text-white border-[#006d35]'
                : 'bg-[#edf4ff] text-[#001d32] border-transparent hover:border-[#006d35]/20'"
            >
              {{ cond.label }}
            </button>
          </div>

          <div class="bg-[rgba(27,136,72,0.1)] rounded-xl p-4 sm:p-6 flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:justify-between">
            <div class="flex items-center gap-4">
              <div class="w-6 h-6 bg-[#1b8848] rounded-full shrink-0" />
              <div>
                <p class="text-[#006d35] font-semibold text-base">{{ t('public.depot.impactLabel') }}</p>
                <p class="text-[#40617f] text-sm">{{ estimatedImpact }}</p>
              </div>
            </div>
            <button
              :disabled="submitting || !canSubmit"
              class="w-full sm:w-auto px-8 sm:px-10 py-4 rounded-xl text-white font-semibold text-base sm:text-lg transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,0,0,0.1)] disabled:opacity-50 disabled:cursor-not-allowed"
              style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
              @click="submitForm"
            >
              <span v-if="submitting">{{ t('public.depot.submitting') }}</span>
              <span v-else>{{ t('public.depot.submit') }}</span>
            </button>
          </div>

          <p v-if="submitError" class="text-red-600 text-sm font-medium">{{ submitError }}</p>
        </div>

        <div class="col-span-12 rounded-[24px] overflow-hidden relative min-h-[160px] sm:h-48">
          <div class="absolute inset-0 bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] opacity-50" />
          <div class="absolute inset-0 bg-gradient-to-r from-[#f7f9ff] to-transparent" />
          <div class="relative p-6 sm:p-12 flex flex-col justify-center h-full">
            <h3 class="font-jakarta font-bold text-[#001d32] text-2xl mb-2">{{ t('public.depot.howItWorksTitle') }}</h3>
            <p class="text-[#40617f] text-base leading-relaxed max-w-lg">
              {{ t('public.depot.howItWorksDesc') }}
            </p>
          </div>
        </div>

      </div>
    </div>

    <!-- ── Interactive Tutorial Overlay ── -->
    <Teleport to="body">
      <Transition name="tut-fade">
        <div v-if="showTutorial" class="fixed inset-0 z-50" @keydown.esc="closeTutorial" tabindex="-1">

          <!-- Dark backdrop — click anywhere dark to close -->
          <div class="absolute inset-0" @click="closeTutorial" />

          <!-- Spotlight: box-shadow creates the dark overlay, slides between sections -->
          <div
            v-if="spotlightRect"
            class="fixed rounded-[28px] pointer-events-none z-[51]"
            :style="{
              top: (spotlightRect.top - 12) + 'px',
              left: (spotlightRect.left - 12) + 'px',
              width: (spotlightRect.width + 24) + 'px',
              height: (spotlightRect.height + 24) + 'px',
              boxShadow: '0 0 0 9999px rgba(0,29,50,0.82), 0 0 0 2px rgba(27,136,72,0.9), 0 0 0 6px rgba(27,136,72,0.2), 0 8px 40px rgba(0,0,0,0.4)',
              transition: 'top 0.55s cubic-bezier(0.4,0,0.2,1), left 0.55s cubic-bezier(0.4,0,0.2,1), width 0.55s cubic-bezier(0.4,0,0.2,1), height 0.55s cubic-bezier(0.4,0,0.2,1)',
            }"
          />

          <!-- Tooltip card — new one mounts per step (triggers enter animation) -->
          <Transition name="tut-tooltip">
            <div
              v-if="spotlightRect"
              :key="tutorialStep"
              class="fixed z-[52] bg-white rounded-2xl shadow-2xl p-6 w-80 max-w-[calc(100vw-2rem)]"
              :style="tooltipStyle"
              @click.stop
            >
              <!-- Progress bar -->
              <div class="flex gap-1.5 mb-5">
                <div
                  v-for="(_, i) in tutorialSteps"
                  :key="i"
                  class="h-1.5 rounded-full transition-all duration-500"
                  :class="i < tutorialStep ? 'bg-[#006d35]/40' : i === tutorialStep ? 'bg-[#006d35]' : 'bg-gray-200'"
                  :style="{ flexGrow: i === tutorialStep ? 2 : 1 }"
                />
              </div>

              <!-- Icon -->
              <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#006d35] to-[#1b8848] flex items-center justify-center mb-4 shadow-lg shadow-[#006d35]/20">
                <component :is="tutorialSteps[tutorialStep].icon" class="w-6 h-6 text-white" />
              </div>

              <!-- Step label -->
              <p class="text-xs font-bold text-[#006d35] uppercase tracking-widest mb-1">
                Étape {{ tutorialStep + 1 }} / {{ tutorialSteps.length }}
              </p>

              <!-- Content -->
              <h3 class="font-jakarta font-bold text-[#001d32] text-lg leading-snug mb-2">
                {{ tutorialSteps[tutorialStep].title }}
              </h3>
              <p class="text-[#40617f] text-sm leading-relaxed">
                {{ tutorialSteps[tutorialStep].desc }}
              </p>

              <!-- Tip -->
              <div v-if="tutorialSteps[tutorialStep].tip" class="mt-4 bg-[#edf4ff] rounded-xl px-3 py-2.5 flex items-start gap-2">
                <SparklesIcon class="w-4 h-4 text-[#006d35] shrink-0 mt-0.5" />
                <p class="text-xs text-[#001d32] leading-relaxed">{{ tutorialSteps[tutorialStep].tip }}</p>
              </div>

              <!-- Navigation -->
              <div class="mt-5 flex items-center justify-between gap-3 pt-4 border-t border-gray-100">
                <button
                  @click="prevTutorialStep"
                  class="flex items-center gap-1 text-sm text-[#40617f] hover:text-[#001d32] font-medium transition"
                  :class="tutorialStep === 0 ? 'opacity-0 pointer-events-none' : ''"
                >
                  <ArrowLeftIcon class="w-4 h-4" />
                  Précédent
                </button>

                <div class="flex gap-2 items-center">
                  <button
                    @click="closeTutorial"
                    class="px-3 py-1.5 rounded-lg text-sm text-[#40617f] hover:bg-gray-100 transition font-medium"
                  >
                    Passer
                  </button>
                  <button
                    @click="nextTutorialStep"
                    class="flex items-center gap-1.5 px-4 py-1.5 rounded-xl text-sm font-semibold text-white transition hover:opacity-90"
                    style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
                  >
                    <span>{{ tutorialStep === tutorialSteps.length - 1 ? 'Terminer' : 'Suivant' }}</span>
                    <CheckCircleIcon v-if="tutorialStep === tutorialSteps.length - 1" class="w-4 h-4" />
                    <ArrowRightIcon v-else class="w-4 h-4" />
                  </button>
                </div>
              </div>
            </div>
          </Transition>

        </div>
      </Transition>
    </Teleport>

  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick, watch } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  PhotoIcon,
  CheckCircleIcon,
  QuestionMarkCircleIcon,
  ExclamationCircleIcon,
  XMarkIcon,
  SparklesIcon,
  PencilSquareIcon,
  TagIcon,
  ArrowLeftIcon,
  ArrowRightIcon,
} from '@heroicons/vue/24/outline'
import { publicGet, userApi } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'

const { t } = useI18n()
const router = useRouter()
const userAuth = useUserAuthStore()

const isLoggedIn = computed(() => userAuth.isLoggedIn)

const form = ref({
  name: '',
  details: '',
  categoryId: null,
  condition: 'good',
  weight: null,
})

const photoSlots = ref([
  { id: 1, file: null },
  { id: 2, file: null },
  { id: 3, file: null },
])

const categories = ref([])
const categoriesLoading = ref(true)
const submitting = ref(false)
const submitError = ref('')

// ── Tutorial state ──
const showTutorial = ref(false)
const tutorialStep = ref(0)
const spotlightRect = ref(null)

// Section refs
const photosCard    = ref(null)
const descCard      = ref(null)
const categoryCard  = ref(null)
const conditionCard = ref(null)

const sectionRefs = [photosCard, descCard, categoryCard, conditionCard]

const tutorialSteps = computed(() => [
  {
    title: t('public.depot.tutorialStep1Title'),
    desc:  t('public.depot.tutorialStep1Desc'),
    tip:   t('public.depot.tutorialStep1Tip'),
    icon:  PhotoIcon,
  },
  {
    title: t('public.depot.tutorialStep2Title'),
    desc:  t('public.depot.tutorialStep2Desc'),
    tip:   t('public.depot.tutorialStep2Tip'),
    icon:  PencilSquareIcon,
  },
  {
    title: t('public.depot.tutorialStep3Title'),
    desc:  t('public.depot.tutorialStep3Desc'),
    icon:  TagIcon,
  },
  {
    title: t('public.depot.tutorialStep4Title'),
    desc:  t('public.depot.tutorialStep4Desc'),
    tip:   t('public.depot.tutorialStep4Tip'),
    icon:  SparklesIcon,
  },
])

function updateSpotlight() {
  const el = sectionRefs[tutorialStep.value]?.value
  if (!el) return
  el.scrollIntoView({ behavior: 'instant', block: 'center' })
  const rect = el.getBoundingClientRect()
  spotlightRect.value = { top: rect.top, left: rect.left, width: rect.width, height: rect.height }
}

const tooltipStyle = computed(() => {
  if (!spotlightRect.value) return { display: 'none' }
  const r   = spotlightRect.value
  const pad = 16
  const tw  = 320
  const th  = 360
  const vw  = window.innerWidth
  const vh  = window.innerHeight

  let top
  if (r.top + r.height + pad + th < vh) {
    top = r.top + r.height + pad
  } else if (r.top - pad - th > 0) {
    top = r.top - th - pad
  } else {
    top = Math.max(pad, Math.min(r.top, vh - th - pad))
  }

  const left = Math.max(pad, Math.min(r.left + r.width / 2 - tw / 2, vw - tw - pad))
  return { top: top + 'px', left: left + 'px' }
})

function openTutorial() {
  tutorialStep.value = 0
  showTutorial.value = true
  document.body.style.overflow = 'hidden'
  nextTick(updateSpotlight)
}

function closeTutorial() {
  showTutorial.value = false
  document.body.style.overflow = ''
  spotlightRect.value = null
}

function nextTutorialStep() {
  if (tutorialStep.value < tutorialSteps.value.length - 1) {
    tutorialStep.value++
  } else {
    closeTutorial()
  }
}

function prevTutorialStep() {
  if (tutorialStep.value > 0) tutorialStep.value--
}

watch(tutorialStep, () => nextTick(updateSpotlight))

function onResize() {
  if (showTutorial.value) nextTick(updateSpotlight)
}

// ── Form helpers ──
function slotLabel(id) {
  return {
    1: t('public.depot.photoMain'),
    2: t('public.depot.photoDetail'),
    3: t('public.depot.photoAngle'),
  }[id] || ''
}

const conditions = computed(() => [
  { key: 'good',  label: t('public.depot.condGood') },
  { key: 'fair',  label: t('public.depot.condUsed') },
  { key: 'poor',  label: t('public.depot.condBroken') },
])

const canSubmit = computed(() =>
  isLoggedIn.value
  && form.value.name.trim().length > 0
  && form.value.details.trim().length > 0
)

const estimatedImpact = computed(() => {
  const w = parseFloat(form.value.weight)
  if (!w || w <= 0) return t('public.depot.impactValue')
  return t('public.depot.impactDynamic', { co2: (w * 2.5).toFixed(1) })
})

const triggerUpload = (id) => {
  const slot = photoSlots.value.find(s => s.id === id)
  if (slot) slot.file = slot.file ? null : 'placeholder'
}

async function loadCategories() {
  categoriesLoading.value = true
  try {
    const res = await publicGet('/categories')
    categories.value = res.data || []
  } catch {
    categories.value = []
  } finally {
    categoriesLoading.value = false
  }
}

async function submitForm() {
  if (!isLoggedIn.value) {
    router.push('/connexion?redirect=/depot')
    return
  }
  if (!canSubmit.value) {
    submitError.value = t('public.depot.errorRequired')
    return
  }
  submitError.value = ''
  submitting.value  = true
  try {
    const payload = {
      title:            form.value.name.trim(),
      description:      form.value.details.trim(),
      condition:        form.value.condition,
      category_id:      form.value.categoryId || null,
      estimated_weight: form.value.weight && form.value.weight > 0 ? Number(form.value.weight) : null,
    }
    const res = await userApi('/deposits', { method: 'POST', body: JSON.stringify(payload) })
    router.push({ path: '/confirmation-depot', query: { id: res.id } })
  } catch (e) {
    submitError.value = e?.message || t('public.depot.errorGeneric')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadCategories()
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  document.body.style.overflow = ''
})
</script>

<style scoped>
/* Overlay fade */
.tut-fade-enter-active { transition: opacity 0.3s ease; }
.tut-fade-leave-active { transition: opacity 0.25s ease; }
.tut-fade-enter-from,
.tut-fade-leave-to    { opacity: 0; }

/* Tooltip pop per step */
.tut-tooltip-enter-active { transition: opacity 0.22s ease, transform 0.22s ease; }
.tut-tooltip-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.tut-tooltip-enter-from   { opacity: 0; transform: translateY(10px) scale(0.96); }
.tut-tooltip-leave-to     { opacity: 0; transform: translateY(-6px) scale(0.97); }
</style>
