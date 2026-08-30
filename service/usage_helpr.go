package service

import (
	"fmt"
	"math"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/logger"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/gin-gonic/gin"
)

const (
	// These are protocol safety ceilings, not model context-window declarations.
	// Requests approaching them require a matching local prompt estimate; an
	// upstream cannot turn a small request into an unbounded bill by forging usage.
	maxUpstreamPromptTokens     = 8 * 1024 * 1024
	maxUpstreamCompletionTokens = 1024 * 1024
	usagePromptSlack            = 64 * 1024
	usageCompletionSlack        = 128 * 1024
)

//func GetPromptTokens(textRequest dto.GeneralOpenAIRequest, relayMode int) (int, error) {
//	switch relayMode {
//	case constant.RelayModeChatCompletions:
//		return CountTokenMessages(textRequest.Messages, textRequest.Model)
//	case constant.RelayModeCompletions:
//		return CountTokenInput(textRequest.Prompt, textRequest.Model), nil
//	case constant.RelayModeModerations:
//		return CountTokenInput(textRequest.Input, textRequest.Model), nil
//	}
//	return 0, errors.New("unknown relay mode")
//}

func ResponseText2Usage(c *gin.Context, responseText string, modeName string, promptTokens int) *dto.Usage {
	common.SetContextKey(c, constant.ContextKeyLocalCountTokens, true)
	usage := &dto.Usage{}
	usage.PromptTokens = promptTokens
	usage.CompletionTokens = EstimateTokenByModel(modeName, responseText)
	usage.TotalTokens = usage.PromptTokens + usage.CompletionTokens
	return usage
}

func ValidUsage(usage *dto.Usage) bool {
	if !HasUsage(usage) {
		return false
	}
	_, reason := normalizeAndValidateUsage(usage, maxUpstreamPromptTokens, maxUpstreamCompletionTokens)
	return reason == ""
}

// HasUsage reports whether an upstream supplied any token count, including an
// invalid negative count. It is intentionally separate from ValidUsage so a
// malformed first snapshot cannot be skipped in favor of a later forged one.
func HasUsage(usage *dto.Usage) bool {
	if usage == nil {
		return false
	}
	if usage.PromptTokens != 0 || usage.CompletionTokens != 0 || usage.InputTokens != 0 ||
		usage.OutputTokens != 0 || usage.TotalTokens != 0 || usage.PromptCacheHitTokens != 0 ||
		usage.ClaudeCacheCreation5mTokens != 0 || usage.ClaudeCacheCreation1hTokens != 0 {
		return true
	}
	input := usage.PromptTokensDetails
	output := usage.CompletionTokenDetails
	if input.CachedTokens != 0 || input.CachedCreationTokens != 0 || input.CacheWriteTokens != 0 ||
		input.TextTokens != 0 || input.AudioTokens != 0 || input.ImageTokens != 0 ||
		output.TextTokens != 0 || output.AudioTokens != 0 || output.ImageTokens != 0 || output.ReasoningTokens != 0 {
		return true
	}
	if usage.InputTokensDetails == nil {
		return false
	}
	responsesInput := usage.InputTokensDetails
	return responsesInput.CachedTokens != 0 || responsesInput.CachedCreationTokens != 0 ||
		responsesInput.CacheWriteTokens != 0 || responsesInput.TextTokens != 0 ||
		responsesInput.AudioTokens != 0 || responsesInput.ImageTokens != 0
}

func UsageAnomalies(ctx *gin.Context) []string {
	if ctx == nil {
		return nil
	}
	anomalies, _ := common.GetContextKeyType[[]string](ctx, constant.ContextKeyUsageAnomalies)
	return anomalies
}

func RecordUsageAnomaly(ctx *gin.Context, reason string) {
	if ctx == nil || reason == "" {
		return
	}
	anomalies := UsageAnomalies(ctx)
	for _, existing := range anomalies {
		if existing == reason {
			return
		}
	}
	if len(anomalies) >= 16 {
		return
	}
	common.SetContextKey(ctx, constant.ContextKeyUsageAnomalies, append(anomalies, reason))
	logger.LogWarn(ctx, "upstream usage rejected: "+reason)
}

