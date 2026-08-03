package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/pkg/billingexpr"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/QuantumNous/new-api/setting/billing_setting"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestModelPriceHelperTieredUsesPreloadedRequestInput(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := map[string]string{}
	require.NoError(t, config.GlobalConfig.SaveToDB(func(key, value string) error {
		saved[key] = value
		return nil
	}))
	t.Cleanup(func() {
		require.NoError(t, config.GlobalConfig.LoadFromDB(saved))
	})

	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"billing_setting.billing_mode": `{"tiered-test-model":"tiered_expr"}`,
		"billing_setting.billing_expr": `{"tiered-test-model":"param(\"stream\") == true ? tier(\"stream\", p * 3) : tier(\"base\", p * 2)"}`,
	}))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	req := httptest.NewRequest(http.MethodPost, "/api/channel/test/1", nil)
	req.Body = nil
	req.ContentLength = 0
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req
	ctx.Set("group", "default")

	info := &relaycommon.RelayInfo{
		OriginModelName: "tiered-test-model",
		UserGroup:       "default",
		UsingGroup:      "default",
		RequestHeaders:  map[string]string{"Content-Type": "application/json"},
		BillingRequestInput: &billingexpr.RequestInput{
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    []byte(`{"stream":true}`),
		},
	}

	priceData, err := ModelPriceHelper(ctx, info, 1000, &types.TokenCountMeta{
		BillingRatios: map[string]float64{"n": 3},
	})
	require.NoError(t, err)
	require.Equal(t, 1500, priceData.QuotaToPreConsume)
	require.NotNil(t, info.TieredBillingSnapshot)
	require.Equal(t, "stream", info.TieredBillingSnapshot.EstimatedTier)
	require.Equal(t, billing_setting.BillingModeTieredExpr, info.TieredBillingSnapshot.BillingMode)
	require.Equal(t, common.QuotaPerUnit, info.TieredBillingSnapshot.QuotaPerUnit)
}

func TestModelPriceHelperTieredPreConsumeMaxTokensFallback(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := map[string]string{}
	require.NoError(t, config.GlobalConfig.SaveToDB(func(key, value string) error {
		saved[key] = value
		return nil
	}))
	t.Cleanup(func() {
		require.NoError(t, config.GlobalConfig.LoadFromDB(saved))
	})

	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"billing_setting.billing_mode":    `{"tiered-fallback-model":"tiered_expr"}`,
		"billing_setting.billing_expr":    `{"tiered-fallback-model":"tier(\"base\", p * 3 + c * 15)"}`,
		"group_ratio_setting.group_ratio": `{"default":1,"free":0}`,
	}))

	const promptTokens = 1000

	cases := []struct {
		name      string
		group     string
		maxTokens int
		expected  int
	}{
		{
			// max_tokens omitted in a paid group -> fall back to 8192 completion tokens.
			// p*3 + c*15 = 1000*3 + 8192*15 = 125880 -> /1e6 * 500000 = 62940
			name:      "non-free group falls back to 8192 completion tokens",
			group:     "default",
			maxTokens: 0,
			expected:  62940,
		},
		{
			// explicit max_tokens is used verbatim, no fallback.
			// 1000*3 + 100*15 = 4500 -> /1e6 * 500000 = 2250
			name:      "explicit max_tokens is used verbatim",
			group:     "default",
			maxTokens: 100,
			expected:  2250,
		},
		{
			// free group (ratio 0) stays zero; fallback is gated on non-zero group ratio.
			name:      "free group stays zero without fallback",
			group:     "free",
			maxTokens: 0,
			expected:  0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			req := httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
			req.Header.Set("Content-Type", "application/json")
			ctx.Request = req
			ctx.Set("group", tc.group)

			info := &relaycommon.RelayInfo{
				OriginModelName: "tiered-fallback-model",
				UserGroup:       tc.group,
				UsingGroup:      tc.group,
				RequestHeaders:  map[string]string{"Content-Type": "application/json"},
				BillingRequestInput: &billingexpr.RequestInput{
					Headers: map[string]string{"Content-Type": "application/json"},
					Body:    []byte(`{}`),
				},
			}

			priceData, err := ModelPriceHelper(ctx, info, promptTokens, &types.TokenCountMeta{MaxTokens: tc.maxTokens})
			require.NoError(t, err)
			require.Equal(t, tc.expected, priceData.QuotaToPreConsume)
		})
	}
}

