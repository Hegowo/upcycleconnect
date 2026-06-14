<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[1100px] mx-auto px-4 sm:px-6 py-10">

      <div class="mb-8">
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">
          {{ t('public.community.title') }}
        </h1>
        <p class="text-[#40617f] text-sm mt-1">{{ t('public.community.subtitle') }}</p>
      </div>

      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else>

        <RouterLink
          to="/communaute/conseils"
          class="group flex items-center justify-between gap-4 bg-gradient-to-br from-[#006d35] to-[#1b8848] rounded-2xl p-5 mb-6 hover:shadow-lg transition"
        >
          <div class="flex items-center gap-4 min-w-0">
            <div class="w-12 h-12 rounded-xl bg-white/15 flex items-center justify-center shrink-0">
              <LightBulbIcon class="w-6 h-6 text-white" />
            </div>
            <div class="min-w-0">
              <p class="font-jakarta font-extrabold text-white text-lg leading-tight">Conseils & astuces</p>
              <p class="text-white/80 text-xs mt-0.5">Articles publiés par l'équipe UpcycleConnect pour vous accompagner.</p>
            </div>
          </div>
          <ArrowRightIcon class="w-5 h-5 text-white opacity-70 group-hover:translate-x-1 transition" />
        </RouterLink>

        <section class="relative rounded-2xl overflow-hidden mb-10 bg-gradient-to-br from-[#001d32] to-[#003060] p-8 sm:p-12">
          <div class="relative z-10">
            <p class="text-[#7dd3b0] text-xs font-bold uppercase tracking-widest mb-2">{{ t('public.community.heroTag') }}</p>
            <h2 class="font-jakarta font-extrabold text-white text-3xl sm:text-4xl leading-tight max-w-lg">
              {{ t('public.community.heroTitle') }}
            </h2>
            <p class="text-[#94a3b8] text-base mt-3 max-w-md">{{ t('public.community.heroSubtitle') }}</p>
            <button
              v-if="userAuth.isLoggedIn"
              @click="openNewThread(null)"
              class="mt-6 inline-flex items-center gap-2 px-5 py-2.5 bg-[#006d35] hover:bg-[#1b8848] text-white text-sm font-semibold rounded-xl transition">
              <PencilSquareIcon class="w-4 h-4" />
              {{ t('public.community.newThread') }}
            </button>
            <RouterLink v-else to="/connexion" class="mt-6 inline-flex items-center gap-2 px-5 py-2.5 bg-[#006d35] hover:bg-[#1b8848] text-white text-sm font-semibold rounded-xl transition">
              {{ t('public.community.loginToPost') }}
            </RouterLink>
          </div>
          <div class="absolute right-8 top-8 w-40 h-40 rounded-full bg-[#006d35]/20 blur-2xl pointer-events-none" />
        </section>

        <section class="mb-10">
          <div class="flex items-center justify-between mb-5">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.community.categoriesTitle') }}</h2>
          </div>

          <div v-if="categories.length === 0" class="text-center py-10 text-[#40617f] text-sm">
            {{ t('public.community.noCategories') }}
          </div>

          <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            <RouterLink
              v-for="cat in categories"
              :key="cat.id"
              :to="`/communaute/forum/${cat.slug}`"
              class="bg-white rounded-2xl p-5 shadow-sm hover:shadow-md transition flex flex-col gap-3 group">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
                  :style="`background: ${cat.color || '#edf4ff'}`">
                  <ChatBubbleLeftEllipsisIcon class="w-5 h-5 text-[#006d35]" />
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-jakarta font-bold text-[#001d32] text-sm truncate group-hover:text-[#006d35] transition">{{ cat.name }}</p>
                </div>
              </div>
              <p v-if="cat.description" class="text-[#64748b] text-xs leading-relaxed line-clamp-2">{{ cat.description }}</p>
              <div class="flex items-center gap-4 text-xs text-[#94a3b8] mt-auto pt-2 border-t border-[#f1f5f9]">
                <span class="flex items-center gap-1">
                  <DocumentTextIcon class="w-3.5 h-3.5" />
                  {{ cat.thread_count }} {{ t('public.community.threads') }}
                </span>
                <span class="flex items-center gap-1">
                  <ChatBubbleLeftIcon class="w-3.5 h-3.5" />
                  {{ cat.reply_count }} {{ t('public.community.replies') }}
                </span>
              </div>
            </RouterLink>
          </div>
        </section>

        <section>
          <div class="flex items-center justify-between mb-5">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.community.recentTitle') }}</h2>
          </div>

          <div v-if="recentThreads.length === 0" class="text-center py-10 text-[#40617f] text-sm">
            {{ t('public.community.noThreads') }}
          </div>

          <div v-else class="flex flex-col gap-3">
            <RouterLink
              v-for="thread in recentThreads"
              :key="thread.id"
              :to="`/communaute/forum/${thread.category?.slug}/${thread.id}`"
              class="bg-white rounded-xl p-4 shadow-sm hover:shadow-md transition flex items-center gap-4 group">

              <div class="w-9 h-9 rounded-full bg-[#edf4ff] flex items-center justify-center shrink-0 text-sm font-bold text-[#40617f]">
                {{ authorInitial(thread.author) }}
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 flex-wrap mb-1">
                  <span v-if="thread.pinned" class="text-xs font-bold text-[#006d35] flex items-center gap-0.5">
                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20"><path d="M9.75 3A.75.75 0 0 1 10 3c.414 0 .75.336.75.75v.693l3.18 1.272A.75.75 0 0 1 14.5 6.5v.25a.75.75 0 0 1-.75.75H6.25A.75.75 0 0 1 5.5 6.75V6.5a.75.75 0 0 1 .57-.727l3.18-1.272V3.75A.75.75 0 0 1 9.75 3zM6 9.5h8l-1 7H7L6 9.5z"/></svg>
                  </span>
                  <span v-if="thread.locked" class="text-xs font-bold text-[#ef4444] flex items-center gap-0.5">
                    <LockClosedIcon class="w-3 h-3" />
                  </span>
                  <span class="text-xs px-2 py-0.5 rounded-full bg-[#edf4ff] text-[#40617f] font-medium">
                    {{ thread.category?.name }}
                  </span>
                </div>
                <p class="font-semibold text-[#001d32] text-sm truncate group-hover:text-[#006d35] transition">{{ thread.title }}</p>
                <p class="text-[#94a3b8] text-xs mt-0.5">
                  {{ thread.author?.first_name || thread.author?.email?.split('@')[0] }}
                  · {{ formatDate(thread.created_at) }}
                </p>
              </div>

              <div class="shrink-0 text-center">
                <p class="font-bold text-[#001d32] text-sm">{{ thread.reply_count }}</p>
                <p class="text-[#94a3b8] text-[10px] uppercase tracking-wide">{{ t('public.community.repliesLabel') }}</p>
              </div>

              <ChevronRightIcon class="w-4 h-4 text-[#94a3b8] shrink-0" />
            </RouterLink>
          </div>
        </section>

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
            <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">{{ t('public.community.modalCategory') }}</label>
            <select v-model="form.categoryId" class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 bg-white">
              <option value="">{{ t('public.community.modalCategoryPlaceholder') }}</option>
              <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
            </select>
          </div>

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
            <button @click="showModal = false"
              class="px-4 py-2 text-sm text-[#40617f] hover:text-[#001d32] font-medium">
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
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import TiptapEditor from '@/components/TiptapEditor.vue'
import {
  ChatBubbleLeftEllipsisIcon,
  ChatBubbleLeftIcon,
  DocumentTextIcon,
  ChevronRightIcon,
  PencilSquareIcon,
  XMarkIcon,
  LockClosedIcon,
  LightBulbIcon,
  ArrowRightIcon,
} from '@heroicons/vue/24/outline'

const { t, locale } = useI18n()
const router = useRouter()
const userAuth = useUserAuthStore()

const loading = ref(true)
const categories = ref([])
const recentThreads = ref([])

const showModal = ref(false)
const submitting = ref(false)
const formError = ref('')
const form = ref({ categoryId: '', title: '', content: '' })

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

function openNewThread(categoryId) {
  form.value = { categoryId: categoryId || '', title: '', content: '' }
  formError.value = ''
  showModal.value = true
}

async function submitThread() {
  formError.value = ''
  if (!form.value.categoryId) { formError.value = t('public.community.errorCategory'); return }
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
        category_id: Number(form.value.categoryId),
        title: form.value.title,
        content: form.value.content,
      }),
    })
    if (res.ok) {
      const data = await res.json()
      const cat = categories.value.find(c => c.id === data.thread.category_id)
      showModal.value = false
      if (cat) {
        router.push(`/communaute/forum/${cat.slug}/${data.thread.id}`)
      } else {
        await loadData()
      }
    } else {
      formError.value = t('public.community.errorSubmit')
    }
  } finally {
    submitting.value = false
  }
}

async function loadData() {
  const [catRes, threadRes] = await Promise.all([
    fetch('/api/public/v1/forum/categories').catch(() => null),
    fetch('/api/public/v1/forum/categories').catch(() => null),
  ])

  if (catRes?.ok) {
    const data = await catRes.json()
    categories.value = data.data || []
  }

  const threads = []
  for (const cat of categories.value) {
    const r = await fetch(`/api/public/v1/forum/categories/${cat.slug}?page=1`).catch(() => null)
    if (r?.ok) {
      const d = await r.json()
      threads.push(...(d.threads || []).slice(0, 3).map(t => ({ ...t, category: cat })))
    }
  }
  threads.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
  recentThreads.value = threads.slice(0, 10)
}

onMounted(async () => {
  await loadData()
  loading.value = false
})
</script>
