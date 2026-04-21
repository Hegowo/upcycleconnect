<template>
  <div class="min-h-[70vh] flex items-center justify-center py-16 px-4">
    <div class="w-full max-w-md">

      
      <div class="text-center mb-8">
        <RouterLink to="/" class="inline-block mb-6">
          <img src="/logoentier.png" class="h-12 mx-auto" alt="UpcycleConnect" />
        </RouterLink>
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Bon retour !</h1>
        <p class="text-[#40617f] text-sm mt-2">Connectez-vous à votre espace personnel</p>
      </div>

      
      <div class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] border border-[#edf4ff] p-8 space-y-5">

        <div>
          <label class="block text-xs font-semibold text-[#001d32] mb-1.5 uppercase tracking-wide">Email</label>
          <input
            v-model="form.email"
            type="email"
            required
            placeholder="votre@email.fr"
            class="w-full px-4 py-3 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/20 focus:border-[#006d35] transition bg-white"
          />
        </div>

        <div>
          <div class="flex items-center justify-between mb-1.5">
            <label class="block text-xs font-semibold text-[#001d32] uppercase tracking-wide">Mot de passe</label>
            <a href="#" class="text-xs text-[#006d35] hover:underline underline-offset-2">Mot de passe oublié ?</a>
          </div>
          <div class="relative">
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              placeholder="••••••••"
              class="w-full px-4 py-3 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/20 focus:border-[#006d35] transition bg-white pr-11"
            />
            <button type="button" @click="showPassword = !showPassword" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition">
              <EyeSlashIcon v-if="showPassword" class="w-5 h-5" />
              <EyeIcon v-else class="w-5 h-5" />
            </button>
          </div>
        </div>

        
        <Transition name="fade">
          <p v-if="verified" class="text-[#006d35] text-sm bg-[#f0fdf4] p-3 rounded-xl border border-[#bbf7d0] flex items-start gap-2">
            <CheckCircleIcon class="w-4 h-4 shrink-0 mt-0.5" />
            <span>Votre email a bien été vérifié. Vous pouvez vous connecter.</span>
          </p>
        </Transition>

        
        <Transition name="fade">
          <div v-if="pendingVerification" class="text-[#854d0e] text-sm bg-[#fef9c3] p-4 rounded-xl border border-[#fde047] space-y-1">
            <p class="font-semibold flex items-center gap-2">
              <ShieldCheckIcon class="w-4 h-4 shrink-0" />
              Vérification requise
            </p>
            <p>Nous avons détecté une connexion depuis un nouvel emplacement. Un email de confirmation vient d'être envoyé à votre adresse.</p>
            <p class="text-xs opacity-70">Le lien est valable 15 minutes.</p>
          </div>
        </Transition>

        
        <Transition name="fade">
          <p v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded-xl border border-red-200 flex items-start gap-2">
            <ExclamationTriangleIcon class="w-4 h-4 shrink-0 mt-0.5" />
            <span>{{ error }}</span>
          </p>
        </Transition>

        <button
          @click.prevent="handleLogin"
          :disabled="loading || !form.email || !form.password"
          class="w-full py-3.5 rounded-xl font-jakarta font-bold text-white text-base transition-all duration-200 hover:opacity-90 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          style="background: linear-gradient(135deg, #006d35, #1b8848);"
        >
          <div v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
          <ArrowRightCircleIcon v-else class="w-4 h-4" />
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>

        
        <div class="flex items-center gap-3">
          <div class="flex-1 h-px bg-[#e2e8f0]"></div>
          <span class="text-xs text-gray-400">ou</span>
          <div class="flex-1 h-px bg-[#e2e8f0]"></div>
        </div>

        <RouterLink
          to="/inscription"
          class="w-full py-3 rounded-xl font-semibold text-[#006d35] text-sm border-2 border-[#006d35]/20 hover:bg-[#f0fdf4] transition flex items-center justify-center gap-2"
        >
          Créer un compte gratuitement
        </RouterLink>
      </div>

      <p class="text-center text-xs text-gray-400 mt-6">
        En vous connectant, vous acceptez nos
        <a href="#" class="underline underline-offset-2 hover:text-[#006d35]">CGU</a>
        et notre
        <a href="#" class="underline underline-offset-2 hover:text-[#006d35]">politique de confidentialité</a>.
      </p>
    </div>
  </div>
</template>

<script setup>import { ref, computed } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { EyeIcon, EyeSlashIcon, ExclamationTriangleIcon, ArrowRightCircleIcon, CheckCircleIcon, ShieldCheckIcon } from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'

const router   = useRouter()
const route    = useRoute()
const userAuth = useUserAuthStore()

const form               = ref({ email: '', password: '' })
const loading            = ref(false)
const error              = ref('')
const showPassword       = ref(false)
const pendingVerification = ref(false)

const verified = computed(() => route.query.verified === '1')

async function handleLogin() {
  loading.value             = true
  error.value               = ''
  pendingVerification.value = false
  try {
    await userAuth.login(form.value.email, form.value.password)
    const redirect = route.query.redirect
    router.push(redirect ? String(redirect) : '/')
  } catch (e) {
    if (e.status === 202 || e.pendingVerification) {
      router.push({
        path: '/connexion/verification-requise',
        query: { email: form.value.email },
      })
      return
    }
    error.value = e.message || 'Identifiants invalides.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
