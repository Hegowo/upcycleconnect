<template>
  <div class="flex h-screen bg-gray-50 overflow-hidden">
    <TheSidebar />
    <div class="flex-1 flex flex-col overflow-hidden">
      <TheHeader />
      <main :class="route.meta.fullscreen ? 'flex-1 overflow-hidden flex flex-col' : 'flex-1 overflow-y-auto p-6'">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import TheSidebar from '@/components/layout/TheSidebar.vue'
import TheHeader from '@/components/layout/TheHeader.vue'

const route = useRoute()
</script>
