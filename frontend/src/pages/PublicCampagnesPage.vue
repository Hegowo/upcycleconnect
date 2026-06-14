<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/profil/pro" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" /> Retour profil pro
      </RouterLink>

      <div class="flex items-center justify-between flex-wrap gap-3 mb-8">
        <div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Campagnes publicitaires</h1>
          <p class="text-[#40617f] text-sm mt-1">Mettez en avant vos projets auprès de la communauté (100–500 €/mois).</p>
        </div>
        <button @click="openForm(null)" :disabled="!hasPremium"
          class="inline-flex items-center gap-2 px-4 py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 transition disabled:opacity-40 disabled:cursor-not-allowed" style="background:linear-gradient(135deg,#006d35,#1b8848);">
          <PlusIcon class="w-4 h-4" /> Nouvelle campagne
        </button>
      </div>

      <div v-if="success" class="bg-green-50 border border-green-200 rounded-2xl p-4 mb-6 flex items-center gap-3">
        <CheckCircleIcon class="w-5 h-5 text-green-600 shrink-0" />
        <p class="text-green-800 text-sm font-medium">Campagne soumise et paiement confirmé ! Elle sera activée après validation par l'équipe.</p>
      </div>

      <div v-if="!loading && !hasPremium" class="bg-white border-2 border-[#006d35] rounded-2xl p-5 mb-6 flex items-center justify-between gap-4 flex-wrap">
        <div class="flex items-center gap-3">
          <div class="w-11 h-11 rounded-xl flex items-center justify-center" style="background:linear-gradient(135deg,#006d35,#1b8848);">
            <LockClosedIcon class="w-5 h-5 text-white" />
          </div>
          <div>
            <p class="font-jakarta font-bold text-[#001d32] text-sm">Fonctionnalité Premium</p>
            <p class="text-[#40617f] text-xs mt-0.5">La création de campagnes publicitaires nécessite un abonnement Premium.</p>
          </div>
        </div>
        <RouterLink to="/profil/pro/abonnement" class="px-4 py-2 rounded-xl text-white text-sm font-bold hover:opacity-90 transition shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">
          Passer à Premium →
        </RouterLink>
      </div>

      <div v-if="loading" class="py-20 text-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>

      <div v-else-if="!campaigns.length" class="bg-white rounded-2xl border border-[#edf4ff] p-12 text-center">
        <MegaphoneIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <p class="font-semibold text-[#001d32]">Aucune campagne</p>
        <p class="text-[#40617f] text-sm mt-1">Crée ta première campagne pour te mettre en avant.</p>
      </div>

      <div v-else class="space-y-4">
        <div v-for="c in campaigns" :key="c.id" class="bg-white rounded-2xl border border-[#edf4ff] p-5 flex items-start justify-between gap-4 flex-wrap hover:border-[#006d35] transition">
          <div class="flex items-start gap-4 flex-1 min-w-0">
            <div v-if="c.image_url" class="w-16 h-16 rounded-xl overflow-hidden shrink-0">
              <img :src="c.image_url" class="w-full h-full object-cover" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span :class="statusBadge(c.status)" class="text-xs font-bold px-2 py-0.5 rounded-full">{{ statusLabel(c.status) }}</span>
              </div>
              <h3 class="font-jakarta font-bold text-[#001d32] text-base truncate">{{ c.title }}</h3>
              <p class="text-[#40617f] text-xs mt-0.5 line-clamp-2">{{ c.description }}</p>
              <p class="text-[#006d35] font-semibold text-sm mt-1">{{ formatEUR(c.budget_cents) }}/mois</p>
            </div>
          </div>
          <div class="flex flex-col gap-2 shrink-0">
            <button v-if="['draft','rejected'].includes(c.status)" @click="openForm(c)"
              class="inline-flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg border border-[#e5e7eb] hover:border-[#006d35] hover:text-[#006d35] transition">
              <PencilSquareIcon class="w-3.5 h-3.5" /> Éditer
            </button>
            <button v-if="['draft','rejected'].includes(c.status)" @click="submit(c)"
              class="inline-flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg text-white font-bold transition hover:opacity-90"
              style="background: linear-gradient(135deg,#006d35,#1b8848);">
              Soumettre & Payer
            </button>
            <button v-if="['draft','rejected'].includes(c.status)" @click="destroy(c)"
              class="inline-flex items-center gap-1 px-3 py-1.5 text-xs rounded-lg border border-red-200 text-red-600 hover:bg-red-50 transition">
              <TrashIcon class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>

      <Teleport to="body">
        <div v-if="showForm" class="fixed inset-0 z-50 flex items-center justify-center p-4">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeForm" />
          <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-lg max-h-[92vh] flex flex-col overflow-hidden">
            <div class="px-6 py-4 border-b flex items-center justify-between shrink-0" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <h2 class="text-white font-jakarta font-bold text-lg">{{ form.id ? 'Modifier' : 'Nouvelle campagne' }}</h2>
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
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Description</label>
                <textarea v-model="form.description" rows="3"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Image (URL)</label>
                <input v-model="form.image_url" type="text"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Budget mensuel (100–500 €) *</label>
                <input v-model.number="form.budget_euros" type="number" min="100" max="500"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>

              <div v-if="myPrestations.length">
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Prestations à mettre en avant</label>
                <div class="space-y-1.5 max-h-36 overflow-y-auto bg-[#f8fafc] rounded-xl p-3 border border-[#e5e7eb]">
                  <label v-for="p in myPrestations" :key="p.id" class="flex items-center gap-2.5 text-sm cursor-pointer hover:bg-white rounded-lg px-2 py-1 transition">
                    <input type="checkbox" :value="p.id" v-model="form.prestation_ids" class="w-4 h-4 accent-[#006d35]" />
                    <span class="text-[#001d32] truncate">{{ p.title }}</span>
                  </label>
                </div>
              </div>

              <div v-if="myEvents.length">
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1.5">Événements à mettre en avant</label>
                <div class="space-y-1.5 max-h-36 overflow-y-auto bg-[#f8fafc] rounded-xl p-3 border border-[#e5e7eb]">
                  <label v-for="e in myEvents" :key="e.id" class="flex items-center gap-2.5 text-sm cursor-pointer hover:bg-white rounded-lg px-2 py-1 transition">
                    <input type="checkbox" :value="e.id" v-model="form.event_ids" class="w-4 h-4 accent-[#006d35]" />
                    <span class="text-[#001d32] truncate">{{ e.title }}</span>
                  </label>
                </div>
              </div>

              <p v-if="!myPrestations.length && !myEvents.length" class="text-xs text-[#94a3b8]">
                Crée d'abord des prestations ou événements publiés pour pouvoir les mettre en avant dans cette campagne.
              </p>
              <p v-if="formError" class="text-red-600 text-sm">{{ formError }}</p>
            </div>
            <div class="px-6 py-4 border-t bg-[#f8fafc] flex items-center justify-between gap-3 shrink-0">
              <button @click="closeForm" class="px-4 py-2.5 rounded-xl border border-gray-200 bg-white text-sm font-semibold text-gray-700 hover:bg-gray-50 transition">Annuler</button>
              <button @click="save" :disabled="saving" class="px-6 py-2.5 rounded-xl text-white font-bold text-sm flex items-center gap-2 transition hover:opacity-90 disabled:opacity-50"
                style="background:linear-gradient(135deg,#006d35,#1b8848);">
                <div v-if="saving" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                {{ saving ? 'Enregistrement...' : form.id ? 'Mettre à jour' : 'Créer' }}
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
import { RouterLink, useRoute } from 'vue-router'
import { ArrowLeftIcon, PlusIcon, PencilSquareIcon, TrashIcon, XMarkIcon, CheckCircleIcon, MegaphoneIcon, LockClosedIcon } from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const route = useRoute()
const campaigns = ref([])
const loading = ref(true)
const showForm = ref(false)
const saving = ref(false)
const formError = ref('')
const success = ref(route.query.success === '1')

