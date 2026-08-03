const CHIP_CLASS =
  'inline-flex items-center rounded-full border border-zinc-200/70 bg-white/55 px-4 py-2 text-[11px] font-semibold tracking-tight text-zinc-600 shadow-sm backdrop-blur-md dark:border-white/10 dark:bg-white/[0.04] dark:text-zinc-400'

const LINK_HOVER_CLASS =
  'transition-colors hover:border-zinc-300 hover:text-zinc-950 dark:hover:border-white/20 dark:hover:text-white'

interface LandingFooterProps {
  siteName: string
  navSubtitle: string
  docUrl: string
}

export function LandingFooter({
  siteName,
  navSubtitle,
  docUrl,
}: LandingFooterProps) {
  const currentYear = new Date().getFullYear()

  return (
    <footer className='border-t border-zinc-200/30 bg-white/45 px-4 py-10 backdrop-blur-xl dark:border-white/10 dark:bg-black/35'>
      <div className='mx-auto flex max-w-7xl flex-col items-center justify-between gap-5 text-center sm:flex-row sm:text-left'>
        <div className='space-y-1'>
          <p className='text-xs font-semibold tracking-tight text-zinc-700 dark:text-zinc-200'>
            &copy; {currentYear} {siteName}
          </p>
          <p className='max-w-xl text-[11px] leading-5 text-zinc-500 dark:text-zinc-500'>
            {navSubtitle} · 保留所有权利。
          </p>
        </div>
        <div className='flex flex-wrap items-center justify-center gap-2.5 text-xs font-semibold sm:justify-end'>
          <span className={CHIP_CLASS}>CCAI v1.0.0</span>
          <span className={`${CHIP_CLASS} gap-1.5`}>
            <span className='text-zinc-400 dark:text-zinc-500'>Author</span>
            <span className='text-zinc-800 dark:text-zinc-200'>canger</span>
          </span>
          <a
            href='https://t.me/cangerx'
            target='_blank'
            rel='noopener noreferrer'
            className={`${CHIP_CLASS} gap-1.5 ${LINK_HOVER_CLASS}`}
          >
            <span className='text-zinc-400 dark:text-zinc-500'>TG</span>
            <span>@cangerx</span>
          </a>
          {docUrl && (
            <a
              href={docUrl}
              target='_blank'
              rel='noopener noreferrer'
              className={`${CHIP_CLASS} ${LINK_HOVER_CLASS}`}
            >
              文档
            </a>
          )}
        </div>
      </div>
    </footer>
  )
}
