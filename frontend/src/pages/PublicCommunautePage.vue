<template>
  <div class="bg-[#f7f9ff] min-h-screen">
    <div class="max-w-[1280px] mx-auto px-6 py-8 flex gap-8 items-start">

      <aside class="bg-[#edf4ff] rounded-[24px] p-4 w-[256px] shrink-0 flex flex-col gap-2 sticky top-24">
        <div class="px-2 pb-6">
          <p class="font-jakarta font-bold text-[#001d32] text-lg">{{ t('public.community.navTitle') }}</p>
          <p class="text-[#64748b] text-xs font-medium mt-0.5">{{ t('public.community.navSubtitle') }}</p>
        </div>

        <nav class="flex flex-col gap-2">
          <button
            v-for="item in navItems"
            :key="item.key"
            @click="activeNav = item.key"
            class="flex items-center gap-3 px-4 py-3 rounded-xl text-sm font-medium transition-all"
            :class="activeNav === item.key
              ? 'bg-white text-[#006d35] shadow-sm'
              : 'text-[#475569] hover:bg-white/60'"
          >
            <component :is="item.icon" class="w-5 h-5 shrink-0" />
            {{ item.label }}
          </button>
        </nav>

        <div class="pt-8">
          <button
            class="w-full py-3 rounded-xl text-white font-semibold text-sm transition hover:opacity-90"
            style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
          >
            {{ t('public.community.newPost') }}
          </button>
        </div>
      </aside>

      <div class="flex-1 min-w-0 flex flex-col gap-10">

        <section
          class="relative h-64 rounded-[24px] overflow-hidden"
          style="background: #b9dbfe;"
        >
          <div class="absolute inset-0 bg-gradient-to-r from-[rgba(185,219,254,0.9)] to-[rgba(185,219,254,0.1)] flex flex-col justify-center p-12">
            <h1 class="font-jakarta font-extrabold text-[#3f607e] text-5xl leading-[1.1] tracking-tight">
              {{ t('public.community.heroTitle1') }}<br />{{ t('public.community.heroTitle2') }}
            </h1>
            <p class="text-[#3f607e] text-lg mt-4 max-w-md">
              {{ t('public.community.heroSubtitle') }}
            </p>
          </div>

          <div class="absolute right-12 top-8 w-40 h-40 rounded-full bg-[rgba(185,219,254,0.5)] blur-2xl" />
          <div class="absolute right-32 bottom-4 w-24 h-24 rounded-full bg-[rgba(0,109,53,0.1)] blur-xl" />
        </section>

        <section>
          <div class="flex items-end justify-between mb-6">
            <div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.community.projectsTitle') }}</h2>
              <p class="text-[#64748b] text-base">{{ t('public.community.projectsSubtitle') }}</p>
            </div>
            <button class="text-[#006d35] font-semibold text-sm hover:underline">{{ t('public.community.seeAll') }}</button>
          </div>

          <div class="grid grid-cols-12 gap-6">

            <div class="col-span-8 bg-white rounded-[24px] shadow-sm p-6 flex gap-6">
              <div class="w-[280px] h-[280px] rounded-xl overflow-hidden shrink-0 bg-gradient-to-br from-[#d1fae5] to-[#a7f3d0] flex items-center justify-center">
                <HomeModernIcon class="w-24 h-24 text-[#006d35]/30" />
              </div>
              <div class="flex flex-col justify-between flex-1">
                <div>
                  <span class="inline-block bg-[#b0f2bd] text-[#0f522b] text-[10px] font-bold uppercase tracking-wide px-3 py-1 rounded-full mb-3">
                    {{ t('public.community.tagSustainable') }}
                  </span>
                  <h3 class="font-jakarta font-bold text-[#001d32] text-2xl mb-3">{{ t('public.community.demoProjectTitle') }}</h3>
                  <p class="text-[#475569] text-sm leading-relaxed">
                    {{ t('public.community.demoProjectDesc') }}
                  </p>
                </div>
                <div class="bg-[#edf4ff] rounded-xl p-4 mt-4">
                  <div class="flex justify-between text-xs font-semibold mb-2">
                    <span class="text-[#006d35]">{{ t('public.community.impactLabel') }}</span>
                    <span class="text-[#001d32]">{{ t('public.community.impactValue') }}</span>
                  </div>
                  <div class="bg-white h-2 rounded-full overflow-hidden">
                    <div class="h-full bg-[#006d35] rounded-full" style="width: 85%;" />
                  </div>
                </div>
              </div>
            </div>

            <div class="col-span-4 bg-white rounded-[24px] shadow-sm p-6 flex flex-col justify-between">
              <div class="w-full h-36 rounded-xl overflow-hidden bg-gradient-to-br from-[#d1fae5] to-[#bbf7d0] flex items-center justify-center mb-4">
                <SparklesIcon class="w-16 h-16 text-[#006d35]/30" />
              </div>
              <div>
                <h3 class="font-jakarta font-bold text-[#001d32] text-lg mb-1">{{ t('public.community.demoSecondaryTitle') }}</h3>
                <p class="text-[#475569] text-sm leading-relaxed">{{ t('public.community.demoSecondaryDesc') }}</p>
              </div>
              <div class="flex items-center justify-between pt-4 border-t border-gray-100 mt-4">
                <div class="flex items-center gap-1.5 text-[#94a3b8] text-xs">
                  <HeartIcon class="w-3.5 h-3.5" />
                  <span>124</span>
                </div>
                <div class="flex items-center gap-1.5 text-[#94a3b8] text-xs">
                  <ChatBubbleLeftIcon class="w-3.5 h-3.5" />
                  <span>18</span>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section class="bg-[#edf4ff] rounded-[24px] p-8">
          <div class="flex items-center justify-between mb-8">
            <div>
              <h2 class="font-jakarta font-bold text-[#001d32] text-2xl">{{ t('public.community.forumTitle') }}</h2>
              <p class="text-[#64748b] text-sm">{{ t('public.community.forumSubtitle') }}</p>
            </div>
            <button class="bg-white border border-[rgba(0,109,53,0.1)] text-[#006d35] font-semibold text-sm px-4 py-2 rounded-xl shadow-sm hover:bg-[#edf4ff] transition">
              {{ t('public.community.viewCategories') }}
            </button>
          </div>

          <div class="flex flex-col gap-3">
            <div
              v-for="thread in forumThreads"
              :key="thread.id"
              class="bg-white rounded-xl shadow-sm pl-6 pr-5 py-5 flex items-center justify-between cursor-pointer hover:shadow-md transition-shadow"
              :style="`border-left-color: ${thread.color};`"
              style="border-left-width: 4px; border-left-style: solid;"
            >
              <div class="flex items-center gap-4 flex-1 min-w-0">
                <div
                  class="w-10 h-10 rounded-full flex items-center justify-center shrink-0"
                  :style="`background: ${thread.avatarBg};`"
                >
                  <UserCircleIcon class="w-5 h-5 text-[#40617f]" />
                </div>
                <div class="min-w-0">
                  <p class="font-semibold text-[#001d32] text-base leading-snug truncate">{{ thread.title }}</p>
                  <p class="text-[#94a3b8] text-xs mt-0.5">
                    {{ t('public.community.postedBy') }}
                    <span class="font-semibold text-[#475569]">{{ thread.author }}</span>
                    {{ t('public.community.postedIn') }}
                    <span class="font-medium text-[#40617f]">{{ thread.category }}</span>
                  </p>
                </div>
              </div>
              <div class="flex items-center gap-6 shrink-0 ml-4">
                <div class="text-center">
                  <p class="font-semibold text-[#001d32] text-sm text-center">{{ thread.replies }}</p>
                  <p class="text-[#94a3b8] text-[10px] uppercase tracking-wide font-semibold">{{ t('public.community.repliesLabel') }}</p>
                </div>
                <ChevronRightIcon class="w-4 h-4 text-[#94a3b8]" />
              </div>
            </div>
          </div>
        </section>

        <section>
          <div class="flex items-center gap-4 mb-8">
            <h2 class="font-jakarta font-bold text-[#001d32] text-2xl shrink-0">{{ t('public.community.newsTitle') }}</h2>
            <div class="flex-1 h-0.5 bg-[#d8eaff]" />
          </div>

          <div class="grid grid-cols-2 gap-6">
            <article
              v-for="article in articles"
              :key="article.id"
              class="bg-white rounded-[24px] shadow-sm p-4 flex gap-4 hover:shadow-md transition-shadow cursor-pointer"
            >
              <div
                class="w-32 h-32 rounded-xl flex items-center justify-center shrink-0"
                :style="`background: linear-gradient(135deg, ${article.colorFrom}, ${article.colorTo});`"
              >
                <component :is="article.icon" class="w-10 h-10 text-white/40" />
              </div>
              <div class="flex flex-col justify-center flex-1 min-w-0">
                <p class="text-[10px] font-semibold uppercase tracking-wider mb-1" :style="`color: ${article.tagColor};`">
                  {{ article.tag }}
                </p>
                <h3 class="font-jakarta font-bold text-[#001d32] text-base leading-snug mb-2">{{ article.title }}</h3>
                <p class="text-[#64748b] text-xs leading-relaxed line-clamp-2">{{ article.excerpt }}</p>
              </div>
            </article>
          </div>
        </section>

      </div>
    </div>
  </div>
