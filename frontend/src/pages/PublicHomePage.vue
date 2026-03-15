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
            class="px-8 py-3 rounded-xl font-semibold text-white transition"
            style="background-color: #1B8848;"
          >
            Voir les prestations
          </RouterLink>
          <RouterLink
            to="/inscription"
            class="px-8 py-3 rounded-xl font-semibold bg-white transition"
            style="color: #103652;"
          >
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
          class="rounded-xl p-6 text-center cursor-pointer hover:shadow-md transition border"
          style="border-color: #e5e7eb;"
        >
          <div class="text-4xl mb-3">{{ cat.icon }}</div>
          <p class="font-medium text-sm" style="color: #10522b;">{{ cat.label }}</p>
        </div>
      </div>
    </section>

    <section class="py-16 px-4" style="background-color: #f0faf4;">
      <div class="max-w-2xl mx-auto text-center">
        <h2 class="text-2xl font-bold mb-3" style="color: #10522b;">Vous êtes artisan ou prestataire ?</h2>
        <p class="text-gray-600 mb-6">
          Inscrivez-vous gratuitement et proposez vos services d'upcycling à des milliers de clients.
        </p>
        <RouterLink
          to="/inscription?type=prestataire"
          class="inline-block px-8 py-3 rounded-xl font-semibold text-white transition"
          style="background-color: #1B8848;"
        >
          Créer mon profil prestataire
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup>
const categories = [
  { icon: '🪑', label: 'Mobilier' },
  { icon: '👗', label: 'Textile' },
  { icon: '📦', label: 'Emballages' },
  { icon: '🔩', label: 'Matériaux' },
  { icon: '🖼️', label: 'Décoration' },
  { icon: '💡', label: 'Électronique' },
  { icon: '🌿', label: 'Végétal' },
  { icon: '🎨', label: 'Créations' },
]
</script>
