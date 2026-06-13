<template>
  <div>
    <label v-if="label" class="block text-xs font-medium text-[#001d32] mb-1.5">{{ label }}</label>
    <div class="relative">
      <LanguageIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#94a3b8] pointer-events-none" />
      <select
        :value="current"
        @change="onChange($event.target.value)"
        class="w-full pl-10 pr-4 py-2.5 border border-[#e2e8f0] rounded-xl text-sm bg-white focus:outline-none focus:border-[#006d35] transition-colors appearance-none cursor-pointer"
      >
        <option v-for="loc in locales" :key="loc.code" :value="loc.code">{{ loc.name }}</option>
      </select>
      <ChevronDownIcon class="w-4 h-4 absolute right-3 top-1/2 -translate-y-1/2 text-[#94a3b8] pointer-events-none" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { LanguageIcon, ChevronDownIcon } from '@heroicons/vue/24/outline'
import { availableLocales, setLocale } from '@/utils/i18n.js'

defineProps({
  label: { type: String, default: '' },
})

const { locale } = useI18n()
const locales = availableLocales
const current = computed(() => locale.value)

async function onChange(code) {
  await setLocale(code)
}
</script>
