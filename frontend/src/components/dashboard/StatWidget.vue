<template>
  <component
    :is="to ? 'RouterLink' : 'div'"
    :to="to || undefined"
    class="bg-white rounded-2xl p-4 sm:p-5 border border-[#f1f5f9] shadow-sm block h-full"
    :class="to ? 'hover:shadow-md transition-shadow' : ''"
  >
    <div class="flex items-center justify-between mb-2 sm:mb-3">
      <div class="w-8 h-8 sm:w-10 sm:h-10 rounded-xl flex items-center justify-center" :style="{ backgroundColor: iconBg }">
        <component :is="icon" class="w-4 h-4 sm:w-5 sm:h-5" :style="{ color: iconColor }" />
      </div>
      <div class="flex items-center gap-1">
        <span v-if="cornerLabel" class="text-[10px] sm:text-xs font-medium text-gray-400 text-right leading-tight">{{ cornerLabel }}</span>
        <div v-if="tooltip" class="relative group">
          <InformationCircleIcon class="w-4 h-4 text-gray-400 cursor-help" />
          <div class="absolute right-0 top-5 w-52 bg-[#001d32] text-white text-xs rounded-xl p-3 z-10 hidden group-hover:block shadow-lg leading-relaxed">
            {{ tooltip }}
          </div>
        </div>
      </div>
    </div>
    <p class="text-xl sm:text-2xl font-bold text-[#001d32]">
      {{ value }}<span v-if="suffix" class="text-xs sm:text-sm font-normal text-gray-400">{{ suffix }}</span>
    </p>
    <p class="text-[10px] sm:text-xs text-gray-400 mt-0.5 sm:mt-1 leading-tight">{{ label }}</p>
  </component>
</template>

<script setup>
import { InformationCircleIcon } from '@heroicons/vue/24/outline'

defineProps({
  icon:        { type: [Object, Function], required: true },
  iconBg:      { type: String, default: '#dcfce7' },
  iconColor:   { type: String, default: '#006d35' },
  value:       { type: [String, Number], default: 0 },
  suffix:      { type: String, default: '' },
  label:       { type: String, default: '' },
  cornerLabel: { type: String, default: '' },
  tooltip:     { type: String, default: '' },
  to:          { type: String, default: '' },
})
</script>
