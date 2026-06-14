<template>
  <div class="space-y-6">
    <div>
      <h2 class="text-xl sm:text-2xl font-bold text-[#001d32]">Monétisation & Finances</h2>
      <p class="text-sm text-[#40617f] mt-0.5">Revenus de la plateforme, commissions, campagnes et abonnements pro.</p>
    </div>

    <div class="flex gap-1 border-b border-[#e5e7eb]">
      <button v-for="tab in tabs" :key="tab.key" @click="activeTab = tab.key"
        :class="['px-4 py-2.5 text-sm font-semibold transition border-b-2 -mb-px flex items-center gap-2',
          activeTab === tab.key ? 'border-[#006d35] text-[#006d35]' : 'border-transparent text-[#64748b] hover:text-[#001d32]']">
        {{ tab.label }}
        <span v-if="tab.key === 'campaigns' && pendingCount > 0"
          class="bg-[#ef4444] text-white text-[10px] font-bold px-1.5 py-0.5 rounded-full min-w-[18px] text-center">
          {{ pendingCount }}
        </span>
      </button>
    </div>

    <div v-if="activeTab === 'finance'">
      <div v-if="loadingFinance" class="py-16 text-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>
      <template v-else-if="finance">

        <div class="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
          <div class="card p-5">
            <BanknotesIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ formatEUR(finance.platform_revenue_cents) }}</p>
            <p class="text-[#64748b] text-xs">Revenus plateforme</p>
          </div>
          <div class="card p-5">
            <ReceiptPercentIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ formatEUR(finance.transactions.commission_cents) }}</p>
            <p class="text-[#64748b] text-xs">Commissions ({{ finance.commission_rate_percent }}%) · {{ finance.transactions.count }} ventes</p>
          </div>
          <div class="card p-5">
            <StarIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ formatEUR(finance.subscriptions.mrr_cents) }}</p>
            <p class="text-[#64748b] text-xs">Abonnements / mois · {{ finance.subscriptions.active_count }} actifs</p>
          </div>
          <div class="card p-5">
            <MegaphoneIcon class="w-6 h-6 text-[#006d35] mb-2" />
            <p class="font-jakarta font-extrabold text-[#001d32] text-2xl">{{ formatEUR(finance.campaigns.total_cents) }}</p>
            <p class="text-[#64748b] text-xs">Campagnes · {{ finance.campaigns.paid_count }} payées</p>
          </div>
        </div>

        <div class="card p-5 mb-6">
          <h3 class="font-bold text-[#001d32] mb-4">Volume de transactions (prestations)</h3>
          <div class="grid grid-cols-3 gap-4 text-center">
            <div>
              <p class="font-jakarta font-extrabold text-[#001d32] text-xl">{{ formatEUR(finance.transactions.gross_cents) }}</p>
              <p class="text-[#64748b] text-xs mt-1">Montant brut encaissé</p>
            </div>
            <div>
              <p class="font-jakarta font-extrabold text-[#006d35] text-xl">{{ formatEUR(finance.transactions.commission_cents) }}</p>
              <p class="text-[#64748b] text-xs mt-1">Commission UpcycleConnect</p>
            </div>
            <div>
              <p class="font-jakarta font-extrabold text-[#40617f] text-xl">{{ formatEUR(finance.transactions.net_to_providers_cents) }}</p>
              <p class="text-[#64748b] text-xs mt-1">Reversé aux prestataires</p>
            </div>
          </div>
        </div>

        <div class="card overflow-hidden">
          <div class="px-5 py-3 border-b border-[#f1f5f9]">
            <h3 class="font-bold text-[#001d32] text-sm">Dernières transactions</h3>
          </div>
          <div v-if="!transactions.length" class="p-8 text-center text-[#64748b] text-sm">Aucune transaction payée.</div>
          <table v-else class="w-full text-sm">
            <thead class="bg-[#f8fafc] text-[#64748b] text-xs uppercase">
              <tr>
                <th class="text-left px-5 py-2.5">Prestation</th>
                <th class="text-left px-5 py-2.5 hidden sm:table-cell">Client</th>
                <th class="text-right px-5 py-2.5">Montant</th>
                <th class="text-right px-5 py-2.5">Commission</th>
                <th class="text-right px-5 py-2.5 hidden md:table-cell">Net pro</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="t in transactions" :key="t.id" class="border-t border-[#f1f5f9]">
                <td class="px-5 py-2.5 font-medium text-[#001d32]">{{ t.title }}</td>
                <td class="px-5 py-2.5 text-[#40617f] hidden sm:table-cell">{{ t.customer }}</td>
                <td class="px-5 py-2.5 text-right">{{ formatEUR(t.amount_cents) }}</td>
                <td class="px-5 py-2.5 text-right text-[#006d35] font-semibold">{{ formatEUR(t.commission_cents) }}</td>
                <td class="px-5 py-2.5 text-right text-[#40617f] hidden md:table-cell">{{ formatEUR(t.net_cents) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </div>

    <div v-if="activeTab === 'campaigns'">
      <div class="flex items-center gap-2 mb-4">
        <button v-for="f in campaignFilters" :key="f.value"
          @click="campaignStatus = f.value; fetchCampaigns()"
          :class="['px-3 py-1.5 rounded-full text-xs font-semibold border transition',
            campaignStatus === f.value ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]']">
          {{ f.label }}
        </button>
      </div>

      <div v-if="loadingCampaigns" class="py-16 text-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>
      <div v-else-if="!campaigns.length" class="card p-10 text-center text-[#64748b] text-sm">
        Aucune campagne.
      </div>
      <div v-else class="space-y-3">
        <div v-for="c in campaigns" :key="c.id" class="card p-5">
          <div class="flex items-start justify-between gap-4 flex-wrap">
            <div class="flex items-start gap-4 flex-1 min-w-0">
              <div v-if="c.image_url" class="w-16 h-16 rounded-xl overflow-hidden shrink-0">
                <img :src="c.image_url" class="w-full h-full object-cover" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1 flex-wrap">
                  <span :class="badgeClass(c.status)" class="text-xs font-bold px-2 py-0.5 rounded-full">{{ statusLabel(c.status) }}</span>
                  <span class="text-xs text-gray-400">par {{ c.provider_name }}</span>
                  <span v-if="c.paid_at" class="text-xs text-green-600 font-semibold">✓ Payée</span>
                </div>
                <h3 class="font-bold text-[#001d32]">{{ c.title }}</h3>
                <p class="text-[#40617f] text-sm mt-0.5 line-clamp-2">{{ c.description }}</p>
                <p class="text-[#006d35] font-semibold text-sm mt-1">{{ formatEUR(c.budget_cents) }}/mois</p>
                <p v-if="c.rejection_reason" class="text-red-600 text-xs mt-1">Motif de rejet : {{ c.rejection_reason }}</p>
              </div>
            </div>
            <div v-if="c.status === 'pending'" class="flex flex-col gap-2 shrink-0">
              <button @click="validate(c, 'active')"
                class="px-4 py-2 rounded-xl text-white text-sm font-semibold hover:opacity-90 transition" style="background-color:#1B8848;">
                Approuver
              </button>
              <button @click="openReject(c)"
                class="px-4 py-2 rounded-xl text-white text-sm font-semibold hover:opacity-90 transition" style="background-color:#dc2626;">
                Rejeter
              </button>
            </div>
            <div v-else-if="c.status === 'active'" class="shrink-0">
              <button @click="validate(c, 'ended')"
                class="px-4 py-2 rounded-xl border border-gray-300 text-gray-600 text-sm font-semibold hover:bg-gray-50 transition">
                Terminer
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-else-if="activeTab === 'subscriptions'">
      <div v-if="loadingSubs" class="py-16 text-center">
        <div class="w-8 h-8 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>
      <div v-else-if="!subscriptions.length" class="card p-10 text-center text-[#64748b] text-sm">
        Aucun abonnement souscrit.
      </div>
      <div v-else class="card overflow-hidden">
        <table class="w-full text-sm">
          <thead class="bg-[#f8fafc] text-[#64748b] text-xs uppercase">
            <tr>
              <th class="text-left px-5 py-3">Professionnel</th>
              <th class="text-left px-5 py-3">Plan</th>
              <th class="text-left px-5 py-3">Statut</th>
              <th class="text-left px-5 py-3 hidden sm:table-cell">Montant</th>
              <th class="text-left px-5 py-3 hidden md:table-cell">Renouvellement</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in subscriptions" :key="s.id" class="border-t border-[#f1f5f9]">
              <td class="px-5 py-3">
                <p class="font-semibold text-[#001d32]">{{ s.user_name }}</p>
                <p class="text-xs text-gray-400">{{ s.user_email }}</p>
              </td>
              <td class="px-5 py-3 font-medium">{{ s.plan_label }}</td>
              <td class="px-5 py-3">
                <span :class="subBadge(s.status)" class="text-xs font-bold px-2 py-0.5 rounded-full">{{ subLabel(s.status) }}</span>
              </td>
              <td class="px-5 py-3 hidden sm:table-cell">{{ formatEUR(s.amount_cents) }}/mois</td>
              <td class="px-5 py-3 hidden md:table-cell text-[#40617f]">{{ s.current_period_end ? formatDate(s.current_period_end) : '—' }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <Teleport to="body">
      <div v-if="rejectModal.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="rejectModal.show = false" />
        <div class="relative bg-white rounded-2xl shadow-2xl w-full max-w-md p-6">
          <h3 class="font-bold text-[#001d32] text-lg mb-3">Rejeter la campagne</h3>
          <p class="text-sm text-[#40617f] mb-3">« {{ rejectModal.campaign?.title }} »</p>
          <textarea v-model="rejectModal.reason" rows="3" placeholder="Motif du rejet (visible par le pro)..."
            class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none" />
          <div class="flex gap-3 justify-end mt-4">
            <button @click="rejectModal.show = false" class="px-4 py-2 text-sm text-[#40617f] hover:text-[#001d32] font-medium">Annuler</button>
            <button @click="confirmReject" class="px-5 py-2 text-white text-sm font-semibold rounded-xl hover:opacity-90 transition" style="background-color:#dc2626;">
              Confirmer le rejet
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { BanknotesIcon, ReceiptPercentIcon, StarIcon, MegaphoneIcon } from '@heroicons/vue/24/outline'

const BASE = '/api/admin/v1'
function authHeaders() {
  const token = localStorage.getItem('admin_token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const finance = ref(null)
const transactions = ref([])
const loadingFinance = ref(true)
async function fetchFinance() {
  loadingFinance.value = true
  try {
    const [ov, tx] = await Promise.all([
      fetch(`${BASE}/finance/overview`, { headers: authHeaders() }).then(r => r.json()),
      fetch(`${BASE}/finance/transactions`, { headers: authHeaders() }).then(r => r.json()),
    ])
    finance.value = ov
    transactions.value = tx.data || []
  } catch { finance.value = null } finally { loadingFinance.value = false }
}

const activeTab = ref('finance')
const tabs = [
  { key: 'finance', label: 'Finances' },
  { key: 'campaigns', label: 'Campagnes pub' },
  { key: 'subscriptions', label: 'Abonnements' },
]

const campaigns = ref([])
const loadingCampaigns = ref(true)
const campaignStatus = ref('pending')
const campaignFilters = [
  { value: '', label: 'Toutes' },
  { value: 'pending', label: 'À valider' },
  { value: 'active', label: 'Actives' },
  { value: 'rejected', label: 'Rejetées' },
  { value: 'ended', label: 'Terminées' },
]
const pendingCount = ref(0)

const subscriptions = ref([])
const loadingSubs = ref(true)

const rejectModal = ref({ show: false, campaign: null, reason: '' })

async function fetchCampaigns() {
  loadingCampaigns.value = true
  try {
    const params = new URLSearchParams()
    if (campaignStatus.value) params.set('status', campaignStatus.value)
    const res = await fetch(`${BASE}/campaigns?${params}`, { headers: authHeaders() })
    const j = await res.json()
    campaigns.value = j.data || []
  } catch { campaigns.value = [] } finally { loadingCampaigns.value = false }

  try {
    const res = await fetch(`${BASE}/campaigns?status=pending`, { headers: authHeaders() })
    const j = await res.json()
    pendingCount.value = (j.data || []).length
  } catch {}
}

async function fetchSubs() {
  loadingSubs.value = true
  try {
    const res = await fetch(`${BASE}/subscriptions`, { headers: authHeaders() })
    const j = await res.json()
    subscriptions.value = j.data || []
  } catch { subscriptions.value = [] } finally { loadingSubs.value = false }
}

async function validate(c, status) {
  await fetch(`${BASE}/campaigns/${c.id}/status`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify({ status }),
  })
  fetchCampaigns()
}

function openReject(c) { rejectModal.value = { show: true, campaign: c, reason: '' } }
async function confirmReject() {
  const c = rejectModal.value.campaign
  await fetch(`${BASE}/campaigns/${c.id}/status`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json', ...authHeaders() },
    body: JSON.stringify({ status: 'rejected', rejection_reason: rejectModal.value.reason || null }),
  })
  rejectModal.value.show = false
  fetchCampaigns()
}

function formatEUR(cents) { return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format((cents || 0) / 100) }
function formatDate(iso) { return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'short', year: 'numeric' }) }
function statusLabel(s) { return { draft: 'Brouillon', pending: 'À valider', active: 'Active', rejected: 'Rejetée', ended: 'Terminée' }[s] || s }
function badgeClass(s) {
  if (s === 'active') return 'bg-green-100 text-green-800'
  if (s === 'pending') return 'bg-yellow-100 text-yellow-800'
  if (s === 'rejected') return 'bg-red-100 text-red-700'
  if (s === 'ended') return 'bg-gray-100 text-gray-600'
  return 'bg-blue-100 text-blue-700'
}
function subLabel(s) { return { active: 'Actif', cancelled: 'Résilié', inactive: 'Inactif', past_due: 'Impayé' }[s] || s }
function subBadge(s) {
  if (s === 'active') return 'bg-green-100 text-green-800'
  if (s === 'cancelled') return 'bg-gray-100 text-gray-600'
  if (s === 'past_due') return 'bg-red-100 text-red-700'
  return 'bg-yellow-100 text-yellow-800'
}

onMounted(() => { fetchFinance(); fetchCampaigns(); fetchSubs() })
</script>
