<template>
  <Teleport to="body">
    <Transition name="modal-fade">
      <div v-if="modelValue" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="decline" />
        <div class="relative bg-white rounded-2xl shadow-2xl p-6 w-full max-w-sm space-y-4 border border-[#edf4ff]">

          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 rounded-full flex items-center justify-center shrink-0"
              :class="provider === 'apple' ? 'bg-black' : 'bg-white border border-[#e2e8f0]'"
            >
              <svg v-if="provider === 'apple'" class="w-5 h-5 text-white" viewBox="0 0 384 512" fill="currentColor">
                <path d="M318.7 268.7c-.2-36.7 16.4-64.4 50-84.8-18.8-26.9-47.2-41.7-84.7-44.6-35.5-2.8-74.3 20.7-88.5 20.7-15 0-49.4-19.7-76.4-19.7C63.3 141.2 4 184.8 4 273.5q0 39.3 14.4 81.2c12.8 36.7 59 126.7 107.2 125.2 25.2-.6 43-17.9 75.8-17.9 31.8 0 48.3 17.9 76.4 17.9 48.6-.7 90.4-82.5 102.6-119.3-65.2-30.7-61.7-90-61.7-91.9zm-56.6-164.2c27.3-32.4 24.8-61.9 24-72.5-24.1 1.4-52 16.4-67.9 34.9-17.5 19.8-27.8 44.3-25.6 71.9 26.1 2 49.9-11.4 69.5-34.3z"/>
              </svg>
              <svg v-else class="w-5 h-5" viewBox="0 0 24 24">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
              </svg>
            </div>
            <div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-base leading-tight">Compte existant détecté</h2>
              <p class="text-xs text-[#40617f] mt-0.5 truncate max-w-[220px]">{{ email }}</p>
            </div>
          </div>

          <p class="text-sm text-[#40617f] leading-relaxed">
            Un compte existe déjà avec cette adresse. Voulez-vous lier votre compte
            <span class="font-semibold text-[#001d32]">{{ providerLabel }}</span>
            à votre compte existant ?
          </p>

          <div>
            <label class="block text-xs font-semibold text-[#001d32] mb-1.5 uppercase tracking-wide">Mot de passe actuel</label>
            <div class="relative">
              <input
                v-model="password"
                :type="showPwd ? 'text' : 'password'"
                placeholder="Confirmez votre identité"
                class="w-full px-4 py-3 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/20 focus:border-[#006d35] transition bg-white pr-11"
                @keyup.enter="confirm"
              />
              <button type="button" @click="showPwd = !showPwd" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition">
                <EyeSlashIcon v-if="showPwd" class="w-4 h-4" />
                <EyeIcon v-else class="w-4 h-4" />
              </button>
            </div>
          </div>

          <Transition name="fade">
            <p v-if="error" class="text-red-600 text-xs bg-red-50 p-2.5 rounded-lg border border-red-200 flex items-center gap-1.5">
              <ExclamationTriangleIcon class="w-3.5 h-3.5 shrink-0" />
              {{ error }}
            </p>
          </Transition>

          <div class="flex gap-3 pt-1">
            <button
              @click="decline"
              class="flex-1 py-2.5 rounded-xl text-sm font-semibold text-[#40617f] border border-[#e2e8f0] hover:bg-gray-50 transition"
            >
              Annuler
            </button>
            <button
              @click="confirm"
              :disabled="linking || !password"
              class="flex-1 py-2.5 rounded-xl text-sm font-bold text-white transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-1.5"
              style="background: linear-gradient(135deg, #006d35, #1b8848);"
            >
              <div v-if="linking" class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin" />
              {{ linking ? '...' : 'Lier mon compte' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon } from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'

const props = defineProps({
  modelValue: Boolean,
  provider: String,
  email: String,
  linkToken: String,
})

const emit = defineEmits(['update:modelValue', 'declined'])

const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const password = ref('')
const showPwd = ref(false)
const linking = ref(false)
const error = ref('')

const providerLabel = computed(() => props.provider === 'apple' ? 'Apple' : 'Google')

function decline() {
  password.value = ''
  error.value = ''
  emit('update:modelValue', false)
  emit('declined')
}

async function confirm() {
  if (!password.value || linking.value) return
  linking.value = true
  error.value = ''
  try {
    const res = await fetch('/api/v1/auth/oauth/link', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ link_token: props.linkToken, password: password.value }),
    })
    const json = await res.json()
    if (!res.ok) {
      error.value = json.message || 'Erreur lors de la liaison.'
      return
    }
    userAuth.token = json.token
    userAuth.user = json.user
    localStorage.setItem('user_token', json.token)
    localStorage.setItem('user_data', JSON.stringify(json.user))
    emit('update:modelValue', false)
    const redirect = route.query.redirect
    router.push(redirect ? String(redirect) : '/')
  } catch {
    error.value = 'Erreur de connexion. Veuillez réessayer.'
  } finally {
    linking.value = false
  }
}
</script>

<style scoped>
.modal-fade-enter-active, .modal-fade-leave-active { transition: opacity 0.2s ease; }
.modal-fade-enter-from, .modal-fade-leave-to { opacity: 0; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.15s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
