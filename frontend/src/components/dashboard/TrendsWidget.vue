<template>
  <div class="bg-white rounded-2xl p-4 sm:p-6 border border-[#f1f5f9] shadow-sm h-full">
    <div class="flex items-center justify-between mb-4 sm:mb-5">
      <div>
        <h3 class="font-semibold text-[#001d32] text-sm sm:text-base">{{ t('dashboard.trends') }}</h3>
        <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5">{{ t('dashboard.trendsSubtitle') }} — {{ year }}</p>
      </div>
      <span class="text-xs border border-[#e5e7eb] rounded-lg px-2 py-1 text-gray-500 bg-[#f8fafc]">{{ year }}</span>
    </div>
    <div class="flex items-end gap-1 sm:gap-2 h-24 sm:h-32">
      <div v-for="(bar, i) in monthlyData" :key="i" class="flex-1 flex flex-col items-center gap-0.5 sm:gap-1">
        <div
          class="w-full rounded-t-md transition-all duration-500"
          :style="{
            height: (bar.value / maxMonthly * 100) + '%',
            backgroundColor: i === currentMonth ? '#006d35' : '#dcfce7',
            minHeight: '4px',
          }"
        ></div>
        <span class="text-[8px] sm:text-[9px] text-gray-400">{{ bar.label }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const year = new Date().getFullYear()
const currentMonth = new Date().getMonth()

const monthlyData = [
  { label: 'Jan', value: 42 }, { label: 'Fév', value: 58 }, { label: 'Mar', value: 71 },
  { label: 'Avr', value: 65 }, { label: 'Mai', value: 89 }, { label: 'Jun', value: 103 },
  { label: 'Jul', value: 94 }, { label: 'Aoû', value: 78 }, { label: 'Sep', value: 112 },
  { label: 'Oct', value: 87 }, { label: 'Nov', value: 134 }, { label: 'Déc', value: 98 },
]
const maxMonthly = computed(() => Math.max(...monthlyData.map(d => d.value)))
</script>
