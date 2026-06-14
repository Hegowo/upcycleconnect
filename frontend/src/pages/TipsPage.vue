<template>
  <div :class="embedded ? '' : 'max-w-6xl mx-auto py-8 px-4'">

    <div class="flex items-center justify-between flex-wrap gap-3 mb-6">
      <div v-if="!embedded">
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl tracking-tight">Conseils</h1>
        <p class="text-[#40617f] text-sm mt-1">Articles publiés dans l'Espace Conseils côté adhérents.</p>
      </div>
      <div v-else>
        <p class="text-sm text-[#64748b]">{{ tips.length }} conseil(s)</p>
      </div>
      <button
        @click="openForm(null)"
        class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm transition hover:opacity-90"
        style="background: linear-gradient(135deg, #006d35, #1b8848);"
      >
        <PlusIcon class="w-4 h-4" />
        Nouveau conseil
      </button>
    </div>

    <div class="flex items-center gap-2 mb-5">
      <button
        v-for="f in filters"
        :key="f.value"
        @click="statusFilter = f.value; fetchTips()"
        :class="['px-3 py-1.5 rounded-full text-xs font-semibold border transition',
                 statusFilter === f.value ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]']"
      >
        {{ f.label }}
      </button>
    </div>

    <div v-if="loading" class="py-20 text-center">
      <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>

    <div v-else-if="!tips.length" class="bg-white rounded-2xl border border-[#edf4ff] p-12 text-center">
      <LightBulbIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-[#001d32] font-semibold">Aucun conseil</p>
      <p class="text-[#40617f] text-sm mt-1">Crée le premier conseil pour démarrer l'Espace Conseils.</p>
    </div>

    <div v-else class="bg-white rounded-2xl border border-[#edf4ff] overflow-hidden">
      <table class="w-full text-sm">
        <thead class="bg-[#f8fafc] text-[#40617f] text-xs uppercase">
          <tr>
            <th class="text-left px-5 py-3">Titre</th>
            <th class="text-left px-5 py-3 hidden md:table-cell">Catégorie</th>
            <th class="text-left px-5 py-3 hidden md:table-cell">Statut</th>
            <th class="text-left px-5 py-3 hidden lg:table-cell">Auteur</th>
            <th class="text-right px-5 py-3">Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="t in tips" :key="t.id" class="border-t border-[#f1f5f9] hover:bg-[#f8fafc]">
            <td class="px-5 py-3">
              <p class="font-semibold text-[#001d32] line-clamp-1">{{ t.title }}</p>
              <p v-if="t.summary" class="text-[#94a3b8] text-xs line-clamp-1">{{ t.summary }}</p>
            </td>
            <td class="px-5 py-3 hidden md:table-cell text-xs text-[#40617f]">{{ t.category || '—' }}</td>
            <td class="px-5 py-3 hidden md:table-cell">
              <span :class="['text-xs font-semibold px-2 py-0.5 rounded-full', statusBadge(t.status)]">{{ statusLabel(t.status) }}</span>
            </td>
            <td class="px-5 py-3 hidden lg:table-cell text-xs text-[#40617f]">{{ t.author_name || '—' }}</td>
            <td class="px-5 py-3 text-right space-x-1">
              <button @click="openForm(t)" class="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs rounded-lg border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35] transition">
                <PencilSquareIcon class="w-3.5 h-3.5" /> Éditer
              </button>
              <button @click="destroy(t)" class="inline-flex items-center gap-1 px-2.5 py-1.5 text-xs rounded-lg border border-red-200 text-red-600 hover:bg-red-50 transition">
                <TrashIcon class="w-3.5 h-3.5" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <Teleport to="body">
      <div v-if="showForm" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeForm" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-3xl max-h-[92vh] flex flex-col overflow-hidden">
          <div class="px-6 py-4 border-b border-gray-100 flex items-center justify-between shrink-0" style="background: linear-gradient(135deg, #006d35, #1b8848);">
            <h2 class="text-white font-jakarta font-bold text-lg">{{ form.id ? 'Modifier le conseil' : 'Nouveau conseil' }}</h2>
            <button @click="closeForm" class="text-white/70 hover:text-white p-2 rounded-lg hover:bg-white/10">
              <XMarkIcon class="w-5 h-5" />
            </button>
          </div>

          <div class="flex-1 overflow-y-auto p-6 space-y-4">
            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Titre *</label>
              <input v-model="form.title" type="text" maxlength="200"
                class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Catégorie</label>
                <input v-model="form.category" type="text" maxlength="60" placeholder="Ex: Bricolage"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Statut</label>
                <select v-model="form.status"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                  <option value="draft">Brouillon</option>
                  <option value="published">Publié</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">URL image de couverture</label>
              <input v-model="form.cover_image" type="text" placeholder="https://..."
                class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
            </div>

            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Résumé</label>
              <textarea v-model="form.summary" rows="2" maxlength="500" placeholder="Phrase d'accroche affichée dans la liste"
                class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none"></textarea>
            </div>

            <div>
              <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Contenu *</label>
              <TiptapEditor v-model="form.content" placeholder="Rédigez ici le contenu du conseil..." minHeight="200px" />
            </div>

            <p v-if="error" class="text-red-600 text-sm">{{ error }}</p>
          </div>

          <div class="px-6 py-4 border-t border-gray-100 bg-[#f8fafc] flex items-center justify-between gap-3 shrink-0">
            <button @click="closeForm" :disabled="saving"
              class="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700 hover:bg-gray-50 transition disabled:opacity-50">
              Annuler
            </button>
            <button @click="submit" :disabled="saving || !form.title || !form.content"
              class="px-6 py-2.5 rounded-xl text-white font-bold text-sm transition hover:opacity-90 disabled:opacity-50 flex items-center gap-2"
              style="background: linear-gradient(135deg, #006d35, #1b8848);">
              <div v-if="saving" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
              {{ saving ? 'Enregistrement…' : (form.id ? 'Mettre à jour' : 'Créer') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  PlusIcon,
  PencilSquareIcon,
  TrashIcon,
  XMarkIcon,
  LightBulbIcon,
} from '@heroicons/vue/24/outline'
import TiptapEditor from '@/components/TiptapEditor.vue'

defineProps({
  embedded: { type: Boolean, default: false },
})

const BASE = '/api/admin/v1'
function authHeaders() {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const tips = ref([])
const loading = ref(true)
const statusFilter = ref('')
const showForm = ref(false)
const saving = ref(false)
const error = ref('')

const filters = [
  { value: '', label: 'Tous' },
  { value: 'published', label: 'Publiés' },
  { value: 'draft', label: 'Brouillons' },
]

const emptyForm = () => ({
  id: null, title: '', slug: '', summary: '', content: '',
  cover_image: '', category: '', status: 'draft',
})
const form = ref(emptyForm())

async function fetchTips() {
  loading.value = true
  try {
    const params = new URLSearchParams()
    if (statusFilter.value) params.set('status', statusFilter.value)
    const res = await fetch(`${BASE}/tips?${params}`, { headers: authHeaders() })
    if (!res.ok) throw new Error()
    const j = await res.json()
    tips.value = j.data || []
  } catch {
    tips.value = []
  } finally {
    loading.value = false
  }
}

function openForm(t) {
  error.value = ''
  if (t) {
    fetch(`${BASE}/tips/${t.id}`, { headers: authHeaders() })
      .then(r => r.json())
      .then(j => {
        form.value = {
          id: j.id, title: j.title, slug: j.slug, summary: j.summary,
          content: j.content, cover_image: j.cover_image || '',
          category: j.category, status: j.status,
        }
        showForm.value = true
      })
  } else {
    form.value = emptyForm()
    showForm.value = true
  }
}

function closeForm() {
  if (saving.value) return
  showForm.value = false
}

async function submit() {
  saving.value = true
  error.value = ''
  try {
    const payload = { ...form.value }
    if (!payload.cover_image) payload.cover_image = null
    const url = form.value.id ? `${BASE}/tips/${form.value.id}` : `${BASE}/tips`
    const method = form.value.id ? 'PUT' : 'POST'
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json', ...authHeaders() },
      body: JSON.stringify(payload),
    })
    if (!res.ok) {
      const j = await res.json().catch(() => ({}))
      throw new Error(j.message || 'Erreur lors de l\'enregistrement')
    }
    showForm.value = false
    fetchTips()
  } catch (e) {
    error.value = e.message
  } finally {
    saving.value = false
  }
}

async function destroy(t) {
  if (!confirm(`Supprimer le conseil "${t.title}" ?`)) return
  await fetch(`${BASE}/tips/${t.id}`, { method: 'DELETE', headers: authHeaders() })
  fetchTips()
}

function statusLabel(s) { return s === 'published' ? 'Publié' : 'Brouillon' }
function statusBadge(s) {
  return s === 'published'
    ? 'bg-[#dcfce7] text-[#166534]'
    : 'bg-[#fef3c7] text-[#92400e]'
}

onMounted(fetchTips)
</script>
