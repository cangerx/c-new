<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div v-else class="home-shell min-h-screen text-zinc-900 dark:text-zinc-100 font-sans antialiased selection:bg-zinc-200 dark:selection:bg-zinc-800 transition-colors duration-300">
    
    <!-- Header iOS style -->
    <header 
      class="glass-header fixed top-0 inset-x-0 z-40 transition-all duration-500 pointer-events-none"
      :class="[isScrolled ? 'pt-4 px-4' : 'pt-0 px-0']"
    >
      <nav 
        class="liquid-nav mx-auto flex items-center justify-between transition-all duration-500 pointer-events-auto"
        :class="[
          isScrolled 
            ? 'max-w-5xl rounded-full px-6 py-2 scale-95 md:scale-100 liquid-nav-compact' 
            : 'max-w-7xl rounded-none px-4 py-5 sm:px-6 w-full liquid-nav-expanded'
        ]"
      >
        <div 
          class="flex min-w-0 items-center gap-3 transition-all duration-500"
          :class="[isScrolled ? 'translate-x-1 scale-95 md:scale-100' : 'translate-x-0 scale-100']"
        >
          <div class="flex h-9 w-9 shrink-0 items-center justify-center overflow-hidden rounded-lg border border-zinc-200/20 bg-white/10 shadow-sm dark:border-white/10">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain p-1" />
          </div>
          <div class="min-w-0">
            <div class="truncate text-sm font-semibold leading-5 text-zinc-900 dark:text-white tracking-tight">{{ siteName }}</div>
            <div class="hidden max-w-[180px] truncate text-[10px] text-zinc-500 dark:text-zinc-400 sm:block tracking-tight font-medium">
              {{ navSubtitle }}
            </div>
          </div>
        </div>

        <!-- Floating Capsule Pill Menu -->
        <div 
          class="liquid-menu hidden md:flex items-center gap-10 tracking-tight transition-all duration-500"
          :class="[
            isScrolled 
              ? 'text-xs rounded-full px-7 py-2.5 text-zinc-700 dark:text-zinc-300 liquid-menu-capsule' 
              : 'text-sm rounded-full px-9 py-3 text-zinc-700 dark:text-zinc-300 font-semibold liquid-menu-capsule'
          ]"
        >
          <a href="#" class="hover:text-zinc-900 dark:hover:text-white transition-colors">首页</a>
          <a href="#features" class="hover:text-zinc-900 dark:hover:text-white transition-colors">核心优势</a>
          <a href="#domestic-models" class="hover:text-zinc-900 dark:hover:text-white transition-colors">国产模型</a>
          <a href="#image-models" class="hover:text-zinc-900 dark:hover:text-white transition-colors">图片模型</a>
          <a href="#video-models" class="hover:text-zinc-900 dark:hover:text-white transition-colors">视频模型</a>
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hover:text-zinc-900 dark:hover:text-white transition-colors"
          >
            开发文档
          </a>
        </div>

        <div 
          class="flex items-center gap-1.5 sm:gap-2.5 transition-all duration-500"
          :class="[isScrolled ? '-translate-x-1 scale-95 md:scale-100' : 'translate-x-0 scale-100']"
        >
          <LocaleSwitcher class="shrink-0" />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="nav-icon-btn hidden sm:inline-flex shrink-0 text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="sm" />
          </a>
          <button
            class="nav-icon-btn hidden sm:inline-flex shrink-0 text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="sm" />
            <Icon v-else name="moon" size="sm" />
          </button>
          
          <!-- Capsule Action Button -->
          <a
            :href="isAuthenticated ? dashboardPath : '/sign-in'"
            class="liquid-action inline-flex h-9 items-center rounded-full border px-5 text-xs font-semibold backdrop-blur-md transition-all duration-500 active:scale-95 whitespace-nowrap"
            :class="[
              isScrolled
                ? 'text-white dark:text-zinc-950 liquid-action-solid'
                : 'text-zinc-800 dark:text-white liquid-action-glass'
            ]"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </a>
        </div>
      </nav>
    </header>

    <main class="relative">
      
      <!-- Hero Section with background image -->
      <section class="hero-section flex min-h-screen flex-col items-center justify-center px-4 text-center relative overflow-hidden">
        <!-- Physical DOM background wrapper for strict isolation (only inside Hero Section) -->
        <div class="home-bg-overlay">
          <div ref="grainShaderEl" class="hero-grain-shader"></div>
        </div>
        
        <!-- Mask overlay to ensure text readability -->
        <div class="hero-gradient-scrim pointer-events-none z-0"></div>

        <div class="relative z-10 max-w-4xl mx-auto flex flex-col items-center pt-24 pb-16">
          <!-- Dynamic Stepfun style Badge -->
          <div class="hero-badge animate-fade-up animate-float inline-flex items-center gap-2.5 px-4 h-10 mb-8 text-xs font-medium rounded-full bg-zinc-150/80 text-zinc-805 dark:bg-white/[0.04] dark:text-zinc-200 border border-zinc-200/50 dark:border-white/[0.12] backdrop-blur-md shadow-sm hover:border-zinc-350 dark:hover:border-white/[0.2] transition-colors cursor-default">
            <span class="flex h-1.5 w-1.5 relative">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-1.5 w-1.5 bg-emerald-500"></span>
            </span>
            <span>已全面适配 DeepSeek-R1 与视频生成全系列大模型</span>
          </div>

          <!-- Large heading with restrained line reveal -->
          <h1 class="hero-title leading-[1.1] tracking-[-0.03em] text-4xl sm:text-6xl md:text-7xl lg:text-8xl font-black max-w-4xl mb-6">
            <span class="hero-title-line hero-title-line-1">
              <span>连接全球智能</span>
            </span>
            <br class="hidden sm:inline"/>
            <span class="hero-title-line hero-title-line-2">
              <span>赋能 API 应用创新</span>
            </span>
          </h1>
          
          <div class="hero-kicker animate-fade-up animate-fade-up-3 flex items-center justify-center gap-2 mb-8">
            <span class="h-px w-6 bg-zinc-300 dark:bg-white/20"></span>
            <span class="text-xs font-semibold uppercase tracking-widest text-zinc-400 dark:text-zinc-500">Build with models, create with agents</span>
            <span class="h-px w-6 bg-zinc-300 dark:bg-white/20"></span>
          </div>
          
          <p class="animate-fade-up animate-fade-up-3 text-zinc-650 dark:text-white/76 text-sm sm:text-base md:text-lg max-w-2xl leading-relaxed mb-10 tracking-tight font-medium">
            {{ siteSubtitle || '专为开发者与企业打造的下一代 AI API 网关底座。统一聚合全球主流大模型，提供毫秒级智能调度、渠道灾备与精细计费。' }}
          </p>

          <div class="hero-actions animate-fade-up animate-fade-up-4 flex flex-wrap justify-center items-center gap-4 sm:gap-6">
            <a
              :href="isAuthenticated ? dashboardPath : '/sign-in'"
              class="shimmer-btn inline-flex items-center justify-center gap-1.5 h-[52px] px-10 rounded-[60px] text-xs font-semibold text-white bg-zinc-950 hover:bg-zinc-850 dark:bg-white dark:text-zinc-950 dark:hover:bg-zinc-100 transition-all duration-200 shadow-md hover:-translate-y-0.5 active:translate-y-0 overflow-hidden"
            >
              <span class="relative z-10 flex items-center gap-1.5">
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="xs" />
              </span>
            </a>
            
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="inline-flex items-center justify-center gap-1.5 h-[52px] px-10 rounded-[60px] text-xs font-semibold text-zinc-700 hover:text-zinc-950 dark:text-zinc-300 dark:hover:text-white border border-zinc-200 dark:border-white/[0.12] bg-white/40 dark:bg-white/[0.04] hover:bg-white/80 dark:hover:bg-white/[0.08] backdrop-blur-md transition-all duration-200 shadow-sm hover:-translate-y-0.5 active:translate-y-0"
            >
              <span>查看开发文档</span>
              <Icon name="chevronRight" size="xs" />
            </a>
          </div>
        </div>
      </section>

      <!-- Features Section (Simplified and Refined) -->
      <section id="features" class="py-24 border-t border-zinc-200/20 dark:border-zinc-900/40 relative">
        <div class="mx-auto max-w-7xl px-4">
          <div class="reveal-element section-heading max-w-2xl mx-auto text-center space-y-4 mb-16">
            <p class="eyebrow text-[10px] font-bold tracking-widest uppercase text-zinc-400 dark:text-zinc-500">Core Features</p>
            <h2 class="text-3xl sm:text-4xl lg:text-5xl font-bold text-zinc-900 dark:text-white tracking-tight leading-tight">卓越网关性能，赋能业务增长</h2>
            <p class="text-zinc-500 dark:text-zinc-400 text-sm tracking-tight">为您提供大模型集成与交付的一站式基础设施服务。</p>
          </div>

          <div class="reveal-element delay-100 grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">
            <!-- Card 1 -->
            <div class="feature-card p-8 rounded-3xl border border-zinc-200/50 dark:border-zinc-800 bg-white/50 dark:bg-zinc-950/40 backdrop-blur-sm shadow-sm transition-all duration-300 hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:hover:shadow-[0_0_25px_rgba(59,130,246,0.1)] hover:scale-[1.01]">
              <div class="inline-flex p-3 rounded-2xl bg-blue-50 dark:bg-blue-950/40 text-blue-650 dark:text-blue-400 border border-blue-100/30 dark:border-blue-800/30 shadow-sm mb-6">
                <Icon name="server" size="md" />
              </div>
              <h3 class="text-lg font-bold text-zinc-900 mb-3 dark:text-white">统一聚合 简易集成</h3>
              <p class="text-zinc-600 dark:text-zinc-400 text-xs sm:text-sm leading-relaxed tracking-tight">
                统一多种主流大模型协议。支持流式传输与函数调用，通过单一 API 密钥，在几行代码内接入全球顶尖 AI 能力。
              </p>
            </div>

            <!-- Card 2 -->
            <div class="feature-card p-8 rounded-3xl border border-zinc-200/50 dark:border-zinc-800 bg-white/50 dark:bg-zinc-950/40 backdrop-blur-sm shadow-sm transition-all duration-300 hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:hover:shadow-[0_0_25px_rgba(16,185,129,0.1)] hover:scale-[1.01]">
              <div class="inline-flex p-3 rounded-2xl bg-emerald-50 dark:bg-emerald-950/40 text-emerald-650 dark:text-emerald-400 border border-emerald-100/30 dark:border-emerald-800/30 shadow-sm mb-6">
                <Icon name="arrowsUpDown" size="md" />
              </div>
              <h3 class="text-lg font-bold text-zinc-900 mb-3 dark:text-white">智能调度 稳定灾备</h3>
              <p class="text-zinc-600 dark:text-zinc-400 text-xs sm:text-sm leading-relaxed tracking-tight">
                基于并发数、可用性与耗时自动分发。首字延迟深度优化，支持故障自动重试与上游负载均衡，保障服务始终在线。
              </p>
            </div>

            <!-- Card 3 -->
            <div class="feature-card p-8 rounded-3xl border border-zinc-200/50 dark:border-zinc-800 bg-white/50 dark:bg-zinc-950/40 backdrop-blur-sm shadow-sm transition-all duration-300 hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:hover:shadow-[0_0_25px_rgba(139,92,246,0.1)] hover:scale-[1.01]">
              <div class="inline-flex p-3 rounded-2xl bg-violet-50 dark:bg-violet-950/40 text-violet-650 dark:text-violet-400 border border-violet-100/30 dark:border-violet-800/30 shadow-sm mb-6">
                <Icon name="shield" size="md" />
              </div>
              <h3 class="text-lg font-bold text-zinc-900 mb-3 dark:text-white">精细管控 安全合规</h3>
              <p class="text-zinc-600 dark:text-zinc-400 text-xs sm:text-sm leading-relaxed tracking-tight">
                支持密钥级额度限制、RPM/TPM 频率控制与 Token 细粒度审计。内置合规治理模块，有效规避滥用风险。
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- Playground / Studio CTA Section (Stepfun Model Experience Center inspired) -->
      <section class="py-24 border-t border-zinc-200/20 dark:border-zinc-900/40 bg-zinc-50/50 dark:bg-[#070708]/30 relative overflow-hidden">
        <div class="absolute -top-40 -right-40 w-96 h-96 bg-blue-500/10 dark:bg-blue-600/5 rounded-full blur-[100px] pointer-events-none"></div>
        <div class="absolute -bottom-40 -left-40 w-96 h-96 bg-purple-500/10 dark:bg-purple-600/5 rounded-full blur-[100px] pointer-events-none"></div>

        <div class="mx-auto max-w-7xl px-4">
          <div class="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
            
            <!-- Left Info Panel (5 columns) -->
            <div class="reveal-element lg:col-span-5 space-y-6 text-left">
              <div class="inline-flex items-center gap-1.5 px-3 py-1 text-[11px] font-bold rounded-full bg-blue-50 text-blue-650 dark:bg-blue-950/30 dark:text-blue-400 border border-blue-200/30 dark:border-blue-800/30">
                <Icon name="terminal" size="xs" />
                <span>MODEL PLAYGROUND</span>
              </div>
              
              <h2 class="text-3xl sm:text-4xl lg:text-5xl font-black text-zinc-950 dark:text-white tracking-tight leading-[1.15]">
                无需编写任何代码<br/>直接在浏览器调试
              </h2>
              
              <p class="text-zinc-650 dark:text-zinc-400 text-sm leading-relaxed tracking-tight font-medium">
                我们为开发者与企业打造了沉浸式的大模型体验中心。您可以在可视化控制台中一键切换、对比多款旗舰模型，观测实时延迟、吞吐量和运行成本，快速筛选最适配业务场景的智能底座。
              </p>
              
              <ul class="space-y-3 pt-2 text-xs sm:text-sm font-semibold text-zinc-700 dark:text-zinc-300">
                <li class="flex items-center gap-2.5">
                  <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-950/30 dark:text-emerald-450 border border-emerald-250/20">
                    <Icon name="check" size="xs" />
                  </span>
                  <span>统一 Payload 格式，支持单请求多模型横向分发</span>
                </li>
                <li class="flex items-center gap-2.5">
                  <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-950/30 dark:text-emerald-450 border border-emerald-250/20">
                    <Icon name="check" size="xs" />
                  </span>
                  <span>毫秒级首字延迟（TTFT）与 Token 生成速率可视化</span>
                </li>
                <li class="flex items-center gap-2.5">
                  <span class="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-emerald-50 text-emerald-600 dark:bg-emerald-950/30 dark:text-emerald-450 border border-emerald-250/20">
                    <Icon name="check" size="xs" />
                  </span>
                  <span>直接导出符合 OpenAI 规范的 cURL/Python/JS 代码</span>
                </li>
              </ul>
              
              <div class="pt-4 flex flex-wrap gap-4">
                <a
                  :href="isAuthenticated ? dashboardPath : '/sign-in'"
                  class="inline-flex items-center justify-center gap-1.5 h-[52px] px-8 rounded-[60px] text-xs font-bold text-white bg-zinc-950 hover:bg-zinc-850 dark:bg-white dark:text-zinc-950 dark:hover:bg-zinc-100 transition-all duration-200 shadow-md hover:-translate-y-0.5 active:translate-y-0"
                >
                  <span>立即进入体验中心</span>
                  <Icon name="arrowRight" size="xs" />
                </a>
              </div>
            </div>

            <!-- Right macOS Console Mockup (7 columns) -->
            <div class="reveal-element delay-100 lg:col-span-7 w-full">
              <div class="playground-console w-full rounded-2xl border border-zinc-200/80 dark:border-white/[0.08] bg-white/70 dark:bg-[#0b0b0d]/80 backdrop-blur-md shadow-2xl overflow-hidden text-left flex flex-col h-[520px]">
                
                <!-- macOS Window Header -->
                <div class="playground-console-header h-12 border-b border-zinc-200/60 dark:border-white/[0.06] bg-zinc-100/50 dark:bg-[#121215]/50 px-4 flex items-center justify-between shrink-0 select-none">
                  <!-- Traffic Lights -->
                  <div class="flex items-center gap-1.5 w-16">
                    <span class="w-3 h-3 rounded-full bg-[#ff5f56] border border-[#e0443e]"></span>
                    <span class="w-3 h-3 rounded-full bg-[#ffbd2e] border border-[#dea123]"></span>
                    <span class="w-3 h-3 rounded-full bg-[#27c93f] border border-[#1aab29]"></span>
                  </div>
                  <!-- Tabs -->
                  <div class="playground-tabs flex items-center gap-1 text-[11px] font-semibold text-zinc-500 dark:text-zinc-400">
                    <span class="px-3 py-1.5 rounded-md bg-white dark:bg-[#1c1c21] text-zinc-800 dark:text-white border border-zinc-200/40 dark:border-white/[0.04] flex items-center gap-1.5">
                      <Icon name="chat" size="xs" class="text-blue-500"/> Chat Playground
                    </span>
                    <span class="px-3 py-1.5 rounded-md hover:bg-zinc-200/40 dark:hover:bg-white/[0.03] cursor-pointer transition-colors flex items-center gap-1.5">
                      <Icon name="terminal" size="xs"/> cURL
                    </span>
                    <span class="hidden sm:inline-flex px-3 py-1.5 rounded-md hover:bg-zinc-200/40 dark:hover:bg-white/[0.03] cursor-pointer transition-colors flex items-center gap-1.5">
                      <Icon name="chart" size="xs"/> Metrics
                    </span>
                  </div>
                  <!-- Status Indicator -->
                  <div class="playground-status flex items-center gap-1.5 text-[10px] font-bold text-emerald-500 dark:text-emerald-400 bg-emerald-500/10 px-2 py-0.5 rounded-md border border-emerald-500/20">
                    <span class="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse"></span>
                    <span>ONLINE</span>
                  </div>
                </div>

                <!-- Window Content (Split Workspace) -->
                <div class="flex flex-1 min-h-0 overflow-hidden divide-x divide-zinc-200/60 dark:divide-white/[0.06]">
                  
                  <!-- Left Sidebar Panel (Control Center) -->
                  <div class="hidden sm:flex flex-col w-56 p-4 space-y-5 bg-zinc-50/30 dark:bg-[#09090b]/40 shrink-0 overflow-y-auto">
                    <div class="space-y-1.5">
                      <label class="text-[10px] font-bold tracking-wider uppercase text-zinc-400 dark:text-zinc-550">Target Model</label>
                      <div class="h-9 px-3 rounded-lg bg-white dark:bg-[#16161a] border border-zinc-200/60 dark:border-white/[0.08] flex items-center justify-between text-xs font-semibold text-zinc-800 dark:text-zinc-200 shadow-sm cursor-pointer hover:border-zinc-350 dark:hover:border-white/[0.15]">
                        <span class="flex items-center gap-1.5">
                          <span class="w-2 h-2 rounded-full bg-blue-500"></span>
                          DeepSeek-R1
                        </span>
                        <Icon name="chevronDown" size="xs" class="opacity-60" />
                      </div>
                    </div>

                    <div class="space-y-2">
                      <div class="flex items-center justify-between text-[10px] font-bold tracking-wider uppercase text-zinc-400 dark:text-zinc-550">
                        <span>Temperature</span>
                        <span class="text-blue-500 font-mono">0.70</span>
                      </div>
                      <div class="relative w-full h-1.5 bg-zinc-200 dark:bg-zinc-800 rounded-full">
                        <div class="absolute left-0 top-0 h-full w-[70%] bg-blue-500 rounded-full"></div>
                        <div class="absolute left-[70%] top-1/2 -translate-x-1/2 -translate-y-1/2 w-3.5 h-3.5 rounded-full bg-white dark:bg-zinc-150 shadow border border-zinc-300 dark:border-zinc-750 cursor-pointer"></div>
                      </div>
                    </div>

                    <div class="space-y-1.5">
                      <label class="text-[10px] font-bold tracking-wider uppercase text-zinc-400 dark:text-zinc-550">System Prompt</label>
                      <div class="p-2.5 rounded-lg bg-white dark:bg-[#16161a] border border-zinc-200/60 dark:border-white/[0.08] text-[11px] font-medium leading-relaxed text-zinc-500 dark:text-zinc-450 h-32 overflow-hidden shadow-inner">
                        You are a professional software architect. Explain how API gateways optimize latency and channel load balancing.
                      </div>
                    </div>
                  </div>

                  <!-- Right Main Panel (Chat Console) -->
                  <div class="flex-1 flex flex-col min-w-0 bg-white/40 dark:bg-[#070709]/20 overflow-hidden">
                    <!-- Messages Area -->
                    <div class="flex-1 p-4 space-y-4 overflow-y-auto text-xs font-semibold scrollbar-thin">
                      
                      <!-- User Message -->
                      <div class="flex items-start gap-2.5 max-w-[85%]">
                        <div class="w-6 h-6 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-[10px] shrink-0 border border-zinc-200/40 dark:border-white/[0.04]">U</div>
                        <div class="px-3.5 py-2.5 rounded-2xl rounded-tl-none bg-zinc-105 dark:bg-[#1a1a20] text-zinc-850 dark:text-zinc-200 leading-relaxed shadow-sm font-medium">
                          CCAPI 平台如何保障多渠道大模型的延迟与稳定性？
                        </div>
                      </div>

                      <!-- Assistant Message (Streaming) -->
                      <div class="flex items-start gap-2.5 max-w-[90%] ml-auto flex-row-reverse">
                        <div class="w-6 h-6 rounded-full bg-blue-600 text-white flex items-center justify-center text-[10px] shrink-0">AI</div>
                        <div class="space-y-2.5 w-full">
                          
                          <!-- Model Stats Badge Group -->
                          <div class="flex flex-wrap gap-1.5 items-center justify-end">
                            <span class="text-[9px] font-bold px-2 py-0.5 rounded-md bg-blue-500/10 text-blue-500 dark:text-blue-400 border border-blue-500/20">DeepSeek-R1</span>
                            <span class="text-[9px] font-mono px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800/80 text-zinc-500 dark:text-zinc-400 border border-zinc-200/30 dark:border-white/[0.04]">TTFT: 14ms</span>
                            <span class="text-[9px] font-mono px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800/80 text-zinc-500 dark:text-zinc-400 border border-zinc-200/30 dark:border-white/[0.04]">120 tok/s</span>
                            <span class="text-[9px] font-mono px-2 py-0.5 rounded-md bg-emerald-500/10 text-emerald-500 dark:text-emerald-450 border border-emerald-500/20">节省 80% 成本</span>
                          </div>

                          <div class="px-3.5 py-2.5 rounded-2xl rounded-tr-none bg-blue-50/40 dark:bg-blue-950/15 border border-blue-100/30 dark:border-blue-900/10 text-zinc-805 dark:text-zinc-250 leading-relaxed shadow-sm space-y-2 font-normal">
                            <p class="font-medium text-blue-650 dark:text-blue-400 text-xs">CCAPI 通过三大底层核心技术确保旗舰级稳定性：</p>
                            <ol class="list-decimal pl-4 text-[11px] space-y-1 text-zinc-650 dark:text-zinc-350">
                              <li><strong class="text-zinc-850 dark:text-zinc-200">毫秒级智能测速路由</strong>：实时监控全球百余个核心节点的延迟，自动为用户分发响应最快的通道。</li>
                              <li><strong class="text-zinc-850 dark:text-zinc-200">多渠道自动无感灾备</strong>：某上游服务商（如 OpenAI）发生拥堵或限流时，网关在 0 毫秒内自动热重试其他备份提供商，外部访问完全不中断。</li>
                              <li><strong class="text-zinc-850 dark:text-zinc-200">多模型高并发调度引擎</strong>：自研的并发控制与缓冲队列，能够轻松承载百万 QPS 级高频调用。</li>
                            </ol>
                            <div class="flex items-center gap-1 pt-1.5 border-t border-blue-200/20 dark:border-blue-900/10 text-[10px] text-zinc-400 dark:text-zinc-500">
                              <span class="w-1.5 h-1.5 rounded-full bg-blue-500 animate-ping"></span>
                              <span>回答已生成完毕 (Total Tokens: 348)</span>
                            </div>
                          </div>
                        </div>
                      </div>

                    </div>

                    <!-- Input Bar -->
                    <div class="p-3.5 border-t border-zinc-200/60 dark:border-white/[0.06] bg-zinc-50/40 dark:bg-[#0c0c0f]/60 flex items-center gap-2 shrink-0">
                      <div class="flex-1 h-9 px-3 rounded-full bg-white dark:bg-[#141417] border border-zinc-200/65 dark:border-white/[0.08] flex items-center justify-between text-xs text-zinc-400 shadow-sm cursor-text hover:border-zinc-350 dark:hover:border-white/[0.15]">
                        <span class="font-medium truncate flex items-center">输入 Prompt，例如"用 Go 写一个高性能 API 转发代理"<span class="typing-cursor">|</span></span>
                        <Icon name="chat" size="xs" class="opacity-40" />
                      </div>
                      <button class="w-9 h-9 rounded-full bg-zinc-950 dark:bg-white text-white dark:text-black flex items-center justify-center shadow-md active:scale-95 transition-transform">
                        <Icon name="arrowRight" size="xs" />
                      </button>
                    </div>

                  </div>
                </div>

                <!-- Footer Console Stats -->
                <div class="playground-console-footer h-8 border-t border-zinc-200/60 dark:border-white/[0.06] bg-zinc-150/40 dark:bg-[#0b0b0e]/80 px-4 flex items-center justify-between shrink-0 select-none text-[10px] font-bold text-zinc-450 dark:text-zinc-500">
                  <div class="flex items-center gap-3">
                    <span class="flex items-center gap-1.5"><span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Latency: 12ms</span>
                    <span class="hidden sm:inline-flex items-center gap-1.5"><span class="w-1.5 h-1.5 rounded-full bg-emerald-500"></span> Uptime: 99.99%</span>
                  </div>
                  <div class="flex items-center gap-3">
                    <span>1,824 Requests/min</span>
                    <span>HTTPS / API v1</span>
                  </div>
                </div>

              </div>
            </div>

          </div>
        </div>
      </section>

      <!-- Domestic Model Section -->
      <section id="domestic-models" class="py-24 border-t border-zinc-200/20 dark:border-zinc-900/40 relative">
        <div class="mx-auto max-w-7xl px-4">
          <!-- Top Header -->
          <div class="reveal-element flex flex-col md:flex-row justify-between items-start md:items-end mb-8 gap-4">
            <div class="inline-flex items-center gap-2 px-6 py-2 border border-zinc-900 dark:border-white text-sm font-semibold tracking-wide text-zinc-900 dark:text-white cursor-default">
              文本大模型 <Icon name="arrowRight" size="xs" />
            </div>
            <div class="text-left md:text-right">
              <h2 class="text-3xl sm:text-4xl font-bold text-zinc-900 dark:text-white tracking-tight leading-tight mb-2">国产旗舰大语言模型矩阵</h2>
              <p class="text-zinc-500 dark:text-zinc-400 text-sm tracking-tight max-w-xl md:ml-auto">聚合国内最强推理与长文本模型，全链路高速调度，提供稳定可靠的底层支撑体系。</p>
            </div>
          </div>

          <!-- Hero Banner -->
          <div class="model-banner reveal-element delay-100 relative w-full rounded-2xl overflow-hidden mb-12 bg-gradient-to-br from-blue-950 via-slate-900 to-indigo-950 text-white shadow-xl flex flex-col justify-between p-8 md:p-12">
            <!-- Decorative diagonal lines -->
            <div class="absolute inset-0 opacity-10 bg-[linear-gradient(45deg,transparent_48%,rgba(255,255,255,0.8)_50%,transparent_52%)] bg-[length:40px_40px]"></div>
            
            <div class="relative z-10 flex flex-col gap-6 max-w-sm mb-12 md:mb-16">
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">100万长文本处理</div>
                <div class="text-xs text-blue-100/70">深度优化百万字长文本理解和精准提炼能力</div>
              </div>
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">20% 推理能力提升</div>
                <div class="text-xs text-blue-100/70">相比传统接口，逻辑分析与数学推理能力大幅增强</div>
              </div>
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">强大的 Agent 编排</div>
                <div class="text-xs text-blue-100/70">支持高频逻辑推理任务与复杂的工具调用链</div>
              </div>
            </div>

            <div class="relative z-10 flex flex-wrap gap-8 md:gap-16 items-end mt-4 pt-8 border-t border-white/10 counter-section">
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1"><span class="counter-value" data-target="99.9" data-decimals="1">0</span><span class="text-2xl">%</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-blue-100/60">SLA 稳定性保障</div>
              </div>
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1"><span class="counter-value" data-target="100" data-decimals="0">0</span><span class="text-2xl">+</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-blue-100/60">核心并发节点</div>
              </div>
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1">&lt;<span class="counter-value" data-target="50" data-decimals="0">0</span><span class="text-2xl">ms</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-blue-100/60">平均网络延迟</div>
              </div>
            </div>
          </div>

          <div class="reveal-element delay-100 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-6">
            <div
              v-for="(model, idx) in domesticModels"
              :key="model.name"
              class="card-stagger feature-card p-6 rounded-2xl border border-zinc-200/50 dark:border-zinc-800 bg-white/50 dark:bg-zinc-950/40 backdrop-blur-sm shadow-sm flex flex-col justify-between hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:hover:shadow-[0_0_25px_rgba(99,102,241,0.15)] hover:scale-[1.02] transition-all duration-300"
              :style="{ animationDelay: (idx * 0.08) + 's' }"
            >
              <div>
                <div class="flex items-center justify-between mb-4">
                  <div class="inline-flex p-2 rounded-xl bg-indigo-50 dark:bg-indigo-950/40 text-indigo-650 dark:text-indigo-400 border border-indigo-100/30 dark:border-indigo-800/30">
                    <Icon :name="model.icon" size="sm" />
                  </div>
                  <span class="text-[10px] font-bold text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-950/30 px-2 py-0.5 rounded-md mb-3 border border-indigo-200/30 dark:border-indigo-800/30">
                    {{ model.tag }}
                  </span>
                </div>
                <h3 class="text-base font-bold text-zinc-950 dark:text-white mb-2">{{ model.name }}</h3>
                <p class="text-zinc-500 dark:text-zinc-400 text-xs leading-relaxed tracking-tight">{{ model.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Video Model Section -->
      <!-- Image Models Section -->
      <section id="image-models" class="py-24 border-t border-zinc-200/20 dark:border-zinc-900/40 relative">
        <div class="mx-auto max-w-7xl px-4">
          <!-- Top Header -->
          <div class="reveal-element flex flex-col md:flex-row justify-between items-start md:items-end mb-8 gap-4">
            <div class="inline-flex items-center gap-2 px-6 py-2 border border-zinc-900 dark:border-white text-sm font-semibold tracking-wide text-zinc-900 dark:text-white cursor-default">
              图像大模型 <Icon name="arrowRight" size="xs" />
            </div>
            <div class="text-left md:text-right">
              <h2 class="text-3xl sm:text-4xl font-bold text-zinc-900 dark:text-white tracking-tight leading-tight mb-2">高保真艺术与写实图像生成</h2>
              <p class="text-zinc-500 dark:text-zinc-400 text-sm tracking-tight max-w-xl md:ml-auto">丰富的提示词控制，满足全场景艺术创作与专业级商业视觉设计。</p>
            </div>
          </div>

          <!-- Hero Banner -->
          <div class="model-banner reveal-element delay-100 relative w-full rounded-2xl overflow-hidden mb-12 bg-gradient-to-br from-fuchsia-950 via-zinc-900 to-purple-950 text-white shadow-xl flex flex-col justify-between p-8 md:p-12">
            <div class="absolute inset-0 opacity-[0.03] bg-[radial-gradient(ellipse_at_center,rgba(255,255,255,1)_0%,transparent_100%)] bg-[length:120px_120px] bg-repeat"></div>
            
            <div class="relative z-10 flex flex-col gap-6 max-w-sm mb-12 md:mb-16">
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">电影级写实与光影</div>
                <div class="text-xs text-purple-100/70">极强的照片写实度，细节纹理完美呈现</div>
              </div>
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">高精准的指令遵循</div>
                <div class="text-xs text-purple-100/70">完美响应复杂的文本描述与多重元素构图</div>
              </div>
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">ControlNet 深度控制</div>
                <div class="text-xs text-purple-100/70">支持动作、线稿、深度图等条件生成限制</div>
              </div>
            </div>

            <div class="relative z-10 flex flex-wrap gap-8 md:gap-16 items-end mt-4 pt-8 border-t border-white/10 counter-section">
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1">4K</div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-purple-100/60">超高分辨率输出</div>
              </div>
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1">&lt;<span class="counter-value" data-target="2" data-decimals="0">0</span><span class="text-2xl">s</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-purple-100/60">单图毫秒级出图</div>
              </div>
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1"><span class="counter-value" data-target="50" data-decimals="0">0</span><span class="text-2xl">+</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-purple-100/60">风格化微调模型</div>
              </div>
            </div>
          </div>

          <div class="reveal-element delay-100 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-6">
            <div
              v-for="(model, idx) in imageModels"
              :key="model.name"
              class="card-stagger feature-card p-6 rounded-2xl border border-zinc-200/50 dark:border-zinc-800 bg-white/50 dark:bg-zinc-950/40 backdrop-blur-sm shadow-sm flex flex-col justify-between hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:hover:shadow-[0_0_25px_rgba(168,85,247,0.15)] hover:scale-[1.02] transition-all duration-300"
              :style="{ animationDelay: (idx * 0.08) + 's' }"
            >
              <div>
                <div class="flex items-center justify-between mb-4">
                  <div class="inline-flex p-2 rounded-xl bg-purple-50 dark:bg-purple-950/40 text-purple-600 dark:text-purple-400 border border-purple-100/30 dark:border-purple-800/30">
                    <Icon :name="model.icon" size="sm" />
                  </div>
                  <span class="text-[10px] font-bold text-purple-600 dark:text-purple-400 bg-purple-50 dark:bg-purple-950/30 px-2 py-0.5 rounded-md border border-purple-200/30 dark:border-purple-800/30">
                    {{ model.tag }}
                  </span>
                </div>
                <h3 class="text-base font-bold text-zinc-900 mb-2 dark:text-white">{{ model.name }}</h3>
                <p class="text-zinc-500 dark:text-zinc-400 text-xs leading-relaxed tracking-tight">{{ model.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      <section id="video-models" class="py-24 border-t border-zinc-200/20 dark:border-zinc-900/40 relative">
        <div class="mx-auto max-w-7xl px-4">
          <!-- Top Header -->
          <div class="reveal-element flex flex-col md:flex-row justify-between items-start md:items-end mb-8 gap-4">
            <div class="inline-flex items-center gap-2 px-6 py-2 border border-zinc-900 dark:border-white text-sm font-semibold tracking-wide text-zinc-900 dark:text-white cursor-default">
              视频大模型 <Icon name="arrowRight" size="xs" />
            </div>
            <div class="text-left md:text-right">
              <h2 class="text-3xl sm:text-4xl font-bold text-zinc-900 dark:text-white tracking-tight leading-tight mb-2">影视级物理时空视频生成</h2>
              <p class="text-zinc-500 dark:text-zinc-400 text-sm tracking-tight max-w-xl md:ml-auto">统一支持主流文生视频与图生视频服务，突破想象力限制，重塑视觉边界。</p>
            </div>
          </div>

          <!-- Hero Banner -->
          <div class="model-banner reveal-element delay-100 relative w-full rounded-2xl overflow-hidden mb-12 bg-gradient-to-br from-zinc-900 via-[#111111] to-stone-900 text-white shadow-xl flex flex-col justify-between p-8 md:p-12 border border-zinc-800">
            <div class="absolute inset-0 opacity-10" style="background-image: linear-gradient(0deg, transparent 24%, rgba(255, 255, 255, .3) 25%, rgba(255, 255, 255, .3) 26%, transparent 27%, transparent 74%, rgba(255, 255, 255, .3) 75%, rgba(255, 255, 255, .3) 76%, transparent 77%, transparent), linear-gradient(90deg, transparent 24%, rgba(255, 255, 255, .3) 25%, rgba(255, 255, 255, .3) 26%, transparent 27%, transparent 74%, rgba(255, 255, 255, .3) 75%, rgba(255, 255, 255, .3) 76%, transparent 77%, transparent); background-size: 30px 30px;"></div>
            
            <div class="relative z-10 flex flex-col gap-6 max-w-sm mb-12 md:mb-16">
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">原生物理世界模拟</div>
                <div class="text-xs text-zinc-400">极佳的三维空间一致性与复杂的物理引擎真实度</div>
              </div>
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">高分辨率多帧渲染</div>
                <div class="text-xs text-zinc-400">支持 1080P 超高清晰度与 60FPS 丝滑帧率输出</div>
              </div>
              <div>
                <div class="text-xl md:text-2xl font-bold mb-1">超长运镜与场景延伸</div>
                <div class="text-xs text-zinc-400">提供长达 60s+ 的影视级镜头语言与多场景转场</div>
              </div>
            </div>

            <div class="relative z-10 flex flex-wrap gap-8 md:gap-16 items-end mt-4 pt-8 border-t border-white/10 counter-section">
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1"><span class="counter-value" data-target="60" data-decimals="0">0</span><span class="text-2xl">s</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-zinc-500">超长视频序列</div>
              </div>
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1"><span class="counter-value" data-target="1080" data-decimals="0">0</span><span class="text-2xl">P</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-zinc-500">原生生成分辨率</div>
              </div>
              <div>
                <div class="text-3xl md:text-5xl font-black tracking-tighter mb-1"><span class="counter-value" data-target="99" data-decimals="0">0</span><span class="text-2xl">%</span></div>
                <div class="text-[10px] font-bold uppercase tracking-widest text-zinc-500">时间线物理一致性</div>
              </div>
            </div>
          </div>

          <div class="reveal-element delay-100 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-6">
            <div
              v-for="(model, idx) in videoModels"
              :key="model.name"
              class="card-stagger feature-card p-6 rounded-2xl border border-zinc-200/50 dark:border-zinc-800 bg-white/50 dark:bg-zinc-950/40 backdrop-blur-sm shadow-sm flex flex-col justify-between hover:shadow-[0_8px_30px_rgb(0,0,0,0.04)] dark:hover:shadow-[0_0_25px_rgba(236,72,153,0.15)] hover:scale-[1.02] transition-all duration-300"
              :style="{ animationDelay: (idx * 0.08) + 's' }"
            >
              <div>
                <div class="flex items-center justify-between mb-4">
                  <div class="inline-flex p-2 rounded-xl bg-pink-50 dark:bg-pink-950/40 text-pink-650 dark:text-pink-400 border border-pink-100/30 dark:border-pink-800/30">
                    <Icon :name="model.icon" size="sm" />
                  </div>
                  <span class="text-[10px] font-bold text-pink-600 dark:text-pink-400 bg-pink-50 dark:bg-pink-950/30 px-2 py-0.5 rounded-md border border-pink-200/30 dark:border-pink-800/30">
                    {{ model.tag }}
                  </span>
                </div>
                <h3 class="text-base font-bold text-zinc-900 mb-2 dark:text-white">{{ model.name }}</h3>
                <p class="text-zinc-500 dark:text-zinc-400 text-xs leading-relaxed tracking-tight">{{ model.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-zinc-200/30 bg-white/45 px-4 py-10 backdrop-blur-xl dark:border-white/10 dark:bg-black/35">
      <div class="mx-auto flex max-w-7xl flex-col items-center justify-between gap-5 text-center sm:flex-row sm:text-left">
        <div class="space-y-1">
          <p class="text-xs font-semibold tracking-tight text-zinc-700 dark:text-zinc-200">
            &copy; {{ currentYear }} {{ siteName }}
          </p>
          <p class="max-w-xl text-[11px] leading-5 text-zinc-500 dark:text-zinc-500">
            {{ navSubtitle }} · {{ t('home.footer.allRightsReserved') }}
          </p>
        </div>
        <div class="flex flex-wrap items-center justify-center gap-2.5 text-xs font-semibold sm:justify-end">
          <span class="inline-flex items-center rounded-full border border-zinc-200/70 bg-white/55 px-4 py-2 text-[11px] font-semibold tracking-tight text-zinc-600 shadow-sm backdrop-blur-md dark:border-white/10 dark:bg-white/[0.04] dark:text-zinc-400">
            CCAI v1.0.0
          </span>
          <span class="inline-flex items-center gap-1.5 rounded-full border border-zinc-200/70 bg-white/55 px-4 py-2 text-[11px] font-semibold tracking-tight text-zinc-600 shadow-sm backdrop-blur-md dark:border-white/10 dark:bg-white/[0.04] dark:text-zinc-400">
            <span class="text-zinc-400 dark:text-zinc-500">Author</span>
            <span class="text-zinc-800 dark:text-zinc-200">canger</span>
          </span>
          <a
            href="https://t.me/cangerx"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex items-center gap-1.5 rounded-full border border-zinc-200/70 bg-white/55 px-4 py-2 text-[11px] font-semibold tracking-tight text-zinc-600 shadow-sm backdrop-blur-md transition-colors hover:border-zinc-300 hover:text-zinc-950 dark:border-white/10 dark:bg-white/[0.04] dark:text-zinc-400 dark:hover:border-white/20 dark:hover:text-white"
          >
            <span class="text-zinc-400 dark:text-zinc-500">TG</span>
            <span>@cangerx</span>
          </a>
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-full border border-zinc-200/70 bg-white/50 px-4 py-2 text-zinc-600 transition-colors hover:border-zinc-300 hover:text-zinc-950 dark:border-white/10 dark:bg-white/[0.04] dark:text-zinc-400 dark:hover:border-white/20 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'

type PaperShadersModule = typeof import('@paper-design/shaders')
type HeroShaderInstance = InstanceType<PaperShadersModule['ShaderMount']>

const { t } = useI18n()

const isScrolled = ref(false)
const grainShaderEl = ref<HTMLElement | null>(null)
let shaderModule: PaperShadersModule | null = null
let heroShader: HeroShaderInstance | null = null
let heroShaderInitHandle: number | null = null
let heroShaderInitUsesIdleCallback = false

function handleScroll() {
  isScrolled.value = window.scrollY > 80
}

const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'CCAPI')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || '')
const navSubtitle = computed(() => siteSubtitle.value || 'AI API Gateway')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))

const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
// new-api has a single dashboard route for both admins and regular users;
// there is no /admin/dashboard equivalent.
const dashboardPath = computed(() => '/dashboard')
const currentYear = computed(() => new Date().getFullYear())

const domesticModels = [
  { name: 'DeepSeek-R1', tag: '最强推理', desc: '深度思考，逻辑推理与代码顶峰模型', icon: 'brain' },
  { name: '通义千问 Qwen-Max', tag: '旗舰模型', desc: '全能基座大模型，极佳的长文本处理能力', icon: 'sparkles' },
  { name: '智谱 GLM-4', tag: '清华系核心', desc: '高精度多模态支持，提供复杂的 Agent 编排', icon: 'cpu' },
  { name: '豆包 Doubao', tag: '超强性价比', desc: '极速响应，适合大规模交互与文本处理场景', icon: 'bolt' },
  { name: '文心一言 ERNIE', tag: '经典国产', desc: '中文语境原生适配，深度集成的搜索增强', icon: 'globe' }
] as const

const imageModels = [
  { name: 'Midjourney', tag: '设计利器', desc: '业界领先的艺术图像生成，支持丰富的指令风格', icon: 'sparkles' },
  { name: 'Stable Diffusion', tag: '开源生态', desc: '强大的图片生成生态，海量微调与精准控制', icon: 'grid' },
  { name: 'DALL·E 3', tag: '原日语境', desc: '精准遵循复杂提示词，极高的语义理解能力', icon: 'cube' },
  { name: '智谱 CogView', tag: '国产新秀', desc: '强大的中文语义理解与多模态交互生成能力', icon: 'eye' },
  { name: '腾讯混元 DiT', tag: '写实王者', desc: '全链路自研大模型，优异的中文写实与二次元生成', icon: 'sun' }
] as const

const videoModels = [
  { name: '可灵 Kling', tag: '国产视频之光', desc: '支持超强运镜控制与极高的物理引擎真实度', icon: 'video' },
  { name: '智谱 CogVideo', tag: '开源力作', desc: '高分辨率多帧渲染，支持丰富的光影艺术表现', icon: 'cloud' },
  { name: 'Sora', tag: '电影级视频', desc: '提供长达 60 秒的影视级运镜与复杂场景生成', icon: 'play' },
  { name: 'Runway Gen-3', tag: '行业标杆', desc: '专业级别视频后期与高度可控的运镜设计', icon: 'trendingUp' },
  { name: 'Luma Dream Machine', tag: '写实动作', desc: '极佳的三维空间一致性，打造丝滑动作生成', icon: 'cube' }
] as const

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  applyHeroShaderTheme()
}

function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  } else {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  }
}

