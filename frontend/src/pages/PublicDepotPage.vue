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
          @click="showTutorial = true"
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

        <div class="col-span-12 lg:col-span-8 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
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

        <div class="col-span-12 lg:col-span-4 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
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

        <div class="col-span-12 lg:col-span-5 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
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

        <div class="col-span-12 lg:col-span-7 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
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

    <div
      v-if="showTutorial"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-[#001d32]/60"
      @click.self="showTutorial = false"
    >
      <div class="bg-white rounded-[24px] max-w-2xl w-full max-h-[90vh] overflow-y-auto shadow-2xl">
        <div class="sticky top-0 bg-white border-b border-[#f1f5f9] px-6 sm:px-8 py-5 flex items-center justify-between rounded-t-[24px]">
          <h2 class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ t('public.depot.tutorialTitle') }}</h2>
          <button @click="showTutorial = false" class="p-2 rounded-lg hover:bg-gray-100 transition">
            <XMarkIcon class="w-5 h-5 text-[#40617f]" />
          </button>
        </div>

        <div class="p-6 sm:p-8 flex flex-col gap-6">
          <p class="text-[#40617f] leading-relaxed">{{ t('public.depot.tutorialIntro') }}</p>

          <div v-for="(step, i) in tutorialSteps" :key="i" class="flex gap-4">
            <div class="w-10 h-10 shrink-0 rounded-full bg-gradient-to-br from-[#006d35] to-[#1b8848] text-white font-bold flex items-center justify-center">
              {{ i + 1 }}
            </div>
            <div class="flex-1">
              <h3 class="font-jakarta font-bold text-[#001d32] text-base">{{ step.title }}</h3>
              <p class="text-[#40617f] text-sm leading-relaxed mt-1">{{ step.desc }}</p>
            </div>
          </div>

          <div class="bg-[#edf4ff] rounded-xl p-4 flex items-start gap-3">
            <SparklesIcon class="w-5 h-5 text-[#006d35] shrink-0 mt-0.5" />
            <p class="text-sm text-[#001d32]">{{ t('public.depot.tutorialTip') }}</p>
          </div>

          <button
            @click="showTutorial = false"
            class="w-full py-3 rounded-xl text-white font-semibold transition hover:opacity-90"
            style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
          >
            {{ t('public.depot.tutorialClose') }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  PhotoIcon,
  CheckCircleIcon,
  QuestionMarkCircleIcon,
  ExclamationCircleIcon,
  XMarkIcon,
  SparklesIcon,
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
const showTutorial = ref(false)

function slotLabel(id) {
  return {
    1: t('public.depot.photoMain'),
    2: t('public.depot.photoDetail'),
    3: t('public.depot.photoAngle'),
  }[id] || ''
}

const conditions = computed(() => [
  { key: 'good', label: t('public.depot.condGood') },
  { key: 'fair', label: t('public.depot.condUsed') },
  { key: 'poor', label: t('public.depot.condBroken') },
])

const tutorialSteps = computed(() => [
  { title: t('public.depot.tutorialStep1Title'), desc: t('public.depot.tutorialStep1Desc') },
  { title: t('public.depot.tutorialStep2Title'), desc: t('public.depot.tutorialStep2Desc') },
  { title: t('public.depot.tutorialStep3Title'), desc: t('public.depot.tutorialStep3Desc') },
  { title: t('public.depot.tutorialStep4Title'), desc: t('public.depot.tutorialStep4Desc') },
])

const canSubmit = computed(() => {
  return isLoggedIn.value
    && form.value.name.trim().length > 0
    && form.value.details.trim().length > 0
})

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
  submitting.value = true
  try {
    const payload = {
      title: form.value.name.trim(),
      description: form.value.details.trim(),
      condition: form.value.condition,
      category_id: form.value.categoryId || null,
      estimated_weight: form.value.weight && form.value.weight > 0 ? Number(form.value.weight) : null,
    }
    const res = await userApi('/deposits', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
    router.push({ path: '/confirmation-depot', query: { id: res.id } })
  } catch (e) {
    submitError.value = e?.message || t('public.depot.errorGeneric')
  } finally {
    submitting.value = false
  }
}

onMounted(loadCategories)
</script>
