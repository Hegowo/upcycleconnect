import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/admin/dashboard',
    },
    {
      path: '/login',
      component: () => import('@/pages/auth/LoginPage.vue'),
      meta: { layout: 'auth', title: 'Connexion' },
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '',        redirect: 'dashboard' },
        { path: 'dashboard',   component: () => import('@/pages/dashboard/DashboardPage.vue'),    meta: { title: 'Tableau de bord' } },
        { path: 'users',       component: () => import('@/pages/users/UserListPage.vue'),          meta: { title: 'Utilisateurs' } },
        { path: 'users/:id',   component: () => import('@/pages/users/UserDetailPage.vue'),        meta: { title: 'Fiche utilisateur' } },
        { path: 'providers',   component: () => import('@/pages/providers/ProviderListPage.vue'),  meta: { title: 'Prestataires' } },
        { path: 'providers/:id', component: () => import('@/pages/providers/ProviderDetailPage.vue'), meta: { title: 'Fiche prestataire' } },
        {
          path: 'admins',
          component: () => import('@/pages/admins/AdminListPage.vue'),
          meta: { title: 'Administrateurs', requiresRole: 'super_admin' },
        },
        {
          path: 'admins/create',
          component: () => import('@/pages/admins/AdminFormPage.vue'),
          meta: { title: 'Nouvel administrateur', requiresRole: 'super_admin' },
        },
        {
          path: 'admins/:id/edit',
          component: () => import('@/pages/admins/AdminFormPage.vue'),
          meta: { title: 'Modifier l\'administrateur', requiresRole: 'super_admin' },
        },
        { path: 'categories',        component: () => import('@/pages/categories/CategoryListPage.vue'),  meta: { title: 'Catégories' } },
        { path: 'categories/create', component: () => import('@/pages/categories/CategoryFormPage.vue'),  meta: { title: 'Nouvelle catégorie' } },
        { path: 'categories/:id/edit', component: () => import('@/pages/categories/CategoryFormPage.vue'), meta: { title: 'Modifier la catégorie' } },
        { path: 'prestations',        component: () => import('@/pages/prestations/PrestationListPage.vue'), meta: { title: 'Prestations' } },
        { path: 'prestations/create', component: () => import('@/pages/prestations/PrestationFormPage.vue'), meta: { title: 'Nouvelle prestation' } },
        { path: 'prestations/:id/edit', component: () => import('@/pages/prestations/PrestationFormPage.vue'), meta: { title: 'Modifier la prestation' } },
        { path: 'events',        component: () => import('@/pages/events/EventListPage.vue'),  meta: { title: 'Événements' } },
        { path: 'events/create', component: () => import('@/pages/events/EventFormPage.vue'),  meta: { title: 'Nouvel événement' } },
        { path: 'events/:id/edit', component: () => import('@/pages/events/EventFormPage.vue'), meta: { title: 'Modifier l\'événement' } },
        {
          path: 'database',
          component: () => import('@/pages/database/DatabasePage.vue'),
          meta: { title: 'Base de données', requiresRole: 'super_admin', fullscreen: true },
        },
      ],
    },
    { path: '/:pathMatch(.*)*', redirect: '/admin/dashboard' },
  ],
})

let authInitialized = false

router.beforeEach(async (to, from, next) => {
  const auth = useAuthStore()

  if (!authInitialized) {
    authInitialized = true
    await auth.initAuth()
  }

  if (to.path === '/login' && auth.isAuthenticated) {
    return next('/admin/dashboard')
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return next('/login')
  }

  if (to.meta.requiresRole === 'super_admin' && !auth.isSuperAdmin) {
    return next('/admin/dashboard')
  }

  next()
})

export default router