function animateCounter(el: HTMLElement, target: number, decimals: number) {
  const duration = 1800
  const startTime = performance.now()
  const easeOutExpo = (t: number) => t === 1 ? 1 : 1 - Math.pow(2, -10 * t)

  function tick(now: number) {
    const elapsed = now - startTime
    const progress = Math.min(elapsed / duration, 1)
    const easedProgress = easeOutExpo(progress)
    const current = easedProgress * target
    el.textContent = current.toFixed(decimals)
    if (progress < 1) {
      requestAnimationFrame(tick)
    } else {
      el.textContent = target.toFixed(decimals)
    }
  }
  requestAnimationFrame(tick)
}

async function initHeroShader() {
  if (!grainShaderEl.value || heroShader) return

  try {
    await nextTick()

    shaderModule = shaderModule || await import('@paper-design/shaders')
    const {
      ShaderFitOptions,
      ShaderMount,
      defaultObjectSizing,
      grainGradientFragmentShader,
      GrainGradientShapes,
      getShaderColorFromString,
      getShaderNoiseTexture
    } = shaderModule

    const noiseTexture = getShaderNoiseTexture()
    if (noiseTexture) {
      await waitForShaderImage(noiseTexture)
    }

    const sizing = defaultObjectSizing
    const theme = getHeroShaderTheme()

    heroShader = new ShaderMount(
      grainShaderEl.value,
      grainGradientFragmentShader,
      {
        u_fit: ShaderFitOptions.cover,
        u_scale: 1,
        u_rotation: sizing.rotation,
        u_originX: sizing.originX,
        u_originY: sizing.originY,
        u_offsetX: sizing.offsetX,
        u_offsetY: sizing.offsetY,
        u_worldWidth: sizing.worldWidth,
        u_worldHeight: sizing.worldHeight,
        u_colorBack: getShaderColorFromString(theme.colorBack),
        u_colors: theme.colors.map(getShaderColorFromString),
        u_colorsCount: theme.colors.length,
        u_softness: 1,
        u_intensity: theme.intensity,
        u_noise: theme.noise,
        u_shape: GrainGradientShapes.corners,
        u_noiseTexture: noiseTexture
      },
      { alpha: true, antialias: true, premultipliedAlpha: false },
      1.85,
      0,
      1,
      1920 * 1080 * 2
    )
    grainShaderEl.value?.classList.add('is-ready')
  } catch (error) {
    console.warn('[home] Paper shader background disabled:', error)
    grainShaderEl.value?.classList.remove('is-ready')
    heroShader = null
  }
}

