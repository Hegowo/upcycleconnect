<template>
  <div class="bg-[#f7f9ff] pb-16">
    <div class="max-w-[1280px] mx-auto px-6 py-10">

      <div class="flex items-center gap-2 text-sm text-[#40617f] mb-8">
        <RouterLink to="/evenements" class="hover:text-[#006d35] transition">{{ t('public.eventDetail.breadcrumbEvents') }}</RouterLink>
        <ChevronRightIcon class="w-3.5 h-3.5" />
        <span class="text-[#001d32] font-medium truncate max-w-[300px]">{{ event?.title || '…' }}</span>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-32">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <!-- Not found -->
      <div v-else-if="!event" class="flex flex-col items-center justify-center py-32 text-center">
        <CalendarDaysIcon class="w-16 h-16 text-[#40617f]/30 mb-4" />
        <p class="font-jakarta font-bold text-[#001d32] text-xl mb-2">{{ t('public.eventDetail.notFound') }}</p>
        <RouterLink to="/evenements" class="text-[#006d35] font-semibold hover:underline text-sm mt-2">
          {{ t('public.eventDetail.back') }}
        </RouterLink>
      </div>

      <template v-else>
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">

          <!-- Main content -->
          <div class="lg:col-span-8 flex flex-col gap-8">

            <!-- Hero -->
            <div class="relative h-[360px] sm:h-[440px] rounded-[32px] overflow-hidden flex items-center justify-center shadow-sm"
              :style="`background: ${cardGradient(event.id)};`">
              <CalendarDaysIcon class="w-40 h-40 text-white/20" />
              <div class="absolute top-6 left-6 flex gap-2 flex-wrap">
                <span v-if="event.category"
                  class="bg-white/90 backdrop-blur-sm font-bold text-sm px-4 py-2 rounded-full shadow-sm"
                  :class="categoryColor(event.category.name)">
                  {{ event.category.name }}
                </span>
                <span v-if="isPast" class="bg-gray-700/80 text-white text-sm font-bold px-4 py-2 rounded-full">
                  {{ t('public.eventDetail.ctaPast') }}
                </span>
              </div>
              <div class="absolute top-6 right-6">
                <span v-if="isFull && !isPast" class="bg-red-500 text-white text-sm font-bold px-4 py-2 rounded-full">
                  {{ t('public.eventDetail.badgeFull') }}
                </span>
                <span v-else-if="!isPast && spotsLeft !== null" class="bg-[#006d35] text-white text-sm font-bold px-4 py-2 rounded-full">
                  {{ spotsLeft }} {{ spotsLeft === 1 ? 'place' : 'places' }}
                </span>
              </div>
            </div>

            <!-- Title + meta -->
            <div class="flex flex-col gap-3">
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl leading-tight tracking-tight">
                {{ event.title }}
              </h1>
              <div class="flex flex-wrap gap-5 text-[#40617f] text-base mt-1">
                <div class="flex items-center gap-2">
                  <CalendarDaysIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ formatDate(event.start_date) }}</span>
                </div>
                <div v-if="sameDay" class="flex items-center gap-2">
                  <ClockIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ formatTime(event.start_date) }} – {{ formatTime(event.end_date) }}</span>
                </div>
                <div v-else class="flex items-center gap-2">
                  <ClockIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ formatDate(event.start_date) }} – {{ formatDate(event.end_date) }}</span>
                </div>
                <div v-if="event.location" class="flex items-center gap-2">
                  <MapPinIcon class="w-5 h-5 shrink-0 text-[#006d35]" />
                  <span>{{ event.location }}</span>
                </div>
              </div>
            </div>

            <!-- Description -->
            <div v-if="event.description" class="bg-white rounded-[32px] p-10 shadow-sm flex flex-col gap-4">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.eventDetail.aboutTitle') }}</h2>
              <p class="text-[#40617f] text-lg leading-relaxed whitespace-pre-wrap">{{ event.description }}</p>
            </div>

            <!-- Organizer -->
            <div v-if="event.creator" class="bg-[#edf4ff] rounded-[32px] p-8 flex gap-6 items-start">
              <div class="w-20 h-20 rounded-2xl overflow-hidden bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] flex items-center justify-center shrink-0 shadow-md">
                <img v-if="event.creator.avatar_url" :src="event.creator.avatar_url" class="w-full h-full object-cover" alt="" />
                <UserCircleIcon v-else class="w-12 h-12 text-[#006d35]/40" />
              </div>
              <div class="flex flex-col gap-1">
                <p class="text-[#006d35] text-xs font-bold uppercase tracking-widest">{{ t('public.eventDetail.organizerTag') }}</p>
                <h3 class="font-jakarta font-bold text-[#001d32] text-xl">
                  {{ event.creator.first_name }} {{ event.creator.last_name }}
                </h3>
              </div>
            </div>

          </div>

          <!-- Sidebar -->
          <div class="lg:col-span-4 flex flex-col gap-6">

            <!-- Registration card -->
            <div class="bg-white border border-[#d8eaff] rounded-[40px] p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] relative overflow-hidden sticky top-6">
              <div class="absolute -right-8 -top-8 w-32 h-32 rounded-full bg-[rgba(0,109,53,0.05)]" />

              <div class="flex items-baseline justify-between mb-2">
                <span class="font-jakarta font-extrabold text-[#001d32] text-4xl">{{ t('public.eventDetail.priceFree') }}</span>
              </div>
              <p class="text-[#40617f] text-sm mb-6">{{ t('public.eventDetail.freeLabel') }}</p>

              <div v-if="!isPast" class="mb-6">
                <div v-if="event.max_participants" class="flex items-center gap-2">
                  <div class="flex-1 bg-gray-100 rounded-full h-2 overflow-hidden">
                    <div class="h-full bg-[#006d35] rounded-full transition-all"
                      :style="`width: ${Math.min(100, (event.registrations_count / event.max_participants) * 100)}%`" />
                  </div>
                  <span class="text-sm font-semibold text-[#40617f] shrink-0">
                    {{ event.registrations_count }}/{{ event.max_participants }}
                  </span>
                </div>
                <p v-else class="text-sm text-[#40617f]">{{ t('public.eventDetail.spotsUnlimited') }}</p>
              </div>

              <!-- Organizer badge -->
              <div v-if="isOrganizer" class="flex items-center gap-2 bg-[#d1fae5] text-[#065f46] px-4 py-2 rounded-xl mb-4 text-sm font-semibold">
                <CheckBadgeIcon class="w-4 h-4 shrink-0" />
                {{ t('public.planning.organizerBadge') }}
              </div>

              <!-- Registered badge -->
              <div v-else-if="isRegistered && !isPast" class="flex items-center gap-2 bg-[#d1fae5] text-[#065f46] px-4 py-2 rounded-xl mb-4 text-sm font-semibold">
                <CheckCircleIcon class="w-4 h-4 shrink-0" />
                {{ t('public.eventDetail.badgeRegistered') }}
              </div>

              <!-- CTA -->
              <button v-if="isOrganizer" disabled
                class="w-full py-4 rounded-2xl text-[#065f46] font-bold text-base bg-[#d1fae5] cursor-default">
                {{ t('public.planning.organizerBadge') }}
              </button>
              <button v-else-if="isPast" disabled
                class="w-full py-4 rounded-2xl text-[#40617f] font-bold text-base bg-gray-100 cursor-default">
                {{ t('public.eventDetail.ctaPast') }}
              </button>
              <button v-else-if="!userAuth.isLoggedIn"
                @click="router.push(`/connexion?redirect=/evenements/${event.id}`)"
                class="w-full py-4 rounded-2xl text-white font-bold text-base flex items-center justify-center gap-2 transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,109,53,0.2)]"
                style="background: linear-gradient(168deg, #006d35 0%, #1b8848 100%);">
                {{ t('public.eventDetail.ctaLogin') }}
                <ArrowRightIcon class="w-4 h-4" />
              </button>
              <button v-else-if="isRegistered"
                @click="unregister" :disabled="actionLoading"
                class="w-full py-4 rounded-2xl font-bold text-base transition border-2 border-red-300 text-red-600 hover:bg-red-50 disabled:opacity-60">
                {{ actionLoading ? '…' : t('public.eventDetail.ctaUnregister') }}
              </button>
              <button v-else-if="isFull" disabled
                class="w-full py-4 rounded-2xl text-[#40617f] font-bold text-base bg-gray-100 cursor-default">
                {{ t('public.eventDetail.ctaFull') }}
              </button>
              <button v-else
                @click="register" :disabled="actionLoading"
                class="w-full py-4 rounded-2xl text-white font-bold text-base flex items-center justify-center gap-2 transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,109,53,0.2)] disabled:opacity-60"
                style="background: linear-gradient(168deg, #006d35 0%, #1b8848 100%);">
                {{ actionLoading ? t('public.eventDetail.registering') : t('public.eventDetail.ctaReserve') }}
                <ArrowRightIcon v-if="!actionLoading" class="w-4 h-4" />
              </button>

              <p v-if="feedback" class="mt-3 text-center text-sm font-medium"
                :class="feedbackError ? 'text-red-600' : 'text-[#006d35]'">
                {{ feedback }}
              </p>
            </div>

            <!-- Location card -->
            <div v-if="event.location" class="bg-white border border-[#d8eaff] rounded-[32px] p-6 flex items-center justify-between gap-4">
              <div>
                <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ event.location }}</p>
              </div>
              <div class="w-12 h-12 bg-[#d8eaff] rounded-xl flex items-center justify-center shrink-0">
                <MapPinIcon class="w-5 h-5 text-[#40617f]" />
              </div>
            </div>

          </div>
        </div>

        <!-- ── Chat section ─────────────────────────────────── -->
        <div class="mt-10">

          <!-- Not logged in -->
          <div v-if="!userAuth.isLoggedIn" class="bg-white rounded-[32px] p-8 shadow-sm flex items-center gap-4">
            <ChatBubbleLeftRightIcon class="w-8 h-8 text-[#40617f]/40 shrink-0" />
            <div>
              <p class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.eventDetail.chatTitle') }}</p>
              <p class="text-[#40617f] text-sm mt-1">{{ t('public.eventDetail.chatLoginHint') }}</p>
            </div>
          </div>

          <!-- Not registered (and not organizer) -->
          <div v-else-if="!canChat" class="bg-white rounded-[32px] p-8 shadow-sm flex items-center gap-4">
            <ChatBubbleLeftRightIcon class="w-8 h-8 text-[#40617f]/40 shrink-0" />
            <div>
              <p class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.eventDetail.chatTitle') }}</p>
              <p class="text-[#40617f] text-sm mt-1">{{ t('public.eventDetail.chatRegisterHint') }}</p>
            </div>
          </div>

          <!-- Chat UI -->
          <div v-else class="bg-white rounded-[32px] shadow-sm overflow-hidden">
            <!-- Header -->
            <div class="px-8 py-5 border-b border-gray-100 flex items-center gap-3">
              <ChatBubbleLeftRightIcon class="w-5 h-5 text-[#006d35]" />
              <h2 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.eventDetail.chatTitle') }}</h2>
              <span v-if="isOrganizer"
                class="ml-auto text-xs font-bold px-2.5 py-1 rounded-full bg-[#d1fae5] text-[#065f46] flex items-center gap-1">
                <CheckBadgeIcon class="w-3 h-3" />
                {{ t('public.eventDetail.chatOrganizerBadge') }}
              </span>
            </div>

            <!-- Messages -->
            <div ref="messagesEl" class="h-[420px] overflow-y-auto px-6 py-4 flex flex-col gap-4 scroll-smooth">
              <div v-if="!messages.length" class="flex-1 flex items-center justify-center text-[#40617f] text-sm">
                {{ t('public.eventDetail.chatEmpty') }}
              </div>

              <div v-for="msg in messages" :key="msg.id"
                :class="['flex gap-3', msg.is_me ? 'flex-row-reverse' : 'flex-row']">

                <!-- Avatar -->
                <div class="shrink-0 w-9 h-9 rounded-full flex items-center justify-center text-xs font-bold text-white"
                  :style="msg.is_me ? 'background: linear-gradient(135deg,#006d35,#1b8848)' : 'background: #40617f'">
                  <img v-if="msg.avatar" :src="msg.avatar" class="w-full h-full object-cover rounded-full" alt="" />
                  <span v-else>{{ initials(msg.user_name) }}</span>
                </div>

                <!-- Bubble -->
                <div :class="['max-w-[70%] flex flex-col gap-1', msg.is_me ? 'items-end' : 'items-start']">
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-semibold text-[#40617f]">
                      {{ msg.is_me ? 'Vous' : msg.user_name }}
                    </span>
                    <span v-if="msg.is_organizer"
                      class="text-[10px] font-bold px-1.5 py-0.5 rounded-full bg-[#d1fae5] text-[#065f46] flex items-center gap-0.5">
                      <CheckBadgeIcon class="w-2.5 h-2.5" />
                      {{ t('public.eventDetail.chatOrganizerBadge') }}
                    </span>
                    <span class="text-[10px] text-gray-400">{{ formatMsgTime(msg.created_at) }}</span>
                  </div>

                  <!-- Text -->
                  <div v-if="msg.content"
                    :class="[
                      'px-4 py-2.5 rounded-2xl text-sm leading-relaxed whitespace-pre-wrap break-words',
                      msg.is_me
                        ? 'bg-[#006d35] text-white rounded-tr-sm'
                        : 'bg-[#f1f5f9] text-[#001d32] rounded-tl-sm'
                    ]">
                    {{ msg.content }}
                  </div>

                  <!-- Image -->
                  <img v-if="msg.image_url" :src="msg.image_url"
                    :alt="t('public.eventDetail.chatImageAlt')"
                    class="max-w-[280px] rounded-2xl object-cover cursor-pointer hover:opacity-90 transition"
                    @click="lightboxSrc = msg.image_url" />
                </div>
              </div>
            </div>

            <!-- Input -->
            <div class="border-t border-gray-100 px-6 py-4 flex gap-3 items-end">
              <input type="file" ref="fileInput" accept="image/*" class="hidden" @change="onFileChange" />

              <button @click="fileInput?.click()"
                :disabled="sendLoading"
                class="shrink-0 w-10 h-10 rounded-xl bg-[#edf4ff] text-[#1a73e8] flex items-center justify-center hover:bg-[#d0e8ff] transition disabled:opacity-40"
                :title="t('public.eventDetail.chatImageUpload')">
                <PhotoIcon class="w-5 h-5" />
              </button>

              <textarea
                v-model="chatInput"
                @keydown.enter.exact.prevent="sendMessage"
                :placeholder="t('public.eventDetail.chatPlaceholder')"
                rows="1"
                class="flex-1 resize-none bg-[#f7f9ff] border border-gray-200 rounded-xl px-4 py-2.5 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-gray-400 max-h-32 overflow-auto"
                style="field-sizing: content;"
              />

              <button @click="sendMessage"
                :disabled="sendLoading || (!chatInput.trim() && !pendingImage)"
                class="shrink-0 w-10 h-10 rounded-xl flex items-center justify-center transition disabled:opacity-40"
                style="background: linear-gradient(135deg,#006d35,#1b8848)">
                <PaperAirplaneIcon class="w-4 h-4 text-white" />
              </button>
            </div>

            <!-- Pending image preview -->
            <div v-if="pendingImagePreview" class="px-6 pb-4 flex items-center gap-3">
              <img :src="pendingImagePreview" class="h-16 w-16 object-cover rounded-xl border border-gray-200" alt="" />
              <button @click="clearPendingImage" class="text-xs text-red-500 hover:underline">Supprimer</button>
            </div>
          </div>
        </div>

        <!-- Lightbox -->
        <Teleport to="body">
          <div v-if="lightboxSrc" @click="lightboxSrc = null"
            class="fixed inset-0 z-50 bg-black/80 flex items-center justify-center p-4 cursor-zoom-out">
            <img :src="lightboxSrc" class="max-w-full max-h-full rounded-2xl object-contain" alt="" />
          </div>
        </Teleport>

      </template>

    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ChevronRightIcon,
  CalendarDaysIcon,
  ClockIcon,
  MapPinIcon,
  UserCircleIcon,
  ArrowRightIcon,
  CheckCircleIcon,
  CheckBadgeIcon,
  ChatBubbleLeftRightIcon,
  PhotoIcon,
  PaperAirplaneIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()
