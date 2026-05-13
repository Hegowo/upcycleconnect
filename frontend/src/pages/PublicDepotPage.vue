<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-4 sm:px-8 pt-8 sm:pt-12 pb-12 sm:pb-16">

      <div class="mb-8 flex items-start justify-between gap-4 flex-wrap">
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl sm:text-5xl tracking-tight">{{ t('public.depot.title') }}</h1>
          <p class="text-[#40617f] text-lg mt-4 max-w-2xl leading-relaxed">{{ t('public.depot.subtitle') }}</p>
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

        <!-- Photos -->
        <div ref="photosCard" class="col-span-12 lg:col-span-8 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <div class="flex items-center justify-between">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.photosTitle') }}</h2>
            <span class="text-[#40617f] text-xs uppercase tracking-wide">{{ t('public.depot.photosMin') }}</span>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
            <div v-for="slot in photoSlots" :key="slot.id">
              <!-- Slot vide -->
              <div
                v-if="!slot.preview"
                class="bg-[#edf4ff] border-2 border-dashed border-[#becabc] rounded-xl flex flex-col items-center justify-center cursor-pointer hover:bg-[#d8eaff] transition h-40 sm:h-52 group"
                @click="triggerUpload(slot.id)"
              >
                <input
                  :ref="el => inputRefs[slot.id] = el"
                  type="file"
                  accept="image/*"
                  class="hidden"
                  @change="onFileChange($event, slot.id)"
                />
                <PhotoIcon class="w-10 h-10 text-[#40617f]/40 mx-auto mb-2 group-hover:text-[#40617f]/60" />
                <p class="text-[#40617f] text-sm text-center font-medium px-2">{{ slotLabel(slot.id) }}</p>
              </div>

              <!-- Slot avec photo -->
              <div v-else class="relative rounded-xl overflow-hidden h-40 sm:h-52 group">
                <img :src="slot.preview" class="w-full h-full object-cover" />
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition flex items-center justify-center gap-2">
                  <button
                    @click="triggerUpload(slot.id)"
                    class="p-2 bg-white/90 rounded-lg hover:bg-white transition"
                    title="Changer"
                  >
                    <input
                      :ref="el => inputRefs[slot.id] = el"
                      type="file"
                      accept="image/*"
                      class="hidden"
                      @change="onFileChange($event, slot.id)"
                    />
                    <ArrowPathIcon class="w-4 h-4 text-[#001d32]" />
                  </button>
                  <button
                    @click="removePhoto(slot.id)"
                    class="p-2 bg-red-500/90 rounded-lg hover:bg-red-500 transition"
                    title="Supprimer"
                  >
                    <TrashIcon class="w-4 h-4 text-white" />
                  </button>
                </div>
                <div class="absolute bottom-2 left-2 right-2 flex items-center gap-1 bg-[#006d35]/90 text-white text-xs font-semibold px-2 py-1 rounded-lg">
                  <CheckCircleIcon class="w-3.5 h-3.5 shrink-0" />
                  {{ t('public.depot.photoAdded') }}
                </div>
              </div>
            </div>
          </div>

          <p class="text-[#40617f] text-sm">{{ t('public.depot.photoTip') }}</p>
        </div>

        <!-- Description -->
        <div ref="descCard" class="col-span-12 lg:col-span-4 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.descTitle') }}</h2>
          <div class="flex flex-col gap-5">
            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldName') }}</label>
              <input v-model="form.name" type="text" :placeholder="t('public.depot.fieldNamePlaceholder')" class="bg-[#edf4ff] px-3 py-3.5 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280]" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldDetails') }}</label>
              <textarea v-model="form.details" rows="5" :placeholder="t('public.depot.fieldDetailsPlaceholder')" class="bg-[#edf4ff] px-3 py-3 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280] resize-none" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">{{ t('public.depot.fieldWeight') }}</label>
              <input v-model.number="form.weight" type="number" step="0.1" min="0" :placeholder="t('public.depot.fieldWeightPlaceholder')" class="bg-[#edf4ff] px-3 py-3 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280]" />
            </div>
          </div>
        </div>

        <!-- Catégorie -->
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
                :class="form.categoryId === cat.id ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-[#edf4ff] text-[#001d32] border-transparent hover:border-[#006d35]/20'"
              >
                {{ cat.name }}
              </button>
            </div>
          </div>
        </div>

        <!-- État + impact + point de collecte + soumettre -->
        <div ref="conditionCard" class="col-span-12 lg:col-span-7 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.depot.conditionTitle') }}</h2>

          <div class="flex flex-wrap gap-3">
            <button
              v-for="cond in conditions"
              :key="cond.key"
              @click="form.condition = cond.key"
              class="px-5 py-2.5 rounded-full text-sm font-medium transition-all border-2"
              :class="form.condition === cond.key ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-[#edf4ff] text-[#001d32] border-transparent hover:border-[#006d35]/20'"
            >
              {{ cond.label }}
            </button>
          </div>

          <!-- Impact estimé dynamique -->
          <div class="bg-gradient-to-br from-[#f0fdf4] to-[#dcfce7] rounded-2xl p-5 border border-[#bbf7d0]">
            <div class="flex items-center gap-2 mb-3">
              <SparklesIcon class="w-4 h-4 text-[#006d35]" />
              <p class="text-[#006d35] font-semibold text-sm">{{ t('public.depot.impactLabel') }}</p>
            </div>
            <div class="grid grid-cols-3 gap-3">
              <div class="bg-white/70 rounded-xl p-3 text-center">
                <p class="text-xl font-bold text-[#006d35]">{{ impactData.co2 }}</p>
                <p class="text-[10px] text-[#40617f] mt-0.5 font-medium">kg CO₂ évité</p>
              </div>
              <div class="bg-white/70 rounded-xl p-3 text-center">
                <p class="text-xl font-bold text-[#006d35]">{{ impactData.score }}</p>
                <p class="text-[10px] text-[#40617f] mt-0.5 font-medium">pts impact</p>
              </div>
              <div class="bg-white/70 rounded-xl p-3 text-center">
                <p class="text-sm font-bold text-[#006d35] leading-tight">{{ impactData.level }}</p>
                <p class="text-[10px] text-[#40617f] mt-0.5 font-medium">niveau</p>
              </div>
            </div>
            <p class="text-xs text-[#40617f] mt-3 leading-relaxed">{{ impactData.tip }}</p>
          </div>

          <!-- Point de collecte -->
          <div class="flex flex-col gap-3">
            <label class="text-[#001d32] font-semibold text-sm flex items-center gap-1.5">
              <MapPinIcon class="w-4 h-4 text-[#006d35]" />
              Point de collecte
            </label>

            <div v-if="collectPointsLoading" class="text-[#40617f] text-sm">{{ t('common.loading') }}</div>
            <div v-else-if="collectPoints.length === 0" class="text-[#40617f] text-sm italic">Aucun point de collecte disponible pour le moment.</div>
            <div v-else class="flex flex-col gap-2">
              <!-- Badge adresse utilisateur -->
              <div v-if="userAddress" class="text-xs text-[#40617f] bg-[#edf4ff] rounded-lg px-3 py-2 flex items-center gap-1.5">
                <MapPinIcon class="w-3.5 h-3.5 text-[#006d35] shrink-0" />
                Points triés par distance depuis votre adresse
              </div>
              <div
                v-for="pt in sortedCollectPoints"
                :key="pt.id"
                @click="form.collectionPointId = pt.id"
                class="flex items-start gap-3 p-3 rounded-xl border-2 cursor-pointer transition-all"
                :class="form.collectionPointId === pt.id ? 'border-[#006d35] bg-[#f0fdf4]' : 'border-[#e5e7eb] hover:border-[#006d35]/30 bg-white'"
              >
                <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 mt-0.5" :class="form.collectionPointId === pt.id ? 'bg-[#006d35]' : 'bg-[#edf4ff]'">
                  <MapPinIcon class="w-4 h-4" :class="form.collectionPointId === pt.id ? 'text-white' : 'text-[#40617f]'" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-semibold text-[#001d32] text-sm">{{ pt.name }}</p>
                  <p class="text-xs text-[#40617f] mt-0.5">{{ pt.address }}, {{ pt.postal_code }} {{ pt.city }}</p>
                  <p v-if="pt.opening_hours" class="text-xs text-[#40617f] mt-0.5 flex items-center gap-1">
                    <ClockIcon class="w-3 h-3 shrink-0" />{{ pt.opening_hours }}
                  </p>
                </div>
                <div v-if="pt._distance != null" class="text-xs font-semibold text-[#006d35] shrink-0">
                  {{ pt._distance < 1 ? (pt._distance * 1000).toFixed(0) + 'm' : pt._distance.toFixed(1) + 'km' }}
                </div>
              </div>
            </div>
          </div>

          <!-- Bouton soumettre -->
          <div class="flex items-center justify-end pt-2 border-t border-[#f1f5f9]">
            <button
              :disabled="submitting || !canSubmit"
              class="px-8 py-4 rounded-xl text-white font-semibold text-base transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,0,0,0.1)] disabled:opacity-50 disabled:cursor-not-allowed"
              style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
              @click="submitForm"
            >
              <span v-if="submitting">{{ t('public.depot.submitting') }}</span>
              <span v-else>{{ t('public.depot.submit') }}</span>
            </button>
          </div>

          <p v-if="submitError" class="text-red-600 text-sm font-medium">{{ submitError }}</p>
        </div>

        <!-- Bannière bas -->
        <div class="col-span-12 rounded-[24px] overflow-hidden relative min-h-[160px] sm:h-48">
          <div class="absolute inset-0 bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] opacity-50" />
          <div class="absolute inset-0 bg-gradient-to-r from-[#f7f9ff] to-transparent" />
          <div class="relative p-6 sm:p-12 flex flex-col justify-center h-full">
            <h3 class="font-jakarta font-bold text-[#001d32] text-2xl mb-2">{{ t('public.depot.howItWorksTitle') }}</h3>
            <p class="text-[#40617f] text-base leading-relaxed max-w-lg">{{ t('public.depot.howItWorksDesc') }}</p>
          </div>
        </div>

      </div>
    </div>

    <!-- Tutorial overlay -->
    <Teleport to="body">
      <Transition name="tut-fade">
        <div v-if="showTutorial" class="fixed inset-0 z-50" @keydown.esc="closeTutorial" tabindex="-1">
          <div class="absolute inset-0" @click="closeTutorial" />

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

          <Transition name="tut-tooltip">
            <div
              v-if="spotlightRect"
              :key="tutorialStep"
              class="fixed z-[52] bg-white rounded-2xl shadow-2xl p-6 w-80 max-w-[calc(100vw-2rem)]"
              :style="tooltipStyle"
              @click.stop
            >
              <div class="flex gap-1.5 mb-5">
                <div v-for="(_, i) in tutorialSteps" :key="i" class="h-1.5 rounded-full transition-all duration-500" :class="i < tutorialStep ? 'bg-[#006d35]/40' : i === tutorialStep ? 'bg-[#006d35]' : 'bg-gray-200'" :style="{ flexGrow: i === tutorialStep ? 2 : 1 }" />
              </div>
              <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#006d35] to-[#1b8848] flex items-center justify-center mb-4 shadow-lg shadow-[#006d35]/20">
                <component :is="tutorialSteps[tutorialStep].icon" class="w-6 h-6 text-white" />
              </div>
              <p class="text-xs font-bold text-[#006d35] uppercase tracking-widest mb-1">Étape {{ tutorialStep + 1 }} / {{ tutorialSteps.length }}</p>
              <h3 class="font-jakarta font-bold text-[#001d32] text-lg leading-snug mb-2">{{ tutorialSteps[tutorialStep].title }}</h3>
              <p class="text-[#40617f] text-sm leading-relaxed">{{ tutorialSteps[tutorialStep].desc }}</p>
              <div v-if="tutorialSteps[tutorialStep].tip" class="mt-4 bg-[#edf4ff] rounded-xl px-3 py-2.5 flex items-start gap-2">
                <SparklesIcon class="w-4 h-4 text-[#006d35] shrink-0 mt-0.5" />
                <p class="text-xs text-[#001d32] leading-relaxed">{{ tutorialSteps[tutorialStep].tip }}</p>
              </div>
              <div class="mt-5 flex items-center justify-between gap-3 pt-4 border-t border-gray-100">
                <button @click="prevTutorialStep" class="flex items-center gap-1 text-sm text-[#40617f] hover:text-[#001d32] font-medium transition" :class="tutorialStep === 0 ? 'opacity-0 pointer-events-none' : ''">
                  <ArrowLeftIcon class="w-4 h-4" /> Précédent
                </button>
                <div class="flex gap-2 items-center">
                  <button @click="closeTutorial" class="px-3 py-1.5 rounded-lg text-sm text-[#40617f] hover:bg-gray-100 transition font-medium">Passer</button>
                  <button @click="nextTutorialStep" class="flex items-center gap-1.5 px-4 py-1.5 rounded-xl text-sm font-semibold text-white transition hover:opacity-90" style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);">
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
  PhotoIcon, CheckCircleIcon, QuestionMarkCircleIcon, ExclamationCircleIcon,
  SparklesIcon, PencilSquareIcon, TagIcon, ArrowLeftIcon, ArrowRightIcon,
  ArrowPathIcon, TrashIcon, MapPinIcon, ClockIcon,
} from '@heroicons/vue/24/outline'
import { publicGet, userApi } from '@/services/publicApi'
import { useUserAuthStore } from '@/stores/userAuth'

