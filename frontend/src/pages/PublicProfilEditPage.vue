<template>
  <div class="bg-[#f7f9ff] min-h-screen py-10">
    <div class="max-w-2xl mx-auto px-6 space-y-6">

      <div class="flex gap-1 bg-white rounded-2xl p-1.5 shadow-sm border border-[#edf4ff]">
        <RouterLink to="/profil" class="flex-1 py-2.5 px-4 rounded-xl text-sm font-semibold text-center transition text-[#40617f] hover:bg-[#f8fafc]">
          Mon profil
        </RouterLink>
        <RouterLink to="/profil/parametres" class="flex-1 py-2.5 px-4 rounded-xl text-sm font-semibold text-center transition text-white shadow-sm" style="background:#006d35;">
          Paramètres
        </RouterLink>
      </div>

      <div class="bg-white rounded-2xl p-6 flex items-center gap-6">
        <div class="relative cursor-pointer group" @click="$refs.avatarInput.click()">
          <div class="w-24 h-24 rounded-2xl overflow-hidden bg-gradient-to-br from-[#006d35] to-[#1b8848] flex items-center justify-center shadow-sm">
            <img
              v-if="avatarPreview"
              :src="avatarPreview"
              class="w-full h-full object-cover"
              @error="onAvatarLoadError"
            />
            <span v-else class="text-white font-jakarta font-extrabold text-3xl">{{ userAuth.initials }}</span>
          </div>
          <div class="absolute inset-0 rounded-2xl bg-black/40 opacity-0 group-hover:opacity-100 transition flex items-center justify-center">
            <CameraIcon class="w-6 h-6 text-white" />
          </div>
        </div>
        <input ref="avatarInput" type="file" accept=".jpg,.jpeg,.png,.webp" class="hidden" @change="handleAvatarChange" />
        <div class="flex-1">
          <p class="font-semibold text-[#001d32] text-sm">{{ t('public.profileEdit.avatarLabel') }}</p>
          <p class="text-xs text-[#40617f] mt-0.5">{{ t('public.profileEdit.avatarHint') }}</p>
          <div class="flex items-center gap-3 mt-3">
            <button
              v-if="avatarFile"
              @click="uploadAvatar"
              :disabled="avatarLoading"
              class="px-4 py-1.5 text-xs font-semibold rounded-lg text-white transition"
              style="background-color:#006d35;"
            >
              {{ avatarLoading ? t('public.profileEdit.avatarUploading') : t('public.profileEdit.avatarUpload') }}
            </button>
            <span v-if="avatarSuccess" class="text-xs text-[#006d35] font-medium flex items-center gap-1">
              <CheckIcon class="w-3.5 h-3.5" /> {{ t('public.profileEdit.avatarSuccess') }}
            </span>
            <span v-if="avatarError" class="text-xs text-red-500">{{ avatarError }}</span>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 space-y-4">
        <h2 class="font-semibold text-[#001d32]">{{ t('public.profileEdit.infoTitle') }}</h2>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.fieldFirstName') }}</label>
            <input v-model="infoForm.first_name" type="text" class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm text-[#001d32] focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.fieldLastName') }}</label>
            <input v-model="infoForm.last_name" type="text" class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm text-[#001d32] focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div class="col-span-2">
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.fieldPhone') }}</label>
            <input v-model="infoForm.phone" type="tel" :placeholder="t('public.profileEdit.fieldPhonePlaceholder')" class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm text-[#001d32] focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
        </div>
        <div class="flex items-center gap-3 pt-1">
          <button @click="saveInfo" :disabled="infoLoading" class="px-5 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 disabled:opacity-60" style="background-color:#006d35;">
            {{ infoLoading ? t('public.profileEdit.savingInfo') : t('public.profileEdit.saveInfo') }}
          </button>
          <span v-if="infoSuccess" class="text-xs text-[#006d35] font-medium flex items-center gap-1">
            <CheckIcon class="w-3.5 h-3.5" /> {{ t('public.profileEdit.infoSuccess') }}
          </span>
          <span v-if="infoError" class="text-xs text-red-500">{{ infoError }}</span>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 space-y-4">
        <h2 class="font-semibold text-[#001d32]">{{ t('public.profileEdit.languageTitle') }}</h2>
        <p class="text-xs text-[#40617f] -mt-2">{{ t('public.profileEdit.languageHint') }}</p>
        <div class="max-w-xs">
          <LanguageSwitcher />
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 space-y-4">
        <h2 class="font-semibold text-[#001d32]">{{ t('public.profileEdit.passwordTitle') }}</h2>
        <div class="space-y-3">
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.fieldCurrentPassword') }}</label>
            <input v-model="pwForm.current" type="password" class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.fieldNewPassword') }}</label>
            <input v-model="pwForm.new" type="password" class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.fieldConfirmPassword') }}</label>
            <input v-model="pwForm.confirm" type="password" class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
        </div>
        <div class="flex items-center gap-3 pt-1">
          <button @click="changePassword" :disabled="pwLoading" class="px-5 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 disabled:opacity-60" style="background-color:#006d35;">
            {{ pwLoading ? t('public.profileEdit.savingPassword') : t('public.profileEdit.savePassword') }}
          </button>
          <span v-if="pwSuccess" class="text-xs text-[#006d35] font-medium flex items-center gap-1">
            <CheckIcon class="w-3.5 h-3.5" /> {{ t('public.profileEdit.passwordSuccess') }}
          </span>
          <span v-if="pwError" class="text-xs text-red-500">{{ pwError }}</span>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 space-y-4">
        <h2 class="font-semibold text-[#001d32]">{{ t('public.profileEdit.emailTitle') }}</h2>
        <p class="text-sm text-[#40617f]">{{ t('public.profileEdit.currentEmail') }} <span class="font-medium text-[#001d32]">{{ userAuth.user?.email }}</span></p>

        <div v-if="emailStep === 0">
          <button @click="emailStart" :disabled="emailLoading" class="px-5 py-2 text-sm font-semibold rounded-lg border-2 border-[#006d35] text-[#006d35] hover:bg-[#f0fdf4] transition disabled:opacity-60">
            {{ emailLoading ? t('public.profileEdit.startingEmailChange') : t('public.profileEdit.startEmailChange') }}
          </button>
          <span v-if="emailError" class="block text-xs text-red-500 mt-2">{{ emailError }}</span>
        </div>

        <div v-else-if="emailStep === 1" class="space-y-3">
          <div class="flex items-start gap-3 p-3 bg-blue-50 rounded-xl text-sm text-[#40617f]">
            <EnvelopeIcon class="w-4 h-4 mt-0.5 shrink-0 text-[#006d35]" />
            <span>{{ t('public.profileEdit.codeNotice1Pre') }} <strong>{{ userAuth.user?.email }}</strong>{{ t('public.profileEdit.codeNotice1Post') }}</span>
          </div>
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.codeCurrentLabel') }}</label>
            <input v-model="emailCode1" type="text" maxlength="6" placeholder="000000"
              class="mt-1 w-40 px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm text-center tracking-widest font-mono focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.newEmailLabel') }}</label>
            <input v-model="emailNew" type="email" :placeholder="t('public.profileEdit.newEmailPlaceholder')"
              class="mt-1 w-full px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div class="flex items-center gap-3">
            <button @click="emailVerifyCurrent" :disabled="emailLoading" class="px-5 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 disabled:opacity-60" style="background-color:#006d35;">
              {{ emailLoading ? t('public.profileEdit.verifyingCurrent') : t('public.profileEdit.verifyCurrent') }}
            </button>
            <button @click="emailStep = 0; emailError = ''" class="text-xs text-[#40617f] hover:underline">{{ t('public.profileEdit.cancel') }}</button>
            <span v-if="emailError" class="text-xs text-red-500">{{ emailError }}</span>
          </div>
        </div>

        <div v-else-if="emailStep === 2" class="space-y-3">
          <div class="flex items-start gap-3 p-3 bg-blue-50 rounded-xl text-sm text-[#40617f]">
            <EnvelopeIcon class="w-4 h-4 mt-0.5 shrink-0 text-[#006d35]" />
            <span>{{ t('public.profileEdit.codeNotice2Pre') }} <strong>{{ emailNew }}</strong>{{ t('public.profileEdit.codeNotice2Post') }}</span>
          </div>
          <div>
            <label class="text-xs font-medium text-[#40617f] uppercase tracking-wide">{{ t('public.profileEdit.codeNewLabel') }}</label>
            <input v-model="emailCode2" type="text" maxlength="6" placeholder="000000"
              class="mt-1 w-40 px-3 py-2 rounded-lg border border-[#e5e7eb] text-sm text-center tracking-widest font-mono focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div class="flex items-center gap-3">
            <button @click="emailVerifyNew" :disabled="emailLoading" class="px-5 py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 disabled:opacity-60" style="background-color:#006d35;">
              {{ emailLoading ? t('public.profileEdit.verifyingNew') : t('public.profileEdit.verifyNew') }}
            </button>
            <button @click="emailStep = 0; emailError = ''" class="text-xs text-[#40617f] hover:underline">{{ t('public.profileEdit.cancel') }}</button>
            <span v-if="emailError" class="text-xs text-red-500">{{ emailError }}</span>
          </div>
        </div>

        <div v-else-if="emailStep === 3" class="flex items-center gap-2 text-[#006d35] text-sm font-medium">
          <CheckCircleIcon class="w-5 h-5" />
          {{ t('public.profileEdit.doneEmailChange1') }} <strong>{{ emailNew }}</strong>.
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 space-y-3">
        <div class="flex items-start justify-between gap-4">
          <div class="flex-1">
            <h2 class="font-semibold text-[#001d32]">Notifications push</h2>
            <p class="text-xs text-[#40617f] mt-0.5">Recevez les alertes importantes (paiement, contrat, devis) directement sur votre appareil.</p>
          </div>
          <button
            v-if="pushStatus !== 'granted'"
            @click="enablePush"
            :disabled="pushLoading"
            class="shrink-0 inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white font-bold text-sm transition hover:opacity-90 disabled:opacity-60"
            style="background: linear-gradient(135deg, #006d35, #1b8848);"
          >
            <div v-if="pushLoading" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
            <BellAlertIcon v-else class="w-4 h-4" />
            Activer
          </button>
          <span v-else class="shrink-0 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-[#dcfce7] text-[#166534] text-xs font-bold">
            <CheckCircleIcon class="w-3.5 h-3.5" />
            Activées
          </span>
        </div>
        <p v-if="pushError" class="text-red-600 text-xs bg-red-50 p-2.5 rounded-lg border border-red-200">{{ pushError }}</p>
        <p v-if="pushStatus === 'denied'" class="text-amber-700 text-xs bg-amber-50 p-2.5 rounded-lg border border-amber-200">
          Les notifications ont été bloquées dans votre navigateur. Pour les réactiver : icône cadenas dans la barre d'adresse → Autoriser les notifications.
        </p>
      </div>

      <div class="bg-white rounded-2xl p-6 space-y-4">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="font-semibold text-[#001d32]">Clés d'accès (Passkeys)</h2>
            <p class="text-xs text-[#40617f] mt-0.5">Connectez-vous sans mot de passe avec votre empreinte digitale, Face ID ou PIN.</p>
          </div>
          <button
            @click="registerPasskey"
            :disabled="addingPasskey"
            class="shrink-0 px-4 py-2 rounded-xl text-sm font-semibold text-white transition disabled:opacity-50 flex items-center gap-1.5"
            style="background: linear-gradient(135deg, #006d35, #1b8848);"
          >
            <div v-if="addingPasskey" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
            <PlusIcon v-else class="w-3.5 h-3.5" />
            Ajouter
          </button>
        </div>

        <Transition name="fade">
          <p v-if="passkeyError" class="text-red-600 text-xs bg-red-50 p-2.5 rounded-lg border border-red-200">{{ passkeyError }}</p>
        </Transition>

        <div v-if="passkeys.length === 0" class="text-sm text-[#94a3b8] text-center py-6 border border-dashed border-[#e2e8f0] rounded-xl">
          Aucune clé d'accès enregistrée.
        </div>
        <div v-else class="space-y-2">
          <div v-for="pk in passkeys" :key="pk.id" class="flex items-center gap-3 p-3 rounded-xl bg-[#f8fafc] border border-[#e2e8f0]">
            <KeyIcon class="w-4 h-4 text-[#006d35] shrink-0" />
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-[#001d32] truncate">{{ pk.name }}</p>
              <p class="text-xs text-[#40617f]">
                Ajoutée le {{ formatDate(pk.created_at) }}
                <span v-if="pk.last_used_at"> · Dernière utilisation {{ formatDate(pk.last_used_at) }}</span>
              </p>
            </div>
            <button @click="deletePasskey(pk.id)" class="p-1.5 rounded-lg hover:bg-red-50 text-[#94a3b8] hover:text-red-500 transition">
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  CameraIcon,
  CheckIcon,
  CheckCircleIcon,
  EnvelopeIcon,
  KeyIcon,
  TrashIcon,
  PlusIcon,
  BellAlertIcon,
} from '@heroicons/vue/24/outline'
import { enablePushNotifications } from '@/utils/onesignal'
import LanguageSwitcher from '@/components/LanguageSwitcher.vue'

const { t } = useI18n()
const router   = useRouter()
const userAuth = useUserAuthStore()

const pushLoading = ref(false)
const pushError   = ref('')
const pushStatus  = ref('default')

function refreshPushStatus() {
  if (typeof Notification === 'undefined') return
  pushStatus.value = Notification.permission || 'default'
}

async function enablePush() {
  pushLoading.value = true
  pushError.value = ''
  try {
    const id = await enablePushNotifications()
    refreshPushStatus()
    if (!id && pushStatus.value !== 'granted') {
      pushError.value = "Activation refusée ou impossible. Vérifiez que les notifications ne sont pas bloquées."
    }
  } catch (e) {
    pushError.value = e.message || "Erreur lors de l'activation."
  } finally {
    pushLoading.value = false
  }
}

onMounted(() => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/parametres')
    return
  }
  infoForm.value.first_name = userAuth.user?.first_name || ''
  infoForm.value.last_name  = userAuth.user?.last_name  || ''
  infoForm.value.phone      = userAuth.user?.phone      || ''
  avatarPreview.value       = userAuth.user?.avatar_url || null
  fetchPasskeys()
  refreshPushStatus()
})

const avatarPreview = ref(null)
const avatarFile    = ref(null)
const avatarLoading = ref(false)
const avatarSuccess = ref(false)
const avatarError   = ref('')

function onAvatarLoadError() {
  avatarPreview.value = null
  if (userAuth.user) userAuth.user.avatar_url = null
}

function handleAvatarChange(e) {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    avatarError.value = t('public.profileEdit.avatarTooBig')
    return
  }
  avatarFile.value    = file
  avatarSuccess.value = false
  avatarError.value   = ''
  const reader = new FileReader()
  reader.onload = ev => { avatarPreview.value = ev.target.result }
  reader.readAsDataURL(file)
}

