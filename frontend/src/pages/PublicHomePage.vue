<template>
  <div class="bg-[#f8fafc]">

    <section class="relative overflow-hidden px-4 sm:px-8 pt-20 sm:pt-32 pb-16 sm:pb-24">

      <div class="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] rounded-full bg-[rgba(28,130,65,0.05)] blur-[32px] pointer-events-none" />
      <div class="absolute right-[-96px] top-[-96px] w-[400px] h-[400px] rounded-full bg-[rgba(25,86,135,0.05)] blur-[32px] pointer-events-none" />

      <div class="max-w-[1280px] mx-auto relative">
        <div class="flex flex-col items-center gap-6 max-w-[896px] mx-auto text-center">

          <h1 class="font-jakarta font-bold text-[#0f172a] tracking-[-0.025em] text-4xl sm:text-5xl md:text-[72px] leading-tight md:leading-[72px]">
            {{ t('public.home.heroTitlePart1') }}
            <span class="text-[#1c8241]"> {{ t('public.home.heroTitleHighlight') }}</span>
            <br />{{ t('public.home.heroTitlePart2') }}
          </h1>

          <p class="text-[#334155] text-xl md:text-2xl leading-[32px] max-w-[672px]">
            {{ t('public.home.heroSubtitle') }}
          </p>

          <div class="flex flex-col sm:flex-row gap-6 items-center pt-4">
            <RouterLink
              to="/prestations"
              class="inline-flex items-center justify-center px-8 py-[17px] rounded-full bg-[#1c8241] text-white text-base font-medium transition hover:opacity-90 shadow-lg"
            >
              {{ t('public.home.ctaBrowse') }}
            </RouterLink>
            <RouterLink
              to="/inscription"
              class="inline-flex items-center justify-center px-8 py-[17px] rounded-full border-2 border-[#1c8241] text-[#1c8241] text-base font-medium bg-transparent transition hover:bg-[#1c8241]/5 shadow-sm"
            >
              {{ t('public.home.ctaBecomeProvider') }}
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <section class="bg-white border-t border-b border-[#e2e8f0] py-16 sm:py-24 px-4 sm:px-8">
      <div class="max-w-[1280px] mx-auto flex flex-col gap-12">

        <div class="flex flex-col items-center gap-4 text-center">
          <h2 class="font-jakarta font-bold text-[#0f172a] text-[36px] leading-[40px]">
            {{ t('public.home.categoriesTitle') }}
          </h2>
          <p class="text-[#334155] text-lg leading-[28px] max-w-[672px]">
            {{ t('public.home.categoriesSubtitle') }}
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6">

          <RouterLink
            to="/prestations"
            class="sm:col-span-2 md:col-span-2 md:row-span-2 min-h-[280px] relative flex flex-col justify-between p-8 rounded-[32px] overflow-hidden border border-[#bbf7d0] cursor-pointer hover:shadow-lg transition-shadow"
            style="background: linear-gradient(143deg, #f0fdf4 0%, #d1fae5 100%);"
          >
            <div class="absolute top-0 right-0 w-40 h-36 opacity-20">
              <HomeModernIcon class="w-full h-full text-[#006d35]" />
            </div>

            <div class="relative z-10">
              <div class="w-12 h-12 bg-white rounded-full shadow-sm flex items-center justify-center mb-12">
                <HomeModernIcon class="w-5 h-5 text-[#006d35]" />
              </div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-2xl leading-[32px]">
                {{ t('public.home.catFurnitureTitle') }}
              </h3>
              <p class="text-[#334155] text-base leading-[24px] mt-4 max-w-xs">
                {{ t('public.home.catFurnitureDesc') }}
              </p>
            </div>

            <div class="flex items-center gap-1 text-[#1c8241] font-medium text-base relative z-10">
              {{ t('public.home.catFurnitureCta') }}
              <ChevronRightIcon class="w-2.5 h-2.5" />
            </div>
          </RouterLink>

          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col gap-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="w-10 h-10 bg-[#eff6ff] rounded-full flex items-center justify-center">
              <ScissorsIcon class="w-4 h-4 text-[#3b82f6]" />
            </div>
            <div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">{{ t('public.home.catTextileTitle') }}</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">{{ t('public.home.catTextileDesc') }}</p>
            </div>
          </RouterLink>

          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col gap-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="w-10 h-10 bg-[#fff7ed] rounded-full flex items-center justify-center">
              <ArchiveBoxIcon class="w-4 h-4 text-[#f97316]" />
            </div>
            <div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">{{ t('public.home.catPackagingTitle') }}</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">{{ t('public.home.catPackagingDesc') }}</p>
            </div>
          </RouterLink>

          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col gap-4 shadow-sm hover:shadow-md transition-shadow">
            <div class="w-10 h-10 bg-[#faf5ff] rounded-full flex items-center justify-center">
              <ComputerDesktopIcon class="w-4 h-4 text-[#a855f7]" />
            </div>
            <div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">{{ t('public.home.catElectronicsTitle') }}</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">{{ t('public.home.catElectronicsDesc') }}</p>
            </div>
          </RouterLink>

          <RouterLink to="/prestations" class="bg-white border border-[#e2e8f0] rounded-[32px] p-6 flex flex-col justify-between shadow-sm hover:shadow-md transition-shadow">
            <div>
              <div class="w-10 h-10 bg-[#f3f4f6] rounded-full flex items-center justify-center mb-4">
                <EllipsisHorizontalIcon class="w-4 h-4 text-[#6b7280]" />
              </div>
              <h3 class="font-jakarta font-bold text-[#0f172a] text-xl">{{ t('public.home.catOtherTitle') }}</h3>
              <p class="text-[#334155] text-sm mt-1 leading-[20px]">{{ t('public.home.catOtherDesc') }}</p>
            </div>
            <div class="flex items-center gap-1 text-[#4b5563] font-medium text-sm">
              {{ t('public.home.catSeeAll') }} <ChevronRightIcon class="w-2 h-2" />
            </div>
          </RouterLink>

        </div>
      </div>
    </section>

    <section class="py-12 sm:py-20 px-4 sm:px-8">
      <div class="max-w-[1280px] mx-auto">
        <div
          class="relative overflow-hidden rounded-[32px] px-6 sm:px-16 py-12 sm:py-20 flex flex-col md:flex-row items-center justify-between gap-8 md:gap-10 shadow-xl text-center md:text-left"
          style="background: #195687;"
        >
          <div class="flex flex-col gap-4 max-w-lg">
            <h2 class="font-jakarta font-bold text-white text-[36px] leading-[40px]">
              {{ t('public.home.providerCtaTitle') }}
            </h2>
            <p class="text-[#dbeafe] text-lg leading-[28px]">
              {{ t('public.home.providerCtaText') }}
            </p>
          </div>
          <RouterLink
            to="/inscription"
            class="shrink-0 inline-flex items-center justify-center bg-white font-jakarta font-bold text-[#195687] text-base px-8 py-[17px] rounded-full shadow-lg transition hover:opacity-90"
          >
            {{ t('public.home.providerCtaButton') }}
          </RouterLink>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup>import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  HomeModernIcon,
  ScissorsIcon,
  ArchiveBoxIcon,
  ComputerDesktopIcon,
  EllipsisHorizontalIcon,
  ChevronRightIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
</script>
