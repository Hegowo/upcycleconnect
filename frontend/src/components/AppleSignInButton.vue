<template>
  <div>
    <button
      @click="handleSignIn"
      :disabled="loading"
      class="w-full flex items-center justify-center gap-2.5 bg-black hover:bg-[#1a1a1a] active:bg-[#111] text-white font-semibold text-sm py-[11px] px-4 rounded-xl transition disabled:opacity-60 cursor-pointer disabled:cursor-not-allowed"
    >
      <svg v-if="!loading" class="w-4 h-4 shrink-0" viewBox="0 0 814 1000" fill="currentColor">
        <path d="M788.1 340.9c-5.8 4.5-108.2 62.2-108.2 190.5 0 148.4 130.3 200.9 134.2 202.2-.6 3.2-20.7 71.9-68.7 141.9-42.8 61.6-87.5 123.1-155.5 123.1s-85.5-39.5-164-39.5c-76 0-103.7 40.8-165.9 40.8s-105-57.8-155.5-127.4C46 702 0 562.4 0 426.2c0-227.1 148.4-347 294.3-347 74.8 0 136.9 49.2 184.3 49.2 45.2 0 116.5-52 201.4-52 32.4 0 134.9 2.9 209.6 103.3z"/>
        <path d="M491.4 85.4C520.8 51.3 542 4.9 542 0c0-.6 0-1.3-.1-1.9-28.4 1-61.8 19.9-87.3 45.3-23.8 23.4-44.6 59.2-44.6 103.4 0 .6.1 1.3.1 1.9.3 0 .6.1.9.1 30.6 0 64.4-20.3 79.4-63.4z" transform="translate(0 0)"/>
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
