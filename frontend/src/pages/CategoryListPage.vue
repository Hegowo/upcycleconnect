<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('categories.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ t('categories.subtitle') }}</p>
      </div>
      <RouterLink to="/admin/categories/create" class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90" style="background-color:#1B8848;">
        {{ t('categories.createBtn') }}
      </RouterLink>
    </div>

    <div class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100 flex items-center justify-between">
        <span class="font-medium text-sm text-gray-700">{{ t('categories.title') }}</span>
        <span class="text-xs px-2 py-0.5 rounded-full bg-blue-100 text-blue-700 font-semibold">{{ categories.length }} {{ t('common.total') }}</span>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('categories.colName') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prestations</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('categories.colOrder') }}</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">{{ t('categories.colStatus') }}</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">{{ t('categories.colActions') }}</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading" v-for="n in 5" :key="n">
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-28"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-24"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-10"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-10"></div></td>
              <td class="px-4 py-3"><div class="h-5 w-9 bg-gray-100 rounded-full animate-pulse"></div></td>
              <td class="px-4 py-3"><div class="h-4 bg-gray-100 rounded animate-pulse w-20 ml-auto"></div></td>
            </tr>
            <tr v-else-if="!categories.length">
              <td :colspan="6" class="px-4 py-16 text-center">
                <div class="text-4xl mb-2">🔍</div>
                <p class="text-gray-500 font-medium">{{ t('common.noResults') }}</p>
              </td>
            </tr>
            <tr v-else v-for="cat in categories" :key="cat.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ cat.name }}<span class="text-gray-400 font-mono text-xs ml-2">{{ cat.slug }}</span></td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ cat.prestations_count ?? 0 }}</td>
              <td class="px-4 py-3">
                <span class="text-xs px-2 py-0.5 rounded-full bg-gray-100 text-gray-600 font-medium">#{{ cat.sort_order }}</span>
              </td>
              <td class="px-4 py-3">
                <button
                  @click="toggleCategory(cat)"
                  :class="['relative inline-flex h-5 w-9 items-center rounded-full transition-colors', cat.is_active ? 'bg-green-500' : 'bg-gray-300']"
                >
                  <span :class="['inline-block h-3 w-3 transform rounded-full bg-white transition-transform shadow', cat.is_active ? 'translate-x-5' : 'translate-x-1']" />
                </button>
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-2">
                  <RouterLink :to="`/admin/categories/${cat.id}/edit`" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#DBEAFE; color:#1d4ed8;">
                    {{ t('categories.actionEdit') }}
                  </RouterLink>
                  <button @click="openDelete(cat)" class="text-xs px-3 py-1 rounded-full font-medium transition" style="background:#FEE2E2; color:#dc2626;">
                    {{ t('categories.actionDelete') }}
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      :title="t('categories.confirmDeleteTitle')"
      :message="t('categories.confirmDeleteMsg', { name: deleteConfirm.category?.name })"
      :confirm-label="t('common.delete')"
      confirm-variant="danger"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/utils/useToast'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'

const { t } = useI18n()
const toast      = useToast()
const categories = ref([])
const loading    = ref(false)
const deleteConfirm = reactive({ show: false, category: null, loading: false })

async function fetchCategories() {
  loading.value = true
  try {
    const data = await categoryService.list()
    categories.value = data.data
  } catch {
    toast.showError(t('categories.toastLoadError'))
  } finally {
    loading.value = false
  }
}

async function toggleCategory(cat) {
  try {
    await categoryService.toggle(cat.id)
    cat.is_active = !cat.is_active
    toast.showSuccess(t('categories.toastToggled'))
  } catch {
    toast.showError(t('categories.toastLoadError'))
  }
}

function openDelete(cat) {
  deleteConfirm.category = cat
  deleteConfirm.show     = true
}

async function executeDelete() {
  deleteConfirm.loading = true
  try {
    await categoryService.remove(deleteConfirm.category.id)
    categories.value = categories.value.filter((c) => c.id !== deleteConfirm.category.id)
    toast.showSuccess(t('categories.toastDeleted'))
    deleteConfirm.show = false
  } catch {
    toast.showError(t('categories.toastDeleteError'))
  } finally {
    deleteConfirm.loading = false
  }
}

onMounted(fetchCategories)
</script>
