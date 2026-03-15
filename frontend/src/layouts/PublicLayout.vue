<template>
  <div class="min-h-screen flex flex-col bg-white">
    <header class="shadow-sm sticky top-0 z-50 bg-white">
      <div class="max-w-6xl mx-auto px-4 py-3 flex items-center justify-between">
        <RouterLink to="/">
          <img src="/logoentier.png" class="h-12 w-auto" alt="UpcycleConnect" />
        </RouterLink>
        <nav class="hidden md:flex items-center gap-6 text-sm font-medium" style="color: #10522b;">
          <RouterLink to="/prestations" class="hover:underline">Prestations</RouterLink>
          <RouterLink to="/evenements" class="hover:underline">Événements</RouterLink>
          <RouterLink to="/prestataires" class="hover:underline">Prestataires</RouterLink>
        </nav>
        <div class="flex items-center gap-3">
          <RouterLink
            to="/connexion"
            class="text-sm font-medium px-4 py-2 rounded-lg border transition"
            style="color: #103652; border-color: #103652;"
          >
            Connexion
          </RouterLink>
          <RouterLink
            to="/inscription"
            class="text-sm font-medium px-4 py-2 rounded-lg text-white transition"
            style="background-color: #1B8848;"
          >
            S'inscrire
          </RouterLink>
        </div>
      </div>
    </header>

    <main class="flex-1">
      <RouterView />
    </main>

    <footer class="py-8 mt-12" style="background-color: #103652; color: #cbd5e1;">
      <div class="max-w-6xl mx-auto px-4 flex flex-col md:flex-row justify-between items-center gap-4">
        <img src="/logoentier.png" class="h-10 w-auto brightness-0 invert opacity-80" alt="UpcycleConnect" />
        <p class="text-sm opacity-60">© {{ new Date().getFullYear() }} UpcycleConnect — Tous droits réservés</p>
        <div class="flex gap-4 text-sm opacity-70">
          <RouterLink to="/mentions-legales" class="hover:opacity-100">Mentions légales</RouterLink>
          <RouterLink to="/contact" class="hover:opacity-100">Contact</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
</script>
