package model

import (
	"testing"

	"github.com/QuantumNous/new-api/setting/billing_setting"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 模型广场的定价类型必须反映真实计费模式。按秒计费的每秒价格与 per-token
// 倍率共用 ModelRatio 这一个 map，只能靠 BillingMode 区分；若不区分，按秒
// 模型会在广场上显示成「按 Token」并被换算成每 1M token 价格。
func TestApplyTaskBillingPricingReflectsBilledMode(t *testing.T) {
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	savedModelRatios := ratio_setting.ModelRatio2JSONString()
	savedDefaultPrice := operation_setting.GetQuotaSetting().DefaultTaskPrice
	savedTaskMode := billing_setting.GetTaskBillingModeCopy()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(savedModelRatios))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = savedDefaultPrice
		config.GlobalConfig.Get("billing_setting").(*billing_setting.BillingSetting).TaskBillingMode = savedTaskMode
	})

	const (
		perCallModel   = "videos-4-mini-480p"
		perSecondModel = "videos-4-mini-720p"
	)

	setTaskModes := func(modes map[string]string) {
		config.GlobalConfig.Get("billing_setting").(*billing_setting.BillingSetting).TaskBillingMode = modes
	}

	t.Run("按次: 价格进 ModelPrice, QuotaType=1", func(t *testing.T) {
		setTaskModes(map[string]string{perCallModel: billing_setting.TaskBillingModePerCall})
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{"`+perCallModel+`": 0.5}`))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{}`))

		var pricing Pricing
		require.True(t, applyTaskBillingPricing(perCallModel, &pricing))

		assert.Equal(t, 1, pricing.QuotaType)
		assert.Equal(t, 0.5, pricing.ModelPrice)
		assert.Equal(t, billing_setting.TaskBillingModePerCall, pricing.BillingMode)
		assert.Zero(t, pricing.ModelRatio, "按次不得写 ModelRatio")
	})

	t.Run("按秒: 每秒价格进 ModelRatio 且标记 BillingMode", func(t *testing.T) {
		setTaskModes(map[string]string{perSecondModel: billing_setting.TaskBillingModePerSecond})
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{}`))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{"`+perSecondModel+`": 0.04}`))

		var pricing Pricing
		require.True(t, applyTaskBillingPricing(perSecondModel, &pricing))

		// QuotaType 仍是 0，前端必须凭 BillingMode 才能与按 Token 区分开
		assert.Equal(t, 0, pricing.QuotaType)
		assert.Equal(t, 0.04, pricing.ModelRatio)
		assert.Equal(t, billing_setting.TaskBillingModePerSecond, pricing.BillingMode)
		assert.Zero(t, pricing.ModelPrice, "按秒不得写 ModelPrice")
		// 没有 BillingMode 就无法与 per-token 倍率 0.04 区分，这正是原 bug
		assert.NotEmpty(t, pricing.BillingMode)
	})

	t.Run("未配置模式时不接管, 交回原有反推逻辑", func(t *testing.T) {
		setTaskModes(map[string]string{})
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{"gpt-4o": 1.25}`))

		var pricing Pricing
		assert.False(t, applyTaskBillingPricing("gpt-4o", &pricing))
		assert.Empty(t, pricing.BillingMode)
		assert.Zero(t, pricing.ModelRatio, "未接管时不应写入任何定价字段")
	})

	t.Run("按次未配价时回落系统默认", func(t *testing.T) {
		setTaskModes(map[string]string{perCallModel: billing_setting.TaskBillingModePerCall})
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.3

		var pricing Pricing
		require.True(t, applyTaskBillingPricing(perCallModel, &pricing))

		assert.Equal(t, 1, pricing.QuotaType)
		assert.Equal(t, 0.3, pricing.ModelPrice)
	})

	t.Run("按秒未配倍率时回落系统默认, 而非 GetModelRatio 的 37.5", func(t *testing.T) {
		setTaskModes(map[string]string{perSecondModel: billing_setting.TaskBillingModePerSecond})
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.08

		var pricing Pricing
		require.True(t, applyTaskBillingPricing(perSecondModel, &pricing))

		assert.Equal(t, 0.08, pricing.ModelRatio)
		assert.NotEqual(t, 37.5, pricing.ModelRatio,
			"不得使用 GetModelRatio 的自用模式兜底值")
	})

	t.Run("系统默认价为 0 时不接管", func(t *testing.T) {
		setTaskModes(map[string]string{perSecondModel: billing_setting.TaskBillingModePerSecond})
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0

		var pricing Pricing
		assert.False(t, applyTaskBillingPricing(perSecondModel, &pricing))
		assert.Empty(t, pricing.BillingMode)
	})
}
