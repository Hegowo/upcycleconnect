<template>
  <div class="bg-[#f7f9ff] min-h-screen flex">

    <aside class="w-64 shrink-0 bg-[#f8fafc] border-r border-[#e2e8f0] flex flex-col py-6 px-4 sticky top-[65px] h-[calc(100vh-65px)] overflow-y-auto">
      <div class="px-4 py-6">
        <p class="text-[#15803d] font-semibold text-sm">{{ t('public.forumCategories.navTitle') }}</p>
        <p class="text-[#64748b] text-xs mt-0.5">{{ t('public.forumCategories.navSubtitle') }}</p>
      </div>

      <nav class="flex flex-col gap-1">
        <RouterLink
          v-for="item in navItems"
          :key="item.key"
          :to="item.to"
          class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-semibold transition"
          :class="item.active
            ? 'bg-[#f0fdf4] text-[#15803d]'
            : 'text-[#64748b] hover:bg-[#f0fdf4] hover:text-[#15803d]'"
        >
          <component :is="item.icon" class="w-5 h-5 shrink-0" />
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="mt-8 px-4">
        <button
          class="w-full py-3 rounded-xl text-white font-semibold text-base text-center transition hover:opacity-90"
          style="background: linear-gradient(165deg, #006d35 0%, #1b8848 100%);"
        >
          {{ t('public.forumCategories.newThread') }}
        </button>
      </div>

      <div class="mt-auto flex flex-col gap-1 pb-4">
        <button class="flex items-center gap-3 px-4 py-3 rounded-xl text-[#64748b] text-sm font-semibold hover:bg-[#f0fdf4] hover:text-[#15803d] transition">
          <Cog6ToothIcon class="w-5 h-5 shrink-0" />
          {{ t('public.forumCategories.settings') }}
        </button>
        <button class="flex items-center gap-3 px-4 py-3 rounded-xl text-[#64748b] text-sm font-semibold hover:bg-[#f0fdf4] hover:text-[#15803d] transition">
          <QuestionMarkCircleIcon class="w-5 h-5 shrink-0" />
          {{ t('public.forumCategories.help') }}
        </button>
      </div>
    </aside>

    <main class="flex-1 min-w-0 p-8">
      <div class="max-w-[1152px] flex flex-col gap-12">

        <div class="flex flex-col gap-2">
          <p class="text-[#006d35] text-xs font-semibold uppercase tracking-[1.2px]">{{ t('public.forumCategories.headerTag') }}</p>
          <h1 class="font-jakarta font-extrabold text-[#001d32] text-5xl tracking-tight">{{ t('public.forumCategories.headerTitle') }}</h1>
          <p class="text-[#40617f] text-lg mt-2 max-w-2xl leading-relaxed">
            {{ t('public.forumCategories.headerSubtitle') }}
          </p>
        </div>

        <div class="grid grid-cols-12 gap-6">

          <div class="col-span-7 bg-white rounded-[24px] p-10 relative flex flex-col justify-between min-h-[378px]">
            <div class="flex items-start justify-between">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
                <CubeIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <span class="bg-[#edf4ff] text-[#40617f] text-sm font-semibold px-3 py-1 rounded-full">{{ t('public.forumCategories.topicsCount') }}</span>
            </div>
            <div class="mt-12">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.forumCategories.cat1Title') }}</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3 max-w-lg">
                {{ t('public.forumCategories.cat1Desc') }}
              </p>
              <div class="flex flex-wrap gap-2 mt-4">
                <span v-for="tag in tags" :key="tag"
                  class="bg-[#e3efff] text-[#001d32] text-xs font-semibold px-3 py-1 rounded-lg">
                  {{ tag }}
                </span>
              </div>
            </div>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
              {{ t('public.forumCategories.explore') }} <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <div class="col-span-5 bg-white rounded-[24px] p-10 relative flex flex-col justify-between min-h-[378px]">
            <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
              <WrenchScrewdriverIcon class="w-7 h-7 text-[#40617f]" />
            </div>
            <div class="mt-12">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.forumCategories.cat2Title') }}</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3">
                {{ t('public.forumCategories.cat2Desc') }}
              </p>
              <p class="text-[#40617f] text-sm font-medium mt-4">{{ t('public.forumCategories.cat2Active') }}</p>
            </div>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
              {{ t('public.forumCategories.explore') }} <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <div class="col-span-4 bg-white rounded-[24px] p-10 relative flex flex-col justify-between min-h-[374px] overflow-hidden">

            <div class="absolute -right-4 -top-4 w-24 h-32 opacity-5 pointer-events-none">
              <div class="w-full h-full bg-[#006d35] rounded-xl" />
            </div>
            <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
              <LightBulbIcon class="w-7 h-7 text-[#40617f]" />
            </div>
            <div class="mt-12">
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.forumCategories.cat3Title1') }}<br>{{ t('public.forumCategories.cat3Title2') }}</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3">
                {{ t('public.forumCategories.cat3Desc') }}
              </p>
            </div>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
              {{ t('public.forumCategories.explore') }} <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <div class="col-span-8 bg-white rounded-[24px] p-10 flex gap-8 items-center min-h-[374px]">
            <div class="flex-1">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center">
                <GlobeAltIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl mt-12">{{ t('public.forumCategories.cat4Title1') }}<br>{{ t('public.forumCategories.cat4Title2') }}</h2>
              <p class="text-[#3f4a3f] text-base leading-relaxed mt-3">
                {{ t('public.forumCategories.cat4Desc') }}
              </p>
              <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base mt-6">
                {{ t('public.forumCategories.explore') }} <ChevronRightIcon class="w-3.5 h-3.5" />
              </button>
            </div>

            <div class="bg-[#edf4ff] rounded-xl w-64 h-64 flex items-center justify-center shrink-0 relative">

              <svg class="absolute inset-0 w-full h-full -rotate-90" viewBox="0 0 256 256">
                <circle cx="128" cy="128" r="100" fill="none" stroke="#d8eaff" stroke-width="12" />
                <circle cx="128" cy="128" r="100" fill="none" stroke="#006d35" stroke-width="12"
                  stroke-dasharray="628" stroke-dashoffset="38" stroke-linecap="round" />
              </svg>
              <div class="relative text-center">
                <p class="font-semibold text-[#006d35] text-4xl">{{ t('public.forumCategories.cat4Score') }}</p>
                <p class="text-[#40617f] text-xs uppercase tracking-[1.2px] font-semibold mt-1">{{ t('public.forumCategories.cat4ScoreLabel') }}</p>
              </div>
            </div>
          </div>

          <div class="col-span-6 bg-white rounded-[24px] p-10 flex flex-col gap-6 self-start">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center shrink-0">
                <ChatBubbleLeftEllipsisIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <div>
                <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.forumCategories.cat5Title') }}</h2>
                <p class="text-[#40617f] text-sm mt-0.5">{{ t('public.forumCategories.cat5Subtitle') }}</p>
              </div>
            </div>
            <p class="text-[#3f4a3f] text-base leading-relaxed">
              {{ t('public.forumCategories.cat5Desc') }}
            </p>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base pt-2">
              {{ t('public.forumCategories.explore') }} <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <div class="col-span-6 bg-white rounded-[24px] p-10 flex flex-col gap-6 self-start">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 rounded-xl bg-[#edf4ff] flex items-center justify-center shrink-0">
                <ShoppingBagIcon class="w-7 h-7 text-[#40617f]" />
              </div>
              <div>
                <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.forumCategories.cat6Title') }}</h2>
                <p class="text-[#40617f] text-sm mt-0.5">{{ t('public.forumCategories.cat6Subtitle') }}</p>
              </div>
            </div>
            <p class="text-[#3f4a3f] text-base leading-relaxed">
              {{ t('public.forumCategories.cat6Desc') }}
            </p>
            <button class="flex items-center gap-2 text-[#006d35] font-semibold text-base pt-2">
              {{ t('public.forumCategories.explore') }} <ChevronRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>

          <div class="col-span-12 bg-[#0f172a] rounded-[32px] px-12 py-16 flex items-center justify-between relative overflow-hidden">

            <div class="absolute -right-20 -top-16 w-64 h-64 rounded-full bg-[rgba(0,109,53,0.2)] blur-[32px] pointer-events-none" />
            <div>
              <h3 class="font-jakarta font-bold text-white text-3xl">{{ t('public.forumCategories.ctaTitle') }}</h3>
              <p class="text-[#94a3b8] text-base mt-2">{{ t('public.forumCategories.ctaSubtitle') }}</p>
            </div>
            <button
              class="bg-[#006d35] text-white font-semibold text-base px-8 py-4 rounded-xl hover:bg-[#1b8848] transition shrink-0"
            >
              {{ t('public.forumCategories.ctaButton') }}
            </button>
          </div>

        </div>
      </div>
    </main>

  </div>
</template>

<script setup>import { RouterLink } from 'vue-router'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  ChatBubbleLeftEllipsisIcon,
  LightBulbIcon,
  WrenchScrewdriverIcon,
  CubeIcon,
  GlobeAltIcon,
  ShoppingBagIcon,
  ChevronRightIcon,
  Cog6ToothIcon,
  QuestionMarkCircleIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()

const navItems = computed(() => [
  { key: 'forum',    label: t('public.forumCategories.navForum'),    icon: ChatBubbleLeftEllipsisIcon, to: '/communaute/categories-forum', active: true },
  { key: 'projects', label: t('public.forumCategories.navProjects'), icon: LightBulbIcon,              to: '/communaute',                  active: false },
  { key: 'news',     label: t('public.forumCategories.navNews'),     icon: GlobeAltIcon,               to: '/communaute',                  active: false },
  { key: 'tips',     label: t('public.forumCategories.navTips'),     icon: WrenchScrewdriverIcon,      to: '/communaute',                  active: false },
])

const tags = computed(() => [
  t('public.forumCategories.tagWood'),
  t('public.forumCategories.tagMetal'),
  t('public.forumCategories.tagPlastic'),
  t('public.forumCategories.tagTextile'),
])
</script>
