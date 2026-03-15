<template>
  <RouterLink :to="link" class="block group rounded-2xl p-5 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-lg" :style="cardStyle">
    <div class="flex items-start justify-between mb-3">
      <div class="w-10 h-10 rounded-xl flex items-center justify-center text-xl" :style="iconBgStyle">
        {{ icon }}
      </div>
      <span v-if="alert" class="text-xs font-semibold px-2 py-0.5 rounded-full bg-orange-100 text-orange-700 animate-pulse">
        Action requise
      </span>
    </div>
    <p class="text-3xl font-bold mb-1" :style="{ color: textColor }">{{ value }}</p>
    <p class="text-sm font-medium opacity-70" :style="{ color: textColor }">{{ label }}</p>
  </RouterLink>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label:  String,
  value:  [Number, String],
  icon:   { type: String, default: '📊' },
  color:  { type: String, default: 'blue' },
  link:   { type: String, default: '#' },
  alert:  { type: Boolean, default: false },
})

const palettes = {
  blue:   { bg: '#EFF6FF', iconBg: '#DBEAFE', text: '#1e40af' },
  green:  { bg: '#F0FDF4', iconBg: '#DCFCE7', text: '#15803d' },
  orange: { bg: '#FFF7ED', iconBg: '#FFEDD5', text: '#c2410c' },
  purple: { bg: '#FAF5FF', iconBg: '#EDE9FE', text: '#7c3aed' },
  teal:   { bg: '#F0FDFA', iconBg: '#CCFBF1', text: '#0f766e' },
  navy:   { bg: '#EFF6FF', iconBg: '#DBEAFE', text: '#103652' },
}

const palette = computed(() => palettes[props.color] || palettes.blue)
const cardStyle   = computed(() => ({ backgroundColor: palette.value.bg }))
const iconBgStyle = computed(() => ({ backgroundColor: palette.value.iconBg }))
const textColor   = computed(() => palette.value.text)
</script>
