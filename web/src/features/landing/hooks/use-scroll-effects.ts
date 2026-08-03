import { useEffect, useState } from 'react'

/** Nav collapses into a pill once the hero starts scrolling away. */
export function useIsScrolled(threshold = 80) {
  const [isScrolled, setIsScrolled] = useState(false)

  useEffect(() => {
    const handleScroll = () => setIsScrolled(window.scrollY > threshold)
    handleScroll()
    window.addEventListener('scroll', handleScroll, { passive: true })
    return () => window.removeEventListener('scroll', handleScroll)
  }, [threshold])

  return isScrolled
}

/**
 * Adds `reveal-active` to every `.reveal-element` as it enters the viewport.
 *
 * Kept as a DOM query rather than per-component refs to match the Vue original
 * exactly: the CSS keys off the class, and elements opt in by carrying
 * `.reveal-element` without needing to know about this hook.
 */
export function useRevealOnScroll() {
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return
          entry.target.classList.add('reveal-active')
          observer.unobserve(entry.target)
        })
      },
      { root: null, rootMargin: '0px', threshold: 0.1 }
    )

    // Defer so the first paint has the initial (hidden) state applied,
    // otherwise elements already in view skip their transition.
    const handle = window.setTimeout(() => {
      document
        .querySelectorAll('.reveal-element')
        .forEach((el) => observer.observe(el))
    }, 100)

    return () => {
      window.clearTimeout(handle)
      observer.disconnect()
    }
  }, [])
}

function animateCounter(el: HTMLElement, target: number, decimals: number) {
  const duration = 1800
  const startTime = performance.now()
  const easeOutExpo = (t: number) => (t === 1 ? 1 : 1 - Math.pow(2, -10 * t))
  let frame = 0

  const tick = (now: number) => {
    const progress = Math.min((now - startTime) / duration, 1)
    el.textContent = (easeOutExpo(progress) * target).toFixed(decimals)
    if (progress < 1) {
      frame = requestAnimationFrame(tick)
    } else {
      el.textContent = target.toFixed(decimals)
    }
  }

  frame = requestAnimationFrame(tick)
  return () => cancelAnimationFrame(frame)
}

/**
 * Counts `.counter-value` numbers up when their `.counter-section` scrolls in.
 * Targets are read from `data-target` / `data-decimals` on each element.
 */
export function useCounterAnimation() {
  useEffect(() => {
    const cancels: Array<() => void> = []

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (!entry.isIntersecting) return
          entry.target.querySelectorAll('.counter-value').forEach((el) => {
            const target = Number.parseFloat(
              el.getAttribute('data-target') || '0'
            )
            const decimals = Number.parseInt(
              el.getAttribute('data-decimals') || '0'
            )
            cancels.push(animateCounter(el as HTMLElement, target, decimals))
          })
          observer.unobserve(entry.target)
        })
      },
      { threshold: 0.3 }
    )

    const handle = window.setTimeout(() => {
      document
        .querySelectorAll('.counter-section')
        .forEach((el) => observer.observe(el))
    }, 100)

    return () => {
      window.clearTimeout(handle)
      observer.disconnect()
      cancels.forEach((cancel) => cancel())
    }
  }, [])
}
