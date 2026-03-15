<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold text-gray-900">{{ t('logs.title') }}</h2>
        <p class="text-sm text-gray-500 mt-0.5">{{ tab === 'admin' ? t('logs.adminSubtitle') : t('logs.platformSubtitle') }}</p>
      </div>
    </div>

    <div class="flex gap-1 p-1 bg-gray-100 rounded-xl w-fit">
      <button @click="tab='admin'" :class="tab==='admin' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm font-medium transition">
        🔐 {{ t('logs.tabAdmin') }}
      </button>
      <button @click="tab='platform'" :class="tab==='platform' ? 'bg-white shadow text-gray-900' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm font-medium transition">
        👥 {{ t('logs.tabPlatform') }}
      </button>
    </div>

    <div v-if="tab === 'admin'" class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100">
        <span class="text-sm font-medium text-gray-700">{{ t('logs.adminSubtitle') }}</span>
      </div>
      <div v-if="adminLoading" class="p-6 space-y-3">
        <div v-for="n in 8" :key="n" class="h-10 bg-gray-100 rounded-lg animate-pulse"></div>
      </div>
      <div v-else-if="!adminLogs.length" class="py-16 text-center text-gray-400">
        <div class="text-4xl mb-2">📋</div>
        <p>{{ t('logs.adminEmpty') }}</p>
      </div>
      <div v-else class="divide-y divide-gray-50">
        <div v-for="log in adminLogs" :key="log.id" class="flex items-start gap-4 px-4 py-3 hover:bg-gray-50 transition-colors">
          <div class="w-9 h-9 rounded-full flex items-center justify-center shrink-0 text-lg" :class="actionBgClass(log.action)">
            {{ actionIcon(log.action) }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800">
              <span class="font-semibold">{{ log.admin_name?.trim() || t('logs.system') }}</span>
              — {{ actionLabel(log.action) }}
              <span v-if="log.resource_id" class="text-gray-400 text-xs"> ({{ t('logs.resourceId', { id: log.resource_id }) }})</span>
            </p>
            <p class="text-xs text-gray-400 mt-0.5">{{ log.ip_address }} · {{ formatDate(log.created_at) }}</p>
          </div>
          <span class="shrink-0 text-xs px-2 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
            {{ log.resource_type }}
          </span>
        </div>
      </div>
    </div>

    <div v-if="tab === 'platform'" class="card overflow-hidden">
      <div class="px-4 py-3 border-b border-gray-100">
        <span class="text-sm font-medium text-gray-700">{{ t('logs.platformSubtitle') }}</span>
      </div>
      <div v-if="platformLoading" class="p-6 space-y-3">
        <div v-for="n in 8" :key="n" class="h-10 bg-gray-100 rounded-lg animate-pulse"></div>
      </div>
      <div v-else-if="!platformLogs.length" class="py-16 text-center text-gray-400">
        <div class="text-4xl mb-2">👥</div>
        <p>{{ t('logs.platformEmpty') }}</p>
      </div>
      <div v-else class="divide-y divide-gray-50">
        <div v-for="entry in platformLogs" :key="entry.id + entry.type" class="flex items-start gap-4 px-4 py-3 hover:bg-gray-50 transition-colors">
          <div class="w-9 h-9 rounded-full flex items-center justify-center shrink-0 text-lg" :class="entry.type === 'provider' ? 'bg-purple-50' : 'bg-blue-50'">
            {{ entry.type === 'provider' ? '🔨' : '👤' }}
          </div>
          <div class="flex-1 min-w-0">
            <p class="text-sm text-gray-800">
              <span class="font-semibold">{{ entry.name }}</span>
              — {{ entry.type === 'provider' ? t('logs.registrationProvider') : t('logs.registrationUser') }}
              <span v-if="entry.company" class="text-gray-500"> ({{ entry.company }})</span>
            </p>
            <p class="text-xs text-gray-400 mt-0.5">{{ entry.email }} · {{ formatDate(entry.created_at) }}</p>
          </div>
          <span class="shrink-0 text-xs px-2 py-0.5 rounded-full font-medium" :class="entry.type === 'provider' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'">
            {{ entry.type === 'provider' ? t('logs.typeProvider') : t('logs.typeUser') }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/services/api'

const { t } = useI18n()
const tab = ref('admin')
const adminLogs = ref([])
const platformLogs = ref([])
const adminLoading = ref(false)
const platformLoading = ref(false)

async function loadAdminLogs() {
  adminLoading.value = true
  try {
    const { data } = await api.get('/logs')
    adminLogs.value = data.data || []
  } finally {
    adminLoading.value = false
  }
}

async function loadPlatformLogs() {
  platformLoading.value = true
  try {
    const { data } = await api.get('/logs/activity')
    platformLogs.value = data.data || []
  } finally {
    platformLoading.value = false
  }
}

onMounted(() => {
  loadAdminLogs()
  loadPlatformLogs()
})

function formatDate(dateStr) {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now - d) / 1000)
  if (diff < 60) return t('dashboard.actionInstant')
  if (diff < 3600) return t('dashboard.actionMinutesAgo', { n: Math.floor(diff / 60) })
  if (diff < 86400) return t('dashboard.actionHoursAgo', { n: Math.floor(diff / 3600) })
  return d.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function actionIcon(action) {
  if (action.includes('created')) return '✅'
  if (action.includes('deleted')) return '🗑️'
  if (action.includes('banned')) return '🚫'
  if (action.includes('approved')) return '✔️'
  if (action.includes('suspended')) return '⏸️'
  if (action.includes('activated')) return '🔓'
  return '✏️'
}
function actionBgClass(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'bg-red-50'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'bg-green-50'
  return 'bg-gray-50'
}
function actionLabel(action) {
  const labels = {
    'user.banned':        t('logs.actionUserBanned'),
    'user.activated':     t('logs.actionUserActivated'),
    'user.deleted':       t('logs.actionUserDeleted'),
    'provider.approved':  t('logs.actionProviderApproved'),
    'provider.suspended': t('logs.actionProviderSuspended'),
    'admin.created':      t('logs.actionAdminCreated'),
    'admin.deleted':      t('logs.actionAdminDeleted'),
    'category.created':   t('logs.actionCategoryCreated'),
    'category.deleted':   t('logs.actionCategoryDeleted'),
    'category.toggled':   t('logs.actionCategoryToggled'),
  }
  return labels[action] || action
}
function actionBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'bg-red-100 text-red-700'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'bg-green-100 text-green-700'
  if (action.includes('suspended')) return 'bg-orange-100 text-orange-700'
  return 'bg-gray-100 text-gray-600'
}
</script>