async function uploadAvatar() {
  if (!avatarFile.value) return
  avatarLoading.value = true
  avatarError.value   = ''
  avatarSuccess.value = false
  try {
    const fd = new FormData()
    fd.append('avatar', avatarFile.value)
    const res = await fetch('/api/v1/profile/avatar', {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
      body: fd,
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || t('public.profileEdit.errorGeneric'))
    avatarSuccess.value = true
    avatarFile.value    = null
    if (userAuth.user) userAuth.user.avatar_url = json.avatar_url
    await userAuth.fetchMe()
  } catch (err) {
    avatarError.value = err.message
  } finally {
    avatarLoading.value = false
  }
}

const infoForm    = ref({ first_name: '', last_name: '', phone: '' })
const infoLoading = ref(false)
const infoSuccess = ref(false)
const infoError   = ref('')

async function saveInfo() {
  infoLoading.value = true
  infoSuccess.value = false
  infoError.value   = ''
  try {
    const res = await fetch('/api/v1/profile/info', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({
        first_name: infoForm.value.first_name,
        last_name:  infoForm.value.last_name,
        phone:      infoForm.value.phone || null,
      }),
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || t('public.profileEdit.errorGeneric'))
    infoSuccess.value  = true
    userAuth.user = { ...userAuth.user, ...json }
    localStorage.setItem('user_data', JSON.stringify(userAuth.user))
  } catch (err) {
    infoError.value = err.message
  } finally {
    infoLoading.value = false
  }
}

const pwForm    = ref({ current: '', new: '', confirm: '' })
const pwLoading = ref(false)
const pwSuccess = ref(false)
const pwError   = ref('')

async function changePassword() {
  pwError.value   = ''
  pwSuccess.value = false
  if (pwForm.value.new !== pwForm.value.confirm) {
    pwError.value = t('public.profileEdit.passwordMismatch')
    return
  }
  if (pwForm.value.new.length < 8) {
    pwError.value = t('public.profileEdit.passwordTooShort')
    return
  }
  pwLoading.value = true
  try {
    const res = await fetch('/api/v1/profile/change-password', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({ current: pwForm.value.current, new: pwForm.value.new }),
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || t('public.profileEdit.errorGeneric'))
    pwSuccess.value = true
    pwForm.value = { current: '', new: '', confirm: '' }
  } catch (err) {
    pwError.value = err.message
  } finally {
    pwLoading.value = false
  }
}

const emailStep    = ref(0)
const emailLoading = ref(false)
const emailError   = ref('')
const emailCode1   = ref('')
const emailCode2   = ref('')
const emailNew     = ref('')

async function emailStart() {
  emailLoading.value = true
  emailError.value   = ''
  try {
    const res = await fetch('/api/v1/profile/email/start', {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || t('public.profileEdit.errorGeneric'))
    emailStep.value = 1
  } catch (err) {
    emailError.value = err.message
  } finally {
    emailLoading.value = false
  }
}

async function emailVerifyCurrent() {
  const newEmail = emailNew.value.trim().toLowerCase()
  if (!emailCode1.value.trim() || !newEmail) {
    emailError.value = t('public.profileEdit.missingFields')
    return
  }
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(newEmail)) {
    emailError.value = 'Adresse email invalide.'
    return
  }
  if (newEmail === (userAuth.user?.email || '').toLowerCase()) {
    emailError.value = 'Cette adresse est identique à votre email actuel.'
    return
  }
  emailLoading.value = true
  emailError.value   = ''
  try {
    const res = await fetch('/api/v1/profile/email/verify-current', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({ code: emailCode1.value.trim(), new_email: emailNew.value.trim() }),
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || t('public.profileEdit.errorGeneric'))
    emailStep.value = 2
    emailCode1.value = ''
  } catch (err) {
    emailError.value = err.message
  } finally {
    emailLoading.value = false
  }
}

async function emailVerifyNew() {
  if (!emailCode2.value) {
    emailError.value = t('public.profileEdit.missingCode')
    return
  }
  emailLoading.value = true
  emailError.value   = ''
  try {
    const res = await fetch('/api/v1/profile/email/verify-new', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({ code: emailCode2.value.trim() }),
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || t('public.profileEdit.errorGeneric'))
    emailStep.value = 3
    if (userAuth.user) userAuth.user.email = json.email
    localStorage.setItem('user_data', JSON.stringify(userAuth.user))
  } catch (err) {
    emailError.value = err.message
  } finally {
    emailLoading.value = false
  }
}

const passkeys      = ref([])
const addingPasskey = ref(false)
const passkeyError  = ref('')

function b64urlToBuffer(s) {
  const b = s.replace(/-/g, '+').replace(/_/g, '/')
  const padded = b.padEnd(b.length + (4 - b.length % 4) % 4, '=')
  const bin = atob(padded)
  const buf = new ArrayBuffer(bin.length)
  const view = new Uint8Array(buf)
  for (let i = 0; i < bin.length; i++) view[i] = bin.charCodeAt(i)
  return buf
}

function bufferToB64url(buf) {
  const bytes = new Uint8Array(buf)
  let s = ''
  for (const b of bytes) s += String.fromCharCode(b)
  return btoa(s).replace(/\+/g, '-').replace(/\//g, '_').replace(/=/g, '')
}

async function fetchPasskeys() {
  try {
    const res = await fetch('/api/v1/passkeys', {
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) passkeys.value = await res.json()
  } catch {}
}

async function registerPasskey() {
  if (addingPasskey.value) return
  passkeyError.value = ''
  addingPasskey.value = true
  try {
    const beginRes = await fetch('/api/v1/passkeys/register/begin', {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (!beginRes.ok) throw new Error('Erreur lors de la génération du challenge.')
    const opts = await beginRes.json()

    opts.publicKey.challenge = b64urlToBuffer(opts.publicKey.challenge)
    opts.publicKey.user.id   = b64urlToBuffer(opts.publicKey.user.id)
    if (opts.publicKey.excludeCredentials) {
      opts.publicKey.excludeCredentials = opts.publicKey.excludeCredentials.map(c => ({
        ...c, id: b64urlToBuffer(c.id),
      }))
    }

    const cred = await navigator.credentials.create(opts)

    const body = {
      id:    cred.id,
      rawId: bufferToB64url(cred.rawId),
      type:  cred.type,
      response: {
        attestationObject: bufferToB64url(cred.response.attestationObject),
        clientDataJSON:    bufferToB64url(cred.response.clientDataJSON),
      },
      clientExtensionResults: cred.getClientExtensionResults(),
    }

    const name = `Clé du ${new Date().toLocaleDateString('fr-FR')}`
    const completeRes = await fetch(`/api/v1/passkeys/register/complete?name=${encodeURIComponent(name)}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify(body),
    })
    if (!completeRes.ok) {
      const json = await completeRes.json()
      throw new Error(json.message || 'Enregistrement échoué.')
    }
    await fetchPasskeys()
  } catch (e) {
    if (e?.name === 'NotAllowedError') return
    passkeyError.value = e.message || "Erreur lors de l'enregistrement."
  } finally {
    addingPasskey.value = false
  }
}

async function deletePasskey(id) {
  try {
    const res = await fetch(`/api/v1/passkeys/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) passkeys.value = passkeys.value.filter(p => p.id !== id)
  } catch {}
}

function formatDate(dateStr) {
  return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric' })
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
