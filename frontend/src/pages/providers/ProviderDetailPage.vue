<template>
  <div class="space-y-4">
    <button @click="router.back()" class="text-sm text-gray-500 hover:text-gray-700">← Retour</button>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
    </div>

    <template v-else-if="provider">
      <div class="card p-6">
        <div class="flex items-start justify-between mb-4">
          <div>
            <h2 class="text-xl font-bold text-gray-900">{{ provider.profile?.company_name }}</h2>
            <p class="text-gray-500 text-sm">{{ provider.first_name }} {{ provider.last_name }} — {{ provider.email }}</p>
          </div>
          <AppBadge :label="provider.profile?.status || '—'" />
        </div>

        <dl class="grid grid-cols-2 gap-4 text-sm mb-6">
          <div>
            <dt class="text-gray-500">SIRET</dt>
            <dd class="text-gray-900 font-medium">{{ provider.profile?.siret || '—' }}</dd>
          </div>
          <div>
            <dt class="text-gray-500">Site web</dt>
            <dd class="text-gray-900 font-medium">
              <a v-if="provider.profile?.website" :href="provider.profile.website" target="_blank" rel="noopener" class="text-blue-600 hover:underline">
                {{ provider.profile.website }}
              </a>
              <span v-else>—</span>
            </dd>
          </div>
          <div class="col-span-2">
            <dt class="text-gray-500">Description</dt>
            <dd class="text-gray-900">{{ provider.profile?.description || '—' }}</dd>
          </div>
        </dl>

        <div class="flex gap-3 pt-4 border-t border-gray-100">
          <button
            v-if="provider.profile?.status !== 'approved'"
            @click="changeStatus('approved')"
            class="btn-primary text-sm px-4 py-2 rounded-md"
          >
            Approuver
          </button>
          <button
            v-if="provider.profile?.status !== 'suspended'"
            @click="changeStatus('suspended')"
            class="btn-danger text-sm px-4 py-2 rounded-md"
          >
            Suspendre
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { providerService } from '@/services/userService'
import { useToast } from '@/composables/useToast'
import AppBadge from '@/components/AppBadge.vue'

const route    = useRoute()
const router   = useRouter()
const toast    = useToast()
const provider = ref(null)
const loading  = ref(true)

async function fetchProvider() {
  try {
    provider.value = await providerService.get(route.params.id)
  } catch {
    toast.showError('Prestataire introuvable.')
    router.push('/admin/providers')
  } finally {
    loading.value = false
  }
}

async function changeStatus(status) {
  try {
    await providerService.updateStatus(provider.value.id, status)
    provider.value.profile.status = status
    toast.showSuccess(`Statut mis à jour : ${status}`)
  } catch {
    toast.showError('Impossible de modifier le statut.')
  }
}

onMounted(fetchProvider)
</script>
