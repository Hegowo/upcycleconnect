<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[860px] mx-auto px-4 sm:px-6 py-10">

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else-if="thread">
        <!-- Breadcrumb -->
        <div class="flex items-center gap-2 text-sm text-[#40617f] mb-6">
          <RouterLink to="/communaute" class="hover:text-[#006d35] transition">{{ t('public.community.title') }}</RouterLink>
          <ChevronRightIcon class="w-3.5 h-3.5" />
          <RouterLink :to="`/communaute/forum/${thread.category?.slug}`" class="hover:text-[#006d35] transition">
            {{ thread.category?.name }}
          </RouterLink>
          <ChevronRightIcon class="w-3.5 h-3.5" />
          <span class="text-[#001d32] font-medium truncate max-w-[200px]">{{ thread.title }}</span>
        </div>

        <!-- Thread header -->
        <div class="bg-white rounded-2xl p-6 shadow-sm mb-6">
          <div class="flex items-start gap-4">
            <div class="w-11 h-11 rounded-full bg-[#edf4ff] flex items-center justify-center shrink-0 text-base font-bold text-[#40617f]">
              {{ authorInitial(thread.author) }}
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap mb-2">
                <span v-if="thread.pinned"
                  class="text-xs font-bold text-[#006d35] bg-[#d1fae5] px-2 py-0.5 rounded-full flex items-center gap-1">
                  {{ t('public.community.pinned') }}
                </span>
                <span v-if="thread.locked"
                  class="text-xs font-bold text-[#ef4444] bg-[#fee2e2] px-2 py-0.5 rounded-full flex items-center gap-1">
                  <LockClosedIcon class="w-3 h-3" />
                  {{ t('public.community.locked') }}
                </span>
              </div>
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-xl sm:text-2xl mb-3">{{ thread.title }}</h1>
              <div class="prose prose-sm max-w-none text-[#374151] whitespace-pre-wrap leading-relaxed">{{ thread.content }}</div>
            </div>
          </div>

          <div class="flex items-center justify-between mt-5 pt-4 border-t border-[#f1f5f9]">
            <div class="flex items-center gap-2">
              <span class="text-sm font-semibold text-[#001d32]">
                {{ thread.author?.first_name || thread.author?.email?.split('@')[0] }}
              </span>
              <span class="text-[#94a3b8] text-xs">· {{ formatDate(thread.created_at) }}</span>
            </div>
            <button v-if="canDeleteThread"
              @click="deleteThread"
              class="text-xs text-[#ef4444] hover:underline flex items-center gap-1">
              <TrashIcon class="w-3.5 h-3.5" />
              {{ t('public.community.delete') }}
            </button>
          </div>
        </div>

        <!-- Replies -->
        <div class="mb-6">
          <p class="font-jakarta font-bold text-[#001d32] text-lg mb-4">
            {{ replies.length }} {{ t('public.community.repliesTitle') }}
          </p>

          <div class="flex flex-col gap-4">
            <div v-for="reply in replies" :key="reply.id"
              class="bg-white rounded-xl p-5 shadow-sm flex gap-4 items-start">
              <div class="w-9 h-9 rounded-full bg-[#edf4ff] flex items-center justify-center shrink-0 text-sm font-bold text-[#40617f]">
                {{ authorInitial(reply.author) }}
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center justify-between mb-2">
                  <div class="flex items-center gap-2">
                    <span class="text-sm font-semibold text-[#001d32]">
                      {{ reply.author?.first_name || reply.author?.email?.split('@')[0] }}
                    </span>
                    <span class="text-[#94a3b8] text-xs">· {{ formatDate(reply.created_at) }}</span>
                  </div>
                  <button v-if="canDeleteReply(reply)"
                    @click="deleteReply(reply.id)"
                    class="text-xs text-[#ef4444] hover:underline flex items-center gap-1">
                    <TrashIcon class="w-3.5 h-3.5" />
                  </button>
                </div>
                <div class="text-sm text-[#374151] whitespace-pre-wrap leading-relaxed">{{ reply.content }}</div>
              </div>
            </div>

            <div v-if="replies.length === 0" class="text-center py-8 text-[#40617f] text-sm">
              {{ t('public.community.noReplies') }}
            </div>
          </div>
        </div>

        <!-- Reply form -->
        <div class="bg-white rounded-2xl p-6 shadow-sm">
          <h3 class="font-jakarta font-bold text-[#001d32] mb-4">{{ t('public.community.replyTitle') }}</h3>

          <template v-if="!userAuth.isLoggedIn">
            <p class="text-[#40617f] text-sm">
              {{ t('public.community.loginToReply') }}
              <RouterLink to="/connexion" class="text-[#006d35] font-semibold hover:underline">{{ t('public.community.loginLink') }}</RouterLink>
            </p>
          </template>

          <template v-else-if="thread.locked">
            <p class="text-[#94a3b8] text-sm flex items-center gap-2">
              <LockClosedIcon class="w-4 h-4" />
              {{ t('public.community.threadLocked') }}
            </p>
          </template>

          <template v-else>
            <textarea v-model="replyContent" rows="4"
              :placeholder="t('public.community.replyPlaceholder')"
              class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none mb-3" />
            <p v-if="replyError" class="text-red-500 text-xs mb-3">{{ replyError }}</p>
            <div class="flex justify-end">
              <button @click="submitReply" :disabled="replySubmitting"
                class="px-5 py-2.5 bg-[#006d35] text-white text-sm font-semibold rounded-xl hover:bg-[#1b8848] transition disabled:opacity-50">
                {{ replySubmitting ? t('common.loading') : t('public.community.replySend') }}
              </button>
            </div>
          </template>
        </div>
      </template>

      <!-- Not found -->
      <div v-else class="text-center py-24 text-[#40617f]">
        <p>{{ t('public.community.threadNotFound') }}</p>
        <RouterLink to="/communaute" class="text-[#006d35] font-semibold hover:underline text-sm mt-2 block">
          {{ t('public.community.backToCommunity') }}
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ChevronRightIcon,
  LockClosedIcon,
  TrashIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const loading = ref(true)
