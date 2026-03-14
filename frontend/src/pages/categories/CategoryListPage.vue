<template>
  <div class="space-y-4">
    <div class="flex justify-between items-center">
      <h2 class="text-base font-medium text-gray-700">{{ categories.length }} catégorie(s)</h2>
      <RouterLink to="/admin/categories/create" class="btn-primary text-sm px-4 py-2 rounded-md">
        + Nouvelle catégorie
      </RouterLink>
    </div>

    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Nom</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Slug</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Prestations</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Ordre</th>
              <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actif</th>
              <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="loading">
              <td colspan="6" class="px-4 py-8 text-center">
                <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-green-600 mx-auto"></div>
              </td>
            </tr>
            <tr v-else-if="!categories.length">
              <td colspan="6" class="px-4 py-8 text-center text-gray-400">Aucune catégorie.</td>
            </tr>
            <tr v-for="cat in categories" :key="cat.id" class="hover:bg-gray-50">
              <td class="px-4 py-3 text-sm font-medium text-gray-900">{{ cat.name }}</td>
              <td class="px-4 py-3 text-sm text-gray-400 font-mono">{{ cat.slug }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ cat.prestations_count ?? 0 }}</td>
              <td class="px-4 py-3 text-sm text-gray-500">{{ cat.sort_order }}</td>
              <td class="px-4 py-3">
                <button
                  @click="toggleCategory(cat)"
                  :class="['relative inline-flex h-5 w-9 items-center rounded-full transition-colors', cat.is_active ? 'bg-green-600' : 'bg-gray-300']"
                >
                  <span :class="['inline-block h-3 w-3 transform rounded-full bg-white transition-transform', cat.is_active ? 'translate-x-5' : 'translate-x-1']" />
                </button>
              </td>
              <td class="px-4 py-3 text-right">
                <div class="flex justify-end gap-3">
                  <RouterLink :to="`/admin/categories/${cat.id}/edit`" class="text-xs text-blue-600 hover:text-blue-800 font-medium">
                    Modifier
                  </RouterLink>
                  <button @click="openDelete(cat)" class="text-xs text-red-600 hover:text-red-800 font-medium">
                    Supprimer
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
      title="Supprimer la catégorie"
      :message="`Êtes-vous sûr de vouloir supprimer la catégorie &quot;${deleteConfirm.category?.name}&quot; ?`"
      confirm-label="Supprimer"
      confirm-variant="danger"
      :loading="deleteConfirm.loading"
      @confirm="executeDelete"
      @cancel="deleteConfirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { categoryService } from '@/services/categoryService'
import { useToast } from '@/composables/useToast'
import AppConfirmDialog from '@/components/common/AppConfirmDialog.vue'

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
    toast.showError('Impossible de charger les catégories.')
  } finally {
    loading.value = false
  }
}

async function toggleCategory(cat) {
  try {
    await categoryService.toggle(cat.id)
    cat.is_active = !cat.is_active
    toast.showSuccess(`Catégorie ${cat.is_active ? 'activée' : 'désactivée'}.`)
  } catch {
    toast.showError('Impossible de modifier le statut.')
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
    toast.showSuccess('Catégorie supprimée.')
    deleteConfirm.show = false
  } catch {
    toast.showError('Impossible de supprimer la catégorie.')
  } finally {
    deleteConfirm.loading = false
  }
}

onMounted(fetchCategories)
</script>
