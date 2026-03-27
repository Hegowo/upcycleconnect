<template>
  <span :class="['inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium', variantClass]">
    {{ label }}
  </span>
</template>

<script setup>import { computed } from 'vue'

const props = defineProps({
  variant: { type: String, default: 'neutral' },
  label:   { type: String, required: true },
})

const STATUS_MAP = {
  active:    'success',
  approved:  'success',
  published: 'success',
  pending:   'warning',
  draft:     'warning',
  inactive:  'neutral',
  suspended: 'danger',
  banned:    'danger',
  cancelled: 'danger',
  archived:  'neutral',
}

const resolvedVariant = computed(() => STATUS_MAP[props.label] || STATUS_MAP[props.variant] || props.variant)

const variantClass = computed(() => ({
  success: 'bg-green-100 text-green-800',
  warning: 'bg-yellow-100 text-yellow-800',
  danger:  'bg-red-100 text-red-800',
  info:    'bg-blue-100 text-blue-800',
  neutral: 'bg-gray-100 text-gray-700',
}[resolvedVariant.value] || 'bg-gray-100 text-gray-700'))
</script>
