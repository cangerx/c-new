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

func TestSoraOpenAIVideoExpandsVendorMediaReferenceAliases(t *testing.T) {
	plugin := compileSoraPlugin(t)
	references := []any{"https://assets.example/one.png", "https://assets.example/two.png"}
	videos := []any{"https://assets.example/guide.mp4"}
	audios := []any{"https://assets.example/voice.mp3"}
	request := map[string]any{
		"model":           "public-video-alias",
		"prompt":          "animate the references",
		"seconds":         15,
		"ratio":           "9:16",
		"resolution":      "720p",
		"referenceImages": references,
		"referenceVideos": videos,
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
	submitBody := submit["body"].(map[string]any)
	for _, key := range []string{"images", "referenceImages", "reference_images", "image_urls"} {
		assert.Equal(t, references, submitBody[key], key)
	}
	for _, key := range []string{"videos", "referenceVideos", "reference_videos", "video_urls"} {
		assert.Equal(t, videos, submitBody[key], key)
	}
	for _, key := range []string{"audios", "referenceAudios", "reference_audios", "audio_urls"} {
		assert.Equal(t, audios, submitBody[key], key)
	}
}

func TestSoraSubmitNormalizesAllImageArrayAliases(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody": map[string]any{
			"prompt":           "animate",
			"images":           []any{"https://assets.example/one.png"},
			"referenceImages":  []any{"https://assets.example/two.png"},
			"reference_images": []any{map[string]any{"url": "https://assets.example/three.png"}},
			"image_urls": []any{
				"https://assets.example/one.png",
				map[string]any{"image_url": map[string]any{"url": "https://assets.example/four.png"}},
			},
		},
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, value)
	body := submit["body"].(map[string]any)
	assert.Equal(t, []any{
		"https://assets.example/one.png",
		"https://assets.example/two.png",
		"https://assets.example/three.png",
		"https://assets.example/four.png",
	}, body["images"])
	for _, key := range []string{"images", "referenceImages", "reference_images", "image_urls"} {
		assert.Equal(t, body["images"], body[key], key)
	}
}

func TestSoraSubmitNormalizesVideoAndAudioArrayAliases(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody": map[string]any{
			"prompt":           "animate",
			"referenceVideos":  []any{"https://assets.example/one.mp4"},
			"reference_videos": []any{map[string]any{"video_url": "https://assets.example/two.mp4"}},
			"audio_urls":       []any{"https://assets.example/one.mp3"},
			"referenceAudios":  []any{map[string]any{"url": "https://assets.example/two.mp3"}},
		},
	})
	require.NoError(t, err)
	body := decodeSoraPluginMap(t, value)["body"].(map[string]any)
	for _, key := range []string{"videos", "referenceVideos", "reference_videos", "video_urls"} {
		assert.Equal(t, []any{"https://assets.example/one.mp4", "https://assets.example/two.mp4"}, body[key], key)
	}
	for _, key := range []string{"audios", "referenceAudios", "reference_audios", "audio_urls"} {
		assert.Equal(t, []any{"https://assets.example/two.mp3", "https://assets.example/one.mp3"}, body[key], key)
	}
}

func TestSoraMultipartWritesNormalizedImageURLs(t *testing.T) {
	plugin := compileSoraPlugin(t)
	file := map[string]any{
		"ref":      "request_file:input_reference",
		"field":    "input_reference",
		"filename": "frame.png",
		"mimeType": "image/png",
		"size":     12,
	}
	value, err := plugin.Engine.Call(t.Context(), "buildSubmitRequest", map[string]any{
		"baseUrl":       "https://provider.example",
		"apiKey":        "secret",
		"upstreamModel": "vendor-video",
		"requestBody": map[string]any{
			"prompt":          "animate",
			"referenceImages": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
			"referenceVideos": []any{"https://assets.example/guide.mp4"},
			"referenceAudios": []any{"https://assets.example/voice.mp3"},
		},
		"files": []any{file},
	})
	require.NoError(t, err)
	submit := decodeSoraPluginMap(t, value)
	assert.Equal(t, "multipart", submit["bodyType"])
	parts := submit["parts"].([]any)
	for _, key := range []string{"images", "referenceImages", "reference_images", "image_urls"} {
		assert.Contains(t, parts, map[string]any{"name": key, "value": "https://assets.example/one.png"})
		assert.Contains(t, parts, map[string]any{"name": key, "value": "https://assets.example/two.png"})
	}
	for _, key := range []string{"videos", "referenceVideos", "reference_videos", "video_urls"} {
		assert.Contains(t, parts, map[string]any{"name": key, "value": "https://assets.example/guide.mp4"})
	}
	for _, key := range []string{"audios", "referenceAudios", "reference_audios", "audio_urls"} {
		assert.Contains(t, parts, map[string]any{"name": key, "value": "https://assets.example/voice.mp3"})
	}
	assert.Contains(t, parts, map[string]any{
		"name":     "input_reference",
		"fileRef":  "request_file:input_reference",
		"filename": "frame.png",
	})
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

	for _, model := range []string{"grok-imagine-video-1.5-preview", "minimax-h3-f", "seedance2.5-720p", "vendor-video"} {
		preserved := decode(t, model, map[string]any{"prompt": "vendor", "seconds": 6})
		assert.Equal(t, float64(6), preserved["seconds"], model)
		assert.NotContains(t, preserved, "duration", model)
	}
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
			"referenceVideos": []any{"https://assets.example/guide.mp4"},
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
	assert.Equal(t, []any{"https://assets.example/guide.mp4"}, body["referenceVideos"])
	assert.Equal(t, []any{"https://assets.example/audio.mp3"}, body["referenceAudios"])
}

func TestSoraResponsesPreservesAllInputImages(t *testing.T) {
	plugin := compileSoraPlugin(t)
	value, err := plugin.Engine.CallPath(t.Context(), "protocols", []string{"openai_responses", "decodeRequest"}, map[string]any{
		"model":         "public-video-alias",
		"upstreamModel": "alibaba/wan-3.0",
		"body": map[string]any{"kind": "json", "value": map[string]any{
			"model":  "public-video-alias",
			"input":  "animate",
			"images": []any{"https://assets.example/one.png", "https://assets.example/two.png"},
		}},
	})
	require.NoError(t, err)
	resolved := decodeSoraPluginMap(t, value)
	body := resolved["requestBody"].(map[string]any)
	assert.Equal(t, "image_to_video", resolved["action"])
	assert.Equal(t, "https://assets.example/one.png", body["input_reference"])
	assert.Equal(t, []any{"https://assets.example/one.png", "https://assets.example/two.png"}, body["images"])
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
	assert.Equal(t, float64(5), body["seconds"])
	assert.NotContains(t, body, "duration")

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
