// Shim for sub2api's @/i18n module. Only the pieces HomeView and
// LocaleSwitcher actually use are provided: createI18n setup, setLocale,
// and availableLocales.
import { createI18n } from 'vue-i18n'

import en from './locales/en/landing'
import zh from './locales/zh/landing'

export const availableLocales = [
  { code: 'en', name: 'English', flag: '🇺🇸' },
  { code: 'zh', name: '中文', flag: '🇨🇳' },
]

const STORAGE_KEY = 'homepage-plugin:locale'

function detectLocale(): string {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored && availableLocales.some((l) => l.code === stored)) {
    return stored
  }
  return navigator.language.toLowerCase().startsWith('zh') ? 'zh' : 'en'
}

export const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'zh',
  messages: { en, zh },
})

export function setLocale(code: string): void {
  if (!availableLocales.some((l) => l.code === code)) {
    return
  }
  // `locale` is a WritableComputedRef in composition mode.
  ;(i18n.global.locale as unknown as { value: string }).value = code
  localStorage.setItem(STORAGE_KEY, code)
  document.documentElement.lang = code
}

export default i18n
