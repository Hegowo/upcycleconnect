<template>
  <div v-if="lastPage > 1" class="flex items-center justify-between mt-4">
    <p class="text-sm text-gray-500">
      Page {{ currentPage }} sur {{ lastPage }}
    </p>
    <div class="flex gap-1">
      <button
        @click="$emit('page-change', currentPage - 1)"
        :disabled="currentPage <= 1"
        class="px-3 py-1.5 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        Précédent
      </button>
      <button
        v-for="page in visiblePages"
        :key="page"
        @click="page !== '...' && $emit('page-change', page)"
        :class="[
          'px-3 py-1.5 text-sm border rounded-md',
          page === currentPage
            ? 'bg-green-600 text-white border-green-600'
            : page === '...'
              ? 'border-transparent cursor-default'
              : 'border-gray-300 hover:bg-gray-50',
        ]"
      >
        {{ page }}
      </button>
      <button
        @click="$emit('page-change', currentPage + 1)"
        :disabled="currentPage >= lastPage"
        class="px-3 py-1.5 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-40 disabled:cursor-not-allowed"
      >
        Suivant
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  currentPage: { type: Number, required: true },
  lastPage:    { type: Number, required: true },
  total:       { type: Number, default: 0 },
})

defineEmits(['page-change'])

const visiblePages = computed(() => {
  const pages = []
  const delta = 2
  for (let i = 1; i <= props.lastPage; i++) {
    if (i === 1 || i === props.lastPage || (i >= props.currentPage - delta && i <= props.currentPage + delta)) {
      pages.push(i)
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...')
    }
  }
  return pages
})
</script>
