<template>
  <section v-if="campaigns.length" :class="wrapperClass">
    <div class="max-w-[1280px] mx-auto px-4 sm:px-8">
      <div class="flex items-center gap-2 mb-5">
        <SparklesIcon class="w-5 h-5 text-[#006d35]" />
        <h2 class="font-jakarta font-bold text-[#001d32] text-lg sm:text-xl">À la une</h2>
        <span class="text-[10px] font-bold uppercase tracking-wider text-[#94a3b8] bg-[#f1f5f9] px-2 py-0.5 rounded-full">Sponsorisé</span>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <RouterLink
          v-for="c in campaigns"
          :key="c.id"
          :to="`/campagnes/${c.id}`"
          class="group bg-white rounded-2xl border border-[#edf4ff] overflow-hidden hover:border-[#006d35] hover:shadow-md transition flex flex-col"
        >
          <div class="aspect-[16/9] bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] overflow-hidden relative">
            <img v-if="c.image_url" :src="c.image_url" :alt="c.title" class="w-full h-full object-cover group-hover:scale-105 transition duration-500" />
            <MegaphoneIcon v-else class="absolute inset-0 m-auto w-12 h-12 text-white/40" />
            <span class="absolute top-3 left-3 bg-[#006d35] text-white text-[10px] font-bold uppercase tracking-wider px-2 py-1 rounded-full">
              Mis en avant
            </span>
          </div>
          <div class="p-5 flex flex-col gap-1.5 flex-1">
            <h3 class="font-jakarta font-bold text-[#001d32] text-base leading-snug line-clamp-2">{{ c.title }}</h3>
            <p v-if="c.description" class="text-[#40617f] text-xs line-clamp-3 leading-relaxed">{{ c.description }}</p>
            <p class="text-[11px] text-[#94a3b8] mt-auto pt-3 border-t border-[#f1f5f9] flex items-center gap-1.5">
              <UserIcon class="w-3.5 h-3.5" />
              Proposé par {{ c.provider_name || 'un artisan UpcycleConnect' }}
            </p>
          </div>
        </RouterLink>
      </div>
    </div>
  </section>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { SparklesIcon, MegaphoneIcon, UserIcon } from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

defineProps({
  wrapperClass: { type: String, default: 'py-10' },
})

const campaigns = ref([])

onMounted(async () => {
  try {
    const data = await publicGet('/campaigns/active')
    campaigns.value = data.data || []
  } catch {
    campaigns.value = []
  }
})
</script>