func TestModelPriceHelperTieredRejectsPreConsumeOverflow(t *testing.T) {
	gin.SetMode(gin.TestMode)

	saved := map[string]string{}
	require.NoError(t, config.GlobalConfig.SaveToDB(func(key, value string) error {
		saved[key] = value
		return nil
	}))
	t.Cleanup(func() {
		require.NoError(t, config.GlobalConfig.LoadFromDB(saved))
	})

	require.NoError(t, config.GlobalConfig.LoadFromDB(map[string]string{
		"billing_setting.billing_mode":    `{"tiered-overflow-model":"tiered_expr"}`,
		"billing_setting.billing_expr":    `{"tiered-overflow-model":"tier(\"overflow\", p * 1000000000000000)"}`,
		"group_ratio_setting.group_ratio": `{"default":1}`,
	}))

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/chat/completions", nil)
	ctx.Set("group", "default")
	info := &relaycommon.RelayInfo{
		OriginModelName: "tiered-overflow-model",
		UserGroup:       "default",
		UsingGroup:      "default",
		BillingRequestInput: &billingexpr.RequestInput{
			Body: []byte(`{}`),
		},
	}

	_, err := ModelPriceHelper(ctx, info, 1000, &types.TokenCountMeta{})

	var clamp *common.QuotaClamp
	require.ErrorAs(t, err, &clamp)
	require.Equal(t, "QuotaRound", clamp.Op)
	require.Equal(t, common.QuotaClampOverflow, clamp.Kind)
}

func TestModelPriceHelperRequestBillingRatiosOnlyApplyToFixedPrice(t *testing.T) {
	gin.SetMode(gin.TestMode)
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	savedModelRatios := ratio_setting.ModelRatio2JSONString()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(savedModelRatios))
	})

	modelPrices, err := common.Marshal(map[string]float64{
		"fixed-image-price":      0.04,
		"fractional-image-price": 0.0000012,
		"overflow-image-price":   float64(common.MaxQuota) / common.QuotaPerUnit / 2,
	})
	require.NoError(t, err)
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(string(modelPrices)))
	modelRatios, err := common.Marshal(map[string]float64{"ratio-image-price": 15})
	require.NoError(t, err)
	require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(string(modelRatios)))

	tests := []struct {
		name           string
		model          string
		wantQuota      int
		wantUsePrice   bool
		wantImageCount bool
	}{
		{
			name:           "fixed price applies image count",
			model:          "fixed-image-price",
			wantQuota:      180000,
			wantUsePrice:   true,
			wantImageCount: true,
		},
		{
			name:         "ratio price ignores request billing ratios",
			model:        "ratio-image-price",
			wantQuota:    15000,
			wantUsePrice: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
			ctx.Set("group", "default")
			info := &relaycommon.RelayInfo{
				OriginModelName: tt.model,
				UserGroup:       "default",
				UsingGroup:      "default",
			}
			meta := &types.TokenCountMeta{
				ImagePriceRatio: 3,
				BillingRatios:   map[string]float64{"n": 3},
			}

			priceData, err := ModelPriceHelper(ctx, info, 1000, meta)

			require.NoError(t, err)
			require.Equal(t, tt.wantQuota, priceData.QuotaToPreConsume)
			require.Equal(t, tt.wantUsePrice, priceData.UsePrice)
			require.Equal(t, tt.wantImageCount, priceData.HasOtherRatio("n"))
			require.Equal(t, priceData.OtherRatios(), info.PriceData.OtherRatios())
		})
	}

	newInfo := func(model string) (*gin.Context, *relaycommon.RelayInfo) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("group", "default")
		return ctx, &relaycommon.RelayInfo{
			OriginModelName: model,
			UserGroup:       "default",
			UsingGroup:      "default",
		}
	}
	meta := &types.TokenCountMeta{BillingRatios: map[string]float64{"n": 3}}

	ctx, info := newInfo("fractional-image-price")
	priceData, err := ModelPriceHelper(ctx, info, 0, meta)
	require.NoError(t, err)
	// 0.0000012 * 500000 * 3 = 1.8, then truncate once to 1.
	require.Equal(t, 1, priceData.QuotaToPreConsume)

	ctx, info = newInfo("overflow-image-price")
	_, err = ModelPriceHelper(ctx, info, 0, meta)
	var clamp *common.QuotaClamp
	require.ErrorAs(t, err, &clamp)
	require.Equal(t, "QuotaFromFloat", clamp.Op)
	require.Equal(t, common.QuotaClampOverflow, clamp.Kind)
	require.Nil(t, info.Billing)
}