const userAuth = useUserAuthStore()

// Event
const loading = ref(true)
const event = ref(null)
const isRegistered = ref(false)
const actionLoading = ref(false)
const feedback = ref('')
const feedbackError = ref(false)

// Chat
const messages = ref([])
const chatInput = ref('')
const sendLoading = ref(false)
const messagesEl = ref(null)
const fileInput = ref(null)
const pendingImage = ref(null)      // URL returned by upload
const pendingImagePreview = ref(null) // local blob preview
const lightboxSrc = ref(null)
let pollTimer = null

const GRADIENTS = [
  ['#dbeafe', '#bfdbfe'], ['#f3e8ff', '#e9d5ff'],
  ['#d1fae5', '#a7f3d0'], ['#fef3c7', '#fde68a'],
  ['#fee2e2', '#fecaca'], ['#e0e7ff', '#c7d2fe'],
]

function cardGradient(id) {
  const [from, to] = GRADIENTS[id % GRADIENTS.length]
  return `linear-gradient(135deg, ${from}, ${to})`
}

function categoryColor(name = '') {
  const n = name.toLowerCase()
  if (n.includes('atelier') || n.includes('workshop')) return 'text-[#1d4ed8]'
  if (n.includes('formation') || n.includes('training')) return 'text-[#065f46]'
  if (n.includes('conf')) return 'text-[#7c3aed]'
  return 'text-[#40617f]'
}

