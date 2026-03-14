<template>
  <div class="space-y-4">
    <div class="flex items-center gap-3">
      <button @click="router.back()" class="text-sm text-gray-500 hover:text-gray-700">← Retour</button>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
    </div>

    <template v-else-if="user">
      <div class="card p-6">
        <div class="flex items-start justify-between mb-6">
          <div>
            <h2 class="text-xl font-bold text-gray-900">{{ user.first_name }} {{ user.last_name }}</h2>
            <p class="text-gray-500 text-sm">{{ user.email }}</p>
          </div>
          <AppBadge :label="user.status" />
        </div>

        <dl class="grid grid-cols-2 gap-4 text-sm">
          <div>
            <dt class="text-gray-500">Téléphone</dt>
            <dd class="text-gray-900 font-medium">{{ user.phone || '—' }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Email vérifié</dt>
            <dd class="text-gray-900 font-medium">{{ user.email_verified_at ? 'Oui' : 'Non' }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Inscrit le</dt>
            <dd class="text-gray-900 font-medium">{{ formatDate(user.created_at) }}</dd>
          </div>
        </dl>

        <div class="mt-6 flex gap-3 pt-4 border-t border-gray-100">
          <button
            v-if="user.status !== 'banned'"
            @click="confirmAction('ban')"
            class="btn-danger text-sm px-4 py-2 rounded-md"
          >
            Bannir l'utilisateur
          </button>
          <button
            v-else
            @click="confirmAction('activate')"
            class="btn-primary text-sm px-4 py-2 rounded-md"
          >
            Activer l'utilisateur
          </button>
        </div>
      </div>
    </template>

    <AppConfirmDialog
      :show="confirm.show"
      :title="confirm.action === 'ban' ? 'Bannir l\'utilisateur' : 'Activer l\'utilisateur'"
      :message="confirm.action === 'ban' ? 'Cette action bannit l\'utilisateur de la plateforme.' : 'Cette action réactive l\'accès de l\'utilisateur.'"
      :confirm-label="confirm.action === 'ban' ? 'Bannir' : 'Activer'"
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
import { userService } from '@/services/userService'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/common/AppBadge.vue'
import AppConfirmDialog from '@/components/common/AppConfirmDialog.vue'

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
    toast.showError('Utilisateur introuvable.')
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
      toast.showSuccess('Utilisateur banni.')
    } else {
      await userService.activate(user.value.id)
      user.value.status = 'active'
      toast.showSuccess('Utilisateur activé.')
    }
    confirm.show = false
  } catch {
    toast.showError('L\'action a échoué.')
  } finally {
    confirm.loading = false
  }
}

function formatDate(iso) {
  return new Date(iso).toLocaleDateString('fr-FR', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

onMounted(fetchUser)
</script>
