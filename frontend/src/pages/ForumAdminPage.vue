<template>
  <div class="space-y-6">

    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Forum Communauté</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Gestion des catégories et modération des discussions</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="flex gap-2 border-b border-[#e5e7eb]">
      <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key"
        :class="[
          'px-4 py-2 text-sm font-semibold transition border-b-2 -mb-px',
          activeTab === tab.key
            ? 'border-[#006d35] text-[#006d35]'
            : 'border-transparent text-[#64748b] hover:text-[#001d32]'
        ]">
        {{ tab.label }}
      </button>
    </div>

    <!-- Categories tab -->
    <div v-if="activeTab === 'categories'">
      <div class="flex items-center justify-between mb-4">
        <p class="text-sm text-[#64748b]">{{ categories.length }} catégorie(s)</p>
        <button @click="openCatModal(null)"
          class="inline-flex items-center gap-2 px-4 py-2 bg-[#006d35] text-white text-sm font-semibold rounded-xl hover:bg-[#1b8848] transition">
          <PlusIcon class="w-4 h-4" />
          Nouvelle catégorie
        </button>
      </div>

      <div v-if="catLoading" class="flex items-center justify-center py-12">
        <div class="w-7 h-7 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else-if="categories.length === 0" class="text-center py-12 text-[#64748b] text-sm">
        Aucune catégorie. Créez-en une pour lancer le forum !
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="cat in categories" :key="cat.id"
          class="bg-white rounded-2xl p-5 shadow-sm border border-[#f1f5f9] flex flex-col gap-3">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-xl flex items-center justify-center shrink-0"
                :style="`background: ${cat.color || '#edf4ff'}`">
                <ChatBubbleLeftEllipsisIcon class="w-5 h-5 text-[#006d35]" />
              </div>
              <p class="font-bold text-[#001d32] text-sm">{{ cat.name }}</p>
            </div>
            <div class="flex gap-2">
              <button @click="openCatModal(cat)" class="text-[#40617f] hover:text-[#006d35]">
                <PencilIcon class="w-4 h-4" />
              </button>
              <button @click="deleteCategory(cat.id)" class="text-[#40617f] hover:text-[#ef4444]">
                <TrashIcon class="w-4 h-4" />
              </button>
            </div>
          </div>
          <p v-if="cat.description" class="text-[#64748b] text-xs line-clamp-2">{{ cat.description }}</p>
          <div class="flex items-center gap-4 text-xs text-[#94a3b8]">
            <span>{{ cat.thread_count }} sujets</span>
            <span>{{ cat.reply_count }} réponses</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Threads tab (moderation) -->
    <div v-if="activeTab === 'threads'">
      <div class="flex items-center justify-between mb-4">
        <p class="text-sm text-[#64748b]">{{ totalThreads }} discussion(s)</p>
      </div>

      <div v-if="threadLoading" class="flex items-center justify-center py-12">
        <div class="w-7 h-7 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
      </div>

      <div v-else class="bg-white rounded-2xl shadow-sm border border-[#f1f5f9] overflow-hidden">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-[#f1f5f9] bg-[#f8fafc]">
              <th class="text-left px-4 py-3 text-xs font-semibold text-[#64748b] uppercase tracking-wide">Titre</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-[#64748b] uppercase tracking-wide hidden sm:table-cell">Catégorie</th>
              <th class="text-left px-4 py-3 text-xs font-semibold text-[#64748b] uppercase tracking-wide hidden md:table-cell">Auteur</th>
              <th class="text-center px-4 py-3 text-xs font-semibold text-[#64748b] uppercase tracking-wide">Rép.</th>
              <th class="text-right px-4 py-3 text-xs font-semibold text-[#64748b] uppercase tracking-wide">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="thread in threads" :key="thread.id"
              class="border-b border-[#f1f5f9] hover:bg-[#f8fafc] transition"
              :class="thread.deleted_at ? 'opacity-50' : ''">
              <td class="px-4 py-3">
                <div class="flex items-center gap-2">
                  <span v-if="thread.pinned" class="text-[#006d35]" title="Épinglé">📌</span>
                  <span v-if="thread.locked" class="text-[#ef4444]" title="Verrouillé">🔒</span>
                  <p class="font-medium text-[#001d32] truncate max-w-[200px]">{{ thread.title }}</p>
                </div>
              </td>
              <td class="px-4 py-3 text-[#64748b] hidden sm:table-cell">{{ thread.category?.name }}</td>
              <td class="px-4 py-3 text-[#64748b] hidden md:table-cell">
                {{ thread.author?.first_name || thread.author?.email?.split('@')[0] }}
              </td>
              <td class="px-4 py-3 text-center text-[#64748b]">{{ thread.reply_count }}</td>
              <td class="px-4 py-3 text-right">
                <div class="flex items-center justify-end gap-2">
                  <button @click="togglePin(thread)"
                    :class="['text-xs font-medium px-2 py-1 rounded-lg transition', thread.pinned ? 'bg-[#d1fae5] text-[#006d35]' : 'bg-[#f1f5f9] text-[#64748b] hover:bg-[#e2e8f0]']">
                    {{ thread.pinned ? 'Désépingler' : 'Épingler' }}
                  </button>
                  <button @click="toggleLock(thread)"
                    :class="['text-xs font-medium px-2 py-1 rounded-lg transition', thread.locked ? 'bg-[#fee2e2] text-[#ef4444]' : 'bg-[#f1f5f9] text-[#64748b] hover:bg-[#e2e8f0]']">
                    {{ thread.locked ? 'Déverrouiller' : 'Verrouiller' }}
                  </button>
                  <button @click="deleteThread(thread.id)" class="text-[#ef4444] hover:text-[#dc2626]">
                    <TrashIcon class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="threads.length === 0">
              <td colspan="5" class="px-4 py-8 text-center text-[#64748b] text-sm">Aucune discussion</td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div v-if="totalThreadPages > 1" class="flex justify-center gap-2 p-4 border-t border-[#f1f5f9]">
          <button v-for="p in totalThreadPages" :key="p" @click="goThreadPage(p)"
            :class="[
              'w-8 h-8 rounded-lg text-xs font-medium transition',
              p === threadPage
                ? 'bg-[#006d35] text-white'
                : 'bg-[#f1f5f9] text-[#64748b] hover:bg-[#e2e8f0]'
            ]">
            {{ p }}
          </button>
        </div>
      </div>
    </div>

    <!-- Category modal -->
    <div v-if="catModal.show" class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4" @click.self="catModal.show = false">
      <div class="bg-white rounded-2xl p-6 w-full max-w-md shadow-xl">
        <div class="flex items-center justify-between mb-5">
          <h3 class="font-bold text-[#001d32] text-lg">{{ catModal.id ? 'Modifier la catégorie' : 'Nouvelle catégorie' }}</h3>
          <button @click="catModal.show = false" class="text-[#94a3b8] hover:text-[#001d32]">
            <XMarkIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="flex flex-col gap-4">
          <div>
            <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">Nom *</label>
            <input v-model="catModal.name" type="text" placeholder="Ex: Techniques & Matériaux"
              class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
          </div>
          <div>
            <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">Description</label>
            <textarea v-model="catModal.description" rows="3" placeholder="Description de la catégorie..."
              class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" />
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">Couleur fond</label>
              <div class="flex items-center gap-2">
                <input v-model="catModal.color" type="color" class="w-10 h-10 rounded-lg cursor-pointer border border-[#e5e7eb]" />
                <input v-model="catModal.color" type="text" placeholder="#edf4ff"
                  class="flex-1 border border-[#e5e7eb] rounded-xl px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>
            </div>
            <div>
              <label class="text-xs font-semibold text-[#40617f] uppercase tracking-wide mb-1 block">Ordre</label>
              <input v-model.number="catModal.sort_order" type="number" min="0"
                class="w-full border border-[#e5e7eb] rounded-xl px-3 py-2.5 text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>
          </div>

          <p v-if="catModal.error" class="text-red-500 text-xs">{{ catModal.error }}</p>

          <div class="flex gap-3 justify-end">
            <button @click="catModal.show = false" class="px-4 py-2 text-sm text-[#40617f] hover:text-[#001d32] font-medium">
              Annuler
            </button>
            <button @click="saveCategory" :disabled="catModal.saving"
              class="px-5 py-2 bg-[#006d35] text-white text-sm font-semibold rounded-xl hover:bg-[#1b8848] transition disabled:opacity-50">
              {{ catModal.saving ? 'Enregistrement...' : 'Enregistrer' }}
            </button>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import {
  ChatBubbleLeftEllipsisIcon,
  PlusIcon,
  PencilIcon,
  TrashIcon,
  XMarkIcon,
} from '@heroicons/vue/24/outline'

const auth = useAuthStore()

const activeTab = ref('categories')
const tabs = [
  { key: 'categories', label: 'Catégories' },
  { key: 'threads',    label: 'Discussions (modération)' },
]

// --- Categories ---
const catLoading = ref(true)
const categories = ref([])

const catModal = ref({
  show: false, id: null, name: '', description: '', color: '#edf4ff', sort_order: 0, saving: false, error: '',
})

function openCatModal(cat) {
  catModal.value = {
    show: true,
    id: cat?.id || null,
    name: cat?.name || '',
    description: cat?.description || '',
    color: cat?.color || '#edf4ff',
    sort_order: cat?.sort_order ?? 0,
    saving: false,
    error: '',
  }
}

async function loadCategories() {
  catLoading.value = true
  const res = await fetch('/api/admin/v1/forum/categories', {
    headers: { Authorization: `Bearer ${auth.token}` },
  }).catch(() => null)
  if (res?.ok) {
    const data = await res.json()
    categories.value = data.data || []
  }
  catLoading.value = false
}

async function saveCategory() {
  catModal.value.error = ''
  if (!catModal.value.name.trim()) { catModal.value.error = 'Le nom est requis.'; return }
  catModal.value.saving = true
  try {
    const url = catModal.value.id
      ? `/api/admin/v1/forum/categories/${catModal.value.id}`
      : '/api/admin/v1/forum/categories'
    const method = catModal.value.id ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${auth.token}` },
      body: JSON.stringify({
        name: catModal.value.name.trim(),
        description: catModal.value.description.trim(),
        color: catModal.value.color,
        sort_order: catModal.value.sort_order,
      }),
    })
    if (res.ok) {
      catModal.value.show = false
      await loadCategories()
    } else {
      catModal.value.error = 'Une erreur est survenue.'
    }
  } finally {
    catModal.value.saving = false
  }
}

async function deleteCategory(id) {
  if (!confirm('Supprimer cette catégorie ? Toutes ses discussions seront supprimées.')) return
  await fetch(`/api/admin/v1/forum/categories/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  await loadCategories()
}

// --- Threads ---
const threadLoading = ref(true)
const threads = ref([])
const totalThreads = ref(0)
const threadPage = ref(1)
const threadPageSize = 30
const totalThreadPages = computed(() => Math.ceil(totalThreads.value / threadPageSize))

async function loadThreads() {
  threadLoading.value = true
  const res = await fetch(`/api/admin/v1/forum/threads?page=${threadPage.value}`, {
    headers: { Authorization: `Bearer ${auth.token}` },
  }).catch(() => null)
  if (res?.ok) {
    const data = await res.json()
    threads.value = data.data || []
    totalThreads.value = data.total || 0
  }
  threadLoading.value = false
}

async function goThreadPage(p) {
  threadPage.value = p
  await loadThreads()
}

async function togglePin(thread) {
  const res = await fetch(`/api/admin/v1/forum/threads/${thread.id}/pin`, {
    method: 'PUT',
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  if (res.ok) {
    const data = await res.json()
    thread.pinned = data.pinned
  }
}

async function toggleLock(thread) {
  const res = await fetch(`/api/admin/v1/forum/threads/${thread.id}/lock`, {
    method: 'PUT',
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  if (res.ok) {
    const data = await res.json()
    thread.locked = data.locked
  }
}

async function deleteThread(id) {
  if (!confirm('Supprimer cette discussion ? Toutes ses réponses seront supprimées.')) return
  await fetch(`/api/admin/v1/forum/threads/${id}`, {
    method: 'DELETE',
    headers: { Authorization: `Bearer ${auth.token}` },
  })
  await loadThreads()
}

onMounted(() => {
  loadCategories()
  loadThreads()
})
</script>
