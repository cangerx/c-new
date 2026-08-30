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
import { act, render } from '@testing-library/react'
import { afterEach, beforeEach, describe, expect, test, vi } from 'vitest'

import { useCounterAnimation, useRevealOnScroll } from '../use-scroll-effects'

interface ObserverRecord {
  callback: IntersectionObserverCallback
  observe: ReturnType<typeof vi.fn>
  options?: IntersectionObserverInit
}

const observerRecords: ObserverRecord[] = []

class IntersectionObserverMock {
  readonly root = null
  readonly rootMargin = '0px'
  readonly thresholds = [0]
  readonly observe = vi.fn()
  readonly unobserve = vi.fn()
  readonly disconnect = vi.fn()
  readonly takeRecords = () => []

  constructor(
    callback: IntersectionObserverCallback,
    options?: IntersectionObserverInit
  ) {
    observerRecords.push({ callback, observe: this.observe, options })
  }
}

function ScrollEffectsHarness(props: { enabled: boolean }) {
  useRevealOnScroll(props.enabled)
  useCounterAnimation(props.enabled)

  return (
    <>
      <div className='reveal-element'>Content</div>
      <div className='counter-section'>
        <span className='counter-value' data-target='42' data-decimals='0'>
          0
        </span>
      </div>
    </>
  )
}

describe('landing scroll effects', () => {
  beforeEach(() => {
    observerRecords.length = 0
    vi.useFakeTimers()
    vi.stubGlobal('IntersectionObserver', IntersectionObserverMock)
  })

  afterEach(() => {
    vi.useRealTimers()
    vi.unstubAllGlobals()
  })

  test('registers landing elements when async homepage loading enables the effects', () => {
    const rendered = render(<ScrollEffectsHarness enabled={false} />)

    expect(observerRecords).toHaveLength(0)

    rendered.rerender(<ScrollEffectsHarness enabled />)

    expect(observerRecords).toHaveLength(2)

    act(() => {
      vi.advanceTimersByTime(100)
    })

    const revealElement = rendered.container.querySelector('.reveal-element')
    const counterSection = rendered.container.querySelector('.counter-section')
    const revealObserver = observerRecords.find(
      (record) => record.options?.threshold === 0.1
    )
    const counterObserver = observerRecords.find(
      (record) => record.options?.threshold === 0.3
    )

    expect(revealObserver?.observe).toHaveBeenCalledWith(revealElement)
    expect(counterObserver?.observe).toHaveBeenCalledWith(counterSection)
  })
})
