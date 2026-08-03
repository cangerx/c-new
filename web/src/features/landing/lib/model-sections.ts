import type { ModelSectionConfig } from '../components/model-section'

import { domesticModels, imageModels, videoModels } from './models'

export const modelSections: ModelSectionConfig[] = [
  {
    id: 'domestic-models',
    pill: '文本大模型',
    heading: '国产旗舰大语言模型矩阵',
    subheading:
      '聚合国内最强推理与长文本模型，全链路高速调度，提供稳定可靠的底层支撑体系。',
    bannerClass:
      'bg-gradient-to-br from-blue-950 via-slate-900 to-indigo-950',
    bannerMutedClass: 'text-blue-100/70',
    showHatch: true,
    highlights: [
      {
        title: '100万长文本处理',
        desc: '深度优化百万字长文本理解和精准提炼能力',
      },
      {
        title: '20% 推理能力提升',
        desc: '相比传统接口，逻辑分析与数学推理能力大幅增强',
      },
      {
        title: '强大的 Agent 编排',
        desc: '支持高频逻辑推理任务与复杂的工具调用链',
      },
    ],
    stats: [
      { target: 99.9, decimals: 1, suffix: '%', label: 'SLA 稳定性保障' },
      { target: 100, decimals: 0, suffix: '+', label: '核心并发节点' },
      {
        prefix: '<',
        target: 50,
        decimals: 0,
        suffix: 'ms',
        label: '平均网络延迟',
      },
    ],
    cardAccentClass:
      'bg-indigo-50 text-indigo-600 border-indigo-100/30 dark:bg-indigo-950/40 dark:text-indigo-400 dark:border-indigo-800/30',
    cardHoverShadow: 'dark:hover:shadow-[0_0_25px_rgba(99,102,241,0.15)]',
    models: domesticModels,
  },
  {
    id: 'image-models',
    pill: '图像大模型',
    heading: '高保真艺术与写实图像生成',
    subheading:
      '丰富的提示词控制，满足全场景艺术创作与专业级商业视觉设计。',
    bannerClass:
      'bg-gradient-to-br from-fuchsia-950 via-zinc-900 to-purple-950',
    bannerMutedClass: 'text-purple-100/70',
    highlights: [
      { title: '电影级写实与光影', desc: '极强的照片写实度，细节纹理完美呈现' },
      {
        title: '高精准的指令遵循',
        desc: '完美响应复杂的文本描述与多重元素构图',
      },
      {
        title: 'ControlNet 深度控制',
        desc: '支持动作、线稿、深度图等条件生成限制',
      },
    ],
    stats: [
      {
        target: 4,
        decimals: 0,
        suffix: 'K',
        label: '超高分辨率输出',
        static: true,
      },
      {
        prefix: '<',
        target: 2,
        decimals: 0,
        suffix: 's',
        label: '单图毫秒级出图',
      },
      { target: 50, decimals: 0, suffix: '+', label: '风格化微调模型' },
    ],
    cardAccentClass:
      'bg-purple-50 text-purple-600 border-purple-100/30 dark:bg-purple-950/40 dark:text-purple-400 dark:border-purple-800/30',
    cardHoverShadow: 'dark:hover:shadow-[0_0_25px_rgba(168,85,247,0.15)]',
    models: imageModels,
  },
  {
    id: 'video-models',
    pill: '视频大模型',
    heading: '影视级物理时空视频生成',
    subheading:
      '统一支持主流文生视频与图生视频服务，突破想象力限制，重塑视觉边界。',
    bannerClass:
      'border border-zinc-800 bg-gradient-to-br from-zinc-900 via-[#111111] to-stone-900',
    bannerMutedClass: 'text-zinc-400',
    highlights: [
      {
        title: '原生物理世界模拟',
        desc: '极佳的三维空间一致性与复杂的物理引擎真实度',
      },
      {
        title: '高分辨率多帧渲染',
        desc: '支持 1080P 超高清晰度与 60FPS 丝滑帧率输出',
      },
      {
        title: '超长运镜与场景延伸',
        desc: '提供长达 60s+ 的影视级镜头语言与多场景转场',
      },
    ],
    stats: [
      { target: 60, decimals: 0, suffix: 's', label: '超长视频序列' },
      { target: 1080, decimals: 0, suffix: 'P', label: '原生生成分辨率' },
      { target: 99, decimals: 0, suffix: '%', label: '时间线物理一致性' },
    ],
    cardAccentClass:
      'bg-pink-50 text-pink-600 border-pink-100/30 dark:bg-pink-950/40 dark:text-pink-400 dark:border-pink-800/30',
    cardHoverShadow: 'dark:hover:shadow-[0_0_25px_rgba(236,72,153,0.15)]',
    models: videoModels,
  },
]
