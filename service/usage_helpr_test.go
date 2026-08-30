package service

import (
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizeUsageForBillingAcceptsAndNormalizesValidUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	info := &relaycommon.RelayInfo{}
	info.SetEstimatePromptTokens(100)

	got := SanitizeUsageForBilling(ctx, info, &dto.Usage{
		InputTokens:  120,
		OutputTokens: 30,
		TotalTokens:  150,
	}, nil)

	require.NotNil(t, got)
	assert.Equal(t, 120, got.PromptTokens)
	assert.Equal(t, 30, got.CompletionTokens)
	assert.Equal(t, 150, got.TotalTokens)
	assert.Empty(t, UsageAnomalies(ctx))
}

func TestSanitizeUsageForBillingRejectsMalformedOrImplausibleUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	fallback := &dto.Usage{PromptTokens: 100, CompletionTokens: 20, TotalTokens: 120}

	tests := []struct {
		name  string
		usage *dto.Usage
	}{
		{
			name:  "negative field",
			usage: &dto.Usage{PromptTokens: 100, CompletionTokens: -1, TotalTokens: 99},
		},
		{
			name:  "inconsistent total",
			usage: &dto.Usage{PromptTokens: 100, CompletionTokens: 20, TotalTokens: 999},
		},
		{
			name:  "negative detail",
			usage: &dto.Usage{PromptTokens: 100, CompletionTokens: 20, TotalTokens: 120, PromptTokensDetails: dto.InputTokenDetails{CachedTokens: -1}},
		},
		{
			name:  "absurd prompt",
			usage: &dto.Usage{PromptTokens: maxUpstreamPromptTokens + 1, CompletionTokens: 20, TotalTokens: maxUpstreamPromptTokens + 21},
		},
		{
			name:  "absurd completion",
			usage: &dto.Usage{PromptTokens: 100, CompletionTokens: maxUpstreamCompletionTokens + 1, TotalTokens: maxUpstreamCompletionTokens + 101},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			info := &relaycommon.RelayInfo{}
			info.SetEstimatePromptTokens(100)

			got := SanitizeUsageForBilling(ctx, info, tt.usage, fallback)

			require.NotNil(t, got)
			assert.Equal(t, 100, got.PromptTokens)
			assert.Equal(t, 20, got.CompletionTokens)
			assert.Equal(t, 120, got.TotalTokens)
			assert.NotEmpty(t, UsageAnomalies(ctx))
		})
	}
}

func TestSanitizeRealtimeUsageRejectsForgedSnapshot(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	local := &dto.RealtimeUsage{
		TotalTokens:  30,
		InputTokens:  20,
		OutputTokens: 10,
		InputTokenDetails: dto.InputTokenDetails{
			TextTokens: 20,
		},
		OutputTokenDetails: dto.OutputTokenDetails{
			TextTokens: 10,
		},
	}
	forged := &dto.RealtimeUsage{
		TotalTokens:  maxUpstreamPromptTokens + 10,
		InputTokens:  maxUpstreamPromptTokens,
		OutputTokens: 10,
	}

	got := SanitizeRealtimeUsage(ctx, forged, local)

	require.NotNil(t, got)
	assert.Equal(t, local, got)
	assert.NotEmpty(t, UsageAnomalies(ctx))
}

func TestSanitizeUsageForBillingPreservesAnthropicCacheSemantics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	info := &relaycommon.RelayInfo{}
	info.SetEstimatePromptTokens(1000)
	usage := &dto.Usage{
		PromptTokens:     100,
		CompletionTokens: 20,
		TotalTokens:      120,
		InputTokens:      1000,
		OutputTokens:     20,
		UsageSemantic:    dto.BillingUsageSemanticAnthropic,
		PromptTokensDetails: dto.InputTokenDetails{
			CachedTokens:         700,
			CachedCreationTokens: 200,
		},
	}

	got := SanitizeUsageForBilling(ctx, info, usage, nil)

	require.NotNil(t, got)
	assert.Equal(t, 100, got.PromptTokens)
	assert.Equal(t, 1000, got.InputTokens)
	assert.Equal(t, 700, got.PromptTokensDetails.CachedTokens)
	assert.Empty(t, UsageAnomalies(ctx))
}

func TestSanitizeUsageForBillingRejectsInflatedBillingUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	info := &relaycommon.RelayInfo{}
	info.SetEstimatePromptTokens(100)
	outer := &dto.Usage{
		PromptTokens:     100,
		CompletionTokens: 10,
		TotalTokens:      110,
		BillingUsage: dto.NewOpenAIChatBillingUsage(&dto.Usage{
			PromptTokens:     maxUpstreamPromptTokens + 1,
			CompletionTokens: 10,
			TotalTokens:      maxUpstreamPromptTokens + 11,
		}),
	}

	got := SanitizeUsageForBilling(ctx, info, effectiveBillingUsage(outer), outer)

	require.NotNil(t, got)
	assert.Equal(t, 100, got.PromptTokens)
	assert.Equal(t, 10, got.CompletionTokens)
	assert.NotEmpty(t, UsageAnomalies(ctx))
}

func TestSanitizeUsageForBillingRejectsInflatedDetailWithSmallParent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	info := &relaycommon.RelayInfo{}
	info.SetEstimatePromptTokens(10)

	got := SanitizeUsageForBilling(ctx, info, &dto.Usage{
		PromptTokens:     10,
		CompletionTokens: 1,
		TotalTokens:      11,
		PromptTokensDetails: dto.InputTokenDetails{
			AudioTokens: 1000,
		},
	}, &dto.Usage{PromptTokens: 10, CompletionTokens: 1, TotalTokens: 11})

	require.NotNil(t, got)
	assert.Zero(t, got.PromptTokensDetails.AudioTokens)
	assert.NotEmpty(t, UsageAnomalies(ctx))
}
