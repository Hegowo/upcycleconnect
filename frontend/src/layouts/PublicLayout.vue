<template>
  <div class="min-h-screen flex flex-col" :style="{ backgroundColor: route.meta.hideNav ? '#f7f9ff' : '#f8fafc' }">

    <header v-if="!route.meta.hideNav" class="sticky top-0 z-50 bg-white/80 backdrop-blur-md border-b border-[#e2e8f0]">
      <div class="max-w-[1280px] mx-auto px-6 py-4 flex items-center justify-between">

        <RouterLink to="/" class="flex items-center">
          <img src="/logoentier.png" :alt="t('public.layout.logoAlt')" class="h-12 w-auto" />
        </RouterLink>

        <nav class="hidden md:flex items-center gap-8">
          <RouterLink
            v-for="link in navLinks"
            :key="link.path"
            :to="link.path"
            class="text-[14px] font-semibold font-jakarta tracking-[-0.025em] pb-1.5 transition-colors"
            :class="isActive(link.path)
              ? 'text-[#006d35] border-b-2 border-[#006d35]'
              : 'text-[#475569] hover:text-[#006d35]'"
          >
            {{ link.label }}
          </RouterLink>
        </nav>

        <div class="flex items-center gap-2">
          <div class="relative hidden lg:block">
            <input
              type="text"
              :placeholder="t('public.layout.searchPlaceholder')"
              class="bg-[#edf4ff] pl-10 pr-4 py-2.5 rounded-xl text-sm text-gray-500 w-56 outline-none focus:ring-2 focus:ring-[#006d35]/20"
            />
            <MagnifyingGlassIcon class="w-[18px] h-[18px] absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          </div>

          <RouterLink
            to="/depot"
            class="hidden md:flex items-center gap-1.5 text-sm font-bold font-jakarta px-4 py-2 rounded-full text-[#006d35] border-2 border-[#006d35]/15 hover:border-[#006d35]/40 hover:bg-[#006d35]/5 transition"
          >
            <PlusCircleIcon class="w-4 h-4" />
            {{ t('public.layout.btnDeposit') }}
          </RouterLink>

          <template v-if="!userAuth.isLoggedIn">
            <RouterLink
              to="/connexion"
              class="hidden sm:block text-sm font-medium px-4 py-2 rounded-lg text-[#40617f] hover:bg-[#edf4ff] transition"
            >
              {{ t('public.layout.btnLogin') }}
            </RouterLink>
            <RouterLink
              to="/inscription"
              class="text-sm font-bold font-jakarta px-5 py-2.5 rounded-full text-white transition hover:opacity-90"
              style="background: #006d35;"
            >
              {{ t('public.layout.btnRegister') }}
            </RouterLink>
          </template>

          <template v-else>
            <button class="hidden sm:flex w-9 h-9 rounded-full items-center justify-center text-[#40617f] hover:bg-[#edf4ff] transition">
              <BellIcon class="w-5 h-5" />
            </button>

            <div class="relative">
              <button
                @click="profileMenuOpen = !profileMenuOpen"
                class="w-9 h-9 rounded-full flex items-center justify-center text-sm font-bold text-white transition hover:opacity-90 focus:outline-none"
                style="background: linear-gradient(135deg, #006d35, #1b8848);"
              >
                {{ userAuth.initials || '?' }}
              </button>

              <Transition name="dropdown">
                <div
                  v-if="profileMenuOpen"
                  v-click-outside="closeMenu"
                  class="absolute right-0 top-12 w-56 bg-white rounded-2xl shadow-xl border border-[#e2e8f0] py-2 z-50"
                >
                  <div class="px-4 py-3 border-b border-[#f1f5f9]">
                    <p class="text-sm font-semibold text-[#001d32] truncate">{{ userAuth.fullName }}</p>
                    <p class="text-xs text-gray-400 truncate">{{ userAuth.user?.email }}</p>
                  </div>

                  <RouterLink
                    to="/profil"
                    @click="closeMenu"
                    class="flex items-center gap-3 px-4 py-2.5 text-sm text-[#334155] hover:bg-[#f8fafc] transition"
                  >
                    <UserCircleIcon class="w-4 h-4 text-[#40617f]" />
                    {{ t('public.layout.menuMyProfile') }}
                  </RouterLink>

                  <RouterLink
                    to="/depot"
                    @click="closeMenu"
                    class="flex items-center gap-3 px-4 py-2.5 text-sm text-[#334155] hover:bg-[#f8fafc] transition"
                  >
                    <PlusCircleIcon class="w-4 h-4 text-[#006d35]" />
                    {{ t('public.layout.menuDeposit') }}
                  </RouterLink>

                  <RouterLink
                    to="/profil/factures"
                    @click="closeMenu"
                    class="flex items-center gap-3 px-4 py-2.5 text-sm text-[#334155] hover:bg-[#f8fafc] transition"
                  >
                    <DocumentTextIcon class="w-4 h-4 text-[#40617f]" />
                    {{ t('public.layout.menuMyInvoices') }}
                  </RouterLink>

                  <div class="border-t border-[#f1f5f9] mt-1 pt-1">
                    <button
                      @click="handleLogout"
                      class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-500 hover:bg-red-50 transition"
                    >
                      <svg class="w-4 h-4" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                      </svg>
                      {{ t('public.layout.menuLogout') }}
                    </button>
                  </div>
                </div>
              </Transition>
            </div>
          </template>
        </div>

        <button
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="md:hidden p-2 rounded-lg text-[#40617f] hover:bg-[#edf4ff] transition"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <Transition name="dropdown">
        <div v-if="mobileMenuOpen" class="md:hidden border-t border-[#e2e8f0] bg-white px-6 py-4 space-y-1">
          <RouterLink
            v-for="link in navLinks"
            :key="link.path"
            :to="link.path"
            @click="mobileMenuOpen = false"
            class="block py-2.5 text-[15px] font-semibold font-jakarta transition-colors"
            :class="isActive(link.path) ? 'text-[#006d35]' : 'text-[#475569]'"
          >
            {{ link.label }}
          </RouterLink>
          <div class="pt-3 border-t border-[#e2e8f0]">
            <RouterLink to="/depot" @click="mobileMenuOpen = false" class="block py-2 text-sm font-bold text-[#006d35]">+ {{ t('public.layout.btnDeposit') }}</RouterLink>
            <template v-if="!userAuth.isLoggedIn">
              <RouterLink to="/connexion" @click="mobileMenuOpen = false" class="block py-2 text-sm font-medium text-[#40617f]">{{ t('public.layout.btnLogin') }}</RouterLink>
              <RouterLink to="/inscription" @click="mobileMenuOpen = false" class="block py-2 text-sm font-bold text-[#006d35]">{{ t('public.layout.btnRegister') }}</RouterLink>
            </template>
            <template v-else>
              <RouterLink to="/profil" @click="mobileMenuOpen = false" class="block py-2 text-sm font-medium text-[#334155]">{{ t('public.layout.menuMyProfile') }}</RouterLink>
              <RouterLink to="/profil/factures" @click="mobileMenuOpen = false" class="block py-2 text-sm font-medium text-[#334155]">{{ t('public.layout.menuMyInvoices') }}</RouterLink>
              <button @click="handleLogout(); mobileMenuOpen = false" class="block py-2 text-sm font-medium text-red-500">{{ t('public.layout.menuLogout') }}</button>
            </template>
          </div>
        </div>
      </Transition>
    </header>

    <main class="flex-1">
      <RouterView />
    </main>

    <footer v-if="!route.meta.hideNav" class="bg-white border-t border-[#e2e8f0]">
      <div class="max-w-[1280px] mx-auto px-8 pt-16 pb-8">

        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
          <div class="space-y-4">
            <RouterLink to="/" class="block">
              <img src="/logoentier.png" :alt="t('public.layout.logoAlt')" class="h-20 w-auto" />
            </RouterLink>
            <p class="text-[#334155] text-sm leading-relaxed">
              {{ t('public.layout.footerTagline') }}
            </p>
            <div class="flex items-center gap-4">
              <a href="#" class="text-[#334155] hover:text-[#006d35] transition" aria-label="Facebook">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/></svg>
              </a>
              <a href="#" class="text-[#334155] hover:text-[#006d35] transition" aria-label="Instagram">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zm0 5.838c-3.403 0-6.162 2.759-6.162 6.162s2.759 6.163 6.162 6.163 6.162-2.759 6.162-6.163c0-3.403-2.759-6.162-6.162-6.162zm0 10.162c-2.209 0-4-1.79-4-4 0-2.209 1.791-4 4-4s4 1.791 4 4c0 2.21-1.791 4-4 4zm6.406-11.845c-.796 0-1.441.645-1.441 1.44s.645 1.44 1.441 1.44c.795 0 1.439-.645 1.439-1.44s-.644-1.44-1.439-1.44z"/></svg>
              </a>
              <a href="#" class="text-[#334155] hover:text-[#006d35] transition" aria-label="LinkedIn">
                <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/></svg>
              </a>
            </div>
          </div>

          <div class="space-y-4">
            <h4 class="text-[#0f172a] text-xs font-semibold uppercase tracking-[0.7px]">{{ t('public.layout.footerColPlatform') }}</h4>
            <ul class="space-y-3">
              <li v-for="link in footerPlateforme" :key="link.label">
                <RouterLink :to="link.path" class="text-[#334155] text-base hover:text-[#006d35] transition">
                  {{ link.label }}
                </RouterLink>
              </li>
            </ul>
          </div>

          <div class="space-y-4">
            <h4 class="text-[#0f172a] text-xs font-semibold uppercase tracking-[0.7px]">{{ t('public.layout.footerColCommunity') }}</h4>
            <ul class="space-y-3">
              <li v-for="link in footerCommunaute" :key="link.label">
                <RouterLink :to="link.path" class="text-[#334155] text-base hover:text-[#006d35] transition">
                  {{ link.label }}
                </RouterLink>
              </li>
            </ul>
          </div>

          <div class="space-y-4">
            <h4 class="text-[#0f172a] text-xs font-semibold uppercase tracking-[0.7px]">{{ t('public.layout.footerColLegal') }}</h4>
            <ul class="space-y-3">
              <li v-for="link in footerLegal" :key="link.label">
                <RouterLink :to="link.path" class="text-[#334155] text-base hover:text-[#006d35] transition">
                  {{ link.label }}
                </RouterLink>
              </li>
            </ul>
          </div>
        </div>

        <div class="border-t border-[#e2e8f0] pt-8 flex flex-col sm:flex-row items-center justify-between gap-3">
          <p class="text-[#334155] text-base">{{ t('public.layout.footerCopyright', { year: new Date().getFullYear() }) }}</p>
          <p class="text-[#334155] text-sm flex items-center gap-1">
            {{ t('public.layout.footerMadeWith') }} <span class="text-red-500 mx-0.5">♥</span> {{ t('public.layout.footerForPlanet') }}
          </p>
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup>import { ref, onMounted, computed } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { MagnifyingGlassIcon, BellIcon, UserCircleIcon, DocumentTextIcon, PlusCircleIcon } from '@heroicons/vue/24/outline'
import { useUserAuthStore } from '@/stores/userAuth'

