import {
  ArrowRight,
  ChartNoAxesColumn,
  Check,
  ChevronDown,
  MessageSquare,
  Terminal,
} from 'lucide-react'

const BULLETS = [
  '统一 Payload 格式，支持单请求多模型横向分发',
  '毫秒级首字延迟（TTFT）与 Token 生成速率可视化',
  '直接导出符合 OpenAI 规范的 cURL/Python/JS 代码',
]

interface PlaygroundSectionProps {
  isAuthenticated: boolean
}

/**
 * Static macOS-style console mockup. Every control here is decorative — the
 * Vue original had no event handlers in this block either.
 */
export function PlaygroundSection({ isAuthenticated }: PlaygroundSectionProps) {
  return (
    <section className='relative overflow-hidden border-t border-zinc-200/20 bg-zinc-50/50 py-24 dark:border-zinc-900/40 dark:bg-[#070708]/30'>
      <div className='pointer-events-none absolute -top-40 -right-40 h-96 w-96 rounded-full bg-blue-500/10 blur-[100px] dark:bg-blue-600/5' />
      <div className='pointer-events-none absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-purple-500/10 blur-[100px] dark:bg-purple-600/5' />

      <div className='mx-auto max-w-7xl px-4'>
        <div className='grid grid-cols-1 items-center gap-12 lg:grid-cols-12'>
          <div className='reveal-element space-y-6 text-left lg:col-span-5'>
            <div className='inline-flex items-center gap-1.5 rounded-full border border-blue-200/30 bg-blue-50 px-3 py-1 text-[11px] font-bold text-blue-600 dark:border-blue-800/30 dark:bg-blue-950/30 dark:text-blue-400'>
              <Terminal className='h-3 w-3' />
              <span>MODEL PLAYGROUND</span>
            </div>

            <h2 className='text-3xl leading-[1.15] font-black tracking-tight text-zinc-950 sm:text-4xl lg:text-5xl dark:text-white'>
              无需编写任何代码
              <br />
              直接在浏览器调试
            </h2>

            <p className='text-sm leading-relaxed font-medium tracking-tight text-zinc-600 dark:text-zinc-400'>
              我们为开发者与企业打造了沉浸式的大模型体验中心。您可以在可视化控制台中一键切换、对比多款旗舰模型，观测实时延迟、吞吐量和运行成本，快速筛选最适配业务场景的智能底座。
            </p>

            <ul className='space-y-3 pt-2 text-xs font-semibold text-zinc-700 sm:text-sm dark:text-zinc-300'>
              {BULLETS.map((bullet) => (
                <li key={bullet} className='flex items-center gap-2.5'>
                  <span className='flex h-5 w-5 shrink-0 items-center justify-center rounded-full border border-blue-200/20 bg-blue-50 text-blue-600 dark:bg-blue-950/30 dark:text-blue-400'>
                    <Check className='h-3 w-3' />
                  </span>
                  <span>{bullet}</span>
                </li>
              ))}
            </ul>

            <div className='flex flex-wrap gap-4 pt-4'>
              <a
                href={isAuthenticated ? '/dashboard' : '/sign-in'}
                className='inline-flex h-[52px] items-center justify-center gap-1.5 rounded-[60px] bg-zinc-950 px-8 text-xs font-bold text-white shadow-md transition-all duration-200 hover:-translate-y-0.5 hover:bg-zinc-800 active:translate-y-0 dark:bg-white dark:text-zinc-950 dark:hover:bg-zinc-100'
              >
                <span>立即进入体验中心</span>
                <ArrowRight className='h-3 w-3' />
              </a>
            </div>
          </div>

          <div className='reveal-element delay-100 w-full lg:col-span-7'>
            <div className='playground-console flex h-[520px] w-full flex-col overflow-hidden rounded-2xl border border-zinc-200/80 bg-white/70 text-left shadow-2xl backdrop-blur-md dark:border-white/[0.08] dark:bg-[#0b0b0d]/80'>
              <div className='playground-console-header flex h-12 shrink-0 items-center justify-between border-b border-zinc-200/60 bg-zinc-100/50 px-4 select-none dark:border-white/[0.06] dark:bg-[#121215]/50'>
                <div className='flex w-16 items-center gap-1.5'>
                  <span className='h-3 w-3 rounded-full border border-[#e0443e] bg-[#ff5f56]' />
                  <span className='h-3 w-3 rounded-full border border-[#dea123] bg-[#ffbd2e]' />
                  <span className='h-3 w-3 rounded-full border border-[#1aab29] bg-[#27c93f]' />
                </div>
                <div className='playground-tabs flex items-center gap-1 text-[11px] font-semibold text-zinc-500 dark:text-zinc-400'>
                  <span className='flex items-center gap-1.5 rounded-md border border-zinc-200/40 bg-white px-3 py-1.5 text-zinc-800 dark:border-white/[0.04] dark:bg-[#1c1c21] dark:text-white'>
                    <MessageSquare className='h-3 w-3 text-blue-500' /> Chat
                    Playground
                  </span>
                  <span className='flex cursor-pointer items-center gap-1.5 rounded-md px-3 py-1.5 transition-colors hover:bg-zinc-200/40 dark:hover:bg-white/[0.03]'>
                    <Terminal className='h-3 w-3' /> cURL
                  </span>
                  <span className='hidden cursor-pointer items-center gap-1.5 rounded-md px-3 py-1.5 transition-colors hover:bg-zinc-200/40 sm:inline-flex dark:hover:bg-white/[0.03]'>
                    <ChartNoAxesColumn className='h-3 w-3' /> Metrics
                  </span>
                </div>
                <div className='playground-status flex items-center gap-1.5 rounded-md border border-blue-500/20 bg-blue-500/10 px-2 py-0.5 text-[10px] font-bold text-blue-500 dark:text-blue-400'>
                  <span className='h-1.5 w-1.5 animate-pulse rounded-full bg-blue-500' />
                  <span>ONLINE</span>
                </div>
              </div>

              <div className='flex min-h-0 flex-1 divide-x divide-zinc-200/60 overflow-hidden dark:divide-white/[0.06]'>
                <div className='hidden w-56 shrink-0 flex-col space-y-5 overflow-y-auto bg-zinc-50/30 p-4 sm:flex dark:bg-[#09090b]/40'>
                  <div className='space-y-1.5'>
                    <label className='text-[10px] font-bold tracking-wider text-zinc-400 uppercase'>
                      Target Model
                    </label>
                    <div className='flex h-9 cursor-pointer items-center justify-between rounded-lg border border-zinc-200/60 bg-white px-3 text-xs font-semibold text-zinc-800 shadow-sm dark:border-white/[0.08] dark:bg-[#16161a] dark:text-zinc-200 dark:hover:border-white/[0.15]'>
                      <span className='flex items-center gap-1.5'>
                        <span className='h-2 w-2 rounded-full bg-blue-500' />
                        DeepSeek-R1
                      </span>
                      <ChevronDown className='h-3 w-3 opacity-60' />
                    </div>
                  </div>

                  <div className='space-y-2'>
                    <div className='flex items-center justify-between text-[10px] font-bold tracking-wider text-zinc-400 uppercase'>
                      <span>Temperature</span>
                      <span className='font-mono text-blue-500'>0.70</span>
                    </div>
                    <div className='relative h-1.5 w-full rounded-full bg-zinc-200 dark:bg-zinc-800'>
                      <div className='absolute top-0 left-0 h-full w-[70%] rounded-full bg-blue-500' />
                      <div className='absolute top-1/2 left-[70%] h-3.5 w-3.5 -translate-x-1/2 -translate-y-1/2 cursor-pointer rounded-full border border-zinc-300 bg-white shadow' />
                    </div>
                  </div>

                  <div className='space-y-1.5'>
                    <label className='text-[10px] font-bold tracking-wider text-zinc-400 uppercase'>
                      System Prompt
                    </label>
                    <div className='h-32 overflow-hidden rounded-lg border border-zinc-200/60 bg-white p-2.5 text-[11px] leading-relaxed font-medium text-zinc-500 shadow-inner dark:border-white/[0.08] dark:bg-[#16161a]'>
                      You are a professional software architect. Explain how API
                      gateways optimize latency and channel load balancing.
                    </div>
                  </div>
                </div>

                <PlaygroundChat />
              </div>

              <div className='playground-console-footer flex h-8 shrink-0 items-center justify-between border-t border-zinc-200/60 bg-zinc-100/40 px-4 text-[10px] font-bold text-zinc-400 select-none dark:border-white/[0.06] dark:bg-[#0b0b0e]/80 dark:text-zinc-500'>
                <div className='flex items-center gap-3'>
                  <span className='flex items-center gap-1.5'>
                    <span className='h-1.5 w-1.5 rounded-full bg-blue-500' />{' '}
                    Latency: 12ms
                  </span>
                  <span className='hidden items-center gap-1.5 sm:inline-flex'>
                    <span className='h-1.5 w-1.5 rounded-full bg-blue-500' />{' '}
                    Uptime: 99.99%
                  </span>
                </div>
                <div className='flex items-center gap-3'>
                  <span>1,824 Requests/min</span>
                  <span>HTTPS / API v1</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  )
}

