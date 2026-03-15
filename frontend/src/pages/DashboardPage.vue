<template>
  <div class="space-y-6">

    <div class="rounded-2xl p-6 text-white flex items-center justify-between" style="background: linear-gradient(135deg, #103652, #1B5B88);">
      <div>
        <p class="text-blue-200 text-sm mb-1">{{ greeting }}</p>
        <h2 class="text-2xl font-bold">{{ auth.fullName || 'Administrateur' }}</h2>
        <p class="text-blue-200 text-sm mt-1">{{ todayFormatted }}</p>
      </div>
      <div class="text-5xl opacity-30 hidden md:block">🌿</div>
    </div>

    <div v-if="(stats?.providers_pending ?? 0) > 0"
      class="flex items-center justify-between rounded-xl p-4 border"
      style="background:#FFF7ED; border-color:#fed7aa;"
    >
      <div class="flex items-center gap-3">
        <span class="text-2xl">⚠️</span>
        <div>
          <p class="text-sm font-semibold text-orange-800">
            {{ t('dashboard.pendingAlert', { n: stats.providers_pending }) }}
          </p>
          <p class="text-xs text-orange-600">{{ t('dashboard.pendingAlertSub') }}</p>
        </div>
      </div>
      <RouterLink to="/admin/providers?status=pending"
        class="text-sm font-semibold px-4 py-2 rounded-lg text-white shrink-0"
        style="background-color: #c2410c;"
      >
        {{ t('dashboard.pendingAlertBtn') }}
      </RouterLink>
    </div>

    <div v-if="loading" class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
      <div v-for="n in 5" :key="n" class="rounded-2xl p-5 bg-gray-100 animate-pulse h-28"></div>
    </div>
    <div v-else class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
      <StatCard :label="t('dashboard.statsUsers')"         :value="stats?.users_count ?? 0"       icon="👥" color="blue"   link="/admin/users" />
      <StatCard :label="t('dashboard.statsProviders')" :value="stats?.providers_count ?? 0"    icon="🔨" color="green"  link="/admin/providers?status=approved" />
      <StatCard :label="t('dashboard.statsPending')"           :value="stats?.providers_pending ?? 0"  icon="⏳" color="orange" link="/admin/providers?status=pending" :alert="(stats?.providers_pending ?? 0) > 0" />
      <StatCard :label="t('dashboard.statsPrestations')"          :value="stats?.prestations_count ?? 0"  icon="📋" color="purple" link="/admin/prestations" />
      <StatCard :label="t('dashboard.statsEvents')"    :value="stats?.events_count ?? 0"       icon="📅" color="teal"   link="/admin/events" />
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-3 gap-6">

      <div class="xl:col-span-2 card p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-semibold text-gray-900">{{ t('dashboard.activityTitle') }}</h3>
          <span class="text-xs text-gray-400">{{ t('dashboard.activitySubtitle') }}</span>
        </div>

        <div v-if="logsLoading" class="space-y-3">
          <div v-for="n in 5" :key="n" class="h-10 bg-gray-100 rounded-lg animate-pulse"></div>
        </div>

        <div v-else-if="logs.length === 0" class="text-center py-8 text-gray-400 text-sm">
          {{ t('dashboard.activityEmpty') }}
        </div>

        <div v-else class="space-y-1 max-h-96 overflow-y-auto pr-1">
          <div
            v-for="log in logs"
            :key="log.id"
            class="flex items-start gap-3 px-3 py-2.5 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <span class="text-lg shrink-0 mt-0.5">{{ actionIcon(log.action) }}</span>
            <div class="flex-1 min-w-0">
              <p class="text-sm text-gray-800">
                <span class="font-medium">{{ log.admin_name?.trim() || t('logs.system') }}</span>
                — {{ actionLabel(log.action) }}
                <span v-if="log.resource_id" class="text-gray-400"> #{{ log.resource_id }}</span>
              </p>
              <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(log.created_at) }}</p>
            </div>
            <span class="shrink-0 text-xs px-2 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
              {{ log.resource_type }}
            </span>
          </div>
        </div>
      </div>

      <div class="space-y-6">

        <div class="card p-6">
          <h3 class="font-semibold text-gray-900 mb-4">{{ t('dashboard.prestationsTitle') }}</h3>
          <div v-if="!stats || !Object.keys(stats.prestations_by_status || {}).length" class="text-sm text-gray-400">
            {{ t('common.noResults') }}
          </div>
          <div v-else class="space-y-3">
            <div v-for="(count, status) in stats.prestations_by_status" :key="status">
              <div class="flex justify-between items-center mb-1">
                <span class="text-xs font-medium capitalize text-gray-600">{{ status }}</span>
                <span class="text-xs font-bold text-gray-800">{{ count }}</span>
              </div>
              <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :style="{ width: barWidth(count) + '%', backgroundColor: statusColor(status) }"
                ></div>
              </div>
            </div>
          </div>
        </div>

        <div class="card p-6">
          <h3 class="font-semibold text-gray-900 mb-4">{{ t('dashboard.quickActionsTitle') }}</h3>
          <div class="grid grid-cols-2 gap-2">
            <RouterLink
              v-for="action in quickActions"
              :key="action.to"
              :to="action.to"
              class="flex flex-col items-center gap-1.5 p-3 rounded-xl text-center transition hover:opacity-80"
              :style="{ backgroundColor: action.bg }"
            >
              <span class="text-xl">{{ action.icon }}</span>
              <span class="text-xs font-medium leading-tight" :style="{ color: action.color }">{{ action.label }}</span>
            </RouterLink>
          </div>
        </div>

      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import StatCard from '@/components/StatCard.vue'

