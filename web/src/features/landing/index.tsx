/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useCallback, useEffect, useRef } from 'react'
import { useTranslation } from 'react-i18next'

import { PublicLayout } from '@/components/layout'
import { RichContent } from '@/components/rich-content'
import { useTheme } from '@/context/theme-provider'
import { useHomePageContent } from '@/features/home/hooks'
import { useStatus } from '@/hooks/use-status'
import { useSystemConfig } from '@/hooks/use-system-config'
import { isLikelyHtml } from '@/lib/content-format'
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
  const { i18n, t } = useTranslation()
  const iframeRef = useRef<HTMLIFrameElement>(null)
  const { resolvedTheme } = useTheme()
  const { status } = useStatus()
  const { systemName, logo } = useSystemConfig()
  const { auth } = useAuthStore()

  const isDark = resolvedTheme === 'dark'
  const isAuthenticated = !!auth.user
  const docUrl = (status?.docs_link as string | undefined) || ''
  const siteSubtitle = ''
  const navSubtitle = siteSubtitle || 'AI API Gateway'
  const { content, isLoaded, isUrl } = useHomePageContent()
  const showDefaultLanding = isLoaded && !content

  const isScrolled = useIsScrolled()
  useRevealOnScroll(showDefaultLanding)
  useCounterAnimation(showDefaultLanding)

  // An admin-configured home page wins over this landing page, matching both
  // the Vue original and upstream's features/home. Without this branch the
  // "Home Page Content" setting would silently do nothing.
  const syncIframePreferences = useCallback(() => {
    try {
      iframeRef.current?.contentWindow?.postMessage(
        { themeMode: resolvedTheme },
        '*'
      )
      iframeRef.current?.contentWindow?.postMessage(
        { lang: i18n.language },
        '*'
      )
    } catch {
      // Cross-origin frames may reject access while navigating.
    }
  }, [i18n.language, resolvedTheme])

  useEffect(() => {
    if (isUrl) {
      syncIframePreferences()
    }
  }, [isUrl, syncIframePreferences])

  if (!isLoaded) {
    return (
      <PublicLayout showMainContainer={false}>
        <main className='flex min-h-screen items-center justify-center'>
          <div className='text-muted-foreground'>{t('Loading...')}</div>
        </main>
      </PublicLayout>
    )
  }

  if (content) {
    if (isUrl) {
      return (
        <PublicLayout showMainContainer={false}>
          {/*
            allow-top-navigation-by-user-activation: the custom home page URL is
            admin-configured (trusted); this lets its target="_top" links
            navigate the top-level window on user click. This token does not
            grant same-origin access.
          */}
          <iframe
            ref={iframeRef}
            src={content}
            className='h-screen w-full border-none'
            title={t('Custom Home Page')}
            sandbox='allow-forms allow-popups allow-popups-to-escape-sandbox allow-scripts allow-top-navigation-by-user-activation'
            allowFullScreen
            onLoad={syncIframePreferences}
          />
        </PublicLayout>
      )
    }

    const contentIsHtml = isLikelyHtml(content)

    if (contentIsHtml) {
      return (
        <PublicLayout showMainContainer={false}>
          <RichContent
            mode='html'
            htmlVariant='isolated'
            content={content}
            className='custom-home-content'
          />
        </PublicLayout>
      )
    }

    return (
      <PublicLayout>
        <div className='mx-auto max-w-6xl px-4 py-8'>
          <RichContent
            mode='markdown'
            content={content}
            className='custom-home-content'
          />
        </div>
      </PublicLayout>
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
