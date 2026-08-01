// Shim for sub2api's pinia stores, backed by new-api's endpoints.
//
// HomeView.vue was written against sub2api's `useAppStore` / `useAuthStore`.
// Rather than editing 1700 lines of view code, this shim exposes the same
// surface HomeView reads and maps it onto new-api's API:
//
//   sub2api field   ->  new-api source
//   site_name       ->  /api/status  .system_name
//   site_logo       ->  /api/status  .logo
//   doc_url         ->  /api/status  .docs_link
//   home_content    ->  /api/home_page_content  .data
//   site_subtitle   ->  (no equivalent; always empty)
//
// Keeping the shim here means new-api upstream changes can only ever break
// this one file, not the ported view.
import { reactive, readonly } from 'vue'

export type PublicSettings = {
  site_name: string
  site_logo: string
  site_subtitle: string
  doc_url: string
  home_content: string
}

const state = reactive({
  settings: null as PublicSettings | null,
  loaded: false,
  role: 0,
  authenticated: false,
})

// new-api role constants (common/constants.go): admin is 10, root is 100.
const ROLE_ADMIN = 10

async function getJson(url: string): Promise<any> {
  const res = await fetch(url, { credentials: 'same-origin' })
  if (!res.ok) {
    throw new Error(`${url} -> HTTP ${res.status}`)
  }
  return res.json()
}

export function useAppStore() {
  return {
    get cachedPublicSettings() {
      return state.settings
    },
    get publicSettingsLoaded() {
      return state.loaded
    },
    get siteName() {
      return state.settings?.site_name ?? ''
    },
    get siteLogo() {
      return state.settings?.site_logo ?? ''
    },
    get docUrl() {
      return state.settings?.doc_url ?? ''
    },

    async fetchPublicSettings(): Promise<PublicSettings | null> {
      try {
        const [status, home] = await Promise.all([
          getJson('/api/status'),
          getJson('/api/home_page_content').catch(() => null),
        ])
        const d = status?.data ?? {}
        state.settings = {
          site_name: d.system_name ?? '',
          site_logo: d.logo ?? '',
          site_subtitle: '',
          doc_url: d.docs_link ?? '',
          home_content: typeof home?.data === 'string' ? home.data : '',
        }
      } catch {
        // The homepage must render even if the API is unreachable; HomeView
        // falls back to its own defaults when fields are empty.
        state.settings = {
          site_name: '',
          site_logo: '',
          site_subtitle: '',
          doc_url: '',
          home_content: '',
        }
      }
      state.loaded = true
      return state.settings
    },
  }
}

export function useAuthStore() {
  return {
    // Plain getters, not computed refs: HomeView wraps these in its own
    // computed(), and a ref inside a computed would not auto-unwrap in the
    // template — isAuthenticated would be a truthy object no matter what.
    get isAuthenticated() {
      return state.authenticated
    },
    get isAdmin() {
      return state.role >= ROLE_ADMIN
    },

    // Resolves login state so the nav can show "Dashboard" vs "Sign in".
    // /api/user/self returns 401 for anonymous visitors, which is expected
    // and simply means "not logged in".
    async probe(): Promise<void> {
      try {
        const res = await fetch('/api/user/self', { credentials: 'same-origin' })
        if (!res.ok) {
          return
        }
        const body = await res.json()
        if (body?.success && body?.data) {
          state.authenticated = true
          state.role = Number(body.data.role ?? 0)
        }
      } catch {
        // Stay anonymous on any failure.
      }
    },
  }
}

export const homepageState = readonly(state)
