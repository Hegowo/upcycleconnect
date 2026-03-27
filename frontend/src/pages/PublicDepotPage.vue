<template>
  <div class="bg-[#f7f9ff] pb-0">
    <div class="max-w-[1280px] mx-auto px-4 sm:px-8 pt-8 sm:pt-12 pb-12 sm:pb-16">

      <!-- Hero -->
      <div class="mb-12">
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl sm:text-5xl tracking-tight">Déposer un objet</h1>
        <p class="text-[#40617f] text-lg mt-4 max-w-2xl leading-relaxed">
          Donnez une seconde vie à vos objets. Remplissez les détails ci-dessous pour initier votre demande de dépôt dans notre écosystème circulaire.
        </p>
      </div>

      <!-- Bento Grid Form -->
      <div class="grid grid-cols-12 gap-4 sm:gap-6">

        <!-- Photos upload -->
        <div class="col-span-12 lg:col-span-8 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <div class="flex items-center justify-between">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">Photos de l'objet</h2>
            <span class="text-[#40617f] text-xs uppercase tracking-wide">Min. 3 slots</span>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-3 gap-4">
            <div
              v-for="slot in photoSlots"
              :key="slot.id"
              class="bg-[#edf4ff] border-2 border-dashed border-[#becabc] rounded-xl flex flex-col items-center justify-center cursor-pointer hover:bg-[#d8eaff] transition h-40 sm:h-60 group"
              @click="triggerUpload(slot.id)"
            >
              <div v-if="!slot.file">
                <PhotoIcon class="w-10 h-10 text-[#40617f]/40 mx-auto mb-2 group-hover:text-[#40617f]/60" />
                <p class="text-[#40617f] text-sm text-center font-medium">{{ slot.label }}</p>
              </div>
              <div v-else class="text-center">
                <CheckCircleIcon class="w-10 h-10 text-[#006d35] mx-auto mb-2" />
                <p class="text-[#006d35] text-sm font-semibold">Photo ajoutée</p>
              </div>
            </div>
          </div>

          <p class="text-[#40617f] text-sm">
            Conseil : Des photos nettes et bien éclairées accélèrent le processus de validation.
          </p>
        </div>

        <!-- Description -->
        <div class="col-span-12 lg:col-span-4 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">Description</h2>

          <div class="flex flex-col gap-5">
            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">Nom de l'objet</label>
              <input
                v-model="form.name"
                type="text"
                placeholder="Ex: Chaise en bois scandinave"
                class="bg-[#edf4ff] px-3 py-3.5 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280]"
              />
            </div>

            <div class="flex flex-col gap-2">
              <label class="text-[#001d32] font-semibold text-sm">Détails</label>
              <textarea
                v-model="form.details"
                rows="6"
                placeholder="Décrivez les dimensions, les matériaux ou l'histoire de l'objet..."
                class="bg-[#edf4ff] px-3 py-3 rounded-xl text-base text-[#001d32] outline-none focus:ring-2 focus:ring-[#006d35]/20 placeholder-[#6b7280] resize-none"
              />
            </div>
          </div>
        </div>

        <!-- Classification -->
        <div class="col-span-12 lg:col-span-5 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">Classification</h2>

          <div class="flex flex-col gap-3">
            <label class="text-[#001d32] font-semibold text-sm">Catégorie</label>
            <div class="grid grid-cols-2 gap-3">
              <button
                v-for="cat in categories"
                :key="cat"
                @click="form.category = cat"
                class="px-5 py-3 rounded-xl text-sm font-medium border-2 transition-all"
                :class="form.category === cat
                  ? 'bg-[#006d35] text-white border-[#006d35]'
                  : 'bg-[#edf4ff] text-[#001d32] border-transparent hover:border-[#006d35]/20'"
              >
                {{ cat }}
              </button>
            </div>
          </div>
        </div>

        <!-- Condition + Submit -->
        <div class="col-span-12 lg:col-span-7 bg-white rounded-[24px] p-6 sm:p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] flex flex-col gap-6">
          <h2 class="font-jakarta font-bold text-[#001d32] text-xl">État de l'objet</h2>

          <div class="flex flex-wrap gap-3">
            <button
              v-for="cond in conditions"
              :key="cond"
              @click="form.condition = cond"
              class="px-5 py-2.5 rounded-full text-sm font-medium transition-all border-2"
              :class="form.condition === cond
                ? 'bg-[#006d35] text-white border-[#006d35]'
                : 'bg-[#edf4ff] text-[#001d32] border-transparent hover:border-[#006d35]/20'"
            >
              {{ cond }}
            </button>
          </div>

          <div class="bg-[rgba(27,136,72,0.1)] rounded-xl p-4 sm:p-6 flex flex-col sm:flex-row items-start sm:items-center gap-4 sm:justify-between">
            <div class="flex items-center gap-4">
              <div class="w-6 h-6 bg-[#1b8848] rounded-full shrink-0" />
              <div>
                <p class="text-[#006d35] font-semibold text-base">Impact estimé</p>
                <p class="text-[#40617f] text-sm">Réduction de ~4.5kg CO2</p>
              </div>
            </div>
            <button
              class="w-full sm:w-auto px-8 sm:px-10 py-4 rounded-xl text-white font-semibold text-base sm:text-lg transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,0,0,0.1)]"
              style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
              @click="submitForm"
            >
              Soumettre la demande
            </button>
          </div>
        </div>

        <!-- How it works (full width) -->
        <div class="col-span-12 rounded-[24px] overflow-hidden relative min-h-[160px] sm:h-48">
          <div class="absolute inset-0 bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] opacity-50" />
          <div class="absolute inset-0 bg-gradient-to-r from-[#f7f9ff] to-transparent" />
          <div class="relative p-6 sm:p-12 flex flex-col justify-center h-full">
            <h3 class="font-jakarta font-bold text-[#001d32] text-2xl mb-2">Comment ça marche ?</h3>
            <p class="text-[#40617f] text-base leading-relaxed max-w-lg">
              Une fois soumis, nos experts valideront votre objet sous 24h. Vous recevrez ensuite un bon de dépôt pour le point de collecte le plus proche.
            </p>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<script setup>
import { ref } from 'vue'
import {
  PhotoIcon,
  CheckCircleIcon,
} from '@heroicons/vue/24/outline'

const form = ref({
  name: '',
  details: '',
  category: '',
  condition: '',
})

const photoSlots = ref([
  { id: 1, label: 'Photo principale', file: null },
  { id: 2, label: 'Détail', file: null },
  { id: 3, label: 'Angle différent', file: null },
])

const categories = ['Mobilier', 'Textile', 'Électro', 'Autre']

const conditions = ['Neuf', 'Très bon état', 'Bon état', 'Usé', 'Cassé/Pour pièces']

const footerLinks = ['Privacy Policy', 'Terms of Service', 'Sustainability Report', 'Contact']

const triggerUpload = (id) => {
  // Demo: toggle file state
  const slot = photoSlots.value.find(s => s.id === id)
  if (slot) slot.file = slot.file ? null : 'placeholder'
}

const submitForm = () => {
  alert('Demande soumise ! Nos experts vous répondront sous 24h.')
}
</script>