function initials(name = '') {
  return name.split(' ').map(w => w[0]).join('').slice(0, 2).toUpperCase()
}

const localeCode = computed(() => locale.value === 'en' ? 'en-US' : 'fr-FR')

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode.value, { day: 'numeric', month: 'long', year: 'numeric' })
}

function formatTime(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleTimeString(localeCode.value, { hour: '2-digit', minute: '2-digit' })
}

function formatMsgTime(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const today = new Date()
  if (d.toDateString() === today.toDateString()) {
    return d.toLocaleTimeString(localeCode.value, { hour: '2-digit', minute: '2-digit' })
  }
  return d.toLocaleDateString(localeCode.value, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}

const sameDay = computed(() => {
  if (!event.value) return false
  return new Date(event.value.start_date).toDateString() === new Date(event.value.end_date).toDateString()
})

const isPast = computed(() => event.value ? new Date(event.value.end_date) < new Date() : false)

const isFull = computed(() => {
  if (!event.value?.max_participants) return false
  return event.value.registrations_count >= event.value.max_participants
})

const spotsLeft = computed(() => {
  if (!event.value?.max_participants) return null
  return Math.max(0, event.value.max_participants - event.value.registrations_count)
})

const isOrganizer = computed(() =>
  !!(event.value?.creator && userAuth.user && event.value.creator.id === userAuth.user.id)
)

const canChat = computed(() => isOrganizer.value || isRegistered.value)

function showFeedback(msg, isError = false) {
  feedback.value = msg
  feedbackError.value = isError
  setTimeout(() => { feedback.value = '' }, 4000)
}

async function register() {
  actionLoading.value = true
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/register`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      isRegistered.value = true
      event.value.registrations_count++
      showFeedback(t('public.eventDetail.registerSuccess'))
      startPolling()
    } else {
      const data = await res.json().catch(() => ({}))
      showFeedback(data.message || t('public.eventDetail.registerError'), true)
    }
  } catch {
    showFeedback(t('public.eventDetail.registerError'), true)
  }
  actionLoading.value = false
}

async function unregister() {
  if (!confirm(t('public.eventDetail.unregisterConfirm'))) return
  actionLoading.value = true
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/register`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      isRegistered.value = false
      event.value.registrations_count = Math.max(0, event.value.registrations_count - 1)
      showFeedback(t('public.eventDetail.unregisterSuccess'))
      stopPolling()
      messages.value = []
    }
  } catch {}
  actionLoading.value = false
}

