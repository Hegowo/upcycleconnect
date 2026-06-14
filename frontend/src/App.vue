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

onMounted(() => {
  if (userAuth.isLoggedIn) syncExistingSubscription()
})
watch(() => userAuth.isLoggedIn, (val) => {
  if (val) syncExistingSubscription()
})
</script>