function waitForShaderImage(image: HTMLImageElement) {
  if (image.complete && image.naturalWidth > 0) {
    return Promise.resolve()
  }

  return new Promise<void>((resolve, reject) => {
    image.onload = () => resolve()
    image.onerror = () => reject(new Error('Paper shader noise texture failed to load'))
  })
}

function getHeroShaderTheme() {
  return isDark.value
    ? {
        colorBack: '#030305',
        colors: ['#00D5FF', '#7C3AED', '#FF2E88', '#2563EB00'],
        intensity: 1,
        noise: 0.72
      }
    : {
        colorBack: '#F7FAFF',
        colors: ['#0EA5E9', '#8B5CF6', '#F43F5E', '#22D3EE00'],
        intensity: 0.96,
        noise: 0.5
      }
}

function applyHeroShaderTheme() {
  if (!heroShader || !shaderModule) return

  const { getShaderColorFromString } = shaderModule
  const theme = getHeroShaderTheme()
  heroShader.setUniforms({
    u_colorBack: getShaderColorFromString(theme.colorBack),
    u_colors: theme.colors.map(getShaderColorFromString),
    u_colorsCount: theme.colors.length,
    u_intensity: theme.intensity,
    u_noise: theme.noise
  })
  heroShader.setSpeed(1.85)
}

