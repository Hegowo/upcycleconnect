<template>
  <div class="flex flex-col h-full -m-4 md:-m-6">
    <div class="flex items-center justify-between px-6 py-4 bg-white border-b border-[#e5e7eb] shrink-0">
      <div>
        <h2 class="text-xl font-bold text-[#001d32]">Documentation API</h2>
        <p class="text-sm text-[#40617f] mt-0.5">OpenAPI 3.0 — UpcycleConnect</p>
      </div>
      <span class="text-xs font-semibold px-3 py-1 rounded-full bg-[#dcfce7] text-[#166534]">v1.0.0</span>
    </div>

    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-3">
        <div class="w-8 h-8 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin mx-auto"></div>
        <p class="text-sm text-gray-400">Chargement de la documentation…</p>
      </div>
    </div>

    <div v-else-if="error" class="flex-1 flex items-center justify-center">
      <div class="text-center space-y-3">
        <p class="text-red-500 font-medium">{{ error }}</p>
        <button @click="initSwagger" class="px-4 py-2 text-sm font-semibold rounded-lg text-white" style="background-color:#006d35;">
          Réessayer
        </button>
      </div>
    </div>

    <div v-show="!loading && !error" class="flex-1 overflow-y-auto swagger-wrapper">
      <div id="swagger-ui-root"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const loading = ref(true)
const error = ref(null)

const API_BASE = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/admin/v1'

const loadedAssets = []

function loadAsset(tag, attrs) {
  return new Promise((resolve, reject) => {
    const el = document.createElement(tag)
    Object.entries(attrs).forEach(([k, v]) => el.setAttribute(k, v))
    el.onload = resolve
    el.onerror = reject
    if (tag === 'link') document.head.appendChild(el)
    else document.body.appendChild(el)
    loadedAssets.push(el)
  })
}

async function initSwagger() {
  loading.value = true
  error.value = null

  try {
    await loadAsset('link', {
      rel: 'stylesheet',
      href: 'https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui.css',
    })

    await loadAsset('script', {
      src: 'https://unpkg.com/swagger-ui-dist@5.17.14/swagger-ui-bundle.js',
    })

    const specRes = await fetch(`${API_BASE}/docs/spec`, {
      headers: { Authorization: `Bearer ${auth.token}` },
    })

    if (!specRes.ok) throw new Error(`HTTP ${specRes.status}`)
    const spec = await specRes.json()

    loading.value = false

    await new Promise(r => setTimeout(r, 50))

    const ui = window.SwaggerUIBundle({
      spec,
      dom_id: '#swagger-ui-root',
      deepLinking: true,
      presets: [
        window.SwaggerUIBundle.presets.apis,
        window.SwaggerUIBundle.SwaggerUIStandalonePreset,
      ],
      layout: 'BaseLayout',
      defaultModelsExpandDepth: 1,
      defaultModelExpandDepth: 2,
      docExpansion: 'list',
      filter: true,
      requestInterceptor(req) {
        req.headers['Authorization'] = `Bearer ${auth.token}`
        return req
      },
    })

    ui.preauthorizeApiKey('bearerAuth', auth.token)
  } catch (e) {
    loading.value = false
    error.value = 'Impossible de charger la documentation. Vérifiez votre connexion.'
  }
}

onMounted(initSwagger)

onUnmounted(() => {
  loadedAssets.forEach(el => el.parentNode?.removeChild(el))
})
</script>

<style>
.swagger-wrapper .swagger-ui .topbar { display: none !important; }
.swagger-wrapper .swagger-ui { font-family: inherit; }
.swagger-wrapper .swagger-ui .info { margin: 24px 0 16px; }
.swagger-wrapper .swagger-ui .info .title { font-size: 1.5rem; color: #001d32; }
.swagger-wrapper .swagger-ui .opblock-tag { font-size: 0.95rem; color: #001d32; }
.swagger-wrapper .swagger-ui .opblock.opblock-get .opblock-summary { border-color: #006d35; }
.swagger-wrapper .swagger-ui .opblock.opblock-get { background: #f0fdf4; border-color: #006d35; }
.swagger-wrapper .swagger-ui .opblock.opblock-post { background: #eff6ff; border-color: #3b82f6; }
.swagger-wrapper .swagger-ui .opblock.opblock-put { background: #fffbeb; border-color: #f59e0b; }
.swagger-wrapper .swagger-ui .opblock.opblock-delete { background: #fff1f2; border-color: #f43f5e; }
.swagger-wrapper .swagger-ui .btn.execute { background-color: #006d35; border-color: #006d35; }
.swagger-wrapper .swagger-ui .btn.execute:hover { background-color: #005a2b; }
.swagger-wrapper .swagger-ui .auth-wrapper .authorize { background-color: #006d35; border-color: #006d35; color: #fff; }
.swagger-wrapper .swagger-ui section.models { display: none; }
</style>
