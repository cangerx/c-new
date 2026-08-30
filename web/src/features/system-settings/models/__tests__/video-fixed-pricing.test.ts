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
import { describe, expect, test } from 'vitest'

import {
  buildModelSnapshots,
  isBasePricingUnset,
} from '../model-pricing-snapshots'

describe('video fixed pricing', () => {
  test('keeps video per-request pricing separate from the existing fixed request price', () => {
    const rows = buildModelSnapshots({
      modelPrice: '{"standard-fixed":0.1}',
      videoModelPrice: '{"video-fixed":0.25}',
      modelRatio: '{}',
      cacheRatio: '{}',
      createCacheRatio: '{}',
      completionRatio: '{}',
      imageRatio: '{}',
      audioRatio: '{}',
      audioCompletionRatio: '{}',
      billingMode: '{}',
      billingExpr: '{}',
    })

    const standard = rows.find((row) => row.name === 'standard-fixed')
    const video = rows.find((row) => row.name === 'video-fixed')

    expect(standard?.billingMode).toBe('per-request')
    expect(standard?.price).toBe('0.1')
    expect(standard?.videoPrice).toBe('')
    expect(video?.billingMode).toBe('video-per-request')
    expect(video?.videoPrice).toBe('0.25')
    expect(video?.price).toBe('')
    expect(isBasePricingUnset(video)).toBe(false)
  })
})
