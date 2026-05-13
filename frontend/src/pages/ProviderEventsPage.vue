<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[960px] mx-auto px-4 sm:px-6 py-10">

      <!-- Header -->
      <div class="mb-8">
        <button @click="router.push('/profil/pro')"
          class="text-sm text-[#40617f] hover:text-[#006d35] flex items-center gap-1.5 mb-4 transition">
          <ArrowLeftIcon class="w-4 h-4" />
          {{ t('public.providerEvents.back') }}
        </button>
        <div class="flex items-center justify-between">
          <div>
            <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">
              {{ t('public.providerEvents.title') }}
            </h1>
            <p class="text-[#40617f] text-sm mt-1">{{ t('public.providerEvents.subtitle') }}</p>
          </div>
          <button @click="openCreateForm"
            class="flex items-center gap-2 bg-[#006d35] text-white font-bold text-sm px-5 py-3 rounded-xl hover:bg-[#1b8848] transition">
            <PlusIcon class="w-4 h-4" />
            {{ t('public.providerEvents.btnNew') }}
          </button>
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="flex items-center justify-center py-24">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <template v-else>

        <!-- Create / Edit form -->
        <div v-if="showForm" class="bg-white rounded-[32px] p-8 mb-8 shadow-sm">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl mb-6">
            {{ editingId ? t('public.providerEvents.formTitleEdit') : t('public.providerEvents.formTitleCreate') }}
          </h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-5">
            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldTitle') }}
              </label>
              <input v-model="form.title" type="text"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldStartDate') }}
              </label>
              <input v-model="form.start_date" type="datetime-local"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldEndDate') }}
              </label>
              <input v-model="form.end_date" type="datetime-local"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldLocation') }}
              </label>
              <input v-model="form.location" type="text"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldMaxPart') }}
              </label>
              <input v-model.number="form.max_participants" type="number" min="1"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20" />
            </div>
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldCategory') }}
              </label>
              <select v-model="form.category_id"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 bg-white">
                <option :value="null">{{ t('public.providerEvents.fieldCategoryNone') }}</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.name }}</option>
              </select>
            </div>
            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1">
                {{ t('public.providerEvents.fieldDesc') }}
              </label>
              <textarea v-model="form.description" rows="4"
                class="w-full border border-gray-200 rounded-xl px-4 py-3 text-sm text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 resize-none" />
            </div>
          </div>

          <p v-if="formError" class="mt-4 text-sm text-red-600 font-medium">{{ formError }}</p>

          <div class="flex gap-3 mt-6">
            <button @click="saveForm" :disabled="formLoading"
              class="bg-[#006d35] text-white font-bold text-sm px-6 py-3 rounded-xl hover:bg-[#1b8848] transition disabled:opacity-60">
              {{ formLoading ? t('public.providerEvents.btnSaving') : t('public.providerEvents.btnSave') }}
            </button>
            <button @click="closeForm"
              class="bg-gray-100 text-[#40617f] font-bold text-sm px-6 py-3 rounded-xl hover:bg-gray-200 transition">
              {{ t('public.providerEvents.btnCancel') }}
            </button>
          </div>
        </div>

        <!-- Empty state -->
        <div v-if="!events.length && !showForm" class="flex flex-col items-center justify-center py-20 text-center">
          <CalendarDaysIcon class="w-14 h-14 text-[#40617f]/30 mb-4" />
          <p class="font-jakarta font-bold text-[#001d32] text-lg mb-1">{{ t('public.providerEvents.emptyTitle') }}</p>
          <p class="text-[#40617f] text-sm">{{ t('public.providerEvents.emptyHint') }}</p>
        </div>

        <!-- Events list -->
        <div v-if="events.length" class="flex flex-col gap-4">
          <div v-for="ev in events" :key="ev.id"
            class="bg-white rounded-2xl p-6 shadow-sm flex flex-col sm:flex-row sm:items-center gap-4">

            <!-- Date badge -->
            <div class="shrink-0 w-14 h-14 rounded-xl flex flex-col items-center justify-center text-white font-jakarta font-bold"
              :class="statusBgClass(ev.status)">
              <span class="text-xs uppercase leading-none">{{ monthShort(ev.start_date) }}</span>
              <span class="text-xl leading-tight">{{ dayNum(ev.start_date) }}</span>
            </div>

            <!-- Info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap mb-1">
                <span :class="['text-xs font-bold px-2.5 py-1 rounded-full', statusClass(ev.status)]">
                  {{ statusLabel(ev.status) }}
                </span>
                <span v-if="ev.registrations_count > 0" class="text-xs text-[#40617f]">
                  {{ ev.registrations_count }} inscrit{{ ev.registrations_count > 1 ? 's' : '' }}
                </span>
              </div>
              <p class="font-jakarta font-bold text-[#001d32] text-base truncate">{{ ev.title }}</p>
              <p class="text-[#40617f] text-xs mt-0.5">{{ formatDate(ev.start_date) }}</p>
            </div>

            <!-- Actions -->
            <div class="flex items-center gap-2 shrink-0 flex-wrap">
              <button v-if="ev.status === 'draft' || ev.status === 'pending'"
                @click="editEvent(ev)"
                class="text-xs bg-[#edf4ff] text-[#1a73e8] font-semibold px-3 py-2 rounded-xl hover:bg-[#d0e8ff] transition">
                {{ t('public.providerEvents.btnEdit') }}
              </button>
              <button v-if="ev.status === 'draft'"
                @click="submitEvent(ev)"
                class="text-xs bg-[#d1fae5] text-[#065f46] font-semibold px-3 py-2 rounded-xl hover:bg-[#a7f3d0] transition">
                {{ t('public.providerEvents.btnSubmit') }}
              </button>
              <button v-if="ev.status !== 'published'"
                @click="deleteEvent(ev)"
                class="text-xs bg-red-50 text-red-600 font-semibold px-3 py-2 rounded-xl hover:bg-red-100 transition">
                {{ t('public.providerEvents.btnDelete') }}
              </button>
            </div>
          </div>
        </div>

      </template>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import {
  ArrowLeftIcon,
  PlusIcon,
  CalendarDaysIcon,
} from '@heroicons/vue/24/outline'

