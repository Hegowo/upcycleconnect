<template>
  <div class="flex h-[100dvh] overflow-hidden" style="background-color:#f8fafc;">
    <TheSidebar />
    <div class="flex-1 flex flex-col overflow-hidden min-w-0">
      <TheHeader />
      <main :class="route.meta.fullscreen
        ? 'flex-1 overflow-hidden flex flex-col'
        : 'flex-1 overflow-y-auto p-4 md:p-6 main-pad'"
      >
        <RouterView />
      </main>
    </div>
    <AdminBottomNav />
  </div>
</template>

<script setup>import { useRoute } from 'vue-router'
import TheSidebar from '@/components/Sidebar.vue'
import TheHeader from '@/components/Header.vue'
import AdminBottomNav from '@/components/AdminBottomNav.vue'

const route = useRoute()
</script>

<style scoped>
@media (max-width: 1023px) {
  .main-pad {
    padding-bottom: calc(64px + env(safe-area-inset-bottom, 0px));
  }
}
@media (min-width: 1024px) {
  .main-pad {
    padding-bottom: 1.5rem;
  }
}
</style>
