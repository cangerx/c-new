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
import { useEffect, useRef } from 'react'

type PaperShadersModule = typeof import('@paper-design/shaders')
type HeroShaderInstance = InstanceType<PaperShadersModule['ShaderMount']>

const SHADER_SPEED = 1.85

function getHeroShaderTheme(isDark: boolean) {
  return isDark
    ? {
        colorBack: '#030305',
        colors: ['#00D5FF', '#7C3AED', '#FF2E88', '#2563EB00'],
        intensity: 1,
        noise: 0.72,
      }
    : {
        colorBack: '#F7FAFF',
        colors: ['#0EA5E9', '#8B5CF6', '#F43F5E', '#22D3EE00'],
        intensity: 0.96,
        noise: 0.5,
      }
}

function waitForShaderImage(image: HTMLImageElement) {
  if (image.complete && image.naturalWidth > 0) return Promise.resolve()

  return new Promise<void>((resolve, reject) => {
    image.addEventListener('load', () => resolve(), { once: true })
    image.addEventListener(
      'error',
      () => reject(new Error('Paper shader noise texture failed to load')),
      { once: true }
    )
  })
}

function prefersReducedMotion() {
  return window.matchMedia('(prefers-reduced-motion: reduce)').matches
}

/**
 * Mounts the WebGL grain-gradient background onto a plain DOM node.
 *
 * `@paper-design/shaders` is framework-agnostic (ShaderMount attaches to an
 * element), so this is the Vue original's logic moved into an effect. Init is
 * deferred to idle time because compiling the shader blocks the main thread,
 * and skipped entirely under prefers-reduced-motion.
 */
export function useHeroShader(isDark: boolean) {
  const containerRef = useRef<HTMLDivElement | null>(null)
  const shaderRef = useRef<HeroShaderInstance | null>(null)
  const moduleRef = useRef<PaperShadersModule | null>(null)

  // Mount once. `isDark` is applied through the separate uniform effect below
  // so theme flips never rebuild the WebGL context.
  useEffect(() => {
    const container = containerRef.current
    if (!container || prefersReducedMotion()) return

    let cancelled = false
    let idleHandle: number | null = null
    let timeoutHandle: number | null = null

    const init = async () => {
      try {
        const shaderModule = await import('@paper-design/shaders')
        if (cancelled) return
        moduleRef.current = shaderModule

        const {
          ShaderFitOptions,
          ShaderMount,
          defaultObjectSizing,
          grainGradientFragmentShader,
          GrainGradientShapes,
          getShaderColorFromString,
          getShaderNoiseTexture,
        } = shaderModule

        const noiseTexture = getShaderNoiseTexture()
        if (noiseTexture) await waitForShaderImage(noiseTexture)
        if (cancelled) return

        const sizing = defaultObjectSizing
        const theme = getHeroShaderTheme(isDark)

        shaderRef.current = new ShaderMount(
          container,
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
            u_noiseTexture: noiseTexture,
          },
          { alpha: true, antialias: true, premultipliedAlpha: false },
          SHADER_SPEED,
          0,
          1,
          1920 * 1080 * 2
        )
        container.classList.add('is-ready')
      } catch (error) {
        // A missing WebGL context must not take the landing page down.
        // eslint-disable-next-line no-console
        console.warn('[landing] Paper shader background disabled:', error)
        container.classList.remove('is-ready')
        shaderRef.current = null
      }
    }

    if (typeof window.requestIdleCallback === 'function') {
      idleHandle = window.requestIdleCallback(() => void init(), {
        timeout: 1800,
      })
    } else {
      timeoutHandle = window.setTimeout(() => void init(), 700)
    }

    return () => {
      cancelled = true
      if (
        idleHandle !== null &&
        typeof window.cancelIdleCallback === 'function'
      ) {
        window.cancelIdleCallback(idleHandle)
      }
      if (timeoutHandle !== null) window.clearTimeout(timeoutHandle)
      shaderRef.current?.dispose()
      shaderRef.current = null
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  // Recolor in place on theme change.
  useEffect(() => {
    const shader = shaderRef.current
    const shaderModule = moduleRef.current
    if (!shader || !shaderModule) return

    const { getShaderColorFromString } = shaderModule
    const theme = getHeroShaderTheme(isDark)
    shader.setUniforms({
      u_colorBack: getShaderColorFromString(theme.colorBack),
      u_colors: theme.colors.map(getShaderColorFromString),
      u_colorsCount: theme.colors.length,
      u_intensity: theme.intensity,
      u_noise: theme.noise,
    })
    shader.setSpeed(SHADER_SPEED)
  }, [isDark])

  return containerRef
}
