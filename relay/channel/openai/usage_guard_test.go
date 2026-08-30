package openai

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/QuantumNous/new-api/relaykit/types"
	"github.com/QuantumNous/new-api/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUsageGuardStream(t *testing.T, path, body string) (*gin.Context, *httptest.ResponseRecorder, *http.Response) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	oldTimeout := constant.StreamingTimeout
	constant.StreamingTimeout = 30
	t.Cleanup(func() { constant.StreamingTimeout = oldTimeout })

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, path, nil)
	ctx.Set(common.RequestIdKey, "usage-guard-test")
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"text/event-stream"}},
	}
	return ctx, recorder, resp
}

func TestOaiStreamHandlerUsesOneSanitizedUsageSnapshot(t *testing.T) {
	body := strings.Join([]string{
		`data: {"id":"chat_1","object":"chat.completion.chunk","model":"test-model","choices":[{"index":0,"delta":{"content":"hello"}}]}`,
		``,
		`data: {"id":"chat_1","object":"chat.completion.chunk","model":"test-model","choices":[],"usage":{"prompt_tokens":10,"completion_tokens":2,"total_tokens":12}}`,
		``,
		`data: {"id":"chat_1","object":"chat.completion.chunk","model":"test-model","choices":[],"usage":{"prompt_tokens":1000,"completion_tokens":200,"total_tokens":1200}}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")
	ctx, recorder, resp := newUsageGuardStream(t, "/v1/chat/completions", body)
	info := &relaycommon.RelayInfo{
		RelayFormat:        types.RelayFormatOpenAI,
		ShouldIncludeUsage: true,
		DisablePing:        true,
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "test-model",
		},
	}
	info.SetEstimatePromptTokens(10)

	usage, apiErr := OaiStreamHandler(ctx, info, resp)

	require.Nil(t, apiErr)
	require.NotNil(t, usage)
	assert.Equal(t, 10, usage.PromptTokens)
	assert.Equal(t, 2, usage.CompletionTokens)
	assert.Equal(t, 12, usage.TotalTokens)
	assert.Equal(t, 1, strings.Count(recorder.Body.String(), `"usage"`))
	assert.NotEmpty(t, service.UsageAnomalies(ctx))
}

func TestOaiResponsesStreamHandlerDoesNotAccumulateTerminalUsage(t *testing.T) {
	body := strings.Join([]string{
		`data: {"type":"response.output_text.delta","delta":"hello"}`,
		``,
		`data: {"type":"response.completed","response":{"id":"resp_1","status":"completed","usage":{"input_tokens":10,"output_tokens":2,"total_tokens":12}}}`,
		``,
		`data: {"type":"response.done","response":{"id":"resp_1","status":"completed","usage":{"input_tokens":1000,"output_tokens":200,"total_tokens":1200}}}`,
		``,
		`data: [DONE]`,
		``,
	}, "\n")
	ctx, recorder, resp := newUsageGuardStream(t, "/v1/responses", body)
	info := &relaycommon.RelayInfo{
		DisablePing: true,
		ChannelMeta: &relaycommon.ChannelMeta{
			UpstreamModelName: "test-model",
		},
	}
	info.SetEstimatePromptTokens(10)

	usage, apiErr := OaiResponsesStreamHandler(ctx, info, resp)

	require.Nil(t, apiErr)
	require.NotNil(t, usage)
	assert.Equal(t, 10, usage.PromptTokens)
	assert.Equal(t, 2, usage.CompletionTokens)
	assert.Equal(t, 12, usage.TotalTokens)
	assert.Equal(t, 1, strings.Count(recorder.Body.String(), `"usage"`))
	assert.NotEmpty(t, service.UsageAnomalies(ctx))
}

func TestRealtimeUsageTrackerRejectsDuplicateResponseDone(t *testing.T) {
	tracker := newRealtimeUsageTracker()
	event := &dto.RealtimeEvent{
		EventId: "event_1",
		Type:    dto.RealtimeEventTypeResponseDone,
		Response: &dto.RealtimeResponse{
			Id: "resp_1",
		},
	}

	assert.True(t, tracker.Accept(event))
	assert.False(t, tracker.Accept(event))
	event.EventId = "event_2"
	assert.False(t, tracker.Accept(event), "response id, not event id, defines billing identity")
}
