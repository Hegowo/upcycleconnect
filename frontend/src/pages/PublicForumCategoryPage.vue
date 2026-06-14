<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[900px] mx-auto px-4 sm:px-6 py-10">

      <div class="mb-8">
        <button @click="router.push('/communaute')"
          class="text-sm text-[#40617f] hover:text-[#006d35] flex items-center gap-1.5 mb-4 transition">
          <ArrowLeftIcon class="w-4 h-4" />
          {{ t('public.community.backToCommunity') }}
        </button>

        <div v-if="loading" class="flex items-center justify-center py-24">
          <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
        </div>

        <template v-else-if="category">
          <div class="flex items-center gap-4 mb-1">
            <div class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
              :style="`background: ${category.color || '#edf4ff'}`">
              <ChatBubbleLeftEllipsisIcon class="w-5 h-5 text-[#006d35]" />
            </div>
            <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl tracking-tight">{{ category.name }}</h1>
          </div>
          <p v-if="category.description" class="text-[#40617f] text-sm mt-1 ml-14">{{ category.description }}</p>
        </template>
      </div>

      <template v-if="!loading">

        <div class="flex items-center justify-between mb-5">
          <p class="text-[#64748b] text-sm">{{ total }} {{ t('public.community.threads') }}</p>
          <button v-if="userAuth.isLoggedIn" @click="openNewThread"
            class="inline-flex items-center gap-2 px-4 py-2 bg-[#006d35] text-white text-sm font-semibold rounded-xl hover:bg-[#1b8848] transition">
            <PencilSquareIcon class="w-4 h-4" />
            {{ t('public.community.newThread') }}
          </button>
          <RouterLink v-else to="/connexion"
            class="inline-flex items-center gap-2 px-4 py-2 bg-[#006d35] text-white text-sm font-semibold rounded-xl hover:bg-[#1b8848] transition">
            {{ t('public.community.loginToPost') }}
          </RouterLink>
        </div>

        <div v-if="threads.length === 0" class="text-center py-16 text-[#40617f]">
          <ChatBubbleLeftEllipsisIcon class="w-12 h-12 mx-auto mb-3 text-[#d1d5db]" />
          <p class="text-sm">{{ t('public.community.noThreadsInCategory') }}</p>
        </div>

        <div v-else class="flex flex-col gap-3 mb-8">
          <RouterLink
            v-for="thread in threads"
            :key="thread.id"
            :to="`/communaute/forum/${slug}/${thread.id}`"
            class="bg-white rounded-xl p-4 shadow-sm hover:shadow-md transition flex items-center gap-4 group">

            <div class="shrink-0 w-9 h-9 rounded-full bg-[#f1f5f9] flex items-center justify-center text-sm font-bold text-[#40617f]">
              {{ authorInitial(thread.author) }}
            </div>

            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span v-if="thread.pinned" title="Épinglé" class="text-[#006d35]">
                  <svg class="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 20 20"><path d="M9.75 3A.75.75 0 0 1 10 3c.414 0 .75.336.75.75v.693l3.18 1.272A.75.75 0 0 1 14.5 6.5v.25a.75.75 0 0 1-.75.75H6.25A.75.75 0 0 1 5.5 6.75V6.5a.75.75 0 0 1 .57-.727l3.18-1.272V3.75A.75.75 0 0 1 9.75 3zM6 9.5h8l-1 7H7L6 9.5z"/></svg>
                </span>
                <LockClosedIcon v-if="thread.locked" class="w-3.5 h-3.5 text-[#ef4444]" />
              </div>
              <p class="font-semibold text-[#001d32] text-sm group-hover:text-[#006d35] transition leading-snug truncate">{{ thread.title }}</p>
              <p class="text-[#94a3b8] text-xs mt-0.5">
                {{ thread.author?.first_name || thread.author?.email?.split('@')[0] }}
                · {{ formatDate(thread.created_at) }}
              </p>
            </div>

            <div class="shrink-0 text-center min-w-[40px]">
              <p class="font-bold text-[#001d32] text-sm">{{ thread.reply_count }}</p>
              <p class="text-[#94a3b8] text-[10px] uppercase tracking-wide">{{ t('public.community.repliesLabel') }}</p>
            </div>

            <ChevronRightIcon class="w-4 h-4 text-[#94a3b8] shrink-0" />
          </RouterLink>
        </div>

        <div v-if="totalPages > 1" class="flex justify-center gap-2">
          <button v-for="p in totalPages" :key="p" @click="goPage(p)"
            :class="[
              'w-9 h-9 rounded-lg text-sm font-medium transition',
              p === page
                ? 'bg-[#006d35] text-white'
                : 'bg-white text-[#40617f] border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35]'
            ]">
            {{ p }}
          </button>
        </div>
      </template>
    </div>

    <div v-if="showModal" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="showModal = false">
      <div class="bg-white rounded-2xl p-6 w-full max-w-lg shadow-xl">
        <div class="flex items-center justify-between mb-5">
          <h3 class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.community.modalTitle') }}</h3>
          <button @click="showModal = false" class="text-[#94a3b8] hover:text-[#001d32]">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="flex flex-col gap-4">
          <div>
            <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">{{ t('public.community.modalThreadTitle') }}</label>
            <input v-model="form.title" type="text"
              :placeholder="t('public.community.modalTitlePlaceholder')"
              class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>

          <div>
            <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">{{ t('public.community.modalContent') }}</label>
            <TiptapEditor v-model="form.content" :placeholder="t('public.community.modalContentPlaceholder')" minHeight="160px" />
          </div>

          <p v-if="formError" class="text-red-500 text-xs">{{ formError }}</p>

          <div class="flex gap-3 justify-end">
            <button @click="showModal = false" class="px-4 py-2 text-sm text-[#40617f] hover:text-[#001d32] font-medium">
              {{ t('common.cancel') }}
            </button>
            <button @click="submitThread" :disabled="submitting"
              class="px-5 py-2 bg-[#006d35] text-white text-sm font-semibold rounded-xl hover:bg-[#1b8848] transition disabled:opacity-50">
              {{ submitting ? t('common.loading') : t('public.community.modalSubmit') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import TiptapEditor from '@/components/TiptapEditor.vue'
import {
  ArrowLeftIcon,
  ChatBubbleLeftEllipsisIcon,
  ChevronRightIcon,
  PencilSquareIcon,
  XMarkIcon,
  LockClosedIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const userAuth = useUserAuthStore()

const slug = computed(() => route.params.slug)
const loading = ref(true)
const category = ref(null)
const threads = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const showModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = ref({ title: '', content: '' })

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

function openNewThread() {
  form.value = { title: '', content: '' }
  formError.value = ''
  showModal.value = true
}

async function submitThread() {
  formError.value = ''
  if (form.value.title.length < 5) { formError.value = t('public.community.errorTitle'); return }
  const stripped = form.value.content.replace(/<[^>]*>/g, '').trim()
  if (stripped.length < 10) { formError.value = t('public.community.errorContent'); return }

  submitting.value = true
  try {
    const res = await fetch('/api/v1/forum/threads', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify({
        category_id: category.value.id,
        title: form.value.title,
        content: form.value.content,
      }),
    })
    if (res.ok) {
      const data = await res.json()
      showModal.value = false
      router.push(`/communaute/forum/${slug.value}/${data.thread.id}`)
    } else {
      formError.value = t('public.community.errorSubmit')
    }
  } finally {
    submitting.value = false
  }
}

async function goPage(p) {
  page.value = p
  await loadThreads()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function loadThreads() {
  const res = await fetch(`/api/public/v1/forum/categories/${slug.value}?page=${page.value}`).catch(() => null)
  if (res?.ok) {
    const data = await res.json()
    category.value = data.category
    threads.value = data.threads || []
    total.value = data.total || 0
  }
}

onMounted(async () => {
  await loadThreads()
  loading.value = false
})

watch(slug, async () => {
  loading.value = true
  page.value = 1
  await loadThreads()
  loading.value = false
})
</script>