const { t } = useI18n()
const router = useRouter()
const userAuth = useUserAuthStore()

const isLoggedIn = computed(() => userAuth.isLoggedIn)
const userAddress = computed(() => userAuth.user?.address || null)

const form = ref({
  name: '',
  details: '',
  categoryId: null,
  condition: 'good',
  weight: null,
  collectionPointId: null,
})

// ── Photos ──
const photoSlots = ref([{ id: 1, file: null, preview: null }, { id: 2, file: null, preview: null }, { id: 3, file: null, preview: null }])
const inputRefs = ref({})

function triggerUpload(id) {
  inputRefs.value[id]?.click()
}

function onFileChange(event, id) {
  const file = event.target.files?.[0]
  if (!file) return
  const slot = photoSlots.value.find(s => s.id === id)
  if (!slot) return
  slot.file = file
  slot.preview = URL.createObjectURL(file)
}

function removePhoto(id) {
  const slot = photoSlots.value.find(s => s.id === id)
  if (!slot) return
  if (slot.preview) URL.revokeObjectURL(slot.preview)
  slot.file = null
  slot.preview = null
  if (inputRefs.value[id]) inputRefs.value[id].value = ''
}

function slotLabel(id) {
  return { 1: t('public.depot.photoMain'), 2: t('public.depot.photoDetail'), 3: t('public.depot.photoAngle') }[id] || ''
}