// SanitizeUsageForBilling treats upstream usage as untrusted input. It returns
// a normalized copy when the snapshot is structurally valid and plausible for
// the request, otherwise it returns the locally counted fallback.
func SanitizeUsageForBilling(ctx *gin.Context, info *relaycommon.RelayInfo, reported, fallback *dto.Usage) *dto.Usage {
	safeFallback := localUsageFallback(info, fallback)
	if reported == nil {
		return safeFallback
	}
	if !HasUsage(reported) {
		common.SetContextKey(ctx, constant.ContextKeyLocalCountTokens, true)
		return safeFallback
	}

	promptLimit := maxUpstreamPromptTokens
	if info != nil && info.GetEstimatePromptTokens() > 0 {
		promptLimit = scaledUsageLimit(info.GetEstimatePromptTokens(), 8, usagePromptSlack, maxUpstreamPromptTokens)
	}
	completionLimit := maxUpstreamCompletionTokens
	if safeFallback.CompletionTokens > 0 {
		completionLimit = scaledUsageLimit(safeFallback.CompletionTokens, 32, usageCompletionSlack, maxUpstreamCompletionTokens)
	}

	normalized, reason := normalizeAndValidateUsage(reported, promptLimit, completionLimit)
	if reason == "" {
		return normalized
	}

	RecordUsageAnomaly(ctx, fmt.Sprintf("%s (prompt=%d completion=%d total=%d)",
		reason, reported.PromptTokens, reported.CompletionTokens, reported.TotalTokens))
	common.SetContextKey(ctx, constant.ContextKeyLocalCountTokens, true)
	return safeFallback
}

func SanitizeRealtimeUsage(ctx *gin.Context, reported, local *dto.RealtimeUsage) *dto.RealtimeUsage {
	fallback := cloneRealtimeUsage(local)
	if reported == nil {
		return fallback
	}
	reason := validateRealtimeUsage(reported)
	if reason == "" && fallback.TotalTokens > 0 {
		inputLimit := scaledUsageLimit(fallback.InputTokens, 32, 256*1024, maxUpstreamPromptTokens)
		outputLimit := scaledUsageLimit(fallback.OutputTokens, 32, usageCompletionSlack, maxUpstreamCompletionTokens)
		if reported.InputTokens > inputLimit {
			reason = "input_tokens_exceed_local_limit"
		} else if reported.OutputTokens > outputLimit {
			reason = "output_tokens_exceed_local_limit"
		}
	}
	if reason != "" {
		RecordUsageAnomaly(ctx, fmt.Sprintf("realtime_%s (input=%d output=%d total=%d)",
			reason, reported.InputTokens, reported.OutputTokens, reported.TotalTokens))
		common.SetContextKey(ctx, constant.ContextKeyLocalCountTokens, true)
		return fallback
	}
	return cloneRealtimeUsage(reported)
}

func localUsageFallback(info *relaycommon.RelayInfo, fallback *dto.Usage) *dto.Usage {
	local := &dto.Usage{}
	if fallback != nil {
		*local = *fallback
		local.BillingUsage = nil
	}
	if local.PromptTokens < 0 {
		local.PromptTokens = 0
	}
	if local.CompletionTokens < 0 {
		local.CompletionTokens = 0
	}
	if local.PromptTokens == 0 && info != nil && info.GetEstimatePromptTokens() > 0 {
		local.PromptTokens = info.GetEstimatePromptTokens()
	}
	local.InputTokens = local.PromptTokens
	local.OutputTokens = local.CompletionTokens
	local.TotalTokens = safeUsageSum(local.PromptTokens, local.CompletionTokens)
	local.UsageSource = "local_usage_guard"
	return local
}

