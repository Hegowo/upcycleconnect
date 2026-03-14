<template>
  <RouterLink
    :to="link"
    :class="['card p-5 hover:shadow-md transition-shadow cursor-pointer block', alert ? 'ring-2 ring-yellow-400' : '']"
  >
    <p class="text-sm text-gray-500 mb-1">{{ label }}</p>
    <p :class="['text-3xl font-bold', colorClass]">{{ value }}</p>
    <p v-if="alert" class="text-xs text-yellow-600 mt-1 font-medium">Action requise</p>
  </RouterLink>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: String,
  value: [Number, String],
  color: { type: String, default: 'gray' },
  link:  { type: String, default: '#' },
  alert: { type: Boolean, default: false },
})

const colorClass = computed(() => ({
  blue:   'text-blue-600',
  green:  'text-green-600',
  orange: 'text-orange-500',
  purple: 'text-purple-600',
  teal:   'text-teal-600',
  gray:   'text-gray-700',
}[props.color] || 'text-gray-700'))
</script>
