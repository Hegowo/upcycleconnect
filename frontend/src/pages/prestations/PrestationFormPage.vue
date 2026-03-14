<template>
  <div class="max-w-xl space-y-4">
    <div class="flex items-center gap-3">
      <RouterLink to="/admin/prestations" class="text-sm text-gray-500 hover:text-gray-700">← Retour</RouterLink>
      <h2 class="text-base font-semibold text-gray-900">
        {{ isEditing ? 'Modifier la prestation' : 'Nouvelle prestation' }}
      </h2>
    </div>

    <div v-if="loadingData" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="card p-6 space-y-4">
      <div v-if="errors.length" class="p-3 bg-red-50 border border-red-200 rounded-md">
        <p v-for="err in errors" :key="err" class="text-sm text-red-700">{{ err }}</p>
      </div>

      <div>
        <label class="label">Titre <span class="text-red-500">*</span></label>
        <input v-model="form.title" type="text" class="input" required />
      </div>

      <div>
        <label class="label">Description <span class="text-red-500">*</span></label>
        <textarea v-model="form.description" class="input min-h-24" rows="4" required />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Catégorie</label>
          <select v-model="form.category_id" class="input">
            <option :value="null">Aucune</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
        <div>
          <label class="label">Prestataire <span class="text-red-500">*</span></label>
          <select v-model="form.provider_id" class="input" required>
            <option :value="null" disabled>Sélectionner...</option>
            <option v-for="p in providers" :key="p.id" :value="p.id">
              {{ p.profile?.company_name || `${p.first_name} ${p.last_name}` }}
            </option>
          </select>
        </div>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Prix (€)</label>
          <input v-model.number="form.price" type="number" step="0.01" min="0" class="input" placeholder="0.00" />
        </div>
        <div>
          <label class="label">Type de tarif <span class="text-red-500">*</span></label>
          <select v-model="form.price_type" class="input" required>
            <option value="fixed">Fixe</option>
            <option value="hourly">Horaire</option>
            <option value="quote">Sur devis</option>
          </select>
        </div>
      </div>

      <div>
        <label class="label">Statut <span class="text-red-500">*</span></label>
        <select v-model="form.status" class="input" required>
          <option value="draft">Brouillon</option>
          <option value="published">Publié</option>
          <option value="suspended">Suspendu</option>
          <option value="archived">Archivé</option>
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
        <RouterLink to="/admin/prestations" class="btn-secondary text-sm px-4 py-2 rounded-md inline-flex items-center">
          Annuler
        </RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { prestationService } from '@/services/prestationService'
import { categoryService } from '@/services/categoryService'
import { providerService } from '@/services/userService'
import { useToast } from '@/composables/useToast'

const route  = useRoute()
const router = useRouter()
const toast  = useToast()

const isEditing   = computed(() => !!route.params.id)
const loading     = ref(false)
const loadingData = ref(false)
const errors      = ref([])
const categories  = ref([])
const providers   = ref([])

const form = reactive({
  title: '',
  description: '',
  category_id: null,
  provider_id: null,
  price: '',
  price_type: 'fixed',
  status: 'draft',
})

async function handleSubmit() {
  errors.value  = []
  loading.value = true
  try {
    const payload = { ...form, price: form.price === '' ? null : form.price }
    if (isEditing.value) {
      await prestationService.update(route.params.id, payload)
      toast.showSuccess('Prestation mise à jour.')
    } else {
      await prestationService.create(payload)
      toast.showSuccess('Prestation créée.')
    }
    router.push('/admin/prestations')
  } catch (err) {
    errors.value = err.response?.status === 422
      ? Object.values(err.response.data.errors || {}).flat()
      : ['Une erreur est survenue.']
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const [catRes, provRes] = await Promise.all([
    categoryService.list(),
    providerService.list({ status: 'approved', per_page: 100 }),
  ])
  categories.value = catRes.data
  providers.value  = provRes.data

  if (isEditing.value) {
    loadingData.value = true
    try {
      const data = await prestationService.get(route.params.id)
      Object.assign(form, {
        title:       data.title,
        description: data.description,
        category_id: data.category?.id || null,
        provider_id: data.provider?.id || null,
        price:       data.price,
        price_type:  data.price_type,
        status:      data.status,
      })
    } finally {
      loadingData.value = false
    }
  }
})
</script>
