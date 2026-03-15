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

        <dl class="grid grid-cols-2 gap-x-6 gap-y-4 text-sm border-t border-gray-100 pt-6">
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

async function fetchUser() {
  try {
    const data = await userService.get(route.params.id)
    user.value = data
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

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

onMounted(fetchUser)
</script>
