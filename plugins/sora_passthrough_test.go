package plugins_test

import (
	"testing"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/pkg/jsplugin"
	builtinplugins "github.com/QuantumNous/new-api/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func compileSoraPlugin(t *testing.T) *jsplugin.LoadedPlugin {
	t.Helper()
	source, err := builtinplugins.Source("sora")
	require.NoError(t, err)
	plugin, err := jsplugin.NewRegistry().RegisterFactory(source, jsplugin.Options{Key: "sora"})
	require.NoError(t, err)
	return plugin
}

func decodeSoraPluginMap(t *testing.T, value any) map[string]any {
	t.Helper()
	encoded, err := common.Marshal(value)
	require.NoError(t, err)
	var decoded map[string]any
	require.NoError(t, common.Unmarshal(encoded, &decoded))
	return decoded
}

func TestSoraOpenAIVideoPreservesVendorParameters(t *testing.T) {
	plugin := compileSoraPlugin(t)
	references := []any{"https://assets.example/one.png", "https://assets.example/two.png"}
	audios := []any{"https://assets.example/voice.mp3"}
	request := map[string]any{
		"model":           "public-video-alias",
		"prompt":          "animate the references",
		"seconds":         15,
		"ratio":           "9:16",
		"resolution":      "720p",
		"referenceImages": references,
		"referenceAudios": audios,
	}

	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
		"model":         "public-video-alias",
		"upstreamModel": "alibaba/wan-3.0",
		"body":          map[string]any{"kind": "json", "value": request},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	assert.Equal(t, "image_to_video", resolved["action"])
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, float64(15), body["duration"])
	assert.NotContains(t, body, "seconds")
	assert.Equal(t, "9:16", body["ratio"])
	assert.Equal(t, "720p", body["resolution"])
	assert.Equal(t, references, body["referenceImages"])
	assert.Equal(t, audios, body["referenceAudios"])

	submitValue, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "alibaba/wan-3.0",
		"requestBody":   body,
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, submitValue)
	assert.Equal(t, body["referenceImages"], submit["body"].(map[string]any)["referenceImages"])
	assert.Equal(t, body["referenceAudios"], submit["body"].(map[string]any)["referenceAudios"])
}

func TestSoraDurationUsesNativeOrVendorWireFormat(t *testing.T) {
	plugin := compileSoraPlugin(t)
	decode := func(t *testing.T, model string, request map[string]any) map[string]any {
		t.Helper()
		value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
			"model": model,
			"body":  map[string]any{"kind": "json", "value": request},
		})
		require.NoError(t, err)
		return decodeSoraPluginMap(t, value)["requestBody"].(map[string]any)
	}

	native := decode(t, "sora-2", map[string]any{"prompt": "native", "duration": 8})
	assert.Equal(t, float64(8), native["seconds"])
	assert.NotContains(t, native, "duration")

	vendor := decode(t, "seedance-2-5-1080p", map[string]any{"prompt": "vendor", "seconds": 10})
	assert.Equal(t, float64(10), vendor["duration"])
	assert.NotContains(t, vendor, "seconds")
}

func TestSoraResponsesPreservesVideoExtensions(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_responses", "decodeRequest"}, map[string]any{
		"model":         "public-video-alias",
		"upstreamModel": "seedance-2-5-1080p",
		"body": map[string]any{"kind": "json", "value": map[string]any{
			"model":           "public-video-alias",
			"input":           "animate",
			"seconds":         10,
			"ratio":           "16:9",
			"resolution":      "1080p",
			"referenceImages": []any{"https://assets.example/frame.png"},
			"referenceAudios": []any{"https://assets.example/audio.mp3"},
		}},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, "image_to_video", resolved["action"])
	assert.Equal(t, float64(10), body["duration"])
	assert.NotContains(t, body, "seconds")
	assert.Equal(t, "16:9", body["ratio"])
	assert.Equal(t, "1080p", body["resolution"])
	assert.Equal(t, []any{"https://assets.example/frame.png"}, body["referenceImages"])
	assert.Equal(t, []any{"https://assets.example/audio.mp3"}, body["referenceAudios"])
}

func TestSoraMultipartKeepsInputReferenceFile(t *testing.T) {
	plugin := compileSoraPlugin(t)
	file := map[string]any{
		"ref":      "request_file:input_reference",
		"field":    "input_reference",
		"filename": "frame.png",
		"mimeType": "image/png",
		"size":     12,
	}
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_video", "decodeRequest"}, map[string]any{
		"model":         "vendor-video",
		"upstreamModel": "vendor-video",
		"body": map[string]any{
			"kind":   "multipart",
			"fields": map[string]any{"prompt": []any{"animate"}, "seconds": []any{"5"}},
			"files":  []any{file},
		},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, float64(5), body["duration"])
	assert.NotContains(t, body, "seconds")

	submitValue, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody":   body,
		"files":         []any{file},
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, submitValue)
	assert.Equal(t, "multipart", submit["bodyType"])
	parts := submit["parts"].([]any)
	require.NotEmpty(t, parts)
	assert.Contains(t, parts, map[string]any{
		"name":     "input_reference",
		"fileRef":  "request_file:input_reference",
		"filename": "frame.png",
	})
}
