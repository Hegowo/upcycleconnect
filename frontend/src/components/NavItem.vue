<template>
  <RouterLink
    :to="to"
    class="nav-item flex items-center gap-3 px-3 py-2 rounded-md text-sm font-medium transition-colors"
    :class="isActive ? 'nav-item--active' : 'nav-item--default'"
  >
    <slot />
  </RouterLink>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps({
  to: { type: String, required: true },
})

const route = useRoute()
const isActive = computed(() => route.path.startsWith(props.to))
</script>

<style scoped>
.nav-item--default {
  color: #c5d8e8;
}
.nav-item--default:hover {
  background-color: #174d72;
  color: #ffffff;
}
.nav-item--active {
  background-color: #1B5B88;
  color: #ffffff;
}
</style>
