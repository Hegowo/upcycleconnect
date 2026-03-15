<template>
  <div class="min-h-[70vh] flex items-center justify-center py-12 px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <img src="/logoentier.png" class="h-16 mx-auto mb-4" alt="UpcycleConnect" />
        <h1 class="text-2xl font-bold" style="color: #103652;">Connexion</h1>
        <p class="text-gray-500 text-sm mt-1">Accédez à votre espace personnel</p>
      </div>

      <form @submit.prevent="handleLogin" class="bg-white rounded-2xl shadow-sm border p-8 space-y-5">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            v-model="form.email"
            type="email"
            required
            placeholder="votre@email.fr"
            class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2"
            style="border-color: #d1d5db; focus-ring-color: #1B8848;"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Mot de passe</label>
          <input
            v-model="form.password"
            type="password"
            required
            placeholder="••••••••"
            class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2"
          />
        </div>

        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2.5 rounded-lg text-white font-semibold transition"
          style="background-color: #1B8848;"
        >
          {{ loading ? 'Connexion...' : 'Se connecter' }}
        </button>
      </form>

      <p class="text-center text-sm text-gray-500 mt-4">
        Pas encore de compte ?
        <RouterLink to="/inscription" class="font-medium" style="color: #1B8848;">S'inscrire</RouterLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const form = ref({ email: '', password: '' })
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  loading.value = false
}
</script>
