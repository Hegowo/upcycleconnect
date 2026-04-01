<template>
  <div class="space-y-6">
    
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-[#001d32]">{{ t('logs.title') }}</h2>
      </div>
    </div>

    
    <div>
      <div>
        <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
          <div class="flex items-center justify-between px-6 py-4 border-b border-[#f1f5f9]">
            <h3 class="font-semibold text-[#001d32]">{{ t('logs.activityFeed') }}</h3>
            <div class="flex items-center gap-2">
              <button @click="exportLogs" class="px-3 py-1.5 text-xs font-semibold rounded-lg border border-[#e5e7eb] text-[#374151] hover:bg-gray-50 transition flex items-center gap-1.5">
                <ArrowDownTrayIcon class="w-3.5 h-3.5" />
                {{ t('common.exportCsv') }}
              </button>
            </div>
          </div>

          
          <div class="flex gap-1 p-3 bg-[#f8fafc] border-b border-[#f1f5f9]">
            <button @click="tab='admin'" :class="tab==='admin' ? 'bg-white shadow text-[#001d32] font-semibold' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm transition">
              {{ t('logs.tabAdmin') }}
            </button>
            <button @click="tab='platform'" :class="tab==='platform' ? 'bg-white shadow text-[#001d32] font-semibold' : 'text-gray-500 hover:text-gray-700'" class="px-4 py-2 rounded-lg text-sm transition">
              {{ t('logs.tabPlatform') }}
            </button>
          </div>

          
          <div v-if="tab === 'admin'" class="overflow-x-auto">
            <div v-if="adminLoading" class="p-4 space-y-3">
              <div v-for="n in 6" :key="n" class="h-12 bg-gray-50 rounded-lg animate-pulse"></div>
            </div>
            <div v-else-if="!adminLogs.length" class="py-12 text-center text-gray-400 text-sm">
              <ClockIcon class="w-8 h-8 mx-auto mb-2 text-gray-300" />
              {{ t('logs.adminEmpty') }}
            </div>
            <table v-else class="min-w-full hidden lg:table">
              <thead>
                <tr class="bg-[#f8fafc]">
                  <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colTimestamp') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colAdmin') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colAction') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colEntity') }}</th>
                  <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colStatus') }}</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-[#f8fafc]">
                <tr
                  v-for="log in adminLogs"
                  :key="log.id"
                  class="hover:bg-[#f8fafc] cursor-pointer transition-colors"
                  @click="openModal(log, 'admin')"
                >
                  <td class="px-6 py-3.5 text-xs font-mono text-gray-500 whitespace-nowrap">
                    {{ formatDateShort(log.created_at) }}
                  </td>
                  <td class="px-6 py-3.5">
                    <div class="flex items-center gap-2">
                      <div class="w-7 h-7 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0"
                        :style="{ backgroundColor: avatarColor(log.admin_name) }">
                        {{ (log.admin_name?.trim() || 'S').charAt(0).toUpperCase() }}
                      </div>
                      <span class="text-sm font-medium text-[#001d32]">{{ log.admin_name?.trim() || t('logs.system') }}</span>
                    </div>
                  </td>
                  <td class="px-6 py-3.5 text-sm text-gray-600">{{ actionLabel(log.action) }}</td>
                  <td class="px-6 py-3.5">
                    <span v-if="log.resource_type" class="text-xs font-medium px-2 py-0.5 rounded-full" :class="actionBadgeClass(log.action)">
                      {{ log.resource_type }}
                      <span v-if="log.resource_id" class="opacity-60"> #{{ log.resource_id }}</span>
                    </span>
                    <span v-else class="text-xs text-gray-400">—</span>
                  </td>
                  <td class="px-6 py-3.5">
                    <span class="text-xs font-semibold px-2 py-0.5 rounded-full" :class="statusBadge(log.action)">
                      {{ statusText(log.action) }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>

            
            <div v-if="!adminLoading && adminLogs.length" class="lg:hidden divide-y divide-[#f8fafc]">
              <div
                v-for="log in adminLogs"
                :key="log.id"
                class="flex items-start gap-3 p-4 hover:bg-[#f8fafc] cursor-pointer transition-colors"
                @click="openModal(log, 'admin')"
              >
                <div class="w-8 h-8 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0"
                  :style="{ backgroundColor: avatarColor(log.admin_name) }">
                  {{ (log.admin_name?.trim() || 'S').charAt(0).toUpperCase() }}
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between gap-2 mb-1">
                    <p class="text-sm font-medium text-[#001d32] truncate">{{ log.admin_name?.trim() || t('logs.system') }}</p>
                    <span class="text-xs font-semibold px-2 py-0.5 rounded-full shrink-0" :class="statusBadge(log.action)">
                      {{ statusText(log.action) }}
                    </span>
                  </div>
                  <p class="text-xs text-gray-600 mb-1">{{ actionLabel(log.action) }}</p>
                  <div class="flex items-center gap-2">
                    <span v-if="log.resource_type" class="text-xs font-medium px-1.5 py-0.5 rounded-full" :class="actionBadgeClass(log.action)">
                      {{ log.resource_type }}<span v-if="log.resource_id" class="opacity-60"> #{{ log.resource_id }}</span>
                    </span>
                    <span class="text-xs text-gray-400 font-mono">{{ formatDateShort(log.created_at) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          
          <div v-if="tab === 'platform'">
            <div v-if="platformLoading" class="p-4 space-y-3">
              <div v-for="n in 5" :key="n" class="h-12 bg-gray-50 rounded-lg animate-pulse"></div>
            </div>
            <div v-else-if="!platformLogs.length" class="py-12 text-center text-gray-400 text-sm">
              <UsersIcon class="w-8 h-8 mx-auto mb-2 text-gray-300" />
              {{ t('logs.platformEmpty') }}
            </div>
            <div v-else class="divide-y divide-[#f8fafc]">
              <div
                v-for="entry in platformLogs"
                :key="entry.id + entry.type"
                class="flex items-center gap-4 px-6 py-3.5 hover:bg-[#f8fafc] transition-colors cursor-pointer"
                @click="openModal(entry, 'platform')"
              >
                <div class="w-8 h-8 rounded-full flex items-center justify-center shrink-0"
                  :class="entry.type === 'provider' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'">
                  <span class="text-xs font-bold">{{ entry.type === 'provider' ? 'P' : 'U' }}</span>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-[#001d32]">{{ entry.name }}</p>
                  <p class="text-xs text-gray-400">{{ entry.email }} · {{ formatDate(entry.created_at) }}</p>
                </div>
                <span class="text-xs font-semibold px-2 py-0.5 rounded-full"
                  :class="entry.type === 'provider' ? 'bg-purple-100 text-purple-700' : 'bg-blue-100 text-blue-700'">
                  {{ entry.type === 'provider' ? t('logs.typeProvider') : t('logs.typeUser') }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    
    <div class="bg-white rounded-2xl border border-[#f1f5f9] shadow-sm overflow-hidden">
      <div class="px-6 py-4 border-b border-[#f1f5f9] flex items-center justify-between">
        <div>
          <h3 class="font-semibold text-[#001d32]">{{ t('logs.staffGovernance') }}</h3>
          <p class="text-xs text-gray-400 mt-0.5">{{ t('logs.staffSubtitle') }}</p>
        </div>
      </div>

      <div v-if="!auth.isSuperAdmin" class="px-6 py-8 text-center text-gray-400 text-sm">
        <ShieldCheckIcon class="w-8 h-8 mx-auto mb-2 text-gray-300" />
        {{ t('logs.restrictedAccess') }}
      </div>

      <template v-else>
        <div v-if="adminsLoading" class="p-4 space-y-3 overflow-x-auto">
          <div v-for="n in 3" :key="n" class="h-12 bg-gray-50 rounded-lg animate-pulse"></div>
        </div>
        <div v-else-if="admins.length">
          <table class="min-w-full hidden lg:table">
            <thead>
              <tr class="bg-[#f8fafc]">
                <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colAdministrator') }}</th>
                <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colLevel') }}</th>
                <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('users.colEmail') }}</th>
                <th class="px-6 py-3 text-left text-xs font-semibold text-[#6b7280] uppercase tracking-wider">{{ t('logs.colSince') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#f8fafc]">
              <tr v-for="admin in admins" :key="admin.id" class="hover:bg-[#f8fafc] transition-colors">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="w-9 h-9 rounded-full flex items-center justify-center text-white text-xs font-bold" style="background-color:#001d32;">
                      {{ ((admin.first_name?.[0] || '') + (admin.last_name?.[0] || '')).toUpperCase() }}
                    </div>
                    <span class="text-sm font-semibold text-[#001d32]">{{ admin.first_name }} {{ admin.last_name }}</span>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <span class="text-xs font-semibold px-2.5 py-1 rounded-full capitalize"
                    :class="admin.role === 'super_admin' ? 'bg-[#fef9c3] text-[#854d0e]' : 'bg-[#dcfce7] text-[#166534]'">
                    {{ admin.role === 'super_admin' ? 'Super Admin' : 'Admin' }}
                  </span>
                </td>
                <td class="px-6 py-4 text-sm text-gray-400">{{ admin.email }}</td>
                <td class="px-6 py-4 text-sm text-gray-400">{{ formatDateShort(admin.created_at) }}</td>
              </tr>
            </tbody>
          </table>

          
          <div class="lg:hidden divide-y divide-[#f8fafc]">
            <div v-for="admin in admins" :key="admin.id" class="flex items-center gap-3 px-5 py-4">
              <div class="w-10 h-10 rounded-full flex items-center justify-center text-white text-xs font-bold shrink-0" style="background-color:#001d32;">
                {{ ((admin.first_name?.[0] || '') + (admin.last_name?.[0] || '')).toUpperCase() }}
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm font-semibold text-[#001d32]">{{ admin.first_name }} {{ admin.last_name }}</p>
                <p class="text-xs text-gray-400 truncate">{{ admin.email }}</p>
              </div>
              <div class="flex flex-col items-end gap-1">
                <span class="text-xs font-semibold px-2 py-0.5 rounded-full capitalize"
                  :class="admin.role === 'super_admin' ? 'bg-[#fef9c3] text-[#854d0e]' : 'bg-[#dcfce7] text-[#166534]'">
                  {{ admin.role === 'super_admin' ? 'Super Admin' : 'Admin' }}
                </span>
                <span class="text-xs text-gray-400">{{ formatDateShort(admin.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="px-6 py-8 text-center text-gray-400 text-sm">{{ t('logs.noAdmins') }}</div>
      </template>
    </div>

    
    <Teleport to="body">
      <Transition name="modal">
        <div v-if="selectedLog" class="fixed inset-0 z-50 flex items-end sm:items-center justify-center sm:p-4" @click.self="closeModal">
          <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>
          <div class="relative w-full sm:max-w-2xl rounded-t-2xl sm:rounded-xl overflow-hidden shadow-2xl border border-gray-700 max-h-[85vh] overflow-y-auto">
            <div class="flex items-center gap-2 px-4 py-3 bg-gray-800 border-b border-gray-700">
              <button @click="closeModal" class="w-3 h-3 rounded-full bg-red-500 hover:bg-red-400 transition-colors"></button>
              <div class="w-3 h-3 rounded-full bg-yellow-500"></div>
              <div class="w-3 h-3 rounded-full bg-green-500"></div>
              <span class="ml-3 text-xs font-mono text-gray-400">upcycleconnect — log#{{ selectedLog.id }}</span>
              <button @click="closeModal" class="ml-auto text-gray-500 hover:text-gray-300 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
              </button>
            </div>
            <div class="bg-gray-900 p-5 font-mono text-sm leading-relaxed">
              <div class="flex items-center gap-2 mb-4 pb-3 border-b border-gray-700">
                <span class="text-green-400">root@upcycleconnect</span>
                <span class="text-gray-500">:</span>
                <span class="text-blue-400">~/logs</span>
                <span class="text-gray-500">$</span>
                <span class="text-white">cat log.{{ selectedLog.id }}.json</span>
              </div>
              <template v-if="selectedLogType === 'admin'">
                <div class="space-y-1.5">
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">id</span><span class="text-yellow-300">{{ selectedLog.id }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">timestamp</span><span class="text-cyan-300">{{ formatDateFull(selectedLog.created_at) }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">ip_address</span><span class="text-green-300">{{ selectedLog.ip_address || 'unknown' }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">admin</span><span class="text-white">{{ selectedLog.admin_name?.trim() || 'system' }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">admin_id</span><span class="text-yellow-300">{{ selectedLog.admin_id ?? '—' }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">action</span><span :class="actionTerminalColor(selectedLog.action)" class="font-bold">{{ selectedLog.action }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">resource_type</span><span class="text-purple-300">{{ selectedLog.resource_type || '—' }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">resource_id</span><span class="text-yellow-300">{{ selectedLog.resource_id ?? '—' }}</span></div>
                  <template v-if="selectedLog.old_values">
                    <div class="flex gap-3 mt-2 pt-2 border-t border-gray-700">
                      <span class="text-gray-500 w-28 shrink-0">old_values</span>
                      <span class="text-red-300 break-all text-xs font-mono whitespace-pre">{{ formatJson(selectedLog.old_values) }}</span>
                    </div>
                  </template>
                  <template v-if="selectedLog.new_values">
                    <div class="flex gap-3">
                      <span class="text-gray-500 w-28 shrink-0">new_values</span>
                      <span class="text-green-300 break-all text-xs font-mono whitespace-pre">{{ formatJson(selectedLog.new_values) }}</span>
                    </div>
                  </template>
                </div>
              </template>
              <template v-else>
                <div class="space-y-1.5">
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">id</span><span class="text-yellow-300">{{ selectedLog.id }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">timestamp</span><span class="text-cyan-300">{{ formatDateFull(selectedLog.created_at) }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">type</span><span :class="selectedLog.type === 'provider' ? 'text-purple-300' : 'text-blue-300'" class="font-bold">{{ selectedLog.type }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">name</span><span class="text-white">{{ selectedLog.name }}</span></div>
                  <div class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">email</span><span class="text-green-300">{{ selectedLog.email }}</span></div>
                  <div v-if="selectedLog.company" class="flex gap-3"><span class="text-gray-500 w-28 shrink-0">company</span><span class="text-gray-300">{{ selectedLog.company }}</span></div>
                </div>
              </template>
              <div class="mt-4 pt-3 border-t border-gray-700 flex items-center gap-1">
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
  </div>
</template>

<script setup>import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import {
  ClockIcon, UsersIcon, ShieldCheckIcon,
  ArrowDownTrayIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const auth = useAuthStore()
const tab = ref('admin')
const adminLogs = ref([])
const platformLogs = ref([])
const adminLoading = ref(false)
const platformLoading = ref(false)
const admins = ref([])
const adminsLoading = ref(false)
const selectedLog = ref(null)
const selectedLogType = ref('admin')

function openModal(log, type) { selectedLog.value = log; selectedLogType.value = type }
function closeModal() { selectedLog.value = null }
function onKeydown(e) { if (e.key === 'Escape') closeModal() }
onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

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

async function fetchAdmins() {
  if (!auth.isSuperAdmin) return
  adminsLoading.value = true
  try {
    const { data } = await api.get('/admins')
    admins.value = data.data || []
  } finally {
    adminsLoading.value = false
  }
}

onMounted(() => {
  loadAdminLogs()
  loadPlatformLogs()
  fetchAdmins()
})

function exportLogs() {
  const source = tab.value === 'admin' ? adminLogs.value : platformLogs.value
  if (!source.length) return
  let rows
  if (tab.value === 'admin') {
    rows = [
      ['ID', 'Timestamp', 'Admin', 'Action', 'Ressource', 'Ressource ID', 'IP'],
      ...source.map(l => [l.id, l.created_at, l.admin_name || 'Système', l.action, l.resource_type || '', l.resource_id || '', l.ip_address || '']),
    ]
  } else {
    rows = [
      ['ID', 'Type', 'Nom', 'Email', 'Date'],
      ...source.map(l => [l.id, l.type, l.name, l.email, l.created_at]),
    ]
  }
  const csv = rows.map(r => r.map(v => `"${String(v).replace(/"/g, '""')}"`).join(',')).join('\n')
  const blob = new Blob(['\uFEFF' + csv], { type: 'text/csv;charset=utf-8;' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob)
  a.download = `logs_${tab.value}_${new Date().toISOString().slice(0,10)}.csv`
  a.click()
  URL.revokeObjectURL(a.href)
}

const avatarColors = ['#006d35', '#1b8848', '#40617f', '#001d32', '#7c3aed', '#0891b2', '#b45309']
function avatarColor(name) {
  if (!name) return '#94a3b8'
  return avatarColors[name.charCodeAt(0) % avatarColors.length]
}

function formatDate(dateStr) {
  const d = new Date(dateStr)
  const now = new Date()
  const diff = Math.floor((now - d) / 1000)
  const locale = t('common.locale') || 'fr-FR'
  if (diff < 60) return t('dashboard.actionInstant')
  if (diff < 3600) return t('dashboard.actionMinutesAgo', { n: Math.floor(diff / 60) })
  if (diff < 86400) return t('dashboard.actionHoursAgo', { n: Math.floor(diff / 3600) })
  return d.toLocaleDateString(locale, { day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit' })
}
function formatDateShort(iso) {
  const locale = t('common.locale') || 'fr-FR'
  return new Date(iso).toLocaleDateString(locale, { day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}
function formatDateFull(dateStr) {
  const d = new Date(dateStr)
  const locale = t('common.locale') || 'fr-FR'
  return d.toLocaleDateString(locale, {
    weekday: 'short', day: '2-digit', month: 'short', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit',
  }) + ` (UTC${d.getTimezoneOffset() <= 0 ? '+' : '-'}${Math.abs(d.getTimezoneOffset() / 60)})`
}

function actionLabel(action) {
  const map = {
    'user.banned':        'logs.actionUserBanned',
    'user.activated':     'logs.actionUserActivated',
    'user.deleted':       'logs.actionUserDeleted',
    'provider.approved':  'logs.actionProviderApproved',
    'provider.suspended': 'logs.actionProviderSuspended',
    'admin.created':      'logs.actionAdminCreated',
    'admin.deleted':      'logs.actionAdminDeleted',
    'category.created':   'logs.actionCategoryCreated',
    'category.deleted':   'logs.actionCategoryDeleted',
    'category.toggled':   'logs.actionCategoryToggled',
  }
  return map[action] ? t(map[action]) : action
}
function actionBadgeClass(action) {
  if (action.includes('deleted') || action.includes('banned')) return 'bg-[#fee2e2] text-[#991b1b]'
  if (action.includes('approved') || action.includes('activated') || action.includes('created')) return 'bg-[#dcfce7] text-[#166534]'
  if (action.includes('suspended')) return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#f1f5f9] text-[#475569]'
}
function statusBadge(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'bg-[#fef9c3] text-[#854d0e]'
  return 'bg-[#dcfce7] text-[#166534]'
}
function statusText(action) {
  if (action.includes('deleted') || action.includes('banned') || action.includes('suspended')) return 'Warning'
  return 'Success'
}
function formatJson(raw) {
  try {
    return JSON.stringify(JSON.parse(raw), null, 2)
  } catch {
    return raw
  }
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
