<template>
  <div class="bg-[#f7f9ff] min-h-screen pb-16">

    <div v-if="loading" class="flex items-center justify-center py-40">
      <div class="w-10 h-10 border-4 border-[#006d35] border-t-transparent rounded-full animate-spin" />
    </div>

    <template v-else>

      <div v-if="!profile" class="max-w-2xl mx-auto px-6 py-20 text-center">
        <BuildingStorefrontIcon class="w-16 h-16 text-[#006d35]/30 mx-auto mb-4" />
        <h1 class="font-jakarta font-extrabold text-[#001d32] text-2xl mb-2">Espace professionnel</h1>
        <p class="text-[#40617f] text-sm mb-6">Votre demande de compte professionnel est en cours de traitement par l'équipe UpcycleConnect.</p>
        <div class="bg-amber-50 border border-amber-200 rounded-2xl px-6 py-4 inline-flex items-start gap-3 text-left">
          <ExclamationTriangleIcon class="w-5 h-5 text-amber-500 shrink-0 mt-0.5" />
          <div>
            <p class="font-semibold text-amber-800 text-sm">Validation en attente</p>
            <p class="text-amber-700 text-xs mt-0.5">Statut actuel : <strong>{{ statusLabel(profile?.status || 'pending') }}</strong>. Vous serez notifié par email dès qu'un administrateur aura examiné votre dossier.</p>
          </div>
        </div>
      </div>

      <template v-else>

        <section class="bg-[#edf4ff] px-6 py-10 relative overflow-hidden">
          <div class="max-w-5xl mx-auto flex items-start gap-8 flex-wrap">

            <div class="bg-white rounded-[24px] p-3 w-28 h-28 flex items-center justify-center shrink-0 shadow-lg">
              <img v-if="userAuth.user?.avatar_url" :src="userAuth.user.avatar_url" class="w-full h-full object-cover rounded-[20px]" />
              <BuildingStorefrontIcon v-else class="w-14 h-14 text-[#006d35]/40" />
            </div>

            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-3 mb-2 flex-wrap">
                <span v-if="profile.status === 'approved'" class="bg-[#006d35] text-white text-xs font-bold uppercase tracking-wide px-3 py-1 rounded-full">
                  Certifié
                </span>
                <span v-else class="bg-amber-100 text-amber-700 text-xs font-bold uppercase tracking-wide px-3 py-1 rounded-full">
                  {{ statusLabel(profile.status) }}
                </span>
                <span v-if="profile.siret" class="text-[#40617f] text-xs font-mono">SIRET {{ profile.siret }}</span>
              </div>

              <h1 class="font-jakarta font-extrabold text-[#001d32] text-4xl sm:text-5xl tracking-tight leading-tight mb-1">
                {{ profile.company_name }}
              </h1>
              <p v-if="userAuth.fullName" class="text-[#40617f] text-sm mb-3">{{ userAuth.fullName }}</p>
              <p v-if="profile.description" class="text-[#001d32] text-sm leading-relaxed max-w-xl mb-4">{{ profile.description }}</p>

              <div class="flex gap-2 flex-wrap">
                <RouterLink to="/profil/pro/dashboard"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl text-white text-sm font-bold hover:opacity-90 transition"
                  style="background:linear-gradient(135deg,#006d35,#1b8848);">
                  <ChartBarIcon class="w-4 h-4" /> Tableau de bord
                </RouterLink>
                <RouterLink to="/profil/pro/projets"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border-2 border-[#006d35] text-[#006d35] text-sm font-bold hover:bg-[#edf4ff] transition bg-white">
                  <WrenchScrewdriverIcon class="w-4 h-4" /> Mes projets
                </RouterLink>
                <RouterLink to="/profil/pro/evenements"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border-2 border-[#006d35] text-[#006d35] text-sm font-bold hover:bg-[#edf4ff] transition bg-white">
                  <CalendarDaysIcon class="w-4 h-4" /> Événements
                </RouterLink>
                <RouterLink to="/profil/pro/collecte"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border-2 border-[#006d35] text-[#006d35] text-sm font-bold hover:bg-[#edf4ff] transition bg-white">
                  <QrCodeIcon class="w-4 h-4" /> Collecte
                </RouterLink>
                <RouterLink to="/profil/pro/abonnement"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border-2 border-[#006d35] text-[#006d35] text-sm font-bold hover:bg-[#edf4ff] transition bg-white">
                  <StarIcon class="w-4 h-4" /> Abonnement
                </RouterLink>
                <RouterLink to="/profil/pro/campagnes"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border-2 border-[#006d35] text-[#006d35] text-sm font-bold hover:bg-[#edf4ff] transition bg-white">
                  <MegaphoneIcon class="w-4 h-4" /> Campagnes
                </RouterLink>
                <RouterLink to="/profil/parametres"
                  class="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-[#e5e7eb] text-[#40617f] text-sm font-semibold hover:bg-white transition bg-white">
                  <Cog6ToothIcon class="w-4 h-4" /> Paramètres
                </RouterLink>
              </div>
            </div>

            <div class="flex flex-col gap-3 shrink-0 w-52">
              <div class="bg-white rounded-xl p-4 flex items-center gap-3 shadow-sm">
                <div class="w-10 h-10 rounded-xl bg-[#d1fae5] flex items-center justify-center">
                  <WrenchScrewdriverIcon class="w-5 h-5 text-[#006d35]" />
                </div>
                <div>
                  <p class="font-jakarta font-bold text-[#001d32] text-xl">{{ projectCount }}</p>
                  <p class="text-[#40617f] text-xs">Projets</p>
                </div>
              </div>
              <div class="bg-white rounded-xl p-4 flex items-center gap-3 shadow-sm">
                <div class="w-10 h-10 rounded-xl bg-[#fef3c7] flex items-center justify-center">
                  <CheckBadgeIcon class="w-5 h-5 text-[#d97706]" />
                </div>
                <div>
                  <p class="font-jakarta font-bold text-[#001d32] text-xl">{{ prestationCount }}</p>
                  <p class="text-[#40617f] text-xs">Prestations</p>
                </div>
              </div>
            </div>
          </div>
        </section>

        <div class="max-w-5xl mx-auto px-6 py-10 grid grid-cols-1 lg:grid-cols-3 gap-6">

          <div class="lg:col-span-1 bg-white rounded-[24px] border border-[#edf4ff] p-6 flex flex-col gap-5">
            <div class="flex items-center justify-between">
              <h2 class="font-jakarta font-bold text-[#001d32] text-lg">Informations</h2>
              <button @click="editing = !editing" class="text-xs text-[#006d35] font-semibold hover:underline">
                {{ editing ? 'Annuler' : 'Modifier' }}
              </button>
            </div>

            <template v-if="!editing">
              <dl class="space-y-3 text-sm">
                <div>
                  <dt class="text-[#94a3b8] text-xs uppercase tracking-wider mb-0.5">Société</dt>
                  <dd class="font-semibold text-[#001d32]">{{ profile.company_name }}</dd>
                </div>
                <div v-if="profile.siret">
                  <dt class="text-[#94a3b8] text-xs uppercase tracking-wider mb-0.5">SIRET</dt>
                  <dd class="font-mono text-[#001d32]">{{ profile.siret }}</dd>
                </div>
                <div v-if="profile.website">
                  <dt class="text-[#94a3b8] text-xs uppercase tracking-wider mb-0.5">Site web</dt>
                  <dd><a :href="profile.website" target="_blank" class="text-[#006d35] hover:underline truncate block">{{ profile.website }}</a></dd>
                </div>
                <div>
                  <dt class="text-[#94a3b8] text-xs uppercase tracking-wider mb-0.5">Email</dt>
                  <dd class="text-[#001d32]">{{ userAuth.user?.email }}</dd>
                </div>
                <div v-if="userAuth.user?.phone">
                  <dt class="text-[#94a3b8] text-xs uppercase tracking-wider mb-0.5">Téléphone</dt>
                  <dd class="text-[#001d32]">{{ userAuth.user.phone }}</dd>
                </div>
              </dl>
            </template>

            <form v-else @submit.prevent="saveProfile" class="space-y-3">
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Nom de l'entreprise *</label>
                <input v-model="editForm.company_name" type="text" required
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Description</label>
                <textarea v-model="editForm.description" rows="3"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30 resize-none"
                  placeholder="Décrivez votre activité..." />
              </div>
              <div>
                <label class="block text-xs font-semibold text-[#40617f] uppercase mb-1">Site web</label>
                <input v-model="editForm.website" type="url" placeholder="https://"
                  class="w-full px-3 py-2.5 bg-[#f8fafc] border border-[#e5e7eb] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#006d35]/30" />
              </div>
              <p v-if="saveError" class="text-red-600 text-xs">{{ saveError }}</p>
              <button type="submit" :disabled="saving"
                class="w-full py-2.5 rounded-xl text-white font-bold text-sm hover:opacity-90 disabled:opacity-50 transition"
                style="background:linear-gradient(135deg,#006d35,#1b8848);">
                {{ saving ? 'Enregistrement...' : 'Enregistrer' }}
              </button>
            </form>
          </div>

          <div class="lg:col-span-2 flex flex-col gap-6">
            <div class="bg-white rounded-[24px] border border-[#edf4ff] p-6">
              <div class="flex items-center justify-between mb-5">
                <h2 class="font-jakarta font-bold text-[#001d32] text-xl">Projets d'upcycling</h2>
                <RouterLink to="/profil/pro/projets" class="text-xs text-[#006d35] font-semibold hover:underline">
                  Gérer mes projets →
                </RouterLink>
              </div>

              <div v-if="!realProjects.length" class="text-center py-8">
                <WrenchScrewdriverIcon class="w-10 h-10 text-gray-200 mx-auto mb-2" />
                <p class="text-[#94a3b8] text-sm">Aucun projet pour l'instant.</p>
                <RouterLink to="/profil/pro/projets" class="text-[#006d35] text-sm font-semibold hover:underline mt-1 inline-block">
                  Créer mon premier projet
                </RouterLink>
              </div>

              <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div v-for="p in realProjects.slice(0,4)" :key="p.id"
                  class="flex flex-col rounded-2xl overflow-hidden border border-[#f1f5f9] hover:border-[#006d35] transition">
                  <div class="h-40 bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] relative overflow-hidden">
                    <img v-if="p.cover_image" :src="p.cover_image" class="w-full h-full object-cover" />
                    <WrenchScrewdriverIcon v-else class="absolute inset-0 m-auto w-14 h-14 text-white/30" />
                    <span v-if="p.category" class="absolute top-3 left-3 bg-white/90 text-[#006d35] text-xs font-semibold px-2 py-0.5 rounded-full">{{ p.category }}</span>
                    <span v-if="p.impact_label" class="absolute bottom-3 right-3 bg-[#006d35] text-white text-xs font-semibold px-2 py-0.5 rounded-full">{{ p.impact_label }}</span>
                  </div>
                  <div class="p-3 flex-1">
                    <p class="font-jakarta font-bold text-[#001d32] text-sm line-clamp-1">{{ p.title }}</p>
                    <p class="text-[#40617f] text-xs mt-0.5 line-clamp-2">{{ p.description }}</p>
                    <p class="text-xs text-[#94a3b8] mt-1">{{ p.step_count }} étape{{ p.step_count !== 1 ? 's' : '' }}</p>
                  </div>
                </div>
              </div>
            </div>

            <div class="bg-white rounded-[24px] border border-[#edf4ff] p-6">
              <div class="flex items-center justify-between mb-5">
                <h2 class="font-jakarta font-bold text-[#001d32] text-xl">Mes prestations</h2>
                <RouterLink to="/prestations" class="text-xs text-[#006d35] font-semibold hover:underline">
                  Voir tout →
                </RouterLink>
              </div>
              <div v-if="!prestations.length" class="text-center py-8">
                <p class="text-[#94a3b8] text-sm">Aucune prestation publiée.</p>
              </div>
              <div v-else class="divide-y divide-[#f1f5f9]">
                <div v-for="s in prestations.slice(0,5)" :key="s.id" class="py-3 flex items-center justify-between gap-3">
                  <div class="flex-1 min-w-0">
                    <p class="font-semibold text-[#001d32] text-sm truncate">{{ s.title }}</p>
                    <p class="text-[#40617f] text-xs">{{ s.category?.name || 'Sans catégorie' }}</p>
                  </div>
                  <div class="text-right shrink-0">
                    <p class="text-[#006d35] font-semibold text-sm">{{ s.price ? formatEUR(s.price) : 'Sur devis' }}</p>
                    <span :class="prestationBadge(s.status)" class="text-xs font-semibold px-2 py-0.5 rounded-full">{{ prestationLabel(s.status) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </template>

    <ProOnboardingWizard
      :show="showOnboarding"
      :company-name="profile?.company_name || ''"
      @done="onOnboardingDone"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import {
  BuildingStorefrontIcon,
  WrenchScrewdriverIcon,
  CalendarDaysIcon,
  QrCodeIcon,
  StarIcon,
  MegaphoneIcon,
  ChartBarIcon,
  Cog6ToothIcon,
  CheckBadgeIcon,
  ExclamationTriangleIcon,
} from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'
import { userApi } from '@/services/publicApi'
import ProOnboardingWizard from '@/components/ProOnboardingWizard.vue'

const router   = useRouter()
const userAuth = useUserAuthStore()

const loading      = ref(true)
const profile      = ref(null)
const realProjects = ref([])
const prestations  = ref([])
const projectCount   = ref(0)
const prestationCount = ref(0)
const editing   = ref(false)
const saving    = ref(false)
const saveError = ref('')
const showOnboarding = ref(false)

function onOnboardingDone() {
  showOnboarding.value = false
  loadData()
}

const editForm = ref({ company_name: '', description: '', website: '' })

function statusLabel(s) {
  return { pending: 'En attente', approved: 'Approuvé', suspended: 'Suspendu', rejected: 'Rejeté' }[s] || s
}
function formatEUR(price) {
  const n = parseFloat(price)
  return isNaN(n) ? price : new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(n)
}
function prestationLabel(s) {
  return { draft: 'Brouillon', published: 'Publié', archived: 'Archivé' }[s] || s
}
function prestationBadge(s) {
  if (s === 'published') return 'bg-green-100 text-green-800'
  if (s === 'draft') return 'bg-yellow-100 text-yellow-800'
  return 'bg-gray-100 text-gray-600'
}

async function loadData() {
  loading.value = true
  try {
    const [profileData, projectsData, prestationsData] = await Promise.all([
      userApi('/provider/profile').catch(() => null),
      userApi('/projects').catch(() => ({ data: [] })),
      userApi('/provider/prestations').catch(() => ({ data: [] })),
    ])
    profile.value      = profileData
    realProjects.value = projectsData.data || []
    prestations.value  = prestationsData.data || []
    projectCount.value   = realProjects.value.length
    prestationCount.value = prestations.value.length

    if (profile.value) {
      editForm.value = {
        company_name: profile.value.company_name || '',
        description:  profile.value.description  || '',
        website:      profile.value.website       || '',
      }

      if (profile.value.status === 'approved' && profile.value.is_onboarded === false) {
        showOnboarding.value = true
      }
    }
  } catch {}
  loading.value = false
}

async function saveProfile() {
  saving.value   = true
  saveError.value = ''
  try {
    const updated = await userApi('/provider/profile', {
      method: 'PUT',
      body: JSON.stringify({
        company_name: editForm.value.company_name,
        description:  editForm.value.description  || null,
        website:      editForm.value.website       || null,
      }),
    })
    profile.value = updated
    editing.value = false
  } catch (e) {
    saveError.value = e.message || 'Erreur lors de l\'enregistrement.'
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (!userAuth.isLoggedIn) {
    router.replace('/connexion?redirect=/profil/pro')
    return
  }
  await loadData()
})
</script>