// ─── Chat ────────────────────────────────────────────────
async function fetchMessages(scrollToBottom = false) {
  if (!canChat.value || !userAuth.token) return
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/messages`, {
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      const data = await res.json()
      const newMsgs = data.data || []
      // Only update + scroll if there are new messages
      if (newMsgs.length !== messages.value.length || scrollToBottom) {
        messages.value = newMsgs
        if (scrollToBottom || newMsgs.length > messages.value.length) {
          await nextTick()
          if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight
        }
      }
    }
  } catch {}
}

function startPolling() {
  stopPolling()
  pollTimer = setInterval(() => fetchMessages(false), 5000)
}

function stopPolling() {
  if (pollTimer) { clearInterval(pollTimer); pollTimer = null }
}

async function sendMessage() {
  if (sendLoading.value) return
  const content = chatInput.value.trim()
  if (!content && !pendingImage.value) return

  sendLoading.value = true
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/messages`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({ content, image_url: pendingImage.value || undefined }),
    })
    if (res.ok) {
      const msg = await res.json()
      messages.value.push(msg)
      chatInput.value = ''
      clearPendingImage()
      await nextTick()
      if (messagesEl.value) messagesEl.value.scrollTop = messagesEl.value.scrollHeight
    }
  } catch {}
  sendLoading.value = false
}

