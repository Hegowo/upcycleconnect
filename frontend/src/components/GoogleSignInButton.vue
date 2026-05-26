<template>
  <div>
    <div ref="buttonEl" class="w-full" />
    <p v-if="error" class="text-red-600 text-xs mt-2 text-center">{{ error }}</p>
    <OAuthLinkModal
      v-model="showLinkModal"
      provider="google"
      :email="linkData.email"
      :link-token="linkData.linkToken"
      @declined="onDeclined"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserAuthStore } from '@/stores/userAuth'
import OAuthLinkModal from '@/components/OAuthLinkModal.vue'

const emit = defineEmits(['loading', 'done', 'error'])
const buttonEl = ref(null)
const error = ref('')
const showLinkModal = ref(false)
const linkData = ref({ email: '', linkToken: '' })
const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const CLIENT_ID = '14117163636-mpil6ctjnfpp5qh68n55inrksodv63a3.apps.googleusercontent.com'

function loadGIS() {
  return new Promise((resolve) => {
    if (window.google?.accounts?.id) { resolve(); return }
    const script = document.createElement('script')
    script.src = 'https://accounts.google.com/gsi/client'
    script.async = true
    script.defer = true
    script.onload = resolve
    document.head.appendChild(script)
  })
}

function onDeclined() {
  linkData.value = { email: '', linkToken: '' }
}

async function handleCredential(response) {
  error.value = ''
  emit('loading', true)
  try {
    const res = await fetch('/api/v1/auth/google', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ credential: response.credential }),
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

    if (!res.ok) throw new Error(json.message || 'Erreur Google.')

    userAuth.token = json.token
    userAuth.user = json.user
    localStorage.setItem('user_token', json.token)
    localStorage.setItem('user_data', JSON.stringify(json.user))

    emit('done')
    const redirect = route.query.redirect
    router.push(redirect ? String(redirect) : '/')
  } catch (e) {
    error.value = e.message
    emit('error', e.message)
  } finally {
    emit('loading', false)
  }
}

onMounted(async () => {
  await loadGIS()
  window.google.accounts.id.initialize({
    client_id: CLIENT_ID,
    callback: handleCredential,
    use_fedcm_for_prompt: true,
  })
  window.google.accounts.id.renderButton(buttonEl.value, {
    theme: 'outline',
    size: 'large',
    width: buttonEl.value?.offsetWidth || 400,
    text: 'continue_with',
    shape: 'rectangular',
    logo_alignment: 'center',
  })
})
</script>
