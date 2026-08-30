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
import { BookOpen, Moon, Sun } from 'lucide-react'

import { LanguageSwitcher } from '@/components/language-switcher'
import { useTheme } from '@/context/theme-provider'
import { cn } from '@/lib/utils'

const NAV_LINKS = [
  { href: '#', label: '首页' },
  { href: '#features', label: '核心优势' },
  { href: '#domestic-models', label: '国产模型' },
  { href: '#image-models', label: '图片模型' },
  { href: '#video-models', label: '视频模型' },
] as const

interface LandingHeaderProps {
  isScrolled: boolean
  isAuthenticated: boolean
  siteName: string
  siteLogo: string
  navSubtitle: string
  docUrl: string
}

export function LandingHeader({
  isScrolled,
  isAuthenticated,
  siteName,
  siteLogo,
  navSubtitle,
  docUrl,
}: LandingHeaderProps) {
  const { resolvedTheme, setTheme } = useTheme()
  const isDark = resolvedTheme === 'dark'

  return (
    <header
      className={cn(
        'glass-header pointer-events-none fixed inset-x-0 top-0 z-40 transition-all duration-500',
        isScrolled ? 'px-4 pt-4' : 'px-0 pt-0'
      )}
    >
      <nav
        className={cn(
          'liquid-nav pointer-events-auto mx-auto flex items-center justify-between transition-all duration-500',
          isScrolled
            ? 'liquid-nav-compact max-w-5xl scale-95 rounded-full px-6 py-2 md:scale-100'
            : 'liquid-nav-expanded w-full max-w-7xl rounded-none px-4 py-5 sm:px-6'
        )}
      >
        <div
          className={cn(
            'flex min-w-0 items-center gap-3 transition-all duration-500',
            isScrolled
              ? 'translate-x-1 scale-95 md:scale-100'
              : 'translate-x-0 scale-100'
          )}
        >
          <div className='flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-zinc-200/20 bg-white/10 shadow-sm dark:border-white/10'>
            <img
              src={siteLogo || '/logo.svg'}
              alt='Logo'
              className='h-full w-full object-contain p-1'
            />
          </div>
          <div className='min-w-0'>
            <div className='truncate text-sm leading-5 font-semibold tracking-tight text-zinc-900 dark:text-white'>
              {siteName}
            </div>
            <div className='hidden max-w-[180px] truncate text-[10px] font-medium tracking-tight text-zinc-500 sm:block dark:text-zinc-400'>
              {navSubtitle}
            </div>
          </div>
        </div>

        <div
          className={cn(
            'liquid-menu liquid-menu-capsule hidden items-center gap-10 tracking-tight transition-all duration-500 md:flex',
            isScrolled
              ? 'rounded-full px-7 py-2.5 text-xs text-zinc-700 dark:text-zinc-300'
              : 'rounded-full px-9 py-3 text-sm font-semibold text-zinc-700 dark:text-zinc-300'
          )}
        >
          {NAV_LINKS.map((link) => (
            <a
              key={link.label}
              href={link.href}
              className='transition-colors hover:text-zinc-900 dark:hover:text-white'
            >
              {link.label}
            </a>
          ))}
          {docUrl && (
            <a
              href={docUrl}
              target='_blank'
              rel='noopener noreferrer'
              className='transition-colors hover:text-zinc-900 dark:hover:text-white'
            >
              开发文档
            </a>
          )}
        </div>

        <div
          className={cn(
            'flex items-center gap-1.5 transition-all duration-500 sm:gap-2.5',
            isScrolled
              ? '-translate-x-1 scale-95 md:scale-100'
              : 'translate-x-0 scale-100'
          )}
        >
          <div className='shrink-0'>
            <LanguageSwitcher />
          </div>
          {docUrl && (
            <a
              href={docUrl}
              target='_blank'
              rel='noopener noreferrer'
              className='nav-icon-btn hidden shrink-0 text-zinc-600 hover:text-zinc-900 sm:inline-flex dark:text-zinc-400 dark:hover:text-white'
              title='查看文档'
            >
              <BookOpen className='h-4 w-4' />
            </a>
          )}
          <button
            type='button'
            className='nav-icon-btn hidden shrink-0 text-zinc-600 hover:text-zinc-900 sm:inline-flex dark:text-zinc-400 dark:hover:text-white'
            title={isDark ? '切换到浅色模式' : '切换到深色模式'}
            onClick={() => setTheme(isDark ? 'light' : 'dark')}
          >
            {isDark ? (
              <Sun className='h-4 w-4' />
            ) : (
              <Moon className='h-4 w-4' />
            )}
          </button>

          <a
            href={isAuthenticated ? '/dashboard' : '/sign-in'}
            className={cn(
              'liquid-action inline-flex h-9 items-center rounded-full border px-5 text-xs font-semibold whitespace-nowrap backdrop-blur-md transition-all duration-500 active:scale-95',
              isScrolled
                ? 'liquid-action-solid text-white dark:text-zinc-950'
                : 'liquid-action-glass text-zinc-800 dark:text-white'
            )}
          >
            {isAuthenticated ? '控制台' : '登录'}
          </a>
        </div>
      </nav>
    </header>
  )
}
