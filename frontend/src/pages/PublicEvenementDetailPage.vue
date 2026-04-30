<template>
  <div class="bg-[#f7f9ff] pb-16">
    <div class="max-w-[1280px] mx-auto px-6 py-10">

      <div class="flex items-center gap-2 text-sm text-[#40617f] mb-8">
        <RouterLink to="/evenements" class="hover:text-[#006d35] transition">{{ t('public.eventDetail.breadcrumbEvents') }}</RouterLink>
        <ChevronRightIcon class="w-3.5 h-3.5" />
        <span class="text-[#001d32] font-medium">{{ t('public.eventDetail.demoTitle') }}</span>
      </div>

      <div class="grid grid-cols-12 gap-8">

        <div class="col-span-8 flex flex-col gap-8">

          <div class="relative h-[500px] rounded-[32px] overflow-hidden bg-gradient-to-br from-[#dbeafe] to-[#60a5fa]/40 flex items-center justify-center shadow-sm">
            <ScissorsIcon class="w-48 h-48 text-[#1d4ed8]/20" />

            <div class="absolute top-6 left-6 flex gap-2">
              <span class="bg-white/90 backdrop-blur-sm text-[#006d35] font-bold text-sm px-4 py-2 rounded-full shadow-sm">
                {{ t('public.eventDetail.badgeWorkshop') }}
              </span>
              <span class="bg-[#006d35] text-white font-bold text-sm px-4 py-2 rounded-full shadow-sm">
                {{ t('public.eventDetail.badgePopular') }}
              </span>
            </div>
          </div>

          <div class="flex flex-col gap-4">
            <h1 class="font-jakarta font-extrabold text-[#001d32] text-[60px] leading-[1] tracking-tight">
              {{ t('public.eventDetail.demoTitle') }}
            </h1>
            <div class="flex flex-wrap gap-6 text-[#40617f] text-base">
              <div class="flex items-center gap-2">
                <CalendarDaysIcon class="w-5 h-5 shrink-0" />
                <span>{{ t('public.eventDetail.demoDate') }}</span>
              </div>
              <div class="flex items-center gap-2">
                <ClockIcon class="w-5 h-5 shrink-0" />
                <span>{{ t('public.eventDetail.demoHours') }}</span>
              </div>
              <div class="flex items-center gap-2">
                <MapPinIcon class="w-5 h-5 shrink-0" />
                <span>{{ t('public.eventDetail.demoLocation') }}</span>
              </div>
            </div>
          </div>

          <div class="bg-white rounded-[32px] p-10 shadow-sm flex flex-col gap-6">
            <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.eventDetail.aboutTitle') }}</h2>
            <p class="text-[#40617f] text-lg leading-relaxed">
              {{ t('public.eventDetail.aboutPara1') }}
            </p>
            <p class="text-[#40617f] text-base leading-relaxed">
              {{ t('public.eventDetail.aboutPara2') }}
            </p>
          </div>

          <div class="bg-[#edf4ff] rounded-[32px] p-8 flex gap-6">
            <div class="w-24 h-24 rounded-2xl overflow-hidden bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] flex items-center justify-center shrink-0 shadow-md">
              <UserCircleIcon class="w-16 h-16 text-[#006d35]/40" />
            </div>
            <div class="flex flex-col gap-1">
              <p class="text-[#006d35] text-xs font-bold uppercase tracking-widest">{{ t('public.eventDetail.facilitatorTag') }}</p>
              <h3 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.eventDetail.facilitatorName') }}</h3>
              <p class="text-[#40617f] text-base leading-relaxed mt-1">
                {{ t('public.eventDetail.facilitatorBio') }}
              </p>
            </div>
          </div>

        </div>

        <div class="col-span-4 flex flex-col gap-6">

          <div class="bg-white border border-[#d8eaff] rounded-[40px] p-8 shadow-[0_12px_40px_0_rgba(0,29,50,0.06)] relative overflow-hidden">

            <div class="absolute -right-8 -top-8 w-32 h-32 rounded-full bg-[rgba(0,109,53,0.05)]" />

            <div class="flex items-baseline justify-between mb-3">
              <span class="font-jakarta font-extrabold text-[#001d32] text-4xl">25€</span>
              <span class="text-[#40617f] text-base">{{ t('public.eventDetail.perPerson') }}</span>
            </div>
            <div class="flex items-center gap-2 bg-[rgba(255,218,214,0.3)] px-3 py-1 rounded-full w-fit mb-8">
              <FireIcon class="w-3.5 h-3.5 text-[#ba1a1a]" />
              <span class="text-[#ba1a1a] font-bold text-sm">{{ t('public.eventDetail.spotsLeft') }}</span>
            </div>

            <div class="flex flex-col gap-3 mb-8">
              <div class="bg-[#f7f9ff] rounded-2xl p-3 flex items-center gap-3">
                <WrenchScrewdriverIcon class="w-5 h-5 text-[#006d35] shrink-0" />
                <span class="text-[#001d32] font-semibold text-sm">{{ t('public.eventDetail.infoMaterial') }}</span>
              </div>
              <div class="bg-[#f7f9ff] rounded-2xl p-3 flex items-center gap-3">
                <AcademicCapIcon class="w-5 h-5 text-[#006d35] shrink-0" />
                <span class="text-[#001d32] font-semibold text-sm">{{ t('public.eventDetail.infoLevel') }}</span>
              </div>
              <div class="bg-[#f7f9ff] rounded-2xl p-3 flex items-center gap-3">
                <SparklesIcon class="w-5 h-5 text-[#006d35] shrink-0" />
                <span class="text-[#001d32] font-semibold text-sm">{{ t('public.eventDetail.infoSnack') }}</span>
              </div>
            </div>

            <button
              class="w-full py-5 rounded-2xl text-white font-bold text-lg flex items-center justify-center gap-2 transition hover:opacity-90 shadow-[0_10px_15px_-3px_rgba(0,109,53,0.2)]"
              style="background: linear-gradient(168deg, #006d35 0%, #1b8848 100%);"
            >
              {{ t('public.eventDetail.ctaReserve') }}
              <ArrowRightIcon class="w-4 h-4" />
            </button>
          </div>

          <div class="bg-[rgba(216,234,255,0.5)] border border-white rounded-[32px] p-6">
            <div class="flex items-center justify-between mb-4">
              <h3 class="font-jakarta font-bold text-[#001d32] text-base">{{ t('public.eventDetail.scoreTitle') }}</h3>
              <span class="font-bold text-[#006d35] text-2xl">94/100</span>
            </div>
            <div class="h-3 bg-white rounded-full overflow-hidden mb-3">
              <div class="h-full bg-[#006d35] rounded-full" style="width: 94%;" />
            </div>
            <p class="text-[#40617f] text-xs font-bold uppercase tracking-wide">{{ t('public.eventDetail.scoreImpact') }}</p>
          </div>

          <div class="bg-white border border-[#d8eaff] rounded-[32px] p-6 flex items-center justify-between">
            <div>
              <p class="font-jakarta font-bold text-[#001d32] text-sm">{{ t('public.eventDetail.locationTitle') }}</p>
              <p class="text-[#40617f] text-xs mt-0.5">{{ t('public.eventDetail.locationAddress') }}</p>
            </div>
            <div class="w-12 h-12 bg-[#d8eaff] rounded-xl flex items-center justify-center shrink-0">
              <MapPinIcon class="w-5 h-5 text-[#40617f]" />
            </div>
          </div>

        </div>

      </div>
    </div>
  </div>
</template>

<script setup>import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  ChevronRightIcon,
  CalendarDaysIcon,
  ClockIcon,
  MapPinIcon,
  UserCircleIcon,
  FireIcon,
  WrenchScrewdriverIcon,
  AcademicCapIcon,
  SparklesIcon,
  ArrowRightIcon,
  ScissorsIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
</script>
