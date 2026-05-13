<template>
  <div class="bg-[#f7f9ff] pb-0">

    <section class="bg-[#edf4ff] px-6 py-12 overflow-hidden relative">
      <div class="max-w-[1280px] mx-auto flex items-center gap-8">

        <div class="bg-white rounded-[24px] p-4 w-48 h-48 flex items-center justify-center shrink-0 shadow-xl">
          <BuildingStorefrontIcon class="w-24 h-24 text-[#006d35]/40" />
        </div>

        <div class="flex-1 flex flex-col gap-3">
          <div class="flex items-center gap-3">
            <span class="bg-[#006d35] text-white text-xs font-bold uppercase tracking-wide px-3 py-1 rounded-full">
              {{ t('public.profileProfessionnel.badgeCertified') }}
            </span>
            <span class="text-[#40617f] text-xs uppercase tracking-wider">{{ t('public.profileProfessionnel.siretLabel') }}</span>
          </div>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-5xl tracking-tight">{{ t('public.profileProfessionnel.demoName') }}</h1>
          <p class="font-jakarta font-semibold text-[#006d35] text-2xl">{{ t('public.profileProfessionnel.demoActivity') }}</p>
          <div class="flex gap-4 mt-2 flex-wrap">
            <button
              class="flex items-center gap-2 px-8 py-3 rounded-xl text-white font-bold text-sm transition hover:opacity-90"
              style="background: linear-gradient(168deg, #006d35 0%, #1b8848 100%);"
            >
              <EnvelopeIcon class="w-4 h-4" />
              {{ t('public.profileProfessionnel.ctaContact') }}
            </button>
            <button class="flex items-center gap-2 bg-[#cee5ff] text-[#001d32] font-bold text-sm px-8 py-3 rounded-xl hover:bg-[#b8d8ff] transition">
              <ShareIcon class="w-4 h-4" />
              {{ t('public.profileProfessionnel.ctaShare') }}
            </button>
            <button @click="router.push('/profil/pro/evenements')"
              class="flex items-center gap-2 bg-white text-[#006d35] border-2 border-[#006d35] font-bold text-sm px-8 py-3 rounded-xl hover:bg-[#edf4ff] transition">
              <CalendarDaysIcon class="w-4 h-4" />
              Gérer mes événements
            </button>
          </div>
        </div>

        <div class="flex flex-col gap-4 w-60 shrink-0">
          <div class="bg-white rounded-xl p-4 flex items-center gap-4 shadow-sm">
            <div class="w-11 h-11 bg-[#d1fae5] rounded-xl flex items-center justify-center shrink-0">
              <CheckBadgeIcon class="w-6 h-6 text-[#006d35]" />
            </div>
            <div>
              <p class="font-jakarta font-bold text-[#001d32] text-2xl">142</p>
              <p class="text-[#40617f] text-xs uppercase tracking-wide">{{ t('public.profileProfessionnel.statProjects') }}</p>
            </div>
          </div>
          <div class="bg-white rounded-xl p-4 flex items-center gap-4 shadow-sm">
            <div class="w-11 h-11 bg-[#fef3c7] rounded-xl flex items-center justify-center shrink-0">
              <StarIcon class="w-6 h-6 text-[#d97706]" />
            </div>
            <div>
              <p class="font-jakarta font-bold text-[#001d32] text-2xl">4.9/5</p>
              <p class="text-[#40617f] text-xs uppercase tracking-wide">{{ t('public.profileProfessionnel.statRating') }}</p>
            </div>
          </div>
        </div>

        <div class="absolute -right-0 -top-0 w-60 h-60 opacity-20 pointer-events-none">
          <div class="w-full h-full rounded-full bg-gradient-to-br from-[#006d35] to-transparent blur-3xl" />
        </div>
      </div>
    </section>

    <div class="max-w-[1280px] mx-auto px-6 py-10">
      <div class="grid grid-cols-12 gap-6">

        <div class="col-span-4 flex flex-col gap-6">

          <div class="bg-white rounded-[24px] p-8 flex flex-col gap-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.profileProfessionnel.servicesTitle') }}</h2>
            <div class="flex flex-col gap-5">
              <div
                v-for="service in services"
                :key="service.title"
                class="flex gap-4"
              >
                <div class="w-5 h-5 mt-0.5 shrink-0">
                  <component :is="service.icon" class="w-5 h-5 text-[#006d35]" />
                </div>
                <div>
                  <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ service.title }}</p>
                  <p class="text-[#3f4a3f] text-xs leading-relaxed mt-1">{{ service.desc }}</p>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-[24px] p-8 flex flex-col gap-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.profileProfessionnel.certsTitle') }}</h2>
            <div class="grid grid-cols-2 gap-4">
              <div
                v-for="cert in certs"
                :key="cert.label"
                class="bg-[#edf4ff] rounded-xl p-4 flex flex-col items-center gap-2"
                :class="cert.wide ? 'col-span-2' : ''"
              >
                <component :is="cert.icon" class="w-6 h-6 text-[#006d35]" />
                <span class="text-[#001d32] text-xs font-bold uppercase tracking-wide text-center">{{ cert.label }}</span>
              </div>
            </div>
          </div>

        </div>

        <div class="col-span-8 flex flex-col gap-6">

          <div class="bg-white rounded-[24px] p-8 flex flex-col gap-8">
            <div class="flex items-end justify-between">
              <div>
                <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.profileProfessionnel.projectsTitle') }}</h2>
                <p class="text-[#3f4a3f] text-sm mt-1">{{ t('public.profileProfessionnel.projectsSubtitle') }}</p>
              </div>
              <button class="flex items-center gap-1 text-[#006d35] font-bold text-sm hover:underline">
                {{ t('public.profileProfessionnel.seeAll') }} <ChevronRightIcon class="w-3.5 h-3.5" />
              </button>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div v-for="project in projects" :key="project.id" class="flex flex-col">
                <div class="relative h-64 rounded-xl overflow-hidden bg-gradient-to-br"
                  :style="`background: linear-gradient(135deg, ${project.colorFrom}, ${project.colorTo});`"
                >
                  <div class="absolute inset-0 flex items-center justify-center">
                    <component :is="project.icon" class="w-24 h-24 text-white/20" />
                  </div>
                  <div class="absolute top-4 left-4">
                    <span class="bg-white/90 backdrop-blur-sm text-[#006d35] text-xs font-semibold uppercase px-3 py-1 rounded-full">
                      {{ project.badge }}
                    </span>
                  </div>
                  <div class="absolute bottom-4 right-4">
                    <span class="bg-[#006d35] text-white text-xs font-semibold px-3 py-1 rounded-full">
                      {{ project.impact }}
                    </span>
                  </div>
                </div>
                <h3 class="font-jakarta font-bold text-[#001d32] text-lg mt-3">{{ project.title }}</h3>
                <p class="text-[#3f4a3f] text-sm leading-relaxed mt-1">{{ project.desc }}</p>
              </div>
            </div>
          </div>

          <div class="bg-[#edf4ff] rounded-[24px] p-8 flex flex-col gap-8">
            <div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.profileProfessionnel.shopTitle') }}</h2>
              <p class="text-[#3f4a3f] text-sm mt-1">{{ t('public.profileProfessionnel.shopSubtitle') }}</p>
            </div>

            <div class="grid grid-cols-3 gap-6">
              <div
                v-for="product in products"
                :key="product.id"
                class="bg-white rounded-xl overflow-hidden"
              >
                <div class="relative h-40 bg-gradient-to-br"
                  :style="`background: linear-gradient(135deg, ${product.colorFrom}, ${product.colorTo});`"
                >
                  <div class="absolute inset-0 flex items-center justify-center">
                    <component :is="product.icon" class="w-14 h-14 text-white/30" />
                  </div>
                  <button class="absolute top-2 right-2 w-7 h-7 bg-white/90 rounded-full flex items-center justify-center hover:bg-white transition">
                    <HeartIcon class="w-3.5 h-3.5 text-[#40617f]" />
                  </button>
                </div>
                <div class="p-3">
                  <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ product.title }}</p>
                  <div class="flex items-center justify-between mt-2">
                    <span class="text-[#006d35] font-semibold text-sm">{{ product.price }}</span>
                    <button class="bg-[rgba(0,109,53,0.1)] p-2 rounded-lg hover:bg-[rgba(0,109,53,0.2)] transition">
                      <ShoppingCartIcon class="w-4 h-4 text-[#006d35]" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>

  </div>
