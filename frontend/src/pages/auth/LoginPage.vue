<template>
  <AuthLayout>
    <div class="card p-8">
      <h2 class="text-xl font-bold text-gray-900 mb-6 text-center">Connexion</h2>

      <div v-if="error" class="mb-4 p-3 bg-red-50 border border-red-200 rounded-md text-sm text-red-700">
        {{ error }}
      </div>

      <form @submit.prevent="handleLogin" class="space-y-4">
        <div>
          <label class="label">Adresse email</label>
          <input
            v-model="form.email"
            type="email"
            class="input"
            placeholder="admin@upcycleconnect.fr"
            autocomplete="email"
            required
          />
        </div>
        <div>
          <label class="label">Mot de passe</label>
          <input
            v-model="form.password"
            type="password"
            class="input"
            placeholder="••••••••"
            autocomplete="current-password"
            required
          />
        </div>
        <button type="submit" :disabled="loading" class="btn-primary w-full justify-center">
          <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>
      </form>

      <p class="mt-4 text-center text-xs text-gray-400">
        Accès réservé aux administrateurs UpcycleConnect
      </p>
    </div>
  </AuthLayout>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AuthLayout from '@/layouts/AuthLayout.vue'

const router  = useRouter()
const auth    = useAuthStore()
const loading = ref(false)
const error   = ref('')

const form = reactive({ email: '', password: '' })

async function handleLogin() {
  error.value   = ''
  loading.value = true
  try {
    await auth.login(form.email, form.password)
    router.push('/admin/dashboard')
  } catch (err) {
    const status = err.response?.status
    if (status === 401) error.value = 'Identifiants invalides.'
    else if (status === 403) error.value = err.response?.data?.message || 'Accès refusé.'
    else if (status === 422) error.value = Object.values(err.response?.data?.errors || {}).flat().join(' ')
    else error.value = 'Une erreur est survenue. Veuillez réessayer.'
  } finally {
    loading.value = false
  }
}
</script>
