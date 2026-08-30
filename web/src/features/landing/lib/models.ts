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
import {
  Bolt,
  Box,
  Brain,
  Cloud,
  Cpu,
  Eye,
  Globe,
  Grid2x2,
  Play,
  Sparkles,
  Sun,
  TrendingUp,
  Video,
  type LucideIcon,
} from 'lucide-react'

export interface ModelEntry {
  name: string
  tag: string
  desc: string
  icon: LucideIcon
}

/**
 * Showcase copy, carried over verbatim from the Vue landing page. These are
 * marketing entries, deliberately independent of the channels configured in
 * admin — editing the list is a code change.
 */
export const domesticModels: readonly ModelEntry[] = [
  {
    name: 'DeepSeek-R1',
    tag: '最强推理',
    desc: '深度思考，逻辑推理与代码顶峰模型',
    icon: Brain,
  },
  {
    name: '通义千问 Qwen-Max',
    tag: '旗舰模型',
    desc: '全能基座大模型，极佳的长文本处理能力',
    icon: Sparkles,
  },
  {
    name: '智谱 GLM-4',
    tag: '清华系核心',
    desc: '高精度多模态支持，提供复杂的 Agent 编排',
    icon: Cpu,
  },
  {
    name: '豆包 Doubao',
    tag: '超强性价比',
    desc: '极速响应，适合大规模交互与文本处理场景',
    icon: Bolt,
  },
  {
    name: '文心一言 ERNIE',
    tag: '经典国产',
    desc: '中文语境原生适配，深度集成的搜索增强',
    icon: Globe,
  },
]

export const imageModels: readonly ModelEntry[] = [
  {
    name: 'Midjourney',
    tag: '设计利器',
    desc: '业界领先的艺术图像生成，支持丰富的指令风格',
    icon: Sparkles,
  },
  {
    name: 'Stable Diffusion',
    tag: '开源生态',
    desc: '强大的图片生成生态，海量微调与精准控制',
    icon: Grid2x2,
  },
  {
    name: 'DALL·E 3',
    tag: '原日语境',
    desc: '精准遵循复杂提示词，极高的语义理解能力',
    icon: Box,
  },
  {
    name: '智谱 CogView',
    tag: '国产新秀',
    desc: '强大的中文语义理解与多模态交互生成能力',
    icon: Eye,
  },
  {
    name: '腾讯混元 DiT',
    tag: '写实王者',
    desc: '全链路自研大模型，优异的中文写实与二次元生成',
    icon: Sun,
  },
]

export const videoModels: readonly ModelEntry[] = [
  {
    name: '可灵 Kling',
    tag: '国产视频之光',
    desc: '支持超强运镜控制与极高的物理引擎真实度',
    icon: Video,
  },
  {
    name: '智谱 CogVideo',
    tag: '开源力作',
    desc: '高分辨率多帧渲染，支持丰富的光影艺术表现',
    icon: Cloud,
  },
  {
    name: 'Sora',
    tag: '电影级视频',
    desc: '提供长达 60 秒的影视级运镜与复杂场景生成',
    icon: Play,
  },
  {
    name: 'Runway Gen-3',
    tag: '行业标杆',
    desc: '专业级别视频后期与高度可控的运镜设计',
    icon: TrendingUp,
  },
  {
    name: 'Luma Dream Machine',
    tag: '写实动作',
    desc: '极佳的三维空间一致性，打造丝滑动作生成',
    icon: Box,
  },
]