// ── Catégories ──
const categories = ref([])
const categoriesLoading = ref(true)

// ── Points de collecte ──
const collectPoints = ref([])
const collectPointsLoading = ref(true)
const userCoords = ref(null)

const conditions = computed(() => [
  { key: 'good', label: t('public.depot.condGood') },
  { key: 'fair', label: t('public.depot.condUsed') },
  { key: 'poor', label: t('public.depot.condBroken') },
])

// ── Impact estimé dynamique ──
const conditionMultiplier = { good: 3.0, fair: 2.5, poor: 1.5 }
const conditionLabel = { good: 'Excellent', fair: 'Moyen', poor: 'Faible' }

const impactData = computed(() => {
  const w = parseFloat(form.value.weight) || 0
  const mult = conditionMultiplier[form.value.condition] || 2.5
  const co2 = w > 0 ? (w * mult).toFixed(1) : '—'
  const score = w > 0 ? Math.round(w * mult * 1.5 + 5) : '—'
  const level = conditionLabel[form.value.condition] || '—'
  const tip = w > 0
    ? `En bon état, cet objet de ${w} kg économise jusqu'à ${(w * 3).toFixed(1)} kg de CO₂.`
    : 'Renseignez le poids pour voir votre impact environnemental estimé.'
  return { co2, score, level, tip }
})

