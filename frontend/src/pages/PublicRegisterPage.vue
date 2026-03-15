<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4" style="background: linear-gradient(135deg, #f0fdf4 0%, #f8fafc 55%, #eff6ff 100%);">
    <div class="w-full max-w-lg">

      <div class="text-center mb-8">
        <img src="/logoentier.png" class="h-14 mx-auto mb-3" alt="UpcycleConnect" />
        <h1 class="text-2xl font-bold tracking-tight" style="color: #103652;">Créer un compte</h1>
        <p class="text-gray-400 text-sm mt-1">Rejoignez la communauté UpcycleConnect</p>
      </div>

      <div v-if="accountType === 'provider'" class="flex items-center justify-center gap-2 mb-6">
        <div v-for="(label, idx) in ['Type', 'Entreprise', 'Compte']" :key="idx" class="flex items-center gap-2">
          <div class="flex items-center gap-1.5">
            <div class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold transition-all duration-300"
              :style="stepIndex > idx ? 'background:#1B8848;color:white;' : stepIndex === idx ? 'background:#103652;color:white;' : 'background:#e5e7eb;color:#9ca3af;'">
              <CheckIcon v-if="stepIndex > idx" class="w-3.5 h-3.5" />
              <span v-else>{{ idx + 1 }}</span>
            </div>
            <span class="text-xs hidden sm:block transition-all duration-300" :style="stepIndex === idx ? 'color:#103652;font-weight:600;' : 'color:#9ca3af;'">{{ label }}</span>
          </div>
          <div v-if="idx < 2" class="w-8 h-px transition-all duration-500" :style="stepIndex > idx ? 'background:#1B8848;' : 'background:#e5e7eb;'"></div>
        </div>
      </div>

      <div class="bg-white rounded-3xl shadow-xl border border-gray-100 overflow-hidden">
        <Transition name="step" mode="out-in">

          <div v-if="step === 'type'" key="type" class="p-8">
            <h2 class="text-lg font-bold text-gray-900 mb-1">Quel type de compte ?</h2>
            <p class="text-sm text-gray-400 mb-6">Sélectionnez le profil qui vous correspond</p>
            <div class="grid grid-cols-2 gap-4">
              <button
                @click="accountType = 'user'"
                class="group relative flex flex-col items-center gap-3 p-6 rounded-2xl border-2 transition-all duration-200 hover:scale-[1.02] active:scale-[0.98]"
                :style="accountType === 'user' ? 'border-color:#1B8848;background:#f0fdf4;' : 'border-color:#e5e7eb;background:white;'"
              >
                <div class="w-14 h-14 rounded-2xl flex items-center justify-center transition-transform duration-200 group-hover:scale-110"
                  :style="accountType === 'user' ? 'background:#dcfce7;' : 'background:#f9fafb;'">
                  <UserIcon class="w-7 h-7" :style="accountType === 'user' ? 'color:#1B8848;' : 'color:#9ca3af;'" />
                </div>
                <div>
                  <p class="font-semibold text-gray-900 text-sm">Particulier</p>
                  <p class="text-xs text-gray-400 mt-0.5">Je cherche des services</p>
                </div>
                <Transition name="badge">
                  <div v-if="accountType === 'user'" class="absolute top-2.5 right-2.5 w-5 h-5 rounded-full flex items-center justify-center text-white" style="background:#1B8848;">
                    <CheckIcon class="w-3 h-3" />
                  </div>
                </Transition>
              </button>

              <button
                @click="accountType = 'provider'"
                class="group relative flex flex-col items-center gap-3 p-6 rounded-2xl border-2 transition-all duration-200 hover:scale-[1.02] active:scale-[0.98]"
                :style="accountType === 'provider' ? 'border-color:#1B5B88;background:#eff6ff;' : 'border-color:#e5e7eb;background:white;'"
              >
                <div class="w-14 h-14 rounded-2xl flex items-center justify-center transition-transform duration-200 group-hover:scale-110"
                  :style="accountType === 'provider' ? 'background:#dbeafe;' : 'background:#f9fafb;'">
                  <WrenchScrewdriverIcon class="w-7 h-7" :style="accountType === 'provider' ? 'color:#1B5B88;' : 'color:#9ca3af;'" />
                </div>
                <div>
                  <p class="font-semibold text-gray-900 text-sm">Prestataire</p>
                  <p class="text-xs text-gray-400 mt-0.5">Je propose des services</p>
                </div>
                <Transition name="badge">
                  <div v-if="accountType === 'provider'" class="absolute top-2.5 right-2.5 w-5 h-5 rounded-full flex items-center justify-center text-white" style="background:#1B5B88;">
                    <CheckIcon class="w-3 h-3" />
                  </div>
                </Transition>
              </button>
            </div>

            <button
              @click="goToNext"
              :disabled="!accountType"
              class="mt-6 w-full py-3 rounded-xl text-white font-semibold transition-all duration-200 hover:opacity-90 active:scale-[0.99] disabled:opacity-40 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              style="background: linear-gradient(135deg, #1B8848, #10522b);"
            >
              Continuer
              <ArrowRightIcon class="w-4 h-4" />
            </button>
          </div>

          <div v-else-if="step === 'siret'" key="siret" class="p-8">
            <button @click="step = 'type'" class="inline-flex items-center gap-1 text-sm text-gray-400 hover:text-gray-600 transition mb-5">
              <ArrowLeftIcon class="w-4 h-4" /> Retour
            </button>

            <div class="flex items-center gap-3 mb-6">
              <div class="w-10 h-10 rounded-xl flex items-center justify-center" style="background:#dbeafe;">
                <MagnifyingGlassIcon class="w-5 h-5" style="color:#1B5B88;" />
              </div>
              <div>
                <h2 class="text-lg font-bold text-gray-900">Recherche par SIRET</h2>
                <p class="text-xs text-gray-400">Pré-remplissage automatique de votre entreprise</p>
              </div>
            </div>

            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1.5">Numéro SIRET <span class="text-xs font-normal text-gray-400">(14 chiffres)</span></label>
              <div class="relative">
                <input
                  v-model="siretInput"
                  @input="onSiretInput"
                  type="text"
                  maxlength="14"
                  placeholder="ex : 12345678901234"
                  class="w-full px-4 py-3 border-2 rounded-xl text-sm transition-all duration-300 outline-none font-mono tracking-widest pr-12"
                  :class="{
                    'border-green-400 bg-green-50/60': siretStatus === 'found',
                    'border-red-300 bg-red-50/60': siretStatus === 'error',
                    'border-gray-200 focus:border-blue-400': !siretStatus
                  }"
                />
                <div class="absolute right-4 top-1/2 -translate-y-1/2 w-6 h-6 flex items-center justify-center">
                  <div v-if="siretLoading" class="w-5 h-5 border-2 border-blue-400 border-t-transparent rounded-full animate-spin"></div>
                  <CheckCircleIcon v-else-if="siretStatus === 'found'" class="w-5 h-5 text-green-500" />
                  <XCircleIcon v-else-if="siretStatus === 'error'" class="w-5 h-5 text-red-400" />
                  <span v-else class="text-xs text-gray-300 font-mono tabular-nums">{{ siretInput.length }}/14</span>
                </div>
              </div>
            </div>

            <Transition name="card-pop">
              <div v-if="companyData" class="mb-4 rounded-2xl border-2 border-green-200 overflow-hidden" style="background:linear-gradient(135deg,#f0fdf4,#dcfce7);">
                <div class="p-4">
                  <div class="flex items-start gap-3">
                    <div class="w-11 h-11 rounded-xl flex items-center justify-center shrink-0 bg-white shadow-sm">
                      <BuildingOfficeIcon class="w-6 h-6 text-gray-400" />
                    </div>
                    <div class="flex-1 min-w-0">
                      <p class="font-bold text-gray-900 text-base leading-tight">{{ companyData.infolegales.nomcommercialrne || companyData.infolegales.nomcommercialinsee || companyData.infolegales.denorne || companyData.infolegales.denoinsee }}</p>
                      <p class="text-xs text-gray-500 mt-0.5 truncate">
                        {{ [companyData.infolegales.voieadressageinsee, companyData.infolegales.codepostalinsee, companyData.infolegales.villeinsee].filter(Boolean).join(', ') || 'Adresse non renseignée' }}
                      </p>
                      <div class="flex flex-wrap gap-1.5 mt-2">
                        <span v-if="companyData.infolegales.naflibinsee || companyData.infolegales.nafinsee" class="text-xs px-2 py-0.5 rounded-full bg-white/80 border border-green-200 text-green-700 font-medium">
                          {{ companyData.infolegales.naflibinsee || companyData.infolegales.nafinsee }}
                        </span>
                        <span v-if="companyData.infolegales.catjurlibinsee" class="text-xs px-2 py-0.5 rounded-full bg-white/80 border border-gray-200 text-gray-500">
                          {{ companyData.infolegales.catjurlibinsee }}
                        </span>
                        <span class="text-xs px-2 py-0.5 rounded-full bg-white/80 border border-gray-200 text-gray-500 font-mono">
                          SIRET {{ siretInput }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
                <button
                  @click="useCompanyData"
                  class="w-full py-2.5 text-white text-sm font-semibold transition hover:opacity-90 active:opacity-80 flex items-center justify-center gap-2"
                  style="background:#1B8848;"
                >
                  <SparklesIcon class="w-4 h-4" />
                  Utiliser ces informations
                </button>
              </div>
            </Transition>

            <Transition name="fade">
              <div v-if="siretStatus === 'error'" class="mb-4 p-3 rounded-xl bg-red-50 border border-red-200 flex items-start gap-2 text-sm text-red-600">
                <ExclamationTriangleIcon class="w-4 h-4 shrink-0 mt-0.5" />
                <span>Aucune entreprise trouvée pour ce SIRET. Vérifiez le numéro ou remplissez votre formulaire manuellement.</span>
              </div>
            </Transition>

            <div class="pt-4 border-t border-gray-100 text-center">
              <button @click="skipSiret" class="text-sm text-gray-400 hover:text-gray-600 transition hover:underline underline-offset-2">
                Je n'ai pas mon SIRET — remplir manuellement
              </button>
            </div>
          </div>

          <div v-else-if="step === 'form'" key="form" class="p-8">
            <button @click="goBack" class="inline-flex items-center gap-1 text-sm text-gray-400 hover:text-gray-600 transition mb-5">
              <ArrowLeftIcon class="w-4 h-4" /> Retour
            </button>

            <div class="flex items-center gap-3 mb-5">
              <div class="w-10 h-10 rounded-xl flex items-center justify-center" :style="accountType === 'provider' ? 'background:#dbeafe;' : 'background:#dcfce7;'">
                <WrenchScrewdriverIcon v-if="accountType === 'provider'" class="w-5 h-5" style="color:#1B5B88;" />
                <UserIcon v-else class="w-5 h-5" style="color:#1B8848;" />
              </div>
              <div>
                <h2 class="text-lg font-bold text-gray-900">{{ accountType === 'provider' ? 'Compte prestataire' : 'Vos informations' }}</h2>
                <p v-if="companyData" class="text-xs font-medium flex items-center gap-1" style="color:#1B8848;">
                  <SparklesIcon class="w-3.5 h-3.5" /> Pré-rempli depuis votre SIRET
                </p>
                <p v-else class="text-xs text-gray-400">Remplissez les champs ci-dessous</p>
              </div>
            </div>

            <form @submit.prevent="handleRegister" class="space-y-4">

              <template v-if="accountType === 'provider'">
                <div class="rounded-2xl border border-gray-100 p-4 space-y-3" style="background:#f8fafc;">
                  <p class="text-xs font-bold text-gray-400 uppercase tracking-wider flex items-center gap-1.5">
                    <BuildingOfficeIcon class="w-3.5 h-3.5" /> Entreprise
                  </p>
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">Nom de l'entreprise <span class="text-red-400">*</span></label>
                    <input v-model="form.companyName" type="text" required
                      class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors"
                      :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-gray-200 bg-white'"
                    />
                  </div>
                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <label class="block text-xs font-medium text-gray-600 mb-1">SIRET</label>
                      <input v-model="form.siret" type="text" maxlength="14"
                        class="w-full px-3 py-2.5 border rounded-xl text-xs focus:outline-none focus:border-blue-400 transition-colors font-mono"
                        :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-gray-200 bg-white'"
                        placeholder="14 chiffres"
                      />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-gray-600 mb-1">Activité</label>
                      <input v-model="form.activity" type="text"
                        class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors"
                        :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-gray-200 bg-white'"
                        placeholder="ex: Plomberie"
                      />
                    </div>
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">Adresse</label>
                    <input v-model="form.address" type="text"
                      class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors"
                      :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-gray-200 bg-white'"
                      placeholder="Rue, code postal, ville"
                    />
                  </div>
                </div>
              </template>

              <div class="space-y-3">
                <p v-if="accountType === 'provider'" class="text-xs font-bold text-gray-400 uppercase tracking-wider pt-1 flex items-center gap-1.5">
                  <UserIcon class="w-3.5 h-3.5" /> Contact
                </p>
                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">Prénom <span class="text-red-400">*</span></label>
                    <input v-model="form.firstName" type="text" required class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors bg-white" />
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-gray-600 mb-1">Nom <span class="text-red-400">*</span></label>
                    <input v-model="form.lastName" type="text" required class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors bg-white" />
                  </div>
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Email <span class="text-red-400">*</span></label>
                  <input v-model="form.email" type="email" required class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors bg-white" />
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Téléphone</label>
                  <input v-model="form.phone" type="tel" class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors bg-white" placeholder="06 12 34 56 78" />
                </div>
                <div>
                  <label class="block text-xs font-medium text-gray-600 mb-1">Mot de passe <span class="text-red-400">*</span></label>
                  <div class="relative">
                    <input v-model="form.password" :type="showPassword ? 'text' : 'password'" required minlength="8"
                      class="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm focus:outline-none focus:border-blue-400 transition-colors bg-white pr-10" />
                    <button type="button" @click="showPassword = !showPassword" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition">
                      <EyeSlashIcon v-if="showPassword" class="w-5 h-5" />
                      <EyeIcon v-else class="w-5 h-5" />
                    </button>
                  </div>
                  <div v-if="form.password" class="mt-2 flex gap-1">
                    <div v-for="i in 4" :key="i" class="h-1.5 flex-1 rounded-full transition-all duration-300"
                      :style="passwordStrength >= i ? `background:${strengthColor}` : 'background:#e5e7eb'"></div>
                  </div>
                  <p v-if="form.password" class="text-xs mt-1 transition-all" :style="`color:${strengthColor}`">{{ strengthLabel }}</p>
                </div>
              </div>

              <Transition name="fade">
                <p v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded-xl border border-red-200 flex items-start gap-2">
                  <ExclamationTriangleIcon class="w-4 h-4 shrink-0 mt-0.5" /><span>{{ error }}</span>
                </p>
              </Transition>

              <button type="submit" :disabled="loading"
                class="w-full py-3 rounded-xl text-white font-semibold transition-all duration-200 hover:opacity-90 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                style="background: linear-gradient(135deg, #1B8848, #10522b);"
              >
                <div v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                <CheckCircleIcon v-else class="w-4 h-4" />
                <span>{{ loading ? 'Création en cours...' : 'Créer mon compte' }}</span>
              </button>
            </form>
          </div>

        </Transition>
      </div>

      <p class="text-center text-sm text-gray-400 mt-6">
        Déjà inscrit ?
        <RouterLink to="/connexion" class="font-medium hover:underline underline-offset-2" style="color: #1B8848;">Se connecter</RouterLink>
      </p>
    </div>

    <span class="fixed bottom-3 left-3 text-xs text-gray-300 font-mono select-none">v1.3.0</span>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import {
  UserIcon,
  WrenchScrewdriverIcon,
  CheckIcon,
  CheckCircleIcon,
  XCircleIcon,
  MagnifyingGlassIcon,
  BuildingOfficeIcon,
  EyeIcon,
  EyeSlashIcon,
  SparklesIcon,
  ExclamationTriangleIcon,
  ArrowRightIcon,
  ArrowLeftIcon,
} from '@heroicons/vue/24/outline'

const step = ref('type')
const accountType = ref('')
const siretInput = ref('')
const siretLoading = ref(false)
const siretStatus = ref('')
const companyData = ref(null)
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const form = ref({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  password: '',
  companyName: '',
  siret: '',
  activity: '',
  address: '',
})

const stepIndex = computed(() => ({ type: 0, siret: 1, form: 2 }[step.value] ?? 0))

function goToNext() {
  if (accountType.value === 'provider') {
    step.value = 'siret'
  } else {
    step.value = 'form'
  }
}

function goBack() {
  step.value = accountType.value === 'provider' ? 'siret' : 'type'
}

function skipSiret() {
  companyData.value = null
  step.value = 'form'
}

let siretTimer = null
function onSiretInput() {
  const digits = siretInput.value.replace(/\D/g, '').slice(0, 14)
  siretInput.value = digits
  siretStatus.value = ''
  companyData.value = null
  clearTimeout(siretTimer)
  if (digits.length === 14) {
    siretTimer = setTimeout(() => fetchSiret(digits), 400)
  }
}

async function fetchSiret(siret) {
  siretLoading.value = true
  siretStatus.value = ''
  try {
    const res = await fetch(`/api/public/v1/siret/${siret}`)
    if (res.status === 404) { siretStatus.value = 'error'; return }
    if (!res.ok) throw new Error('server_error')
    const data = await res.json()
    if (data?.infolegales) {
      companyData.value = data
      siretStatus.value = 'found'
    } else {
      siretStatus.value = 'error'
    }
  } catch {
    siretStatus.value = 'error'
  } finally {
    siretLoading.value = false
  }
}

function useCompanyData() {
  const d = companyData.value?.infolegales
  if (!d) return
  form.value.companyName = d.nomcommercialrne || d.nomcommercialinsee || d.denorne || d.denoinsee || ''
  form.value.siret       = siretInput.value
  form.value.activity    = d.naflibinsee || d.nafinsee || ''
  form.value.address     = [d.voieadressageinsee, d.codepostalinsee, d.villeinsee].filter(Boolean).join(', ')
  step.value = 'form'
}

const passwordStrength = computed(() => {
  const p = form.value.password
  if (!p) return 0
  let s = 0
  if (p.length >= 8)          s++
  if (/[A-Z]/.test(p))        s++
  if (/[0-9]/.test(p))        s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return s
})

const strengthColor = computed(() => ['#e5e7eb','#ef4444','#f97316','#eab308','#22c55e'][passwordStrength.value])
const strengthLabel = computed(() => ['','Très faible','Faible','Moyen','Fort'][passwordStrength.value])

async function handleRegister() {
  loading.value = true
  error.value = ''
  loading.value = false
}
</script>

<style scoped>
.step-enter-active,
.step-leave-active {
  transition: opacity 0.2s ease, transform 0.25s ease;
}
.step-enter-from {
  opacity: 0;
  transform: translateX(20px);
}
.step-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

.card-pop-enter-active {
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.card-pop-leave-active {
  transition: all 0.2s ease;
}
.card-pop-enter-from {
  opacity: 0;
  transform: translateY(-8px) scale(0.96);
}
.card-pop-leave-to {
  opacity: 0;
  transform: scale(0.96);
}

.badge-enter-active {
  transition: all 0.25s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.badge-leave-active {
  transition: all 0.15s ease;
}
.badge-enter-from,
.badge-leave-to {
  opacity: 0;
  transform: scale(0);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
