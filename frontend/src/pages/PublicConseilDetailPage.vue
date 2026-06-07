<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 py-10">

      <RouterLink to="/communaute/conseils" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#006d35] mb-6">
        <ArrowLeftIcon class="w-4 h-4" />
        Tous les conseils
      </RouterLink>

      <div v-if="loading" class="py-20 text-center">
        <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto" />
      </div>

      <div v-else-if="!tip" class="bg-white rounded-2xl p-12 text-center border border-[#edf4ff]">
        <ExclamationCircleIcon class="w-12 h-12 text-gray-300 mx-auto mb-3" />
        <p class="text-[#001d32] font-semibold">Conseil introuvable</p>
        <RouterLink to="/communaute/conseils" class="inline-block mt-4 px-4 py-2 rounded-xl bg-[#006d35] text-white text-sm font-semibold">
          Retour aux conseils
        </RouterLink>
      </div>

      <article v-else class="bg-white rounded-[24px] overflow-hidden border border-[#edf4ff]">
        <div v-if="tip.cover_image" class="aspect-[16/9] bg-[#edf4ff]">
          <img :src="tip.cover_image" :alt="tip.title" class="w-full h-full object-cover" />
        </div>

        <div class="p-6 sm:p-10">
          <span v-if="tip.category" class="inline-block text-[10px] font-bold uppercase tracking-widest text-[#006d35] mb-3">{{ tip.category }}</span>

          <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl sm:text-3xl tracking-tight leading-tight mb-4">
            {{ tip.title }}
          </h1>

          <div class="flex items-center gap-4 text-xs text-[#94a3b8] pb-5 mb-6 border-b border-[#f1f5f9]">
            <span class="flex items-center gap-1.5">
              <UserIcon class="w-3.5 h-3.5" />
              {{ tip.author_name || 'UpcycleConnect' }}
            </span>
            <span>·</span>
            <span>{{ formatDate(tip.published_at || tip.created_at) }}</span>
            <span>·</span>
            <span class="flex items-center gap-1.5">
              <EyeIcon class="w-3.5 h-3.5" />
              {{ tip.view_count || 0 }} vues
            </span>
          </div>

          <p v-if="tip.summary" class="text-[#001d32] text-base leading-relaxed font-medium bg-[#edf4ff] rounded-xl px-5 py-4 mb-6">
            {{ tip.summary }}
          </p>

          <div class="tip-content text-[#001d32] text-base leading-relaxed" v-html="tip.content" />
        </div>
      </article>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  ArrowLeftIcon,
  UserIcon,
  EyeIcon,
  ExclamationCircleIcon,
} from '@heroicons/vue/24/outline'
import { publicGet } from '@/services/publicApi'

const route = useRoute()
const tip = ref(null)
const loading = ref(true)

async function fetchTip(slug) {
  loading.value = true
  try {
    tip.value = await publicGet(`/tips/${slug}`)
  } catch {
    tip.value = null
  } finally {
    loading.value = false
  }
}

function formatDate(iso) {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString('fr-FR', { day: 'numeric', month: 'long', year: 'numeric' })
}

onMounted(() => fetchTip(route.params.slug))
watch(() => route.params.slug, (s) => s && fetchTip(s))
</script>

<style scoped>
.tip-content :deep(h1),
.tip-content :deep(h2),
.tip-content :deep(h3) {
  font-family: 'Plus Jakarta Sans', sans-serif;
  font-weight: 800;
  color: #001d32;
  margin-top: 1.6em;
  margin-bottom: 0.6em;
  line-height: 1.25;
}
.tip-content :deep(h1) { font-size: 1.5rem; }
.tip-content :deep(h2) { font-size: 1.3rem; }
.tip-content :deep(h3) { font-size: 1.1rem; }
.tip-content :deep(p) { margin: 0.9em 0; }
.tip-content :deep(ul),
.tip-content :deep(ol) { margin: 0.9em 0; padding-left: 1.4em; }
.tip-content :deep(li) { margin: 0.3em 0; }
.tip-content :deep(a) { color: #006d35; text-decoration: underline; }
.tip-content :deep(blockquote) {
  border-left: 4px solid #006d35;
  padding-left: 1em;
  color: #40617f;
  margin: 1.2em 0;
  font-style: italic;
}
.tip-content :deep(img) { max-width: 100%; border-radius: 12px; margin: 1.2em 0; }
.tip-content :deep(strong) { font-weight: 700; color: #001d32; }
.tip-content :deep(code) { background: #edf4ff; padding: 0.15em 0.4em; border-radius: 4px; font-size: 0.9em; }
</style>
