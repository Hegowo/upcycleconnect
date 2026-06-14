<template>
  <div class="space-y-6 max-w-3xl">

    <div>
      <h2 class="text-2xl font-bold text-[#001d32]">{{ t('settings.title') }}</h2>
      <p class="text-sm text-[#40617f] mt-0.5">{{ t('settings.subtitle') }}</p>
    </div>

    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm p-6 space-y-6">
      <h3 class="font-semibold text-[#001d32] text-base">{{ t('settings.globalSettings') }}</h3>

      <div class="flex items-center justify-between py-3 border-b border-[#f8fafc]">
        <div>
          <p class="text-sm font-medium text-[#001d32]">{{ t('settings.autoValidation') }}</p>
          <p class="text-xs text-gray-400 mt-0.5">{{ t('settings.autoValidationDesc') }}</p>
        </div>
        <button
          @click="toggleAutoValidation"
          :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors', autoValidation ? 'bg-[#006d35]' : 'bg-gray-300']"
        >
          <span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform shadow', autoValidation ? 'translate-x-6' : 'translate-x-1']" />
        </button>
      </div>

      <div class="flex items-center justify-between py-3 border-b border-[#f8fafc]">
        <div>
          <p class="text-sm font-medium text-[#001d32]">{{ t('settings.emailNotif') }}</p>
          <p class="text-xs text-gray-400 mt-0.5">{{ t('settings.emailNotifDesc') }}</p>
        </div>
        <button
          @click="emailNotifications = !emailNotifications"
          :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors', emailNotifications ? 'bg-[#006d35]' : 'bg-gray-300']"
        >
          <span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform shadow', emailNotifications ? 'translate-x-6' : 'translate-x-1']" />
        </button>
      </div>

      <div class="py-3">
        <p class="text-sm font-medium text-[#001d32] mb-3">{{ t('settings.profileInfo') }}</p>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ t('settings.fieldFirstName') }}</label>
            <p class="text-sm text-[#001d32] mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg">{{ auth.user?.first_name || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ t('settings.fieldLastName') }}</label>
            <p class="text-sm text-[#001d32] mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg">{{ auth.user?.last_name || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ t('settings.fieldEmail') }}</label>
            <p class="text-sm text-[#001d32] mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg">{{ auth.user?.email || '—' }}</p>
          </div>
          <div>
            <label class="text-xs font-medium text-gray-400 uppercase tracking-wider">{{ t('settings.fieldRole') }}</label>
            <p class="text-sm capitalize mt-1 px-3 py-2 bg-[#f8fafc] rounded-lg font-medium" style="color:#006d35;">{{ auth.user?.role?.replace('_', ' ') || '—' }}</p>
          </div>
        </div>
      </div>

      <div class="flex items-center gap-3 pt-2">
        <button class="px-6 py-2.5 text-sm font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition">
          {{ t('settings.cancelBtn') }}
        </button>
        <button @click="saveSettings" class="px-6 py-2.5 text-sm font-semibold rounded-lg text-white transition hover:opacity-90" style="background-color:#006d35;">
          {{ t('settings.saveBtn') }}
        </button>
      </div>
    </div>

  </div>
</template>

<script setup>import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/utils/useToast'

const { t } = useI18n()
const auth = useAuthStore()
const toast = useToast()

const autoValidation   = ref(false)
const emailNotifications = ref(true)

function toggleAutoValidation() {
  autoValidation.value = !autoValidation.value
}

function saveSettings() {
  toast.showSuccess(t('settings.toastSaved'))
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
