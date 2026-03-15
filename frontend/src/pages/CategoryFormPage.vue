<template>
  <div class="max-w-2xl space-y-6">
    <RouterLink to="/admin/categories" class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      {{ t('common.back') }}
    </RouterLink>

    <div v-if="loadingData" class="card p-8 space-y-4">
      <div class="h-6 bg-gray-100 rounded animate-pulse w-48 mb-6"></div>
      <div v-for="n in 4" :key="n" class="h-10 bg-gray-100 rounded animate-pulse"></div>
    </div>

    <div v-else class="card p-8 space-y-5">
      <h2 class="text-lg font-bold text-gray-900 mb-6">
        {{ isEditing ? t('categories.editTitle') : t('categories.createTitle') }}
      </h2>

      <div v-if="errors.length" class="p-3 bg-red-50 border border-red-200 rounded-xl">
        <p v-for="err in errors" :key="err" class="text-sm text-red-700">{{ err }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div>
          <label class="label">{{ t('categories.fieldName') }} <span class="text-red-500">*</span></label>
          <input v-model="form.name" @input="autoSlug" type="text" class="input" placeholder="Ex: Mobilier" required />
        </div>

        <div>
          <label class="label">{{ t('categories.fieldSlug') }}</label>
          <input v-model="form.slug" type="text" class="input font-mono" placeholder="mobilier" />
          <p class="text-xs text-gray-400 mt-1">{{ t('categories.fieldSlugHint') }}</p>
        </div>

        <div>
          <label class="label">{{ t('categories.fieldDescription') }}</label>
          <textarea v-model="form.description" class="input min-h-20" rows="3" placeholder="Description optionnelle..." />
        </div>

        <div class="flex items-center gap-6">
          <div>
            <label class="label">{{ t('categories.fieldOrder') }}</label>
            <input v-model.number="form.sort_order" type="number" class="input w-24" min="0" />
          </div>
          <div class="flex items-center gap-2 pt-5">
            <input v-model="form.is_active" id="is_active" type="checkbox" class="h-4 w-4 rounded border-gray-300 text-green-600" />
            <label for="is_active" class="text-sm text-gray-700">{{ t('categories.fieldActive') }}</label>
          </div>
        </div>

        <div class="flex gap-3 pt-4 border-t border-gray-100">
          <button
            type="submit"
            :disabled="loading"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90 disabled:opacity-60"
            style="background-color:#1B8848;"
          >
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            {{ loading ? t('common.saving') : t('common.save') }}
          </button>
          <RouterLink to="/admin/categories" class="inline-flex items-center px-4 py-2 rounded-xl text-sm font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 transition">
            {{ t('common.cancel') }}
          </RouterLink>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/utils/useToast'

const { t } = useI18n()
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
      toast.showSuccess(t('categories.toastUpdated'))
    } else {
      await categoryService.create(form)
      toast.showSuccess(t('categories.toastCreated'))
    }
    router.push('/admin/categories')
  } catch (err) {
    if (err.response?.status === 422) {
      errors.value = Object.values(err.response.data.errors || {}).flat()
    } else if (err.response?.status === 409) {
      errors.value = [err.response.data.message]
    } else {
      errors.value = [t('categories.toastSaveError')]
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
