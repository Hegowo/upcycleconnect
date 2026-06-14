<template>
  <Teleport to="body">
    <Transition name="slide">
      <div v-if="open" class="fixed inset-0 z-50 flex justify-end" @click.self="$emit('update:open', false)">
        <div class="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
        <aside class="relative w-full sm:max-w-sm h-full bg-white shadow-2xl flex flex-col">
          <header class="flex items-center justify-between px-5 py-4 border-b border-[#f1f5f9]">
            <div>
              <h3 class="font-bold text-[#001d32]">{{ t('dashboard.customizeTitle') }}</h3>
              <p class="text-xs text-gray-400 mt-0.5">{{ t('dashboard.customizeHint') }}</p>
            </div>
            <button @click="$emit('update:open', false)" class="text-gray-400 hover:text-gray-600 transition shrink-0">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </header>

          <div class="flex-1 overflow-y-auto p-4 space-y-2">
            <div
              v-for="(it, i) in local"
              :key="it.id"
              draggable="true"
              @dragstart="dragIndex = i"
              @dragover.prevent="onDragOver(i)"
              @dragend="dragIndex = null"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl border transition cursor-grab active:cursor-grabbing"
              :class="[
                it.visible ? 'border-[#e2e8f0] bg-white' : 'border-dashed border-[#e2e8f0] bg-[#f8fafc] opacity-70',
                dragIndex === i ? 'ring-2 ring-[#006d35]' : '',
              ]"
            >
              <Bars3Icon class="w-4 h-4 text-gray-300 shrink-0" />
              <span class="flex-1 text-sm font-medium text-[#001d32] truncate">{{ it.label }}</span>
              <button
                type="button"
                role="switch"
                :aria-checked="it.visible"
                @click="toggle(i)"
                class="relative w-9 h-5 rounded-full transition shrink-0"
                :style="{ backgroundColor: it.visible ? '#006d35' : '#cbd5e1' }"
              >
                <span class="absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full transition-transform" :class="it.visible ? 'translate-x-4' : ''"></span>
              </button>
            </div>
          </div>

          <footer class="flex items-center justify-between gap-3 px-5 py-4 border-t border-[#f1f5f9]">
            <button @click="$emit('reset')" class="text-sm font-medium text-gray-500 hover:text-gray-700 transition">
              {{ t('dashboard.layoutReset') }}
            </button>
            <button @click="$emit('update:open', false)" class="px-4 py-2 rounded-lg text-white text-sm font-semibold transition hover:opacity-90" style="background-color:#006d35;">
              {{ t('common.close') }}
            </button>
          </footer>
        </aside>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { XMarkIcon, Bars3Icon } from '@heroicons/vue/24/outline'

const { t } = useI18n()
const props = defineProps({
  open:  { type: Boolean, default: false },
  items: { type: Array, default: () => [] }, // [{ id, visible, label }]
})
const emit = defineEmits(['update:open', 'update:items', 'reset'])

// Local working copy so drag reordering stays smooth; we emit on every change.
const local = ref([])
watch(() => props.items, (v) => { local.value = v.map(x => ({ ...x })) }, { immediate: true, deep: true })

const dragIndex = ref(null)

function emitItems() {
  emit('update:items', local.value.map(x => ({ id: x.id, visible: x.visible })))
}
function toggle(i) {
  local.value[i].visible = !local.value[i].visible
  emitItems()
}
function onDragOver(i) {
  if (dragIndex.value === null || dragIndex.value === i) return
  const moved = local.value.splice(dragIndex.value, 1)[0]
  local.value.splice(i, 0, moved)
  dragIndex.value = i
  emitItems()
}
</script>

<style scoped>
.slide-enter-active, .slide-leave-active { transition: opacity 0.2s ease; }
.slide-enter-from, .slide-leave-to { opacity: 0; }
.slide-enter-active aside, .slide-leave-active aside { transition: transform 0.25s ease; }
.slide-enter-from aside, .slide-leave-to aside { transform: translateX(100%); }
</style>
