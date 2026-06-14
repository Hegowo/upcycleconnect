import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n, { loadDynamicLocales } from '@/utils/i18n.js'
import './assets/main.css'
import { useUserAuthStore } from '@/stores/userAuth'

window.addEventListener('vite:preloadError', () => {
  window.location.reload()
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(i18n)

const authStore = useUserAuthStore()
Promise.allSettled([authStore.init(), loadDynamicLocales()]).finally(() => {
  app.mount('#app')
})
