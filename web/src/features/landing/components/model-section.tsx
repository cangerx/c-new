import { ArrowRight } from 'lucide-react'

import type { ModelEntry } from '../lib/models'

interface BannerHighlight {
  title: string
  desc: string
}

interface BannerStat {
  /** Rendered before the counter, e.g. the "<" in "<50ms". */
  prefix?: string
  target: number
  decimals: number
  /** Rendered after the counter, e.g. "%", "+", "ms". */
  suffix: string
  label: string
  /**
   * Render the number as plain text instead of a counting animation. The Vue
   * original left "4K" static, so animating it would show "0K" on first paint.
   */
  static?: boolean
}

export interface ModelSectionConfig {
  id: string
  pill: string
  heading: string
  subheading: string
  /** Gradient + optional border for the dark promo banner. */
  bannerClass: string
  /** Muted text colour used by banner copy and stat labels. */
  bannerMutedClass: string
  highlights: BannerHighlight[]
  stats: BannerStat[]
  /** Accent applied to each model card's icon chip and tag. */
  cardAccentClass: string
  cardHoverShadow: string
  models: readonly ModelEntry[]
  /** Diagonal hatch overlay is only used by the text-model banner. */
  showHatch?: boolean
}

export function ModelSection({ config }: { config: ModelSectionConfig }) {
  return (
    <section
      id={config.id}
      className='relative border-t border-zinc-200/20 py-24 dark:border-zinc-900/40'
    >
      <div className='mx-auto max-w-7xl px-4'>
        <div className='reveal-element mb-8 flex flex-col items-start justify-between gap-4 md:flex-row md:items-end'>
          <div className='inline-flex cursor-default items-center gap-2 border border-zinc-900 px-6 py-2 text-sm font-semibold tracking-wide text-zinc-900 dark:border-white dark:text-white'>
            {config.pill} <ArrowRight className='h-3 w-3' />
          </div>
          <div className='text-left md:text-right'>
            <h2 className='mb-2 text-3xl leading-tight font-bold tracking-tight text-zinc-900 sm:text-4xl dark:text-white'>
              {config.heading}
            </h2>
            <p className='max-w-xl text-sm tracking-tight text-zinc-500 md:ml-auto dark:text-zinc-400'>
              {config.subheading}
            </p>
          </div>
        </div>

        <div
          className={`model-banner reveal-element delay-100 relative mb-12 flex w-full flex-col justify-between overflow-hidden rounded-2xl p-8 text-white shadow-xl md:p-12 ${config.bannerClass}`}
        >
          {config.showHatch && (
            <div className='absolute inset-0 bg-[linear-gradient(45deg,transparent_48%,rgba(255,255,255,0.8)_50%,transparent_52%)] bg-[length:40px_40px] opacity-10' />
          )}

          <div className='relative z-10 mb-12 flex max-w-sm flex-col gap-6 md:mb-16'>
            {config.highlights.map((item) => (
              <div key={item.title}>
                <div className='mb-1 text-xl font-bold md:text-2xl'>
                  {item.title}
                </div>
                <div className={`text-xs ${config.bannerMutedClass}`}>
                  {item.desc}
                </div>
              </div>
            ))}
          </div>

          <div className='counter-section relative z-10 mt-4 flex flex-wrap items-end gap-8 border-t border-white/10 pt-8 md:gap-16'>
            {config.stats.map((stat) => (
              <div key={stat.label}>
                <div className='mb-1 text-3xl font-black tracking-tighter md:text-5xl'>
                  {stat.prefix}
                  {stat.static ? (
                    <>
                      {stat.target}
                      {stat.suffix}
                    </>
                  ) : (
                    <>
                      <span
                        className='counter-value'
                        data-target={stat.target}
                        data-decimals={stat.decimals}
                      >
                        0
                      </span>
                      <span className='text-2xl'>{stat.suffix}</span>
                    </>
                  )}
                </div>
                <div
                  className={`text-[10px] font-bold tracking-widest uppercase ${config.bannerMutedClass}`}
                >
                  {stat.label}
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className='reveal-element delay-100 grid grid-cols-1 gap-6 sm:grid-cols-2 md:grid-cols-5'>
          {config.models.map((model, idx) => {
            const Icon = model.icon
            return (
              <div
                key={model.name}
                className={`card-stagger feature-card flex flex-col justify-between rounded-2xl border border-zinc-200/50 bg-white/50 p-6 shadow-sm backdrop-blur-sm transition-all duration-300 hover:scale-[1.02] hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:border-zinc-800 dark:bg-zinc-950/40 ${config.cardHoverShadow}`}
                style={{ animationDelay: `${idx * 0.08}s` }}
              >
                <div>
                  <div className='mb-4 flex items-center justify-between'>
                    <div
                      className={`inline-flex rounded-xl border p-2 ${config.cardAccentClass}`}
                    >
                      <Icon className='h-4 w-4' />
                    </div>
                    <span
                      className={`rounded-md border px-2 py-0.5 text-[10px] font-bold ${config.cardAccentClass}`}
                    >
                      {model.tag}
                    </span>
                  </div>
                  <h3 className='mb-2 text-base font-bold text-zinc-950 dark:text-white'>
                    {model.name}
                  </h3>
                  <p className='text-xs leading-relaxed tracking-tight text-zinc-500 dark:text-zinc-400'>
                    {model.desc}
                  </p>
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </section>
  )
}