async function onFileChange(e) {
  const file = e.target.files?.[0]
  if (!file) return

  // Local preview
  pendingImagePreview.value = URL.createObjectURL(file)

  // Upload
  const fd = new FormData()
  fd.append('image', file)
  try {
    const res = await fetch(`/api/v1/events/${event.value.id}/messages/image`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${userAuth.token}` },
      body: fd,
    })
    if (res.ok) {
      const data = await res.json()
      pendingImage.value = data.url
    } else {
      clearPendingImage()
    }
  } catch {
    clearPendingImage()
  }
  // Reset input so same file can be re-selected
  if (fileInput.value) fileInput.value.value = ''
}

function clearPendingImage() {
  if (pendingImagePreview.value) URL.revokeObjectURL(pendingImagePreview.value)
  pendingImage.value = null
  pendingImagePreview.value = null
}

onMounted(async () => {
  const id = route.params.id
  try {
    const res = await fetch(`/api/public/v1/events/${id}`)
    if (res.ok) event.value = await res.json()
  } catch {}

  if (userAuth.isLoggedIn && event.value) {
    try {
      const res = await fetch(`/api/v1/events/${id}/registration`, {
        headers: { Authorization: `Bearer ${userAuth.token}` },
      })
      if (res.ok) isRegistered.value = (await res.json()).registered
    } catch {}
  }

  loading.value = false

  // Start chat if user has access
  if (canChat.value) {
    await fetchMessages(true)
    startPolling()
  }
})

onUnmounted(() => {
  stopPolling()
  if (pendingImagePreview.value) URL.revokeObjectURL(pendingImagePreview.value)
})
</script>
