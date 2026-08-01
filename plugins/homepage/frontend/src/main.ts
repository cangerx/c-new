import { createApp } from 'vue'

import i18n, { setLocale } from './i18n'
import './style.css'
import HomeView from './views/HomeView.vue'

// Restore the persisted locale so <html lang> matches on first paint.
setLocale(i18n.global.locale.value as string)

createApp(HomeView).use(i18n).mount('#app')
