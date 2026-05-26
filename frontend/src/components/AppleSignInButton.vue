<template>
  <div>
    <button
      @click="handleSignIn"
      :disabled="loading"
      class="w-full flex items-center justify-center gap-2.5 bg-black hover:bg-[#1a1a1a] active:bg-[#111] text-white font-semibold text-sm py-[11px] px-4 rounded-xl transition disabled:opacity-60 cursor-pointer disabled:cursor-not-allowed"
    >
      <svg v-if="!loading" class="w-4 h-4 shrink-0" viewBox="0 0 384 512" fill="currentColor">
        <path d="M318.7 268.7c-.2-36.7 16.4-64.4 50-84.8-18.8-26.9-47.2-41.7-84.7-44.6-35.5-2.8-74.3 20.7-88.5 20.7-15 0-49.4-19.7-76.4-19.7C63.3 141.2 4 184.8 4 273.5q0 39.3 14.4 81.2c12.8 36.7 59 126.7 107.2 125.2 25.2-.6 43-17.9 75.8-17.9 31.8 0 48.3 17.9 76.4 17.9 48.6-.7 90.4-82.5 102.6-119.3-65.2-30.7-61.7-90-61.7-91.9zm-56.6-164.2c27.3-32.4 24.8-61.9 24-72.5-24.1 1.4-52 16.4-67.9 34.9-17.5 19.8-27.8 44.3-25.6 71.9 26.1 2 49.9-11.4 69.5-34.3z"/>
      </svg>
      <span v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin shrink-0" />
      {{ loading ? 'Connexion...' : 'Continuer avec Apple' }}
    </button>
    <p v-if="error" class="text-red-600 text-xs mt-2 text-center">{{ error }}</p>
    <OAuthLinkModal
      v-model="showLinkModal"
      provider="apple"
      :email="linkData.email"
      :link-token="linkData.linkToken"
      @declined="onDeclined"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserAuthStore } from '@/stores/userAuth'
import OAuthLinkModal from '@/components/OAuthLinkModal.vue'

const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const loading = ref(false)
const error = ref('')
const showLinkModal = ref(false)
const linkData = ref({ email: '', linkToken: '' })

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

function onDeclined() {
  linkData.value = { email: '', linkToken: '' }
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

    if (res.status === 409) {
      if (json.status === 'email_conflict') {
        linkData.value = { email: json.email, linkToken: json.link_token }
        showLinkModal.value = true
        return
      }
      error.value = json.message || 'Un compte existe déjà avec cet email.'
      return
    }

    if (!res.ok) throw new Error(json.message || 'Erreur Apple.')

    userAuth.token = json.token
    userAuth.user = json.user
    localStorage.setItem('user_token', json.token)
    localStorage.setItem('user_data', JSON.stringify(json.user))

    const redirect = route.query.redirect
    router.push(redirect ? String(redirect) : '/')
  } catch (e) {
    if (e?.error === 'popup_closed_by_user' || e?.error === 'user_trigger_new_signin_flow') return
    error.value = e.message || 'Connexion Apple échouée.'
  } finally {
    loading.value = false
  }
}
</script>
