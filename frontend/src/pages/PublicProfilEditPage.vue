<template>
  <div class="bg-[#f7f9ff] min-h-screen py-10">
    <div class="max-w-2xl mx-auto px-6 space-y-6">

      <div class="flex items-center gap-4">
        <button @click="router.push('/profil')" class="p-2 rounded-xl hover:bg-white transition">
          <ArrowLeftIcon class="w-5 h-5 text-[#40617f]" />
        </button>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ t('public.profileEdit.title') }}</h1>
          <p class="text-sm text-[#40617f]">{{ t('public.profileEdit.subtitle') }}</p>
        </div>
      </div>

      <div class="bg-white rounded-2xl p-6 flex items-center gap-6">
        <div class="relative cursor-pointer group" @click="$refs.avatarInput.click()">
          <div class="w-24 h-24 rounded-2xl overflow-hidden bg-gradient-to-br from-[#006d35] to-[#1b8848] flex items-center justify-center shadow-sm">
            <img v-if="avatarPreview" :src="avatarPreview" class="w-full h-full object-cover" />
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

    </div>
  </div>
</template>

<script setup>import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ArrowLeftIcon,
  CameraIcon,
  CheckIcon,
  CheckCircleIcon,
  EnvelopeIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const router   = useRouter()
const userAuth = useUserAuthStore()

onMounted(() => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/modifier')
    return
  }
  infoForm.value.first_name = userAuth.user?.first_name || ''
  infoForm.value.last_name  = userAuth.user?.last_name  || ''
  infoForm.value.phone      = userAuth.user?.phone      || ''
  avatarPreview.value       = userAuth.user?.avatar_url || null
})

const avatarPreview = ref(null)
const avatarFile    = ref(null)
const avatarLoading = ref(false)
const avatarSuccess = ref(false)
const avatarError   = ref('')

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
  if (!emailCode1.value || !emailNew.value) {
    emailError.value = t('public.profileEdit.missingFields')
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
      body: JSON.stringify({ code: emailCode1.value, new_email: emailNew.value }),
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
      body: JSON.stringify({ code: emailCode2.value }),
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
</script>