function shouldSkipHeroShader() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

function scheduleHeroShaderInit() {
  if (heroShaderInitHandle !== null || shouldSkipHeroShader()) return

  if (typeof window.requestIdleCallback === 'function') {
    heroShaderInitUsesIdleCallback = true
    heroShaderInitHandle = window.requestIdleCallback(
      () => {
        heroShaderInitHandle = null
        void initHeroShader()
      },
      { timeout: 1800 }
    )
    return
  }

  heroShaderInitUsesIdleCallback = false
  heroShaderInitHandle = window.setTimeout(() => {
    heroShaderInitHandle = null
    void initHeroShader()
  }, 700)
}

onMounted(() => {
  initTheme()
  scheduleHeroShaderInit()

  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }

  // Resolve login state so the nav shows "Dashboard" instead of "Sign in"
  // for visitors who already have a session.
  authStore.probe()

  window.addEventListener('scroll', handleScroll)

  // Set up intersection observer for scroll reveal animations
  const observer = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        entry.target.classList.add('reveal-active')
        observer.unobserve(entry.target)
      }
    })
  }, {
    root: null,
    rootMargin: '0px',
    threshold: 0.1
  })

  setTimeout(() => {
    document.querySelectorAll('.reveal-element').forEach((el) => {
      observer.observe(el)
    })
  }, 100)

  // Animated counter observer - triggers count-up when stats scroll into view
  const counterObserver = new IntersectionObserver((entries) => {
    entries.forEach(entry => {
      if (entry.isIntersecting) {
        const counters = entry.target.querySelectorAll('.counter-value')
        counters.forEach((el) => {
          const target = parseFloat(el.getAttribute('data-target') || '0')
          const decimals = parseInt(el.getAttribute('data-decimals') || '0')
          animateCounter(el as HTMLElement, target, decimals)
        })
        counterObserver.unobserve(entry.target)
      }
    })
  }, { threshold: 0.3 })

  setTimeout(() => {
    document.querySelectorAll('.counter-section').forEach((el) => {
      counterObserver.observe(el)
    })
  }, 100)
})

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll)
  if (heroShaderInitHandle !== null) {
    if (heroShaderInitUsesIdleCallback && typeof window.cancelIdleCallback === 'function') {
      window.cancelIdleCallback(heroShaderInitHandle)
    } else {
      window.clearTimeout(heroShaderInitHandle)
    }
    heroShaderInitHandle = null
  }
  heroShader?.dispose()
  heroShader = null
})
</script>

