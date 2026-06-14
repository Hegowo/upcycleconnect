<template>
  <div class="relative" ref="rootRef">
    <button
      @click="toggle"
      :class="[
        'relative p-2 rounded-lg transition',
        isDark
          ? 'text-[#94a3b8] hover:text-[#f1f5f9] hover:bg-[#263244]'
          : 'text-gray-400 hover:text-gray-600 hover:bg-gray-50'
      ]"
      :title="'Notifications'"
    >
      <BellIcon class="w-5 h-5" />
      <span
        v-if="unreadCount > 0"
        class="absolute -top-0.5 -right-0.5 min-w-[18px] h-[18px] px-1 rounded-full bg-red-500 text-white text-[10px] font-bold flex items-center justify-center"
      >
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
    </button>

    <Transition name="dropdown">
      <div
        v-if="open"
        :class="[
          'absolute right-0 mt-2 w-80 sm:w-96 rounded-xl shadow-2xl border z-50 overflow-hidden',
          isDark ? 'bg-[#1e293b] border-[#334155]' : 'bg-white border-gray-200'
        ]"
      >
        <div :class="['px-4 py-3 flex items-center justify-between border-b', isDark ? 'border-[#334155]' : 'border-gray-100']">
          <h3 :class="['font-jakarta font-bold text-sm', isDark ? 'text-[#f1f5f9]' : 'text-gray-900']">Notifications</h3>
          <button
            v-if="unreadCount > 0"
            @click="markAll"
            :class="['text-xs font-medium hover:underline', isDark ? 'text-[#4ade80]' : 'text-[#006d35]']"
          >
            Tout marquer lu
          </button>
        </div>

        <div class="max-h-[440px] overflow-y-auto">
          <div v-if="loading" class="px-4 py-8 text-center">
            <div :class="['w-6 h-6 border-2 rounded-full animate-spin mx-auto', isDark ? 'border-[#4ade80] border-t-transparent' : 'border-[#006d35] border-t-transparent']" />
          </div>

          <div v-else-if="items.length === 0" :class="['px-4 py-10 text-center text-sm', isDark ? 'text-[#64748b]' : 'text-gray-400']">
            <BellSlashIcon class="w-8 h-8 mx-auto mb-2 opacity-40" />
            Aucune notification pour le moment.
          </div>

          <div v-else>
            <button
              v-for="n in items"
              :key="n.id"
              @click="onItemClick(n)"
              :class="[
                'w-full text-left px-4 py-3 border-b transition flex gap-3 items-start',
                isDark ? 'border-[#1e3352] hover:bg-[#263244]' : 'border-gray-50 hover:bg-gray-50',
                !n.read ? (isDark ? 'bg-[#006d35]/8' : 'bg-[#f0fdf4]/60') : ''
              ]"
            >
              <div
                class="w-8 h-8 rounded-full flex items-center justify-center shrink-0 mt-0.5"
                :style="{ backgroundColor: iconBg(n.type) }"
              >
                <component :is="iconFor(n.type)" class="w-4 h-4 text-white" />
              </div>
              <div class="flex-1 min-w-0">
                <p :class="['text-sm font-semibold leading-tight', isDark ? 'text-[#f1f5f9]' : 'text-gray-900']">
                  {{ n.title }}
                </p>
                <p :class="['text-xs leading-snug mt-0.5 line-clamp-2', isDark ? 'text-[#94a3b8]' : 'text-gray-500']">
                  {{ n.body }}
                </p>
                <p :class="['text-[10px] mt-1', isDark ? 'text-[#64748b]' : 'text-gray-400']">{{ formatTime(n.created_at) }}</p>
              </div>
              <span v-if="!n.read" class="w-2 h-2 rounded-full bg-red-500 shrink-0 mt-2" />
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import {
  BellIcon,
  BellSlashIcon,
  CheckCircleIcon,
  BanknotesIcon,
  DocumentTextIcon,
  ShieldCheckIcon,
  XCircleIcon,
  UserPlusIcon,
  BriefcaseIcon,
  PencilSquareIcon,
} from '@heroicons/vue/24/outline'
import { adminNotifications, userNotifications } from '@/services/notifications'

