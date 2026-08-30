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
import { ArrowUpDown, Server, Shield, type LucideIcon } from 'lucide-react'

interface FeatureCard {
  icon: LucideIcon
  title: string
  desc: string
  /** Per-card accent, kept inline so each card keeps its own hover glow. */
  iconClass: string
  hoverShadow: string
}

const CARDS: FeatureCard[] = [
  {
    icon: Server,
    title: '统一聚合 简易集成',
    desc: '统一多种主流大模型协议。支持流式传输与函数调用，通过单一 API 密钥，在几行代码内接入全球顶尖 AI 能力。',
    iconClass:
      'bg-blue-50 text-blue-600 border-blue-100/30 dark:bg-blue-950/40 dark:text-blue-400 dark:border-blue-800/30',
    hoverShadow: 'dark:hover:shadow-[0_0_25px_rgba(59,130,246,0.1)]',
  },
  {
    icon: ArrowUpDown,
    title: '智能调度 稳定灾备',
    desc: '基于并发数、可用性与耗时自动分发。首字延迟深度优化，支持故障自动重试与上游负载均衡，保障服务始终在线。',
    iconClass:
      'bg-blue-50 text-blue-600 border-blue-100/30 dark:bg-blue-950/40 dark:text-blue-400 dark:border-blue-800/30',
    hoverShadow: 'dark:hover:shadow-[0_0_25px_rgba(16,185,129,0.1)]',
  },
  {
    icon: Shield,
    title: '精细管控 安全合规',
    desc: '支持密钥级额度限制、RPM/TPM 频率控制与 Token 细粒度审计。内置合规治理模块，有效规避滥用风险。',
    iconClass:
      'bg-violet-50 text-violet-600 border-violet-100/30 dark:bg-violet-950/40 dark:text-violet-400 dark:border-violet-800/30',
    hoverShadow: 'dark:hover:shadow-[0_0_25px_rgba(139,92,246,0.1)]',
  },
]

export function FeaturesSection() {
  return (
    <section
      id='features'
      className='relative border-t border-zinc-200/20 py-24 dark:border-zinc-900/40'
    >
      <div className='mx-auto max-w-7xl px-4'>
        <div className='reveal-element section-heading mx-auto mb-16 max-w-2xl space-y-4 text-center'>
          <p className='eyebrow text-[10px] font-bold tracking-widest text-zinc-400 uppercase dark:text-zinc-500'>
            Core Features
          </p>
          <h2 className='text-3xl leading-tight font-bold tracking-tight text-zinc-900 sm:text-4xl lg:text-5xl dark:text-white'>
            卓越网关性能，赋能业务增长
          </h2>
          <p className='text-sm tracking-tight text-zinc-500 dark:text-zinc-400'>
            为您提供大模型集成与交付的一站式基础设施服务。
          </p>
        </div>

        <div className='reveal-element grid grid-cols-1 gap-6 delay-100 md:grid-cols-3 lg:gap-8'>
          {CARDS.map((card) => {
            const Icon = card.icon
            return (
              <div
                key={card.title}
                className={`feature-card rounded-3xl border border-zinc-200/50 bg-white/50 p-8 shadow-sm backdrop-blur-sm transition-all duration-300 hover:scale-[1.01] hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:border-zinc-800 dark:bg-zinc-950/40 ${card.hoverShadow}`}
              >
                <div
                  className={`mb-6 inline-flex rounded-2xl border p-3 shadow-sm ${card.iconClass}`}
                >
                  <Icon className='h-5 w-5' />
                </div>
                <h3 className='mb-3 text-lg font-bold text-zinc-900 dark:text-white'>
                  {card.title}
                </h3>
                <p className='text-xs leading-relaxed tracking-tight text-zinc-600 sm:text-sm dark:text-zinc-400'>
                  {card.desc}
                </p>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
