<template>
  <div class="max-w-2xl space-y-6">
    <RouterLink to="/admin/admins" class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      {{ t('common.back') }}
    </RouterLink>

    <div v-if="loading && isEditing" class="card p-8 space-y-4">
      <div class="h-6 bg-gray-100 rounded animate-pulse w-48 mb-6"></div>
      <div v-for="n in 5" :key="n" class="h-10 bg-gray-100 rounded animate-pulse"></div>
    </div>

    <div v-else class="card p-8 space-y-5">
      <h2 class="text-lg font-bold text-gray-900 mb-6">
        {{ isEditing ? t('admins.editTitle') : t('admins.createTitle') }}
      </h2>

      <div v-if="errors.length" class="p-3 bg-red-50 border border-red-200 rounded-xl">
        <p v-for="err in errors" :key="err" class="text-sm text-red-700">{{ err }}</p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="label">{{ t('admins.fieldFirstName') }} <span class="text-red-500">*</span></label>
            <input v-model="form.first_name" type="text" class="input" required />
          </div>
          <div>
            <label class="label">{{ t('admins.fieldLastName') }} <span class="text-red-500">*</span></label>
            <input v-model="form.last_name" type="text" class="input" required />
          </div>
        </div>

        <div>
          <label class="label">{{ t('admins.fieldEmail') }} <span class="text-red-500">*</span></label>
          <input v-model="form.email" type="email" class="input" :disabled="isEditing" required />
        </div>

        <div>
          <label class="label">
            {{ t('admins.fieldPassword') }}
            {{ isEditing ? `(${t('admins.fieldPasswordHint')})` : '' }}
            <span v-if="!isEditing" class="text-red-500">*</span>
          </label>
          <input v-model="form.password" type="password" class="input" :required="!isEditing" minlength="8" />
          <p v-if="!isEditing" class="text-xs text-gray-400 mt-1">{{ t('admins.fieldPasswordNew') }}</p>
        </div>

        <div>
          <label class="label">{{ t('admins.fieldRole') }} <span class="text-red-500">*</span></label>
          <select v-model="form.role" class="input" required>
            <option value="admin">{{ t('admins.roleAdmin') }}</option>
            <option value="super_admin">{{ t('admins.roleSuperAdmin') }}</option>
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
          <RouterLink to="/admin/admins" class="inline-flex items-center px-4 py-2 rounded-xl text-sm font-semibold text-gray-600 bg-gray-100 hover:bg-gray-200 transition">
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
import { adminUserService } from '@/services/userService'
import { useToast } from '@/utils/useToast'

const { t } = useI18n()
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
    errors.value = [t('admins.toastLoadError')]
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
      toast.showSuccess(t('admins.toastUpdated'))
    } else {
      await adminUserService.create(payload)
      toast.showSuccess(t('admins.toastCreated'))
    }
    router.push('/admin/admins')
  } catch (err) {
    errors.value = err.response?.status === 422
      ? Object.values(err.response.data.errors || {}).flat()
      : [t('admins.toastSaveError')]
  } finally {
    loading.value = false
  }
}
</script>
