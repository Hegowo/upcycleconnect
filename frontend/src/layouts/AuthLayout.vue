<template>
  <div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <img
          src="/logoentier.png"
          alt="UpcycleConnect"
          class="mx-auto h-40 w-auto"
          @error="onLogoError"
          ref="logoImg"
        />
        <p class="text-gray-500 text-sm mt-2">{{ t('nav.administration') }}</p>
      </div>
      <slot />
    </div>
  </div>
</template>

<script setup>import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const logoImg = ref(null)

function onLogoError() {
  if (logoImg.value) logoImg.value.src = '/logoentier.svg'
}
</script>
