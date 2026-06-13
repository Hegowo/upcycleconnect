import { createI18n } from 'vue-i18n'
import { ref } from 'vue'
import fr from '@/locales/fr.json'
import en from '@/locales/en.json'

const savedLang = localStorage.getItem('admin_lang') || 'fr'

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'fr',
  messages: { fr, en },
})

// Built-in languages always available (shipped with the app).
export const builtinLocales = [
  { code: 'fr', name: 'Français', builtin: true },
  { code: 'en', name: 'English', builtin: true },
]

// availableLocales is reactive: built-in + any language added via the back-office.
export const availableLocales = ref([...builtinLocales])

// fetchAndRegister pulls one dynamic locale's messages and registers them into vue-i18n.
export async function fetchAndRegister(code) {
  if (code === 'fr' || code === 'en') return true
  try {
    const res = await fetch(`/api/public/v1/locales/${code}`, { headers: { Accept: 'application/json' } })
    if (!res.ok) return false
    const messages = await res.json()
    i18n.global.setLocaleMessage(code, messages)
    return true
  } catch {
    return false
  }
}

// loadDynamicLocales fetches admin-managed languages from the API and registers the
// active one's messages — so a language added in the back-office works WITHOUT any
// code change or rebuild (spec requirement, page 7).
export async function loadDynamicLocales() {
  try {
    const res = await fetch('/api/public/v1/locales', { headers: { Accept: 'application/json' } })
    if (!res.ok) return
    const json = await res.json()
    const list = json.data || []
    const merged = [...builtinLocales]
    for (const loc of list) {
      if (loc.builtin) continue
      merged.push({ code: loc.code, name: loc.name, builtin: false })
      if (loc.code === i18n.global.locale.value) {
        await fetchAndRegister(loc.code)
      }
    }
    availableLocales.value = merged
  } catch {
    /* offline / not configured — keep built-in only */
  }
}

// setLocale switches the active language, loading its messages on demand if dynamic.
export async function setLocale(code) {
  if (code !== 'fr' && code !== 'en') {
    await fetchAndRegister(code)
  }
  i18n.global.locale.value = code
  localStorage.setItem('admin_lang', code)
}

export default i18n
export const { t } = i18n.global