const { t, locale } = useI18n()
const auth = useAuthStore()
const stats = ref(null)
const logs  = ref([])
const loading     = ref(true)
const logsLoading = ref(true)

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 12) return t('dashboard.greetingMorning')
  if (h < 18) return t('dashboard.greetingAfternoon')
  return t('dashboard.greetingEvening')
})

const dateLocale = computed(() => locale.value === 'en' ? 'en-GB' : 'fr-FR')

const todayFormatted = computed(() => {
  return new Date().toLocaleDateString(dateLocale.value, { weekday: 'long', day: 'numeric', month: 'long', year: 'numeric' })
})

onMounted(async () => {
  try {
    const { data } = await api.get('/dashboard/stats')
    stats.value = data
  } finally {
    loading.value = false
  }

  try {
    const { data } = await api.get('/logs')
    logs.value = data.data || []
  } finally {
    logsLoading.value = false
  }
})

const totalPrestations = computed(() => {
  if (!stats.value) return 1
  return Object.values(stats.value.prestations_by_status || {}).reduce((a, b) => a + b, 0) || 1
})
function barWidth(count) {
  return Math.round((count / totalPrestations.value) * 100)
}
function statusColor(status) {
  return { published: '#1B8848', draft: '#94a3b8', suspended: '#f97316', archived: '#e11d48' }[status] || '#94a3b8'
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
function actionLabel(action) {
  const labels = {
    'user.banned':       t('logs.actionUserBanned'),
    'user.activated':    t('logs.actionUserActivated'),
    'user.deleted':      t('logs.actionUserDeleted'),
    'provider.approved': t('logs.actionProviderApproved'),
    'provider.suspended':t('logs.actionProviderSuspended'),
    'admin.created':     t('logs.actionAdminCreated'),
    'admin.deleted':     t('logs.actionAdminDeleted'),
    'category.created':  t('logs.actionCategoryCreated'),
    'category.deleted':  t('logs.actionCategoryDeleted'),
    'category.toggled':  t('logs.actionCategoryToggled'),
  }
  return labels[action] || action
}
function actionBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'bg-red-100 text-red-700'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'bg-green-100 text-green-700'
  if (action.includes('suspended')) return 'bg-orange-100 text-orange-700'
  return 'bg-gray-100 text-gray-600'
}
function formatDate(dateStr) {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now - d) / 1000)
  if (diff < 60) return t('dashboard.actionInstant')
  if (diff < 3600) return t('dashboard.actionMinutesAgo', { n: Math.floor(diff / 60) })
  if (diff < 86400) return t('dashboard.actionHoursAgo', { n: Math.floor(diff / 3600) })
  return d.toLocaleDateString(dateLocale.value, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}

const quickActions = computed(() => [
  { to: '/admin/categories/create',  icon: '🏷️',  label: t('dashboard.quickCategory'),  bg: '#EDE9FE', color: '#6d28d9' },
  { to: '/admin/events/create',      icon: '📅',  label: t('dashboard.quickEvent'),    bg: '#CCFBF1', color: '#0f766e' },
  { to: '/admin/prestations',        icon: '📋',  label: t('dashboard.quickPrestations'),    bg: '#DCFCE7', color: '#15803d' },
  { to: '/admin/providers',          icon: '🔨',  label: t('dashboard.quickProviders'),        bg: '#DBEAFE', color: '#1e40af' },
])
</script>
