import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/PublicLayout.vue'),
      children: [
        { path: '',            component: () => import('@/pages/PublicHomePage.vue'),          meta: { title: 'Accueil' } },
        { path: 'connexion',    component: () => import('@/pages/PublicLoginPage.vue'),          meta: { title: 'Connexion' } },
        { path: 'connexion/verification-requise', component: () => import('@/pages/PublicLoginPendingPage.vue'), meta: { title: 'Vérification requise' } },
        { path: 'inscription',  component: () => import('@/pages/PublicRegisterPage.vue'),       meta: { title: 'Inscription', hideNav: true } },
        { path: 'verify-email', component: () => import('@/pages/PublicVerifyEmailPage.vue'),    meta: { title: 'Activation du compte' } },
        { path: 'verify-login', component: () => import('@/pages/PublicVerifyLoginPage.vue'),    meta: { title: 'Confirmation de connexion' } },
        { path: 'prestations',       component: () => import('@/pages/PublicPrestationsPage.vue'),     meta: { title: 'Prestations' } },
        { path: 'prestations/:id',   component: () => import('@/pages/PublicPrestationDetailPage.vue'), meta: { title: 'Détail Prestation' } },
        { path: 'paiement/succes',   component: () => import('@/pages/PublicPaymentSuccessPage.vue'),   meta: { title: 'Paiement confirmé' } },
        { path: 'paiement/annule',   component: () => import('@/pages/PublicPaymentCancelPage.vue'),    meta: { title: 'Paiement annulé' } },
        { path: 'profil/factures',   component: () => import('@/pages/PublicInvoicesPage.vue'),         meta: { title: 'Mes factures' } },
        { path: 'profil/reservations/:id', component: () => import('@/pages/PublicReservationDetailPage.vue'), meta: { title: 'Détail Réservation' } },
        { path: 'evenements',       component: () => import('@/pages/PublicEvenementsPage.vue'),      meta: { title: 'Événements' } },
        { path: 'evenements/:id',   component: () => import('@/pages/PublicEvenementDetailPage.vue'), meta: { title: 'Détail Événement' } },
        { path: 'communaute',       component: () => import('@/pages/PublicCommunautePage.vue'),      meta: { title: 'Communauté' } },
        { path: 'profil',               component: () => import('@/pages/PublicProfilPage.vue'),                  meta: { title: 'Mon Profil' } },
        { path: 'profil/modifier',      component: () => import('@/pages/PublicProfilEditPage.vue'),              meta: { title: 'Modifier mon profil' } },
        { path: 'profil/pro',           component: () => import('@/pages/PublicProfilProfessionnelPage.vue'),  meta: { title: 'Profil Professionnel' } },
        { path: 'depot',                component: () => import('@/pages/PublicDepotPage.vue'),                meta: { title: 'Déposer un objet' } },
        { path: 'confirmation-depot',   component: () => import('@/pages/PublicConfirmationDepotPage.vue'),    meta: { title: 'Confirmation de Dépôt' } },
        { path: 'communaute/categories-forum', component: () => import('@/pages/PublicCategoriesForumPage.vue'), meta: { title: 'Catégories Forum' } },
      ],
    },

    {
      path: '/admin/login',
      component: () => import('@/pages/AdminLoginPage.vue'),
      meta: { layout: 'auth', title: 'Administration' },
    },

    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '',        redirect: 'dashboard' },
        { path: 'dashboard',   component: () => import('@/pages/DashboardPage.vue'),    meta: { title: 'Tableau de bord' } },
        { path: 'users',       component: () => import('@/pages/UserListPage.vue'),          meta: { title: 'Utilisateurs' } },
        { path: 'users/:id',   component: () => import('@/pages/UserDetailPage.vue'),        meta: { title: 'Fiche utilisateur' } },
        { path: 'providers',   component: () => import('@/pages/ProviderListPage.vue'),  meta: { title: 'Prestataires' } },
        { path: 'providers/:id', component: () => import('@/pages/ProviderDetailPage.vue'), meta: { title: 'Fiche prestataire' } },
        { path: 'admins',        component: () => import('@/pages/AdminListPage.vue'),      meta: { title: 'Administrateurs', requiresRole: 'super_admin' } },
        { path: 'admins/create', component: () => import('@/pages/AdminFormPage.vue'),      meta: { title: 'Nouvel administrateur', requiresRole: 'super_admin' } },
        { path: 'admins/:id/edit', component: () => import('@/pages/AdminFormPage.vue'),    meta: { title: 'Modifier l\'administrateur', requiresRole: 'super_admin' } },
        { path: 'categories',        component: () => import('@/pages/CategoryListPage.vue'),    meta: { title: 'Catégories' } },
        { path: 'categories/create', component: () => import('@/pages/CategoryFormPage.vue'),    meta: { title: 'Nouvelle catégorie' } },
        { path: 'categories/:id/edit', component: () => import('@/pages/CategoryFormPage.vue'), meta: { title: 'Modifier la catégorie' } },
        { path: 'prestations',         component: () => import('@/pages/PrestationListPage.vue'),  meta: { title: 'Prestations' } },
        { path: 'prestations/create',  component: () => import('@/pages/PrestationFormPage.vue'),  meta: { title: 'Nouvelle prestation' } },
        { path: 'prestations/:id/edit', component: () => import('@/pages/PrestationFormPage.vue'), meta: { title: 'Modifier la prestation' } },
        { path: 'events',        component: () => import('@/pages/EventListPage.vue'),  meta: { title: 'Événements' } },
        { path: 'events/create', component: () => import('@/pages/EventFormPage.vue'),  meta: { title: 'Nouvel événement' } },
        { path: 'events/:id/edit', component: () => import('@/pages/EventFormPage.vue'), meta: { title: 'Modifier l\'événement' } },
        {
          path: 'logs',
          component: () => import('@/pages/LogsPage.vue'),
          meta: { title: 'Journaux d\'activité' },
        },
        {
          path: 'database',
          component: () => import('@/pages/DatabasePage.vue'),
          meta: { title: 'Base de données', requiresRole: 'super_admin', fullscreen: true },
        },
        { path: 'settings', component: () => import('@/pages/SettingsPage.vue'), meta: { title: 'Paramètres' } },
        { path: 'box-requests', component: () => import('@/pages/BoxRequestsPage.vue'), meta: { title: 'Demandes de Dépôt' } },
        { path: 'docs', component: () => import('@/pages/SwaggerPage.vue'), meta: { title: 'Documentation API', fullscreen: true } },
      ],
    },

    { path: '/:pathMatch(.*)*', redirect: '/' },
  ],
})

let authInitialized = false

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()

  if (!authInitialized) {
    authInitialized = true
    await auth.initAuth()
  }

  if (to.path === '/admin/login' && auth.isAuthenticated) {
    return next('/admin/dashboard')
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next('/admin/login')
  }

  if (to.meta.requiresRole === 'super_admin' && !auth.isSuperAdmin) {
    return next('/admin/dashboard')
  }

  next()
})

router.afterEach((to) => {
  const base = 'UpcycleConnect'
  const title = to.meta?.title
  document.title = title ? `${title} — ${base}` : base
})

export default router