func normalizeAndValidateUsage(usage *dto.Usage, promptLimit, completionLimit int) (*dto.Usage, string) {
	if usage == nil {
		return nil, "missing_usage"
	}
	normalized := *usage

	fields := []struct {
		name  string
		value int
	}{
		{"prompt_tokens", usage.PromptTokens},
		{"completion_tokens", usage.CompletionTokens},
		{"total_tokens", usage.TotalTokens},
		{"input_tokens", usage.InputTokens},
		{"output_tokens", usage.OutputTokens},
		{"prompt_cache_hit_tokens", usage.PromptCacheHitTokens},
		{"cached_tokens", usage.PromptTokensDetails.CachedTokens},
		{"cached_creation_tokens", usage.PromptTokensDetails.CachedCreationTokens},
		{"cache_write_tokens", usage.PromptTokensDetails.CacheWriteTokens},
		{"input_text_tokens", usage.PromptTokensDetails.TextTokens},
		{"input_audio_tokens", usage.PromptTokensDetails.AudioTokens},
		{"input_image_tokens", usage.PromptTokensDetails.ImageTokens},
		{"output_text_tokens", usage.CompletionTokenDetails.TextTokens},
		{"output_audio_tokens", usage.CompletionTokenDetails.AudioTokens},
		{"output_image_tokens", usage.CompletionTokenDetails.ImageTokens},
		{"reasoning_tokens", usage.CompletionTokenDetails.ReasoningTokens},
		{"claude_cache_creation_5m_tokens", usage.ClaudeCacheCreation5mTokens},
		{"claude_cache_creation_1h_tokens", usage.ClaudeCacheCreation1hTokens},
	}
	if usage.InputTokensDetails != nil {
		if normalized.PromptTokensDetails.CachedTokens == 0 {
			normalized.PromptTokensDetails.CachedTokens = usage.InputTokensDetails.CachedTokens
		}
		if normalized.PromptTokensDetails.CachedCreationTokens == 0 {
			normalized.PromptTokensDetails.CachedCreationTokens = usage.InputTokensDetails.CachedCreationTokens
		}
		if normalized.PromptTokensDetails.CacheWriteTokens == 0 {
			normalized.PromptTokensDetails.CacheWriteTokens = usage.InputTokensDetails.CacheWriteTokens
		}
		if normalized.PromptTokensDetails.TextTokens == 0 {
			normalized.PromptTokensDetails.TextTokens = usage.InputTokensDetails.TextTokens
		}
		if normalized.PromptTokensDetails.AudioTokens == 0 {
			normalized.PromptTokensDetails.AudioTokens = usage.InputTokensDetails.AudioTokens
		}
		if normalized.PromptTokensDetails.ImageTokens == 0 {
			normalized.PromptTokensDetails.ImageTokens = usage.InputTokensDetails.ImageTokens
		}
		fields = append(fields,
			struct {
				name  string
				value int
			}{"responses_cached_tokens", usage.InputTokensDetails.CachedTokens},
			struct {
				name  string
				value int
			}{"responses_cache_write_tokens", usage.InputTokensDetails.CacheWriteTokens},
			struct {
				name  string
				value int
			}{"responses_cached_creation_tokens", usage.InputTokensDetails.CachedCreationTokens},
			struct {
				name  string
				value int
			}{"responses_text_tokens", usage.InputTokensDetails.TextTokens},
			struct {
				name  string
				value int
			}{"responses_audio_tokens", usage.InputTokensDetails.AudioTokens},
			struct {
				name  string
				value int
			}{"responses_image_tokens", usage.InputTokensDetails.ImageTokens},
		)
	}
	for _, field := range fields {
		if field.value < 0 {
			return nil, "negative_" + field.name
		}
		if field.value > maxUpstreamPromptTokens {
			return nil, "oversized_" + field.name
		}
	}

	if usage.UsageSemantic != dto.BillingUsageSemanticAnthropic &&
		usage.PromptTokens > 0 && usage.InputTokens > 0 && usage.PromptTokens != usage.InputTokens {
		return nil, "inconsistent_prompt_and_input_tokens"
	}
	if usage.CompletionTokens > 0 && usage.OutputTokens > 0 && usage.CompletionTokens != usage.OutputTokens {
		return nil, "inconsistent_completion_and_output_tokens"
	}
	if normalized.PromptTokens == 0 {
		normalized.PromptTokens = normalized.InputTokens
	}
	if normalized.CompletionTokens == 0 {
		normalized.CompletionTokens = normalized.OutputTokens
	}
	if normalized.PromptTokens > promptLimit {
		return nil, "prompt_tokens_exceed_request_limit"
	}
	if normalized.CompletionTokens > completionLimit {
		return nil, "completion_tokens_exceed_request_limit"
	}
	total := safeUsageSum(normalized.PromptTokens, normalized.CompletionTokens)
	if usage.TotalTokens != 0 && usage.TotalTokens != total {
		return nil, "inconsistent_total_tokens"
	}

	if normalized.UsageSemantic != dto.BillingUsageSemanticAnthropic {
		inputDetails := []int{
			normalized.PromptTokensDetails.CachedTokens,
			normalized.PromptTokensDetails.CachedCreationTokens,
			normalized.PromptTokensDetails.CacheWriteTokens,
			normalized.PromptTokensDetails.TextTokens,
			normalized.PromptTokensDetails.AudioTokens,
			normalized.PromptTokensDetails.ImageTokens,
			normalized.PromptCacheHitTokens,
		}
		for _, detail := range inputDetails {
			if detail > normalized.PromptTokens {
				return nil, "input_detail_exceeds_prompt_tokens"
			}
		}
	} else {
		contextTokens := safeUsageSum(normalized.PromptTokens, normalized.PromptTokensDetails.CachedTokens)
		contextTokens = safeUsageSum(contextTokens, normalized.PromptTokensDetails.CacheCreationTokensTotal())
		if contextTokens > promptLimit {
			return nil, "anthropic_context_tokens_exceed_request_limit"
		}
		if normalized.InputTokens > 0 && normalized.InputTokens != contextTokens {
			return nil, "inconsistent_anthropic_input_tokens"
		}
	}
	outputDetails := []int{
		normalized.CompletionTokenDetails.TextTokens,
		normalized.CompletionTokenDetails.AudioTokens,
		normalized.CompletionTokenDetails.ImageTokens,
		normalized.CompletionTokenDetails.ReasoningTokens,
	}
	for _, detail := range outputDetails {
		if detail > normalized.CompletionTokens {
			return nil, "output_detail_exceeds_completion_tokens"
		}
	}

	if normalized.UsageSemantic != dto.BillingUsageSemanticAnthropic {
		normalized.InputTokens = normalized.PromptTokens
	}
	normalized.OutputTokens = normalized.CompletionTokens
	normalized.TotalTokens = total
	return &normalized, ""
}

