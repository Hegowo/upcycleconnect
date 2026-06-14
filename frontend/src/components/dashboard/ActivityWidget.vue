<template>
  <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden h-full">
    <div class="flex items-center justify-between px-4 sm:px-6 py-3 sm:py-4 border-b border-[#f1f5f9]">
      <h3 class="font-semibold text-[#001d32] text-sm sm:text-base">{{ t('dashboard.recentActivity') }}</h3>
      <RouterLink to="/admin/logs" class="text-xs font-medium hover:underline" style="color:#006d35;">{{ t('common.seeAll') }}</RouterLink>
    </div>

    <div v-if="loading" class="p-4 space-y-3">
      <div v-for="n in 5" :key="n" class="h-10 bg-gray-50 rounded-lg animate-pulse"></div>
    </div>
    <div v-else-if="logs.length === 0" class="text-center py-12 text-gray-400 text-sm">
      {{ t('dashboard.activityEmpty') }}
    </div>
    <div v-else class="divide-y divide-[#f8fafc]">
      <div
        v-for="log in logs.slice(0, 8)"
        :key="log.id"
        class="flex items-start gap-3 px-4 sm:px-6 py-3 sm:py-3.5 hover:bg-[#f8fafc] transition-colors cursor-pointer"
        @click="openModal(log)"
      >
        <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full flex items-center justify-center shrink-0 text-xs font-bold text-white mt-0.5"
          :style="{ backgroundColor: avatarColor(log.admin_name) }">
          {{ (log.admin_name?.trim() || 'S').charAt(0).toUpperCase() }}
        </div>

        <div class="flex-1 min-w-0">
          <div class="hidden sm:flex items-center gap-3">
            <p class="text-sm text-gray-800 flex-1 min-w-0 truncate">
              <span class="font-medium">{{ log.admin_name?.trim() || 'Système' }}</span>
              — {{ actionLabel(log.action) }}
              <span v-if="log.resource_id" class="text-gray-400"> #{{ log.resource_id }}</span>
            </p>
            <div class="flex items-center gap-2 shrink-0">
              <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
                {{ log.resource_type || '—' }}
              </span>
              <span class="text-xs px-2 py-0.5 rounded-full font-medium" :class="statusBadgeClass(log.action)">
                {{ statusLabel(log.action) }}
              </span>
            </div>
          </div>
          <div class="sm:hidden">
            <p class="text-xs text-gray-800 leading-snug">
              <span class="font-semibold">{{ log.admin_name?.trim() || 'Système' }}</span>
              — {{ actionLabel(log.action) }}
            </p>
            <div class="flex items-center gap-1.5 mt-1 flex-wrap">
              <span class="text-[10px] text-gray-400">{{ formatDate(log.created_at) }}</span>
              <span class="text-[10px] px-1.5 py-0.5 rounded-full font-medium" :class="actionBadgeClass(log.action)">
                {{ log.resource_type || '—' }}
              </span>
              <span class="text-[10px] px-1.5 py-0.5 rounded-full font-medium" :class="statusBadgeClass(log.action)">
                {{ statusLabel(log.action) }}
              </span>
            </div>
          </div>
          <p class="hidden sm:block text-xs text-gray-400 mt-0.5">{{ formatDate(log.created_at) }}</p>
        </div>
      </div>
    </div>
  </div>

  <Teleport to="body">
    <Transition name="modal">
      <div v-if="selectedLog" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center sm:p-4" @click.self="closeModal">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
        <div class="relative w-full sm:max-w-2xl rounded-t-2xl sm:rounded-xl overflow-hidden shadow-2xl border border-gray-700">
          <div class="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-gray-700">
            <button @click="closeModal" class="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors"></button>
            <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
            <div class="w-3 h-3 rounded-full bg-green-500"></div>
            <span class="ml-3 text-xs font-mono text-gray-400 truncate">upcycleconnect — log#{{ selectedLog.id }}</span>
            <button @click="closeModal" class="ml-auto text-gray-500 hover:text-gray-300 transition-colors shrink-0">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
          </div>
          <div class="bg-gray-900 p-4 sm:p-5 font-mono text-xs sm:text-sm leading-relaxed max-h-[70vh] overflow-y-auto">
            <div class="flex items-center gap-2 mb-4 pb-3 border-b border-gray-700 flex-wrap">
              <span class="text-green-400">root@upcycleconnect</span>
              <span class="text-gray-500">:</span>
              <span class="text-blue-400">~/logs</span>
              <span class="text-gray-500">$</span>
              <span class="text-white">cat log.{{ selectedLog.id }}.json</span>
            </div>
            <div class="space-y-1.5">
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">id</span><span class="text-yellow-300">{{ selectedLog.id }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">timestamp</span><span class="text-cyan-300 break-all">{{ formatDateFull(selectedLog.created_at) }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">ip_address</span><span class="text-green-300">{{ selectedLog.ip_address || 'unknown' }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">admin</span><span class="text-white">{{ selectedLog.admin_name?.trim() || 'system' }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">admin_id</span><span class="text-yellow-300">{{ selectedLog.admin_id ?? '—' }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">action</span><span :class="actionTerminalColor(selectedLog.action)" class="font-bold break-all">{{ selectedLog.action }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">resource_type</span><span class="text-purple-300">{{ selectedLog.resource_type || '—' }}</span></div>
              <div class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">resource_id</span><span class="text-yellow-300">{{ selectedLog.resource_id ?? '—' }}</span></div>
              <template v-if="selectedLog.old_values">
                <div class="flex gap-2 sm:gap-3 mt-2 pt-2 border-t border-gray-700">
                  <span class="text-gray-500 w-24 sm:w-28 shrink-0">old_values</span>
                  <span class="text-red-300 break-all text-xs font-mono whitespace-pre">{{ formatJson(selectedLog.old_values) }}</span>
                </div>
              </template>
              <template v-if="selectedLog.new_values">
                <div class="flex gap-2 sm:gap-3">
                  <span class="text-gray-500 w-24 sm:w-28 shrink-0">new_values</span>
                  <span class="text-green-300 break-all text-xs font-mono whitespace-pre">{{ formatJson(selectedLog.new_values) }}</span>
                </div>
              </template>
              <div v-if="selectedLog.details" class="flex gap-2 sm:gap-3"><span class="text-gray-500 w-24 sm:w-28 shrink-0">details</span><span class="text-gray-300 break-all">{{ selectedLog.details }}</span></div>
            </div>
            <div class="mt-4 pt-3 border-t border-gray-700 flex items-center gap-1 flex-wrap">
              <span class="text-green-400">root@upcycleconnect</span>
              <span class="text-gray-500">:</span>
              <span class="text-blue-400">~/logs</span>
              <span class="text-gray-500">$</span>
              <span class="inline-block w-2 h-4 bg-green-400 ml-1 animate-pulse"></span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
defineProps({
  logs:    { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const selectedLog = ref(null)
function openModal(log) { selectedLog.value = log }
function closeModal() { selectedLog.value = null }
function onKeydown(e) { if (e.key === 'Escape') closeModal() }
onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

const avatarColors = ['#006d35', '#1b8848', '#40617f', '#001d32', '#7c3aed', '#0891b2', '#b45309']
function avatarColor(name) {
  if (!name) return '#94a3b8'
  return avatarColors[name.charCodeAt(0) % avatarColors.length]
}
function actionLabel(action) {
  const labels = {
    'user.banned': 'a banni un utilisateur', 'user.activated': 'a activé un utilisateur',
    'user.deleted': 'a supprimé un utilisateur', 'provider.approved': 'a approuvé un prestataire',
    'provider.suspended': 'a suspendu un prestataire', 'admin.created': 'a créé un administrateur',
    'admin.deleted': 'a supprimé un administrateur', 'category.created': 'a créé une catégorie',
    'category.deleted': 'a supprimé une catégorie', 'category.toggled': 'a modifié une catégorie',
  }
  return labels[action] || action
}
function actionBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'bg-[#fee2e2] text-[#991b1b]'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'bg-[#dcfce7] text-[#166534]'
  if (action.includes('suspended')) return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function statusBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#dcfce7] text-[#166534]'
}
function statusLabel(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'Warning'
  return 'Success'
}
function formatDate(dateStr) {
  const d = new Date(dateStr), now = new Date(), diff = Math.floor((now - d) / 1000)
  if (diff < 60) return "à l'instant"
  if (diff < 3600) return `il y a ${Math.floor(diff / 60)}min`
  if (diff < 86400) return `il y a ${Math.floor(diff / 3600)}h`
  return d.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
function formatDateFull(dateStr) {
  const d = new Date(dateStr)
  return d.toLocaleDateString('fr-FR', {
    weekday: 'short', day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  }) + ` (UTC${d.getTimezoneOffset() <= 0 ? '+' : '-'}${Math.abs(d.getTimezoneOffset() / 60)})`
}
function formatJson(raw) {
  try { return JSON.stringify(JSON.parse(raw), null, 2) } catch { return raw }
}
function actionTerminalColor(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'text-red-400'
  if (action.includes('suspended')) return 'text-orange-400'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'text-green-400'
  return 'text-gray-300'
}
</script>

<style scoped>
.modal-enter-active, .modal-leave-active { transition: opacity 0.2s ease; }
.modal-enter-from, .modal-leave-to { opacity: 0; }
</style>
