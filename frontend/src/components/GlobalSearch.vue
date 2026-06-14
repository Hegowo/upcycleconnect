<template>
  <Teleport to="body">
    <Transition name="search-fade">
      <div v-if="isOpen" class="fixed inset-0 z-[200] flex items-start justify-center pt-[8vh] px-4">

        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm" @click="close" />

        <div
          ref="panelRef"
          class="search-panel relative w-full max-w-[660px] bg-white rounded-2xl shadow-2xl overflow-hidden flex flex-col"
          style="max-height: 80vh;"
        >

          <div class="flex items-center gap-3 px-5 py-4 border-b border-[#e5e7eb] shrink-0">
            <MagnifyingGlassIcon class="w-5 h-5 text-[#94a3b8] shrink-0" />
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              :placeholder="t('public.search.placeholder')"
              class="flex-1 text-[15px] text-[#001d32] placeholder:text-[#94a3b8] outline-none bg-transparent"
              @keydown="handleKey"
            />
            <div class="flex items-center gap-2 shrink-0">
              <span v-if="loading" class="w-4 h-4 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin" />
              <button @click="close" class="text-[#94a3b8] hover:text-[#001d32] transition">
                <XMarkIcon class="w-4 h-4" />
              </button>
            </div>
          </div>

          <div ref="scrollRef" class="overflow-y-auto flex-1">

            <div v-if="!query.trim()" class="p-3">

              <div v-if="recent.length > 0" class="mb-5">
                <div class="flex items-center justify-between px-3 mb-2">
                  <p class="text-[10px] font-bold text-[#94a3b8] uppercase tracking-widest">{{ t('public.search.recentTitle') }}</p>
                  <button @click="clearRecent" class="text-[10px] text-[#94a3b8] hover:text-[#ef4444] transition">{{ t('public.search.clearRecent') }}</button>
                </div>
                <div class="flex flex-col gap-0.5">
                  <button
                    v-for="(r, i) in recent"
                    :key="i"
                    @click="query = r; inputRef?.focus()"
                    class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-[#f8fafc] transition text-left group"
                  >
                    <ClockIcon class="w-4 h-4 text-[#94a3b8] shrink-0" />
                    <span class="text-sm text-[#374151] flex-1">{{ r }}</span>
                    <XMarkIcon class="w-3.5 h-3.5 text-[#c4c4c4] opacity-0 group-hover:opacity-100 transition" @click.stop="removeRecent(i)" />
                  </button>
                </div>
              </div>

              <p class="text-[10px] font-bold text-[#94a3b8] uppercase tracking-widest px-3 mb-2">{{ t('public.search.quickNav') }}</p>
              <div class="grid grid-cols-2 gap-1.5">
                <button
                  v-for="(action, i) in quickActions"
                  :key="i"
                  @click="go(action.url)"
                  class="flex items-center gap-3 px-3 py-2.5 rounded-xl border border-[#f1f5f9] hover:border-[#006d35]/30 hover:bg-[#f0fdf4] transition text-left group"
                >
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0" :style="`background: ${action.bg}`">
                    <component :is="action.icon" class="w-4 h-4" :style="`color: ${action.color}`" />
                  </div>
                  <div class="min-w-0">
                    <p class="text-sm font-semibold text-[#001d32] truncate">{{ action.label }}</p>
                    <p class="text-xs text-[#94a3b8] truncate">{{ action.sub }}</p>
                  </div>
                </button>
              </div>

              <p class="text-center text-[11px] text-[#c4c4c4] mt-5 mb-1">
                <kbd class="bg-[#f1f5f9] px-1.5 py-0.5 rounded text-[#64748b] font-mono">⌘K</kbd>
                {{ t('public.search.shortcutHint') }}
              </p>
            </div>

            <div v-else-if="hasResults" class="p-2">
              <template v-for="section in sections" :key="section.key">
                <div v-if="results[section.key]?.length > 0" class="mb-3">
                  <p class="text-[10px] font-bold text-[#94a3b8] uppercase tracking-widest px-3 py-1.5">{{ section.label }}</p>
                  <div class="flex flex-col gap-0.5">
                    <button
                      v-for="(item, j) in results[section.key]"
                      :key="item.id"
                      :ref="el => setRef(el, section.key, j)"
                      @click="selectItem(item)"
                      @mouseenter="hoveredKey = `${section.key}-${j}`"
                      :class="[
                        'w-full flex items-center gap-3 px-3 py-2.5 rounded-xl transition text-left',
                        hoveredKey === `${section.key}-${j}` || flatIndex(section.key, j) === focusedIndex
                          ? 'bg-[#006d35] text-white'
                          : 'hover:bg-[#f8fafc]'
                      ]"
                    >

                      <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition"
                        :class="hoveredKey === `${section.key}-${j}` || flatIndex(section.key, j) === focusedIndex ? 'bg-white/20' : section.iconBg">
                        <component :is="section.icon" class="w-4 h-4 transition"
                          :class="hoveredKey === `${section.key}-${j}` || flatIndex(section.key, j) === focusedIndex ? 'text-white' : section.iconColor" />
                      </div>

                      <div class="flex-1 min-w-0">
                        <p class="text-sm font-semibold truncate"
                          :class="hoveredKey === `${section.key}-${j}` || flatIndex(section.key, j) === focusedIndex ? 'text-white' : 'text-[#001d32]'"
                          v-html="highlight(item.title)" />
                        <p v-if="item.subtitle" class="text-xs truncate"
                          :class="hoveredKey === `${section.key}-${j}` || flatIndex(section.key, j) === focusedIndex ? 'text-white/70' : 'text-[#94a3b8]'">
                          {{ item.subtitle }}
                        </p>
                      </div>

                      <span v-if="item.meta" class="text-xs shrink-0 font-medium px-2 py-0.5 rounded-full transition"
                        :class="hoveredKey === `${section.key}-${j}` || flatIndex(section.key, j) === focusedIndex ? 'bg-white/20 text-white' : 'bg-[#f1f5f9] text-[#64748b]'">
                        {{ item.meta }}
                      </span>

                      <ChevronRightIcon class="w-4 h-4 shrink-0 opacity-40" />
                    </button>
                  </div>
                </div>
              </template>
            </div>

            <div v-else-if="query.trim().length >= 2 && !loading" class="py-14 text-center">
              <div class="w-14 h-14 rounded-2xl bg-[#f1f5f9] flex items-center justify-center mx-auto mb-3">
                <MagnifyingGlassIcon class="w-7 h-7 text-[#94a3b8]" />
              </div>
              <p class="text-[#374151] font-semibold text-sm">{{ t('public.search.noResults') }}</p>
              <p class="text-[#94a3b8] text-xs mt-1">{{ t('public.search.noResultsSub', { q: query }) }}</p>
            </div>

            <div v-else-if="query.trim().length === 1" class="py-10 text-center text-[#94a3b8] text-sm">
              {{ t('public.search.typeMore') }}
            </div>
          </div>

          <div class="shrink-0 border-t border-[#e5e7eb] bg-[#f8fafc] px-5 py-2.5 flex items-center gap-5">
            <span class="text-[11px] text-[#94a3b8] flex items-center gap-1">
              <kbd class="bg-white border border-[#e5e7eb] px-1.5 py-0.5 rounded text-[10px] font-mono shadow-sm">↑↓</kbd>
              Naviguer
            </span>
            <span class="text-[11px] text-[#94a3b8] flex items-center gap-1">
              <kbd class="bg-white border border-[#e5e7eb] px-1.5 py-0.5 rounded text-[10px] font-mono shadow-sm">↵</kbd>
              Ouvrir
            </span>
            <span class="text-[11px] text-[#94a3b8] flex items-center gap-1">
              <kbd class="bg-white border border-[#e5e7eb] px-1.5 py-0.5 rounded text-[10px] font-mono shadow-sm">Esc</kbd>
              Fermer
            </span>
            <span class="ml-auto text-[11px] text-[#94a3b8]">UpcycleConnect</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  MagnifyingGlassIcon,
  XMarkIcon,
  ClockIcon,
  ChevronRightIcon,
  CalendarDaysIcon,
  WrenchScrewdriverIcon,
  ChatBubbleLeftEllipsisIcon,
  BuildingStorefrontIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const router = useRouter()

