<template>
  <div class="min-h-[70vh] flex items-center justify-center py-12 px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <img src="/logoentier.png" class="h-16 mx-auto mb-4" alt="UpcycleConnect" />
        <h1 class="text-2xl font-bold" style="color: #103652;">Créer un compte</h1>
        <p class="text-gray-500 text-sm mt-1">Rejoignez la communauté UpcycleConnect</p>
      </div>

      <div class="flex gap-3 mb-6">
        <button
          @click="form.type = 'user'"
          :class="form.type === 'user' ? 'border-2 font-semibold' : 'border text-gray-500'"
          class="flex-1 py-3 rounded-xl text-sm transition"
          :style="form.type === 'user' ? 'border-color: #1B8848; color: #1B8848;' : ''"
        >
          👤 Particulier
        </button>
        <button
          @click="form.type = 'provider'"
          :class="form.type === 'provider' ? 'border-2 font-semibold' : 'border text-gray-500'"
          class="flex-1 py-3 rounded-xl text-sm transition"
          :style="form.type === 'provider' ? 'border-color: #1B8848; color: #1B8848;' : ''"
        >
          🔨 Prestataire
        </button>
      </div>

      <form @submit.prevent="handleRegister" class="bg-white rounded-2xl shadow-sm border p-8 space-y-4">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Prénom</label>
            <input v-model="form.firstName" type="text" required class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nom</label>
            <input v-model="form.lastName" type="text" required class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none" />
          </div>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input v-model="form.email" type="email" required class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Mot de passe</label>
          <input v-model="form.password" type="password" required minlength="8" class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none" />
          <p class="text-xs text-gray-400 mt-1">Minimum 8 caractères</p>
        </div>

        <template v-if="form.type === 'provider'">
          <hr class="border-gray-100" />
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nom de l'entreprise</label>
            <input v-model="form.companyName" type="text" required class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">SIRET</label>
            <input v-model="form.siret" type="text" maxlength="14" class="w-full px-4 py-2 border rounded-lg text-sm focus:outline-none" placeholder="14 chiffres" />
          </div>
        </template>

        <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2.5 rounded-lg text-white font-semibold transition"
          style="background-color: #1B8848;"
        >
          {{ loading ? 'Création...' : 'Créer mon compte' }}
        </button>
      </form>

      <p class="text-center text-sm text-gray-500 mt-4">
        Déjà inscrit ?
        <RouterLink to="/connexion" class="font-medium" style="color: #1B8848;">Se connecter</RouterLink>
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const form = ref({
  type: 'user',
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  companyName: '',
  siret: '',
})
const loading = ref(false)
const error = ref('')

async function handleRegister() {
  loading.value = true
  error.value = ''
  loading.value = false
}
</script>
