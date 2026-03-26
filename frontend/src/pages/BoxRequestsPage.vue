<template>
  <div class="space-y-6">
    <!-- Dev banner -->
    <div class="flex items-center gap-3 px-4 py-3 rounded-xl border" style="background:#fef9c3; border-color:#fde047;">
      <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0" style="background:#fde047;">
        <WrenchScrewdriverIcon class="w-4 h-4" style="color:#854d0e;" />
      </div>
      <div>
        <p class="text-sm font-semibold" style="color:#854d0e;">Fonctionnalité en cours de développement — Interface preview</p>
        <p class="text-xs" style="color:#a16207;">Les données affichées sont des maquettes. L'API backend n'est pas encore disponible.</p>
      </div>
    </div>

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-[#001d32]">Demandes de Dépôt</h2>
        <p class="text-sm text-[#40617f] mt-0.5">Validez les demandes de dépôt d'objets des utilisateurs</p>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-sm font-semibold px-3 py-1 rounded-full" style="background:#fef9c3; color:#854d0e;">
          {{ mockRequests.filter(r => r.status === 'pending').length }} en attente
        </span>
      </div>
    </div>

    <!-- Split layout -->
    <div class="grid grid-cols-1 xl:grid-cols-5 gap-6">

      <!-- Left col: Queue (2/5) -->
      <div class="xl:col-span-2">
        <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          <div class="px-5 py-4 border-b border-[#f1f5f9]">
            <h3 class="font-semibold text-[#001d32]">File d'attente</h3>
            <p class="text-xs text-gray-400 mt-0.5">{{ mockRequests.length }} demandes au total</p>
          </div>
          <div class="divide-y divide-[#f8fafc]">
            <div
              v-for="req in mockRequests"
              :key="req.id"
              @click="selectedRequest = req"
              class="flex items-start gap-3 p-4 cursor-pointer transition-colors"
              :class="selectedRequest?.id === req.id ? 'bg-[#f0fdf4] border-l-2 border-[#006d35]' : 'hover:bg-[#f8fafc] border-l-2 border-transparent'"
            >
              <div class="w-10 h-10 rounded-xl overflow-hidden shrink-0 bg-gray-100 flex items-center justify-center">
                <PhotoIcon class="w-5 h-5 text-gray-400" />
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-semibold text-[#001d32] truncate">{{ req.object_name }}</p>
                <p class="text-xs text-gray-400 mt-0.5">{{ req.user_name }} · {{ req.category }}</p>
                <div class="flex items-center justify-between mt-1.5">
                  <span class="text-[10px] text-gray-400">{{ req.date }}</span>
                  <span class="text-xs font-semibold px-2 py-0.5 rounded-full" :class="requestStatusBadge(req.status)">
                    {{ requestStatusText(req.status) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Right col: Detail (3/5) -->
      <div class="xl:col-span-3">
        <div v-if="!selectedRequest" class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm flex items-center justify-center h-96">
          <div class="text-center text-gray-400">
            <InboxIcon class="w-12 h-12 mx-auto mb-3 text-gray-300" />
            <p class="text-sm font-medium">Sélectionnez une demande</p>
            <p class="text-xs mt-1">Cliquez sur une demande pour voir les détails</p>
          </div>
        </div>

        <div v-else class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          <!-- Header -->
          <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
            <div>
              <h3 class="font-semibold text-[#001d32]">{{ selectedRequest.object_name }}</h3>
              <p class="text-xs text-gray-400 mt-0.5">Demande #{{ selectedRequest.id }} · {{ selectedRequest.date }}</p>
            </div>
            <span class="text-xs font-semibold px-2.5 py-1 rounded-full" :class="requestStatusBadge(selectedRequest.status)">
              {{ requestStatusText(selectedRequest.status) }}
            </span>
          </div>

          <div class="p-6 space-y-5">
            <!-- Photos mock -->
            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Photos de l'objet</p>
              <div class="grid grid-cols-3 gap-2">
                <div v-for="n in 3" :key="n" class="aspect-square rounded-xl bg-[#f8fafc] border border-[#e5e7eb] flex items-center justify-center">
                  <div class="text-center">
                    <PhotoIcon class="w-6 h-6 text-gray-300 mx-auto" />
                    <p class="text-[10px] text-gray-300 mt-1">Photo {{ n }}</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Info grid -->
            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-3">Informations</p>
              <div class="grid grid-cols-2 gap-3">
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Utilisateur</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.user_name }}</p>
                </div>
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Catégorie</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.category }}</p>
                </div>
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Poids estimé</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.weight }}</p>
                </div>
                <div class="bg-[#f8fafc] rounded-xl p-3">
                  <p class="text-xs text-gray-400 mb-1">Dimensions</p>
                  <p class="text-sm font-semibold text-[#001d32]">{{ selectedRequest.dimensions }}</p>
                </div>
              </div>
            </div>

            <!-- Description -->
            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Description</p>
              <p class="text-sm text-gray-600 bg-[#f8fafc] rounded-xl p-3 leading-relaxed">{{ selectedRequest.description }}</p>
            </div>

            <!-- Impact metrics -->
            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-3">Métriques d'Impact</p>
              <div class="grid grid-cols-3 gap-3">
                <div class="rounded-xl p-3 text-center" style="background:#f0fdf4;">
                  <p class="text-lg font-bold" style="color:#006d35;">{{ selectedRequest.co2_saved }}</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">CO₂ économisé</p>
                </div>
                <div class="rounded-xl p-3 text-center" style="background:#dbeafe;">
                  <p class="text-lg font-bold text-blue-700">{{ selectedRequest.water_saved }}</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Eau préservée</p>
                </div>
                <div class="rounded-xl p-3 text-center" style="background:#fef9c3;">
                  <p class="text-lg font-bold" style="color:#854d0e;">{{ selectedRequest.score }}</p>
                  <p class="text-[10px] text-gray-400 mt-0.5">Score éco</p>
                </div>
              </div>
            </div>

            <!-- Validation note -->
            <div>
              <p class="text-xs font-semibold text-[#6b7280] uppercase tracking-wider mb-2">Note de validation</p>
              <textarea
                v-model="validationNote"
                class="w-full px-3 py-2 text-sm bg-[#f8fafc] border border-[#e5e7eb] rounded-xl focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition resize-none"
                rows="3"
                placeholder="Ajouter une note (optionnel)..."
              ></textarea>
            </div>

            <!-- Action buttons -->
            <div class="flex items-center gap-3 pt-2">
              <button class="flex-1 py-3 rounded-xl text-sm font-semibold border-2 border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition flex items-center justify-center gap-2">
                <XCircleIcon class="w-4 h-4 text-red-400" />
                Rejeter avec Commentaire
              </button>
              <button class="flex-1 py-3 rounded-xl text-sm font-semibold text-white transition hover:opacity-90 flex items-center justify-center gap-2" style="background-color:#006d35;">
                <CheckCircleIcon class="w-4 h-4" />
                Approuver et Générer Code
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import {
  InboxIcon, PhotoIcon, WrenchScrewdriverIcon,
  CheckCircleIcon, XCircleIcon,
} from '@heroicons/vue/24/outline'