const { t } = useI18n()
const route    = useRoute()
const router   = useRouter()
const userAuth = useUserAuthStore()

const profileMenuOpen = ref(false)
const mobileMenuOpen = ref(false)

onMounted(() => userAuth.init())

function closeMenu() { profileMenuOpen.value = false }

async function handleLogout() {
  profileMenuOpen.value = false
  await userAuth.logout()
  router.push('/')
}

const navLinks = computed(() => [
  { label: t('public.layout.navHome'),         path: '/' },
  { label: t('public.layout.navPrestations'),  path: '/prestations' },
  { label: t('public.layout.navEvents'),       path: '/evenements' },
  { label: t('public.layout.navCommunity'),    path: '/communaute' },
])

function isActive(path) {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

const footerPlateforme = computed(() => [
  { label: t('public.layout.footerLinkDeposit'),        path: '/depot' },
  { label: t('public.layout.footerLinkFindPrestation'), path: '/prestations' },
  { label: t('public.layout.footerLinkEvents'),         path: '/evenements' },
  { label: t('public.layout.footerLinkBecomeProvider'), path: '/inscription' },
  { label: t('public.layout.footerLinkHowItWorks'),     path: '/' },
])

const footerCommunaute = computed(() => [
  { label: t('public.layout.footerLinkBlog'),         path: '/' },
  { label: t('public.layout.footerLinkTestimonials'), path: '/' },
  { label: t('public.layout.footerLinkResources'),    path: '/' },
])

const footerLegal = computed(() => [
  { label: t('public.layout.footerLinkLegalNotice'), path: '/' },
  { label: t('public.layout.footerLinkTerms'),       path: '/' },
  { label: t('public.layout.footerLinkPrivacy'),     path: '/' },
  { label: t('public.layout.footerLinkContact'),     path: '/' },
])

const vClickOutside = {
  mounted(el, binding) {
    el._clickOutside = (e) => { if (!el.contains(e.target)) binding.value(e) }
    document.addEventListener('mousedown', el._clickOutside)
  },
  unmounted(el) {
    document.removeEventListener('mousedown', el._clickOutside)
  },
}
</script>

<style scoped>
.dropdown-enter-active { transition: opacity 0.15s ease, transform 0.15s ease; }
.dropdown-leave-active { transition: opacity 0.1s ease, transform 0.1s ease; }
.dropdown-enter-from, .dropdown-leave-to { opacity: 0; transform: translateY(-6px) scale(0.97); }
</style>
