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
import { render, screen } from '@testing-library/react'
import type { ReactNode } from 'react'
import { beforeEach, describe, expect, test, vi } from 'vitest'

const mockState = vi.hoisted(() => ({
  homePage: {
    content: '',
    isLoaded: true,
    isUrl: false,
  },
}))

vi.mock('@/components/layout', () => ({
  PublicLayout: (props: { children: ReactNode }) => (
    <div data-testid='public-layout'>{props.children}</div>
  ),
}))

vi.mock('@/components/rich-content', () => ({
  RichContent: (props: { content: string; mode?: string }) => (
    <div data-testid='rich-content' data-mode={props.mode}>
      {props.content}
    </div>
  ),
}))

vi.mock('@/context/theme-provider', () => ({
  useTheme: () => ({ resolvedTheme: 'light', setTheme: vi.fn() }),
}))

vi.mock('@/features/home/hooks', () => ({
  useHomePageContent: () => mockState.homePage,
}))

vi.mock('@/hooks/use-status', () => ({
  useStatus: () => ({ status: { docs_link: 'https://docs.example.com' } }),
}))

vi.mock('@/hooks/use-system-config', () => ({
  useSystemConfig: () => ({ systemName: 'Test API', logo: '/logo.png' }),
}))

vi.mock('@/stores/auth-store', () => ({
  useAuthStore: () => ({ auth: { user: null } }),
}))

vi.mock('../components/features-section', () => ({
  FeaturesSection: () => <section data-testid='features-section' />,
}))

vi.mock('../components/hero-section', () => ({
  HeroSection: () => <section data-testid='hero-section' />,
}))

vi.mock('../components/landing-footer', () => ({
  LandingFooter: () => <footer data-testid='landing-footer' />,
}))

vi.mock('../components/landing-header', () => ({
  LandingHeader: () => <header data-testid='landing-header' />,
}))

vi.mock('../components/model-section', () => ({
  ModelSection: () => <section data-testid='model-section' />,
}))

vi.mock('../components/playground-section', () => ({
  PlaygroundSection: () => <section data-testid='playground-section' />,
}))

vi.mock('../hooks/use-scroll-effects', () => ({
  useCounterAnimation: vi.fn(),
  useIsScrolled: () => false,
  useRevealOnScroll: vi.fn(),
}))

const { Landing } = await import('../index')

describe('landing page routing behavior', () => {
  beforeEach(() => {
    mockState.homePage = {
      content: '',
      isLoaded: true,
      isUrl: false,
    }
  })

  test('shows the new landing page when no custom homepage is configured', () => {
    render(<Landing />)

    expect(screen.getByTestId('landing-header')).toBeInTheDocument()
    expect(screen.getByTestId('hero-section')).toBeInTheDocument()
    expect(screen.getByTestId('playground-section')).toBeInTheDocument()
    expect(screen.getByTestId('landing-footer')).toBeInTheDocument()
  })

  test('keeps configured Markdown content ahead of the default landing page', () => {
    mockState.homePage = {
      content: '# Private portal',
      isLoaded: true,
      isUrl: false,
    }

    render(<Landing />)

    expect(screen.getByTestId('rich-content')).toHaveAttribute(
      'data-mode',
      'markdown'
    )
    expect(screen.getByText('# Private portal')).toBeInTheDocument()
    expect(screen.queryByTestId('landing-header')).not.toBeInTheDocument()
  })

  test('keeps configured homepage URLs in the existing sandboxed iframe', () => {
    mockState.homePage = {
      content: 'https://portal.example.com',
      isLoaded: true,
      isUrl: true,
    }

    render(<Landing />)

    const iframe = screen.getByTitle('Custom Home Page')
    expect(iframe).toHaveAttribute('src', 'https://portal.example.com')
    expect(iframe).toHaveAttribute(
      'sandbox',
      'allow-forms allow-popups allow-popups-to-escape-sandbox allow-scripts allow-top-navigation-by-user-activation'
    )
    expect(screen.queryByTestId('landing-header')).not.toBeInTheDocument()
  })
})