</template>

<script setup>import {
  ChatBubbleLeftEllipsisIcon,
  FolderOpenIcon,
  NewspaperIcon,
  LightBulbIcon,
  HomeModernIcon,
  SparklesIcon,
  HeartIcon,
  ChatBubbleLeftIcon,
  UserCircleIcon,
  ChevronRightIcon,
  ScissorsIcon,
  AcademicCapIcon,
} from '@heroicons/vue/24/outline'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const activeNav = ref('forum')

const navItems = computed(() => [
  { key: 'forum',    label: t('public.community.navForum'),    icon: ChatBubbleLeftEllipsisIcon },
  { key: 'projects', label: t('public.community.navProjects'), icon: FolderOpenIcon },
  { key: 'news',     label: t('public.community.navNews'),     icon: NewspaperIcon },
  { key: 'tips',     label: t('public.community.navTips'),     icon: LightBulbIcon },
])

const forumThreads = computed(() => [
  {
    id: 1,
    title: t('public.community.threadTitle1'),
    author: 'Marc_B',
    category: t('public.community.threadCategory1'),
    replies: 42,
    color: '#006d35',
    avatarBg: '#b9dbfe',
  },
  {
    id: 2,
    title: t('public.community.threadTitle2'),
    author: 'Sarah_Art',
    category: t('public.community.threadCategory2'),
    replies: 15,
    color: '#40617f',
    avatarBg: '#cee5ff',
  },
  {
    id: 3,
    title: t('public.community.threadTitle3'),
    author: 'EcoWarrior',
    category: t('public.community.threadCategory3'),
    replies: 128,
    color: '#006d35',
    avatarBg: '#b9dbfe',
  },
])