func validateRealtimeUsage(usage *dto.RealtimeUsage) string {
	if usage == nil {
		return "missing_usage"
	}
	fields := []struct {
		name  string
		value int
	}{
		{"total_tokens", usage.TotalTokens},
		{"input_tokens", usage.InputTokens},
		{"output_tokens", usage.OutputTokens},
		{"input_cached_tokens", usage.InputTokenDetails.CachedTokens},
		{"input_text_tokens", usage.InputTokenDetails.TextTokens},
		{"input_audio_tokens", usage.InputTokenDetails.AudioTokens},
		{"output_text_tokens", usage.OutputTokenDetails.TextTokens},
		{"output_audio_tokens", usage.OutputTokenDetails.AudioTokens},
	}
	for _, field := range fields {
		if field.value < 0 {
			return "negative_" + field.name
		}
	}
	if usage.InputTokens > maxUpstreamPromptTokens {
		return "input_tokens_exceed_limit"
	}
	if usage.OutputTokens > maxUpstreamCompletionTokens {
		return "output_tokens_exceed_limit"
	}
	total := safeUsageSum(usage.InputTokens, usage.OutputTokens)
	if usage.TotalTokens != 0 && usage.TotalTokens != total {
		return "inconsistent_total_tokens"
	}
	if usage.InputTokenDetails.TextTokens > usage.InputTokens ||
		usage.InputTokenDetails.AudioTokens > usage.InputTokens ||
		usage.InputTokenDetails.CachedTokens > usage.InputTokens ||
		usage.InputTokenDetails.TextTokens > usage.InputTokens-usage.InputTokenDetails.AudioTokens {
		return "input_detail_exceeds_input_tokens"
	}
	if usage.OutputTokenDetails.TextTokens > usage.OutputTokens ||
		usage.OutputTokenDetails.AudioTokens > usage.OutputTokens ||
		usage.OutputTokenDetails.TextTokens > usage.OutputTokens-usage.OutputTokenDetails.AudioTokens {
		return "output_detail_exceeds_output_tokens"
	}
	return ""
}

func cloneRealtimeUsage(usage *dto.RealtimeUsage) *dto.RealtimeUsage {
	if usage == nil {
		return &dto.RealtimeUsage{}
	}
	clone := *usage
	if clone.InputTokens < 0 {
		clone.InputTokens = 0
	}
	if clone.OutputTokens < 0 {
		clone.OutputTokens = 0
	}
	clone.TotalTokens = safeUsageSum(clone.InputTokens, clone.OutputTokens)
	return &clone
}

func scaledUsageLimit(base, multiplier, slack, ceiling int) int {
	if base <= 0 {
		return ceiling
	}
	if base > (ceiling-slack)/multiplier {
		return ceiling
	}
	limit := base*multiplier + slack
	if limit > ceiling {
		return ceiling
	}
	return limit
}

func safeUsageSum(left, right int) int {
	if left < 0 || right < 0 {
		return 0
	}
	if left > math.MaxInt-right {
		return math.MaxInt
	}
	return left + right
}