</template>

<script setup>import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import {
  BuildingStorefrontIcon,
  EnvelopeIcon,
  ShareIcon,
  CheckBadgeIcon,
  StarIcon,
  ChevronRightIcon,
  HomeModernIcon,
  ScissorsIcon,
  LightBulbIcon,
  SparklesIcon,
  HeartIcon,
  ShoppingCartIcon,
  GlobeAltIcon,
  SparklesIcon as LeafIcon,
  TrophyIcon,
  CalendarDaysIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const router = useRouter()

const services = computed(() => [
  {
    icon: HomeModernIcon,
    title: t('public.profileProfessionnel.service1Title'),
    desc:  t('public.profileProfessionnel.service1Desc'),
  },
  {
    icon: ScissorsIcon,
    title: t('public.profileProfessionnel.service2Title'),
    desc:  t('public.profileProfessionnel.service2Desc'),
  },
  {
    icon: LightBulbIcon,
    title: t('public.profileProfessionnel.service3Title'),
    desc:  t('public.profileProfessionnel.service3Desc'),
  },
])

const certs = computed(() => [
  { icon: GlobeAltIcon, label: t('public.profileProfessionnel.certNetwork'),    wide: false },
  { icon: LeafIcon,     label: t('public.profileProfessionnel.certZeroWaste'),  wide: false },
  { icon: TrophyIcon,   label: t('public.profileProfessionnel.certArtisan'),    wide: true  },
])

const projects = computed(() => [
  {
    id: 1,
    title: t('public.profileProfessionnel.project1Title'),
    desc:  t('public.profileProfessionnel.project1Desc'),
    badge: t('public.profileProfessionnel.project1Badge'),
    impact: t('public.profileProfessionnel.project1Impact'),
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: HomeModernIcon,
  },
  {
    id: 2,
    title: t('public.profileProfessionnel.project2Title'),
    desc:  t('public.profileProfessionnel.project2Desc'),
    badge: t('public.profileProfessionnel.project2Badge'),
    impact: t('public.profileProfessionnel.project2Impact'),
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: LightBulbIcon,
  },
])

const products = computed(() => [
  { id: 1, title: t('public.profileProfessionnel.product1Title'), price: '85,00 €',  colorFrom: '#fef3c7', colorTo: '#fde68a', icon: LightBulbIcon },
  { id: 2, title: t('public.profileProfessionnel.product2Title'), price: '210,00 €', colorFrom: '#d1fae5', colorTo: '#a7f3d0', icon: HomeModernIcon },
  { id: 3, title: t('public.profileProfessionnel.product3Title'), price: '45,00 €',  colorFrom: '#f3e8ff', colorTo: '#e9d5ff', icon: SparklesIcon },
])
</script>