const articles = computed(() => [
  {
    id: 1,
    tag: t('public.community.tagTrend'),
    tagColor: '#006d35',
    title: t('public.community.articleTitle1'),
    excerpt: t('public.community.articleExcerpt1'),
    colorFrom: '#dbeafe',
    colorTo: '#bfdbfe',
    icon: SparklesIcon,
  },
  {
    id: 2,
    tag: t('public.community.tagPractical'),
    tagColor: '#40617f',
    title: t('public.community.articleTitle2'),
    excerpt: t('public.community.articleExcerpt2'),
    colorFrom: '#d1fae5',
    colorTo: '#a7f3d0',
    icon: LightBulbIcon,
  },
  {
    id: 3,
    tag: t('public.community.tagInspiration'),
    tagColor: '#7c3aed',
    title: t('public.community.articleTitle3'),
    excerpt: t('public.community.articleExcerpt3'),
    colorFrom: '#f3e8ff',
    colorTo: '#e9d5ff',
    icon: AcademicCapIcon,
  },
  {
    id: 4,
    tag: t('public.community.tagTechnique'),
    tagColor: '#b45309',
    title: t('public.community.articleTitle4'),
    excerpt: t('public.community.articleExcerpt4'),
    colorFrom: '#fef3c7',
    colorTo: '#fde68a',
    icon: ScissorsIcon,
  },
])
</script>
