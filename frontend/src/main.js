import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from '@/utils/i18n.js'
import './assets/main.css'

// When a lazy-loaded chunk fails to load (stale deployment), reload to get fresh chunks
window.addEventListener('vite:preloadError', () => {
  window.location.reload()
})

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

app.mount('#app')
