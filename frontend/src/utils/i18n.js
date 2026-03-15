import { createI18n } from 'vue-i18n'
import fr from '@/locales/fr.json'
import en from '@/locales/en.json'

const savedLang = localStorage.getItem('admin_lang') || 'fr'

const i18n = createI18n({
  legacy: false,
  locale: savedLang,
  fallbackLocale: 'fr',
  messages: { fr, en },
})

export default i18n
export const { t } = i18n.global
