<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-[1100px] mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/communaute" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" />
        Retour à la communauté
      </RouterLink>

      <div class="mb-8">
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-3xl tracking-tight">
          Conseils & astuces
        </h1>
        <p class="text-[#40617f] text-sm mt-1">
          Les articles publiés par l'équipe UpcycleConnect pour vous accompagner dans votre démarche d'upcycling.
        </p>
      </div>

      <div class="flex flex-col sm:flex-row gap-3 mb-6">
        <div class="relative flex-1">
          <MagnifyingGlassIcon class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#94a3b8]" />
          <input
            v-model="search"
            @input="debouncedFetch"
            type="search"
            placeholder="Rechercher un conseil..."
            class="w-full pl-10 pr-4 py-2.5 bg-white border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 focus:border-[#006d35] transition"
          />
        </div>
      </div>

      <div v-if="categories.length" class="flex flex-wrap items-center gap-2 mb-8">
        <button
          @click="setCategory('')"
          :class="['px-3 py-1.5 rounded-full text-xs font-semibold border transition',
                   category === '' ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]']"
        >
          Tout
        </button>
        <button
          v-for="cat in categories"
          :key="cat"
          @click="setCategory(cat)"
          :class="['px-3 py-1.5 rounded-full text-xs font-semibold border transition',
                   category === cat ? 'bg-[#006d35] text-white border-[#006d35]' : 'bg-white text-[#40617f] border-[#e5e7eb] hover:border-[#006d35]']"
        >
          {{ cat }}
        </button>
      </div>

      <div v-if="loading" class="py-20 text-center">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>

      <div v-else-if="!tips.length" class="bg-white rounded-2xl p-12 text-center border border-[#edf4ff]">
        <LightBulbIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <p class="text-[#001d32] font-semibold">Aucun conseil pour l'instant</p>
        <p class="text-[#40617f] text-sm mt-1">L'équipe publiera bientôt de nouveaux articles ici.</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <RouterLink
          v-for="tip in tips"
          :key="tip.id"
          :to="`/communaute/conseils/${tip.slug}`"
          class="group bg-white rounded-2xl overflow-hidden border border-[#edf4ff] hover:border-[#006d35] hover:shadow-md transition flex flex-col"
        >
          <div class="aspect-[16/9] bg-[#edf4ff] overflow-hidden">
            <img
              v-if="tip.cover_image"
              :src="tip.cover_image"
              :alt="tip.title"
              class="w-full h-full object-cover group-hover:scale-105 transition duration-500"
            />
            <div v-else class="w-full h-full flex items-center justify-center">
              <LightBulbIcon class="w-12 h-12 text-[#006d35]/30" />
            </div>
          </div>
          <div class="p-5 flex flex-col gap-2 flex-1">
            <span v-if="tip.category" class="text-[10px] font-bold uppercase tracking-wider text-[#006d35]">{{ tip.category }}</span>
            <h3 class="font-jakarta font-bold text-[#001d32] text-base leading-snug group-hover:text-[#006d35] transition line-clamp-2">{{ tip.title }}</h3>
            <p v-if="tip.summary" class="text-[#40617f] text-xs line-clamp-3 leading-relaxed">{{ tip.summary }}</p>
            <div class="flex items-center justify-between text-[11px] text-[#94a3b8] mt-auto pt-3 border-t border-[#f1f5f9]">
              <span class="flex items-center gap-1.5">
                <UserIcon class="w-3.5 h-3.5" />
                {{ tip.author_name || 'UpcycleConnect' }}
              </span>
              <span>{{ formatDate(tip.published_at || tip.created_at) }}</span>
            </div>
          </div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import {
  ArrowLeftIcon,
  MagnifyingGlassIcon,
  LightBulbIcon,
  UserIcon,
} from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

const tips = ref([])
const categories = ref([])
const category = ref('')
const search = ref('')
const loading = ref(true)
let searchTimer = null

async function fetchTips() {
  loading.value = true
  try {
    const data = await publicGet('/tips', {
      category: category.value,
      search: search.value,
    })
    tips.value = data.data || []
    categories.value = data.categories || []
  } catch {
    tips.value = []
  } finally {
    loading.value = false
  }
}

function debouncedFetch() {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(fetchTips, 350)
}

function setCategory(c) {
  category.value = c
  fetchTips()
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(fetchTips)
</script>