<style scoped>
/* Marquee Animation styles */
.marquee-track {
  display: flex;
  width: max-content;
}
.animate-marquee {
  animation: marquee 30s linear infinite;
}
@keyframes marquee {
  0% {
    transform: translateX(0);
  }
  100% {
    transform: translateX(calc(-100% - 4rem)); /* Shift left by 100% of marquee list width plus gap-16 (4rem) */
  }
}

/* Custom Scrollbar for Playground mock console */
.scrollbar-thin::-webkit-scrollbar {
  width: 4px;
}
.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}
.scrollbar-thin::-webkit-scrollbar-thumb {
  background: rgba(100, 116, 139, 0.2);
  border-radius: 999px;
}
.dark .scrollbar-thin::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
}

.nav-icon-btn {
  display: inline-flex;
  height: 2rem;
  width: 2rem;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  color: rgb(113 113 122);
  border: 1px solid rgba(255, 255, 255, 0.42);
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.42), rgba(255, 255, 255, 0.16)),
    rgba(255, 255, 255, 0.2);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.6),
    inset 0 -1px 0 rgba(255, 255, 255, 0.14),
    0 8px 22px rgba(15, 23, 42, 0.08);
  -webkit-backdrop-filter: blur(18px) saturate(1.55) brightness(1.04);
  backdrop-filter: blur(18px) saturate(1.55) brightness(1.04);
  transition: transform 180ms cubic-bezier(0.2, 0.8, 0.2, 1), box-shadow 180ms ease, background 180ms ease, color 180ms ease;
}

