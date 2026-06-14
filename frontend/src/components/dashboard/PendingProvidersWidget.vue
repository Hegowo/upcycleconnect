<template>
  <div class="bg-white rounded-2xl p-4 border border-[#f1f5f9] shadow-sm h-full flex flex-col justify-center" style="border-top: 3px solid #006d35;">
    <div class="flex items-start justify-between mb-2">
      <div class="w-9 h-9 rounded-xl flex items-center justify-center" style="background-color:#dcfce7;">
        <BriefcaseIcon class="w-5 h-5" style="color:#006d35;" />
      </div>
      <span v-if="pending > 0" class="text-xs font-semibold px-2 py-0.5 rounded-full" style="background:#dcfce7; color:#166534;">
        {{ pending }} {{ t('dashboard.statsPending').toLowerCase() }}
      </span>
    </div>
    <h3 class="font-bold text-sm text-[#001d32]">{{ t('dashboard.validateProviders') }}</h3>
    <p class="text-xs text-gray-500 mt-0.5 leading-snug">{{ t('dashboard.validateProvidersDesc') }}</p>
    <RouterLink to="/admin/providers" class="inline-flex items-center gap-1.5 text-sm font-semibold transition mt-2" style="color:#006d35;">
      {{ t('dashboard.seeProviders') }}
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/></svg>
    </RouterLink>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import { BriefcaseIcon } from '@heroicons/vue/24/outline'

const { t } = useI18n()
defineProps({ pending: { type: Number, default: 0 } })
</script>
