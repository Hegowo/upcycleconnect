<template>
  <div>
    <section class="py-20 px-4 text-center" style="background: linear-gradient(135deg, #103652 0%, #1B5B88 100%);">
      <div class="max-w-3xl mx-auto">
        <h1 class="text-4xl md:text-5xl font-bold text-white mb-4">
          Donnez une seconde vie à vos objets
        </h1>
        <p class="text-lg text-blue-100 mb-8">
          Trouvez des artisans et prestataires spécialisés dans l'upcycling près de chez vous.
        </p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <RouterLink
            to="/prestations"
            class="inline-flex items-center justify-center gap-2 px-8 py-3 rounded-xl font-semibold text-white transition hover:opacity-90"
            style="background-color: #1B8848;"
          >
            <MagnifyingGlassIcon class="w-5 h-5" />
            Voir les prestations
          </RouterLink>
          <RouterLink
            to="/inscription"
            class="inline-flex items-center justify-center gap-2 px-8 py-3 rounded-xl font-semibold bg-white transition hover:opacity-90"
            style="color: #103652;"
          >
            <BriefcaseIcon class="w-5 h-5" />
            Devenir prestataire
          </RouterLink>
        </div>
      </div>
    </section>

    <section class="py-16 px-4 max-w-6xl mx-auto">
      <h2 class="text-2xl font-bold mb-8 text-center" style="color: #10522b;">
        Nos catégories
      </h2>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-6">
        <div
          v-for="cat in categories"
          :key="cat.label"
          class="rounded-2xl p-6 text-center cursor-pointer hover:shadow-lg transition-all duration-200 hover:scale-[1.03] border border-transparent"
          :style="`background: ${cat.bg};`"
        >
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center mx-auto mb-3" :style="`background: ${cat.iconBg};`">
            <component :is="cat.icon" class="w-7 h-7" :style="`color: ${cat.color};`" />
          </div>
          <p class="font-semibold text-sm" :style="`color: ${cat.color};`">{{ cat.label }}</p>
        </div>
      </div>
    </section>

    <section class="py-16 px-4" style="background-color: #f0faf4;">
      <div class="max-w-2xl mx-auto text-center">
        <div class="w-16 h-16 rounded-2xl flex items-center justify-center mx-auto mb-4" style="background:#dcfce7;">
          <UserPlusIcon class="w-8 h-8" style="color:#1B8848;" />
        </div>
        <h2 class="text-2xl font-bold mb-3" style="color: #10522b;">Vous êtes artisan ou prestataire ?</h2>
        <p class="text-gray-600 mb-6">
          Inscrivez-vous gratuitement et proposez vos services d'upcycling à des milliers de clients.
        </p>
        <RouterLink
          to="/inscription?type=prestataire"
          class="inline-flex items-center justify-center gap-2 px-8 py-3 rounded-xl font-semibold text-white transition hover:opacity-90"
          style="background-color: #1B8848;"
        >
          <SparklesIcon class="w-5 h-5" />
          Créer mon profil prestataire
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup>
import {
  MagnifyingGlassIcon,
  BriefcaseIcon,
  UserPlusIcon,
  SparklesIcon,
  HomeModernIcon,
  ScissorsIcon,
  ArchiveBoxIcon,
  WrenchIcon,
  PhotoIcon,
  LightBulbIcon,
  SunIcon,
  PaintBrushIcon,
} from '@heroicons/vue/24/outline'

const categories = [
  {
    label: 'Mobilier',
    icon: HomeModernIcon,
    color: '#92400e',
    bg: '#fffbeb',
    iconBg: '#fef3c7',
  },
  {
    label: 'Textile',
    icon: ScissorsIcon,
    color: '#9d174d',
    bg: '#fff1f2',
    iconBg: '#ffe4e6',
  },
  {
    label: 'Emballages',
    icon: ArchiveBoxIcon,
    color: '#1e40af',
    bg: '#eff6ff',
    iconBg: '#dbeafe',
  },
  {
    label: 'Matériaux',
    icon: WrenchIcon,
    color: '#374151',
    bg: '#f9fafb',
    iconBg: '#f3f4f6',
  },
  {
    label: 'Décoration',
    icon: PhotoIcon,
    color: '#6d28d9',
    bg: '#f5f3ff',
    iconBg: '#ede9fe',
  },
  {
    label: 'Électronique',
    icon: LightBulbIcon,
    color: '#b45309',
    bg: '#fffbeb',
    iconBg: '#fef9c3',
  },
  {
    label: 'Végétal',
    icon: SunIcon,
    color: '#166534',
    bg: '#f0fdf4',
    iconBg: '#dcfce7',
  },
  {
    label: 'Créations',
    icon: PaintBrushIcon,
    color: '#c2410c',
    bg: '#fff7ed',
    iconBg: '#ffedd5',
  },
]
</script>
