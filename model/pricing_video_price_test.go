package model

import (
	"testing"

	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/setting/ratio_setting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPricingExposesDedicatedVideoFixedPrice(t *testing.T) {
	resetPricingEndpointTestTables(t)

	savedVideoPrices := ratio_setting.VideoModelPrice2JSONString()
	t.Cleanup(func() {
		require.NoError(t, ratio_setting.UpdateVideoModelPriceByJSONString(savedVideoPrices))
		InvalidatePricingCache()
	})

	require.NoError(t, ratio_setting.UpdateVideoModelPriceByJSONString(`{"video-marketplace-model":0.25}`))
	insertPricingEndpointChannel(t, 501, constant.ChannelTypeOpenAI, dto.ChannelOtherSettings{})
	insertPricingEndpointAbility(t, 501, "video-marketplace-model")
	InitChannelCache()

	var got *Pricing
	for _, pricing := range GetPricing() {
		if pricing.ModelName == "video-marketplace-model" {
			pricingCopy := pricing
			got = &pricingCopy
			break
		}
	}

	require.NotNil(t, got)
	assert.Equal(t, 1, got.QuotaType)
	assert.Equal(t, "video_per_request", got.BillingMode)
	require.NotNil(t, got.VideoModelPrice)
	assert.Equal(t, 0.25, *got.VideoModelPrice)
	assert.Zero(t, got.ModelPrice)
}
