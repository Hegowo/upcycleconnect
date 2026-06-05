<template>
  <RouterView />
  <AppToast />
</template>

<script setup>
import { onMounted, watch } from 'vue'
import AppToast from '@/components/AppToast.vue'
import { useUserAuthStore } from '@/stores/userAuth'
import { syncExistingSubscription } from '@/utils/onesignal'

const userAuth = useUserAuthStore()

// On boot or after login, re-sync the OneSignal player ID with the backend
// so push delivery keeps working even if the ID was rotated by the browser.
onMounted(() => {
  if (userAuth.isLoggedIn) syncExistingSubscription()
})
watch(() => userAuth.isLoggedIn, (val) => {
  if (val) syncExistingSubscription()
})
</script>
