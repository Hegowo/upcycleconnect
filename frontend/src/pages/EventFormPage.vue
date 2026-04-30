<template>
  <div class="max-w-2xl space-y-6">
    <RouterLink to="/admin/events" class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      {{ t('common.back') }}
    </RouterLink>

    <div v-if="loadingData" class="card p-8 space-y-4">
      <div class="h-6 bg-gray-100 rounded animate-pulse w-48 mb-6"></div>
      <div v-for="n in 5" :key="n" class="h-10 bg-gray-100 rounded animate-pulse"></div>
    </div>

    <div v-else class="card p-8 space-y-5">
      <h2 class="text-lg font-bold text-gray-900 mb-6">
        {{ isEditing ? t('events.editTitle') : t('events.createTitle') }}
      </h2>

      <div v-if="errors.length" class="p-3 bg-red-50 border border-red-200 rounded-xl">
        <p v-for="err in errors" :key="err" class="text-sm text-red-700">{{ err }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div>
          <label class="label">{{ t('events.fieldTitle') }} <span class="text-red-500">*</span></label>
          <input v-model="form.title" type="text" class="input" required />
        </div>

        <div>
          <label class="label">{{ t('events.fieldDescription') }} <span class="text-red-500">*</span></label>
          <textarea v-model="form.description" class="input min-h-24" rows="4" required />
        </div>

        <div>
          <label class="label">{{ t('events.fieldLocation') }}</label>
          <input v-model="form.location" type="text" class="input" placeholder="Ex: Paris 11e — Salle communautaire" />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('events.fieldStartDate') }} <span class="text-red-500">*</span></label>
            <input v-model="form.start_date" type="datetime-local" class="input" required :min="isEditing ? undefined : todayMin" />
          </div>
          <div>
            <label class="label">{{ t('events.fieldEndDate') }} <span class="text-red-500">*</span></label>
            <input v-model="form.end_date" type="datetime-local" class="input" required :min="isEditing ? undefined : todayMin" />
          </div>
        </div>

        <div v-if="dateError" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">{{ dateError }}</div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('events.fieldMaxParticipants') }}</label>
            <input v-model.number="form.max_participants" type="number" min="1" class="input" placeholder="Illimité si vide" />
          </div>
          <div>
            <label class="label">{{ t('events.fieldCategory') }}</label>
            <select v-model="form.category_id" class="input">
              <option :value="null">{{ t('events.selectCategory') }}</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
        </div>

        <div>
          <label class="label">{{ t('events.fieldStatus') }} <span class="text-red-500">*</span></label>
          <select v-model="form.status" class="input" required>
            <option value="draft">Brouillon</option>
            <option value="published">Publié</option>
            <option value="cancelled">Annulé</option>
          </select>
        </div>

        <div class="flex gap-3 pt-4 border-t border-gray-100">
          <button
            type="submit"
            :disabled="loading || !!dateError"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90 disabled:opacity-60"
            style="background-color:#1B8848;"
          >
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            {{ loading ? t('common.saving') : t('common.save') }}
          </button>
          <RouterLink to="/admin/events" class="inline-flex items-center px-4 py-2 rounded-xl text-sm font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 transition">
            {{ t('common.cancel') }}
          </RouterLink>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { eventService } from '@/services/eventService'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/utils/useToast'

const { t } = useI18n()
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

const todayMin = computed(() => {
  const now = new Date()
  now.setSeconds(0, 0)
  return now.toISOString().slice(0, 16)
})

const dateError = computed(() => {
  if (!isEditing.value && form.start_date && form.start_date < todayMin.value) {
    return t('events.datePastError')
  }
  if (form.start_date && form.end_date && form.end_date <= form.start_date) {
    return t('events.dateError')
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
      toast.showSuccess(t('events.toastUpdated'))
    } else {
      await eventService.create(payload)
      toast.showSuccess(t('events.toastCreated'))
    }
    router.push('/admin/events')
  } catch (err) {
    errors.value = err.response?.status === 422
      ? Object.values(err.response.data.errors || {}).flat()
      : [t('events.toastSaveError')]
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