function PlaygroundChat() {
  return (
    <div className='flex min-w-0 flex-1 flex-col overflow-hidden bg-white/40 dark:bg-[#070709]/20'>
      <div className='scrollbar-thin flex-1 space-y-4 overflow-y-auto p-4 text-xs font-semibold'>
        <div className='flex max-w-[85%] items-start gap-2.5'>
          <div className='flex h-6 w-6 shrink-0 items-center justify-center rounded-full border border-zinc-200/40 bg-zinc-100 text-[10px] dark:border-white/[0.04] dark:bg-zinc-800'>
            U
          </div>
          <div className='rounded-2xl rounded-tl-none bg-zinc-100 px-3.5 py-2.5 leading-relaxed font-medium text-zinc-800 shadow-sm dark:bg-[#1a1a20] dark:text-zinc-200'>
            CCAPI 平台如何保障多渠道大模型的延迟与稳定性？
          </div>
        </div>

        <div className='ml-auto flex max-w-[90%] flex-row-reverse items-start gap-2.5'>
          <div className='flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-blue-600 text-[10px] text-white'>
            AI
          </div>
          <div className='w-full space-y-2.5'>
            <div className='flex flex-wrap items-center justify-end gap-1.5'>
              <span className='rounded-md border border-blue-500/20 bg-blue-500/10 px-2 py-0.5 text-[9px] font-bold text-blue-500 dark:text-blue-400'>
                DeepSeek-R1
              </span>
              <span className='rounded-md border border-zinc-200/30 bg-zinc-100 px-2 py-0.5 font-mono text-[9px] text-zinc-500 dark:border-white/[0.04] dark:bg-zinc-800/80 dark:text-zinc-400'>
                TTFT: 14ms
              </span>
              <span className='rounded-md border border-zinc-200/30 bg-zinc-100 px-2 py-0.5 font-mono text-[9px] text-zinc-500 dark:border-white/[0.04] dark:bg-zinc-800/80 dark:text-zinc-400'>
                120 tok/s
              </span>
              <span className='rounded-md border border-blue-500/20 bg-blue-500/10 px-2 py-0.5 font-mono text-[9px] text-blue-500 dark:text-blue-400'>
                节省 80% 成本
              </span>
            </div>

            <div className='space-y-2 rounded-2xl rounded-tr-none border border-blue-100/30 bg-blue-50/40 px-3.5 py-2.5 leading-relaxed font-normal text-zinc-800 shadow-sm dark:border-blue-900/10 dark:bg-blue-950/15 dark:text-zinc-200'>
              <p className='text-xs font-medium text-blue-600 dark:text-blue-400'>
                CCAPI 通过三大底层核心技术确保旗舰级稳定性：
              </p>
              <ol className='list-decimal space-y-1 pl-4 text-[11px] text-zinc-600 dark:text-zinc-400'>
                <li>
                  <strong className='text-zinc-800 dark:text-zinc-200'>
                    毫秒级智能测速路由
                  </strong>
                  ：实时监控全球百余个核心节点的延迟，自动为用户分发响应最快的通道。
                </li>
                <li>
                  <strong className='text-zinc-800 dark:text-zinc-200'>
                    多渠道自动无感灾备
                  </strong>
                  ：某上游服务商（如 OpenAI）发生拥堵或限流时，网关在 0
                  毫秒内自动热重试其他备份提供商，外部访问完全不中断。
                </li>
                <li>
                  <strong className='text-zinc-800 dark:text-zinc-200'>
                    多模型高并发调度引擎
                  </strong>
                  ：自研的并发控制与缓冲队列，能够轻松承载百万 QPS 级高频调用。
                </li>
              </ol>
              <div className='flex items-center gap-1 border-t border-blue-200/20 pt-1.5 text-[10px] text-zinc-400 dark:border-blue-900/10 dark:text-zinc-500'>
                <span className='h-1.5 w-1.5 animate-ping rounded-full bg-blue-500' />
                <span>回答已生成完毕 (Total Tokens: 348)</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className='flex shrink-0 items-center gap-2 border-t border-zinc-200/60 bg-zinc-50/40 p-3.5 dark:border-white/[0.06] dark:bg-[#0c0c0f]/60'>
        <div className='flex h-9 flex-1 cursor-text items-center justify-between rounded-full border border-zinc-200/65 bg-white px-3 text-xs text-zinc-400 shadow-sm dark:border-white/[0.08] dark:bg-[#141417] dark:hover:border-white/[0.15]'>
          <span className='flex items-center truncate font-medium'>
            输入 Prompt，例如"用 Go 写一个高性能 API 转发代理"
            <span className='typing-cursor'>|</span>
          </span>
          <MessageSquare className='h-3 w-3 opacity-40' />
        </div>
        <button
          type='button'
          className='flex h-9 w-9 items-center justify-center rounded-full bg-zinc-950 text-white shadow-md transition-transform active:scale-95 dark:bg-white dark:text-black'
        >
          <ArrowRight className='h-3 w-3' />
        </button>
      </div>
    </div>
  )
}
