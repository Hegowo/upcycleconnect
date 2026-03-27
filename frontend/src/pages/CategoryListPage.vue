<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between gap-3">
      <div class="min-w-0">
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32] truncate">{{ t('categories.objectTitle') }}</h2>
        <p class="text-sm text-[#40617f] mt-0.5 hidden sm:block">{{ t('categories.objectSubtitle') }}</p>
      </div>
      <RouterLink
        to="/admin/categories/create"
        class="p-2 sm:px-4 sm:py-2 text-sm font-semibold rounded-lg text-white transition hover:opacity-90 flex items-center gap-2 shrink-0"
        style="background-color:#006d35;"
      >
        <PlusIcon class="w-4 h-4" />
        <span class="hidden sm:inline">{{ t('categories.newCategoryBtn') }}</span>
      </RouterLink>
    </div>

    <!-- Loading grid -->
    <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="n in 8" :key="n" class="bg-white rounded-2xl border border-[#f1f5f9] h-48 animate-pulse"></div>
    </div>

    <!-- Empty -->
    <div v-else-if="!categories.length" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm py-16 text-center">
      <div class="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto mb-3">
        <TagIcon class="w-6 h-6 text-gray-400" />
      </div>
      <p class="text-gray-500 font-medium">{{ t('categories.noCategories') }}</p>
      <p class="text-gray-400 text-sm mt-1">{{ t('categories.createFirst') }}</p>
    </div>

    <!-- Categories grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div
        v-for="cat in categories"
        :key="cat.id"
        class="bg-white rounded-2xl border shadow-sm hover:shadow-md transition-shadow overflow-hidden"
        :class="cat.is_active ? 'border-[#f1f5f9]' : 'border-gray-200 opacity-75'"
      >
        <div class="p-5">
          <!-- Icon + toggle -->
          <div class="flex items-start justify-between mb-4">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center" :style="{ backgroundColor: catIconBg(cat) }">
              <component :is="catIcon(cat)" class="w-6 h-6" :style="{ color: catIconColor(cat) }" />
            </div>
            <button
              @click="toggleCategory(cat)"
              :class="['relative inline-flex h-5 w-9 items-center rounded-full transition-colors shrink-0', cat.is_active ? 'bg-[#006d35]' : 'bg-gray-300']"
              :title="cat.is_active ? t('categories.deactivate') : t('categories.activate')"
            >
              <span :class="['inline-block h-3 w-3 transform rounded-full bg-white transition-transform shadow', cat.is_active ? 'translate-x-5' : 'translate-x-1']" />
            </button>
          </div>

          <!-- Name + desc -->
          <h4 class="text-sm font-bold text-[#001d32] mb-1">{{ cat.name }}</h4>
          <p class="text-xs text-gray-400 mb-3 line-clamp-2">{{ cat.description || t('categories.noDescription') }}</p>

          <!-- Count -->
          <div class="flex items-center justify-between mb-4">
            <span class="text-xs font-semibold px-2 py-0.5 rounded-full" :class="cat.is_active ? 'bg-[#dcfce7] text-[#166534]' : 'bg-gray-100 text-gray-500'">
              {{ cat.is_active ? t('categories.active') : t('categories.inactive') }}
            </span>
            <span class="text-xs text-gray-400 font-medium">
              {{ cat.prestations_count ?? 0 }} {{ t('categories.articles') }}
            </span>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2 pt-3 border-t border-[#f8fafc]">
            <RouterLink
              :to="`/admin/categories/${cat.id}/edit`"
              class="flex-1 py-1.5 rounded-lg text-xs font-semibold text-center transition"
              style="background:#f1f5f9; color:#475569;"
            >
              {{ t('categories.actionEdit') }}
            </RouterLink>
            <button
              @click="openDelete(cat)"
              class="p-1.5 rounded-lg text-gray-400 hover:text-red-500 hover:bg-red-50 transition"
              :title="t('categories.actionDelete')"
            >
              <TrashIcon class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Summary bar -->
    <div v-if="categories.length" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-4">
      <div class="grid grid-cols-3 gap-3 sm:gap-6">
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">{{ t('categories.totalCategories') }}</p>
          <p class="text-lg font-bold text-[#001d32]">{{ categories.length }}</p>
        </div>
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">{{ t('categories.activeCount') }}</p>
          <p class="text-lg font-bold text-[#006d35]">{{ categories.filter(c => c.is_active).length }}</p>
        </div>
        <div class="text-center">
          <p class="text-xs text-gray-400 uppercase tracking-wider mb-1">{{ t('categories.inactiveCount') }}</p>
          <p class="text-lg font-bold text-gray-400">{{ categories.filter(c => !c.is_active).length }}</p>
        </div>
      </div>
    </div>

    <AppConfirmDialog
      :show="deleteConfirm.show"
      :title="t('categories.confirmDeleteCategoryTitle')"
      :message="t('categories.confirmDeleteCategoryMsg', { name: deleteConfirm.category?.name })"
      :confirm-label="t('categories.actionDelete')"
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
import {
  PlusIcon, TagIcon, TrashIcon,
  CubeIcon, SparklesIcon, WrenchScrewdriverIcon,
  GlobeAltIcon, HeartIcon, TruckIcon,
} from '@heroicons/vue/24/outline'

const toast      = useToast()
const { t }      = useI18n()
const categories = ref([])
const loading    = ref(false)
const deleteConfirm = reactive({ show: false, category: null, loading: false })

const categoryIcons = [CubeIcon, SparklesIcon, WrenchScrewdriverIcon, GlobeAltIcon, HeartIcon, TruckIcon, TagIcon]
const categoryColors = [
  { bg: '#dcfce7', color: '#006d35' },
  { bg: '#dbeafe', color: '#1d4ed8' },
  { bg: '#fef9c3', color: '#854d0e' },
  { bg: '#fce7f3', color: '#9d174d' },
  { bg: '#e0e7ff', color: '#4338ca' },
  { bg: '#fff7ed', color: '#c2410c' },
  { bg: '#f0fdf4', color: '#166534' },
]

function catIcon(cat) {
  return categoryIcons[cat.id % categoryIcons.length]
}
function catIconBg(cat) {
  return categoryColors[cat.id % categoryColors.length].bg
}
function catIconColor(cat) {
  return categoryColors[cat.id % categoryColors.length].color
}

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
    toast.showSuccess(t('categories.toastToggledMsg'))
  } catch {
    toast.showError(t('categories.toastToggleError'))
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
    toast.showError(t('categories.toastDeleteErrorMsg'))
  } finally {
    deleteConfirm.loading = false
  }
}

onMounted(fetchCategories)
</script>