const emptyForm = () => ({ id: null, title: '', description: '', image_url: '', budget_euros: 100, prestation_ids: [], event_ids: [] })
const form = ref(emptyForm())
const hasPremium = ref(true)
const myPrestations = ref([])
const myEvents = ref([])

async function fetchCampaigns() {
  loading.value = true
  try {
    const [campData, subData, prestaData, eventData] = await Promise.all([
      userApi('/campaigns'),
      userApi('/subscription').catch(() => ({ subscription: null })),
      userApi('/provider/prestations?status=published').catch(() => ({ data: [] })),
      userApi('/provider/events?status=published').catch(() => ({ data: [] })),
    ])
    campaigns.value = campData.data || []
    hasPremium.value = subData.subscription?.plan === 'premium' && subData.subscription?.status === 'active'
    myPrestations.value = prestaData.data || []
    myEvents.value = eventData.data || []
  } finally { loading.value = false }
}

function openForm(c) {
  formError.value = ''
  form.value = c ? {
    id: c.id, title: c.title, description: c.description || '',
    image_url: c.image_url || '', budget_euros: c.budget_cents / 100,
    prestation_ids: (c.items || []).filter(i => i.type === 'prestation').map(i => i.id),
    event_ids: (c.items || []).filter(i => i.type === 'event').map(i => i.id),
  } : emptyForm()
  showForm.value = true
}
function closeForm() { if (!saving.value) showForm.value = false }

async function save() {
  if (!form.value.title) { formError.value = 'Titre requis'; return }
  if (form.value.budget_euros < 100 || form.value.budget_euros > 500) { formError.value = 'Budget entre 100 et 500 €'; return }
  saving.value = true
  formError.value = ''
  try {
    const payload = { ...form.value, image_url: form.value.image_url || null }
    if (form.value.id) await userApi(`/campaigns/${form.value.id}`, { method: 'PUT', body: JSON.stringify(payload) })
    else await userApi('/campaigns', { method: 'POST', body: JSON.stringify(payload) })
    showForm.value = false
    await fetchCampaigns()
  } catch (e) { formError.value = e.message } finally { saving.value = false }
}

async function submit(c) {
  try {
    const res = await userApi(`/campaigns/${c.id}/submit`, { method: 'POST' })
    if (res.checkout_url) window.location.href = res.checkout_url
    else { await fetchCampaigns() }
  } catch (e) { alert(e.message) }
}

async function destroy(c) {
  if (!confirm(`Supprimer la campagne "${c.title}" ?`)) return
  await userApi(`/campaigns/${c.id}`, { method: 'DELETE' })
  await fetchCampaigns()
}

function formatEUR(cents) {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100)
}
function statusLabel(s) {
  return { draft: 'Brouillon', pending: 'En attente', active: 'Active', rejected: 'Rejetée', ended: 'Terminée' }[s] || s
}
function statusBadge(s) {
  if (s === 'active') return 'bg-green-100 text-green-800'
  if (s === 'pending') return 'bg-yellow-100 text-yellow-800'
  if (s === 'rejected') return 'bg-red-100 text-red-700'
  if (s === 'ended') return 'bg-gray-100 text-gray-600'
  return 'bg-blue-100 text-blue-700'
}

onMounted(fetchCampaigns)
</script>
