import { useTheme } from '@/context/theme-provider'
import { useHomePageContent } from '@/features/home/hooks'
import { useStatus } from '@/hooks/use-status'
import { useSystemConfig } from '@/hooks/use-system-config'
import { useAuthStore } from '@/stores/auth-store'

import { FeaturesSection } from './components/features-section'
import { HeroSection } from './components/hero-section'
import { LandingFooter } from './components/landing-footer'
import { LandingHeader } from './components/landing-header'
import { ModelSection } from './components/model-section'
import { PlaygroundSection } from './components/playground-section'
import {
  useCounterAnimation,
  useIsScrolled,
  useRevealOnScroll,
} from './hooks/use-scroll-effects'
import { modelSections } from './lib/model-sections'

import './landing.css'

/**
 * Landing page ported from the Vue homepage plugin
 * (plugins/homepage/frontend). Kept visually identical on purpose, so it uses
 * the original's raw zinc palette rather than the app's semantic tokens.
 *
 * Structural note: sections must stay inside <main> and the footer must stay a
 * <footer>, because landing.css targets both in its mobile media queries.
 */
export function Landing() {
  const { resolvedTheme } = useTheme()
  const { status } = useStatus()
  const { systemName, logo } = useSystemConfig()
  const { auth } = useAuthStore()

  const isDark = resolvedTheme === 'dark'
  const isAuthenticated = !!auth.user
  const docUrl = (status?.docs_link as string | undefined) || ''
  const siteSubtitle = ''
  const navSubtitle = siteSubtitle || 'AI API Gateway'

  const isScrolled = useIsScrolled()
  useRevealOnScroll()
  useCounterAnimation()

  // An admin-configured home page wins over this landing page, matching both
  // the Vue original and upstream's features/home. Without this branch the
  // "Home Page Content" setting would silently do nothing.
  const { content, isLoaded, isUrl } = useHomePageContent()

  if (isLoaded && content) {
    if (isUrl) {
      return (
        <div className='min-h-screen'>
          {/*
            allow-top-navigation-by-user-activation: the custom home page URL is
            admin-configured (trusted); this lets its target="_top" links
            navigate the top-level window on user click. Matches upstream's
            features/home sandbox. Does NOT grant same-origin access.
          */}
          <iframe
            src={content.trim()}
            className='h-screen w-full border-0'
            title='Custom Home Page'
            sandbox='allow-forms allow-popups allow-popups-to-escape-sandbox allow-scripts allow-top-navigation-by-user-activation'
            allowFullScreen
          />
        </div>
      )
    }
    return (
      <div
        className='min-h-screen'
        // SECURITY: same trade-off as the Vue original and upstream — this is an
        // admin-only setting, so the XSS surface is accepted deliberately.
        // eslint-disable-next-line react/no-danger
        dangerouslySetInnerHTML={{ __html: content }}
      />
    )
  }

  return (
    <div className='home-shell min-h-screen font-sans text-zinc-900 antialiased transition-colors duration-300 selection:bg-zinc-200 dark:text-zinc-100 dark:selection:bg-zinc-800'>
      <LandingHeader
        isScrolled={isScrolled}
        isAuthenticated={isAuthenticated}
        siteName={systemName}
        siteLogo={logo}
        navSubtitle={navSubtitle}
        docUrl={docUrl}
      />

      <main className='relative'>
        <HeroSection
          isDark={isDark}
          isAuthenticated={isAuthenticated}
          siteSubtitle={siteSubtitle}
          docUrl={docUrl}
        />
        <FeaturesSection />
        <PlaygroundSection isAuthenticated={isAuthenticated} />
        {modelSections.map((config) => (
          <ModelSection key={config.id} config={config} />
        ))}
      </main>

      <LandingFooter
        siteName={systemName}
        navSubtitle={navSubtitle}
        docUrl={docUrl}
      />
    </div>
  )
}