// ── Points de collecte triés par distance ──
function haversineKm(lat1, lon1, lat2, lon2) {
  const R = 6371
  const dLat = (lat2 - lat1) * Math.PI / 180
  const dLon = (lon2 - lon1) * Math.PI / 180
  const a = Math.sin(dLat / 2) ** 2 + Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) * Math.sin(dLon / 2) ** 2
  return R * 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a))
}

const sortedCollectPoints = computed(() => {
  if (!userCoords.value) return collectPoints.value.map(p => ({ ...p, _distance: null }))
  return collectPoints.value
    .map(p => ({ ...p, _distance: haversineKm(userCoords.value.lat, userCoords.value.lon, p.latitude, p.longitude) }))
    .sort((a, b) => a._distance - b._distance)
})

async function geocodeUserAddress(address) {
  try {
    const res = await fetch(`https://api-adresse.data.gouv.fr/search/?q=${encodeURIComponent(address)}&limit=1`)
    const data = await res.json()
    const f = data.features?.[0]
    if (f) {
      userCoords.value = { lat: f.geometry.coordinates[1], lon: f.geometry.coordinates[0] }
      if (!form.value.collectionPointId && sortedCollectPoints.value.length > 0) {
        form.value.collectionPointId = sortedCollectPoints.value[0].id
      }
    }
  } catch {}
}

// ── Submit ──
const submitting = ref(false)
const submitError = ref('')

const canSubmit = computed(() =>
  isLoggedIn.value && form.value.name.trim().length > 0 && form.value.details.trim().length > 0
)

