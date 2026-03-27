<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/50" @click="$emit('cancel')" />
      <div class="relative bg-white rounded-lg shadow-xl w-full max-w-md mx-4 p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ title }}</h3>
        <p class="text-sm text-gray-500 mb-6">{{ message }}</p>
        <div class="flex justify-end gap-3">
          <button
            @click="$emit('cancel')"
            class="btn-secondary"
          >
            Annuler
          </button>
          <button
            @click="$emit('confirm')"
            :disabled="loading"
            :class="['btn', confirmVariantClass]"
          >
            <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
            </svg>
            {{ confirmLabel }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>import { computed } from 'vue'

const props = defineProps({
  show:         { type: Boolean, default: false },
  title:        { type: String, default: 'Confirmer l\'action' },
  message:      { type: String, default: 'Êtes-vous sûr de vouloir effectuer cette action ?' },
  confirmLabel: { type: String, default: 'Confirmer' },
  confirmVariant: { type: String, default: 'danger' },
  loading:      { type: Boolean, default: false },
})

defineEmits(['confirm', 'cancel'])

const confirmVariantClass = computed(() => ({
  danger:  'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 px-4 py-2 text-sm rounded-md',
  primary: 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500 px-4 py-2 text-sm rounded-md',
}[props.confirmVariant] || 'bg-red-600 text-white hover:bg-red-700 px-4 py-2 text-sm rounded-md'))
</script>
