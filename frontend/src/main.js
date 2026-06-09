import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from '@/utils/i18n.js'
import './assets/main.css'
import { useUserAuthStore } from '@/stores/userAuth'

// When a lazy-loaded chunk fails to load (stale deployment), reload to get fresh chunks
window.addEventListener('vite:preloadError', () => {
  window.location.reload()
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)

// Initialize auth BEFORE mounting so that isProvider / profilePath are correct
// on the very first render — avoids the "menu briefly shows /profil for pros" race.
const authStore = useUserAuthStore()
authStore.init().finally(() => {
  app.mount('#app')
})