const validationNote = ref('')
const selectedRequest = ref(null)

const mockRequests = ref([
  {
    id: 'REQ-001',
    object_name: 'Chaise de bureau ergonomique',
    user_name: 'Marie Dupont',
    category: 'Mobilier',
    status: 'pending',
    date: '25 mars 2026',
    weight: '8 kg',
    dimensions: '60×60×120 cm',
    description: 'Chaise de bureau en bon état, quelques rayures mineures sur les accoudoirs. Mécanisme d\'inclinaison fonctionnel.',
    co2_saved: '12 kg',
    water_saved: '340 L',
    score: '87/100',
  },
  {
    id: 'REQ-002',
    object_name: 'Lampe de bureau LED',
    user_name: 'Thomas Martin',
    category: 'Électronique',
    status: 'pending',
    date: '24 mars 2026',
    weight: '1.2 kg',
    dimensions: '15×15×40 cm',
    description: 'Lampe LED flexible en parfait état de fonctionnement. Câble USB inclus.',
    co2_saved: '2.1 kg',
    water_saved: '45 L',
    score: '94/100',
  },
  {
    id: 'REQ-003',
    object_name: 'Bibliothèque 5 étagères',
    user_name: 'Sophie Bernard',
    category: 'Mobilier',
    status: 'approved',
    date: '23 mars 2026',
    weight: '25 kg',
    dimensions: '80×30×180 cm',
    description: 'Bibliothèque en bois massif, très bonne condition. Légère usure au bas.',
    co2_saved: '45 kg',
    water_saved: '1 200 L',
    score: '78/100',
  },
  {
    id: 'REQ-004',
    object_name: 'Vélo de route aluminium',
    user_name: 'Paul Leroy',
    category: 'Sport',
    status: 'pending',
    date: '22 mars 2026',
    weight: '9 kg',
    dimensions: '160×60×90 cm',
    description: 'Vélo de route Decathlon, taille M. Pneus usés à remplacer, cadre en bon état.',
    co2_saved: '18 kg',
    water_saved: '560 L',
    score: '82/100',
  },
  {
    id: 'REQ-005',
    object_name: 'Machine à café expresso',
    user_name: 'Claire Petit',
    category: 'Électroménager',
    status: 'rejected',
    date: '21 mars 2026',
    weight: '3.5 kg',
    dimensions: '25×35×30 cm',
    description: 'Machine expresso Delonghi, problème de pompe. Ne fonctionne plus correctement.',
    co2_saved: '4.8 kg',
    water_saved: '120 L',
    score: '45/100',
  },
])

function requestStatusBadge(status) {
  if (status === 'pending')  return 'bg-[#fef9c3] text-[#854d0e]'
  if (status === 'approved') return 'bg-[#dcfce7] text-[#166534]'
  if (status === 'rejected') return 'bg-[#fee2e2] text-[#991b1b]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function requestStatusText(status) {
  const map = { pending: 'En attente', approved: 'Approuvé', rejected: 'Rejeté' }
  return map[status] || status
}
</script>