.nav-icon-btn:hover {
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.58), rgba(255, 255, 255, 0.22)),
    rgba(255, 255, 255, 0.28);
  color: rgb(9 9 11);
  transform: translateY(-1px);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.72),
    inset 0 -1px 0 rgba(255, 255, 255, 0.18),
    0 12px 28px rgba(15, 23, 42, 0.12);
}

.nav-icon-btn:active {
  transform: scale(0.94);
}

.dark .nav-icon-btn {
  color: rgb(161 161 170);
  border-color: rgba(255, 255, 255, 0.13);
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0.035)),
    rgba(12, 12, 14, 0.32);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.18),
    inset 0 -1px 0 rgba(255, 255, 255, 0.05),
    0 10px 24px rgba(0, 0, 0, 0.28);
}

.dark .nav-icon-btn:hover {
  background:
    linear-gradient(145deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.06)),
    rgba(20, 20, 24, 0.42);
  color: white;
}

/* Glass Header base style */
.glass-header {
  position: fixed !important;
  inset: 0 0 auto 0;
  z-index: 40;
}

.liquid-nav {
  position: relative;
}

.liquid-nav-compact {
  background: transparent;
  border: 0;
  box-shadow: none;
}

.dark .liquid-nav-compact {
  background: transparent;
  border-color: transparent;
  box-shadow: none;
}

