// 外部测试包：sora adaptor 间接依赖 relay/helper，
// 放在 package helper 里会成导入环，helper_test 单独编译可打破。
package helper_test

import (
	"net/http/httptest"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/relay/channel/task/sora"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relay/helper"
	"github.com/QuantumNous/new-api/setting/billing_setting"
	"github.com/QuantumNous/new-api/setting/config"
	"github.com/QuantumNous/new-api/setting/operation_setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	hosttypes "github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 按次计费的固定价格不得再乘时长，按秒计费必须乘时长。
// 这是 relay_task.go 里 `if !info.PriceData.UsePrice` 那道闸门的契约：
// ModelPriceHelperPerCall 决定 UsePrice，闸门决定 OtherRatios 是否进 quota。
// 两者错配会让「按次」随时长涨价，或让「按秒」无视时长收固定价。
func TestTaskBillingPerCallVsPerSecondAppliesSecondsCorrectly(t *testing.T) {
	gin.SetMode(gin.TestMode)

	savedModelPrices := ratio_setting.ModelPrice2JSONString()
	savedModelRatios := ratio_setting.ModelRatio2JSONString()
	savedMode := operation_setting.GetQuotaSetting().DefaultTaskBillingMode
	savedPrice := operation_setting.GetQuotaSetting().DefaultTaskPrice
	savedTaskMode := billing_setting.GetTaskBillingModeCopy()
	savedPatches := constant.TaskPricePatches
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(savedModelPrices))
		require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(savedModelRatios))
		operation_setting.GetQuotaSetting().DefaultTaskBillingMode = savedMode
		operation_setting.GetQuotaSetting().DefaultTaskPrice = savedPrice
		config.GlobalConfig.Get("billing_setting").(*billing_setting.BillingSetting).TaskBillingMode = savedTaskMode
		constant.TaskPricePatches = savedPatches
	})

	const (
		perCallModel   = "videos-4-mini-480p"
		perSecondModel = "videos-4-mini-720p"
		perCallPrice   = 0.5  // $/次
		perSecondPrice = 0.04 // $/秒
	)

	// videos-4-mini-480p → 按次；videos-4-mini-720p → 按秒
	config.GlobalConfig.Get("billing_setting").(*billing_setting.BillingSetting).TaskBillingMode = map[string]string{
		perCallModel:   billing_setting.TaskBillingModePerCall,
		perSecondModel: billing_setting.TaskBillingModePerSecond,
	}
	require.NoError(t, ratio_setting.UpdateModelPriceByJSONString(
		`{"`+perCallModel+`": 0.5}`))
	require.NoError(t, ratio_setting.UpdateModelRatioByJSONString(
		`{"`+perSecondModel+`": 0.04}`))
	// 系统默认设成明显不同的值，确保下面断言读的是模型级配置而非兜底。
	operation_setting.GetQuotaSetting().DefaultTaskBillingMode = "per_second"
	operation_setting.GetQuotaSetting().DefaultTaskPrice = 9.99
	constant.TaskPricePatches = nil

	// submitTask 复刻 RelayTaskSubmit 的计费段（步骤 5-6）：
	// ModelPriceHelperPerCall → EstimateBilling → UsePrice 闸门。
	submitTask := func(t *testing.T, model, seconds, size string) hosttypes.PriceData {
		t.Helper()

		ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctx.Set("group", "default")
		info := &relaycommon.RelayInfo{
			OriginModelName: model,
			UserGroup:       "default",
			UsingGroup:      "default",
			// Action 是 *TaskRelayInfo 提升上来的字段，生产代码在
			// InitChannelMeta 里初始化；不初始化会在 EstimateBilling 里 nil panic。
			TaskRelayInfo: &relaycommon.TaskRelayInfo{},
		}

		priceData, err := helper.ModelPriceHelperPerCall(ctx, info)
		require.NoError(t, err)
		info.PriceData = priceData

		// 走真实 adaptor 的 EstimateBilling，而不是自己编 seconds 比率。
		ctx.Set("task_request", relaycommon.TaskSubmitReq{
			Model:   model,
			Prompt:  "a cat",
			Seconds: seconds,
			Size:    size,
		})
		adaptor := &sora.TaskAdaptor{}
		if ratios := adaptor.EstimateBilling(ctx, info); len(ratios) > 0 {
			for k, v := range ratios {
				info.PriceData.AddOtherRatio(k, v)
			}
		}

		// relay_task.go:199 的闸门原样复刻
		if !info.PriceData.UsePrice &&
			!common.StringsContains(constant.TaskPricePatches, model) {
			quotaWithRatios := info.PriceData.ApplyOtherRatiosToFloat(float64(info.PriceData.Quota))
			quota, _ := common.QuotaFromFloatChecked(quotaWithRatios)
			info.PriceData.Quota = quota
		}
		return info.PriceData
	}

	t.Run("按次: 时长不参与计费", func(t *testing.T) {
		short := submitTask(t, perCallModel, "4", "720x1280")
		long := submitTask(t, perCallModel, "12", "720x1280")

		assert.True(t, short.UsePrice, "按次必须走 UsePrice")
		// 0.5 * 500000 = 250000，与秒数无关
		wantQuota := int(perCallPrice * common.QuotaPerUnit)
		assert.Equal(t, wantQuota, short.Quota)
		assert.Equal(t, wantQuota, long.Quota, "按次计费下 4 秒和 12 秒必须同价")
	})

	t.Run("按秒: 额度随时长线性增长", func(t *testing.T) {
		four := submitTask(t, perSecondModel, "4", "720x1280")
		eight := submitTask(t, perSecondModel, "8", "720x1280")

		assert.False(t, four.UsePrice, "按秒必须走 ratio 而非 UsePrice")
		// 预扣 = ratio/2 * QuotaPerUnit * seconds
		assert.Equal(t, int(perSecondPrice/2*common.QuotaPerUnit*4), four.Quota)
		assert.Equal(t, int(perSecondPrice/2*common.QuotaPerUnit*8), eight.Quota)
		assert.Equal(t, four.Quota*2, eight.Quota, "时长翻倍则额度翻倍")
	})

	t.Run("同一时长下两种模式各自正确", func(t *testing.T) {
		call := submitTask(t, perCallModel, "8", "720x1280")
		second := submitTask(t, perSecondModel, "8", "720x1280")

		assert.Equal(t, 250000, call.Quota)  // 0.5/次
		assert.Equal(t, 80000, second.Quota) // 0.04/秒 × 8 秒，预扣半价
	})

	t.Run("按次无视尺寸比率, 按秒受尺寸比率影响", func(t *testing.T) {
		callWide := submitTask(t, perCallModel, "8", "1792x1024")
		secondWide := submitTask(t, perSecondModel, "8", "1792x1024")

		assert.Equal(t, 250000, callWide.Quota, "按次固定价不受尺寸影响")
		assert.Equal(t, 133333, secondWide.Quota) // 80000 * 1.666667
	})
}