async function submitForm() {
  if (!isLoggedIn.value) { router.push('/connexion?redirect=/depot'); return }
  if (!canSubmit.value) { submitError.value = t('public.depot.errorRequired'); return }
  submitError.value = ''
  submitting.value = true
  try {
    const payload = {
      title:               form.value.name.trim(),
      description:         form.value.details.trim(),
      condition:           form.value.condition,
      category_id:         form.value.categoryId || null,
      estimated_weight:    form.value.weight && form.value.weight > 0 ? Number(form.value.weight) : null,
      collection_point_id: form.value.collectionPointId || null,
    }
    const res = await userApi('/deposits', { method: 'POST', body: JSON.stringify(payload) })
    router.push({ path: '/confirmation-depot', query: { id: res.id } })
  } catch (e) {
    submitError.value = e?.message || t('public.depot.errorGeneric')
  } finally {
    submitting.value = false
  }
}

// ── Tutorial ──
const showTutorial = ref(false)
const tutorialStep = ref(0)
const spotlightRect = ref(null)
const photosCard = ref(null)
const descCard = ref(null)
const categoryCard = ref(null)
const conditionCard = ref(null)
const sectionRefs = [photosCard, descCard, categoryCard, conditionCard]

const tutorialSteps = computed(() => [
  { title: t('public.depot.tutorialStep1Title'), desc: t('public.depot.tutorialStep1Desc'), tip: t('public.depot.tutorialStep1Tip'), icon: PhotoIcon },
  { title: t('public.depot.tutorialStep2Title'), desc: t('public.depot.tutorialStep2Desc'), tip: t('public.depot.tutorialStep2Tip'), icon: PencilSquareIcon },
  { title: t('public.depot.tutorialStep3Title'), desc: t('public.depot.tutorialStep3Desc'), icon: TagIcon },
  { title: t('public.depot.tutorialStep4Title'), desc: t('public.depot.tutorialStep4Desc'), tip: t('public.depot.tutorialStep4Tip'), icon: SparklesIcon },
])

function updateSpotlight() {
  const el = sectionRefs[tutorialStep.value]?.value
  if (!el) return
  el.scrollIntoView({ behavior: 'instant', block: 'center' })
  const r = el.getBoundingClientRect()
  spotlightRect.value = { top: r.top, left: r.left, width: r.width, height: r.height }
}

const tooltipStyle = computed(() => {
  if (!spotlightRect.value) return { display: 'none' }
  const r = spotlightRect.value, pad = 16, tw = 320, th = 360
  const vw = window.innerWidth, vh = window.innerHeight
  const top = r.top + r.height + pad + th < vh ? r.top + r.height + pad : r.top - th - pad > 0 ? r.top - th - pad : Math.max(pad, Math.min(r.top, vh - th - pad))
  const left = Math.max(pad, Math.min(r.left + r.width / 2 - tw / 2, vw - tw - pad))
  return { top: top + 'px', left: left + 'px' }
})

function openTutorial() { tutorialStep.value = 0; showTutorial.value = true; document.body.style.overflow = 'hidden'; nextTick(updateSpotlight) }
function closeTutorial() { showTutorial.value = false; document.body.style.overflow = ''; spotlightRect.value = null }
function nextTutorialStep() { tutorialStep.value < tutorialSteps.value.length - 1 ? tutorialStep.value++ : closeTutorial() }
function prevTutorialStep() { if (tutorialStep.value > 0) tutorialStep.value-- }
watch(tutorialStep, () => nextTick(updateSpotlight))
function onResize() { if (showTutorial.value) nextTick(updateSpotlight) }

// ── Init ──
async function loadCategories() {
  try { const res = await publicGet('/categories'); categories.value = res.data || [] }
  catch { categories.value = [] }
  finally { categoriesLoading.value = false }
}

async function loadCollectPoints() {
  try { collectPoints.value = await publicGet('/collection-points') }
  catch { collectPoints.value = [] }
  finally { collectPointsLoading.value = false }
}

onMounted(async () => {
  await Promise.all([loadCategories(), loadCollectPoints()])
  if (userAddress.value) {
    await geocodeUserAddress(userAddress.value)
  }
  window.addEventListener('resize', onResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', onResize)
  document.body.style.overflow = ''
  photoSlots.value.forEach(s => { if (s.preview) URL.revokeObjectURL(s.preview) })
})
</script>

<style scoped>
.tut-fade-enter-active { transition: opacity 0.3s ease; }
.tut-fade-leave-active { transition: opacity 0.25s ease; }
.tut-fade-enter-from, .tut-fade-leave-to { opacity: 0; }
.tut-tooltip-enter-active { transition: opacity 0.22s ease, transform 0.22s ease; }
.tut-tooltip-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.tut-tooltip-enter-from { opacity: 0; transform: translateY(10px) scale(0.96); }
.tut-tooltip-leave-to { opacity: 0; transform: translateY(-6px) scale(0.97); }
</style>
