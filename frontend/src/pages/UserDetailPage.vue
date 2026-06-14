<template>
  <div class="max-w-2xl space-y-6">
    <RouterLink to="/admin/users" class="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700">
      {{ t('common.back') }}
    </RouterLink>

    <div v-if="loading" class="card p-8 space-y-4">
      <div class="flex items-center gap-4">
        <div class="w-16 h-16 rounded-full bg-gray-100 animate-pulse"></div>
        <div class="space-y-2">
          <div class="h-5 bg-gray-100 rounded animate-pulse w-40"></div>
          <div class="h-4 bg-gray-100 rounded animate-pulse w-56"></div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4 pt-4">
        <div v-for="n in 4" :key="n" class="h-10 bg-gray-100 rounded animate-pulse"></div>
      </div>
    </div>

    <template v-else-if="user">

      <div class="card p-8">
        <div class="flex items-start justify-between mb-6">
          <div class="flex items-center gap-4">
            <div class="w-16 h-16 rounded-full flex items-center justify-center text-white text-xl font-bold shrink-0" style="background-color:#1B5B88;">
              {{ (user.first_name?.[0] || '') + (user.last_name?.[0] || '') }}
            </div>
            <div>
              <h2 class="text-xl font-bold text-gray-900">{{ user.first_name }} {{ user.last_name }}</h2>
              <p class="text-sm text-gray-400 mt-0.5">{{ user.email }}</p>
            </div>
          </div>
          <AppBadge :label="user.status" />
        </div>

        <dl class="grid grid-cols-1 sm:grid-cols-2 gap-x-6 gap-y-4 text-sm border-t border-gray-100 pt-6">
          <div>
            <dt class="text-xs text-gray-400 uppercase font-medium mb-1">{{ t('users.fieldPhone') }}</dt>
            <dd class="text-gray-900 font-medium">{{ user.phone || '—' }}</dd>
          </div>
          <div>
            <dt class="text-xs text-gray-400 uppercase font-medium mb-1">{{ t('users.fieldVerified') }}</dt>
            <dd class="font-medium" :class="user.email_verified_at ? 'text-green-600' : 'text-gray-400'">
              {{ user.email_verified_at ? `✓ ${t('users.verified')}` : t('users.notVerified') }}
            </dd>
          </div>
          <div>
            <dt class="text-xs text-gray-400 uppercase font-medium mb-1">{{ t('users.fieldJoined') }}</dt>
            <dd class="text-gray-900 font-medium">{{ formatDate(user.created_at) }}</dd>
          </div>
          <div>
            <dt class="text-xs text-gray-400 uppercase font-medium mb-1">{{ t('users.fieldStatus') }}</dt>
            <dd>
              <span class="text-xs px-2 py-1 rounded-full font-semibold"
                :class="{
                  'bg-green-100 text-green-700': user.status === 'active',
                  'bg-red-100 text-red-700': user.status === 'banned',
                  'bg-gray-100 text-gray-500': user.status === 'inactive',
                }">
                {{ user.status }}
              </span>
            </dd>
          </div>
        </dl>

        <div class="mt-6 flex gap-3 pt-4 border-t border-gray-100">
          <button
            v-if="user.status !== 'banned'"
            @click="confirmAction('ban')"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90"
            style="background-color:#dc2626;"
          >
            {{ t('users.actionBan') }}
          </button>
          <button
            v-else
            @click="confirmAction('activate')"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-semibold transition hover:opacity-90"
            style="background-color:#1B8848;"
          >
            {{ t('users.actionActivate') }}
          </button>
        </div>
      </div>

      <div class="card p-6 space-y-4">
        <div class="flex items-center gap-2 mb-1">
          <EnvelopeIcon class="w-4 h-4 text-[#006d35]" />
          <h3 class="font-semibold text-gray-800 text-sm">Modifier l'adresse email</h3>
        </div>
        <div class="flex gap-3">
          <input
            v-model="newEmail"
            type="email"
            placeholder="nouvelle@email.com"
            class="flex-1 px-4 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/20 focus:border-[#006d35] transition bg-white"
          />
          <button
            @click="doUpdateEmail"
            :disabled="emailLoading || !newEmail || newEmail === user.email"
            class="px-4 py-2.5 rounded-xl text-white text-sm font-semibold transition hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            style="background-color:#006d35;"
          >
            <div v-if="emailLoading" class="w-3 h-3 border-2 border-white border-t-transparent rounded-full animate-spin" />
            <CheckIcon v-else class="w-4 h-4" />
            Enregistrer
          </button>
        </div>
        <Transition name="fade">
          <p v-if="emailSuccess" class="text-green-600 text-xs bg-green-50 p-2.5 rounded-xl border border-green-200">{{ emailSuccess }}</p>
        </Transition>
        <Transition name="fade">
          <p v-if="emailError" class="text-red-600 text-xs bg-red-50 p-2.5 rounded-xl border border-red-200">{{ emailError }}</p>
        </Transition>
      </div>

      <div class="card p-6 space-y-4">
        <h3 class="font-semibold text-gray-800 text-sm mb-2">Actions administrateur</h3>

        <div class="flex flex-col sm:flex-row gap-3">

          <button
            @click="doSendReset"
            :disabled="resetLoading"
            class="flex-1 px-4 py-3 rounded-xl border border-[#e2e8f0] bg-white text-sm font-semibold text-gray-700 hover:bg-gray-50 transition flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <div v-if="resetLoading" class="w-3 h-3 border-2 border-gray-500 border-t-transparent rounded-full animate-spin" />
            <KeyIcon v-else class="w-4 h-4 text-[#40617f]" />
            Envoyer un email de récupération
          </button>

          <button
            @click="doExportPdf"
            :disabled="exportLoading"
            class="flex-1 px-4 py-3 rounded-xl border border-[#e2e8f0] bg-white text-sm font-semibold text-gray-700 hover:bg-gray-50 transition flex items-center justify-center gap-2 disabled:opacity-50"
          >
            <div v-if="exportLoading" class="w-3 h-3 border-2 border-gray-500 border-t-transparent rounded-full animate-spin" />
            <ArrowDownTrayIcon v-else class="w-4 h-4 text-[#40617f]" />
            Exporter en PDF
          </button>
        </div>

        <Transition name="fade">
          <p v-if="resetSuccess" class="text-green-600 text-xs bg-green-50 p-2.5 rounded-xl border border-green-200">{{ resetSuccess }}</p>
        </Transition>
        <Transition name="fade">
          <p v-if="resetError" class="text-red-600 text-xs bg-red-50 p-2.5 rounded-xl border border-red-200">{{ resetError }}</p>
        </Transition>
      </div>
    </template>

    <AppConfirmDialog
      :show="confirm.show"
      :title="confirm.action === 'ban' ? t('users.confirmBanTitle') : t('users.confirmActivateTitle')"
      :message="confirm.action === 'ban' ? t('users.confirmBanMsg', { name: `${user?.first_name} ${user?.last_name}` }) : t('users.confirmActivateMsg', { name: `${user?.first_name} ${user?.last_name}` })"
      :confirm-label="confirm.action === 'ban' ? t('users.actionBan') : t('users.actionActivate')"
      :confirm-variant="confirm.action === 'ban' ? 'danger' : 'primary'"
      :loading="confirm.loading"
      @confirm="executeAction"
      @cancel="confirm.show = false"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { EnvelopeIcon, CheckIcon, KeyIcon, ArrowDownTrayIcon } from '@heroicons/vue/24/outline'
