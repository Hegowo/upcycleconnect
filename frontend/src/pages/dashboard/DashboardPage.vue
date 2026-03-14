<template>
  <div class="space-y-6">
    <!-- Stats cards -->
    <div class="grid grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-4">
      <StatCard
        label="Utilisateurs"
        :value="stats?.users_count ?? '—'"
        color="blue"
        link="/admin/users"
      />
      <StatCard
        label="Prestataires validés"
        :value="stats?.providers_count ?? '—'"
        color="green"
        link="/admin/providers?status=approved"
      />
      <StatCard
        label="En attente validation"
        :value="stats?.providers_pending ?? '—'"
        color="orange"
        link="/admin/providers?status=pending"
        :alert="(stats?.providers_pending ?? 0) > 0"
      />
      <StatCard
        label="Prestations"
        :value="stats?.prestations_count ?? '—'"
        color="purple"
        link="/admin/prestations"
      />
      <StatCard
        label="Événements actifs"
        :value="stats?.events_count ?? '—'"
        color="teal"
        link="/admin/events"
      />
    </div>

    <!-- Prestations by status -->
    <div v-if="stats" class="card p-6">
      <h2 class="text-base font-semibold text-gray-900 mb-4">Répartition des prestations par statut</h2>
      <div class="space-y-3">
        <div
          v-for="(count, status) in stats.prestations_by_status"
          :key="status"
          class="flex items-center justify-between"
        >
          <div class="flex items-center gap-2">
            <AppBadge :label="status" />
          </div>
          <span class="text-sm font-medium text-gray-700">{{ count }}</span>
        </div>
        <p v-if="!Object.keys(stats.prestations_by_status || {}).length" class="text-sm text-gray-400">
          Aucune prestation.
        </p>
      </div>
    </div>

    <!-- Prestataires en attente -->
    <div v-if="(stats?.providers_pending ?? 0) > 0" class="card p-4 border-l-4 border-yellow-400 bg-yellow-50">
      <p class="text-sm text-yellow-800 font-medium">
        {{ stats.providers_pending }} prestataire(s) en attente de validation.
      </p>
      <RouterLink to="/admin/providers?status=pending" class="text-sm text-yellow-700 underline mt-1 inline-block">
        Voir les demandes →
      </RouterLink>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'
import AppBadge from '@/components/AppBadge.vue'
import StatCard from './StatCard.vue'

const stats   = ref(null)
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await api.get('/dashboard/stats')
    stats.value = data
  } finally {
    loading.value = false
  }
})
</script>
