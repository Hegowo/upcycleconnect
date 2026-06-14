<template>
  <div class="bg-white rounded-2xl p-4 sm:p-6 border border-[#f1f5f9] shadow-sm h-full flex flex-col">
    <div class="flex items-center justify-between mb-4 sm:mb-5 shrink-0">
      <div>
        <h3 class="font-semibold text-[#001d32] text-sm sm:text-base">{{ t('dashboard.trends') }}</h3>
        <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5">{{ t('dashboard.trendsSubtitle') }} — {{ year }}</p>
      </div>
      <span class="text-xs border border-[#e5e7eb] rounded-lg px-2 py-1 text-gray-500 bg-[#f8fafc]">{{ year }}</span>
    </div>

    <div class="flex-1 flex items-end gap-1 sm:gap-2 min-h-[6rem]">
      <div v-for="(bar, i) in monthlyData" :key="i" class="flex-1 flex flex-col items-center gap-0.5 sm:gap-1 h-full justify-end group relative">
        <span v-if="bar.value > 0" class="text-[9px] font-semibold text-[#006d35] opacity-0 group-hover:opacity-100 transition">{{ bar.value }}</span>
        <div
          class="w-full rounded-t-md transition-all duration-500"
          :style="{
            height: (bar.value / maxMonthly * 100) + '%',
            backgroundColor: i === currentMonth ? '#006d35' : (bar.value > 0 ? '#86efac' : '#eef2f7'),
            minHeight: bar.value > 0 ? '6px' : '3px',
          }"
          :title="`${bar.label} : ${bar.value}`"
        ></div>
        <span class="text-[8px] sm:text-[9px] text-gray-400">{{ bar.label }}</span>
      </div>
    </div>

    <p v-if="!loading && totalSignups === 0" class="text-center text-[11px] text-gray-400 mt-3 shrink-0">
      {{ t('dashboard.trendsEmpty') }}
    </p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'

const { t } = useI18n()
const LABELS = ['Jan', 'Fév', 'Mar', 'Avr', 'Mai', 'Jun', 'Jul', 'Aoû', 'Sep', 'Oct', 'Nov', 'Déc']

const year = ref(new Date().getFullYear())
const currentMonth = new Date().getMonth()
const counts = ref(Array(12).fill(0))
const loading = ref(true)

const monthlyData = computed(() => LABELS.map((label, i) => ({ label, value: counts.value[i] || 0 })))
const maxMonthly = computed(() => Math.max(1, ...counts.value))
const totalSignups = computed(() => counts.value.reduce((a, b) => a + b, 0))

onMounted(async () => {
  try {
    const { data } = await api.get('/dashboard/trends')
    if (Array.isArray(data.monthly_signups)) counts.value = data.monthly_signups
    if (data.year) year.value = data.year
  } catch {  } finally {
    loading.value = false
  }
})
</script>
