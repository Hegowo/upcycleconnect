<template>
  <div class="relative" ref="container">
    <div class="relative">
      <input
        ref="inputEl"
        :value="modelValue"
        @input="onInput"
        @keydown="onKeydown"
        @focus="onFocus"
        type="text"
        :placeholder="placeholder"
        :disabled="disabled"
        autocomplete="off"
        v-bind="$attrs"
      />
      <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
        <div v-if="loading" class="w-4 h-4 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin"></div>
        <MapPinIcon v-else-if="!modelValue" class="w-4 h-4 text-gray-300" />
        <CheckCircleIcon v-else class="w-4 h-4 text-[#006d35]" />
      </div>
    </div>

    <Transition name="dropdown">
      <ul
        v-if="showDropdown && suggestions.length"
        class="absolute z-50 w-full mt-1 bg-white border border-[#e2e8f0] rounded-xl shadow-lg overflow-hidden"
      >
        <li
          v-for="(s, i) in suggestions"
          :key="s.properties.id || i"
          @mousedown.prevent="selectSuggestion(s)"
          @mouseover="activeIndex = i"
          class="flex items-start gap-3 px-4 py-3 cursor-pointer transition-colors text-sm"
          :class="activeIndex === i ? 'bg-[#f0fdf4]' : 'hover:bg-[#f7f9ff]'"
        >
          <MapPinIcon class="w-4 h-4 text-[#006d35] shrink-0 mt-0.5" />
          <div class="min-w-0">
            <span class="font-medium text-[#001d32] block truncate">{{ s.properties.name }}</span>
            <span class="text-[#40617f] text-xs">{{ s.properties.postcode }} {{ s.properties.city }}</span>
          </div>
        </li>
      </ul>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { MapPinIcon, CheckCircleIcon } from '@heroicons/vue/24/outline'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue', 'select'])

const container  = ref(null)
const inputEl    = ref(null)
const suggestions = ref([])
const loading    = ref(false)
const showDropdown = ref(false)
const activeIndex  = ref(-1)

let debounceTimer = null

function onInput(e) {
  const val = e.target.value
  emit('update:modelValue', val)
  activeIndex.value = -1
  clearTimeout(debounceTimer)
  if (val.length < 3) {
    suggestions.value = []
    showDropdown.value = false
    return
  }
  debounceTimer = setTimeout(() => fetchSuggestions(val), 300)
}

function onFocus() {
  if (suggestions.value.length) showDropdown.value = true
}

async function fetchSuggestions(query) {
  loading.value = true
  try {
    const url = `https://api-adresse.data.gouv.fr/search/?q=${encodeURIComponent(query)}&limit=5&autocomplete=1`
    const res = await fetch(url)
    const data = await res.json()
    suggestions.value = data.features || []
    showDropdown.value = suggestions.value.length > 0
  } catch {
    suggestions.value = []
    showDropdown.value = false
  } finally {
    loading.value = false
  }
}

function selectSuggestion(s) {
  const label = s.properties.label
  emit('update:modelValue', label)
  emit('select', s.properties)
  suggestions.value = []
  showDropdown.value = false
  activeIndex.value = -1
}

function onKeydown(e) {
  if (!showDropdown.value || !suggestions.value.length) return
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    activeIndex.value = Math.min(activeIndex.value + 1, suggestions.value.length - 1)
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    activeIndex.value = Math.max(activeIndex.value - 1, 0)
  } else if (e.key === 'Enter' && activeIndex.value >= 0) {
    e.preventDefault()
    selectSuggestion(suggestions.value[activeIndex.value])
  } else if (e.key === 'Escape') {
    showDropdown.value = false
    activeIndex.value = -1
  }
}

function onClickOutside(e) {
  if (container.value && !container.value.contains(e.target)) {
    showDropdown.value = false
    activeIndex.value = -1
  }
}

onMounted(() => document.addEventListener('mousedown', onClickOutside))
onBeforeUnmount(() => document.removeEventListener('mousedown', onClickOutside))
</script>

<style scoped>
.dropdown-enter-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.dropdown-leave-active { transition: opacity 0.1s ease, transform 0.1s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
