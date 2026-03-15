<template>
  <div class="space-y-6 max-w-2xl">
    <div>
      <h2 class="text-xl font-bold text-gray-900">{{ t('settings.title') }}</h2>
      <p class="text-sm text-gray-500 mt-0.5">{{ t('settings.subtitle') }}</p>
    </div>

    <div class="card p-6 space-y-4">
      <h3 class="font-semibold text-gray-900">{{ t('settings.sectionLanguage') }}</h3>
      <p class="text-sm text-gray-500">{{ t('settings.langDescription') }}</p>
      <div class="flex gap-3">
        <button
          @click="setLang('fr')"
          :class="currentLang === 'fr' ? 'ring-2 font-semibold' : 'border text-gray-500 hover:border-gray-400'"
          class="flex items-center gap-2 px-5 py-3 rounded-xl transition text-sm"
          :style="currentLang === 'fr' ? 'ring-color:#1B8848; border-color:#1B8848; color:#1B8848; border: 2px solid #1B8848;' : 'border: 1px solid #d1d5db;'"
        >
          🇫🇷 {{ t('settings.langFr') }}
        </button>
        <button
          @click="setLang('en')"
          :class="currentLang === 'en' ? 'font-semibold' : 'text-gray-500 hover:border-gray-400'"
          class="flex items-center gap-2 px-5 py-3 rounded-xl transition text-sm"
          :style="currentLang === 'en' ? 'border: 2px solid #1B8848; color:#1B8848;' : 'border: 1px solid #d1d5db;'"
        >
          🇬🇧 {{ t('settings.langEn') }}
        </button>
      </div>
    </div>

    <div class="card p-6 space-y-4">
      <h3 class="font-semibold text-gray-900">{{ t('settings.sectionProfile') }}</h3>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="label">{{ t('settings.fieldFirstName') }}</label>
          <p class="text-sm text-gray-800 mt-1 px-3 py-2 bg-gray-50 rounded-lg">{{ auth.user?.first_name || '—' }}</p>
        </div>
        <div>
          <label class="label">{{ t('settings.fieldLastName') }}</label>
          <p class="text-sm text-gray-800 mt-1 px-3 py-2 bg-gray-50 rounded-lg">{{ auth.user?.last_name || '—' }}</p>
        </div>
        <div>
          <label class="label">{{ t('settings.fieldEmail') }}</label>
          <p class="text-sm text-gray-800 mt-1 px-3 py-2 bg-gray-50 rounded-lg">{{ auth.user?.email || '—' }}</p>
        </div>
        <div>
          <label class="label">{{ t('settings.fieldRole') }}</label>
          <p class="text-sm capitalize mt-1 px-3 py-2 bg-gray-50 rounded-lg font-medium" style="color:#1B8848;">{{ auth.user?.role?.replace('_', ' ') || '—' }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/utils/useToast'

const { t, locale } = useI18n()
const auth = useAuthStore()
const toast = useToast()

const currentLang = computed(() => locale.value)

function setLang(lang) {
  locale.value = lang
  localStorage.setItem('admin_lang', lang)
  toast.showSuccess(t('settings.toastLangUpdated'))
}
</script>
