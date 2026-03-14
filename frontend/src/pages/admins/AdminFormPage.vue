<template>
  <div class="max-w-md space-y-4">
    <div class="flex items-center gap-3">
      <RouterLink to="/admin/admins" class="text-sm text-gray-500 hover:text-gray-700">← Retour</RouterLink>
      <h2 class="text-base font-semibold text-gray-900">
        {{ isEditing ? 'Modifier l\'administrateur' : 'Nouvel administrateur' }}
      </h2>
    </div>

    <form @submit.prevent="handleSubmit" class="card p-6 space-y-4">
      <div v-if="errors.length" class="p-3 bg-red-50 border border-red-200 rounded-md">
        <p v-for="err in errors" :key="err" class="text-sm text-red-700">{{ err }}</p>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Prénom <span class="text-red-500">*</span></label>
          <input v-model="form.first_name" type="text" class="input" required />
        </div>
        <div>
          <label class="label">Nom <span class="text-red-500">*</span></label>
          <input v-model="form.last_name" type="text" class="input" required />
        </div>
      </div>

      <div>
        <label class="label">Email <span class="text-red-500">*</span></label>
        <input v-model="form.email" type="email" class="input" :disabled="isEditing" required />
      </div>

      <div>
        <label class="label">Mot de passe {{ isEditing ? '(laisser vide pour conserver)' : '' }} <span v-if="!isEditing" class="text-red-500">*</span></label>
        <input v-model="form.password" type="password" class="input" :required="!isEditing" minlength="8" />
      </div>

      <div>
        <label class="label">Rôle <span class="text-red-500">*</span></label>
        <select v-model="form.role" class="input" required>
          <option value="admin">Administrateur</option>
          <option value="super_admin">Super Administrateur</option>
        </select>
      </div>

      <div class="flex gap-3 pt-2 border-t border-gray-100">
        <button type="submit" :disabled="loading" class="btn-primary text-sm px-4 py-2 rounded-md">
          <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ loading ? 'Enregistrement...' : 'Enregistrer' }}
        </button>
        <RouterLink to="/admin/admins" class="btn-secondary text-sm px-4 py-2 rounded-md inline-flex items-center">
          Annuler
        </RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { adminUserService } from '@/services/userService'
import { useToast } from '@/composables/useToast'

const route  = useRoute()
const router = useRouter()
const toast  = useToast()

const isEditing = computed(() => !!route.params.id)
const loading   = ref(false)
const errors    = ref([])

const form = reactive({
  first_name: '',
  last_name:  '',
  email:      '',
  password:   '',
  role:       'admin',
})

onMounted(async () => {
  if (!isEditing.value) return
  try {
    loading.value = true
    const admin = await adminUserService.getById(route.params.id)
    form.first_name = admin.first_name
    form.last_name  = admin.last_name
    form.email      = admin.email
    form.role       = admin.role ?? 'admin'
  } catch {
    errors.value = ['Impossible de charger les données de l\'administrateur.']
  } finally {
    loading.value = false
  }
})

async function handleSubmit() {
  errors.value  = []
  loading.value = true
  try {
    const payload = { ...form }
    if (isEditing.value && !payload.password) delete payload.password

    if (isEditing.value) {
      await adminUserService.update(route.params.id, payload)
      toast.showSuccess('Administrateur mis à jour.')
    } else {
      await adminUserService.create(payload)
      toast.showSuccess('Administrateur créé.')
    }
    router.push('/admin/admins')
  } catch (err) {
    errors.value = err.response?.status === 422
      ? Object.values(err.response.data.errors || {}).flat()
      : ['Une erreur est survenue.']
  } finally {
    loading.value = false
  }
}
</script>