const isOpen = ref(false)
const query = ref('')
const loading = ref(false)
const results = ref({ events: [], prestations: [], threads: [], providers: [] })
const focusedIndex = ref(-1)
const hoveredKey = ref('')
const inputRef = ref(null)
const panelRef = ref(null)
const scrollRef = ref(null)
const recent = ref(JSON.parse(localStorage.getItem('uc_recent_searches') || '[]'))

const itemRefs = {}

const sections = computed(() => [
  { key: 'events',      label: t('public.search.sectionEvents'),      icon: CalendarDaysIcon,              iconBg: 'bg-[#edf4ff]',  iconColor: 'text-[#1a73e8]' },
  { key: 'prestations', label: t('public.search.sectionPrestations'), icon: WrenchScrewdriverIcon,         iconBg: 'bg-[#d1fae5]',  iconColor: 'text-[#006d35]' },
  { key: 'threads',     label: t('public.search.sectionForum'),       icon: ChatBubbleLeftEllipsisIcon,    iconBg: 'bg-[#fef3c7]',  iconColor: 'text-[#92400e]' },
  { key: 'providers',   label: t('public.search.sectionProviders'),   icon: BuildingStorefrontIcon,        iconBg: 'bg-[#f3e8ff]',  iconColor: 'text-[#7c3aed]' },
])

