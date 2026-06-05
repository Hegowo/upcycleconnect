<template>
  <div class="max-w-4xl mx-auto py-10 px-4">

    <RouterLink to="/profil" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
      <ArrowLeftIcon class="w-4 h-4" />
      Retour
    </RouterLink>

    <div class="mb-8">
      <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl sm:text-3xl tracking-tight">Contrats reçus</h1>
      <p class="text-[#40617f] text-sm mt-1">Tous les contrats signés par vos clients sur vos prestations.</p>
    </div>

    <div v-if="loading" class="py-20 text-center">
      <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
    </div>

    <div v-else-if="error" class="p-4 rounded-xl bg-red-50 text-red-600 text-sm border border-red-200">
      {{ error }}
    </div>

    <div v-else-if="!contracts.length" class="bg-white rounded-2xl border border-[#edf4ff] p-12 text-center">
      <DocumentTextIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
      <p class="text-[#001d32] font-semibold">Aucun contrat pour l'instant</p>
      <p class="text-[#40617f] text-sm mt-1">Quand un client signera un contrat sur une de vos prestations, il apparaîtra ici.</p>
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="c in contracts"
        :key="c.id"
        class="bg-white rounded-2xl border border-[#edf4ff] p-5 hover:border-[#006d35] transition"
      >
        <div class="flex items-start justify-between gap-4 flex-wrap">
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1.5">
              <span class="text-xs font-semibold px-2 py-0.5 rounded-full" :class="statusBadgeClass(c.status)">
                {{ statusLabel(c.status) }}
              </span>
              <span class="text-xs text-gray-400 font-mono">{{ c.number }}</span>
            </div>
            <h3 class="font-jakarta font-bold text-[#001d32] text-base truncate">{{ c.prestation_title }}</h3>
            <div class="grid grid-cols-2 sm:grid-cols-3 gap-3 mt-3 text-xs">
              <div>
                <p class="text-gray-400 uppercase font-semibold tracking-wide">Client</p>
                <p class="text-[#001d32] font-medium truncate">{{ c.customer_name }}</p>
              </div>
              <div>
                <p class="text-gray-400 uppercase font-semibold tracking-wide">Email</p>
                <p class="text-[#001d32] font-medium truncate">{{ c.customer_email }}</p>
              </div>
              <div>
                <p class="text-gray-400 uppercase font-semibold tracking-wide">Signé le</p>
                <p class="text-[#001d32] font-medium">{{ formatDate(c.signed_at) }}</p>
              </div>
            </div>
          </div>
          <div class="flex flex-col items-end gap-2 shrink-0">
            <p class="text-xl font-bold text-[#006d35]">{{ formatEUR(c.amount_cents) }}</p>
            <button
              v-if="c.has_pdf"
              @click="downloadContract(c)"
              :disabled="downloadingId === c.id"
              class="inline-flex items-center gap-1.5 px-3 py-2 rounded-lg border border-[#006d35] text-[#006d35] text-xs font-bold transition hover:bg-[#f0fdf4] disabled:opacity-60"
            >
              <ArrowDownTrayIcon class="w-3.5 h-3.5" />
              {{ downloadingId === c.id ? '…' : 'PDF' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeftIcon, ArrowDownTrayIcon, DocumentTextIcon } from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'
import { userApi } from '@/services/publicApi'

const router   = useRouter()
const userAuth = useUserAuthStore()

const contracts     = ref([])
const loading       = ref(true)
const error         = ref('')
const downloadingId = ref(null)

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.push('/connexion?redirect=/profil/contrats-recus')
    return
  }
  try {
    const data = await userApi('/provider/contracts')
    contracts.value = data.data || []
  } catch (e) {
    error.value = e.message || 'Impossible de charger les contrats.'
  } finally {
    loading.value = false
  }
})

async function downloadContract(c) {
  downloadingId.value = c.id
  try {
    const res = await fetch(`/api/v1/contracts/${c.id}/download`, {
      headers: { Authorization: `Bearer ${userAuth.token}` },
    })
    if (!res.ok) throw new Error('failed')
    const blob = await res.blob()
    const url  = URL.createObjectURL(blob)
    const a    = document.createElement('a')
    a.href     = url
    a.download = `${c.number}.pdf`
    a.click()
    URL.revokeObjectURL(url)
  } catch {} finally {
    downloadingId.value = null
  }
}

function formatEUR(cents) {
  return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format((cents || 0) / 100)
}
function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: 'long', year: 'numeric' })
}
function statusLabel(s) {
  return { active: 'Payé', signed_pending_payment: 'En attente de paiement', cancelled: 'Annulé' }[s] || s
}
function statusBadgeClass(s) {
  if (s === 'active') return 'bg-green-100 text-green-800'
  if (s === 'signed_pending_payment') return 'bg-yellow-100 text-yellow-800'
  return 'bg-gray-100 text-gray-600'
}
</script>
