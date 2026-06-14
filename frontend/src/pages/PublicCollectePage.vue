<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-lg mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/profil/pro" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" /> Retour profil pro
      </RouterLink>

      <div class="mb-8">
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">Collecte d'objets</h1>
        <p class="text-[#40617f] text-sm mt-1">Saisissez le code-barres d'un bon de dépôt pour marquer l'objet comme collecté.</p>
      </div>

      <div class="bg-white rounded-[24px] border border-[#edf4ff] p-8">
        <div class="flex flex-col gap-5">

          <div class="flex flex-col items-center gap-3 mb-2">
            <div class="w-16 h-16 rounded-2xl flex items-center justify-center" style="background:linear-gradient(135deg,#006d35,#1b8848);">
              <QrCodeIcon class="w-8 h-8 text-white" />
            </div>
            <p class="text-[#40617f] text-sm text-center max-w-xs">
              Entrez le code imprimé sur le bon de dépôt du particulier (format <span class="font-mono font-semibold">UP-123-ABCDEF</span>).
            </p>
          </div>

          <div>
            <label class="block text-xs font-semibold text-[#40617f] uppercase mb-2">Code-barres du bon de dépôt</label>
            <input
              v-model="barcode"
              type="text"
              placeholder="UP-123-ABCDEF"
              class="w-full px-4 py-3 bg-[#f8fafc] border-2 border-[#e5e7eb] rounded-xl text-base font-mono focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition uppercase tracking-widest"
              @keydown.enter="collect"
              autofocus
            />
          </div>

          <button
            @click="collect"
            :disabled="!barcode.trim() || loading"
            class="w-full py-3.5 rounded-xl text-white font-jakarta font-bold text-base transition hover:opacity-90 disabled:opacity-50 flex items-center justify-center gap-2"
            style="background:linear-gradient(135deg,#006d35,#1b8848);"
          >
            <div v-if="loading" class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
            <CheckCircleIcon v-else class="w-5 h-5" />
            {{ loading ? 'Vérification...' : 'Valider la collecte' }}
          </button>

          <div v-if="result && !error" class="bg-[#f0fdf4] rounded-2xl border border-[#bbf7d0] p-5">
            <div class="flex items-center gap-3 mb-3">
              <CheckCircleIcon class="w-6 h-6 text-[#006d35] shrink-0" />
              <p class="font-jakarta font-bold text-[#006d35]">Objet collecté avec succès !</p>
            </div>
            <div class="space-y-1 text-sm text-[#001d32]">
              <p><span class="text-[#40617f]">Objet :</span> <strong>{{ result.deposit.title }}</strong></p>
              <p><span class="text-[#40617f]">Référence :</span> <span class="font-mono">#UP-{{ result.deposit.id }}</span></p>
              <p v-if="result.deposit.estimated_weight"><span class="text-[#40617f]">Poids :</span> {{ result.deposit.estimated_weight }} kg</p>
              <p v-if="result.deposit.category"><span class="text-[#40617f]">Catégorie :</span> {{ result.deposit.category?.name }}</p>
            </div>
          </div>

          <div v-if="error" class="bg-red-50 rounded-2xl border border-red-200 p-5 flex items-start gap-3">
            <ExclamationCircleIcon class="w-5 h-5 text-red-500 shrink-0 mt-0.5" />
            <p class="text-red-700 text-sm">{{ error }}</p>
          </div>

        </div>
      </div>

      <div v-if="recent.length" class="mt-8">
        <h2 class="font-jakarta font-bold text-[#001d32] text-lg mb-4">Collectes récentes</h2>
        <div class="space-y-3">
          <div v-for="d in recent" :key="d.id" class="bg-white rounded-xl border border-[#edf4ff] p-4 flex items-center justify-between gap-4">
            <div class="flex-1 min-w-0">
              <p class="font-semibold text-[#001d32] text-sm truncate">{{ d.title }}</p>
              <p class="text-[#40617f] text-xs font-mono">#UP-{{ d.id }} · {{ d.qr_code }}</p>
            </div>
            <span class="text-xs font-bold px-2 py-0.5 rounded-full bg-[#dcfce7] text-[#166534]">Collecté</span>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowLeftIcon, QrCodeIcon, CheckCircleIcon, ExclamationCircleIcon } from '@heroicons/vue/24/outline'
import { userApi } from '@/services/publicApi'

const barcode = ref('')
const loading = ref(false)
const result  = ref(null)
const error   = ref('')
const recent  = ref([])

async function collect() {
  if (!barcode.value.trim() || loading.value) return
  loading.value = true
  result.value  = null
  error.value   = ''
  try {
    const res = await userApi('/deposits/collect', {
      method: 'POST',
      body: JSON.stringify({ barcode: barcode.value.trim().toUpperCase() }),
    })
    result.value = res

    if (!recent.value.find(r => r.id === res.deposit.id)) {
      recent.value.unshift(res.deposit)
    }
    barcode.value = ''
  } catch (e) {
    error.value = e.message || 'Erreur lors de la collecte.'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  try {
    const data = await userApi('/deposits?status=collected')
    recent.value = (data.data || []).slice(0, 10)
  } catch {}
})
</script>
