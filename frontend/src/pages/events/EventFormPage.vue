<template>
  <div class="max-w-xl space-y-4">
    <div class="flex items-center gap-3">
      <RouterLink to="/admin/events" class="text-sm text-gray-500 hover:text-gray-700">← Retour</RouterLink>
      <h2 class="text-base font-semibold text-gray-900">
        {{ isEditing ? 'Modifier l\'événement' : 'Nouvel événement' }}
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

      <div>
        <label class="label">Lieu</label>
        <input v-model="form.location" type="text" class="input" placeholder="Ex: Paris 11e — Salle communautaire" />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Date et heure de début <span class="text-red-500">*</span></label>
          <input v-model="form.start_date" type="datetime-local" class="input" required />
        </div>
        <div>
          <label class="label">Date et heure de fin <span class="text-red-500">*</span></label>
          <input v-model="form.end_date" type="datetime-local" class="input" required />
        </div>
      </div>

      <div v-if="dateError" class="text-sm text-red-600">{{ dateError }}</div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">Participants max</label>
          <input v-model.number="form.max_participants" type="number" min="1" class="input" placeholder="Illimité si vide" />
        </div>
        <div>
          <label class="label">Catégorie</label>
          <select v-model="form.category_id" class="input">
            <option :value="null">Aucune</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
          </select>
        </div>
      </div>

      <div>
        <label class="label">Statut <span class="text-red-500">*</span></label>
        <select v-model="form.status" class="input" required>
          <option value="draft">Brouillon</option>
          <option value="published">Publié</option>
          <option value="cancelled">Annulé</option>
        </select>
      </div>

      <div class="flex gap-3 pt-2 border-t border-gray-100">
        <button type="submit" :disabled="loading || !!dateError" class="btn-primary text-sm px-4 py-2 rounded-md">
          <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
          </svg>
          {{ loading ? 'Enregistrement...' : 'Enregistrer' }}
        </button>
        <RouterLink to="/admin/events" class="btn-secondary text-sm px-4 py-2 rounded-md inline-flex items-center">
          Annuler
        </RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { eventService } from '@/services/eventService'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/composables/useToast'

const route  = useRoute()
const router = useRouter()
const toast  = useToast()

const isEditing   = computed(() => !!route.params.id)
const loading     = ref(false)
const loadingData = ref(false)
const errors      = ref([])
const categories  = ref([])

const form = reactive({
  title: '',
  description: '',
  location: '',
  start_date: '',
  end_date: '',
  max_participants: '',
  category_id: null,
  status: 'draft',
})

const dateError = computed(() => {
  if (form.start_date && form.end_date && form.end_date <= form.start_date) {
    return 'La date de fin doit être postérieure à la date de début.'
  }
  return ''
})

async function handleSubmit() {
  if (dateError.value) return
  errors.value  = []
  loading.value = true
  try {
    const payload = {
      ...form,
      max_participants: form.max_participants === '' ? null : form.max_participants,
    }
    if (isEditing.value) {
      await eventService.update(route.params.id, payload)
      toast.showSuccess('Événement mis à jour.')
    } else {
      await eventService.create(payload)
      toast.showSuccess('Événement créé.')
    }
    router.push('/admin/events')
  } catch (err) {
    errors.value = err.response?.status === 422
      ? Object.values(err.response.data.errors || {}).flat()
      : ['Une erreur est survenue.']
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const catData = await categoryService.list()
  categories.value = catData.data

  if (isEditing.value) {
    loadingData.value = true
    try {
      const data = await eventService.get(route.params.id)
      Object.assign(form, {
        title:            data.title,
        description:      data.description,
        location:         data.location || '',
        start_date:       data.start_date ? data.start_date.slice(0, 16) : '',
        end_date:         data.end_date ? data.end_date.slice(0, 16) : '',
        max_participants: data.max_participants ?? '',
        category_id:      data.category?.id || null,
        status:           data.status,
      })
    } finally {
      loadingData.value = false
    }
  }
})
</script>