.liquid-nav-expanded {
  background: transparent;
}

.dark .liquid-nav-expanded {
  background: transparent;
}

.liquid-menu,
.liquid-action {
  position: relative;
  overflow: hidden;
  isolation: isolate;
  -webkit-backdrop-filter: blur(22px) saturate(1.65) contrast(1.08) brightness(1.04);
  backdrop-filter: blur(22px) saturate(1.65) contrast(1.08) brightness(1.04);
  transform: translateZ(0);
  transition:
    transform 220ms cubic-bezier(0.2, 0.8, 0.2, 1),
    box-shadow 220ms ease,
    background 220ms ease,
    border-color 220ms ease;
}

.liquid-menu::before,
.liquid-action::before {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 0;
  border-radius: inherit;
  pointer-events: none;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.72) 0%, rgba(255, 255, 255, 0.18) 28%, rgba(255, 255, 255, 0.06) 58%, rgba(255, 255, 255, 0.32) 100%),
    radial-gradient(circle at 18% 8%, rgba(255, 255, 255, 0.78), transparent 28%),
    radial-gradient(circle at 78% 92%, rgba(255, 255, 255, 0.26), transparent 32%);
  opacity: 0.62;
  mix-blend-mode: screen;
}

.liquid-menu::after,
.liquid-action::after {
  content: '';
  position: absolute;
  inset: 1px;
  z-index: 0;
  border-radius: inherit;
  pointer-events: none;
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.78),
    inset 0 -1px 0 rgba(255, 255, 255, 0.2),
    inset 1px 0 0 rgba(255, 255, 255, 0.24),
    inset -1px 0 0 rgba(255, 255, 255, 0.14),
    inset 0 0 18px rgba(255, 255, 255, 0.12);
}

.liquid-menu-capsule {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.42), rgba(255, 255, 255, 0.18)),
    rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.56);
  box-shadow:
    0 16px 36px rgba(15, 23, 42, 0.12),
    0 1px 1px rgba(255, 255, 255, 0.38),
    inset 0 0 0 0.5px rgba(255, 255, 255, 0.42);
}

.liquid-menu-capsule:hover,
.liquid-action:hover {
  transform: translateY(-1px) translateZ(0);
}

.dark .liquid-menu-capsule {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.13), rgba(255, 255, 255, 0.035)),
    rgba(10, 10, 12, 0.38);
  border-color: rgba(255, 255, 255, 0.16);
  box-shadow:
    0 16px 38px rgba(0, 0, 0, 0.34),
    0 1px 1px rgba(255, 255, 255, 0.06),
    inset 0 0 0 0.5px rgba(255, 255, 255, 0.1);
}

.dark .liquid-menu::before,
.dark .liquid-action::before {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.22), rgba(255, 255, 255, 0.045) 48%, rgba(255, 255, 255, 0.1)),
    radial-gradient(circle at 18% 8%, rgba(255, 255, 255, 0.28), transparent 34%),
    radial-gradient(circle at 76% 92%, rgba(255, 255, 255, 0.1), transparent 34%);
  opacity: 0.76;
}

.dark .liquid-menu::after,
.dark .liquid-action::after {
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.26),
    inset 0 -1px 0 rgba(255, 255, 255, 0.08),
    inset 1px 0 0 rgba(255, 255, 255, 0.1),
    inset -1px 0 0 rgba(255, 255, 255, 0.05),
    inset 0 0 18px rgba(255, 255, 255, 0.04);
}

.liquid-menu a {
  position: relative;
  z-index: 1;
}

.liquid-menu a::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: -0.45rem;
  width: 0.32rem;
  height: 0.32rem;
  border-radius: 999px;
  background: currentColor;
  opacity: 0;
  transform: translateX(-50%) scale(0.4);
  transition: opacity 180ms ease, transform 180ms ease;
}

.liquid-menu a:hover::after {
  opacity: 0.42;
  transform: translateX(-50%) scale(1);
}

.liquid-action-glass {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.48), rgba(255, 255, 255, 0.2)),
    rgba(255, 255, 255, 0.24);
  border-color: rgba(255, 255, 255, 0.58);
  box-shadow:
    0 12px 28px rgba(15, 23, 42, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.72),
    inset 0 -1px 0 rgba(255, 255, 255, 0.18);
}

.liquid-action-glass:hover {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.6), rgba(255, 255, 255, 0.26)),
    rgba(255, 255, 255, 0.3);
}

.dark .liquid-action-glass {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.15), rgba(255, 255, 255, 0.045)),
    rgba(16, 16, 20, 0.38);
  border-color: rgba(255, 255, 255, 0.16);
  box-shadow:
    0 12px 30px rgba(0, 0, 0, 0.34),
    inset 0 1px 0 rgba(255, 255, 255, 0.22),
    inset 0 -1px 0 rgba(255, 255, 255, 0.06);
}

.liquid-action-solid {
  background:
    linear-gradient(135deg, rgba(24, 24, 27, 0.82), rgba(9, 9, 11, 0.62)),
    rgba(9, 9, 11, 0.56);
  border-color: rgba(255, 255, 255, 0.22);
  box-shadow:
    0 12px 30px rgba(15, 23, 42, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    inset 0 -1px 0 rgba(255, 255, 255, 0.08);
}

.dark .liquid-action-solid {
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.86), rgba(255, 255, 255, 0.58)),
    rgba(255, 255, 255, 0.68);
  border-color: rgba(255, 255, 255, 0.5);
  box-shadow:
    0 12px 30px rgba(0, 0, 0, 0.34),
    inset 0 1px 0 rgba(255, 255, 255, 0.74),
    inset 0 -1px 0 rgba(255, 255, 255, 0.28);
}

.liquid-action > span,
.liquid-action > svg {
  position: relative;
  z-index: 1;
}

.liquid-menu-capsule::before,
.liquid-action::before {
  animation: liquidSpecularShift 8s ease-in-out infinite alternate;
}

@keyframes liquidSpecularShift {
  0% {
    background-position: 0% 0%, 18% 8%, 78% 92%;
  }
  100% {
    background-position: 100% 100%, 28% 16%, 68% 82%;
  }
}

/* Full screen background for the entire home shell */
.home-shell {
  position: relative;
  background-color: #f7f8fb;
}

.dark .home-shell {
  background-color: #050505;
}

.home-bg-overlay {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;
  background:
    radial-gradient(circle at 16% 18%, rgba(14, 165, 233, 0.88), transparent 34%),
    radial-gradient(circle at 78% 20%, rgba(139, 92, 246, 0.8), transparent 36%),
    radial-gradient(circle at 58% 76%, rgba(244, 63, 94, 0.72), transparent 38%),
    #f8fafc;
  pointer-events: none;
}

.dark .home-bg-overlay {
  background:
    radial-gradient(circle at 18% 18%, rgba(0, 213, 255, 0.68), transparent 34%),
    radial-gradient(circle at 78% 24%, rgba(124, 58, 237, 0.62), transparent 36%),
    radial-gradient(circle at 58% 76%, rgba(255, 46, 136, 0.54), transparent 38%),
    #050505;
}