const router = useRouter()
const { t, locale } = useI18n()
const userAuth = useUserAuthStore()

const loading = ref(true)
const events = ref([])
const categories = ref([])
const showForm = ref(false)
const editingId = ref(null)
const formLoading = ref(false)
const formError = ref('')

const form = ref({
  title: '',
  description: '',
  location: '',
  start_date: '',
  end_date: '',
  max_participants: null,
  category_id: null,
})

function toDatetimeLocal(iso) {
  if (!iso) return ''
  const d = new Date(iso)
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function openCreateForm() {
  editingId.value = null
  form.value = { title: '', description: '', location: '', start_date: '', end_date: '', max_participants: null, category_id: null }
  formError.value = ''
  showForm.value = true
}

function editEvent(ev) {
  editingId.value = ev.id
  form.value = {
    title: ev.title || '',
    description: ev.description || '',
    location: ev.location || '',
    start_date: toDatetimeLocal(ev.start_date),
    end_date: toDatetimeLocal(ev.end_date),
    max_participants: ev.max_participants || null,
    category_id: ev.category?.id || null,
  }
  formError.value = ''
  showForm.value = true
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function closeForm() {
  showForm.value = false
  editingId.value = null
  formError.value = ''
}

function toISO(datetimeLocal) {
  if (!datetimeLocal) return ''
  return new Date(datetimeLocal).toISOString()
}

async function saveForm() {
  formError.value = ''
  if (!form.value.title.trim()) {
    formError.value = 'Le titre est requis.'
    return
  }
  if (!form.value.start_date || !form.value.end_date) {
    formError.value = 'Les dates de début et de fin sont requises.'
    return
  }

  const payload = {
    title: form.value.title.trim(),
    description: form.value.description || null,
    location: form.value.location || null,
    start_date: toISO(form.value.start_date),
    end_date: toISO(form.value.end_date),
    max_participants: form.value.max_participants || null,
    category_id: form.value.category_id || null,
  }

  formLoading.value = true
  try {
    const url = editingId.value
      ? `/api/v1/provider/events/${editingId.value}`
      : '/api/v1/provider/events'
    const method = editingId.value ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${userAuth.token}`,
      },
      body: JSON.stringify(payload),
    })
    if (res.ok) {
      const saved = await res.json()
      if (editingId.value) {
        const idx = events.value.findIndex(e => e.id === editingId.value)
        if (idx !== -1) events.value[idx] = saved
      } else {
        events.value.unshift(saved)
      }
      closeForm()
    } else {
      const data = await res.json().catch(() => ({}))
      formError.value = data.message || 'Erreur lors de la sauvegarde.'
    }
  } catch {
    formError.value = 'Erreur réseau.'
  }
  formLoading.value = false
}

async function submitEvent(ev) {
  if (!confirm(t('public.providerEvents.confirmSubmit'))) return
  try {
    const res = await fetch(`/api/v1/provider/events/${ev.id}/submit`, {
      method: 'PUT',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok) {
      const idx = events.value.findIndex(e => e.id === ev.id)
      if (idx !== -1) events.value[idx] = { ...events.value[idx], status: 'pending' }
    }
  } catch {}
}

async function deleteEvent(ev) {
  if (!confirm(t('public.providerEvents.confirmDelete'))) return
  try {
    const res = await fetch(`/api/v1/provider/events/${ev.id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (res.ok || res.status === 204) {
      events.value = events.value.filter(e => e.id !== ev.id)
    }
  } catch {}
}

function statusLabel(s) {
  const map = {
    draft: t('public.providerEvents.statusDraft'),
    pending: t('public.providerEvents.statusPending'),
    published: t('public.providerEvents.statusPublished'),
    rejected: t('public.providerEvents.statusRejected'),
  }
  return map[s] || s
}

function statusClass(s) {
  const map = {
    draft: 'bg-gray-100 text-gray-600',
    pending: 'bg-yellow-50 text-yellow-700',
    published: 'bg-green-50 text-green-700',
    rejected: 'bg-red-50 text-red-600',
  }
  return map[s] || 'bg-gray-100 text-gray-600'
}

function statusBgClass(s) {
  const map = {
    draft: 'bg-gray-400',
    pending: 'bg-yellow-500',
    published: 'bg-gradient-to-br from-[#006d35] to-[#1b8848]',
    rejected: 'bg-red-400',
  }
  return map[s] || 'bg-gray-400'
}

const localeCode = () => locale.value === 'en' ? 'en-US' : 'fr-FR'

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode(), { day: 'numeric', month: 'long', year: 'numeric' })
}

function monthShort(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(localeCode(), { month: 'short' }).replace('.', '')
}

function dayNum(iso) {
  if (!iso) return ''
  return new Date(iso).getDate()
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/pro/evenements')
    return
  }

  const headers = { Authorization: `Bearer ${userAuth.token}`, Accept: 'application/json' }

  const [evRes, catRes] = await Promise.all([
    fetch('/api/v1/provider/events?per_page=50', { headers }).catch(() => null),
    fetch('/api/public/v1/categories', { headers }).catch(() => null),
  ])

  if (evRes?.ok) {
    const data = await evRes.json()
    events.value = data.data || []
  }

  if (catRes?.ok) {
    const data = await catRes.json()
    categories.value = data.data || data || []
  }

  loading.value = false
})
</script>
