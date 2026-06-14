<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[860px] mx-auto px-4 sm:px-6 py-10">

      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else-if="thread">

        <div class="flex items-center gap-2 text-sm text-[#40617f] mb-6 flex-wrap">
          <RouterLink to="/communaute" class="hover:text-[#006d35] transition">{{ t('public.community.title') }}</RouterLink>
          <ChevronRightIcon class="w-3.5 h-3.5 shrink-0" />
          <RouterLink :to="`/communaute/forum/${thread.category?.slug}`" class="hover:text-[#006d35] transition">
            {{ thread.category?.name }}
          </RouterLink>
          <ChevronRightIcon class="w-3.5 h-3.5 shrink-0" />
          <span class="text-[#001d32] font-medium truncate max-w-[180px]">{{ thread.title }}</span>
        </div>

        <div class="bg-white rounded-2xl shadow-sm mb-6 overflow-hidden">

          <div class="flex items-center gap-2 px-6 pt-5 pb-3 flex-wrap">
            <span v-if="thread.pinned" class="text-xs font-bold text-[#006d35] bg-[#d1fae5] px-2 py-0.5 rounded-full">
              📌 {{ t('public.community.pinned') }}
            </span>
            <span v-if="thread.locked" class="text-xs font-bold text-[#ef4444] bg-[#fee2e2] px-2 py-0.5 rounded-full flex items-center gap-1">
              <LockClosedIcon class="w-3 h-3" />{{ t('public.community.locked') }}
            </span>
          </div>

          <div class="px-6 pb-5">
            <h1 class="font-jakarta font-extrabold text-[#001d32] text-xl sm:text-2xl mb-5">{{ thread.title }}</h1>

            <div class="flex items-start gap-4">

              <div class="w-11 h-11 rounded-full flex items-center justify-center shrink-0 text-base font-bold text-white"
                :style="`background: ${avatarColor(thread.author)}`">
                {{ authorInitial(thread.author) }}
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-3">
                  <span class="font-semibold text-[#001d32] text-sm">
                    {{ authorName(thread.author) }}
                  </span>
                  <span class="text-[#94a3b8] text-xs">· {{ formatDate(thread.created_at) }}</span>
                  <span v-if="thread.updated_at !== thread.created_at" class="text-[#94a3b8] text-xs">(modifié)</span>
                </div>

                <div v-if="!editingThread" class="forum-content" v-html="thread.content" />

                <div v-else class="mb-3">
                  <TiptapEditor v-model="editThreadContent" :placeholder="t('public.community.modalContentPlaceholder')" minHeight="150px" />
                  <div class="flex gap-2 mt-2 justify-end">
                    <button @click="editingThread = false" class="px-3 py-1.5 text-xs text-[#40617f] font-medium hover:text-[#001d32]">
                      {{ t('common.cancel') }}
                    </button>
                    <button @click="saveThreadEdit" :disabled="saving"
                      class="px-4 py-1.5 bg-[#006d35] text-white text-xs font-semibold rounded-lg hover:bg-[#1b8848] transition disabled:opacity-50">
                      {{ saving ? t('common.loading') : t('common.save') }}
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-3 mt-5 pt-4 border-t border-[#f1f5f9] flex-wrap">

              <template v-if="isOwner && !editingThread">
                <button @click="startEditThread"
                  class="flex items-center gap-1 text-xs text-[#40617f] hover:text-[#006d35] transition font-medium">
                  <PencilSquareIcon class="w-3.5 h-3.5" />
                  {{ t('public.community.edit') }}
                </button>
                <button @click="toggleClose"
                  :class="['flex items-center gap-1 text-xs font-medium transition', thread.locked ? 'text-[#ef4444] hover:text-[#dc2626]' : 'text-[#40617f] hover:text-[#ef4444]']">
                  <LockClosedIcon class="w-3.5 h-3.5" />
                  {{ thread.locked ? t('public.community.unlock') : t('public.community.close') }}
                </button>
                <button @click="deleteThread"
                  class="flex items-center gap-1 text-xs text-[#40617f] hover:text-[#ef4444] transition font-medium">
                  <TrashIcon class="w-3.5 h-3.5" />
                  {{ t('public.community.delete') }}
                </button>
              </template>

              <button v-if="userAuth.isLoggedIn && !isOwner" @click="report('thread', thread.id)"
                class="flex items-center gap-1 text-xs text-[#94a3b8] hover:text-[#ef4444] transition font-medium ml-auto">
                <FlagIcon class="w-3.5 h-3.5" />
                {{ t('public.community.report') }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="isOwner && bans.length > 0" class="bg-[#fef2f2] border border-[#fecaca] rounded-xl p-4 mb-6">
          <p class="text-sm font-semibold text-[#991b1b] mb-3 flex items-center gap-2">
            <ShieldExclamationIcon class="w-4 h-4" />
            {{ t('public.community.bannedUsers') }} ({{ bans.length }})
          </p>
          <div class="flex flex-wrap gap-2">
            <div v-for="ban in bans" :key="ban.id"
              class="flex items-center gap-2 bg-white rounded-lg px-3 py-1.5 border border-[#fecaca]">
              <span class="text-sm text-[#991b1b] font-medium">{{ authorName(ban.user) }}</span>
              <button @click="unbanUser(ban.user_id)"
                class="text-[#94a3b8] hover:text-[#ef4444] transition">
                <XMarkIcon class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>

        <div class="mb-6">
          <p class="font-jakarta font-bold text-[#001d32] text-lg mb-4">
            {{ replies.length }} {{ t('public.community.repliesTitle') }}
          </p>

          <div class="flex flex-col gap-4">
            <div v-if="replies.length === 0" class="text-center py-8 text-[#40617f] text-sm">
              {{ t('public.community.noReplies') }}
            </div>

            <div v-for="reply in replies" :key="reply.id"
              class="bg-white rounded-xl shadow-sm overflow-hidden">
              <div class="p-5 flex gap-4 items-start">

                <div class="w-9 h-9 rounded-full flex items-center justify-center shrink-0 text-sm font-bold text-white"
                  :style="`background: ${avatarColor(reply.author)}`">
                  {{ authorInitial(reply.author) }}
                </div>

                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-2 flex-wrap">
                    <span class="font-semibold text-[#001d32] text-sm">{{ authorName(reply.author) }}</span>
                    <span class="text-[#94a3b8] text-xs">· {{ formatDate(reply.created_at) }}</span>
                  </div>

                  <div v-if="editingReply !== reply.id" class="forum-content" v-html="reply.content" />

                  <div v-else>
                    <TiptapEditor v-model="editReplyContent" :placeholder="t('public.community.replyPlaceholder')" minHeight="100px" />
                    <div class="flex gap-2 mt-2 justify-end">
                      <button @click="editingReply = null" class="px-3 py-1.5 text-xs text-[#40617f] font-medium hover:text-[#001d32]">
                        {{ t('common.cancel') }}
                      </button>
                      <button @click="saveReplyEdit(reply)" :disabled="saving"
                        class="px-4 py-1.5 bg-[#006d35] text-white text-xs font-semibold rounded-lg hover:bg-[#1b8848] transition disabled:opacity-50">
                        {{ saving ? t('common.loading') : t('common.save') }}
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <div class="flex items-center gap-3 px-5 py-2 border-t border-[#f8fafc] bg-[#fafafa] flex-wrap">
                <template v-if="canEditReply(reply) && editingReply !== reply.id">
                  <button @click="startEditReply(reply)"
                    class="flex items-center gap-1 text-xs text-[#40617f] hover:text-[#006d35] transition font-medium">
                    <PencilSquareIcon class="w-3.5 h-3.5" />
                    {{ t('public.community.edit') }}
                  </button>
                </template>
                <template v-if="canDeleteReply(reply) && editingReply !== reply.id">
                  <button @click="deleteReply(reply.id)"
                    class="flex items-center gap-1 text-xs text-[#40617f] hover:text-[#ef4444] transition font-medium">
                    <TrashIcon class="w-3.5 h-3.5" />
                    {{ t('public.community.delete') }}
                  </button>
                </template>

                <button v-if="isOwner && reply.user_id !== userAuth.user?.id && !isBanned(reply.user_id)"
                  @click="banUser(reply.user_id)"
                  class="flex items-center gap-1 text-xs text-[#40617f] hover:text-[#ef4444] transition font-medium">
                  <ShieldExclamationIcon class="w-3.5 h-3.5" />
                  {{ t('public.community.ban') }}
                </button>
                <span v-if="isOwner && isBanned(reply.user_id)" class="text-xs text-[#ef4444] font-medium flex items-center gap-1">
                  <ShieldExclamationIcon class="w-3.5 h-3.5" />
                  {{ t('public.community.bannedLabel') }}
                </span>

                <button v-if="userAuth.isLoggedIn && !canEditReply(reply)"
                  @click="report('reply', reply.id)"
                  class="flex items-center gap-1 text-xs text-[#94a3b8] hover:text-[#ef4444] transition font-medium ml-auto">
                  <FlagIcon class="w-3.5 h-3.5" />
                  {{ t('public.community.report') }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm p-6">
          <h3 class="font-jakarta font-bold text-[#001d32] mb-4">{{ t('public.community.replyTitle') }}</h3>

          <template v-if="!userAuth.isLoggedIn">
            <p class="text-[#40617f] text-sm">
              {{ t('public.community.loginToReply') }}
              <RouterLink to="/connexion" class="text-[#006d35] font-semibold hover:underline">{{ t('public.community.loginLink') }}</RouterLink>
            </p>
          </template>

          <template v-else-if="isBanned(userAuth.user?.id)">
            <p class="text-[#ef4444] text-sm flex items-center gap-2">
              <ShieldExclamationIcon class="w-4 h-4" />
              {{ t('public.community.youAreBanned') }}
            </p>
          </template>

          <template v-else-if="thread.locked">
            <p class="text-[#94a3b8] text-sm flex items-center gap-2">
              <LockClosedIcon class="w-4 h-4" />
              {{ t('public.community.threadLocked') }}
            </p>
          </template>

          <template v-else>
            <TiptapEditor
              v-model="replyContent"
              :placeholder="t('public.community.replyPlaceholder')"
              minHeight="140px"
              class="mb-3"
            />
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

      <div v-else class="text-center py-24 text-[#40617f]">
        <p>{{ t('public.community.threadNotFound') }}</p>
        <RouterLink to="/communaute" class="text-[#006d35] font-semibold hover:underline text-sm mt-2 block">
          {{ t('public.community.backToCommunity') }}
        </RouterLink>
      </div>
    </div>

    <div v-if="reportModal.show" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="reportModal.show = false">
      <div class="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl">
        <h3 class="font-bold text-[#001d32] text-lg mb-4">{{ t('public.community.reportTitle') }}</h3>
        <textarea v-model="reportModal.reason" rows="4"
          :placeholder="t('public.community.reportPlaceholder')"
          class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none mb-4" />
        <p v-if="reportModal.error" class="text-red-500 text-xs mb-3">{{ reportModal.error }}</p>
        <div class="flex gap-3 justify-end">
          <button @click="reportModal.show = false" class="px-4 py-2 text-sm text-[#40617f] font-medium">{{ t('common.cancel') }}</button>
          <button @click="submitReport" :disabled="reportModal.sending"
            class="px-5 py-2 bg-[#ef4444] text-white text-sm font-semibold rounded-xl hover:bg-[#dc2626] transition disabled:opacity-50">
            {{ reportModal.sending ? t('common.loading') : t('public.community.reportSend') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import TiptapEditor from '@/components/TiptapEditor.vue'
import {
  ChevronRightIcon,
  LockClosedIcon,
  TrashIcon,
  PencilSquareIcon,
  FlagIcon,
  ShieldExclamationIcon,
  XMarkIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const loading = ref(true)
const thread = ref(null)
const replies = ref([])
const bans = ref([])
const replyContent = ref('')
const replyError = ref('')
const replySubmitting = ref(false)
const saving = ref(false)

const editingThread = ref(false)
const editThreadContent = ref('')
const editingReply = ref(null)
const editReplyContent = ref('')

const threadId = computed(() => route.params.id)
const isOwner = computed(() => userAuth.isLoggedIn && thread.value?.author?.id === userAuth.user?.id)

const reportModal = ref({ show: false, type: '', targetId: null, reason: '', sending: false, error: '' })

const COLORS = ['#006d35','#1a73e8','#7c3aed','#b45309','#ef4444','#0891b2','#059669']
function avatarColor(author) {
  if (!author) return '#94a3b8'
  return COLORS[(author.id || 0) % COLORS.length]
}
function authorInitial(author) {
  if (!author) return '?'
  return (author.first_name || author.email || '?')[0].toUpperCase()
}
function authorName(author) {
  if (!author) return 'Anonyme'
  if (author.first_name && author.last_name) return `${author.first_name} ${author.last_name}`
  if (author.first_name) return author.first_name
  return author.email?.split('@')[0] || 'Anonyme'
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

function isBanned(userId) {
  return bans.value.some(b => b.user_id === userId)
}

function canEditReply(reply) {
  return userAuth.isLoggedIn && reply.author?.id === userAuth.user?.id
}
function canDeleteReply(reply) {
  return userAuth.isLoggedIn && (reply.author?.id === userAuth.user?.id || isOwner.value)
}

async function loadThread() {
  const res = await fetch(`/api/public/v1/forum/threads/${threadId.value}`).catch(() => null)
  if (res?.ok) {
    const data = await res.json()
    thread.value = data.thread
    replies.value = data.replies || []
    bans.value = data.bans || []
  }
}

function startEditThread() {
  editThreadContent.value = thread.value.content
  editingThread.value = true
}
async function saveThreadEdit() {
  saving.value = true
  try {
    const res = await fetch(`/api/v1/forum/threads/${threadId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify({ content: editThreadContent.value }),
    })
    if (res.ok) {
      thread.value.content = editThreadContent.value
      editingThread.value = false
    }
  } finally { saving.value = false }
}

function startEditReply(reply) {
  editReplyContent.value = reply.content
  editingReply.value = reply.id
}
async function saveReplyEdit(reply) {
  saving.value = true
  try {
    const res = await fetch(`/api/v1/forum/replies/${reply.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify({ content: editReplyContent.value }),
    })
    if (res.ok) {
      reply.content = editReplyContent.value
      editingReply.value = null
    }
  } finally { saving.value = false }
}

async function toggleClose() {
  const res = await fetch(`/api/v1/forum/threads/${threadId.value}/close`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  if (res.ok) {
    const data = await res.json()
    thread.value.locked = data.locked
  }
}

async function deleteThread() {
  if (!confirm(t('public.community.confirmDelete'))) return
  const res = await fetch(`/api/v1/forum/threads/${threadId.value}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  if (res.ok) router.push(`/communaute/forum/${thread.value.category?.slug}`)
}

async function banUser(userId) {
  if (!confirm(t('public.community.confirmBan'))) return
  const res = await fetch(`/api/v1/forum/threads/${threadId.value}/ban/${userId}`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  if (res.ok) {
    const data = await res.json()
    bans.value.push(data.ban)
  }
}
async function unbanUser(userId) {
  await fetch(`/api/v1/forum/threads/${threadId.value}/ban/${userId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  bans.value = bans.value.filter(b => b.user_id !== userId)
}

async function submitReply() {
  replyError.value = ''
  const stripped = replyContent.value.replace(/<[^>]*>/g, '').trim()
  if (stripped.length < 2) {
    replyError.value = t('public.community.errorReplyContent')
    return
  }
  replySubmitting.value = true
  try {
    const res = await fetch(`/api/v1/forum/threads/${threadId.value}/replies`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify({ content: replyContent.value }),
    })
    if (res.ok) {
      const data = await res.json()
      replies.value.push(data.reply)
      replyContent.value = ''
    } else {
      const err = await res.json()
      replyError.value = err.error === 'banned' ? t('public.community.youAreBanned') : t('public.community.errorSubmit')
    }
  } finally { replySubmitting.value = false }
}

async function deleteReply(replyId) {
  if (!confirm(t('public.community.confirmDelete'))) return
  const res = await fetch(`/api/v1/forum/replies/${replyId}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${userAuth.token}` },
  })
  if (res.ok) replies.value = replies.value.filter(r => r.id !== replyId)
}

function report(type, id) {
  reportModal.value = { show: true, type, targetId: id, reason: '', sending: false, error: '' }
}
async function submitReport() {
  reportModal.value.error = ''
  reportModal.value.sending = true
  try {
    const res = await fetch('/api/v1/forum/reports', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${userAuth.token}` },
      body: JSON.stringify({ type: reportModal.value.type, target_id: reportModal.value.targetId, reason: reportModal.value.reason }),
    })
    if (res.ok) {
      reportModal.value.show = false
      alert(t('public.community.reportSent'))
    } else {
      reportModal.value.error = t('public.community.errorSubmit')
    }
  } finally { reportModal.value.sending = false }
}

onMounted(async () => {
  await loadThread()
  loading.value = false
})
</script>

<style>

.forum-content { line-height: 1.7; color: #1e293b; }
.forum-content p { margin-bottom: 0.75em; }
.forum-content p:last-child { margin-bottom: 0; }
.forum-content h1 { font-size: 1.5rem; font-weight: 800; margin-bottom: 0.5em; color: #001d32; }
.forum-content h2 { font-size: 1.25rem; font-weight: 700; margin-bottom: 0.5em; color: #001d32; }
.forum-content h3 { font-size: 1.1rem; font-weight: 700; margin-bottom: 0.5em; color: #001d32; }
.forum-content ul { list-style: disc; padding-left: 1.5em; margin-bottom: 0.75em; }
.forum-content ol { list-style: decimal; padding-left: 1.5em; margin-bottom: 0.75em; }
.forum-content li { margin-bottom: 0.25em; }
.forum-content blockquote { border-left: 4px solid #006d35; padding-left: 1em; color: #64748b; font-style: italic; margin: 0.75em 0; }
.forum-content code { background: #f1f5f9; border-radius: 4px; padding: 2px 6px; font-size: 0.875em; font-family: monospace; }
.forum-content pre { background: #1e293b; color: #e2e8f0; border-radius: 8px; padding: 12px 16px; overflow-x: auto; margin-bottom: 0.75em; }
.forum-content pre code { background: none; padding: 0; }
.forum-content img { max-width: 100%; border-radius: 8px; margin: 0.5em 0; }
.forum-content a { color: #006d35; text-decoration: underline; }
.forum-content hr { border: none; border-top: 2px solid #e5e7eb; margin: 1em 0; }
.forum-content mark { border-radius: 3px; padding: 0 2px; }
</style>