const thread = ref(null)
const replies = ref([])
const replyContent = ref('')
const replyError = ref('')
const replySubmitting = ref(false)

const threadId = computed(() => route.params.id)

const canDeleteThread = computed(() => {
  if (!userAuth.isLoggedIn || !thread.value) return false
  return thread.value.author?.id === userAuth.user?.id
})

function canDeleteReply(reply) {
  if (!userAuth.isLoggedIn) return false
  return reply.author?.id === userAuth.user?.id
}

function authorInitial(author) {
  if (!author) return '?'
  const n = author.first_name || author.email || '?'
  return n[0].toUpperCase()
}

function formatDate(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const now = new Date()
  const diff = now - d
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return t('public.community.justNow')
  if (mins < 60) return t('public.community.minutesAgo', { n: mins })
  const hours = Math.floor(mins / 60)
  if (hours < 24) return t('public.community.hoursAgo', { n: hours })
  const days = Math.floor(hours / 24)
  if (days < 7) return t('public.community.daysAgo', { n: days })
  const loc = locale.value === 'en' ? 'en-US' : 'fr-FR'
  return d.toLocaleDateString(loc, { day: 'numeric', month: 'short', year: 'numeric' })
}

async function loadThread() {
  const res = await fetch(`/api/public/v1/forum/threads/${threadId.value}`).catch(() => null)
  if (res?.ok) {
    const data = await res.json()
    thread.value = data.thread
    replies.value = data.replies || []
  }
}

async function submitReply() {
  replyError.value = ''
  if (replyContent.value.trim().length < 2) {
    replyError.value = t('public.community.errorReplyContent')
    return
  }
  replySubmitting.value = true
  try {
    const res = await fetch(`/api/v1/forum/threads/${threadId.value}/replies`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({ content: replyContent.value.trim() }),
    })
    if (res.ok) {
      const data = await res.json()
      replies.value.push(data.reply)
      replyContent.value = ''
    } else {
      replyError.value = t('public.community.errorSubmit')
    }
  } finally {
    replySubmitting.value = false
  }
}

async function deleteThread() {
  if (!confirm(t('public.community.confirmDelete'))) return
  const res = await fetch(`/api/v1/forum/threads/${threadId.value}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  if (res.ok) {
    router.push(`/communaute/forum/${thread.value.category?.slug}`)
  }
}

async function deleteReply(replyId) {
  if (!confirm(t('public.community.confirmDelete'))) return
  const res = await fetch(`/api/v1/forum/replies/${replyId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  if (res.ok) {
    replies.value = replies.value.filter(r => r.id !== replyId)
  }
}

onMounted(async () => {
  await loadThread()
  loading.value = false
})
</script>