import { userService } from '@/services/userService'
import { useToast } from '@/utils/useToast'
import AppBadge from '@/components/AppBadge.vue'
import AppConfirmDialog from '@/components/AppConfirmDialog.vue'

const { t } = useI18n()
const route  = useRoute()
const router = useRouter()
const toast  = useToast()

const user    = ref(null)
const loading = ref(true)
const confirm = reactive({ show: false, action: '', loading: false })

const newEmail     = ref('')
const emailLoading = ref(false)
const emailSuccess = ref('')
const emailError   = ref('')

const resetLoading = ref(false)
const resetSuccess = ref('')
const resetError   = ref('')

const exportLoading = ref(false)

async function fetchUser() {
  try {
    const data = await userService.get(route.params.id)
    user.value  = data
    newEmail.value = data.email
  } catch {
    toast.showError(t('users.toastLoadError'))
    router.push('/admin/users')
  } finally {
    loading.value = false
  }
}

function confirmAction(action) {
  confirm.action = action
  confirm.show   = true
}

async function executeAction() {
  confirm.loading = true
  try {
    if (confirm.action === 'ban') {
      await userService.ban(user.value.id)
      user.value.status = 'banned'
      toast.showSuccess(t('users.toastBanned', { name: `${user.value.first_name} ${user.value.last_name}` }))
    } else {
      await userService.activate(user.value.id)
      user.value.status = 'active'
      toast.showSuccess(t('users.toastActivated', { name: `${user.value.first_name} ${user.value.last_name}` }))
    }
    confirm.show = false
  } catch {
    toast.showError(t('users.toastError'))
  } finally {
    confirm.loading = false
  }
}

async function doUpdateEmail() {
  if (!newEmail.value || newEmail.value === user.value.email) return
  emailLoading.value = true
  emailSuccess.value = ''
  emailError.value   = ''
  try {
    const updated = await userService.updateEmail(user.value.id, newEmail.value)
    user.value.email = updated.email
    user.value.email_verified_at = updated.email_verified_at
    emailSuccess.value = 'Email mis à jour avec succès.'
    setTimeout(() => { emailSuccess.value = '' }, 4000)
  } catch (e) {
    emailError.value = e?.response?.data?.message || 'Erreur lors de la mise à jour.'
    setTimeout(() => { emailError.value = '' }, 5000)
  } finally {
    emailLoading.value = false
  }
}

async function doSendReset() {
  resetLoading.value = true
  resetSuccess.value = ''
  resetError.value   = ''
  try {
    await userService.sendPasswordReset(user.value.id)
    resetSuccess.value = `Email de récupération envoyé à ${user.value.email}.`
    setTimeout(() => { resetSuccess.value = '' }, 5000)
  } catch (e) {
    resetError.value = e?.response?.data?.message || "Erreur lors de l'envoi."
    setTimeout(() => { resetError.value = '' }, 5000)
  } finally {
    resetLoading.value = false
  }
}

async function doExportPdf() {
  exportLoading.value = true
  try {
    const blob = await userService.exportPdf(user.value.id)
    const url  = URL.createObjectURL(blob)
    const a    = document.createElement('a')
    a.href     = url
    a.download = `utilisateur_${user.value.id}_${new Date().toISOString().slice(0,10)}.pdf`
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    toast.showError('Erreur lors de la génération du PDF.')
  } finally {
    exportLoading.value = false
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

onMounted(fetchUser)
</script>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