const quickActions = computed(() => [
  { label: t('public.layout.navEvents'),      sub: t('public.search.qaEvents'),      url: '/evenements',  icon: CalendarDaysIcon,           bg: '#edf4ff', color: '#1a73e8' },
  { label: t('public.layout.navPrestations'), sub: t('public.search.qaPrestations'), url: '/prestations', icon: WrenchScrewdriverIcon,      bg: '#d1fae5', color: '#006d35' },
  { label: t('public.layout.navCommunity'),   sub: t('public.search.qaForum'),       url: '/communaute',  icon: ChatBubbleLeftEllipsisIcon, bg: '#fef3c7', color: '#92400e' },
  { label: t('public.layout.btnDeposit'),     sub: t('public.search.qaDeposit'),     url: '/depot',       icon: BuildingStorefrontIcon,     bg: '#f3e8ff', color: '#7c3aed' },
])

const hasResults = computed(() =>
  sections.value.some(s => results.value[s.key]?.length > 0)
)

const flatResults = computed(() => {
  const flat = []
  for (const section of sections.value) {
    for (const item of (results.value[section.key] || [])) {
      flat.push(item)
    }
  }
  return flat
})

function flatIndex(sectionKey, j) {
  let idx = 0
  for (const section of sections.value) {
    if (section.key === sectionKey) return idx + j
    idx += results.value[section.key]?.length || 0
  }
  return -1
}

function setRef(el, sectionKey, j) {
  if (!itemRefs[sectionKey]) itemRefs[sectionKey] = []
  itemRefs[sectionKey][j] = el
}

let debounceTimer = null
watch(query, (val) => {
  focusedIndex.value = -1
  hoveredKey.value = ''
  clearTimeout(debounceTimer)
  if (val.trim().length < 2) {
    results.value = { events: [], prestations: [], threads: [], providers: [] }
    loading.value = false
    return
  }
  loading.value = true
  debounceTimer = setTimeout(() => doSearch(val.trim()), 300)
})

async function doSearch(q) {
  try {
    const res = await fetch(`/api/public/v1/search?q=${encodeURIComponent(q)}`).catch(() => null)
    if (res?.ok) {
      results.value = await res.json()
    }
  } finally {
    loading.value = false
  }
}

function highlight(text) {
  if (!query.value.trim()) return text
  const escaped = query.value.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  return text.replace(new RegExp(`(${escaped})`, 'gi'), '<mark class="bg-yellow-200 text-[#001d32] rounded px-0.5 not-italic">$1</mark>')
}

function selectItem(item) {
  saveRecent(query.value.trim())
  go(item.url)
}

function go(url) {
  close()
  router.push(url)
}

function saveRecent(q) {
  if (!q) return
  const updated = [q, ...recent.value.filter(r => r !== q)].slice(0, 5)
  recent.value = updated
  localStorage.setItem('uc_recent_searches', JSON.stringify(updated))
}
function removeRecent(i) {
  recent.value.splice(i, 1)
  localStorage.setItem('uc_recent_searches', JSON.stringify(recent.value))
}
function clearRecent() {
  recent.value = []
  localStorage.removeItem('uc_recent_searches')
}

function handleKey(e) {
  const total = flatResults.value.length
  if (e.key === 'Escape') { close(); return }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    focusedIndex.value = (focusedIndex.value + 1) % Math.max(total, 1)
    scrollToFocused()
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    focusedIndex.value = focusedIndex.value <= 0 ? total - 1 : focusedIndex.value - 1
    scrollToFocused()
  } else if (e.key === 'Enter' && focusedIndex.value >= 0) {
    e.preventDefault()
    const item = flatResults.value[focusedIndex.value]
    if (item) selectItem(item)
  }
}

function scrollToFocused() {
  nextTick(() => {
    let idx = 0
    for (const section of sections.value) {
      const sectionItems = results.value[section.key] || []
      for (let j = 0; j < sectionItems.length; j++) {
        if (idx === focusedIndex.value) {
          itemRefs[section.key]?.[j]?.scrollIntoView({ block: 'nearest' })
          return
        }
        idx++
      }
    }
  })
}

function open() {
  isOpen.value = true
  nextTick(() => inputRef.value?.focus())
}
function close() {
  isOpen.value = false
  query.value = ''
  results.value = { events: [], prestations: [], threads: [], providers: [] }
  focusedIndex.value = -1
  hoveredKey.value = ''
}

function onKeydown(e) {
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    isOpen.value ? close() : open()
  }
  if (e.key === 'Escape' && isOpen.value) close()
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onUnmounted(() => document.removeEventListener('keydown', onKeydown))

defineExpose({ open, close })
</script>

<style scoped>
.search-fade-enter-active { transition: opacity 0.15s ease; }
.search-fade-leave-active { transition: opacity 0.1s ease; }
.search-fade-enter-from, .search-fade-leave-to { opacity: 0; }

.search-fade-enter-active .search-panel {
  animation: panelIn 0.15s ease;
}
.search-fade-leave-active .search-panel {
  animation: panelOut 0.1s ease forwards;
}

@keyframes panelIn {
  from { transform: translateY(-8px) scale(0.98); opacity: 0; }
  to   { transform: translateY(0)    scale(1);    opacity: 1; }
}
@keyframes panelOut {
  from { transform: translateY(0)    scale(1);    opacity: 1; }
  to   { transform: translateY(-8px) scale(0.98); opacity: 0; }
}
</style>
