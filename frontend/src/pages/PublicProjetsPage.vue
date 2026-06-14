<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-5xl mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/profil/pro" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" /> Retour profil pro
      </RouterLink>

      <div class="flex items-center justify-between flex-wrap gap-3 mb-8">
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Projets d'upcycling</h1>
          <p class="text-[#40617f] text-sm mt-1">Documentez vos projets de transformation et partagez-les avec la communauté.</p>
        </div>
        <button @click="openProjectForm(null)" class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
          <PlusIcon class="w-4 h-4" /> Nouveau projet
        </button>
      </div>

      <div v-if="loading" class="py-20 text-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>

      <div v-else-if="!projects.length" class="bg-white rounded-2xl border border-[#edf4ff] p-12 text-center">
        <WrenchScrewdriverIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <p class="font-semibold text-[#001d32]">Aucun projet</p>
        <p class="text-[#40617f] text-sm mt-1">Crée ton premier projet pour documenter ton travail de transformation.</p>
      </div>

      <div v-else class="space-y-6">
        <div v-for="p in projects" :key="p.id" class="bg-white rounded-[24px] border border-[#edf4ff] overflow-hidden">

          <div class="p-5 flex items-start justify-between gap-4 flex-wrap border-b border-[#f1f5f9]">
            <div class="flex items-start gap-4 flex-1 min-w-0">
              <div v-if="p.cover_image" class="w-16 h-16 rounded-xl overflow-hidden shrink-0">
                <img :src="p.cover_image" class="w-full h-full object-cover" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1 flex-wrap">
                  <span :class="projectStatusBadge(p.status)" class="text-xs font-bold px-2 py-0.5 rounded-full">{{ projectStatusLabel(p.status) }}</span>
                  <span v-if="p.category" class="text-xs text-[#40617f] bg-[#f7f9ff] px-2 py-0.5 rounded-full">{{ p.category }}</span>
                </div>
                <h3 class="font-jakarta font-bold text-[#001d32] text-lg">{{ p.title }}</h3>
                <p class="text-[#40617f] text-sm mt-0.5 line-clamp-2">{{ p.description }}</p>
                <p v-if="p.impact_label" class="text-[#006d35] text-xs font-semibold mt-1">🌱 {{ p.impact_label }}</p>
              </div>
            </div>
            <div class="flex gap-2 shrink-0">
              <button @click="openProjectForm(p)" class="inline-flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35] transition">
                <PencilSquareIcon class="w-3.5 h-3.5" /> Éditer
              </button>
              <button @click="addStep(p)" class="inline-flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg text-white font-bold transition hover:opacity-90" style="background:linear-gradient(135deg,#006d35,#1b8848);">
                <PlusIcon class="w-3.5 h-3.5" /> Étape
              </button>
              <button @click="deleteProject(p)" class="inline-flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg border border-red-200 text-red-600 hover:bg-red-50 transition">
                <TrashIcon class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>

          <div v-if="p.steps && p.steps.length" class="p-5">
            <p class="text-xs font-semibold uppercase tracking-wider text-[#40617f] mb-3">Étapes de transformation</p>
            <div class="space-y-3">
              <div v-for="(step, idx) in p.steps" :key="step.id"
                class="flex items-start gap-3 p-3 rounded-xl"
                :class="step.completed_at ? 'bg-[#f0fdf4]' : 'bg-[#f8fafc]'"
              >
                <div class="w-6 h-6 rounded-full flex items-center justify-center shrink-0 text-xs font-bold"
                  :class="step.completed_at ? 'bg-[#006d35] text-white' : 'bg-[#e5e7eb] text-[#40617f]'">
                  {{ idx + 1 }}
                </div>
                <div class="flex-1 min-w-0">
                  <p class="font-semibold text-[#001d32] text-sm">{{ step.title }}</p>
                  <p v-if="step.description" class="text-[#40617f] text-xs mt-0.5">{{ step.description }}</p>
                </div>
                <div class="flex gap-1 shrink-0">
                  <button @click="editStep(p, step)" class="p-1.5 rounded-lg hover:bg-white transition text-[#40617f] hover:text-[#006d35]">
                    <PencilSquareIcon class="w-3.5 h-3.5" />
                  </button>
                  <button @click="deleteStep(p, step)" class="p-1.5 rounded-lg hover:bg-red-50 transition text-red-400">
                    <TrashIcon class="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="px-5 pb-5">
            <p class="text-[#94a3b8] text-sm">Pas encore d'étapes — clique sur "+ Étape" pour commencer à documenter ce projet.</p>
          </div>
        </div>
      </div>

      <Teleport to="body">
        <div v-if="showProjectForm" class="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeProjectForm" />
          <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[92vh] flex flex-col overflow-hidden">
            <div class="px-6 py-4 border-b flex items-center justify-between shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <h2 class="text-white font-jakarta font-bold text-lg">{{ projectForm.id ? 'Modifier le projet' : 'Nouveau projet' }}</h2>
              <button @click="closeProjectForm" class="text-white/70 hover:text-white p-2 rounded-lg hover:bg-white/10">
                <XMarkIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="flex-1 overflow-y-auto p-6 space-y-4">
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Titre *</label>
                <input v-model="projectForm.title" type="text" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" /></div>
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Description</label>
                <textarea v-model="projectForm.description" rows="3" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" /></div>
              <div class="grid grid-cols-2 gap-3">
                <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Catégorie</label>
                  <input v-model="projectForm.category" type="text" placeholder="Ex: Textile" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" /></div>
                <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Statut</label>
                  <select v-model="projectForm.status" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30">
                    <option value="in_progress">En cours</option>
                    <option value="completed">Terminé</option>
                    <option value="showcased">Mis en avant</option>
                  </select></div>
              </div>
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Image de couverture (URL)</label>
                <input v-model="projectForm.cover_image" type="text" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" /></div>
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Impact écologique</label>
                <input v-model="projectForm.impact_label" type="text" placeholder="Ex: 150 kg de déchets évités" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" /></div>
              <p v-if="projectFormError" class="text-red-600 text-sm">{{ projectFormError }}</p>
            </div>
            <div class="px-6 py-4 border-t bg-[#f8fafc] flex items-center justify-between gap-3 shrink-0">
              <button @click="closeProjectForm" class="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700">Annuler</button>
              <button @click="saveProject" :disabled="savingProject" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm flex items-center gap-2 hover:opacity-90 disabled:opacity-50 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
                <div v-if="savingProject" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                {{ savingProject ? 'Enregistrement...' : projectForm.id ? 'Mettre à jour' : 'Créer' }}
              </button>
            </div>
          </div>
        </div>

        <div v-if="showStepForm" class="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeStepForm" />
          <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md max-h-[90vh] flex flex-col overflow-hidden">
            <div class="px-6 py-4 border-b flex items-center justify-between shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <h2 class="text-white font-jakarta font-bold text-lg">{{ stepForm.id ? 'Modifier l\'étape' : 'Nouvelle étape' }}</h2>
              <button @click="closeStepForm" class="text-white/70 hover:text-white p-2 rounded-lg hover:bg-white/10"><XMarkIcon class="w-5 h-5" /></button>
            </div>
            <div class="flex-1 overflow-y-auto p-6 space-y-4">
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Titre *</label>
                <input v-model="stepForm.title" type="text" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" /></div>
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Description</label>
                <textarea v-model="stepForm.description" rows="3" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" /></div>
              <div><label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Image (URL)</label>
                <input v-model="stepForm.image_url" type="text" class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" /></div>
              <div class="flex items-center gap-3">
                <input v-model="stepForm.completed" type="checkbox" id="stepCompleted" class="w-4 h-4 accent-[#006d35]" />
                <label for="stepCompleted" class="text-sm text-[#001d32] font-medium">Étape complétée</label>
              </div>
            </div>
            <div class="px-6 py-4 border-t bg-[#f8fafc] flex items-center justify-between gap-3 shrink-0">
              <button @click="closeStepForm" class="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700">Annuler</button>
              <button @click="saveStep" :disabled="savingStep" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm flex items-center gap-2 hover:opacity-90 disabled:opacity-50 transition" style="background:linear-gradient(135deg,#006d35,#1b8848);">
                <div v-if="savingStep" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                {{ savingStep ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </div>
          </div>
        </div>
      </Teleport>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowLeftIcon, PlusIcon, PencilSquareIcon, TrashIcon, XMarkIcon, WrenchScrewdriverIcon } from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const projects = ref([])
const loading = ref(true)
const showProjectForm = ref(false)
const savingProject = ref(false)
const projectFormError = ref('')
const showStepForm = ref(false)
const savingStep = ref(false)
let currentProjectId = null

const emptyProject = () => ({ id: null, title: '', description: '', category: '', cover_image: '', impact_label: '', status: 'in_progress' })
const emptyStep    = () => ({ id: null, title: '', description: '', image_url: '', completed: false, step_order: 0 })
const projectForm = ref(emptyProject())
const stepForm    = ref(emptyStep())

async function fetchProjects() {
  loading.value = true
  try { projects.value = (await userApi('/projects')).data || [] }
  finally { loading.value = false }
}

function openProjectForm(p) {
  projectFormError.value = ''
  projectForm.value = p ? {
    id: p.id, title: p.title, description: p.description || '', category: p.category || '',
    cover_image: p.cover_image || '', impact_label: p.impact_label || '', status: p.status,
  } : emptyProject()
  showProjectForm.value = true
}
function closeProjectForm() { if (!savingProject.value) showProjectForm.value = false }
async function saveProject() {
  if (!projectForm.value.title) { projectFormError.value = 'Titre requis'; return }
  savingProject.value = true
  projectFormError.value = ''
  try {
    const payload = { ...projectForm.value, cover_image: projectForm.value.cover_image || null, impact_label: projectForm.value.impact_label || null }
    if (projectForm.value.id) await userApi(`/projects/${projectForm.value.id}`, { method: 'PUT', body: JSON.stringify(payload) })
    else await userApi('/projects', { method: 'POST', body: JSON.stringify(payload) })
    showProjectForm.value = false
    await fetchProjects()
  } catch (e) { projectFormError.value = e.message } finally { savingProject.value = false }
}
async function deleteProject(p) {
  if (!confirm(`Supprimer le projet "${p.title}" et toutes ses étapes ?`)) return
  await userApi(`/projects/${p.id}`, { method: 'DELETE' })
  await fetchProjects()
}

function addStep(p) {
  currentProjectId = p.id
  stepForm.value = { ...emptyStep(), step_order: (p.steps?.length || 0) }
  showStepForm.value = true
}
function editStep(p, step) {
  currentProjectId = p.id
  stepForm.value = { id: step.id, title: step.title, description: step.description || '', image_url: step.image_url || '', completed: !!step.completed_at, step_order: step.step_order }
  showStepForm.value = true
}
function closeStepForm() { if (!savingStep.value) showStepForm.value = false }
async function saveStep() {
  if (!stepForm.value.title) return
  savingStep.value = true
  try {
    const payload = { ...stepForm.value, image_url: stepForm.value.image_url || null }
    if (stepForm.value.id) await userApi(`/projects/${currentProjectId}/steps/${stepForm.value.id}`, { method: 'PUT', body: JSON.stringify(payload) })
    else await userApi(`/projects/${currentProjectId}/steps`, { method: 'POST', body: JSON.stringify(payload) })
    showStepForm.value = false
    await fetchProjects()
  } catch (e) { alert(e.message) } finally { savingStep.value = false }
}
async function deleteStep(p, step) {
  if (!confirm(`Supprimer l'étape "${step.title}" ?`)) return
  await userApi(`/projects/${p.id}/steps/${step.id}`, { method: 'DELETE' })
  await fetchProjects()
}

function projectStatusLabel(s) { return { in_progress: 'En cours', completed: 'Terminé', showcased: 'Mis en avant' }[s] || s }
function projectStatusBadge(s) {
  if (s === 'showcased') return 'bg-[#dcfce7] text-[#166534]'
  if (s === 'completed') return 'bg-blue-100 text-blue-700'
  return 'bg-yellow-100 text-yellow-800'
}

onMounted(fetchProjects)
</script>