const props = defineProps({
  variant: { type: String, default: 'admin' },
  isDark:  { type: Boolean, default: false },
})

const router = useRouter()
const svc    = computed(() => props.variant === 'admin' ? adminNotifications : userNotifications)

const open        = ref(false)
const loading     = ref(false)
const items       = ref([])
const unreadCount = ref(0)
const rootRef     = ref(null)

let pollTimer = null
const POLL_MS = 60_000

async function refreshCount() {
  try { unreadCount.value = await svc.value.unreadCount() } catch {}
}

async function fetchList() {
  loading.value = true
  try {
    const data = await svc.value.list(30)
    items.value = data.data || []
    unreadCount.value = data.unread_count || 0
  } catch {
    items.value = []
  } finally {
    loading.value = false
  }
}

async function toggle() {
  open.value = !open.value
  if (open.value) await fetchList()
}

async function markAll() {
  await svc.value.markAllRead()
  items.value = items.value.map(n => ({ ...n, read: true }))
  unreadCount.value = 0
}

async function onItemClick(n) {
  if (!n.read) {
    try { await svc.value.markRead(n.id) } catch {}
    n.read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }
  if (n.link) {
    open.value = false
    const link = String(n.link).trim()

    if (/^(https?:)?\/\//i.test(link) || link.startsWith('mailto:') || link.startsWith('tel:')) {
      window.open(link, '_blank', 'noopener')
    } else {
      router.push(link.startsWith('/') ? link : `/${link}`)
    }
  }
}

function onDocClick(e) {
  if (rootRef.value && !rootRef.value.contains(e.target)) open.value = false
}

onMounted(() => {
  refreshCount()
  pollTimer = setInterval(refreshCount, POLL_MS)
  document.addEventListener('mousedown', onDocClick)
})
onUnmounted(() => {
  clearInterval(pollTimer)
  document.removeEventListener('mousedown', onDocClick)
})

function iconFor(type) {
  if (type.startsWith('payment'))    return BanknotesIcon
  if (type.startsWith('quote'))      return DocumentTextIcon
  if (type.startsWith('reservation')) return BriefcaseIcon
  if (type === 'account.banned' || type === 'provider.suspended') return XCircleIcon
  if (type === 'account.activated' || type === 'provider.approved') return CheckCircleIcon
  if (type === 'password.reset_requested') return ShieldCheckIcon
  if (type === 'provider.applied' || type === 'provider.application_received') return UserPlusIcon
  if (type === 'prestation.published') return PencilSquareIcon
  return BellIcon
}
function iconBg(type) {
  if (type.startsWith('payment'))    return '#006d35'
  if (type.startsWith('quote'))      return '#0284c7'
  if (type.startsWith('reservation')) return '#7c3aed'
  if (type === 'account.banned' || type === 'provider.suspended') return '#dc2626'
  if (type === 'account.activated' || type === 'provider.approved') return '#15803d'
  if (type === 'password.reset_requested') return '#d97706'
  if (type === 'prestation.published') return '#0d9488'
  return '#64748b'
}

function formatTime(iso) {
  if (!iso) return ''
  const d   = new Date(iso)
  const now = new Date()
  const diff = Math.floor((now - d) / 1000)
  if (diff < 60) return "À l'instant"
  if (diff < 3600) return `il y a ${Math.floor(diff/60)} min`
  if (diff < 86400) return `il y a ${Math.floor(diff/3600)} h`
  if (diff < 604800) return `il y a ${Math.floor(diff/86400)} j`
  return d.toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit' })
}

watch(() => router.currentRoute.value.fullPath, () => { open.value = false })
</script>

<style scoped>
.dropdown-enter-active, .dropdown-leave-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-4px); }
.line-clamp-2 { display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }
</style>
