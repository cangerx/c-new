package relay

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/QuantumNous/new-api/common"
	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestBuildTaskRequestLogInputCapturesVideoRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set("task_request", relaycommon.TaskSubmitReq{
		Prompt:   "a city at night",
		Model:    "video-model",
		Size:     "1280x720",
		Duration: 6,
		Images: []string{
			"https://example.com/reference.png",
			"data:image/png;base64,abcdef",
		},
		Metadata: map[string]interface{}{
			"aspect_ratio": "16:9",
			"api_key":      "must-not-be-logged",
		},
	})

	input := BuildTaskRequestLogInput(c)
	var snapshot map[string]any
	require.NoError(t, common.Unmarshal([]byte(input), &snapshot))
	require.Equal(t, "a city at night", snapshot["prompt"])
	require.Equal(t, "video-model", snapshot["model"])
	require.Equal(t, "1280x720", snapshot["size"])
	require.Equal(t, float64(6), snapshot["duration"])

	images := snapshot["images"].([]any)
	require.Equal(t, "https://example.com/reference.png", images[0])
	require.Contains(t, images[1], "embedded data omitted")
	metadata := snapshot["metadata"].(map[string]any)
	require.Equal(t, "[redacted]", metadata["api_key"])
	require.NotContains(t, input, "must-not-be-logged")
}

func TestBuildTaskRequestLogInputIsBounded(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Set("task_request", relaycommon.TaskSubmitReq{
		Prompt: strings.Repeat("x", maxTaskRequestLogBytes*2),
		Metadata: map[string]interface{}{
			"large": strings.Repeat("y", maxTaskRequestLogBytes*2),
		},
	})

	input := BuildTaskRequestLogInput(c)
	require.LessOrEqual(t, len(input), maxTaskRequestLogBytes)
	var snapshot map[string]any
	require.NoError(t, common.Unmarshal([]byte(input), &snapshot))
}

func TestBuildTaskRequestLogInputReturnsEmptyWithoutParsedRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	require.Empty(t, BuildTaskRequestLogInput(c))
}

func TestBuildTaskRequestLogInputPreservesUnknownRawJSONFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)
	c.Request = httptest.NewRequest("POST", "/v1/videos", strings.NewReader(`{
		"model":"video-model",
		"prompt":"city at night",
		"aspect_ratio":"16:9",
		"resolution":"1080p",
		"custom_parameter":{"camera_motion":"dolly-in","access_token":"must-not-be-logged"}
	}`))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("task_request", relaycommon.TaskSubmitReq{
		Model:  "video-model",
		Prompt: "city at night",
	})
	t.Cleanup(func() { common.CleanupBodyStorage(c) })

	input := BuildTaskRequestLogInput(c)
	var snapshot map[string]any
	require.NoError(t, common.Unmarshal([]byte(input), &snapshot))
	require.Equal(t, "16:9", snapshot["aspect_ratio"])
	require.Equal(t, "1080p", snapshot["resolution"])
	custom := snapshot["custom_parameter"].(map[string]any)
	require.Equal(t, "dolly-in", custom["camera_motion"])
	require.Equal(t, "[redacted]", custom["access_token"])
	require.NotContains(t, input, "must-not-be-logged")
}
