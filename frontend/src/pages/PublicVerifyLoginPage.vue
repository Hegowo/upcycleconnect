<template>
  <div class="min-h-[70vh] flex items-center justify-center py-16 px-4">
    <div class="w-full max-w-md text-center">

      <div v-if="state === 'loading'" class="space-y-4">
        <div class="w-14 h-14 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto"></div>
        <p class="text-[#40617f] text-sm">{{ t('public.verifyLogin.loading') }}</p>
      </div>

      <div v-else-if="state === 'success'" class="space-y-6">
        <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#dcfce7;">
          <ShieldCheckIcon class="w-10 h-10 text-[#006d35]" />
        </div>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">{{ t('public.verifyLogin.successTitle') }}</h1>
          <p class="text-[#40617f] mt-2">{{ t('public.verifyLogin.successText') }}</p>
        </div>
        <p class="text-xs text-gray-400">{{ t('public.verifyLogin.redirecting') }}</p>
      </div>

      <div v-else-if="state === 'expired'" class="space-y-6">
        <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#fef9c3;">
          <ClockIcon class="w-10 h-10" style="color:#854d0e;" />
        </div>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ t('public.verifyLogin.expiredTitle') }}</h1>
          <p class="text-[#40617f] mt-2">{{ t('public.verifyLogin.expiredText') }}</p>
        </div>
        <RouterLink to="/connexion" class="inline-flex items-center gap-2 px-8 py-3 rounded-xl font-bold text-white transition hover:opacity-90" style="background:#006d35;">
          {{ t('public.verifyLogin.expiredCta') }}
        </RouterLink>
      </div>

      <div v-else class="space-y-6">
        <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#fee2e2;">
          <XCircleIcon class="w-10 h-10 text-red-500" />
        </div>
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ t('public.verifyLogin.errorTitle') }}</h1>
          <p class="text-[#40617f] mt-2">{{ errorMsg }}</p>
        </div>
        <RouterLink to="/connexion" class="inline-flex items-center gap-2 px-8 py-3 rounded-xl font-bold text-[#006d35] border-2 border-[#006d35]/30 hover:bg-[#f0fdf4] transition">
          {{ t('public.verifyLogin.errorCta') }}
        </RouterLink>
      </div>

    </div>
  </div>
</template>

<script setup>import { ref, onMounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ShieldCheckIcon, ClockIcon, XCircleIcon } from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'

const { t } = useI18n()
const route    = useRoute()
const router   = useRouter()
const userAuth = useUserAuthStore()

const state    = ref('loading')
const errorMsg = ref(t('public.verifyLogin.errorDefault'))

onMounted(async () => {
  const token = route.query.token
  if (!token) {
    state.value = 'error'
    return
  }
  try {
    const res = await fetch(`/api/v1/auth/verify-login?token=${encodeURIComponent(token)}`)
    const json = await res.json().catch(() => ({}))

    if (res.ok) {
      userAuth.token = json.token
      userAuth.user  = json.user
      localStorage.setItem('user_token', json.token)
      localStorage.setItem('user_data', JSON.stringify(json.user))
      state.value = 'success'
      setTimeout(() => router.push('/'), 1500)
    } else if (res.status === 410) {
      state.value = 'expired'
    } else {
      errorMsg.value = json.message || errorMsg.value
      state.value = 'error'
    }
  } catch {
    errorMsg.value = t('public.verifyLogin.errorNetwork')
    state.value = 'error'
  }
})
</script>