func TestModelPriceHelperPerCallDefaultTaskBilling(t *testing.T) {
	gin.SetMode(gin.TestMode)
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	savedModelRatios := ratio_setting.ModelRatio2JSONString()
	savedMode := operation_setting.GetQuotaSetting().DefaultTaskBillingMode
	savedPrice := operation_setting.GetQuotaSetting().DefaultTaskPrice
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(savedModelRatios))
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = savedMode
		operation_setting.GetQuotaSetting().DefaultTaskPrice = savedPrice
	})

	// 清除所有模型价格/倍率配置，只留系统默认
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString("{}"))
	require.NoError(t, ratio_setting.UpdateModelRatioByJSONString("{}"))

	newInfo := func(model string) (*gin.Context, *relaycommon.RelayInfo) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("group", "default")
		return ctx, &relaycommon.RelayInfo{
			OriginModelName: model,
			UserGroup:       "default",
			UsingGroup:      "default",
		}
	}

	t.Run("per_call default applies fixed price once", func(t *testing.T) {
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = "per_call"
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.1

		ctx, info := newInfo("unconfigured-model")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.True(t, priceData.UsePrice)
		// 0.1 * 500000 = 50000，固定价不随 seconds 放大
		require.Equal(t, 50000, priceData.Quota)

		// 模拟 adaptor.EstimateBilling 注入 seconds 倍率
		priceData.AddOtherRatio("seconds", 5)
		require.Equal(t, 50000, priceData.Quota)
	})

	t.Run("per_second default applies ratio times seconds", func(t *testing.T) {
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = "per_second"
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.04

		ctx, info := newInfo("unconfigured-model")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.False(t, priceData.UsePrice)
		// 基础预扣：0.04 / 2 * 500000 = 10000
		require.Equal(t, 10000, priceData.Quota)

		// 模拟 relay_task step 6：按量才把 OtherRatios 乘入
		priceData.AddOtherRatio("seconds", 5)
		quotaWithRatios := priceData.ApplyOtherRatiosToFloat(float64(priceData.Quota))
		require.Equal(t, 50000.0, quotaWithRatios)
	})

	t.Run("default price zero keeps error for unconfigured model", func(t *testing.T) {
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = "per_call"
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0

		ctx, info := newInfo("unconfigured-model")
		_, err := ModelPriceHelperPerCall(ctx, info)

		require.Error(t, err)
		require.Contains(t, err.Error(), "has not been priced")
	})
}

func TestModelPriceHelperPerCallExplicitMode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	savedModelRatios := ratio_setting.ModelRatio2JSONString()
	savedMode := operation_setting.GetQuotaSetting().DefaultTaskBillingMode
	savedPrice := operation_setting.GetQuotaSetting().DefaultTaskPrice
	savedTaskMode := billing_setting.GetTaskBillingModeCopy()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(savedModelRatios))
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = savedMode
		operation_setting.GetQuotaSetting().DefaultTaskPrice = savedPrice
		bs := config.GlobalConfig.Get("billing_setting").(*billing_setting.BillingSetting)
		bs.TaskBillingMode = savedTaskMode
	})

	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString("{}"))
	require.NoError(t, ratio_setting.UpdateModelRatioByJSONString("{}"))

	bs := config.GlobalConfig.Get("billing_setting").(*billing_setting.BillingSetting)
	bs.TaskBillingMode = map[string]string{
		"model-per-call":   billing_setting.TaskBillingModePerCall,
		"model-per-second": billing_setting.TaskBillingModePerSecond,
	}

	newInfo := func(model string) (*gin.Context, *relaycommon.RelayInfo) {
		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("group", "default")
		return ctx, &relaycommon.RelayInfo{
			OriginModelName: model,
			UserGroup:       "default",
			UsingGroup:      "default",
		}
	}

	t.Run("per_call with model price uses model price", func(t *testing.T) {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{"model-per-call": 0.2}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.1

		ctx, info := newInfo("model-per-call")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.True(t, priceData.UsePrice)
		require.Equal(t, 100000, priceData.Quota) // 0.2 * 500000
	})

	t.Run("per_call without model price falls back to system default", func(t *testing.T) {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.1

		ctx, info := newInfo("model-per-call")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.True(t, priceData.UsePrice)
		require.Equal(t, 50000, priceData.Quota) // 0.1 * 500000
	})

	t.Run("per_call falls back to built-in default model price", func(t *testing.T) {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.5
		bs.TaskBillingMode = map[string]string{
			"sora-2": billing_setting.TaskBillingModePerCall,
		}

		ctx, info := newInfo("sora-2")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.True(t, priceData.UsePrice)
		require.Equal(t, 150000, priceData.Quota) // 内置默认价 0.3 * 500000，优先于系统默认 0.5
	})

	t.Run("per_second with model ratio uses model ratio", func(t *testing.T) {
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{"model-per-second": 0.04}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.08

		ctx, info := newInfo("model-per-second")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.False(t, priceData.UsePrice)
		require.Equal(t, 10000, priceData.Quota) // 0.04 / 2 * 500000 预扣一半
	})

	t.Run("per_second without model ratio uses system default not 37.5", func(t *testing.T) {
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.08

		ctx, info := newInfo("model-per-second")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.False(t, priceData.UsePrice)
		// 回落系统默认价 0.08，而不是 GetModelRatio 的 37.5 兜底
		require.Equal(t, 20000, priceData.Quota) // 0.08 / 2 * 500000
	})

	t.Run("unconfigured mode keeps legacy fallback", func(t *testing.T) {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(`{}`))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(`{}`))
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = "per_second"
		operation_setting.GetQuotaSetting().DefaultTaskPrice = 0.04

		ctx, info := newInfo("no-mode-model")
		priceData, err := ModelPriceHelperPerCall(ctx, info)

		require.NoError(t, err)
		require.False(t, priceData.UsePrice)
		require.Equal(t, 10000, priceData.Quota) // 系统默认 per_second 0.04 / 2 * 500000
	})
}
