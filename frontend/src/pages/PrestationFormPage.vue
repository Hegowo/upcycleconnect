<template>
  <div class="max-w-2xl space-y-6">
    <RouterLink to="/admin/prestations" class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      {{ t('common.back') }}
    </RouterLink>

    <div v-if="loadingData" class="card p-8 space-y-4">
      <div class="h-6 bg-gray-100 rounded animate-pulse w-48 mb-6"></div>
      <div v-for="n in 5" :key="n" class="h-10 bg-gray-100 rounded animate-pulse"></div>
    </div>

    <div v-else class="card p-8 space-y-5">
      <h2 class="text-lg font-bold text-gray-900 mb-6">
        {{ isEditing ? t('prestations.editTitle') : t('prestations.createTitle') }}
      </h2>

      <div v-if="errors.length" class="p-3 bg-red-50 border border-red-200 rounded-xl">
        <p v-for="err in errors" :key="err" class="text-sm text-red-700">{{ err }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div>
          <label class="label">{{ t('prestations.fieldTitle') }} <span class="text-red-500">*</span></label>
          <input v-model="form.title" type="text" class="input" required />
        </div>

        <div>
          <label class="label">{{ t('prestations.fieldDescription') }} <span class="text-red-500">*</span></label>
          <textarea v-model="form.description" class="input min-h-24" rows="4" required />
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('prestations.fieldCategory') }}</label>
            <select v-model="form.category_id" class="input">
              <option :value="null">{{ t('prestations.selectCategory') }}</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>
          <div>
            <label class="label">{{ t('prestations.fieldProvider') }} <span class="text-red-500">*</span></label>
            <select v-model="form.provider_id" class="input" required>
              <option :value="null" disabled>{{ t('prestations.selectProvider') }}</option>
              <option v-for="p in providers" :key="p.id" :value="p.id">
                {{ p.profile?.company_name || `${p.first_name} ${p.last_name}` }}
              </option>
            </select>
          </div>
        </div>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('prestations.fieldPrice') }}</label>
            <input v-model.number="form.price" type="number" step="0.01" min="0" class="input" placeholder="0.00" />
          </div>
          <div>
            <label class="label">{{ t('prestations.fieldPriceType') }} <span class="text-red-500">*</span></label>
            <select v-model="form.price_type" class="input" required>
              <option value="fixed">{{ t('prestations.priceTypeFixed') }}</option>
              <option value="hourly">{{ t('prestations.priceTypeHourly') }}</option>
              <option value="quote">{{ t('prestations.priceTypeQuote') }}</option>
            </select>
          </div>
        </div>

        <div>
          <label class="label">{{ t('prestations.fieldStatus') }} <span class="text-red-500">*</span></label>
          <select v-model="form.status" class="input" required>
            <option value="draft">Brouillon</option>
            <option value="published">Publié</option>
            <option value="suspended">Suspendu</option>
            <option value="archived">Archivé</option>
          </select>
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
          <RouterLink to="/admin/prestations" class="inline-flex items-center px-4 py-2 rounded-xl text-sm font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 transition">
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
import { prestationService } from '@/services/prestationService'
import { categoryService } from '@/services/categoryService'
import { providerService } from '@/services/userService'
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
      toast.showSuccess(t('prestations.toastUpdated'))
    } else {
      await prestationService.create(payload)
      toast.showSuccess(t('prestations.toastCreated'))
    }
    router.push('/admin/prestations')
  } catch (err) {
    errors.value = err.response?.status === 422
      ? Object.values(err.response.data.errors || {}).flat()
      : [t('prestations.toastSaveError')]
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