.hero-grain-shader {
  position: absolute;
  inset: 0;
  z-index: 0;
  opacity: 1;
  mix-blend-mode: normal;
  pointer-events: none;
}

.hero-grain-shader canvas {
  display: block;
  width: 100%;
  height: 100%;
}

.hero-gradient-scrim {
  content: '';
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 50% 42%, rgba(255, 255, 255, 0.02), rgba(255, 255, 255, 0.24) 72%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.1), rgba(248, 250, 252, 0.34));
}

.dark .hero-gradient-scrim {
  background:
    radial-gradient(circle at 50% 42%, rgba(0, 0, 0, 0.06), rgba(0, 0, 0, 0.34) 68%),
    linear-gradient(180deg, rgba(0, 0, 0, 0.08), rgba(0, 0, 0, 0.38));
}

/* Hero Section */
.hero-section {
  position: relative;
  width: 100%;
  overflow: hidden;
  isolation: isolate;
}

/* Feature card styles */
.feature-card {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.feature-card:hover {
  transform: translateY(-4px);
  border-color: rgba(161, 161, 170, 0.4);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.04);
}

.dark .feature-card:hover {
  border-color: rgba(161, 161, 170, 0.2);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.dark .feature-card {
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.02);
}

/* Scroll Reveal Animations */
.reveal-element {
  opacity: 0;
  transform: translateY(20px);
  transition: all 0.8s cubic-bezier(0.16, 1, 0.3, 1);
}
.reveal-active {
  opacity: 1;
  transform: translateY(0);
}
.delay-100 { transition-delay: 100ms; }

/* Hero Load Animations */
.animate-fade-up {
  animation: fadeUp 1s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
}
.animate-fade-up-1 { animation-delay: 0.1s; }
.animate-fade-up-2 { animation-delay: 0.2s; }
.animate-fade-up-3 { animation-delay: 0.3s; }
.animate-fade-up-4 { animation-delay: 0.4s; }

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Micro-animations */
.animate-float {
  animation: float 6s ease-in-out infinite;
}
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

.shimmer-btn {
  position: relative;
}
.shimmer-btn::after {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(to right, transparent, rgba(255,255,255,0.15), transparent);
  transform: skewX(-20deg);
  animation: shimmer 3.5s infinite;
}
.dark .shimmer-btn::after {
  background: linear-gradient(to right, transparent, rgba(255,255,255,0.1), transparent);
}
@keyframes shimmer {
  0% { left: -100%; }
  20%, 100% { left: 200%; }
}

.hero-title {
  color: #09090b;
  text-wrap: balance;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.72), 0 22px 64px rgba(255, 255, 255, 0.65);
}
.dark .hero-title {
  color: #ffffff;
  text-shadow: 0 1px 0 rgba(255, 255, 255, 0.06), 0 22px 64px rgba(0, 0, 0, 0.55);
}

.hero-title-line {
  display: inline-block;
  overflow: hidden;
  vertical-align: top;
}

.hero-title-line > span {
  display: inline-block;
  opacity: 0;
  transform: translate3d(0, 1.08em, 0);
  animation: heroLineReveal 1120ms cubic-bezier(0.19, 1, 0.22, 1) forwards;
  will-change: transform, opacity;
}

.hero-title-line-2 > span {
  animation-delay: 120ms;
}

@keyframes heroLineReveal {
  0% {
    opacity: 0;
    transform: translate3d(0, 1.08em, 0);
  }
  42% {
    opacity: 1;
  }
  100% {
    opacity: 1;
    transform: translate3d(0, 0, 0);
  }
}

/* Card Stagger Entrance */
.reveal-active .card-stagger {
  animation: cardSlideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
  opacity: 0;
}
@keyframes cardSlideUp {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.96);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

/* Typing Cursor Blink */
.typing-cursor {
  display: inline-block;
  margin-left: 1px;
  font-weight: 400;
  color: #3b82f6;
  animation: cursorBlink 1s step-end infinite;
}
@keyframes cursorBlink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

/* Counter value tabular nums for smooth counting */
.counter-value {
  font-variant-numeric: tabular-nums;
}

@media (max-width: 767px) {
  .home-shell {
    overflow-x: hidden;
  }

  .glass-header {
    padding-top: max(0.75rem, env(safe-area-inset-top)) !important;
    padding-left: 0.75rem !important;
    padding-right: 0.75rem !important;
  }

  .liquid-nav {
    width: 100%;
    min-height: 3.25rem;
    gap: 0.625rem;
    border-radius: 999px !important;
    padding: 0.5rem 0.625rem !important;
    background:
      linear-gradient(135deg, rgba(255, 255, 255, 0.64), rgba(255, 255, 255, 0.22)),
      rgba(255, 255, 255, 0.18);
    border: 1px solid rgba(255, 255, 255, 0.54);
    box-shadow:
      0 14px 34px rgba(15, 23, 42, 0.12),
      inset 0 1px 0 rgba(255, 255, 255, 0.72);
    -webkit-backdrop-filter: blur(22px) saturate(1.45);
    backdrop-filter: blur(22px) saturate(1.45);
  }

  .dark .liquid-nav {
    background:
      linear-gradient(135deg, rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0.04)),
      rgba(10, 10, 12, 0.42);
    border-color: rgba(255, 255, 255, 0.14);
    box-shadow:
      0 14px 34px rgba(0, 0, 0, 0.36),
      inset 0 1px 0 rgba(255, 255, 255, 0.17);
  }

  .liquid-nav > div:first-child {
    gap: 0.5rem;
    min-width: 0;
  }

  .liquid-nav > div:first-child > div:first-child {
    width: 2rem;
    height: 2rem;
    border-radius: 0.75rem;
  }

  .liquid-nav .liquid-action {
    height: 2.25rem;
    padding-left: 0.9rem;
    padding-right: 0.9rem;
    font-size: 0.72rem;
  }

  .hero-section {
    min-height: 100svh;
    padding-left: 1.125rem;
    padding-right: 1.125rem;
  }

  .hero-section > .relative {
    padding-top: 5.5rem;
    padding-bottom: 4rem;
  }

  .hero-badge {
    max-width: 100%;
    height: auto;
    min-height: 2.25rem;
    margin-bottom: 1.35rem;
    padding: 0.55rem 0.8rem;
    line-height: 1.45;
    text-align: left;
  }

  .hero-title {
    max-width: 22rem;
    margin-bottom: 1rem;
    font-size: clamp(2.55rem, 12vw, 3.45rem);
    line-height: 1.04;
    letter-spacing: -0.025em;
  }

  .hero-kicker {
    width: 100%;
    margin-bottom: 1.35rem;
  }

  .hero-kicker span:nth-child(2) {
    max-width: 15rem;
    line-height: 1.4;
    white-space: normal;
  }

  .hero-section p {
    max-width: 21rem;
    margin-bottom: 1.75rem;
    font-size: 0.9rem;
    line-height: 1.75;
  }

  .hero-actions {
    width: 100%;
    gap: 0.75rem;
  }

  .hero-actions > a {
    width: min(100%, 20.5rem);
    height: 3rem;
    padding-left: 1.25rem;
    padding-right: 1.25rem;
  }

  main section:not(.hero-section) {
    padding-top: 4rem;
    padding-bottom: 4rem;
  }

  .section-heading {
    margin-bottom: 2.25rem;
  }

  .feature-card {
    border-radius: 1.25rem !important;
  }

  .playground-console {
    height: 30rem !important;
    border-radius: 1.25rem;
  }

  .playground-console-header {
    height: auto;
    min-height: 3.25rem;
    gap: 0.5rem;
    padding: 0.75rem;
  }

  .playground-tabs {
    min-width: 0;
    flex: 1;
    justify-content: center;
  }

  .playground-tabs span {
    padding-left: 0.55rem;
    padding-right: 0.55rem;
  }

  .playground-status {
    display: none;
  }

  .playground-console-footer {
    height: auto;
    min-height: 2rem;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
  }

  .playground-console-footer > div:last-child {
    display: none;
  }

  .model-banner {
    min-height: auto;
    padding: 1.5rem !important;
    border-radius: 1.25rem;
  }

  .model-banner .counter-section {
    display: grid;
    grid-template-columns: 1fr;
    gap: 1rem;
    padding-top: 1.25rem;
  }

  footer .max-w-7xl {
    align-items: stretch;
  }

  footer .flex-wrap {
    width: 100%;
    gap: 0.625rem;
  }

  footer .flex-wrap > * {
    flex: 1 1 calc(50% - 0.625rem);
    justify-content: center;
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }
}

@media (max-width: 380px) {
  .liquid-nav .liquid-action {
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }

  .hero-title {
    font-size: clamp(2.25rem, 11.2vw, 2.85rem);
  }

  .hero-actions > a,
  .auth-panel {
    width: 100%;
  }
}
</style>
