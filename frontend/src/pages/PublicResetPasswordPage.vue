<template>
  <div class="min-h-[70vh] flex items-center justify-center py-16 px-4">
    <div class="w-full max-w-md">

      <div class="text-center mb-8">
        <RouterLink to="/" class="inline-block mb-6">
          <img src="/logoentier.png" class="h-12 mx-auto" alt="UpcycleConnect" />
        </RouterLink>
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Nouveau mot de passe</h1>
        <p class="text-[#40617f] text-sm mt-2">Choisissez un mot de passe sécurisé pour votre compte.</p>
      </div>

      <div v-if="!token" class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] border border-[#edf4ff] p-8 text-center space-y-4">
        <div class="w-16 h-16 rounded-full flex items-center justify-center mx-auto" style="background:#fee2e2;">
          <ExclamationTriangleIcon class="w-8 h-8 text-red-500" />
        </div>
        <p class="text-gray-700 font-medium">Lien invalide ou expiré.</p>
        <RouterLink to="/connexion" class="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl font-bold text-[#006d35] border-2 border-[#006d35]/30 hover:bg-[#f0fdf4] transition text-sm">
          Retour à la connexion
        </RouterLink>
      </div>

      <div v-else-if="success" class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] border border-[#edf4ff] p-8 text-center space-y-6">
        <div class="w-16 h-16 rounded-full flex items-center justify-center mx-auto" style="background:#dcfce7;">
          <CheckCircleIcon class="w-8 h-8 text-[#006d35]" />
        </div>
        <div>
          <h2 class="font-bold text-[#001d32] text-xl">Mot de passe mis à jour !</h2>
          <p class="text-[#40617f] text-sm mt-2">Vous pouvez maintenant vous connecter avec votre nouveau mot de passe.</p>
        </div>
        <RouterLink
          to="/connexion"
          class="inline-flex items-center gap-2 px-8 py-3.5 rounded-xl font-jakarta font-bold text-white text-base transition hover:opacity-90"
          style="background: linear-gradient(135deg, #006d35, #1b8848);"
        >
          <ArrowRightCircleIcon class="w-5 h-5" />
          Se connecter
        </RouterLink>
      </div>

      <div v-else class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] border border-[#edf4ff] p-8 space-y-5">

        <div>
          <label class="block text-xs font-semibold text-[#001d32] mb-1.5 uppercase tracking-wide">Nouveau mot de passe</label>
          <div class="relative">
            <input
              v-model="password"
              :type="showPwd ? 'text' : 'password'"
              placeholder="8 caractères minimum"
              class="w-full px-4 py-3 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/20 focus:border-[#006d35] transition bg-white pr-11"
            />
            <button type="button" @click="showPwd = !showPwd" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition">
              <EyeSlashIcon v-if="showPwd" class="w-5 h-5" />
              <EyeIcon v-else class="w-5 h-5" />
            </button>
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-[#001d32] mb-1.5 uppercase tracking-wide">Confirmer le mot de passe</label>
          <input
            v-model="confirm"
            type="password"
            placeholder="Répétez le mot de passe"
            class="w-full px-4 py-3 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/20 focus:border-[#006d35] transition bg-white"
          />
        </div>

        <Transition name="fade">
          <p v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded-xl border border-red-200 flex items-start gap-2">
            <ExclamationTriangleIcon class="w-4 h-4 shrink-0 mt-0.5" />
            <span>{{ error }}</span>
          </p>
        </Transition>

        <button
          @click="submit"
          :disabled="loading || !password || !confirm"
          class="w-full py-3.5 rounded-xl font-jakarta font-bold text-white text-base transition-all hover:opacity-90 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          style="background: linear-gradient(135deg, #006d35, #1b8848);"
        >
          <div v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
          <ArrowRightCircleIcon v-else class="w-4 h-4" />
          {{ loading ? 'Enregistrement…' : 'Enregistrer le mot de passe' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon, CheckCircleIcon, ArrowRightCircleIcon } from '@heroicons/vue/24/outline'
import api from '@/services/api'

const route = useRoute()

const token   = ref('')
const password = ref('')
const confirm  = ref('')
const showPwd  = ref(false)
const loading  = ref(false)
const error    = ref('')
const success  = ref(false)

onMounted(() => {
  token.value = route.query.token || ''
})

async function submit() {
  error.value = ''
  if (password.value.length < 8) {
    error.value = 'Le mot de passe doit contenir au moins 8 caractères.'
    return
  }
  if (password.value !== confirm.value) {
    error.value = 'Les mots de passe ne correspondent pas.'
    return
  }
  loading.value = true
  try {
    await api.post('/auth/reset-password', { token: token.value, password: password.value })
    success.value = true
  } catch (e) {
    error.value = e?.response?.data?.message || 'Lien invalide ou expiré.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
