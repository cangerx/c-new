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
import assert from 'node:assert/strict'
import { describe, test } from 'node:test'

import {
  getModeLabel,
  getPriceDetail,
  getPriceSummary,
  type ModelPricingSnapshot,
} from '../model-pricing-snapshots'

const t = (key: string) => key

const snapshot = (
  overrides: Partial<ModelPricingSnapshot>
): ModelPricingSnapshot => ({
  name: 'sora-2',
  hasConflict: false,
  ...overrides,
})

// per_second 的价格存在 ModelRatio，per-token 倍率也存在 ModelRatio。
// 只有 taskBillingMode 能区分两者，所以展示层必须优先看它，否则按秒模型
// 会被显示成「按 Token」并按 per-1M-token 换算出一个无意义的价格。
describe('task billing mode drives the pricing display', () => {
  test('per_second reads ModelRatio as dollars per second, not a token ratio', () => {
    const row = snapshot({ ratio: '0.04', taskBillingMode: 'per_second' })

    assert.equal(
      getModeLabel(row.billingMode, row.taskBillingMode),
      'Per second'
    )
    assert.equal(getPriceSummary(row, t), '$0.04 / second')
    assert.equal(getPriceDetail(row, t), 'Priced per second of output')
  })

  test('per_call reads ModelPrice', () => {
    const row = snapshot({ price: '0.5', taskBillingMode: 'per_call' })

    assert.equal(getModeLabel(row.billingMode, row.taskBillingMode), 'Per call')
    assert.equal(getPriceSummary(row, t), '$0.5 / call')
    assert.equal(getPriceDetail(row, t), 'Fixed price per task')
  })

  test('an explicit task mode with no price falls back to the system default', () => {
    assert.equal(
      getPriceSummary(snapshot({ taskBillingMode: 'per_second' }), t),
      'System default'
    )
    assert.equal(
      getPriceSummary(snapshot({ taskBillingMode: 'per_call' }), t),
      'System default'
    )
  })

  test('the same ModelRatio without a task mode stays a per-token ratio', () => {
    const row = snapshot({ ratio: '0.04' })

    assert.equal(
      getModeLabel(row.billingMode, row.taskBillingMode),
      'Per-token'
    )
    assert.equal(getPriceSummary(row, t), 'Input $0.08')
  })

  test('tiered_expr outranks a stale task mode', () => {
    const row = snapshot({
      billingMode: 'tiered_expr',
      billingExpr: 'tier(0, 100, 1.0)',
      taskBillingMode: 'per_second',
    })

    assert.equal(getPriceSummary(row, t), 'Tiered pricing · 1 tiers')
  })
})
