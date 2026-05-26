<template>
  <div>
    <button
      @click="handleSignIn"
      :disabled="loading"
      class="w-full flex items-center justify-center gap-2.5 bg-black hover:bg-[#1a1a1a] active:bg-[#111] text-white font-semibold text-sm py-[11px] px-4 rounded-xl transition disabled:opacity-60 cursor-pointer disabled:cursor-not-allowed"
    >
      <svg v-if="!loading" class="w-4 h-4 shrink-0" viewBox="0 0 24 24" fill="currentColor">
        <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.8-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
      </svg>
      <span v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin shrink-0" />
      {{ loading ? 'Connexion...' : 'Continuer avec Apple' }}
    </button>
    <p v-if="error" class="text-red-600 text-xs mt-2 text-center">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserAuthStore } from '@/stores/userAuth'

const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const loading = ref(false)
const error = ref('')

const SERVICES_ID = 'xyz.upcycleconnect.signin'
const REDIRECT_URI = 'https://upcycleconnect.xyz'

function loadAppleSDK() {
  return new Promise((resolve) => {
    if (window.AppleID) { resolve(); return }
    const script = document.createElement('script')
    script.src = 'https://appleid.cdn-apple.com/appleauth/static/jsapi/appleid/1/en_US/appleid.auth.js'
    script.async = true
    script.onload = resolve
    document.head.appendChild(script)
  })
}

async function handleSignIn() {
  error.value = ''
  loading.value = true
  try {
    await loadAppleSDK()

    window.AppleID.auth.init({
      clientId: SERVICES_ID,
      scope: 'name email',
      redirectURI: REDIRECT_URI,
      usePopup: true,
    })

    const response = await window.AppleID.auth.signIn()

    const idToken = response.authorization?.id_token
    if (!idToken) throw new Error('Réponse Apple invalide.')

    const firstName = response.user?.name?.firstName || ''
    const lastName = response.user?.name?.lastName || ''

    const res = await fetch('/api/v1/auth/apple', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id_token: idToken, first_name: firstName, last_name: lastName }),
    })
    const json = await res.json()
    if (!res.ok) throw new Error(json.message || 'Erreur Apple.')

    userAuth.token = json.token
    userAuth.user = json.user
    localStorage.setItem('user_token', json.token)
    localStorage.setItem('user_data', JSON.stringify(json.user))

    const redirect = route.query.redirect
    router.push(redirect ? String(redirect) : '/')
  } catch (e) {
    // L'utilisateur a fermé le popup Apple → pas d'erreur à afficher
    if (e?.error === 'popup_closed_by_user' || e?.error === 'user_trigger_new_signin_flow') return
    error.value = e.message || 'Connexion Apple échouée.'
  } finally {
    loading.value = false
  }
}
</script>
