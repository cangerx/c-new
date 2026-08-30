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
import { ArrowRight, ChevronRight } from 'lucide-react'

import { useHeroShader } from '../hooks/use-hero-shader'

const DEFAULT_SUBTITLE =
  '专为开发者与企业打造的下一代 AI API 网关底座。统一聚合全球主流大模型，提供毫秒级智能调度、渠道灾备与精细计费。'

interface HeroSectionProps {
  isDark: boolean
  isAuthenticated: boolean
  siteSubtitle: string
  docUrl: string
}

export function HeroSection({
  isDark,
  isAuthenticated,
  siteSubtitle,
  docUrl,
}: HeroSectionProps) {
  const shaderRef = useHeroShader(isDark)

  return (
    <section className='hero-section relative flex min-h-screen flex-col items-center justify-center overflow-hidden px-4 text-center'>
      {/* Physical DOM wrapper keeps the WebGL canvas isolated to the hero. */}
      <div className='home-bg-overlay'>
        <div ref={shaderRef} className='hero-grain-shader' />
      </div>

      <div className='hero-gradient-scrim pointer-events-none z-0' />

      <div className='relative z-10 mx-auto flex max-w-4xl flex-col items-center pt-24 pb-16'>
        <div className='hero-badge animate-fade-up animate-float mb-8 inline-flex h-10 cursor-default items-center gap-2.5 rounded-full border border-zinc-200/50 bg-zinc-100/80 px-4 text-xs font-medium text-zinc-800 shadow-sm backdrop-blur-md transition-colors dark:border-white/[0.12] dark:bg-white/[0.04] dark:text-zinc-200 dark:hover:border-white/[0.2]'>
          <span className='relative flex h-1.5 w-1.5'>
            <span className='absolute inline-flex h-full w-full animate-ping rounded-full bg-blue-400 opacity-75' />
            <span className='relative inline-flex h-1.5 w-1.5 rounded-full bg-blue-500' />
          </span>
          <span>已全面适配 DeepSeek-R1 与视频生成全系列大模型</span>
        </div>

        <h1 className='hero-title mb-6 max-w-4xl text-4xl leading-[1.1] font-black tracking-[-0.03em] sm:text-6xl md:text-7xl lg:text-8xl'>
          <span className='hero-title-line hero-title-line-1'>
            <span>连接全球智能</span>
          </span>
          <br className='hidden sm:inline' />
          <span className='hero-title-line hero-title-line-2'>
            <span>赋能 API 应用创新</span>
          </span>
        </h1>

        <div className='hero-kicker animate-fade-up animate-fade-up-3 mb-8 flex items-center justify-center gap-2'>
          <span className='h-px w-6 bg-zinc-300 dark:bg-white/20' />
          <span className='text-xs font-semibold tracking-widest text-zinc-400 uppercase dark:text-zinc-500'>
            Build with models, create with agents
          </span>
          <span className='h-px w-6 bg-zinc-300 dark:bg-white/20' />
        </div>

        <p className='animate-fade-up animate-fade-up-3 mb-10 max-w-2xl text-sm leading-relaxed font-medium tracking-tight text-zinc-600 sm:text-base md:text-lg dark:text-white/75'>
          {siteSubtitle || DEFAULT_SUBTITLE}
        </p>

        <div className='hero-actions animate-fade-up animate-fade-up-4 flex flex-wrap items-center justify-center gap-4 sm:gap-6'>
          <a
            href={isAuthenticated ? '/dashboard' : '/sign-in'}
            className='shimmer-btn inline-flex h-[52px] items-center justify-center gap-1.5 overflow-hidden rounded-[60px] bg-zinc-950 px-10 text-xs font-semibold text-white shadow-md transition-all duration-200 hover:-translate-y-0.5 hover:bg-zinc-800 active:translate-y-0 dark:bg-white dark:text-zinc-950 dark:hover:bg-zinc-100'
          >
            <span className='relative z-10 flex items-center gap-1.5'>
              {isAuthenticated ? '进入控制台' : '立即开始'}
              <ArrowRight className='h-3 w-3' />
            </span>
          </a>

          {docUrl && (
            <a
              href={docUrl}
              target='_blank'
              rel='noopener noreferrer'
              className='inline-flex h-[52px] items-center justify-center gap-1.5 rounded-[60px] border border-zinc-200 bg-white/40 px-10 text-xs font-semibold text-zinc-700 shadow-sm backdrop-blur-md transition-all duration-200 hover:-translate-y-0.5 hover:bg-white/80 hover:text-zinc-950 active:translate-y-0 dark:border-white/[0.12] dark:bg-white/[0.04] dark:text-zinc-300 dark:hover:bg-white/[0.08] dark:hover:text-white'
            >
              <span>查看开发文档</span>
              <ChevronRight className='h-3 w-3' />
            </a>
          )}
        </div>
      </div>
    </section>
  )
}
