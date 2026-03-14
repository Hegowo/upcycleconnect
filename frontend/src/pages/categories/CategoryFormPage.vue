<template>
  <div class="max-w-xl space-y-4">
    <div class="flex items-center gap-3">
      <RouterLink to="/admin/categories" class="text-sm text-gray-500 hover:text-gray-700">← Retour</RouterLink>
      <h2 class="text-base font-semibold text-gray-900">
        {{ isEditing ? 'Modifier la catégorie' : 'Nouvelle catégorie' }}
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
        <label class="label">Nom <span class="text-red-500">*</span></label>
        <input v-model="form.name" @input="autoSlug" type="text" class="input" placeholder="Ex: Mobilier" required />
      </div>

      <div>
        <label class="label">Slug</label>
        <input v-model="form.slug" type="text" class="input font-mono" placeholder="mobilier" />
        <p class="text-xs text-gray-400 mt-1">Généré automatiquement depuis le nom si laissé vide.</p>
      </div>

      <div>
        <label class="label">Description</label>
        <textarea v-model="form.description" class="input min-h-20" rows="3" placeholder="Description optionnelle..." />
      </div>

      <div class="flex items-center gap-6">
        <div>
          <label class="label">Ordre d'affichage</label>
          <input v-model.number="form.sort_order" type="number" class="input w-24" min="0" />
        </div>
        <div class="flex items-center gap-2 pt-5">
          <input v-model="form.is_active" id="is_active" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-green-600" />
          <label for="is_active" class="text-sm text-gray-700">Catégorie active</label>
        </div>
      </div>

      <div class="flex gap-3 pt-2 border-t border-gray-100">
        <button type="submit" :disabled="loading" class="btn-primary text-sm px-4 py-2 rounded-md">
          <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ loading ? 'Enregistrement...' : 'Enregistrer' }}
        </button>
        <RouterLink to="/admin/categories" class="btn-secondary text-sm px-4 py-2 rounded-md inline-flex items-center">
          Annuler
        </RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/composables/useToast'

const route  = useRoute()
const router = useRouter()
const toast  = useToast()

const isEditing  = computed(() => !!route.params.id)
const loading    = ref(false)
const loadingData = ref(false)
const errors     = ref([])

const form = reactive({
  name: '',
  slug: '',
  description: '',
  is_active: true,
  sort_order: 0,
})

function autoSlug() {
  form.slug = form.name
    .toLowerCase()
    .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '')
}

async function handleSubmit() {
  errors.value  = []
  loading.value = true

  try {
    if (isEditing.value) {
      await categoryService.update(route.params.id, form)
      toast.showSuccess('Catégorie mise à jour.')
    } else {
      await categoryService.create(form)
      toast.showSuccess('Catégorie créée.')
    }
    router.push('/admin/categories')
  } catch (err) {
    if (err.response?.status === 422) {
      errors.value = Object.values(err.response.data.errors || {}).flat()
    } else if (err.response?.status === 409) {
      errors.value = [err.response.data.message]
    } else {
      errors.value = ['Une erreur est survenue.']
    }
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  if (isEditing.value) {
    loadingData.value = true
    try {
      const data = await categoryService.get(route.params.id)
      Object.assign(form, data)
    } finally {
      loadingData.value = false
    }
  }
})
</script>
